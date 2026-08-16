package types

import (
	"encoding/hex"
	"errors"
	"strings"
)

// ParseAttestationID accepts the canonical 32-byte hexadecimal replay key.
func ParseAttestationID(value string) ([32]byte, string, error) {
	normalized := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), "0x")
	if len(normalized) != 64 {
		return [32]byte{}, "", errors.New("attestation ID must contain 32 hexadecimal bytes")
	}
	decoded, err := hex.DecodeString(normalized)
	if err != nil {
		return [32]byte{}, "", errors.New("attestation ID must contain 32 hexadecimal bytes")
	}
	var id [32]byte
	copy(id[:], decoded)
	return id, "0x" + normalized, nil
}
