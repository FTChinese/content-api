package models

import (
	"strings"
)

const imageBaseURL = "http://i.ftimg.net/"

// GalleryItem is a image url with its caption
type GalleryItem struct {
	// Image URL should be prepended with http://i.ftimg.net/
	ImageURL string `json:"imageUrl" db:"image_url"`
	Caption  string `json:"caption" db:"caption"`
}

// Gallery is a piece of photo news
type Gallery struct {
	Teaser
	Body  string        `json:"body" db:"body"`
	Items []GalleryItem `json:"items"`
}

// RawGallery is the scan target of DB.
type RawGallery struct {
	RawContentBase
	Body string `db:"body"`
}

func (r RawGallery) Build() Gallery {
	return Gallery{
		Teaser: r.Teaser(),
		Body:   strings.TrimSpace(r.Body),
		Items:  []GalleryItem{},
	}
}
