package advance

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/campaign"
	"github.com/tribeshq/tribes/pkg/deploy"
)

type CampaignAdvanceHandlers struct {
	OrderRepository    repository.OrderRepository
	UserRepository     repository.UserRepository
	CampaignRepository repository.CampaignRepository
}

func NewCampaignAdvanceHandlers(
	orderRepository repository.OrderRepository,
	userRepository repository.UserRepository,
	campaignRepository repository.CampaignRepository,
) *CampaignAdvanceHandlers {
	return &CampaignAdvanceHandlers{
		OrderRepository:    orderRepository,
		UserRepository:     userRepository,
		CampaignRepository: campaignRepository,
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

	ctx := context.Background()
	createCampaign := campaign.NewCreateCampaignUseCase(
		h.CampaignRepository,
		h.UserRepository,
	)

	res, err := createCampaign.Execute(ctx, &input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}

	bytecode, err := deploy.GetBytecodeFromJSON("../../skel/Badge.json", "bytecode")
	if err != nil {
		return fmt.Errorf("failed to get bytecode: %w", err)
	}

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(env.AppAddress())
	if err != nil {
		return fmt.Errorf("failed to encode constructor args: %w", err)
	}
	initCode := append(bytecode, constructorArgs...)

	abiJson := `[{
		"type": "function",
		"name": "deploy2",
		"inputs": [
			{"type": "bytes"},
			{"type": "bytes32"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	deploy2Payload, err := abiInterface.Pack(
		"deploy2",
		initCode,
		common.HexToHash(strconv.Itoa(int(metadata.BlockTimestamp))),
	)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %w", err)
	}
	env.Voucher(common.Address(res.DeployerAddress), big.NewInt(0), deploy2Payload)

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

	ctx := context.Background()
	closeCampaign := campaign.NewCloseCampaignUseCase(h.UserRepository, h.CampaignRepository, h.OrderRepository)
	res, err := closeCampaign.Execute(ctx, &input, metadata)
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
		"name":"mint",
		"inputs":[
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
			mintPayload, err := abiInterface.Pack(
				"mint",
				common.Address(order.Investor),
				big.NewInt(1),
				big.NewInt(1),
				[]byte{},
			)
			if err != nil {
				return fmt.Errorf("failed to pack ABI: %w", err)
			}
			env.Voucher(common.Address(res.BadgeAddress), big.NewInt(0), mintPayload)
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

	ctx := context.Background()
	settleCampaign := campaign.NewSettleCampaignUseCase(
		h.UserRepository,
		h.CampaignRepository,
		h.OrderRepository,
	)

	res, err := settleCampaign.Execute(ctx, &input, deposit, metadata)
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

	ctx := context.Background()
	executeCampaignCollateral := campaign.NewExecuteCampaignCollateralUseCase(h.UserRepository, h.CampaignRepository, h.OrderRepository)
	res, err := executeCampaignCollateral.Execute(ctx, &input, metadata)
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
