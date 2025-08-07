package repository

import (
	"context"
	"ecommerce/services/users/models"
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

func (r *RedisRepo) Set(user *models.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, data, 5*time.Minute).Err()
}

func (r *RedisRepo) Get(id int) (*models.User, error) {
	key := fmt.Sprintf("user:%d", id)
	data, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // cache miss - возвращаем nil, nil
	}
	if err != nil {
		return nil, err
	}

	var user models.User
	err = json.Unmarshal([]byte(data), &user)
	return &user, err
}

func (r *RedisRepo) Delete(id int) error {
	key := fmt.Sprintf("user:%d", id)
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisRepo) Close() error {
	return r.client.Close()
}
