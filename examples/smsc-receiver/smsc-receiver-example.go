package main

import (
    "flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
    "path"
	"path/filepath"
	"smpp"
    "strconv"
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
	octets, err := conn.Read(*buf)

	if err != nil {
		if err == io.EOF {
			return nil
		}

		logger.Fatalln("Failed on read: ", err)
	}

	logger.Printf("Read (%d) octets", octets)

	pdu, err := smpp.DecodePDU(*buf)

	if err != nil {
		logger.Fatalln("Failed to encode incoming PDU: ", err)
	}

	return pdu
}

func main() {
    var bindIP string
    var bindPort int

    flag.StringVar(&bindIP, "addr", "127.0.0.1", "Listener IP")
    flag.IntVar(&bindPort, "port", 2775, "Listener TCP port")

    os.Args[0] = path.Base(os.Args[0])

    flag.Parse()

	bindAddr := bindIP + ":" + strconv.Itoa(bindPort)

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

	logger.Println("Waiting for bind-transmitter from peer")

	receivedPDU := recvPDU(conn, &readbuf, logger)

	if receivedPDU.CommandID != smpp.CommandBindTransmitter {
		logger.Fatalf("Expected bind_transmitter from peer, but received (%08x)\n", receivedPDU.CommandID)
	}

	logger.Println("bind-transmitter received; sending bind-transmitter-resp")

	seqNum := receivedPDU.SequenceNumber

	responsePDU := smpp.NewPDU(smpp.CommandBindTransmitterResp, 0, seqNum, []*smpp.Parameter{
		smpp.NewCOctetStringParameter("smsc01"), // system_id
	}, []*smpp.Parameter{})

	sendPDU(conn, responsePDU, logger)

	logger.Println("Responded with transmitter_resp")

	logger.Println("Entering listen loop")

	msgID := uint32(1)
	for {
		pdu := recvPDU(conn, &readbuf, logger)

		if pdu == nil {
			logger.Println("Peer closed connection")
			break
		}

		logger.Printf("Received PDU from peer, type is (%s)", pdu.CommandName())

		switch pdu.CommandID {
		case smpp.CommandSubmitSm:
			responsePDU = smpp.NewPDU(smpp.CommandSubmitSmResp, 0, pdu.SequenceNumber, []*smpp.Parameter{
				smpp.NewCOctetStringParameter(fmt.Sprintf("msg-%d", msgID)),
			}, []*smpp.Parameter{})
			msgID++

			sendPDU(conn, responsePDU, logger)
		default:
			logger.Println("No response for this message type")
		}
	}

	logger.Println("Done")
}
