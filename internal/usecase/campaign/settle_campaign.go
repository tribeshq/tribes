package campaign

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/order"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type SettleCampaignInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type SettleCampaignOutputDTO struct {
	Id                uint                    `json:"id"`
	Title             string                  `json:"title,omitempty"`
	Description       string                  `json:"description,omitempty"`
	Promotion         string                  `json:"promotion,omitempty"`
	Token             types.Address           `json:"token"`
	Creator           *user.UserOutputDTO     `json:"creator"`
	CollateralAddress types.Address           `json:"collateral"`
	CollateralAmount  *uint256.Int            `json:"collateral_amount"`
	BadgeAddress      types.Address           `json:"badge_address"`
	DebtIssued        *uint256.Int            `json:"debt_issued"`
	MaxInterestRate   *uint256.Int            `json:"max_interest_rate"`
	TotalObligation   *uint256.Int            `json:"total_obligation"`
	TotalRaised       *uint256.Int            `json:"total_raised"`
	State             string                  `json:"state"`
	Orders            []*order.OrderOutputDTO `json:"orders"`
	CreatedAt         int64                   `json:"created_at"`
	ClosesAt          int64                   `json:"closes_at"`
	MaturityAt        int64                   `json:"maturity_at"`
	UpdatedAt         int64                   `json:"updated_at"`
}

type SettleCampaignUseCase struct {
	UserRepository     repository.UserRepository
	CampaignRepository repository.CampaignRepository
	OrderRepository    repository.OrderRepository
}

func NewSettleCampaignUseCase(
	UserRepository repository.UserRepository,
	CampaignRepository repository.CampaignRepository,
	OrderRepository repository.OrderRepository,
) *SettleCampaignUseCase {
	return &SettleCampaignUseCase{
		UserRepository:     UserRepository,
		CampaignRepository: CampaignRepository,
		OrderRepository:    OrderRepository,
	}
}

func (uc *SettleCampaignUseCase) Execute(
	input *SettleCampaignInputDTO,
	deposit rollmelette.Deposit,
	metadata rollmelette.Metadata,
) (*SettleCampaignOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit types: %T", deposit)
	}

	campaign, err := uc.CampaignRepository.FindCampaignById(input.Id)
	if err != nil {
		return nil, fmt.Errorf("error finding campaign: %w", err)
	}

	if err := uc.Validate(campaign, erc20Deposit, metadata); err != nil {
		return nil, err
	}

	var ordersToUpdate []*entity.Order
	for _, order := range campaign.Orders {
		if order.State == entity.OrderStateAccepted || order.State == entity.OrderStatePartiallyAccepted {
			order.State = entity.OrderStateSettled
			order.UpdatedAt = metadata.BlockTimestamp
			ordersToUpdate = append(ordersToUpdate, order)
		}
	}
	for _, order := range ordersToUpdate {
		if _, err := uc.OrderRepository.UpdateOrder(order); err != nil {
			return nil, fmt.Errorf("error updating order: %w", err)
		}
	}

	campaign.State = entity.CampaignStateSettled
	campaign.UpdatedAt = metadata.BlockTimestamp
	res, err := uc.CampaignRepository.UpdateCampaign(campaign)
	if err != nil {
		return nil, fmt.Errorf("error updating campaign: %w", err)
	}

	creator, err := uc.UserRepository.FindUserByAddress(res.CreatorAddress)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}

	orderDTOs := make([]*order.OrderOutputDTO, len(res.Orders))
	for i, o := range res.Orders {
		investor, err := uc.UserRepository.FindUserByAddress(o.InvestorAddress)
		if err != nil {
			return nil, fmt.Errorf("error finding investor: %w", err)
		}
		orderDTOs[i] = &order.OrderOutputDTO{
			Id:         o.Id,
			CampaignId: o.CampaignId,
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

	return &SettleCampaignOutputDTO{
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
		Orders:            orderDTOs,
		CreatedAt:         res.CreatedAt,
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}

func (uc *SettleCampaignUseCase) Validate(
	Campaign *entity.Campaign,
	deposit *rollmelette.ERC20Deposit,
	metadata rollmelette.Metadata,
) error {
	if metadata.BlockTimestamp > Campaign.MaturityAt {
		return fmt.Errorf("the maturity date of the campaign campaign has passed")
	}

	if Campaign.State == entity.CampaignStateSettled {
		return fmt.Errorf("campaign campaign already settled")
	}

	if Campaign.State != entity.CampaignStateClosed {
		return fmt.Errorf("campaign campaign not closed")
	}

	if deposit.Value.Cmp(Campaign.TotalObligation.ToBig()) < 0 {
		return fmt.Errorf("deposit amount is lower than the total obligation")
	}

	if Campaign.CreatorAddress != types.Address(deposit.Sender) {
		return fmt.Errorf("only the campaign creator can settle the campaign")
	}
	return nil
}
