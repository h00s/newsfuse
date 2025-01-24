package services

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-raptor/raptor/v3"
	"github.com/redis/go-redis/v9"
)

type Memstore struct {
	raptor.Service
	client *redis.Client
}

func NewMemstore(c *raptor.Config) *Memstore {
	db, err := strconv.Atoi(c.AppConfig["redis_db"])
	if err != nil {
		db = 0
	}

	return &Memstore{
		client: redis.NewClient(&redis.Options{
			Addr:     c.AppConfig["redis_address"],
			Password: c.AppConfig["redis_password"],
			DB:       db,
		}),
	}
}

func (m *Memstore) Get(key string) (string, error) {
	value := m.client.Get(context.Background(), key)
	if value.Err() == redis.Nil {
		return "", nil
	}

	data, err := value.Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

func (m *Memstore) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.client.Set(context.Background(), key, data, 0).Err()
}
