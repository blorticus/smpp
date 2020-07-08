package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strconv"

	smpp "github.com/blorticus/smpp-go"
)

type outputter struct {
	outputDebug bool
}

func (o *outputter) die(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "[FATAL] ")
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func (o *outputter) debug(format string, a ...interface{}) {
	if o.outputDebug {
		fmt.Print("[DEBUG] ")
		fmt.Printf(format, a...)
		fmt.Print("\n")
	}
}

func (o *outputter) sayThatTransportWasReceived(transportConnection net.Conn) {
	fmt.Printf("Transport Opened From: %s\n", transportConnection.RemoteAddr().String())
}

func (o *outputter) sayThatBindWasReceived(bindPDU *smpp.PDU) {
	fmt.Printf("transceiver-bind Received from ESME\n")
}

func (o *outputter) sayThatBindResponseWasSent() {
	fmt.Printf("Sent transceiver-bind-resp to ESME\n\n")
}

func (o *outputter) sayThatPeerClosedTransport() {
	fmt.Printf("Peer closed transport\n")
}

func (o *outputter) sayThatPduWasReceived(pdu *smpp.PDU) {
	if pdu.CommandID == smpp.CommandSubmitSm {
		fmt.Printf("Received submit-sm from ESME : Message = [%s]\n", pdu.MandatoryParameters[17].Value.([]uint8))
	} else {
		fmt.Printf("Received PDU from ESME       : Command = %-21.21s\n", pdu.CommandName())
	}
}

func (o *outputter) sayThatUnexpectedPduTypeWasReceived() {
	fmt.Printf("PDU type is unexpected.  Closing connection.\n")
}

func (o *outputter) sayThatSubmitSmResponseWasSent(systemID string, responseNumber uint32) {
	fmt.Printf("Sent submit-sm-resp to ESME  : Id = %s-%03d\n", systemID, responseNumber)
}

func (o *outputter) sayThatWeAreWaitingForConnection(bindIP string, bindPort int) {
	fmt.Printf("Listening for ESME connection on %s:%d\n", bindIP, bindPort)
}

type trasceiverBindResponseInfo struct {
	systemID       string
	sequenceNumber uint32
}

type smscTransceiverApp struct {
	outputter  *outputter
	readBuffer []byte
}

func newSmscTransceiverApp() *smscTransceiverApp {
	outputter := &outputter{outputDebug: false}
	return &smscTransceiverApp{outputter: outputter, readBuffer: make([]byte, 9200)}
}

func (app *smscTransceiverApp) processCommandLineParameters() (bindIP string, bindPort int, systemID string) {
	flag.StringVar(&bindIP, "addr", "127.0.0.1", "IP address on which to listen for incoming connection")
	flag.IntVar(&bindPort, "port", 2775, "Listening TCP port for incoming connection")
	flag.StringVar(&systemID, "id", "smsc01", "SMSC system-id")

	os.Args[0] = path.Base(os.Args[0])

	flag.Parse()

	if bindPort < 1 || bindPort > 65535 {
		app.outputter.die("Port must be between 1 and 65535, inclusive")
	}

	return bindIP, bindPort, systemID
}

func (app *smscTransceiverApp) awaitIncomingEMSEConnection(bindIP string, bindPort int) net.Conn {
	bindAddr := bindIP + ":" + strconv.Itoa(bindPort)

	listener, err := net.Listen("tcp", bindAddr)

	if err != nil {
		app.outputter.die("Failed to start listener: %s", err)
	}

	conn, err := listener.Accept()

	if err != nil {
		app.outputter.die("Failed to Accept incoming socket: %s", err)
	}

	listener.Close()

	return conn
}

func (app *smscTransceiverApp) receiveBindFromESME(fromTransport net.Conn) *smpp.PDU {
	pdu := app.recvPDU(fromTransport)

	if pdu.CommandID != smpp.CommandBindTransceiver {
		app.outputter.die("Expected BindTransceiver from ESME, received %s", pdu.CommandName())
	}

	return pdu
}

func (app *smscTransceiverApp) sendBindResponseToESME(onTransport net.Conn, bindInfo trasceiverBindResponseInfo) {
	transceiverBindResponseMessage := smpp.NewPDU(smpp.CommandBindTransceiverResp, 0, bindInfo.sequenceNumber, []*smpp.Parameter{
		smpp.NewCOctetStringParameter(bindInfo.systemID),
	}, []*smpp.Parameter{})

	app.sendPDU(onTransport, transceiverBindResponseMessage)
}

func (app *smscTransceiverApp) sendPDU(conn net.Conn, pdu *smpp.PDU) {
	encoded, err := pdu.Encode()

	if err != nil {
		app.outputter.die("Failed to encode PDU: %s\n", err)
	}

	_, err = conn.Write(encoded)

	if err != nil {
		app.outputter.die("Failed to write: %s\n", err)
	}
}

func (app *smscTransceiverApp) recvPDU(conn net.Conn) *smpp.PDU {
	_, err := conn.Read(app.readBuffer)

	if err != nil {
		if err == io.EOF {
			return nil
		}

		app.outputter.die("Failed on read: %s", err)
	}

	pdu, err := smpp.DecodePDU(app.readBuffer)

	if err != nil {
		app.outputter.die("Failed to encode incoming PDU: %s", err)
	}

	return pdu
}

func (app *smscTransceiverApp) sendSubmitSmResponse(submitSmPDU *smpp.PDU, transportToESME net.Conn, systemID string, responseNumber uint32) {
	respPDU := smpp.NewPDU(smpp.CommandBindTransceiverResp, 0, submitSmPDU.SequenceNumber, []*smpp.Parameter{
		smpp.NewCOctetStringParameter(fmt.Sprintf("%s-%03d", systemID, responseNumber)),
	}, []*smpp.Parameter{})

	app.sendPDU(transportToESME, respPDU)
}

func (app *smscTransceiverApp) closeConnectionToESME(transport net.Conn) {
	transport.Close()
}

func main() {
	app := newSmscTransceiverApp()

	bindIP, bindPort, systemID := app.processCommandLineParameters()

	app.outputter.sayThatWeAreWaitingForConnection(bindIP, bindPort)

	transport := app.awaitIncomingEMSEConnection(bindIP, bindPort)
	defer transport.Close()
	app.outputter.sayThatTransportWasReceived(transport)

	bindPDU := app.receiveBindFromESME(transport)
	app.outputter.sayThatBindWasReceived(bindPDU)

	app.sendBindResponseToESME(transport, trasceiverBindResponseInfo{systemID: systemID, sequenceNumber: bindPDU.SequenceNumber})
	app.outputter.sayThatBindResponseWasSent()

	for responseNumber := uint32(1); true; responseNumber++ {
		pdu := app.recvPDU(transport)

		if pdu == nil {
			app.outputter.sayThatPeerClosedTransport()
			os.Exit(0)
		}

		app.outputter.sayThatPduWasReceived(pdu)

		switch pdu.CommandID {
		case smpp.CommandSubmitSm:
			app.sendSubmitSmResponse(pdu, transport, systemID, responseNumber)
			app.outputter.sayThatSubmitSmResponseWasSent(systemID, responseNumber)
		default:
			app.outputter.sayThatUnexpectedPduTypeWasReceived()
			app.closeConnectionToESME(transport)
			os.Exit(1)
		}
	}
}
