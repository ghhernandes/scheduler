package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ghhernandes/scheduler/storage/etcd"
	"github.com/ghhernandes/scheduler/watcher"
)

func main() {
	logger := log.New(os.Stdout, "watcher: ", log.LstdFlags)

	storage, err := etcd.New(logger, etcd.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	defer storage.Close()

	watcher := watcher.New(logger, watcher.Config{
		HandleTimeout: 10 * time.Second,

		Locker:     nil,
		Watcher:    nil,
		Dispatcher: nil,
		Storage:    storage,
	})

	watcher.Watch(context.Background())

	<-watcher.Done()
}
