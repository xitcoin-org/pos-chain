package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper exposes only the module-account lookup required by validator
// incentives. The treasury account must be registered without mint or burn
// permissions.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper exposes the minimum bank operations required to inspect the
// treasury and pay a recipient from the prefunded module account.
type BankKeeper interface {
	GetBalance(
		ctx context.Context,
		addr sdk.AccAddress,
		denom string,
	) sdk.Coin
	SendCoinsFromModuleToAccount(
		ctx context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
}
