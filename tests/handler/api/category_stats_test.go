package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	apigw "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	stats_handler "github.com/MamangRust/microservice-ecommerce-grpc-stats-reader/handler"
	stats_repo "github.com/MamangRust/microservice-ecommerce-grpc-stats-reader/repository"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	pb "github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// noopLogger satisfies the logger.LoggerInterface used by stats_reader handlers.
type noopLogger struct{ zl *zap.Logger }

func (l *noopLogger) Info(msg string, f ...zap.Field)  { l.zl.Info(msg, f...) }
func (l *noopLogger) Fatal(msg string, f ...zap.Field) { l.zl.Fatal(msg, f...) }
func (l *noopLogger) Debug(msg string, f ...zap.Field) { l.zl.Debug(msg, f...) }
func (l *noopLogger) Error(msg string, f ...zap.Field) { l.zl.Error(msg, f...) }
func (l *noopLogger) Warn(msg string, f ...zap.Field)  { l.zl.Warn(msg, f...) }
func (l *noopLogger) Check(lv zapcore.Level, msg string) *zapcore.CheckedEntry {
	return nil
}
func (l *noopLogger) With(f ...zap.Field) logger.LoggerInterface { return l }
func (l *noopLogger) Sync() error                        { return nil }

type CategoryStatsTestSuite struct {
	tests.BaseTestSuite
	chContainer testcontainers.Container
	chConn      clickhouse.Conn
}

func (s *CategoryStatsTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.setupCategoryStatsService()
}

func (s *CategoryStatsTestSuite) TearDownSuite() {
	if s.chConn != nil {
		s.chConn.Close()
	}
	if s.chContainer != nil {
		s.chContainer.Terminate(s.Ctx)
	}
	s.BaseTestSuite.TearDownSuite()
}

func (s *CategoryStatsTestSuite) setupCategoryStatsService() {
	ctx := context.Background()

	// Start ClickHouse container
	chContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "clickhouse/clickhouse-server:23.3.8.21-alpine",
			ExposedPorts: []string{"9000/tcp"},
			WaitingFor:   wait.ForListeningPort("9000/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	s.Require().NoError(err)
	s.chContainer = chContainer

	host, err := chContainer.Host(ctx)
	s.Require().NoError(err)
	port, err := chContainer.MappedPort(ctx, "9000")
	s.Require().NoError(err)
	chAddr := fmt.Sprintf("%s:%s", host, port.Port())

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{chAddr},
		Auth: clickhouse.Auth{Database: "default", Username: "default", Password: ""},
	})
	s.Require().NoError(err)
	s.chConn = conn

	// Create tables
	for _, ddl := range categoryDDLs {
		s.Require().NoError(conn.Exec(ctx, ddl))
	}

	// Seed data
	s.seedCategoryData(ctx, conn)

	// Start stats reader gRPC service
	zl, _ := zap.NewDevelopment()
	l := &noopLogger{zl: zl}
	repo := stats_repo.NewRepository(conn)
	catHandler := stats_handler.NewCategoryStatsHandler(repo, l)

	server := grpc.NewServer()
	pb.RegisterCategoryStatsServiceServer(server, catHandler)
	pb.RegisterCategoryStatsByIdServiceServer(server, catHandler)
	pb.RegisterCategoryStatsByMerchantServiceServer(server, catHandler)

	addr, err := tests.RunGRPCServer(server)
	s.Require().NoError(err)
	s.Servers = append(s.Servers, server)

	conn2, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["stats_reader"] = conn2
}

var categoryDDLs = []string{
	`CREATE TABLE IF NOT EXISTS order_events (
		order_id UInt64, merchant_id UInt32, total_price Int64, created_at DateTime
	) ENGINE = MergeTree() ORDER BY (order_id)`,
	`CREATE TABLE IF NOT EXISTS order_item_events (
		order_id UInt64, product_id UInt64, category_id UInt32,
		merchant_id UInt32, category_name String, quantity Int32, price Int64, created_at DateTime
	) ENGINE = MergeTree() ORDER BY (order_id)`,
}

