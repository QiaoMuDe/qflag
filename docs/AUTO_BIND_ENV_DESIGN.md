# 环境变量自动绑定功能设计文档

## 一、背景和目标

### 1.1 当前机制

QFlag 当前支持通过 `BindEnv(name string)` 方法显式绑定环境变量：

```go
host := cmd.String("host", "h", "主机地址", "localhost")
host.BindEnv("DATABASE_HOST")  // 显式指定环境变量名
```

**优点**：
- ✅ 灵活性高，可以绑定到任意环境变量名
- ✅ 明确性，一眼就能看出绑定到哪个环境变量

**缺点**：
- ❌ 需要手动指定环境变量名，代码冗余
- ❌ 容易拼写错误
- ❌ 样板代码多

### 1.2 目标

新增自动绑定功能，简化环境变量绑定流程，提升开发体验。

---

## 二、方案设计

### 2.1 核心思路

**保持现有方法不变，新增自动绑定方法**，实现职责分离：

- `BindEnv(name string)` - 显式指定环境变量名（现有方法，保持不变）
- `AutoBindEnv()` - 自动使用标志名作为环境变量名（新增方法）

### 2.2 方法设计

#### 2.2.1 AutoBindEnv 方法

```go
// AutoBindEnv 自动绑定环境变量
//
// 功能说明:
//   - 自动使用标志的长名称作为环境变量名（转为大写）
//   - 如果没有设置长名称，会触发 panic
//   - 环境变量前缀（EnvPrefix）在解析时自动拼接，无需手动处理
//
// 使用示例:
//   - 标志 "host" -> 绑定到 "HOST"
//   - 标志 "PORT" -> 绑定到 "PORT"
//   - 标志 "db-host" -> 绑定到 "DB-HOST"
//
// 注意事项:
//   - 环境变量的优先级低于命令行参数
//   - 必须设置长名称，否则会 panic
//   - 短名称不会被使用，避免冲突
//   - 自动转为大写，确保环境变量命名规范
func (f *BaseFlag[T]) AutoBindEnv() {
    f.mu.Lock()
    defer f.mu.Unlock()
    
    // 必须有长名称
    if f.longName == "" {
        panic(types.NewError("EMPTY_LONG_NAME", "flag must have a long name for AutoBindEnv", nil))
    }
    
    // 使用长名称并转为大写
    f.envVar = strings.ToUpper(f.longName)
}
```

#### 2.2.2 实现位置

文件路径：`internal/flag/base_flag.go`

---

## 三、使用示例

### 3.1 基础使用

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 设置环境变量（大写）
    os.Setenv("HOST", "localhost")
    os.Setenv("PORT", "8080")
    
    // 创建标志并自动绑定环境变量
    host := qflag.Root.String("host", "h", "主机地址", "127.0.0.1")
    host.AutoBindEnv()  // 自动绑定到 "HOST"（大写）
    
    port := qflag.Root.Uint("port", "p", "端口号", 3000)
    port.AutoBindEnv()  // 自动绑定到 "PORT"（大写）
    
    // 解析参数
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析错误: %v\n", err)
        return
    }
    
    // 使用参数
    fmt.Printf("主机: %s\n", host.Get())  // 输出: localhost（来自环境变量 HOST）
    fmt.Printf("端口: %d\n", port.Get())  // 输出: 8080（来自环境变量 PORT）
}
```

### 3.2 对比现有方法

```go
// 方式1：显式绑定（现有方法，保持不变）
host := cmd.String("host", "h", "主机地址", "localhost")
host.BindEnv("DATABASE_HOST")  // 绑定到 DATABASE_HOST

// 方式2：自动绑定（新增方法）
host := cmd.String("host", "h", "主机地址", "localhost")
host.AutoBindEnv()  // 自动绑定到 HOST
```

### 3.3 环境变量优先级

```
命令行参数 > 环境变量 > 默认值
```

示例：

```bash
# 设置环境变量（大写）
export HOST=192.168.1.1

# 方式1：使用环境变量
./myapp
# 输出：主机: 192.168.1.1（来自环境变量 HOST）

