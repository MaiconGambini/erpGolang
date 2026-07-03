package sales

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	sharedaudit "github.com/MaiconGambini/erpGolang/backend/internal/shared/audit"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool  *pgxpool.Pool
	audit sharedaudit.Recorder
}

func NewService(pool *pgxpool.Pool, audit sharedaudit.Recorder) *Service {
	return &Service{pool: pool, audit: audit}
}

type SaleItemDTO struct {
	ID          string `json:"id"`
	ProductID   string `json:"productId"`
	ProductName string `json:"productName"`
	ProductSku  string `json:"productSku"`
	Quantity    int    `json:"quantity"`
	UnitPrice   string `json:"unitPrice"`
	LineTotal   string `json:"lineTotal"`
}

type SaleDTO struct {
	ID           string        `json:"id"`
	CustomerID   string        `json:"customerId"`
	CustomerName string        `json:"customerName"`
	Status       string        `json:"status"`
	Total        string        `json:"total"`
	Notes        *string       `json:"notes,omitempty"`
	Items        []SaleItemDTO `json:"items,omitempty"`
	CreatedAt    string        `json:"createdAt"`
	UpdatedAt    string        `json:"updatedAt"`
}

type ListParams struct {
	TenantID string
	Search   string
	Status   *string
	Limit    int
	Offset   int
}

type ItemInput struct {
	ProductID string
	Quantity  int
}

type CreateInput struct {
	CustomerID string
	Notes      *string
	Items      []ItemInput
}

var (
	errNotFound          = errors.New("not found")
	errValidation        = errors.New("validation error")
	errInvalidStatus     = errors.New("invalid status")
	errInsufficientStock = errors.New("insufficient stock")
)

