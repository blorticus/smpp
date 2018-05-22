package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"smpp"
)

func main() {
	logger := log.New(os.Stderr, filepath.Base(os.Args[0])+": ", 0)

	logger.Println("Creating PDU object")

	pdu := smpp.NewPDU(smpp.CommandDataSm, 0, 0x419, []*smpp.Parameter{
		smpp.NewASCIIParameter("WAP"),
		smpp.NewFLParameter(uint8(0)),
		smpp.NewFLParameter(uint8(1)),
		smpp.NewASCIIParameter("10597"),
		smpp.NewFLParameter(uint8(1)),
		smpp.NewFLParameter(uint8(1)),
		smpp.NewASCIIParameter("+18809990011"),
		smpp.NewFLParameter(uint8(0)),
		smpp.NewFLParameter(uint8(4)),
	}, []*smpp.Parameter{})

	encoded, _ := pdu.Encode()

	for i := 0; i < len(encoded); i++ {
		fmt.Printf("%02x ", encoded[i])
	}

	fmt.Println("")
}
