package models

import (
	"encoding/json"
	"github.com/guregu/null"
	"github.com/tidwall/gjson"
	"strings"
)

type AudioTimeline struct {
	Start null.Float `json:"start"`
	Text  string     `json:"text"`
}

// 英语电台,FT Arts,音乐,音乐之生,interactive_search
// Standfirst - clongleadbody
// BodyJSON - cbody
// BodyXML - ebody
type AudioArticle struct {
	Teaser
	AudioURL    null.String       `json:"audioUrl" db:"audio_url"`
	Byline      null.String       `json:"byline" db:"byline"`
	RawBodyJSON string            `json:"-" db:"raw_body_json"`
	Body        [][]AudioTimeline `json:"body"`
	RawBodyXML  string            `json:"-" db:"raw_body_xml"`
	BodyXML     []string          `json:"bodyXml"`
}

func (a *AudioArticle) Normalize() {
	a.Teaser.Normalize()

	a.BodyXML = strings.Split(a.RawBodyXML, "\n")

	result := gjson.Get(a.RawBodyJSON, "text")

	if result.Exists() {
		_ = json.Unmarshal([]byte(result.String()), &a.Body)
	}
}

type Word struct {
	Term        string `json:"term"`
	Description string `json:"description"`
}

type SpeedReading struct {
	Teaser
	RawVocab     string `json:"-" db:"raw_vocab"`
	Vocabularies []Word `json:"vocabularies"`
	TitleEN      string `json:"titleEn" db:"title_en"`
	Body         string `json:"body" db:"raw_body"`
	Quiz         string `json:"quiz" db:"quiz"`
}

func (s *SpeedReading) Normalize() {
	s.Teaser.Normalize()

	entries := strings.Split(s.RawVocab, "\n")

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

		s.Vocabularies = append(s.Vocabularies, word)
	}
}
