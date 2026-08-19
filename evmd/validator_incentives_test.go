package evmd

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	evmconfig "github.com/xitcoin-org/pos-chain/evmd/config"
	validatorincentivestypes "github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

func TestValidatorIncentivesWiring(t *testing.T) {
	app, genesis := setup(
		true,
		0,
		"xitcoin-validator-incentives-test",
		101089,
	)

	raw, found := genesis[validatorincentivestypes.ModuleName]
	require.True(t, found)
	require.NotEmpty(t, raw)

	var state validatorincentivestypes.GenesisState
	app.AppCodec().MustUnmarshalJSON(raw, &state)
	require.NoError(t, state.Validate())
	require.Empty(t, state.Authority)
	require.Equal(
		t,
		validatorincentivestypes.DefaultParams(),
		state.Params,
	)

	require.NotNil(
		t,
		app.GetKey(validatorincentivestypes.StoreKey),
	)

	permissions := evmconfig.GetMaccPerms()
	modulePermissions, found := permissions[
		validatorincentivestypes.TreasuryAccountName
	]
	require.True(t, found)
	require.Empty(t, modulePermissions)

	treasuryAddress := authtypes.NewModuleAddress(
		validatorincentivestypes.TreasuryAccountName,
	)
	require.NotEmpty(t, treasuryAddress)
	require.True(
		t,
		evmconfig.BlockedAddresses()[
			sdk.AccAddress(treasuryAddress).String()
		],
	)
}
