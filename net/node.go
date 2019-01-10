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
	"time"
)

var log = elog.Log

const (
	Protocol           = "/eager7/test/1.0.0"
	sendMessageTimeout = time.Minute * 10
)

type Instance struct {
	ctx        context.Context
	Host       host.Host
	ID         peer.ID
	Address    []string
	SenderMap  PeerMap
	BootStrap  *BootStrap
	RouteTable *RouteTable
}

func NewInstance(ctx context.Context, b64Pri string, address ...string) (*Instance, error) {
	i := new(Instance)
	i.SenderMap.Initialize()
	if ctx == nil {
		ctx = context.Background()
	}
	i.ctx = ctx
	i.Address = address
	if err := i.initNetwork(b64Pri); err != nil {
		return nil, err
	}
	return i, nil
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

	libp2p.ConnectionManager()
	ps := peerstore.NewPeerstore()
	if err := ps.AddPrivKey(i.ID, private); err != nil {
		return errors.New(err.Error())
	}
	if err := ps.AddPubKey(i.ID, private.GetPublic()); err != nil {
		return errors.New(err.Error())
	}
	options = append(options, libp2p.Peerstore(ps))
	options = append(options, libp2p.ListenAddrStrings(i.Address...))

	i.Host, err = libp2p.New(i.ctx, options...)
	if err != nil {
		return err
	}

	i.Host.SetStreamHandler(Protocol, func(stream net.Stream) {
		log.Debug("receive stream from:", stream.Conn().RemotePeer().Pretty(), stream.Conn().RemoteMultiaddr(), "| topic is:", stream.Protocol())
		i.RouteTable.RouteUpdate(stream.Conn().RemotePeer())
		go i.receive(stream)
	})
	i.Host.Network().Notify(i)
	i.RouteTable = RouteInitialize(i.Host)
	listen, err := i.Host.Network().InterfaceListenAddresses()
	if err != nil {
		log.Error(err)
	}
	log.Debug("i am listen", listen)

	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", i.Host.ID().Pretty()))
	addresses := i.Host.Addrs()
	var addrM multiaddr.Multiaddr
	for _, i := range addresses {
		if strings.HasPrefix(i.String(), "/ip4") {
			addrM = i
			break
		}
	}
	log.Debug("self ipfs addr:", addrM.Encapsulate(hostAddr))
	return nil
}

func (i *Instance) SendMessage(b64Pub, address, port string, message *mpb.Message) (err error) {
	var s net.Stream
	if s, err = i.newStream(b64Pub, address, port); err != nil {
		return err
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

func (i *Instance) newStream(b64Pub, address, port string) (net.Stream, error) {
	id, err := IdFromPublicKey(b64Pub)
	if err != nil {
		return nil, err
	}
	if p := i.SenderMap.Get(id); p != nil {
		log.Warn("the stream is created already:", p.ID.Pretty())
		return p.s, nil
	}

	addr, err := multiaddr.NewMultiaddr(NewAddrInfo(address, port))
	if err != nil {
		return nil, errors.New(err.Error())
	}
	if len(i.Host.Peerstore().Addrs(id)) == 0 {
		i.Host.Peerstore().AddAddr(id, addr, time.Minute*10)
	}
	log.Info("create new stream:", id.Pretty(), address, port)
	s, err := i.Host.NewStream(i.ctx, id, Protocol)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	log.Info("add stream:", s)
	i.SenderMap.Add(id, s, addr)
	return s, nil
}
func (i *Instance) reset(s net.Stream) error {
	if err := s.Reset(); err != nil {
		return err
	}
	i.SenderMap.Del(s.Conn().RemotePeer())
	return nil
}
func (i *Instance) receive(s net.Stream) {
	log.Debug("start receive thread")
	reader := io.NewDelimitedReader(s, net.MessageSizeMax)
	for {
		msg := new(mpb.Message)
		if err := reader.ReadMsg(msg); err != nil {
			log.Error("the peer ", s.Conn().RemotePeer().Pretty(), i.Host.Peerstore().Addrs(s.Conn().RemotePeer()), "is disconnected:", err)
			if err := s.Reset(); err != nil {
				log.Error(err)
			}
			return
		}
		log.Info("receive msg:", msg.String())
	}
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
	return nil
}
