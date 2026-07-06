package querytime

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseOptional(s string) pgtype.Timestamptz {
	var t pgtype.Timestamptz
	if s == "" {
		return t
	}
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		parsed, err = time.Parse("2006-01-02", s)
		if err != nil {
			return t
		}
	}
	_ = t.Scan(parsed)
	return t
}

func DefaultRange() (time.Time, time.Time) {
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -30)
	return from, to
}

func ParseRange(fromStr, toStr string) (time.Time, time.Time) {
	from, to := DefaultRange()
	if fromStr != "" {
		if t := ParseOptional(fromStr); t.Valid {
			from = t.Time
		}
	}
	if toStr != "" {
		if t := ParseOptional(toStr); t.Valid {
			to = t.Time
		}
	}
	return from, to
}
