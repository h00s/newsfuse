// Package controllers holds the HTTP actions: parse the request, call a service, map the
// result to a response DTO.
package controllers

import (
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
)

// pageCursor reads ?beforeId=, the last headline of the previous page. Absent or empty means the
// first page.
func pageCursor(ctx *raptor.Context) (int64, error) {
	if ctx.QueryParam("beforeId") == "" {
		return 0, nil
	}
	return ctx.QueryInt64("beforeId")
}

// queryTime reads a required RFC 3339 time from the query string.
func queryTime(ctx *raptor.Context, name string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, ctx.QueryParam(name))
	if err != nil {
		return time.Time{}, errs.NewErrorBadRequest("Invalid " + name + ": expected an RFC 3339 time")
	}
	return t, nil
}