func (s *CategoryStatsTestSuite) seedCategoryData(ctx context.Context, conn clickhouse.Conn) {
	ts := time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC)
	for i := 0; i < 10; i++ {
		_ = conn.Exec(ctx, `INSERT INTO order_events (order_id, merchant_id, total_price, created_at) VALUES (?, ?, ?, ?)`,
			uint64(i+1), uint32(1), int64(50000), ts)
		_ = conn.Exec(ctx, `INSERT INTO order_item_events (order_id, product_id, category_id, merchant_id, category_name, quantity, price, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			uint64(i+1), uint64(i+1), uint32(1), uint32(1), "Electronics", int32(2), int64(25000), ts)
	}
}

func (s *CategoryStatsTestSuite) graphqlHandler() http.Handler {
	return func() http.Handler { resolver := apigw.NewResolver(s.Conns, s.Log, s.GetCacheStore()); return apigw.NewHandler(resolver) }()
}

func (s *CategoryStatsTestSuite) doQuery(query string, variables map[string]interface{}) map[string]interface{} {
	body := map[string]interface{}{"query": query, "variables": variables}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.graphqlHandler().ServeHTTP(rec, req)
	var result map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &result)
	return result
}

func (s *CategoryStatsTestSuite) TestFindMonthPrice() {
	q := `query ($input: FindYearInput!) { findMonthPrice(input: $input) { status data { month category_id total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025}})
	d := r["data"].(map[string]interface{})["findMonthPrice"].(map[string]interface{})
	s.Equal("success", d["status"])
	s.NotEmpty(d["data"])
}

func (s *CategoryStatsTestSuite) TestFindYearPrice() {
	q := `query ($input: FindYearInput!) { findYearPrice(input: $input) { status data { year category_id total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025}})
	d := r["data"].(map[string]interface{})["findYearPrice"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *CategoryStatsTestSuite) TestFindMonthPriceById() {
	q := `query ($input: FindYearCategoryByIdInput!) { findMonthPriceById(input: $input) { status data { month category_id } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "category_id": 1}})
	d := r["data"].(map[string]interface{})["findMonthPriceById"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *CategoryStatsTestSuite) TestFindYearPriceById() {
	q := `query ($input: FindYearCategoryByIdInput!) { findYearPriceById(input: $input) { status data { year category_id } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "category_id": 1}})
	d := r["data"].(map[string]interface{})["findYearPriceById"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *CategoryStatsTestSuite) TestFindMonthPriceByMerchant() {
	q := `query ($input: FindYearCategoryByMerchantInput!) { findMonthPriceByMerchant(input: $input) { status data { month category_id } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findMonthPriceByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *CategoryStatsTestSuite) TestFindYearPriceByMerchant() {
	q := `query ($input: FindYearCategoryByMerchantInput!) { findYearPriceByMerchant(input: $input) { status data { year category_id } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findYearPriceByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *CategoryStatsTestSuite) TestFindMonthlyTotalPrices() {
	q := `query ($input: FindYearMonthTotalPricesInput!) { findMonthlyTotalPrices(input: $input) { status data { year month total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6}})
	d := r["data"].(map[string]interface{})["findMonthlyTotalPrices"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *CategoryStatsTestSuite) TestFindYearlyTotalPrices() {
	q := `query ($input: FindYearTotalPricesInput!) { findYearlyTotalPrices(input: $input) { status data { year total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025}})
	d := r["data"].(map[string]interface{})["findYearlyTotalPrices"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *CategoryStatsTestSuite) TestFindMonthlyTotalPricesById() {
	q := `query ($input: FindYearMonthTotalPriceByIdInput!) { findMonthlyTotalPricesById(input: $input) { status data { year month total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6, "category_id": 1}})
	d := r["data"].(map[string]interface{})["findMonthlyTotalPricesById"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *CategoryStatsTestSuite) TestFindYearlyTotalPricesById() {
	q := `query ($input: FindYearCategoryByIdInput!) { findYearlyTotalPricesById(input: $input) { status data { year total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "category_id": 1}})
	d := r["data"].(map[string]interface{})["findYearlyTotalPricesById"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *CategoryStatsTestSuite) TestFindMonthlyTotalPricesByMerchant() {
	q := `query ($input: FindYearMonthTotalPriceByMerchantInput!) { findMonthlyTotalPricesByMerchant(input: $input) { status data { year month total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findMonthlyTotalPricesByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *CategoryStatsTestSuite) TestFindYearlyTotalPricesByMerchant() {
	q := `query ($input: FindYearTotalPriceByMerchantInput!) { findYearlyTotalPricesByMerchant(input: $input) { status data { year total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findYearlyTotalPricesByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func TestCategoryStatsGraphqlHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryStatsTestSuite))
}
