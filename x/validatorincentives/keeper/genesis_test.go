package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

func TestKeeperGenesisRoundTrip(t *testing.T) {
	ctx, k := keeperTestContext(t)
	state := types.DefaultGenesisState()
	state.Authority = authorityAddress(1)
	state.Params.AnnualRateBasisPoints = 900

	require.NoError(t, k.InitGenesis(ctx, state))
	require.Equal(t, state, k.ExportGenesis(ctx))

	total, err := k.GetTotalDistributed(ctx)
	require.NoError(t, err)
	require.True(t, total.IsZero())

	_, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.False(t, found)
}

func TestKeeperGenesisRejectsInvalidState(t *testing.T) {
	ctx, k := keeperTestContext(t)
	state := types.DefaultGenesisState()
	state.Params.RewardPeriodBlocks = 0

	require.Error(t, k.InitGenesis(ctx, state))
	require.Equal(t, "", k.GetAuthority(ctx))
	require.Equal(t, types.DefaultParams(), k.GetParams(ctx))
}
