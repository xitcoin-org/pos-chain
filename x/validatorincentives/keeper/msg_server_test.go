package keeper

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

func testMsgServer(
	t *testing.T,
) (context.Context, Keeper, types.MsgServer) {
	t.Helper()

	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)
	authority := authorityAddress(1)
	k.SetAuthority(ctx, authority)

	treasury, _ := fundedTreasury(
		keeperTestXTC(100_000_000),
	)
	staking := activationStakingKeeper{
		bonded: keeperTestXTC(2_000_000_000),
	}
	server := NewMsgServerImpl(k, staking, treasury)
	return sdk.WrapSDKContext(ctx), k, server
}

func TestMsgServerUpdateParams(t *testing.T) {
	goCtx, k, server := testMsgServer(t)
	ctx := sdk.UnwrapSDKContext(goCtx)
	authority := k.GetAuthority(ctx)

	_, err := server.UpdateParams(
		goCtx,
		&types.MsgUpdateParams{
			Authority:             authority,
			AnnualRateBasisPoints: 900,
			BlocksPerYear:         types.DefaultBlocksPerYear,
			RewardPeriodBlocks:    types.DefaultBlocksPerYear / 4,
		},
	)
	require.NoError(t, err)
	require.Equal(
		t,
		uint32(900),
		k.GetParams(ctx).AnnualRateBasisPoints,
	)
}

func TestMsgServerActivateFundedPeriod(t *testing.T) {
	goCtx, k, server := testMsgServer(t)
	ctx := sdk.UnwrapSDKContext(goCtx)
	authority := k.GetAuthority(ctx)

	response, err := server.ActivateFundedPeriod(
		goCtx,
		&types.MsgActivateFundedPeriod{
			Authority: authority,
			CommittedAnnualBudgetAtomic: keeperTestXTC(
				160_000_000,
			).String(),
		},
	)
	require.NoError(t, err)
	require.Equal(t, uint64(100), response.StartBlock)
	require.Greater(t, response.EndBlock, response.StartBlock)
	require.NotEmpty(t, response.PeriodProvisionAtomic)
}

func TestMsgServerRejectsInvalidRequests(t *testing.T) {
	goCtx, _, server := testMsgServer(t)

	_, err := server.UpdateParams(goCtx, nil)
	require.ErrorContains(t, err, "nil")

	_, err = server.ActivateFundedPeriod(goCtx, nil)
	require.ErrorContains(t, err, "nil")

	_, err = server.ActivateFundedPeriod(
		goCtx,
		&types.MsgActivateFundedPeriod{
			Authority:                   "invalid",
			CommittedAnnualBudgetAtomic: "1",
		},
	)
	require.ErrorContains(t, err, "invalid authority")

	_, err = server.ActivateFundedPeriod(
		goCtx,
		&types.MsgActivateFundedPeriod{
			Authority:                   authorityAddress(1),
			CommittedAnnualBudgetAtomic: "0",
		},
	)
	require.ErrorContains(t, err, "must be positive")
}
