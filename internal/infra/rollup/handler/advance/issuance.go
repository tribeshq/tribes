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
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/issuance"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

// BasisPointsDivisor is the divisor used for basis point calculations.
// 10000 basis points = 100%, so 1 basis point = 0.01%.
// Example: 900 basis points = 9%, 950 basis points = 9.5%
var (
	BasisPointsDivisor = uint256.NewInt(10000)
)

type IssuanceAdvanceHandlers struct {
	Config             *configs.RollupConfig
	OrderRepository    repository.OrderRepository
	UserRepository     repository.UserRepository
	IssuanceRepository repository.IssuanceRepository
}

func NewIssuanceAdvanceHandlers(
	cfg *configs.RollupConfig,
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	issuanceRepo repository.IssuanceRepository,
) *IssuanceAdvanceHandlers {
	return &IssuanceAdvanceHandlers{
		Config:             cfg,
		OrderRepository:    orderRepo,
		UserRepository:     userRepo,
		IssuanceRepository: issuanceRepo,
	}
}

func (h *IssuanceAdvanceHandlers) CreateIssuance(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input issuance.CreateIssuanceInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	createIssuance := issuance.NewCreateIssuanceUseCase(
		h.Config.BadgeFactoryAddress,
		h.IssuanceRepository,
		h.UserRepository,
	)

	res, err := createIssuance.Execute(&input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to create issuance: %w", err)
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
		common.Address(h.Config.BadgeFactoryAddress),
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

	issuance, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("issuance created - "), issuance...))
	return nil
}

func (h *IssuanceAdvanceHandlers) CloseIssuance(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input issuance.CloseIssuanceInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	closeIssuance := issuance.NewCloseIssuanceUseCase(h.UserRepository, h.IssuanceRepository, h.OrderRepository)
	res, err := closeIssuance.Execute(&input, metadata)
	if err != nil && res == nil {
		return fmt.Errorf("failed to close issuance: %w", err)
	}

	token := common.Address(res.Token)

	// Process orders
	for _, order := range res.Orders {
		if order.State == string(entity.OrderStateRejected) {
			if err = env.ERC20Transfer(
				token,
				env.AppAddress(),
				common.Address(order.Investor.Address),
				order.Amount.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer rejected order: %w", err)
			}
		}
	}

	findAdminUseCase := user.NewFindUsersByRoleUseCase(h.UserRepository)
	admins, err := findAdminUseCase.Execute(&user.FindUserByRoleInputDTO{Role: "admin"})
	if err != nil {
		return fmt.Errorf("failed to find admin user: %w", err)
	}
	if len(admins) == 0 {
		return fmt.Errorf("no admin user found to receive fee")
	}
	adminUser := admins[0]

	// Calculate fee distribution: fee goes to admin, rest goes to creator
	// IssuanceFee is in basis points (e.g., 500 = 5%, 250 = 2.5%)
	feeBasisPoints := uint256.NewInt(h.Config.IssuanceFee)
	basisPointsDivisor := uint256.NewInt(10000)

	// Calculate admin fee: (totalRaised * feeBasisPoints) / 10000
	adminFee := new(uint256.Int).Mul(res.TotalRaised, feeBasisPoints)
	adminFee.Div(adminFee, basisPointsDivisor)

	// Calculate creator amount: totalRaised - adminFee
	creatorAmount := new(uint256.Int).Sub(res.TotalRaised, adminFee)

	// Transfer fee to admin
	if adminFee.Sign() > 0 {
		if err := env.ERC20Transfer(token, env.AppAddress(), common.Address(adminUser.Address), adminFee.ToBig()); err != nil {
			return fmt.Errorf("failed to transfer fee to admin: %w", err)
		}
	}

	// Transfer remaining amount to creator
	if err := env.ERC20Transfer(token, env.AppAddress(), common.Address(res.Creator.Address), creatorAmount.ToBig()); err != nil {
		return fmt.Errorf("failed to transfer amount to creator: %w", err)
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

	// Mint Bond Certificates
	for _, order := range res.Orders {
		if order.State != string(entity.OrderStateRejected) {
			safeMintPayload, err := abiInterface.Pack(
				"safeMint",
				common.Address(res.BadgeAddress),
				common.Address(order.Investor.Address),
				big.NewInt(1),
				big.NewInt(1),
				[]byte{},
			)
			if err != nil {
				return fmt.Errorf("failed to pack ABI: %w", err)
			}
			env.DelegateCallVoucher(common.Address(h.Config.SafeErc1155MintAddress), safeMintPayload)
		}
	}

	issuance, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte(fmt.Sprintf("issuance %v - ", res.State)), issuance...))
	return nil
}

