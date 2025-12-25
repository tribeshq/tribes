package advance

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/configs"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/campaign"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CampaignAdvanceHandlers struct {
	cfg                *configs.RollupConfig
	OrderRepository    repository.OrderRepository
	UserRepository     repository.UserRepository
	campaignRepository repository.CampaignRepository
}

func NewCampaignAdvanceHandlers(
	cfg *configs.RollupConfig,
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	campaignRepo repository.CampaignRepository,
) *CampaignAdvanceHandlers {
	return &CampaignAdvanceHandlers{
		cfg:                cfg,
		OrderRepository:    orderRepo,
		UserRepository:     userRepo,
		campaignRepository: campaignRepo,
	}
}

func (h *CampaignAdvanceHandlers) CreateCampaign(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input campaign.CreateCampaignInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	createCampaign := campaign.NewCreateCampaignUseCase(
		h.cfg,
		h.campaignRepository,
		h.UserRepository,
	)

	res, err := createCampaign.Execute(&input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}

	abiJson := `[{
		"type": "function",
		"name": "newBadge",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	deployBadgePayload, err := abiInterface.Pack(
		"newBadge",
		env.AppAddress(),
		common.HexToHash(strconv.Itoa(metadata.Index)),
	)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %w", err)
	}
	env.Voucher(
		common.Address(h.cfg.BadgeFactoryAddress),
		big.NewInt(0),
		deployBadgePayload,
	)

	erc20Deposit := deposit.(*rollmelette.ERC20Deposit)
	if err := env.ERC20Transfer(
		erc20Deposit.Token,
		erc20Deposit.Sender,
		env.AppAddress(),
		erc20Deposit.Value,
	); err != nil {
		return fmt.Errorf("failed to transfer ERC20: %w", err)
	}

	campaign, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("campaign created - "), campaign...))
	return nil
}

func (h *CampaignAdvanceHandlers) CloseCampaign(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input campaign.CloseCampaignInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	closeCampaign := campaign.NewCloseCampaignUseCase(h.UserRepository, h.campaignRepository, h.OrderRepository)
	res, err := closeCampaign.Execute(&input, metadata)
	if err != nil && res == nil {
		return fmt.Errorf("failed to close campaign: %w", err)
	}

	token := common.Address(res.Token)

	// Process orders
	for _, order := range res.Orders {
		if order.State == entity.OrderStateRejected {
			if err = env.ERC20Transfer(
				token,
				env.AppAddress(),
				common.Address(order.Investor),
				order.Amount.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer rejected order: %w", err)
			}
		}
	}

	abiJSON := `[{
		"type":"function",
		"name":"safeMint",
		"inputs":[
			{"type":"address"},
			{"type":"address"},
			{"type":"uint256"},
			{"type":"uint256"},
			{"type":"bytes"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	for _, order := range res.Orders {
		if order.State != entity.OrderStateRejected {
			safeMintPayload, err := abiInterface.Pack(
				"safeMint",
				common.Address(res.BadgeAddress),
				common.Address(order.Investor),
				big.NewInt(1),
				big.NewInt(1),
				[]byte{},
			)
			if err != nil {
				return fmt.Errorf("failed to pack ABI: %w", err)
			}
			env.DelegateCallVoucher(common.Address(h.cfg.SafeErc1155MintAddress), safeMintPayload)
		}
	}

	if err := env.ERC20Transfer(token, env.AppAddress(), common.Address(res.Creator.Address), res.TotalRaised.ToBig()); err != nil {
		return fmt.Errorf("failed to transfer total raised: %w", err)
	}

	campaign, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte(fmt.Sprintf("campaign %v - ", res.State)), campaign...))
	return nil
}

func (h *CampaignAdvanceHandlers) SettleCampaign(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input campaign.SettleCampaignInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	settleCampaign := campaign.NewSettleCampaignUseCase(
		h.UserRepository,
		h.campaignRepository,
		h.OrderRepository,
	)

	res, err := settleCampaign.Execute(&input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to settle campaign: %w", err)
	}

	contractAddr := common.Address(res.Token)
	creatorAddr := common.Address(res.Creator.Address)

	// Process settled orders
	for _, order := range res.Orders {
		if order.State == entity.OrderStateSettled {
			// Calculate interest for this order
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))

			// Calculate total payment
			totalPayment := new(uint256.Int).Add(order.Amount, interest)

			if err := env.ERC20Transfer(
				contractAddr,
				creatorAddr,
				common.Address(order.Investor),
				totalPayment.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer settled order: %w", err)
			}
		}
	}

	campaign, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("campaign settled - "), campaign...))
	return nil
}

func (h *CampaignAdvanceHandlers) ExecuteCampaignCollateral(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input campaign.ExecuteCampaignCollateralInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	executeCampaignCollateral := campaign.NewExecuteCampaignCollateralUseCase(h.UserRepository, h.campaignRepository, h.OrderRepository)
	res, err := executeCampaignCollateral.Execute(&input, metadata)
	if err != nil {
		return fmt.Errorf("failed to execute campaign collateral: %w", err)
	}

	totalFinalValue := uint256.NewInt(0)
	orderFinalValues := make(map[uint]*uint256.Int)
	for _, order := range res.Orders {
		if order.State == entity.OrderStateSettledByCollateral {
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))
			finalValue := new(uint256.Int).Add(order.Amount, interest)
			orderFinalValues[order.Id] = finalValue
			totalFinalValue.Add(totalFinalValue, finalValue)
		}
	}

	for _, order := range res.Orders {
		if order.State == entity.OrderStateSettledByCollateral {
			finalValue := orderFinalValues[order.Id]
			orderShare := new(uint256.Int).Mul(finalValue, res.CollateralAmount)
			orderShare.Div(orderShare, totalFinalValue)

			if err = env.ERC20Transfer(
				common.Address(res.CollateralAddress),
				env.AppAddress(),
				common.Address(order.Investor),
				orderShare.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer collateral to investor: %w", err)
			}
		}
	}

	campaign, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("campaign collateral executed - "), campaign...))
	return nil
}
