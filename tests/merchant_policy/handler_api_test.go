package merchant_policy_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantPolicyApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler

}

func (s *MerchantPolicyApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupMerchantService()
	s.SetupMerchantPolicyService()
	s.SetupOrderService()
	s.SetupTransactionService()
	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *MerchantPolicyApiTestSuite) TestMerchantpolicyApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllMerchantPolicies(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllMerchantPolicies"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findActiveMerchantPolicies(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findActiveMerchantPolicies"].(map[string]interface{})["status"])
}

func TestMerchantpolicyApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(MerchantPolicyApiTestSuite))
}
