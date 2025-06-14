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
	auctionId, ok := request.Params.Arguments["auction_id"].(string)
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

	execLayerData := fmt.Sprintf(`{"path":"order/create","data":{"auction_id":%s,"interest_rate":"%s"}}`, auctionId, interestRate)

	execLayerDataBytes := []byte(execLayerData)
	amountBig := new(big.Int)
	amountBig.SetString(amount, 10)

	receipt, err := ethutil.ERC20Deposit(ctx, t.Client, t.TxOpts, t.PortalAddress, t.TokenAddress, t.AppAddress, amountBig, execLayerDataBytes)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(receipt.TxHash.Hex()), nil
}

// func (t *AdvanceStateTool) CreateSocialAccount(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	userId, ok := request.Params.Arguments["user_id"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("user_id must be a number"), nil
// 	}
// 	username, ok := request.Params.Arguments["username"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("username must be a string"), nil
// 	}
// 	followers, ok := request.Params.Arguments["followers"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("followers must be a number"), nil
// 	}
// 	platform, ok := request.Params.Arguments["platform"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("platform must be a string"), nil
// 	}
// 	createdAt, ok := request.Params.Arguments["created_at"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("created_at must be a number"), nil
// 	}
// 	input := map[string]interface{}{
// 		"user_id":    uint(userId),
// 		"username":   username,
// 		"followers":  uint(followers),
// 		"platform":   platform,
// 		"created_at": int64(createdAt),
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// func (t *AdvanceStateTool) DeleteSocialAccount(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	socialAccountId, ok := request.Params.Arguments["social_account_id"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("social_account_id must be a number"), nil
// 	}
// 	input := map[string]interface{}{
// 		"social_account_id": uint(socialAccountId),
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// func (t *AdvanceStateTool) CreateUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	role, ok := request.Params.Arguments["role"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("role must be a string"), nil
// 	}
// 	address, ok := request.Params.Arguments["address"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("address must be a string"), nil
// 	}
// 	input := map[string]interface{}{
// 		"role":    role,
// 		"address": address,
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// func (t *AdvanceStateTool) UpdateUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	role, ok := request.Params.Arguments["role"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("role must be a string"), nil
// 	}
// 	address, ok := request.Params.Arguments["address"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("address must be a string"), nil
// 	}
// 	investmentLimit, _ := request.Params.Arguments["investment_limit"].(string)      // opcional
// 	debtIssuanceLimit, _ := request.Params.Arguments["debt_issuance_limit"].(string) // opcional
// 	input := map[string]interface{}{
// 		"role":                role,
// 		"address":             address,
// 		"investment_limit":    investmentLimit,
// 		"debt_issuance_limit": debtIssuanceLimit,
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// func (t *AdvanceStateTool) DeleteUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	address, ok := request.Params.Arguments["address"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("address must be a string"), nil
// 	}
// 	input := map[string]interface{}{
// 		"address": address,
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// func (t *AdvanceStateTool) Withdraw(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	token, ok := request.Params.Arguments["token"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("token must be a string"), nil
// 	}
// 	amount, ok := request.Params.Arguments["amount"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("amount must be a string"), nil
// 	}
// 	input := map[string]interface{}{
// 		"token":  token,
// 		"amount": amount,
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// // CONTRACT ADVANCE HANDLERS
// func (t *AdvanceStateTool) CreateContract(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	symbol, ok := request.Params.Arguments["symbol"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("symbol must be a string"), nil
// 	}
// 	address, ok := request.Params.Arguments["address"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("address must be a string"), nil
// 	}
// 	input := map[string]interface{}{
// 		"symbol":  symbol,
// 		"address": address,
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// func (t *AdvanceStateTool) UpdateContract(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	id, ok := request.Params.Arguments["id"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("id must be a number"), nil
// 	}
// 	symbol, ok := request.Params.Arguments["symbol"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("symbol must be a string"), nil
// 	}
// 	address, ok := request.Params.Arguments["address"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("address must be a string"), nil
// 	}
// 	input := map[string]interface{}{
// 		"id":      uint(id),
// 		"symbol":  symbol,
// 		"address": address,
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// func (t *AdvanceStateTool) DeleteContract(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	symbol, ok := request.Params.Arguments["symbol"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("symbol must be a string"), nil
// 	}
// 	input := map[string]interface{}{
// 		"symbol": symbol,
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// // CROWDFUNDING ADVANCE HANDLERS
// func (t *AdvanceStateTool) CreateAuction(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	debtIssued, ok := request.Params.Arguments["debt_issued"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("debt_issued must be a string"), nil
// 	}
// 	maxInterestRate, ok := request.Params.Arguments["max_interest_rate"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("max_interest_rate must be a string"), nil
// 	}
// 	fundraisingDuration, ok := request.Params.Arguments["fundraising_duration"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("fundraising_duration must be a number"), nil
// 	}
// 	closesAt, ok := request.Params.Arguments["closes_at"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("closes_at must be a number"), nil
// 	}
// 	maturityAt, ok := request.Params.Arguments["maturity_at"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("maturity_at must be a number"), nil
// 	}
// 	proof, ok := request.Params.Arguments["proof"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("proof must be a string"), nil
// 	}
// 	input := map[string]interface{}{
// 		"debt_issued":          debtIssued,
// 		"max_interest_rate":    maxInterestRate,
// 		"fundraising_duration": int64(fundraisingDuration),
// 		"closes_at":            int64(closesAt),
// 		"maturity_at":          int64(maturityAt),
// 		"proof":                proof,
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	// Aqui você deve usar ERC20Deposit, montando execLayerData com o input acima
// 	// Exemplo:
// 	amount, ok := request.Params.Arguments["amount"].(string)
// 	if !ok {
// 		return mcp.NewToolResultError("amount must be a string (deposit amount)"), nil
// 	}
// 	execLayerData := fmt.Sprintf(`{"path":"auction/create","data":%s}`, string(inputBytes))
// 	execLayerDataBytes := []byte(execLayerData)
// 	amountBig := new(big.Int)
// 	amountBig.SetString(amount, 10)
// 	receipt, err := ethutil.ERC20Deposit(ctx, t.Client, t.TxOpts, t.PortalAddress, t.TokenAddress, t.AppAddress, amountBig, execLayerDataBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText(receipt.TxHash.Hex()), nil
// }

