package backend

import (
	"errors"
	"fmt"
	"time"

	"github.com/YuleiGong/g_task/message"

	"github.com/go-redis/redis"
)

type Redis struct {
	url        string
	poolSize   int
	db         int
	password   string
	expireTime time.Duration
	client     *redis.Client
}

func NewRedis(cfg *RedisConf) *Redis {
	return &Redis{
		url:        cfg.url,
		password:   cfg.password,
		poolSize:   cfg.poolSize,
		db:         cfg.db,
		expireTime: cfg.expireTime,
	}
}

func (r *Redis) Clone() Backend {
	return &Redis{
		url:        r.url,
		password:   r.password,
		poolSize:   r.poolSize,
		db:         r.db,
		expireTime: r.expireTime,
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

func (r *Redis) GetResult(taskID string) (msg *message.MessageResult, err error) {
	key := fmt.Sprintf("%s-result", taskID)
	var m []byte
	if m, err = r.client.Get(key).Bytes(); err != nil {
		if errors.Is(err, redis.Nil) {
			err = ErrBackendNil
		}
		return
	}

	if msg, err = message.DefaultMessageResult.Deserialize(m); err != nil {
		return
	}

	return msg, err
}

func (r *Redis) SetResult(taskID string, msg *message.MessageResult) (err error) {
	var m string
	if m, err = msg.Serializa(); err != nil {
		return
	}
	key := fmt.Sprintf("%s-result", taskID)
	if _, err = r.client.Set(key, m, r.expireTime).Result(); err != nil {
		return
	}

	return err
}
