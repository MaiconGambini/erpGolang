package suppliers

import (
	"context"
	"errors"
	"strconv"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	sharedaudit "github.com/MaiconGambini/erpGolang/backend/internal/shared/audit"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *db.Queries
	audit   sharedaudit.Recorder
}

func NewService(pool *pgxpool.Pool, audit sharedaudit.Recorder) *Service {
	return &Service{queries: db.New(pool), audit: audit}
}

type SupplierDTO struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Document  *string `json:"document,omitempty"`
	Email     *string `json:"email,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Active    bool    `json:"active"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type ListParams struct {
	TenantID string
	Search   string
	Active   *bool
	Limit    int
	Offset   int
}

func (s *Service) List(ctx context.Context, p ListParams) ([]SupplierDTO, int64, error) {
	tid, err := uuid.Parse(p.TenantID)
	if err != nil {
		return nil, 0, err
	}
	pgTenant := pgutil.UUIDToPg(tid)
	rows, err := s.queries.ListSuppliers(ctx, db.ListSuppliersParams{
		TenantID:    pgTenant,
		Search:      p.Search,
		Active:      pgutil.BoolFromPtr(p.Active),
		LimitCount:  int32(p.Limit),
		OffsetCount: int32(p.Offset),
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := s.queries.CountSuppliers(ctx, db.CountSuppliersParams{
		TenantID: pgTenant,
		Search:   p.Search,
		Active:   pgutil.BoolFromPtr(p.Active),
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]SupplierDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, toDTO(row))
	}
	return out, total, nil
}

func (s *Service) Get(ctx context.Context, tenantID, id string) (SupplierDTO, error) {
	tid, _ := uuid.Parse(tenantID)
	sid, _ := uuid.Parse(id)
	row, err := s.queries.GetSupplier(ctx, db.GetSupplierParams{
		ID:       pgutil.UUIDToPg(sid),
		TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SupplierDTO{}, errNotFound
		}
		return SupplierDTO{}, err
	}
	return toDTO(row), nil
}

type CreateInput struct {
	Name     string
	Document *string
	Email    *string
	Phone    *string
	Active   bool
}

func (s *Service) Create(ctx context.Context, tenantID, actorID string, in CreateInput) (SupplierDTO, error) {
	tid, _ := uuid.Parse(tenantID)
	row, err := s.queries.CreateSupplier(ctx, db.CreateSupplierParams{
		TenantID: pgutil.UUIDToPg(tid),
		Name:     in.Name,
		Document: pgutil.TextFromPtr(in.Document),
		Email:    pgutil.TextFromPtr(in.Email),
		Phone:    pgutil.TextFromPtr(in.Phone),
		Active:   in.Active,
	})
	if err != nil {
		return SupplierDTO{}, err
	}
	dto := toDTO(row)
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "create", EntityType: "supplier", EntityID: dto.ID,
	}, nil, dto)
	return dto, nil
}

func (s *Service) Update(ctx context.Context, tenantID, actorID, id string, in CreateInput) (SupplierDTO, error) {
	tid, _ := uuid.Parse(tenantID)
	sid, _ := uuid.Parse(id)
	before, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return SupplierDTO{}, err
	}
	row, err := s.queries.UpdateSupplier(ctx, db.UpdateSupplierParams{
		ID:       pgutil.UUIDToPg(sid),
		TenantID: pgutil.UUIDToPg(tid),
		Name:     in.Name,
		Document: pgutil.TextFromPtr(in.Document),
		Email:    pgutil.TextFromPtr(in.Email),
		Phone:    pgutil.TextFromPtr(in.Phone),
		Active:   in.Active,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SupplierDTO{}, errNotFound
		}
		return SupplierDTO{}, err
	}
	dto := toDTO(row)
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "update", EntityType: "supplier", EntityID: dto.ID,
	}, before, dto)
	return dto, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, actorID, id string) error {
	tid, _ := uuid.Parse(tenantID)
	sid, _ := uuid.Parse(id)
	before, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return err
	}
	_, err = s.queries.SoftDeleteSupplier(ctx, db.SoftDeleteSupplierParams{
		ID:       pgutil.UUIDToPg(sid),
		TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errNotFound
		}
		return err
	}
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "delete", EntityType: "supplier", EntityID: id,
	}, before, nil)
	return nil
}

var errNotFound = errors.New("not found")

func toDTO(row db.Supplier) SupplierDTO {
	id, _ := pgutil.PgToUUID(row.ID)
	return SupplierDTO{
		ID:        id.String(),
		Name:      row.Name,
		Document:  pgutil.TextPtr(row.Document),
		Email:     pgutil.TextPtr(row.Email),
		Phone:     pgutil.TextPtr(row.Phone),
		Active:    row.Active,
		CreatedAt: row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: row.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ParseBoolQuery(v string) *bool {
	if v == "" {
		return nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return nil
	}
	return &b
}
