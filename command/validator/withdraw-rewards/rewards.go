package rewards

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-edge/command"
	bridgeHelper "github.com/0xPolygon/polygon-edge/command/bridge/helper"
	"github.com/0xPolygon/polygon-edge/command/helper"
	polybftsecrets "github.com/0xPolygon/polygon-edge/command/secrets/init"
	validatorHelper "github.com/0xPolygon/polygon-edge/command/validator/helper"
	"github.com/0xPolygon/polygon-edge/consensus/polybft/contractsapi"
	"github.com/0xPolygon/polygon-edge/contracts"
	"github.com/0xPolygon/polygon-edge/helper/common"
	"github.com/0xPolygon/polygon-edge/txrelayer"
	"github.com/0xPolygon/polygon-edge/types"
)

var params withdrawRewardsParams

func GetCommand() *cobra.Command {
	unstakeCmd := &cobra.Command{
		Use:     "withdraw-rewards",
		Short:   "Withdraws pending rewards on child chain for given validator",
		PreRunE: runPreRun,
		RunE:    runCommand,
	}

	helper.RegisterJSONRPCFlag(unstakeCmd)
	setFlags(unstakeCmd)

	return unstakeCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.accountDir,
		polybftsecrets.AccountDirFlag,
		"",
		polybftsecrets.AccountDirFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.accountConfig,
		polybftsecrets.AccountConfigFlag,
		"",
		polybftsecrets.AccountConfigFlagDesc,
	)

	cmd.Flags().DurationVar(
		&params.txTimeout,
		helper.TxTimeoutFlag,
		150*time.Second,
		helper.TxTimeoutDesc,
	)

	cmd.MarkFlagsMutuallyExclusive(polybftsecrets.AccountDirFlag, polybftsecrets.AccountConfigFlag)
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	params.jsonRPC = helper.GetJSONRPCAddress(cmd)

	return params.validateFlags()
}

func runCommand(cmd *cobra.Command, _ []string) error {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	validatorAccount, err := validatorHelper.GetAccount(params.accountDir, params.accountConfig)
	if err != nil {
		return err
	}

	validatorAddr := validatorAccount.Ecdsa.Address()

	txRelayer, err := txrelayer.NewTxRelayer(txrelayer.WithIPAddress(params.jsonRPC),
		txrelayer.WithReceiptsTimeout(params.txTimeout))
	if err != nil {
		return err
	}

	encoded, err := contractsapi.EpochManager.Abi.Methods["pendingRewards"].Encode([]interface{}{validatorAddr})
	if err != nil {
		return err
	}

	response, err := txRelayer.Call(validatorAddr, contracts.EpochManagerContract, encoded)
	if err != nil {
		return err
	}

	amount, err := common.ParseUint256orHex(&response)
	if err != nil {
		return err
	}

	encoded, err = contractsapi.EpochManager.Abi.Methods["withdrawReward"].Encode([]interface{}{})
	if err != nil {
		return err
	}

	txn := bridgeHelper.CreateTransaction(validatorAddr, &contracts.EpochManagerContract, encoded, nil, false)

	receipt, err := txRelayer.SendTransaction(txn, validatorAccount.Ecdsa)
	if err != nil {
		return err
	}

	if receipt.Status != uint64(types.ReceiptSuccess) {
		return fmt.Errorf("withdraw transaction failed on block: %d", receipt.BlockNumber)
	}

	result := &withdrawRewardResult{
		ValidatorAddress: validatorAccount.Ecdsa.Address().String(),
		RewardAmount:     amount.Uint64(),
	}

	outputter.WriteCommandResult(result)

	return nil
}
