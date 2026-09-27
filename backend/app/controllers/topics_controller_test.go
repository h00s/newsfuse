package controllers_test

import (
	"net/http"
	"slices"
	"testing"
)

func TestTopicsIndexListsSeededTopicsInOrder(t *testing.T) {
	topics := get[[]topicJSON](t, "/api/v1/topics", http.StatusOK)

	want := []topicJSON{{1, "BBŽ"}, {2, "Hrvatska"}, {3, "Svijet"}, {4, "Tech"}}
	if !slices.Equal(topics, want) {
		t.Errorf("topics = %v, want %v", topics, want)
	}
}
