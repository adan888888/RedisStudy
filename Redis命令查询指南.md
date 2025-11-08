# Redis命令查询指南

## 查询 user:1001 的常用命令

### 1. 连接到Redis

```bash
redis-cli
```

### 2. 基本查询命令

#### 检查key是否存在
```bash
redis-cli EXISTS user:1001
```
返回：`1` 表示存在，`0` 表示不存在

#### 查看key的类型
```bash
redis-cli TYPE user:1001
```
返回：`hash`（哈希类型）

#### 查看key的过期时间（TTL命令）
```bash
redis-cli TTL user:1001
```

**返回值说明：**
- `(integer) -2` - **key不存在**
- `(integer) -1` - key存在，但**永不过期**
- `(integer) 正数` - key存在，返回**剩余过期时间（秒）**

**详细示例：**

```bash
# 示例1：key不存在
redis-cli TTL nonexistent_key
# 返回: (integer) -2  ← key不存在

# 示例2：key存在，但永不过期
redis-cli SET name "张三"
redis-cli TTL name
# 返回: (integer) -1  ← 永不过期

# 示例3：key存在，有过期时间（60秒）
redis-cli SETEX token 60 "abc123"
redis-cli TTL token
# 返回: (integer) 60  ← 还有60秒过期

# 等待几秒后再次查看
redis-cli TTL token
# 返回: (integer) 55  ← 还有55秒过期

# 示例4：查看哈希类型key的过期时间
redis-cli HSET user:1001 name "李四" age "25"
redis-cli EXPIRE user:1001 120
redis-cli TTL user:1001
# 返回: (integer) 120  ← 还有120秒过期
```

**常见问题：**

1. **返回 -2 怎么办？**
   - 说明key不存在，需要先创建key
   ```bash
   # 先创建key
   redis-cli HSET user:1001 name "李四"
   # 然后再查看TTL
   redis-cli TTL user:1001
   ```

2. **返回 -1 是什么意思？**
   - key存在但没有设置过期时间，永不过期
   - 如果想设置过期时间，使用 `EXPIRE key seconds`

3. **如何设置过期时间？**
   ```bash
   # 设置60秒后过期
   redis-cli EXPIRE user:1001 60
   
   # 或者在创建时设置过期时间
   redis-cli SETEX keyname 60 "value"
   ```

**一行命令执行：**
```bash
redis-cli TTL user:1001
```

**Go代码对应：**
```go
// Go代码
ttl, err := rdb.TTL(ctx, "user:1001").Result()
if err != nil {
    fmt.Printf("TTL操作失败: %v\n", err)
} else {
    if ttl == -1 {
        fmt.Println("过期时间: 永不过期")
    } else if ttl == -2 {
        fmt.Println("过期时间: key不存在")
    } else {
        fmt.Printf("过期时间: %v (剩余 %d 秒)\n", ttl, int(ttl.Seconds()))
    }
}

// 对应Redis命令: TTL user:1001
```

#### 获取字符串类型key的值（GET命令）
```bash
redis-cli GET keyname
```
**说明：**
- 用于获取字符串类型（string）的key的值
- 如果key不存在，返回 `(nil)`
- 如果key不是字符串类型，返回错误

**示例：**
```bash
# 设置一个字符串类型的key
redis-cli SET name "张三"

# 获取值
redis-cli GET name
# 返回: "张三"

# 获取不存在的key
redis-cli GET nonexistent
# 返回: (nil)

# 如果key是其他类型（如hash、list、set、zset），使用GET会报错
# 错误: WRONGTYPE Operation against a key holding the wrong kind of value
```

**常见错误处理：**

