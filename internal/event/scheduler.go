package event

import (
	"time"

	"github.com/hibiken/asynq"
)

func (evt *Event) eventScheduler(rdb asynq.RedisConnOpt) error {
	evt.scheduler = asynq.NewScheduler(
		rdb,
		&asynq.SchedulerOpts{Location: time.Local},
	)

	// Run blocks and waits for os signal to terminate the program.
	if err := evt.scheduler.Run(); err != nil {
		return err
	}

	return nil
}
