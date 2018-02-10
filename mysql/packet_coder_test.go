package mysql

import (
	"bytes"
	"testing"
)

func TestEncodeOK(t *testing.T) {
	c := PacketCoder{capabilities: CLIENT_PROTOCOL_41}
	payload := c.encodeOK(0, 233, 233)
	expectedPayload := []byte{0, 233, 233, 0, 0, 0, 0}
	if len(payload) != len(expectedPayload) {
		t.Fatalf("bad result %v, expected: %v", len(payload), len(expectedPayload))
	}
	if !bytes.Equal(payload, expectedPayload) {
		t.Fatalf("bad result %v, expected: %v", payload, expectedPayload)
	}
}
