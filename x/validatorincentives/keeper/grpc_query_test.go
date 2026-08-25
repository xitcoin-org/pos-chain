package keeper

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

func TestQueryServerParams(t *testing.T) {
	ctx, k := keeperTestContext(t)
	treasury, _ := fundedTreasury(keeperTestXTC(40_000_000))
	server := NewQueryServer(k, treasury)

	authority := sdk.AccAddress(
		bytes.Repeat([]byte{2}, 20),
	).String()
	k.SetAuthority(ctx, authority)
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	response, err := server.Params(
		sdk.WrapSDKContext(ctx),
		&types.QueryParamsRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, authority, response.Authority)
	require.Equal(
		t,
		types.DefaultAnnualRateBasisPoints,
		response.AnnualRateBasisPoints,
	)
	require.Equal(
		t,
		types.DefaultBlocksPerYear,
		response.BlocksPerYear,
	)
	require.Equal(
		t,
		types.DefaultBlocksPerYear/4,
		response.RewardPeriodBlocks,
	)
}

func TestQueryServerPeriodAndTreasury(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)

	treasuryBalance := keeperTestXTC(40_000_000)
	treasury, _ := fundedTreasury(treasuryBalance)
	server := NewQueryServer(k, treasury)

	empty, err := server.Period(
		sdk.WrapSDKContext(ctx),
		&types.QueryPeriodRequest{},
	)
	require.NoError(t, err)
	require.False(t, empty.Active)
	require.Equal(t, "0", empty.TotalDistributedAtomic)

	state := activeTreasuryPeriod(t, ctx, k)
	distributed := keeperTestXTC(1_000)
	require.NoError(t, k.SetTotalDistributed(ctx, distributed))

	period, err := server.Period(
		sdk.WrapSDKContext(ctx),
		&types.QueryPeriodRequest{},
	)
	require.NoError(t, err)
	require.True(t, period.Active)
	require.Equal(t, state.StartBlock, period.StartBlock)
	require.Equal(t, state.EndBlock, period.EndBlock)
	require.Equal(
		t,
		state.PeriodProvisionAtomic,
		period.PeriodProvisionAtomic,
	)
	require.Equal(
		t,
		distributed.String(),
		period.TotalDistributedAtomic,
	)

	treasuryResponse, err := server.Treasury(
		sdk.WrapSDKContext(ctx),
		&types.QueryTreasuryRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, IncentiveDenom, treasuryResponse.Denom)
	require.Equal(
		t,
		treasuryBalance.String(),
		treasuryResponse.BalanceAtomic,
	)
}
