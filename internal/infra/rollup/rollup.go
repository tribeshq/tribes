package rollup

import (
	"log/slog"
	"os"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/configs"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/rollup/middleware"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/router"
)

type CreateInfo struct {
	Repo   repository.Repository
	Config *configs.RollupConfig
}

func Create(c *CreateInfo) *router.Router {
	handlers, err := NewHandlers(c.Repo, c.Config)
	if err != nil {
		slog.Error("Failed to initialize handlers", "error", err)
		os.Exit(1)
	}

	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)
	r.Use(router.ErrorHandlingMiddleware)

	rbacFactory := middleware.NewRBACFactory(c.Repo)

	orderInvestorGroup := r.Group("order")
	orderInvestorGroup.Use(rbacFactory.InvestorOnly())
	{
		// restricted operations
		orderInvestorGroup.HandleAdvance("create", handlers.OrderAdvanceHandlers.CreateOrder)
		orderInvestorGroup.HandleAdvance("cancel", handlers.OrderAdvanceHandlers.CancelOrder)

		// Public operations
		orderInvestorGroup.HandleInspect("", handlers.OrderInspectHandlers.FindAllOrders)
		orderInvestorGroup.HandleInspect("id", handlers.OrderInspectHandlers.FindOrderById)
		orderInvestorGroup.HandleInspect("issuance", handlers.OrderInspectHandlers.FindBidsByIssuanceId)
		orderInvestorGroup.HandleInspect("investor", handlers.OrderInspectHandlers.FindOrdersByInvestorAddress)
	}

	issuanceGroup := r.Group("issuance")
	issuanceCreatorGroup := issuanceGroup.Group("creator")
	issuanceCreatorGroup.Use(rbacFactory.CreatorOnly())
	{
		// restricted operations
		issuanceCreatorGroup.HandleAdvance("create", handlers.IssuanceAdvanceHandlers.CreateIssuance)
		issuanceCreatorGroup.HandleAdvance("settle", handlers.IssuanceAdvanceHandlers.SettleIssuance)

		// Public operations
		issuanceGroup.HandleInspect("", handlers.IssuanceInspectHandlers.FindAllIssuances)
		issuanceGroup.HandleInspect("id", handlers.IssuanceInspectHandlers.FindIssuanceById)
		issuanceGroup.HandleAdvance("close", handlers.IssuanceAdvanceHandlers.CloseIssuance)
		issuanceGroup.HandleInspect("creator", handlers.IssuanceInspectHandlers.FindIssuancesByCreatorAddress)
		issuanceGroup.HandleInspect("investor", handlers.IssuanceInspectHandlers.FindIssuancesByInvestorAddress)
		issuanceGroup.HandleAdvance("execute-collateral", handlers.IssuanceAdvanceHandlers.ExecuteIssuanceCollateral)
	}

	userGroup := r.Group("user")
	adminUserGroup := userGroup.Group("admin")
	adminUserGroup.Use(rbacFactory.AdminOnly())
	{
		// restricted operations
		adminUserGroup.HandleAdvance("create", handlers.UserAdvanceHandlers.CreateUser)
		adminUserGroup.HandleAdvance("delete", handlers.UserAdvanceHandlers.DeleteUser)
		adminUserGroup.HandleAdvance("emergency-erc20-withdraw", handlers.UserAdvanceHandlers.EmergencyERC20Withdraw)
		adminUserGroup.HandleAdvance("emergency-ether-withdraw", handlers.UserAdvanceHandlers.EmergencyEtherWithdraw)

		// Public operations
		userGroup.HandleInspect("", handlers.UserInspectHandlers.FindAllUsers)
		userGroup.HandleInspect("address", handlers.UserInspectHandlers.FindUserByAddress)
		userGroup.HandleInspect("balance", handlers.UserInspectHandlers.ERC20BalanceOf)
		userGroup.HandleAdvance("withdraw", handlers.UserAdvanceHandlers.ERC20Withdraw)
	}

	socialGroup := r.Group("social")

	verifierGroup := socialGroup.Group("verifier")
	verifierGroup.Use(rbacFactory.VerifierOnly())

	socialAdminGroup := socialGroup.Group("admin")
	socialAdminGroup.Use(rbacFactory.AdminOnly())
	{
		// restricted operations
		verifierGroup.HandleAdvance("create", handlers.SocialAccountsHandlers.CreateSocialAccount)
		socialAdminGroup.HandleAdvance("delete", handlers.SocialAccountsHandlers.DeleteSocialAccount)

		// Public operations
		socialGroup.HandleInspect("id", handlers.SocialAccountHandlers.FindSocialAccountById)
		socialGroup.HandleInspect("user/id", handlers.SocialAccountHandlers.FindSocialAccountsByUserId)
	}
	return r
}
