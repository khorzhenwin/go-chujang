package watchlist

import "gorm.io/gorm"

type Storage interface {
	Create(ticker Ticker) error
	GetAll() []Ticker
	GetByID(id string) (*Ticker, error)
	Update(id string, t Ticker) error
	Delete(id string) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// AutoMigrate sets up the schema
func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&Ticker{})
}
