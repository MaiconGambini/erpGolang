package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/config"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/database"
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

