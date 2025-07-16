package integration

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
	"github.com/tribeshq/tribes/assets"
	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/internal/infra/cartesi"
	"github.com/tribeshq/tribes/internal/infra/repository/factory"
)

// TribesRollupSuite is the base suite for all integration tests
type TribesRollupSuite struct {
	suite.Suite
	Bytecode []byte
	Tester   *rollmelette.Tester
}

// SetupTest initializes the test environment
func (s *TribesRollupSuite) SetupTest() {
	cfg, err := configs.LoadRollupConfig()
	if err != nil {
		slog.Error("Failed to load rollup config", "error", err)
		os.Exit(1)
	}

	s.Bytecode, err = assets.GetBadgeBytecode()
	if err != nil {
		slog.Error("Failed to get badge bytecode", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	repo, err := factory.NewRepositoryFromConnectionString(ctx, "sqlite://:memory:")
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database initialized")

	createInfo := cartesi.CreateInfo{
		Repo:   repo,
		Config: cfg,
	}

	dapp := cartesi.Create(&createInfo)
	s.Tester = rollmelette.NewTester(dapp)
}

// setupCommonAddresses returns common addresses used in tests
func (s *TribesRollupSuite) setupCommonAddresses() (
	admin common.Address,
	token common.Address,
	creator common.Address,
	factory common.Address,
	verifier common.Address,
	collateral common.Address,
	safeERC1155MintAddress common.Address,
	applicationAddress common.Address,
) {
	admin = common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	token = common.HexToAddress("0x0000000000000000000000000000000000000009")
	creator = common.HexToAddress("0x0000000000000000000000000000000000000007")
	factory = common.HexToAddress("0x0000000000000000000000000000000000000013")
	verifier = common.HexToAddress("0x0000000000000000000000000000000000000025")
	safeERC1155MintAddress = common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral = common.HexToAddress("0x0000000000000000000000000000000000000008")
	applicationAddress = common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e")
	return
}

// setupInvestorAddresses returns investor addresses for tests
func (s *TribesRollupSuite) setupInvestorAddresses() (
	investor01 common.Address,
	investor02 common.Address,
	investor03 common.Address,
	investor04 common.Address,
	investor05 common.Address,
) {
	investor01 = common.HexToAddress("0x0000000000000000000000000000000000000001")
	investor02 = common.HexToAddress("0x0000000000000000000000000000000000000002")
	investor03 = common.HexToAddress("0x0000000000000000000000000000000000000003")
	investor04 = common.HexToAddress("0x0000000000000000000000000000000000000004")
	investor05 = common.HexToAddress("0x0000000000000000000000000000000000000005")
	return
}

// setupTimeValues returns common time values for tests
func (s *TribesRollupSuite) setupTimeValues() (baseTime int64, closesAt int64, maturityAt int64) {
	baseTime = time.Now().Unix()
	closesAt = baseTime + 5
	maturityAt = baseTime + 10
	return
}
