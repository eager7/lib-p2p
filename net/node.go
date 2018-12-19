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

type PeerInfo struct {
	s         net.Stream
	MultiAddr multiaddr.Multiaddr
	PublicKey string
}

type Instance struct {
	ctx     context.Context
	Host    host.Host
	ID      peer.ID
	Address string
	Port    string
	Peers   map[peer.ID]PeerInfo
	lock    sync.RWMutex
}

func New(ctx context.Context, b64Pri, address, port string) (*Instance, error) {
	i := new(Instance)
	if err := i.Init(ctx, b64Pri, address, port); err != nil {
		return nil, err
	}
	return i, nil
}

func (i *Instance) Init(ctx context.Context, b64Pri, address, port string) error {
	i.Peers = make(map[peer.ID]PeerInfo, 0)
	i.ctx = ctx
	i.Address = address
	i.Port = port
	return i.InitNetwork(b64Pri)
}

func (i *Instance) InitNetwork(b64Pri string) (err error) {
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

	if i.ctx == nil {
		i.ctx = context.Background()
	}
	i.Host, err = libp2p.New(i.ctx, options...)
	if err != nil {
		return err
	}

	i.Host.SetStreamHandler(Protocol, i.NetworkHandler)
	//h.Network().Notify()

	addr := fmt.Sprintf("/ip4/%s/tcp/%s", i.Address, i.Port)
	mAddr, err := multiaddr.NewMultiaddr(addr)
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
	log.Debug("I am ", fullAddr)
	return nil
}

//每个连接只会触发一次这个回调函数，之后需要在线程中做收发
func (i *Instance) NetworkHandler(s net.Stream) {
	log.Debug("receive msg from:", s.Conn().RemotePeer().Pretty(), s.Conn().RemoteMultiaddr(), "| topic is:", s.Protocol())

	pub, err := s.Conn().RemotePublicKey().Bytes()
	if err != nil {
		log.Warn(err)
		return
	}
	i.PeerAdd(s.Conn().RemotePeer(), s, s.Conn().RemoteMultiaddr(), crypto.ConfigEncodeKey(pub))

	go i.ReceiveMessage(s)
}

func (i *Instance) ConnectPeer(b64Pub, address, port string) (net.Stream, error) {
	id, err := IdFromPublicKey(b64Pub)
	if err != nil {
		return nil, err
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
	i.PeerAdd(id, s, addr, b64Pub)
	return s, nil
}

func (i *Instance) ReceiveMessage(s net.Stream) {
	reader := io.NewDelimitedReader(s, net.MessageSizeMax)
	for {
		msg := new(pnet.Message)
		err := reader.ReadMsg(msg)
		if err != nil {
			s.Reset()
			i.PeerDel(s.Conn().RemotePeer())
			log.Warn("the peer ", s.Conn().RemotePeer().Pretty(), "is disconnected:", err)
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
	info := i.PeerGet(id)
	if info == nil {
		return errors.New(fmt.Sprintf("the peer %s is not connect", id.Pretty()))
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

func (i *Instance) PeerAdd(id peer.ID, s net.Stream, addr multiaddr.Multiaddr, b64Pub string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	if _, ok := i.Peers[id]; ok {
		return
	}
	i.Peers[id] = PeerInfo{s: s, MultiAddr: addr, PublicKey: b64Pub}
}

func (i *Instance) PeerGet(id peer.ID) *PeerInfo {
	i.lock.RLock()
	defer i.lock.RUnlock()
	if info, ok := i.Peers[id]; ok {
		return &info
	}
	return nil
}

func (i *Instance) PeerDel(id peer.ID) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	if _, ok := i.Peers[id]; ok {
		delete(i.Peers, id)
		return nil
	}
	return errors.New(fmt.Sprintf("can't find stream by id:%s", id))
}

func (i *Instance) ResetStream(s net.Stream) error {
	id := s.Conn().RemotePeer()
	if err := s.Reset(); err != nil {
		return err
	}
	i.PeerDel(id)
	return nil
}
