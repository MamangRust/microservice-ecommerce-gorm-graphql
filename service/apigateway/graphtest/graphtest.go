package graphtest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-pkg/upload_image"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	pb "github.com/MamangRust/microservice-ecommerce-shared/pb"
	graph "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/internal/handler"
	graphqlmapper "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/internal/mapper"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/grpc"
)

// Resolver is an alias for the internal graph.Resolver type.
type Resolver = graph.Resolver

// ServiceConnections is an alias for the internal ServiceConnections type.
type ServiceConnections = graph.ServiceConnections

// ConnMap converts a *ServiceConnections to a map[string]*grpc.ClientConn.
func ConnMap(c *ServiceConnections) map[string]*grpc.ClientConn {
	if c == nil {
		return nil
	}
	m := map[string]*grpc.ClientConn{
		"auth":               c.AuthClient,
		"role":               c.RoleClient,
		"user":               c.UserClient,
		"category":           c.CategoryClient,
		"merchant":           c.MerchantClient,
		"order-item":         c.OrderItemClient,
		"order":              c.OrderClient,
		"product":            c.ProductClient,
		"transaction":        c.TransactionClient,
		"cart":               c.CartClient,
		"review":             c.ReviewClient,
		"slider":             c.SliderClient,
		"shipping-address":   c.ShippingClient,
		"banner":             c.BannerClient,
		"merchant_award":     c.MerchantAwardClient,
		"merchant_business":  c.MerchantBusinessClient,
		"merchant_detail":    c.MerchantDetailClient,
		"merchant_policy":    c.MerchantPolicyClient,
		"review-detail":      c.ReviewDetailClient,
		"merchant-social":    c.MerchantSocialLinkClient,
		"stats_reader":       c.StatsReaderClient,
	}
	return m
}

// NewResolver creates a new Resolver from the given connections, logger, and cache.
func NewResolver(conns map[string]*grpc.ClientConn, log logger.LoggerInterface, cacheStore *cache.CacheStore) *Resolver {
	mapper := graphqlmapper.NewGraphqlMapper()
	imageUpload := upload_image.NewImageUpload(log)
	clients := buildGRPCClients(conns)

	return graph.NewResolver(&graph.Deps{
		Clients:     clients,
		Logger:      log,
		Mapping:     mapper,
		Cache:       cacheStore,
		ImageUpload: imageUpload,
	})
}

// NewHandler creates a GraphQL HTTP handler from a resolver.
func NewHandler(resolver *Resolver) *handler.Server {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	return srv
}

// GraphQLQuery represents a GraphQL request payload.
type GraphQLQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse wraps the standard GraphQL response structure.
type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []GraphQLError         `json:"errors,omitempty"`
}

// GraphQLError represents a single error in the GraphQL response.
type GraphQLError struct {
	Message string `json:"message"`
}

// ExecuteGraphQL executes a query against the handler.
func ExecuteGraphQL(srv http.Handler, query string, variables map[string]interface{}, authToken string) (*GraphQLResponse, error) {
	gqlReq := GraphQLQuery{
		Query:     query,
		Variables: variables,
	}
	body, err := json.Marshal(gqlReq)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	respBody, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		return nil, err
	}

	var resp GraphQLResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// buildGRPCClients creates GRPCClients from a connection map.
