package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Authority string `json:"authority"`
	Params    Params `json:"params"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Authority: "",
		Params:    DefaultParams(),
	}
}

func (g GenesisState) Validate() error {
	if err := g.Params.Validate(); err != nil {
		return fmt.Errorf("invalid genesis parameters: %w", err)
	}

	authority := strings.TrimSpace(g.Authority)
	if authority == "" {
		return nil
	}
	if authority != g.Authority {
		return fmt.Errorf("incentive authority contains surrounding spaces")
	}
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		return fmt.Errorf("invalid incentive authority address: %w", err)
	}
	return nil
}
