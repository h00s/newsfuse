package controllers_test

import (
	"net/http"
	"slices"
	"testing"
)

func TestSourcesIndexListsEverySource(t *testing.T) {
	sources := get[[]sourceJSON](t, "/api/v1/sources", http.StatusOK)

	if len(sources) != 11 {
		t.Fatalf("got %d sources, want 11", len(sources))
	}
	hcl := sourceJSON{ID: sourceHCL, Name: "HCL", TopicID: topicTech, IsScrapable: true}
	if sources[10] != hcl {
		t.Errorf("last source = %+v, want %+v", sources[10], hcl)
	}
}

func TestSourcesIndexFiltersByTopic(t *testing.T) {
	sources := get[[]sourceJSON](t, "/api/v1/sources?topicId=4", http.StatusOK)

	got := make([]int64, len(sources))
	for i, s := range sources {
		got[i] = s.ID
	}
	if want := []int64{sourceHackerNews, sourceBug, sourceHCL}; !slices.Equal(got, want) {
		t.Errorf("source ids = %v, want %v", got, want)
	}
}

func TestSourcesIndexUnknownTopicIs404(t *testing.T) {
	get[errorJSON](t, "/api/v1/sources?topicId=999", http.StatusNotFound)
}

func TestSourcesIndexInvalidTopicIs400(t *testing.T) {
	get[errorJSON](t, "/api/v1/sources?topicId=tech", http.StatusBadRequest)
}
