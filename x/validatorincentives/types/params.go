package types

import "fmt"

// Params defines the governance-controlled operating parameters for funded
// validator incentives. Hard protocol ceilings remain compile-time safeguards.
type Params struct {
	AnnualRateBasisPoints uint32 `json:"annual_rate_basis_points"`
	BlocksPerYear         uint64 `json:"blocks_per_year"`
	RewardPeriodBlocks    uint64 `json:"reward_period_blocks"`
}

func DefaultParams() Params {
	return Params{
		AnnualRateBasisPoints: DefaultAnnualRateBasisPoints,
		BlocksPerYear:         DefaultBlocksPerYear,
		RewardPeriodBlocks:    DefaultBlocksPerYear / 4,
	}
}

func (p Params) Validate() error {
	if p.AnnualRateBasisPoints > MaxAnnualRateBasisPoints {
		return fmt.Errorf("annual rate exceeds the protocol ceiling")
	}
	if p.BlocksPerYear == 0 {
		return fmt.Errorf("blocks per year must be greater than zero")
	}
	if p.RewardPeriodBlocks == 0 {
		return fmt.Errorf("reward period must contain at least one block")
	}
	if p.RewardPeriodBlocks > p.BlocksPerYear {
		return fmt.Errorf("reward period cannot exceed one year")
	}
	if p.BlocksPerYear%p.RewardPeriodBlocks != 0 {
		return fmt.Errorf("reward period must divide blocks per year exactly")
	}

	return nil
}

// ValidateUpdate applies both structural validation and the rate-of-change
// safeguard between consecutive governance-controlled parameter sets.
func (p Params) ValidateUpdate(previous Params) error {
	if err := previous.Validate(); err != nil {
		return fmt.Errorf("invalid previous parameters: %w", err)
	}
	if err := p.Validate(); err != nil {
		return err
	}
	if err := ValidateRateTransition(
		previous.AnnualRateBasisPoints,
		p.AnnualRateBasisPoints,
	); err != nil {
		return err
	}

	return nil
}
