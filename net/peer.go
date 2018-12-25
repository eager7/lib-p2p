package net

import (
	"fmt"
	"github.com/eager7/lib-p2p/common/errors"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"sync"
)

type Peer struct {
	ID        peer.ID
	s         net.Stream
	PeerInfo  peerstore.PeerInfo
	PublicKey string
}

type PeerMap struct {
	Peers map[peer.ID]Peer
	P sync.Map
	lock  sync.RWMutex
}

func (p *PeerMap) Initialize() {
	p.Peers = make(map[peer.ID]Peer)
}

func (p *PeerMap) Add(id peer.ID, s net.Stream, addr multiaddr.Multiaddr, b64Pub string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, ok := p.Peers[id]; ok {
		return
	}
	peerInfo := peerstore.PeerInfo{ID: id, Addrs: []multiaddr.Multiaddr{addr}}
	p.Peers[id] = Peer{ID: id, s: s, PeerInfo: peerInfo, PublicKey: b64Pub}
}

func (p *PeerMap) Del(id peer.ID) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, ok := p.Peers[id]; ok {
		delete(p.Peers, id)
		return nil
	}
	return errors.New(fmt.Sprintf("can't find stream by id:%s", id))
}

func (p *PeerMap) Get(id peer.ID) *Peer {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if info, ok := p.Peers[id]; ok {
		return &info
	}
	return nil
}

func (p *PeerMap) Iterator() <-chan Peer {
	channel := make(chan Peer)
	go func() {
		p.lock.RLock()
		defer p.lock.RUnlock()
		for _, v := range p.Peers {
			channel <- v
		}
		close(channel)
	}()
	return channel
}