1. **错误：`WRONGTYPE Operation against a key holding the wrong kind of value`**
   
   这个错误表示你使用了错误的命令类型。例如：
   ```bash
   # 错误示例：scores是zset（有序集合）类型，不能用GET命令
   redis-cli GET "scores"
   # 错误: WRONGTYPE Operation against a key holding the wrong kind of value
   ```

   **解决方法：**
   ```bash
   # 步骤1：先查看key的类型
   redis-cli TYPE "scores"
   # 返回: zset（说明是有序集合类型）
   
   # 步骤2：根据类型使用正确的命令
   # 对于zset类型，应该使用ZRANGE或ZREVRANGE
   redis-cli ZRANGE "scores" 0 -1 WITHSCORES
   
   # 对于hash类型，应该使用HGET或HGETALL
   redis-cli HGETALL "user:1001"
   
   # 对于list类型，应该使用LRANGE
   redis-cli LRANGE "tasks" 0 -1
   
   # 对于set类型，应该使用SMEMBERS
   redis-cli SMEMBERS "tags"
   ```

2. **类型与命令对应表：**
   | key类型 | 正确的查询命令 | 错误的命令示例 |
   |---------|---------------|---------------|
   | string | `GET keyname` | - |
   | hash | `HGET keyname field` 或 `HGETALL keyname` | `GET keyname` ❌ |
   | list | `LRANGE keyname 0 -1` | `GET keyname` ❌ |
   | set | `SMEMBERS keyname` | `GET keyname` ❌ |
   | zset | `ZRANGE keyname 0 -1 WITHSCORES` | `GET keyname` ❌ |

3. **实际案例：**
   ```bash
   # 案例1：scores是有序集合（zset）
   redis-cli TYPE "scores"
   # 返回: zset
   
   # 错误操作
   redis-cli GET "scores"
   # 错误: WRONGTYPE Operation against a key holding the wrong kind of value
   
   # 正确操作
   redis-cli ZRANGE "scores" 0 -1 WITHSCORES
   # 返回: 所有成员和分数
   
   # 案例2：user:1001是哈希（hash）
   redis-cli TYPE "user:1001"
   # 返回: hash
   
   # 错误操作
   redis-cli GET "user:1001"
   # 错误: WRONGTYPE Operation against a key holding the wrong kind of value
   
   # 正确操作
   redis-cli HGETALL "user:1001"
   # 返回: 所有字段和值
   ```

**一行命令执行：**
```bash
redis-cli GET name
```

**带密码连接：**
```bash
redis-cli -a yourpassword GET name
```

**连接远程服务器：**
```bash
redis-cli -h hostname -p port GET name
```

**Go代码对应：**
```go
// Go代码
val, err := rdb.Get(ctx, "name").Result()
if err != nil {
    fmt.Printf("GET操作失败: %v\n", err)
} else {
    fmt.Printf("GET name = %s\n", val)
}

// 对应Redis命令: GET name
```

**Go代码说明：**
- `rdb.Get(ctx, "name")` - 对应 Redis 命令 `GET name`
- `.Result()` - 执行命令并获取结果
- 返回两个值：`(string, error)`
  - 第一个值：key的值（字符串）
  - 第二个值：错误信息（如果key不存在，返回 `redis.Nil`）

**其他Go获取方式：**
```go
// 方式1：获取字符串结果
val, err := rdb.Get(ctx, "name").Result()

// 方式2：获取整数结果（自动转换）
age, err := rdb.Get(ctx, "age").Int()

// 方式3：获取浮点数结果（自动转换）
score, err := rdb.Get(ctx, "score").Float64()

// 方式4：获取布尔值结果（自动转换）
flag, err := rdb.Get(ctx, "flag").Bool()

// 方式5：检查key是否存在
val, err := rdb.Get(ctx, "name").Result()
if err == redis.Nil {
    fmt.Println("key不存在")
} else if err != nil {
    fmt.Printf("错误: %v\n", err)
} else {
    fmt.Printf("值: %s\n", val)
}
```

### 3. 哈希类型操作（user:1001是哈希类型）

#### 查询所有字段和值
```bash
redis-cli HGETALL user:1001
```
返回所有字段和对应的值

#### 查询单个字段
```bash
redis-cli HGET user:1001 name
redis-cli HGET user:1001 age
redis-cli HGET user:1001 email
```

