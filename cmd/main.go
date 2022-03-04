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

	"github.com/xdorro/golang-grpc-base-project/config"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

// operation is a cleanup function on shutting down
type operation func(ctx context.Context) error

const (
	defaultShutdownTimeout = 10 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create logging and config
	log := logger.NewLogger()
	config.NewConfig(log)

	// init server
	srv := initializeServer(ctx, log)

	// wait for termination signal and register client & http server clean-up operations
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
			log.Panic(fmt.Sprintf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds()))
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for innerKey, innerOp := range ops {
			wg.Add(1)
			func() {
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
