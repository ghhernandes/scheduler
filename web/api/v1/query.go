package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ghhernandes/scheduler"
	"github.com/ghhernandes/scheduler/storage"
)

type ScheduleQuery struct {
	log       *log.Logger
	queryable storage.Queryable
}

func NewScheduleQuery(log *log.Logger, query storage.Queryable) *ScheduleQuery {
	return &ScheduleQuery{log: log, queryable: query}
}

func (s ScheduleQuery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.query(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s ScheduleQuery) query(w http.ResponseWriter, r *http.Request) {
	var results []ScheduleResult

	if len(r.URL.Query()) < 1 {
		http.Error(w, "no matchers specified", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	scheduleSet := s.queryable.Querier(ctx).Select(ctx, queryMatchers(r.URL)...)

	for scheduleSet.Next() {
		results = append(results, ScheduleResult{
			ValueResult: ValueResult{
				Date:    scheduleSet.At().Value.Date,
				Webhook: scheduleSet.At().Value.Webhook,
			},
			Ref:      scheduleSet.At().Ref.String(),
			Versions: nil,
		})
	}
	data, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "error on query marshalling", http.StatusInternalServerError)
	}

	w.Write(data)
}

func queryMatchers(url *url.URL) []scheduler.Matcher {
	var matchers []scheduler.Matcher

	params := url.Query()

	if prefix, ok := params["prefix"]; ok {
		matchers = append(matchers, scheduler.Matcher{
			Type:  scheduler.MatchPrefix,
			Value: prefix[0],
		})
	}

	if ref, ok := params["ref"]; ok {
		matchers = append(matchers, scheduler.Matcher{
			Type:  scheduler.MatchRef,
			Value: ref[0],
		})

	}

	return matchers
}
