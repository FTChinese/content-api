package models

import (
	"github.com/guregu/null"
)

type BilingualStory struct {
	Teaser
	StoryBase
	Body              []Bilingual       `json:"body"`
	Translator        null.String       `json:"translator"`
	AlternativeTitles AlternativeTitles `json:"alternativeTitles"`
	Related           []ArticleMeta     `json:"related"`
}
