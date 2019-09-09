package models

// GalleryID is the photonnewsid column of photonews
type GalleryID int

// GalleryItem is a image url with its caption
type GalleryItem struct {
	// Image URL should be prepended with http://i.ftimg.net/
	ImageURL string `json:"imageUrl"`
	Caption  string `json:"caption"`
}

// Gallery is a piece of photo news
type Gallery struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Standfirst string `json:"standfirst"`
	Body       string `json:"body"`
	// Cover URL from db: photonews/2017/10/59eeb427456bb1.69227275.jpg
	// should be normalized to:
	// http://i.ftimg.net/photonews/2017/10/59eeb427456bb1.69227275.jpg
	CoverURL  string        `json:"coverUrl"`
	Items     []GalleryItem `json:"items"`
	Tags      []string      `json:"tags"`
	UpdatedAt string        `json:"updatedAt"`
}
