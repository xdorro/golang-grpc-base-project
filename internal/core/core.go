package core

import (
	"sync"

	"github.com/xdorro/golang-grpc-base-project/internal/server"
)

type App interface {
	// Options returns the current options
	Options() Options
	// Run the service
	Run() error
}

type app struct {
	opts Options

	once sync.Once
}

// NewApp creates and returns a new App based on the packages within.
func NewApp(opts ...Option) App {
	a := &app{
		opts: newOptions(opts...),
	}

	return a
}

func (a *app) Options() Options {
	return a.opts
}

func (a *app) Stop() error {
	return nil
}

func (a *app) Run() error {
	srv := server.NewServer(a.Options().Address)
	return srv.Run()
}
