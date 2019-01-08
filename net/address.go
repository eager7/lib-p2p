package net

import (
	"fmt"
	"github.com/eager7/lib-p2p/common/errors"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-peer"
	"strconv"
	"strings"
)

func IdFromPublicKey(pubKey string) (peer.ID, error) {
	key, err := crypto.ConfigDecodeKey(pubKey)
	if err != nil {
		return "", errors.New(err.Error())
	}
	pk, err := crypto.UnmarshalPublicKey(key)
	if err != nil {
		return "", errors.New(err.Error())
	}
	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return id, nil
}

func NewAddrInfo(ip, port string) (addr string) {
	tcpPort, _ := strconv.Atoi(port)
	if strings.Contains(ip, ":") {
		addr = fmt.Sprintf("/ip6/%s/tcp/%d", ip, tcpPort)
	} else {
		addr = fmt.Sprintf("/ip4/%s/tcp/%d", ip, tcpPort)
	}
	return
}
