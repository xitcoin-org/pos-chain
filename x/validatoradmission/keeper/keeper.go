package keeper

import (
	"bytes"
	"sort"
	"strings"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xitcoin-org/pos-chain/x/validatoradmission/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
}

func NewKeeper(storeKey storetypes.StoreKey) Keeper {
	return Keeper{storeKey: storeKey}
}

func (k Keeper) SetAuthority(ctx sdk.Context, authority string) {
	ctx.KVStore(k.storeKey).Set(types.KeyAuthority, []byte(authority))
}

func (k Keeper) GetAuthority(ctx sdk.Context) string {
	return string(ctx.KVStore(k.storeKey).Get(types.KeyAuthority))
}

func (k Keeper) SetApprovedValidator(ctx sdk.Context, validatorAddress string, approved bool) {
	store := ctx.KVStore(k.storeKey)
	key := append(append([]byte{}, types.KeyApprovedValidatorPrefix...), []byte(validatorAddress)...)

	if approved {
		store.Set(key, []byte{1})
		return
	}

	store.Delete(key)
}

func (k Keeper) IsApprovedValidator(ctx sdk.Context, validatorAddress string) bool {
	key := append(append([]byte{}, types.KeyApprovedValidatorPrefix...), []byte(strings.TrimSpace(validatorAddress))...)
	return ctx.KVStore(k.storeKey).Has(key)
}

func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	k.SetAuthority(ctx, state.Authority)
	for _, validatorAddress := range state.ApprovedValidators {
		k.SetApprovedValidator(ctx, validatorAddress, true)
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.KeyApprovedValidatorPrefix, nil)
	defer iterator.Close()

	approvedValidators := []string{}
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		if !bytes.HasPrefix(key, types.KeyApprovedValidatorPrefix) {
			break
		}
		approvedValidators = append(
			approvedValidators,
			string(key[len(types.KeyApprovedValidatorPrefix):]),
		)
	}
	sort.Strings(approvedValidators)

	return types.GenesisState{
		Authority:          k.GetAuthority(ctx),
		ApprovedValidators: approvedValidators,
	}
}
