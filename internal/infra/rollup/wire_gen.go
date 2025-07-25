// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package rollup

import (
	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/infra/rollup/handler/advance"
	"github.com/tribeshq/tribes/internal/infra/rollup/handler/inspect"
)

// Injectors from wire.go:

func NewHandlers(repo repository.Repository, cfg *configs.RollupConfig) (*Handlers, error) {
	orderAdvanceHandlers := advance.NewOrderAdvanceHandlers(repo, repo, repo)
	userAdvanceHandlers := advance.NewUserAdvanceHandlers(cfg, repo)
	socialAccountAdvanceHandlers := advance.NewSocialAccountAdvanceHandlers(repo, repo)
	campaignAdvanceHandlers := advance.NewCampaignAdvanceHandlers(cfg, repo, repo, repo)
	orderInspectHandlers := inspect.NewOrderInspectHandlers(repo, repo)
	userInspectHandlers := inspect.NewUserInspectHandlers(repo)
	socialAccountInspectHandlers := inspect.NewSocialAccountInspectHandlers(repo)
	campaignInspectHandlers := inspect.NewCampaignInspectHandlers(repo, repo)
	handlers := &Handlers{
		OrderAdvanceHandlers:    orderAdvanceHandlers,
		UserAdvanceHandlers:     userAdvanceHandlers,
		SocialAccountsHandlers:  socialAccountAdvanceHandlers,
		CampaignAdvanceHandlers: campaignAdvanceHandlers,
		OrderInspectHandlers:    orderInspectHandlers,
		UserInspectHandlers:     userInspectHandlers,
		SocialAccountHandlers:   socialAccountInspectHandlers,
		CampaignInspectHandlers: campaignInspectHandlers,
	}
	return handlers, nil
}

// wire.go:

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