#### 查询所有字段名
```bash
redis-cli HKEYS user:1001
```

#### 查询所有字段值
```bash
redis-cli HVALS user:1001
```

#### 查询字段数量
```bash
redis-cli HLEN user:1001
```

#### 检查字段是否存在
```bash
redis-cli HEXISTS user:1001 name
```
返回：`1` 表示存在，`0` 表示不存在

### 4. 一行命令执行（不进入交互式模式）

```bash
# 查询key是否存在
redis-cli EXISTS user:1001

# 查看key类型
redis-cli TYPE user:1001

# 获取字符串类型key的值
redis-cli GET keyname

# 查询所有字段和值
redis-cli HGETALL user:1001

# 查询单个字段
redis-cli HGET user:1001 name

# 查看过期时间
redis-cli TTL user:1001

# 查询所有字段名
redis-cli HKEYS user:1001

# 查询所有字段值
redis-cli HVALS user:1001

# 查询字段数量
redis-cli HLEN user:1001
```

### 5. 带密码的Redis连接

如果Redis设置了密码，需要添加 `-a` 参数：

```bash
redis-cli -a yourpassword HGETALL user:1001
```

### 6. 连接远程Redis服务器

如果Redis不在本地或端口不是6379：

```bash
redis-cli -h hostname -p port HGETALL user:1001
```

示例：
```bash
redis-cli -h 192.168.1.100 -p 6379 HGETALL user:1001
```

### 7. 其他常用查询命令

#### 查看所有匹配的key
```bash
# 查看匹配特定模式的key
redis-cli KEYS "user:*"

# 查看所有key（注意：必须用引号包裹 *，否则会报错）
redis-cli KEYS "*"
```

**重要提示：**
- 使用 `KEYS *` 时，**必须用引号包裹 `*`**，否则shell会将 `*` 解释为通配符，导致错误
- 错误示例：`redis-cli KEYS *` → 会报错：`ERR wrong number of arguments for 'keys' command`
- 正确示例：`redis-cli KEYS "*"` → 正常执行

**原因说明：**
在shell中，`*` 会被展开为当前目录下的所有文件名，所以：
- `redis-cli KEYS *` 实际变成了 `redis-cli KEYS file1 file2 file3 ...`（参数过多）
- `redis-cli KEYS "*"` 才是正确的命令（`*` 作为字符串传递给Redis）

**性能警告：**
- `KEYS *` 会扫描所有key，在生产环境中可能阻塞Redis
- 建议使用 `SCAN` 命令代替（非阻塞，但需要多次调用）

#### 删除key
```bash
redis-cli DEL user:1001
```

#### 设置key的过期时间
```bash
redis-cli EXPIRE user:1001 60
```
设置60秒后过期

#### 移除key的过期时间
```bash
redis-cli PERSIST user:1001
```

## 其他数据类型查询命令参考

### 字符串类型 (string)
```bash
# 获取值
redis-cli GET keyname

# 设置值
redis-cli SET keyname "value"

# 设置值并指定过期时间（秒）
redis-cli SETEX keyname 60 "value"

# 设置值（仅当key不存在时）
redis-cli SETNX keyname "value"

# 获取并设置值（原子操作）
redis-cli GETSET keyname "newvalue"

# 获取字符串的一部分
redis-cli GETRANGE keyname 0 5

# 获取字符串长度
redis-cli STRLEN keyname

# 追加字符串
redis-cli APPEND keyname "追加的内容"
```

### 列表类型 (list)
```bash
redis-cli LRANGE keyname 0 -1    # 获取所有元素
redis-cli LLEN keyname          # 获取列表长度
```

### 集合类型 (set)
```bash
redis-cli SMEMBERS keyname      # 获取所有成员
redis-cli SCARD keyname         # 获取成员数量
```

### 有序集合类型 (zset)
```bash
redis-cli ZRANGE keyname 0 -1 WITHSCORES    # 获取所有成员和分数（从低到高）
redis-cli ZREVRANGE keyname 0 -1 WITHSCORES # 获取所有成员和分数（从高到低）
redis-cli ZCARD keyname                     # 获取成员数量
```

