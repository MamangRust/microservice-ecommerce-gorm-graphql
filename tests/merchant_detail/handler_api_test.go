package merchant_detail_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantDetailApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler

}

func (s *MerchantDetailApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupMerchantService()
	s.SetupMerchantDetailService()
	s.SetupOrderService()
	s.SetupTransactionService()
	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *MerchantDetailApiTestSuite) TestMerchantdetailApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllMerchantDetails(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllMerchantDetails"].(map[string]interface{})["status"])
}

func TestMerchantdetailApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(MerchantDetailApiTestSuite))
}
