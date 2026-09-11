package reports

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/pgutil"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/querytime"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jung-kurt/gofpdf/v2"
)

type Service struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool, queries: db.New(pool)}
}

type SalesByDayDTO struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
	Total string `json:"total"`
}

type TopProductDTO struct {
	ProductName string `json:"productName"`
	ProductSku  string `json:"productSku"`
	Quantity    int64  `json:"quantity"`
	Revenue     string `json:"revenue"`
}

type RangeParams struct {
	TenantID string
	From     time.Time
	To       time.Time
	Limit    int32
}

func (s *Service) SalesByDay(ctx context.Context, p RangeParams) ([]SalesByDayDTO, error) {
	tid, _ := uuid.Parse(p.TenantID)
	rows, err := s.queries.SalesByDay(ctx, db.SalesByDayParams{
		TenantID: pgutil.UUIDToPg(tid),
		FromDate: pgutil.Timestamptz(p.From),
		ToDate:   pgutil.Timestamptz(p.To),
	})
	if err != nil {
		return nil, err
	}
	out := make([]SalesByDayDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, SalesByDayDTO{
			Date:  row.SaleDate.Time.Format("2006-01-02"),
			Count: row.SaleCount,
			Total: pgutil.NumericToString(row.SaleTotal),
		})
	}
	return out, nil
}

func (s *Service) TopProducts(ctx context.Context, p RangeParams) ([]TopProductDTO, error) {
	tid, _ := uuid.Parse(p.TenantID)
	limit := p.Limit
	if limit <= 0 {
		limit = 5
	}
	rows, err := s.queries.TopProducts(ctx, db.TopProductsParams{
		TenantID:    pgutil.UUIDToPg(tid),
		FromDate:    pgutil.Timestamptz(p.From),
		ToDate:      pgutil.Timestamptz(p.To),
		LimitCount:  limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]TopProductDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, TopProductDTO{
			ProductName: row.ProductName,
			ProductSku:  row.ProductSku,
			Quantity:    row.QuantitySold,
			Revenue:     pgutil.NumericToString(row.Revenue),
		})
	}
	return out, nil
}

func (s *Service) SalePDF(ctx context.Context, tenantID, saleID string) ([]byte, error) {
	tid, _ := uuid.Parse(tenantID)
	sid, _ := uuid.Parse(saleID)
	q := db.New(s.pool)
	sale, err := q.GetSale(ctx, db.GetSaleParams{
		ID: pgutil.UUIDToPg(sid), TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		return nil, err
	}
	items, err := q.ListSaleItems(ctx, db.ListSaleItemsParams{
		SaleID: pgutil.UUIDToPg(sid), TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		return nil, err
	}
	id, _ := pgutil.PgToUUID(sale.ID)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "goERP - Pedido de Venda")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(40, 8, fmt.Sprintf("ID: %s", id.String()))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("Cliente: %s", sale.CustomerName))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("Status: %s", sale.Status))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("Data: %s", sale.CreatedAt.Time.Format("02/01/2006 15:04")))
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(80, 8, "Produto")
	pdf.Cell(30, 8, "Qtd")
	pdf.Cell(40, 8, "Preco")
	pdf.Cell(40, 8, "Total")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	for _, it := range items {
		pdf.Cell(80, 7, it.ProductName)
		pdf.Cell(30, 7, fmt.Sprintf("%d", it.Quantity))
		pdf.Cell(40, 7, pgutil.NumericToString(it.UnitPrice))
		pdf.Cell(40, 7, pgutil.NumericToString(it.LineTotal))
		pdf.Ln(7)
	}
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, fmt.Sprintf("Total: R$ %s", pgutil.NumericToString(sale.Total)))
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Service) SalesSummaryPDF(ctx context.Context, tenantID string, from, to time.Time) ([]byte, error) {
	tid, _ := uuid.Parse(tenantID)
	rows, err := s.queries.ListConfirmedSalesForReport(ctx, db.ListConfirmedSalesForReportParams{
		TenantID: pgutil.UUIDToPg(tid),
		FromDate: pgutil.Timestamptz(from),
		ToDate:   pgutil.Timestamptz(to),
	})
	if err != nil {
		return nil, err
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "goERP - Resumo de Vendas")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(40, 8, fmt.Sprintf("Periodo: %s a %s", from.Format("02/01/2006"), to.Format("02/01/2006")))
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(60, 8, "Cliente")
	pdf.Cell(40, 8, "Data")
	pdf.Cell(40, 8, "Total")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	var grand float64
	for _, row := range rows {
		pdf.Cell(60, 7, row.CustomerName)
		pdf.Cell(40, 7, row.CreatedAt.Time.Format("02/01/2006"))
		total := pgutil.NumericToString(row.Total)
		pdf.Cell(40, 7, "R$ "+total)
		pdf.Ln(7)
		if v, err := row.Total.Float64Value(); err == nil && v.Valid {
			grand += v.Float64
		}
	}
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, fmt.Sprintf("Total confirmado: R$ %.2f", grand))
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ParseRangeQuery(fromStr, toStr string) (time.Time, time.Time) {
	return querytime.ParseRange(fromStr, toStr)
}
