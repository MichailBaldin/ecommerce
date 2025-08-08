package repository

import (
	"context"
	"ecommerce/services/products/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepo struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepo(addr string) CacheRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	return &RedisRepo{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisRepo) Set(product *models.Product) error {
	key := fmt.Sprintf("product:%d", product.ID)
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, data, 5*time.Minute).Err()
}

func (r *RedisRepo) Get(id int) (*models.Product, error) {
	key := fmt.Sprintf("product:%d", id)
	data, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // cache miss - возвращаем nil, nil
	}
	if err != nil {
		return nil, err
	}

	var product models.Product
	err = json.Unmarshal([]byte(data), &product)
	return &product, err
}

func (r *RedisRepo) Delete(id int) error {
	key := fmt.Sprintf("product:%d", id)
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisRepo) Close() error {
	return r.client.Close()
}
