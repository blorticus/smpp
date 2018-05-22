package smpp

import (
	"encoding/binary"
)

// ParameterType is an enumeration of parameter types
type ParameterType uint32

const (
	// TypeUint8 encoded as single byte, unsigned 8-bit integer
	TypeUint8 ParameterType = iota
	// TypeUint16 encoded as two bytes, unsigned 16-bit integer
	TypeUint16
	// TypeUint32 encoded as four bytes, unsigned 32-bit integer
	TypeUint32
	// TypeASCII encoded as variable length field of bytes, must be NULL (0) terminated
	TypeASCII
	// TypeOctetString encoded as variable length field of bytes
	TypeOctetString
	// TypeTLV encodes as octet string
	TypeTLV
)

// TLV represents the value in a Parameter struct for TLV type Parameters
type TLV struct {
	Tag     uint16
	VLength uint16
	Value   interface{}
}

// Parameter is mandatory or optional parameter.  If the type is TypeTLV then
// Value is an instance of TLV
type Parameter struct {
	Type         ParameterType
	EncodeLength uint32
	Value        interface{}
}

// NewFLParameter creates a Parameter where the length is fixed by the type (e.g., TypeUint32).
// The stored Type will be introspected from the type of 'value'
func NewFLParameter(value interface{}) *Parameter {
	param := new(Parameter)

	param.Value = value

	switch value.(type) {
	case uint8:
		param.Type = TypeUint8
		param.EncodeLength = 1

	case uint16:
		param.Type = TypeUint16
		param.EncodeLength = 2

	case uint32:
		param.Type = TypeUint32
		param.EncodeLength = 4

	default:
		return nil
	}

	return param
}

// NewASCIIParameter creates a Parameter where the value is a C-Octet String
// (a null terminated string).  It is up to the caller to pass a correct value
// if Decimal or Hex is required.  The passed string must contain only ASCII,
// or this won't work the way you expect.
func NewASCIIParameter(value string) *Parameter {
	return &Parameter{TypeASCII, uint32(len(value)) + 1, value}
}

// NewTLVParameter creates a new Parameter with the provided tag.  The length
// is introspected from the value, which may be a uint8, a uint16, a uint32,
// a string, or a []byte.
func NewTLVParameter(tag uint16, value interface{}) *Parameter {
	switch value.(type) {
	case uint8:
		return &Parameter{TypeTLV, 5, TLV{tag, 1, value}}

	case uint16:
		return &Parameter{TypeTLV, 6, TLV{tag, 2, value}}

	case uint32:
		return &Parameter{TypeTLV, 8, TLV{tag, 4, value}}

	case string:
		return &Parameter{TypeTLV, 4 + uint32(len(value.(string))), TLV{tag, uint16(len(value.(string))), value}}

	case []byte:
		return &Parameter{TypeTLV, 4 + uint32(len(value.([]byte))), TLV{tag, uint16(len(value.([]byte))), value}}
	}

	return nil
}

// Encode converts the 'param' object into a byte stream appropriate for
// network transmission
func (param Parameter) Encode() []byte {
	encoded := make([]byte, param.EncodeLength)

	switch param.Type {
	case TypeUint8:
		encoded[0] = param.Value.(uint8)

	case TypeUint16:
		binary.BigEndian.PutUint16(encoded[0:2], param.Value.(uint16))

	case TypeUint32:
		binary.BigEndian.PutUint32(encoded[0:4], param.Value.(uint32))

	case TypeASCII:
		b := make([]byte, param.EncodeLength)
		copy(b[0:param.EncodeLength], param.Value.(string))
		return b

	case TypeTLV:
		binary.BigEndian.PutUint16(encoded[0:2], param.Value.(TLV).Tag)
		binary.BigEndian.PutUint16(encoded[2:4], param.Value.(TLV).VLength)

		v := param.Value.(TLV).Value
		switch v.(type) {
		case uint8:
			encoded[4] = v.(uint8)

		case uint16:
			binary.BigEndian.PutUint16(encoded[4:6], v.(uint16))

		case uint32:
			binary.BigEndian.PutUint32(encoded[4:8], v.(uint32))

		case string:
			copy(encoded[4:4+len(v.(string))], v.(string))

		case []byte:
			copy(encoded[4:4+len(v.([]byte))], v.([]byte))
		}
	}

	return encoded
}

