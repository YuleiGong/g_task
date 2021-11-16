package backend

import (
	"errors"
	"fmt"

	"github.com/YuleiGong/g_task/message"

	"github.com/go-redis/redis"
)

type Redis struct {
	url      string
	poolSize int
	db       int
	password string
	client   *redis.Client
}

func NewRedis(url, password string, db, poolSize int) *Redis {
	return &Redis{
		url:      url,
		password: password,
		poolSize: poolSize,
		db:       db,
	}
}

func (r *Redis) Clone() Backend {
	return &Redis{
		url:      r.url,
		password: r.password,
		poolSize: r.poolSize,
		db:       r.db,
	}
}

func (r *Redis) Activate() (err error) {
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.url,
		PoolSize: r.poolSize,
		DB:       r.db,
		Password: r.password,
	})

	return r.client.Ping().Err()
}

func (r *Redis) GetResult(taskID string) (msg string, err error) {
	var m []byte
	if m, err = r.client.Get(taskID).Bytes(); err != nil {
		if errors.Is(redis.Nil, err) {
			return msg, nil
		}
		return
	}

	return string(m), err
}

func (r *Redis) SetResult(taskID string, msg *message.MessageResult) (err error) {
	var m string
	if m, err = msg.Serializa(); err != nil {
		return
	}
	key := fmt.Sprintf("%s-result", taskID)
	if _, err = r.client.Set(key, m, expireTime).Result(); err != nil {
		return
	}

	return err
}
