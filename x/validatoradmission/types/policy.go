package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultMaxApprovedValidators uint32 = 258
	DefaultMinimumSelfDelegation        = "1000000000000000000000000axtc"
)

func ValidatePolicy(maxApprovedValidators uint32, minimumSelfDelegation string) error {
	if maxApprovedValidators == 0 {
		return fmt.Errorf("max approved validators must be greater than zero")
	}

	coin, err := sdk.ParseCoinNormalized(minimumSelfDelegation)
	if err != nil {
		return fmt.Errorf("invalid minimum self delegation: %w", err)
	}
	if coin.Denom != "axtc" {
		return fmt.Errorf("minimum self delegation must use axtc")
	}
	if !coin.IsValid() || !coin.IsPositive() {
		return fmt.Errorf("minimum self delegation must be positive")
	}

	return nil
}

func DefaultPolicy() (uint32, string) {
	return DefaultMaxApprovedValidators, DefaultMinimumSelfDelegation
}
