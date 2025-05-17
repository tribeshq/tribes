package root

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/configs/auth"
	"github.com/tribeshq/tribes/internal/infra/mcp/tool"
)

const (
	CMD_NAME = "mcp"
)

var (
	cfg                    *configs.McpConfig
	authKind               string
	appAddress             string
	privateKey             string
	tokenAddress           string
	mnemonicPhrase         string
	mnemonicIndex          int
	jsonrpcEndpoint        string
	inspectEndpoint        string
	erc20portalAddress     string
	blockchainHttpEndpoint string
	blockchainId           int
)

var Cmd = &cobra.Command{
	Use:   "tribes-" + CMD_NAME,
	Short: "Runs Tribes MCP",
	Long:  `MCP server for debt issuance through crowdfunding w/ collateralized tokenization of receivables`,
	Run:   run,
}

func init() {
	Cmd.Flags().StringVar(&tokenAddress, "token-address", "", "Token contract address")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_CONTRACTS_TOKEN_ADDRESS, Cmd.Flags().Lookup("token-address")))
	cobra.CheckErr(Cmd.MarkFlagRequired("token-address"))

	Cmd.Flags().StringVar(&appAddress, "app-address", "", "Application contract address")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_CONTRACTS_APPLICATION_ADDRESS, Cmd.Flags().Lookup("app-address")))
	cobra.CheckErr(Cmd.MarkFlagRequired("app-address"))

	Cmd.Flags().StringVar(&authKind, "auth-kind", "", "Authentication kind")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_AUTH_KIND, Cmd.Flags().Lookup("auth-kind")))

	Cmd.Flags().StringVar(&privateKey, "private-key", "", "Private key")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_AUTH_PRIVATE_KEY, Cmd.Flags().Lookup("private-key")))

	Cmd.Flags().StringVar(&mnemonicPhrase, "mnemonic-phrase", "", "Mnemonic phrase")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_AUTH_MNEMONIC, Cmd.Flags().Lookup("mnemonic-phrase")))

	Cmd.Flags().IntVar(&mnemonicIndex, "mnemonic-index", 0, "Mnemonic index")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_AUTH_MNEMONIC_ACCOUNT_INDEX, Cmd.Flags().Lookup("mnemonic-index")))

	Cmd.Flags().StringVar(&inspectEndpoint, "inspect-endpoint", "http://127.0.0.1:8080/inspect", "Inspect endpoint")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_INSPECT_ENDPOINT, Cmd.Flags().Lookup("inspect-endpoint")))

	Cmd.Flags().StringVar(&erc20portalAddress, "erc20-portal-address", "0x05355c2F9bA566c06199DEb17212c3B78C1A3C31", "ERC20 Portal contract address")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_CONTRACTS_ERC20_PORTAL_ADDRESS, Cmd.Flags().Lookup("erc20-portal-address")))

	Cmd.Flags().IntVar(&blockchainId, "blockchain-id", 13370, "Blockchain ID")
	cobra.CheckErr(viper.BindPFlag(configs.TRIBES_BLOCKCHAIN_ID, Cmd.Flags().Lookup("blockchain-id")))

	Cmd.Flags().StringVar(&blockchainHttpEndpoint, "blockchain-http-endpoint", "http://127.0.0.1:8080/anvil", "Blockchain HTTP endpoint")
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
	ctx := cmd.Context()
	s := server.NewMCPServer(
		"Tribes MCP Server",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	//--------------------------------
	// MCP Tools
	//--------------------------------

	listAllCrowdfundings := mcp.NewTool("list_all_crowdfundings",
		mcp.WithDescription("List all crowdfundings in the system by creators"),
	)

	listAllOrders := mcp.NewTool("list_all_orders",
		mcp.WithDescription("List all orders in the system"),
	)

	listCrowdfundingByCreator := mcp.NewTool("list_crowdfunding_by_creator",
		mcp.WithDescription("List all crowdfundings created by a specific creator"),
		mcp.WithString("creator",
			mcp.Required(),
			mcp.Description("Address of the creator"),
		),
	)

	createOrder := mcp.NewTool("create_order",
		mcp.WithDescription("Create a new order for a crowdfunding"),
		mcp.WithString("crowdfunding_id",
			mcp.Required(),
			mcp.Description("ID of the crowdfunding"),
		),
		mcp.WithString("amount",
			mcp.Required(),
			mcp.Description("Amount to invest"),
		),
		mcp.WithString("interest_rate",
			mcp.Required(),
			mcp.Description("Interest rate for the investment"),
		),
	)

	//--------------------------------
	// Setup signer
	//--------------------------------

	client, err := ethclient.DialContext(ctx, blockchainHttpEndpoint)
	if err != nil {
		slog.Error("Failed to connect to blockchain",
			"error", err,
			"endpoint", blockchainHttpEndpoint,
		)
		cobra.CheckErr(err)
	}

	chainId, err := client.ChainID(ctx)
	if err != nil {
		slog.Error("Failed to get chain ID", "error", err)
		cobra.CheckErr(err)
	}

	txOpts, err := auth.GetTransactOpts(chainId)
	if err != nil {
		slog.Error("Failed to configure transaction options", "error", err)
		cobra.CheckErr(err)
	}

	//--------------------------------
	// MCP Handlers
	//--------------------------------

	inspectStateTool := tool.NewInspectStateTool(
		inspectEndpoint,
		common.HexToAddress(appAddress),
	)

	advanceStateTool := tool.NewAdvanceStateTool(
		client,
		txOpts,
		common.HexToAddress(appAddress),
		common.HexToAddress(tokenAddress),
		common.HexToAddress(erc20portalAddress),
	)

	//--------------------------------
	// Setup MCP server
	//--------------------------------

	s.AddTool(listAllCrowdfundings, inspectStateTool.ListAllCrowdfundings)
	s.AddTool(listAllOrders, inspectStateTool.ListAllOrders)
	s.AddTool(listCrowdfundingByCreator, inspectStateTool.ListCrowdfundingByCreator)
	s.AddTool(createOrder, advanceStateTool.CreateOrder)

	slog.Info("Starting MCP server",
		"blockchainId", blockchainId,
		"blockchainHttpEndpoint", blockchainHttpEndpoint,
		"inspectEndpoint", inspectEndpoint,
		"jsonrpcEndpoint", jsonrpcEndpoint,
		"appAddress", appAddress,
		"tokenAddress", tokenAddress,
	)

	if err := server.ServeStdio(s); err != nil {
		slog.Error("Server error",
			"error", err,
			"endpoint", jsonrpcEndpoint,
		)
		fmt.Printf("Server error: %v\n", err)
	}
}
