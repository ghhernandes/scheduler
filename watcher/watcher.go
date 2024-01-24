package watcher

import (
	"context"
	"log"
	"time"

	"github.com/ghhernandes/scheduler"
	"github.com/ghhernandes/scheduler/dispatcher"
	"github.com/ghhernandes/scheduler/storage"
)

type Handler struct {
	log *log.Logger

	chQuit chan struct{}

	watcher    Watcher
	dispatcher dispatcher.Dispatcher
	storage    storage.Storage
}

func New(log *log.Logger, w Watcher, d dispatcher.Dispatcher) *Handler {
	return &Handler{
		log:     log,
		watcher: w,
		chQuit:  make(chan struct{}),
	}
}

func (h *Handler) Watch(ctx context.Context) {
	defer close(h.chQuit)
	chWatch := h.watcher.Watch(ctx)

	for {
		select {
		case <-ctx.Done():
			h.log.Println("watcher context done.")
			return
		case schedule := <-chWatch:
			h.handle(ctx, schedule)
		}
	}
}

func (h *Handler) Done() <-chan struct{} {
	return h.chQuit
}

func (h *Handler) handle(ctx context.Context, s scheduler.Schedule) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)

	chDispatch := make(chan error)
	go func() {
		err := h.dispatcher.Dispatch(ctxTimeout, s)
		chDispatch <- err
	}()

	go func() {
		err := <-chDispatch
		cancel()

		if err != nil {
			h.log.Printf("dispatch: error: %s", err.Error())
			return
		}

		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := storage.ScheduledCommit(ctxTimeout, h.storage, s); err != nil {
			h.log.Printf("watcher: scheduled commit: error: %s", err.Error())
		}
	}()
}