func buildGRPCClients(conns map[string]*grpc.ClientConn) *graph.GRPCClients {
	get := func(key string) *grpc.ClientConn {
		if c, ok := conns[key]; ok && c != nil {
			return c
		}
		return nil
	}

	return &graph.GRPCClients{
		AuthClient:                       newAuthServiceClient(get("auth")),
		RoleCommandClient:                newRoleCommandServiceClient(get("role")),
		RoleQueryClient:                  newRoleQueryServiceClient(get("role")),
		UserCommandClient:                newUserCommandServiceClient(get("user")),
		UserQueryClient:                  newUserQueryServiceClient(get("user")),
		BannerCommandClient:              newBannerCommandServiceClient(get("banner")),
		BannerQueryClient:                newBannerQueryServiceClient(get("banner")),
		CartCommandClient:                newCartCommandServiceClient(get("cart")),
		CartQueryClient:                  newCartQueryServiceClient(get("cart")),
		CategoryCommandClient:            newCategoryCommandServiceClient(get("category")),
		CategoryQueryClient:              newCategoryQueryServiceClient(get("category")),
		CategoryStatsClient:              newCategoryStatsServiceClient(get("stats_reader")),
		CategoryStatsByMerchantClient:    newCategoryStatsByMerchantServiceClient(get("stats_reader")),
		CategoryStatsByIdClient:          newCategoryStatsByIdServiceClient(get("stats_reader")),
		MerchantCommandClient:            newMerchantCommandServiceClient(get("merchant")),
		MerchantQueryClient:              newMerchantQueryServiceClient(get("merchant")),
		MerchantAwardCommandClient:       newMerchantAwardCommandServiceClient(get("merchant_award")),
		MerchantAwardQueryClient:         newMerchantAwardQueryServiceClient(get("merchant_award")),
		MerchantBusinessCommandClient:    newMerchantBusinessCommandServiceClient(get("merchant_business")),
		MerchantBusinessQueryClient:      newMerchantBusinessQueryServiceClient(get("merchant_business")),
		MerchantDetailCommandClient:      newMerchantDetailCommandServiceClient(get("merchant_detail")),
		MerchantDetailQueryClient:        newMerchantDetailQueryServiceClient(get("merchant_detail")),
		MerchantPolicyCommandClient:      newMerchantPolicyCommandServiceClient(get("merchant_policy")),
		MerchantPolicyQueryClient:        newMerchantPolicyQueryServiceClient(get("merchant_policy")),
		MerchantSocialLinkClient:         newMerchantSocialCommandServiceClient(get("merchant-social")),
		OrderCommandClient:               newOrderCommandServiceClient(get("order")),
		OrderQueryClient:                 newOrderQueryServiceClient(get("order")),
		OrderStatsClient:                 newOrderStatsServiceClient(get("stats_reader")),
		OrderItemCommandClient:           newOrderItemCommandServiceClient(get("order-item")),
		OrderItemQueryClient:             newOrderItemQueryServiceClient(get("order-item")),
		ProductCommandClient:             newProductCommandServiceClient(get("product")),
		ProductQueryClient:               newProductQueryServiceClient(get("product")),
		ReviewCommandClient:              newReviewCommandServiceClient(get("review")),
		ReviewQueryClient:                newReviewQueryServiceClient(get("review")),
		ReviewDetailCommandClient:        newReviewDetailCommandServiceClient(get("review-detail")),
		ReviewDetailQueryClient:          newReviewDetailQueryServiceClient(get("review-detail")),
		ShippingCommandClient:            newShippingCommandServiceClient(get("shipping-address")),
		ShippingQueryClient:              newShippingQueryServiceClient(get("shipping-address")),
		SliderCommandClient:              newSliderCommandServiceClient(get("slider")),
		SliderQueryClient:                newSliderQueryServiceClient(get("slider")),
		TransactionCommandClient:         newTransactionCommandServiceClient(get("transaction")),
		TransactionQueryClient:           newTransactionQueryServiceClient(get("transaction")),
		TransactionStatsClient:           newTransactionStatsServiceClient(get("stats_reader")),
		TransactionStatsByMerchantClient: newTransactionStatsByMerchantServiceClient(get("stats_reader")),
	}
}

// Nil-safe gRPC client constructors

func newAuthServiceClient(c *grpc.ClientConn) pb.AuthServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewAuthServiceClient(c)
}

func newRoleCommandServiceClient(c *grpc.ClientConn) pb.RoleCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewRoleCommandServiceClient(c)
}

func newRoleQueryServiceClient(c *grpc.ClientConn) pb.RoleQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewRoleQueryServiceClient(c)
}

func newUserCommandServiceClient(c *grpc.ClientConn) pb.UserCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewUserCommandServiceClient(c)
}

func newUserQueryServiceClient(c *grpc.ClientConn) pb.UserQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewUserQueryServiceClient(c)
}

func newBannerCommandServiceClient(c *grpc.ClientConn) pb.BannerCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewBannerCommandServiceClient(c)
}

func newBannerQueryServiceClient(c *grpc.ClientConn) pb.BannerQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewBannerQueryServiceClient(c)
}

func newCartCommandServiceClient(c *grpc.ClientConn) pb.CartCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewCartCommandServiceClient(c)
}

func newCartQueryServiceClient(c *grpc.ClientConn) pb.CartQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewCartQueryServiceClient(c)
}

func newCategoryCommandServiceClient(c *grpc.ClientConn) pb.CategoryCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewCategoryCommandServiceClient(c)
}

func newCategoryQueryServiceClient(c *grpc.ClientConn) pb.CategoryQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewCategoryQueryServiceClient(c)
}

func newCategoryStatsServiceClient(c *grpc.ClientConn) pb.CategoryStatsServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewCategoryStatsServiceClient(c)
}

func newCategoryStatsByMerchantServiceClient(c *grpc.ClientConn) pb.CategoryStatsByMerchantServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewCategoryStatsByMerchantServiceClient(c)
}

