package integration

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

type UserSuite struct {
	DCMRollupSuite
}

func (s *UserSuite) TestCreateUser() {
	admin, _, creator, _, _, _, _, _ := s.setupCommonAddresses()
	baseTime, _, _ := s.setupTimeValues()

	// create creator user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":3,"role":"creator","address":"%s","social_accounts":[],"created_at":%d}`, creator, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))
}

func (s *UserSuite) TestCreateInvestorUser() {
	admin, _, _, _, _, _, _, _ := s.setupCommonAddresses()
	investor01, _, _, _, _ := s.setupInvestorAddresses()
	baseTime, _, _ := s.setupTimeValues()

	// create investor user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":3,"role":"investor","address":"%s","social_accounts":[],"created_at":%d}`, investor01, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))
}

func (s *UserSuite) TestFindAllUsers() {
	admin, _, creator, _, _, _, _, _ := s.setupCommonAddresses()
	investor01, _, _, _, _ := s.setupInvestorAddresses()
	baseTime, _, _ := s.setupTimeValues()

	// create creator user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	// create investor user
	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	s.Tester.Advance(admin, createUserInput)

	// find all users
	findAllUsersInput := []byte(`{"path":"user"}`)
	findAllUsersOutput := s.Tester.Inspect(findAllUsersInput)
	s.Len(findAllUsersOutput.Reports, 1)

	expectedFindAllUsersOutput := fmt.Sprintf(`[{"id":1,"role":"admin","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},{"id":2,"role":"verifier","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},{"id":3,"role":"creator","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0}]`,
		admin,
		baseTime,
		common.HexToAddress("0x0000000000000000000000000000000000000025"),
		baseTime,
		creator,
		baseTime,
		investor01,
		baseTime)
	s.Equal(expectedFindAllUsersOutput, string(findAllUsersOutput.Reports[0].Payload))
}

func (s *UserSuite) TestFindUserByAddress() {
	admin, _, creator, _, _, _, _, _ := s.setupCommonAddresses()
	baseTime, _, _ := s.setupTimeValues()

	// create creator user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	// find user by address
	findUserByAddressInput := []byte(fmt.Sprintf(`{"path":"user/address","data":{"address":"%s"}}`, creator))
	findUserByAddressOutput := s.Tester.Inspect(findUserByAddressInput)
	s.Len(findUserByAddressOutput.Reports, 1)

	expectedFindUserByAddressOutput := fmt.Sprintf(`{"id":3,"role":"creator","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0}`,
		creator,
		baseTime)
	s.Equal(expectedFindUserByAddressOutput, string(findUserByAddressOutput.Reports[0].Payload))
}

func (s *UserSuite) TestDeleteUser() {
	admin, _, _, _, _, _, _, _ := s.setupCommonAddresses()
	investor01, _, _, _, _ := s.setupInvestorAddresses()

	// create investor user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	s.Tester.Advance(admin, createUserInput)

	// delete user
	deleteUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/delete","data":{"address":"%s"}}`, investor01))
	deleteUserOutput := s.Tester.Advance(admin, deleteUserInput)
	s.Len(deleteUserOutput.Notices, 1)

	expectedDeleteUserOutput := fmt.Sprintf(`user deleted - {"address":"%s"}`, investor01)
	s.Equal(expectedDeleteUserOutput, string(deleteUserOutput.Notices[0].Payload))
}
