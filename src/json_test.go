package smpp

import (
	"smpp"
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

	document, err := smpp.UnmarshallJSON(data)

	if err != nil {
		t.Errorf("On UnmarshallJSON for header-only message, received error: %s", err)
	}

	if document == nil {
		t.Errorf("On UnmarshallJSON for header-only message, document returned is nil")
	}

	if len(document.Messages) != 1 {
		t.Errorf("On UnmarshallJSON for header-only message, expected 1 message in document, got %d", len(document.Messages))
	} else {
		message := document.Messages[0]

		if message.CommandID != 1 {
			t.Errorf("On UnmarshallJSON for header-only message, expected CommandID = 1, got = %d", message.CommandID)
		}

		if message.CommandStatus != 0 {
			t.Errorf("On UnmarshallJSON for header-only message, expected CommandStatus = 0, got = %d", message.CommandStatus)
		}

		if message.SequenceNumber != 1 {
			t.Errorf("On UnmarshallJSON for header-only message, expected SequenceNumber = 1, got = %d", message.SequenceNumber)
		}

		if message.EncodedLength != 0 {
			t.Errorf("On UnmarshallJSON for header-only message, expected EncodedLength = 0, got = %d", message.EncodedLength)
		}
	}
}

func compareJSONDocumentObjects(expected *smpp.JSONDocument, received *smpp.JSONDocument) error {
	return nil
}
