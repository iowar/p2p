package message

import (
	"encoding/json"
	"sync"

	"github.com/iowar/p2p/config"
)

type Builder struct {
	sync.RWMutex
}

func NewMsgBuilder() *Builder {
	return &Builder{}
}

func (mb *Builder) OutboundPing() []byte {
	mb.Lock()
	defer mb.Unlock()

	OutboundMessage.SetOp(Ping)
	OutboundMessage.SetTimestamp()
	msgBytes, _ := json.Marshal(OutboundMessage)
	return msgBytes
}

func (mb *Builder) OutboundPong() []byte {
	mb.Lock()
	defer mb.Unlock()

	OutboundMessage.SetOp(Pong)
	OutboundMessage.SetTimestamp()
	msgBytes, _ := json.Marshal(OutboundMessage)
	return msgBytes
}

func (mb *Builder) OutboundGetVersion() []byte {
	mb.Lock()
	defer mb.Unlock()

	OutboundMessage.SetOp(Version)
	OutboundMessage.SetTimestamp()
	msgBytes, _ := json.Marshal(OutboundMessage)
	return msgBytes
}

func (mb *Builder) OutboundVersion() []byte {
	mb.Lock()
	defer mb.Unlock()

	OutboundMessage.SetOp(Version)
	OutboundMessage.SetTimestamp()
	OutboundMessage.SetBytes([]byte(config.Version))
	msgBytes, _ := json.Marshal(OutboundMessage)
	return msgBytes
}

func (mb *Builder) OutboundGetPeerList() []byte {
	mb.Lock()
	defer mb.Unlock()

	OutboundMessage.SetOp(GetPeerList)
	OutboundMessage.SetTimestamp()
	msgBytes, _ := json.Marshal(OutboundMessage)
	return msgBytes
}

func (mb *Builder) OutboundPeerList() []byte {
	mb.Lock()
	defer mb.Unlock()

	OutboundMessage.SetOp(PeerList)
	OutboundMessage.SetTimestamp()
	//TODO
	msgBytes, _ := json.Marshal(OutboundMessage)
	return msgBytes
}
