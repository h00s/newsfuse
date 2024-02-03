package models

import "github.com/go-raptor/raptor"

type Source struct {
	raptor.Model `json:"-"`
	Name         string
	Headlines    []Headline `json:"-"`
}

type Sources []Source
