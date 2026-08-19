package types

const (
	// ModuleName is the canonical validator incentive module name.
	ModuleName = "validator_incentives"

	// StoreKey is the module KV-store key.
	StoreKey = ModuleName

	// RouterKey is the module message-routing key.
	RouterKey = ModuleName

	// TreasuryAccountName is the funded module account. It must be registered
	// without mint or burn permissions.
	TreasuryAccountName = ModuleName
)

var (
	// ParamsKey stores the active incentive parameters.
	ParamsKey = []byte{0x01}

	// PeriodStateKey stores the current deterministic distribution period.
	PeriodStateKey = []byte{0x02}

	// TotalDistributedKey stores the cumulative treasury-funded distribution.
	TotalDistributedKey = []byte{0x03}

	// AuthorityKey stores the governance or security authority permitted to
	// update parameters and activate a fully funded reward period.
	AuthorityKey = []byte{0x04}
)
