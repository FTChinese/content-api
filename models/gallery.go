package models

import "strings"

const imageBaseURL = "http://i.ftimg.net/"

// GalleryItem is a image url with its caption
type GalleryItem struct {
	// Image URL should be prepended with http://i.ftimg.net/
	ImageURL string `json:"imageUrl" db:"imageUrl"`
	Caption  string `json:"caption" db:"caption"`
}

func (g *GalleryItem) Normalize() {
	g.ImageURL = strings.TrimSpace(g.ImageURL)
	g.Caption = strings.TrimSpace(g.Caption)

	g.ImageURL = imageBaseURL + g.ImageURL
}

// Gallery is a piece of photo news
type Gallery struct {
	ID         int    `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	Standfirst string `json:"standfirst" db:"standfirst"`
	Body       string `json:"body" db:"body"`
	// Cover URL from db: photonews/2017/10/59eeb427456bb1.69227275.jpg
	// should be normalized to:
	// http://i.ftimg.net/photonews/2017/10/59eeb427456bb1.69227275.jpg
	CoverURL  string        `json:"coverUrl" db:"coverUrl"`
	Items     []GalleryItem `json:"items"`
	Tag       string        `json:"-" db:"tag"`
	Tags      []string      `json:"tags"`
	UpdatedAt string        `json:"updatedAt" db:"updatedAt"`
}

func (g *Gallery) Normalize() {
	g.Standfirst = strings.TrimSpace(g.Standfirst)
	g.Body = strings.TrimSpace(g.Body)

	g.Tags = strings.Split(g.Tag, ",")
	if !strings.HasPrefix(g.CoverURL, "http") {
		g.CoverURL = imageBaseURL + g.CoverURL
	}
}
