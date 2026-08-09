package cli

import (
	"github.com/spf13/cobra"

	"github.com/xitcoin-org/pos-chain/x/validatoradmission/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

// NewTxCmd returns Xitcoin Validator Admission transaction commands.
func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "validator-admission",
		Short:                      "Xitcoin validator admission commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewApproveValidatorCmd(),
		NewRevokeValidatorCmd(),
	)
	return cmd
}

func addXitcoinTxFlags(cmd *cobra.Command) {
	flags.AddTxFlagsToCmd(cmd)

	if flag := cmd.Flags().Lookup(flags.FlagFees); flag != nil {
		flag.Usage = "Fees in xits, for example 1000000000000000xits"
	}
	if flag := cmd.Flags().Lookup(flags.FlagGasPrices); flag != nil {
		flag.Usage = "Gas prices in xits, for example 1000000000xits"
	}
}

// NewApproveValidatorCmd approves a validator address on Xitcoin.
func NewApproveValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-validator VALIDATOR_ADDRESS",
		Short: "Approve a validator address on Xitcoin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgApproveValidator{
				Authority:        clientCtx.GetFromAddress().String(),
				ValidatorAddress: args[0],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addXitcoinTxFlags(cmd)
	return cmd
}

// NewRevokeValidatorCmd revokes and suspends a validator address on Xitcoin.
func NewRevokeValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-validator VALIDATOR_ADDRESS",
		Short: "Revoke and suspend a validator address on Xitcoin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRevokeValidator{
				Authority:        clientCtx.GetFromAddress().String(),
				ValidatorAddress: args[0],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addXitcoinTxFlags(cmd)
	return cmd
}
