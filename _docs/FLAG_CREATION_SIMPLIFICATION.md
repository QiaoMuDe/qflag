# 标志创建与注册简化方案

## 现状分析

### 当前使用方式 (繁琐) 

```go
// 步骤1: 创建标志
portFlag := flag.NewIntFlag("port", "p", "服务端口号", 8080)

// 步骤2: 添加到命令
if err := cmd.AddFlag(portFlag); err != nil {
    panic(err)
}
```

### 问题
1. 创建和注册分离, 需要两步操作
2. 代码冗长, 不够直观
3. 需要导入 `flag` 包

---

## 优化方案

### 目标使用方式 (简洁) 

```go
// 方式1: 直接创建并注册
cmd.Int("port", "p", "服务端口号", 8080)

// 方式2: 创建后获取标志对象进行后续操作
verboseFlag := cmd.String("verbose", "v", "详细输出模式", "info")
verboseFlag.AddAlias("vv") // 进阶操作
```

### 设计原则
1. **链式调用**: `cmd.Int(...)` 返回标志对象, 支持链式调用
2. **向后兼容**: 保留 `AddFlag` 方法和原有创建函数
3. **类型安全**: 每个类型对应一个方法
4. **一致性**: 所有方法遵循相同签名模式

---

## 实现方案

### 在 Cmd 中添加的方法

#### 1. 基础类型方法

```go
// 整数类型
func (c *Cmd) Int(name, shortName, description string, default_ int) *IntFlag

// 字符串类型
func (c *Cmd) String(name, shortName, description string, default_ string) *StringFlag

// 布尔类型
func (c *Cmd) Bool(name, shortName, description string, default_ bool) *BoolFlag

// 浮点数类型
func (c *Cmd) Float64(name, shortName, description string, default_ float64) *Float64Flag

// 持续时间类型
func (c *Cmd) Duration(name, shortName, description string, default_ time.Duration) *DurationFlag

// 64位整数类型
func (c *Cmd) Int64(name, shortName, description string, default_ int64) *Int64Flag

// 无符号整数类型
func (c *Cmd) Uint(name, shortName, description string, default_ uint) *UintFlag

// 64位无符号整数类型
func (c *Cmd) Uint64(name, shortName, description string, default_ uint64) *Uint64Flag

// 大小类型 (如 1KB, 100MB) 
func (c *Cmd) Size(name, shortName, description string, default_ uint64) *SizeFlag
```

#### 2. 枚举类型方法

```go
// 枚举类型
func (c *Cmd) Enum(name, shortName, description string, default_ string, options ...string) *EnumFlag
```

#### 3. 切片类型方法

```go
// 字符串切片类型
func (c *Cmd) StringSlice(name, shortName, description string, default_ []string) *StringSliceFlag

// 整数切片类型
func (c *Cmd) IntSlice(name, shortName, description string, default_ []int) *IntSliceFlag

// 64位整数切片类型
func (c *Cmd) Int64Slice(name, shortName, description string, default_ []int64) *Int64SliceFlag
```

#### 4. 映射类型方法

```go
// 映射类型
func (c *Cmd) Map(name, shortName, description string, default_ map[string]string) *MapFlag
```

---

## 内部实现

### 方法内部流程

```go
func (c *Cmd) Int(name, shortName, description string, default_ int) *IntFlag {
    // 1. 创建标志
    flag := flag.NewIntFlag(name, shortName, description, default_)
    
    // 2. 注册到命令
    if err := c.flagRegistry.Register(flag); err != nil {
        // 处理错误: 可以 panic、记录日志或返回 nil
        panic(err)
    }
    
    // 3. 返回标志对象供后续操作
    return flag
}
```

### 错误处理策略

#### 方案A: Panic (推荐) 
- 标志名冲突时 panic
- 符合 "快速失败" 原则
- 开发阶段即可发现配置错误

```go
func (c *Cmd) Int(name, shortName, description string, default_ int) *IntFlag {
    if c.HasFlag(name) {
        panic(fmt.Sprintf("flag '%s' already exists", name))
    }
    flag := flag.NewIntFlag(name, shortName, description, default_)
    c.flagRegistry.Register(flag) // 注册不会失败, 因为已检查
    return flag
}
```

