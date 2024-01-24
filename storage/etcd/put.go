package etcd

import (
	"context"
	"encoding/json"

	"github.com/ghhernandes/scheduler"
	"github.com/ghhernandes/scheduler/storage"
)

func (etcd etcd) Appender(ctx context.Context) storage.Appender {
	return etcd
}

func (etcd etcd) Append(ctx context.Context, ref *scheduler.ScheduleRef, v scheduler.Value) (scheduler.ScheduleRef, error) {
	var key string
	var newRef scheduler.ScheduleRef

	if ref == nil {
		newRef = scheduler.NewRef()
	} else {
		newRef = *ref
	}

	value, err := json.Marshal(v)
	if err != nil {
		return newRef, err
	}

	//TODO: error handling
	_, err = etcd.client.Put(ctx, key, string(value))
	return newRef, err
}

func (etcd etcd) Scheduled(ctx context.Context, s scheduler.Schedule) error {
	return nil
}

func (etcd etcd) Commit() error {
	return nil
}

func (etcd etcd) Rollback() error {
	return nil
}
