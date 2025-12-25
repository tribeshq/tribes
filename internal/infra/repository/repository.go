package repository

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	types "github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
)

type CampaignRepository interface {
	CreateCampaign(campaign *entity.Campaign) (*entity.Campaign, error)
	FindCampaignsByCreatorAddress(creator types.Address) ([]*entity.Campaign, error)
	FindCampaignsByInvestorAddress(investor types.Address) ([]*entity.Campaign, error)
	FindCampaignById(id uint) (*entity.Campaign, error)
	FindAllCampaigns() ([]*entity.Campaign, error)
	UpdateCampaign(Campaign *entity.Campaign) (*entity.Campaign, error)
}

type OrderRepository interface {
	CreateOrder(order *entity.Order) (*entity.Order, error)
	FindOrderById(id uint) (*entity.Order, error)
	FindOrdersByCampaignId(id uint) ([]*entity.Order, error)
	FindOrdersByState(campaignId uint, state string) ([]*entity.Order, error)
	FindOrdersByInvestorAddress(investor types.Address) ([]*entity.Order, error)
	FindAllOrders() ([]*entity.Order, error)
	UpdateOrder(order *entity.Order) (*entity.Order, error)
	DeleteOrder(id uint) error
}

type SocialAccountRepository interface {
	CreateSocialAccount(socialAccount *entity.SocialAccount) (*entity.SocialAccount, error)
	FindSocialAccountById(id uint) (*entity.SocialAccount, error)
	FindSocialAccountsByUserId(userID uint) ([]*entity.SocialAccount, error)
	DeleteSocialAccount(id uint) error
}

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	FindUsersByRole(role string) ([]*entity.User, error)
	FindUserByAddress(address types.Address) (*entity.User, error)
	FindAllUsers() ([]*entity.User, error)
	DeleteUser(address types.Address) error
}

type Repository interface {
	CampaignRepository
	OrderRepository
	SocialAccountRepository
	UserRepository
	Close() error
}
