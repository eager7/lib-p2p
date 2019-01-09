package net

import (
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-kbucket"
	"github.com/libp2p/go-libp2p-peer"
	"time"
)

type RouteTable struct {
	route *kbucket.RoutingTable
}

func RouteInitialize(host host.Host) *RouteTable {
	id := kbucket.ConvertPeerID(host.ID())
	route := kbucket.NewRoutingTable(20, id, time.Minute, host.Peerstore())
	route.PeerAdded = func(id peer.ID) {
		//当更新一个节点信息时的回调
	}
	route.PeerRemoved = func(id peer.ID) {
		//当移除一个节点信息时的回调
	}
	return &RouteTable{route: route}
}

func (r *RouteTable) RouteUpdate(id peer.ID) {
	r.route.Update(id)
}
