package smpp

import (
	"fmt"
	"testing"
)

func TestHeaderOnlyMessage(t *testing.T) {
	data := []byte(`
{
	"messages": [
		{
			"command_id": 1,
			"command_status": 0,
			"sequence_number": 1,
			"encoded_length": 0
		}
	]
}	
`)

	document, err := UnmarshallJSON(data)

	if err != nil {
		t.Errorf("On UnmarshallJSON for header-only message, received error: %s", err)
	}

	err = compareJSONDocumentObjects(&JSONDocument{
		Messages: []JSONMessage{
			{
				CommandID:      1,
				CommandStatus:  0,
				EncodedLength:  0,
				SequenceNumber: 1,
			},
		},
	}, document)

	if err != nil {
		t.Errorf("On UnmarshallJSON for header-only message, %s", err)
	}
}

func compareJSONDocumentObjects(expected *JSONDocument, received *JSONDocument) error {
	if received == nil {
		return fmt.Errorf("expected document, got nil document")
	}

	if len(expected.Messages) != len(received.Messages) {
		return fmt.Errorf("expected %d messages, got %d", len(expected.Messages), len(received.Messages))
	}

	for i, expectedMsg := range expected.Messages {
		receivedMsg := received.Messages[i]

		if expectedMsg.CommandID != receivedMsg.CommandID {
			return fmt.Errorf("on message [%d], expected CommandID = %d, got = %d", i, expectedMsg.CommandID, receivedMsg.CommandID)
		}

		if expectedMsg.CommandStatus != receivedMsg.CommandStatus {
			return fmt.Errorf("on message [%d], expected CommandStatus = %d, got = %d", i, expectedMsg.CommandStatus, receivedMsg.CommandStatus)
		}

		if expectedMsg.SequenceNumber != receivedMsg.SequenceNumber {
			return fmt.Errorf("on message [%d], expected SequenceNumber = %d, got = %d", i, expectedMsg.SequenceNumber, receivedMsg.SequenceNumber)
		}

		if expectedMsg.EncodedLength != receivedMsg.EncodedLength {
			return fmt.Errorf("on message [%d], expected EncodedLength = %d, got = %d", i, expectedMsg.EncodedLength, receivedMsg.EncodedLength)
		}

		if err := compareJSONMandatoryParameters(expectedMsg.MandatoryParameters, receivedMsg.MandatoryParameters); err != nil {
			return fmt.Errorf("on message [%d], %s", i, err)
		}

		if err := compareJSONOptionalParameters(expectedMsg.OptionalParameters, receivedMsg.OptionalParameters); err != nil {
			return fmt.Errorf("on message [%d], %s", i, err)
		}

	}

	return nil
}

