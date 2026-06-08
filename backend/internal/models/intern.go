package models

type InternHealthResponse struct {
	ActiveInterns    int            `json:"activeInterns"`
	CompletedInterns int            `json:"completedInternships"`
	InactiveInterns  int            `json:"inactiveInterns"`
	TopColleges      map[string]int `json:"topColleges"`
	FastestDomains   map[string]int `json:"fastestGrowingDomains"`
}
