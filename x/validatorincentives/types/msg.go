package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgActivateFundedPeriod{}
)

// ValidateBasic validates a funded incentive parameter update request.
func (m *MsgUpdateParams) ValidateBasic() error {
	if m == nil {
		return errors.New("update params message is nil")
	}

	return (UpdateParamsCommand{
		Authority: m.Authority,
		Params: Params{
			AnnualRateBasisPoints: m.AnnualRateBasisPoints,
			BlocksPerYear:         m.BlocksPerYear,
			RewardPeriodBlocks:    m.RewardPeriodBlocks,
		},
	}).ValidateBasic()
}

// GetSignBytes returns legacy sign bytes.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(&m)
}

// ValidateBasic validates a prefunded period activation request. Eligible
// bonded stake and treasury balance are deliberately not accepted from the
// transaction and are verified by the keeper from on-chain state.
func (m *MsgActivateFundedPeriod) ValidateBasic() error {
	if m == nil {
		return errors.New("activate funded period message is nil")
	}
	if err := validateCommandAuthority(m.Authority); err != nil {
		return err
	}

	budget, err := ParseStoredAtomicAmount(
		m.CommittedAnnualBudgetAtomic,
	)
	if err != nil {
		return fmt.Errorf(
			"invalid committed annual budget: %w",
			err,
		)
	}
	if !budget.IsPositive() {
		return errors.New(
			"committed annual budget must be positive",
		)
	}
	return nil
}

// GetSignBytes returns legacy sign bytes.
func (m MsgActivateFundedPeriod) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(&m)
}
