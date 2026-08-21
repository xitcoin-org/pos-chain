package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisStateValidation(t *testing.T) {
	state := DefaultGenesisState()
	require.NoError(t, state.Validate())
	require.Equal(t, DefaultParams(), state.Params)

	state.Authority = commandAuthority()
	require.NoError(t, state.Validate())

	state.Authority = " " + state.Authority
	require.ErrorContains(t, state.Validate(), "surrounding spaces")

	state = DefaultGenesisState()
	state.Params.AnnualRateBasisPoints =
		MaxAnnualRateBasisPoints + 1
	require.ErrorContains(t, state.Validate(), "protocol ceiling")
}
