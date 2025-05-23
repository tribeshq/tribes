package contract_usecase

import (
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type FindContractOutputDTO struct {
	Id        uint    `json:"id"`
	Symbol    string  `json:"symbol"`
	Address   Address `json:"address"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}
