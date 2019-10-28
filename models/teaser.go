package models

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"strings"
)

type Teaser struct {
	ID         string      `json:"id" db:"id"`
	CreatedAt  chrono.Time `json:"createdAt" db:"created_utc"`
	UpdatedAt  chrono.Time `json:"updatedAt" db:"updated_utc"`
	RawTag     string      `json:"-" db:"tag"` // Only used as db scan target.
	Tags       []string    `json:"tags"`
	MemberTier enum.Tier   `json:"tier"`
	Title      string      `json:"title" db:"title"`
	Standfirst string      `json:"standfirst" db:"standfirst"`
	CoverURL   null.String `json:"coverUrl" db:"cover_url"`
}

func (t *Teaser) Normalize() {
	t.Tags = strings.Split(t.RawTag, ",")

	if strings.Contains(t.RawTag, "会员专享") {
		t.MemberTier = enum.TierStandard
	}
}
