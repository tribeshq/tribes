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

func TestCampaignSuite(t *testing.T) {
	suite.Run(t, new(CampaignSuite))
}

type CampaignSuite struct {
	TribesRollupSuite
}

func (s *CampaignSuite) TestCreateCampaign() {
	admin, token, creator, deployer, verifier, collateral, badgeAddress := s.setupCommonAddresses()
	baseTime, closesAt, maturityAt := s.setupTimeValues()

	// create creator user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":3,"role":"creator","address":"%s","social_accounts":[],"created_at":%d}`, creator, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// verify social account
	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	createSocialAccountOutput := s.Tester.Advance(verifier, createSocialAccountInput)
	s.Len(createSocialAccountOutput.Notices, 1)

	expectedCreateSocialAccountOutput := fmt.Sprintf(`social account created - {"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}`, baseTime)
	s.Equal(expectedCreateSocialAccountOutput, string(createSocialAccountOutput.Notices[0].Payload))

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
}

func (s *CampaignSuite) TestFindAllCampaigns() {
	admin, token, creator, deployer, verifier, collateral, badgeAddress := s.setupCommonAddresses()
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

	findAllCampaignsInput := []byte(`{"path":"campaign"}`)

	findAllCampaignsOutput := s.Tester.Inspect(findAllCampaignsInput)
	s.Len(findAllCampaignsOutput.Reports, 1)

	expectedFindAllCampaignsOutput := fmt.Sprintf(`[{"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral_address":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}]`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime,
		closesAt,
		maturityAt,
	)
	s.Equal(string(findAllCampaignsOutput.Reports[0].Payload), expectedFindAllCampaignsOutput)
}

func (s *CampaignSuite) TestFindCampaignById() {
	admin, token, creator, deployer, verifier, collateral, badgeAddress := s.setupCommonAddresses()
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

	findCampaignByIdInput := []byte(`{"path":"campaign/id", "data":{"id":1}}`)

	findCampaignByIdOutput := s.Tester.Inspect(findCampaignByIdInput)
	s.Len(findCampaignByIdOutput.Reports, 1)

	expectedFindCampaignByIdOutput := fmt.Sprintf(`{"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral_address":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(string(findCampaignByIdOutput.Reports[0].Payload), expectedFindCampaignByIdOutput)
}

func (s *CampaignSuite) TestFindCampaignsByCreatorAddress() {
	admin, token, creator, deployer, verifier, collateral, badgeAddress := s.setupCommonAddresses()
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

	findCampaignsByCreatorInput := []byte(fmt.Sprintf(`{"path":"campaign/creator", "data":{"creator_address":"%s"}}`, creator))

	findCampaignsByCreatorOutput := s.Tester.Inspect(findCampaignsByCreatorInput)
	s.Len(findCampaignsByCreatorOutput.Reports, 1)

	expectedFindCampaignsByCreatorAddressOutput := fmt.Sprintf(`[{"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral_address":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}]`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime,
		closesAt,
		maturityAt,
	)
	s.Equal(string(findCampaignsByCreatorOutput.Reports[0].Payload), expectedFindCampaignsByCreatorAddressOutput)
}

func (s *CampaignSuite) TestFindCampaignsByInvestorAddress() {
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
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt)
	s.Equal(string(closeCampaignOutput.Notices[0].Payload), expectedCloseCampaignOutput)

	// Withdraw raised amount
	withdrawRaisedAmountInput := []byte(fmt.Sprintf(`{"path":"user/withdraw","data":{"token":"%s","amount":"100000"}}`, token.Hex()))
	withdrawRaisedAmountOutput := s.Tester.Advance(creator, withdrawRaisedAmountInput)
	s.Len(withdrawRaisedAmountOutput.Notices, 1)

	expectedWithdrawRaisedAmountOutput := fmt.Sprintf(`ERC20 withdrawn - token: %s, amount: 100000, user: %s`, token.Hex(), creator.Hex())
	s.Equal(string(withdrawRaisedAmountOutput.Notices[0].Payload), expectedWithdrawRaisedAmountOutput)

	expectedFindCampaignByCreatorOutput := fmt.Sprintf(`[{"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral_address":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}]`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt)

	findCampaignsByCreatorInput := []byte(fmt.Sprintf(`{"path":"campaign/creator", "data":{"creator_address":"%s"}}`, creator))

	findCampaignsByCreatorOutput := s.Tester.Inspect(findCampaignsByCreatorInput)
	s.Len(findCampaignsByCreatorOutput.Reports, 1)
	s.Equal(string(findCampaignsByCreatorOutput.Reports[0].Payload), expectedFindCampaignByCreatorOutput)
}
