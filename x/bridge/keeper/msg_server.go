package keeper

import (
	"context"
	"encoding/hex"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xitcoin-org/pos-chain/x/bridge/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServer creates the testnet-only bridge admission message server.
func NewMsgServer(k Keeper) types.MsgServer {
	return msgServer{keeper: k}
}

// SubmitAttestation records a validated attestation ID only. It has no bank,
// token, minting, reserve, custody, relayer, or settlement capability.
func (s msgServer) SubmitAttestation(goCtx context.Context, req *types.MsgSubmitAttestation) (*types.MsgSubmitAttestationResponse, error) {
	if req == nil {
		return nil, errors.New("bridge: empty attestation submission")
	}
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	config, found, err := s.keeper.GetRouteConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrRouteDisabled
	}

	id, err := s.keeper.AdmitAttestation(ctx, config, req.Attestation(), req.Signatures)
	if err != nil {
		return nil, err
	}
	return &types.MsgSubmitAttestationResponse{AttestationId: hex.EncodeToString(id[:])}, nil
}
