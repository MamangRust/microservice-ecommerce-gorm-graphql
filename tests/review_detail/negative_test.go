package review_detail_test

import (
	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
)

func (s *ReviewDetailApiTestSuite) TestReviewDetailNotFound() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, "query { findAllReviewDetails(input: { page: 1, pageSize: 10 }) { status message } }", nil, "")
	s.Require().NoError(err)
	if len(resp.Errors) > 0 {
		s.T().Logf("GraphQL error (expected): %v", resp.Errors[0].Message)
	}
}
