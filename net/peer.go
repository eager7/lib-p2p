package net

import (
	"fmt"
	"github.com/eager7/lib-p2p/common/errors"

	"sync"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
	"gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	"gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
)

type Peer struct {
	ID        peer.ID
	s         net.Stream
	PeerInfo  peerstore.PeerInfo
}

type PeerMap struct {
	Peers map[peer.ID]Peer
	P sync.Map
	lock  sync.RWMutex
}

func (p *PeerMap) Initialize() {
	p.Peers = make(map[peer.ID]Peer)
}

func (p *PeerMap) Add(id peer.ID, s net.Stream, addr multiaddr.Multiaddr) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, ok := p.Peers[id]; ok {
		return
	}
	//peerInfo := peerstore.PeerInfo{ID: id, Addrs: []multiaddr.Multiaddr{addr}}
	//p.Peers[id] = Peer{ID: id, s: s, PeerInfo: peerInfo}
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
