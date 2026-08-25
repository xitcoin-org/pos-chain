package types

import (
	"strings"
	"testing"
)

func TestParseAttestationIDCanonicalizesHex(t *testing.T) {
	input := "  0x" + strings.Repeat("AB", 32) + "  "
	id, canonical, err := ParseAttestationID(input)
	if err != nil {
		t.Fatalf("valid ID rejected: %v", err)
	}
	if canonical != "0x"+strings.Repeat("ab", 32) {
		t.Fatalf("canonical ID = %q", canonical)
	}
	if id[0] != 0xab || id[31] != 0xab {
		t.Fatalf("decoded ID is incorrect: %x", id)
	}
}

func TestParseAttestationIDRejectsMalformedValues(t *testing.T) {
	for _, value := range []string{"", "0x01", strings.Repeat("z", 64), strings.Repeat("a", 66)} {
		if _, _, err := ParseAttestationID(value); err == nil {
			t.Fatalf("malformed ID %q accepted", value)
		}
	}
}
