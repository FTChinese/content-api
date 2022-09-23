package pkg

import (
	"github.com/guregu/null"
)

type AudioTimeline struct {
	Start null.Float `json:"start"`
	Text  string     `json:"text"`
}

type Word struct {
	Term        string `json:"term"`
	Description string `json:"description"`
}

type AlternativeTitles struct {
	English   null.String `json:"english"`
	Promotion null.String `json:"promotion"`
}

type Interactive struct {
	Teaser
	Byline            null.String       `json:"byline"`
	BodyXML           string            `json:"bodyXml"`
	AlternativeTitles AlternativeTitles `json:"alternativeTitles"`
	Timeline          [][]AudioTimeline `json:"timeline"` // For those with audio and subtitles
	Vocabularies      []Word            `json:"vocabularies"`
	Quiz              null.String       `json:"quiz"` // For Speed reading.
}
