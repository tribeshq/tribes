package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSocialAccountSuite(t *testing.T) {
	suite.Run(t, new(SocialAccountSuite))
}

type SocialAccountSuite struct {
	DCMRollupSuite
}

func (s *SocialAccountSuite) TestCreateSocialAccount() {
	admin, _, creator, _, verifier, _, _, _ := s.setupCommonAddresses()
	baseTime, _, _ := s.setupTimeValues()

	// create creator user first
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	// create social account
	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	createSocialAccountOutput := s.Tester.Advance(verifier, createSocialAccountInput)
	s.Len(createSocialAccountOutput.Notices, 1)

	expectedCreateSocialAccountOutput := fmt.Sprintf(`social account created - {"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}`, baseTime)
	s.Equal(expectedCreateSocialAccountOutput, string(createSocialAccountOutput.Notices[0].Payload))
}

func (s *SocialAccountSuite) TestFindSocialAccountById() {
	admin, _, creator, _, verifier, _, _, _ := s.setupCommonAddresses()
	baseTime, _, _ := s.setupTimeValues()

	// create creator user first
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	// create social account
	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	s.Tester.Advance(verifier, createSocialAccountInput)

	// find social account by id
	findSocialAccountByIdInput := []byte(`{"path":"social/id","data":{"social_account_id":1}}`)
	findSocialAccountByIdOutput := s.Tester.Inspect(findSocialAccountByIdInput)
	s.Len(findSocialAccountByIdOutput.Reports, 1)

	expectedFindSocialAccountByIdOutput := fmt.Sprintf(`{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d,"updated_at":0}`, baseTime)
	s.Equal(expectedFindSocialAccountByIdOutput, string(findSocialAccountByIdOutput.Reports[0].Payload))
}

func (s *SocialAccountSuite) TestFindSocialAccountsByUserId() {
	admin, _, creator, _, verifier, _, _, _ := s.setupCommonAddresses()
	baseTime, _, _ := s.setupTimeValues()

	// create creator user first
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	// create social account
	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	s.Tester.Advance(verifier, createSocialAccountInput)

	// create another social account for the same user
	createSocialAccountInput = []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test2","platform":"instagram"}}`, creator))
	s.Tester.Advance(verifier, createSocialAccountInput)

	// find social accounts by user id
	findSocialAccountsByUserIdInput := []byte(`{"path":"social/user/id","data":{"user_id":3}}`)
	findSocialAccountsByUserIdOutput := s.Tester.Inspect(findSocialAccountsByUserIdInput)
	s.Len(findSocialAccountsByUserIdOutput.Reports, 1)

	expectedFindSocialAccountsByUserIdOutput := fmt.Sprintf(`[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d,"updated_at":0},{"id":2,"user_id":3,"username":"test2","platform":"instagram","created_at":%d,"updated_at":0}]`, baseTime, baseTime)
	s.Equal(expectedFindSocialAccountsByUserIdOutput, string(findSocialAccountsByUserIdOutput.Reports[0].Payload))
}

func (s *SocialAccountSuite) TestDeleteSocialAccount() {
	admin, _, creator, _, verifier, _, _, _ := s.setupCommonAddresses()

	// create creator user first
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	// create social account
	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	s.Tester.Advance(verifier, createSocialAccountInput)

	// delete social account
	deleteSocialAccountInput := []byte(`{"path":"social/admin/delete","data":{"social_account_id":1}}`)
	deleteSocialAccountOutput := s.Tester.Advance(admin, deleteSocialAccountInput)
	s.Len(deleteSocialAccountOutput.Notices, 1)

	expectedDeleteSocialAccountOutput := `social account deleted - {"social_account_id":1}`
	s.Equal(expectedDeleteSocialAccountOutput, string(deleteSocialAccountOutput.Notices[0].Payload))
}
