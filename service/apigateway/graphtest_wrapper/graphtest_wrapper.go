// Package graphtest_wrapper re-exports graphtest for external test packages.
package graphtest_wrapper

import (
	"net/http"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	gt "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest"
	"google.golang.org/grpc"
)

// Resolver is an alias for graphtest.Resolver.
type Resolver = gt.Resolver

// ServiceConnections is an alias for graphtest.ServiceConnections.
type ServiceConnections = gt.ServiceConnections

// ConnMap converts a *ServiceConnections to a map[string]*grpc.ClientConn.
func ConnMap(c *ServiceConnections) map[string]*grpc.ClientConn {
	return gt.ConnMap(c)
}

// NewResolver creates a new Resolver from the given connections.
func NewResolver(conns map[string]*grpc.ClientConn, log logger.LoggerInterface, cacheStore *cache.CacheStore) *Resolver {
	return gt.NewResolver(conns, log, cacheStore)
}

// NewHandler creates a GraphQL HTTP handler from a resolver.
func NewHandler(resolver *Resolver) http.Handler {
	return gt.NewHandler(resolver)
}

// ExecuteGraphQL executes a query against the handler.
func ExecuteGraphQL(handler http.Handler, query string, variables map[string]interface{}, authToken string) (*gt.GraphQLResponse, error) {
	return gt.ExecuteGraphQL(handler, query, variables, authToken)
}
