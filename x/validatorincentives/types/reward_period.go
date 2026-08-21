package types

import (
	"errors"

	sdkmath "cosmossdk.io/math"
)

const (
	// MaxQuarterlyRateIncreaseBasisPoints limits an APR increase to one
	// percentage point between consecutive reward periods.
	MaxQuarterlyRateIncreaseBasisPoints uint32 = 100
)

// RequiredAnnualProvision calculates the fully funded annual amount required
// to advertise rateBasisPoints against eligible bonded stake.
func RequiredAnnualProvision(
	eligibleBonded sdkmath.Int,
	rateBasisPoints uint32,
) (sdkmath.Int, error) {
	if eligibleBonded.IsNegative() {
		return sdkmath.Int{}, errors.New("eligible bonded stake cannot be negative")
	}
	if rateBasisPoints > MaxAnnualRateBasisPoints {
		return sdkmath.Int{}, errors.New("annual rate exceeds the protocol ceiling")
	}

	return eligibleBonded.
		MulRaw(int64(rateBasisPoints)).
		QuoRaw(int64(BasisPointDenominator)), nil
}

// ValidateRateTransition prevents abrupt APR increases. Reductions may be
// immediate so an authority can preserve solvency or pause distributions.
func ValidateRateTransition(previous, next uint32) error {
	if previous > MaxAnnualRateBasisPoints || next > MaxAnnualRateBasisPoints {
		return errors.New("annual rate exceeds the protocol ceiling")
	}
	if next > previous &&
		next-previous > MaxQuarterlyRateIncreaseBasisPoints {
		return errors.New("annual rate increase exceeds the quarterly limit")
	}

	return nil
}

// ValidateFundedRewardPeriod returns the provision committed to one reward
// period. Activation fails unless the advertised APR has a fully committed
// annual budget and the treasury can fund the complete period.
//
// periodBlocks must be no longer than one year. Reward execution can use
// CumulativeProvision and BlockProvision to release this amount without
// rounding drift.
func ValidateFundedRewardPeriod(
	eligibleBonded sdkmath.Int,
	treasuryBalance sdkmath.Int,
	committedAnnualBudget sdkmath.Int,
	rateBasisPoints uint32,
	periodBlocks uint64,
	blocksPerYear uint64,
) (sdkmath.Int, error) {
	if treasuryBalance.IsNegative() {
		return sdkmath.Int{}, errors.New("treasury balance cannot be negative")
	}
	if committedAnnualBudget.IsNegative() {
		return sdkmath.Int{}, errors.New("committed annual budget cannot be negative")
	}
	if periodBlocks == 0 {
		return sdkmath.Int{}, errors.New("reward period must contain at least one block")
	}

	requiredAnnual, err := RequiredAnnualProvision(
		eligibleBonded,
		rateBasisPoints,
	)
	if err != nil {
		return sdkmath.Int{}, err
	}
	if committedAnnualBudget.LT(requiredAnnual) {
		return sdkmath.Int{}, errors.New("committed annual budget does not fully fund the active APR")
	}

	periodProvision, err := CumulativeProvision(
		requiredAnnual,
		periodBlocks,
		blocksPerYear,
	)
	if err != nil {
		return sdkmath.Int{}, err
	}
	if treasuryBalance.LT(periodProvision) {
		return sdkmath.Int{}, errors.New("treasury balance does not fund the complete reward period")
	}

	return periodProvision, nil
}
