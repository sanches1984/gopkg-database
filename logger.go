package database

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/sanches1984/gopkg-logger"
	"time"
)

const queryStartTime = "StartTime"

type IDBLogger interface {
	pg.QueryHook
	SetDuration(duration time.Duration)
}

type dbLogger struct {
	duration time.Duration
}

func (d *dbLogger) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	if event.Stash == nil {
		event.Stash = make(map[interface{}]interface{})
	}
	event.Stash[queryStartTime] = time.Now()
	return ctx, nil
}

func (d *dbLogger) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	query, err := event.FormattedQuery()
	if err == nil {
		var duration time.Duration
		if event.Stash != nil {
			if v, ok := event.Stash[queryStartTime]; ok {
				duration = time.Now().Sub(v.(time.Time))
			}
		}
		logLevel := logger.LevelInfo
		if d.duration != 0 {
			if d.duration > duration {
				return nil
			}
			logLevel = logger.LevelWarning
		}
		txt := "query: " + query
		if duration != 0 {
			txt += fmt.Sprintf(" [%d ms]", duration.Nanoseconds()/1000000)
		}
		if event.Err != nil {
			txt += "\nerror: " + event.Err.Error()
		}
		if ctx == nil {
			logger.Log(logger.App, logLevel, txt)
		} else {
			logger.Log(ctx, logLevel, txt)
		}
	}
	return nil
}

func (d *dbLogger) SetDuration(duration time.Duration) {
	d.duration = duration
}
