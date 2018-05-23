package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"smpp"
)

func main() {
	peerIP := "127.0.0.1"
	peerPort := "2775"
	peerAddr := peerIP + ":" + peerPort

	logger := log.New(os.Stderr, filepath.Base(os.Args[0])+": ", 0)

	logger.Println("Creating PDU object")

	bindTransmitter := smpp.NewPDU(smpp.CommandBindTransmitter, 0, 1, []*smpp.Parameter{
		smpp.NewCOctetStringParameter("esme01"),   // system_id
		smpp.NewCOctetStringParameter("password"), // password
		smpp.NewCOctetStringParameter("generic"),  // system_type
		smpp.NewFLParameter(uint8(0x34)),          // interface_version
		smpp.NewFLParameter(uint8(0x0)),           // addr_ton
		smpp.NewFLParameter(uint8(0x0)),           // addr_npi
		smpp.NewCOctetStringParameter(""),         // address_range
	}, []*smpp.Parameter{})

	logger.Println("Starting client connection to [", peerAddr, "]")

	conn, err := net.Dial("tcp", peerAddr)

	if err != nil {
		logger.Fatalln("Failed to establish outbound connection to [", peerAddr, "]: ", err)
	}

	defer conn.Close()

	encoded, err := bindTransmitter.Encode()

	if err != nil {
		logger.Fatalln("Failed to encode bind_transmitter message: ", err)
	}

	logger.Printf("Encoded bind_transmitter message.  Encoded size = (%d)\n", len(encoded))

	_, err = conn.Write(encoded)

	if err != nil {
		logger.Fatalln("Failed to write bind_transmitter message: ", err)
	}

	readbuf := make([]byte, 4096)

	bytesRead, err := conn.Read(readbuf)

	if err != nil {
		logger.Fatalln("Failed to read from peer: ", err)
	}

	logger.Printf("Received (%d) octets from peer", bytesRead)

	_, err = smpp.DecodePDU(readbuf)

	if err != nil {
		logger.Fatalln("Failed to decode incoming PDU: ", err)
	}

	logger.Println("Closing peer connection")

	logger.Println("Done")
}
