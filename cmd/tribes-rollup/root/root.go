package root

import (
	"log/slog"
	"os"

	"github.com/rollmelette/rollmelette"
	"github.com/spf13/cobra"
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
		orderGroup.HandleInspect("auction", handlers.OrderInspectHandlers.FindBidsByAuctionId)
		orderGroup.HandleInspect("investor", handlers.OrderInspectHandlers.FindOrdersByInvestor)
	}

	auctionGroup := r.Group("auction")
	{
		creatorGroup := auctionGroup.Group("creator")
		creatorGroup.Use(rbacFactory.CreatorOnly())
		creatorGroup.HandleAdvance("create", handlers.AuctionAdvanceHandlers.CreateAuction)
		creatorGroup.HandleAdvance("settle", handlers.AuctionAdvanceHandlers.SettleAuction)

		// Public operations
		auctionGroup.HandleInspect("", handlers.AuctionInspectHandlers.FindAllAuctions)
		auctionGroup.HandleInspect("id", handlers.AuctionInspectHandlers.FindAuctionById)
		auctionGroup.HandleAdvance("close", handlers.AuctionAdvanceHandlers.CloseAuction)
		auctionGroup.HandleInspect("creator", handlers.AuctionInspectHandlers.FindAuctionsByCreator)
		auctionGroup.HandleInspect("investor", handlers.AuctionInspectHandlers.FindAuctionsByInvestor)
		auctionGroup.HandleAdvance("execute-collateral", handlers.AuctionAdvanceHandlers.ExecuteAuctionCollateral)
	}

	userGroup := r.Group("user")
	{
		adminGroup := userGroup.Group("admin")
		adminGroup.Use(rbacFactory.AdminOnly())
		adminGroup.HandleAdvance("create", handlers.UserAdvanceHandlers.CreateUser)
		adminGroup.HandleAdvance("update", handlers.UserAdvanceHandlers.UpdateUser)
		adminGroup.HandleAdvance("delete", handlers.UserAdvanceHandlers.DeleteUser)

		// Public operations
		userGroup.HandleInspect("", handlers.UserInspectHandlers.FindAllUsers)
		userGroup.HandleInspect("address", handlers.UserInspectHandlers.FindUserByAddress)
		userGroup.HandleInspect("erc20-balance", handlers.UserInspectHandlers.ERC20BalanceOf)
		userGroup.HandleAdvance("erc20-withdraw", handlers.UserAdvanceHandlers.ERC20Withdraw)
	}

	socialGroup := r.Group("social")
	{
		creatorGroup := socialGroup.Group("creator")
		creatorGroup.Use(rbacFactory.CreatorOnly())
		creatorGroup.HandleAdvance("create", handlers.SocialAccountsHandlers.CreateSocialAccount)
		creatorGroup.HandleAdvance("delete", handlers.SocialAccountsHandlers.DeleteSocialAccount)

		// Public operations
		socialGroup.HandleInspect("id", handlers.SocialAccountHandlers.FindSocialAccountById)
		socialGroup.HandleInspect("user/id", handlers.SocialAccountHandlers.FindSocialAccountsByUserId)
	}
	return r
}