func (s *Service) List(ctx context.Context, p ListParams) ([]SaleDTO, int64, error) {
	tid, err := uuid.Parse(p.TenantID)
	if err != nil {
		return nil, 0, err
	}
	pgTenant := pgutil.UUIDToPg(tid)
	q := db.New(s.pool)
	rows, err := q.ListSales(ctx, db.ListSalesParams{
		TenantID:    pgTenant,
		Search:      p.Search,
		Status:      textFromPtr(p.Status),
		LimitCount:  int32(p.Limit),
		OffsetCount: int32(p.Offset),
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := q.CountSales(ctx, db.CountSalesParams{
		TenantID: pgTenant,
		Search:   p.Search,
		Status:   textFromPtr(p.Status),
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]SaleDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, listRowToDTO(row))
	}
	return out, total, nil
}

func (s *Service) Get(ctx context.Context, tenantID, id string) (SaleDTO, error) {
	tid, _ := uuid.Parse(tenantID)
	sid, _ := uuid.Parse(id)
	q := db.New(s.pool)
	row, err := q.GetSale(ctx, db.GetSaleParams{
		ID:       pgutil.UUIDToPg(sid),
		TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SaleDTO{}, errNotFound
		}
		return SaleDTO{}, err
	}
	items, err := q.ListSaleItems(ctx, db.ListSaleItemsParams{
		SaleID:   pgutil.UUIDToPg(sid),
		TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		return SaleDTO{}, err
	}
	return saleRowToDTO(row, items), nil
}

func (s *Service) Create(ctx context.Context, tenantID, actorID string, in CreateInput) (SaleDTO, error) {
	if err := validateCreateInput(in); err != nil {
		return SaleDTO{}, err
	}
	tid, _ := uuid.Parse(tenantID)
	pgTenant := pgutil.UUIDToPg(tid)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return SaleDTO{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	qtx := db.New(tx)

	cid, _ := uuid.Parse(in.CustomerID)
	if _, err := qtx.GetCustomer(ctx, db.GetCustomerParams{
		ID:       pgutil.UUIDToPg(cid),
		TenantID: pgTenant,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SaleDTO{}, fmt.Errorf("%w: customer not found", errValidation)
		}
		return SaleDTO{}, err
	}

	resolved, totalNum, err := s.resolveItems(ctx, qtx, pgTenant, in.Items)
	if err != nil {
		return SaleDTO{}, err
	}
	total, err := pgutil.NumericFromString(fmt.Sprintf("%.2f", totalNum))
	if err != nil {
		return SaleDTO{}, err
	}

	sale, err := qtx.CreateSale(ctx, db.CreateSaleParams{
		TenantID:   pgTenant,
		CustomerID: pgutil.UUIDToPg(cid),
		Status:     "draft",
		Total:      total,
		Notes:      pgutil.TextFromPtr(in.Notes),
	})
	if err != nil {
		return SaleDTO{}, err
	}
	saleID, _ := pgutil.PgToUUID(sale.ID)

	for _, item := range resolved {
		if _, err := qtx.CreateSaleItem(ctx, db.CreateSaleItemParams{
			TenantID:  pgTenant,
			SaleID:    sale.ID,
			ProductID: item.productID,
			Quantity:  item.quantity,
			UnitPrice: item.unitPrice,
			LineTotal: item.lineTotal,
		}); err != nil {
			return SaleDTO{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return SaleDTO{}, err
	}

	dto, err := s.Get(ctx, tenantID, saleID.String())
	if err != nil {
		return SaleDTO{}, err
	}
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "create", EntityType: "sale", EntityID: dto.ID,
	}, nil, dto)
	return dto, nil
}

func (s *Service) Update(ctx context.Context, tenantID, actorID, id string, in CreateInput) (SaleDTO, error) {
	if err := validateCreateInput(in); err != nil {
		return SaleDTO{}, err
	}
	tid, _ := uuid.Parse(tenantID)
	sid, _ := uuid.Parse(id)
	pgTenant := pgutil.UUIDToPg(tid)
	pgSale := pgutil.UUIDToPg(sid)

	before, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return SaleDTO{}, err
	}
	if before.Status != "draft" {
		return SaleDTO{}, errInvalidStatus
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return SaleDTO{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	qtx := db.New(tx)

	cid, _ := uuid.Parse(in.CustomerID)
	resolved, totalNum, err := s.resolveItems(ctx, qtx, pgTenant, in.Items)
	if err != nil {
		return SaleDTO{}, err
	}
	total, err := pgutil.NumericFromString(fmt.Sprintf("%.2f", totalNum))
	if err != nil {
		return SaleDTO{}, err
	}

	if _, err := qtx.UpdateSaleDraft(ctx, db.UpdateSaleDraftParams{
		ID:         pgSale,
		TenantID:   pgTenant,
		CustomerID: pgutil.UUIDToPg(cid),
		Total:      total,
		Notes:      pgutil.TextFromPtr(in.Notes),
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SaleDTO{}, errInvalidStatus
		}
		return SaleDTO{}, err
	}

	if err := qtx.DeleteSaleItems(ctx, db.DeleteSaleItemsParams{
		SaleID: pgSale, TenantID: pgTenant,
	}); err != nil {
		return SaleDTO{}, err
	}
	for _, item := range resolved {
		if _, err := qtx.CreateSaleItem(ctx, db.CreateSaleItemParams{
			TenantID:  pgTenant,
			SaleID:    pgSale,
			ProductID: item.productID,
			Quantity:  item.quantity,
			UnitPrice: item.unitPrice,
			LineTotal: item.lineTotal,
		}); err != nil {
			return SaleDTO{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return SaleDTO{}, err
	}

	dto, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return SaleDTO{}, err
	}
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "update", EntityType: "sale", EntityID: dto.ID,
	}, before, dto)
	return dto, nil
}

func (s *Service) Confirm(ctx context.Context, tenantID, actorID, id string) (SaleDTO, error) {
	before, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return SaleDTO{}, err
	}
	if before.Status != "draft" {
		return SaleDTO{}, errInvalidStatus
	}

	tid, _ := uuid.Parse(tenantID)
	sid, _ := uuid.Parse(id)
	pgTenant := pgutil.UUIDToPg(tid)
	pgSale := pgutil.UUIDToPg(sid)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return SaleDTO{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	qtx := db.New(tx)

	items, err := qtx.ListSaleItems(ctx, db.ListSaleItemsParams{
		SaleID: pgSale, TenantID: pgTenant,
	})
	if err != nil {
		return SaleDTO{}, err
	}
	for _, item := range items {
		if _, err := qtx.DecrementProductStock(ctx, db.DecrementProductStockParams{
			ID:       item.ProductID,
			TenantID: pgTenant,
			Quantity: item.Quantity,
		}); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return SaleDTO{}, errInsufficientStock
			}
			return SaleDTO{}, err
		}
	}

	if _, err := qtx.UpdateSaleStatus(ctx, db.UpdateSaleStatusParams{
		ID: pgSale, TenantID: pgTenant, Status: "confirmed",
	}); err != nil {
		return SaleDTO{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return SaleDTO{}, err
	}

	dto, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return SaleDTO{}, err
	}
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "confirm", EntityType: "sale", EntityID: dto.ID,
	}, before, dto)
	return dto, nil
}