# 方式2：命令行参数优先
./myapp --host 10.0.0.1
# 输出：主机: 10.0.0.1（命令行参数优先级更高）
```

### 3.4 错误处理示例

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/internal/types"
)

func main() {
    // ❌ 错误示例：没有长名称会 panic
    host := qflag.Root.String("", "h", "主机地址", "localhost")
    
    defer func() {
        if r := recover(); r != nil {
            err := r.(*types.Error)
            fmt.Printf("错误码: %s\n", err.Code)      // 输出: EMPTY_LONG_NAME
            fmt.Printf("错误消息: %s\n", err.Message)  // 输出: flag must have a long name for AutoBindEnv
        }
    }()
    
    host.AutoBindEnv()  // panic: 结构化错误
    
    // ✅ 正确示例：必须有长名称
    host = qflag.Root.String("host", "h", "主机地址", "localhost")
    host.AutoBindEnv()  // 正常工作，绑定到 HOST
    
    qflag.Parse()
}
```

### 3.5 复杂标志名示例

```go
// 连字符分隔的标志名
dbHost := cmd.String("db-host", "", "数据库主机", "localhost")
dbHost.AutoBindEnv()  // 绑定到 "DB-HOST"

// 下划线分隔的标志名
db_port := cmd.String("db_port", "", "数据库端口", "3306")
db_port.AutoBindEnv()  // 绑定到 "DB_PORT"

// 驼峰命名的标志名
dbUserName := cmd.String("dbUserName", "", "数据库用户名", "root")
dbUserName.AutoBindEnv()  // 绑定到 "DBUSERNAME"
```

---

## 四、批量自动绑定（可选）

### 4.1 批量自动绑定方法

```go
// internal/cmd/cmd.go

// AutoBindAllEnv 为所有标志自动绑定环境变量
//
// 功能说明:
//   - 遍历命令的所有标志
//   - 为每个标志调用 AutoBindEnv() 方法
//   - 批量设置环境变量绑定
//
// 使用示例:
//   cmd.String("host", "h", "主机地址", "localhost")
//   cmd.Uint("port", "p", "端口号", 8080)
//   cmd.AutoBindAllEnv()  // 自动绑定 host 和 port
func (c *Cmd) AutoBindAllEnv() error {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    for _, f := range c.flagRegistry.List() {
        // 调用标志的 AutoBindEnv 方法
        if bf, ok := f.(interface{ AutoBindEnv() }); ok {
            bf.AutoBindEnv()
        }
    }
    
    return nil
}
```

