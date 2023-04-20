package pkg

import (
	"strings"
)

const imageBaseURL = "https://i.ftimg.net/"

// GalleryItem is a image url with its caption
type GalleryItem struct {
	// Image URL should be prepended with http://i.ftimg.net/
	ImageURL string `json:"imageUrl" db:"image_url"`
	Caption  string `json:"caption" db:"caption"`
}

const StmtGalleryItem = `
SELECT TRIM(BOTH FROM pic_url) AS image_url,
    TRIM(BOTH FROM pbody) AS caption
FROM cmstmp01.photonews_picture
WHERE photonewsid = ?
ORDER BY orders`

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

const StmtGallery = `
SELECT photonewsid AS id,
    add_times AS created_at,
    accessright AS access_right,
    cn_title AS title_cn,
    TRIM(BOTH FROM shortlead) AS long_lead_cn,
    leadbody AS body,
    cover AS cover_url
FROM cmstmp01.photonews
WHERE photonewsid = ?
LIMIT 1`
