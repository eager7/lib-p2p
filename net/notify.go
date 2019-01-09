package net

import (
	"github.com/libp2p/go-libp2p-net"
	"github.com/multiformats/go-multiaddr"
)

func (i *Instance) Listen(n net.Network, a multiaddr.Multiaddr)      { log.Debug("Listen") }
func (i *Instance) ListenClose(n net.Network, a multiaddr.Multiaddr) { log.Debug("ListenClose") }
func (i *Instance) Connected(n net.Network, v net.Conn) {
	id := v.RemotePeer()
	log.Debug("Connected ID:", id.Pretty(), "Addr:", n.Peerstore().Addrs(v.RemotePeer()))
	if i.Host.Network().Connectedness(id) == net.Connected {
		i.RouteTable.RouteUpdate(id)
	}
}
func (i *Instance) Disconnected(n net.Network, v net.Conn) {
	log.Debug("Disconnected:", n.Peerstore().Addrs(v.RemotePeer()))
	i.RouteTable.RouteRemove(v.RemotePeer())

}
func (i *Instance) OpenedStream(n net.Network, v net.Stream) {}
func (i *Instance) ClosedStream(n net.Network, v net.Stream) {}
