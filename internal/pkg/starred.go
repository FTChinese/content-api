package pkg

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTchinese/content-api/pkg/validator"
	"github.com/guregu/null"
)

type StarParams struct {
	UserID string      `json:"userId" db:"user_id"`
	ID     string      `json:"id" db:"id"`
	Kind   ContentKind `json:"kind" db:"kind"`
}

func (p *StarParams) Validate() *render.ValidationError {
	return validator.New("id").Required().Validate(p.ID)
}

// Starred is an article user starred.
type Starred struct {
	StarParams
	StarredAt  int64  `json:"starredAt"`
	Title      string `json:"title"`
	Standfirst string `json:"standfirst"`
}

type RawStarred struct {
	StarParams
	StarredCST       string      `db:"starred_at"` // This is a string of UTC+8 time.
	StoryTitle       null.String `db:"story_title"`
	StoryLead        null.String `db:"story_lead"`
	InteractiveTitle null.String `db:"interact_title"`
	InteractiveLead  null.String `db:"interact_lead"`
	VideoTitle       null.String `db:"video_title"`
	VideoLead        null.String `db:"video_lead"`
	PhotoTitle       null.String `db:"photo_title"`
	PhotoLead        null.String `db:"photo_lead"`
}

func (r *RawStarred) Build() Starred {
	t, _ := chrono.ParseDateTime(r.StarredCST, chrono.TZShanghai)

	var title, standfist string

	switch r.Kind {
	case ContentKindStory:
		title = r.StoryTitle.String
		standfist = r.StoryLead.String

	case ContentKindInteractive:
		title = r.InteractiveTitle.String
		standfist = r.InteractiveLead.String

	case ContentKindVideo:
		title = r.VideoTitle.String
		standfist = r.VideoLead.String

	case ContentKindPhoto:
		title = r.PhotoTitle.String
		standfist = r.PhotoLead.String
	}

	return Starred{
		StarParams: StarParams{
			UserID: r.UserID,
			ID:     r.ID,
			Kind:   r.Kind,
		},
		StarredAt:  t.Unix(),
		Title:      title,
		Standfirst: standfist,
	}
}
