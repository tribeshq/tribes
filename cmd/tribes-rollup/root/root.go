// *************************************************************************************
// *                           PLATFORM FUNCTIONAL REQUIREMENTS                        *
// *************************************************************************************

// 1. Registration of Small Business Entities:
//    1.1 Ensure that entities are legally constituted and meet the regulatory requirements.

// 2. Management of Public Offerings:
//    2.1 Define minimum and maximum target amounts for fundraising (maximum of R$ 15 million).
//    2.2 Set fundraising duration to no more than 180 days.
//    2.3 Guarantee a 5-day withdrawal period for investors after confirming their participation.
//    2.4 A company must wait 120 days after the close of a successful crowdfunding campaign
//        before starting a new campaign.

// 3. Investor Control:
//    3.1 Verify investor profiles (e.g., lead investors or qualified investors).
//    3.2 Limit the annual investment amount to R$ 20,000, except for higher income or qualified investors.

// 4. Publication of Essential Information:
//    4.1 Maintain a dedicated page for offers with clear and objective investment details.
//    4.2 Publish relevant documents, such as:
//        4.2.1 The corporate charter.
//        4.2.2 Investment agreements.
//        4.2.3 Financial statements.

// 5. Investment Processing:
//    5.1 Transfer collected funds directly to the small business's accounts after the offer closes.
//    5.2 Prohibit fund transit through accounts linked to the platform or its stakeholders.

// 6. Reporting and Audits:
//    6.1 Provide monthly reports on transaction volumes and prices.
//    6.2 Ensure financial statements are audited for offerings above R$ 10 million.

// 7. Promotion and Disclosure:
//    7.1 Allow wide promotion with content and language restrictions.
//    7.2 Enable events and interactions with investors, adhering to regulatory guidelines.

// 8. Intermediation of Subsequent Transactions:
//    8.1 Ensure secure transfer of security ownership.
//    8.2 Support buying and selling of already issued securities when authorized.

// 9. Regulatory Compliance:
//    9.1 Fulfill CVM registration requirements, including a minimum capital of R$ 200,000.
//    9.2 Develop a code of conduct addressing conflicts of interest for partners and administrators.

package root

import (
	"log/slog"
	"os"

	"github.com/rollmelette/rollmelette"
	"github.com/spf13/cobra"
	"github.com/tribeshq/tribes/internal/infra/cartesi/middleware"
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
		Long:  `Cartesi Rollup Application for debt issuance through crowdfunding w/ collateralized tokenization of receivables`,
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
		cmd.Context(),
		map[bool]string{true: "sqlite://:memory:", false: "sqlite://tribes.db"}[useMemoryDB],
	)
	if err != nil {
		slog.Error("Failed to setup database", "error", err, "type", map[bool]string{true: "in-memory", false: "persistent"}[useMemoryDB])
		os.Exit(1)
	}
	slog.Info("Database initialized", "type", map[bool]string{true: "in-memory", false: "persistent"}[useMemoryDB])

	defer repo.Close()
	handlers, err := NewHandlers(repo)
	if err != nil {
		slog.Error("Failed to initialize handlers", "error", err)
		os.Exit(1)
	}
	slog.Info("Handlers initialized")

	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)
	r.Use(router.ValidationMiddleware)
	r.Use(router.ErrorHandlingMiddleware)

	rbacFactory := middleware.NewRBACFactory(repo)

	contractGroup := r.Group("contract")
	{
		contractGroup.Use(rbacFactory.AdminOnly())
		contractGroup.HandleAdvance("create", handlers.ContractAdvanceHandlers.CreateContract)
		contractGroup.HandleAdvance("update", handlers.ContractAdvanceHandlers.UpdateContract)
		contractGroup.HandleAdvance("delete", handlers.ContractAdvanceHandlers.DeleteContract)

		// Public operations
		contractGroup.HandleInspect("", handlers.ContractInspectHandlers.FindAllContracts)
		contractGroup.HandleInspect("symbol", handlers.ContractInspectHandlers.FindContractBySymbol)
		contractGroup.HandleInspect("address", handlers.ContractInspectHandlers.FindContractByAddress)
	}

	orderGroup := r.Group("order")
	{
		orderGroup.Use(rbacFactory.InvestorOnly())
		orderGroup.HandleAdvance("create", handlers.OrderAdvanceHandlers.CreateOrder)
		orderGroup.HandleAdvance("cancel", handlers.OrderAdvanceHandlers.CancelOrder)

		// Public operations
		orderGroup.HandleInspect("", handlers.OrderInspectHandlers.FindAllOrders)
		orderGroup.HandleInspect("id", handlers.OrderInspectHandlers.FindOrderById)
		orderGroup.HandleInspect("investor", handlers.OrderInspectHandlers.FindOrdersByInvestor)
		orderGroup.HandleInspect("crowdfunding", handlers.OrderInspectHandlers.FindBisdByCrowdfundingId)
	}

	crowdfundingGroup := r.Group("crowdfunding")
	{
		adminGroup := crowdfundingGroup.Group("admin")

		adminGroup.Use(rbacFactory.AdminOnly())
		adminGroup.HandleAdvance("delete", handlers.CrowdfundingAdvanceHandlers.DeleteCrowdfunding)
		adminGroup.HandleAdvance("update", handlers.CrowdfundingAdvanceHandlers.UpdateCrowdfunding)

		creatorGroup := crowdfundingGroup.Group("creator")
		creatorGroup.Use(rbacFactory.CreatorOnly())
		creatorGroup.HandleAdvance("create", handlers.CrowdfundingAdvanceHandlers.CreateCrowdfunding)
		creatorGroup.HandleAdvance("settle", handlers.CrowdfundingAdvanceHandlers.SettleCrowdfunding)

		// Public operations
		crowdfundingGroup.HandleAdvance("close", handlers.CrowdfundingAdvanceHandlers.CloseCrowdfunding)
		crowdfundingGroup.HandleInspect("", handlers.CrowdfundingInspectHandlers.FindAllCrowdfundings)
		crowdfundingGroup.HandleInspect("id", handlers.CrowdfundingInspectHandlers.FindCrowdfundingById)
		crowdfundingGroup.HandleInspect("creator", handlers.CrowdfundingInspectHandlers.FindCrowdfundingsByCreator)
		crowdfundingGroup.HandleInspect("investor", handlers.CrowdfundingInspectHandlers.FindCrowdfundingsByInvestor)
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
		userGroup.HandleInspect("balance", handlers.UserInspectHandlers.Balance)
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

	opts := rollmelette.NewRunOpts()
	if err := rollmelette.Run(cmd.Context(), opts, r); err != nil {
		slog.Error("Failed to run rollmelette", "error", err)
		os.Exit(1)
	}
}
