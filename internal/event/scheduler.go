package event

import (
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

func (evt *Event) eventScheduler(rdb asynq.RedisConnOpt) error {
	evt.scheduler = asynq.NewScheduler(
		rdb,
		&asynq.SchedulerOpts{Location: time.Local},
	)

	// Run blocks and waits for os signal to terminate the program.
	if err := evt.scheduler.Run(); err != nil {
		logger.Error("evt.scheduler.Run()", zap.Error(err))
		return err
	}

	return nil
}
