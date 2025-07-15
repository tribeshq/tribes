package cartesi

import (
	"log/slog"
	"os"

	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/internal/infra/cartesi/middleware"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/router"
)

type CreateInfo struct {
	Repo   repository.Repository
	Config *configs.RollupConfig
}

func Create(info *CreateInfo) *router.Router {
	handlers, err := NewHandlers(info.Repo, info.Config)
	if err != nil {
		slog.Error("Failed to initialize handlers", "error", err)
		os.Exit(1)
	}
	slog.Info("Handlers initialized")

	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)
	r.Use(router.ErrorHandlingMiddleware)

	rbacFactory := middleware.NewRBACFactory(info.Repo)

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
