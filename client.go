package gredis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client 是 Redis 操作的主体（类似 Figure）
type Client struct {
	connections map[string]*Connection
	current     string // 当前使用的连接名
	ctx         context.Context
	results     []interface{} // 存储操作结果
}

// Connection 表示单个 Redis 连接（类似 Axis）
type Connection struct {
	client *redis.Client
	config ConnectionConfig
	name   string
}

func NewConnection(config ConnectionConfig) *Connection {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	return &Connection{
		client: rdb,
		config: config,
		name:   config.Name,
	}
}

// Client 的生命周期管理方法
func (c *Client) initRedis() {
	c.connections = make(map[string]*Connection)
	c.ctx = context.Background()
	c.results = make([]interface{}, 0)

	// 创建默认连接
	defaultConfig := ConnectionConfig{
		Name:     "default",
		Host:     "14.103.115.251",
		Port:     6379,
		Password: "200228",
		DB:       0,
		PoolSize: 10,
	}
	c.connections["default"] = NewConnection(defaultConfig)
	c.current = "default"
}

func (c *Client) finishRedis() {
	// 关闭所有连接
	for _, conn := range c.connections {
		if conn.client != nil {
			conn.client.Close()
		}
	}

	// 打印操作结果（可选）
	if len(c.results) > 0 {
		fmt.Println("Redis操作结果:")
		for i, result := range c.results {
			fmt.Printf("操作 %d: %v\n", i+1, result)
		}
	}
}

// 连接管理方法
func (c *Client) AddConnection(name string, conn *Connection) {
	c.connections[name] = conn
}

func (c *Client) UseConnection(name string) {
	if _, exists := c.connections[name]; exists {
		c.current = name
	} else {
		log.Printf("连接 %s 不存在", name)
	}
}

func (c *Client) getCurrentConnection() *Connection {
	return c.connections[c.current]
}

// Redis 基本操作方法（委托给当前连接）
func (c *Client) Set(key string, value interface{}, expiration time.Duration) {
	conn := c.getCurrentConnection()
	if conn == nil {
		log.Printf("没有可用的连接")
		return
	}

	result := conn.client.Set(c.ctx, key, value, expiration)
	c.results = append(c.results, fmt.Sprintf("SET %s: %v", key, result.Val()))
}

func (c *Client) Get(key string) string {
	conn := c.getCurrentConnection()
	if conn == nil {
		log.Printf("没有可用的连接")
		return ""
	}

	result := conn.client.Get(c.ctx, key)
	value := result.Val()
	c.results = append(c.results, fmt.Sprintf("GET %s: %s", key, value))
	return value
}

func (c *Client) Del(keys ...string) int64 {
	conn := c.getCurrentConnection()
	if conn == nil {
		log.Printf("没有可用的连接")
		return 0
	}

	result := conn.client.Del(c.ctx, keys...)
	count := result.Val()
	c.results = append(c.results, fmt.Sprintf("DEL %v: %d keys deleted", keys, count))
	return count
}

func (c *Client) Expire(key string, expiration time.Duration) bool {
	conn := c.getCurrentConnection()
	if conn == nil {
		log.Printf("没有可用的连接")
		return false
	}

	result := conn.client.Expire(c.ctx, key, expiration)
	success := result.Val()
	c.results = append(c.results, fmt.Sprintf("EXPIRE %s: %v", key, success))
	return success
}

// 高级操作
func (c *Client) Pipeline(commands ...*RedisCommand) []interface{} {
	conn := c.getCurrentConnection()
	if conn == nil {
		log.Printf("没有可用的连接")
		return nil
	}

	pipe := conn.client.Pipeline()
	for _, cmd := range commands {
		switch cmd.Operation {
		case "SET":
			pipe.Set(c.ctx, cmd.Key, cmd.Value, cmd.TTL)
		case "GET":
			pipe.Get(c.ctx, cmd.Key)
		case "DEL":
			pipe.Del(c.ctx, cmd.Keys...)
		case "EXPIRE":
			pipe.Expire(c.ctx, cmd.Key, cmd.TTL)
		}
	}

	results, err := pipe.Exec(c.ctx)
	if err != nil {
		log.Printf("Pipeline执行错误: %v", err)
		return nil
	}

	var values []interface{}
	for _, result := range results {
		values = append(values, result)
	}

	c.results = append(c.results, fmt.Sprintf("Pipeline执行了 %d 个命令", len(commands)))
	return values
}

