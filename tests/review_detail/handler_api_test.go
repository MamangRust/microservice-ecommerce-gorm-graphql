package review_detail_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ReviewDetailApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler

}

func (s *ReviewDetailApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupReviewService()
	s.SetupReviewDetailService()
	s.SetupOrderService()
	s.SetupTransactionService()
	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *ReviewDetailApiTestSuite) TestReviewdetailApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllReviewDetails(input: { page: 1, pageSize: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllReviewDetails"].(map[string]interface{})["status"])
}

func TestReviewdetailApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(ReviewDetailApiTestSuite))
}
