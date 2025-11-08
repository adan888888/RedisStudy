package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		Password: "",                // 密码，如果没有设置密码则为空
		DB:       0,                 // 使用默认数据库
	})

	// 确保在程序结束时关闭连接
	defer rdb.Close()

	// 创建上下文
	ctx := context.Background()

	// 测试连接
	if !testConnection(rdb, ctx) {
		return
	}

	// 执行各种Redis操作
	//1.stringOperations - 字符串操作（SET、GET、带过期时间的SET、检查key存在）
	// stringOperations(rdb, ctx)
	//2.hashOperations - 哈希操作（HSET、HGET、HGETALL）
	hashOperations(rdb, ctx)
	//3.listOperations - 列表操作（LPUSH、LRANGE、RPOP）
	// listOperations(rdb, ctx)
	// //4.setOperations - 集合操作（SADD、SMEMBERS、SISMEMBER）
	// setOperations(rdb, ctx)
	// //5.sortedSetOperations - 有序集合操作（ZADD、ZRANGE、ZREVRANGE）
	// sortedSetOperations(rdb, ctx)
	// //6.keyOperations - 键操作（KEYS、DEL、EXPIRE、TTL）			
	// keyOperations(rdb, ctx)
	// //7.pipelineOperations - 管道操作（批量执行）
	// pipelineOperations(rdb, ctx)
	// //8.transactionOperations - 事务操作（Watch + Multi + Exec）
	// transactionOperations(rdb, ctx)

	fmt.Println("\n========== Redis操作示例完成 ==========")
}

// testConnection 测试Redis连接
func testConnection(rdb *redis.Client, ctx context.Context) bool {
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("连接Redis失败: %v\n", err)
		return false
	}
	fmt.Printf("Redis连接成功: %s\n", pong)
	return true
}

// stringOperations 字符串操作示例
func stringOperations(rdb *redis.Client, ctx context.Context) {
	fmt.Println("\n========== 字符串操作 ==========")

	// SET操作
	err := rdb.Set(ctx, "name", "张三", 0).Err() // 设置key为name，值为张三，过期时间为0，即永不过期
	if err != nil {
		fmt.Printf("SET操作失败: %v\n", err)
	} else {
		fmt.Println("SET name = 张三")
	}

	// GET操作
	val, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		fmt.Printf("GET操作失败: %v\n", err)
	} else {
		fmt.Printf("GET name = %s\n", val)
	}

	// SET带过期时间（10秒）
	err = rdb.Set(ctx, "token", "abc123xyz", 10*time.Second).Err()
	if err != nil {
		fmt.Printf("SET带过期时间失败: %v\n", err)
	} else {
		fmt.Println("SET token = abc123xyz (10秒后过期)")
	}

	// 检查key是否存在
	exists, err := rdb.Exists(ctx, "name").Result()
	if err != nil {
		fmt.Printf("检查key存在性失败: %v\n", err)
	} else {
		fmt.Printf("key 'name' 是否存在: %d (1=存在, 0=不存在)\n", exists)
	}
}

// hashOperations 哈希操作示例
func hashOperations(rdb *redis.Client, ctx context.Context) {
	fmt.Println("\n========== 哈希操作 ==========")

	// HSET操作
	err := rdb.HSet(ctx, "user:1001", map[string]interface{}{
		"name":  "李四",
		"age":   "25",
		"email": "lisi@example.com",
	}).Err()
	if err != nil {
		fmt.Printf("HSET操作失败: %v\n", err)
	} else {
		fmt.Println("HSET user:1001 {name: 李四, age: 25, email: lisi@example.com}")
	}

	// HGET操作
	name, err := rdb.HGet(ctx, "user:1001", "name").Result()
	if err != nil {
		fmt.Printf("HGET操作失败: %v\n", err)
	} else {
		fmt.Printf("HGET user:1001 name = %s\n", name)
	}

	// HGETALL操作
	userInfo, err := rdb.HGetAll(ctx, "user:1001").Result()
	if err != nil {
		fmt.Printf("HGETALL操作失败: %v\n", err)
	} else {
		fmt.Printf("HGETALL user:1001 = %v\n", userInfo)
	}
}

// listOperations 列表操作示例
func listOperations(rdb *redis.Client, ctx context.Context) {
	fmt.Println("\n========== 列表操作 ==========")

	// LPUSH操作（从左边推入）
	err := rdb.LPush(ctx, "tasks", "任务1", "任务2", "任务3").Err()
	if err != nil {
		fmt.Printf("LPUSH操作失败: %v\n", err)
	} else {
		fmt.Println("LPUSH tasks [任务1, 任务2, 任务3]")
	}

	// LRANGE操作（获取列表范围）
	tasks, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		fmt.Printf("LRANGE操作失败: %v\n", err)
	} else {
		fmt.Printf("LRANGE tasks 0 -1 = %v\n", tasks)
	}

	// RPOP操作（从右边弹出）
	task, err := rdb.RPop(ctx, "tasks").Result()
	if err != nil {
		fmt.Printf("RPOP操作失败: %v\n", err)
	} else {
		fmt.Printf("RPOP tasks = %s\n", task)
	}
}

