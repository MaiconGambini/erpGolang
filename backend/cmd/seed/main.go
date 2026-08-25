package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/config"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := db.New(pool)

	tenants := []struct {
		slug, name, email, password, userName string
	}{
		{"acme", "Acme Corp", "admin@acme.com", "admin123", "Acme Admin"},
		{"beta", "Beta Ltd", "admin@beta.com", "admin123", "Beta Admin"},
	}

	for _, t := range tenants {
		if err := seedTenant(ctx, queries, t.slug, t.name, t.email, t.password, t.userName, cfg.BcryptCost); err != nil {
			log.Printf("seed %s: %v", t.slug, err)
		} else {
			fmt.Printf("seeded tenant %s (%s)\n", t.slug, t.email)
		}
	}
	if err := seedViewer(ctx, queries, "acme", "viewer@acme.com", "admin123", "Acme Viewer", cfg.BcryptCost); err != nil {
		log.Printf("seed viewer: %v", err)
	} else {
		fmt.Println("seeded viewer@acme.com (role viewer)")
	}
	if err := seedDemoData(ctx, pool, "acme"); err != nil {
		log.Printf("seed demo data: %v", err)
	} else {
		fmt.Println("seeded demo dataset (acme)")
	}
}

func seedViewer(ctx context.Context, q *db.Queries, slug, email, password, userName string, cost int) error {
	tenant, err := q.GetTenantBySlug(ctx, slug)
	if err != nil {
		return err
	}
	_, err = q.GetUserByEmail(ctx, email)
	if err == nil {
		return fmt.Errorf("already exists")
	}
	hash, err := auth.HashPassword(password, cost)
	if err != nil {
		return err
	}
	_, err = q.CreateUser(ctx, db.CreateUserParams{
		TenantID:     tenant.ID,
		Email:        email,
		PasswordHash: hash,
		Name:         userName,
		Role:         "viewer",
	})
	return err
}

func seedTenant(ctx context.Context, q *db.Queries, slug, name, email, password, userName string, cost int) error {
	_, err := q.GetTenantBySlug(ctx, slug)
	if err == nil {
		return fmt.Errorf("already exists")
	}
	tenant, err := q.CreateTenant(ctx, db.CreateTenantParams{Slug: slug, Name: name})
	if err != nil {
		return err
	}
	hash, err := auth.HashPassword(password, cost)
	if err != nil {
		return err
	}
	_, err = q.CreateUser(ctx, db.CreateUserParams{
		TenantID:     tenant.ID,
		Email:        email,
		PasswordHash: hash,
		Name:         userName,
		Role:         "admin",
	})
	return err
}

