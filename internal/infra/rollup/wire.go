//go:build wireinject
// +build wireinject

package rollup

import (
	"github.com/google/wire"
	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/infra/rollup/handler/advance"
	"github.com/tribeshq/tribes/internal/infra/rollup/handler/inspect"
)

func NewHandlers(repo repository.Repository, cfg *configs.RollupConfig) (*Handlers, error) {
	wire.Build(
		// Bind repository interfaces
		wire.Bind(new(repository.UserRepository), new(repository.Repository)),
		wire.Bind(new(repository.OrderRepository), new(repository.Repository)),
		wire.Bind(new(repository.CampaignRepository), new(repository.Repository)),
		wire.Bind(new(repository.SocialAccountRepository), new(repository.Repository)),
		// Advance handlers
		advance.NewOrderAdvanceHandlers,
		advance.NewUserAdvanceHandlers,
		advance.NewSocialAccountAdvanceHandlers,
		advance.NewCampaignAdvanceHandlers,
		// Inspect handlers
		inspect.NewOrderInspectHandlers,
		inspect.NewUserInspectHandlers,
		inspect.NewSocialAccountInspectHandlers,
		inspect.NewCampaignInspectHandlers,
		wire.Struct(new(Handlers), "*"),
	)
	return &Handlers{}, nil
}

// Handlers contains all handler dependencies
type Handlers struct {
	// Advance handlers
	OrderAdvanceHandlers    *advance.OrderAdvanceHandlers
	UserAdvanceHandlers     *advance.UserAdvanceHandlers
	SocialAccountsHandlers  *advance.SocialAccountAdvanceHandlers
	CampaignAdvanceHandlers *advance.CampaignAdvanceHandlers

	// Inspect handlers
	OrderInspectHandlers    *inspect.OrderInspectHandlers
	UserInspectHandlers     *inspect.UserInspectHandlers
	SocialAccountHandlers   *inspect.SocialAccountInspectHandlers
	CampaignInspectHandlers *inspect.CampaignInspectHandlers
}
