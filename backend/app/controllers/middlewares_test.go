package controllers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/go-raptor/raptor/v4"
)

var securityHeaders = map[string]string{
	"X-Content-Type-Options":     "nosniff",
	"X-Frame-Options":            "DENY",
	"Content-Security-Policy":    "frame-ancestors 'none'",
	"Cross-Origin-Opener-Policy": "same-origin",
	"Referrer-Policy":            "strict-origin-when-cross-origin",
}

// requestid and secure run before csrf, so even its 403 carries a request id and the headers.
func TestResponsesCarryARequestIDAndSecurityHeaders(t *testing.T) {
	responses := map[string]*httptest.ResponseRecorder{
		"a list":     app.TestGet("/api/v1/topics", newClient()),
		"a JSON 404": app.TestGet("/api/v1/nope", newClient()),
		"a cross-site 403": app.TestPost(summarizePath(1), nil, newClient(),
			raptor.WithHeader("Sec-Fetch-Site", "cross-site"),
			raptor.WithHeader("Origin", "https://evil.example")),
	}
	for name, rec := range responses {
		if rec.Header().Get("X-Request-Id") == "" {
			t.Errorf("%s (%d) has no X-Request-Id", name, rec.Code)
		}
		for header, want := range securityHeaders {
			if got := rec.Header().Get(header); got != want {
				t.Errorf("%s (%d): %s = %q, want %q", name, rec.Code, header, got, want)
			}
		}
	}
}
