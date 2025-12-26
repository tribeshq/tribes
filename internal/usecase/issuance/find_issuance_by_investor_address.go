package issuance

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/order"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
)

type FindIssuancesByInvestorAddressInputDTO struct {
	InvestorAddress types.Address `json:"investor_address" validate:"required"`
}

type FindIssuancesByInvestorAddressOutputDTO []*IssuanceOutputDTO

type FindIssuancesByInvestorAddressUseCase struct {
	UserRepository     repository.UserRepository
	IssuanceRepository repository.IssuanceRepository
}

func NewFindIssuancesByInvestorAddressUseCase(userRepo repository.UserRepository, issuanceRepo repository.IssuanceRepository) *FindIssuancesByInvestorAddressUseCase {
	return &FindIssuancesByInvestorAddressUseCase{
		UserRepository:     userRepo,
		IssuanceRepository: issuanceRepo,
	}
}

func (f *FindIssuancesByInvestorAddressUseCase) Execute(input *FindIssuancesByInvestorAddressInputDTO) (*FindIssuancesByInvestorAddressOutputDTO, error) {
	res, err := f.IssuanceRepository.FindIssuancesByInvestorAddress(input.InvestorAddress)
	if err != nil {
		return nil, err
	}
	output := make(FindIssuancesByInvestorAddressOutputDTO, len(res))
	for i, issuance := range res {
		orders := make([]*order.OrderOutputDTO, len(issuance.Orders))
		for j, o := range issuance.Orders {
			investor, err := f.UserRepository.FindUserByAddress(o.InvestorAddress)
			if err != nil {
				return nil, fmt.Errorf("error finding investor: %w", err)
			}
			orders[j] = &order.OrderOutputDTO{
				Id:         o.Id,
				IssuanceId: o.IssuanceId,
				Investor: &user.UserOutputDTO{
					Id:             investor.Id,
					Role:           string(investor.Role),
					Address:        investor.Address,
					SocialAccounts: investor.SocialAccounts,
					CreatedAt:      investor.CreatedAt,
					UpdatedAt:      investor.UpdatedAt,
				},
				Amount:       o.Amount,
				InterestRate: o.InterestRate,
				State:        string(o.State),
				CreatedAt:    o.CreatedAt,
				UpdatedAt:    o.UpdatedAt,
			}
		}
		creator, err := f.UserRepository.FindUserByAddress(issuance.CreatorAddress)
		if err != nil {
			return nil, fmt.Errorf("error finding creator: %w", err)
		}
		output[i] = &IssuanceOutputDTO{
			Id:          issuance.Id,
			Title:       issuance.Title,
			Description: issuance.Description,
			Promotion:   issuance.Promotion,
			Token:       issuance.Token,
			Creator: &user.UserOutputDTO{
				Id:             creator.Id,
				Role:           string(creator.Role),
				Address:        creator.Address,
				SocialAccounts: creator.SocialAccounts,
				CreatedAt:      creator.CreatedAt,
				UpdatedAt:      creator.UpdatedAt,
			},
			CollateralAddress: issuance.CollateralAddress,
			CollateralAmount:  issuance.CollateralAmount,
			BadgeAddress:      issuance.BadgeAddress,
			DebtIssued:        issuance.DebtIssued,
			MaxInterestRate:   issuance.MaxInterestRate,
			TotalObligation:   issuance.TotalObligation,
			TotalRaised:       issuance.TotalRaised,
			State:             string(issuance.State),
			Orders:            orders,
			CreatedAt:         issuance.CreatedAt,
			ClosesAt:          issuance.ClosesAt,
			MaturityAt:        issuance.MaturityAt,
			UpdatedAt:         issuance.UpdatedAt,
		}
	}
	return &output, nil
}
