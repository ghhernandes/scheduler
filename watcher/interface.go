package watcher

import (
	"context"

	"github.com/ghhernandes/scheduler"
)

type Watcher interface {
	Watch(ctx context.Context) <-chan scheduler.Schedule
}
