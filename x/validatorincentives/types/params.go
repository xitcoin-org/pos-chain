package types

import "fmt"

// Params defines the governance-controlled operating parameters for funded
// validator incentives.
type Params struct {
	TreasuryReleaseRateBasisPoints uint32 `json:"treasury_release_rate_basis_points"`
	BlocksPerYear                  uint64 `json:"blocks_per_year"`
	CalculationPeriodBlocks        uint64 `json:"calculation_period_blocks"`
}

func DefaultParams() Params {
	return Params{
		TreasuryReleaseRateBasisPoints: DefaultTreasuryReleaseRateBasisPoints,
		BlocksPerYear:                  DefaultBlocksPerYear,
		CalculationPeriodBlocks:        DefaultCalculationPeriodBlocks,
	}
}

func (p Params) Validate() error {
	if p.TreasuryReleaseRateBasisPoints > BasisPointDenominator {
		return fmt.Errorf("treasury release rate exceeds 100 percent")
	}
	if p.BlocksPerYear == 0 {
		return fmt.Errorf("blocks per year must be greater than zero")
	}
	if p.CalculationPeriodBlocks == 0 {
		return fmt.Errorf("calculation period must contain at least one block")
	}
	if p.CalculationPeriodBlocks > p.BlocksPerYear {
		return fmt.Errorf("calculation period cannot exceed one year")
	}
	return nil
}

// ValidateUpdate validates a complete governance-controlled parameter set.
func (p Params) ValidateUpdate(previous Params) error {
	if err := previous.Validate(); err != nil {
		return fmt.Errorf("invalid previous parameters: %w", err)
	}
	if err := p.Validate(); err != nil {
		return err
	}
	return nil
}
