package models

import (
	"github.com/FTChinese/go-rest/chrono"
	"strings"
)

type Teaser struct {
	ID         string      `json:"id" db:"id"`
	Tag        string      `json:"-"` // Only used as db scan target.
	Tags       []string    `json:"tags"`
	Title      string      `json:"title" db:"title"`
	Standfirst string      `json:"standfirst" db:"standfirst"`
	CoverURL   string      `json:"coverUrl" db:"cover_url"`
	CreatedAt  chrono.Time `json:"createdAt" db:"created_utc"`
	UpdatedAt  chrono.Time `json:"updatedAt" db:"updated_utc"`
}

func (t *Teaser) Normalize() {
	t.Tags = strings.Split(t.Tag, ",")
}

type FrontPage struct {
	Date chrono.Date `json:"date"`
	Data []Teaser    `json:"data"`
}

type ArchivedFrontPage struct {
	Date string   `json:"date"`
	Data []Teaser `json:"data"`
}

type ChannelPage struct {
	Data []Teaser `json:"data"`
}
