package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"smpp"
)

func sendPDU(conn net.Conn, pdu *smpp.PDU, logger *log.Logger) {
	encoded, err := pdu.Encode()

	if err != nil {
		logger.Fatalln("Failed to encode PDU: ", err)
	}

	_, err = conn.Write(encoded)

	if err != nil {
		logger.Fatalln("Failed to write: ", err)
	}
}

func recvPDU(conn net.Conn, buf *[]byte, logger *log.Logger) *smpp.PDU {
	_, err := conn.Read(*buf)

	if err != nil {
		logger.Fatalln("Failed on read: ", err)
	}

	encoded, err := smpp.DecodePDU(*buf)

	if err != nil {
		logger.Fatalln("Failed to encode incoming PDU: ", err)
	}

	return encoded
}

func main() {
	peerIP := "127.0.0.1"
	peerPort := "2775"
	peerAddr := peerIP + ":" + peerPort

	logger := log.New(os.Stderr, filepath.Base(os.Args[0])+": ", 0)

	logger.Println("Starting client connection to [", peerAddr, "]")

	conn, err := net.Dial("tcp", peerAddr)

	if err != nil {
		logger.Fatalln("Failed to establish outbound connection to [", peerAddr, "]: ", err)
	}

	defer conn.Close()

	readbuf := make([]byte, 4096)

	logger.Println("Sending bind-transmitter PDU")

	bindTransmitter := smpp.NewPDU(smpp.CommandBindTransmitter, 0, 1, []*smpp.Parameter{
		smpp.NewCOctetStringParameter("esme01"),   // system_id
		smpp.NewCOctetStringParameter("password"), // password
		smpp.NewCOctetStringParameter("generic"),  // system_type
		smpp.NewFLParameter(uint8(0x34)),          // interface_version
		smpp.NewFLParameter(uint8(0x0)),           // addr_ton
		smpp.NewFLParameter(uint8(0x0)),           // addr_npi
		smpp.NewCOctetStringParameter(""),         // address_range
	}, []*smpp.Parameter{})

	sendPDU(conn, bindTransmitter, logger)

	logger.Println("Waiting for bind-transmitter-resp")

	_ = recvPDU(conn, &readbuf, logger)

	logger.Println("PDU received")

	submitSmPDU := smpp.NewPDU(smpp.CommandSubmitSm, 0, 0x5e, []*smpp.Parameter{
		smpp.NewFLParameter(uint8(0)),
		smpp.NewFLParameter(uint8(0)),
		smpp.NewFLParameter(uint8(1)),
		smpp.NewCOctetStringParameter("28809090"),
		smpp.NewFLParameter(uint8(1)),
		smpp.NewFLParameter(uint8(1)),
		smpp.NewCOctetStringParameter("13139591463"),
		smpp.NewFLParameter(uint8(0)),
		smpp.NewFLParameter(uint8(0)),
		smpp.NewFLParameter(uint8(0)),
		smpp.NewFLParameter(uint8(0)),
		smpp.NewCOctetStringParameter("000000000500000R"),
		smpp.NewFLParameter(uint8(0)),
		smpp.NewFLParameter(uint8(0)),
		smpp.NewFLParameter(uint8(0xf0)),
		smpp.NewFLParameter(uint8(0)),
		smpp.NewFLParameter(uint8(0x8d)),
		smpp.NewOctetStringFromString("This is a test short message, though it is somewhat longer than short, being > 50 characters! Don't get excited :@ :# :$ :% :^) emoji like..."),
	}, []*smpp.Parameter{
		smpp.NewTLVParameter(0x020c, uint16(5)),
		smpp.NewTLVParameter(0x020e, uint8(2)),
		smpp.NewTLVParameter(0x020f, uint8(1)),
	})

	encoded, _ := submitSmPDU.Encode()

	logger.Printf("encoded length = (%d)\n", len(encoded))

	logger.Println("Attempting to write first submit-sm PDU")

	sendPDU(conn, submitSmPDU, logger)

	logger.Println("PDU sent; waiting response")

	_ = recvPDU(conn, &readbuf, logger)

	logger.Println("Received PDU")

	logger.Println("Closing peer connection")
	logger.Println("Done")
}
