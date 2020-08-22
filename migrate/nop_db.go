package migrate

import (
	"context"
	"github.com/go-pg/pg/v9"
	"io"
	"strings"

	"github.com/go-pg/pg/v9/orm"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type nopDB struct {
	logger      Logger
	ignoreTable string
}

func (n *nopDB) Begin() (*pg.Tx, error) {
	return nil, nil
}

func (n *nopDB) ModelContext(c context.Context, model ...interface{}) *orm.Query {
	return nil
}

func (n *nopDB) ExecContext(c context.Context, query interface{}, params ...interface{}) (orm.Result, error) {
	return nil, nil
}

func (n *nopDB) ExecOneContext(c context.Context, query interface{}, params ...interface{}) (orm.Result, error) {
	return nil, nil
}

func (n *nopDB) QueryContext(c context.Context, model interface{}, query interface{}, params ...interface{}) (orm.Result, error) {
	return nil, nil
}

func (n *nopDB) QueryOneContext(c context.Context, model interface{}, query interface{}, params ...interface{}) (orm.Result, error) {
	return nil, nil
}

// Model ...
func (n *nopDB) Model(model ...interface{}) *orm.Query {
	return orm.NewQuery(n, model...)
}

// Select ...
func (n *nopDB) Select(model interface{}) error {
	return orm.Select(n, model)
}

// Insert ...
func (n *nopDB) Insert(model ...interface{}) error {
	return orm.Insert(n, model...)
}

// Delete ...
func (n *nopDB) Delete(model interface{}) error {
	return orm.Delete(n, model)
}

// ForceDelete ...
func (n *nopDB) ForceDelete(model interface{}) error {
	return orm.ForceDelete(n, model)
}

// Update ...
func (n *nopDB) Update(model interface{}) error {
	return orm.Update(n, model)
}

// Exec ...
func (n *nopDB) Exec(query interface{}, params ...interface{}) (orm.Result, error) {
	n.logQuery(query, params...)
	return nil, nil
}

// ExecOne ...
func (n *nopDB) ExecOne(query interface{}, params ...interface{}) (orm.Result, error) {
	n.logQuery(query, params...)
	return nil, nil
}

// Query ...
func (n *nopDB) Query(coll, query interface{}, params ...interface{}) (orm.Result, error) {
	n.logQuery(query, params...)
	return nil, nil
}

// QueryOne ...
func (n *nopDB) QueryOne(model, query interface{}, params ...interface{}) (orm.Result, error) {
	n.logQuery(query, params...)
	return nil, nil
}

// CopyFrom ...
func (n *nopDB) CopyFrom(r io.Reader, query interface{}, params ...interface{}) (orm.Result, error) {
	n.logQuery(query, params...)
	return nil, nil
}

// CopyTo ...
func (n *nopDB) CopyTo(w io.Writer, query interface{}, params ...interface{}) (orm.Result, error) {
	n.logQuery(query, params...)
	return nil, nil
}

// Formatter ...
func (n *nopDB) Formatter() orm.QueryFormatter {
	return n.Formatter()
}

// Context ...
func (*nopDB) Context() context.Context {
	return context.Background()
}

// FormatQuery ...
func (n *nopDB) FormatQuery(dst []byte, query string, params ...interface{}) []byte {
	return n.FormatQuery(dst, query, params...)
}

func (n *nopDB) logQuery(query interface{}, params ...interface{}) {
	var lastQuery string

	switch qry := query.(type) {
	case orm.QueryAppender:
		var dst []byte
		if b, err := qry.AppendQuery(n.Formatter(), dst); err == nil {
			lastQuery = string(b)
		}
	case string:
		var dst []byte
		lastQuery = string(n.FormatQuery(dst, qry, params...))
	}

	if n.logger != nil && !strings.Contains(lastQuery, n.ignoreTable) {
		n.logger.Printf("=== Unapplied SQL query ===\n%s\n\n", lastQuery)
	}
}
