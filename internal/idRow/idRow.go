package idRow

// Captures multiple rows of id fields only
type IDField struct {
	ID int `json:"id"`
}

// ToIntSlice converts capture struct to []int
func ToIntSlice(idFields []IDField) []int {
	ids := make([]int, len(idFields))
	for i, row := range idFields {
		ids[i] = row.ID
	}
	return ids
}
