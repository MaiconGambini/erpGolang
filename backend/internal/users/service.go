package users

import (
	"context"
	"errors"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *db.Queries
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{queries: db.New(pool)}
}

type UserDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListParams struct {
	TenantID string
	Limit    int
	Offset   int
}

type UpdateInput struct {
	Name   string
	Role   string
	Active bool
}

var (
	errNotFound    = errors.New("not found")
	errValidation  = errors.New("validation error")
	errSelfDisable = errors.New("cannot deactivate self")
)

var validRoles = map[string]struct{}{
	"admin": {}, "manager": {}, "operator": {}, "viewer": {},
}

func (s *Service) List(ctx context.Context, p ListParams) ([]UserDTO, int64, error) {
	tid, err := uuid.Parse(p.TenantID)
	if err != nil {
		return nil, 0, err
	}
	pgTenant := pgutil.UUIDToPg(tid)
	rows, err := s.queries.ListUsers(ctx, db.ListUsersParams{
		TenantID:    pgTenant,
		LimitCount:  int32(p.Limit),
		OffsetCount: int32(p.Offset),
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := s.queries.CountUsers(ctx, pgTenant)
	if err != nil {
		return nil, 0, err
	}
	out := make([]UserDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, rowToDTO(row.ID, row.Email, row.Name, row.Role, row.Active, row.CreatedAt, row.UpdatedAt))
	}
	return out, total, nil
}

func (s *Service) Get(ctx context.Context, tenantID, id string) (UserDTO, error) {
	tid, _ := uuid.Parse(tenantID)
	uid, _ := uuid.Parse(id)
	row, err := s.queries.GetUser(ctx, db.GetUserParams{
		ID:       pgutil.UUIDToPg(uid),
		TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserDTO{}, errNotFound
		}
		return UserDTO{}, err
	}
	return rowToDTO(row.ID, row.Email, row.Name, row.Role, row.Active, row.CreatedAt, row.UpdatedAt), nil
}

func (s *Service) Update(ctx context.Context, tenantID, actorID, id string, in UpdateInput) (UserDTO, error) {
	if in.Name == "" {
		return UserDTO{}, errValidation
	}
	if _, ok := validRoles[in.Role]; !ok {
		return UserDTO{}, errValidation
	}
	if actorID == id && !in.Active {
		return UserDTO{}, errSelfDisable
	}
	tid, _ := uuid.Parse(tenantID)
	uid, _ := uuid.Parse(id)
	row, err := s.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       pgutil.UUIDToPg(uid),
		TenantID: pgutil.UUIDToPg(tid),
		Name:     in.Name,
		Role:     in.Role,
		Active:   in.Active,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserDTO{}, errNotFound
		}
		return UserDTO{}, err
	}
	return rowToDTO(row.ID, row.Email, row.Name, row.Role, row.Active, row.CreatedAt, row.UpdatedAt), nil
}

func rowToDTO(id pgtype.UUID, email, name, role string, active bool, createdAt, updatedAt pgtype.Timestamptz) UserDTO {
	uid, _ := pgutil.PgToUUID(id)
	return UserDTO{
		ID:        uid.String(),
		Email:     email,
		Name:      name,
		Role:      role,
		Active:    active,
		CreatedAt: createdAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: updatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}