// CommandIDType is an enumeration of defined  Command IDs
type CommandIDType uint32

const (
	_                          CommandIDType = iota
	CommandGenericNack                       = 0x80000000
	CommandBindReceiver                      = 0x00000001
	CommandBindReceiverResp                  = 0x80000001
	CommandBindTransmitter                   = 0x00000002
	CommandBindTransmitterResp               = 0x80000002
	CommandQuerySm                           = 0x00000003
	CommandQuerySmResp                       = 0x80000003
	CommandSubmitSm                          = 0x00000004
	CommandSubmitSmResp                      = 0x80000004
	CommandDeliverSm                         = 0x00000005
	CommandDeliverSmResp                     = 0x80000005
	CommandUnbind                            = 0x00000006
	CommandUnbindResp                        = 0x80000006
	CommandReplaceSm                         = 0x00000007
	CommandReplaceSmResp                     = 0x80000007
	CommandCancelSm                          = 0x00000008
	CommandCancelSmResp                      = 0x80000008
	CommandBindTransceiver                   = 0x00000009
	CommandBindTransceiverResp               = 0x80000009
	CommandOutbind                           = 0x0000000B
	CommandEnquireLink                       = 0x00000015
	CommandEnquireLinkResp                   = 0x80000015
	CommandSubmitMulti                       = 0x00000021
	CommandSubmitMultiResp                   = 0x80000021
	CommandAlertNotification                 = 0x00000102
	CommandDataSm                            = 0x00000103
	CommandDataSmResp                        = 0x80000103
)

// PDU is a PDU for
type PDU struct {
	CommandLength       uint32
	CommandID           CommandIDType
	CommandStatus       uint32
	SequenceNumber      uint32
	MandatoryParameters []*Parameter
	OptionalParameters  []*Parameter
}

// NewPDU creates a new PDU object
func NewPDU(id CommandIDType, status uint32, sequence uint32, mandatoryParams []*Parameter, optionalParams []*Parameter) *PDU {
	pdu := PDU{0, id, status, sequence, mandatoryParams, optionalParams}

	length := uint32(16)

	for i := 0; i < len(mandatoryParams); i++ {
		length += mandatoryParams[i].EncodeLength
	}

	for i := 0; i < len(optionalParams); i++ {
		length += optionalParams[i].EncodeLength
	}

	pdu.CommandLength = length

	return &pdu
}

// ComputeLength computes the enocde length of the PDU and
// returns it
func (pdu PDU) ComputeLength() uint32 {
	length := uint32(16) // header

	for _, mparam := range pdu.MandatoryParameters {
		length += mparam.EncodeLength
	}

	for _, oparam := range pdu.OptionalParameters {
		length += oparam.EncodeLength
	}

	return length
}

// Encode converts the 'pdu' object into a byte stream appropriate for network transmission
func (pdu PDU) Encode() ([]byte, error) {
	if pdu.CommandLength < 1 {
		return []byte{}, nil
	}

	pdu.CommandLength = pdu.ComputeLength()

	encoded := make([]byte, pdu.CommandLength)

	binary.BigEndian.PutUint32(encoded[0:4], pdu.CommandLength)
	binary.BigEndian.PutUint32(encoded[4:8], uint32(pdu.CommandID))
	binary.BigEndian.PutUint32(encoded[8:12], pdu.CommandStatus)
	binary.BigEndian.PutUint32(encoded[12:16], pdu.SequenceNumber)

	s := uint32(16)
	for _, mparam := range pdu.MandatoryParameters {
		e := mparam.EncodeLength + s
		copy(encoded[s:e], mparam.Encode())
		s = e
	}

	for _, oparam := range pdu.OptionalParameters {
		e := oparam.EncodeLength + s
		copy(encoded[s:e], oparam.Encode())
		s = e
	}

	return encoded, nil
}
