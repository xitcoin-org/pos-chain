package types

const (
	ModuleName = "validator_admission"
	StoreKey   = ModuleName
)

var (
	KeyAuthority               = []byte{0x01}
	KeyApprovedValidatorPrefix = []byte{0x02}
)
