// Package keeper stores bridge execution state.
// It has no authority to mint, reserve, or transfer tokens.
package keeper

import (
	"errors"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xitcoin-org/pos-chain/x/bridge/types"
)

var ErrAttestationAlreadyProcessed = errors.New("bridge attestation already processed")

type Keeper struct {
	storeKey storetypes.StoreKey
}

func NewKeeper(storeKey storetypes.StoreKey) Keeper {
	return Keeper{storeKey: storeKey}
}

func (k Keeper) IsProcessed(ctx sdk.Context, id [32]byte) bool {
	return ctx.KVStore(k.storeKey).Has(processedKey(id))
}

func (k Keeper) MarkProcessed(ctx sdk.Context, id [32]byte) {
	ctx.KVStore(k.storeKey).Set(processedKey(id), []byte{1})
}

// ConsumeAttestation atomically reserves a valid attestation identifier in the
// transaction state. A second call for the same ID fails before settlement.
func (k Keeper) ConsumeAttestation(ctx sdk.Context, attestation types.Attestation) ([32]byte, error) {
	id, err := attestation.ID()
	if err != nil {
		return [32]byte{}, err
	}
	if k.IsProcessed(ctx, id) {
		return [32]byte{}, ErrAttestationAlreadyProcessed
	}
	k.MarkProcessed(ctx, id)
	return id, nil
}

func processedKey(id [32]byte) []byte {
	key := make([]byte, len(types.KeyProcessedAttestationPrefix)+len(id))
	copy(key, types.KeyProcessedAttestationPrefix)
	copy(key[len(types.KeyProcessedAttestationPrefix):], id[:])
	return key
}
