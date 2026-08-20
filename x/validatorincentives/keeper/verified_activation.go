package keeper

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

// ActivatePeriodFromChainState derives eligible bonded XTC and the treasury
// balance from canonical on-chain keepers. The caller supplies only the
// authority and committed annual budget.
func (k Keeper) ActivatePeriodFromChainState(
	ctx sdk.Context,
	caller string,
	committedAnnualBudget sdkmath.Int,
	stakingKeeper types.StakingKeeper,
	treasury Treasury,
) (types.PeriodState, error) {
	if !committedAnnualBudget.IsPositive() {
		return types.PeriodState{},
			errors.New("committed annual budget must be positive")
	}
	if err := k.RequireAuthority(ctx, caller); err != nil {
		return types.PeriodState{}, err
	}

	eligibleBonded, err := stakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		return types.PeriodState{}, err
	}
	if !eligibleBonded.IsPositive() {
		return types.PeriodState{},
			errors.New("eligible bonded stake must be positive")
	}

	balance, err := treasury.Balance(ctx)
	if err != nil {
		return types.PeriodState{}, err
	}
	if balance.Denom != IncentiveDenom {
		return types.PeriodState{},
			errors.New("unexpected incentive treasury denomination")
	}
	if !balance.Amount.IsPositive() {
		return types.PeriodState{},
			errors.New("incentive treasury must be funded")
	}

	return k.ActivateFundedPeriod(
		ctx,
		eligibleBonded,
		balance.Amount,
		committedAnnualBudget,
	)
}
