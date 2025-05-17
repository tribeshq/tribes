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
