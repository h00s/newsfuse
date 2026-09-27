package controllers_test

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/services"
	"github.com/uptrace/bun"
)

// Topics and sources are seeded by the migrations and never change, so tests read them by id.
const (
	topicBBZ  = 1 // sources 1 klikni.hr, 2 MojPortal.hr, 3 Radio Daruvar
	topicTech = 4 // sources 8 Hacker News, 9 Bug, 11 HCL

	sourceKlikni     = 1
	sourceHackerNews = 8
	sourceBug        = 9
	sourceHCL        = 11
)

// The API's JSON, spelled out here rather than borrowed from the models, so a changed json tag
// fails a test instead of changing both sides at once.
type (
	topicJSON struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	sourceJSON struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		TopicID     int64  `json:"topicId"`
		IsScrapable bool   `json:"isScrapable"`
	}
	headlineJSON struct {
		ID          int64      `json:"id"`
		Title       string     `json:"title"`
		URL         string     `json:"url"`
		PublishedAt time.Time  `json:"publishedAt"`
		Source      sourceJSON `json:"source"`
	}
	storyJSON struct {
		ID         int64  `json:"id"`
		HeadlineID int64  `json:"headlineId"`
		Content    string `json:"content"`
		Summary    string `json:"summary"`
	}
	errorJSON struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func db(t *testing.T) *bun.DB {
	t.Helper()
	service := raptor.GetService[services.DatabaseService](app)
	if service == nil {
		t.Fatal("DatabaseService is not registered")
	}
	return service.Conn()
}

func headlinesService(t *testing.T) *services.HeadlinesService {
	t.Helper()
	service := raptor.GetService[services.HeadlinesService](app)
	if service == nil {
		t.Fatal("HeadlinesService is not registered")
	}
	return service
}

var clientIPs atomic.Uint32

// newClient simulates a separate browser, with its own address and so its own limiter bucket.
func newClient() raptor.TestRequestOption {
	n := clientIPs.Add(1)
	return raptor.WithRemoteAddr(fmt.Sprintf("10.%d.%d.%d", byte(n>>16), byte(n>>8), byte(n)))
}

// decode asserts the status, then decodes the body into T.
func decode[T any](t *testing.T, rec *httptest.ResponseRecorder, want int) T {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status = %d, want %d; body %s", rec.Code, want, rec.Body)
	}
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("decode %T: %v; body %s", v, err, rec.Body)
	}
	return v
}

func get[T any](t *testing.T, path string, want int) T {
	t.Helper()
	return decode[T](t, app.TestGet(path, newClient()), want)
}

var suffixes atomic.Uint64

func uniqueSuffix() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + "-" + strconv.FormatUint(suffixes.Add(1), 10)
}

// seedHeadlines stores n headlines of one source through Ingest, the path the scrapers use, and
// returns them oldest first with their ids. Like a scraper's batch, the first title is the newest
// (published now) and each next one a minute older. The rows are deleted when the test ends.
func seedHeadlines(t *testing.T, sourceID int64, n int, titles ...string) models.Headlines {
	t.Helper()
	prefix := "https://example.test/" + uniqueSuffix() + "/"
	now := time.Now()
	batch := make(models.Headlines, n)
	for i := range n {
		title := fmt.Sprintf("Naslov %s %d", prefix, n-i)
		if i < len(titles) {
			title = titles[i]
		}
		batch[i] = models.Headline{
			SourceID:    sourceID,
			Title:       title,
			URL:         prefix + strconv.Itoa(n-i),
			PublishedAt: now.Add(-time.Duration(i) * time.Minute),
		}
	}

	t.Cleanup(func() {
		_, _ = db(t).NewDelete().Model((*models.Headline)(nil)).
			Where("url LIKE ?", prefix+"%").Exec(context.Background()) // cascades to stories
	})
	inserted, err := headlinesService(t).Ingest(batch)
	if err != nil {
		t.Fatalf("seed headlines: %v", err)
	}
	if inserted != n {
		t.Fatalf("seed headlines: inserted %d, want %d", inserted, n)
	}

	var stored models.Headlines
	if err := db(t).NewSelect().Model(&stored).Where("url LIKE ?", prefix+"%").Order("id").Scan(t.Context()); err != nil {
		t.Fatalf("load seeded headlines: %v", err)
	}
	return stored
}

// seedStory stores a story for a headline, as if it had been scraped before.
func seedStory(t *testing.T, headlineID int64, content string) *models.Story {
	t.Helper()
	story := &models.Story{HeadlineID: headlineID, Content: content}
	if _, err := db(t).NewInsert().Model(story).Returning("*").Exec(t.Context()); err != nil {
		t.Fatalf("seed story: %v", err)
	}
	return story
}

func ids(headlines []headlineJSON) []int64 {
	out := make([]int64, len(headlines))
	for i, h := range headlines {
		out[i] = h.ID
	}
	return out
}
