package memory

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/ArtShib/urlshortener/internal/model"
)

type MemoryRepository struct {
	listURLs map[string]*model.URL
	mu       sync.RWMutex
	fileName string
}

func NewMemoryRepository(ctx context.Context, fileName string) (*MemoryRepository, error) {

	repo := &MemoryRepository{
		listURLs: make(map[string]*model.URL),
		fileName: fileName,
	}
	if err := repo.LoadingRepository(ctx); err != nil {
		return repo, err
	}
	return repo, nil
}

func (r *MemoryRepository) Save(ctx context.Context, url *model.URL) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.listURLs[url.UUID]
	if ok {
		return nil, model.ErrURLConflict
	}
	r.listURLs[url.UUID] = url
	return url, nil
}

func (r *MemoryRepository) Get(ctx context.Context, uuid string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.listURLs[uuid]

	if !ok {
		return nil, errors.New("longUrl is not found")
	}

	return url, nil
}

func (r *MemoryRepository) LoadingRepository(ctx context.Context) error {

	info, err := os.Stat(r.fileName)
	if os.IsNotExist(err) {
		return err
	}
	if info.Size() == 0 {
		return errors.New("file is empty")
	}
	data, err := os.ReadFile(r.fileName)
	if err != nil {
		return err
	}

	urls, err := r.unmarshalURL(data)
	if err != nil {
		return err
	}

	r.loadData(ctx, urls)

	return nil
}

func (r *MemoryRepository) unmarshalURL(data []byte) ([]*model.URL, error) {
	var urls []*model.URL
	if err := json.Unmarshal(data, &urls); err != nil {
		return nil, err
	}
	return urls, nil
}

func (r *MemoryRepository) loadData(ctx context.Context, urls []*model.URL) {
	for _, url := range urls {
		r.Save(ctx, url)
	}
}

func (r *MemoryRepository) Close() error {

	if len(r.listURLs) == 0 {
		return errors.New("listURLs is empty")
	}

	file, err := os.OpenFile(r.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() error {
		if err := file.Close(); err != nil {
			return err
		}
		return nil
	}()

	urls := make([]*model.URL, 0, len(r.listURLs))
	for _, v := range r.listURLs {
		urls = append(urls, v)
	}

	data, err := json.Marshal(urls)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (r *MemoryRepository) Ping(ctx context.Context) error {
	_, err := os.Stat(r.fileName)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

func (r *MemoryRepository) GetBatch(ctx context.Context, userId string) (model.URLUserBatch, error) {
	return nil, nil
}