// seedDemoData populates the acme tenant with a believable dataset so a fresh
// checkout shows a live dashboard (charts need confirmed sales spread across
// recent days). Idempotent: skipped once the tenant has any products, since
// E2E tests and manual use create their own rows.
func seedDemoData(ctx context.Context, pool *pgxpool.Pool, slug string) error {
	queries := db.New(pool)
	tenant, err := queries.GetTenantBySlug(ctx, slug)
	if err != nil {
		return err
	}
	var productCount int64
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM products WHERE tenant_id = $1 AND deleted_at IS NULL`,
		tenant.ID).Scan(&productCount); err != nil {
		return err
	}
	if productCount > 0 {
		return fmt.Errorf("already has data")
	}

	customers := []struct{ name, document, docType, email, city, state string }{
		{"Ana Ribeiro", "12345678909", "cpf", "ana.ribeiro@example.com", "São Paulo", "SP"},
		{"Construtora Vale Verde", "12345678000195", "cnpj", "compras@valeverde.example.com", "Campinas", "SP"},
		{"Bruno Almeida", "98765432100", "cpf", "bruno.almeida@example.com", "Belo Horizonte", "MG"},
		{"Padaria Pão Dourado", "11222333000181", "cnpj", "pedidos@paodourado.example.com", "Curitiba", "PR"},
		{"Carla Mendes", "45678912300", "cpf", "carla.mendes@example.com", "Porto Alegre", "RS"},
		{"Mercado Bom Preço", "22333444000122", "cnpj", "estoque@bompreco.example.com", "Santos", "SP"},
		{"Diego Fontes", "32165498700", "cpf", "diego.fontes@example.com", "Recife", "PE"},
		{"Clínica Vida Saudável", "55666777000133", "cnpj", "financeiro@vidasaudavel.example.com", "Ribeirão Preto", "SP"},
	}
	customerIDs := make([]uuid.UUID, 0, len(customers))
	for i, c := range customers {
		var id uuid.UUID
		// created_at backdated so the "new customers (30d)" KPI has signal.
		ageDays := 40 - i*5
		if ageDays < 5 {
			ageDays = 5 + i
		}
		if err := pool.QueryRow(ctx,
			`INSERT INTO customers (tenant_id, name, document, document_type, email, city, state, active, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, true, now() - make_interval(days => $8), now())
			 RETURNING id`,
			tenant.ID, c.name, c.document, c.docType, c.email, c.city, c.state, ageDays).Scan(&id); err != nil {
			return fmt.Errorf("customer %s: %w", c.name, err)
		}
		customerIDs = append(customerIDs, id)
	}

	products := []struct {
		name, sku string
		price     string
		stock     int32
	}{
		{"Teclado Mecânico TKL", "DEMO-TEC-001", "289.90", 42},
		{"Mouse Sem Fio Pro", "DEMO-MOU-002", "149.90", 3},
		{"Monitor 24\" IPS", "DEMO-MON-003", "899.00", 18},
		{"Headset USB ANC", "DEMO-HEAD-004", "379.90", 27},
		{"Cadeira Ergonômica", "DEMO-CAD-005", "1250.00", 9},
		{"Webcam Full HD", "DEMO-WEB-006", "219.90", 34},
	}
	productIDs := make([]uuid.UUID, 0, len(products))
	productPrices := make([]string, 0, len(products))
	for _, p := range products {
		var id uuid.UUID
		if err := pool.QueryRow(ctx,
			`INSERT INTO products (tenant_id, name, sku, price, stock, unit, active, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, 'UN', true, now() - interval '50 days', now())
			 RETURNING id`,
			tenant.ID, p.name, p.sku, p.price, p.stock).Scan(&id); err != nil {
			return fmt.Errorf("product %s: %w", p.sku, err)
		}
		productIDs = append(productIDs, id)
		productPrices = append(productPrices, p.price)
	}

	// 30 confirmed sales spread over the last 45 days (deterministic spread),
	// plus 3 drafts. Totals = sum of line totals, rounded to 2 decimals.
	insertSale := `INSERT INTO sales (tenant_id, customer_id, status, total, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NULL, $5, $5) RETURNING id`
	insertItem := `INSERT INTO sale_items (tenant_id, sale_id, product_id, quantity, unit_price, line_total, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	for i := 0; i < 33; i++ {
		status := "confirmed"
		if i >= 30 {
			status = "draft"
		}
		daysAgo := (i * 3) % 45
		if status == "draft" {
			daysAgo = i - 30
		}
		saleAt := time.Now().AddDate(0, 0, -daysAgo).Truncate(time.Second)
		cust := customerIDs[i%len(customerIDs)]

		line1Idx := i % len(productIDs)
		line2Idx := (i*2 + 1) % len(productIDs)
		q1 := int32(1 + i%3)
		q2 := int32(1 + (i+1)%2)
		t1 := lineTotal(productPrices[line1Idx], int(q1))
		t2 := lineTotal(productPrices[line2Idx], int(q2))
		total := round2(t1 + t2)

		var saleID uuid.UUID
		if err := pool.QueryRow(ctx, insertSale,
			tenant.ID, cust, status, total, saleAt).Scan(&saleID); err != nil {
			return fmt.Errorf("sale %d: %w", i, err)
		}
		items := []struct {
			productID uuid.UUID
			qty       int32
			price     string
			total     string
		}{
			{productIDs[line1Idx], q1, productPrices[line1Idx], format2(t1)},
			{productIDs[line2Idx], q2, productPrices[line2Idx], format2(t2)},
		}
		for _, it := range items {
			if _, err := pool.Exec(ctx, insertItem,
				tenant.ID, saleID, it.productID, it.qty, it.price, it.total, saleAt); err != nil {
				return fmt.Errorf("sale %d item: %w", i, err)
			}
		}
	}
	return nil
}

func parsePrice(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func lineTotal(price string, qty int) float64 {
	return math.Round(parsePrice(price)*float64(qty)*100) / 100
}

func round2(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', 2, 64)
}

func format2(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}
