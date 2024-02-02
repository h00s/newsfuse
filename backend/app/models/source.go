package models

import "github.com/go-raptor/raptor"

type Source struct {
	raptor.Model
	Name      string
	Headlines []Headline
}

type Sources []Source
