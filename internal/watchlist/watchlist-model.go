package watchlist

type Ticker struct {
	ID     string `json:"id"`     // UUID or string
	Symbol string `json:"symbol"` // e.g., AAPL
	Notes  string `json:"notes,omitempty"`
}
