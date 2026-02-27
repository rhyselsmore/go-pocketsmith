package pocketsmith

import "time"

// Institution represents the financial institution details.
type Institution struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	CurrencyCode   string    `json:"currency_code"`
	Colour         string    `json:"colour"`
	LogoURL        string    `json:"logo_url"`
	FaviconDataURI string    `json:"favicon_data_uri"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
