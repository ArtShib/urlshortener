package requestdeletion

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ArtShib/urlshortener/internal/model"
)

// URLService описывает интерфейс запросов на удаление.
type URLService interface {
	DeleteBatch(ctx context.Context, batch model.URLUserRequestArray) error
}

// DeletePool структура WorkerPool
type DeletePool struct {
	logger         *slog.Logger
	wg             sync.WaitGroup
	mu             sync.Mutex
	cancel         context.CancelFunc
	deleteRequests chan model.DeleteRequest
	inputCh        chan model.URLUserRequest
	buffer         model.URLUserRequestArray
	activeWorkers  int32
	URLService     URLService
	config         *model.WorkerPoolDelete
}

// NewWorkerPool конструктор WorkerPool
func NewWorkerPool(svc URLService, log *slog.Logger, cfg *model.WorkerPoolDelete) *DeletePool {
	return &DeletePool{
		logger:         log,
		deleteRequests: make(chan model.DeleteRequest, cfg.CountWorkers),
		inputCh:        make(chan model.URLUserRequest, cfg.InputChainSize),
		buffer:         make(model.URLUserRequestArray, 0, cfg.BufferSize),
		URLService:     svc,
		config:         cfg,
	}
}

// Start заускает WorkerPool
func (p *DeletePool) Start(ctx context.Context) {
	p.logger.Info("Sart DeletePool")
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	p.wg.Add(1)
	go p.scaleWorkers(ctx)

	p.wg.Add(1)
	go p.batcher(ctx)

}

func (p *DeletePool) scaleWorkers(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			queueLen := len(p.deleteRequests)
			activeWorkers := atomic.LoadInt32(&p.activeWorkers)

			if queueLen > 0 && activeWorkers < p.config.CountWorkers {

				atomic.AddInt32(&p.activeWorkers, 1)
				p.wg.Add(1)
				go p.worker(ctx, int(activeWorkers)+1, p.deleteRequests)
				p.logger.Debug("Add DeleteWorker",
					"ID", atomic.LoadInt32(&p.activeWorkers))
			}
		}
	}
}

// Stop останавливает WorkerPool
func (p *DeletePool) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
	p.wg.Wait()
}

func (p *DeletePool) worker(ctx context.Context, id int, requests <-chan model.DeleteRequest) {
	defer func() {
		p.wg.Done()
		atomic.AddInt32(&p.activeWorkers, -1)

	}()
	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-requests:
			if !ok {
				return
			}
			for _, uuid := range req.UUIDs {
				select {
				case p.inputCh <- model.URLUserRequest{
					UUID:   uuid,
					UserID: req.UserID,
				}:
				default:
					p.logger.Error("Input queue full, dropping request",
						"worker_id", id,
						"uuid", uuid)
				}
			}
		}
	}
}

// AddRequest добавляет запрос на удаление в очередь обработки.
func (p *DeletePool) AddRequest(req model.DeleteRequest) {
	p.deleteRequests <- req
}

func (p *DeletePool) batcher(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.flushBuffer(ctx)
			return

		case task := <-p.inputCh:
			p.addToBuffer(ctx, task)

		case <-ticker.C:
			p.flushIfNeeded(ctx)
		}
	}
}

func (p *DeletePool) flushBuffer(ctx context.Context) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.buffer) > 0 {
		batch := make(model.URLUserRequestArray, len(p.buffer))
		copy(batch, p.buffer)
		p.buffer = p.buffer[:0]

		p.logger.Info("Flushing buffer on shutdown",
			"batch_size", len(batch),
		)
		p.processBatch(ctx, batch)
	}
}

func (p *DeletePool) addToBuffer(ctx context.Context, task model.URLUserRequest) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.buffer = append(p.buffer, task)

	if len(p.buffer) >= p.config.BatchSize {
		go p.flushBufferAsync(ctx, p.getBufferCopy())
		p.buffer = p.buffer[:0]
	}
}

func (p *DeletePool) flushIfNeeded(ctx context.Context) {
	p.mu.Lock()
	if len(p.buffer) == 0 {
		p.mu.Unlock()
		return
	}

	batch := p.getBufferCopy()
	p.buffer = p.buffer[:0]
	p.mu.Unlock()

	p.processBatch(ctx, batch)
}

func (p *DeletePool) getBufferCopy() model.URLUserRequestArray {
	batch := make(model.URLUserRequestArray, len(p.buffer))
	copy(batch, p.buffer)
	return batch
}

func (p *DeletePool) flushBufferAsync(ctx context.Context, batch model.URLUserRequestArray) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	p.processBatch(ctx, batch)
}

func (p *DeletePool) processBatch(ctx context.Context, batch model.URLUserRequestArray) {
	if err := p.URLService.DeleteBatch(ctx, batch); err != nil {
		p.logger.Error("Batch processing failed",
			"error", err,
			"batch_size", len(batch))
	} else {
		p.logger.Info("Batch processed successfully",
			"batch_size", len(batch))
	}
}
