package watchlist

import "gorm.io/gorm"

type Ticker struct {
	gorm.Model
	Symbol string `json:"symbol"` // e.g., AAPL
	Notes  string `json:"notes,omitempty"`
}
