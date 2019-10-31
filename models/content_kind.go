package models

type ContentKind int

const (
	ContentKindStory ContentKind = iota
	ContentKindVideo
	ContentKindGallery
	ContentKindAudio
	ContentKindArticle
	ContentKindSpeedReading
	ContentKindReport
	ContentKindSponsor
)

func (k ContentKind) String() string {
	names := [...]string{
		"Story",
		"Video",
		"Gallery",
		"Audio",
		"Article",
		"SpeedReading",
		"FTAReport",
		"Sponsor",
	}

	if k < ContentKindStory || k > ContentKindSponsor {
		return ""
	}

	return names[k]
}

func (k ContentKind) MarshalJSON() ([]byte, error) {
	s := k.String()

	if s == "" {
		return []byte("null"), nil
	}

	return []byte(`"` + s + `"`), nil
}
