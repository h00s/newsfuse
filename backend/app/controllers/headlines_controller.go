package controllers

import (
	"strings"
	"unicode/utf8"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/services"
)

const (
	minSearchLength = 3
	maxSearchLength = 100
)

type HeadlinesController struct {
	raptor.Controller

	Headlines *services.HeadlinesService
	Stories   *services.StoriesService
}

// Index pages through a topic's headlines, newest first: ?topicId= is required, ?sourceId=
// narrows the list to one of the topic's sources, and ?beforeId= continues after a page.
func (c *HeadlinesController) Index(ctx *raptor.Context) error {
	topicID, err := requiredQueryID(ctx, "topicId")
	if err != nil {
		return err
	}
	sourceID, hasSource, err := queryID(ctx, "sourceId")
	if err != nil {
		return err
	}
	beforeID, _, err := queryID(ctx, "beforeId")
	if err != nil {
		return err
	}

	var source *int64
	if hasSource {
		source = &sourceID
	}
	headlines, err := c.Headlines.List(topicID, source, beforeID)
	if err != nil {
		return err
	}
	return ctx.Data(models.NewHeadlineResponses(headlines))
}

// Search pages through the headlines whose title contains ?query=, newest first.
func (c *HeadlinesController) Search(ctx *raptor.Context) error {
	query := strings.TrimSpace(ctx.QueryParam("query"))
	if n := utf8.RuneCountInString(query); n < minSearchLength || n > maxSearchLength {
		return errs.NewErrorBadRequest("Invalid query",
			"minLength", minSearchLength, "maxLength", maxSearchLength)
	}
	beforeID, _, err := queryID(ctx, "beforeId")
	if err != nil {
		return err
	}

	headlines, err := c.Headlines.Search(query, beforeID)
	if err != nil {
		return err
	}
	return ctx.Data(models.NewHeadlineResponses(headlines))
}

// Count is how many of a topic's headlines were published after ?since= (RFC 3339).
func (c *HeadlinesController) Count(ctx *raptor.Context) error {
	topicID, err := requiredQueryID(ctx, "topicId")
	if err != nil {
		return err
	}
	since, err := queryTime(ctx, "since")
	if err != nil {
		return err
	}

	count, err := c.Headlines.CountSince(topicID, since)
	if err != nil {
		return err
	}
	return ctx.Data(models.HeadlineCountResponse{Count: count})
}

// Story is the article behind a headline, scraped from its site on the first read.
func (c *HeadlinesController) Story(ctx *raptor.Context) error {
	id, err := pathID(ctx)
	if err != nil {
		return err
	}
	story, err := c.Stories.GetByHeadline(ctx.Request().Context(), id)
	if err != nil {
		return err
	}
	return ctx.Data(models.NewStoryResponse(story))
}
