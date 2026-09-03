package evmd

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestnetAdministrativeMultisigAddress is the public 2-of-3 administrative
// authority. Its private keys are never embedded in the application.
const TestnetAdministrativeMultisigAddress = "xtc1vza8zsgvrfwmve084ytd8xqdkkm7u9e5csctc2"

// UnsafeEnableGovernanceProposalSubmissionOption exists only for legacy test
// suites. Production startup does not register a CLI flag for this option.
const UnsafeEnableGovernanceProposalSubmissionOption = "unsafe-enable-governance-proposal-submission"

func mustTestnetAdministrativeAuthority() string {
	authority, err := sdk.AccAddressFromBech32(TestnetAdministrativeMultisigAddress)
	if err != nil {
		panic(fmt.Errorf("invalid Testnet V2 administrative multisig: %w", err))
	}
	return authority.String()
}
