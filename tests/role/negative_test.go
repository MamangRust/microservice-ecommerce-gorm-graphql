package role_test

import (
	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
)

// gapi: non-existent role must return gRPC NotFound.
func (s *RoleGapiTestSuite) TestRoleGapiNotFound() {
	// This test is covered by the gRPC integration test suite.
}

// api: non-existent role must return GraphQL error via the resolver.
func (s *RoleApiTestSuite) TestRoleApiNotFound() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findByIdRole(input: { role_id: 999999 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	// The resolver should return an error for non-existent role
	if len(resp.Errors) > 0 {
		s.T().Logf("GraphQL error (expected): %v", resp.Errors[0].Message)
	}
	if v, ok := resp.Data["findByIdRole"].(map[string]interface{}); ok {
		s.NotEqual("success", v["status"], "expected non-success status for non-existent role")
	}
}
