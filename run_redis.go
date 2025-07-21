package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		Password: "",               // 密码，如果没有设置则为空
		DB:       0,                // 使用默认数据库
	})

	// 创建上下文
	ctx := context.Background()

	// 测试连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("连接Redis失败:", err)
	}
	fmt.Println("连接Redis成功:", pong)

	// 基本操作示例

	// 1. 设置键值对
	err = rdb.Set(ctx, "name", "张三", 0).Err()
	if err != nil {
		log.Fatal("设置键值失败:", err)
	}
	fmt.Println("设置 name = 张三")

	// 2. 获取值
	val, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		log.Fatal("获取键值失败:", err)
	}
	fmt.Println("获取 name =", val)

	// 3. 设置带过期时间的键值对
	err = rdb.Set(ctx, "session", "abc123", 30*time.Second).Err()
	if err != nil {
		log.Fatal("设置带过期时间的键值失败:", err)
	}
	fmt.Println("设置 session = abc123 (30秒后过期)")

	// 4. 检查键是否存在
	exists, err := rdb.Exists(ctx, "name").Result()
	if err != nil {
		log.Fatal("检查键存在性失败:", err)
	}
	fmt.Printf("键 'name' 是否存在: %t\n", exists == 1)

	// 5. 删除键
	err = rdb.Del(ctx, "name").Err()
	if err != nil {
		log.Fatal("删除键失败:", err)
	}
	fmt.Println("删除键 'name'")

	// 6. 列表操作
	err = rdb.LPush(ctx, "mylist", "item1", "item2", "item3").Err()
	if err != nil {
		log.Fatal("列表推送失败:", err)
	}

	listItems, err := rdb.LRange(ctx, "mylist", 0, -1).Result()
	if err != nil {
		log.Fatal("获取列表失败:", err)
	}
	fmt.Println("列表内容:", listItems)

	// 7. 哈希操作
	err = rdb.HSet(ctx, "user:1", "name", "李四", "age", "25", "city", "北京").Err()
	if err != nil {
		log.Fatal("设置哈希失败:", err)
	}

	userInfo, err := rdb.HGetAll(ctx, "user:1").Result()
	if err != nil {
		log.Fatal("获取哈希失败:", err)
	}
	fmt.Println("用户信息:", userInfo)

	// 8. 集合操作
	err = rdb.SAdd(ctx, "tags", "go", "redis", "database", "nosql").Err()
	if err != nil {
		log.Fatal("添加集合元素失败:", err)
	}

	members, err := rdb.SMembers(ctx, "tags").Result()
	if err != nil {
		log.Fatal("获取集合成员失败:", err)
	}
	fmt.Println("标签集合:", members)

	// 关闭连接
	err = rdb.Close()
	if err != nil {
		log.Fatal("关闭Redis连接失败:", err)
	}
	fmt.Println("已关闭Redis连接")
}
