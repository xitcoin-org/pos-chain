package cli

import (
	"strings"
	"testing"
)

func TestNewTxCmdIncludesBridgeOperations(t *testing.T) {
	cmd := NewTxCmd()
	for _, name := range []string{"submit-attestation", "initiate-outbound", "initialize-route"} {
		child, _, err := cmd.Find([]string{name})
		if err != nil {
			t.Fatalf("find %s: %v", name, err)
		}
		if child == cmd || child.Name() != name {
			t.Fatalf("missing %s command", name)
		}
	}
}

func TestDecodeSignatures(t *testing.T) {
	valid := strings.Repeat("ab", 65)
	decoded, err := decodeSignatures([]string{"0x" + valid, valid})
	if err != nil {
		t.Fatalf("decode valid signatures: %v", err)
	}
	if len(decoded) != 2 || len(decoded[0]) != 65 || len(decoded[1]) != 65 {
		t.Fatalf("unexpected decoded signatures: %d", len(decoded))
	}
}

func TestDecodeSignaturesRejectsInvalidInput(t *testing.T) {
	valid := strings.Repeat("ab", 65)
	tests := [][]string{
		{valid},
		{valid, valid, valid, valid},
		{valid, "not-hex"},
		{valid, strings.Repeat("ab", 64)},
	}
	for _, values := range tests {
		if _, err := decodeSignatures(values); err == nil {
			t.Fatalf("expected error for %d signatures", len(values))
		}
	}
}
