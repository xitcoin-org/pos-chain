package keeper

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
}

func NewKeeper(storeKey storetypes.StoreKey) Keeper {
	return Keeper{storeKey: storeKey}
}

func (k Keeper) SetAuthority(ctx sdk.Context, authority string) {
	ctx.KVStore(k.storeKey).Set(types.AuthorityKey, []byte(authority))
}

func (k Keeper) GetAuthority(ctx sdk.Context) string {
	return string(ctx.KVStore(k.storeKey).Get(types.AuthorityKey))
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	value := make([]byte, 20)
	binary.BigEndian.PutUint32(value[0:4], params.AnnualRateBasisPoints)
	binary.BigEndian.PutUint64(value[4:12], params.BlocksPerYear)
	binary.BigEndian.PutUint64(value[12:20], params.RewardPeriodBlocks)
	ctx.KVStore(k.storeKey).Set(types.ParamsKey, value)
	return nil
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	value := ctx.KVStore(k.storeKey).Get(types.ParamsKey)
	if len(value) != 20 {
		return types.DefaultParams()
	}

	return types.Params{
		AnnualRateBasisPoints: binary.BigEndian.Uint32(value[0:4]),
		BlocksPerYear:         binary.BigEndian.Uint64(value[4:12]),
		RewardPeriodBlocks:    binary.BigEndian.Uint64(value[12:20]),
	}
}

func (k Keeper) UpdateParams(ctx sdk.Context, next types.Params) error {
	previous := k.GetParams(ctx)
	if err := next.ValidateUpdate(previous); err != nil {
		return err
	}
	return k.SetParams(ctx, next)
}

func (k Keeper) SetPeriodState(
	ctx sdk.Context,
	state types.PeriodState,
) error {
	if err := state.Validate(); err != nil {
		return err
	}
	encoded, err := json.Marshal(state)
	if err != nil {
		return err
	}
	ctx.KVStore(k.storeKey).Set(types.PeriodStateKey, encoded)
	return nil
}

func (k Keeper) GetPeriodState(
	ctx sdk.Context,
) (types.PeriodState, bool, error) {
	value := ctx.KVStore(k.storeKey).Get(types.PeriodStateKey)
	if len(value) == 0 {
		return types.PeriodState{}, false, nil
	}

	var state types.PeriodState
	if err := json.Unmarshal(value, &state); err != nil {
		return types.PeriodState{}, false, err
	}
	if err := state.Validate(); err != nil {
		return types.PeriodState{}, false, err
	}
	return state, true, nil
}

func (k Keeper) SetTotalDistributed(
	ctx sdk.Context,
	amount sdkmath.Int,
) error {
	if amount.IsNegative() {
		return errors.New("total distributed amount cannot be negative")
	}
	ctx.KVStore(k.storeKey).Set(
		types.TotalDistributedKey,
		[]byte(amount.String()),
	)
	return nil
}

func (k Keeper) GetTotalDistributed(ctx sdk.Context) (sdkmath.Int, error) {
	value := string(
		ctx.KVStore(k.storeKey).Get(types.TotalDistributedKey),
	)
	if value == "" {
		return sdkmath.ZeroInt(), nil
	}

	amount, err := types.ParseStoredAtomicAmount(value)
	if err != nil {
		return sdkmath.Int{}, err
	}
	if amount.IsNegative() {
		return sdkmath.Int{}, errors.New("stored total distributed amount is negative")
	}
	return amount, nil
}
