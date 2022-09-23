package pkg

import (
	"github.com/guregu/null"
)

type Byline struct {
	Organization string `json:"organization"`
	Authors      string `json:"authors"`
	Locations    string `json:"locations"`
}

type StoryDetails struct {
	Title  string `json:"title"`
	Byline Byline `json:"byline"`
	Body   string `json:"body"`
}

// Story is the monolingual version.
type Story struct {
	Teaser
	Bilingual  bool          `json:"bilingual"`
	Regions    string        `json:"regions"`
	Genres     string        `json:"genres"`
	Industries string        `json:"industries"`
	Topics     string        `json:"topics"`
	ContentCN  StoryDetails  `json:"contentCn"`
	ContentEn  StoryDetails  `json:"contentEn"`
	Related    []ArticleMeta `json:"related"` // Deprecated
}

func NewStory(raw RawStory) Story {
	return Story{
		Teaser: Teaser{
			ArticleMeta: raw.ArticleMeta(),
			Standfirst:  raw.LongLeadCN,
			CoverURL:    raw.CoverURL,
			Tags:        raw.Tag,
			AudioURL:    null.String{},
		},
		Bilingual:  raw.isBilingual(),
		Regions:    raw.Area,
		Genres:     raw.Genre,
		Industries: raw.Industry,
		Topics:     raw.Topic,
		ContentCN: StoryDetails{
			Title: raw.TitleCN,
			Byline: Byline{
				Organization: raw.BylineStatusCN,
				Authors:      raw.BylineAuthorCN,
				Locations:    raw.BylineDescCN,
			},
			Body: raw.BodyCN,
		},
		ContentEn: StoryDetails{
			Title: raw.TitleEN,
			Byline: Byline{
				Organization: raw.BylineStatusEN,
				Authors:      raw.BylineAuthorEN,
				Locations:    raw.BylineDescEN,
			},
			Body: raw.BodyCN,
		},
	}
}
