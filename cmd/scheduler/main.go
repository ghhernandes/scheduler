package main

import (
	"log"
	"os"
	"time"

	"github.com/ghhernandes/scheduler/storage/etcd"
	"github.com/ghhernandes/scheduler/web"
)

func main() {
	logger := log.New(os.Stdout, "scheduler: ", log.LstdFlags)

	storage, err := etcd.New(logger, etcd.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	defer storage.Close()

	options := &web.Options{
		ListenAddress: ":8080",
		Storage:       storage,
	}

	handler := web.New(logger, options)

	<-handler.Listen()
}
