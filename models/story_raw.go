package models

// RawStory is used to retrieve a story from db as is.
type RawStory struct {
	ID             string `json:"id" db:"story_id"`
	Bilingual      bool   `json:"bilingual"`
	TitleCN        string `json:"titleCn" db:"title_cn"`
	TitleEN        string `json:"titleEn" db:"title_en"`
	Standfirst     string `json:"standfirst" db:"standfirst"`
	CoverURL       string `json:"coverUrl" db:"coverUrl"`
	BylineDescCN   string `json:"bylineDescCn" db:"byline_desc_cn"`
	BylineDescEN   string `json:"bylineDescEn" db:"byline_desc_en"`
	BylineAuthorCN string `json:"bylineAuthorCn" db:"byline_author_cn"`
	BylineAuthorEN string `json:"bylineAuthorEn" db:"byline_author_en"`
	BylineStatusCN string `json:"bylineStatusCn" db:"byline_status_cn"`
	BylineStatusEN string `json:"bylineStatusEn" db:"byline_status_en"`
	Tags           string `json:"tags" db:"tags"`
	Genre          string `json:"genre" db:"genre"`
	Topic          string `json:"topic" db:"topic"`
	Industry       string `json:"industry" db:"industry"`
	Area           string `json:"area" db:"area"`
	CreatedAt      string `json:"createdAt" db:"created_at"`
	UpdatedAt      string `json:"updatedAt" db:"updated_at"`
	BodyCN         string `json:"bodyCn" db:"body_cn"`
	BodyEN         string `json:"bodyEn" db:"body_en"`
}
