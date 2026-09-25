package handler

type createLinkRequest struct {
	URL string `json:"url"`
}

type linkResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	Clicks      int    `json:"clicks"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
