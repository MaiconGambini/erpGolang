package audit

import (
	"context"
	"encoding/json"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	sharedaudit "github.com/MaiconGambini/erpGolang/backend/internal/shared/audit"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *db.Queries
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{queries: db.New(pool)}
}

func (s *Service) Record(ctx context.Context, event sharedaudit.Event) error {
	return s.RecordWithPayload(ctx, event, nil, nil)
}

func (s *Service) RecordWithPayload(ctx context.Context, event sharedaudit.Event, before, after any) error {
	tenantID, err := uuid.Parse(event.TenantID)
	if err != nil {
		return err
	}
	entityID, err := uuid.Parse(event.EntityID)
	if err != nil {
		return err
	}
	var actor pgtype.UUID
	if event.ActorID != "" {
		aid, err := uuid.Parse(event.ActorID)
		if err != nil {
			return err
		}
		actor = pgutil.UUIDToPg(aid)
	}
	var b, a []byte
	if before != nil {
		b, _ = json.Marshal(before)
	}
	if after != nil {
		a, _ = json.Marshal(after)
	}
	_, err = s.queries.InsertAuditLog(ctx, db.InsertAuditLogParams{
		TenantID:    pgutil.UUIDToPg(tenantID),
		ActorUserID: actor,
		Action:      event.Action,
		EntityType:  event.EntityType,
		EntityID:    pgutil.UUIDToPg(entityID),
		Before:      b,
		After:       a,
	})
	return err
}
