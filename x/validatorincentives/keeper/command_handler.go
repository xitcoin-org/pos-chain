package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

// HandleUpdateParamsCommand validates the domain command before delegating to
// the authority-protected keeper operation.
func (k Keeper) HandleUpdateParamsCommand(
	ctx sdk.Context,
	command types.UpdateParamsCommand,
) error {
	if err := command.ValidateBasic(); err != nil {
		return err
	}
	return k.UpdateParamsAuthorized(
		ctx,
		command.Authority,
		command.Params,
	)
}

// HandleActivatePeriodCommand validates and parses canonical atomic amounts
// before delegating to the authority-protected period activation.
func (k Keeper) HandleActivatePeriodCommand(
	ctx sdk.Context,
	command types.ActivatePeriodCommand,
) (types.PeriodState, error) {
	if err := command.ValidateBasic(); err != nil {
		return types.PeriodState{}, err
	}

	eligible, treasury, budget, err := command.Amounts()
	if err != nil {
		return types.PeriodState{}, err
	}
	return k.ActivateFundedPeriodAuthorized(
		ctx,
		command.Authority,
		eligible,
		treasury,
		budget,
	)
}
