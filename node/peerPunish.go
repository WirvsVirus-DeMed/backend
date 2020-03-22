package node

import (
	"net"
	"time"
)

type PeerPunish struct {
	Peer  net.TCPAddr
	until time.Time
}
