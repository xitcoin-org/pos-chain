package types

import (
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateParamsCommand is the deterministic domain command used by the future
// protobuf message server to request an authorized parameter update.
type UpdateParamsCommand struct {
	Authority string
	Params    Params
}

func (c UpdateParamsCommand) ValidateBasic() error {
	if err := validateCommandAuthority(c.Authority); err != nil {
		return err
	}
	if err := c.Params.Validate(); err != nil {
		return fmt.Errorf("invalid incentive parameters: %w", err)
	}
	return nil
}

// ActivatePeriodCommand defines the prefunded values required to activate one
// reward period. Amounts are canonical base-10 atomic axtc strings.
type ActivatePeriodCommand struct {
	Authority                   string
	EligibleBondedAtomic         string
	TreasuryBalanceAtomic        string
	CommittedAnnualBudgetAtomic string
}

func (c ActivatePeriodCommand) ValidateBasic() error {
	if err := validateCommandAuthority(c.Authority); err != nil {
		return err
	}

	eligible, treasury, budget, err := c.Amounts()
	if err != nil {
		return err
	}
	if !eligible.IsPositive() {
		return errors.New("eligible bonded amount must be positive")
	}
	if !treasury.IsPositive() {
		return errors.New("treasury balance must be positive")
	}
	if !budget.IsPositive() {
		return errors.New("committed annual budget must be positive")
	}
	return nil
}

func (c ActivatePeriodCommand) Amounts() (
	sdkmath.Int,
	sdkmath.Int,
	sdkmath.Int,
	error,
) {
	eligible, err := ParseStoredAtomicAmount(c.EligibleBondedAtomic)
	if err != nil {
		return sdkmath.Int{}, sdkmath.Int{}, sdkmath.Int{},
			fmt.Errorf("invalid eligible bonded amount: %w", err)
	}
	treasury, err := ParseStoredAtomicAmount(c.TreasuryBalanceAtomic)
	if err != nil {
		return sdkmath.Int{}, sdkmath.Int{}, sdkmath.Int{},
			fmt.Errorf("invalid treasury balance: %w", err)
	}
	budget, err := ParseStoredAtomicAmount(
		c.CommittedAnnualBudgetAtomic,
	)
	if err != nil {
		return sdkmath.Int{}, sdkmath.Int{}, sdkmath.Int{},
			fmt.Errorf("invalid committed annual budget: %w", err)
	}
	return eligible, treasury, budget, nil
}

func validateCommandAuthority(authority string) error {
	if authority == "" {
		return errors.New("authority is required")
	}
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}
	return nil
}
