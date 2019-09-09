package models

import (
	"errors"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type StoryMeta struct {
	ID         string      `json:"id"`
	Areas      []string    `json:"areas"`
	Genres     []string    `json:"genres"`
	Industries []string    `json:"industries"`
	MemberTier enum.Tier   `json:"tier"`
	Tags       []string    `json:"tags"`
	Topics     []string    `json:"topics"`
	CreatedAt  chrono.Time `json:"createdAt"`
	UpdatedAt  chrono.Time `json:"updatedAt"`
}

// Story is the monolingual version.
type Story struct {
	StoryMeta
	Bilingual  bool        `json:"bilingual"`
	Title      string      `json:"title"`
	Standfirst string      `json:"standfirst"`
	CoverURL   string      `json:"coverUrl"`
	Byline     Byline      `json:"byline"`
	Body       []string    `json:"body"`
	Translator null.String `json:"translator"`
}

func NewStoryCN(raw RawStory) Story {
	b, t := raw.splitCN()
	return Story{
		StoryMeta:  raw.MetaData(),
		Bilingual:  raw.Bilingual,
		Title:      raw.TitleCN,
		Standfirst: raw.Standfirst,
		CoverURL:   raw.CoverURL,
		Byline:     raw.BylineCN(),
		Body:       b,
		Translator: null.NewString(t, t != ""),
	}
}

func NewStoryEN(raw RawStory) (Story, error) {
	if !raw.HasEN() {
		return Story{}, errors.New("not found")
	}

	return Story{
		StoryMeta:  raw.MetaData(),
		Bilingual:  raw.Bilingual,
		Title:      raw.TitleEN,
		Standfirst: raw.Standfirst,
		CoverURL:   raw.CoverURL,
		Byline:     raw.BylineEN(),
		Body:       raw.splitEN(),
		Translator: null.String{},
	}, nil
}
