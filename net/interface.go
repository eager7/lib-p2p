package net

import "github.com/eager7/lib-p2p/message"

type Network interface {
	SendMessage(b64Pub, address, port string, payload *pnet.Message) error
}

type Payload interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}