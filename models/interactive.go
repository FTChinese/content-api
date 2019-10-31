package models

import (
	"encoding/json"
	"github.com/guregu/null"
	"github.com/tidwall/gjson"
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
	InteractiveTeaser
	Type              ContentKind       `json:"type"`
	Byline            null.String       `json:"byline"`
	BodyXML           string            `json:"bodyXml"`
	AlternativeTitles AlternativeTitles `json:"alternativeTitles"`
	Timeline          [][]AudioTimeline `json:"timeline"` // For those with audio and subtitles
	Vocabularies      []Word            `json:"vocabularies"`
	Quiz              null.String       `json:"quiz"` // For Speed reading.
}

// NewPlainInteractive is used to build data for:
// 企业公告,interactive_search,2019吴晓波青年午餐会,去广告
// FT研究院,报告,置顶,去广告,会员专享,interactive_search
// 一周新闻,教程,入门级
func NewPlainInteractive(raw *RawInteractive) Interactive {
	return Interactive{
		InteractiveTeaser: InteractiveTeaser{
			ArticleMeta: raw.ArticleMeta(),
			TeaserBase: TeaserBase{
				Title:      raw.TitleCN,
				Standfirst: raw.LongLeadCN,
				CoverURL:   raw.CoverURL,
			},
			AudioURL: null.String{},
		},
		Byline:  null.NewString(raw.BylineCN, raw.BylineCN != ""),
		BodyXML: raw.BodyCN,
	}
}

// NewAudioArticle is used to build contents for:
// 一波好书,音频,会员专享
// 每日一词,音频,会员专享
// 麦可林学英语,音频,会员专享,interactive_search,英语电台
// 英语电台,interactive_search,
// 英语电台,FT Arts,音乐,音乐之生,interactive_search
// BoomEar艺术播客,音频
func NewAudioArticle(raw *RawInteractive) Interactive {

	var timeline [][]AudioTimeline
	if len(raw.BodyCN) > 0 {
		var result = gjson.Get(raw.BodyCN, "text")
		if result.Exists() {
			_ = json.Unmarshal([]byte(result.String()), &timeline)
		}
	}

	return Interactive{
		InteractiveTeaser: InteractiveTeaser{
			ArticleMeta: raw.ArticleMeta(),
			TeaserBase: TeaserBase{
				Title:      raw.TitleCN,
				Standfirst: raw.LongLeadCN,
				CoverURL:   raw.CoverURL,
			},
			AudioURL: null.NewString(raw.ShortLeadCN, raw.ShortLeadCN != ""),
		},
		Byline:            null.NewString(raw.BylineCN, raw.BylineCN != ""),
		BodyXML:           raw.BodyEN,
		AlternativeTitles: AlternativeTitles{},
		Timeline:          timeline,
		Quiz:              null.String{},
	}
}

// NewSpeedReading is used to build data for:
// 速读,interactive_search,
func NewSpeedReading(raw *RawInteractive) Interactive {
	return Interactive{
		InteractiveTeaser: InteractiveTeaser{
			ArticleMeta: raw.ArticleMeta(),
			TeaserBase: TeaserBase{
				Title:      raw.TitleCN,
				Standfirst: raw.ShortLeadCN,
				CoverURL:   raw.CoverURL,
			},
		},
		Byline:  null.String{},
		BodyXML: raw.BodyEN,
		AlternativeTitles: AlternativeTitles{
			English: raw.TitleEN,
		},
		Timeline:     nil,
		Vocabularies: raw.Vocabularies(),
		Quiz:         null.StringFrom(raw.BodyCN),
	}
}
