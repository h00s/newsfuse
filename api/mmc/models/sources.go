package models

import "github.com/h00s/newsfuse/api/services"

type Sources struct {
	services *services.Services
}

func NewSources(services *services.Services) *Sources {
	return &Sources{
		services: services,
	}
}
