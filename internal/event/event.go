package event

import (
	"strings"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/pkg/client"
)

// Event is an event that can be published to an event bus.
type Event struct {
	log    *zap.Logger
	client *client.Client

	event     *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
}

// NewEvent creates a new event.
func NewEvent(log *zap.Logger, client *client.Client) *Event {
	evt := &Event{
		log:    log,
		client: client,
	}

	redisURL := strings.Trim(viper.GetString("REDIS_URL"), " ")
	rdb := asynq.RedisClientOpt{Addr: redisURL}

	evt.event = asynq.NewClient(rdb)

	go func() {
		if err := evt.eventWorker(rdb); err != nil {
			evt.log.Fatal("error starting event worker", zap.Error(err))
			return
		}
	}()

	go func() {
		if err := evt.eventScheduler(rdb); err != nil {
			evt.log.Fatal("error starting event scheduler", zap.Error(err))
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
