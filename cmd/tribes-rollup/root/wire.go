//go:build wireinject
// +build wireinject

package root

import (
	"github.com/google/wire"
	"github.com/tribeshq/tribes/internal/infra/cartesi/handler/advance_handler"
	"github.com/tribeshq/tribes/internal/infra/cartesi/handler/inspect_handler"
	"github.com/tribeshq/tribes/internal/infra/repository"
)

// NewHandlers creates a new instance of all handlers with their dependencies
func NewHandlers(repo repository.Repository) (*Handlers, error) {
	wire.Build(
		// Bind repository interfaces
		wire.Bind(new(repository.UserRepository), new(repository.Repository)),
		wire.Bind(new(repository.OrderRepository), new(repository.Repository)),
		wire.Bind(new(repository.ContractRepository), new(repository.Repository)),
		wire.Bind(new(repository.CrowdfundingRepository), new(repository.Repository)),
		wire.Bind(new(repository.SocialAccountRepository), new(repository.Repository)),
		// Advance handlers
		advance_handler.NewOrderAdvanceHandlers,
		advance_handler.NewUserAdvanceHandlers,
		advance_handler.NewSocialAccountAdvanceHandlers,
		advance_handler.NewCrowdfundingAdvanceHandlers,
		advance_handler.NewContractAdvanceHandlers,
		// Inspect handlers
		inspect_handler.NewOrderInspectHandlers,
		inspect_handler.NewUserInspectHandlers,
		inspect_handler.NewSocialAccountInspectHandlers,
		inspect_handler.NewCrowdfundingInspectHandlers,
		inspect_handler.NewContractInspectHandlers,
		wire.Struct(new(Handlers), "*"),
	)
	return &Handlers{}, nil
}

// Handlers contains all handler dependencies
type Handlers struct {
	// Advance handlers
	OrderAdvanceHandlers        *advance_handler.OrderAdvanceHandlers
	UserAdvanceHandlers         *advance_handler.UserAdvanceHandlers
	SocialAccountsHandlers      *advance_handler.SocialAccountAdvanceHandlers
	CrowdfundingAdvanceHandlers *advance_handler.CrowdfundingAdvanceHandlers
	ContractAdvanceHandlers     *advance_handler.ContractAdvanceHandlers

	// Inspect handlers
	OrderInspectHandlers        *inspect_handler.OrderInspectHandlers
	UserInspectHandlers         *inspect_handler.UserInspectHandlers
	SocialAccountHandlers       *inspect_handler.SocialAccountInspectHandlers
	CrowdfundingInspectHandlers *inspect_handler.CrowdfundingInspectHandlers
	ContractInspectHandlers     *inspect_handler.ContractInspectHandlers
}
