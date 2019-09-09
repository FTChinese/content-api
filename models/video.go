package models

// Video represents a video news
type Video struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Standfirst string `json:"standfirst"`
	Byline     string `json:"byline"`
	CcID       string `json:"ccId"`
	PosterURL  string `json:"posterUrl"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}
