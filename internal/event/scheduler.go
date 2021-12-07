package event

import (
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func (evt *Event) eventScheduler(rdb asynq.RedisConnOpt) error {
	evt.scheduler = asynq.NewScheduler(
		rdb,
		&asynq.SchedulerOpts{Location: time.Local},
	)

	// Run blocks and waits for os signal to terminate the program.
	if err := evt.scheduler.Run(); err != nil {
		evt.log.Error("evt.scheduler.Run()", zap.Error(err))
		return err
	}

	return nil
}
