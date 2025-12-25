package campaign

import (
	"fmt"
	"sort"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CloseCampaignInputDTO struct {
	CreatorAddress types.Address `json:"creator_address" validate:"required"`
}

type CloseCampaignOutputDTO struct {
	Id                uint                `json:"id"`
	Title             string              `json:"title,omitempty"`
	Description       string              `json:"description,omitempty"`
	Promotion         string              `json:"promotion,omitempty"`
	Token             types.Address       `json:"token,omitempty"`
	Creator           *user.UserOutputDTO `json:"creator,omitempty"`
	CollateralAddress types.Address       `json:"collateral,omitempty"`
	CollateralAmount  *uint256.Int        `json:"collateral_amount,omitempty"`
	BadgeAddress      types.Address       `json:"badge_address,omitempty"`
	DebtIssued        *uint256.Int        `json:"debt_issued,omitempty"`
	MaxInterestRate   *uint256.Int        `json:"max_interest_rate,omitempty"`
	TotalObligation   *uint256.Int        `json:"total_obligation,omitempty"`
	TotalRaised       *uint256.Int        `json:"total_raised,omitempty"`
	State             string              `json:"state,omitempty"`
	Orders            []*entity.Order     `json:"orders,omitempty"`
	CreatedAt         int64               `json:"created_at,omitempty"`
	ClosesAt          int64               `json:"closes_at,omitempty"`
	MaturityAt        int64               `json:"maturity_at,omitempty"`
	UpdatedAt         int64               `json:"updated_at,omitempty"`
}

type CloseCampaignUseCase struct {
	UserRepository     repository.UserRepository
	OrderRepository    repository.OrderRepository
	campaignRepository repository.CampaignRepository
}

func NewCloseCampaignUseCase(userRepo repository.UserRepository, campaignRepo repository.CampaignRepository, orderRepo repository.OrderRepository) *CloseCampaignUseCase {
	return &CloseCampaignUseCase{
		UserRepository:     userRepo,
		campaignRepository: campaignRepo,
		OrderRepository:    orderRepo,
	}
}

func (u *CloseCampaignUseCase) Execute(input *CloseCampaignInputDTO, metadata rollmelette.Metadata) (*CloseCampaignOutputDTO, error) {
	// -------------------------------------------------------------------------
	// 1. Find ongoing campaign for the creator
	// -------------------------------------------------------------------------
	campaigns, err := u.campaignRepository.FindCampaignsByCreatorAddress(input.CreatorAddress)
	if err != nil {
		return nil, err
	}
	var ongoingCampaign *entity.Campaign
	for _, campaign := range campaigns {
		if campaign.State == entity.CampaignStateOngoing {
			ongoingCampaign = campaign
			break
		}
	}
	if ongoingCampaign == nil {
		return nil, fmt.Errorf("no ongoing campaign found, cannot close it")
	}

	// -------------------------------------------------------------------------
	// 2. Validate campaign expiration
	// -------------------------------------------------------------------------
	if metadata.BlockTimestamp < ongoingCampaign.ClosesAt {
		return nil, fmt.Errorf("campaign not expired yet, cannot close it")
	}

	// -------------------------------------------------------------------------
	// 3. Fetch and sort campaign orders
	// -------------------------------------------------------------------------
	orders, err := u.OrderRepository.FindOrdersByCampaignId(ongoingCampaign.Id)
	if err != nil {
		return nil, err
	}
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].InterestRate.Cmp(orders[j].InterestRate) == 0 {
			return orders[i].Amount.Cmp(orders[j].Amount) > 0
		}
		return orders[i].InterestRate.Cmp(orders[j].InterestRate) < 0
	})

	// -------------------------------------------------------------------------
	// 4. Select winning orders and calculate obligations
	// -------------------------------------------------------------------------
	debtRemaining := new(uint256.Int).Set(ongoingCampaign.DebtIssued)
	totalCollected := uint256.NewInt(0)
	totalObligation := uint256.NewInt(0)

	for _, order := range orders {
		if debtRemaining.IsZero() {
			// Reject surplus orders
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			if _, err := u.OrderRepository.UpdateOrder(order); err != nil {
				return nil, err
			}
			continue
		}

		// Accept full or partial order
		acceptAmount := order.Amount
		if debtRemaining.Lt(order.Amount) {
			acceptAmount = new(uint256.Int).Set(debtRemaining)
		}
		interest := new(uint256.Int).Mul(acceptAmount, order.InterestRate)
		interest.Div(interest, uint256.NewInt(100))

		orderObligation := new(uint256.Int).Add(acceptAmount, interest)
		totalCollected.Add(totalCollected, acceptAmount)
		totalObligation.Add(totalObligation, orderObligation)

		if debtRemaining.Cmp(order.Amount) >= 0 {
			order.State = entity.OrderStateAccepted
			debtRemaining.Sub(debtRemaining, order.Amount)
		} else {
			order.State = entity.OrderStatePartiallyAccepted
			// Create rejected order for the surplus
			rejectedAmount := new(uint256.Int).Sub(order.Amount, acceptAmount)
			_, err := u.OrderRepository.CreateOrder(&entity.Order{
				CampaignId:   order.CampaignId,
				Investor:     order.Investor,
				Amount:       rejectedAmount,
				InterestRate: order.InterestRate,
				State:        entity.OrderStateRejected,
				CreatedAt:    order.CreatedAt,
				UpdatedAt:    metadata.BlockTimestamp,
			})
			if err != nil {
				return nil, err
			}
			debtRemaining.Clear()
		}
		order.Amount = acceptAmount
		order.UpdatedAt = metadata.BlockTimestamp
		if _, err := u.OrderRepository.UpdateOrder(order); err != nil {
			return nil, err
		}
	}

	// -------------------------------------------------------------------------
	// 5. Check if minimum funding (2/3) was reached
	// -------------------------------------------------------------------------
	twoThirds := new(uint256.Int).Mul(ongoingCampaign.DebtIssued, uint256.NewInt(2))
	twoThirds.Div(twoThirds, uint256.NewInt(3))
	if totalCollected.Lt(twoThirds) {
		// Cancel campaign and reject all orders
		for _, order := range orders {
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			if _, err := u.OrderRepository.UpdateOrder(order); err != nil {
				return nil, err
			}
		}
		ongoingCampaign.State = entity.CampaignStateCanceled
		ongoingCampaign.UpdatedAt = metadata.BlockTimestamp
		if _, err := u.campaignRepository.UpdateCampaign(ongoingCampaign); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("campaign canceled due to insufficient funds collected, expected at least 2/3 of the debt issued: %s, got: %s", twoThirds.String(), totalCollected.String())
	}

	// -------------------------------------------------------------------------
	// 6. Close campaign and return result
	// -------------------------------------------------------------------------
	ongoingCampaign.State = entity.CampaignStateClosed
	ongoingCampaign.TotalObligation = totalObligation
	ongoingCampaign.TotalRaised = totalCollected
	ongoingCampaign.UpdatedAt = metadata.BlockTimestamp
	res, err := u.campaignRepository.UpdateCampaign(ongoingCampaign)
	if err != nil {
		return nil, err
	}

	creator, err := u.UserRepository.FindUserByAddress(res.Creator)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}

	return &CloseCampaignOutputDTO{
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
		Orders:            res.Orders,
		State:             string(res.State),
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		CreatedAt:         res.CreatedAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
