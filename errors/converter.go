package errors

import (
	"context"
	"strings"

	"github.com/severgroup-tt/gopkg-errors"

	"github.com/go-pg/pg/v9"
)

const (
	msgInternal   = "Internal error"
	msgNotFound   = "Entity not found"
	msgConflict   = "Entity already exists"
	msgBadRequest = "Found too many entities"
)

const (
	pgDuplicateErr = "duplicate key value"
	pgCodeField    = 'C'
	pgStatusField  = 'S'
	pgMessageField = 'M'
)

// Converter ...
func Converter() errors.ErrorConverter {
	return func(ctx context.Context, err error) (*errors.Error, bool) {
		for {
			if err == pg.ErrNoRows {
				return errors.NotFound.ErrWithStack(ctx, msgNotFound).WithCause(err), true
			}
			if err == pg.ErrMultiRows {
				return errors.BadRequest.ErrWithStack(ctx, msgBadRequest).WithCause(err), true
			}

			if errTyped, ok := err.(pg.Error); ok {
				return convert(ctx, errTyped), true
			}

			errC, ok := err.(errors.Causer)
			if !ok {
				return nil, false
			}

			err = errC.Cause()
		}
	}
}

func convert(ctx context.Context, err pg.Error) *errors.Error {
	var result *errors.Error
	message := err.Field(pgMessageField)

	if strings.Contains(message, pgDuplicateErr) {
		result = errors.Conflict.ErrWithStack(ctx, msgConflict).WithCause(err.(error))
	} else {
		result = errors.Internal.ErrWithStack(ctx, msgInternal).WithCause(err.(error))
	}

	return result.WithLogKV(
		"code", err.Field(pgCodeField),
		"status", err.Field(pgStatusField),
		"message", message,
	)
}
