package smpp

import (
	"testing"
)

func TestSimplePDUs(t *testing.T) {
	runPDUTest(t, "PDU 01", CommandDataSm, 0, 0x0f, []*Parameter{}, []*Parameter{}, 16, []byte{
		0x00, 0x00, 0x00, 0x10,
		0x00, 0x00, 0x01, 0x03,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x0f,
	})
}

func TestPDUsWithParameters(t *testing.T) {
	runPDUTest(t, "PDU with Mandatory Params Only", CommandDataSm, 0, 0x419, []*Parameter{
		NewASCIIParameter("WAP"),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(1)),
		NewASCIIParameter("10597"),
		NewFLParameter(uint8(1)),
		NewFLParameter(uint8(1)),
		NewASCIIParameter("+18809990011"),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(4)),
	}, []*Parameter{}, 45, []byte{
		0x0, 0x0, 0x0, 0x2d, // length
		0x0, 0x0, 0x1, 0x03, // command ID
		0x0, 0x0, 0x0, 0x0, // status
		0x0, 0x0, 0x4, 0x19, // sequence number
		0x57, 0x41, 0x50, 0x0,
		0x0,
		0x1,
		0x31, 0x30, 0x35, 0x39, 0x37, 0x0,
		0x1,
		0x1,
		0x2b, 0x31, 0x38, 0x38, 0x30, 0x39, 0x39, 0x39, 0x30, 0x30, 0x31, 0x31, 0x0,
		0x0,
		0x4,
	})

	runPDUTest(t, "PDU with Mandatory Params Only", CommandDataSm, 0, 0x419, []*Parameter{
		NewASCIIParameter("WAP"),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(1)),
		NewASCIIParameter("10597"),
		NewFLParameter(uint8(1)),
		NewFLParameter(uint8(1)),
		NewASCIIParameter("+18809990011"),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(4)),
	}, []*Parameter{
		NewTLVParameter(0x20a, uint16(0x23f0)),
		NewTLVParameter(0x20b, uint16(0x0b84)),
		NewTLVParameter(0x20c, uint16(0x0417)),
		NewTLVParameter(0x20e, uint8(2)),
		NewTLVParameter(0x20f, uint8(2)),
		NewTLVParameter(0x0017, uint32(0x000542e3)),
		NewTLVParameter(0x0019, uint8(0)),
		NewTLVParameter(0x0424, []byte{
			0x05, 0x42, 0xe3, 0x83, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x31, 0x30,
			0x37, 0x2e, 0x32, 0x33, 0x39, 0x2e, 0x31, 0x34, 0x2e, 0x31, 0x34, 0x3a, 0x38,
			0x30, 0x30, 0x33, 0x2f, 0x31, 0x30, 0x32, 0x30, 0x31, 0x36, 0x31, 0x30, 0x32,
			0x34, 0x35, 0x30, 0x30, 0x30, 0x31, 0x31, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30,
			0x30, 0x30, 0x30, 0x00,
		}),
		NewTLVParameter(0x0421, uint8(0)),
	}, 151, []byte{
		0x0, 0x0, 0x0, 0x97, // length
		0x0, 0x0, 0x1, 0x03, // command ID
		0x0, 0x0, 0x0, 0x0, // status
		0x0, 0x0, 0x4, 0x19, // sequence number
		0x57, 0x41, 0x50, 0x0,
		0x0,
		0x1,
		0x31, 0x30, 0x35, 0x39, 0x37, 0x0,
		0x1,
		0x1,
		0x2b, 0x31, 0x38, 0x38, 0x30, 0x39, 0x39, 0x39, 0x30, 0x30, 0x31, 0x31, 0x0,
		0x0,
		0x4,
		0x02, 0x0a, 0x00, 0x02, 0x23, 0xf0,
		0x02, 0x0b, 0x00, 0x02, 0x0b, 0x84,
		0x02, 0x0c, 0x00, 0x02, 0x04, 0x17,
		0x02, 0x0e, 0x00, 0x01, 0x02,
		0x02, 0x0f, 0x00, 0x01, 0x02,
		0x00, 0x17, 0x00, 0x04, 0x00, 0x05, 0x42, 0xe3,
		0x00, 0x19, 0x00, 0x01, 0x0,
		0x04, 0x24, 0x00, 0x38, 0x05, 0x42, 0xe3, 0x83, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x31, 0x30,
		0x37, 0x2e, 0x32, 0x33, 0x39, 0x2e, 0x31, 0x34, 0x2e, 0x31, 0x34, 0x3a, 0x38, 0x30, 0x30, 0x33, 0x2f,
		0x31, 0x30, 0x32, 0x30, 0x31, 0x36, 0x31, 0x30, 0x32, 0x34, 0x35, 0x30, 0x30, 0x30, 0x31, 0x31, 0x30,
		0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x00,
		0x04, 0x21, 0x00, 0x01, 0x0,
	})

}

func runPDUTest(t *testing.T, testname string, commandID CommandIDType, status uint32, sequence uint32, mandatoryParams []*Parameter, optionalParams []*Parameter, expectedLength uint32, expectedEncoding []byte) {
	pdu := NewPDU(commandID, status, sequence, mandatoryParams, optionalParams)

	if pdu == nil {
		t.Errorf("Test %s: NewPDU returns nil", testname)
	} else {
		if pdu.CommandID != commandID {
			t.Errorf("Test %s: CommandID does not match expected value, expected %d, got %d", testname, commandID, pdu.CommandID)
		}

		if pdu.CommandLength != expectedLength {
			t.Errorf("Test %s: CommandLength does not match expected value, expected %d, got %d", testname, expectedLength, pdu.CommandLength)
		}

		if pdu.CommandStatus != status {
			t.Errorf("Test %s: CommandStatus does not match expected value, expected %d, got %d", testname, status, pdu.CommandStatus)
		}

		if pdu.SequenceNumber != sequence {
			t.Errorf("Test %s: SequenceNumber does not match expected value, expected 0x%08x, got 0x%08x", testname, sequence, pdu.SequenceNumber)
		}

		encoded, err := pdu.Encode()

		if err != nil {
			t.Errorf("Test %s: error from pdu.Encode(), error = %s", testname, err)
		} else if pdu.CommandLength != uint32(len(expectedEncoding)) {
			t.Errorf("Test %s: Provided byte stream length does not match CommandLength", testname)
		} else {
			for i := 0; i < len(expectedEncoding); i++ {
				if encoded[i] != expectedEncoding[i] {
					t.Errorf("Test %s: encoding byte (%d) does not match expected, expected (0x%02x), got (0x%02x)", testname, i, expectedEncoding[i], encoded[i])
				}
			}
		}
	}
}
