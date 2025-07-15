package repository

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CampaignRepository interface {
	CreateCampaign(campaign *entity.Campaign) (*entity.Campaign, error)
	FindCampaignsByCreatorAddress(creator custom_type.Address) ([]*entity.Campaign, error)
	FindCampaignsByInvestorAddress(investor custom_type.Address) ([]*entity.Campaign, error)
	FindCampaignById(id uint) (*entity.Campaign, error)
	FindAllCampaigns() ([]*entity.Campaign, error)
	UpdateCampaign(Campaign *entity.Campaign) (*entity.Campaign, error)
}

type OrderRepository interface {
	CreateOrder(order *entity.Order) (*entity.Order, error)
	FindOrderById(id uint) (*entity.Order, error)
	FindOrdersByCampaignId(id uint) ([]*entity.Order, error)
	FindOrdersByState(CampaignId uint, state string) ([]*entity.Order, error)
	FindOrdersByInvestorAddress(investor custom_type.Address) ([]*entity.Order, error)
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
	FindUserByAddress(address custom_type.Address) (*entity.User, error)
	FindAllUsers() ([]*entity.User, error)
	DeleteUser(address custom_type.Address) error
}

type Repository interface {
	CampaignRepository
	OrderRepository
	SocialAccountRepository
	UserRepository
	Close() error
}
