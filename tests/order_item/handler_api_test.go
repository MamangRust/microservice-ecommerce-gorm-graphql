package order_item_test

import (
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderItemApiTestSuite struct {
	tests.BaseTestSuite
	handler  http.Handler
	orderID  int
	userID   int
	merchID  int
	prodID   int
}

func (s *OrderItemApiTestSuite) SetupSuite() {
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

func (s *OrderItemApiTestSuite) TestOrderItemApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllOrderItems(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllOrderItems"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findActiveOrderItems(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findActiveOrderItems"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findOrderItemsByOrder(input: { id: `+strconv.Itoa(s.orderID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findOrderItemsByOrder"].(map[string]interface{})["status"])
}

func TestOrderItemApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(OrderItemApiTestSuite))
}
