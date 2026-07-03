package pgutil

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToPg(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func PgToUUID(id pgtype.UUID) (uuid.UUID, error) {
	if !id.Valid {
		return uuid.Nil, nil
	}
	return uuid.FromBytes(id.Bytes[:])
}

func MustPgToUUID(id pgtype.UUID) uuid.UUID {
	u, err := PgToUUID(id)
	if err != nil {
		return uuid.Nil
	}
	return u
}

func TextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	v := t.String
	return &v
}

func TextFromPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func Timestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func BoolFromPtr(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func NumericFromString(s string) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	if err := n.Scan(s); err != nil {
		return pgtype.Numeric{}, err
	}
	return n, nil
}

func NumericToString(n pgtype.Numeric) string {
	if !n.Valid {
		return "0.00"
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", f.Float64)
}
