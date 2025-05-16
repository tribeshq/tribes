package tool

import (
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

var (
	inputBoxAddress = common.HexToAddress(os.Getenv("INPUT_BOX_ADDRESS"))
	portalAddress   = common.HexToAddress(os.Getenv("ERC20_PORTAL_ADDRESS"))
	tokenAddress    = common.HexToAddress(os.Getenv("STABLECOIN_ADDRESS"))
	appAddress      = common.HexToAddress(os.Getenv("APP_ADDRESS"))
	amount          = big.NewInt(0)
)
