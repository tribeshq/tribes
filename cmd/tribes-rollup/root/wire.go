//go:build wireinject
// +build wireinject

package root

import (
	"github.com/google/wire"
	"github.com/tribeshq/tribes/internal/infra/cartesi/handler/advance"
	"github.com/tribeshq/tribes/internal/infra/cartesi/handler/inspect"
	"github.com/tribeshq/tribes/internal/infra/repository"
)

func NewHandlers(repo repository.Repository) (*Handlers, error) {
	wire.Build(
		// Bind repository interfaces
		wire.Bind(new(repository.UserRepository), new(repository.Repository)),
		wire.Bind(new(repository.OrderRepository), new(repository.Repository)),
		wire.Bind(new(repository.AuctionRepository), new(repository.Repository)),
		wire.Bind(new(repository.SocialAccountRepository), new(repository.Repository)),
		// Advance handlers
		advance.NewOrderAdvanceHandlers,
		advance.NewUserAdvanceHandlers,
		advance.NewSocialAccountAdvanceHandlers,
		advance.NewAuctionAdvanceHandlers,
		// Inspect handlers
		inspect.NewOrderInspectHandlers,
		inspect.NewUserInspectHandlers,
		inspect.NewSocialAccountInspectHandlers,
		inspect.NewAuctionInspectHandlers,
		wire.Struct(new(Handlers), "*"),
	)
	return &Handlers{}, nil
}

// Handlers contains all handler dependencies
type Handlers struct {
	// Advance handlers
	OrderAdvanceHandlers   *advance.OrderAdvanceHandlers
	UserAdvanceHandlers    *advance.UserAdvanceHandlers
	SocialAccountsHandlers *advance.SocialAccountAdvanceHandlers
	AuctionAdvanceHandlers *advance.AuctionAdvanceHandlers

	// Inspect handlers
	OrderInspectHandlers   *inspect.OrderInspectHandlers
	UserInspectHandlers    *inspect.UserInspectHandlers
	SocialAccountHandlers  *inspect.SocialAccountInspectHandlers
	AuctionInspectHandlers *inspect.AuctionInspectHandlers
}
