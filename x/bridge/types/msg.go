package types

import (
	"errors"
	"math/big"
	"strings"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

var _ types.Msg = &MsgSubmitAttestation{}
var _ types.Msg = &MsgInitiateOutboundTransfer{}

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

func (m *MsgInitiateOutboundTransfer) ValidateBasic() error {
	if m == nil {
		return errors.New("bridge: empty outbound transfer")
	}
	if _, err := types.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if !validRouteID(m.RouteId) {
		return errors.New("invalid route ID")
	}
	destination := strings.TrimSpace(m.Destination)
	if !common.IsHexAddress(destination) || common.HexToAddress(destination) == (common.Address{}) {
		return errors.New("invalid Cronos destination")
	}
	amount, ok := new(big.Int).SetString(m.Amount, 10)
	if !ok || amount.Sign() <= 0 {
		return errors.New("amount must be a positive integer in atomic units")
	}
	return nil
}

func (m MsgInitiateOutboundTransfer) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(&m)
}
