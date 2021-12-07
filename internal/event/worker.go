package event

import (
	"github.com/hibiken/asynq"
)

func (evt *Event) eventWorker(rdb asynq.RedisConnOpt) error {
	evt.server = asynq.NewServer(
		rdb,
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	// mux.HandleFunc("email:welcome", sendWelcomeEmail)
	// mux.HandleFunc("email:reminder", sendReminderEmail)

	if err := evt.server.Run(mux); err != nil {
		return err
	}

	return nil
}
