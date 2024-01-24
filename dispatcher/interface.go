package dispatcher

import (
	"context"

	"github.com/ghhernandes/scheduler"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, schedule scheduler.Schedule) error
}
