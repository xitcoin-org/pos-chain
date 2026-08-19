package keeper

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

const IncentiveDenom = "axtc"

// Treasury provides the narrow, non-minting bank boundary used by validator
// incentives. Its module account is registered without mint or burn rights.
type Treasury struct {
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

func NewTreasury(
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) Treasury {
	return Treasury{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

func (t Treasury) Balance(ctx sdk.Context) (sdk.Coin, error) {
	address := t.accountKeeper.GetModuleAddress(
		types.TreasuryAccountName,
	)
	if len(address) == 0 {
		return sdk.Coin{}, errors.New("incentive treasury module account is missing")
	}

	return t.bankKeeper.GetBalance(
		ctx,
		address,
		IncentiveDenom,
	), nil
}

func (t Treasury) Send(
	ctx sdk.Context,
	recipient sdk.AccAddress,
	amount sdkmath.Int,
) error {
	if len(recipient) == 0 {
		return errors.New("reward recipient is empty")
	}
	if !amount.IsPositive() {
		return errors.New("reward amount must be positive")
	}

	balance, err := t.Balance(ctx)
	if err != nil {
		return err
	}
	if balance.Amount.LT(amount) {
		return errors.New("incentive treasury balance is insufficient")
	}

	return t.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.TreasuryAccountName,
		recipient,
		sdk.NewCoins(sdk.NewCoin(IncentiveDenom, amount)),
	)
}

// DistributeFunded validates all period constraints before requesting a bank
// transfer. Accounting is advanced only after the bank transfer succeeds.
// Cosmos transaction caching rolls back both operations if a later error is
// returned before transaction commit.
func (k Keeper) DistributeFunded(
	ctx sdk.Context,
	treasury Treasury,
	recipient sdk.AccAddress,
	amount sdkmath.Int,
) error {
	if err := k.validateDistribution(ctx, amount); err != nil {
		return err
	}
	if err := treasury.Send(ctx, recipient, amount); err != nil {
		return err
	}
	return k.RecordDistribution(ctx, amount)
}

func (k Keeper) validateDistribution(
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
	return nil
}
