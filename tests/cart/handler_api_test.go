package cart_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CartApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler

}

func (s *CartApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupCartService()
	s.SetupOrderService()
	s.SetupTransactionService()
	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *CartApiTestSuite) TestCartApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllCarts(input: { user_id: 1, page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllCarts"].(map[string]interface{})["status"])
}

func TestCartApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(CartApiTestSuite))
}
