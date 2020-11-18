package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// error
var (
	ErrNotFound = errors.New("not found")
)

// Redis info struct
type Redis struct {
	Address     string `json:"Address"      envconfig:"REDIS_ADDRESS"       validate:"required"`
	Password    string `json:"Password"     envconfig:"REDIS_PASSWORD"`
	DB          int    `json:"DB"           envconfig:"REDIS_DB"`
	DiscountTTL int    `json:"DiscountTTL"  envconfig:"REDIS_DISCOUNT_TTL"`
}

// Cache cache implementation based on Redis
type Cache struct {
	redis *redis.Client
}

// New init cache client
func New(cfg Redis) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping().Result()

	c := Cache{
		redis: client,
	}
	return &c, err
}

// Get retrieves value from cache
func (m *Cache) Get(key string) (string, error) {
	value, err := m.redis.Get(key).Result()

	if err == redis.Nil {
		return value, ErrNotFound
	}

	return value, err
}

// Set stores value to cache
func (m *Cache) Set(key string, value interface{}, exp time.Duration) error {
	return m.redis.Set(key, value, exp).Err()
}

// Del invalidates value in cache
func (m *Cache) Del(key string) error {
	deletedAt, err := m.redis.Del(key).Result()
	if err == nil && deletedAt != 1 {
		return fmt.Errorf("failed to remove record, key: %s", key)
	}
	return err
}