func newCategoryStatsByIdServiceClient(c *grpc.ClientConn) pb.CategoryStatsByIdServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewCategoryStatsByIdServiceClient(c)
}

func newMerchantCommandServiceClient(c *grpc.ClientConn) pb.MerchantCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantCommandServiceClient(c)
}

func newMerchantQueryServiceClient(c *grpc.ClientConn) pb.MerchantQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantQueryServiceClient(c)
}

func newMerchantAwardCommandServiceClient(c *grpc.ClientConn) pb.MerchantAwardCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantAwardCommandServiceClient(c)
}

func newMerchantAwardQueryServiceClient(c *grpc.ClientConn) pb.MerchantAwardQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantAwardQueryServiceClient(c)
}

func newMerchantBusinessCommandServiceClient(c *grpc.ClientConn) pb.MerchantBusinessCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantBusinessCommandServiceClient(c)
}

func newMerchantBusinessQueryServiceClient(c *grpc.ClientConn) pb.MerchantBusinessQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantBusinessQueryServiceClient(c)
}

func newMerchantDetailCommandServiceClient(c *grpc.ClientConn) pb.MerchantDetailCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantDetailCommandServiceClient(c)
}

func newMerchantDetailQueryServiceClient(c *grpc.ClientConn) pb.MerchantDetailQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantDetailQueryServiceClient(c)
}

func newMerchantPolicyCommandServiceClient(c *grpc.ClientConn) pb.MerchantPolicyCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantPolicyCommandServiceClient(c)
}

func newMerchantPolicyQueryServiceClient(c *grpc.ClientConn) pb.MerchantPolicyQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantPolicyQueryServiceClient(c)
}

func newMerchantSocialCommandServiceClient(c *grpc.ClientConn) pb.MerchantSocialCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewMerchantSocialCommandServiceClient(c)
}

func newOrderCommandServiceClient(c *grpc.ClientConn) pb.OrderCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewOrderCommandServiceClient(c)
}

func newOrderQueryServiceClient(c *grpc.ClientConn) pb.OrderQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewOrderQueryServiceClient(c)
}

func newOrderStatsServiceClient(c *grpc.ClientConn) pb.OrderStatsServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewOrderStatsServiceClient(c)
}

func newOrderItemCommandServiceClient(c *grpc.ClientConn) pb.OrderItemCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewOrderItemCommandServiceClient(c)
}

func newOrderItemQueryServiceClient(c *grpc.ClientConn) pb.OrderItemQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewOrderItemQueryServiceClient(c)
}

func newProductCommandServiceClient(c *grpc.ClientConn) pb.ProductCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewProductCommandServiceClient(c)
}

func newProductQueryServiceClient(c *grpc.ClientConn) pb.ProductQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewProductQueryServiceClient(c)
}

func newReviewCommandServiceClient(c *grpc.ClientConn) pb.ReviewCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewReviewCommandServiceClient(c)
}

func newReviewQueryServiceClient(c *grpc.ClientConn) pb.ReviewQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewReviewQueryServiceClient(c)
}

func newReviewDetailCommandServiceClient(c *grpc.ClientConn) pb.ReviewDetailCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewReviewDetailCommandServiceClient(c)
}

func newReviewDetailQueryServiceClient(c *grpc.ClientConn) pb.ReviewDetailQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewReviewDetailQueryServiceClient(c)
}

func newShippingCommandServiceClient(c *grpc.ClientConn) pb.ShippingCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewShippingCommandServiceClient(c)
}

func newShippingQueryServiceClient(c *grpc.ClientConn) pb.ShippingQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewShippingQueryServiceClient(c)
}

func newSliderCommandServiceClient(c *grpc.ClientConn) pb.SliderCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewSliderCommandServiceClient(c)
}

func newSliderQueryServiceClient(c *grpc.ClientConn) pb.SliderQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewSliderQueryServiceClient(c)
}

func newTransactionCommandServiceClient(c *grpc.ClientConn) pb.TransactionCommandServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewTransactionCommandServiceClient(c)
}

func newTransactionQueryServiceClient(c *grpc.ClientConn) pb.TransactionQueryServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewTransactionQueryServiceClient(c)
}

func newTransactionStatsServiceClient(c *grpc.ClientConn) pb.TransactionStatsServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewTransactionStatsServiceClient(c)
}

func newTransactionStatsByMerchantServiceClient(c *grpc.ClientConn) pb.TransactionStatsByMerchantServiceClient {
	if c == nil {
		return nil
	}
	return pb.NewTransactionStatsByMerchantServiceClient(c)
}
