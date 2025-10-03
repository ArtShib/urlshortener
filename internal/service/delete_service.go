package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ArtShib/urlshortener/internal/model"
)

type ConfigDelete struct {
	BatchSize    int
	BatchTimeout time.Duration
	WorkersCount int
}

func DefaultConfig() ConfigDelete {
	return ConfigDelete{
		BatchSize:    10,
		BatchTimeout: time.Second * 1,
		WorkersCount: 3,
	}
}

type DeleteService struct {
	repo    URLRepository
	config  ConfigDelete
	inputCh chan model.URLUserRequest
	buffer  model.URLUserRequestArray

	wg     sync.WaitGroup
	mu     sync.Mutex
	cancel context.CancelFunc
}

func NewDeleteService(config ConfigDelete) *DeleteService {
	return &DeleteService{
		config:  config,
		inputCh: make(chan model.URLUserRequest, config.BatchSize*2),
		buffer:  make(model.URLUserRequestArray, 0, config.BatchSize),
	}
}

func (s *DeleteService) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	for i := 0; i < s.config.WorkersCount; i++ {
		s.wg.Add(1)
		go s.worker(ctx, i)
	}
	s.wg.Add(1)
	go s.batcher(ctx)
}

func (s *DeleteService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}

func (s *DeleteService) AddQueueDelete(deleteRequest model.DeleteRequest) {
	if len(deleteRequest.UUIDs) == 0 {
		return
	}
	for _, uuid := range deleteRequest.UUIDs {
		s.QueueDelete(uuid, deleteRequest.UserID)
	}
}

func (s *DeleteService) QueueDelete(uuid string, userID string) {
	task := model.URLUserRequest{
		UUID:   uuid,
		UserID: userID,
	}
	select {
	case s.inputCh <- task:
	default:
		log.Printf("warning: delete queue is full, task dropped")
	}
}

func (s *DeleteService) batcher(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(s.config.BatchTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.flushBuffer()
			return

		case task := <-s.inputCh:
			s.mu.Lock()
			s.buffer = append(s.buffer, task)

			if len(s.buffer) >= s.config.BatchSize {
				s.flushAndUnlock()
			} else {
				s.mu.Unlock()
			}

		case <-ticker.C:
			s.mu.Lock()
			if len(s.buffer) > 0 {
				s.flushAndUnlock()
			} else {
				s.mu.Unlock()
			}
		}
	}
}

func (s *DeleteService) flushBuffer() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.buffer) > 0 {
		s.distributeTasks(s.buffer)
		s.buffer = s.buffer[:0]
	}
}

func (s *DeleteService) flushAndUnlock() {
	tasks := make(model.URLUserRequestArray, len(s.buffer))
	copy(tasks, s.buffer)
	s.buffer = s.buffer[:0]
	s.mu.Unlock()

	s.distributeTasks(tasks)
}

func (s *DeleteService) distributeTasks(tasks model.URLUserRequestArray) {
	for _, task := range tasks {
		select {
		case s.inputCh <- task:
		default:
			log.Printf("warning: failed to redistribute task")
		}
	}
}

func (s *DeleteService) worker(ctx context.Context, id int) {
	defer s.wg.Done()

	log.Printf("delete worker %d started", id)

	batch := make(model.URLUserRequestArray, 0, s.config.BatchSize)
	deleteCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			if len(batch) > 0 {
				s.repo.DeleteBatch(deleteCtx, batch)
			}
			log.Printf("delete worker %d stopped", id)
			return

		case task := <-s.inputCh:
			batch = append(batch, task)
			if len(batch) >= s.config.BatchSize {
				s.repo.DeleteBatch(deleteCtx, batch)
				batch = batch[:0]
			}
		}
	}
}
