# 锁架构重构方案

## 问题分析

### 当前架构的锁问题

当前代码存在**锁嵌套**问题, 导致死锁风险: 

```
Cmd.Parse()
  └── Cmd.FlagRegistry().List()
        └── FlagRegistry.List()
              └── registry.List()
                    └── for flag in items
                          └── flag.Name()  ← 需要获取 Flag 的锁
```

**具体死锁原因**: 
- `registry.List()` 持有 `registry.mu` 锁
- 遍历 `flag.Name()` 时尝试获取 `flag.mu` 锁
- 如果同一个 flag 正在被修改 (持有 `flag.mu` 并等待 `registry.mu`) , 就形成死锁

### 当前锁的使用位置

| 组件 | 锁类型 | 问题 |
|------|--------|------|
| `Cmd` | `sync.RWMutex` | 顶层, 合理 |
| `FlagRegistryImpl` | `sync.RWMutex` | 冗余, 与 registry 锁嵌套 |
| `registry[T]` | `sync.RWMutex` | 内部组件, 不需要 |
| `BaseFlag[T]` | `sync.RWMutex` | 顶层, 合理 |

## 重构方案: 分层锁架构

### 设计原则

```
┌─────────────────────────────────────────────────────────┐
│  Cmd 层 (读写锁)                                     │
│  - 管理子命令和标志集合                                   │
│  - 暴露给用户调用的 API                                   │
│  - 持有 Registry 的引用                                  │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  Registry 层 (无锁)                                      │
│  - 作为纯数据结构, 不管理锁                                │
│  - 所有操作由调用方 (Cmd) 的锁保护                      │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  Flag 层 (读写锁)                                        │
│  - 标志的值需要保护                                       │
│  - 支持并发读取                                          │
└─────────────────────────────────────────────────────────┘
```

### 方案: 移除 Registry 的锁

#### 核心思路
Registry 作为**纯数据结构**, 不应该有自己的锁。所有对 Registry 的访问都由调用方的锁保护。

#### 修改内容

##### 1. 移除 `registry[T]` 的锁

```go
// 修改前
type registry[T any] struct {
    mu        sync.RWMutex  // 删除
    items     map[string]T
    nameIndex map[string]string
}

// 修改后
type registry[T any] struct {
    items     map[string]T
    nameIndex map[string]string
}

// 创建时
func NewRegistry[T any]() *registry[T] {
    return &registry[T]{
        items:     make(map[string]T),
        nameIndex: make(map[string]string),
    }
}
```

##### 2. 移除 `registry[T]` 方法中的锁

所有方法移除 `mu.Lock()` 和 `mu.RLock()`: 

```go
// 修改前
func (r *registry[T]) Register(name string, item T) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    // ...
}

// 修改后
func (r *registry[T]) Register(name string, item T) error {
    // 无锁, 纯数据结构操作
    // ...
}
```

##### 3. 移除 `FlagRegistryImpl` 和 `CmdRegistryImpl` 的锁

```go
// 修改前
type FlagRegistryImpl struct {
    *registry[types.Flag]
    mu sync.RWMutex  // 删除
}

type CmdRegistryImpl struct {
    *registry[types.Cmd]
    mu sync.RWMutex  // 删除
}

// 修改后
type FlagRegistryImpl struct {
    *registry[types.Flag]
    // 无锁
}

type CmdRegistryImpl struct {
    *registry[types.Cmd]
    // 无锁
}
```

##### 4. Cmd 层保证并发安全

所有访问 Registry 的操作由 Cmd 的锁保护: 

```go
func (c *Cmd) AddFlag(flag types.Flag) error {
    c.mu.Lock()         // ← 获取 Cmd 的锁
    defer c.mu.Unlock()

    // 此时可以安全访问 c.flagRegistry
    // 因为没有其他锁会嵌套
    return c.flagRegistry.Register(flag.Name(), flag)
}

func (c *Cmd) Flags() []types.Flag {
    c.mu.RLock()        // ← 读锁
    defer c.mu.RUnlock()

    // 安全访问
    return c.flagRegistry.List()
}
```

##### 5. Flag 层保持锁

```go
type BaseFlag[T any] struct {
    mu          sync.RWMutex  // 保留, 用于保护值
    longName    string
    shortName   string
    // ...
}

// 所有修改值的方法需要加锁
func (f *BaseFlag[T]) Set(value string) error {
    f.mu.Lock()
    defer f.mu.Unlock()
    // ...
}
```

### 方案优势

1. **消除死锁**: 锁不会嵌套
2. **简化调试**: 只需要关注 Cmd 层的锁
3. **性能更好**: Registry 操作无锁开销
4. **职责清晰**: 
   - Registry = 纯数据结构
   - Flag = 值的安全访问
   - Cmd = 顶层并发控制

### 需要修改的文件

| 文件 | 修改内容 |
|------|----------|
| `internal/registry/impl.go` | 删除 `mu` 字段和方法中的锁 |
| `internal/registry/flag_registry.go` | 删除 `mu` 字段 |
| `internal/registry/cmd_registry.go` | 删除 `mu` 字段 |
| `internal/cmd/base_cmd.go` | 确保所有访问 Registry 的操作在锁保护下 |

### 风险评估

| 风险 | 级别 | 缓解措施 |
|------|------|----------|
| Registry 公开访问 | 低 | Go 没有真正的私有字段, 但约定 Registry 只通过 Cmd 访问 |
| Flag.Name() 调用 | 低 | Flag 的锁与 Cmd 锁不嵌套 |
| Registry 内部方法并发 | 低 | 所有调用方都有锁保护 |

### 对比

| 指标 | 当前架构 | 重构后 |
|------|----------|--------|
| 锁的数量 | 4 | 2 |
| 死锁风险 | 高 | 无 |
| 代码复杂度 | 高 | 低 |
| 性能 | 中 | 高 |
| 可维护性 | 低 | 高 |

## 总结

通过将锁集中在 Cmd 和 Flag 层, Registry 作为纯数据结构, 可以: 
- ✅ 消除死锁风险
- ✅ 简化并发控制
- ✅ 提高代码可维护性
