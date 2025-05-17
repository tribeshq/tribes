package tool

import (
	"context"
	"fmt"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tribeshq/tribes/pkg/ethutil"
)

type AdvanceStateTool struct {
	Client        *ethclient.Client
	TxOpts        *bind.TransactOpts
	AppAddress    common.Address
	TokenAddress  common.Address
	PortalAddress common.Address
}

func NewAdvanceStateTool(client *ethclient.Client, txOpts *bind.TransactOpts, appAddress common.Address, tokenAddress common.Address, portalAddress common.Address) *AdvanceStateTool {
	return &AdvanceStateTool{
		Client:        client,
		TxOpts:        txOpts,
		AppAddress:    appAddress,
		TokenAddress:  tokenAddress,
		PortalAddress: portalAddress,
	}
}

func (t *AdvanceStateTool) CreateOrder(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	crowdfundingId, ok := request.Params.Arguments["crowdfunding_id"].(string)
	if !ok {
		return mcp.NewToolResultError("name must be a string"), nil
	}
	amount, ok := request.Params.Arguments["amount"].(string)
	if !ok {
		return mcp.NewToolResultError("amount must be a string"), nil
	}
	interestRate, ok := request.Params.Arguments["interest_rate"].(string)
	if !ok {
		return mcp.NewToolResultError("interest_rate must be a string"), nil
	}

	execLayerData := fmt.Sprintf(`{"path":"create_order","payload":{"crowdfunding_id":%s,"interest_rate":"%s"}}`, crowdfundingId, interestRate)

	execLayerDataBytes := []byte(execLayerData)
	amountBig := new(big.Int)
	amountBig.SetString(amount, 10)

	receipt, err := ethutil.ERC20Deposit(ctx, t.Client, t.TxOpts, t.PortalAddress, t.TokenAddress, t.AppAddress, amountBig, execLayerDataBytes)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(receipt.TxHash.Hex()), nil
}
