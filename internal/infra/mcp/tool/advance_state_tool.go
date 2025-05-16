package tool

import (
	"context"
	"errors"
	"fmt"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tribeshq/tribes/pkg/ethutil"
)

type AdvanceStateTool struct {
	Client *ethclient.Client
	TxOpts *bind.TransactOpts
}

func NewAdvanceStateTool(client *ethclient.Client, txOpts *bind.TransactOpts) *AdvanceStateTool {
	return &AdvanceStateTool{
		Client: client,
		TxOpts: txOpts,
	}
}

func (t *AdvanceStateTool) CreateOrder(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	crowdfundingId, ok := request.Params.Arguments["crowdfunding_id"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}
	amount, ok := request.Params.Arguments["amount"].(string)
	if !ok {
		return nil, errors.New("amount must be a string")
	}
	interestRate, ok := request.Params.Arguments["interest_rate"].(string)
	if !ok {
		return nil, errors.New("interest_rate must be a string")
	}

	execLayerData := fmt.Sprintf(`{"path":"create_order","payload":{"crowdfunding_id":%s,"interest_rate":"%s"}}`, crowdfundingId, interestRate)

	execLayerDataBytes := []byte(execLayerData)
	amountBig := new(big.Int)
	amountBig.SetString(amount, 10)

	receipt, err := ethutil.ERC20Deposit(ctx, t.Client, t.TxOpts, portalAddress, tokenAddress, appAddress, amountBig, execLayerDataBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deposit ERC20: %w", err)
	}

	return mcp.NewToolResultText(receipt.TxHash.Hex()), nil
}
