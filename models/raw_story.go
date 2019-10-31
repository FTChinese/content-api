package models

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"strings"
)

// RawStory is used to retrieve an article from db as is.
type RawStory struct {
	ID             string      `json:"id" db:"id"`
	Bilingual      bool        `json:"bilingual"`
	TitleCN        string      `json:"titleCn" db:"title_cn"`
	TitleEN        string      `json:"titleEn" db:"title_en"`
	LongLeadCN     string      `json:"standfirst" db:"long_lead_cn"`
	CoverURL       null.String `json:"coverUrl" db:"cover_url"`
	BylineDescCN   string      `json:"bylineDescCn" db:"byline_desc_cn"`
	BylineDescEN   string      `json:"bylineDescEn" db:"byline_desc_en"`
	BylineAuthorCN string      `json:"bylineAuthorCn" db:"byline_author_cn"`
	BylineAuthorEN string      `json:"bylineAuthorEn" db:"byline_author_en"`
	BylineStatusCN string      `json:"bylineStatusCn" db:"byline_status_cn"`
	BylineStatusEN string      `json:"bylineStatusEn" db:"byline_status_en"`
	AccessRight    int64       `json:"accessRight" db:"access_right"`
	Tag            string      `json:"tags" db:"tag"`
	Genre          string      `json:"genre" db:"genre"`
	Topic          string      `json:"topic" db:"topic"`
	Industry       string      `json:"industry" db:"industry"`
	Area           string      `json:"area" db:"area"`
	CreatedAt      chrono.Time `json:"createdAt" db:"created_utc"`
	UpdatedAt      chrono.Time `json:"updatedAt" db:"updated_utc"`
	RawBody
}

func (r *RawStory) Normalize() {
	//r.CoverURL = imageBaseURL + r.CoverURL
	r.Bilingual = r.BodyCN != "" && r.BodyEN != ""
}

func (r *RawStory) isBilingual() bool {
	return r.BodyCN != "" && r.BodyEN != ""
}

func (r *RawStory) Sanitize() {
	r.BodyCN = strings.TrimSpace(r.BodyCN)
	r.BodyEN = strings.TrimSpace(r.BodyEN)
}

func (r *RawStory) BylineCN() Byline {
	var authors []Authors

	placeGroups := strings.Split(r.BylineStatusCN, ",")

	// Handle irregular format.
	if len(placeGroups) == 1 && !strings.Contains(r.BylineAuthorCN, ";") {
		return Byline{
			Organization: r.BylineDescCN,
			Authors: []Authors{
				{
					Names: strings.Split(r.BylineAuthorCN, ","),
					Place: r.BylineStatusCN,
				},
			},
		}
	}

	nameGroups := strings.Split(r.BylineAuthorCN, ",")
	pairs := AlignStringPairs(nameGroups, placeGroups)

	for _, v := range pairs {

		authors = append(authors, Authors{
			Names: strings.Split(v.First, ";"),
			Place: v.Second,
		})
	}

	return Byline{
		Organization: r.BylineDescCN,
		Authors:      authors,
	}
}

func (r *RawStory) BylineEN() Byline {
	var authors []Authors

	placeGroups := strings.Split(r.BylineStatusEN, ",")
	// Handle irregular format.
	if len(placeGroups) == 1 && !strings.Contains(r.BylineAuthorEN, ";") {
		return Byline{
			Organization: r.BylineDescEN,
			Authors: []Authors{
				{
					Names: strings.Split(r.BylineAuthorEN, ","),
					Place: r.BylineStatusEN,
				},
			},
		}
	}

	nameGroups := strings.Split(r.BylineAuthorEN, ",")
	pairs := AlignStringPairs(nameGroups, placeGroups)

	for _, v := range pairs {

		authors = append(authors, Authors{
			Names: strings.Split(v.First, ";"),
			Place: v.Second,
		})
	}

	return Byline{
		Organization: r.BylineDescEN,
		Authors:      authors,
	}
}

func (r *RawStory) ArticleMeta() ArticleMeta {
	var tier enum.Tier
	switch r.AccessRight {
	case 1:
		tier = enum.TierStandard
	case 2:
		tier = enum.TierPremium
	default:
		tier = enum.InvalidTier
	}

	return ArticleMeta{
		ID:         r.ID,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		Tags:       strings.Split(r.Tag, ","),
		MemberTier: tier,
	}
}

func (r *RawStory) TeaserBaseCN() TeaserBase {
	return TeaserBase{
		Title:      r.TitleCN,
		Standfirst: r.LongLeadCN,
		CoverURL:   r.CoverURL,
	}
}

func (r *RawStory) TeaserBaseEN() TeaserBase {
	return TeaserBase{
		Title:      r.TitleEN,
		Standfirst: r.LongLeadCN,
		CoverURL:   r.CoverURL,
	}
}

func (r *RawStory) StoryBase() StoryBase {
	return StoryBase{
		Bilingual:  r.isBilingual(),
		Byline:     r.BylineCN(),
		Areas:      strings.Split(r.Area, ","),
		Genres:     strings.Split(r.Genre, ","),
		Industries: strings.Split(r.Industry, ","),
		Topics:     strings.Split(r.Topic, ","),
	}
}

func (r *RawStory) Teaser() Teaser {
	return Teaser{
		ArticleMeta: r.ArticleMeta(),
		TeaserBase:  r.TeaserBaseCN(),
	}
}
