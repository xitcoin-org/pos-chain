package keeper

import (
	"context"
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

type treasuryAccountKeeper struct {
	address sdk.AccAddress
}

func (k treasuryAccountKeeper) GetModuleAddress(string) sdk.AccAddress {
	return k.address
}

type treasuryBankKeeper struct {
	balance sdk.Coin
	sendErr error
	sent    sdk.Coins
}

func (k *treasuryBankKeeper) GetBalance(
	context.Context,
	sdk.AccAddress,
	string,
) sdk.Coin {
	return k.balance
}

func (k *treasuryBankKeeper) SendCoinsFromModuleToAccount(
	_ context.Context,
	_ string,
	_ sdk.AccAddress,
	amount sdk.Coins,
) error {
	if k.sendErr != nil {
		return k.sendErr
	}
	k.sent = amount
	return nil
}

func fundedTreasury(balance sdkmath.Int) (Treasury, *treasuryBankKeeper) {
	bank := &treasuryBankKeeper{
		balance: sdk.NewCoin(IncentiveDenom, balance),
	}
	account := treasuryAccountKeeper{
		address: sdk.AccAddress(make([]byte, 20)),
	}
	return NewTreasury(account, bank), bank
}

func activeTreasuryPeriod(
	t *testing.T,
	ctx sdk.Context,
	k Keeper,
) types.PeriodState {
	t.Helper()

	state, err := k.ActivateFundedPeriod(
		ctx,
		keeperTestXTC(2_000_000_000),
		keeperTestXTC(100_000_000),
		keeperTestXTC(160_000_000),
	)
	require.NoError(t, err)
	return state
}

func TestDistributeFundedTransfersThenAccounts(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)
	state := activeTreasuryPeriod(t, ctx, k)

	provision, err := types.ParseStoredAtomicAmount(
		state.PeriodProvisionAtomic,
	)
	require.NoError(t, err)
	amount := provision.QuoRaw(4)

	treasury, bank := fundedTreasury(provision)
	recipient := sdk.AccAddress(make([]byte, 20))
	recipient[0] = 1

	require.NoError(
		t,
		k.DistributeFunded(ctx, treasury, recipient, amount),
	)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewCoin(IncentiveDenom, amount)),
		bank.sent,
	)

	stored, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, amount.String(), stored.DistributedAtomic)
}

func TestDistributeFundedBankFailureDoesNotAccount(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)
	state := activeTreasuryPeriod(t, ctx, k)

	provision, err := types.ParseStoredAtomicAmount(
		state.PeriodProvisionAtomic,
	)
	require.NoError(t, err)
	amount := provision.QuoRaw(4)

	treasury, bank := fundedTreasury(provision)
	bank.sendErr = errors.New("bank transfer rejected")
	recipient := sdk.AccAddress(make([]byte, 20))
	recipient[0] = 1

	require.ErrorContains(
		t,
		k.DistributeFunded(ctx, treasury, recipient, amount),
		"bank transfer rejected",
	)

	stored, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, "0", stored.DistributedAtomic)

	total, err := k.GetTotalDistributed(ctx)
	require.NoError(t, err)
	require.True(t, total.IsZero())
}

func TestDistributeFundedRejectsInsufficientTreasury(t *testing.T) {
	ctx, k := keeperTestContext(t)
	ctx = ctx.WithBlockHeight(100)
	state := activeTreasuryPeriod(t, ctx, k)

	provision, err := types.ParseStoredAtomicAmount(
		state.PeriodProvisionAtomic,
	)
	require.NoError(t, err)
	amount := provision.QuoRaw(2)

	treasury, bank := fundedTreasury(amount.SubRaw(1))
	recipient := sdk.AccAddress(make([]byte, 20))
	recipient[0] = 1

	require.ErrorContains(
		t,
		k.DistributeFunded(ctx, treasury, recipient, amount),
		"insufficient",
	)
	require.Empty(t, bank.sent)

	stored, found, err := k.GetPeriodState(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, "0", stored.DistributedAtomic)
}
