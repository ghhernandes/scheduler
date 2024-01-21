package web

import (
	"log"
	"net/http"

	"github.com/ghhernandes/scheduler/storage"
)

type API struct {
	log *log.Logger

	storage storage.Storage
}

func New(log *log.Logger, storage storage.Storage) *API {
	return &API{log: log, storage: storage}
}

func (api *API) Register(mux *http.ServeMux) {
	mux.Handle("/v1/schedule", NewScheduleCreate(api.log, api.storage))
	mux.Handle("/v1/schedule/keys", NewScheduleQuery(api.log, api.storage))
}
