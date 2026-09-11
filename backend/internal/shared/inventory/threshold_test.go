package inventory

import "testing"

func TestNormalizeLowStockThreshold(t *testing.T) {
	tests := []struct {
		in   int
		want int32
	}{
		{0, DefaultLowStockThreshold},
		{-1, DefaultLowStockThreshold},
		{5, 5},
		{10, 10},
		{maxLowStockThreshold + 1, maxLowStockThreshold},
	}
	for _, tc := range tests {
		if got := NormalizeLowStockThreshold(tc.in); got != tc.want {
			t.Fatalf("NormalizeLowStockThreshold(%d) = %d, want %d", tc.in, got, tc.want)
		}
	}
}
