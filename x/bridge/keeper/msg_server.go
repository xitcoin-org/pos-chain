package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xitcoin-org/pos-chain/x/bridge/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServer creates the bridge settlement message server.
func NewMsgServer(k Keeper) types.MsgServer {
	return msgServer{keeper: k}
}

// SubmitAttestation validates and settles one inbound Cronos transfer.
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

	attestation := req.Attestation()
	if attestation.Direction != types.DirectionCronosToXitcoin {
		return nil, ErrSettlementDirection
	}
	id, err := s.keeper.AdmitAttestation(ctx, config, attestation, req.Signatures)
	if err != nil {
		return nil, err
	}
	if err := s.keeper.SettleInbound(ctx, config, attestation); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"bridge_inbound_settled",
		sdk.NewAttribute("attestation_id", hex.EncodeToString(id[:])),
		sdk.NewAttribute("route_id", config.RouteID),
		sdk.NewAttribute("destination", attestation.Destination),
		sdk.NewAttribute("amount", attestation.Amount),
	))
	return &types.MsgSubmitAttestationResponse{AttestationId: hex.EncodeToString(id[:])}, nil
}

func (s msgServer) InitiateOutboundTransfer(goCtx context.Context, req *types.MsgInitiateOutboundTransfer) (*types.MsgInitiateOutboundTransferResponse, error) {
	if req == nil {
		return nil, errors.New("bridge: empty outbound transfer")
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
	if req.RouteId != config.RouteID {
		return nil, ErrAttestationRouteMismatch
	}
	sender, err := sdk.AccAddressFromBech32(req.Sender)
	if err != nil {
		return nil, err
	}
	id, nonce, err := s.keeper.InitiateOutboundTransfer(ctx, config, sender, req.Destination, req.Amount)
	if err != nil {
		return nil, err
	}
	requestID := hex.EncodeToString(id[:])
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"bridge_outbound_burned",
		sdk.NewAttribute("request_id", requestID),
		sdk.NewAttribute("route_id", config.RouteID),
		sdk.NewAttribute("sender", sender.String()),
		sdk.NewAttribute("destination", req.Destination),
		sdk.NewAttribute("amount", req.Amount),
		sdk.NewAttribute("nonce", fmt.Sprintf("%d", nonce)),
	))
	return &types.MsgInitiateOutboundTransferResponse{RequestId: requestID, Nonce: nonce}, nil
}
