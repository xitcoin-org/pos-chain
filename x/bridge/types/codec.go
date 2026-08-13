package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// RegisterInterfaces is intentionally empty until bridge transaction messages
// are independently designed, reviewed and added.
func RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

func RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}
