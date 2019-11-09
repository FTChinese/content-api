package models

import (
	"strings"
)

const imageBaseURL = "http://i.ftimg.net/"

// RawGallery is the scan target of DB.
type RawGallery struct {
	RawContentBase
	Body string `db:"body"`
}

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

// NewGallery creates a Gallery instance from raw data retrieved from DB.
func NewGallery(raw *RawGallery) Gallery {
	return Gallery{
		Teaser: raw.Teaser(),
		Body:   strings.TrimSpace(raw.Body),
		Items:  []GalleryItem{},
	}
}
