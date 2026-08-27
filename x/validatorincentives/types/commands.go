package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateParamsCommand is the deterministic domain command used by the future
// protobuf message server to request an authorized parameter update.
type UpdateParamsCommand struct {
	Authority string
	Params    Params
}

func (c UpdateParamsCommand) ValidateBasic() error {
	if err := validateCommandAuthority(c.Authority); err != nil {
		return err
	}
	if err := c.Params.Validate(); err != nil {
		return fmt.Errorf("invalid incentive parameters: %w", err)
	}
	return nil
}

func validateCommandAuthority(authority string) error {
	if authority == "" {
		return errors.New("authority is required")
	}
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}
	return nil
}
