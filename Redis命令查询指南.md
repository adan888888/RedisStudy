# Redis命令查询指南

## 查询 user:1001 的常用命令

### 1. 连接到Redis

```bash
redis-cli
```

### 2. 基本查询命令

#### 检查key是否存在
```bash
EXISTS user:1001
```
返回：`1` 表示存在，`0` 表示不存在

#### 查看key的类型
```bash
TYPE user:1001
```
返回：`hash`（哈希类型）

#### 查看key的过期时间
```bash
TTL user:1001
```
返回：
- `-1` 表示永不过期
- `-2` 表示key不存在
- 正数表示剩余秒数

### 3. 哈希类型操作（user:1001是哈希类型）

#### 查询所有字段和值
```bash
HGETALL user:1001
```
返回所有字段和对应的值

#### 查询单个字段
```bash
HGET user:1001 name
HGET user:1001 age
HGET user:1001 email
```

#### 查询所有字段名
```bash
HKEYS user:1001
```

#### 查询所有字段值
```bash
HVALS user:1001
```

#### 查询字段数量
```bash
HLEN user:1001
```

#### 检查字段是否存在
```bash
HEXISTS user:1001 name
```
返回：`1` 表示存在，`0` 表示不存在

### 4. 一行命令执行（不进入交互式模式）

```bash
# 查询key是否存在
redis-cli EXISTS user:1001

# 查看key类型
redis-cli TYPE user:1001

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
KEYS user:*
KEYS *
```

#### 删除key
```bash
DEL user:1001
```

#### 设置key的过期时间
```bash
EXPIRE user:1001 60
```
设置60秒后过期

#### 移除key的过期时间
```bash
PERSIST user:1001
```

## 其他数据类型查询命令参考

### 字符串类型 (string)
```bash
GET keyname
```

### 列表类型 (list)
```bash
LRANGE keyname 0 -1    # 获取所有元素
LLEN keyname          # 获取列表长度
```

### 集合类型 (set)
```bash
SMEMBERS keyname      # 获取所有成员
SCARD keyname         # 获取成员数量
```

### 有序集合类型 (zset)
```bash
ZRANGE keyname 0 -1 WITHSCORES    # 获取所有成员和分数（从低到高）
ZREVRANGE keyname 0 -1 WITHSCORES # 获取所有成员和分数（从高到低）
ZCARD keyname                     # 获取成员数量
```

