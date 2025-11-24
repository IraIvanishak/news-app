package unsplash

// SearchResponse represents the Unsplash API search response
type SearchResponse struct {
	Total      int           `json:"total"`
	TotalPages int           `json:"total_pages"`
	Results    []PhotoResult `json:"results"`
}

// PhotoResult represents a single photo in search results
type PhotoResult struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	AltDesc     string    `json:"alt_description"`
	URLs        PhotoURLs `json:"urls"`
	Links       Links     `json:"links"`
	User        User      `json:"user"`
}

// PhotoURLs contains different sizes of photo URLs
type PhotoURLs struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

// Links contains various links related to the photo
type Links struct {
	DownloadLocation string `json:"download_location"`
}

// User represents the photographer
type User struct {
	Name  string    `json:"name"`
	Links UserLinks `json:"links"`
}

// UserLinks contains links to the photographer's profile
type UserLinks struct {
	HTML string `json:"html"`
}
