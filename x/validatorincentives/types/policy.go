package types

import (
	"errors"
	"math/big"

	sdkmath "cosmossdk.io/math"
)

const (
	// BasisPointDenominator represents 100%.
	BasisPointDenominator uint32 = 10_000

	// MaxAnnualRateBasisPoints is the hard treasury-funded annual-rate ceiling:
	// 20% of eligible bonded stake. Exceeding it requires a software upgrade.
	MaxAnnualRateBasisPoints uint32 = 2_000

	// DefaultAnnualRateBasisPoints starts at the conservative 8% operating
	// target. Governance may adjust it within the hard safety ceiling.
	DefaultAnnualRateBasisPoints uint32 = 800

	// DefaultBlocksPerYear matches the verified candidate-testnet mint params.
	DefaultBlocksPerYear uint64 = 6_311_520
)

var (

	// DefaultInitialTreasuryReference is the non-binding mainnet planning
	// reference of 100,000,000 XTC at 18 decimals.
	DefaultInitialTreasuryReference = sdkmath.NewIntFromBigInt(mustInt(
		"100000000000000000000000000",
	))
)

// AnnualProvision returns the maximum treasury-funded provision for a year.
//
// The result is the minimum of:
//   - bonded stake multiplied by the configured annual rate;
//   - the annual budget already funded and committed;
//   - the available funded treasury balance.
//
// The function never mints tokens and always rounds down to atomic units.
func AnnualProvision(
	bonded sdkmath.Int,
	treasuryBalance sdkmath.Int,
	annualBudget sdkmath.Int,
	rateBasisPoints uint32,
) (sdkmath.Int, error) {
	if bonded.IsNegative() {
		return sdkmath.Int{}, errors.New("bonded stake cannot be negative")
	}
	if treasuryBalance.IsNegative() {
		return sdkmath.Int{}, errors.New("treasury balance cannot be negative")
	}
	if annualBudget.IsNegative() {
		return sdkmath.Int{}, errors.New("annual funded budget cannot be negative")
	}
	if rateBasisPoints > MaxAnnualRateBasisPoints {
		return sdkmath.Int{}, errors.New("annual rate exceeds the protocol ceiling")
	}

	rateLimited := bonded.
		MulRaw(int64(rateBasisPoints)).
		QuoRaw(int64(BasisPointDenominator))

	return minInt(rateLimited, annualBudget, treasuryBalance), nil
}

// CumulativeProvision returns the cumulative amount that may have been released
// after elapsedBlocks in a fixed annual period.
//
// Using a cumulative target prevents rounding drift: at blocksPerYear the
// cumulative amount equals annualProvision exactly.
func CumulativeProvision(
	annualProvision sdkmath.Int,
	elapsedBlocks uint64,
	blocksPerYear uint64,
) (sdkmath.Int, error) {
	if annualProvision.IsNegative() {
		return sdkmath.Int{}, errors.New("annual provision cannot be negative")
	}
	if blocksPerYear == 0 {
		return sdkmath.Int{}, errors.New("blocks per year must be greater than zero")
	}
	if elapsedBlocks > blocksPerYear {
		return sdkmath.Int{}, errors.New("elapsed blocks exceed the annual period")
	}

	return annualProvision.
		Mul(sdkmath.NewIntFromUint64(elapsedBlocks)).
		Quo(sdkmath.NewIntFromUint64(blocksPerYear)), nil
}

// BlockProvision returns the amount released for one 1-indexed block in a
// fixed annual period. The sum across the full period equals annualProvision.
func BlockProvision(
	annualProvision sdkmath.Int,
	blockIndex uint64,
	blocksPerYear uint64,
) (sdkmath.Int, error) {
	if blockIndex == 0 {
		return sdkmath.Int{}, errors.New("block index must be 1-indexed")
	}

	current, err := CumulativeProvision(
		annualProvision,
		blockIndex,
		blocksPerYear,
	)
	if err != nil {
		return sdkmath.Int{}, err
	}

	previous, err := CumulativeProvision(
		annualProvision,
		blockIndex-1,
		blocksPerYear,
	)
	if err != nil {
		return sdkmath.Int{}, err
	}

	return current.Sub(previous), nil
}

func minInt(values ...sdkmath.Int) sdkmath.Int {
	result := values[0]
	for _, value := range values[1:] {
		if value.LT(result) {
			result = value
		}
	}
	return result
}

func mustInt(value string) *big.Int {
	result, ok := new(big.Int).SetString(value, 10)
	if !ok {
		panic("invalid integer constant")
	}
	return result
}
