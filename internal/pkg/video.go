package pkg

import "github.com/FTChinese/go-rest/enum"

// Video represents a video news
type Video struct {
	ID         int64       `json:"id" db:"id"`
	Kind       ContentKind `json:"type"`
	CreatedAt  string      `json:"createdAt" db:"created_at"`
	UpdatedAt  string      `json:"updatedAt" db:"updated_at"`
	MemberTier enum.Tier   `json:"tier"`
	Title      string      `json:"title" db:"title"`
	Standfirst string      `json:"standfirst" db:"standfirst"`
	CoverURL   string      `json:"coverUrl" db:"poster_url"`
	CcID       string      `json:"ccId" db:"ccId"`
	Byline     string      `json:"byline" db:"byline"`
}

// RawVideo is the scan target in SQL.
type RawVideo struct {
	ID        int64  `db:"id"`
	CreatedAt string `db:"created_date"`
	UpdateAt  string `db:"updated_date"`
	RawPerm
	Title        string `db:"title"`
	LongLeadCN   string `db:"long_lead_cn"`
	PostURL      string `db:"poster_url"`
	CcID         string `db:"cc_id"`
	BylineDescCN string `db:"byline_desc_cn"`
	BylineCN     string `db:"byline_cn"`
}

// Build converts the raw data retrieved from DB to Video type.
func (r RawVideo) Build() Video {
	return Video{
		ID:         r.ID,
		Kind:       ContentKindVideo,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdateAt,
		MemberTier: r.MemberTier(),
		Title:      r.Title,
		Standfirst: r.LongLeadCN,
		CoverURL:   r.PostURL,
		CcID:       r.CcID,
		Byline:     r.BylineDescCN + " " + r.BylineCN,
	}
}
