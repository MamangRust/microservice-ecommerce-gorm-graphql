package slider_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type SliderApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler

}

func (s *SliderApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupSliderService()
	s.SetupOrderService()
	s.SetupTransactionService()
	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *SliderApiTestSuite) TestSliderApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllSliders(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllSliders"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findActiveSliders(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findActiveSliders"].(map[string]interface{})["status"])
}

func TestSliderApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(SliderApiTestSuite))
}
