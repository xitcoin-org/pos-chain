package types

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = &MsgSubmitAttestation{}

// Attestation returns the internal attestation represented by this message.
func (m MsgSubmitAttestation) Attestation() Attestation {
	return Attestation{
		RouteID:       m.RouteId,
		Direction:     Direction(m.Direction),
		SourceChainID: m.SourceChainId,
		SourceRef:     m.SourceRef,
		Nonce:         m.Nonce,
		Destination:   m.Destination,
		Amount:        m.Amount,
		DeadlineUnix:  m.DeadlineUnix,
	}
}

// ValidateBasic validates the submitting account and unsigned message fields.
func (m *MsgSubmitAttestation) ValidateBasic() error {
	if m == nil {
		return errors.New("bridge: empty attestation submission")
	}
	if _, err := types.AccAddressFromBech32(m.Submitter); err != nil {
		return err
	}
	return m.Attestation().Validate()
}

// GetSignBytes returns legacy sign bytes for compatibility with the SDK.
func (m MsgSubmitAttestation) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(&m)
}
