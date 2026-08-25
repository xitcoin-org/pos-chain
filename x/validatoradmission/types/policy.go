package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultMaxApprovedValidators uint32 = 258
	DefaultMinimumSelfDelegation        = "5000000000000000000000000axtc"
)

func ValidatePolicy(maxApprovedValidators uint32, minimumSelfDelegation string) error {
	coin, err := validatePolicyAmount(maxApprovedValidators, minimumSelfDelegation)
	if err != nil {
		return err
	}
	if coin.Denom != "axtc" {
		return fmt.Errorf("minimum self delegation must use axtc")
	}

	return nil
}

// ValidateGenesisPolicy permits an isolated development chain to use its
// configured genesis denomination. Runtime policy updates remain restricted
// to axtc by ValidatePolicy.
func ValidateGenesisPolicy(maxApprovedValidators uint32, minimumSelfDelegation string) error {
	_, err := validatePolicyAmount(maxApprovedValidators, minimumSelfDelegation)
	return err
}

func validatePolicyAmount(maxApprovedValidators uint32, minimumSelfDelegation string) (sdk.Coin, error) {
	if maxApprovedValidators == 0 {
		return sdk.Coin{}, fmt.Errorf("max approved validators must be greater than zero")
	}

	coin, err := sdk.ParseCoinNormalized(minimumSelfDelegation)
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("invalid minimum self delegation: %w", err)
	}
	if !coin.IsValid() || !coin.IsPositive() {
		return sdk.Coin{}, fmt.Errorf("minimum self delegation must be positive")
	}

	return coin, nil
}

func DefaultPolicy() (uint32, string) {
	return DefaultMaxApprovedValidators, DefaultMinimumSelfDelegation
}
