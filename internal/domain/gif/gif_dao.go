package gif

// Gif model definition from Giphy
type Gif struct {
	Title  string `json:"title"`
	Images Images `json:"images"`
}

// Images type definition from Giphy
type Images struct {
	Original Original `json:"original"`
}

// Original image type definition from Giphy
type Original struct {
	URL string `json:"url"`
}
