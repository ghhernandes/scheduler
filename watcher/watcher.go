package watcher

import (
	"context"
	"log"
	"time"

	"github.com/ghhernandes/scheduler"
	"github.com/ghhernandes/scheduler/dispatcher"
	"github.com/ghhernandes/scheduler/lock"
	"github.com/ghhernandes/scheduler/storage"
)

type Config struct {
	HandleTimeout time.Duration

	Locker     lock.Locker
	Watcher    Watcher
	Dispatcher dispatcher.Dispatcher
	Storage    storage.Storage
}

type Handler struct {
	log    *log.Logger
	config *Config

	chQuit chan struct{}

	locker     lock.Locker
	watcher    Watcher
	dispatcher dispatcher.Dispatcher
	storage    storage.Storage
}

func New(log *log.Logger, cfg Config) *Handler {
	return &Handler{
		log:    log,
		config: &cfg,
		chQuit: make(chan struct{}),

		locker:     cfg.Locker,
		watcher:    cfg.Watcher,
		dispatcher: cfg.Dispatcher,
		storage:    cfg.Storage,
	}
}

func (h *Handler) Watch(ctx context.Context) {
	defer close(h.chQuit)

	chWatch, err := h.watcher.Watch(ctx)
	if err != nil {
		h.log.Printf("watcher start: error: %s", err.Error())
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				h.log.Println("watcher context done.")
				return
			case schedule, ok := <-chWatch:
				if !ok {
					h.log.Println("watch channel closed.")
					return
				}
				go h.handle(ctx, schedule)
			}
		}
	}()
}

func (h *Handler) Done() <-chan struct{} {
	return h.chQuit
}

func (h *Handler) handle(ctx context.Context, s scheduler.Schedule) {
	if !h.acquireLock(s) {
		return
	}
	defer h.unlock(s.Ref)

	ctxHandle, cancelHandle := context.WithTimeout(ctx, h.config.HandleTimeout)
	defer cancelHandle()

	err := h.dispatcher.Dispatch(ctxHandle, s)
	if err != nil {
		h.log.Printf("dispatch: error: %s", err.Error())
		return
	}

	if err := storage.ScheduledCommit(ctxHandle, h.storage, s); err != nil {
		h.log.Printf("watcher: scheduled commit: error: %s", err.Error())
	}
}

func (h *Handler) acquireLock(s scheduler.Schedule) bool {
	err := h.locker.Lock(s.Ref)
	if err != nil {
		h.log.Printf("watcher: cannot lock ref %s", s.Ref)
		return false
	}
	return true
}

func (h *Handler) unlock(ref scheduler.ScheduleRef) {
	err := h.locker.Unlock(ref)
	if err != nil {
		h.log.Printf("watch: handler: unlock error: %s", err.Error())
	}
}
