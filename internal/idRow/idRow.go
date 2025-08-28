package idRow

// Captures only the id field of a row
type IDRow struct {
	ID int `json:"id"`
}

// Captures multiple rows of id fields only
type IDRows struct {
	rows []IDRow
}

// ToIntSlice converts capture struct to []int
func (idRows *IDRows) ToIntSlice() []int {
	ids := make([]int, len(idRows.rows))
	for i, row := range idRows.rows {
		ids[i] = row.ID
	}
	return ids
}
