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

const submitAttestationName = "xitcoin/bridge/MsgSubmitAttestation"
const initiateOutboundTransferName = "xitcoin/bridge/MsgInitiateOutboundTransfer"

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}

// RegisterInterfaces registers testnet-only bridge messages.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSubmitAttestation{},
		&MsgInitiateOutboundTransfer{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterLegacyAminoCodec supports legacy sign bytes.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSubmitAttestation{}, submitAttestationName, nil)
	cdc.RegisterConcrete(&MsgInitiateOutboundTransfer{}, initiateOutboundTransferName, nil)
}
