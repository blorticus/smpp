package smpp

import "net"

type peerStates int

const (
	peerDisconnected peerStates = iota
	peerUnbound
	peerBound
)

// Peer represents a peer for an SMPP entity (an ESME or SMSC)
type Peer struct {
	connectionToRemotePeer *net.TCPConn
	state                  peerStates
}

// NewPeerWithConnection instantiates a Peer object, providing it with an already
// created connection.  This can be useful if you wish to use local binds for a
// Peer connection
func NewPeerWithConnection(c *net.TCPConn) *Peer {
	return &Peer{connectionToRemotePeer: c, state: peerUnbound}
}
