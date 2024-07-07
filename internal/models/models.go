// Package models contain types for service
package models

// URLCreate - body for creation
type URLCreate struct {
	URL string `json:"URL"`
}

// URLResponse - response after creation
type URLResponse struct {
	Result string `json:"result"`
}

// URLs - batch creation request
type URLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// FullURLs - batch creation
type FullURLs struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// ShortURLs - response after batch creation
type ShortURLs struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserURL - full link info
type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
}
