package models

import (
	"errors"
	"github.com/guregu/null"
)

type StoryBase struct {
	Bilingual  bool     `json:"bilingual"`
	Byline     Byline   `json:"byline"`
	Areas      []string `json:"areas"`
	Genres     []string `json:"genres"`
	Industries []string `json:"industries"`
	Topics     []string `json:"topics"`
}

// Story is the monolingual version.
type Story struct {
	ArticleMeta
	TeaserBase
	StoryBase
	Body       []string    `json:"body"`
	Translator null.String `json:"translator"`
}

func NewStoryCN(raw *RawStory) Story {
	b, t := raw.splitCN()
	return Story{
		ArticleMeta: raw.ArticleMeta(),
		TeaserBase:  raw.TeaserBaseCN(),
		StoryBase:   raw.StoryBase(),
		Body:        b,
		Translator:  null.NewString(t, t != ""),
	}
}

func NewStoryEN(raw *RawStory) (Story, error) {
	if !raw.HasEN() {
		return Story{}, errors.New("not found")
	}

	return Story{
		ArticleMeta: raw.ArticleMeta(),
		TeaserBase:  raw.TeaserBaseEN(),
		StoryBase:   raw.StoryBase(),
		Body:        raw.splitEN(),
		Translator:  null.String{},
	}, nil
}
