package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/xitcoin-org/pos-chain/x/validatorincentives/types"
)

// NewTxCmd returns Xitcoin validator incentive transaction commands.
func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "validator-incentives",
		Short:                      "Xitcoin prefunded validator incentive commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewUpdateParamsCmd(),
		NewActivateFundedPeriodCmd(),
	)
	return cmd
}

func addXitcoinTxFlags(cmd *cobra.Command) {
	flags.AddTxFlagsToCmd(cmd)

	if flag := cmd.Flags().Lookup(flags.FlagFees); flag != nil {
		flag.Usage = "Fees in axtc, for example 1000000000000000axtc"
	}
	if flag := cmd.Flags().Lookup(flags.FlagGasPrices); flag != nil {
		flag.Usage = "Gas prices in axtc, for example 1000000000axtc"
	}
}

// NewUpdateParamsCmd updates funded incentive operating parameters.
func NewUpdateParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "update-params ANNUAL_RATE_BASIS_POINTS " +
			"BLOCKS_PER_YEAR REWARD_PERIOD_BLOCKS",
		Short: "Update validator incentive parameters",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			rate, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}
			blocksPerYear, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			rewardPeriodBlocks, err := strconv.ParseUint(
				args[2],
				10,
				64,
			)
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateParams{
				Authority:             clientCtx.GetFromAddress().String(),
				AnnualRateBasisPoints: uint32(rate),
				BlocksPerYear:         blocksPerYear,
				RewardPeriodBlocks:    rewardPeriodBlocks,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(
				clientCtx,
				cmd.Flags(),
				msg,
			)
		},
	}

	addXitcoinTxFlags(cmd)
	return cmd
}

// NewActivateFundedPeriodCmd activates a period using a committed annual
// budget. Eligible bonded stake and the treasury balance are read on-chain.
func NewActivateFundedPeriodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-funded-period COMMITTED_ANNUAL_BUDGET_ATOMIC",
		Short: "Activate a prefunded validator incentive period",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgActivateFundedPeriod{
				Authority:                   clientCtx.GetFromAddress().String(),
				CommittedAnnualBudgetAtomic: args[0],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(
				clientCtx,
				cmd.Flags(),
				msg,
			)
		},
	}

	addXitcoinTxFlags(cmd)
	return cmd
}
