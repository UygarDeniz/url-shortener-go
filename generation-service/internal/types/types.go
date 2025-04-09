package types

type URL struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
}

type URLRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
