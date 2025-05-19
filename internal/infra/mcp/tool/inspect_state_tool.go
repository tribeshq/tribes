package tool

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mark3labs/mcp-go/mcp"
)

// 1. list all open crowdfundings;
// 2. list all orders;
// 3. list crowdfunding by creator;

type InspectStateTool struct {
	BaseURL    string
	Client     *http.Client
	AppAddress common.Address
}

func NewInspectStateTool(baseURL string, appAddress common.Address) *InspectStateTool {
	return &InspectStateTool{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		BaseURL:    baseURL,
		AppAddress: appAddress,
	}
}

func (t *InspectStateTool) ListAllCrowdfundings(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	payload := map[string]interface{}{
		"path": "crowdfunding",
		"data": map[string]interface{}{},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get crowdfundings: status %d, body: %s", resp.StatusCode, string(body))), nil
	}

	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}

	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(decoded)), nil
}

func (t *InspectStateTool) ListAllOrders(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	payload := map[string]interface{}{
		"path": "order",
		"data": map[string]interface{}{},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get orders: status %d, body: %s", resp.StatusCode, string(body))), nil
	}

	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}

	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(decoded)), nil
}

func (t *InspectStateTool) ListCrowdfundingByCreator(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	creator, ok := request.Params.Arguments["creator"].(string)
	if !ok {
		return mcp.NewToolResultError("creator must be a string"), nil
	}

	payload := map[string]interface{}{
		"path": "crowdfunding/creator",
		"data": map[string]interface{}{
			"creator": creator,
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get crowdfundings by creator: status %d, body: %s", resp.StatusCode, string(body))), nil
	}

	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}

	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(decoded)), nil
}

// SocialAccount MCP inspect handlers
func (t *InspectStateTool) ListAllSocialAccounts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	payload := map[string]interface{}{
		"path": "social_account",
		"data": map[string]interface{}{},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get social accounts: status %d, body: %s", resp.StatusCode, string(body))), nil
	}
	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}
	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(decoded)), nil
}

func (t *InspectStateTool) FindSocialAccountById(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.Params.Arguments["id"].(string)
	if !ok {
		return mcp.NewToolResultError("id must be a string"), nil
	}
	payload := map[string]interface{}{
		"path": "social_account/id",
		"data": map[string]interface{}{"id": id},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get social account by id: status %d, body: %s", resp.StatusCode, string(body))), nil
	}
	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}
	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(decoded)), nil
}

func (t *InspectStateTool) FindSocialAccountsByUserId(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	userId, ok := request.Params.Arguments["user_id"].(string)
	if !ok {
		return mcp.NewToolResultError("user_id must be a string"), nil
	}
	payload := map[string]interface{}{
		"path": "social_account/user_id",
		"data": map[string]interface{}{"user_id": userId},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get social accounts by user id: status %d, body: %s", resp.StatusCode, string(body))), nil
	}
	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}
	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(decoded)), nil
}

// User MCP inspect handlers
func (t *InspectStateTool) ListAllUsers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	payload := map[string]interface{}{
		"path": "user",
		"data": map[string]interface{}{},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get users: status %d, body: %s", resp.StatusCode, string(body))), nil
	}
	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}
	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(decoded)), nil
}

func (t *InspectStateTool) FindUserByAddress(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return mcp.NewToolResultError("address must be a string"), nil
	}
	payload := map[string]interface{}{
		"path": "user/address",
		"data": map[string]interface{}{"address": address},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get user by address: status %d, body: %s", resp.StatusCode, string(body))), nil
	}
	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}
	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(decoded)), nil
}

func (t *InspectStateTool) UserBalance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return mcp.NewToolResultError("address must be a string"), nil
	}
	payload := map[string]interface{}{
		"path": "user/balance",
		"data": map[string]interface{}{"address": address},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get user balance: status %d, body: %s", resp.StatusCode, string(body))), nil
	}
	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}
	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(decoded)), nil
}

// Contract MCP inspect handlers
func (t *InspectStateTool) ListAllContracts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	payload := map[string]interface{}{
		"path": "contract",
		"data": map[string]interface{}{},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get contracts: status %d, body: %s", resp.StatusCode, string(body))), nil
	}
	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}
	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(decoded)), nil
}

func (t *InspectStateTool) FindContractBySymbol(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol, ok := request.Params.Arguments["symbol"].(string)
	if !ok {
		return mcp.NewToolResultError("symbol must be a string"), nil
	}
	payload := map[string]interface{}{
		"path": "contract/symbol",
		"data": map[string]interface{}{"symbol": symbol},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get contract by symbol: status %d, body: %s", resp.StatusCode, string(body))), nil
	}
	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}
	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(decoded)), nil
}

func (t *InspectStateTool) FindContractByAddress(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return mcp.NewToolResultError("address must be a string"), nil
	}
	payload := map[string]interface{}{
		"path": "contract/address",
		"data": map[string]interface{}{"address": address},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.BaseURL, t.AppAddress.Hex()), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.Client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcp.NewToolResultError(fmt.Sprintf("failed to get contract by address: status %d, body: %s", resp.StatusCode, string(body))), nil
	}
	var response struct {
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(response.Reports) == 0 {
		return mcp.NewToolResultError("no reports in response"), nil
	}
	hexStr := strings.TrimPrefix(response.Reports[0].Payload, "0x")
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(decoded)), nil
}
