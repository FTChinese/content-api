package models

import (
	"errors"
	"github.com/guregu/null"
)

type BilingualStory struct {
	ArticleMeta
	Title      Bilingual   `json:"title"`
	Standfirst string      `json:"standfirst"`
	CoverURL   null.String `json:"coverUrl"`
	StoryBase
	Body       []Bilingual `json:"body"`
	Translator null.String `json:"translator"`
}

func NewBilingualStory(raw *RawStory) (BilingualStory, error) {

	if !raw.Bilingual {
		return BilingualStory{}, errors.New("not found")
	}

	s := BilingualStory{
		ArticleMeta: raw.ArticleMeta(),
		Title: Bilingual{
			CN: raw.TitleCN,
			EN: raw.TitleEN,
		},
		Standfirst: raw.LongLeadCN,
		CoverURL:   raw.CoverURL,
		StoryBase:  raw.StoryBase(),
		Body:       []Bilingual{},
		Translator: null.String{},
	}

	cnParas, translator := raw.splitCN()
	s.Translator = null.NewString(translator, translator != "")

	pairs := AlignStringPairs(cnParas, raw.splitEN())

	for _, v := range pairs {
		s.Body = append(s.Body, Bilingual{
			CN: v.First,
			EN: v.Second,
		})
	}

	return s, nil
}
