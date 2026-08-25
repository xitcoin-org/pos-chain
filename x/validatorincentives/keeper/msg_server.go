package keeper

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper        Keeper
	stakingKeeper types.StakingKeeper
	treasury      Treasury
}

func NewMsgServerImpl(
	keeper Keeper,
	stakingKeeper types.StakingKeeper,
	treasury Treasury,
) types.MsgServer {
	return msgServer{
		keeper:        keeper,
		stakingKeeper: stakingKeeper,
		treasury:      treasury,
	}
}

func (s msgServer) UpdateParams(
	goCtx context.Context,
	message *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	if message == nil {
		return nil, errors.New("update params message is nil")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	command := types.UpdateParamsCommand{
		Authority: message.Authority,
		Params: types.Params{
			AnnualRateBasisPoints: message.AnnualRateBasisPoints,
			BlocksPerYear:         message.BlocksPerYear,
			RewardPeriodBlocks:    message.RewardPeriodBlocks,
		},
	}
	if err := s.keeper.HandleUpdateParamsCommand(
		ctx,
		command,
	); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}

func (s msgServer) ActivateFundedPeriod(
	goCtx context.Context,
	message *types.MsgActivateFundedPeriod,
) (*types.MsgActivateFundedPeriodResponse, error) {
	if message == nil {
		return nil, errors.New("activate funded period message is nil")
	}
	if _, err := sdk.AccAddressFromBech32(
		message.Authority,
	); err != nil {
		return nil, fmt.Errorf("invalid authority address: %w", err)
	}

	budget, err := types.ParseStoredAtomicAmount(
		message.CommittedAnnualBudgetAtomic,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid committed annual budget: %w",
			err,
		)
	}
	if !budget.IsPositive() {
		return nil, errors.New(
			"committed annual budget must be positive",
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	state, err := s.keeper.ActivatePeriodFromChainState(
		ctx,
		message.Authority,
		budget,
		s.stakingKeeper,
		s.treasury,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgActivateFundedPeriodResponse{
		StartBlock:            state.StartBlock,
		EndBlock:              state.EndBlock,
		PeriodProvisionAtomic: state.PeriodProvisionAtomic,
	}, nil
}
