package root

import (
	"log/slog"
	"os"
	"time"

	"github.com/rollmelette/rollmelette"
	"github.com/spf13/cobra"
	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/pkg/router"
	"gorm.io/gorm"
)

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

// *************************************************************************************

const (
	CMD_NAME = "tribes-rollup"
)

var (
	verbose     bool
	useMemoryDB bool
	Cmd         = &cobra.Command{
		Use:   CMD_NAME,
		Short: "Runs Tribes Rollup",
		Long:  `Debt issuance through crowdfunding w/ collateralized tokenization of receivables`,
		Run:   run,
	}
)

func init() {
	Cmd.Flags().BoolVar(&verbose, "verbose", false, "Show detailed logs")

	Cmd.PersistentFlags().BoolVar(
		&useMemoryDB,
		"memory-db",
		false,
		"Use in-memory SQLite database instead of persistent",
	)

	Cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if verbose {
			configs.ConfigureLogger(slog.LevelDebug)
		} else {
			configs.ConfigureLogger(slog.LevelInfo)
		}
	}
}

func run(cmd *cobra.Command, args []string) {
	startTime := time.Now()

	ctx := cmd.Context()

	var db *gorm.DB
	var err error
	if useMemoryDB {
		db, err = configs.SetupSQlite(":memory:")
		if err != nil {
			slog.Error("Failed to setup in-memory SQLite database", "error", err)
			os.Exit(1)
		}
		slog.Info("In-memory database initialized")
	} else {
		db, err = configs.SetupSQlite("tribes.db")
		if err != nil {
			slog.Error("Failed to setup SQLite database", "error", err)
			os.Exit(1)
		}
		slog.Info("Persistent database initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get SQL DB from GORM", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	ah, err := NewAdvanceHandlers(db)
	if err != nil {
		slog.Error("Failed to initialize advance handlers", "error", err)
		os.Exit(1)
	}
	slog.Info("Advance handlers initialized")

	ih, err := NewInspectHandlers(db)
	if err != nil {
		slog.Error("Failed to initialize inspect handlers", "error", err)
		os.Exit(1)
	}
	slog.Info("Inspect handlers initialized")

	ms, err := NewMiddlewares(db)
	if err != nil {
		slog.Error("Failed to initialize middlewares", "error", err)
		os.Exit(1)
	}
	slog.Info("Middlewares initialized")

	r := NewDApp(ah, ih, ms)
	slog.Info("Router initialized")

	opts := rollmelette.NewRunOpts()
	if rollupUrl, isSet := os.LookupEnv("ROLLUP_HTTP_SERVER_URL"); isSet {
		opts.RollupURL = rollupUrl
	}

	ready := make(chan struct{}, 1)
	go func() {
		select {
		case <-ready:
			duration := time.Since(startTime)
			slog.Info("DApp is ready", "after", duration)
		case <-ctx.Done():
		}
	}()

	if err := rollmelette.Run(ctx, opts, r); err != nil {
		slog.Error("Application exited with an error", "error", err)
		os.Exit(1)
	}
}

func NewDApp(ah *AdvanceHandlers, ih *InspectHandlers, ms *Middlewares) *router.Router {
	r := router.NewRouter()

	r.HandleAdvance("create_contract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.CreateContractHandler, []string{"admin"}))
	r.HandleAdvance("update_contract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.UpdateContractHandler, []string{"admin"}))
	r.HandleAdvance("delete_contract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.DeleteContractHandler, []string{"admin"}))

	r.HandleAdvance("create_order", ms.RBAC.Middleware(ah.OrderAdvanceHandlers.CreateOrderHandler, []string{"non_qualified_investor", "qualified_investor"}))
	r.HandleAdvance("cancel_order", ms.RBAC.Middleware(ah.OrderAdvanceHandlers.CancelOrderHandler, []string{"non_qualified_investor", "qualified_investor"}))

	r.HandleAdvance("create_crowdfunding", ms.RBAC.Middleware(ah.CrowdfundingAdvanceHandlers.CreateCrowdfundingHandler, []string{"creator"}))
	r.HandleAdvance("delete_crowdfunding", ms.RBAC.Middleware(ah.CrowdfundingAdvanceHandlers.DeleteCrowdfundingHandler, []string{"admin"}))
	r.HandleAdvance("update_crowdfunding", ms.RBAC.Middleware(ah.CrowdfundingAdvanceHandlers.UpdateCrowdfundingHandler, []string{"admin"}))
	r.HandleAdvance("close_crowdfunding", ah.CrowdfundingAdvanceHandlers.CloseCrowdfundingHandler)
	r.HandleAdvance("settle_crowdfunding", ms.RBAC.Middleware(ah.CrowdfundingAdvanceHandlers.SettleCrowdfundingHandler, []string{"creator"}))

	r.HandleAdvance("create_user", ms.RBAC.Middleware(ah.UserAdvanceHandlers.CreateUserHandler, []string{"admin"}))
	r.HandleAdvance("update_user", ms.RBAC.Middleware(ah.UserAdvanceHandlers.UpdateUserHandler, []string{"admin"}))
	r.HandleAdvance("delete_user", ms.RBAC.Middleware(ah.UserAdvanceHandlers.DeleteUserHandler, []string{"admin"}))
	r.HandleAdvance("withdraw", ah.UserAdvanceHandlers.WithdrawHandler)

	r.HandleAdvance("create_social_account", ms.RBAC.Middleware(ah.SocialAccountsHandlers.CreateSocialAccountHandler, []string{"admin"}))
	r.HandleAdvance("delete_social_account", ms.RBAC.Middleware(ah.SocialAccountsHandlers.DeleteSocialAccountHandler, []string{"admin"}))

	r.HandleInspect("find_crowdfunding", ih.CrowdfundingInspectHandlers.FindAllCrowdfundingsHandler)
	r.HandleInspect("find_crowdfunding_by_id", ih.CrowdfundingInspectHandlers.FindCrowdfundingByIdHandler)
	r.HandleInspect("find_crowdfunding_by_creator", ih.CrowdfundingInspectHandlers.FindCrowdfundingsByCreatorHandler)
	r.HandleInspect("find_crowdfunding_by_investor", ih.CrowdfundingInspectHandlers.FindCrowdfundingsByInvestorHandler)

	r.HandleInspect("find_order", ih.OrderInspectHandlers.FindAllOrdersHandler)
	r.HandleInspect("find_order_by_id", ih.OrderInspectHandlers.FindOrderByIdHandler)
	r.HandleInspect("find_order_by_investor", ih.OrderInspectHandlers.FindOrdersByInvestorHandler)
	r.HandleInspect("find_order_by_crowdfunding", ih.OrderInspectHandlers.FindBisdByCrowdfundingIdHandler)

	r.HandleInspect("find_contract", ih.ContractInspectHandlers.FindAllContractsHandler)
	r.HandleInspect("find_contract_by_symbol", ih.ContractInspectHandlers.FindContractBySymbolHandler)
	r.HandleInspect("find_contract_by_address", ih.ContractInspectHandlers.FindContractByAddressHandler)

	r.HandleInspect("find_all_users", ih.UserInspectHandlers.FindAllUsersHandler)
	r.HandleInspect("find_user_by_address", ih.UserInspectHandlers.FindUserByAddressHandler)
	r.HandleInspect("find_balance_by_address", ih.UserInspectHandlers.BalanceHandler)

	r.HandleInspect("find_social_account_by_id", ih.SocialAccountHandlers.FindSocialAccountById)
	r.HandleInspect("find_social_account_by_user", ih.SocialAccountHandlers.FindSocialAccountsByUserId)

	return r
}
