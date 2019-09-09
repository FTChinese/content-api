package models

// Channel contains article summary for each channel
type Channel struct {
	ID         string   `json:"id"`
	TitleCN    string   `json:"titleCn"`
	Standfirst string   `json:"standfirst"`
	Tags       []string `json:"tags"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
}
