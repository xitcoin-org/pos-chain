package types

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	sdkmath "cosmossdk.io/math"
)

// PeriodState is the deterministic, auditable snapshot used for one funded
// reward period. Amounts are stored as base-10 atomic axtc strings.
type PeriodState struct {
	StartBlock                  uint64 `json:"start_block"`
	EndBlock                    uint64 `json:"end_block"`
	AnnualRateBasisPoints       uint32 `json:"annual_rate_basis_points"`
	EligibleBondedAtomic        string `json:"eligible_bonded_atomic"`
	CommittedAnnualBudgetAtomic string `json:"committed_annual_budget_atomic"`
	PeriodProvisionAtomic       string `json:"period_provision_atomic"`
	DistributedAtomic           string `json:"distributed_atomic"`
}

func NewPeriodState(
	startBlock uint64,
	eligibleBonded sdkmath.Int,
	treasuryBalance sdkmath.Int,
	committedAnnualBudget sdkmath.Int,
	params Params,
) (PeriodState, error) {
	if err := params.Validate(); err != nil {
		return PeriodState{}, err
	}
	if startBlock > math.MaxUint64-params.RewardPeriodBlocks {
		return PeriodState{}, errors.New("reward period end block overflows")
	}

	provision, err := ValidateFundedRewardPeriod(
		eligibleBonded,
		treasuryBalance,
		committedAnnualBudget,
		params.AnnualRateBasisPoints,
		params.RewardPeriodBlocks,
		params.BlocksPerYear,
	)
	if err != nil {
		return PeriodState{}, err
	}

	return PeriodState{
		StartBlock:                  startBlock,
		EndBlock:                    startBlock + params.RewardPeriodBlocks,
		AnnualRateBasisPoints:       params.AnnualRateBasisPoints,
		EligibleBondedAtomic:        eligibleBonded.String(),
		CommittedAnnualBudgetAtomic: committedAnnualBudget.String(),
		PeriodProvisionAtomic:       provision.String(),
		DistributedAtomic:           "0",
	}, nil
}

func (p PeriodState) Validate() error {
	if p.EndBlock <= p.StartBlock {
		return errors.New("reward period end block must follow start block")
	}
	if p.AnnualRateBasisPoints > MaxAnnualRateBasisPoints {
		return errors.New("reward period rate exceeds the protocol ceiling")
	}

	eligible, err := parseAtomicAmount(p.EligibleBondedAtomic)
	if err != nil {
		return fmt.Errorf("invalid eligible bonded amount: %w", err)
	}
	budget, err := parseAtomicAmount(p.CommittedAnnualBudgetAtomic)
	if err != nil {
		return fmt.Errorf("invalid committed annual budget: %w", err)
	}
	provision, err := parseAtomicAmount(p.PeriodProvisionAtomic)
	if err != nil {
		return fmt.Errorf("invalid period provision: %w", err)
	}
	distributed, err := parseAtomicAmount(p.DistributedAtomic)
	if err != nil {
		return fmt.Errorf("invalid distributed amount: %w", err)
	}

	if eligible.IsNegative() || budget.IsNegative() ||
		provision.IsNegative() || distributed.IsNegative() {
		return errors.New("reward period amounts cannot be negative")
	}
	if distributed.GT(provision) {
		return errors.New("distributed amount exceeds period provision")
	}

	return nil
}

func (p PeriodState) RemainingProvision() (sdkmath.Int, error) {
	if err := p.Validate(); err != nil {
		return sdkmath.Int{}, err
	}

	provision, _ := parseAtomicAmount(p.PeriodProvisionAtomic)
	distributed, _ := parseAtomicAmount(p.DistributedAtomic)
	return provision.Sub(distributed), nil
}

// ParseStoredAtomicAmount parses a canonical base-10 atomic axtc state value.
func ParseStoredAtomicAmount(value string) (sdkmath.Int, error) {
	return parseAtomicAmount(value)
}

func parseAtomicAmount(value string) (sdkmath.Int, error) {
	parsed, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return sdkmath.Int{}, errors.New("amount is not a base-10 integer")
	}
	return sdkmath.NewIntFromBigInt(parsed), nil
}
