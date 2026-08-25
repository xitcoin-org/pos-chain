package types

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisStateValidate(t *testing.T) {
	authority := sdk.AccAddress(bytes.Repeat([]byte{1}, 20)).String()
	validator := sdk.ValAddress(bytes.Repeat([]byte{2}, 20)).String()

	valid := GenesisState{
		Authority:          authority,
		ApprovedValidators: []string{validator},
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid state rejected: %v", err)
	}

	cases := []GenesisState{
		{ApprovedValidators: []string{validator}},
		{Authority: authority},
		{Authority: "invalid-address", ApprovedValidators: []string{validator}},
		{Authority: authority, ApprovedValidators: []string{"invalid-validator"}},
		{Authority: authority, ApprovedValidators: []string{validator, validator}},
	}

	for _, state := range cases {
		if err := state.Validate(); err == nil {
			t.Fatalf("invalid state accepted: %#v", state)
		}
	}
}

func TestGenesisStateAllowsLocalTestDenomination(t *testing.T) {
	authority := sdk.AccAddress(bytes.Repeat([]byte{3}, 20)).String()
	validator := sdk.ValAddress(bytes.Repeat([]byte{4}, 20)).String()
	state := GenesisState{
		Authority:             authority,
		ApprovedValidators:    []string{validator},
		MaxApprovedValidators: 1,
		MinimumSelfDelegation: "1000000000000000000atest",
	}
	if err := state.Validate(); err != nil {
		t.Fatalf("local genesis denomination rejected: %v", err)
	}
	if err := ValidatePolicy(state.MaxApprovedValidators, state.MinimumSelfDelegation); err == nil {
		t.Fatal("runtime policy must remain restricted to axtc")
	}
}
