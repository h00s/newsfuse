package services

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/core"
)

// testResources are what Raptor injects into a service, plus the cancel its shutdown runs once
// requests have drained.
func testResources() (res *raptor.Resources, shutdown func()) {
	res = raptor.NewTestResources()
	return res, (&core.Core{Resources: res}).CancelAppContext
}
