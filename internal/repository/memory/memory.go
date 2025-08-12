package memory

import (
	"errors"
	"sync"
	
	"github.com/ArtShib/urlshortener/internal/model"
)

type MemoryRepository struct{
	listURLs map[string] *model.URL
	mu sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository{
	return &MemoryRepository{
		listURLs: make(map[string]*model.URL),
	}
}

func (r *MemoryRepository) Store(url *model.URL) error{
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.listURLs[url.ShortCode]
	if ok {
		return errors.New("link already exists")
	}
	r.listURLs[url.ShortCode] = url
	return nil
}

func (r *MemoryRepository) FindByShortCode(shortCode string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	url, ok := r.listURLs[shortCode]
	
	if !ok {
		return nil, errors.New("longUrl is not found")
	}	
	
	return url, nil
}
