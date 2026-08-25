package keeper

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

func TestActivateFundedPeriodRejectsOverlap(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)

	state, err := k.ActivateFundedPeriod(
		ctx,
		keeperTestXTC(2_000_000_000),
		keeperTestXTC(100_000_000),
		keeperTestXTC(160_000_000),
	)
	require.NoError(t, err)
	require.Equal(t, uint64(100), state.StartBlock)
	require.NotEqual(t, "0", state.PeriodProvisionAtomic)

	_, err = k.ActivateFundedPeriod(
		ctx,
		keeperTestXTC(2_000_000_000),
		keeperTestXTC(100_000_000),
		keeperTestXTC(160_000_000),
	)
	require.ErrorContains(t, err, "already active")
}

func TestRecordDistributionIsBoundedAndAuditable(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)

	state, err := k.ActivateFundedPeriod(
		ctx,
		keeperTestXTC(2_000_000_000),
		keeperTestXTC(100_000_000),
		keeperTestXTC(160_000_000),
	)
	require.NoError(t, err)

	provision, err := types.ParseStoredAtomicAmount(
		state.PeriodProvisionAtomic,
	)
	require.NoError(t, err)
	first := provision.QuoRaw(2)
	require.True(t, first.IsPositive())

	require.NoError(t, k.RecordDistribution(ctx, first))

	stored, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, first.String(), stored.DistributedAtomic)

	total, err := k.GetTotalDistributed(ctx)
	require.NoError(t, err)
	require.True(t, total.Equal(first))

	require.ErrorContains(
		t,
		k.RecordDistribution(ctx, provision),
		"exceeds",
	)
	require.Error(
		t,
		k.RecordDistribution(ctx, sdkmath.ZeroInt()),
	)
}

func TestRecordDistributionRequiresActivePeriod(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)

	require.ErrorContains(
		t,
		k.RecordDistribution(ctx, keeperTestXTC(1)),
		"no incentive period",
	)

	state, err := k.ActivateFundedPeriod(
		ctx,
		keeperTestXTC(1_000_000),
		keeperTestXTC(100_000),
		keeperTestXTC(80_000),
	)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(state.EndBlock))
	require.ErrorContains(
		t,
		k.RecordDistribution(ctx, keeperTestXTC(1)),
		"outside the active period",
	)
}
