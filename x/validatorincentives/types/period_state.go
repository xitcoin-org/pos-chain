package types

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	sdkmath "cosmossdk.io/math"
)

// PeriodState is the deterministic snapshot used for one daily calculation
// period. Amounts are canonical base-10 atomic axtc strings.
type PeriodState struct {
	StartBlock                    uint64 `json:"start_block"`
	EndBlock                      uint64 `json:"end_block"`
	TreasuryReleaseRateBasisPoints uint32 `json:"treasury_release_rate_basis_points"`
	TreasuryBalanceAtomic         string `json:"treasury_balance_atomic"`
	EligibleBondedAtomic          string `json:"eligible_bonded_atomic"`
	AnnualizedCapacityAtomic      string `json:"annualized_capacity_atomic"`
	DerivedAPYBasisPoints         string `json:"derived_apy_basis_points"`
	PeriodProvisionAtomic         string `json:"period_provision_atomic"`
	DistributedAtomic             string `json:"distributed_atomic"`
}

func NewPeriodState(
	startBlock uint64,
	eligibleBonded sdkmath.Int,
	treasuryBalance sdkmath.Int,
	params Params,
) (PeriodState, error) {
	if err := params.Validate(); err != nil {
		return PeriodState{}, err
	}
	if startBlock > math.MaxUint64-params.CalculationPeriodBlocks {
		return PeriodState{}, errors.New("calculation period end block overflows")
	}
	if eligibleBonded.IsNegative() {
		return PeriodState{}, errors.New("eligible bonded stake cannot be negative")
	}

	annualized, err := AnnualizedRewardCapacity(
		treasuryBalance,
		params.TreasuryReleaseRateBasisPoints,
	)
	if err != nil {
		return PeriodState{}, err
	}
	apy, err := DerivedAPYBasisPoints(annualized, eligibleBonded)
	if err != nil {
		return PeriodState{}, err
	}

	provision := sdkmath.ZeroInt()
	if eligibleBonded.IsPositive() {
		provision, err = PeriodProvision(
			annualized,
			params.CalculationPeriodBlocks,
			params.BlocksPerYear,
		)
		if err != nil {
			return PeriodState{}, err
		}
	}

	return PeriodState{
		StartBlock:                     startBlock,
		EndBlock:                       startBlock + params.CalculationPeriodBlocks,
		TreasuryReleaseRateBasisPoints: params.TreasuryReleaseRateBasisPoints,
		TreasuryBalanceAtomic:          treasuryBalance.String(),
		EligibleBondedAtomic:           eligibleBonded.String(),
		AnnualizedCapacityAtomic:       annualized.String(),
		DerivedAPYBasisPoints:          apy.String(),
		PeriodProvisionAtomic:          provision.String(),
		DistributedAtomic:              "0",
	}, nil
}

func (p PeriodState) Validate() error {
	if p.EndBlock <= p.StartBlock {
		return errors.New("calculation period end block must follow start block")
	}
	if p.TreasuryReleaseRateBasisPoints > BasisPointDenominator {
		return errors.New("treasury release rate exceeds 100 percent")
	}

	fields := []struct {
		name  string
		value string
	}{
		{"treasury balance", p.TreasuryBalanceAtomic},
		{"eligible bonded", p.EligibleBondedAtomic},
		{"annualized capacity", p.AnnualizedCapacityAtomic},
		{"derived APY", p.DerivedAPYBasisPoints},
		{"period provision", p.PeriodProvisionAtomic},
		{"distributed", p.DistributedAtomic},
	}
	parsed := make([]sdkmath.Int, len(fields))
	for index, field := range fields {
		amount, err := parseAtomicAmount(field.value)
		if err != nil {
			return fmt.Errorf("invalid %s amount: %w", field.name, err)
		}
		if amount.IsNegative() {
			return fmt.Errorf("%s amount cannot be negative", field.name)
		}
		parsed[index] = amount
	}
	if parsed[5].GT(parsed[4]) {
		return errors.New("distributed amount exceeds period provision")
	}
	if parsed[4].GT(parsed[0]) {
		return errors.New("period provision exceeds treasury snapshot")
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
