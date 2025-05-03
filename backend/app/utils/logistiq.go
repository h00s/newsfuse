package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-logistiq/handler"
	"github.com/go-raptor/raptor/v4/core"
)

func NewLogistiqHandler(c *core.Core) (*handler.Handler, error) {
	var err error
	batchSize, err := strconv.Atoi(c.Resources.Config.AppConfig["logistiq_batch_size"])
	if err != nil {
		return nil, errors.New("logistiq_batch_size must be an integer")
	}
	timeoutSeconds, err := strconv.Atoi(c.Resources.Config.AppConfig["logistiq_timeout_seconds"])
	if err != nil {
		return nil, errors.New("logistiq_timeout_seconds must be an integer")
	}

	opts := handler.Options{
		Level:      core.ParseLogLevel(c.Resources.Config.GeneralConfig.LogLevel),
		BatchSize:  batchSize,
		Timeout:    time.Duration(timeoutSeconds) * time.Second,
		NATSURL:    c.Resources.Config.AppConfig["logistiq_nats_url"],
		Subject:    c.Resources.Config.AppConfig["logistiq_subject"],
		SetHandler: c.Resources.SetLogHandler,
	}

	return handler.New(opts), nil
}
