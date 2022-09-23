package pkg

import (
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
