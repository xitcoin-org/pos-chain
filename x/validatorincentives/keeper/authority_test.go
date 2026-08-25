package keeper

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

func authorityAddress(fill byte) string {
	return sdk.AccAddress(bytes.Repeat([]byte{fill}, 20)).String()
}

func TestRequireAuthority(t *testing.T) {
	ctx, k := keeperTestContext(t)
	authority := authorityAddress(1)

	require.ErrorContains(
		t,
		k.RequireAuthority(ctx, authority),
		"not configured",
	)

	k.SetAuthority(ctx, authority)
	require.NoError(t, k.RequireAuthority(ctx, authority))
	require.ErrorContains(
		t,
		k.RequireAuthority(ctx, authorityAddress(2)),
		"not authorized",
	)
	require.ErrorContains(
		t,
		k.RequireAuthority(ctx, ""),
		"empty",
	)
}

func TestAuthorizedParameterUpdate(t *testing.T) {
	ctx, k := keeperTestContext(t)
	authority := authorityAddress(1)
	k.SetAuthority(ctx, authority)

	next := types.DefaultParams()
	next.AnnualRateBasisPoints = 900

	require.Error(
		t,
		k.UpdateParamsAuthorized(
			ctx,
			authorityAddress(2),
			next,
		),
	)
	require.Equal(t, types.DefaultParams(), k.GetParams(ctx))

	require.NoError(
		t,
		k.UpdateParamsAuthorized(ctx, authority, next),
	)
	require.Equal(t, next, k.GetParams(ctx))
}

func TestAuthorizedPeriodActivation(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)
	authority := authorityAddress(1)
	k.SetAuthority(ctx, authority)

	_, err := k.ActivateFundedPeriodAuthorized(
		ctx,
		authorityAddress(2),
		keeperTestXTC(2_000_000_000),
		keeperTestXTC(100_000_000),
		keeperTestXTC(160_000_000),
	)
	require.ErrorContains(t, err, "not authorized")

	_, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.False(t, found)

	state, err := k.ActivateFundedPeriodAuthorized(
		ctx,
		authority,
		keeperTestXTC(2_000_000_000),
		keeperTestXTC(100_000_000),
		keeperTestXTC(160_000_000),
	)
	require.NoError(t, err)
	require.Equal(t, uint64(100), state.StartBlock)
}
