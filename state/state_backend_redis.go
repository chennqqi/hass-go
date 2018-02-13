package state

import (
	"time"

	"gopkg.in/redis.v5"
)

type redisStateDB struct {
	redis *redis.Client
}

func newRedisStateDB(configfilepath string) (db *redisStateDB, err error) {
	r := &redisStateDB{}

	return r, err
}

func (r *redisStateDB) open() error {
	r.redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return nil
}

func (r *redisStateDB) close() error {
	return r.redis.Close()
}

func (r *redisStateDB) setString(key string, value string) error {
	status := r.redis.Set(key, value, time.Duration(0))
	return status.Err()
}

func (r *redisStateDB) getString(key string) (value string, err error) {
	status := r.redis.Get(key)
	return status.String(), status.Err()
}
