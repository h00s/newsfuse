package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-logistiq/handler"
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/core"
)

func NewLogistiqHandler(config *raptor.Config) (*handler.Handler, error) {
	var err error
	batchSize, err := strconv.Atoi(config.AppConfig["logistiq_batch_size"])
	if err != nil {
		return nil, errors.New("logistiq_batch_size must be an integer")
	}
	timeoutSeconds, err := strconv.Atoi(config.AppConfig["logistiq_timeout_seconds"])
	if err != nil {
		return nil, errors.New("logistiq_timeout_seconds must be an integer")
	}

	opts := handler.Options{
		Level:     core.ParseLogLevel(config.GeneralConfig.LogLevel),
		BatchSize: batchSize,
		Timeout:   time.Duration(timeoutSeconds) * time.Second,
		NATSURL:   config.AppConfig["logistiq_nats_url"],
		Subject:   config.AppConfig["logistiq_subject"],
	}

	return handler.New(opts), nil
}
