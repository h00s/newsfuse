package controllers_test

import (
	"net/http"
	"testing"

	"github.com/go-raptor/raptor/v4"
)

// The probes are public and JSON. A browser-like Accept header must still get the probe, not the
// SPA shell that the / catch-all serves to navigations.
func TestProbesAnswerOK(t *testing.T) {
	for _, path := range []string{"/healthz", "/readyz"} {
		rec := app.TestGet(path, newClient(), raptor.WithHeader("Accept", "text/html"))
		res := raptor.DecodeJSON[struct {
			Status string `json:"status"`
		}](t, rec, http.StatusOK)
		if res.Status != "ok" {
			t.Errorf("GET %s status = %q, want ok", path, res.Status)
		}
	}
}
