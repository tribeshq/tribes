package integration

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
)

func TestCampaignSuite(t *testing.T) {
	suite.Run(t, new(CampaignSuite))
}

type CampaignSuite struct {
	DCMRollupSuite
}

func (s *CampaignSuite) TestCreateCampaign() {
	admin, token, creator, factory, verifier, collateral, _, applicationAddress := s.setupCommonAddresses()
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

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(applicationAddress)
	s.Require().NoError(err)

	badgeAddress := crypto.CreateAddress2(
		factory,
		common.HexToHash(strconv.Itoa(2)),
		crypto.Keccak256(append(s.Bytecode, constructorArgs...)),
	)

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	createCampaignOutput := s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	s.Len(createCampaignOutput.Vouchers, 1)
	s.Equal(factory, createCampaignOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newBadge",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newBadge"].Inputs.Unpack(createCampaignOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(applicationAddress, unpacked[0])
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(2)), common.BytesToHash(saltBytes[:]))
}

func (s *CampaignSuite) TestFindAllCampaigns() {
	admin, token, creator, factory, verifier, collateral, _, applicationAddress := s.setupCommonAddresses()
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

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(applicationAddress)
	s.Require().NoError(err)

	badgeAddress := crypto.CreateAddress2(
		factory,
		common.HexToHash(strconv.Itoa(2)),
		crypto.Keccak256(append(s.Bytecode, constructorArgs...)),
	)

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	createCampaignOutput := s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	s.Len(createCampaignOutput.Vouchers, 1)
	s.Equal(factory, createCampaignOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newBadge",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newBadge"].Inputs.Unpack(createCampaignOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(applicationAddress, unpacked[0])
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(2)), common.BytesToHash(saltBytes[:]))

	findAllCampaignsInput := []byte(`{"path":"campaign"}`)

	findAllCampaignsOutput := s.Tester.Inspect(findAllCampaignsInput)
	s.Len(findAllCampaignsOutput.Reports, 1)

	expectedFindAllCampaignsOutput := fmt.Sprintf(`[{"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}]`,
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
	s.Equal(expectedFindAllCampaignsOutput, string(findAllCampaignsOutput.Reports[0].Payload))
}

