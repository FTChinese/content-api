package models

import (
	"errors"
	"github.com/guregu/null"
)

type BilingualStory struct {
	StoryMeta
	Title      Bilingual   `json:"title"`
	Standfirst string      `json:"standfirst"`
	CoverURL   string      `json:"coverUrl"`
	BylineCN   Byline      `json:"bylineCn"`
	BylineEN   Byline      `json:"bylineEn"`
	Body       []Bilingual `json:"body"`
	Translator null.String `json:"translator"`
}

func NewBilingualStory(raw RawStory) (BilingualStory, error) {

	if !raw.Bilingual {
		return BilingualStory{}, errors.New("not found")
	}

	s := BilingualStory{
		StoryMeta: raw.MetaData(),
		Title: Bilingual{
			CN: raw.TitleCN,
			EN: raw.TitleEN,
		},
		Standfirst: raw.Standfirst,
		CoverURL:   raw.CoverURL,
		BylineCN:   raw.BylineCN(),
		BylineEN:   raw.BylineEN(),
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
