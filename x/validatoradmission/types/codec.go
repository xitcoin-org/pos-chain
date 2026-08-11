package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc is reserved for tests and JSON encoding.
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	// AminoCdc supports legacy sign bytes.
	AminoCdc = codec.NewAminoCodec(amino) //nolint:staticcheck
)

const (
	approveValidatorName = "xitcoin/validatoradmission/MsgApproveValidator"
	revokeValidatorName  = "xitcoin/validatoradmission/MsgRevokeValidator"
	updateParamsName     = "xitcoin/validatoradmission/MsgUpdateParams"
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}

// RegisterInterfaces registers Xitcoin Validator Admission messages.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgApproveValidator{},
		&MsgRevokeValidator{},
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterLegacyAminoCodec supports legacy sign bytes.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgApproveValidator{}, approveValidatorName, nil)
	cdc.RegisterConcrete(&MsgRevokeValidator{}, revokeValidatorName, nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, updateParamsName, nil)
}
