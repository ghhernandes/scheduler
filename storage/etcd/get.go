package etcd

import (
	"context"
	"encoding/json"

	"github.com/ghhernandes/scheduler"
	"github.com/ghhernandes/scheduler/storage"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (etcd etcd) Querier(ctx context.Context) storage.Querier {
	return etcd
}

func (etcd etcd) Select(ctx context.Context, matches ...scheduler.Matcher) storage.ScheduleSet {
	for _, m := range matches {
		if m.Type == scheduler.MatchRef {
			r, err := etcd.client.Get(ctx, m.Value)
			return &scheduleSet{response: r, err: err}
		}
	}
	return nil
}

type scheduleSet struct {
	response *clientv3.GetResponse

	pos      int
	at       scheduler.Schedule
	versions []scheduler.ValueVer
	err      error
}

func (s *scheduleSet) Next() bool {
	var value scheduler.Value

	if s.pos >= len(s.response.Kvs) {
		return false
	}

	key := s.response.Kvs[s.pos].Key
	kvalue := s.response.Kvs[s.pos].Value

	json.Unmarshal(kvalue, &value)

	s.at = scheduler.Schedule{
		Ref:   scheduler.NewRefString(string(key)),
		Value: value,
	}

	s.pos++
	return true
}

func (s scheduleSet) At() scheduler.Schedule {
	return s.at
}

func (s scheduleSet) Versions() []scheduler.ValueVer {
	return s.versions
}

func (s scheduleSet) Err() error {
	return s.err
}
