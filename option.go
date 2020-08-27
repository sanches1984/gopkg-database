package database

import (
	"context"
	logger "github.com/sanches1984/gopkg-logger"
	"time"
)

type Option func(ctx context.Context) context.Context

func WithLongLogger(duration time.Duration) Option {
	logger.Info(logger.App, "DB Logger: log query over %v", duration)
	return func(ctx context.Context) context.Context {
		dbLogger := &dbLogger{duration: duration}
		dbc := FromContext(ctx)
		if dbc == nil {
			return ctx
		}
		dbc.Db().AddQueryHook(dbLogger)
		return context.WithValue(ctx, &dbLoggerKey, dbLogger)
	}
}
