package smpp

import "encoding/json"

/**
{
	messages: [
	{
		command_id: 1,
		sequence_number: 1,
		command_status: 0,
		encoded_length: 80,
		mandatory_parameters: {

		},
		optional_parameters: {

		}
	},
	{ ... },
	]
}
**/

// JSONMandatoryParameterMap describes the layout of the SMPP JSON file mandatory_parameters field
type JSONMandatoryParameterMap struct {
	AddrTon              uint8  `json:"addr_ton"`
	AddrNpi              uint8  `json:"addr_npi"`
	AddressRange         string `json:"address_range"`
	DataCoding           uint8  `json:"data_coding"`
	DestinationAddr      string `json:"destination_addr"`
	DestinationAddrNpi   uint8  `json:"destination_addr_npi"`
	DestinationAddrTon   uint8  `json:"destination_addr_ton"`
	EsmClass             uint8  `json:"esm_class"`
	InterfaceVersion     uint8  `json:"interface_version"`
	Password             string `json:"password"`
	MessageID            string `json:"message_id"`
	PriorityFlag         uint8  `json:"priority_flag"`
	ProtocolID           uint8  `json:"protocol_id"`
	RegisteredDelivery   uint8  `json:"registered_delivery"`
	ReplaceIfPresentFlag uint8  `json:"replace_if_present_flag"`
	ScheduleDeliveryTime string `json:"schedule_delivery_time"`
	ServiceType          string `json:"service_type"`
	ShortMessage         string `json:"short_message"`
	SmDefaultMsgID       uint8  `json:"sm_default_msg_id"`
	SmLength             uint8  `json:"sm_length"`
	SourceAddrNpi        uint8  `json:"source_addr_npi"`
	SourceAddrTon        uint8  `json:"source_addr_ton"`
	SourceAddr           string `json:"source_addr"`
	SystemID             string `json:"system_id"`
	SystemType           string `json:"system_type"`
	ValidityPeriod       string `json:"validity_period"`
}

// JSONOptionalParameterMap describes the layout of the SMPP JSON file optional_parameters field
type JSONOptionalParameterMap struct {
	SCInterfaceVersion       uint32 `yaml:"SC_interface_version"`
	AdditionalStatusInfoText uint32 `yaml:"additional_status_info_text"`
	AlertOnMessageDelivery   uint32 `yaml:"alert_on_message_delivery"`
	CallbackNum              uint32 `yaml:"callback_num"`
	CallbackNumAtag          uint32 `yaml:"callback_num_atag"`
	CallbackNumPresInd       uint32 `yaml:"callback_num_pres_ind"`
	DeliveryFailureReason    uint32 `yaml:"delivery_failure_reason"`
	DestAddrSubunit          uint32 `yaml:"dest_addr_subunit"`
	DestBearerType           uint32 `yaml:"dest_bearer_type"`
	DestNetworkType          uint32 `yaml:"dest_network_type"`
	DestSubaddress           uint32 `yaml:"dest_subaddress"`
	DestTelematicsID         uint32 `yaml:"dest_telematics_id"`
	DestinationPort          uint32 `yaml:"destination_port"`
	DisplayTime              uint32 `yaml:"display_time"`
	DpfResult                uint32 `yaml:"dpf_result"`
	ItsReplyType             uint32 `yaml:"its_reply_type"`
	ItsSessionInfo           uint32 `yaml:"its_session_info"`
	LanguageIndicator        uint32 `yaml:"language_indicator"`
	MessagePayload           uint32 `yaml:"message_payload"`
	MessageState             uint32 `yaml:"message_state"`
	MoreMessagesToSend       uint32 `yaml:"more_messages_to_send"`
	MsAvailabilityStatus     uint32 `yaml:"ms_availability_status"`
	MsMsgWaitFacilities      uint32 `yaml:"ms_msg_wait_facilities"`
	MsValidity               uint32 `yaml:"ms_validity"`
	NetworkErrorCode         uint32 `yaml:"network_error_code"`
	NumberOfMessages         uint32 `yaml:"number_of_messages"`
	PayloadType              uint32 `yaml:"payload_type"`
	PrivacyIndicator         uint32 `yaml:"privacy_indicator"`
	QosTimeToLive            uint32 `yaml:"qos_time_to_live"`
	ReceiptedMessageID       uint32 `yaml:"receipted_message_id"`
	SarMsgRefNum             uint32 `yaml:"sar_msg_ref_num"`
	SarSegmentSeqnum         uint32 `yaml:"sar_segment_seqnum"`
	SarTotalSegments         uint32 `yaml:"sar_total_segments"`
	SetDpf                   uint32 `yaml:"set_dpf"`
	SmsSignal                uint32 `yaml:"sms_signal"`
	SourceAddrSubunit        uint32 `yaml:"source_addr_subunit"`
	SourceBearerType         uint32 `yaml:"source_bearer_type"`
	SourceNetworkType        uint32 `yaml:"source_network_type"`
	SourcePort               uint32 `yaml:"source_port"`
	SourceSubaddress         uint32 `yaml:"source_subaddress"`
	SourceTelematicsID       uint32 `yaml:"source_telematics_id"`
	UserMessageReference     uint32 `yaml:"user_message_reference"`
	UserResponseCode         uint32 `yaml:"user_response_code"`
	UssdServiceOp            uint32 `yaml:"ussd_service_op"`
}

// JSONMessage describes the SMPP JSON message structure
type JSONMessage struct {
	CommandID           uint32                    `json:"command_id"`
	SequenceNumber      uint32                    `json:"sequence_number"`
	CommandStatus       uint32                    `json:"command_status"`
	EncodedLength       uint32                    `json:"encoded_length"`
	MandatoryParameters JSONMandatoryParameterMap `json:"mandatory_parameters"`
	OptionalParameters  JSONOptionalParameterMap  `json:"optional_parameters"`
}

// JSONDocument describes the SMPP JSON document structure
type JSONDocument struct {
	Messages []JSONMessage `json:"messages"`
}

// UnmarshallJSON treats 'data' as a stream of characters in a JSON struct, and
// attempts to convert them into a JSONMessage
func UnmarshallJSON(data []byte) (*JSONDocument, error) {
	var jsonDocument JSONDocument

	err := json.Unmarshal(data, &jsonDocument)

	return &jsonDocument, err
}

// ConvertJSONToPDUs creates PDUs from well-formed JSONMessage objects
func ConvertJSONToPDUs(jsonMessages []*JSONMessage) ([]*PDU, error) {
	pdus := make([]*PDU, len(jsonMessages))

	for i, nextJSONMessage := range jsonMessages {
		mandatoryParameterList, err := produceParametersFromJSONMandatoryParameters(&nextJSONMessage.MandatoryParameters)

		if err != nil {
			return nil, err
		}

		optionalParameterList, err := produceParametersFromJSONMandatoryParameters(&nextJSONMessage.MandatoryParameters)

		if err != nil {
			return nil, err
		}

		pdus[i] = NewPDU(
			CommandIDType(nextJSONMessage.CommandID),
			nextJSONMessage.CommandStatus,
			nextJSONMessage.SequenceNumber,
			mandatoryParameterList,
			optionalParameterList,
		)
	}

	return pdus, nil
}

func produceParametersFromJSONMandatoryParameters(json *JSONMandatoryParameterMap) ([]*Parameter, error) {
	return nil, nil
}

func produceParametersFromJSONOptionalParameters(json *JSONOptionalParameterMap) ([]*Parameter, error) {
	return nil, nil
}
