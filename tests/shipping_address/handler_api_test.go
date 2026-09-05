package shipping_address_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ShippingAddressApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler

}

func (s *ShippingAddressApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupMerchantService()
	s.SetupOrderService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()
	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *ShippingAddressApiTestSuite) TestShippingApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllShipping(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllShipping"].(map[string]interface{})["status"])
}

func TestShippingAddressApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(ShippingAddressApiTestSuite))
}
