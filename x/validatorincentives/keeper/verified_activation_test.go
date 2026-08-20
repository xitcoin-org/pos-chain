package keeper

import (
	"context"
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

type activationStakingKeeper struct {
	bonded sdkmath.Int
	err    error
}

func (k activationStakingKeeper) TotalBondedTokens(
	context.Context,
) (sdkmath.Int, error) {
	return k.bonded, k.err
}

func TestActivatePeriodReadsCanonicalChainState(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)
	authority := authorityAddress(1)
	k.SetAuthority(ctx, authority)

	treasury, _ := fundedTreasury(keeperTestXTC(100_000_000))
	staking := activationStakingKeeper{
		bonded: keeperTestXTC(2_000_000_000),
	}

	state, err := k.ActivatePeriodFromChainState(
		ctx,
		authority,
		keeperTestXTC(160_000_000),
		staking,
		treasury,
	)
	require.NoError(t, err)
	require.Equal(
		t,
		staking.bonded.String(),
		state.EligibleBondedAtomic,
	)
	require.Equal(
		t,
		keeperTestXTC(160_000_000).String(),
		state.CommittedAnnualBudgetAtomic,
	)
}

func TestActivatePeriodRejectsUnfundedTreasury(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)
	authority := authorityAddress(1)
	k.SetAuthority(ctx, authority)

	treasury, _ := fundedTreasury(sdkmath.ZeroInt())
	staking := activationStakingKeeper{
		bonded: keeperTestXTC(20_000_000),
	}

	_, err := k.ActivatePeriodFromChainState(
		ctx,
		authority,
		keeperTestXTC(1_000_000),
		staking,
		treasury,
	)
	require.ErrorContains(t, err, "must be funded")

	_, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.False(t, found)
}

func TestActivatePeriodRejectsStakingQueryFailure(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)
	authority := authorityAddress(1)
	k.SetAuthority(ctx, authority)

	treasury, _ := fundedTreasury(keeperTestXTC(100_000_000))
	staking := activationStakingKeeper{
		bonded: sdkmath.ZeroInt(),
		err:    errors.New("staking query failed"),
	}

	_, err := k.ActivatePeriodFromChainState(
		ctx,
		authority,
		keeperTestXTC(1_000_000),
		staking,
		treasury,
	)
	require.ErrorContains(t, err, "staking query failed")

	_, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.False(t, found)
}
