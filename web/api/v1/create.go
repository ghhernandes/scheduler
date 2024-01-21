package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ghhernandes/scheduler/storage"
)

type ScheduleCreate struct {
	log *log.Logger

	storage storage.Storage
}

func NewScheduleCreate(log *log.Logger, storage storage.Storage) *ScheduleCreate {
	return &ScheduleCreate{log: log, storage: storage}
}

func (s ScheduleCreate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.post(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s ScheduleCreate) post(w http.ResponseWriter, r *http.Request) {
	var createRequest ScheduleCreateRequest

	err := json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ref, err := storage.AppendAndCommit(ctx, s.storage, nil, createRequest.ToValue())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	result := ScheduleCreateResult{
		Ref: ref.String(),
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}
