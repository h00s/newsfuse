package controllers_test

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/h00s/newsfuse/app/models"
)

func TestHeadlinesIndexIsNewestFirstWithItsSource(t *testing.T) {
	seeded := seedHeadlines(t, sourceBug, 3)

	headlines := get[[]headlineJSON](t, "/api/v1/headlines?topicId=4", http.StatusOK)

	if len(headlines) < 3 {
		t.Fatalf("got %d headlines, want at least 3", len(headlines))
	}
	want := []int64{seeded[2].ID, seeded[1].ID, seeded[0].ID}
	if got := ids(headlines[:3]); !slices.Equal(got, want) {
		t.Errorf("first ids = %v, want %v (newest first)", got, want)
	}
	bug := sourceJSON{ID: sourceBug, Name: "Bug", TopicID: topicTech, IsScrapable: true}
	if headlines[0].Source != bug {
		t.Errorf("source = %+v, want %+v", headlines[0].Source, bug)
	}
	if headlines[0].Title != seeded[2].Title || headlines[0].URL != seeded[2].URL || headlines[0].PublishedAt.IsZero() {
		t.Errorf("headline = %+v, want title, url and publishedAt of %+v", headlines[0], seeded[2])
	}
}

func TestHeadlinesIndexPagesWithBeforeID(t *testing.T) {
	seeded := seedHeadlines(t, sourceHackerNews, 35)

	first := get[[]headlineJSON](t, "/api/v1/headlines?topicId=4", http.StatusOK)
	if len(first) != 30 {
		t.Fatalf("first page has %d headlines, want 30", len(first))
	}
	next := get[[]headlineJSON](t, fmt.Sprintf("/api/v1/headlines?topicId=4&beforeId=%d", first[29].ID), http.StatusOK)

	want := []int64{seeded[4].ID, seeded[3].ID, seeded[2].ID, seeded[1].ID, seeded[0].ID}
	if len(next) < 5 || !slices.Equal(ids(next[:5]), want) {
		t.Errorf("second page starts with %v, want %v", ids(next[:min(5, len(next))]), want)
	}
}

func TestHeadlinesIndexFiltersBySource(t *testing.T) {
	seedHeadlines(t, sourceHackerNews, 2)
	bug := seedHeadlines(t, sourceBug, 2)

	headlines := get[[]headlineJSON](t, "/api/v1/headlines?topicId=4&sourceId=9", http.StatusOK)

	for _, h := range headlines {
		if h.Source.ID != sourceBug {
			t.Fatalf("headline %d is from source %d, want only %d", h.ID, h.Source.ID, sourceBug)
		}
	}
	if got := ids(headlines); !slices.Contains(got, bug[0].ID) || !slices.Contains(got, bug[1].ID) {
		t.Errorf("ids = %v, want them to include %d and %d", got, bug[0].ID, bug[1].ID)
	}
}

func TestHeadlinesIndexSourceOfAnotherTopicIs404(t *testing.T) {
	get[errorJSON](t, fmt.Sprintf("/api/v1/headlines?topicId=4&sourceId=%d", sourceKlikni), http.StatusNotFound)
}

func TestHeadlinesIndexUnknownTopicIs404(t *testing.T) {
	get[errorJSON](t, "/api/v1/headlines?topicId=999", http.StatusNotFound)
}

func TestHeadlinesIndexRejectsInvalidParameters(t *testing.T) {
	for _, query := range []string{"", "?topicId=", "?topicId=tech", "?topicId=99999999999999999999", "?topicId=4&beforeId=x", "?topicId=4&sourceId=x"} {
		rec := app.TestGet("/api/v1/headlines"+query, newClient())
		if rec.Code != http.StatusBadRequest {
			t.Errorf("GET /api/v1/headlines%s = %d, want 400", query, rec.Code)
		}
	}
}

// Ids are parsed, not range-checked: an id no row can have is unknown like any other, and an
// empty or non-positive beforeId is the first page.
func TestHeadlinesIndexTreatsNonPositiveIDsAsUnknown(t *testing.T) {
	seedHeadlines(t, sourceBug, 2)
	for _, query := range []string{"?topicId=0", "?topicId=-1", "?topicId=4&sourceId=0"} {
		get[errorJSON](t, "/api/v1/headlines"+query, http.StatusNotFound)
	}
	first := ids(get[[]headlineJSON](t, "/api/v1/headlines?topicId=4", http.StatusOK))
	for _, query := range []string{"?topicId=4&beforeId=0", "?topicId=4&beforeId=-1", "?topicId=4&beforeId=", "?topicId=4&sourceId="} {
		if got := ids(get[[]headlineJSON](t, "/api/v1/headlines"+query, http.StatusOK)); !slices.Equal(got, first) {
			t.Errorf("%s ids = %v, want the first page %v", query, got, first)
		}
	}
}

func TestHeadlinesIndexServesNewHeadlinesAfterIngest(t *testing.T) {
	seedHeadlines(t, sourceHCL, 1)
	get[[]headlineJSON](t, "/api/v1/headlines?topicId=4", http.StatusOK) // fills the first-page cache

	later := seedHeadlines(t, sourceHCL, 1)

	headlines := get[[]headlineJSON](t, "/api/v1/headlines?topicId=4", http.StatusOK)
	if len(headlines) == 0 || headlines[0].ID != later[0].ID {
		t.Errorf("first headline = %v, want the one ingested last (%d)", ids(headlines[:min(1, len(headlines))]), later[0].ID)
	}
}

