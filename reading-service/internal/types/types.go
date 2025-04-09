package types

type URLResponse struct {
	OriginalURL string `json:"original_url"`
}

type URL struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	CreatedAt   int64  `json:"created_at"`
}
