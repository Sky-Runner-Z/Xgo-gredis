package gredis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const (
	GopPackage = true
	Gop_game   = "Redis"
	Gop_sprite = "Client"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

type Connecter interface {
	InitRedis()
	FinishRedis()
}

// Gopt_Redis_Main 是 Go+ 编译器识别 .redis 项目的入口
func Gopt_Redis_Main(conn Connecter) {
	conn.InitRedis()
	defer conn.FinishRedis()
	conn.(interface{ MainEntry() }).MainEntry()
}

func NewClient(addr string, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisClient) initRedis() {
	// 连接测试
	_, err := r.client.Ping(r.ctx).Result()
	if err != nil {
		panic("Redis connection failed: " + err.Error())
	}
}

func (r *RedisClient) finishRedis() {
	if r.client != nil {
		r.client.Close()
	}
}
