package smpp

import (
	"encoding/binary"
	"net"
)

// NetworkStreamReader provides a mechanism for reading PDUs from an incoming TCP stream connection, breaking
// the stream into PDUs
type NetworkStreamReader struct {
	connectionFromWhichToRead net.Conn
	readBuffer                []byte
	pduBuffer                 []byte
}

// NewNetworkStreamReader creates a NetworkStreamReader that operates on the identified connection
func NewNetworkStreamReader(fromConnection net.Conn) *NetworkStreamReader {
	return &NetworkStreamReader{connectionFromWhichToRead: fromConnection, readBuffer: make([]byte, 65536), pduBuffer: make([]byte, 0, 65536)}
}

// Read performs a read of the associated TCP stream and attempts to extract one or more PDUs from the
// stream.  If there are data left over after extracting zero or more PDUs, those data are saved, and
// subsequent Read values are appended to those data
func (reader *NetworkStreamReader) Read() ([]*PDU, error) {
	bytesRead, err := reader.connectionFromWhichToRead.Read(reader.readBuffer)

	if err != nil {
		return nil, err
	}

	reader.pduBuffer = append(reader.pduBuffer, reader.readBuffer[:bytesRead]...)

	extractedPDUs := make([]*PDU, 0, 3)

	for len(reader.pduBuffer) >= 16 {
		pduLength := uint32(binary.BigEndian.Uint32(reader.pduBuffer[0:4]))

		if len(reader.pduBuffer) >= int(pduLength) {
			pdu, err := DecodePDU(reader.pduBuffer[:pduLength])

			copy(reader.pduBuffer[0:len(reader.pduBuffer)-int(pduLength)], reader.pduBuffer[pduLength:])
			reader.pduBuffer = reader.pduBuffer[:len(reader.pduBuffer)-int(pduLength)]

			if err != nil {
				return extractedPDUs, err
			}

			extractedPDUs = append(extractedPDUs, pdu)
		} else {
			return extractedPDUs, nil
		}
	}

	return extractedPDUs, nil
}

// ExtractNextPDUs repeatedly reads from the TCP stream until there is at least one PDU.
// It returns the set of extracted PDUs, and like Read(), stores any remaining data for
// subsequent calls
func (reader *NetworkStreamReader) ExtractNextPDUs() ([]*PDU, error) {
	for {
		pdus, err := reader.Read()

		if err != nil {
			return nil, err
		}

		if len(pdus) > 0 {
			return pdus, nil
		}
	}
}
