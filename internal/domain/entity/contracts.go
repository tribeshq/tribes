package entity

import (
	"errors"
	"fmt"

	. "github.com/tribeshq/tribes/pkg/custom_type"
)

var (
	ErrInvalidContract  = errors.New("invalid contract")
	ErrContractNotFound = errors.New("contract not found")
)

type Contract struct {
	Id        uint    `json:"id" gorm:"primaryKey"`
	Symbol    string  `json:"symbol,omitempty" gorm:"uniqueIndex;not null"`
	Address   Address `json:"address,omitempty" gorm:"custom_type:text;not null"`
	CreatedAt int64   `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt int64   `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewContract(symbol string, address Address, createdAt int64) (*Contract, error) {
	contract := &Contract{
		Symbol:    symbol,
		Address:   address,
		CreatedAt: createdAt,
	}
	if err := contract.validate(); err != nil {
		return nil, err
	}
	return contract, nil
}

func (c *Contract) validate() error {
	if c.Symbol == "" {
		return fmt.Errorf("%w: symbol cannot be empty", ErrInvalidContract)
	}
	if c.Address == (Address{}) {
		return fmt.Errorf("%w: address cannot be empty", ErrInvalidContract)
	}
	return nil
}
