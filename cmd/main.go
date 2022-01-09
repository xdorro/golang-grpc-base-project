package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/internal/config"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

const (
	defaultShutdownTimeout = 10 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.NewLogger()
	config.NewConfig(log)

	// create new server
	srv, err := initializeServer(ctx, log)
	if err != nil {
		log.Fatal("initializeServer()", zap.Error(err))
	}

	// wait for termination signal and register database & http server clean-up operations
	wait := gracefulShutdown(ctx, log, defaultShutdownTimeout, map[string]operation{
		"server": func(ctx context.Context) error {
			return srv.Close()
		},
		"logger": func(ctx context.Context) error {
			return log.Sync()
		},
	})

	<-wait
}

// operation is a cleanup function on shutting down
type operation func(ctx context.Context) error

func gracefulShutdown(
	ctx context.Context, log *zap.Logger, timeout time.Duration, ops map[string]operation,
) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscall that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Info("Shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Fatal(fmt.Sprintf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds()))
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Info(fmt.Sprintf("cleaning up: %s", innerKey))
				if err := innerOp(ctx); err != nil {
					log.Info(fmt.Sprintf("%s: clean up failed: %s", innerKey, err.Error()))
					return
				}

				log.Info(fmt.Sprintf("%s was shutdown gracefully", innerKey))
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}
