package models

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
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

// ChannelSetting represents one row in the channel table.
// This controls the configuration of a channel.
type ChannelSetting struct {
	ID          int64       `json:"id" db:"id"`
	ParentID    int64       `json:"parentId" db:"parent_id"`
	KeyName     string      `json:"keyName" db:"key_name"`
	Name        string      `json:"name"`
	Title       string      `json:"title" db:"title"`
	Description null.String `json:"description" db:"description"`
	CreatedAt   chrono.Time `json:"createdAt" db:"created_utc"`
	UpdatedAt   chrono.Time `json:"updatedAt" db:"updated_utc"`
}

// ChannelPage contains all the data a channel required, including
// its configuration and a list of article summaries.
type ChannelPage struct {
	ChannelSetting
	Data []Teaser `json:"data"`
}
