package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/go-raptor/raptor/v4/errs"
)

func TestDatabaseQueriesStopWithTheApp(t *testing.T) {
	res, shutdown := testResources()
	s := &DatabaseService{}
	if err := s.Init(res); err != nil {
		t.Fatal(err)
	}
	if err := s.Setup(); err != nil {
		t.Fatal(err)
	}

	shutdown()

	if s.Ctx.Err() == nil {
		t.Error("DB.Ctx is still live after shutdown, so a stray query could hold the pool open")
	}
}

// Shutdown cancels the app context under a scrape's ingest; that, or a client going away, isn't
// a database fault: a 503, and no error-level line.
func TestHandleErrorTreatsCancellationAsNoFault(t *testing.T) {
	var logs bytes.Buffer
	res, _ := testResources()
	res.Log = slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))
	s := &DatabaseService{}
	if err := s.Init(res); err != nil {
		t.Fatal(err)
	}

	for _, cause := range []error{context.Canceled, fmt.Errorf("bun: %w", context.Canceled)} {
		if e, ok := errors.AsType[*errs.Error](s.HandleError(cause)); !ok || e.Code != http.StatusServiceUnavailable {
			t.Errorf("HandleError(%v) = %v, want a 503", cause, e)
		}
	}
	if strings.Contains(logs.String(), "level=ERROR") {
		t.Errorf("a cancelled query was logged as an error:\n%s", logs.String())
	}
	if e, _ := errors.AsType[*errs.Error](s.HandleError(context.DeadlineExceeded)); e == nil || e.Code != http.StatusInternalServerError {
		t.Errorf("HandleError(DeadlineExceeded) = %v, want a 500: a slow query stays an error", e)
	}
}