func (s *Service) Cancel(ctx context.Context, tenantID, actorID, id string) (SaleDTO, error) {
	before, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return SaleDTO{}, err
	}
	if before.Status != "confirmed" {
		return SaleDTO{}, errInvalidStatus
	}

	tid, _ := uuid.Parse(tenantID)
	sid, _ := uuid.Parse(id)
	pgTenant := pgutil.UUIDToPg(tid)
	pgSale := pgutil.UUIDToPg(sid)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return SaleDTO{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	qtx := db.New(tx)

	items, err := qtx.ListSaleItems(ctx, db.ListSaleItemsParams{
		SaleID: pgSale, TenantID: pgTenant,
	})
	if err != nil {
		return SaleDTO{}, err
	}
	for _, item := range items {
		if err := qtx.IncrementProductStock(ctx, db.IncrementProductStockParams{
			ID: item.ProductID, TenantID: pgTenant, Quantity: item.Quantity,
		}); err != nil {
			return SaleDTO{}, err
		}
	}

	if _, err := qtx.UpdateSaleStatus(ctx, db.UpdateSaleStatusParams{
		ID: pgSale, TenantID: pgTenant, Status: "cancelled",
	}); err != nil {
		return SaleDTO{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return SaleDTO{}, err
	}

	dto, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return SaleDTO{}, err
	}
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "cancel", EntityType: "sale", EntityID: dto.ID,
	}, before, dto)
	return dto, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, actorID, id string) error {
	before, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if before.Status != "draft" {
		return errInvalidStatus
	}

	tid, _ := uuid.Parse(tenantID)
	sid, _ := uuid.Parse(id)
	_, err = db.New(s.pool).SoftDeleteSale(ctx, db.SoftDeleteSaleParams{
		ID: pgutil.UUIDToPg(sid), TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errInvalidStatus
		}
		return err
	}
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "delete", EntityType: "sale", EntityID: id,
	}, before, nil)
	return nil
}

type resolvedItem struct {
	productID pgtype.UUID
	quantity  int32
	unitPrice pgtype.Numeric
	lineTotal pgtype.Numeric
}

func (s *Service) resolveItems(ctx context.Context, q *db.Queries, pgTenant pgtype.UUID, items []ItemInput) ([]resolvedItem, float64, error) {
	if len(items) == 0 {
		return nil, 0, errValidation
	}
	out := make([]resolvedItem, 0, len(items))
	var total float64
	for _, in := range items {
		if in.Quantity <= 0 {
			return nil, 0, errValidation
		}
		pid, err := uuid.Parse(in.ProductID)
		if err != nil {
			return nil, 0, errValidation
		}
		product, err := q.GetProduct(ctx, db.GetProductParams{
			ID: pgutil.UUIDToPg(pid), TenantID: pgTenant,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, 0, fmt.Errorf("%w: product not found", errValidation)
			}
			return nil, 0, err
		}
		if !product.Active {
			return nil, 0, fmt.Errorf("%w: product inactive", errValidation)
		}
		unitStr := pgutil.NumericToString(product.Price)
		unitF, _ := strconv.ParseFloat(unitStr, 64)
		lineF := math.Round(unitF*float64(in.Quantity)*100) / 100
		lineStr := fmt.Sprintf("%.2f", lineF)
		lineNum, err := pgutil.NumericFromString(lineStr)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, resolvedItem{
			productID: pgutil.UUIDToPg(pid),
			quantity:  int32(in.Quantity),
			unitPrice: product.Price,
			lineTotal: lineNum,
		})
		total += lineF
	}
	return out, math.Round(total*100) / 100, nil
}

func validateCreateInput(in CreateInput) error {
	if in.CustomerID == "" || len(in.Items) == 0 {
		return errValidation
	}
	for _, item := range in.Items {
		if item.ProductID == "" || item.Quantity <= 0 {
			return errValidation
		}
	}
	return nil
}

func listRowToDTO(row db.ListSalesRow) SaleDTO {
	id, _ := pgutil.PgToUUID(row.ID)
	cid, _ := pgutil.PgToUUID(row.CustomerID)
	return SaleDTO{
		ID:           id.String(),
		CustomerID:   cid.String(),
		CustomerName: row.CustomerName,
		Status:       row.Status,
		Total:        pgutil.NumericToString(row.Total),
		Notes:        pgutil.TextPtr(row.Notes),
		CreatedAt:    row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    row.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func saleRowToDTO(row db.GetSaleRow, items []db.ListSaleItemsRow) SaleDTO {
	id, _ := pgutil.PgToUUID(row.ID)
	cid, _ := pgutil.PgToUUID(row.CustomerID)
	dto := SaleDTO{
		ID:           id.String(),
		CustomerID:   cid.String(),
		CustomerName: row.CustomerName,
		Status:       row.Status,
		Total:        pgutil.NumericToString(row.Total),
		Notes:        pgutil.TextPtr(row.Notes),
		CreatedAt:    row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    row.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		Items:        make([]SaleItemDTO, 0, len(items)),
	}
	for _, item := range items {
		iid, _ := pgutil.PgToUUID(item.ID)
		pid, _ := pgutil.PgToUUID(item.ProductID)
		dto.Items = append(dto.Items, SaleItemDTO{
			ID:          iid.String(),
			ProductID:   pid.String(),
			ProductName: item.ProductName,
			ProductSku:  item.ProductSku,
			Quantity:    int(item.Quantity),
			UnitPrice:   pgutil.NumericToString(item.UnitPrice),
			LineTotal:   pgutil.NumericToString(item.LineTotal),
		})
	}
	return dto
}

func textFromPtr(s *string) pgtype.Text {
	if s == nil || *s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func ParseStatusQuery(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
