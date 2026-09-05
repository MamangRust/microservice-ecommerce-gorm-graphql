package user_test

import (
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type UserApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
	userID  int
}

func (s *UserApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupMerchantService()
	s.SetupOrderService()
	s.SetupTransactionService()

	s.userID = s.SeedUser(s.Ctx)

	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *UserApiTestSuite) TestUserApiLifecycle() {
	// 1. FindById
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findByIdUser(input: { id: `+strconv.Itoa(s.userID)+` }) { status message data { id } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByIdUser"].(map[string]interface{})["status"])

	// 2. FindAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllUsers(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllUsers"].(map[string]interface{})["status"])

	// 3. FindByActive
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByActiveUsers(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByActiveUsers"].(map[string]interface{})["status"])

	// 4. Create
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { createUser(input: { firstname: "Test", lastname: "User", email: "`+strconv.FormatInt(int64(s.userID), 10)+`-new@test.com", password: "password123", confirm_password: "password123" }) { status message data { id } } }`, nil, "")
	s.Require().NoError(err)
	createResult := resp.Data["createUser"].(map[string]interface{})
	s.Equal("success", createResult["status"])
	data := createResult["data"].(map[string]interface{})
	newUserID := int(data["id"].(float64))

	// 5. Update (all required fields)
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { updateUser(input: { id: `+strconv.Itoa(newUserID)+`, firstname: "Updated", lastname: "User", email: "updated@test.com", password: "password123", confirm_password: "password123" }) { status message data { id } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["updateUser"].(map[string]interface{})["status"])

	// 6. Trashed
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedUser(input: { id: `+strconv.Itoa(newUserID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedUser"].(map[string]interface{})["status"])

	// 7. FindByTrashed
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByTrashedUsers(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByTrashedUsers"].(map[string]interface{})["status"])

	// 8. Restore
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreUser(input: { id: `+strconv.Itoa(newUserID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreUser"].(map[string]interface{})["status"])

	// 9. Re-trash + DeletePermanent
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedUser(input: { id: `+strconv.Itoa(newUserID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedUser"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteUserPermanent(input: { id: `+strconv.Itoa(newUserID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteUserPermanent"].(map[string]interface{})["status"])

	// 10. RestoreAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreAllUser { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreAllUser"].(map[string]interface{})["status"])

	// 11. DeleteAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteAllUserPermanent { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteAllUserPermanent"].(map[string]interface{})["status"])
}

func TestUserApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserApiTestSuite))
}