func (s *CampaignSuite) TestFindCampaignById() {
	admin, token, creator, factory, verifier, collateral, _, applicationAddress := s.setupCommonAddresses()
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

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(applicationAddress)
	s.Require().NoError(err)

	badgeAddress := crypto.CreateAddress2(
		factory,
		common.HexToHash(strconv.Itoa(2)),
		crypto.Keccak256(append(s.Bytecode, constructorArgs...)),
	)

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	createCampaignOutput := s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	s.Len(createCampaignOutput.Vouchers, 1)
	s.Equal(factory, createCampaignOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newBadge",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newBadge"].Inputs.Unpack(createCampaignOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(applicationAddress, unpacked[0])
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(2)), common.BytesToHash(saltBytes[:]))

	findCampaignByIdInput := []byte(`{"path":"campaign/id", "data":{"id":1}}`)

	findCampaignByIdOutput := s.Tester.Inspect(findCampaignByIdInput)
	s.Len(findCampaignByIdOutput.Reports, 1)

	expectedFindCampaignByIdOutput := fmt.Sprintf(`{"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(expectedFindCampaignByIdOutput, string(findCampaignByIdOutput.Reports[0].Payload))
}

func (s *CampaignSuite) TestFindCampaignsByCreatorAddress() {
	admin, token, creator, factory, verifier, collateral, _, applicationAddress := s.setupCommonAddresses()
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

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(applicationAddress)
	s.Require().NoError(err)

	badgeAddress := crypto.CreateAddress2(
		factory,
		common.HexToHash(strconv.Itoa(2)),
		crypto.Keccak256(append(s.Bytecode, constructorArgs...)),
	)

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	createCampaignOutput := s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	s.Len(createCampaignOutput.Vouchers, 1)
	s.Equal(factory, createCampaignOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newBadge",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newBadge"].Inputs.Unpack(createCampaignOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(applicationAddress, unpacked[0])
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(2)), common.BytesToHash(saltBytes[:]))

	findCampaignsByCreatorInput := []byte(fmt.Sprintf(`{"path":"campaign/creator", "data":{"creator":"%s"}}`, creator))

	findCampaignsByCreatorOutput := s.Tester.Inspect(findCampaignsByCreatorInput)
	s.Len(findCampaignsByCreatorOutput.Reports, 1)

	expectedFindCampaignsByCreatorAddressOutput := fmt.Sprintf(`[{"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}]`,
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
	s.Equal(expectedFindCampaignsByCreatorAddressOutput, string(findCampaignsByCreatorOutput.Reports[0].Payload))
}
func (s *CampaignSuite) TestFindCampaignsByInvestorAddress() {
	admin, token, creator, factory, verifier, collateral, safeERC1155MintAddress, applicationAddress := s.setupCommonAddresses()
	investor01, investor02, investor03, investor04, investor05 := s.setupInvestorAddresses()
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

	// create investors users
	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor01, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor02, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor03))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor03, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor04))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor04, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor05))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor05, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(applicationAddress)
	s.Require().NoError(err)

	badgeAddress := crypto.CreateAddress2(
		factory,
		common.HexToHash(strconv.Itoa(7)),
		crypto.Keccak256(append(s.Bytecode, constructorArgs...)),
	)

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	createCampaignOutput := s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	s.Len(createCampaignOutput.Vouchers, 1)
	s.Equal(factory, createCampaignOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newBadge",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newBadge"].Inputs.Unpack(createCampaignOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(applicationAddress, unpacked[0])
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(7)), common.BytesToHash(saltBytes[:]))

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

	expectedCloseCampaignOutput := fmt.Sprintf(`campaign closed - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":{"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":{"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":{"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":{"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt)
	s.Equal(expectedCloseCampaignOutput, string(closeCampaignOutput.Notices[0].Payload))

	// Withdraw raised amount
	withdrawRaisedAmountInput := []byte(fmt.Sprintf(`{"path":"user/withdraw","data":{"token":"%s","amount":"100000"}}`, token.Hex()))
	withdrawRaisedAmountOutput := s.Tester.Advance(creator, withdrawRaisedAmountInput)
	s.Len(withdrawRaisedAmountOutput.Notices, 1)

	expectedWithdrawRaisedAmountOutput := fmt.Sprintf(`ERC20 withdrawn - token: %s, amount: 100000, user: %s`, token.Hex(), creator.Hex())
	s.Equal(expectedWithdrawRaisedAmountOutput, string(withdrawRaisedAmountOutput.Notices[0].Payload))

	expectedFindCampaignByCreatorOutput := fmt.Sprintf(`[{"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":{"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":{"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":{"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":{"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}]`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt)

	findCampaignsByCreatorInput := []byte(fmt.Sprintf(`{"path":"campaign/creator", "data":{"creator":"%s"}}`, creator))

	findCampaignsByCreatorOutput := s.Tester.Inspect(findCampaignsByCreatorInput)
	s.Len(findCampaignsByCreatorOutput.Reports, 1)
	s.Equal(expectedFindCampaignByCreatorOutput, string(findCampaignsByCreatorOutput.Reports[0].Payload))

	// Verify that delegate call vouchers were created for badge minting
	s.Len(closeCampaignOutput.DelegateCallVouchers, 5)

	// Verify delegate call voucher destinations are SafeCall contract
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[0].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[1].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[2].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[3].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[4].Destination)
}

func (s *CampaignSuite) TestCloseCampaign() {
	admin, token, creator, factory, verifier, collateral, safeERC1155MintAddress, applicationAddress := s.setupCommonAddresses()
	investor01, investor02, investor03, investor04, investor05 := s.setupInvestorAddresses()
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

	// create investors users
	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor01, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor02, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor03))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor03, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor04))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor04, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor05))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor05, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(applicationAddress)
	s.Require().NoError(err)

	badgeAddress := crypto.CreateAddress2(
		factory,
		common.HexToHash(strconv.Itoa(7)),
		crypto.Keccak256(append(s.Bytecode, constructorArgs...)),
	)

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	createCampaignOutput := s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	s.Len(createCampaignOutput.Vouchers, 1)
	s.Equal(factory, createCampaignOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newBadge",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newBadge"].Inputs.Unpack(createCampaignOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(applicationAddress, unpacked[0])
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(7)), common.BytesToHash(saltBytes[:]))

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

	expectedCloseCampaignOutput := fmt.Sprintf(`campaign closed - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":{"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":{"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":{"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":{"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, baseTime, closesAt,
		investor02.Hex(), baseTime, baseTime, closesAt,
		investor03.Hex(), baseTime, baseTime, closesAt,
		investor04.Hex(), baseTime, baseTime, closesAt,
		investor05.Hex(), baseTime, baseTime, closesAt,
		investor01.Hex(), baseTime, baseTime, closesAt,
		baseTime, closesAt, maturityAt, closesAt)
	s.Equal(expectedCloseCampaignOutput, string(closeCampaignOutput.Notices[0].Payload))

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

	// verify number of vouchers for badge safeERC1155MintAddress delegate calls
	s.Len(closeCampaignOutput.DelegateCallVouchers, 5)

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

	abiInterface, err = abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	// verify delegate call voucher destination is SafeCall contract
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[0].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[1].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[2].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[3].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[4].Destination)

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor01)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor01, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor02)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[1].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor02, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor03)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[2].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor03, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor04)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[3].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor04, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor05)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[4].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor05, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])
}

func (s *CampaignSuite) TestExecuteCampaignCollateral() {
	admin, token, creator, factory, verifier, collateral, safeERC1155MintAddress, applicationAddress := s.setupCommonAddresses()
	investor01, investor02, investor03, investor04, investor05 := s.setupInvestorAddresses()
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

	// create investors users
	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor01, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor02, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor03))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor03, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor04))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor04, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor05))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor05, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(applicationAddress)
	s.Require().NoError(err)

	badgeAddress := crypto.CreateAddress2(
		factory,
		common.HexToHash(strconv.Itoa(7)),
		crypto.Keccak256(append(s.Bytecode, constructorArgs...)),
	)

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	createCampaignOutput := s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	s.Len(createCampaignOutput.Vouchers, 1)
	s.Equal(factory, createCampaignOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newBadge",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newBadge"].Inputs.Unpack(createCampaignOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(applicationAddress, unpacked[0])
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(7)), common.BytesToHash(saltBytes[:]))

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

	expectedCloseCampaignOutput := fmt.Sprintf(`campaign closed - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":{"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":{"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":{"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":{"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, baseTime, closesAt,
		investor02.Hex(), baseTime, baseTime, closesAt,
		investor03.Hex(), baseTime, baseTime, closesAt,
		investor04.Hex(), baseTime, baseTime, closesAt,
		investor05.Hex(), baseTime, baseTime, closesAt,
		investor01.Hex(), baseTime, baseTime, closesAt,
		baseTime, closesAt, maturityAt, closesAt)
	s.Equal(expectedCloseCampaignOutput, string(closeCampaignOutput.Notices[0].Payload))

	// Withdraw raised amount
	withdrawRaisedAmountInput := []byte(fmt.Sprintf(`{"path":"user/withdraw","data":{"token":"%s","amount":"100000"}}`, token.Hex()))
	withdrawRaisedAmountOutput := s.Tester.Advance(creator, withdrawRaisedAmountInput)
	s.Len(withdrawRaisedAmountOutput.Notices, 1)

	expectedWithdrawRaisedAmountOutput := fmt.Sprintf(`ERC20 withdrawn - token: %s, amount: 100000, user: %s`, token.Hex(), creator.Hex())
	s.Equal(expectedWithdrawRaisedAmountOutput, string(withdrawRaisedAmountOutput.Notices[0].Payload))

	findCampaignByIdInput := []byte(`{"path":"campaign/id", "data":{"id":1}}`)

	findCampaignByIdOutput := s.Tester.Inspect(findCampaignByIdInput)
	s.Len(findCampaignByIdOutput.Reports, 1)

	expectedFindCampaignByCreatorOutput := fmt.Sprintf(`[{"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":{"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":{"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":{"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":{"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}]`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, baseTime, closesAt,
		investor02.Hex(), baseTime, baseTime, closesAt,
		investor03.Hex(), baseTime, baseTime, closesAt,
		investor04.Hex(), baseTime, baseTime, closesAt,
		investor05.Hex(), baseTime, baseTime, closesAt,
		investor01.Hex(), baseTime, baseTime, closesAt,
		baseTime, closesAt, maturityAt, closesAt)

	findCampaignsByCreatorInput := []byte(fmt.Sprintf(`{"path":"campaign/creator", "data":{"creator":"%s"}}`, creator))

	findCampaignsByCreatorOutput := s.Tester.Inspect(findCampaignsByCreatorInput)
	s.Len(findCampaignsByCreatorOutput.Reports, 1)
	s.Equal(expectedFindCampaignByCreatorOutput, string(findCampaignsByCreatorOutput.Reports[0].Payload))

	time.Sleep(6 * time.Second)

	executeCampaignCollateralInput := []byte(`{"path":"campaign/execute-collateral", "data":{"id":1}}`)
	executeCampaignCollateralOutput := s.Tester.Advance(creator, executeCampaignCollateralInput)
	s.Len(executeCampaignCollateralOutput.Notices, 1)

	updatedAt := baseTime + 11

	expectedExecuteCampaignCollateralOutput := fmt.Sprintf(`campaign collateral executed - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"collateral_executed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"59500","interest_rate":"9","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":{"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"28000","interest_rate":"8","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":{"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"2000","interest_rate":"4","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":{"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5000","interest_rate":"6","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":{"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5500","interest_rate":"4","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, baseTime, updatedAt,
		investor02.Hex(), baseTime, baseTime, updatedAt,
		investor03.Hex(), baseTime, baseTime, updatedAt,
		investor04.Hex(), baseTime, baseTime, updatedAt,
		investor05.Hex(), baseTime, baseTime, updatedAt,
		investor01.Hex(), baseTime, baseTime, closesAt,
		baseTime, closesAt, maturityAt, updatedAt)
	s.Equal(expectedExecuteCampaignCollateralOutput, string(executeCampaignCollateralOutput.Notices[0].Payload))

	// Verify final balances after campaign collateral execution
	// The collateral (10000) is distributed proportionally to accepted orders based on their final value
	// Total final value = 59500*1.09 + 28000*1.08 + 2000*1.04 + 5000*1.06 + 5500*1.04 = 64855 + 30240 + 2080 + 5300 + 5720 = 108195
	// investor01: 64855/108195 * 10000 = 5994 (rounded down)
	// investor02: 30240/108195 * 10000 = 2794 (rounded down)
	// investor03: 2080/108195 * 10000 = 192 (rounded down)
	// investor04: 5300/108195 * 10000 = 489 (rounded down)
	// investor05: 5720/108195 * 10000 = 528 (rounded down)
	// Total distributed: 5994 + 2794 + 192 + 489 + 528 = 9997
	// Remaining: 10000 - 9997 = 3 tokens remain in the application (not distributed)
	// Final distribution:
	// investor01: 5994, investor02: 2794, investor03: 192, investor04: 489, investor05: 528
	// creator: no additional deposit, just execution of existing collateral

	// Verify investor01 balance (received 5994 collateral)
	erc20BalanceInput := []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor01.Hex(), collateral.Hex()))
	erc20BalanceOutput := s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"5994"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor02 balance (received 2794 collateral)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor02.Hex(), collateral.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"2794"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor03 balance (received 192 collateral)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor03.Hex(), collateral.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"192"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor04 balance (received 489 collateral)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor04.Hex(), collateral.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"489"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor05 balance (received 528 collateral)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor05.Hex(), collateral.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"528"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify creator balance (no additional deposit, just execution of existing collateral)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, creator.Hex(), collateral.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))

	// verify number of vouchers for badge safeERC1155MintAddress delegate calls
	s.Len(closeCampaignOutput.DelegateCallVouchers, 5)

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

	abiInterface, err = abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	// verify delegate call voucher destination is SafeCall contract
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[0].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[1].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[2].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[3].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[4].Destination)

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor01)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor01, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor02)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[1].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor02, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor03)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[2].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor03, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor04)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[3].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor04, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor05)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[4].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor05, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])
}

func (s *CampaignSuite) TestSettleCampaign() {
	admin, token, creator, factory, verifier, collateral, safeERC1155MintAddress, applicationAddress := s.setupCommonAddresses()
	investor01, investor02, investor03, investor04, investor05 := s.setupInvestorAddresses()
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

	// create investors users
	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor01, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor02, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor03))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor03, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor04))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor04, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor05))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor05, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(applicationAddress)
	s.Require().NoError(err)

	badgeAddress := crypto.CreateAddress2(
		factory,
		common.HexToHash(strconv.Itoa(7)),
		crypto.Keccak256(append(s.Bytecode, constructorArgs...)),
	)

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	createCampaignOutput := s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	s.Len(createCampaignOutput.Vouchers, 1)
	s.Equal(factory, createCampaignOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newBadge",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newBadge"].Inputs.Unpack(createCampaignOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(applicationAddress, unpacked[0])

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

	expectedCloseCampaignOutput := fmt.Sprintf(`campaign closed - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":{"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":{"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":{"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":{"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, baseTime, closesAt,
		investor02.Hex(), baseTime, baseTime, closesAt,
		investor03.Hex(), baseTime, baseTime, closesAt,
		investor04.Hex(), baseTime, baseTime, closesAt,
		investor05.Hex(), baseTime, baseTime, closesAt,
		investor01.Hex(), baseTime, baseTime, closesAt,
		baseTime, closesAt, maturityAt, closesAt)
	s.Equal(expectedCloseCampaignOutput, string(closeCampaignOutput.Notices[0].Payload))

	// Withdraw raised amount
	withdrawRaisedAmountInput := []byte(fmt.Sprintf(`{"path":"user/withdraw","data":{"token":"%s","amount":"100000"}}`, token.Hex()))
	withdrawRaisedAmountOutput := s.Tester.Advance(creator, withdrawRaisedAmountInput)
	s.Len(withdrawRaisedAmountOutput.Notices, 1)

	expectedWithdrawRaisedAmountOutput := fmt.Sprintf(`ERC20 withdrawn - token: %s, amount: 100000, user: %s`, token.Hex(), creator.Hex())
	s.Equal(expectedWithdrawRaisedAmountOutput, string(withdrawRaisedAmountOutput.Notices[0].Payload))

	time.Sleep(5 * time.Second)

	settleCampaignInput := []byte(`{"path":"campaign/creator/settle", "data":{"id":1}}`)
	settleCampaignOutput := s.Tester.DepositERC20(token, creator, big.NewInt(108195), settleCampaignInput)
	s.Len(settleCampaignOutput.Notices, 1)

	settledAt := baseTime + 10

	expectedSettleCampaignOutput := fmt.Sprintf(`campaign settled - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"settled","orders":[`+
		`{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"59500","interest_rate":"9","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":{"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"28000","interest_rate":"8","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":{"id":6,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"2000","interest_rate":"4","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":{"id":7,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5000","interest_rate":"6","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":{"id":8,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"5500","interest_rate":"4","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, baseTime, settledAt,
		investor02.Hex(), baseTime, baseTime, settledAt,
		investor03.Hex(), baseTime, baseTime, settledAt,
		investor04.Hex(), baseTime, baseTime, settledAt,
		investor05.Hex(), baseTime, baseTime, settledAt,
		investor01.Hex(), baseTime, baseTime, closesAt,
		baseTime, closesAt, maturityAt, settledAt)
	s.Equal(expectedSettleCampaignOutput, string(settleCampaignOutput.Notices[0].Payload))

	// Verify final balances after campaign settlement
	// investor01: should receive 59500 + (59500 * 9% = 5355) = 64855
	// investor02: should receive 28000 + (28000 * 8% = 2240) = 30240
	// investor03: should receive 2000 + (2000 * 4% = 80) = 2080
	// investor04: should receive 5000 + (5000 * 6% = 300) = 5300
	// investor05: should receive 5500 + (5500 * 4% = 220) = 5720
	// creator: paid 108195 to settle the campaign

	// Verify investor01 balance (received 64855 + rejected order amount = 65355)
	erc20BalanceInput := []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor01.Hex(), token.Hex()))
	erc20BalanceOutput := s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"65355"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor02 balance (received 30240)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor02.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"30240"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor03 balance (received 2080)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor03.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"2080"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor04 balance (received 5300)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor04.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"5300"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor05 balance (received 5720)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, investor05.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"5720"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify creator balance (had 100000, paid 108195, so should be -8195)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/balance","data":{"address":"%s","token":"%s"}}`, creator.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))

	// verify number of vouchers for badge safeERC1155MintAddress delegate calls
	s.Len(closeCampaignOutput.DelegateCallVouchers, 5)

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

	abiInterface, err = abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	// verify delegate call voucher destination is SafeCall contract
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[0].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[1].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[2].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[3].Destination)
	s.Equal(safeERC1155MintAddress, closeCampaignOutput.DelegateCallVouchers[4].Destination)

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor01)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor01, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor02)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[1].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor02, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor03)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[2].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor03, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor04)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[3].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor04, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])

	// verify delegate call voucher payload for badge safeERC1155MintAddress (investor05)
	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(closeCampaignOutput.DelegateCallVouchers[4].Payload[4:])
	s.Require().NoError(err)
	s.Equal(badgeAddress, unpacked[0]) // target is badgeAddress
	s.Equal(investor05, unpacked[1])
	s.Equal(big.NewInt(1), unpacked[2])
	s.Equal(big.NewInt(1), unpacked[3])
	s.Equal([]byte{}, unpacked[4])
}
