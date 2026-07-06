package customers

import (
	"context"
	"errors"
	"strconv"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	sharedaudit "github.com/MaiconGambini/erpGolang/backend/internal/shared/audit"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/party"
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

type CustomerDTO struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Document     *string `json:"document,omitempty"`
	DocumentType *string `json:"documentType,omitempty"`
	Email        *string `json:"email,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	PostalCode   *string `json:"postalCode,omitempty"`
	Street       *string `json:"street,omitempty"`
	StreetNumber *string `json:"streetNumber,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	Active       bool    `json:"active"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

type ListParams struct {
	TenantID string
	Search   string
	Active   *bool
	Limit    int
	Offset   int
}

type CreateInput struct {
	Name         string
	Document     *string
	DocumentType *string
	Email        *string
	Phone        *string
	PostalCode   *string
	Street       *string
	StreetNumber *string
	City         *string
	State        *string
	Active       bool
}

func (s *Service) List(ctx context.Context, p ListParams) ([]CustomerDTO, int64, error) {
	tid, err := uuid.Parse(p.TenantID)
	if err != nil {
		return nil, 0, err
	}
	pgTenant := pgutil.UUIDToPg(tid)
	rows, err := s.queries.ListCustomers(ctx, db.ListCustomersParams{
		TenantID:    pgTenant,
		Search:      p.Search,
		Active:      pgutil.BoolFromPtr(p.Active),
		LimitCount:  int32(p.Limit),
		OffsetCount: int32(p.Offset),
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := s.queries.CountCustomers(ctx, db.CountCustomersParams{
		TenantID: pgTenant,
		Search:   p.Search,
		Active:   pgutil.BoolFromPtr(p.Active),
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]CustomerDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, toDTO(row))
	}
	return out, total, nil
}

func (s *Service) Get(ctx context.Context, tenantID, id string) (CustomerDTO, error) {
	tid, _ := uuid.Parse(tenantID)
	cid, _ := uuid.Parse(id)
	row, err := s.queries.GetCustomer(ctx, db.GetCustomerParams{
		ID:       pgutil.UUIDToPg(cid),
		TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CustomerDTO{}, errNotFound
		}
		return CustomerDTO{}, err
	}
	return toDTO(row), nil
}

func (s *Service) Create(ctx context.Context, tenantID, actorID string, in CreateInput) (CustomerDTO, error) {
	params, err := buildParams(in)
	if err != nil {
		return CustomerDTO{}, err
	}
	tid, _ := uuid.Parse(tenantID)
	params.TenantID = pgutil.UUIDToPg(tid)
	row, err := s.queries.CreateCustomer(ctx, params)
	if err != nil {
		if isUniqueViolation(err) {
			return CustomerDTO{}, errDuplicateDocument
		}
		return CustomerDTO{}, err
	}
	dto := toDTO(row)
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "create", EntityType: "customer", EntityID: dto.ID,
	}, nil, dto)
	return dto, nil
}

func (s *Service) Update(ctx context.Context, tenantID, actorID, id string, in CreateInput) (CustomerDTO, error) {
	params, err := buildParams(in)
	if err != nil {
		return CustomerDTO{}, err
	}
	tid, _ := uuid.Parse(tenantID)
	cid, _ := uuid.Parse(id)
	before, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return CustomerDTO{}, err
	}
	row, err := s.queries.UpdateCustomer(ctx, db.UpdateCustomerParams{
		ID:           pgutil.UUIDToPg(cid),
		TenantID:     pgutil.UUIDToPg(tid),
		Name:         params.Name,
		Document:     params.Document,
		DocumentType: params.DocumentType,
		Email:        params.Email,
		Phone:        params.Phone,
		PostalCode:   params.PostalCode,
		Street:       params.Street,
		StreetNumber: params.StreetNumber,
		City:         params.City,
		State:        params.State,
		Active:       params.Active,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CustomerDTO{}, errNotFound
		}
		if isUniqueViolation(err) {
			return CustomerDTO{}, errDuplicateDocument
		}
		return CustomerDTO{}, err
	}
	dto := toDTO(row)
	_ = s.audit.RecordWithPayload(ctx, sharedaudit.Event{
		TenantID: tenantID, ActorID: actorID,
		Action: "update", EntityType: "customer", EntityID: dto.ID,
	}, before, dto)
	return dto, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, actorID, id string) error {
	tid, _ := uuid.Parse(tenantID)
	cid, _ := uuid.Parse(id)
	before, err := s.Get(ctx, tenantID, id)
	if err != nil {
		return err
	}
	count, err := s.queries.CountSalesByCustomer(ctx, db.CountSalesByCustomerParams{
		TenantID:   pgutil.UUIDToPg(tid),
		CustomerID: pgutil.UUIDToPg(cid),
	})
	if err != nil {
		return err
	}
	if count > 0 {
		return errHasLinkedSales
	}
	_, err = s.queries.SoftDeleteCustomer(ctx, db.SoftDeleteCustomerParams{
		ID:       pgutil.UUIDToPg(cid),
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
		Action: "delete", EntityType: "customer", EntityID: id,
	}, before, nil)
	return nil
}

var (
	errNotFound          = errors.New("not found")
	errValidation        = errors.New("validation error")
	errDuplicateDocument = errors.New("duplicate document")
	errHasLinkedSales    = errors.New("customer has linked sales")
)

func buildParams(in CreateInput) (db.CreateCustomerParams, error) {
	if err := party.ValidatePartyInput(party.PartyInput{
		Name:         in.Name,
		DocumentType: in.DocumentType,
		Document:     in.Document,
		Email:        in.Email,
		State:        in.State,
	}); err != nil {
		return db.CreateCustomerParams{}, errValidation
	}
	return db.CreateCustomerParams{
		Name:         in.Name,
		Document:     pgutil.TextFromPtr(party.NormalizeDocumentPtr(in.Document)),
		DocumentType: pgutil.TextFromPtr(party.NormalizeDocumentTypePtr(in.DocumentType)),
		Email:        pgutil.TextFromPtr(in.Email),
		Phone:        pgutil.TextFromPtr(in.Phone),
		PostalCode:   pgutil.TextFromPtr(in.PostalCode),
		Street:       pgutil.TextFromPtr(in.Street),
		StreetNumber: pgutil.TextFromPtr(in.StreetNumber),
		City:         pgutil.TextFromPtr(in.City),
		State:        pgutil.TextFromPtr(party.NormalizeStatePtr(in.State)),
		Active:       in.Active,
	}, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func toDTO(row db.Customer) CustomerDTO {
	id, _ := pgutil.PgToUUID(row.ID)
	return CustomerDTO{
		ID:           id.String(),
		Name:         row.Name,
		Document:     pgutil.TextPtr(row.Document),
		DocumentType: pgutil.TextPtr(row.DocumentType),
		Email:        pgutil.TextPtr(row.Email),
		Phone:        pgutil.TextPtr(row.Phone),
		PostalCode:   pgutil.TextPtr(row.PostalCode),
		Street:       pgutil.TextPtr(row.Street),
		StreetNumber: pgutil.TextPtr(row.StreetNumber),
		City:         pgutil.TextPtr(row.City),
		State:        pgutil.TextPtr(row.State),
		Active:       row.Active,
		CreatedAt:    row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    row.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
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
