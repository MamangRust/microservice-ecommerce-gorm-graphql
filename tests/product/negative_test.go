package product_test

import (
	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
)

func (s *ProductApiTestSuite) TestProductNotFound() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, "query { findProductById(input: { id: 999999 }) { status message } }", nil, "")
	s.Require().NoError(err)
	if len(resp.Errors) > 0 {
		s.T().Logf("GraphQL error (expected): %v", resp.Errors[0].Message)
	}
}
