package entity

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	ErrInvalidUser  = errors.New("invalid user")
	ErrUserNotFound = errors.New("user not found")
)

type UserRole string

const (
	UserRoleAdmin                UserRole = "admin"
	UserRoleCreator              UserRole = "creator"
	UserRoleNonQualifiedInvestor UserRole = "non_qualified_investor"
	UserRoleQualifiedInvestor    UserRole = "qualified_investor"
)

type UserRepository interface {
	CreateUser(ctx context.Context, User *User) (*User, error)
	FindUsersByRole(ctx context.Context, role string) ([]*User, error)
	FindUserByAddress(ctx context.Context, address common.Address) (*User, error)
	FindAllUsers(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, User *User) (*User, error)
	DeleteUser(ctx context.Context, address common.Address) error
}

type User struct {
	Id                uint           `json:"id" gorm:"primaryKey"`
	Role              UserRole       `json:"role,omitempty" gorm:"not null"`
	Address           common.Address `json:"address,omitempty" gorm:"type:text;uniqueIndex;not null"`
	InvestmentLimit   *uint256.Int   `json:"investment_limit,omitempty" gorm:"type:text"`
	DebtIssuanceLimit *uint256.Int   `json:"debt_issuance_limit,omitempty" gorm:"type:text"`
	CreatedAt         int64          `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt         int64          `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewUser(role string, address common.Address, created_at int64) (*User, error) {
	user := &User{
		Role:              UserRole(role),
		InvestmentLimit:   setInvestimentLimit(role),
		DebtIssuanceLimit: setDebtIssuanceLimit(role),
		Address:           address,
		CreatedAt:         created_at,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Validate() error {
	if u.Role == "" {
		return fmt.Errorf("%w: role cannot be empty", ErrInvalidUser)
	}
	if u.Address == (common.Address{}) {
		return fmt.Errorf("%w: address cannot be empty", ErrInvalidUser)
	}
	if u.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidCrowdfunding)
	}
	return nil
}

func setInvestimentLimit(role string) *uint256.Int {
	switch role {
	case string(UserRoleQualifiedInvestor):
		return new(uint256.Int).SetAllOne() //According to CVM Resolution 88
	case string(UserRoleNonQualifiedInvestor):
		return uint256.NewInt(20000) //According to CVM Resolution 88
	default:
		return uint256.NewInt(0)
	}
}

func setDebtIssuanceLimit(role string) *uint256.Int {
	switch role {
	case string(UserRoleCreator):
		return uint256.NewInt(15000000) //According to CVM Resolution 88
	default:
		return uint256.NewInt(0)
	}
}
