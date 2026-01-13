package audit

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ArtShib/urlshortener/internal/model"
)

// EventService описывает интерфейс сохранения аудита (сервисный уровень)
type EventService interface {
	SendAuditRecord(ctx context.Context, record *model.Event) error
	Close() error
}

// WorkerPoolEvent структура WorkerPool Event
type WorkerPoolEvent struct {
	logger        *slog.Logger
	wg            sync.WaitGroup
	stopOnce      sync.Once
	cancel        context.CancelFunc
	eventCh       chan *model.Event
	activeWorkers atomic.Int32
	workerID      atomic.Int32
	EventService  EventService
	config        *model.WorkerPoolEvent
	stopped       atomic.Bool
}

// New конструктор WorkerPool Event
func New(svc EventService, log *slog.Logger, cfg *model.WorkerPoolEvent) *WorkerPoolEvent {
	return &WorkerPoolEvent{
		logger:       log,
		eventCh:      make(chan *model.Event, cfg.EventChainSize),
		EventService: svc,
		config:       cfg,
	}
}

// Start запуск WorkerPool Event
func (p *WorkerPoolEvent) Start(ctx context.Context) {
	const op = "WorkerPoolEvent.Start"
	log := p.logger.With(
		slog.String("op", op),
	)
	log.Debug("Starting EventPool")
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.wg.Add(1)
	go p.scaleWorkers(ctx)
	p.addWorker(ctx)
}

func (p *WorkerPoolEvent) scaleWorkers(ctx context.Context) {
	const op = "WorkerPoolEvent.scaleWorkers"
	log := p.logger.With(
		slog.String("op", op),
	)
	log.Debug("Start scaleWorkers")

	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		p.wg.Done()
		ticker.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if p.stopped.Load() {
				return
			}
			queueLen := len(p.eventCh)
			activeWorkers := p.activeWorkers.Load()

			if queueLen > 0 && activeWorkers < p.config.CountWorkers {
				p.addWorker(ctx)
			}
		}
	}
}

// Stop остановка WorkerPool Event
func (p *WorkerPoolEvent) Stop() {
	p.stopOnce.Do(func() {
		p.stopped.Store(true)
		const op = "WorkerPoolEvent.Stop"
		log := p.logger.With(
			slog.String("op", op),
		)
		log.Debug("Stopping WorkerPool")
		if p.cancel != nil {
			p.cancel()
		}

		p.wg.Wait()
		if p.EventService != nil {
			if err := p.EventService.Close(); err != nil {
				log.Error("EventService.Close", "error", err)
			}
		}
		log.Debug("All workers stopped")
	})
}

func (p *WorkerPoolEvent) worker(ctx context.Context, id int) {
	const op = "WorkerPoolEvent.worker"
	log := p.logger.With(
		slog.String("op", op),
		slog.Int("worker_id", id),
	)
	log.Debug("Start Worker")
	defer func() {
		p.wg.Done()
		p.activeWorkers.Add(-1)
	}()
	for {
		select {
		case <-ctx.Done():
			//return
			//case event, ok := <-p.eventCh:
			//	if !ok {
			//		return
			//	}
			//	p.processEvent(ctx, event, id, log)
			for {
				select {
				case event := <-p.eventCh:
					p.processEvent(context.Background(), event, id, log)
				default:
					return
				}
			}
		case event := <-p.eventCh:
			p.processEvent(ctx, event, id, log)
		}
	}
}

func (p *WorkerPoolEvent) processEvent(ctx context.Context, event *model.Event, id int, log *slog.Logger) {
	//pCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()
	if err := p.EventService.SendAuditRecord(ctx, event); err != nil {
		log.Error("Error adding audit record",
			"error", err,
			"worker_id", id)
	}
}
func (p *WorkerPoolEvent) addWorker(ctx context.Context) {
	const op = "WorkerPoolEvent.addWorker"
	log := p.logger.With(
		slog.String("op", op),
	)
	log.Debug("Start addWorker")

	if p.stopped.Load() {
		return
	}

	select {
	case <-ctx.Done():
		return
	default:
		p.wg.Add(1)
		p.activeWorkers.Add(1)
		p.workerID.Add(1)
		id := p.workerID.Load()
		//id := p.activeWorkers.Load()
		go p.worker(ctx, int(id))
		log.Debug("Worker added", "ID", id)
	}
}

// AddEventRecord метод добаления сообщения аудита в очередь
func (p *WorkerPoolEvent) AddEventRecord(event *model.Event) {
	const op = "WorkerPoolEvent.AddEventRecord"

	if p.stopped.Load() {
		p.logger.Debug("WorkerPool is stopped, dropping audit")
		return
	}
	select {
	case p.eventCh <- event:
		p.logger.Debug("Add EventRecord")
	default:
		p.logger.Error(op, "error", fmt.Errorf("event buffer full, dropping audit"))
	}
}
