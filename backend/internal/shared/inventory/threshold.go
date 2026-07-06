package inventory

const DefaultLowStockThreshold = 5

const maxLowStockThreshold = 1_000_000

// NormalizeLowStockThreshold returns a safe int32 threshold for low-stock queries.
// Values <= 0 use DefaultLowStockThreshold; values above maxLowStockThreshold are capped.
func NormalizeLowStockThreshold(n int) int32 {
	if n <= 0 {
		return DefaultLowStockThreshold
	}
	if n > maxLowStockThreshold {
		return maxLowStockThreshold
	}
	return int32(n)
}
