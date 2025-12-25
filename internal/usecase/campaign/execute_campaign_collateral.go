package campaign

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type ExecuteCampaignCollateralInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type ExecuteCampaignCollateralOutputDTO struct {
	Id                uint                `json:"id"`
	Title             string              `json:"title,omitempty"`
	Description       string              `json:"description,omitempty"`
	Promotion         string              `json:"promotion,omitempty"`
	Token             types.Address       `json:"token"`
	Creator           *user.UserOutputDTO `json:"creator"`
	CollateralAddress types.Address       `json:"collateral"`
	CollateralAmount  *uint256.Int        `json:"collateral_amount"`
	BadgeAddress      types.Address       `json:"badge_address"`
	DebtIssued        *uint256.Int        `json:"debt_issued"`
	MaxInterestRate   *uint256.Int        `json:"max_interest_rate"`
	TotalObligation   *uint256.Int        `json:"total_obligation"`
	TotalRaised       *uint256.Int        `json:"total_raised"`
	State             string              `json:"state"`
	Orders            []*entity.Order     `json:"orders"`
	CreatedAt         int64               `json:"created_at"`
	ClosesAt          int64               `json:"closes_at"`
	MaturityAt        int64               `json:"maturity_at"`
	UpdatedAt         int64               `json:"updated_at"`
}

type ExecuteCampaignCollateralUseCase struct {
	UserRepository     repository.UserRepository
	campaignRepository repository.CampaignRepository
	OrderRepository    repository.OrderRepository
}

func NewExecuteCampaignCollateralUseCase(userRepo repository.UserRepository, campaignRepo repository.CampaignRepository, orderRepo repository.OrderRepository) *ExecuteCampaignCollateralUseCase {
	return &ExecuteCampaignCollateralUseCase{
		UserRepository:     userRepo,
		campaignRepository: campaignRepo,
		OrderRepository:    orderRepo,
	}
}

func (uc *ExecuteCampaignCollateralUseCase) Execute(input *ExecuteCampaignCollateralInputDTO, metadata rollmelette.Metadata) (*ExecuteCampaignCollateralOutputDTO, error) {
	campaign, err := uc.campaignRepository.FindCampaignById(input.Id)
	if err != nil {
		return nil, err
	}

	if err := uc.Validate(campaign, metadata); err != nil {
		return nil, err
	}

	var ordersToUpdate []*entity.Order
	for _, order := range campaign.Orders {
		if order.State == entity.OrderStateAccepted || order.State == entity.OrderStatePartiallyAccepted {
			order.State = entity.OrderStateSettledByCollateral
			order.UpdatedAt = metadata.BlockTimestamp
			ordersToUpdate = append(ordersToUpdate, order)
		}
	}
	for _, order := range ordersToUpdate {
		if _, err := uc.OrderRepository.UpdateOrder(order); err != nil {
			return nil, fmt.Errorf("error updating order: %w", err)
		}
	}

	campaign.State = entity.CampaignStateCollateralExecuted
	campaign.UpdatedAt = metadata.BlockTimestamp

	res, err := uc.campaignRepository.UpdateCampaign(campaign)
	if err != nil {
		return nil, err
	}

	creator, err := uc.UserRepository.FindUserByAddress(res.Creator)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}

	return &ExecuteCampaignCollateralOutputDTO{
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
		Orders:            res.Orders,
		CreatedAt:         res.CreatedAt,
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}

func (uc *ExecuteCampaignCollateralUseCase) Validate(campaign *entity.Campaign, metadata rollmelette.Metadata) error {
	if metadata.BlockTimestamp < campaign.MaturityAt {
		return fmt.Errorf("the maturity date of the campaign campaign has not passed")
	}
	if campaign.State != entity.CampaignStateClosed {
		return fmt.Errorf("campaign campaign not closed")
	}
	return nil
}
