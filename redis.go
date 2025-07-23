package gredis

import (
	"log"
	"reflect"
	"time"
)

const (
	GopPackage = true
	Gop_game   = "Client"     // 主要操作对象是Redis客户端
	Gop_sprite = "Connection" // 子对象是连接
)

type RedisOperator interface {
	initRedis()
	finishRedis()
}

// Gopt_Client_Main is required by Go+ compiler as the entry of a .redis project.
func Gopt_Client_Main(operator RedisOperator) {
	operator.initRedis()
	defer operator.finishRedis()
	operator.(interface{ MainEntry() }).MainEntry()
}

// Gopt_Client_Run is required by Go+ compiler for multi-connection operations.
func Gopt_Client_Run(operator RedisOperator, configs ...ConnectionConfig) {
	v := reflect.ValueOf(operator).Elem()
	t := reflect.TypeOf(operator).Elem()
	client := instance(v)

	for _, config := range configs {
		for j, n := 0, v.NumField(); j < n; j++ {
			field := t.Field(j)
			typ := field.Type
			m, ok := reflect.PtrTo(typ).MethodByName("Main")
			if ok {
				parent, conn := instanceConnection(field, config)
				m.Func.Call([]reflect.Value{parent})
				client.AddConnection(config.Name, conn)
			}
		}
	}
}

// Redis utility functions
func Set(key string, value interface{}, expiration time.Duration) *RedisCommand {
	return &RedisCommand{
		Operation: "SET",
		Key:       key,
		Value:     value,
		TTL:       expiration * time.Second,
	}
}

func Get(key string) *RedisCommand {
	return &RedisCommand{
		Operation: "GET",
		Key:       key,
	}
}

func Del(keys ...string) *RedisCommand {
	return &RedisCommand{
		Operation: "DEL",
		Keys:      keys,
	}
}

func Expire(key string, expiration time.Duration) *RedisCommand {
	return &RedisCommand{
		Operation: "EXPIRE",
		Key:       key,
		TTL:       expiration * time.Second,
	}
}

// Pipeline operations
func Pipeline(commands ...*RedisCommand) *RedisPipeline {
	return &RedisPipeline{
		Commands: commands,
	}
}

// Transaction operations
func Multi(commands ...*RedisCommand) *RedisTransaction {
	return &RedisTransaction{
		Commands: commands,
	}
}

func instance(operator reflect.Value) *Client {
	fld := operator.FieldByName("Client")
	if !fld.IsValid() {
		log.Panicf("type %v doesn't have field gredis.Client", operator.Type())
	}
	return fld.Addr().Interface().(*Client)
}

func instanceConnection(field reflect.StructField, config ConnectionConfig) (reflect.Value, *Connection) {
	typ := field.Type
	parent := reflect.New(typ)
	conn := parent.Elem().FieldByName("Connection")
	newConn := NewConnection(config)
	conn.Set(reflect.ValueOf(newConn).Elem())
	return parent, newConn
}

type ConnectionConfig struct {
	Name     string
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type RedisCommand struct {
	Operation string
	Key       string
	Keys      []string
	Value     interface{}
	TTL       time.Duration
	Result    interface{}
	Error     error
}

type RedisPipeline struct {
	Commands []*RedisCommand
	Results  []interface{}
	Error    error
}

type RedisTransaction struct {
	Commands []*RedisCommand
	Results  []interface{}
	Error    error
}
