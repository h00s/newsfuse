// Package controllers holds the HTTP actions: parse the request, call a service, map the
// result to a response DTO.
package controllers

import (
	"strconv"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
)

// pathID parses the {id} path parameter; garbage is a 400.
func pathID(ctx *raptor.Context) (int64, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errs.NewErrorBadRequest("Invalid ID")
	}
	return id, nil
}

// queryID reads an optional id from the query string. ok is false when the parameter is absent;
// anything but a positive integer is a 400.
func queryID(ctx *raptor.Context, name string) (id int64, ok bool, err error) {
	raw := ctx.QueryParam(name)
	if raw == "" {
		return 0, false, nil
	}
	id, err = strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		return 0, false, errs.NewErrorBadRequest("Invalid " + name)
	}
	return id, true, nil
}

// requiredQueryID is queryID for a parameter the action can't do without.
func requiredQueryID(ctx *raptor.Context, name string) (int64, error) {
	id, ok, err := queryID(ctx, name)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errs.NewErrorBadRequest("Missing " + name)
	}
	return id, nil
}

// queryTime reads a required RFC 3339 time from the query string.
func queryTime(ctx *raptor.Context, name string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, ctx.QueryParam(name))
	if err != nil {
		return time.Time{}, errs.NewErrorBadRequest("Invalid " + name + ": expected an RFC 3339 time")
	}
	return t, nil
}
