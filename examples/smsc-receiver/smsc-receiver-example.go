package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"smpp"
)

func main() {
	bindIP := "127.0.0.1"
	bindPort := "2775"
	bindAddr := bindIP + ":" + bindPort

	logger := log.New(os.Stderr, filepath.Base(os.Args[0])+": ", 0)

	logger.Println("Starting listener")

	listener, err := net.Listen("tcp", bindAddr)

	if err != nil {
		logger.Fatalln("Failed to start listener: ", err)
	}

	defer listener.Close()

	conn, err := listener.Accept()

	if err != nil {
		logger.Fatalln("Failed to Accept incoming socket: ", err)
	}

	defer conn.Close()

	logger.Println("Accepted incoming connection from: ", conn.RemoteAddr().String())

	readbuf := make([]byte, 4096)

	bytesRead, err := conn.Read(readbuf)

	if err != nil {
		logger.Fatalln("Failed to Read from incoming socket: ", err)
	}

	logger.Printf("Received (%d) bytes from peer", bytesRead)

	receivedPDU, err := smpp.DecodePDU(readbuf)

	if err != nil {
		logger.Fatalln("Failed to decode incoming PDU: ", err)
	}

	if receivedPDU.CommandID != smpp.CommandBindTransmitter {
		logger.Fatalf("Expected bind_transmitter from peer, but received (%08x)\n", smpp.CommandBindTransmitter)
	}

	seqNum := receivedPDU.SequenceNumber

	responsePDU := smpp.NewPDU(smpp.CommandBindTransmitterResp, 0, seqNum, []*smpp.Parameter{
		smpp.NewCOctetStringParameter("smsc01"), // system_id
	}, []*smpp.Parameter{})

	encoded, err := responsePDU.Encode()

	if err != nil {
		logger.Fatalln("Failed to encode transmitter_resp PDU: ", err)
	}

	_, err = conn.Write(encoded)

	if err != nil {
		logger.Fatalln("Failed to write encoded transmitter_resp PDU: ", err)
	}

	logger.Println("Responded with transmitter_resp")
	logger.Println("Done")
}
