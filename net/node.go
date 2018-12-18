package net

import (
	"github.com/eager7/go/mlog"
	"github.com/eager7/lib-p2p/common/errors"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"context"
	"github.com/libp2p/go-libp2p-net"
	"fmt"
	"github.com/multiformats/go-multiaddr"
)

var log = mlog.L

func New() Network {

	return nil
}

func InitNetwork(ctx context.Context, address, port string ) {
	private, public, err := crypto.GenerateKeyPair(crypto.RSA, 1024)
	errors.CheckErrorPanic(err)

	id, err := peer.IDFromPrivateKey(private)
	errors.CheckErrorPanic(err)

	var options []libp2p.Option
	options = append(options, libp2p.Identity(private))
	ps := peerstore.NewPeerstore()
	ps.AddPrivKey(id, private)
	ps.AddPubKey(id, public)
	options = append(options, libp2p.Peerstore(ps))

	if ctx == nil {
		ctx = context.Background()
	}
	h, err := libp2p.New(ctx, options...)
	errors.CheckErrorPanic(err)

	h.SetStreamHandler("/eager7/test/1.0.0", networkHandler)
	//h.Network().Notify()

	addr := fmt.Sprintf("/ip4/%s/tcp/%s", address, port)
	mAddr, err := multiaddr.NewMultiaddr(addr)
	errors.CheckErrorPanic(err)

	err = h.Network().Listen([]multiaddr.Multiaddr{mAddr}...)
	errors.CheckErrorPanic(err)
	log.Debug(h.Network().InterfaceListenAddresses())
}

func networkHandler(s net.Stream) {

}