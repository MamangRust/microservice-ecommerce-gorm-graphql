package transaction_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type TransactionApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler

}

func (s *TransactionApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *TransactionApiTestSuite) TestTransactionApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllTransaction(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllTransaction"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findMonthStatusSuccess(input: { year: 2026, month: 1 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.T().Logf("findMonthStatusSuccess: %v", resp.Data["findMonthStatusSuccess"])
}

func TestTransactionApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(TransactionApiTestSuite))
}
