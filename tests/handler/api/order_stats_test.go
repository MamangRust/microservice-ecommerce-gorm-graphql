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
	apigw "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	stats_handler "github.com/MamangRust/microservice-ecommerce-grpc-stats-reader/handler"
	stats_repo "github.com/MamangRust/microservice-ecommerce-grpc-stats-reader/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
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

type OrderStatsTestSuite struct {
	tests.BaseTestSuite
	chContainer testcontainers.Container
	chConn      clickhouse.Conn
}

func (s *OrderStatsTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.setupOrderStatsService()
}

func (s *OrderStatsTestSuite) TearDownSuite() {
	if s.chConn != nil {
		s.chConn.Close()
	}
	if s.chContainer != nil {
		s.chContainer.Terminate(s.Ctx)
	}
	s.BaseTestSuite.TearDownSuite()
}

func (s *OrderStatsTestSuite) setupOrderStatsService() {
	ctx := context.Background()

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

	for _, ddl := range orderDDLs {
		s.Require().NoError(conn.Exec(ctx, ddl))
	}
	s.seedOrderData(ctx, conn)

	zl, _ := zap.NewDevelopment()
	l := &noopLoggerOrder{zl: zl}
	repo := stats_repo.NewRepository(conn)
	orderHandler := stats_handler.NewOrderStatsHandler(repo, l)

	server := grpc.NewServer()
	pb.RegisterOrderStatsServiceServer(server, orderHandler)

	addr, err := tests.RunGRPCServer(server)
	s.Require().NoError(err)
	s.Servers = append(s.Servers, server)

	conn2, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["stats_reader"] = conn2
}

type noopLoggerOrder struct{ zl *zap.Logger }

func (l *noopLoggerOrder) Info(msg string, f ...zap.Field)                      { l.zl.Info(msg, f...) }
func (l *noopLoggerOrder) Fatal(msg string, f ...zap.Field)                     { l.zl.Fatal(msg, f...) }
func (l *noopLoggerOrder) Debug(msg string, f ...zap.Field)                     { l.zl.Debug(msg, f...) }
func (l *noopLoggerOrder) Error(msg string, f ...zap.Field)                     { l.zl.Error(msg, f...) }
func (l *noopLoggerOrder) Warn(msg string, f ...zap.Field)                      { l.zl.Warn(msg, f...) }
func (l *noopLoggerOrder) Check(lv zapcore.Level, msg string) *zapcore.CheckedEntry { return nil }
func (l *noopLoggerOrder) With(f ...zap.Field) logger.LoggerInterface           { return l }
func (l *noopLoggerOrder) Sync() error                                          { return nil }

var orderDDLs = []string{
	`CREATE TABLE IF NOT EXISTS order_events (
		order_id UInt64, merchant_id UInt32, total_price Int64, created_at DateTime
	) ENGINE = MergeTree() ORDER BY (order_id)`,
	`CREATE TABLE IF NOT EXISTS order_item_events (
		order_id UInt64, product_id UInt64, category_id UInt32,
		merchant_id UInt32, category_name String, quantity Int32, price Int64, created_at DateTime
	) ENGINE = MergeTree() ORDER BY (order_id)`,
}

func (s *OrderStatsTestSuite) seedOrderData(ctx context.Context, conn clickhouse.Conn) {
	ts := time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC)
	for i := 0; i < 20; i++ {
		_ = conn.Exec(ctx, `INSERT INTO order_events (order_id, merchant_id, total_price, created_at) VALUES (?, ?, ?, ?)`,
			uint64(i+1), uint32(1), int64(100000), ts)
		_ = conn.Exec(ctx, `INSERT INTO order_item_events (order_id, product_id, category_id, merchant_id, category_name, quantity, price, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			uint64(i+1), uint64(i+1), uint32(1), uint32(1), "Books", int32(3), int64(50000), ts)
	}
}

func (s *OrderStatsTestSuite) graphqlHandler() http.Handler {
	return func() http.Handler { resolver := apigw.NewResolver(s.Conns, s.Log, s.GetCacheStore()); return apigw.NewHandler(resolver) }()
}

func (s *OrderStatsTestSuite) doQuery(query string, variables map[string]interface{}) map[string]interface{} {
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

func (s *OrderStatsTestSuite) TestFindMonthlyRevenue() {
	q := `query ($input: FindMonthYearOrderInput!) { findMonthlyRevenue(input: $input) { status data { month order_count total_revenue total_items_sold } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6}})
	d := r["data"].(map[string]interface{})["findMonthlyRevenue"].(map[string]interface{})
	s.Equal("success", d["status"])
	s.NotEmpty(d["data"])
}

func (s *OrderStatsTestSuite) TestFindYearlyRevenue() {
	q := `query ($input: FindYearOrderInput!) { findYearlyRevenue(input: $input) { status data { year order_count total_revenue total_items_sold } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025}})
	d := r["data"].(map[string]interface{})["findYearlyRevenue"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *OrderStatsTestSuite) TestFindMonthlyRevenueByMerchant() {
	q := `query ($input: FindYearOrderByMerchantInput!) { findMonthlyRevenueByMerchant(input: $input) { status data { month order_count } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findMonthlyRevenueByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *OrderStatsTestSuite) TestFindYearlyRevenueByMerchant() {
	q := `query ($input: FindYearOrderByMerchantInput!) { findYearlyRevenueByMerchant(input: $input) { status data { year order_count } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findYearlyRevenueByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *OrderStatsTestSuite) TestFindMonthlyTotalRevenue() {
	q := `query ($input: FindYearMonthTotalRevenue!) { findMonthlyTotalRevenue(input: $input) { status data { year month total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6}})
	d := r["data"].(map[string]interface{})["findMonthlyTotalRevenue"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *OrderStatsTestSuite) TestFindYearlyTotalRevenue() {
	q := `query ($input: FindYearTotalRevenue!) { findYearlyTotalRevenue(input: $input) { status data { year total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025}})
	d := r["data"].(map[string]interface{})["findYearlyTotalRevenue"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *OrderStatsTestSuite) TestFindMonthlyTotalRevenueByMerchant() {
	q := `query ($input: FindYearMonthTotalRevenueByMerchant!) { findMonthlyTotalRevenueByMerchant(input: $input) { status data { year month total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findMonthlyTotalRevenueByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *OrderStatsTestSuite) TestFindYearlyTotalRevenueByMerchant() {
	q := `query ($input: FindYearTotalRevenueByMerchant!) { findYearlyTotalRevenueByMerchant(input: $input) { status data { year total_revenue } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findYearlyTotalRevenueByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func TestOrderStatsGraphqlHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsTestSuite))
}
