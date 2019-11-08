package models

import "github.com/FTChinese/go-rest/enum"

// Video represents a video news
type Video struct {
	ID          int         `json:"id" db:"id"`
	Kind        ContentKind `json:"type"`
	CreatedAt   string      `json:"createdAt" db:"created_at"`
	UpdatedAt   string      `json:"updatedAt" db:"updated_at"`
	MemberTier  enum.Tier   `json:"tier"`
	AccessRight int64       `json:"-" db:"access_right"`
	Title       string      `json:"title" db:"title"`
	Standfirst  string      `json:"standfirst" db:"standfirst"`
	PosterURL   string      `json:"posterUrl" db:"posterUrl"`
	CcID        string      `json:"ccId" db:"ccId"`
	Byline      string      `json:"byline" db:"byline"`
}
