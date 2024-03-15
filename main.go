package main

import (
	"context"
	"time"

	"./logging"
	"./remote"
	"./service"
)

func main() {
	server := remote.Server()
	defer server.Close()

	service := service.SuperService{
		RequestDecorator: logging.Request,
		Logger:           logging.Logger{},
		RemoteUrl:        server.URL,
	}

	//start with an empty ctx and initial call
	ctx := context.Background()
	service.SomeServiceLogic(ctx, "service call from main.go")

	//simulate some time passing, then call again. traceId should be different
	println("\nSome time is passing...\n")
	time.Sleep(1 * time.Second)
	service.SomeServiceLogic(ctx, "another service call from main.go")

}
