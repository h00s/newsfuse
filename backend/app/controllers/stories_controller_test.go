package controllers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/models"
)

func summarizePath(id int64) string {
	return fmt.Sprintf("/api/v1/stories/%d/summarize", id)
}

func TestSummarizeStoresThePlainTextSummary(t *testing.T) {
	headline := seedHeadlines(t, sourceBug, 1)[0]
	stored := seedStory(t, headline.ID, "<p>Prvi odlomak.</p><p>Drugi odlomak.</p>")

	story := raptor.DecodeJSON[storyJSON](t, app.TestPost(summarizePath(stored.ID), nil, newClient()), http.StatusOK)

	if story.Summary == "" || story.ID != stored.ID {
		t.Fatalf("story = %+v, want a summary for story %d", story, stored.ID)
	}
	var row models.Story
	if err := db(t).NewSelect().Model(&row).Where("id = ?", stored.ID).Scan(t.Context()); err != nil {
		t.Fatal(err)
	}
	if row.Summary != story.Summary {
		t.Errorf("stored summary = %q, want %q", row.Summary, story.Summary)
	}

	again := raptor.DecodeJSON[storyJSON](t, app.TestPost(summarizePath(stored.ID), nil, newClient()), http.StatusOK)
	if again.Summary != story.Summary {
		t.Errorf("second summary = %q, want the stored %q", again.Summary, story.Summary)
	}
}

func TestSummarizeUnknownStoryIs404(t *testing.T) {
	raptor.DecodeJSON[errorJSON](t, app.TestPost(summarizePath(999999999), nil, newClient()), http.StatusNotFound)
}

// A GET must never spend an LLM call. The /api/v1/{path...} catch-all answers it, so the status
// is its JSON 404 rather than a 405.
func TestSummarizeIsNotAGet(t *testing.T) {
	headline := seedHeadlines(t, sourceBug, 1)[0]
	stored := seedStory(t, headline.ID, "<p>Tekst.</p>")

	raptor.DecodeJSON[errorJSON](t, app.TestGet(summarizePath(stored.ID), newClient()), http.StatusNotFound)

	var row models.Story
	if err := db(t).NewSelect().Model(&row).Where("id = ?", stored.ID).Scan(t.Context()); err != nil {
		t.Fatal(err)
	}
	if row.Summary != "" {
		t.Errorf("GET stored summary %q, want none", row.Summary)
	}
}

func TestSummarizeRejectsCrossSiteRequests(t *testing.T) {
	rec := app.TestPost(summarizePath(1), nil, newClient(),
		raptor.WithHeader("Sec-Fetch-Site", "cross-site"),
		raptor.WithHeader("Origin", "https://evil.example"))
	if rec.Code != http.StatusForbidden {
		t.Errorf("cross-site POST = %d, want 403", rec.Code)
	}
}

func TestSummarizeIsRateLimited(t *testing.T) {
	client := newClient()
	for i := range 5 {
		if rec := app.TestPost(summarizePath(999999999), nil, client); rec.Code != http.StatusNotFound {
			t.Fatalf("request %d = %d, want 404", i+1, rec.Code)
		}
	}
	if rec := app.TestPost(summarizePath(999999999), nil, client); rec.Code != http.StatusTooManyRequests {
		t.Errorf("sixth request = %d, want 429", rec.Code)
	}
}

func TestSummarizeNonPositiveIDIs404(t *testing.T) {
	raptor.DecodeJSON[errorJSON](t, app.TestPost(summarizePath(0), nil, newClient()), http.StatusNotFound)
}
