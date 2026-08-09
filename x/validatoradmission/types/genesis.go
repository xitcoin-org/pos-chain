package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Authority          string   `json:"authority"`
	ApprovedValidators []string `json:"approved_validators"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Authority:          "",
		ApprovedValidators: []string{},
	}
}

func (g GenesisState) Validate() error {
	authority := strings.TrimSpace(g.Authority)
	if authority == "" {
		return fmt.Errorf("validator admission authority is required")
	}
	if authority != g.Authority {
		return fmt.Errorf("validator admission authority contains surrounding spaces")
	}
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		return fmt.Errorf("invalid Xitcoin authority address: %w", err)
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
