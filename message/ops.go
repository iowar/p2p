package message

type Op byte

const (
	UNDEFINED Op = iota
	Ping
	Pong
	GetVersion
	Version
	GetPeerList
	PeerList
)

func (op Op) String() string {
	switch op {
	case Ping:
		return "ping"
	case Pong:
		return "pong"
	case GetVersion:
		return "getversion"
	case Version:
		return "version"
	case GetPeerList:
		return "getpeerlist"
	case PeerList:
		return "peerlist"
	default:
		return "Unknown Op"
	}
}
