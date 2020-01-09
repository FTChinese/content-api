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

func NewBilingualStory(raw *RawStory) BilingualStory {

	s := BilingualStory{
		Teaser:     raw.Teaser(),
		StoryBase:  raw.StoryBase(),
		Body:       []Bilingual{},
		Translator: null.String{},
		AlternativeTitles: AlternativeTitles{
			English:   null.StringFrom(raw.TitleEN),
			Promotion: null.String{},
		},
		Related: raw.Related,
	}

	cnParas, translator := raw.splitCNWithTranslator()
	s.Translator = null.NewString(translator, translator != "")

	pairs := ZipString(cnParas, raw.splitEN())

	for _, v := range pairs {
		s.Body = append(s.Body, Bilingual{
			CN: v.First,
			EN: v.Second,
		})
	}

	return s
}
