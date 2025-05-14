// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains functions to help using the Go-ethereum library.
// It is not the objective of this package to replace or hide Go-ethereum.
package ethutil

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	contracts "github.com/tribeshq/tribes/pkg/rollups_contracts"
)

// Gas limit when sending transactions.
const GasLimit = 30_000_000

// Dev mnemonic used by Foundry/Anvil.
const FoundryMnemonic = "test test test test test test test test test test test junk"

func AddInput(
	ctx context.Context,
	client *ethclient.Client,
	transactionOpts *bind.TransactOpts,
	inputBoxAddress common.Address,
	application common.Address,
	input []byte,
) (uint64, uint64, error) {
	if client == nil {
		return 0, 0, fmt.Errorf("AddInput: client is nil")
	}
	inputBox, err := contracts.NewIInputBox(inputBoxAddress, client)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to connect to InputBox contract: %v", err)
	}
	receipt, err := sendTransaction(
		ctx, client, transactionOpts, big.NewInt(0), GasLimit,
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return inputBox.AddInput(txOpts, application, input)
		},
	)
	if err != nil {
		return 0, 0, err
	}
	index, err := getInputIndex(inputBoxAddress, inputBox, receipt)
	if err != nil {
		return 0, 0, err
	}
	return index, receipt.BlockNumber.Uint64(), nil
}

func getInputIndex(
	inputBoxAddress common.Address,
	inputBox *contracts.IInputBox,
	receipt *types.Receipt,
) (uint64, error) {
	for _, log := range receipt.Logs {
		if log.Address != inputBoxAddress {
			continue
		}
		inputAdded, err := inputBox.ParseInputAdded(*log)
		if err != nil {
			return 0, fmt.Errorf("failed to parse input added event: %v", err)
		}
		// We assume that uint64 will fit all dapp inputs for now
		return inputAdded.Index.Uint64(), nil
	}
	return 0, fmt.Errorf("input index not found")
}

func ERC20Deposit(
	ctx context.Context,
	client *ethclient.Client,
	transactionOpts *bind.TransactOpts,
	portalAddress common.Address,
	tokenAddress common.Address,
	accountAddress common.Address,
	amount *big.Int,
	execLayerData []byte,
) (*types.Receipt, error) {
	if client == nil {
		return nil, fmt.Errorf("AddERC20Deposit: client is nil")
	}

	portal, err := contracts.NewIERC20Portal(portalAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ERC20Portal contract: %v", err)
	}

	receipt, err := sendTransaction(
		ctx, client, transactionOpts, big.NewInt(0), GasLimit, // Value is 0 for token deposits
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return portal.DepositERC20Tokens(txOpts, tokenAddress, accountAddress, amount, execLayerData)
		},
	)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func ValidateOutput(
	ctx context.Context,
	client *ethclient.Client,
	appAddr common.Address,
	index uint64,
	output []byte,
	outputHashesSiblings []common.Hash,
) error {
	if client == nil {
		return fmt.Errorf("ValidateOutput: client is nil")
	}
	proof := contracts.OutputValidityProof{
		OutputIndex:          index,
		OutputHashesSiblings: make([][32]byte, len(outputHashesSiblings)),
	}

	for i, hash := range outputHashesSiblings {
		copy(proof.OutputHashesSiblings[i][:], hash[:])
	}

	app, err := contracts.NewIApplication(appAddr, client)
	if err != nil {
		return fmt.Errorf("failed to connect to CartesiDapp contract: %v", err)
	}
	return app.ValidateOutput(&bind.CallOpts{Context: ctx}, output, proof)
}

func ExecuteOutput(
	ctx context.Context,
	client *ethclient.Client,
	transactionOpts *bind.TransactOpts,
	appAddr common.Address,
	index uint64,
	output []byte,
	outputHashesSiblings []common.Hash,
) (*common.Hash, error) {
	if client == nil {
		return nil, fmt.Errorf("ExecuteOutput: client is nil")
	}
	proof := contracts.OutputValidityProof{
		OutputIndex:          index,
		OutputHashesSiblings: make([][32]byte, len(outputHashesSiblings)),
	}

	for i, hash := range outputHashesSiblings {
		copy(proof.OutputHashesSiblings[i][:], hash[:])
	}

	app, err := contracts.NewIApplication(appAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to CartesiDapp contract: %v", err)
	}
	receipt, err := sendTransaction(
		ctx, client, transactionOpts, big.NewInt(0), GasLimit,
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return app.ExecuteOutput(txOpts, output, proof)
		},
	)
	if err != nil {
		return nil, err
	}

	return &receipt.TxHash, nil
}
