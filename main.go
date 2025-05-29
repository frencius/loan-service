package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/infrastructure"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx, cancel := context.WithCancel(context.Background())
	app, err := application.SetupApp(ctx)
	if err != nil {
		log.Fatalf("failed to initiate app: %v", err)
	}
	defer app.Close()

	terminateChannel := make(chan os.Signal, 1)
	signal.Notify(terminateChannel,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		msgStr := fmt.Sprintf("system call: %+v", <-terminateChannel)
		log.Println(msgStr)
		cancel()
	}()

	log.Println("service started")
	hs := infrastructure.RunHTTPServer(app)

	defer hs.Close()
	<-ctx.Done()

	log.Println("service finished")
}
