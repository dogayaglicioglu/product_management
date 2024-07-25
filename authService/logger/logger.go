package logger

import (
	"fmt"
	"log"
	"sync"

	"github.com/fluent/fluent-logger-golang/fluent"
)

type LoggerInstance struct {
	Logger *fluent.Fluent
}

var (
	once   sync.Once
	LOGGER *LoggerInstance
)

func Logger() *LoggerInstance {
	once.Do(func() {
		logger, err := fluent.New(fluent.Config{
			FluentPort: 24224,
			FluentHost: "fluentd",
		})
		if err != nil {
			log.Fatalf("Failed to connect to Fluentd: %v", err)
		}

		LOGGER = &LoggerInstance{
			Logger: logger,
		}
	})
	return LOGGER
}
func LogInfo(logger *fluent.Fluent, msg string, traceId string) {
	err := logger.Post("auth-service.info", map[string]string{
		"message": msg,
		"traceId": traceId,
	})
	if err != nil {
		fmt.Printf("Failed to log info: %v\n", err)
	}
}

func LogError(logger *fluent.Fluent, msg string, traceId string, err error) {
	logErr := logger.Post("auth-service.error", map[string]string{
		"message": msg,
		"traceID": traceId,
		"error":   err.Error(),
	})
	if logErr != nil {
		fmt.Printf("Failed to log error: %v\n", logErr)
	}
}
