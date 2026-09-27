package controllers_test

import (
	"net/http"
	"testing"
)

func TestUnknownAPIPathIsJSON404(t *testing.T) {
	res := get[errorJSON](t, "/api/v1/nope", http.StatusNotFound)
	if res.Code != http.StatusNotFound {
		t.Errorf("error body code = %d, want 404", res.Code)
	}
}

func TestOldAPIPathsAreGone(t *testing.T) {
	for _, path := range []string{"/api/v1/topics/1/headlines", "/api/v1/topics/1/headlines/count"} {
		if rec := app.TestGet(path, newClient()); rec.Code != http.StatusNotFound {
			t.Errorf("GET %s = %d, want 404", path, rec.Code)
		}
	}
}
