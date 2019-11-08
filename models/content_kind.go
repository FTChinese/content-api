package models

type ContentKind int

const (
	ContentKindStory        ContentKind = iota // New report
	ContentKindVideo                           // Video
	ContentKindGallery                         // Photo gallery
	ContentKindAudio                           // Interactive audio
	ContentKindArticle                         // interactive plain article
	ContentKindSpeedReading                    // interactive speed reading
	ContentKindReport                          // Interactive fta report
	ContentKindSponsor                         // Interactive sponsor
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
