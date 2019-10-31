package models

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
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

func (r *RawInteractive) Kind() ContentKind {
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

func (r *RawInteractive) Tier() enum.Tier {
	if strings.Contains(r.Tag, "会员专享") {
		return enum.TierStandard
	}

	return enum.InvalidTier
}

func (r *RawInteractive) ArticleMeta() ArticleMeta {
	return ArticleMeta{
		ID:         r.ID,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		Tags:       strings.Split(strings.TrimSpace(r.Tag), ","),
		MemberTier: r.Tier(),
	}
}

func (r *RawInteractive) Vocabularies() []Word {
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

func (r *RawInteractive) AudioTeaser() InteractiveTeaser {
	return InteractiveTeaser{
		ArticleMeta: r.ArticleMeta(),
		TeaserBase: TeaserBase{
			Title:      r.TitleCN,
			Standfirst: r.LongLeadCN,
			CoverURL:   r.CoverURL,
		},
		AudioURL: null.NewString(r.ShortLeadCN, r.ShortLeadCN != ""),
	}
}

func (r *RawInteractive) SpeedReadingTeaser() InteractiveTeaser {
	return InteractiveTeaser{
		ArticleMeta: r.ArticleMeta(),
		TeaserBase: TeaserBase{
			Title:      r.TitleCN,
			Standfirst: r.ShortLeadCN,
			CoverURL:   r.CoverURL,
		},
	}
}

func (r *RawInteractive) Build() Interactive {
	k := r.Kind()
	var i Interactive
	switch k {
	case ContentKindReport, ContentKindSponsor, ContentKindArticle:
		i = NewPlainInteractive(r)

	case ContentKindAudio:
		i = NewAudioArticle(r)

	case ContentKindSpeedReading:
		i = NewSpeedReading(r)

	default:
		i = NewPlainInteractive(r)
	}

	i.Type = k

	return i
}
