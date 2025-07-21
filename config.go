package gredis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// SetArgs 参数结构体
type SetArgs struct {
	NX bool // 只有当 key 不存在时才设置
	EX int  // 过期时间（秒）
}

// 结果类型封装
type Result struct {
	cmd *redis.StatusCmd
}

func (r *Result) Result() (string, error) {
	return r.cmd.Result()
}

func (r *Result) Val() string {
	return r.cmd.Val()
}

func (r *Result) Err() error {
	return r.cmd.Err()
}

type StringResult struct {
	cmd *redis.StringCmd
}

func (r *StringResult) Result() (string, error) {
	return r.cmd.Result()
}

func (r *StringResult) Val() string {
	return r.cmd.Val()
}

func (r *StringResult) Err() error {
	return r.cmd.Err()
}

type IntResult struct {
	cmd *redis.IntCmd
}

func (r *IntResult) Result() (int64, error) {
	return r.cmd.Result()
}

func (r *IntResult) Val() int64 {
	return r.cmd.Val()
}

func (r *IntResult) Err() error {
	return r.cmd.Err()
}

type BoolResult struct {
	cmd *redis.BoolCmd
}

func (r *BoolResult) Result() (bool, error) {
	return r.cmd.Result()
}

func (r *BoolResult) Val() bool {
	return r.cmd.Val()
}

func (r *BoolResult) Err() error {
	return r.cmd.Err()
}

type DurationResult struct {
	cmd *redis.DurationCmd
}

func (r *DurationResult) Result() (time.Duration, error) {
	return r.cmd.Result()
}

func (r *DurationResult) Val() time.Duration {
	return r.cmd.Val()
}

func (r *DurationResult) Err() error {
	return r.cmd.Err()
}
