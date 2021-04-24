package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/suborbital/reactr/rt"
)

type cache struct {
	client *redis.Client
}

func (c *cache) Set(key string, val []byte, ttl int) error {
	err := c.client.Set(key, val, time.Duration(ttl)).Err()
	return err
}

func (c *cache) Get(key string) ([]byte, error) {
	val, err := c.client.Get(key).Result()
	return []byte(val), err
}

func (c *cache) Delete(key string) error {
	err := c.client.Del(key).Err()
	return err
}

func NewCache(host string, port int) rt.Cache {
	c := &cache{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
	return c
}
