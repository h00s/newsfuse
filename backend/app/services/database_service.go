// Package services holds data access, the scrapers and the summarizer.
package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uptrace/bun"
)

type DatabaseService struct {
	raptor.Service

	// Ctx is the context CRUD queries run on: the app's context, live while requests run and
	// drain, and cancelled at shutdown so a query that outlives its request stops instead of
	// holding the pool open. Work that should stop when the client disconnects takes a
	// context.Context parameter instead.
	Ctx context.Context
}

func (s *DatabaseService) Setup() error {
	s.Ctx = s.AppContext()
	return nil
}

func (s *DatabaseService) Conn() *bun.DB {
	return s.Database.Conn().(*bun.DB)
}

// HandleError maps a query error onto an errs value. A missing row is not an error at this
// level: lists come back empty, writes check RowsAffected, and single reads use
// HandleErrorNotFound.
func (s *DatabaseService) HandleError(err error) error {
	if err == nil || errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	// Shutdown cancelled the app context, or the client went away: not a database fault, so not
	// an error line. A deadline is a slow query, and falls through to the error below.
	if errors.Is(err, context.Canceled) {
		s.Log.Debug("Database call cancelled", "error", err)
		return errs.NewErrorServiceUnavailable("Database call cancelled")
	}
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return s.handlePostgresError(pgErr)
	}
	s.Log.Error("Unhandled database error", "error", err)
	return errs.NewErrorInternal("Database error")
}

// HandleErrorNotFound is HandleError for single-row reads, where no row means 404.
func (s *DatabaseService) HandleErrorNotFound(err error, message ...string) error {
	if errors.Is(err, sql.ErrNoRows) {
		msg := "Resource not found"
		if len(message) > 0 {
			msg = message[0]
		}
		return errs.NewErrorNotFound(msg)
	}
	return s.HandleError(err)
}

// HandleAffected finishes an Update or Delete: it maps a query error, and turns a write that
// matched no row into a 404. Neither statement reports sql.ErrNoRows, so without this check a
// write aimed at a missing row would answer 200.
func (s *DatabaseService) HandleAffected(res sql.Result, err error, notFound string) error {
	if err != nil {
		return s.HandleError(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return s.HandleError(err)
	}
	if n == 0 {
		return errs.NewErrorNotFound(notFound)
	}
	return nil
}

// handlePostgresError gives the client a status and at most the constraint name. The driver's
// message and detail echo submitted values and the schema, so they go to the log only.
func (s *DatabaseService) handlePostgresError(pgErr *pgconn.PgError) error {
	var attrs []any
	if pgErr.ConstraintName != "" {
		attrs = []any{"constraint", pgErr.ConstraintName}
	}

	var mapped *errs.Error
	switch pgErr.Code {
	case "23505": // unique_violation
		mapped = errs.NewErrorConflict("Unique constraint violation", attrs...)
	case "23503": // foreign_key_violation
		mapped = errs.NewErrorConflict("Foreign key violation", attrs...)
	case "23502": // not_null_violation
		mapped = errs.NewErrorUnprocessableEntity("Missing required value", attrs...)
	case "23514": // check_violation
		mapped = errs.NewErrorUnprocessableEntity("Check constraint violation", attrs...)
	case "22001": // string_data_right_truncation
		mapped = errs.NewErrorUnprocessableEntity("Value too long", attrs...)
	case "22003": // numeric_value_out_of_range
		mapped = errs.NewErrorUnprocessableEntity("Number out of range", attrs...)
	case "22P02": // invalid_text_representation
		mapped = errs.NewErrorUnprocessableEntity("Invalid value", attrs...)
	default:
		s.Log.Error("Unhandled database error", "code", pgErr.Code, "message", pgErr.Message)
		return errs.NewErrorInternal("Database error")
	}

	s.Log.Warn("Database constraint error",
		"code", pgErr.Code, "constraint", pgErr.ConstraintName,
		"message", pgErr.Message, "detail", pgErr.Detail)
	return mapped
}
