package auth_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type AuthApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
}

func (s *AuthApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupMerchantService()
	s.SetupOrderService()
	s.SetupTransactionService()
	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *AuthApiTestSuite) TestAuthApiLifecycle() {
	// Register
	resp, err := graphtest.ExecuteGraphQL(s.handler, `mutation { registerUser(input: { firstname: "Auth", lastname: "Test", email: "auth@test.com", password: "password123", confirm_password: "password123" }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.T().Logf("registerUser: %v", resp.Data["registerUser"])

	// Login
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { loginUser(input: { email: "auth@test.com", password: "password123" }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.T().Logf("loginUser: %v", resp.Data["loginUser"])
}

func TestAuthApiSuite(t *testing.T) {
	if testing.Short() { t.Skip("skipping integration test") }
	suite.Run(t, new(AuthApiTestSuite))
}
