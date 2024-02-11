package models

import "github.com/go-raptor/raptor"

type Topic struct {
	raptor.Model
	Name    string
	Sources []Source `json:"-"`
}

type Topics []Topic
