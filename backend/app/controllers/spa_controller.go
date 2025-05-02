package controllers

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/go-raptor/raptor/v4"
)

type SPAController struct {
	raptor.Controller

	lock  sync.RWMutex
	files map[string]bool

	directory string
	file      string
}

func NewSPAController(directory, file string) *SPAController {
	return &SPAController{
		files:     make(map[string]bool),
		directory: directory,
		file:      directory + "/" + file,
	}
}

func (sc *SPAController) Index(c *raptor.Context) error {
	requestedPath := c.Request().URL.Path
	filePath := filepath.Join(sc.directory, requestedPath)

	sc.lock.RLock()
	exists, inCache := sc.files[filePath]
	sc.lock.RUnlock()

	if inCache {
		if exists {
			return c.File(filePath)
		}
		return c.File(sc.file)
	}

	fileInfo, err := os.Stat(filePath)
	if err == nil && !fileInfo.IsDir() {
		sc.lock.Lock()
		sc.files[filePath] = true
		sc.lock.Unlock()
		return c.File(filePath)
	}

	sc.lock.Lock()
	sc.files[filePath] = false
	sc.lock.Unlock()
	return c.File(sc.file)
}
