package backend

import "time"

const (
	expireTime = 1 * time.Hour
	poolSize   = 20
)

type RedisConf struct {
	url        string
	poolSize   int
	db         int
	password   string
	expireTime time.Duration
}

func NewRedisConf(url, password string, db int) *RedisConf {
	return &RedisConf{
		url:        url,
		password:   password,
		db:         db,
		expireTime: expireTime,
		poolSize:   poolSize,
	}
}

func (r *RedisConf) SetExpireTime(expireTime time.Duration) {
	r.expireTime = expireTime
}

func (r *RedisConf) SetPoolSize(poolSize int) {
	r.poolSize = poolSize
}
