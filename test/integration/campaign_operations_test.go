package integration

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

func TestCampaignOperationsSuite(t *testing.T) {
	suite.Run(t, new(CampaignOperationsSuite))
}

type CampaignOperationsSuite struct {
	TribesRollupSuite
}

func (s *CampaignOperationsSuite) TestCloseCampaign() {
	badgeId := big.NewInt(1)
	admin, token, creator, deployer, verifier, collateral, badgeAddress := s.setupCommonAddresses()
	investor01, investor02, investor03, investor04, investor05 := s.setupInvestorAddresses()
	baseTime, closesAt, maturityAt := s.setupTimeValues()

	// create creator user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":3,"role":"creator","address":"%s","social_accounts":[],"created_at":%d}`, creator, baseTime)
	s.Equal(string(createUserOutput.Notices[0].Payload), expectedCreateUserOutput)

	// verify social account
	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	createSocialAccountOutput := s.Tester.Advance(verifier, createSocialAccountInput)
	s.Len(createSocialAccountOutput.Notices, 1)

	expectedCreateSocialAccountOutput := fmt.Sprintf(`social account created - {"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}`, baseTime)
	s.Equal(string(createSocialAccountOutput.Notices[0].Payload), expectedCreateSocialAccountOutput)

	// create investors users
	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor01, baseTime)
	s.Equal(string(createUserOutput.Notices[0].Payload), expectedCreateUserOutput)

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor02, baseTime)
	s.Equal(string(createUserOutput.Notices[0].Payload), expectedCreateUserOutput)

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor03))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor03, baseTime)
	s.Equal(string(createUserOutput.Notices[0].Payload), expectedCreateUserOutput)

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor04))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor04, baseTime)
	s.Equal(string(createUserOutput.Notices[0].Payload), expectedCreateUserOutput)

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor05))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor05, baseTime)
	s.Equal(string(createUserOutput.Notices[0].Payload), expectedCreateUserOutput)

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","badge_address":"%s","closes_at":%d,"maturity_at":%d}}`,
		token,
		badgeAddress,
		closesAt,
		maturityAt,
	),
	)
	createCampaignOutput := s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral_address":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(string(createCampaignOutput.Notices[0].Payload), expectedCreateCampaignOutput)

	s.Len(createCampaignOutput.Vouchers, 1)
	s.Equal(deployer, createCampaignOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "deploy2",
		"inputs": [
			{"type": "bytes"},
			{"type": "bytes32"}
		]
	}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(createCampaignOutput.AppContract)
	s.Require().NoError(err)

	initCode := append(s.Bytecode, constructorArgs...)

	unpacked, err := abiInterface.Methods["deploy2"].Inputs.Unpack(createCampaignOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(initCode, unpacked[0])

	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	createOrderOutput := s.Tester.DepositERC20(token, investor01, big.NewInt(60000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"8"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor02, big.NewInt(28000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"4"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor03, big.NewInt(2000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"6"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor04, big.NewInt(5000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"4"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor05, big.NewInt(5500), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	time.Sleep(5 * time.Second)

	anyone := common.HexToAddress("0x0000000000000000000000000000000000000001")
	closeCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/close", "data":{"creator_address":"%s"}}`, creator))
	closeCampaignOutput := s.Tester.Advance(anyone, closeCampaignInput)
	s.Len(closeCampaignOutput.Notices, 1)

	expectedCloseCampaignOutput := fmt.Sprintf(`campaign closed - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral_address":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, closesAt,
		investor02.Hex(), baseTime, closesAt,
		investor03.Hex(), baseTime, closesAt,
		investor04.Hex(), baseTime, closesAt,
		investor05.Hex(), baseTime, closesAt,
		investor01.Hex(), baseTime, closesAt,
		baseTime, closesAt, maturityAt, closesAt)
	s.Equal(string(closeCampaignOutput.Notices[0].Payload), expectedCloseCampaignOutput)

	// Verify final balances after campaign close
	// investor01: deposited 60000, partially accepted 59500, rejected 500
	// investor02: deposited 28000, fully accepted 28000
	// investor03: deposited 2000, fully accepted 2000
	// investor04: deposited 5000, fully accepted 5000
	// investor05: deposited 5500, fully accepted 5500
	// creator: deposited 10000 collateral, received 100000 from investors

	// Verify investor01 balance (60000 - 59500 = 500 rejected should be returned)
	erc20BalanceInput := []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor01.Hex(), token.Hex()))
	erc20BalanceOutput := s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"500"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor02 balance (28000 - 28000 = 0)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor02.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor03 balance (2000 - 2000 = 0)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor03.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor04 balance (5000 - 5000 = 0)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor04.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor05 balance (5500 - 5500 = 0)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor05.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify creator balance (should have received 100000 from investors)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, creator.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"100000"`, string(erc20BalanceOutput.Reports[0].Payload))

	// verify number of vouchers for badge mint calls
	s.Len(closeCampaignOutput.Vouchers, 5)

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

	abiInterface, err = abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	// verify voucher payload for badge mint call (investor01)
	unpacked, err = abiInterface.Methods["mint"].Inputs.Unpack(closeCampaignOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(investor01, unpacked[0])
	s.Equal(badgeId, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal([]byte{}, unpacked[3])

	// verify voucher payload for badge mint call (investor02)
	unpacked, err = abiInterface.Methods["mint"].Inputs.Unpack(closeCampaignOutput.Vouchers[1].Payload[4:])
	s.Require().NoError(err)
	s.Equal(investor02, unpacked[0])
	s.Equal(badgeId, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal([]byte{}, unpacked[3])

	// verify voucher payload for badge mint call (investor03)
	unpacked, err = abiInterface.Methods["mint"].Inputs.Unpack(closeCampaignOutput.Vouchers[2].Payload[4:])
	s.Require().NoError(err)
	s.Equal(investor03, unpacked[0])
	s.Equal(badgeId, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal([]byte{}, unpacked[3])

	// verify voucher payload for badge mint call (investor04)
	unpacked, err = abiInterface.Methods["mint"].Inputs.Unpack(closeCampaignOutput.Vouchers[3].Payload[4:])
	s.Require().NoError(err)
	s.Equal(investor04, unpacked[0])
	s.Equal(badgeId, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal([]byte{}, unpacked[3])

	// verify voucher payload for badge mint call (investor05)
	unpacked, err = abiInterface.Methods["mint"].Inputs.Unpack(closeCampaignOutput.Vouchers[4].Payload[4:])
	s.Require().NoError(err)
	s.Equal(investor05, unpacked[0])
	s.Equal(badgeId, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal([]byte{}, unpacked[3])
}
