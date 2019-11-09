package models

import "github.com/FTChinese/go-rest/enum"

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

// NewVideo converts the raw data retrieved from DB to Video type.
func NewVideo(raw *RawVideo) Video {
	return Video{
		ID:         raw.ID,
		Kind:       ContentKindVideo,
		CreatedAt:  raw.CreatedAt,
		UpdatedAt:  raw.UpdateAt,
		MemberTier: raw.MemberTier(),
		Title:      raw.Title,
		Standfirst: raw.LongLeadCN,
		CoverURL:   raw.PostURL,
		CcID:       raw.CcID,
		Byline:     raw.BylineDescCN + " " + raw.BylineCN,
	}
}
