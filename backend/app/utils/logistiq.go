package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-logistiq/handler"
	"github.com/go-raptor/components"
	"github.com/go-raptor/raptor/v3"
)

func NewLogistiqHandler(utils *raptor.Utils) (*handler.Handler, error) {
	var err error
	batchSize, err := strconv.Atoi(utils.Config.AppConfig["logistiq_batch_size"])
	if err != nil {
		return nil, errors.New("logistiq_batch_size must be an integer")
	}
	timeoutSeconds, err := strconv.Atoi(utils.Config.AppConfig["logistiq_timeout_seconds"])
	if err != nil {
		return nil, errors.New("logistiq_timeout_seconds must be an integer")
	}

	opts := handler.Options{
		Level:      components.ParseLogLevel(utils.Config.GeneralConfig.LogLevel),
		BatchSize:  batchSize,
		Timeout:    time.Duration(timeoutSeconds) * time.Second,
		NATSURL:    utils.Config.AppConfig["logistiq_nats_url"],
		Subject:    utils.Config.AppConfig["logistiq_subject"],
		SetHandler: utils.SetHandler,
	}

	return handler.New(opts), nil
}
