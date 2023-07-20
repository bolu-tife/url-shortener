package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheStorage interface {
	SetShortUrlToLongUrl(string, string) error
	GetLongUrlFromShortUrl(string) (string, error)
}

type CacheStore struct {
	cache *redis.Client
}

func NewRedisStore() (*CacheStore, error) {
	redisUrl := GetConfig().RedisUrl
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return &CacheStore{cache: client}, nil
}

func (s *CacheStore) SetShortUrlToLongUrl(shortUrl, longUrl string) error {
	ctx := context.Background()
	key := "url:" + shortUrl

	return s.cache.Set(ctx, key, longUrl, time.Hour).Err()
}

func (s *CacheStore) GetLongUrlFromShortUrl(shortUrl string) (string, error) {
	ctx := context.Background()
	key := "url:" + shortUrl
	return s.cache.Get(ctx, key).Result()
}
