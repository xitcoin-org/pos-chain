package types

import "testing"

func TestDefaultGenesisStateIsDisabledAndValid(t *testing.T) {
	if err := DefaultGenesisState().Validate(); err != nil {
		t.Fatalf("default disabled genesis must be valid: %v", err)
	}
}

func TestDisabledGenesisRejectsConfiguredValidatorsWithoutAuthority(t *testing.T) {
	state := DefaultGenesisState()
	state.ApprovedValidators = []string{"xitcoinvaloper1placeholder"}
	if err := state.Validate(); err == nil {
		t.Fatal("configured validators must require an authority")
	}
}