func compareJSONMandatoryParameters(expected JSONMandatoryParameterMap, received JSONMandatoryParameterMap) error {
	if expected.AddrTon != received.AddrTon {
		return fmt.Errorf("for MandatoryParameters expected AddrTon = (%d), got = (%d)", expected.AddrTon, received.AddrTon)
	}
	if expected.AddrNpi != received.AddrNpi {
		return fmt.Errorf("for MandatoryParameters expected AddrNpi = (%d), got = (%d)", expected.AddrNpi, received.AddrNpi)
	}
	if expected.AddressRange != received.AddressRange {
		return fmt.Errorf("for MandatoryParameters expected AddressRange = (%s), got = (%s)", expected.AddressRange, received.AddressRange)
	}
	if expected.DataCoding != received.DataCoding {
		return fmt.Errorf("for MandatoryParameters expected DataCoding = (%d), got = (%d)", expected.DataCoding, received.DataCoding)
	}
	if expected.DestinationAddr != received.DestinationAddr {
		return fmt.Errorf("for MandatoryParameters expected DestinationAddr = (%s), got = (%s)", expected.DestinationAddr, received.DestinationAddr)
	}
	if expected.DestinationAddrNpi != received.DestinationAddrNpi {
		return fmt.Errorf("for MandatoryParameters expected DestinationAddrNpi = (%d), got = (%d)", expected.DestinationAddrNpi, received.DestinationAddrNpi)
	}
	if expected.DestinationAddrTon != received.DestinationAddrTon {
		return fmt.Errorf("for MandatoryParameters expected DestinationAddrTon = (%d), got = (%d)", expected.DestinationAddrTon, received.DestinationAddrTon)
	}
	if expected.EsmClass != received.EsmClass {
		return fmt.Errorf("for MandatoryParameters expected EsmClass = (%d), got = (%d)", expected.EsmClass, received.EsmClass)
	}
	if expected.InterfaceVersion != received.InterfaceVersion {
		return fmt.Errorf("for MandatoryParameters expected InterfaceVersion = (%d), got = (%d)", expected.InterfaceVersion, received.InterfaceVersion)
	}
	if expected.Password != received.Password {
		return fmt.Errorf("for MandatoryParameters expected Password = (%s), got = (%s)", expected.Password, received.Password)
	}
	if expected.MessageID != received.MessageID {
		return fmt.Errorf("for MandatoryParameters expected MessageID = (%s), got = (%s)", expected.MessageID, received.MessageID)
	}
	if expected.PriorityFlag != received.PriorityFlag {
		return fmt.Errorf("for MandatoryParameters expected PriorityFlag = (%d), got = (%d)", expected.PriorityFlag, received.PriorityFlag)
	}
	if expected.ProtocolID != received.ProtocolID {
		return fmt.Errorf("for MandatoryParameters expected ProtocolID = (%d), got = (%d)", expected.ProtocolID, received.ProtocolID)
	}
	if expected.RegisteredDelivery != received.RegisteredDelivery {
		return fmt.Errorf("for MandatoryParameters expected RegisteredDelivery = (%d), got = (%d)", expected.RegisteredDelivery, received.RegisteredDelivery)
	}
	if expected.ReplaceIfPresentFlag != received.ReplaceIfPresentFlag {
		return fmt.Errorf("for MandatoryParameters expected ReplaceIfPresentFlag = (%d), got = (%d)", expected.ReplaceIfPresentFlag, received.ReplaceIfPresentFlag)
	}
	if expected.ScheduleDeliveryTime != received.ScheduleDeliveryTime {
		return fmt.Errorf("for MandatoryParameters expected ScheduleDeliveryTime = (%s), got = (%s)", expected.ScheduleDeliveryTime, received.ScheduleDeliveryTime)
	}
	if expected.ServiceType != received.ServiceType {
		return fmt.Errorf("for MandatoryParameters expected ServiceType = (%s), got = (%s)", expected.ServiceType, received.ServiceType)
	}
	if expected.ShortMessage != received.ShortMessage {
		return fmt.Errorf("for MandatoryParameters expected ShortMessage = (%s), got = (%s)", expected.ShortMessage, received.ShortMessage)
	}
	if expected.SmDefaultMsgID != received.SmDefaultMsgID {
		return fmt.Errorf("for MandatoryParameters expected SmDefaultMsgID = (%d), got = (%d)", expected.SmDefaultMsgID, received.SmDefaultMsgID)
	}
	if expected.SmLength != received.SmLength {
		return fmt.Errorf("for MandatoryParameters expected SmLength = (%d), got = (%d)", expected.SmLength, received.SmLength)
	}
	if expected.SourceAddrNpi != received.SourceAddrNpi {
		return fmt.Errorf("for MandatoryParameters expected SourceAddrNpi = (%d), got = (%d)", expected.SourceAddrNpi, received.SourceAddrNpi)
	}
	if expected.SourceAddrTon != received.SourceAddrTon {
		return fmt.Errorf("for MandatoryParameters expected SourceAddrTon = (%d), got = (%d)", expected.SourceAddrTon, received.SourceAddrTon)
	}
	if expected.SourceAddr != received.SourceAddr {
		return fmt.Errorf("for MandatoryParameters expected SourceAddr = (%s), got = (%s)", expected.SourceAddr, received.SourceAddr)
	}
	if expected.SystemID != received.SystemID {
		return fmt.Errorf("for MandatoryParameters expected SystemID = (%s), got = (%s)", expected.SystemID, received.SystemID)
	}
	if expected.SystemType != received.SystemType {
		return fmt.Errorf("for MandatoryParameters expected SystemType = (%s), got = (%s)", expected.SystemType, received.SystemType)
	}
	if expected.ValidityPeriod != received.ValidityPeriod {
		return fmt.Errorf("for MandatoryParameters expected ValidityPeriod = (%s), got = (%s)", expected.ValidityPeriod, received.ValidityPeriod)
	}

	return nil
}

