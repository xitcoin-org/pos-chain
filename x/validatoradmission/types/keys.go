package types

const (
	ModuleName = "validator_admission"
	StoreKey   = ModuleName
)

var (
	KeyAuthority               = []byte{0x01}
	KeyApprovedValidatorPrefix = []byte{0x02}
	KeyMaxApprovedValidators   = []byte{0x03}
	KeyMinimumSelfDelegation   = []byte{0x04}
)
