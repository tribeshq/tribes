package campaign

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	user "github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
)

type FindCampaignByIdInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type FindCampaignByIdUseCase struct {
	UserRepository     repository.UserRepository
	campaignRepository repository.CampaignRepository
}

func NewFindCampaignByIdUseCase(userRepo repository.UserRepository, campaignRepo repository.CampaignRepository) *FindCampaignByIdUseCase {
	return &FindCampaignByIdUseCase{
		UserRepository:     userRepo,
		campaignRepository: campaignRepo,
	}
}

func (f *FindCampaignByIdUseCase) Execute(input *FindCampaignByIdInputDTO) (*CampaignOutputDTO, error) {
	res, err := f.campaignRepository.FindCampaignById(input.Id)
	if err != nil {
		return nil, err
	}
	orders := make([]*entity.Order, len(res.Orders))
	for i, order := range res.Orders {
		orders[i] = &entity.Order{
			Id:         order.Id,
			CampaignId: order.CampaignId,

			Investor:     order.Investor,
			Amount:       order.Amount,
			InterestRate: order.InterestRate,
			State:        order.State,
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
		}
	}
	creator, err := f.UserRepository.FindUserByAddress(res.Creator)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}
	return &CampaignOutputDTO{
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
