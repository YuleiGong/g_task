package broker

import (
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

func (r *Redis) Clone() Broker {
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

func (r *Redis) Pop() (taskID string, msg *message.Message, err error) {

	var vals []string
	vals, err = r.client.BRPop(popTimeout, queueName).Result()
	if err != nil {
		return
	}
	taskID = vals[1]

	var m []byte
	if m, err = r.client.Get(taskID).Bytes(); err != nil {
		return
	}

	if msg, err = message.DefaultMessage.Deserialize(m); err != nil {
		return
	}

	return taskID, msg, err
}

func (r *Redis) Push(taskID string, msg *message.Message) (err error) {
	pipeline := r.client.Pipeline()
	if _, err = pipeline.LPush(queueName, taskID).Result(); err != nil {
		return
	}

	var m string
	if m, err = msg.Serialize(); err != nil {
		return
	}

	if _, err = pipeline.Set(taskID, m, r.expireTime).Result(); err != nil {
		return
	}

	_, err = pipeline.Exec()

	return err
}

func (r *Redis) Del(taskID string) (err error) {
	if _, err = r.client.Del(taskID).Result(); err != nil {
		return
	}

	return err
}

func (r *Redis) Set(taskID string, msg *message.Message) (err error) {
	var m string
	if m, err = msg.Serialize(); err != nil {
		return
	}
	if _, err = r.client.Set(taskID, m, expireTime).Result(); err != nil {
		return
	}

	return err
}
