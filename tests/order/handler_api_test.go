package order_test

import (
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderApiTestSuite struct {
	tests.BaseTestSuite
	handler  http.Handler
	orderID  int
	userID   int
	merchID  int
	prodID   int
}

func (s *OrderApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupOrderService()
	s.SetupTransactionService()

	s.userID = s.SeedUser(s.Ctx)
	catID := s.SeedCategory(s.Ctx)
	s.merchID = s.SeedMerchant(s.Ctx, s.userID)
	s.prodID = s.SeedProduct(s.Ctx, s.merchID, catID)
	s.orderID = s.SeedOrder(s.Ctx, s.userID, s.merchID, s.prodID)

	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *OrderApiTestSuite) TestOrderApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findOrderById(input: { id: `+strconv.Itoa(s.orderID)+` }) { status message data { id } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findOrderById"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllOrders(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllOrders"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findActiveOrders(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findActiveOrders"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findMonthlyRevenue(input: { year: 2026, month: 1 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.T().Logf("findMonthlyRevenue: %v", resp.Data["findMonthlyRevenue"])
}

func TestOrderApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(OrderApiTestSuite))
}
