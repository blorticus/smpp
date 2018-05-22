package smpp

import (
	"testing"
)

// TestFixedLength tests NewFLParameter variants
func TestFixedLength(t *testing.T) {
	runFLTest(t, "NewFLParameter(uint8(0))", uint8(0), TypeUint8, 1, uint8(0), []byte{0x0})
	runFLTest(t, "NewFLParameter(uint8(255))", uint8(255), TypeUint8, 1, uint8(255), []byte{0xff})

	runFLTest(t, "NewFLParameter(uint16(0))", uint16(0), TypeUint16, 2, uint16(0), []byte{0x0, 0x0})
	runFLTest(t, "NewFLParameter(uint16(65535))", uint16(65535), TypeUint16, 2, uint16(65535), []byte{0xff, 0xff})

	runFLTest(t, "NewFLParameter(uint32(0))", uint32(0), TypeUint32, 4, uint32(0), []byte{0x0, 0x0, 0x0, 0x0})
	runFLTest(t, "NewFLParameter(uint32(65535))", uint32(4294967295), TypeUint32, 4, uint32(4294967295), []byte{0xff, 0xff, 0xff, 0xff})
}

func TestASCII(t *testing.T) {
	runASCIITest(t, "NewASCIIParameter('')", "", 1, []byte{0x00})
	runASCIITest(t, "NewASCIIParameter('this')", "this", 5, []byte{0x74, 0x68, 0x69, 0x73, 0x00})
}

func runASCIITest(t *testing.T, testname string, value string, expectedEncodeLength uint32, expectedEncoding []byte) {
	param := NewASCIIParameter(value)

	if param == nil {
		t.Errorf("Test %s:, received nil object", testname)
	} else {
		if param.Type != TypeASCII {
			t.Errorf("Test %s: expected TypeASCII, but got type (%d)", testname, param.Type)
		}

		if param.EncodeLength != uint32(len(value)+1) {
			t.Errorf("Test %s: EncodedLength incorrect, expected (%d), got (%d)", testname, expectedEncodeLength, param.EncodeLength)
		}

		encoded := param.Encode()

		if uint32(len(encoded)) != param.EncodeLength {
			t.Errorf("Test %s: Provided encoded byte stream is length (%d), but EncodeLength (%d)", testname, len(encoded), param.EncodeLength)
		} else {
			for i := 0; i < len(encoded); i++ {
				if encoded[i] != expectedEncoding[i] {
					t.Errorf("Test %s: Encoded byte stream at byte (%d), expected (0x%02x), got (0x%02x)", testname, i, expectedEncoding[i], encoded[i])
				}
			}
		}
	}
}

func runFLTest(t *testing.T, testname string, value interface{}, expectedType ParameterType, expectedEncodeLength uint32, expectedValue interface{}, encodedBytes []byte) {
	param := NewFLParameter(value)

	if param == nil {
		t.Errorf("Test %s:, received nil object", testname)
	} else {
		if param.Type != expectedType {
			t.Errorf("Test %s: expected Type == %d, got %d", testname, expectedType, param.Type)
		}

		if param.EncodeLength != expectedEncodeLength {
			t.Errorf("Test %s: expected EncodeLength == %d, got %d", testname, expectedEncodeLength, param.EncodeLength)
		}

		if param.Value != expectedValue {
			t.Errorf("Test %s: value differs from expected", testname)
		}

		if param.EncodeLength != uint32(len(encodedBytes)) {
			t.Errorf("Test %s: provided encoding length is %d but object encode length is %d", testname, len(encodedBytes), param.EncodeLength)
		}

		e := param.Encode()
		for i := 0; i < len(encodedBytes); i++ {
			if e[i] != encodedBytes[i] {
				t.Errorf("Test %s: encoded byte %d does not match expected byte, got %02x, expected %02x", testname, i, e[i], encodedBytes[i])
			}
		}
	}
}
