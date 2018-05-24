package smpp

import (
	"bytes"
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
		NewCOctetStringParameter("WAP"),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(1)),
		NewCOctetStringParameter("10597"),
		NewFLParameter(uint8(1)),
		NewFLParameter(uint8(1)),
		NewCOctetStringParameter("+18809990011"),
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

	runPDUTest(t, "PDU with Mandatory and Optional Params", CommandDataSm, 0, 0x419, []*Parameter{
		NewCOctetStringParameter("WAP"),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(1)),
		NewCOctetStringParameter("10597"),
		NewFLParameter(uint8(1)),
		NewFLParameter(uint8(1)),
		NewCOctetStringParameter("+18809990011"),
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

func TestCommandDataSmDecode(t *testing.T) {
	_, err := DecodePDU([]byte{
		0x0, 0x0, 0x0, 0x10,
		0x0, 0x0, 0x1, 0x03,
		0x0, 0x0, 0x0, 0x00,
		0x0, 0x0, 0x4, 0x19,
	})

	if err == nil {
		t.Errorf("No Error on DecodePDU when its just Header for CommandDataSm")
	}

	encoded := []byte{
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
	}

	testPDUDecode(t, "CommandDataSm-1", encoded, 0x2d, 0x103, 0, 0x419, 9, 0)
}

func TestCommandBindTransmitterRespPDU(t *testing.T) {
	testname := "Command bind-transmitter-resp-1"

	encoded := []byte{
		0x0, 0x0, 0x0, 0x17,
		0x80, 0x0, 0x0, 0x02,
		0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x1,
		0x73, 0x6d, 0x73, 0x63, 0x30, 0x31, 0x0,
	}

	bindTransmitterResp := NewPDU(CommandBindTransmitterResp, 0, 1, []*Parameter{
		NewCOctetStringParameter("smsc01"), // system_id
	}, []*Parameter{})

	btEncode, err := bindTransmitterResp.Encode()

	if err != nil {
		t.Error("Test ", testname, ": failed to encode PDU: ", err)
	} else {
		if len(encoded) != len(btEncode) {
			t.Errorf("Test %s: length of encoded [%d] != length of Encode() output [%d]", testname, len(encoded), len(btEncode))
		} else {
			for i := 0; i < len(encoded); i++ {
				if encoded[i] != btEncode[i] {
					t.Errorf("Test %s: encoded and Encode() output begin differing at byte (%d), expect (%02x), got (%02x)", testname, i, encoded[i], btEncode[i])
					break
				}
			}
		}
	}

	testPDUDecode(t, testname, encoded, 23, 0x80000002, 0x0, 0x01, 1, 0)
}

func compareByteArrays(t *testing.T, testname string, expected []byte, got []byte) {
	if len(expected) != len(got) {
		t.Errorf("%s: length of expected encoding [%d] != length of Encode() output [%d]", testname, len(expected), len(got))
	} else {
		for i := 0; i < len(expected); i++ {
			if expected[i] != got[i] {
				t.Errorf("%s: expected encoding and Encode() output begin differing at element (%d), expect (%02x), got (%02x)", testname, i, expected[i], got[i])
				return
			}
		}
	}

}

func TestCommandBindTransmitterPDU(t *testing.T) {
	encoded := []byte{
		0x0, 0x0, 0x0, 0x2c, //length
		0x0, 0x0, 0x0, 0x02, // command ID
		0x0, 0x0, 0x0, 0x0, // status
		0x0, 0x0, 0x0, 0x1, // sequence number
		0x65, 0x73, 0x6d, 0x65, 0x30, 0x31, 0x00, // system_id = esme01
		0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x00, // password = password
		0x67, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x0, // system_type = generic
		0x34, // interface_version
		0x0,  // addr_ton
		0x0,  // addr_npi
		0x0,  // address_range = ""
	}

	bindTransmitter := NewPDU(CommandBindTransmitter, 0, 1, []*Parameter{
		NewCOctetStringParameter("esme01"),   // system_id
		NewCOctetStringParameter("password"), // password
		NewCOctetStringParameter("generic"),  // system_type
		NewFLParameter(uint8(0x34)),          // interface_version
		NewFLParameter(uint8(0x0)),           // addr_ton
		NewFLParameter(uint8(0x0)),           // addr_npi
		NewCOctetStringParameter(""),         // address_range
	}, []*Parameter{})

	btEncode, err := bindTransmitter.Encode()

	if err != nil {
		t.Error("Command bind-transmitter-1: failed to encode PDU: ", err)
	} else {
		compareByteArrays(t, "Command bind-transmitter-1", encoded, btEncode)
	}

	testPDUDecode(t, "Command bind-transmitter-1", encoded, 44, 0x02, 0x0, 0x01, 7, 0)
}

func TestCommandSubmitSmPDU(t *testing.T) {
	testname := "Command submit-sm-1"

	encoded := []byte{
		0x00, 0x00, 0x00, 0xe1, // length
		0x00, 0x00, 0x00, 0x04, // command ID
		0x00, 0x00, 0x00, 0x00, // status
		0x00, 0x00, 0x00, 0x5e, // sequence number
		0x00,                                                 // service_type
		0x00,                                                 // source_addr_ton
		0x01,                                                 // source_addr_npi
		0x32, 0x38, 0x38, 0x30, 0x39, 0x30, 0x39, 0x30, 0x00, // source_addr
		0x01,                                                                   // dest_addr_ton
		0x01,                                                                   // dest_addr_npi
		0x31, 0x33, 0x31, 0x33, 0x39, 0x35, 0x39, 0x31, 0x34, 0x36, 0x33, 0x00, // destination_addr
		0x00,                                                                                                 // esm_class
		0x00,                                                                                                 // protocol_id
		0x00,                                                                                                 // priority_level
		0x00,                                                                                                 // schedule_delivery_time
		0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x35, 0x30, 0x30, 0x30, 0x30, 0x30, 0x52, 0x00, // validity_time
		0x00, // registered_delivery
		0x00, // replace_if_present_flag
		0xf0, // data_coding
		0x00, // sm_default_msg_id
		0x8d, // sm_length
		// short_message...
		0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, 0x20, 0x73, 0x68, 0x6f, 0x72, 0x74,
		0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x74, 0x68, 0x6f, 0x75, 0x67, 0x68, 0x20, 0x69, 0x74, 0x20,
		0x69, 0x73, 0x20, 0x73, 0x6f, 0x6d, 0x65, 0x77, 0x68, 0x61, 0x74, 0x20, 0x6c, 0x6f, 0x6e, 0x67, 0x65, 0x72, 0x20, 0x74,
		0x68, 0x61, 0x6e, 0x20, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2c, 0x20, 0x62, 0x65, 0x69, 0x6e, 0x67, 0x20, 0x3e, 0x20, 0x35,
		0x30, 0x20, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x21, 0x20, 0x44, 0x6f, 0x6e, 0x27, 0x74, 0x20,
		0x67, 0x65, 0x74, 0x20, 0x65, 0x78, 0x63, 0x69, 0x74, 0x65, 0x64, 0x20, 0x3a, 0x40, 0x20, 0x3a, 0x23, 0x20, 0x3a, 0x24,
		0x20, 0x3a, 0x25, 0x20, 0x3a, 0x5e, 0x29, 0x20, 0x65, 0x6d, 0x6f, 0x6a, 0x69, 0x20, 0x6c, 0x69, 0x6b, 0x65, 0x2e, 0x2e,
		0x2e,
		0x02, 0x0c, 0x00, 0x02, 0x00, 0x05, // opt param: sar_msg_ref_num
		0x02, 0x0e, 0x00, 0x01, 0x02, // opt param: sar_total_segments
		0x02, 0x0f, 0x00, 0x01, 0x01, // opt param: sar_seg_seqnum1
	}

	submitSmPDU := NewPDU(CommandSubmitSm, 0, 0x5e, []*Parameter{
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(1)),
		NewCOctetStringParameter("28809090"),
		NewFLParameter(uint8(1)),
		NewFLParameter(uint8(1)),
		NewCOctetStringParameter("13139591463"),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(0)),
		NewCOctetStringParameter("000000000500000R"),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(0xf0)),
		NewFLParameter(uint8(0)),
		NewFLParameter(uint8(0x8d)),
		NewOctetStringFromString("This is a test short message, though it is somewhat longer than short, being > 50 characters! Don't get excited :@ :# :$ :% :^) emoji like..."),
	}, []*Parameter{
		NewTLVParameter(0x020c, uint16(5)),
		NewTLVParameter(0x020e, uint8(2)),
		NewTLVParameter(0x020f, uint8(1)),
	})

	pduEncode, err := submitSmPDU.Encode()

	if err != nil {
		t.Error(testname, ": error on Encode of submitSmPDU: ", err)
	} else {
		compareByteArrays(t, testname, encoded, pduEncode)
	}

	testPDUDecode(t, testname, pduEncode, 0xe1, CommandSubmitSm, 0, 0x5e, 18, 3)
}