### 4.2 使用示例

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 设置环境变量
    os.Setenv("host", "localhost")
    os.Setenv("port", "8080")
    
    // 创建标志
    qflag.Root.String("host", "h", "主机地址", "127.0.0.1")
    qflag.Root.Uint("port", "p", "端口号", 3000)
    
    // 批量自动绑定所有标志的环境变量
    qflag.Root.AutoBindAllEnv()
    
    // 解析参数
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析错误: %v\n", err)
        return
    }
    
    // 使用参数
    fmt.Printf("主机: %s\n", qflag.Root.GetFlag("host").GetStr())
    fmt.Printf("端口: %d\n", qflag.Root.GetFlag("port").(*qflag.UintFlag).Get())
}
```

---

## 五、方法对比

| 方法 | 参数 | 用途 | 示例 | 适用场景 |
|------|------|------|------|---------|
| `BindEnv(name)` | 必需 | 显式指定环境变量名 | `BindEnv("DATABASE_HOST")` | 需要自定义环境变量名 |
| `AutoBindEnv()` | 无 | 自动使用长名称（转大写） | `AutoBindEnv()` → `HOST` | 标志名与环境变量名一致 |

**关键区别**：
- `BindEnv(name)` - 完全自定义，灵活性最高
- `AutoBindEnv()` - 自动绑定，强制使用长名称并转大写，规范统一

**设计优势**：
- ✅ 避免短名称冲突
- ✅ 统一环境变量命名规范（大写）
- ✅ 强制要求长名称，提高代码可读性
- ✅ 防止误解析到其他环境变量

---

## 六、优势分析

### 6.1 向后兼容

- ✅ 不修改现有 `BindEnv` 方法
- ✅ 现有代码无需修改
- ✅ 保持API稳定性

### 6.2 职责清晰

- ✅ `BindEnv(name)` - 显式指定，灵活性高
- ✅ `AutoBindEnv()` - 自动绑定，简洁方便
- ✅ 两个方法各司其职，互不干扰

### 6.3 开发体验

**之前**：
```go
host := cmd.String("host", "h", "主机地址", "localhost")
host.BindEnv("HOST")  // 冗余，需要手动指定
port := cmd.Uint("port", "p", "端口号", 8080)
port.BindEnv("PORT")  // 冗余，需要手动指定
```

**之后**：
```go
host := cmd.String("host", "h", "主机地址", "localhost")
host.AutoBindEnv()  // 简洁，自动绑定
port := cmd.Uint("port", "p", "端口号", 8080)
port.AutoBindEnv()  // 简洁，自动绑定
```

### 6.4 灵活性

- ✅ 支持显式绑定（需要自定义名称时）
- ✅ 支持自动绑定（常规场景）
- ✅ 支持混合使用

---

## 七、实现计划

### 7.1 第一阶段：基础实现

1. 在 `internal/flag/base_flag.go` 中添加 `AutoBindEnv()` 方法
2. 添加单元测试
3. 更新文档和示例

### 7.2 第二阶段：增强功能（可选）

1. 支持环境变量前缀
2. 添加批量自动绑定方法 `AutoBindAllEnv()`
3. 添加更多测试用例

### 7.3 第三阶段：文档完善

1. 更新 README.md
2. 更新 APIDOC.md
3. 添加使用示例到 examples 目录

---

## 八、测试用例

### 8.1 基础测试

```go
func TestAutoBindEnv(t *testing.T) {
    // 测试长名称（转为大写）
    flag := NewStringFlag("host", "h", "主机地址", "localhost")
    flag.AutoBindEnv()
    assert.Equal(t, "HOST", flag.GetEnvVar())
    
    // 测试复杂标志名
    flag = NewStringFlag("db-host", "", "数据库主机", "localhost")
    flag.AutoBindEnv()
    assert.Equal(t, "DB-HOST", flag.GetEnvVar())
    
    // 测试下划线分隔
    flag = NewStringFlag("db_port", "", "数据库端口", "3306")
    flag.AutoBindEnv()
    assert.Equal(t, "DB_PORT", flag.GetEnvVar())
    
    // 测试驼峰命名
    flag = NewStringFlag("dbUserName", "", "数据库用户名", "root")
    flag.AutoBindEnv()
    assert.Equal(t, "DBUSERNAME", flag.GetEnvVar())
}

func TestAutoBindEnvPanic(t *testing.T) {
    // 测试没有长名称时 panic
    flag := NewStringFlag("", "h", "主机地址", "localhost")
    
    defer func() {
        if r := recover(); r != nil {
            err := r.(*types.Error)
            assert.Equal(t, "EMPTY_LONG_NAME", err.Code)
            assert.Equal(t, "flag must have a long name for AutoBindEnv", err.Message)
        }
    }()
    
    flag.AutoBindEnv()
    t.Error("Expected panic but didn't get one")
}
```

### 8.2 集成测试

```go
func TestAutoBindEnvIntegration(t *testing.T) {
    // 设置环境变量（大写）
    os.Setenv("HOST", "192.168.1.1")
    defer os.Unsetenv("HOST")
    
    // 创建命令
    cmd := NewCmd("test", "", ContinueOnError)
    host := cmd.String("host", "h", "主机地址", "localhost")
    host.AutoBindEnv()
    
    // 解析
    err := cmd.Parse([]string{})
    assert.NoError(t, err)
    
    // 验证
    assert.Equal(t, "192.168.1.1", host.Get())
}

