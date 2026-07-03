package products

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	sharedaudit "github.com/MaiconGambini/erpGolang/backend/internal/shared/audit"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *db.Queries
	audit   sharedaudit.Recorder
}

func NewService(pool *pgxpool.Pool, audit sharedaudit.Recorder) *Service {
	return &Service{queries: db.New(pool), audit: audit}
}

type ProductDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Sku       string `json:"sku"`
	Price     string `json:"price"`
	Stock     int    `json:"stock"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListParams struct {
	TenantID string
	Search   string
	Active   *bool
	Limit    int
	Offset   int
}

type LowStockParams struct {
	TenantID  string
	Threshold int
	Limit     int
}

func (s *Service) List(ctx context.Context, p ListParams) ([]ProductDTO, int64, error) {
	tid, err := uuid.Parse(p.TenantID)
	if err != nil {
		return nil, 0, err
	}
	pgTenant := pgutil.UUIDToPg(tid)
	rows, err := s.queries.ListProducts(ctx, db.ListProductsParams{
		TenantID:    pgTenant,
		Search:      p.Search,
		Active:      pgutil.BoolFromPtr(p.Active),
		LimitCount:  int32(p.Limit),
		OffsetCount: int32(p.Offset),
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := s.queries.CountProducts(ctx, db.CountProductsParams{
		TenantID: pgTenant,
		Search:   p.Search,
		Active:   pgutil.BoolFromPtr(p.Active),
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]ProductDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, toDTO(row))
	}
	return out, total, nil
}

func (s *Service) LowStock(ctx context.Context, p LowStockParams) ([]ProductDTO, error) {
	tid, err := uuid.Parse(p.TenantID)
	if err != nil {
		return nil, err
	}
	threshold := int32(p.Threshold)
	if threshold < 0 {
		threshold = 0
	}
	limit := int32(p.Limit)
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.queries.ListLowStockProducts(ctx, db.ListLowStockProductsParams{
		TenantID:   pgutil.UUIDToPg(tid),
		Threshold:  threshold,
		LimitCount: limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]ProductDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, toDTO(row))
	}
	return out, nil
}

func (s *Service) Get(ctx context.Context, tenantID, id string) (ProductDTO, error) {
	tid, _ := uuid.Parse(tenantID)
	pid, _ := uuid.Parse(id)
	row, err := s.queries.GetProduct(ctx, db.GetProductParams{
		ID:       pgutil.UUIDToPg(pid),
		TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ProductDTO{}, errNotFound
		}
		return ProductDTO{}, err
	}
	return toDTO(row), nil
}

type CreateInput struct {
	Name   string
	Sku    string
	Price  string
	Stock  int
	Active bool
}

func (s *Service) Create(ctx context.Context, tenantID, actorID string, in CreateInput) (ProductDTO, error) {
	if err := validateInput(in); err != nil {
		return ProductDTO{}, err
	}
	tid, _ := uuid.Parse(tenantID)
	price, err := pgutil.NumericFromString(in.Price)
	if err != nil {
		return ProductDTO{}, fmt.Errorf("%w: invalid price", errValidation)
	}
	row, err := s.queries.CreateProduct(ctx, db.CreateProductParams{
		TenantID: pgutil.UUIDToPg(tid),
		Name:     in.Name,
		Sku:      in.Sku,
		Price:    price,
		Stock:    int32(in.Stock),
		Active:   in.Active,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return ProductDTO{}, errDuplicateSKU
		}
		return ProductDTO{}, err
	}
	dto := toDTO(row)
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "create", EntityType: "product", EntityID: dto.ID,
	}, nil, dto)
	return dto, nil
}

func (s *Service) Update(ctx context.Context, tenantID, actorID, id string, in CreateInput) (ProductDTO, error) {
	if err := validateInput(in); err != nil {
		return ProductDTO{}, err
	}
	tid, _ := uuid.Parse(tenantID)
	pid, _ := uuid.Parse(id)
	before, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return ProductDTO{}, err
	}
	price, err := pgutil.NumericFromString(in.Price)
	if err != nil {
		return ProductDTO{}, fmt.Errorf("%w: invalid price", errValidation)
	}
	row, err := s.queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:       pgutil.UUIDToPg(pid),
		TenantID: pgutil.UUIDToPg(tid),
		Name:     in.Name,
		Sku:      in.Sku,
		Price:    price,
		Stock:    int32(in.Stock),
		Active:   in.Active,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ProductDTO{}, errNotFound
		}
		if isUniqueViolation(err) {
			return ProductDTO{}, errDuplicateSKU
		}
		return ProductDTO{}, err
	}
	dto := toDTO(row)
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "update", EntityType: "product", EntityID: dto.ID,
	}, before, dto)
	return dto, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, actorID, id string) error {
	tid, _ := uuid.Parse(tenantID)
	pid, _ := uuid.Parse(id)
	before, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return err
	}
	_, err = s.queries.SoftDeleteProduct(ctx, db.SoftDeleteProductParams{
		ID:       pgutil.UUIDToPg(pid),
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
		Action: "delete", EntityType: "product", EntityID: id,
	}, before, nil)
	return nil
}

var (
	errNotFound     = errors.New("not found")
	errValidation   = errors.New("validation error")
	errDuplicateSKU = errors.New("duplicate sku")
)

func validateInput(in CreateInput) error {
	if in.Name == "" || in.Sku == "" {
		return errValidation
	}
	price, err := strconv.ParseFloat(in.Price, 64)
	if err != nil || price <= 0 {
		return errValidation
	}
	if in.Stock < 0 {
		return errValidation
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func toDTO(row db.Product) ProductDTO {
	id, _ := pgutil.PgToUUID(row.ID)
	return ProductDTO{
		ID:        id.String(),
		Name:      row.Name,
		Sku:       row.Sku,
		Price:     pgutil.NumericToString(row.Price),
		Stock:     int(row.Stock),
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
