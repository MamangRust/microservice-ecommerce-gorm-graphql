package cart_test

import (
	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
)

func (s *CartApiTestSuite) TestCartNotFound() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, "query { findAllCarts(input: { user_id: 999999, page: 1, page_size: 10 }) { status message } }", nil, "")
	s.Require().NoError(err)
	if len(resp.Errors) > 0 {
		s.T().Logf("GraphQL error (expected): %v", resp.Errors[0].Message)
	}
}
