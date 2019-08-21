package smpp

import (
	"fmt"
	"net"
)

// ESME represents an ESME, which initiates connection to one or more SMSCs
type ESME struct {
}

// ConnectToPeer connects a transport (TCP) to a remote peer
func (esme *ESME) ConnectToPeer(remoteAddr net.IP, remotePort uint16) (peer *Peer, err error) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: remoteAddr, Port: int(remotePort), Zone: ""})

	if err != nil {
		return nil, err
	}

	return NewPeerWithConnection(conn), nil
}

// BindToPeer establishes a bind with a remote peer to which a transport connection is already completed
func (esme *ESME) BindToPeer(peer *Peer, bind BindInfo) error {
	if peer.state == peerDisconnected {
		return fmt.Errorf("Peer has no connected transport")
	}

	if peer.state == peerBound {
		return fmt.Errorf("Peer is already bound")
	}

	return nil
}

// StartListenLoop should be run in a goroutine, and listens for incoming messages from peers
func (esme *ESME) StartListenLoop() {

}

// SendMessageToPeer sends a message to a bound peer
func (esme *ESME) SendMessageToPeer(peer *Peer, pdu *PDU) error {
	return nil
}
