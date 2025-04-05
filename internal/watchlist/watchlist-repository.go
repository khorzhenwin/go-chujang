package watchlist

import "sync"

type Storage interface {
	Create(ticker Ticker) error
	GetAll() []Ticker
	GetByID(id string) (*Ticker, error)
	Update(id string, t Ticker) error
	Delete(id string) error
}

type InMemoryStore struct {
	mu      sync.RWMutex
	tickers map[string]Ticker
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{tickers: make(map[string]Ticker)}
}
