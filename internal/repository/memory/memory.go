package memory

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/ArtShib/urlshortener/internal/model"
)

type MemoryRepository struct{
	listURLs map[string] *model.URL
	mu sync.RWMutex
}

func NewMemoryRepository(fileName string) (*MemoryRepository, error){
	repo := &MemoryRepository{
		listURLs: make(map[string]*model.URL),
	}
	if err := repo.LoadingRepository(fileName); err != nil {
		return repo, err
	}
	return repo, nil
}

func (r *MemoryRepository) Store(url *model.URL) error{
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.listURLs[url.UUID]
	if ok {
		return errors.New("link already exists")
	}
	r.listURLs[url.UUID] = url
	return nil
}

func (r *MemoryRepository) FindByShortCode(uuid string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	url, ok := r.listURLs[uuid]
	
	if !ok {
		return nil, errors.New("longUrl is not found")
	}	
	
	return url, nil
}

func (r *MemoryRepository) LoadingRepository(fileName string) error {
	
	info, err := os.Stat(fileName) 
	if os.IsNotExist(err) {
		return err //stat test3.json: no such file or directory
	}
	if info.Size() == 0{
		return errors.New("file is empty")
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	var urls []model.URL
	if err := json.Unmarshal(data, &urls); err != nil {
		return err
	}
	for _, url := range urls {
		r.Store(&url)
	} 
	return nil
}
func (r *MemoryRepository) SavingRepository(fileName string) error {

	if len(r.listURLs) == 0 {
		return errors.New("listURLs is empty")
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE| os.O_TRUNC, 0777)
	defer func() error{
		if err := file.Close(); err != nil {
			return err
		}
		return nil
	}()

	if err != nil {
		return nil
	}
	urls := make([]*model.URL, 0, len(r.listURLs))
	for _, v := range r.listURLs {
		urls = append(urls, v)
	}

	data, err := json.Marshal(urls)	
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}
