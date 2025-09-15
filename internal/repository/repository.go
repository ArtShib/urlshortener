package repository

import (
	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/ArtShib/urlshortener/internal/repository/memory"
)

type URLRepository interface {
	Store(url *model.URL) error
	FindByShortCode(uuid string) (*model.URL, error)
	SavingRepository(fileName string) error 
}
func NewRepository(fileName string) (URLRepository, error) {
	repo, err := memory.NewMemoryRepository(fileName)
	return repo, err
}
