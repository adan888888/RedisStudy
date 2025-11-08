# Redis Key 和 Value 类型说明

## 重要规则

### 1. Key 必须是字符串
- **Key 只能是字符串类型**，不能是数字、对象等其他类型
- 如果需要使用数字作为key，需要先转换为字符串

### 2. Value 可以是多种数据类型
Redis 支持以下 Value 类型：
- **string** - 字符串
- **hash** - 哈希（类似map）
- **list** - 列表
- **set** - 集合
- **zset** - 有序集合

## 详细说明

### Key 类型限制

```go
// ✅ 正确：key是字符串
rdb.Set(ctx, "user:1001", "value", 0)
rdb.Set(ctx, "123", "value", 0)  // 数字作为字符串

// ❌ 错误：不能直接使用数字作为key
// rdb.Set(ctx, 123, "value", 0)  // 编译错误

// ✅ 正确：将数字转换为字符串
userId := 1001
rdb.Set(ctx, fmt.Sprintf("user:%d", userId), "value", 0)
```

### Value 类型说明

#### 1. String 类型（字符串）
```go
// 存储字符串
rdb.Set(ctx, "name", "张三", 0)

// 存储数字（需要转换为字符串）
rdb.Set(ctx, "age", "25", 0)  // 字符串 "25"
rdb.Set(ctx, "age", strconv.Itoa(25), 0)  // 将int转为字符串

// 存储数字（使用Set方法，会自动转换）
rdb.Set(ctx, "count", 100, 0)  // Go客户端会自动转换为字符串
```

#### 2. Hash 类型（哈希）
```go
// Hash的field和value都必须是字符串
rdb.HSet(ctx, "user:1001", map[string]interface{}{
    "name":  "李四",      // ✅ 字符串
    "age":   "25",       // ✅ 字符串（数字需要转字符串）
    "score": "95.5",     // ✅ 字符串（浮点数需要转字符串）
})
```

#### 3. List 类型（列表）
```go
// List的元素必须是字符串
rdb.LPush(ctx, "tasks", "任务1", "任务2", "任务3")
```

#### 4. Set 类型（集合）
```go
// Set的成员必须是字符串
rdb.SAdd(ctx, "tags", "go", "redis", "golang")
```

#### 5. ZSet 类型（有序集合）
```go
// ZSet的member必须是字符串，score可以是浮点数
rdb.ZAdd(ctx, "scores", redis.Z{
    Score:  95.5,        // ✅ 浮点数（score可以是数字）
    Member: "张三",      // ✅ 字符串（member必须是字符串）
})
```

## 底层存储机制

**重要：Redis 底层存储时，所有值最终都是字符串形式**

- 数字 `123` 存储为字符串 `"123"`
- 浮点数 `95.5` 存储为字符串 `"95.5"`
- 布尔值 `true` 存储为字符串 `"true"` 或 `"1"`

## Go 客户端类型转换

### 存储数字

```go
// 方法1：手动转换为字符串
age := 25
rdb.Set(ctx, "age", strconv.Itoa(age), 0)

// 方法2：使用fmt.Sprintf
age := 25
rdb.Set(ctx, "age", fmt.Sprintf("%d", age), 0)

// 方法3：Go客户端会自动转换（推荐）
age := 25
rdb.Set(ctx, "age", age, 0)  // 客户端内部会转换为字符串

// 方法4：使用SetNX等方法的变体
rdb.Set(ctx, "count", 100, 0)  // 自动转换
```

### 读取数字

```go
// 读取字符串
val, _ := rdb.Get(ctx, "age").Result()  // 返回 "25" (字符串)

// 读取并转换为int
age, _ := rdb.Get(ctx, "age").Int()     // 返回 25 (int)

// 读取并转换为float64
score, _ := rdb.Get(ctx, "score").Float64()  // 返回 95.5 (float64)

// 读取并转换为bool
flag, _ := rdb.Get(ctx, "flag").Bool()   // 返回 true (bool)
```

## 实际示例

### 示例1：存储用户信息（包含数字）

```go
// 存储用户信息
userID := 1001
age := 25
score := 95.5

// 方法1：使用Hash存储（推荐用于对象）
rdb.HSet(ctx, fmt.Sprintf("user:%d", userID), map[string]interface{}{
    "id":    strconv.Itoa(userID),  // 必须转字符串
    "age":   strconv.Itoa(age),      // 必须转字符串
    "score": fmt.Sprintf("%.1f", score),  // 必须转字符串
})

// 方法2：使用多个String key存储
rdb.Set(ctx, fmt.Sprintf("user:%d:age", userID), age, 0)      // 自动转换
rdb.Set(ctx, fmt.Sprintf("user:%d:score", userID), score, 0) // 自动转换
```

### 示例2：存储和读取数字

```go
// 存储
count := 100
rdb.Set(ctx, "count", count, 0)

// 读取
// 方式1：读取为字符串
countStr, _ := rdb.Get(ctx, "count").Result()  // "100"

// 方式2：直接读取为int（推荐）
countInt, _ := rdb.Get(ctx, "count").Int()     // 100

// 方式3：读取为int64
countInt64, _ := rdb.Get(ctx, "count").Int64() // 100
```

### 示例3：Hash中的数字处理

```go
// 存储
rdb.HSet(ctx, "user:1001", map[string]interface{}{
    "age":   "25",      // Hash的value必须是字符串
    "score": "95.5",    // Hash的value必须是字符串
})

// 读取
ageStr, _ := rdb.HGet(ctx, "user:1001", "age").Result()  // "25" (字符串)
ageInt, _ := rdb.HGet(ctx, "user:1001", "age").Int()     // 25 (int)
```

## 总结

| 项目 | 类型要求 | 说明 |
|------|---------|------|
| **Key** | 必须是字符串 | 不能是数字、对象等 |
| **String Value** | 字符串 | 数字需要转换为字符串 |
| **Hash Field** | 字符串 | 字段名必须是字符串 |
| **Hash Value** | 字符串 | 字段值必须是字符串 |
| **List Element** | 字符串 | 列表元素必须是字符串 |
| **Set Member** | 字符串 | 集合成员必须是字符串 |
| **ZSet Member** | 字符串 | 成员必须是字符串 |
| **ZSet Score** | 浮点数 | 分数可以是数字 |

**关键点：**
1. Key 必须是字符串
2. 所有 Value 在底层都是字符串存储
3. Go 客户端提供了便捷的类型转换方法（如 `.Int()`, `.Float64()`）
4. Hash 的 field 和 value 都必须是字符串
5. 数字需要显式转换为字符串（或使用客户端自动转换）

