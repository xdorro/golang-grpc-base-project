package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kucow/golang-grpc-base/internal/common"
)

// operation is a cleanup function on shutting down
type operation func(ctx context.Context) error

func gracefulShutdown(opts *common.Option, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscall that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		opts.Log.Info("Shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			opts.Log.Fatal(fmt.Sprintf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds()))
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

				opts.Log.Info(fmt.Sprintf("cleaning up: %s", innerKey))
				if err := innerOp(opts.Ctx); err != nil {
					opts.Log.Info(fmt.Sprintf("%s: clean up failed: %s", innerKey, err.Error()))
					return
				}

				opts.Log.Info(fmt.Sprintf("%s was shutdown gracefully", innerKey))
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}
