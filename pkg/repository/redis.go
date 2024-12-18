package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	db *redis.Client
}

func NewRedisDb(opt *redis.Options) (*redis.Client, error) {
	ctx := context.Background()
	client := redis.NewClient(opt)
	err := client.Ping(ctx).Err()
	return client, err
}

func NewCacheRedis(db *redis.Client) *Cache {
	return &Cache{
		db: db,
	}
}

func (r *Cache) SaveSecret(secretId string, secret *big.Int) error {
	ctx := context.Background()
	err := r.db.Set(ctx, secretId, secret.String(), time.Hour).Err()
	return err
}

func (r *Cache) GetSecret(secretId string) (*big.Int, error) {
	secretString, err := r.db.Get(context.Background(), secretId).Result()
	if err != nil {
		return nil, err
	}

	secret, ok := new(big.Int).SetString(secretString, 10)
	if !ok {
		return nil, fmt.Errorf("bad secret value: %s", secretString)
	}
	return secret, nil
}
