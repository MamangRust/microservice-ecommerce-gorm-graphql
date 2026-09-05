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

type TransactionStatsTestSuite struct {
	tests.BaseTestSuite
	chContainer testcontainers.Container
	chConn      clickhouse.Conn
}

func (s *TransactionStatsTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.setupTransactionStatsService()
}

func (s *TransactionStatsTestSuite) TearDownSuite() {
	if s.chConn != nil {
		s.chConn.Close()
	}
	if s.chContainer != nil {
		s.chContainer.Terminate(s.Ctx)
	}
	s.BaseTestSuite.TearDownSuite()
}

func (s *TransactionStatsTestSuite) setupTransactionStatsService() {
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

	_ = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS transaction_events (
		transaction_id UInt64, order_id UInt64, merchant_id UInt32,
		payment_method String, amount Int64, status String, created_at DateTime
	) ENGINE = MergeTree() ORDER BY (transaction_id)`)

	s.seedTransactionData(ctx, conn)

	zl, _ := zap.NewDevelopment()
	l := &noopLoggerTxn{zl: zl}
	repo := stats_repo.NewRepository(conn)
	txnHandler := stats_handler.NewTransactionStatsHandler(repo, l)

	server := grpc.NewServer()
	pb.RegisterTransactionStatsServiceServer(server, txnHandler)
	pb.RegisterTransactionStatsByMerchantServiceServer(server, txnHandler)

	addr, err := tests.RunGRPCServer(server)
	s.Require().NoError(err)
	s.Servers = append(s.Servers, server)

	conn2, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["stats_reader"] = conn2
}

type noopLoggerTxn struct{ zl *zap.Logger }

func (l *noopLoggerTxn) Info(msg string, f ...zap.Field)                      { l.zl.Info(msg, f...) }
func (l *noopLoggerTxn) Fatal(msg string, f ...zap.Field)                     { l.zl.Fatal(msg, f...) }
func (l *noopLoggerTxn) Debug(msg string, f ...zap.Field)                     { l.zl.Debug(msg, f...) }
func (l *noopLoggerTxn) Error(msg string, f ...zap.Field)                     { l.zl.Error(msg, f...) }
func (l *noopLoggerTxn) Warn(msg string, f ...zap.Field)                      { l.zl.Warn(msg, f...) }
func (l *noopLoggerTxn) Check(lv zapcore.Level, msg string) *zapcore.CheckedEntry { return nil }
func (l *noopLoggerTxn) With(f ...zap.Field) logger.LoggerInterface           { return l }
func (l *noopLoggerTxn) Sync() error                                          { return nil }

func (s *TransactionStatsTestSuite) seedTransactionData(ctx context.Context, conn clickhouse.Conn) {
	ts := time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC)
	// 5 success + 3 failed transactions
	for i := 0; i < 5; i++ {
		_ = conn.Exec(ctx, `INSERT INTO transaction_events (transaction_id, order_id, merchant_id, payment_method, amount, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			uint64(i+1), uint64(i+1), uint32(1), "bank_transfer", int64(50000), "success", ts)
	}
	for i := 5; i < 8; i++ {
		_ = conn.Exec(ctx, `INSERT INTO transaction_events (transaction_id, order_id, merchant_id, payment_method, amount, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			uint64(i+1), uint64(i+1), uint32(1), "credit_card", int64(30000), "failed", ts)
	}
}

func (s *TransactionStatsTestSuite) graphqlHandler() http.Handler {
	return func() http.Handler { resolver := apigw.NewResolver(s.Conns, s.Log, s.GetCacheStore()); return apigw.NewHandler(resolver) }()
}

func (s *TransactionStatsTestSuite) doQuery(query string, variables map[string]interface{}) map[string]interface{} {
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

func (s *TransactionStatsTestSuite) TestFindMonthStatusSuccess() {
	q := `query ($input: FindMonthlyTransactionStatus!) { findMonthStatusSuccess(input: $input) { status data { year month total_success total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6}})
	d := r["data"].(map[string]interface{})["findMonthStatusSuccess"].(map[string]interface{})
	s.Equal("success", d["status"])
	s.NotEmpty(d["data"])
}

func (s *TransactionStatsTestSuite) TestFindYearStatusSuccess() {
	q := `query ($input: FindYearlyTransactionStatus!) { findYearStatusSuccess(input: $input) { status data { year total_success total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6}})
	d := r["data"].(map[string]interface{})["findYearStatusSuccess"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindMonthStatusFailed() {
	q := `query ($input: FindMonthlyTransactionStatus!) { findMonthStatusFailed(input: $input) { status data { year month total_failed total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6}})
	d := r["data"].(map[string]interface{})["findMonthStatusFailed"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindYearStatusFailed() {
	q := `query ($input: FindYearlyTransactionStatus!) { findYearStatusFailed(input: $input) { status data { year total_failed total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6}})
	d := r["data"].(map[string]interface{})["findYearStatusFailed"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindMonthStatusSuccessByMerchant() {
	q := `query ($input: FindMonthlyTransactionStatusByMerchant!) { findMonthStatusSuccessByMerchant(input: $input) { status data { year month total_success total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findMonthStatusSuccessByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindYearStatusSuccessByMerchant() {
	q := `query ($input: FindYearlyTransactionStatusByMerchant!) { findYearStatusSuccessByMerchant(input: $input) { status data { year total_success total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findYearStatusSuccessByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindMonthStatusFailedByMerchant() {
	q := `query ($input: FindMonthlyTransactionStatusByMerchant!) { findMonthStatusFailedByMerchant(input: $input) { status data { year month total_failed total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findMonthStatusFailedByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindYearStatusFailedByMerchant() {
	q := `query ($input: FindYearlyTransactionStatusByMerchant!) { findYearStatusFailedByMerchant(input: $input) { status data { year total_failed total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findYearStatusFailedByMerchant"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindMonthMethodSuccess() {
	q := `query ($input: MonthTransactionMethod!) { findMonthMethodSuccess(input: $input) { status data { month payment_method total_transactions total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6}})
	d := r["data"].(map[string]interface{})["findMonthMethodSuccess"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindYearMethodSuccess() {
	q := `query ($input: YearTransactionMethod!) { findYearMethodSuccess(input: $input) { status data { year payment_method total_transactions total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025}})
	d := r["data"].(map[string]interface{})["findYearMethodSuccess"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindMonthMethodFailed() {
	q := `query ($input: MonthTransactionMethod!) { findMonthMethodFailed(input: $input) { status data { month payment_method total_transactions total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6}})
	d := r["data"].(map[string]interface{})["findMonthMethodFailed"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindYearMethodFailed() {
	q := `query ($input: YearTransactionMethod!) { findYearMethodFailed(input: $input) { status data { year payment_method total_transactions total_amount } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025}})
	d := r["data"].(map[string]interface{})["findYearMethodFailed"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindMonthMethodByMerchantSuccess() {
	q := `query ($input: MonthTransactionMethodByMerchant!) { findMonthMethodByMerchantSuccess(input: $input) { status data { month payment_method total_transactions } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findMonthMethodByMerchantSuccess"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindYearMethodByMerchantSuccess() {
	q := `query ($input: YearTransactionMethodByMerchant!) { findYearMethodByMerchantSuccess(input: $input) { status data { year payment_method total_transactions } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findYearMethodByMerchantSuccess"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindMonthMethodByMerchantFailed() {
	q := `query ($input: MonthTransactionMethodByMerchant!) { findMonthMethodByMerchantFailed(input: $input) { status data { month payment_method total_transactions } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "month": 6, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findMonthMethodByMerchantFailed"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func (s *TransactionStatsTestSuite) TestFindYearMethodByMerchantFailed() {
	q := `query ($input: YearTransactionMethodByMerchant!) { findYearMethodByMerchantFailed(input: $input) { status data { year payment_method total_transactions } } }`
	r := s.doQuery(q, map[string]interface{}{"input": map[string]interface{}{"year": 2025, "merchant_id": 1}})
	d := r["data"].(map[string]interface{})["findYearMethodByMerchantFailed"].(map[string]interface{})
	s.Equal("success", d["status"])
}

func TestTransactionStatsGraphqlHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsTestSuite))
}
