package web

import (
	"log"
	"net/http"

	"github.com/ghhernandes/scheduler/storage"
	v1 "github.com/ghhernandes/scheduler/web/api/v1"
)

type Options struct {
	ListenAddress string

	Storage storage.Storage
}

type Handler struct {
	log     *log.Logger
	mux     *http.ServeMux
	options *Options
	apiV1   *v1.API
	quitCh  chan struct{}
	storage storage.Storage
}

func New(log *log.Logger, o *Options) *Handler {
	mux := http.NewServeMux()

	return &Handler{
		log:     log,
		mux:     mux,
		options: o,
		apiV1:   v1.New(log, o.Storage),
		quitCh:  make(chan struct{}),

		storage: o.Storage,
	}
}

func (h *Handler) Listen() <-chan struct{} {
	h.log.Print("starting server...")

	h.apiV1.Register(h.mux)

	go func() {
		http.ListenAndServe(h.options.ListenAddress, h.mux)
		close(h.quitCh)
	}()

	return h.quitCh
}
