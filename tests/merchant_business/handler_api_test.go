package merchant_business_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantBusinessApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler

}

func (s *MerchantBusinessApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupMerchantService()
	s.SetupMerchantBusinessService()
	s.SetupOrderService()
	s.SetupTransactionService()
	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *MerchantBusinessApiTestSuite) TestMerchantbusinessApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllMerchantBusinesses(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllMerchantBusinesses"].(map[string]interface{})["status"])
}

func TestMerchantbusinessApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(MerchantBusinessApiTestSuite))
}
