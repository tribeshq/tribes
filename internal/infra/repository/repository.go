package repository

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type ContractRepository interface {
	CreateContract(ctx context.Context, contract *entity.Contract) (*entity.Contract, error)
	FindAllContracts(ctx context.Context) ([]*entity.Contract, error)
	FindContractBySymbol(ctx context.Context, symbol string) (*entity.Contract, error)
	FindContractByAddress(ctx context.Context, address Address) (*entity.Contract, error)
	UpdateContract(ctx context.Context, contract *entity.Contract) (*entity.Contract, error)
	DeleteContract(ctx context.Context, symbol string) error
}

type CrowdfundingRepository interface {
	CreateCrowdfunding(ctx context.Context, crowdfunding *entity.Crowdfunding) (*entity.Crowdfunding, error)
	FindCrowdfundingsByCreator(ctx context.Context, creator Address) ([]*entity.Crowdfunding, error)
	FindCrowdfundingsByInvestor(ctx context.Context, investor Address) ([]*entity.Crowdfunding, error)
	FindCrowdfundingById(ctx context.Context, id uint) (*entity.Crowdfunding, error)
	FindAllCrowdfundings(ctx context.Context) ([]*entity.Crowdfunding, error)
	UpdateCrowdfunding(ctx context.Context, crowdfunding *entity.Crowdfunding) (*entity.Crowdfunding, error)
	DeleteCrowdfunding(ctx context.Context, id uint) error
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	FindOrderById(ctx context.Context, id uint) (*entity.Order, error)
	FindOrdersByCrowdfundingId(ctx context.Context, id uint) ([]*entity.Order, error)
	FindOrdersByState(ctx context.Context, crowdfundingId uint, state string) ([]*entity.Order, error)
	FindOrdersByInvestor(ctx context.Context, investor Address) ([]*entity.Order, error)
	FindAllOrders(ctx context.Context) ([]*entity.Order, error)
	UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	DeleteOrder(ctx context.Context, id uint) error
}

type SocialAccountRepository interface {
	CreateSocialAccount(ctx context.Context, socialAccount *entity.SocialAccount) (*entity.SocialAccount, error)
	FindSocialAccountById(ctx context.Context, id uint) (*entity.SocialAccount, error)
	FindSocialAccountsByUserId(ctx context.Context, userID uint) ([]*entity.SocialAccount, error)
	DeleteSocialAccount(ctx context.Context, id uint) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, User *entity.User) (*entity.User, error)
	FindUsersByRole(ctx context.Context, role string) ([]*entity.User, error)
	FindUserByAddress(ctx context.Context, address Address) (*entity.User, error)
	FindAllUsers(ctx context.Context) ([]*entity.User, error)
	UpdateUser(ctx context.Context, User *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, address Address) error
}

type Repository interface {
	ContractRepository
	CrowdfundingRepository
	OrderRepository
	SocialAccountRepository
	UserRepository
	Close() error
}
