package root

import (
	"log/slog"
	"os"

	"github.com/rollmelette/rollmelette"
	"github.com/spf13/cobra"
	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/internal/infra/cartesi/middleware"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/infra/repository/factory"
	"github.com/tribeshq/tribes/pkg/router"
)

const (
	CMD_NAME = "rollup"
)

var (
	useMemoryDB bool
	Cmd         = &cobra.Command{
		Use:   "tribes-" + CMD_NAME,
		Short: "Runs Tribes Rollup",
		Long:  `A Linux-powered EVM rollup serving as a Debt Capital Market for the creator economy`,
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
	Cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		_, err := configs.LoadRollupConfig()
		if err != nil {
			return err
		}
		return nil
	}
}

func run(cmd *cobra.Command, args []string) {
	repo, err := factory.NewRepositoryFromConnectionString(
		map[bool]string{true: "sqlite://:memory:", false: "sqlite:///mnt/data/tribes.db"}[useMemoryDB],
	)
	if err != nil {
		slog.Error("Failed to setup database", "error", err, "type", map[bool]string{true: "in-memory", false: "persistent"}[useMemoryDB])
		os.Exit(1)
	}
	slog.Info("Database initialized", "type", map[bool]string{true: "in-memory", false: "persistent"}[useMemoryDB])

	defer repo.Close()

	r := NewTribesRollup(repo)
	opts := rollmelette.NewRunOpts()
	if err := rollmelette.Run(cmd.Context(), opts, r); err != nil {
		slog.Error("Failed to run rollmelette", "error", err)
		os.Exit(1)
	}
}

func NewTribesRollup(repo repository.Repository) *router.Router {
	handlers, err := NewHandlers(repo)
	if err != nil {
		slog.Error("Failed to initialize handlers", "error", err)
		os.Exit(1)
	}
	slog.Info("Handlers initialized")

	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)
	r.Use(router.ErrorHandlingMiddleware)

	rbacFactory := middleware.NewRBACFactory(repo)

	orderGroup := r.Group("order")
	{
		orderGroup.Use(rbacFactory.InvestorOnly())
		orderGroup.HandleAdvance("create", handlers.OrderAdvanceHandlers.CreateOrder)
		orderGroup.HandleAdvance("cancel", handlers.OrderAdvanceHandlers.CancelOrder)

		// Public operations
		orderGroup.HandleInspect("", handlers.OrderInspectHandlers.FindAllOrders)
		orderGroup.HandleInspect("id", handlers.OrderInspectHandlers.FindOrderById)
		orderGroup.HandleInspect("campaign", handlers.OrderInspectHandlers.FindBidsByCampaignId)
		orderGroup.HandleInspect("investor", handlers.OrderInspectHandlers.FindOrdersByInvestorAddress)
	}

	campaignGroup := r.Group("campaign")
	{
		creatorGroup := campaignGroup.Group("creator")
		creatorGroup.Use(rbacFactory.CreatorOnly())
		creatorGroup.HandleAdvance("create", handlers.CampaignAdvanceHandlers.CreateCampaign)
		creatorGroup.HandleAdvance("settle", handlers.CampaignAdvanceHandlers.SettleCampaign)

		// Public operations
		campaignGroup.HandleInspect("", handlers.CampaignInspectHandlers.FindAllCampaigns)
		campaignGroup.HandleInspect("id", handlers.CampaignInspectHandlers.FindCampaignById)
		campaignGroup.HandleAdvance("close", handlers.CampaignAdvanceHandlers.CloseCampaign)
		campaignGroup.HandleInspect("creator", handlers.CampaignInspectHandlers.FindCampaignsByCreatorAddress)
		campaignGroup.HandleInspect("investor", handlers.CampaignInspectHandlers.FindCampaignsByInvestorAddress)
		campaignGroup.HandleAdvance("execute-collateral", handlers.CampaignAdvanceHandlers.ExecuteCampaignCollateral)
	}

	userGroup := r.Group("user")
	{
		adminGroup := userGroup.Group("admin")
		adminGroup.Use(rbacFactory.AdminOnly())
		adminGroup.HandleAdvance("create", handlers.UserAdvanceHandlers.CreateUser)
		adminGroup.HandleAdvance("delete", handlers.UserAdvanceHandlers.DeleteUser)
		adminGroup.HandleAdvance("emergency-erc20-withdraw", handlers.UserAdvanceHandlers.EmergencyERC20Withdraw)
		adminGroup.HandleAdvance("emergency-ether-withdraw", handlers.UserAdvanceHandlers.EmergencyEtherWithdraw)

		// Public operations
		userGroup.HandleInspect("", handlers.UserInspectHandlers.FindAllUsers)
		userGroup.HandleInspect("address", handlers.UserInspectHandlers.FindUserByAddress)
		userGroup.HandleInspect("balance", handlers.UserInspectHandlers.ERC20BalanceOf)
		userGroup.HandleAdvance("withdraw", handlers.UserAdvanceHandlers.ERC20Withdraw)
	}

	socialGroup := r.Group("social")
	{
		verifierGroup := socialGroup.Group("verifier")
		verifierGroup.Use(rbacFactory.VerifierOnly())
		verifierGroup.HandleAdvance("create", handlers.SocialAccountsHandlers.CreateSocialAccount)
		verifierGroup.HandleAdvance("delete", handlers.SocialAccountsHandlers.DeleteSocialAccount)

		// Public operations
		socialGroup.HandleInspect("id", handlers.SocialAccountHandlers.FindSocialAccountById)
		socialGroup.HandleInspect("user/id", handlers.SocialAccountHandlers.FindSocialAccountsByUserId)
	}
	return r
}
