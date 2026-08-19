package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

func TestHandleUpdateParamsCommand(t *testing.T) {
	ctx, k := keeperTestContext(t)
	authority := authorityAddress(1)
	k.SetAuthority(ctx, authority)

	next := types.DefaultParams()
	next.AnnualRateBasisPoints = 900
	command := types.UpdateParamsCommand{
		Authority: authority,
		Params:    next,
	}

	require.NoError(
		t,
		k.HandleUpdateParamsCommand(ctx, command),
	)
	require.Equal(t, next, k.GetParams(ctx))

	unauthorized := command
	unauthorized.Authority = authorityAddress(2)
	unauthorized.Params.AnnualRateBasisPoints = 1_000

	require.ErrorContains(
		t,
		k.HandleUpdateParamsCommand(ctx, unauthorized),
		"not authorized",
	)
	require.Equal(t, next, k.GetParams(ctx))
}

func TestHandleActivatePeriodCommand(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)
	authority := authorityAddress(1)
	k.SetAuthority(ctx, authority)

	command := types.ActivatePeriodCommand{
		Authority:                   authority,
		EligibleBondedAtomic:         keeperTestXTC(2_000_000_000).String(),
		TreasuryBalanceAtomic:        keeperTestXTC(100_000_000).String(),
		CommittedAnnualBudgetAtomic: keeperTestXTC(160_000_000).String(),
	}

	state, err := k.HandleActivatePeriodCommand(ctx, command)
	require.NoError(t, err)
	require.Equal(t, uint64(100), state.StartBlock)

	stored, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, state, stored)
}

func TestInvalidCommandCannotMutateState(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)
	authority := authorityAddress(1)
	k.SetAuthority(ctx, authority)

	command := types.ActivatePeriodCommand{
		Authority:                   authority,
		EligibleBondedAtomic:         "not-an-integer",
		TreasuryBalanceAtomic:        keeperTestXTC(100_000_000).String(),
		CommittedAnnualBudgetAtomic: keeperTestXTC(160_000_000).String(),
	}

	_, err := k.HandleActivatePeriodCommand(ctx, command)
	require.ErrorContains(t, err, "invalid eligible bonded amount")

	_, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.False(t, found)
}
