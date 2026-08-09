package contracts

import (
	contractutils "github.com/xitcoin-org/pos-chain/contracts/utils"
	evmtypes "github.com/xitcoin-org/pos-chain/x/vm/types"
)

func LoadIcs20CallerContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("ICS20Caller.json")
}
