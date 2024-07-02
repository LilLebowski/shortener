package models

type URLCreate struct {
	URL string `json:"URL"`
}

type URLResponse struct {
	Result string `json:"result"`
}

type URLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortURLs struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
