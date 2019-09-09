package models

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"strings"
	"time"
)

// RawStory is used to retrieve a story from db as is.
type RawStory struct {
	ID             string `json:"id" db:"story_id"`
	Bilingual      bool   `json:"bilingual"`
	TitleCN        string `json:"titleCn" db:"title_cn"`
	TitleEN        string `json:"titleEn" db:"title_en"`
	Standfirst     string `json:"standfirst" db:"standfirst"`
	CoverURL       string `json:"coverUrl" db:"cover_url"`
	BylineDescCN   string `json:"bylineDescCn" db:"byline_desc_cn"`
	BylineDescEN   string `json:"bylineDescEn" db:"byline_desc_en"`
	BylineAuthorCN string `json:"bylineAuthorCn" db:"byline_author_cn"`
	BylineAuthorEN string `json:"bylineAuthorEn" db:"byline_author_en"`
	BylineStatusCN string `json:"bylineStatusCn" db:"byline_status_cn"`
	BylineStatusEN string `json:"bylineStatusEn" db:"byline_status_en"`
	AccessRight    int64  `json:"accessRight" db:"access_right"`
	Tag            string `json:"tags" db:"tag"`
	Genre          string `json:"genre" db:"genre"`
	Topic          string `json:"topic" db:"topic"`
	Industry       string `json:"industry" db:"industry"`
	Area           string `json:"area" db:"area"`
	CreatedAt      int64  `json:"createdAt" db:"created_at"`
	UpdatedAt      int64  `json:"updatedAt" db:"updated_at"`
	RawBody
}

func (s *RawStory) Sanitize() {
	s.BodyCN = strings.TrimSpace(s.BodyCN)
	s.BodyEN = strings.TrimSpace(s.BodyEN)
}

func (s *RawStory) SetBilingual() {
	s.Bilingual = s.BodyCN != "" && s.BodyEN != ""
}

func (s *RawStory) BylineCN() Byline {
	var authors []Author

	placeGroups := strings.Split(s.BylineStatusCN, ",")

	// Handle irregular format.
	if len(placeGroups) == 1 && !strings.Contains(s.BylineAuthorCN, ";") {
		return Byline{
			Organization: s.BylineDescCN,
			Authors: []Author{
				{
					Names: strings.Split(s.BylineAuthorCN, ","),
					Place: s.BylineStatusCN,
				},
			},
		}
	}

	nameGroups := strings.Split(s.BylineAuthorCN, ",")
	pairs := AlignStringPairs(nameGroups, placeGroups)

	for _, v := range pairs {

		authors = append(authors, Author{
			Names: strings.Split(v.First, ";"),
			Place: v.Second,
		})
	}

	return Byline{
		Organization: s.BylineDescCN,
		Authors:      authors,
	}
}

func (s *RawStory) BylineEN() Byline {
	var authors []Author

	placeGroups := strings.Split(s.BylineStatusEN, ",")
	// Handle irregular format.
	if len(placeGroups) == 1 && !strings.Contains(s.BylineAuthorEN, ";") {
		return Byline{
			Organization: s.BylineDescEN,
			Authors: []Author{
				{
					Names: strings.Split(s.BylineAuthorEN, ","),
					Place: s.BylineStatusEN,
				},
			},
		}
	}

	nameGroups := strings.Split(s.BylineAuthorEN, ",")
	pairs := AlignStringPairs(nameGroups, placeGroups)

	for _, v := range pairs {

		authors = append(authors, Author{
			Names: strings.Split(v.First, ";"),
			Place: v.Second,
		})
	}

	return Byline{
		Organization: s.BylineDescEN,
		Authors:      authors,
	}
}

func (s *RawStory) MetaData() StoryMeta {
	var tier enum.Tier
	switch s.AccessRight {
	case 1:
		tier = enum.TierStandard
	case 2:
		tier = enum.TierPremium
	default:
		tier = enum.InvalidTier
	}

	return StoryMeta{
		ID:         s.ID,
		Areas:      strings.Split(s.Area, ","),
		Genres:     strings.Split(s.Genre, ","),
		Industries: strings.Split(s.Industry, ","),
		MemberTier: tier,
		Tags:       strings.Split(s.Tag, ","),
		Topics:     strings.Split(s.Topic, ","),
		CreatedAt:  chrono.TimeFrom(time.Unix(s.CreatedAt, 0)),
		UpdatedAt:  chrono.TimeFrom(time.Unix(s.UpdatedAt, 0)),
	}
}
