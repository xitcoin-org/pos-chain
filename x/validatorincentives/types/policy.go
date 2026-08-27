package types

import (
	"errors"

	sdkmath "cosmossdk.io/math"
)

const (
	BasisPointDenominator uint32 = 10_000
	DefaultTreasuryReleaseRateBasisPoints uint32 = 1_000
	DefaultBlocksPerYear uint64 = 6_311_520
	DefaultCalculationPeriodBlocks uint64 = 17_280
)

func AnnualizedRewardCapacity(
	treasuryBalance sdkmath.Int,
	releaseRateBasisPoints uint32,
) (sdkmath.Int, error) {
	if treasuryBalance.IsNegative() {
		return sdkmath.Int{}, errors.New("treasury balance cannot be negative")
	}
	if releaseRateBasisPoints > BasisPointDenominator {
		return sdkmath.Int{}, errors.New("treasury release rate exceeds 100 percent")
	}

	return treasuryBalance.
		MulRaw(int64(releaseRateBasisPoints)).
		QuoRaw(int64(BasisPointDenominator)), nil
}

func PeriodProvision(
	annualizedCapacity sdkmath.Int,
	periodBlocks uint64,
	blocksPerYear uint64,
) (sdkmath.Int, error) {
	if annualizedCapacity.IsNegative() {
		return sdkmath.Int{}, errors.New("annualized capacity cannot be negative")
	}
	if periodBlocks == 0 {
		return sdkmath.Int{}, errors.New("calculation period must contain at least one block")
	}
	if blocksPerYear == 0 {
		return sdkmath.Int{}, errors.New("blocks per year must be greater than zero")
	}
	if periodBlocks > blocksPerYear {
		return sdkmath.Int{}, errors.New("calculation period cannot exceed one year")
	}

	return annualizedCapacity.
		Mul(sdkmath.NewIntFromUint64(periodBlocks)).
		Quo(sdkmath.NewIntFromUint64(blocksPerYear)), nil
}

func DerivedAPYBasisPoints(
	annualizedCapacity sdkmath.Int,
	eligibleBonded sdkmath.Int,
) (sdkmath.Int, error) {
	if annualizedCapacity.IsNegative() {
		return sdkmath.Int{}, errors.New("annualized capacity cannot be negative")
	}
	if eligibleBonded.IsNegative() {
		return sdkmath.Int{}, errors.New("eligible bonded stake cannot be negative")
	}
	if eligibleBonded.IsZero() {
		return sdkmath.ZeroInt(), nil
	}

	return annualizedCapacity.
		MulRaw(int64(BasisPointDenominator)).
		Quo(eligibleBonded), nil
}

func CumulativeProvision(
	periodProvision sdkmath.Int,
	elapsedBlocks uint64,
	periodBlocks uint64,
) (sdkmath.Int, error) {
	if periodProvision.IsNegative() {
		return sdkmath.Int{}, errors.New("period provision cannot be negative")
	}
	if periodBlocks == 0 {
		return sdkmath.Int{}, errors.New("calculation period must contain at least one block")
	}
	if elapsedBlocks > periodBlocks {
		return sdkmath.Int{}, errors.New("elapsed blocks exceed the calculation period")
	}

	return periodProvision.
		Mul(sdkmath.NewIntFromUint64(elapsedBlocks)).
		Quo(sdkmath.NewIntFromUint64(periodBlocks)), nil
}
