package issuance

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/order"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
)

type FindIssuanceByIdInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type FindIssuanceByIdUseCase struct {
	UserRepository     repository.UserRepository
	IssuanceRepository repository.IssuanceRepository
}

func NewFindIssuanceByIdUseCase(userRepo repository.UserRepository, issuanceRepo repository.IssuanceRepository) *FindIssuanceByIdUseCase {
	return &FindIssuanceByIdUseCase{
		UserRepository:     userRepo,
		IssuanceRepository: issuanceRepo,
	}
}

func (f *FindIssuanceByIdUseCase) Execute(input *FindIssuanceByIdInputDTO) (*IssuanceOutputDTO, error) {
	res, err := f.IssuanceRepository.FindIssuanceById(input.Id)
	if err != nil {
		return nil, err
	}
	orders := make([]*order.OrderOutputDTO, len(res.Orders))
	for i, o := range res.Orders {
		investor, err := f.UserRepository.FindUserByAddress(o.InvestorAddress)
		if err != nil {
			return nil, fmt.Errorf("error finding investor: %w", err)
		}
		orders[i] = &order.OrderOutputDTO{
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
	creator, err := f.UserRepository.FindUserByAddress(res.CreatorAddress)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}
	return &IssuanceOutputDTO{
		Id:          res.Id,
		Title:       res.Title,
		Description: res.Description,
		Promotion:   res.Promotion,
		Token:       res.Token,
		Creator: &user.UserOutputDTO{
			Id:             creator.Id,
			Role:           string(creator.Role),
			Address:        creator.Address,
			SocialAccounts: creator.SocialAccounts,
			CreatedAt:      creator.CreatedAt,
			UpdatedAt:      creator.UpdatedAt,
		},
		CollateralAddress: res.CollateralAddress,
		CollateralAmount:  res.CollateralAmount,
		BadgeAddress:      res.BadgeAddress,
		DebtIssued:        res.DebtIssued,
		MaxInterestRate:   res.MaxInterestRate,
		TotalObligation:   res.TotalObligation,
		TotalRaised:       res.TotalRaised,
		State:             string(res.State),
		Orders:            orders,
		CreatedAt:         res.CreatedAt,
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