func (c *Client) Transaction(commands ...*RedisCommand) []interface{} {
	conn := c.getCurrentConnection()
	if conn == nil {
		log.Printf("没有可用的连接")
		return nil
	}

	results, err := conn.client.TxPipelined(c.ctx, func(pipe redis.Pipeliner) error {
		for _, cmd := range commands {
			switch cmd.Operation {
			case "SET":
				pipe.Set(c.ctx, cmd.Key, cmd.Value, cmd.TTL)
			case "GET":
				pipe.Get(c.ctx, cmd.Key)
			case "DEL":
				pipe.Del(c.ctx, cmd.Keys...)
			case "EXPIRE":
				pipe.Expire(c.ctx, cmd.Key, cmd.TTL)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("事务执行错误: %v", err)
		return nil
	}

	var values []interface{}
	for _, result := range results {
		values = append(values, result)
	}

	c.results = append(c.results, fmt.Sprintf("事务执行了 %d 个命令", len(commands)))
	return values
}

// Hash 操作
func (c *Client) HSet(key, field string, value interface{}) {
	conn := c.getCurrentConnection()
	if conn == nil {
		return
	}

	result := conn.client.HSet(c.ctx, key, field, value)
	c.results = append(c.results, fmt.Sprintf("HSET %s %s: %v", key, field, result.Val()))
}

func (c *Client) HGet(key, field string) string {
	conn := c.getCurrentConnection()
	if conn == nil {
		return ""
	}

	result := conn.client.HGet(c.ctx, key, field)
	value := result.Val()
	c.results = append(c.results, fmt.Sprintf("HGET %s %s: %s", key, field, value))
	return value
}

// List 操作
func (c *Client) LPush(key string, values ...interface{}) int64 {
	conn := c.getCurrentConnection()
	if conn == nil {
		return 0
	}

	result := conn.client.LPush(c.ctx, key, values...)
	length := result.Val()
	c.results = append(c.results, fmt.Sprintf("LPUSH %s: list length %d", key, length))
	return length
}

func (c *Client) LPop(key string) string {
	conn := c.getCurrentConnection()
	if conn == nil {
		return ""
	}

	result := conn.client.LPop(c.ctx, key)
	value := result.Val()
	c.results = append(c.results, fmt.Sprintf("LPOP %s: %s", key, value))
	return value
}

// Set 操作
func (c *Client) SAdd(key string, members ...interface{}) int64 {
	conn := c.getCurrentConnection()
	if conn == nil {
		return 0
	}

	result := conn.client.SAdd(c.ctx, key, members...)
	count := result.Val()
	c.results = append(c.results, fmt.Sprintf("SADD %s: %d members added", key, count))
	return count
}

func (c *Client) SMembers(key string) []string {
	conn := c.getCurrentConnection()
	if conn == nil {
		return nil
	}

	result := conn.client.SMembers(c.ctx, key)
	members := result.Val()
	c.results = append(c.results, fmt.Sprintf("SMEMBERS %s: %d members", key, len(members)))
	return members
}

// 连接状态检查
func (c *Client) Ping() string {
	conn := c.getCurrentConnection()
	if conn == nil {
		return "PONG"
	}

	result := conn.client.Ping(c.ctx)
	response := result.Val()
	c.results = append(c.results, fmt.Sprintf("PING: %s", response))
	return response
}

// 获取信息
func (c *Client) Info(section ...string) string {
	conn := c.getCurrentConnection()
	if conn == nil {
		return ""
	}

	result := conn.client.Info(c.ctx, section...)
	info := result.Val()
	c.results = append(c.results, "INFO命令已执行")
	return info
}
