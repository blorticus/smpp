package smpp

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
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
	// TypeCOctetString encoded as variable length field of bytes, must be NULL (0) terminated
	TypeCOctetString
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

// ParameterDefinition provides attributes for a Parameter.  MaxLength will
// be set to the fixed length for fixed length types (e.g., TypeUint32).  If there is no
// MaxLength for a variable sized type (e.g., TypeASCII), then MaxLength will be set to zero.
// TagID is set only if the type is TypeTLV; otherwise it is zero.
type ParameterDefinition struct {
	Name      string
	Type      ParameterType
	MaxLength uint16
	TagID     uint16
}

var parameterTypeDefinition = map[string]ParameterDefinition{
	// Mandatory Parameter Set
	"addr_ton":                ParameterDefinition{"addr_ton", TypeUint8, 1, 0},
	"addr_npi":                ParameterDefinition{"addr_npi", TypeUint8, 1, 0},
	"address_range":           ParameterDefinition{"address_range", TypeCOctetString, 41, 0},
	"data_coding":             ParameterDefinition{"data_coding", TypeUint8, 1, 0},
	"destination_addr":        ParameterDefinition{"destination_addr", TypeCOctetString, 21, 0},
	"destination_addr_npi":    ParameterDefinition{"source_addr_npi", TypeUint8, 1, 0},
	"destination_addr_ton":    ParameterDefinition{"source_addr_ton", TypeUint8, 1, 0},
	"esm_class":               ParameterDefinition{"esm_class", TypeUint8, 1, 0},
	"interface_version":       ParameterDefinition{"interface_version", TypeUint8, 1, 0},
	"password":                ParameterDefinition{"password", TypeCOctetString, 9, 0},
	"message_id":              ParameterDefinition{"message_id", TypeCOctetString, 9, 0},
	"priority_flag":           ParameterDefinition{"priority_flag", TypeUint8, 1, 0},
	"protocol_id":             ParameterDefinition{"protocol_id", TypeUint8, 1, 0},
	"registered_delivery":     ParameterDefinition{"registered_delivery", TypeUint8, 1, 0},
	"replace_if_present_flag": ParameterDefinition{"replace_if_present_flag", TypeUint8, 1, 0},
	"schedule_delivery_time":  ParameterDefinition{"schedule_delivery_time", TypeCOctetString, 21, 0},
	"service_type":            ParameterDefinition{"service_type", TypeCOctetString, 9, 0},
	"short_message":           ParameterDefinition{"short_message", TypeCOctetString, 254, 0},
	"sm_default_msg_id":       ParameterDefinition{"sm_default_msg_id", TypeUint8, 1, 0},
	"sm_length":               ParameterDefinition{"sm_length", TypeUint8, 1, 0},
	"source_addr_npi":         ParameterDefinition{"source_addr_npi", TypeUint8, 1, 0},
	"source_addr_ton":         ParameterDefinition{"source_addr_ton", TypeUint8, 1, 0},
	"source_addr":             ParameterDefinition{"source_addr", TypeCOctetString, 21, 0},
	"system_id":               ParameterDefinition{"system_id", TypeCOctetString, 16, 0},
	"system_type":             ParameterDefinition{"system_type", TypeCOctetString, 13, 0},
	"validity_period":         ParameterDefinition{"validity_period", TypeCOctetString, 21, 0},

	// Optional Parameter Set
	"SC_interface_version":        ParameterDefinition{"SC_interface_version", TypeTLV, 0, 0x0210},
	"additional_status_info_text": ParameterDefinition{"additional_status_info_text", TypeTLV, 0, 0x001D},
	"alert_on_message_delivery":   ParameterDefinition{"alert_on_message_delivery", TypeTLV, 0, 0x130C},
	"callback_num":                ParameterDefinition{"callback_num", TypeTLV, 0, 0x0381},
	"callback_num_atag":           ParameterDefinition{"callback_num_atag", TypeTLV, 0, 0x0303},
	"callback_num_pres_ind":       ParameterDefinition{"callback_num_pres_ind", TypeTLV, 0, 0x0302},
	"delivery_failure_reason":     ParameterDefinition{"delivery_failure_reason", TypeTLV, 0, 0x0425},
	"dest_addr_subunit":           ParameterDefinition{"dest_addr_subunit", TypeTLV, 0, 0x0005},
	"dest_bearer_type":            ParameterDefinition{"dest_bearer_type", TypeTLV, 0, 0x0007},
	"dest_network_type":           ParameterDefinition{"dest_network_type", TypeTLV, 0, 0x0006},
	"dest_subaddress":             ParameterDefinition{"dest_subaddress", TypeTLV, 0, 0x0203},
	"dest_telematics_id":          ParameterDefinition{"dest_telematics_id", TypeTLV, 0, 0x0008},
	"destination_port":            ParameterDefinition{"destination_port", TypeTLV, 0, 0x020B},
	"display_time":                ParameterDefinition{"display_time", TypeTLV, 0, 0x1201},
	"dpf_result":                  ParameterDefinition{"dpf_result", TypeTLV, 0, 0x0420},
	"its_reply_type":              ParameterDefinition{"its_reply_type", TypeTLV, 0, 0x1380},
	"its_session_info":            ParameterDefinition{"its_session_info", TypeTLV, 0, 0x1383},
	"language_indicator":          ParameterDefinition{"language_indicator", TypeTLV, 0, 0x020D},
	"message_payload":             ParameterDefinition{"message_payload", TypeTLV, 0, 0x0424},
	"message_state":               ParameterDefinition{"message_state", TypeTLV, 0, 0x0427},
	"more_messages_to_send":       ParameterDefinition{"more_messages_to_send", TypeTLV, 0, 0x0426},
	"ms_availability_status":      ParameterDefinition{"ms_availability_status", TypeTLV, 0, 0x0422},
	"ms_msg_wait_facilities":      ParameterDefinition{"ms_msg_wait_facilities", TypeTLV, 0, 0x0030},
	"ms_validity":                 ParameterDefinition{"ms_validity", TypeTLV, 0, 0x1204},
	"network_error_code":          ParameterDefinition{"network_error_code", TypeTLV, 0, 0x0423},
	"number_of_messages":          ParameterDefinition{"number_of_messages", TypeTLV, 0, 0x0304},
	"payload_type":                ParameterDefinition{"payload_type", TypeTLV, 0, 0x0019},
	"privacy_indicator":           ParameterDefinition{"privacy_indicator", TypeTLV, 0, 0x0201},
	"qos_time_to_live":            ParameterDefinition{"qos_time_to_live", TypeTLV, 0, 0x0017},
	"receipted_message_id":        ParameterDefinition{"receipted_message_id", TypeTLV, 0, 0x001E},
	"sar_msg_ref_num":             ParameterDefinition{"sar_msg_ref_num", TypeTLV, 0, 0x020C},
	"sar_segment_seqnum":          ParameterDefinition{"sar_segment_seqnum", TypeTLV, 0, 0x020F},
	"sar_total_segments":          ParameterDefinition{"sar_total_segments", TypeTLV, 0, 0x020E},
	"set_dpf":                     ParameterDefinition{"set_dpf", TypeTLV, 0, 0x0421},
	"sms_signal":                  ParameterDefinition{"sms_signal", TypeTLV, 0, 0x1203},
	"source_addr_subunit":         ParameterDefinition{"source_addr_subunit", TypeTLV, 0, 0x000D},
	"source_bearer_type":          ParameterDefinition{"source_bearer_type", TypeTLV, 0, 0x000F},
	"source_network_type":         ParameterDefinition{"source_network_type", TypeTLV, 0, 0x000E},
	"source_port":                 ParameterDefinition{"source_port", TypeTLV, 0, 0x020A},
	"source_subaddress":           ParameterDefinition{"source_subaddress", TypeTLV, 0, 0x0202},
	"source_telematics_id":        ParameterDefinition{"source_telematics_id", TypeTLV, 0, 0x0010},
	"user_message_reference":      ParameterDefinition{"user_message_reference", TypeTLV, 0, 0x0204},
	"user_response_code":          ParameterDefinition{"user_response_code", TypeTLV, 0, 0x0205},
	"ussd_service_op":             ParameterDefinition{"ussd_service_op", TypeTLV, 0, 0x0501},
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

// NewCOctetStringParameter creates a Parameter where the value is a C-Octet String
// (a null terminated string).  It is up to the caller to pass a correct value
// if Decimal or Hex is required.  The passed string must contain only ASCII,
// or this won't work the way you expect.
func NewCOctetStringParameter(value string) *Parameter {
	return &Parameter{TypeCOctetString, uint32(len(value)) + 1, value}
}

// NewOctetStringFromString creates a Parameter of type OctetString (not null
// terminated) from a string.
func NewOctetStringFromString(value string) *Parameter {
	return &Parameter{TypeOctetString, uint32(len(value)), []byte(value)}
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
func (param *Parameter) Encode() []byte {
	encoded := make([]byte, param.EncodeLength)

	switch param.Type {
	case TypeUint8:
		encoded[0] = param.Value.(uint8)

	case TypeUint16:
		binary.BigEndian.PutUint16(encoded[0:2], param.Value.(uint16))

	case TypeUint32:
		binary.BigEndian.PutUint32(encoded[0:4], param.Value.(uint32))

	case TypeCOctetString:
		b := make([]byte, param.EncodeLength)
		copy(b[0:param.EncodeLength], param.Value.(string))
		return b

	case TypeOctetString:
		b := make([]byte, param.EncodeLength)
		copy(b[0:], param.Value.([]byte))
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

// These correspond to SMPP Message Type
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

var pduCommandName = map[CommandIDType]string{
	CommandGenericNack:         "generic-nack",
	CommandBindReceiver:        "bind-receiver",
	CommandBindReceiverResp:    "bind-receiver-resp",
	CommandBindTransmitter:     "bind-transmitter",
	CommandBindTransmitterResp: "bind-transmitter-resp",
	CommandQuerySm:             "query-sm",
	CommandQuerySmResp:         "query-sm-resp",
	CommandSubmitSm:            "submit-sm",
	CommandSubmitSmResp:        "submit-sm-resp",
	CommandDeliverSm:           "deliver-sm",
	CommandDeliverSmResp:       "deliver-sm-resp",
	CommandUnbind:              "unbind",
	CommandUnbindResp:          "unbind-resp",
	CommandReplaceSm:           "replace-sm",
	CommandReplaceSmResp:       "replace-sm-resp",
	CommandCancelSm:            "cancel-sm",
	CommandCancelSmResp:        "cancel-sm-resp",
	CommandBindTransceiver:     "bind-tranceiver",
	CommandBindTransceiverResp: "bind-tranceiver-resp",
	CommandOutbind:             "outbind",
	CommandEnquireLink:         "enquire-link",
	CommandEnquireLinkResp:     "enquire-link-resp",
	CommandSubmitMulti:         "submit-multi",
	CommandSubmitMultiResp:     "submit-multi-resp",
	CommandAlertNotification:   "alert-notification",
	CommandDataSm:              "data-sm",
	CommandDataSmResp:          "data-sm-resp",
}

var commandNameToCommandID = map[string]CommandIDType{
	"bind-receiver":         CommandBindReceiver,
	"bind-receiver-resp":    CommandBindReceiverResp,
	"bind-transmitter":      CommandBindTransmitter,
	"bind-transmitter-resp": CommandBindTransmitterResp,
	"query-sm":              CommandQuerySm,
	"query-sm-resp":         CommandQuerySmResp,
	"submit-sm":             CommandSubmitSm,
	"submit-sm-resp":        CommandSubmitSmResp,
	"deliver-sm":            CommandDeliverSm,
	"deliver-sm-resp":       CommandDeliverSmResp,
	"unbind":                CommandUnbind,
	"unbind-resp":           CommandUnbindResp,
	"replace-sm":            CommandReplaceSm,
	"replace-sm-resp":       CommandReplaceSmResp,
	"cancel-sm":             CommandCancelSm,
	"cancel-sm-resp":        CommandCancelSmResp,
	"bind-tranceiver":       CommandBindTransceiver,
	"bind-tranceiver-resp":  CommandBindTransceiverResp,
	"outbind":               CommandOutbind,
	"enquire-link":          CommandEnquireLink,
	"enquire-link-resp":     CommandEnquireLinkResp,
	"submit-multi":          CommandSubmitMulti,
	"submit-multi-resp":     CommandSubmitMultiResp,
	"alert-notification":    CommandAlertNotification,
	"data-sm":               CommandDataSm,
	"data-sm-resp":          CommandDataSmResp,
}

// CommandName returns the string representation for a CommandID
func CommandName(commandID CommandIDType) string {
	return pduCommandName[commandID]
}

// CommandIDFromString takes a command string name and returns the corresponding CommandIDType.
// The boolean is set to true if the commandName is understand; otherwise it is false, and the
// returned value for CommandIDType is undefined
func CommandIDFromString(commandName string) (CommandIDType, bool) {
	commandID, ok := commandNameToCommandID[commandName]
	return commandID, ok
}

// PDU is a PDU for
type PDU struct {
	CommandLength       uint32
	CommandID           CommandIDType
	CommandStatus       uint32
	SequenceNumber      uint32
	MandatoryParameters []*Parameter
	OptionalParameters  []*Parameter
}

// PDUDefinition describes a PDU.  It contains the set of mandatory Parameters and the
// minimum length (including the header and the mandatory Parameters).
type PDUDefinition struct {
	Type                CommandIDType
	MinLength           uint32
	MandatoryParameters []string
}

var pduTypeDefinition = map[CommandIDType]PDUDefinition{
	CommandGenericNack: PDUDefinition{CommandGenericNack, 0, []string{}},
	CommandBindReceiver: PDUDefinition{CommandBindReceiver, 0, []string{
		"system_id", "password", "system_type", "interface_version", "addr_ton",
		"addr_npi", "address_range",
	}},
	CommandBindReceiverResp: PDUDefinition{CommandBindReceiverResp, 0, []string{
		"system_id",
	}},
	CommandBindTransmitter: PDUDefinition{CommandBindTransmitter, 0, []string{
		"system_id", "password", "system_type", "interface_version", "addr_ton",
		"addr_npi", "address_range",
	}},
	CommandBindTransmitterResp: PDUDefinition{CommandBindTransmitterResp, 0, []string{
		"system_id",
	}},
	CommandQuerySm: PDUDefinition{CommandQuerySm, 0, []string{
		"message_id", "source_addr_ton", "source_addr_npi", "source_addr",
	}},
	CommandQuerySmResp: PDUDefinition{CommandQuerySmResp, 0, []string{
		"message_id", "final_date", "message_state", "error_code",
	}},
	CommandSubmitSm: PDUDefinition{CommandSubmitSm, 0, []string{
		"service_type", "source_addr_ton", "source_addr_npi", "source_addr",
		"dest_addr_ton", "dest_addr_npi", "destination_addr", "esm_class",
		"protocol_id", "priority_flag", "schedule_delivery_time", "validity_period",
		"registered_delivery", "replace_if_present_flag", "data_coding",
		"sm_default_msg_id", "sm_length", "short_message",
	}},
	CommandSubmitSmResp: PDUDefinition{CommandSubmitSmResp, 0, []string{
		"message_id",
	}},
	CommandDeliverSm: PDUDefinition{CommandDeliverSm, 0, []string{
		"service_type", "source_addr_ton", "source_addr_npi", "source_addr",
		"dest_addr_ton", "dest_addr_npi", "destination_addr", "esm_class",
		"protocol_id", "priority_flag", "schedule_delivery_time", "validity_period",
		"registered_delivery", "replace_if_present_flag", "data_coding",
		"sm_default_msg_id", "sm_length", "short_message",
	}},
	CommandDeliverSmResp: PDUDefinition{CommandDeliverSmResp, 0, []string{
		"message_id",
	}},
	CommandUnbind:     PDUDefinition{CommandUnbind, 0, []string{}},
	CommandUnbindResp: PDUDefinition{CommandUnbindResp, 0, []string{}},
	CommandReplaceSm: PDUDefinition{CommandReplaceSm, 0, []string{
		"message_id", "source_addr_ton", "source_addr_npi", "source_addr",
		"schedule_delivery_time", "validity_period", "registered_delivery",
		"registered_delivery", "sm_length", "short_message",
	}},
	CommandReplaceSmResp: PDUDefinition{CommandReplaceSmResp, 0, []string{}},
	CommandCancelSm: PDUDefinition{CommandCancelSm, 0, []string{
		"service_type", "message_id", "source_addr_ton", "source_addr_npi", "source_addr",
		"dest_addr_ton", "dest_addr_npi", "destination_addr",
	}},
	CommandCancelSmResp: PDUDefinition{CommandCancelSmResp, 0, []string{}},
	CommandBindTransceiver: PDUDefinition{CommandBindTransceiver, 0, []string{
		"system_id", "password", "system_type", "interface_version", "addr_ton",
		"addr_npi", "address_range",
	}},
	CommandBindTransceiverResp: PDUDefinition{CommandBindTransceiverResp, 0, []string{
		"system_id",
	}},
	CommandOutbind: PDUDefinition{CommandOutbind, 0, []string{
		"system_id", "password",
	}},
	CommandEnquireLink:     PDUDefinition{CommandEnquireLink, 0, []string{}},
	CommandEnquireLinkResp: PDUDefinition{CommandEnquireLinkResp, 0, []string{}},
	CommandSubmitMulti:     PDUDefinition{CommandSubmitMulti, 0, []string{}},
	CommandSubmitMultiResp: PDUDefinition{CommandSubmitMultiResp, 0, []string{}},
	CommandAlertNotification: PDUDefinition{CommandAlertNotification, 0, []string{
		"source_addr_ton", "source_addr_npi", "source_addr", "esme_addr_ton",
		"esme_addr_npi", "esme_addr",
	}},
	CommandDataSm: PDUDefinition{CommandDataSm, 26, []string{
		"service_type", "source_addr_ton", "source_addr_npi", "source_addr", "dest_addr_ton",
		"dest_addr_npi", "destination_addr", "esm_class", "registered_delivery", "data_coding",
	}},
	CommandDataSmResp: PDUDefinition{CommandDataSmResp, 0, []string{}},
}

// LengthOfNextPDU reads a stream that should contain at least a fragment of an SMPP PDU.
// If the length of stream is less than 4, then return 0 (meaning length is not yet known)
func LengthOfNextPDU(stream []byte) uint32 {
	if len(stream) < 4 {
		return uint32(0)
	}

	return binary.BigEndian.Uint32(stream[0:4])
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

// CommandName returns the string name for this PDU's CommandID
func (pdu *PDU) CommandName() string {
	return pduCommandName[pdu.CommandID]
}

// IsRequest returns true if the Command is in the request range (i.e., top-order bit is not set);
// otherwise it returns false
func (pdu *PDU) IsRequest() bool {
	if uint32(pdu.CommandID)&0x80000000 == 0 {
		return true
	}

	return false
}

// ComputeLength computes the enocde length of the PDU and
// returns it
func (pdu *PDU) ComputeLength() uint32 {
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
func (pdu *PDU) Encode() ([]byte, error) {
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

// DecodePDU accepts a byte stream in network byte order, and attempts to convert
// it to a PDU object
func DecodePDU(stream []byte) (*PDU, error) {
	if len(stream) < 16 {
		return nil, fmt.Errorf("Incoming stream invalid length, is (%d) octets", len(stream))
	}

	pduLength := uint32(binary.BigEndian.Uint32(stream[0:4]))

	if pduLength < 16 {
		return nil, fmt.Errorf("Stream length field value (%d) is less than minimum (16)", pduLength)
	}

	if pduLength > uint32(len(stream)) {
		return nil, fmt.Errorf("Stream length field value is (%d) but stream length is (%d)", pduLength, len(stream))
	}

	commandID := CommandIDType(uint32(binary.BigEndian.Uint32(stream[4:8])))

	pduDef, exists := pduTypeDefinition[commandID]

	if exists {
		if pduDef.MinLength > pduLength {
			return nil, fmt.Errorf("Stream length (%d) less than minimum (%d) for command type (%08x)", pduLength, pduDef.MinLength, commandID)
		}
	} else {
		return nil, fmt.Errorf("Stream command-id (%08x) not known", commandID)
	}

	status := uint32(binary.BigEndian.Uint32(stream[8:12]))
	sequenceNumber := uint32(binary.BigEndian.Uint32(stream[12:16]))

	mandatoryPList := list.New()
	optionalPList := list.New()

	s := 16
	smLength := uint8(0)
	smLengthFound := false

	for i := 0; i < len(pduDef.MandatoryParameters); i++ {
		if s >= int(pduLength) {
			break
		}

		paramName := pduDef.MandatoryParameters[i]
		paramDef := parameterTypeDefinition[paramName]

		switch paramDef.Type {
		case TypeUint8:
			mandatoryPList.PushBack(NewFLParameter(uint8(stream[s])))

			if paramName == "sm_length" {
				smLength = uint8(stream[s])
				smLengthFound = true
			}

			s++

		case TypeUint16:
			mandatoryPList.PushBack(NewFLParameter(binary.BigEndian.Uint16(stream[s : s+2])))
			s += 2

		case TypeUint32:
			mandatoryPList.PushBack(NewFLParameter(binary.BigEndian.Uint32(stream[s : s+4])))
			s += 4

		case TypeCOctetString:
			nullOffset := bytes.IndexByte(stream[s:], 0)

			if nullOffset < 0 {
				return nil, fmt.Errorf("Require C-String-Octet type but failed to find null terminator")
			}

			bb := stream[s : s+nullOffset]
			mandatoryPList.PushBack(NewCOctetStringParameter(string(bb)))
			s += nullOffset + 1

		case TypeOctetString:
			if paramName == "short_message" {
				if smLengthFound {
					if smLength > 0 {
						pp := &Parameter{TypeOctetString, uint32(smLength), stream[s : s+int(smLength)]}
						mandatoryPList.PushBack(pp)
						s += int(smLength)
					}
				} else {
					return nil, fmt.Errorf("Found short_message field but no sm_length field")
				}
			} else {
				return nil, fmt.Errorf("Unknown definition for type (%s)", paramName)
			}
		}
	}

	// Optional Parameters are all TLV
	for uint32(s) < pduLength {
		tlvTag := binary.BigEndian.Uint16(stream[s : s+2])
		tlvLen := binary.BigEndian.Uint16(stream[s+2 : s+4])
		tlvVal := stream[s+4 : s+4+int(tlvLen)]

		optionalPList.PushBack(NewTLVParameter(tlvTag, tlvVal))

		s += 4 + int(tlvLen)
	}

	mp := make([]*Parameter, mandatoryPList.Len())
	op := make([]*Parameter, optionalPList.Len())

	i := 0
	for e := mandatoryPList.Front(); e != nil; e = e.Next() {
		mp[i] = e.Value.(*Parameter)
		i++
	}

	i = 0
	for e := optionalPList.Front(); e != nil; e = e.Next() {
		op[i] = e.Value.(*Parameter)
		i++
	}

	return NewPDU(commandID, status, sequenceNumber, mp, op), nil
}
