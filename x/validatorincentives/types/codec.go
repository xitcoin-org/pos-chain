package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino = codec.NewLegacyAmino()

	ModuleCdc = codec.NewProtoCodec(
		codectypes.NewInterfaceRegistry(),
	)

	AminoCdc = codec.NewAminoCodec(amino) //nolint:staticcheck
)

const (
	updateParamsName =
		"xitcoin/validatorincentives/MsgUpdateParams"
	activateFundedPeriodName =
		"xitcoin/validatorincentives/MsgActivateFundedPeriod"
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}

func RegisterInterfaces(
	registry codectypes.InterfaceRegistry,
) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgActivateFundedPeriod{},
	)
	msgservice.RegisterMsgServiceDesc(
		registry,
		&_Msg_serviceDesc,
	)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(
		&MsgUpdateParams{},
		updateParamsName,
		nil,
	)
	cdc.RegisterConcrete(
		&MsgActivateFundedPeriod{},
		activateFundedPeriodName,
		nil,
	)
}
