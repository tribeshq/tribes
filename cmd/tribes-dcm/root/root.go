package root

import (
	"context"
	"log/slog"
	"os"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/configs"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository/factory"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/rollup"
	"github.com/rollmelette/rollmelette"
	"github.com/spf13/cobra"
)

const (
	CMD_NAME = "dcm"
)

var (
	cfg *configs.RollupConfig
)

var Cmd = &cobra.Command{
	Use:   "tribes-" + CMD_NAME,
	Short: "Runs Debt Capital Market as rollup",
	Long:  `A Linux-powered EVM rollup serving as a Debt Capital Market for the creator economy`,
	Run:   run,
}

func init() {
	Cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = configs.LoadRollupConfig()
		if err != nil {
			return err
		}
		return nil
	}
}

func run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MaxStartupTime)

	defer cancel()

	repo, err := factory.NewRepositoryFromConnectionString(ctx, "sqlite:///mnt/data/rollup.db")
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	defer repo.Close()

	createInfo := &rollup.CreateInfo{
		Repo:   repo,
		Config: cfg,
	}

	r := rollup.Create(createInfo)
	opts := rollmelette.NewRunOpts()
	if err := rollmelette.Run(cmd.Context(), opts, r); err != nil {
		slog.Error("Failed to run rollmelette", "error", err)
		os.Exit(1)
	}
}
