package types

import "testing"

func TestDefaultGenesisIsDisabledAndValid(t *testing.T) {
	state := DefaultGenesisState()
	if state.RouteConfig != nil || state.Paused {
		t.Fatal("default bridge genesis must not configure or enable a route")
	}
	if err := state.Validate(); err != nil {
		t.Fatalf("default genesis rejected: %v", err)
	}
}

func TestGenesisRejectsPauseWithoutRoute(t *testing.T) {
	if err := (GenesisState{Paused: true}).Validate(); err == nil {
		t.Fatal("paused empty route accepted")
	}
}
