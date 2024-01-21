package main

import (
	"log"
	"os"

	"github.com/ghhernandes/scheduler/web"
)

func main() {
	logger := log.New(os.Stdout, "scheduler: ", log.LstdFlags)

	options := &web.Options{
		ListenAddress: ":8080",
		Storage:       nil,
	}

	handler := web.New(logger, options)

	<-handler.Listen()
}
