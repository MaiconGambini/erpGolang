package customers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/authctx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/export"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/tenantctx"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(q.Get("offset"))
	if q.Get("format") == "csv" {
		items, _, err := h.svc.List(r.Context(), ListParams{
			TenantID: tenantID,
			Search:   q.Get("search"),
			Active:   ParseBoolQuery(q.Get("active")),
			Limit:    10000,
			Offset:   0,
		})
		if err != nil {
			httpx.Error(w, "INTERNAL_ERROR", "failed to export customers", http.StatusInternalServerError)
			return
		}
		rows := make([][]string, 0, len(items))
		for _, c := range items {
			doc, email, phone := "", "", ""
			if c.Document != nil {
				doc = *c.Document
			}
			if c.Email != nil {
				email = *c.Email
			}
			if c.Phone != nil {
				phone = *c.Phone
			}
			active := "false"
			if c.Active {
				active = "true"
			}
			rows = append(rows, []string{c.Name, doc, email, phone, active})
		}
		_ = export.WriteCSV(w, "clientes.csv", []string{"nome", "documento", "email", "telefone", "ativo"}, rows)
		return
	}
	items, total, err := h.svc.List(r.Context(), ListParams{
		TenantID: tenantID,
		Search:   q.Get("search"),
		Active:   ParseBoolQuery(q.Get("active")),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to list customers", http.StatusInternalServerError)
		return
	}
	httpx.Paginated(w, items, httpx.PaginationMeta{
		Limit: limit, Offset: offset, Total: total,
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	id := chi.URLParam(r, "id")
	item, err := h.svc.Get(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			httpx.Error(w, "NOT_FOUND", "customer not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "INTERNAL_ERROR", "failed to get customer", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

type customerRequest struct {
	Name         string  `json:"name"`
	Document     *string `json:"document"`
	DocumentType *string `json:"documentType"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	PostalCode   *string `json:"postalCode"`
	Street       *string `json:"street"`
	StreetNumber *string `json:"streetNumber"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	Active       bool    `json:"active"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	var req customerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "INVALID_JSON", "invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.svc.Create(r.Context(), tenantID, user.ID, requestToInput(req))
	if err != nil {
		writeServiceError(w, err, "failed to create customer")
		return
	}
	httpx.JSON(w, http.StatusCreated, item)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	var req customerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "INVALID_JSON", "invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.svc.Update(r.Context(), tenantID, user.ID, id, requestToInput(req))
	if err != nil {
		writeServiceError(w, err, "failed to update customer")
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), tenantID, user.ID, id); err != nil {
		writeServiceError(w, err, "failed to delete customer")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func requestToInput(req customerRequest) CreateInput {
	return CreateInput{
		Name:         req.Name,
		Document:     req.Document,
		DocumentType: req.DocumentType,
		Email:        req.Email,
		Phone:        req.Phone,
		PostalCode:   req.PostalCode,
		Street:       req.Street,
		StreetNumber: req.StreetNumber,
		City:         req.City,
		State:        req.State,
		Active:       req.Active,
	}
}

func writeServiceError(w http.ResponseWriter, err error, fallback string) {
	switch {
	case errors.Is(err, errNotFound):
		httpx.Error(w, "NOT_FOUND", "customer not found", http.StatusNotFound)
	case errors.Is(err, errValidation):
		httpx.Error(w, "VALIDATION_ERROR", "invalid customer data", http.StatusBadRequest)
	case errors.Is(err, errDuplicateDocument):
		httpx.Error(w, "DUPLICATE_DOCUMENT", "document already exists for this tenant", http.StatusConflict)
	case errors.Is(err, errHasLinkedSales):
		httpx.Error(w, "CUSTOMER_HAS_SALES", "customer has linked sales", http.StatusConflict)
	default:
		httpx.Error(w, "INTERNAL_ERROR", fallback, http.StatusInternalServerError)
	}
}
