package controllers

import (
	"strconv"

	"github.com/go-raptor/errs"
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/services"
)

type StoriesController struct {
	raptor.Controller

	Stories *services.StoriesService
}

func (sc *StoriesController) Get(c raptor.State) error {
	headlineID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return errs.NewErrorBadRequest("Invalid Headline ID")
	}

	story, err := sc.Stories.Get(headlineID)
	if err != nil {
		return errs.NewErrorNotFound("Story not found")
	}

	return c.JSONResponse(story)
}

func (sc *StoriesController) Summarize(c raptor.State) error {
	storyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return errs.NewErrorBadRequest("Invalid Story ID")
	}

	story, err := sc.Stories.Summarize(storyID)
	if err != nil {
		return errs.NewErrorNotFound("Story not found")
	}

	return c.JSONResponse(story)
}
