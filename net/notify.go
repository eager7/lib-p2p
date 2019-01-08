package net

import (
	"github.com/libp2p/go-libp2p-net"
	"github.com/multiformats/go-multiaddr"
)

func (i *Instance) Listen(n net.Network, a multiaddr.Multiaddr)      { log.Debug("Listen") }
func (i *Instance) ListenClose(n net.Network, a multiaddr.Multiaddr) { log.Debug("ListenClose") }
func (i *Instance) Connected(n net.Network, v net.Conn) {
	log.Debug("Connected ID:", v.RemotePeer().Pretty(), "Addr:", n.Peerstore().Addrs(v.RemotePeer()))
}
func (i *Instance) Disconnected(n net.Network, v net.Conn) {
	log.Debug("Disconnected:", n.Peerstore().Addrs(v.RemotePeer()))

}
func (i *Instance) OpenedStream(n net.Network, v net.Stream) {
	//id := v.Conn().RemotePeer()
	//addresses :=  n.Peerstore().Addrs(v.Conn().RemotePeer())
	//log.Debug("OpenedStream", id.Pretty(), addresses, v.Conn().RemoteMultiaddr(), v)
}
func (i *Instance) ClosedStream(n net.Network, v net.Stream) {
	//log.Debug("ClosedStream", n.Peerstore().Addrs(v.Conn().RemotePeer()), v)
}
