package models

// Video represents a video news
type Video struct {
	ID         int    `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	Standfirst string `json:"standfirst" db:"standfirst"`
	Byline     string `json:"byline" db:"byline"`
	CcID       string `json:"ccId" db:"ccId"`
	PosterURL  string `json:"posterUrl" db:"posterUrl"`
	CreatedAt  string `json:"createdAt" db:"createdAt"`
	UpdatedAt  string `json:"updatedAt" db:"updatedAt"`
}
