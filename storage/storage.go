package storage

import (
	"context"

	"github.com/ghhernandes/scheduler"
)

func AppendAndCommit(ctx context.Context, storage Storage, ref *scheduler.ScheduleRef, v scheduler.Value) (scheduler.ScheduleRef, error) {
	a := storage.Appender(ctx)

	resultRef, err := a.Append(ctx, ref, v)
	if err != nil {
		return scheduler.ScheduleRef{}, err
	}

	if err = a.Commit(); err != nil {
		return scheduler.ScheduleRef{}, err
	}

	return resultRef, nil
}

func ScheduledCommit(ctx context.Context, storage Storage, s scheduler.Schedule) error {
	a := storage.Appender(ctx)

	err := a.Scheduled(ctx, s)
	if err != nil {
		return err
	}

	if err = a.Commit(); err != nil {
		return err
	}

	return nil
}
