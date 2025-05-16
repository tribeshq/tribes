package root

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

const (
	CMD_NAME = "mcp"
)

var (
	useMemoryDB bool
	Cmd         = &cobra.Command{
		Use:   "tribes-" + CMD_NAME,
		Short: "Runs Tribes MCP",
		Long:  `MCP server for debt issuance through crowdfunding w/ collateralized tokenization of receivables`,
		Run:   run,
	}
)

func init() {
	Cmd.PersistentFlags().BoolVar(
		&useMemoryDB,
		"memory-db",
		false,
		"Use in-memory SQLite database instead of persistent",
	)
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