func (h *IssuanceAdvanceHandlers) SettleIssuance(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input issuance.SettleIssuanceInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	settleIssuance := issuance.NewSettleIssuanceUseCase(
		h.UserRepository,
		h.IssuanceRepository,
		h.OrderRepository,
	)

	res, err := settleIssuance.Execute(&input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to settle issuance: %w", err)
	}

	contractAddr := common.Address(res.Token)
	creatorAddr := common.Address(res.Creator.Address)

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

	// Process settled orders
	for _, order := range res.Orders {
		if order.State == string(entity.OrderStateSettled) {
			// Calculate interest for this order using basis points
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, BasisPointsDivisor)

			// Calculate total payment
			totalPayment := new(uint256.Int).Add(order.Amount, interest)

			if err := env.ERC20Transfer(
				contractAddr,
				creatorAddr,
				common.Address(order.Investor.Address),
				totalPayment.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer settled order: %w", err)
			}

			// Mint Discharge Certificates
			safeMintPayload, err := abiInterface.Pack(
				"safeMint",
				common.Address(res.BadgeAddress),
				common.Address(order.Investor.Address),
				big.NewInt(2),
				big.NewInt(1),
				[]byte{},
			)
			if err != nil {
				return fmt.Errorf("failed to pack ABI: %w", err)
			}
			env.DelegateCallVoucher(common.Address(h.Config.SafeErc1155MintAddress), safeMintPayload)
		}
	}

	issuance, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("issuance settled - "), issuance...))
	return nil
}

func (h *IssuanceAdvanceHandlers) ExecuteIssuanceCollateral(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input issuance.ExecuteIssuanceCollateralInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	executeIssuanceCollateral := issuance.NewExecuteIssuanceCollateralUseCase(h.UserRepository, h.IssuanceRepository, h.OrderRepository)
	res, err := executeIssuanceCollateral.Execute(&input, metadata)
	if err != nil {
		return fmt.Errorf("failed to execute issuance collateral: %w", err)
	}

	totalFinalValue := uint256.NewInt(0)
	orderFinalValues := make(map[uint]*uint256.Int)
	for _, order := range res.Orders {
		if order.State == string(entity.OrderStateSettledByCollateral) {
			// Calculate interest using basis points
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, BasisPointsDivisor)
			finalValue := new(uint256.Int).Add(order.Amount, interest)
			orderFinalValues[order.Id] = finalValue
			totalFinalValue.Add(totalFinalValue, finalValue)
		}
	}

	for _, order := range res.Orders {
		if order.State == string(entity.OrderStateSettledByCollateral) {
			finalValue := orderFinalValues[order.Id]
			orderShare := new(uint256.Int).Mul(finalValue, res.CollateralAmount)
			orderShare.Div(orderShare, totalFinalValue)

			if err = env.ERC20Transfer(
				common.Address(res.CollateralAddress),
				env.AppAddress(),
				common.Address(order.Investor.Address),
				orderShare.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer collateral to investor: %w", err)
			}
		}
	}

	issuance, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("issuance collateral executed - "), issuance...))
	return nil
}
