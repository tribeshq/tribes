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
	"github.com/spf13/viper"
)

const (
	CMD_NAME = "dcm"
)

var (
	maxStartupTime         int
	databaseUrl            string
	adminAddress           string
	adminAddressTest       string
	verifierAddress        string
	verifierAddressTest    string
	badgeFactoryAddress    string
	emergencyWithdrawAddr  string
	safeErc1155MintAddress string
	cfg                    *configs.RollupConfig
)

var Cmd = &cobra.Command{
	Use:   "tribes-" + CMD_NAME,
	Short: "Runs Debt Capital Market as rollup",
	Long:  `A Linux-powered EVM rollup serving as a Debt Capital Market for the creator economy`,
	Run:   run,
}

func init() {
	// Rollup flags
	Cmd.Flags().IntVar(&maxStartupTime, "max-startup-time", 10, "Maximum startup time in seconds")
	cobra.CheckErr(viper.BindPFlag(configs.MAX_STARTUP_TIME, Cmd.Flags().Lookup("max-startup-time")))

	// Database flags
	Cmd.Flags().StringVar(&databaseUrl, "database-url", "sqlite:///mnt/data/rollup.db", "SQLite database connection string")
	cobra.CheckErr(viper.BindPFlag(configs.DATABASE_URL, Cmd.Flags().Lookup("database-url")))

	// Contracts flags
	Cmd.Flags().StringVar(&adminAddress, "admin-address", "", "Address of the admin user")
	cobra.CheckErr(viper.BindPFlag(configs.ADMIN_ADDRESS, Cmd.Flags().Lookup("admin-address")))

	Cmd.Flags().StringVar(&verifierAddress, "verifier-address", "", "Address of the verifier contract")
	cobra.CheckErr(viper.BindPFlag(configs.VERIFIER_ADDRESS, Cmd.Flags().Lookup("verifier-address")))

	Cmd.Flags().StringVar(&badgeFactoryAddress, "badge-factory-address", "", "Address of the badge factory contract")
	cobra.CheckErr(viper.BindPFlag(configs.BADGE_FACTORY_ADDRESS, Cmd.Flags().Lookup("badge-factory-address")))

	Cmd.Flags().StringVar(&emergencyWithdrawAddr, "emergency-withdraw-address", "", "Address for emergency withdrawals")
	cobra.CheckErr(viper.BindPFlag(configs.EMERGENCY_WITHDRAW_ADDRESS, Cmd.Flags().Lookup("emergency-withdraw-address")))

	Cmd.Flags().StringVar(&safeErc1155MintAddress, "safe-erc1155-mint-address", "", "Address for safe ERC1155 minting")
	cobra.CheckErr(viper.BindPFlag(configs.SAFE_ERC1155_MINT_ADDRESS, Cmd.Flags().Lookup("safe-erc1155-mint-address")))

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

	repo, err := factory.NewRepositoryFromConnectionString(ctx, cfg.DatabaseUrl)
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
