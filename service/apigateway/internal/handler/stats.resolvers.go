package graph

// This file contains all stats-related GraphQL resolvers, following the same
// pattern as microservice-payment-gateway-grpc/service/apigateway/internal/handler/stats.resolvers.go.
// Stats resolvers are centralized here rather than spread across category/order/transaction files.

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/internal/model"
)

// ─── Category Stats ───────────────────────────────────────────────────────────

func (r *queryResolver) FindMonthPrice(ctx context.Context, input model.FindYearInput) (*model.APIResponseCategoryMonthPrice, error) {
	res, err := r.StatsRead.CategoryStats.FindMonthPrice(ctx, &pb.FindYearCategory{Year: input.Year})
	if err != nil {
		return &model.APIResponseCategoryMonthPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlCategoryMonthlyPrice(res), nil
}

func (r *queryResolver) FindYearPrice(ctx context.Context, input model.FindYearInput) (*model.APIResponseCategoryYearPrice, error) {
	res, err := r.StatsRead.CategoryStats.FindYearPrice(ctx, &pb.FindYearCategory{Year: input.Year})
	if err != nil {
		return &model.APIResponseCategoryYearPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlCategoryYearlyPrice(res), nil
}

func (r *queryResolver) FindMonthPriceByID(ctx context.Context, input model.FindYearCategoryByIDInput) (*model.APIResponseCategoryMonthPrice, error) {
	res, err := r.StatsRead.CategoryStatsById.FindMonthPriceById(ctx, &pb.FindYearCategoryById{
		Year:       input.Year,
		CategoryId: input.CategoryID,
	})
	if err != nil {
		return &model.APIResponseCategoryMonthPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlCategoryMonthlyPrice(res), nil
}

func (r *queryResolver) FindYearPriceByID(ctx context.Context, input model.FindYearCategoryByIDInput) (*model.APIResponseCategoryYearPrice, error) {
	res, err := r.StatsRead.CategoryStatsById.FindYearPriceById(ctx, &pb.FindYearCategoryById{
		Year:       input.Year,
		CategoryId: input.CategoryID,
	})
	if err != nil {
		return &model.APIResponseCategoryYearPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlCategoryYearlyPrice(res), nil
}

func (r *queryResolver) FindMonthPriceByMerchant(ctx context.Context, input model.FindYearCategoryByMerchantInput) (*model.APIResponseCategoryMonthPrice, error) {
	res, err := r.StatsRead.CategoryStatsByMerchant.FindMonthPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year:       input.Year,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseCategoryMonthPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlCategoryMonthlyPrice(res), nil
}

func (r *queryResolver) FindYearPriceByMerchant(ctx context.Context, input model.FindYearCategoryByMerchantInput) (*model.APIResponseCategoryYearPrice, error) {
	res, err := r.StatsRead.CategoryStatsByMerchant.FindYearPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year:       input.Year,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseCategoryYearPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlCategoryYearlyPrice(res), nil
}

func (r *queryResolver) FindMonthlyTotalPrices(ctx context.Context, input model.FindYearMonthTotalPricesInput) (*model.APIResponseCategoryMonthlyTotalPrice, error) {
	res, err := r.StatsRead.CategoryStats.FindMonthlyTotalPrices(ctx, &pb.FindYearMonthTotalPrices{
		Year:  input.Year,
		Month: input.Month,
	})
	if err != nil {
		return &model.APIResponseCategoryMonthlyTotalPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlMonthlyTotalPrice(res), nil
}

func (r *queryResolver) FindYearlyTotalPrices(ctx context.Context, input model.FindYearTotalPricesInput) (*model.APIResponseCategoryYearlyTotalPrice, error) {
	res, err := r.StatsRead.CategoryStats.FindYearlyTotalPrices(ctx, &pb.FindYearTotalPrices{Year: input.Year})
	if err != nil {
		return &model.APIResponseCategoryYearlyTotalPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlYearlyTotalPrice(res), nil
}

func (r *queryResolver) FindMonthlyTotalPricesByID(ctx context.Context, input model.FindYearMonthTotalPriceByIDInput) (*model.APIResponseCategoryMonthlyTotalPrice, error) {
	res, err := r.StatsRead.CategoryStatsById.FindMonthlyTotalPricesById(ctx, &pb.FindYearMonthTotalPriceById{
		Year:       input.Year,
		Month:      input.Month,
		CategoryId: input.CategoryID,
	})
	if err != nil {
		return &model.APIResponseCategoryMonthlyTotalPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlMonthlyTotalPrice(res), nil
}

func (r *queryResolver) FindYearlyTotalPricesByID(ctx context.Context, input model.FindYearCategoryByIDInput) (*model.APIResponseCategoryYearlyTotalPrice, error) {
	res, err := r.StatsRead.CategoryStatsById.FindYearlyTotalPricesById(ctx, &pb.FindYearTotalPriceById{
		Year:       input.Year,
		CategoryId: input.CategoryID,
	})
	if err != nil {
		return &model.APIResponseCategoryYearlyTotalPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlYearlyTotalPrice(res), nil
}

func (r *queryResolver) FindMonthlyTotalPricesByMerchant(ctx context.Context, input model.FindYearMonthTotalPriceByMerchantInput) (*model.APIResponseCategoryMonthlyTotalPrice, error) {
	res, err := r.StatsRead.CategoryStatsByMerchant.FindMonthlyTotalPricesByMerchant(ctx, &pb.FindYearMonthTotalPriceByMerchant{
		Year:       input.Year,
		Month:      input.Month,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseCategoryMonthlyTotalPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlMonthlyTotalPrice(res), nil
}

func (r *queryResolver) FindYearlyTotalPricesByMerchant(ctx context.Context, input model.FindYearTotalPriceByMerchantInput) (*model.APIResponseCategoryYearlyTotalPrice, error) {
	res, err := r.StatsRead.CategoryStatsByMerchant.FindYearlyTotalPricesByMerchant(ctx, &pb.FindYearTotalPriceByMerchant{
		Year:       input.Year,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseCategoryYearlyTotalPrice{Status: "error", Message: err.Error()}, nil
	}
	return r.CategoryGraphql.Mapping.ToGraphqlYearlyTotalPrice(res), nil
}

// ─── Order Stats ──────────────────────────────────────────────────────────────

func (r *queryResolver) FindMonthlyRevenue(ctx context.Context, input model.FindMonthYearOrderInput) (*model.APIResponseOrderMonthly, error) {
	res, err := r.StatsRead.OrderStats.FindMonthlyRevenue(ctx, &pb.FindYearOrder{Year: input.Year})
	if err != nil {
		return &model.APIResponseOrderMonthly{Status: "error", Message: err.Error()}, nil
	}
	return r.OrderGraphql.Mapping.ToGraphqlResponseOrderMonthlyRevenue(res), nil
}

func (r *queryResolver) FindYearlyRevenue(ctx context.Context, input model.FindYearOrderInput) (*model.APIResponseOrderYearly, error) {
	res, err := r.StatsRead.OrderStats.FindYearlyRevenue(ctx, &pb.FindYearOrder{Year: input.Year})
	if err != nil {
		return &model.APIResponseOrderYearly{Status: "error", Message: err.Error()}, nil
	}
	return r.OrderGraphql.Mapping.ToGraphqlResponseOrderYearlyRevenue(res), nil
}

func (r *queryResolver) FindMonthlyRevenueByMerchant(ctx context.Context, input model.FindYearOrderByMerchantInput) (*model.APIResponseOrderMonthly, error) {
	res, err := r.StatsRead.OrderStats.FindMonthlyRevenueByMerchant(ctx, &pb.FindYearOrderByMerchant{
		Year:       input.Year,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseOrderMonthly{Status: "error", Message: err.Error()}, nil
	}
	return r.OrderGraphql.Mapping.ToGraphqlResponseOrderMonthlyRevenue(res), nil
}

func (r *queryResolver) FindYearlyRevenueByMerchant(ctx context.Context, input model.FindYearOrderByMerchantInput) (*model.APIResponseOrderYearly, error) {
	res, err := r.StatsRead.OrderStats.FindYearlyRevenueByMerchant(ctx, &pb.FindYearOrderByMerchant{
		Year:       input.Year,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseOrderYearly{Status: "error", Message: err.Error()}, nil
	}
	return r.OrderGraphql.Mapping.ToGraphqlResponseOrderYearlyRevenue(res), nil
}

func (r *queryResolver) FindMonthlyTotalRevenue(ctx context.Context, input model.FindYearMonthTotalRevenue) (*model.APIResponseOrderMonthlyTotalRevenue, error) {
	res, err := r.StatsRead.OrderStats.FindMonthlyTotalRevenue(ctx, &pb.FindYearMonthTotalRevenue{
		Year:  input.Year,
		Month: input.Month,
	})
	if err != nil {
		return &model.APIResponseOrderMonthlyTotalRevenue{Status: "error", Message: err.Error()}, nil
	}
	return r.OrderGraphql.Mapping.ToGraphqlResponseOrderMonthlyTotalRevenue(res), nil
}

func (r *queryResolver) FindYearlyTotalRevenue(ctx context.Context, input model.FindYearTotalRevenue) (*model.APIResponseOrderYearlyTotalRevenue, error) {
	res, err := r.StatsRead.OrderStats.FindYearlyTotalRevenue(ctx, &pb.FindYearTotalRevenue{Year: input.Year})
	if err != nil {
		return &model.APIResponseOrderYearlyTotalRevenue{Status: "error", Message: err.Error()}, nil
	}
	return r.OrderGraphql.Mapping.ToGraphqlResponseOrderYearlyTotalRevenue(res), nil
}

func (r *queryResolver) FindMonthlyTotalRevenueByMerchant(ctx context.Context, input model.FindYearMonthTotalRevenueByMerchant) (*model.APIResponseOrderMonthlyTotalRevenue, error) {
	res, err := r.StatsRead.OrderStats.FindMonthlyTotalRevenueByMerchant(ctx, &pb.FindYearMonthTotalRevenueByMerchant{
		Year:       input.Year,
		Month:      input.Month,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseOrderMonthlyTotalRevenue{Status: "error", Message: err.Error()}, nil
	}
	return r.OrderGraphql.Mapping.ToGraphqlResponseOrderMonthlyTotalRevenue(res), nil
}

func (r *queryResolver) FindYearlyTotalRevenueByMerchant(ctx context.Context, input model.FindYearTotalRevenueByMerchant) (*model.APIResponseOrderYearlyTotalRevenue, error) {
	res, err := r.StatsRead.OrderStats.FindYearlyTotalRevenueByMerchant(ctx, &pb.FindYearTotalRevenueByMerchant{
		Year:       input.Year,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseOrderYearlyTotalRevenue{Status: "error", Message: err.Error()}, nil
	}
	return r.OrderGraphql.Mapping.ToGraphqlResponseOrderYearlyTotalRevenue(res), nil
}

// ─── Transaction Stats Amount ─────────────────────────────────────────────────

func (r *queryResolver) FindMonthStatusSuccess(ctx context.Context, input model.FindMonthlyTransactionStatus) (*model.APIResponseTransactionMonthAmountSuccess, error) {
	res, err := r.StatsRead.TransactionStats.GetMonthlyAmountSuccess(ctx, &pb.MonthAmountTransactionRequest{
		Year:  input.Year,
		Month: input.Month,
	})
	if err != nil {
		return &model.APIResponseTransactionMonthAmountSuccess{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseMonthAmountSuccess(res), nil
}

func (r *queryResolver) FindYearStatusSuccess(ctx context.Context, input model.FindYearlyTransactionStatus) (*model.APIResponseTransactionYearAmountSuccess, error) {
	res, err := r.StatsRead.TransactionStats.GetYearlyAmountSuccess(ctx, &pb.YearAmountTransactionRequest{Year: input.Year})
	if err != nil {
		return &model.APIResponseTransactionYearAmountSuccess{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseYearAmountSuccess(res), nil
}

func (r *queryResolver) FindMonthStatusFailed(ctx context.Context, input model.FindMonthlyTransactionStatus) (*model.APIResponseTransactionMonthAmountFailed, error) {
	res, err := r.StatsRead.TransactionStats.GetMonthlyAmountFailed(ctx, &pb.MonthAmountTransactionRequest{
		Year:  input.Year,
		Month: input.Month,
	})
	if err != nil {
		return &model.APIResponseTransactionMonthAmountFailed{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseMonthAmountFailed(res), nil
}

func (r *queryResolver) FindYearStatusFailed(ctx context.Context, input model.FindYearlyTransactionStatus) (*model.APIResponseTransactionYearAmountFailed, error) {
	res, err := r.StatsRead.TransactionStats.GetYearlyAmountFailed(ctx, &pb.YearAmountTransactionRequest{Year: input.Year})
	if err != nil {
		return &model.APIResponseTransactionYearAmountFailed{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseYearAmountFailed(res), nil
}

// ─── Transaction Stats Amount By Merchant ─────────────────────────────────────

func (r *queryResolver) FindMonthStatusSuccessByMerchant(ctx context.Context, input model.FindMonthlyTransactionStatusByMerchant) (*model.APIResponseTransactionMonthAmountSuccess, error) {
	res, err := r.StatsRead.TransactionStatsByMerchant.GetMonthlyAmountSuccessByMerchant(ctx, &pb.MonthAmountTransactionMerchantRequest{
		Year:       input.Year,
		Month:      input.Month,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseTransactionMonthAmountSuccess{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseMonthAmountSuccess(res), nil
}

func (r *queryResolver) FindYearStatusSuccessByMerchant(ctx context.Context, input model.FindYearlyTransactionStatusByMerchant) (*model.APIResponseTransactionYearAmountSuccess, error) {
	res, err := r.StatsRead.TransactionStatsByMerchant.GetYearlyAmountSuccessByMerchant(ctx, &pb.YearAmountTransactionMerchantRequest{
		Year:       input.Year,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseTransactionYearAmountSuccess{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseYearAmountSuccess(res), nil
}

func (r *queryResolver) FindMonthStatusFailedByMerchant(ctx context.Context, input model.FindMonthlyTransactionStatusByMerchant) (*model.APIResponseTransactionMonthAmountFailed, error) {
	res, err := r.StatsRead.TransactionStatsByMerchant.GetMonthlyAmountFailedByMerchant(ctx, &pb.MonthAmountTransactionMerchantRequest{
		Year:       input.Year,
		Month:      input.Month,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseTransactionMonthAmountFailed{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseMonthAmountFailed(res), nil
}

func (r *queryResolver) FindYearStatusFailedByMerchant(ctx context.Context, input model.FindYearlyTransactionStatusByMerchant) (*model.APIResponseTransactionYearAmountFailed, error) {
	res, err := r.StatsRead.TransactionStatsByMerchant.GetYearlyAmountFailedByMerchant(ctx, &pb.YearAmountTransactionMerchantRequest{
		Year:       input.Year,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseTransactionYearAmountFailed{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseYearAmountFailed(res), nil
}

// ─── Transaction Stats Method ─────────────────────────────────────────────────

func (r *queryResolver) FindMonthMethodSuccess(ctx context.Context, input model.MonthTransactionMethod) (*model.APIResponseTransactionMonthPaymentMethod, error) {
	res, err := r.StatsRead.TransactionStats.GetMonthlyTransactionMethodSuccess(ctx, &pb.MonthMethodTransactionRequest{
		Year:  input.Year,
		Month: input.Month,
	})
	if err != nil {
		return &model.APIResponseTransactionMonthPaymentMethod{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseMonthMethod(res), nil
}

func (r *queryResolver) FindYearMethodSuccess(ctx context.Context, input model.YearTransactionMethod) (*model.APIResponseTransactionYearPaymentMethod, error) {
	res, err := r.StatsRead.TransactionStats.GetYearlyTransactionMethodSuccess(ctx, &pb.YearMethodTransactionRequest{Year: input.Year})
	if err != nil {
		return &model.APIResponseTransactionYearPaymentMethod{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseYearMethod(res), nil
}

func (r *queryResolver) FindMonthMethodByMerchantSuccess(ctx context.Context, input model.MonthTransactionMethodByMerchant) (*model.APIResponseTransactionMonthPaymentMethod, error) {
	res, err := r.StatsRead.TransactionStatsByMerchant.GetMonthlyTransactionMethodByMerchantSuccess(ctx, &pb.MonthMethodTransactionMerchantRequest{
		Year:       input.Year,
		Month:      input.Month,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseTransactionMonthPaymentMethod{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseMonthMethod(res), nil
}

func (r *queryResolver) FindYearMethodByMerchantSuccess(ctx context.Context, input model.YearTransactionMethodByMerchant) (*model.APIResponseTransactionYearPaymentMethod, error) {
	res, err := r.StatsRead.TransactionStatsByMerchant.GetYearlyTransactionMethodByMerchantSuccess(ctx, &pb.YearMethodTransactionMerchantRequest{
		Year:       input.Year,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseTransactionYearPaymentMethod{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseYearMethod(res), nil
}

func (r *queryResolver) FindMonthMethodFailed(ctx context.Context, input model.MonthTransactionMethod) (*model.APIResponseTransactionMonthPaymentMethod, error) {
	res, err := r.StatsRead.TransactionStats.GetMonthlyTransactionMethodFailed(ctx, &pb.MonthMethodTransactionRequest{
		Year:  input.Year,
		Month: input.Month,
	})
	if err != nil {
		return &model.APIResponseTransactionMonthPaymentMethod{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseMonthMethod(res), nil
}

func (r *queryResolver) FindYearMethodFailed(ctx context.Context, input model.YearTransactionMethod) (*model.APIResponseTransactionYearPaymentMethod, error) {
	res, err := r.StatsRead.TransactionStats.GetYearlyTransactionMethodFailed(ctx, &pb.YearMethodTransactionRequest{Year: input.Year})
	if err != nil {
		return &model.APIResponseTransactionYearPaymentMethod{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseYearMethod(res), nil
}

func (r *queryResolver) FindMonthMethodByMerchantFailed(ctx context.Context, input model.MonthTransactionMethodByMerchant) (*model.APIResponseTransactionMonthPaymentMethod, error) {
	res, err := r.StatsRead.TransactionStatsByMerchant.GetMonthlyTransactionMethodByMerchantFailed(ctx, &pb.MonthMethodTransactionMerchantRequest{
		Year:       input.Year,
		Month:      input.Month,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseTransactionMonthPaymentMethod{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseMonthMethod(res), nil
}

func (r *queryResolver) FindYearMethodByMerchantFailed(ctx context.Context, input model.YearTransactionMethodByMerchant) (*model.APIResponseTransactionYearPaymentMethod, error) {
	res, err := r.StatsRead.TransactionStatsByMerchant.GetYearlyTransactionMethodByMerchantFailed(ctx, &pb.YearMethodTransactionMerchantRequest{
		Year:       input.Year,
		MerchantId: input.MerchantID,
	})
	if err != nil {
		return &model.APIResponseTransactionYearPaymentMethod{Status: "error", Message: err.Error()}, nil
	}
	return r.TransactionGraphql.Mapping.ToGraphqlResponseYearMethod(res), nil
}