func TestAutoBindEnvComplexNames(t *testing.T) {
    // 测试复杂标志名
    os.Setenv("DB-HOST", "192.168.1.100")
    os.Setenv("DB_PORT", "3306")
    defer func() {
        os.Unsetenv("DB-HOST")
        os.Unsetenv("DB_PORT")
    }()
    
    cmd := NewCmd("test", "", ContinueOnError)
    
    dbHost := cmd.String("db-host", "", "数据库主机", "localhost")
    dbHost.AutoBindEnv()
    
    dbPort := cmd.String("db_port", "", "数据库端口", "5432")
    dbPort.AutoBindEnv()
    
    cmd.Parse([]string{})
    
    assert.Equal(t, "192.168.1.100", dbHost.Get())
    assert.Equal(t, "3306", dbPort.Get())
}
```

---

## 九、注意事项

### 9.1 必须设置长名称

`AutoBindEnv()` 方法要求标志必须设置长名称，否则会 panic：

```go
// ❌ 错误：没有长名称会 panic
host := cmd.String("", "h", "主机地址", "localhost")
host.AutoBindEnv()  // panic: AutoBindEnv: flag must have a long name

// ✅ 正确：必须有长名称
host := cmd.String("host", "h", "主机地址", "localhost")
host.AutoBindEnv()  // 正常工作，绑定到 HOST
```

### 9.2 环境变量优先级

```
命令行参数 > 环境变量 > 默认值
```

### 9.3 自动转大写

标志名会自动转换为大写，符合环境变量命名规范：

```go
// 标志名会自动转为大写
host := cmd.String("host", "h", "主机地址", "localhost")
host.AutoBindEnv()  // 绑定到 HOST（大写）

// 复杂标志名也会转为大写
dbHost := cmd.String("db-host", "", "数据库主机", "localhost")
dbHost.AutoBindEnv()  // 绑定到 DB-HOST（大写）
```

### 9.4 环境变量前缀

环境变量前缀（EnvPrefix）在解析时自动拼接，无需在 AutoBindEnv 中处理：

```go
// 设置环境变量前缀
cmd.SetEnvPrefix("MYAPP")

// 创建标志并自动绑定
host := cmd.String("host", "h", "主机地址", "localhost")
host.AutoBindEnv()  // 绑定到 "HOST"

// 解析时会自动尝试：
// 1. MYAPP_HOST（带前缀）
// 2. HOST（不带前缀）
```

### 9.5 避免短名称冲突

`AutoBindEnv()` 只使用长名称，避免短名称冲突：

```go
// 短名称不会被使用
host := cmd.String("host", "h", "主机地址", "localhost")
host.AutoBindEnv()  // 只绑定到 HOST，不会绑定到 H

// 这样可以避免多个标志使用相同短名称导致的冲突
```

---

## 十、总结

本方案通过新增 `AutoBindEnv()` 方法，在保持向后兼容的前提下，简化了环境变量绑定流程，提升了开发体验。核心优势：

1. ✅ **向后兼容** - 不修改现有API
2. ✅ **职责清晰** - 两个方法各司其职
3. ✅ **开发体验好** - 常用场景更简洁
4. ✅ **灵活性高** - 支持两种绑定方式
5. ✅ **易于实现** - 实现简单，测试充分
6. ✅ **强制规范** - 必须使用长名称，自动转大写，统一命名规范
7. ✅ **避免冲突** - 不使用短名称，防止环境变量冲突
8. ✅ **前缀自动处理** - 环境变量前缀在解析时自动拼接，无需手动处理

### 关键设计决策

1. **强制使用长名称**：必须设置长名称，否则会 panic，确保代码可读性和规范性
2. **自动转大写**：标志名自动转换为大写，符合环境变量命名规范
3. **不使用短名称**：避免短名称冲突，防止误解析到其他环境变量
4. **前缀自动处理**：环境变量前缀（EnvPrefix）在解析时自动拼接，`AutoBindEnv()` 方法只负责绑定标志名

### 设计优势

这个设计比之前的方案更加严格和规范：

- **防止冲突**：不使用短名称，避免多个标志使用相同短名称导致的环境变量冲突
- **统一规范**：自动转大写，确保环境变量命名规范统一
- **强制要求**：必须设置长名称，提高代码可读性
- **防止误解析**：避免解析到不相关的环境变量

这是一个渐进式的增强方案，可以根据实际需求逐步实现和完善。
