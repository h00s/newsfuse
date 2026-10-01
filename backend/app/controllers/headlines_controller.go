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
	topicID, err := ctx.QueryInt64("topicId")
	if err != nil {
		return err
	}
	var sourceID *int64
	if ctx.QueryParam("sourceId") != "" {
		id, err := ctx.QueryInt64("sourceId")
		if err != nil {
			return err
		}
		sourceID = &id
	}
	beforeID, err := pageCursor(ctx)
	if err != nil {
		return err
	}

	headlines, err := c.Headlines.List(topicID, sourceID, beforeID)
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
	beforeID, err := pageCursor(ctx)
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
	topicID, err := ctx.QueryInt64("topicId")
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
	id, err := ctx.ParamInt64("id")
	if err != nil {
		return err
	}
	story, err := c.Stories.GetByHeadline(ctx.Request().Context(), id)
	if err != nil {
		return err
	}
	return ctx.Data(models.NewStoryResponse(story))
}
