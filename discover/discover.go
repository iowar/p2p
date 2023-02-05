package discover

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/iowar/p2p/config"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

// optional discovery
type Discovery struct {
	h      host.Host
	ctx    context.Context
	notify *discoveryNotifee
	idht   *dht.IpfsDHT
	//...
}

// Create discover method
func NewDiscovery(h host.Host, ctx context.Context) *Discovery {
	notify := &discoveryNotifee{
		h:   h,
		ctx: ctx,
	}

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	idht, err := dht.New(ctx, h)
	if err != nil {
		panic(err)
	}
	if err = idht.Bootstrap(ctx); err != nil {
		panic(err)
	}

	return &Discovery{
		h:      h,
		ctx:    ctx,
		notify: notify,
		idht:   idht,
	}
}

// StartDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func (d *Discovery) StartMdnsDiscovery() error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(d.h, config.DiscoveryServiceTag, d.notify)
	return s.Start()
}

func (d *Discovery) StartDhtRouting() {
	go d.discoverPeersWithRendezvous()
}

func (d *Discovery) ConnectBootstrap() {
	var wg sync.WaitGroup
	for _, peerAddr := range config.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := d.h.Connect(d.ctx, *peerinfo); err != nil {
				fmt.Println("Bootstrap warning:", err)
			} else {
				log.Println("\033[31m Connection established with bootstrap node: \033[0m", *peerinfo)
			}
		}()
	}

	wg.Wait()
	return
}

// works for a rendezvous peer discovery !!!
// TODO
func (d *Discovery) discoverPeersWithRendezvous() {

	routingDiscovery := drouting.NewRoutingDiscovery(d.idht)
	dutil.Advertise(d.ctx, routingDiscovery, config.Topic)

	// Look for others who have announced and attempt to connect to them
	log.Println("[rendezvous] Searching for peers...")
	for i := 0; i < 90; i++ {
		time.Sleep(time.Second)
		peerChan, err := routingDiscovery.FindPeers(d.ctx, config.Topic)
		if err != nil {
			panic(err)
		}

		for peer := range peerChan {
			if peer.ID == d.h.ID() {
				continue // No self connection
			}
			err := d.h.Connect(d.ctx, peer)
			if err != nil {
				log.Println("Failed connecting to ", peer.ID.Pretty(), ", error:", err)
			} else {
				d.notify.HandlePeerFound(peer)
				break
			}
		}
	}
	log.Println("[rendezvous] Peer discovery complete")
}
