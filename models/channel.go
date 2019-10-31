package models

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

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
	KeyName     string      `json:"keyName" db:"key_name"` // Human-readable string to represent and search this row.
	Name        string      `json:"name"`                  // The last segment of KeyName
	Title       string      `json:"title" db:"title"`
	Description null.String `json:"description" db:"description"`
	KeyWords    null.String `json:"-" db:"key_words"` // Comma-separated tags to find all articles under this channel.
	ContentKind ContentKind `json:"-"`
	CreatedAt   chrono.Time `json:"createdAt" db:"created_utc"`
	UpdatedAt   chrono.Time `json:"updatedAt" db:"updated_utc"`
}

// ChannelPage contains all the data a channel required, including
// its configuration and a list of article summaries.
type ChannelPage struct {
	ChannelSetting
	Data []Teaser `json:"data"`
}

type InteractivePage struct {
	ChannelSetting
	Data []InteractiveTeaser `json:"data"`
}
