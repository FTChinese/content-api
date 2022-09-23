package pkg

import (
	"encoding/json"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"github.com/tidwall/gjson"
	"strings"
)

type RawInteractive struct {
	ID          string      `db:"id"`
	CreatedAt   chrono.Time `db:"created_utc"`
	UpdatedAt   chrono.Time `db:"updated_utc"`
	Tag         string      `db:"tag"`
	TitleCN     string      `db:"title_cn"` // cheadline
	TitleEN     null.String `db:"title_en"`
	LongLeadCN  string      `db:"long_lead_cn"`  // clongleadbody. This is used to hold vocabularies in Speed Reading, and Standfirst for others.
	LongLeadEN  string      `db:"long_lead_en"`  // elongleadbody
	ShortLeadCN string      `db:"short_lead_cn"` // cshortleadbody. This is used to hold Standfirst for Speed Reading, and MP3 url for others.
	BylineCN    string      `db:"byline_cn"`
	CoverURL    null.String `db:"cover_url"`
	// cbody. This is Quiz for Speed Reading, JOSN text for those with MP3 url, plain text for ads.
	// ebody. This is english text for most, but Chinese text for Michael Learn English.
	RawBody
}

func (r RawInteractive) Kind() ContentKind {
	switch {
	case strings.Contains(r.Tag, "企业公告"):
		return ContentKindSponsor

	case strings.Contains(r.Tag, "FT研究院"):
		return ContentKindReport

	case strings.Contains(r.Tag, "一周新闻"):
		return ContentKindArticle

	case strings.Contains(r.Tag, "音频"):
		return ContentKindAudio

	case strings.Contains(r.Tag, "英语电台"):
		return ContentKindAudio

	case strings.Contains(r.Tag, "速读"):
		return ContentKindSpeedReading

	default:
		return ContentKindArticle
	}
}

func (r RawInteractive) Tier() enum.Tier {
	if strings.Contains(r.Tag, "会员专享") {
		return enum.TierStandard
	}

	return enum.TierNull
}

func (r RawInteractive) Tags() []string {
	return strings.Split(strings.TrimSpace(r.Tag), ",")
}

func (r RawInteractive) ArticleMeta() ArticleMeta {
	return ArticleMeta{
		ID:         r.ID,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		MemberTier: r.Tier(),
		Title:      r.TitleCN,
	}
}

// Teaser create a teaser from non-audio articles, excluding
// speed reading, which using MySQL in a perverted way.
func (r RawInteractive) Teaser() Teaser {
	return Teaser{
		ArticleMeta: r.ArticleMeta(),
		Standfirst:  r.LongLeadCN,
		CoverURL:    r.CoverURL,
		Tags:        r.Tags(),
		AudioURL:    null.String{},
	}
}

func (r RawInteractive) AudioTeaser() Teaser {
	return Teaser{
		ArticleMeta: r.ArticleMeta(),
		Standfirst:  r.LongLeadCN,
		CoverURL:    r.CoverURL,
		Tags:        r.Tags(),
		AudioURL:    null.NewString(r.ShortLeadCN, r.ShortLeadCN != ""),
	}
}

// Vocabularies builds the clongleadbody column
// into structured data for speed reading.
func (r RawInteractive) Vocabularies() []Word {
	entries := strings.Split(r.LongLeadCN, "\n")

	var words = make([]Word, 0)
	for _, entry := range entries {
		var word Word
		w := strings.Split(entry, "|")

		if len(w) > 1 {
			word = Word{
				Term:        strings.TrimSpace(w[0]),
				Description: strings.TrimSpace(w[1]),
			}
		} else {
			word = Word{
				Term:        strings.TrimSpace(entry),
				Description: "",
			}
		}

		words = append(words, word)
	}

	return words
}

// PlainInteractive is used to build data for:
// 企业公告,interactive_search,2019吴晓波青年午餐会,去广告
// FT研究院,报告,置顶,去广告,会员专享,interactive_search // Delimited by `\n`
// 一周新闻,教程,入门级
func (r RawInteractive) NewPlainInteractive() Interactive {
	return Interactive{
		Teaser:  r.Teaser(),
		Byline:  null.NewString(r.BylineCN, r.BylineCN != ""),
		BodyXML: r.BodyCN,
	}
}

// NewAudioArticle is used to build contents for:
//                                      Delimiter       Timeline
// ----------------------------------------------------------------
// 一波好书,音频,会员专享                     \n               N
// ----------------------------------------------------------------
// 每日一词,音频,会员专享                     \n               N
// ----------------------------------------------------------------
// 麦可林学英语,音频,会员专享,                 \n              N
// interactive_search,英语电台
// ----------------------------------------------------------------
// 英语电台,interactive_search,             N/A             Y
// ----------------------------------------------------------------
// 英语电台,FT Arts,音乐,音乐之生,            \n for Chinese  Y
// interactive_search                     N/A for English
// ----------------------------------------------------------------
// BoomEar艺术播客,音频                      \n              Y
func (r RawInteractive) NewAudioArticle() Interactive {

	var timeline [][]AudioTimeline
	if len(r.BodyCN) > 0 {
		var result = gjson.Get(r.BodyCN, "text")
		if result.Exists() {
			_ = json.Unmarshal([]byte(result.String()), &timeline)
		}
	}

	return Interactive{
		Teaser:            r.AudioTeaser(),
		Byline:            null.NewString(r.BylineCN, r.BylineCN != ""),
		BodyXML:           r.BodyEN,
		AlternativeTitles: AlternativeTitles{},
		Timeline:          timeline,
		Quiz:              null.String{},
	}
}

// NewSpeedReading is used to build data for:
// 速读,interactive_search,
// Its body is not delimited by any separators.
func (r RawInteractive) NewSpeedReading() Interactive {
	return Interactive{
		Teaser: Teaser{
			ArticleMeta: r.ArticleMeta(),
			Standfirst:  r.ShortLeadCN,
			CoverURL:    r.CoverURL,
			Tags:        r.Tags(),
			AudioURL:    null.String{},
		},
		Byline:  null.String{},
		BodyXML: r.BodyEN,
		AlternativeTitles: AlternativeTitles{
			English: r.TitleEN,
		},
		Timeline:     nil,
		Vocabularies: r.Vocabularies(),
		Quiz:         null.StringFrom(r.BodyCN),
	}
}

func (r RawInteractive) Build() Interactive {
	k := r.Kind()
	var i Interactive
	switch k {
	case ContentKindReport, ContentKindSponsor, ContentKindArticle:
		i = r.NewPlainInteractive()

	case ContentKindAudio:
		i = r.NewAudioArticle()

	case ContentKindSpeedReading:
		i = r.NewSpeedReading()

	default:
		i = r.NewPlainInteractive()
	}

	i.Kind = k

	return i
}
