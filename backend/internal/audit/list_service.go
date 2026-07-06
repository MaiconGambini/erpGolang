package audit

import (
	"context"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListService struct {
	queries *db.Queries
}

func NewListService(pool *pgxpool.Pool) *ListService {
	return &ListService{queries: db.New(pool)}
}

type LogDTO struct {
	ID          string  `json:"id"`
	ActorUserID *string `json:"actorUserId,omitempty"`
	Action      string  `json:"action"`
	EntityType  string  `json:"entityType"`
	EntityID    string  `json:"entityId"`
	CreatedAt   string  `json:"createdAt"`
}

type ListParams struct {
	TenantID string
	Limit    int
	Offset   int
}

func (s *ListService) List(ctx context.Context, p ListParams) ([]LogDTO, int64, error) {
	tid, err := uuid.Parse(p.TenantID)
	if err != nil {
		return nil, 0, err
	}
	pgTenant := pgutil.UUIDToPg(tid)
	rows, err := s.queries.ListAuditLogs(ctx, db.ListAuditLogsParams{
		TenantID:    pgTenant,
		LimitCount:  int32(p.Limit),
		OffsetCount: int32(p.Offset),
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := s.queries.CountAuditLogs(ctx, pgTenant)
	if err != nil {
		return nil, 0, err
	}
	out := make([]LogDTO, 0, len(rows))
	for _, row := range rows {
		id, _ := pgutil.PgToUUID(row.ID)
		entityID, _ := pgutil.PgToUUID(row.EntityID)
		dto := LogDTO{
			ID:         id.String(),
			Action:     row.Action,
			EntityType: row.EntityType,
			EntityID:   entityID.String(),
			CreatedAt:  row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}
		if row.ActorUserID.Valid {
			actorID, _ := pgutil.PgToUUID(row.ActorUserID)
			s := actorID.String()
			dto.ActorUserID = &s
		}
		out = append(out, dto)
	}
	return out, total, nil
}
