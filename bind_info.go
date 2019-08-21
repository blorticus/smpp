package smpp

// BindType is an enumeration of Bind types
type BindType uint

// These are the possible bind types
const (
	TransceiverBind BindType = iota
	ReceiverBind
	TransmitterBind
)

// BindInfo contains information for a bind
type BindInfo struct {
	Type BindType
}
