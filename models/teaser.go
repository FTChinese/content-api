package models

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type ArticleMeta struct {
	ID         string      `json:"id"`
	CreatedAt  chrono.Time `json:"createdAt"`
	UpdatedAt  chrono.Time `json:"updatedAt"`
	Tags       []string    `json:"tags"`
	MemberTier enum.Tier   `json:"tier"`
}

type TeaserBase struct {
	Title      string      `json:"title"`
	Standfirst string      `json:"standfirst"`
	CoverURL   null.String `json:"coverUrl"`
}

type Teaser struct {
	ArticleMeta
	TeaserBase
}

type InteractiveTeaser struct {
	ArticleMeta
	TeaserBase
	AudioURL null.String `json:"audioUrl"`
}
