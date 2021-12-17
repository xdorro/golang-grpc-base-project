package event

import (
	"strings"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"

	"github.com/xdorro/golang-grpc-base-project/pkg/client"
)

// Event is an event that can be published to an event bus.
type Event struct {
	client *client.Client

	event     *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
}

// NewEvent creates a new event.
func NewEvent(client *client.Client) *Event {
	redisURL := strings.Trim(viper.GetString("REDIS_URL"), " ")
	rdb := asynq.RedisClientOpt{Addr: redisURL}

	evt := &Event{
		client: client,
		event:  asynq.NewClient(rdb),
	}

	go func() {
		if err := evt.eventWorker(rdb); err != nil {
			return
		}
	}()

	go func() {
		if err := evt.eventScheduler(rdb); err != nil {
			return
		}
	}()

	return evt
}

// Close closes the event.
func (evt *Event) Close() error {
	evt.server.Shutdown()
	evt.scheduler.Shutdown()

	return evt.event.Close()
}
