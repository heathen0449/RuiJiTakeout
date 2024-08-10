package models

type PageResult struct {
	Total   int64         `json:"total"`
	Records []interface{} `json:"records"`
}

// ToInterfaceSlice 将切片转换为接口切片
func ToInterfaceSlice[T any](slice []T) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
