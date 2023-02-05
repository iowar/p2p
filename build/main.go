package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iowar/p2p/config"
	"github.com/iowar/p2p/discover"
	"github.com/iowar/p2p/message"
	pubsubio "github.com/iowar/p2p/pubsub"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set your own keypair
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	//var idht *dht.IpfsDHT

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		panic(err)
	}
	host, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0", // regular tcp connections
			//"/ip4/0.0.0.0/tcp/0/ws", // websocket endpoint
			//"/ip4/0.0.0.0/udp/0/quic", // a UDP endpoint for the QUIC transport
		),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		//libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
		//	idht, err = dht.New(ctx, h)
		//	return idht, err
		//}),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
	)
	if err != nil {
		panic(err)
	}
	defer host.Close()

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}
	topic, err := ps.Join(config.PubSubTopic)
	if err != nil {
		panic(err)
	}
	defer topic.Close()
	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	// message builder (outbound)
	msgBuilder := message.NewMsgBuilder()

	// create newServer
	pubsubio.NewServer(ctx, host, topic, sub, msgBuilder)

	log.Printf("Connect to me on:\033[31m")
	for _, addr := range host.Addrs() {
		log.Printf("%s/p2p/%s", addr, host.ID().Pretty())
	}
	log.Printf("------------------------------------------------------------\033[0m")

	// trigger discovery package
	discovery := discover.NewDiscovery(host, ctx)

	// setup local mDNS discovery
	if err := discovery.StartMdnsDiscovery(); err != nil {
		panic(err)
	}

	// connect default nodes
	discovery.ConnectBootstrap()

	// It has a weak infrastructure
	// TODO use own techniques
	discovery.StartDhtRouting()

	time.Sleep(time.Millisecond * 1314)
	msgBytes := msgBuilder.OutboundPing()
	topic.Publish(ctx, msgBytes)

	done := make(chan struct{}, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	case <-done:
		host.Close()
	}

}
