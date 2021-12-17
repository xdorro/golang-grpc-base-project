package event

import (
	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
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
		logger.Error("evt.server.Run()", zap.Error(err))
		return err
	}

	return nil
}