#### 方案B: 返回错误
- 更灵活, 但不常用
- 需要处理错误

```go
func (c *Cmd) Int(name, shortName, description string, default_ int) (*IntFlag, error) {
    if c.HasFlag(name) {
        return nil, fmt.Errorf("flag '%s' already exists", name)
    }
    flag := flag.NewIntFlag(name, shortName, description, default_)
    return flag, c.flagRegistry.Register(flag)
}
```

**推荐方案A**, 因为: 
- 标志配置应在程序启动时完成
- 重复注册是编程错误, 应立即发现

---

## 辅助方法

### 检查标志是否存在

```go
// 检查标志是否存在 (长名或短名) 
func (c *Cmd) HasFlag(name string) bool

// 检查长标志名是否存在
func (c *Cmd) HasLongFlag(name string) bool

// 检查短标志名是否存在
func (c *Cmd) HasShortFlag(name string) bool
```

### 获取标志

```go
// 根据名称获取标志
func (c *Cmd) GetFlag(name string) (types.Flag, bool)

// 便捷获取方法 (返回具体类型) 
func (c *Cmd) GetInt(name string) (*IntFlag, bool)
func (c *Cmd) GetString(name string) (*StringFlag, bool)
func (c *Cmd) GetBool(name string) (*BoolFlag, bool)
// ... 其他类型
```

---

## 使用示例

### 完整示例

```go
package main

import (
    "fmt"
    "os"
    "time"
    
    "gitee.com/MM-Q/qflag/internal/cmd"
)

func main() {
    rootCmd := cmd.NewCmd("app", "a", types.ContinueOnError)
    rootCmd.SetDesc("这是一个示例应用")
    
    // 直接创建和注册标志
    port := rootCmd.Int("port", "p", "服务端口号", 8080)
    host := rootCmd.String("host", "H", "服务地址", "localhost")
    debug := rootCmd.Bool("debug", "d", "开启调试模式", false)
    timeout := rootCmd.Duration("timeout", "t", "请求超时时间", 30*time.Second)
    
    // 枚举类型
    logLevel := rootCmd.Enum("log-level", "l", "日志级别", "info", "debug", "info", "warn", "error")
    
    // 切片类型
    tags := rootCmd.StringSlice("tag", "", "标签", nil)
    
    // 设置运行函数
    rootCmd.SetRun(func(c types.Cmd) error {
        fmt.Printf("启动服务: %s:%d\n", host.Get(), port.Get())
        fmt.Printf("调试模式: %v\n", debug.Get())
        fmt.Printf("超时时间: %v\n", timeout.Get())
        fmt.Printf("日志级别: %v\n", logLevel.Get())
        return nil
    })
    
    // 解析并运行
    if err := rootCmd.ParseAndRoute(os.Args[1:]); err != nil {
        fmt.Fprintf(os.Stderr, "错误: %v\n", err)
        os.Exit(1)
    }
}
```

### 子命令中使用

```go
// 创建子命令
serveCmd := cmd.NewCmd("serve", "s", types.ContinueOnError)
serveCmd.SetDesc("启动服务")

// 子命令中直接创建标志
servePort := serveCmd.Int("port", "p", "服务端口", 8080)

// 添加到根命令
rootCmd.AddSubCmd(serveCmd)
```

---

## 向后兼容性

### 保留原有方法

```go
// 原有方法仍然可用
flag := flag.NewIntFlag("port", "p", "端口号", 8080)
cmd.AddFlag(flag)
```

### 标志类型导出

确保所有标志类型 (如 `IntFlag`, `StringFlag` 等) 可从 `flag` 包导出, 供高级用户使用。

---

## 总结

| 特性 | 原来 | 优化后 |
|-----|-----|-------|
| 创建步骤 | 2步 (创建+添加)  | 1步 |
| 代码行数 | 4-6行 | 1行 |
| 导入包 | `flag` + `cmd` | 仅 `cmd` |
| 可读性 | 一般 | 优秀 |
| 灵活性 | 高 | 高 (返回标志对象)  |

这个方案大幅简化了用户的使用方式, 同时保留了所有高级功能。
