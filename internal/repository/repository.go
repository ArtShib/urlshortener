package repository

import (
	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/ArtShib/urlshortener/internal/repository/memory"
)

type URLRepository interface {
	Store(url *model.URL) error
	FindByShortCode(shortCode string) (*model.URL, error)
}
func NewRepository() URLRepository{
	return memory.NewMemoryRepository()
}
