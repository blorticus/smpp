package smpp

import (
	"io"
	"net"
)

// EventMessageType is an enum type for the EventMessage Type field
type EventMessageType int

const (
	// PeerConnectionOpened raised when SMPP peer connection completes
	PeerConnectionOpened EventMessageType = iota
	// PeerConnectionClosed raised when SMPP peer connection completes
	PeerConnectionClosed
	// PeerTransportOpened raised when transport toward SMPP peer opens
	PeerTransportOpened
	// PeerTransportClosed raised when transport toward SMPP peer closes
	PeerTransportClosed
	// PeerReadError raised when an error occurred on transport read attempt
	PeerReadError
)

// EventMessage is a message between an Instance and spawned goroutines
type EventMessage struct {
	Type          EventMessageType
	ReportingPeer *Peer
	Error         error
}

// PeerControlMessageType is an enum type for the PeerControlMessage Type field
type PeerControlMessageType int

const ()

// PeerControlMessage is used by Instances to control Peer objects
type PeerControlMessage struct {
	Type PeerControlMessageType
}

type peerState int

// const (
// 	closed = iota
// 	transportOpen
// 	connected
// )

// Peer represents a connection toward an SMPP peer
type Peer struct {
	SystemID string
	Password string
	//	state    peerState
}

// NewPeer creates a new peer object
// func NewPeer(systemID string, password string) *Peer {
// 	return &Peer{SystemID: systemID, Password: password, state: closed}
// }

// Start launches a Peer
func (peer *Peer) Start(conn net.Conn, events chan<- (*EventMessage), cntrl <-chan (*PeerControlMessage)) {
	defer conn.Close()

	var bytesRead int
	var err error
	readBuf := make([]byte, 65515)

	for {
		bytesRead, err = conn.Read(readBuf)

		if err != nil {
			if err == io.EOF {
				events <- &EventMessage{Type: PeerConnectionClosed, ReportingPeer: peer, Error: err}
				return
			}
		}
	}
}
