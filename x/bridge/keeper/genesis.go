package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xitcoin-org/pos-chain/x/bridge/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	if state.RouteConfig != nil {
		if err := k.SetRouteConfig(ctx, *state.RouteConfig); err != nil {
			panic(err)
		}
		if state.Paused {
			encoded, err := json.Marshal(RouteState{Paused: true})
			if err != nil {
				panic(err)
			}
			ctx.KVStore(k.storeKey).Set(routeStateKey, encoded)
		}
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	config, found, err := k.GetRouteConfig(ctx)
	if err != nil {
		panic(err)
	}
	if !found {
		return types.DefaultGenesisState()
	}
	state, err := k.GetRouteState(ctx)
	if err != nil {
		panic(err)
	}
	return types.GenesisState{RouteConfig: &config, Paused: state.Paused}
}
