package pkg

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type ArticleMeta struct {
	ID         string      `json:"id"`
	Kind       ContentKind `json:"type"`
	CreatedAt  chrono.Time `json:"createdAt"`
	UpdatedAt  chrono.Time `json:"updatedAt"`
	MemberTier enum.Tier   `json:"tier"`
	Title      string      `json:"title"`
}

type Teaser struct {
	ArticleMeta
	Standfirst string      `json:"standfirst"`
	CoverURL   null.String `json:"coverUrl"`
	Tags       []string    `json:"tags"`
	AudioURL   null.String `json:"audioUrl"`
}

//type InteractiveTeaser struct {
//	Teaser
//	AudioURL null.String `json:"audioUrl"`
//}
