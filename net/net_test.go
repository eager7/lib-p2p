package net_test

import (
	"testing"
	"github.com/libp2p/go-libp2p-peer"
	"fmt"
	"github.com/eager7/lib-p2p/net"
)

func TestPeerMap_Iterator(t *testing.T) {
	p := new(net.PeerMap)
	p.Initialize()
	p.Add(peer.ID("test1"), nil, nil)
	p.Add(peer.ID("test2"), nil, nil)
	p.Add(peer.ID("test3"), nil, nil)
	p.Add(peer.ID("test4"), nil, nil)
	p.Add(peer.ID("test5"), nil, nil)

	fmt.Println(p)

	for v := range p.Iterator() {
		fmt.Println(v)
	}
}
