package smpp

import (
	"net"
	"testing"
	"time"
)

type fakeNetConn struct {
	nextReadValue []byte
	nextReadError error

	bindTrasceiver01Msg []byte
	enquireLink01Msg    []byte
}

func newFakeNetConn() *fakeNetConn {
	conn := &fakeNetConn{
		bindTrasceiver01Msg: []byte{
			0, 0, 0, 0x20, // len = 32
			0, 0, 0, 0x02, // command = bind_trasceiver
			0, 0, 0, 0x00, // status code = 0
			0, 0, 0, 0x01, // seq number = 1
			0x66, 0x6f, 0x6f, 0, // systemID = 'foo'
			0x62, 0x61, 0x72, 0, // passwd = 'bar'
			0x62, 0x61, 0x7a, 0, // systemType = 'baz'
			0x34, // version
			0,    // addr_ton
			0,    // addr_npi
			0,    // address_range
		},

		enquireLink01Msg: []byte{
			0, 0, 0, 0x10, // len = 16
			0, 0, 0, 0x15, // command = eqnuire_link
			0, 0, 0, 0x00, // status code = 0
			0, 0, 0, 0x02, // seq number = 2
		},
	}

	return conn
}

func (conn *fakeNetConn) Read(b []byte) (int, error) {
	if conn.nextReadError != nil {
		return 0, conn.nextReadError
	}

	copy(b, conn.nextReadValue)

	return len(conn.nextReadValue), nil
}

func (conn *fakeNetConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (conn *fakeNetConn) Close() error {
	return nil
}

func (conn *fakeNetConn) LocalAddr() net.Addr {
	return nil
}

func (conn *fakeNetConn) RemoteAddr() net.Addr {
	return nil
}

func (conn *fakeNetConn) SetDeadline(t time.Time) error {
	return nil
}

func (conn *fakeNetConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (conn *fakeNetConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestOneMessageOneRead(t *testing.T) {
	conn := newFakeNetConn()
	conn.nextReadError = nil
	conn.nextReadValue = conn.bindTrasceiver01Msg

	reader := NewNetworkStreamReader(conn)

	pdus, err := reader.Read()

	if err != nil {
		t.Fatalf("Expected no error on Read(), but got error = (%s)", err)
	}

	if len(pdus) != 1 {
		t.Fatalf("Expected one PDU from Read(), but got (%d)", len(pdus))
	}
}

func TestOneMessageTwoReads(t *testing.T) {
	conn := newFakeNetConn()
	conn.nextReadError = nil
	conn.nextReadValue = conn.bindTrasceiver01Msg[:20]

	reader := NewNetworkStreamReader(conn)

	pdus, err := reader.Read()

	if err != nil {
		t.Fatalf("Expected no error on first Read(), but got error = (%s)", err)
	}

	if len(pdus) != 0 {
		t.Fatalf("Expected no PDU from first Read(), but got (%d)", len(pdus))
	}

	conn.nextReadValue = conn.bindTrasceiver01Msg[20:]

	pdus, err = reader.Read()

	if err != nil {
		t.Fatalf("Expected no error on second Read(), but got error = (%s)", err)
	}

	if len(pdus) != 1 {
		t.Fatalf("Expected one PDU from second Read(), but got (%d)", len(pdus))
	}
}

func TestTwoMessagesOneRead(t *testing.T) {
	conn := newFakeNetConn()
	conn.nextReadError = nil
	conn.nextReadValue = append(conn.bindTrasceiver01Msg, conn.enquireLink01Msg...)

	reader := NewNetworkStreamReader(conn)

	pdus, err := reader.Read()

	if err != nil {
		t.Fatalf("Expected no error on Read(), but got error = (%s)", err)
	}

	if len(pdus) != 2 {
		t.Fatalf("Expected two PDUs from Read(), but got (%d)", len(pdus))
	}
}

func TestTwoMessagesTwoReads(t *testing.T) {
	conn := newFakeNetConn()
	conn.nextReadError = nil

	msgSet := append(conn.bindTrasceiver01Msg, conn.enquireLink01Msg...)

	reader := NewNetworkStreamReader(conn)

	conn.nextReadValue = msgSet[:12]
	pdus, err := reader.Read()

	if err != nil {
		t.Fatalf("Expected no error on first Read(), but got error = (%s)", err)
	}

	if len(pdus) != 0 {
		t.Fatalf("Expected zero PDUs from first Read(), but got (%d)", len(pdus))
	}

	conn.nextReadValue = msgSet[12:31]
	pdus, err = reader.Read()

	if err != nil {
		t.Fatalf("Expected no error on second Read(), but got error = (%s)", err)
	}

	if len(pdus) != 0 {
		t.Fatalf("Expected zero PDUs from second Read(), but got (%d)", len(pdus))
	}

	conn.nextReadValue = msgSet[31:40]
	pdus, err = reader.Read()

	if err != nil {
		t.Fatalf("Expected no error on third Read(), but got error = (%s)", err)
	}

	if len(pdus) != 1 {
		t.Fatalf("Expected zero PDUs from third Read(), but got (%d)", len(pdus))
	}

	conn.nextReadValue = msgSet[40:]
	pdus, err = reader.Read()

	if err != nil {
		t.Fatalf("Expected no error on fourth Read(), but got error = (%s)", err)
	}

	if len(pdus) != 1 {
		t.Fatalf("Expected zero PDUs from fourth Read(), but got (%d)", len(pdus))
	}
}
