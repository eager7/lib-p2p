package net

import (
	"fmt"
	"github.com/libp2p/go-libp2p-net"
	"github.com/multiformats/go-multiaddr"
)

func (i *Instance) Listen(n net.Network, a multiaddr.Multiaddr)      { fmt.Println("Listen") }
func (i *Instance) ListenClose(n net.Network, a multiaddr.Multiaddr) { fmt.Println("ListenClose") }
func (i *Instance) Connected(n net.Network, v net.Conn) {
	fmt.Println("Connected:", n.Peerstore().Addrs(v.RemotePeer()))
}
func (i *Instance) Disconnected(n net.Network, v net.Conn) {
	fmt.Println("Disconnected:", n.Peerstore().Addrs(v.RemotePeer()))

}
func (i *Instance) OpenedStream(n net.Network, v net.Stream) {
	fmt.Println("OpenedStream", n.Peerstore().Addrs(v.Conn().RemotePeer()))
}
func (i *Instance) ClosedStream(n net.Network, v net.Stream) {
	fmt.Println("ClosedStream", n.Peerstore().Addrs(v.Conn().RemotePeer()))
}
