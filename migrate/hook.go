package migrate

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
	"strings"
)

type loggerHook struct {
	IgnoreTable string
}

func (h loggerHook) BeforeQuery(ctx context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (h loggerHook) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	query, err := event.FormattedQuery()
	if err == nil {
		return err
	}
	if !strings.Contains(query, h.IgnoreTable) {
		fmt.Printf("=== Applied SQL query ===\n%s\n\n", query)
	}
	return nil
}
