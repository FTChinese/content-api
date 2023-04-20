package pkg

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ContentKind int

const (
	ContentKindNull  ContentKind = iota
	ContentKindStory             // News
	ContentKindInteractive
	ContentKindVideo // Video
	ContentKindPhoto // Photo gallery
	ContentKindOther
	ContentKindAudio        // Interactive audio
	ContentKindArticle      // interactive plain article
	ContentKindSpeedReading // interactive speed reading
	ContentKindReport       // Interactive fta report
	ContentKindSponsor      // Interactive sponsor
)

var contentKindNames = [...]string{
	"",
	"story",
	"interactive",
	"video",
	"photo",
	"other",
	"audio",
	"article",
	"speed_reading",
	"fta_report",
	"sponsor",
}

var contentKindMap = map[ContentKind]string{
	ContentKindStory:        contentKindNames[1],
	ContentKindInteractive:  contentKindNames[2],
	ContentKindVideo:        contentKindNames[3],
	ContentKindPhoto:        contentKindNames[4],
	ContentKindOther:        contentKindNames[5],
	ContentKindAudio:        contentKindNames[6],
	ContentKindArticle:      contentKindNames[7],
	ContentKindSpeedReading: contentKindNames[8],
	ContentKindReport:       contentKindNames[9],
	ContentKindSponsor:      contentKindNames[10],
}

var contentKindValue = map[string]ContentKind{
	contentKindNames[1]:  ContentKindStory,
	contentKindNames[2]:  ContentKindInteractive,
	contentKindNames[3]:  ContentKindVideo,
	contentKindNames[4]:  ContentKindPhoto,
	contentKindNames[5]:  ContentKindOther,
	contentKindNames[6]:  ContentKindAudio,
	contentKindNames[7]:  ContentKindArticle,
	contentKindNames[8]:  ContentKindSpeedReading,
	contentKindNames[9]:  ContentKindReport,
	contentKindNames[10]: ContentKindSponsor,
}

func ParseContentKind(name string) (ContentKind, error) {
	if x, ok := contentKindValue[name]; ok {
		return x, nil
	}

	return ContentKindNull, fmt.Errorf("%s is not a valid PayMethod", name)
}

func (x ContentKind) String() string {
	if s, ok := contentKindMap[x]; ok {
		return s
	}

	return ""
}

// UnmarshalJSON implements the Unmarshaler interface.
func (x *ContentKind) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp, _ := ParseContentKind(s)

	*x = tmp

	return nil
}

func (x ContentKind) MarshalJSON() ([]byte, error) {
	s := x.String()

	if s == "" {
		return []byte("null"), nil
	}

	return []byte(`"` + s + `"`), nil
}

// Scan implements sql.Scanner interface to retrieve value from SQL.
// SQL null will be turned into zero value InvalidPay.
func (x *ContentKind) Scan(src interface{}) error {
	if src == nil {
		*x = ContentKindNull
		return nil
	}

	switch s := src.(type) {
	case []byte:
		tmp, _ := ParseContentKind(string(s))
		*x = tmp
		return nil

	default:
		return errors.New("incompatible type to scan")
	}
}

// Value implements driver.Valuer interface to save value into SQL.
func (x ContentKind) Value() (driver.Value, error) {
	s := x.String()
	if s == "" {
		return nil, nil
	}

	return s, nil
}
