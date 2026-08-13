package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Authority             string   `json:"authority"`
	ApprovedValidators    []string `json:"approved_validators"`
	MaxApprovedValidators uint32   `json:"max_approved_validators"`
	MinimumSelfDelegation string   `json:"minimum_self_delegation"`
}

func DefaultGenesisState() GenesisState {
	maxApprovedValidators, minimumSelfDelegation := DefaultPolicy()
	return GenesisState{
		Authority:             "",
		ApprovedValidators:    []string{},
		MaxApprovedValidators: maxApprovedValidators,
		MinimumSelfDelegation: minimumSelfDelegation,
	}
}

func (g GenesisState) Policy() (uint32, string) {
	maxApprovedValidators := g.MaxApprovedValidators
	minimumSelfDelegation := g.MinimumSelfDelegation
	defaultMax, defaultMinimum := DefaultPolicy()

	if maxApprovedValidators == 0 {
		maxApprovedValidators = defaultMax
	}
	if minimumSelfDelegation == "" {
		minimumSelfDelegation = defaultMinimum
	}

	return maxApprovedValidators, minimumSelfDelegation
}

func (g GenesisState) Validate() error {
	authority := strings.TrimSpace(g.Authority)
	if authority == "" {
		// An entirely unconfigured module is disabled. This keeps default
		// application genesis valid without granting any admission authority.
		if len(g.ApprovedValidators) != 0 {
			return fmt.Errorf("validator admission authority is required")
		}
		maxApprovedValidators, minimumSelfDelegation := g.Policy()
		return ValidatePolicy(maxApprovedValidators, minimumSelfDelegation)
	}
	if authority != g.Authority {
		return fmt.Errorf("validator admission authority contains surrounding spaces")
	}
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		return fmt.Errorf("invalid Xitcoin authority address: %w", err)
	}

	maxApprovedValidators, minimumSelfDelegation := g.Policy()
	if err := ValidatePolicy(maxApprovedValidators, minimumSelfDelegation); err != nil {
		return err
	}
	if uint32(len(g.ApprovedValidators)) > maxApprovedValidators {
		return fmt.Errorf("approved validators exceed configured capacity")
	}

	if len(g.ApprovedValidators) == 0 {
		return fmt.Errorf("at least one approved validator is required")
	}

	seen := map[string]struct{}{}
	for _, address := range g.ApprovedValidators {
		trimmed := strings.TrimSpace(address)
		if trimmed == "" {
			return fmt.Errorf("approved validator address is empty")
		}
		if trimmed != address {
			return fmt.Errorf("approved validator address contains surrounding spaces")
		}
		if _, err := sdk.ValAddressFromBech32(address); err != nil {
			return fmt.Errorf("invalid Xitcoin validator address: %w", err)
		}
		if _, exists := seen[address]; exists {
			return fmt.Errorf("approved validator listed twice: %s", address)
		}
		seen[address] = struct{}{}
	}

	return nil
}
