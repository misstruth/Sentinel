package retriever

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache Redis 缓存实现
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr string, ttl time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		ttl:    ttl,
	}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]Document, bool, error) {
	cacheKey := c.buildKey(key)
	val, err := c.client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var docs []Document
	if err := json.Unmarshal([]byte(val), &docs); err != nil {
		return nil, false, err
	}

	return docs, true, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, docs []Document, ttl time.Duration) error {
	cacheKey := c.buildKey(key)
	data, err := json.Marshal(docs)
	if err != nil {
		return err
	}

	if ttl == 0 {
		ttl = c.ttl
	}

	return c.client.Set(ctx, cacheKey, data, ttl).Err()
}

func (c *RedisCache) buildKey(query string) string {
	hash := md5.Sum([]byte(query))
	return "rag:query:" + hex.EncodeToString(hash[:])
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}
