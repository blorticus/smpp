package main

import (
	"flag"
	"fmt"
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

func (o *outputter) sayThatTransportIsOpen(connectedTransport net.Conn) {
	fmt.Printf("Transport connected from %s to %s\n", connectedTransport.LocalAddr(), connectedTransport.RemoteAddr())
}

func (o *outputter) sayThatBindWasSent() {
	fmt.Printf("Bind was sent to remote peer\n")
}

func (o *outputter) sayThatBindWasReceived() {
	fmt.Printf("Bind received from remote peer\n")
}

func (o *outputter) sayThatSubmitSmWasSent(msgNumber int) {
	fmt.Printf("Sent SubmitSM: number %03d\n", msgNumber)
}

func (o *outputter) sayThatPduWasReceived(receivedPDU *smpp.PDU) {
	fmt.Printf("Received PDU : %s\n", receivedPDU.CommandName())
}

type trasceiverBindInfo struct {
	password   string
	systemID   string
	systemType string
}

type esmeTransceiverApp struct {
	outputter  *outputter
	readBuffer []byte
	submitSm   *smpp.PDU
}

func newEsmeTransceiverApp() *esmeTransceiverApp {
	outputter := &outputter{outputDebug: false}
	return &esmeTransceiverApp{outputter: outputter, readBuffer: make([]byte, 9200), submitSm: generateBaseSubmitSM()}
}

func generateBaseSubmitSM() *smpp.PDU {
	return smpp.NewPDU(smpp.CommandSubmitSm, 0, 0x5e, []*smpp.Parameter{
		smpp.NewFLParameter(uint8(0)),                                  // service type
		smpp.NewFLParameter(uint8(0)),                                  // source_addr_ton
		smpp.NewFLParameter(uint8(1)),                                  // source_addr_npi
		smpp.NewCOctetStringParameter("28809090"),                      // source_addr
		smpp.NewFLParameter(uint8(1)),                                  // dest_addr_ton
		smpp.NewFLParameter(uint8(1)),                                  // dest_addr_npi
		smpp.NewCOctetStringParameter("13139591463"),                   // destination_addr
		smpp.NewFLParameter(uint8(0)),                                  // esm_class
		smpp.NewFLParameter(uint8(0)),                                  // protocol_id
		smpp.NewFLParameter(uint8(0)),                                  // priority_flag
		smpp.NewFLParameter(uint8(0)),                                  // schedule_delivery_time
		smpp.NewFLParameter(uint8(0)),                                  // validity_period
		smpp.NewFLParameter(uint8(0)),                                  // registered_delivery
		smpp.NewFLParameter(uint8(0)),                                  // replace_if_flag_present
		smpp.NewFLParameter(uint8(0xf0)),                               // data_coding
		smpp.NewFLParameter(uint8(0)),                                  // sm_default_msg_id
		smpp.NewFLParameter(uint8(0x8d)),                               // sm_length
		smpp.NewOctetStringFromString("This is a test short message."), // short_message
	}, []*smpp.Parameter{})
}

func (app *esmeTransceiverApp) processCommandLineParameters() (peerIP string, peerPort int, systemID string, messagesToSend int) {
	flag.StringVar(&peerIP, "addr", "127.0.0.1", "IP address of peer")
	flag.IntVar(&peerPort, "port", 2775, "Listening TCP port for peer")
	flag.StringVar(&systemID, "id", "esme01", "ESME system-id")
	flag.IntVar(&messagesToSend, "count", 10, "Number of messages to send")

	os.Args[0] = path.Base(os.Args[0])

	flag.Parse()

	if peerPort < 1 || peerPort > 65535 {
		app.outputter.die("Port must be between 1 and 65535, inclusive")
	}

	return peerIP, peerPort, systemID, messagesToSend
}

func (app *esmeTransceiverApp) connectToSMSC(peerIP string, peerPort int) net.Conn {
	peerAddr := peerIP + ":" + strconv.Itoa(peerPort)

	conn, err := net.Dial("tcp", peerAddr)

	if err != nil {
		app.outputter.die("Failed to establish outbound connection to [%s]: %s\n", peerAddr, err)
	}

	return conn
}

func (app *esmeTransceiverApp) sendBindToSMSC(transportConnection net.Conn, bindInfo trasceiverBindInfo) {
	transceiverBindMessage := smpp.NewPDU(smpp.CommandBindTransceiver, 0, 1, []*smpp.Parameter{
		smpp.NewCOctetStringParameter(bindInfo.systemID),
		smpp.NewCOctetStringParameter(bindInfo.password),
		smpp.NewCOctetStringParameter(bindInfo.systemType),
		smpp.NewFLParameter(uint8(0x34)),  // interface_version
		smpp.NewFLParameter(uint8(0x0)),   // addr_ton
		smpp.NewFLParameter(uint8(0x0)),   // addr_npi
		smpp.NewCOctetStringParameter(""), // address_range
	}, []*smpp.Parameter{})

	app.sendPDU(transportConnection, transceiverBindMessage)
}

func (app *esmeTransceiverApp) receiveBindFromSMSC(transportConnection net.Conn) {
	_ = app.recvPDU(transportConnection)
}

func (app *esmeTransceiverApp) sendPDU(conn net.Conn, pdu *smpp.PDU) {
	encoded, err := pdu.Encode()

	if err != nil {
		app.outputter.die("Failed to encode PDU: %s", err)
	}

	_, err = conn.Write(encoded)

	if err != nil {
		app.outputter.die("Failed to write: %s", err)
	}
}

func (app *esmeTransceiverApp) recvPDU(conn net.Conn) *smpp.PDU {
	_, err := conn.Read(app.readBuffer)

	if err != nil {
		app.outputter.die("Failed on read: %s", err)
	}

	pdu, err := smpp.DecodePDU(app.readBuffer)

	if err != nil {
		app.outputter.die("Failed to encode incoming PDU: %s", err)
	}

	return pdu
}

func (app *esmeTransceiverApp) generateSubmitSmWithFromAndMessageNumber(from string, messageNumber int) *smpp.PDU {
	shortMsgString := fmt.Sprintf("FROM: %s | MSG NUMBER: %03d", from, messageNumber)
	app.submitSm.MandatoryParameters[16] = smpp.NewFLParameter(uint8(len(shortMsgString)))
	app.submitSm.MandatoryParameters[17] = smpp.NewOctetStringFromString(shortMsgString)
	return app.submitSm
}

func main() {
	app := newEsmeTransceiverApp()

	peerIP, peerPort, systemID, messagesToSend := app.processCommandLineParameters()

	conn := app.connectToSMSC(peerIP, peerPort)
	defer conn.Close()
	app.outputter.sayThatTransportIsOpen(conn)

	app.sendBindToSMSC(conn, trasceiverBindInfo{systemID: systemID, password: "password", systemType: "generic"})
	app.outputter.sayThatBindWasSent()

	app.receiveBindFromSMSC(conn)
	app.outputter.sayThatBindWasReceived()

	for i := 0; i < messagesToSend; i++ {
		submitSmPDU := app.generateSubmitSmWithFromAndMessageNumber(systemID, i+1)
		app.sendPDU(conn, submitSmPDU)
		app.outputter.sayThatSubmitSmWasSent(i + 1)

		responsePDU := app.recvPDU(conn)
		app.outputter.sayThatPduWasReceived(responsePDU)
	}
}
