package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fakubwoy/go-file-share/internal/config"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return client, nil
}
