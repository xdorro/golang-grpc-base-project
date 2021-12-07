package event

import (
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
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
		evt.log.Error("evt.server.Run()", zap.Error(err))
		return err
	}

	return nil
}
