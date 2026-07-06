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

func (s *Service) Summary(ctx context.Context, params SummaryParams) (db.GetDashboardSummaryRow, error) {
	tid, err := uuid.Parse(params.TenantID)
	if err != nil {
		return db.GetDashboardSummaryRow{}, err
	}
	return s.queries.GetDashboardSummary(ctx, db.GetDashboardSummaryParams{
		TenantID:  pgutil.UUIDToPg(tid),
		Threshold: params.Threshold,
	})
}
