package entity

import (
	"errors"
	"fmt"

	"github.com/holiman/uint256"
	. "github.com/tribeshq/tribes/pkg/custom_type"
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

type User struct {
	Id                uint             `json:"id" gorm:"primaryKey"`
	Role              UserRole         `json:"role,omitempty" gorm:"not null"`
	Address           Address          `json:"address,omitempty" gorm:"custom_type:text;uniqueIndex;not null"`
	InvestmentLimit   *uint256.Int     `json:"investment_limit,omitempty" gorm:"custom_type:text"`
	DebtIssuanceLimit *uint256.Int     `json:"debt_issuance_limit,omitempty" gorm:"custom_type:text"`
	CreatedAt         int64            `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt         int64            `json:"updated_at,omitempty" gorm:"default:0"`
	SocialAccounts    []*SocialAccount `json:"social_accounts,omitempty" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func NewUser(role string, investmentLimit *uint256.Int, debtIssuanceLimit *uint256.Int, address Address, created_at int64) (*User, error) {
	user := &User{
		Role:              UserRole(role),
		InvestmentLimit:   investmentLimit,
		DebtIssuanceLimit: debtIssuanceLimit,
		Address:           address,
		CreatedAt:         created_at,
	}
	if err := user.validate(); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) validate() error {
	if u.Role == "" {
		return fmt.Errorf("%w: role cannot be empty", ErrInvalidUser)
	}
	if u.Address == (Address{}) {
		return fmt.Errorf("%w: address cannot be empty", ErrInvalidUser)
	}
	if u.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidCrowdfunding)
	}
	return nil
}
