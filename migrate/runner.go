package migrate

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/sanches1984/gopkg-database"
	"github.com/sanches1984/gopkg-errors"

	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
)

const (
	ignoreTableDefault = "gopg_migrations"
)

type Runner struct {
	DSN             string
	IgnoreTable     string
	LocalValuesPath string
	Logger          Logger
	DbConfKey       string

	dryRun      bool
	cleanScheme []string
}

func NewRunner(dsn string, options ...OptionFn) *Runner {
	ret := &Runner{
		DSN:         dsn,
		Logger:      log.New(os.Stdout, "", 0),
		IgnoreTable: ignoreTableDefault,
	}
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

func (r *Runner) Run(ctx context.Context) error {
	dbClient, err := r.initDB(r.DSN)
	if err != nil {
		log.Fatalf("Failed to connect database, error: %v", err)
	}
	defer dbClient.Close()

	conn := dbClient.Db()

	conn.AddQueryHook(loggerHook{})

	if len(r.cleanScheme) > 0 {
		for _, scheme := range r.cleanScheme {
			if err := r.cleanDatabase(conn, scheme); err != nil {
				return err
			}
		}
	}

	var oldVersion, newVersion int64
	if r.dryRun {
		nopDB := &nopDB{
			logger:      r.Logger,
			ignoreTable: r.IgnoreTable,
		}

		oldVersion, err = migrations.Version(conn)
		if err != nil {
			return errors.Internal.Err(ctx, "Failed to get current schema version").WithPayloadKV("err", err)
		}

		unappliedMigrations := make([]*migrations.Migration, 0)
		for _, m := range migrations.RegisteredMigrations() {
			if m.Version > oldVersion {
				unappliedMigrations = append(unappliedMigrations, m)
			}
		}

		newVersion = oldVersion
		if len(unappliedMigrations) > 0 {
			for _, m := range unappliedMigrations {
				err = m.Up(nopDB)
				if err != nil {
					return errors.Internal.Err(ctx, "Failed to run migrations").WithPayloadKV("err", err)
				}
				newVersion = m.Version
			}
		}
	} else {
		_, _, err = migrations.Run(conn, "init")
		if err != nil {
			pgErr, ok := err.(pg.Error)
			if !ok || pgErr.Field(67) != "42P07" {
				return err
			}
		}
		oldVersion, newVersion, err = migrations.Run(conn, flag.Args()...)
		if err != nil {
			return errors.Internal.Err(ctx, "Failed to run migrations").WithPayloadKV("err", err)
		}
	}

	if newVersion != oldVersion {
		r.Logger.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		r.Logger.Printf("version is %d\n", oldVersion)
	}

	return nil
}

// Connect connects to database
func (r *Runner) initDB(dsn string) (database.IClient, error) {
	r.Logger.Printf("connect to %s\n", dsn)
	opt, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	return database.Connect(opt.ApplicationName, opt), nil
}

// Clean database public scheme
func (r *Runner) cleanDatabase(db *pg.DB, scheme string) error {
	r.Logger.Printf("clean scheme %s\n", scheme)
	var res bool
	_, err := db.Query(&res, "DROP SCHEMA "+scheme+" CASCADE")
	if err != nil {
		return err
	}
	_, err = db.Query(&res, "CREATE SCHEMA "+scheme)
	if err != nil {
		return err
	}
	return nil
}
