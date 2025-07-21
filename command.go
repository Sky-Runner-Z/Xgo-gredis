package gredis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// SET 命令的各种重载版本

// Set__0: set key value
func (r *RedisClient) Set__0(key string, value interface{}) *Result {
	cmd := r.client.Set(r.ctx, key, value, 0)
	return &Result{cmd: cmd}
}

// Set__1: set key value ex seconds
func (r *RedisClient) Set__1(key string, value interface{}, ex int) *Result {
	expiration := time.Duration(ex) * time.Second
	cmd := r.client.Set(r.ctx, key, value, expiration)
	return &Result{cmd: cmd}
}

// Set__2: set key value nx
func (r *RedisClient) Set__2(key string, value interface{}, nx bool) *Result {
	if nx {
		cmd := r.client.SetNX(r.ctx, key, value, 0)
		return &Result{cmd: cmd}
	}
	return r.Set__0(key, value)
}

// Set__3: set key value nx ex seconds
func (r *RedisClient) Set__3(key string, value interface{}, nx bool, ex int) *Result {
	expiration := time.Duration(ex) * time.Second
	if nx {
		cmd := r.client.SetNX(r.ctx, key, value, expiration)
		return &Result{cmd: cmd}
	}
	cmd := r.client.Set(r.ctx, key, value, expiration)
	return &Result{cmd: cmd}
}

// SetArgs 使用参数结构体的方式
func (r *RedisClient) SetArgs(key string, value interface{}, args *SetArgs) *Result {
	setArgs := &redis.SetArgs{
		Mode: redis.KeepTTL,
		TTL:  0,
	}

	if args != nil {
		if args.NX {
			setArgs.Mode = "NX"
		}
		if args.EX > 0 {
			setArgs.TTL = time.Duration(args.EX) * time.Second
		}
	}

	cmd := r.client.SetArgs(r.ctx, key, value, *setArgs)
	return &Result{cmd: cmd}
}

// GET 命令
func (r *RedisClient) Get(key string) *StringResult {
	cmd := r.client.Get(r.ctx, key)
	return &StringResult{cmd: cmd}
}

// DEL 命令
func (r *RedisClient) Del(keys ...string) *IntResult {
	cmd := r.client.Del(r.ctx, keys...)
	return &IntResult{cmd: cmd}
}

// EXISTS 命令
func (r *RedisClient) Exists(keys ...string) *IntResult {
	cmd := r.client.Exists(r.ctx, keys...)
	return &IntResult{cmd: cmd}
}

// TTL 命令
func (r *RedisClient) TTL(key string) *DurationResult {
	cmd := r.client.TTL(r.ctx, key)
	return &DurationResult{cmd: cmd}
}

// EXPIRE 命令
func (r *RedisClient) Expire(key string, expiration int) *BoolResult {
	cmd := r.client.Expire(r.ctx, key, time.Duration(expiration)*time.Second)
	return &BoolResult{cmd: cmd}
}
