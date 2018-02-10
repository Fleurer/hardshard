package mysql

import (
	"bytes"
	"testing"
)

func testEncodeOK(t *testing.T) {
	c := PacketCoder{capabilities: CLIENT_PROTOCOL_41}
	payload := c.encodeOK(0, 233, 233)
	expectedPayload := []byte{}
	if !bytes.Equal(payload, expectedPayload) {
		t.Fatalf("bad payload %s, expectedPayload: %s", payload, expectedPayload)
	}
}
