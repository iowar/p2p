package config

import "time"

// p2p last version
const Version = "0.0.1"

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other peers.
const DiscoveryServiceTag = "XXXMEV-0.0.1"

// Topic name for DHT
const Topic = "XXXMEV-TEST-43"

// PubSub topic
const PubSubTopic = "XXXMEVTOPIC-0.0.1"
