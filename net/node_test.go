package net_test

import (
	"testing"
	"github.com/eager7/lib-p2p/net"
)

func TestNew(t *testing.T) {
	net.InitNetwork(nil, "127.0.0.1", "9001")
}
