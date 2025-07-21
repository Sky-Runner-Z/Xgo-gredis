package gredis

import (
	"fmt"
	"log"
	"reflect"
)

// 帮助函数：创建 Redis 客户端实例
func Connect(addr string) *RedisClient {
	return NewClient(addr, "", 0)
}

func ConnectWithAuth(addr, password string) *RedisClient {
	return NewClient(addr, password, 0)
}

func ConnectWithDB(addr string, db int) *RedisClient {
	return NewClient(addr, "", db)
}

func ConnectFull(addr, password string, db int) *RedisClient {
	return NewClient(addr, password, db)
}

// 创建 SetArgs 的便捷函数
func NX() *SetArgs {
	return &SetArgs{NX: true}
}

func EX(seconds int) *SetArgs {
	return &SetArgs{EX: seconds}
}

func NXEX(seconds int) *SetArgs {
	return &SetArgs{NX: true, EX: seconds}
}

// 实例化函数（参考 gplot 的实现）
func instance(conn reflect.Value) *RedisClient {
	fld := conn.FieldByName("Redis")
	if !fld.IsValid() {
		log.Panicf("type %v doesn't have field redis.Redis", conn.Type())
	}
	return fld.Addr().Interface().(*RedisClient)
}

// 错误处理辅助函数
func CheckError(err error) {
	if err != nil {
		fmt.Printf("Redis error: %v\n", err)
	}
}

func MustGet(result *StringResult) string {
	val, err := result.Result()
	if err != nil {
		panic("Redis GET failed: " + err.Error())
	}
	return val
}

func MustSet(result *Result) {
	err := result.Err()
	if err != nil {
		panic("Redis SET failed: " + err.Error())
	}
}