func testPDUDecode(t *testing.T, testname string, encodedPDU []byte, expectedCommandLength uint32, expectedCommandID CommandIDType, expectedStatus uint32, expectedSequence uint32, expectedMParamCount uint32, expectedOParamCount uint32) {
	pdu, err := DecodePDU(encodedPDU)

	if err != nil {
		t.Errorf("Test %s: failed to decode, error = %s", testname, err)
	}

	if pdu.CommandLength != expectedCommandLength {
		t.Errorf("Test %s: commandLength should be (%d), is (%d)", testname, expectedCommandLength, pdu.CommandLength)
	}

	if pdu.CommandID != expectedCommandID {
		t.Errorf("Test %s: CommandID should be (%08x), is (%08x)", testname, expectedCommandID, pdu.CommandID)
	}

	if pdu.CommandStatus != expectedStatus {
		t.Errorf("Test %s: CommandStatus should be (%d), is (%d)", testname, expectedStatus, pdu.CommandStatus)
	}

	if pdu.SequenceNumber != expectedSequence {
		t.Errorf("Test %s: SequenceNumber should be (%d), is (%d)", testname, expectedSequence, pdu.SequenceNumber)
	}

	if len(pdu.MandatoryParameters) != int(expectedMParamCount) {
		t.Errorf("Test %s: count of mandatory parameters should be (%d), is (%d)", testname, expectedMParamCount, len(pdu.MandatoryParameters))
	}

	if len(pdu.OptionalParameters) != int(expectedOParamCount) {
		t.Errorf("Test %s: count of optional parameters should be (%d), is (%d)", testname, expectedOParamCount, len(pdu.OptionalParameters))
	}

	reEncoded, err := pdu.Encode()

	if err != nil {
		t.Errorf("Test %s: error on pdu.Encode() = %s", testname, err)
	} else {
		if !bytes.Equal(reEncoded, encodedPDU) {
			t.Errorf("Test %s: pdu.Encode() does not match original encoding", testname)
		}
	}
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
