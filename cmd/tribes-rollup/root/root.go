package root

import (
	"log/slog"
	"os"

	"github.com/rollmelette/rollmelette"
	"github.com/spf13/cobra"
	"github.com/tribeshq/tribes/internal/infra/cartesi/middleware"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/infra/repository/factory"
	"github.com/tribeshq/tribes/pkg/rollups/router"
)

const (
	CMD_NAME = "rollup"
)

var (
	useMemoryDB bool
	Cmd         = &cobra.Command{
		Use:   "tribes-" + CMD_NAME,
		Short: "Runs Tribes Rollup",
		Long:  `Cartesi Rollup Application for debt issuance through auction w/ collateralized tokenization of receivables`,
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
		map[bool]string{true: "sqlite://:memory:", false: "sqlite://tribes.db"}[useMemoryDB],
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
		orderGroup.HandleInspect("investor", handlers.OrderInspectHandlers.FindOrdersByInvestor)
		orderGroup.HandleInspect("auction", handlers.OrderInspectHandlers.FindBidsByAuctionId)
	}

	auctionGroup := r.Group("auction")
	{
		creatorGroup := auctionGroup.Group("creator")
		creatorGroup.Use(rbacFactory.CreatorOnly())
		creatorGroup.HandleAdvance("create", handlers.AuctionAdvanceHandlers.CreateAuction)
		creatorGroup.HandleAdvance("settle", handlers.AuctionAdvanceHandlers.SettleAuction)

		// Public operations
		auctionGroup.HandleAdvance("execute-collateral", handlers.AuctionAdvanceHandlers.ExecuteAuctionCollateral)
		auctionGroup.HandleAdvance("close", handlers.AuctionAdvanceHandlers.CloseAuction)
		auctionGroup.HandleInspect("", handlers.AuctionInspectHandlers.FindAllAuctions)
		auctionGroup.HandleInspect("id", handlers.AuctionInspectHandlers.FindAuctionById)
		auctionGroup.HandleInspect("creator", handlers.AuctionInspectHandlers.FindAuctionsByCreator)
		auctionGroup.HandleInspect("investor", handlers.AuctionInspectHandlers.FindAuctionsByInvestor)
	}

	userGroup := r.Group("user")
	{
		userGroup.Use(rbacFactory.AdminOnly())
		userGroup.HandleAdvance("create", handlers.UserAdvanceHandlers.CreateUser)
		userGroup.HandleAdvance("update", handlers.UserAdvanceHandlers.UpdateUser)
		userGroup.HandleAdvance("delete", handlers.UserAdvanceHandlers.DeleteUser)
		userGroup.HandleAdvance("withdraw", handlers.UserAdvanceHandlers.Withdraw)

		// Public operations
		userGroup.HandleInspect("", handlers.UserInspectHandlers.FindAllUsers)
		userGroup.HandleInspect("address", handlers.UserInspectHandlers.FindUserByAddress)
	}

	socialGroup := r.Group("social")
	{
		socialGroup.Use(rbacFactory.AdminOnly())
		socialGroup.HandleAdvance("create", handlers.SocialAccountsHandlers.CreateSocialAccount)
		socialGroup.HandleAdvance("delete", handlers.SocialAccountsHandlers.DeleteSocialAccount)

		// Public operations
		socialGroup.HandleInspect("id", handlers.SocialAccountHandlers.FindSocialAccountById)
		socialGroup.HandleInspect("user/id", handlers.SocialAccountHandlers.FindSocialAccountsByUserId)
	}

	return r
}
