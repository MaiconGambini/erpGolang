package dashboard

import (
	"context"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *db.Queries
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{queries: db.New(pool)}
}

type SummaryParams struct {
	TenantID  string
	Threshold int32
}

type SummaryDTO struct {
	ActiveCustomers       int64  `json:"activeCustomers"`
	NewCustomers30d       int64  `json:"newCustomers30d"`
	DraftSales            int64  `json:"draftSales"`
	LowStockAlerts        int64  `json:"lowStockAlerts"`
	ConfirmedSalesCount   int64  `json:"confirmedSalesCount"`
	ConfirmedSalesTotal   string `json:"confirmedSalesTotal"`
}

func (s *Service) Summary(ctx context.Context, params SummaryParams) (SummaryDTO, error) {
	tid, err := uuid.Parse(params.TenantID)
	if err != nil {
		return SummaryDTO{}, err
	}
	row, err := s.queries.GetDashboardSummary(ctx, db.GetDashboardSummaryParams{
		TenantID:  pgutil.UUIDToPg(tid),
		Threshold: params.Threshold,
	})
	if err != nil {
		return SummaryDTO{}, err
	}
	return SummaryDTO{
		ActiveCustomers:     row.ActiveCustomers,
		NewCustomers30d:     row.NewCustomers30d,
		DraftSales:          row.DraftSales,
		LowStockAlerts:      row.LowStockAlerts,
		ConfirmedSalesCount: row.ConfirmedSalesCount,
		ConfirmedSalesTotal: pgutil.NumericToString(row.ConfirmedSalesTotal),
	}, nil
}
