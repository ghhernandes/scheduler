package storage

import (
	"context"

	"github.com/ghhernandes/scheduler"
)

type Appender interface {
	Append(ctx context.Context, ref *scheduler.ScheduleRef, v scheduler.Value) (scheduler.ScheduleRef, error)

	Commit() error

	Rollback() error
}

type Appendable interface {
	Appender(ctx context.Context) Appender
}

type Querier interface {
	Select(ctx context.Context, matches ...scheduler.Matcher) ScheduleSet
}

type Queryable interface {
	Querier(ctx context.Context) Querier
}

type Storage interface {
	Appendable
	Queryable
}

type ScheduleSet interface {
	Next() bool
	At() scheduler.Schedule
	Versions() []scheduler.ValueVer
	Err() error
}
