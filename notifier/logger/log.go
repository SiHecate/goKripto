package logger

import (
	"context"
	"time"

	zaploki "github.com/paul-milne/zap-loki"
	"go.uber.org/zap"
)

func InitLogger() (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	loki := zaploki.New(context.Background(), zaploki.Config{
		Url:          "http://loki:3100",
		BatchMaxSize: 1000,
		BatchMaxWait: 10 * time.Second,
		Labels:       map[string]string{"app": "Notifier"},
	})

	return loki.WithCreateLogger(zapConfig)
}
