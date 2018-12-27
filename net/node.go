package net

import (
	"context"
	"fmt"
	"github.com/eager7/go/elog"
	"github.com/eager7/lib-p2p/common/errors"
	"github.com/eager7/lib-p2p/message"
	"github.com/gogo/protobuf/io"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"strings"
	"sync"
	"time"
)

var log = elog.Log

const (
	Protocol           = "/eager7/test/1.0.0"
	sendMessageTimeout = time.Minute * 10
)

type Instance struct {
	ctx     context.Context
	Host    host.Host
	ID      peer.ID
	Address string
	Port    string
	Peers   PeerMap
	lock    sync.RWMutex
}

func New(ctx context.Context, b64Pri, address, port string) (*Instance, error) {
	i := new(Instance)
	if err := i.initialize(ctx, b64Pri, address, port); err != nil {
		return nil, err
	}
	return i, nil
}

func (i *Instance) initialize(ctx context.Context, b64Pri, address, port string) error {
	i.Peers.Initialize()
	if ctx == nil {
		ctx = context.Background()
	}
	i.ctx = ctx
	i.Address = address
	i.Port = port
	return i.initNetwork(b64Pri)
}

func (i *Instance) initNetwork(b64Pri string) (err error) {
	var private crypto.PrivKey
	var public crypto.PubKey
	if b64Pri == "" {
		private, public, err = crypto.GenerateKeyPair(crypto.RSA, 1024)
		if err != nil {
			return err
		}
		b, _ := private.Bytes()
		log.Info("generate private b64 key:", crypto.ConfigEncodeKey(b))
		b, _ = public.Bytes()
		log.Info("generate public b64 key:", crypto.ConfigEncodeKey(b))
	} else {
		key, err := crypto.ConfigDecodeKey(b64Pri)
		if err != nil {
			return err
		}
		private, err = crypto.UnmarshalPrivateKey(key)
		if err != nil {
			return err
		}
	}

	i.ID, err = peer.IDFromPrivateKey(private)
	if err != nil {
		return err
	}
	log.Info("this node id is :", i.ID.String())

	var options []libp2p.Option
	options = append(options, libp2p.Identity(private))
	ps := peerstore.NewPeerstore()
	ps.AddPrivKey(i.ID, private)
	ps.AddPubKey(i.ID, private.GetPublic())
	options = append(options, libp2p.Peerstore(ps))

	i.Host, err = libp2p.New(i.ctx, options...)
	if err != nil {
		return err
	}

	i.Host.SetStreamHandler(Protocol, i.NetworkHandler)
	i.Host.Network().Notify(i)

	mAddr, err := multiaddr.NewMultiaddr(NewAddrInfo(i.Address, i.Port))
	if err != nil {
		return err
	}
	err = i.Host.Network().Listen([]multiaddr.Multiaddr{mAddr}...)
	if err != nil {
		return err
	}
	log.Debug(i.Host.Network().InterfaceListenAddresses())

	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", i.Host.ID().Pretty()))
	addresses := i.Host.Addrs()
	var addrM multiaddr.Multiaddr
	for _, i := range addresses {
		if strings.HasPrefix(i.String(), "/ip4") {
			addrM = i
			break
		}
	}
	fullAddr := addrM.Encapsulate(hostAddr)
	log.Debug("I am ", fullAddr, i.Host.Peerstore().Addrs(i.ID))
	return nil
}

//每个连接只会触发一次这个回调函数，之后需要在线程中做收发
func (i *Instance) NetworkHandler(s net.Stream) {
	log.Debug("receive connect peer from:", s.Conn().RemotePeer().Pretty(), s.Conn().RemoteMultiaddr(), "| topic is:", s.Protocol(), s)
	i.Peers.Add(s.Conn().RemotePeer(), s, s.Conn().RemoteMultiaddr())

	go i.ReceiveMessage(s)
}

func (i *Instance) StreamConnect(b64Pub, address, port string) (net.Stream, error) {
	log.Info("connect to:", address, port)
	id, err := IdFromPublicKey(b64Pub)
	if err != nil {
		return nil, err
	}
	if p := i.Peers.Get(id); p != nil {
		log.Warn("the stream is created:", p.s)
		return p.s, nil
	}

	addr, err := multiaddr.NewMultiaddr(NewAddrInfo(address, port))
	if err != nil {
		return nil, errors.New(err.Error())
	}
	i.Host.Peerstore().AddAddr(id, addr, peerstore.PermanentAddrTTL)
	s, err := i.Host.NewStream(i.ctx, id, Protocol)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	log.Info("add stream:", s)
	i.Peers.Add(id, s, addr)
	//go i.ReceiveMessage(s)
	return s, nil
}

/**
 *  @brief 建立和对端的连接，这个函数会调用Dial函数拨号，可以节省ConnectPeer函数执行时间，因为如果已经拨号成功，那么创建流时就不需要再次拨号，因此此函数可以作为ping函数使用，实时去刷新和节点间的连接
 *  @param b64Pub - the public key
 *  @param address - the address of ip
 *  @param port - the port of ip
 */
func (i *Instance) Connect(b64Pub, address, port string) error {
	id, err := IdFromPublicKey(b64Pub)
	if err != nil {
		return err
	}
	addr, err := multiaddr.NewMultiaddr(NewAddrInfo(address, port))
	if err != nil {
		return errors.New(err.Error())
	}
	pi := peerstore.PeerInfo{ID: id, Addrs: []multiaddr.Multiaddr{addr}}
	if err := i.Host.Connect(i.ctx, pi); err != nil {
		return errors.New(err.Error())
	}
	//i.Peers.Add(id, nil, addr, b64Pub)
	return nil
}

func (i *Instance) ReceiveMessage(s net.Stream) {
	log.Debug("start receive thread")
	reader := io.NewDelimitedReader(s, net.MessageSizeMax)
	for {
		msg := new(pnet.Message)
		err := reader.ReadMsg(msg)
		if err != nil {
			log.Error("the peer ", s.Conn().RemotePeer().Pretty(), i.Host.Peerstore().Addrs(s.Conn().RemotePeer()), "is disconnected:", err)
			s.Reset()
			i.Peers.Del(s.Conn().RemotePeer())
			return
		}
		log.Info("receive msg:", msg.String())
	}
}

func (i *Instance) SendMessage(b64Pub string, message *pnet.Message) error {
	id, err := IdFromPublicKey(b64Pub)
	if err != nil {
		return errors.New(err.Error())
	}
	var s net.Stream
	info := i.Peers.Get(id)
	if info == nil {
		return errors.New(fmt.Sprintf("the peer %s-%s is not connect", id.Pretty(), i.Host.Peerstore().Addrs(id)))
	} else {
		s = info.s
	}

	deadline := time.Now().Add(sendMessageTimeout)
	if dl, ok := i.ctx.Deadline(); ok {
		deadline = dl
		log.Info("set deal line:", deadline)
	}
	if err := s.SetWriteDeadline(deadline); err != nil {
		return errors.New(err.Error())
	}

	writer := io.NewDelimitedWriter(s)
	err = writer.WriteMsg(message)
	if err != nil {
		return errors.New(err.Error())
	}
	if err := s.SetWriteDeadline(time.Time{}); err != nil {
		log.Warn("error resetting deadline: ", err)
	}
	log.Info("send message finished:", message)
	return nil
}

func (i *Instance) ResetStream(s net.Stream) error {
	id := s.Conn().RemotePeer()
	if err := s.Reset(); err != nil {
		return err
	}
	i.Peers.Del(id)
	return nil
}
