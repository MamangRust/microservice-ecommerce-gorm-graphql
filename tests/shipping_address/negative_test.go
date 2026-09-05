package shipping_address_test

import (
	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
)

func (s *ShippingAddressApiTestSuite) TestShippingNotFound() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, "query { findAllShipping(input: { page: 1, page_size: 10 }) { status message } }", nil, "")
	s.Require().NoError(err)
	if len(resp.Errors) > 0 {
		s.T().Logf("GraphQL error (expected): %v", resp.Errors[0].Message)
	}
}
