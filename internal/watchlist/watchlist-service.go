package watchlist

type Service struct {
	store Storage
}

func NewService(store Storage) *Service {
	return &Service{store: store}
}

// Expose methods like AddTicker, GetTickers, UpdateTicker, DeleteTicker, etc.
