package evmd

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// KCALBAdministrativeMultisigAddress is the public 2-of-3 administrative
// authority. Its private keys are never embedded in the application.
const KCALBAdministrativeMultisigAddress = "xtc1e3q4pm23ky0qetnep33j4yezq6c3lc7fcds4je"

// UnsafeEnableGovernanceProposalSubmissionOption exists only for legacy test
// suites. Production startup does not register a CLI flag for this option.
const UnsafeEnableGovernanceProposalSubmissionOption = "unsafe-enable-governance-proposal-submission"

func mustKCALBAdministrativeAuthority() string {
	authority, err := sdk.AccAddressFromBech32(KCALBAdministrativeMultisigAddress)
	if err != nil {
		panic(fmt.Errorf("invalid KCALB administrative multisig: %w", err))
	}
	return authority.String()
}