// func (t *AdvanceStateTool) CloseAuction(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	auctionId, ok := request.Params.Arguments["auction_id"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("auction_id must be a number"), nil
// 	}
// 	input := map[string]interface{}{
// 		"auction_id": uint(auctionId),
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// func (t *AdvanceStateTool) SettleAuction(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	auctionId, ok := request.Params.Arguments["auction_id"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("auction_id must be a number"), nil
// 	}
// 	input := map[string]interface{}{
// 		"auction_id": uint(auctionId),
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// func (t *AdvanceStateTool) UpdateAuction(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	auctionId, ok := request.Params.Arguments["auction_id"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("auction_id must be a number"), nil
// 	}
// 	input := map[string]interface{}{
// 		"auction_id": uint(auctionId),
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }

// func (t *AdvanceStateTool) DeleteAuction(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	auctionId, ok := request.Params.Arguments["auction_id"].(float64)
// 	if !ok {
// 		return mcp.NewToolResultError("auction_id must be a number"), nil
// 	}
// 	input := map[string]interface{}{
// 		"auction_id": uint(auctionId),
// 	}
// 	inputBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return mcp.NewToolResultError("failed to marshal input"), nil
// 	}
// 	_, _, err = ethutil.AddInput(ctx, t.Client, t.TxOpts, t.AppAddress, t.AppAddress, inputBytes)
// 	if err != nil {
// 		return mcp.NewToolResultError(err.Error()), nil
// 	}
// 	return mcp.NewToolResultText("advance state input sent"), nil
// }
