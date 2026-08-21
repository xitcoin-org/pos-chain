package keeper

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

// ActivateFundedPeriod records a deterministic reward period whose provision
// is bounded by eligible stake, the committed budget and available treasury
// funds. This method records accounting state only; it never mints tokens.
func (k Keeper) ActivateFundedPeriod(
	ctx sdk.Context,
	eligibleBonded sdkmath.Int,
	treasuryBalance sdkmath.Int,
	committedAnnualBudget sdkmath.Int,
) (types.PeriodState, error) {
	if ctx.BlockHeight() < 0 {
		return types.PeriodState{}, errors.New("block height cannot be negative")
	}

	current, found, err := k.GetPeriodState(ctx)
	if err != nil {
		return types.PeriodState{}, err
	}
	if found && uint64(ctx.BlockHeight()) < current.EndBlock {
		return types.PeriodState{}, errors.New("an incentive period is already active")
	}

	state, err := types.NewPeriodState(
		uint64(ctx.BlockHeight()),
		eligibleBonded,
		treasuryBalance,
		committedAnnualBudget,
		k.GetParams(ctx),
	)
	if err != nil {
		return types.PeriodState{}, err
	}
	if err := k.SetPeriodState(ctx, state); err != nil {
		return types.PeriodState{}, err
	}

	return state, nil
}

// RecordDistribution advances the funded-period and lifetime accounting after
// a separate bank transfer has succeeded. It cannot exceed the prefunded
// provision and cannot be used outside the active block interval.
func (k Keeper) RecordDistribution(
	ctx sdk.Context,
	amount sdkmath.Int,
) error {
	if !amount.IsPositive() {
		return errors.New("distribution amount must be positive")
	}
	if ctx.BlockHeight() < 0 {
		return errors.New("block height cannot be negative")
	}

	state, found, err := k.GetPeriodState(ctx)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("no incentive period has been configured")
	}

	height := uint64(ctx.BlockHeight())
	if height < state.StartBlock || height >= state.EndBlock {
		return errors.New("distribution attempted outside the active period")
	}

	remaining, err := state.RemainingProvision()
	if err != nil {
		return err
	}
	if amount.GT(remaining) {
		return errors.New("distribution exceeds the remaining period provision")
	}

	distributed, err := types.ParseStoredAtomicAmount(state.DistributedAtomic)
	if err != nil {
		return err
	}
	total, err := k.GetTotalDistributed(ctx)
	if err != nil {
		return err
	}

	state.DistributedAtomic = distributed.Add(amount).String()
	if err := k.SetPeriodState(ctx, state); err != nil {
		return err
	}
	return k.SetTotalDistributed(ctx, total.Add(amount))
}
