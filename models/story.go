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
	Teaser
	StoryBase
	Body       []string      `json:"body"`
	Translator null.String   `json:"translator"`
	Related    []ArticleMeta `json:"related"`
}

func NewStoryCN(raw *RawStory) Story {
	b, t := raw.splitCNWithTranslator()
	return Story{
		Teaser:     raw.Teaser(),
		StoryBase:  raw.StoryBase(),
		Body:       b,
		Translator: null.NewString(t, t != ""),
		Related:    raw.Related,
	}
}

func NewStoryEN(raw *RawStory) (Story, error) {
	if !raw.HasEN() {
		return Story{}, errors.New("not found")
	}

	m := raw.ArticleMeta()
	m.Title = raw.TitleEN

	return Story{
		Teaser: Teaser{
			ArticleMeta: m,
			Standfirst:  raw.LongLeadCN,
			CoverURL:    raw.CoverURL,
			Tags:        raw.Tags(),
		},
		StoryBase:  raw.StoryBase(),
		Body:       raw.splitEN(),
		Translator: null.String{},
		Related:    raw.Related,
	}, nil
}
