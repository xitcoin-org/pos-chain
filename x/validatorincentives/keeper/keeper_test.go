package keeper

import (
	"bytes"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

func keeperTestContext(t *testing.T) (sdk.Context, Keeper) {
	t.Helper()

	key := storetypes.NewKVStoreKey(types.StoreKey)
	ctx := sdktestutil.DefaultContext(
		key,
		storetypes.NewTransientStoreKey("validator_incentives_test"),
	)
	return ctx, NewKeeper(key)
}

func keeperTestXTC(value int64) sdkmath.Int {
	return sdkmath.NewInt(value).MulRaw(1_000_000_000_000_000_000)
}

func TestKeeperParamsAndAuthority(t *testing.T) {
	ctx, k := keeperTestContext(t)

	require.Equal(t, types.DefaultParams(), k.GetParams(ctx))

	authority := sdk.AccAddress(bytes.Repeat([]byte{1}, 20)).String()
	k.SetAuthority(ctx, authority)
	require.Equal(t, authority, k.GetAuthority(ctx))

	next := types.DefaultParams()
	next.AnnualRateBasisPoints = 900
	require.NoError(t, k.UpdateParams(ctx, next))
	require.Equal(t, next, k.GetParams(ctx))

	invalid := next
	invalid.AnnualRateBasisPoints = 1_001
	require.Error(t, k.UpdateParams(ctx, invalid))
	require.Equal(t, next, k.GetParams(ctx))
}

func TestKeeperPeriodState(t *testing.T) {
	ctx, k := keeperTestContext(t)

	state, err := types.NewPeriodState(
		100,
		keeperTestXTC(2_000_000_000),
		keeperTestXTC(100_000_000),
		keeperTestXTC(160_000_000),
		types.DefaultParams(),
	)
	require.NoError(t, err)
	require.NoError(t, k.SetPeriodState(ctx, state))

	stored, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, state, stored)
}

func TestKeeperTotalDistributed(t *testing.T) {
	ctx, k := keeperTestContext(t)

	total, err := k.GetTotalDistributed(ctx)
	require.NoError(t, err)
	require.True(t, total.IsZero())

	require.NoError(t, k.SetTotalDistributed(ctx, keeperTestXTC(12_345)))
	total, err = k.GetTotalDistributed(ctx)
	require.NoError(t, err)
	require.True(t, total.Equal(keeperTestXTC(12_345)))

	require.Error(t, k.SetTotalDistributed(ctx, sdkmath.NewInt(-1)))
}
