package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
)

// ValidateBasic validates a funded incentive parameter update request.
func (m *MsgUpdateParams) ValidateBasic() error {
	if m == nil {
		return errors.New("update params message is nil")
	}

	return (UpdateParamsCommand{
		Authority: m.Authority,
		Params: Params{
			TreasuryReleaseRateBasisPoints: m.TreasuryReleaseRateBasisPoints,
			BlocksPerYear:                  m.BlocksPerYear,
			CalculationPeriodBlocks:        m.CalculationPeriodBlocks,
		},
	}).ValidateBasic()
}

// GetSignBytes returns legacy sign bytes.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(&m)
}