// setOperations 集合操作示例
func setOperations(rdb *redis.Client, ctx context.Context) {
	fmt.Println("\n========== 集合操作 ==========")

	// SADD操作（添加成员）
	err := rdb.SAdd(ctx, "tags", "go", "redis", "golang", "database").Err()
	if err != nil {
		fmt.Printf("SADD操作失败: %v\n", err)
	} else {
		fmt.Println("SADD tags [go, redis, golang, database]")
	}

	// SMEMBERS操作（获取所有成员）
	tags, err := rdb.SMembers(ctx, "tags").Result()
	if err != nil {
		fmt.Printf("SMEMBERS操作失败: %v\n", err)
	} else {
		fmt.Printf("SMEMBERS tags = %v\n", tags)
	}

	// SISMEMBER操作（检查成员是否存在）
	isMember, err := rdb.SIsMember(ctx, "tags", "go").Result()
	if err != nil {
		fmt.Printf("SISMEMBER操作失败: %v\n", err)
	} else {
		fmt.Printf("SISMEMBER tags go = %v\n", isMember)
	}
}

// sortedSetOperations 有序集合操作示例
func sortedSetOperations(rdb *redis.Client, ctx context.Context) {
	fmt.Println("\n========== 有序集合操作 ==========")

	// ZADD操作（添加带分数的成员）
	err := rdb.ZAdd(ctx, "scores", redis.Z{
		Score:  95.5,
		Member: "张三",
	}, redis.Z{
		Score:  88.0,
		Member: "李四",
	}, redis.Z{
		Score:  92.5,
		Member: "王五",
	}).Err()
	if err != nil {
		fmt.Printf("ZADD操作失败: %v\n", err)
	} else {
		fmt.Println("ZADD scores [(95.5, 张三), (88.0, 李四), (92.5, 王五)]")
	}

	// ZRANGE操作（按分数范围获取，从低到高）
	topScores, err := rdb.ZRangeWithScores(ctx, "scores", 0, -1).Result()
	if err != nil {
		fmt.Printf("ZRANGE操作失败: %v\n", err)
	} else {
		fmt.Println("ZRANGE scores 0 -1 (按分数从低到高):")
		for _, z := range topScores {
			fmt.Printf("  %s: %.1f\n", z.Member, z.Score)
		}
	}

	// ZREVRANGE操作（按分数范围获取，从高到低）
	revScores, err := rdb.ZRevRangeWithScores(ctx, "scores", 0, -1).Result()
	if err != nil {
		fmt.Printf("ZREVRANGE操作失败: %v\n", err)
	} else {
		fmt.Println("ZREVRANGE scores 0 -1 (按分数从高到低):")
		for _, z := range revScores {
			fmt.Printf("  %s: %.1f\n", z.Member, z.Score)
		}
	}
}

// keyOperations 键操作示例
func keyOperations(rdb *redis.Client, ctx context.Context) {
	fmt.Println("\n========== 键操作 ==========")

	// 获取所有匹配的key
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		fmt.Printf("KEYS操作失败: %v\n", err)
	} else {
		fmt.Printf("KEYS * = %v (共%d个key)\n", keys, len(keys))
	}

	// 删除key
	deleted, err := rdb.Del(ctx, "name").Result()
	if err != nil {
		fmt.Printf("DEL操作失败: %v\n", err)
	} else {
		fmt.Printf("DEL name (删除了%d个key)\n", deleted)
	}

	// 设置key的过期时间
	err = rdb.Expire(ctx, "user:1001", 60*time.Second).Err()
	if err != nil {
		fmt.Printf("EXPIRE操作失败: %v\n", err)
	} else {
		fmt.Println("EXPIRE user:1001 60 (设置60秒过期)")
	}

	// 获取key的剩余过期时间
	ttl, err := rdb.TTL(ctx, "user:1001").Result()
	if err != nil {
		fmt.Printf("TTL操作失败: %v\n", err)
	} else {
		fmt.Printf("TTL user:1001 = %v\n", ttl)
	}
}

// pipelineOperations 管道操作示例（批量执行）
func pipelineOperations(rdb *redis.Client, ctx context.Context) {
	fmt.Println("\n========== 管道操作（批量执行）==========")

	pipe := rdb.Pipeline()
	pipe.Set(ctx, "key1", "value1", 0)
	pipe.Set(ctx, "key2", "value2", 0)
	pipe.Set(ctx, "key3", "value3", 0)
	pipe.Get(ctx, "key1")
	pipe.Get(ctx, "key2")

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		fmt.Printf("管道操作失败: %v\n", err)
	} else {
		fmt.Printf("管道执行成功，共执行了%d个命令\n", len(cmds))
		// 获取管道中的结果
		for i, cmd := range cmds {
			if i >= 3 { // 前3个是SET命令，后2个是GET命令
				if getCmd, ok := cmd.(*redis.StringCmd); ok {
					val, _ := getCmd.Result()
					fmt.Printf("  管道GET结果: %s\n", val)
				}
			}
		}
	}
}

// transactionOperations 事务操作示例
func transactionOperations(rdb *redis.Client, ctx context.Context) {
	fmt.Println("\n========== 事务操作 ==========")

	// 使用Watch + Multi + Exec实现事务
	key := "balance"
	err := rdb.Set(ctx, key, "100", 0).Err()
	if err != nil {
		fmt.Printf("设置初始值失败: %v\n", err)
	} else {
		fmt.Printf("设置初始值 %s = 100\n", key)
	}

	// Watch监听key
	err = rdb.Watch(ctx, func(tx *redis.Tx) error {
		// 获取当前值
		balance, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		// 开始事务
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// 在事务中执行操作
			pipe.Set(ctx, key, balance+50, 0)
			return nil
		})
		return err
	}, key)

	if err != nil {
		fmt.Printf("事务操作失败: %v\n", err)
	} else {
		finalBalance, _ := rdb.Get(ctx, key).Int()
		fmt.Printf("事务执行成功，%s = %d\n", key, finalBalance)
	}
}