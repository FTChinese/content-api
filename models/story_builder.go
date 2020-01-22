package models

import (
	"errors"
	"github.com/guregu/null"
)

func (r RawStory) BuildCN() Story {
	b, t := r.splitCNWithTranslator()
	return Story{
		Teaser:     r.Teaser(),
		StoryBase:  r.StoryBase(),
		Body:       b,
		Translator: null.NewString(t, t != ""),
		Related:    r.Related,
	}
}

func (r RawStory) BuildEN() (Story, error) {
	if !r.HasEN() {
		return Story{}, errors.New("not found")
	}

	m := r.ArticleMeta()
	m.Title = r.TitleEN

	return Story{
		Teaser: Teaser{
			ArticleMeta: m,
			Standfirst:  r.LongLeadCN,
			CoverURL:    r.CoverURL,
			Tags:        r.Tags(),
		},
		StoryBase:  r.StoryBase(),
		Body:       r.splitEN(),
		Translator: null.String{},
		Related:    r.Related,
	}, nil
}

func (r RawStory) BuildBilingual() BilingualStory {
	s := BilingualStory{
		Teaser:     r.Teaser(),
		StoryBase:  r.StoryBase(),
		Body:       []Bilingual{},
		Translator: null.String{},
		AlternativeTitles: AlternativeTitles{
			English:   null.StringFrom(r.TitleEN),
			Promotion: null.String{},
		},
		Related: r.Related,
	}

	cnParas, translator := r.splitCNWithTranslator()
	s.Translator = null.NewString(translator, translator != "")

	pairs := ZipString(cnParas, r.splitEN())

	for _, v := range pairs {
		s.Body = append(s.Body, Bilingual{
			CN: v.First,
			EN: v.Second,
		})
	}

	return s
}