## Redis 数据类型完整列表

Redis 支持以下数据类型：

### 基本数据类型（5种）

| 类型 | 名称 | 说明 | 查看类型命令 |
|------|------|------|-------------|
| **string** | 字符串 | 最基本的数据类型，可以存储字符串、数字、二进制数据 | `TYPE keyname` 返回 `string` |
| **hash** | 哈希 | 键值对集合，适合存储对象（类似Java的HashMap） | `TYPE keyname` 返回 `hash` |
| **list** | 列表 | 有序的字符串列表，支持从两端操作（类似Java的LinkedList） | `TYPE keyname` 返回 `list` |
| **set** | 集合 | 无序的字符串集合，成员唯一（类似Java的HashSet） | `TYPE keyname` 返回 `set` |
| **zset** | 有序集合 | 带分数的有序集合，按分数排序（类似Java的TreeMap） | `TYPE keyname` 返回 `zset` |

### 高级数据类型

| 类型 | 名称 | 说明 | Redis版本要求 |
|------|------|------|--------------|
| **stream** | 流 | 类似日志的数据结构，支持消息队列功能 | Redis 5.0+ |
| **bitmap** | 位图 | 基于string实现的位操作，用于统计 | 基于string |
| **hyperloglog** | 基数统计 | 用于统计唯一元素数量，占用空间小 | Redis 2.8+ |
| **geospatial** | 地理空间 | 存储地理位置信息，支持距离计算 | Redis 3.2+ |

### 查看key的数据类型

```bash
# 查看key的类型
redis-cli TYPE keyname

# 返回的可能值：
# - string  (字符串)
# - hash    (哈希)
# - list    (列表)
# - set     (集合)
# - zset    (有序集合)
# - stream  (流)
# - none    (key不存在)
```

### 数据类型快速参考

```bash
# 1. String (字符串)
redis-cli SET name "张三"
redis-cli GET name

# 2. Hash (哈希)
redis-cli HSET user:1001 name "李四" age "25"
redis-cli HGETALL user:1001

# 3. List (列表)
redis-cli LPUSH tasks "任务1" "任务2"
redis-cli LRANGE tasks 0 -1

# 4. Set (集合)
redis-cli SADD tags "go" "redis" "golang"
redis-cli SMEMBERS tags

# 5. ZSet (有序集合)
redis-cli ZADD scores 95.5 "张三" 88.0 "李四"
redis-cli ZRANGE scores 0 -1 WITHSCORES

# 6. Stream (流) - Redis 5.0+
redis-cli XADD mystream * field1 value1 field2 value2
redis-cli XREAD COUNT 10 STREAMS mystream 0

# 7. Bitmap (位图) - 基于string
redis-cli SETBIT mybitmap 0 1
redis-cli GETBIT mybitmap 0

# 8. HyperLogLog (基数统计)
redis-cli PFADD visitors user1 user2 user3
redis-cli PFCOUNT visitors

# 9. Geospatial (地理空间)
redis-cli GEOADD cities 116.3974 39.9093 "北京"
redis-cli GEODIST cities "北京" "上海" km
```

### 数据类型选择指南

| 使用场景 | 推荐类型 | 原因 |
|---------|---------|------|
| 存储单个值 | **string** | 简单直接 |
| 存储对象/用户信息 | **hash** | 可以存储多个字段，节省空间 |
| 消息队列/任务队列 | **list** | 支持FIFO/LIFO操作 |
| 标签/去重 | **set** | 自动去重，快速判断存在 |
| 排行榜/排序 | **zset** | 自动按分数排序 |
| 消息流/日志 | **stream** | 支持消息队列，可持久化 |
| 用户签到/统计 | **bitmap** | 节省空间，支持位操作 |
| UV统计 | **hyperloglog** | 占用空间极小 |
| 附近的人/距离计算 | **geospatial** | 内置地理位置功能 |