func compareJSONOptionalParameters(expected JSONOptionalParameterMap, received JSONOptionalParameterMap) error {
	if expected.SCInterfaceVersion != received.SCInterfaceVersion {
		return fmt.Errorf("for OptionalParameters expected SCInterfaceVersion = (%d), got = (%d)", expected.SCInterfaceVersion, received.SCInterfaceVersion)
	}
	if expected.AdditionalStatusInfoText != received.AdditionalStatusInfoText {
		return fmt.Errorf("for OptionalParameters expected AdditionalStatusInfoText = (%d), got = (%d)", expected.AdditionalStatusInfoText, received.AdditionalStatusInfoText)
	}
	if expected.AlertOnMessageDelivery != received.AlertOnMessageDelivery {
		return fmt.Errorf("for OptionalParameters expected AlertOnMessageDelivery = (%d), got = (%d)", expected.AlertOnMessageDelivery, received.AlertOnMessageDelivery)
	}
	if expected.CallbackNum != received.CallbackNum {
		return fmt.Errorf("for OptionalParameters expected CallbackNum = (%d), got = (%d)", expected.CallbackNum, received.CallbackNum)
	}
	if expected.CallbackNumAtag != received.CallbackNumAtag {
		return fmt.Errorf("for OptionalParameters expected CallbackNumAtag = (%d), got = (%d)", expected.CallbackNumAtag, received.CallbackNumAtag)
	}
	if expected.CallbackNumPresInd != received.CallbackNumPresInd {
		return fmt.Errorf("for OptionalParameters expected CallbackNumPresInd = (%d), got = (%d)", expected.CallbackNumPresInd, received.CallbackNumPresInd)
	}
	if expected.DeliveryFailureReason != received.DeliveryFailureReason {
		return fmt.Errorf("for OptionalParameters expected DeliveryFailureReason = (%d), got = (%d)", expected.DeliveryFailureReason, received.DeliveryFailureReason)
	}
	if expected.DestAddrSubunit != received.DestAddrSubunit {
		return fmt.Errorf("for OptionalParameters expected DestAddrSubunit = (%d), got = (%d)", expected.DestAddrSubunit, received.DestAddrSubunit)
	}
	if expected.DestBearerType != received.DestBearerType {
		return fmt.Errorf("for OptionalParameters expected DestBearerType = (%d), got = (%d)", expected.DestBearerType, received.DestBearerType)
	}
	if expected.DestNetworkType != received.DestNetworkType {
		return fmt.Errorf("for OptionalParameters expected DestNetworkType = (%d), got = (%d)", expected.DestNetworkType, received.DestNetworkType)
	}
	if expected.DestSubaddress != received.DestSubaddress {
		return fmt.Errorf("for OptionalParameters expected DestSubaddress = (%d), got = (%d)", expected.DestSubaddress, received.DestSubaddress)
	}
	if expected.DestTelematicsID != received.DestTelematicsID {
		return fmt.Errorf("for OptionalParameters expected DestTelematicsID = (%d), got = (%d)", expected.DestTelematicsID, received.DestTelematicsID)
	}
	if expected.DestinationPort != received.DestinationPort {
		return fmt.Errorf("for OptionalParameters expected DestinationPort = (%d), got = (%d)", expected.DestinationPort, received.DestinationPort)
	}
	if expected.DisplayTime != received.DisplayTime {
		return fmt.Errorf("for OptionalParameters expected DisplayTime = (%d), got = (%d)", expected.DisplayTime, received.DisplayTime)
	}
	if expected.DpfResult != received.DpfResult {
		return fmt.Errorf("for OptionalParameters expected DpfResult = (%d), got = (%d)", expected.DpfResult, received.DpfResult)
	}
	if expected.ItsReplyType != received.ItsReplyType {
		return fmt.Errorf("for OptionalParameters expected ItsReplyType = (%d), got = (%d)", expected.ItsReplyType, received.ItsReplyType)
	}
	if expected.ItsSessionInfo != received.ItsSessionInfo {
		return fmt.Errorf("for OptionalParameters expected ItsSessionInfo = (%d), got = (%d)", expected.ItsSessionInfo, received.ItsSessionInfo)
	}
	if expected.LanguageIndicator != received.LanguageIndicator {
		return fmt.Errorf("for OptionalParameters expected LanguageIndicator = (%d), got = (%d)", expected.LanguageIndicator, received.LanguageIndicator)
	}
	if expected.MessagePayload != received.MessagePayload {
		return fmt.Errorf("for OptionalParameters expected MessagePayload = (%d), got = (%d)", expected.MessagePayload, received.MessagePayload)
	}
	if expected.MessageState != received.MessageState {
		return fmt.Errorf("for OptionalParameters expected MessageState = (%d), got = (%d)", expected.MessageState, received.MessageState)
	}
	if expected.MoreMessagesToSend != received.MoreMessagesToSend {
		return fmt.Errorf("for OptionalParameters expected MoreMessagesToSend = (%d), got = (%d)", expected.MoreMessagesToSend, received.MoreMessagesToSend)
	}
	if expected.MsAvailabilityStatus != received.MsAvailabilityStatus {
		return fmt.Errorf("for OptionalParameters expected MsAvailabilityStatus = (%d), got = (%d)", expected.MsAvailabilityStatus, received.MsAvailabilityStatus)
	}
	if expected.MsMsgWaitFacilities != received.MsMsgWaitFacilities {
		return fmt.Errorf("for OptionalParameters expected MsMsgWaitFacilities = (%d), got = (%d)", expected.MsMsgWaitFacilities, received.MsMsgWaitFacilities)
	}
	if expected.MsValidity != received.MsValidity {
		return fmt.Errorf("for OptionalParameters expected MsValidity = (%d), got = (%d)", expected.MsValidity, received.MsValidity)
	}
	if expected.NetworkErrorCode != received.NetworkErrorCode {
		return fmt.Errorf("for OptionalParameters expected NetworkErrorCode = (%d), got = (%d)", expected.NetworkErrorCode, received.NetworkErrorCode)
	}
	if expected.NumberOfMessages != received.NumberOfMessages {
		return fmt.Errorf("for OptionalParameters expected NumberOfMessages = (%d), got = (%d)", expected.NumberOfMessages, received.NumberOfMessages)
	}
	if expected.PayloadType != received.PayloadType {
		return fmt.Errorf("for OptionalParameters expected PayloadType = (%d), got = (%d)", expected.PayloadType, received.PayloadType)
	}
	if expected.PrivacyIndicator != received.PrivacyIndicator {
		return fmt.Errorf("for OptionalParameters expected PrivacyIndicator = (%d), got = (%d)", expected.PrivacyIndicator, received.PrivacyIndicator)
	}
	if expected.QosTimeToLive != received.QosTimeToLive {
		return fmt.Errorf("for OptionalParameters expected QosTimeToLive = (%d), got = (%d)", expected.QosTimeToLive, received.QosTimeToLive)
	}
	if expected.ReceiptedMessageID != received.ReceiptedMessageID {
		return fmt.Errorf("for OptionalParameters expected ReceiptedMessageID = (%d), got = (%d)", expected.ReceiptedMessageID, received.ReceiptedMessageID)
	}
	if expected.SarMsgRefNum != received.SarMsgRefNum {
		return fmt.Errorf("for OptionalParameters expected SarMsgRefNum = (%d), got = (%d)", expected.SarMsgRefNum, received.SarMsgRefNum)
	}
	if expected.SarSegmentSeqnum != received.SarSegmentSeqnum {
		return fmt.Errorf("for OptionalParameters expected SarSegmentSeqnum = (%d), got = (%d)", expected.SarSegmentSeqnum, received.SarSegmentSeqnum)
	}
	if expected.SarTotalSegments != received.SarTotalSegments {
		return fmt.Errorf("for OptionalParameters expected SarTotalSegments = (%d), got = (%d)", expected.SarTotalSegments, received.SarTotalSegments)
	}
	if expected.SetDpf != received.SetDpf {
		return fmt.Errorf("for OptionalParameters expected SetDpf = (%d), got = (%d)", expected.SetDpf, received.SetDpf)
	}
	if expected.SmsSignal != received.SmsSignal {
		return fmt.Errorf("for OptionalParameters expected SmsSignal = (%d), got = (%d)", expected.SmsSignal, received.SmsSignal)
	}
	if expected.SourceAddrSubunit != received.SourceAddrSubunit {
		return fmt.Errorf("for OptionalParameters expected SourceAddrSubunit = (%d), got = (%d)", expected.SourceAddrSubunit, received.SourceAddrSubunit)
	}
	if expected.SourceBearerType != received.SourceBearerType {
		return fmt.Errorf("for OptionalParameters expected SourceBearerType = (%d), got = (%d)", expected.SourceBearerType, received.SourceBearerType)
	}
	if expected.SourceNetworkType != received.SourceNetworkType {
		return fmt.Errorf("for OptionalParameters expected SourceNetworkType = (%d), got = (%d)", expected.SourceNetworkType, received.SourceNetworkType)
	}
	if expected.SourcePort != received.SourcePort {
		return fmt.Errorf("for OptionalParameters expected SourcePort = (%d), got = (%d)", expected.SourcePort, received.SourcePort)
	}
	if expected.SourceSubaddress != received.SourceSubaddress {
		return fmt.Errorf("for OptionalParameters expected SourceSubaddress = (%d), got = (%d)", expected.SourceSubaddress, received.SourceSubaddress)
	}
	if expected.SourceTelematicsID != received.SourceTelematicsID {
		return fmt.Errorf("for OptionalParameters expected SourceTelematicsID = (%d), got = (%d)", expected.SourceTelematicsID, received.SourceTelematicsID)
	}
	if expected.UserMessageReference != received.UserMessageReference {
		return fmt.Errorf("for OptionalParameters expected UserMessageReference = (%d), got = (%d)", expected.UserMessageReference, received.UserMessageReference)
	}
	if expected.UserResponseCode != received.UserResponseCode {
		return fmt.Errorf("for OptionalParameters expected UserResponseCode = (%d), got = (%d)", expected.UserResponseCode, received.UserResponseCode)
	}
	if expected.UssdServiceOp != received.UssdServiceOp {
		return fmt.Errorf("for OptionalParameters expected UssdServiceOp = (%d), got = (%d)", expected.UssdServiceOp, received.UssdServiceOp)
	}

	return nil
}
