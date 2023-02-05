package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/iowar/p2p/message"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type PubSubServer struct {
	ctx   context.Context
	self  peer.ID
	host  host.Host
	topic *pubsub.Topic
	sub   *pubsub.Subscription
	msb   *message.Builder
}

func (pss PubSubServer) ID() peer.ID { return pss.self }

func NewServer(ctx context.Context, host host.Host, topic *pubsub.Topic, sub *pubsub.Subscription, msb *message.Builder) *PubSubServer {
	var pss = PubSubServer{
		ctx:   ctx,
		self:  host.ID(),
		host:  host,
		topic: topic,
		sub:   sub,
		msb:   msb,
	}

	go pss.Protocol()
	return &pss
}

func (pss *PubSubServer) Protocol() {
	log.Printf("Pubsub protocol started...")
	defer pss.sub.Cancel()
	for {
		msg, err := pss.sub.Next(pss.ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		err = json.Unmarshal(msg.Data, message.InboundMessage)
		if err != nil {
			log.Println(err)
		}

		// check self
		if strings.Compare(pss.self.String(), msg.ReceivedFrom.String()) == 0 {
			continue
		}

		message.InboundMessage.SetPeer(msg.ReceivedFrom)

		copiedMsg := *message.InboundMessage
		go func() {

			log.Printf("[%s] from %s", copiedMsg.Op(), copiedMsg.Peer())
			switch copiedMsg.Op() {

			case message.Ping:
				pss.optPing()
			case message.Pong:
				pss.optPong()
			case message.GetVersion:
				pss.optGetVersion()
			case message.Version:
				pss.optVersion()
			case message.GetPeerList:
				pss.optGetPeerList()
			case message.PeerList:
				pss.optPeerList()
			default:
				log.Printf("Unknown option!")
			}
		}()
	}
}

func (pss *PubSubServer) optPing() {
	// create pong message
	msgBytes := pss.msb.OutboundPong()

	pss.topic.Publish(pss.ctx, msgBytes)
}

func (pss *PubSubServer) optPong() {
	// NULL
	// ...
}

func (pss *PubSubServer) optGetVersion() {
	// NULL
	// ...
}

func (pss *PubSubServer) optVersion() {
	// NULL
	// ...
}

func (pss *PubSubServer) optGetPeerList() {
	// NULL
	// ...
}

func (pss *PubSubServer) optPeerList() {
	// NULL
	// ...
}
