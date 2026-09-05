package review_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ReviewApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler

}

func (s *ReviewApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupReviewService()
	s.SetupOrderService()
	s.SetupTransactionService()
	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *ReviewApiTestSuite) TestReviewApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllReviews(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllReviews"].(map[string]interface{})["status"])
}

func TestReviewApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(ReviewApiTestSuite))
}
