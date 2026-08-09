package evmd

import (
	cmn "github.com/xitcoin-org/pos-chain/precompiles/common"
	evmtypes "github.com/xitcoin-org/pos-chain/x/vm/types"
)

type BankKeeper interface {
	evmtypes.BankKeeper
	cmn.BankKeeper
}