func TestHeadlinesIngestSkipsKnownURLs(t *testing.T) {
	seeded := seedHeadlines(t, sourceHCL, 2)
	again := models.Headlines{
		{SourceID: sourceHCL, Title: "Ponovljeno", URL: seeded[0].URL},
		{SourceID: sourceHCL, Title: "Ponovljeno", URL: seeded[1].URL},
	}

	inserted, err := headlinesService(t).Ingest(again)
	if err != nil {
		t.Fatal(err)
	}
	if inserted != 0 {
		t.Errorf("Ingest of known urls inserted %d, want 0", inserted)
	}
}

func TestHeadlinesSearchFindsTitles(t *testing.T) {
	token := "pretraga" + strings.ReplaceAll(uniqueSuffix(), "-", "")
	seeded := seedHeadlines(t, sourceKlikni, 2, "Prvi "+token, "Nešto drugo")

	headlines := get[[]headlineJSON](t, "/api/v1/headlines/search?query="+token, http.StatusOK)

	if got := ids(headlines); len(got) != 1 || got[0] != seeded[1].ID {
		t.Errorf("search ids = %v, want only %d", got, seeded[1].ID)
	}
	if headlines[0].Source.ID != sourceKlikni {
		t.Errorf("source = %+v, want klikni.hr", headlines[0].Source)
	}
}

func TestHeadlinesSearchTakesWildcardsLiterally(t *testing.T) {
	token := "posto" + strings.ReplaceAll(uniqueSuffix(), "-", "")
	seeded := seedHeadlines(t, sourceKlikni, 2, "Rast od 100% "+token, "Rast od 1000 "+token)

	headlines := get[[]headlineJSON](t, "/api/v1/headlines/search?query="+url.QueryEscape("100% "+token), http.StatusOK)

	if got := ids(headlines); len(got) != 1 || got[0] != seeded[1].ID {
		t.Errorf("search ids = %v, want only %d (the literal 100%%)", got, seeded[1].ID)
	}
}

func TestHeadlinesSearchRejectsShortQueries(t *testing.T) {
	for _, query := range []string{"", "ab", url.QueryEscape("  ab  ")} {
		rec := app.TestGet("/api/v1/headlines/search?query="+query, newClient())
		if rec.Code != http.StatusBadRequest {
			t.Errorf("search %q = %d, want 400", query, rec.Code)
		}
	}
}

func TestHeadlinesCountSince(t *testing.T) {
	seedHeadlines(t, sourceHCL, 3) // published now, 1 and 2 minutes ago
	since := url.QueryEscape(time.Now().Add(-90 * time.Second).Format(time.RFC3339))

	res := get[struct {
		Count int `json:"count"`
	}](t, "/api/v1/headlines/count?topicId=4&since="+since, http.StatusOK)

	if res.Count != 2 {
		t.Errorf("count = %d, want 2", res.Count)
	}
}

func TestHeadlinesCountRejectsInvalidParameters(t *testing.T) {
	since := url.QueryEscape(time.Now().Format(time.RFC3339))
	get[errorJSON](t, "/api/v1/headlines/count?topicId=4", http.StatusBadRequest)
	get[errorJSON](t, "/api/v1/headlines/count?topicId=4&since=1727430000000", http.StatusBadRequest)
	get[errorJSON](t, "/api/v1/headlines/count?since="+since, http.StatusBadRequest)
	get[errorJSON](t, "/api/v1/headlines/count?topicId=999&since="+since, http.StatusNotFound)
}

func TestHeadlineStoryServesTheStoredStorySanitized(t *testing.T) {
	headline := seedHeadlines(t, sourceBug, 1)[0]
	stored := seedStory(t, headline.ID, `<p>Tekst članka</p><script>alert(1)</script>`)

	story := get[storyJSON](t, fmt.Sprintf("/api/v1/headlines/%d/story", headline.ID), http.StatusOK)

	if story.ID != stored.ID || story.HeadlineID != headline.ID {
		t.Errorf("story = %+v, want id %d for headline %d", story, stored.ID, headline.ID)
	}
	if !strings.Contains(story.Content, "Tekst članka") || strings.Contains(story.Content, "script") {
		t.Errorf("content = %q, want the text without the script", story.Content)
	}
}

func TestHeadlineStoryOfUnknownHeadlineIs404(t *testing.T) {
	get[errorJSON](t, "/api/v1/headlines/999999999/story", http.StatusNotFound)
}

func TestHeadlineStoryInvalidIDIs400(t *testing.T) {
	for _, id := range []string{"abc", "99999999999999999999"} {
		get[errorJSON](t, "/api/v1/headlines/"+id+"/story", http.StatusBadRequest)
	}
}

func TestHeadlineStoryOfNonPositiveIDIs404(t *testing.T) {
	get[errorJSON](t, "/api/v1/headlines/0/story", http.StatusNotFound)
}
