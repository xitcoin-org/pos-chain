package types

const (
	ModuleName = "bridge"
	StoreKey   = ModuleName
)

var KeyProcessedAttestationPrefix = []byte{0x01}

var (
	KeyOutstandingAmount = []byte{0x07}
	KeyOutboundNonce     = []byte{0x08}
)
