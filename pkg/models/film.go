package models

type Film struct {
	Title     string         `json:"title" db:"title"`
	CoverUrl  string         `json:"cover_url"`
	AvgRating float64        `json:"avg_rating" db:"avg_rating"`
	Reviews   []ReviewOutput `json:"reviews"`
}
