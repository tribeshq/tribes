package repository

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	. "github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
)

type IssuanceRepository interface {
	CreateIssuance(issuance *entity.Issuance) (*entity.Issuance, error)
	FindIssuancesByCreatorAddress(creator Address) ([]*entity.Issuance, error)
	FindOngoingIssuanceByCreatorAddress(creator Address) (*entity.Issuance, error)
	FindIssuancesByInvestorAddress(investor Address) ([]*entity.Issuance, error)
	FindIssuanceById(id uint) (*entity.Issuance, error)
	FindAllIssuances() ([]*entity.Issuance, error)
	UpdateIssuance(Issuance *entity.Issuance) (*entity.Issuance, error)
}

type OrderRepository interface {
	CreateOrder(order *entity.Order) (*entity.Order, error)
	FindOrderById(id uint) (*entity.Order, error)
	FindOrdersByIssuanceId(id uint) ([]*entity.Order, error)
	FindOrdersByState(issuanceId uint, state string) ([]*entity.Order, error)
	FindOrdersByInvestorAddress(investor Address) ([]*entity.Order, error)
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
	FindUserByAddress(address Address) (*entity.User, error)
	FindAllUsers() ([]*entity.User, error)
	DeleteUser(address Address) error
}

type Repository interface {
	IssuanceRepository
	OrderRepository
	SocialAccountRepository
	UserRepository
	Close() error
}
