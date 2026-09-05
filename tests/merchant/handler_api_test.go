package merchant_test

import (
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantApiTestSuite struct {
	tests.BaseTestSuite
	handler    http.Handler
	merchantID int
	userID     int
}

func (s *MerchantApiTestSuite) SetupSuite() {
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

func (s *MerchantApiTestSuite) TestMerchantApiLifecycle() {
	// 1. Create
	resp, err := graphtest.ExecuteGraphQL(s.handler, `mutation { createMerchant(input: { user_id: `+strconv.Itoa(s.userID)+`, name: "Test Merchant", description: "desc", address: "addr", contact_email: "m@test.com", contact_phone: "123", status: "active" }) { status message data { id } } }`, nil, "")
	s.Require().NoError(err)
	createResult := resp.Data["createMerchant"].(map[string]interface{})
	s.Equal("success", createResult["status"])
	data := createResult["data"].(map[string]interface{})
	s.merchantID = int(data["id"].(float64))

	// 2. FindById
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findMerchantById(input: { id: `+strconv.Itoa(s.merchantID)+` }) { status message data { id name } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findMerchantById"].(map[string]interface{})["status"])

	// 3. FindAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllMerchants(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllMerchants"].(map[string]interface{})["status"])

	// 4. FindByActive
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findActiveMerchants(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findActiveMerchants"].(map[string]interface{})["status"])

	// 5. Update
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { updateMerchant(input: { merchant_id: `+strconv.Itoa(s.merchantID)+`, user_id: `+strconv.Itoa(s.userID)+`, name: "Updated Merchant", description: "desc", address: "addr", contact_email: "m@test.com", contact_phone: "123", status: "active" }) { status message data { id name } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["updateMerchant"].(map[string]interface{})["status"])

	// 6. Trash
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashMerchant(input: { id: `+strconv.Itoa(s.merchantID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashMerchant"].(map[string]interface{})["status"])

	// 7. FindByTrashed
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findTrashedMerchants(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findTrashedMerchants"].(map[string]interface{})["status"])

	// 8. Restore
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreMerchant(input: { id: `+strconv.Itoa(s.merchantID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreMerchant"].(map[string]interface{})["status"])

	// 9. Re-trash + DeletePermanent
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashMerchant(input: { id: `+strconv.Itoa(s.merchantID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteMerchantPermanent(input: { id: `+strconv.Itoa(s.merchantID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteMerchantPermanent"].(map[string]interface{})["status"])

	// 10. RestoreAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreAllMerchants { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreAllMerchants"].(map[string]interface{})["status"])

	// 11. DeleteAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteAllMerchantsPermanent { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteAllMerchantsPermanent"].(map[string]interface{})["status"])
}

func TestMerchantApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantApiTestSuite))
}
