package models

// TopNav is the first level of navigation data
type TopNav struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Title    string   `json:"title"`
	Children []SubNav `json:"children"`
}

// SubNav is the second level of navigation data
type SubNav struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}
