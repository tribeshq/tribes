package root

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tribeshq/tribes/configs"
)

const (
	CMD_NAME = "mcp"
)

var (
	cfg                    *configs.McpConfig
	erc20portalAddress     string
	tokenAddress           string
	appAddress             string
	blockchainHttpEndpoint string
	jsonrpcEndpoint        string
)

var Cmd = &cobra.Command{
	Use:   "tribes-" + CMD_NAME,
	Short: "Runs Tribes MCP",
	Long:  `MCP server for debt issuance through crowdfunding w/ collateralized tokenization of receivables`,
	Run:   run,
}

func init() {
	Cmd.Flags().StringVar(&erc20portalAddress, "erc20-portal-address", "0x05355c2F9bA566c06199DEb17212c3B78C1A3C31", "ERC20 Portal contract address")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_CONTRACTS_ERC20_PORTAL_ADDRESS, Cmd.Flags().Lookup("erc20-portal-address")))

	Cmd.Flags().StringVar(&tokenAddress, "token-address", "", "Token contract address")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_CONTRACTS_TOKEN_ADDRESS, Cmd.Flags().Lookup("token-address")))
	cobra.CheckErr(Cmd.MarkFlagRequired("token-address"))

	Cmd.Flags().StringVar(&appAddress, "app-address", "", "Application contract address")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_CONTRACTS_APPLICATION_ADDRESS, Cmd.Flags().Lookup("app-address")))
	cobra.CheckErr(Cmd.MarkFlagRequired("app-address"))

	Cmd.Flags().StringVar(&blockchainHttpEndpoint, "blockchain-http-endpoint", "http://localhost:8080/anvil", "Blockchain HTTP endpoint")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_BLOCKCHAIN_HTTP_ENDPOINT, Cmd.Flags().Lookup("blockchain-http-endpoint")))

	Cmd.Flags().StringVar(&jsonrpcEndpoint, "jsonrpc-endpoint", "http://127.0.0.1:8080/rpc", "Jsonrpc endpoint")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_JSONRPC_ENDPOINT, Cmd.Flags().Lookup("jsonrpc-endpoint")))

	Cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = configs.LoadMcpConfig()
		if err != nil {
			return err
		}
		return nil
	}
}

func run(cmd *cobra.Command, args []string) {
	startTime := time.Now()

	s := server.NewMCPServer(
		"Tribes MCP Server",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	ready := make(chan struct{}, 1)
	go func() {
		select {
		case <-ready:
			duration := time.Since(startTime)
			slog.Info("DApp is ready", "after", duration)
		case <-cmd.Context().Done():
		}
	}()

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
