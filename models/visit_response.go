package models

type VisitResponse struct {
	Mark      int    `json:"mark"`
	VisitedAt int    `json:"visited_at"`
	Place     string `json:"place"`
}

type ByVisitedAsc []VisitResponse

func (a ByVisitedAsc) Len() int           { return len(a) }
func (a ByVisitedAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVisitedAsc) Less(i, j int) bool { return a[i].VisitedAt < a[j].VisitedAt }
