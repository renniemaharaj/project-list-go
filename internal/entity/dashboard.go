package entity

// The defined dashboard metrics type
type MetricsDashboard struct {
	Projects               int     `json:"projects"`
	Active                 int     `json:"active"`
	Completed              int     `json:"completed"`
	Idle                   int     `json:"idle"`
	OutOfBudget            int     `json:"outOfBudget"`
	TotalDebit             float64 `json:"totalDebit"`
	TotalCredit            float64 `json:"totalCredit"`
	AverageCreditOverDebit float32 `json:"avgCreditOverDebit"`
	EndingSoon             int     `json:"endingSoonCount"`
}
