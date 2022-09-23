package pkg

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"strings"
)

func numToTier(n int64) enum.Tier {
	var tier enum.Tier
	switch n {
	case 1:
		tier = enum.TierStandard
	case 2:
		tier = enum.TierPremium
	default:
		tier = enum.TierNull
	}

	return tier
}

type RawContentBase struct {
	ID          string      `json:"id" db:"id"`
	CreatedAt   int64       `json:"createdAt" db:"created_at"`
	UpdatedAt   int64       `json:"updatedAt" db:"updated_at"`
	AccessRight int64       `json:"accessRight" db:"access_right"`
	TitleCN     string      `json:"titleCn" db:"title_cn"`
	LongLeadCN  string      `json:"standfirst" db:"long_lead_cn"`
	CoverURL    null.String `json:"coverUrl" db:"cover_url"`
	Tag         string      `json:"tags" db:"tag"`
}

// ArticleMeta create the meta data of an article.
func (r RawContentBase) ArticleMeta() ArticleMeta {
	return ArticleMeta{
		ID:         r.ID,
		Kind:       ContentKindStory,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		MemberTier: numToTier(r.AccessRight),
		Title:      r.TitleCN,
	}
}

func (r RawContentBase) Teaser() Teaser {
	return Teaser{
		ArticleMeta: r.ArticleMeta(),
		Standfirst:  r.LongLeadCN,
		CoverURL:    r.CoverURL,
		Tags:        r.Tag,
	}
}

// RawStory is used to retrieve an article from db as is.
type RawStory struct {
	RawContentBase
	Bilingual      bool          `json:"bilingual"`
	TitleEN        string        `json:"titleEn" db:"title_en"`
	BylineDescCN   string        `json:"bylineDescCn" db:"byline_desc_cn"`
	BylineDescEN   string        `json:"bylineDescEn" db:"byline_desc_en"`
	BylineAuthorCN string        `json:"bylineAuthorCn" db:"byline_author_cn"`
	BylineAuthorEN string        `json:"bylineAuthorEn" db:"byline_author_en"`
	BylineStatusCN string        `json:"bylineStatusCn" db:"byline_status_cn"`
	BylineStatusEN string        `json:"bylineStatusEn" db:"byline_status_en"`
	Genre          string        `json:"genre" db:"genre"`
	Topic          string        `json:"topic" db:"topic"`
	Industry       string        `json:"industry" db:"industry"`
	Area           string        `json:"area" db:"area"`
	BodyCN         string        `json:"bodyCn" db:"body_cn"`
	BodyEN         string        `json:"bodyEn" db:"body_en"`
	Related        []ArticleMeta `json:"related"`
}

func (r *RawStory) Normalize() {
	//r.CoverURL = imageBaseURL + r.CoverURL
	r.Bilingual = r.BodyCN != "" && r.BodyEN != ""
}

func (r *RawStory) Sanitize() {
	r.BodyCN = strings.TrimSpace(r.BodyCN)
	r.BodyEN = strings.TrimSpace(r.BodyEN)
}

func (r RawStory) isBilingual() bool {
	return r.BodyCN != "" && r.BodyEN != ""
}
