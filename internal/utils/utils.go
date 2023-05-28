package utils

func CalculateTotalPages(totalCount int64, limit int64) int {
	return int((totalCount + limit - 1) / limit)
}
