# CmdOpts 设计文档

## 一、背景与需求

### 1.1 背景
QFlag 项目中已经存在 `CmdSpec` 结构体，用于通过规格创建新命令。用户在实际使用中，经常遇到以下场景：
- 已经通过 `NewCmd` 创建了命令实例
- 需要批量设置命令的多个属性（描述、版本、示例等）
- 希望有一个结构化的方式来管理这些配置

### 1.2 需求
用户需要一个结构体来配置现有命令的属性，而不是创建新命令实例。需要配置的属性包括：
- 命令描述
- 运行函数
- 版本号
- 是否使用中文
- 环境变量前缀
- 命令使用语法
- Logo 文本
- 示例
- 注意事项
- 子命令列表
- 互斥组

### 1.3 现有结构
- `CmdSpec`：用于**创建**新命令
- `types.CmdConfig`：用于**存储**命令内部配置
- **缺失**：用于**配置**现有命令的结构体

---

## 二、设计目标

### 2.1 核心目标
- 提供一个结构化的方式来配置现有命令
- 与 `CmdSpec` 形成对比：Spec（创建）vs Opts（配置）
- 支持部分配置（未设置的属性不会被修改）
- 提供清晰的错误处理

### 2.2 设计原则
- **简洁性**：命名简洁，易于理解和使用
- **一致性**：与现有代码风格保持一致
- **可扩展性**：易于添加新的配置选项
- **类型安全**：利用 Go 的类型系统保证安全
- **语义清晰**：方法名明确表达其功能

---

## 三、命名选择

### 3.1 结构体命名对比

| 命名 | 长度 | 语义清晰度 | 推荐度 |
|------|------|-----------|--------|
| CmdConfigurer | 12 字符 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **CmdOpts** | **7 字符** | ⭐⭐⭐⭐ | **⭐⭐⭐⭐⭐** |
| CmdOptions | 10 字符 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| CmdCustomizer | 12 字符 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

### 3.2 方法命名对比

| 方法名 | 语义 | 优点 | 推荐度 |
|--------|------|------|--------|
| Configure(opts) | 配置命令 | 更简洁，符合常见习惯 | ⭐⭐⭐⭐⭐ |
| **ApplyOpts(opts)** | **应用选项到命令** | **更明确，与 CmdOpts 呼应** | **⭐⭐⭐⭐⭐** |

### 3.3 最终选择：CmdOpts + ApplyOpts

**选择理由**：
1. **简洁**：`CmdOpts` 只有 7 个字符，易于输入和记忆
2. **语义清晰**：Opts = Options（选项），明确表达配置选项的含义
3. **符合习惯**：Go 社区常用简写（如 `json.Encoder` 中的 `Encoder`）
4. **避免冲突**：不会与 `types.CmdConfig` 冲突
5. **形成对比**：
   - `CmdSpec`：用于创建新命令（Specification）
   - `CmdOpts`：用于配置现有命令（Options）
6. **方法呼应**：
   - `CmdOpts` 结构体
   - `ApplyOpts` 方法：Apply Options = 应用选项
   - 语义清晰，易于理解

---

## 四、结构体设计

### 4.1 CmdOpts 结构体定义

```go
// CmdOpts 命令选项
//
// CmdOpts 提供了配置现有命令的方式，包含命令的所有可配置属性。
// 与 CmdSpec 不同，CmdOpts 用于配置已存在的命令，而不是创建新命令。
//
// 使用场景:
//   - 已有命令实例，需要批量设置属性
//   - 需要结构化的配置管理
//   - 需要部分配置（未设置的属性不会被修改）
//
// 示例:
//   cmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)
//   opts := &cmd.CmdOpts{
//       Desc: "我的应用程序",
//       Version: "1.0.0",
//       UseChinese: true,
//   }
//   cmd.ApplyOpts(opts)
type CmdOpts struct {
    // 基本属性
    Desc string // 命令描述

    // 运行函数
    RunFunc func(types.Command) error // 命令执行函数

    // 配置选项
    Version     string // 版本号
    UseChinese  bool   // 是否使用中文
    EnvPrefix   string // 环境变量前缀
    UsageSyntax string // 命令使用语法
    LogoText    string // Logo文本

    // 示例和说明
    Examples map[string]string // 示例使用, key为描述, value为示例命令
    Notes    []string          // 注意事项

    // 子命令和互斥组
    SubCmds     []types.Command    // 子命令列表, 用于添加到命令中
    MutexGroups []types.MutexGroup // 互斥组列表
}
```

### 4.2 字段说明

| 字段 | 类型 | 说明 | 是否必填 |
|------|------|------|---------|
| Desc | string | 命令描述 | 否 |
| RunFunc | func(types.Command) error | 命令执行函数 | 否 |
| Version | string | 版本号 | 否 |
| UseChinese | bool | 是否使用中文 | 否 |
| EnvPrefix | string | 环境变量前缀 | 否 |
| UsageSyntax | string | 命令使用语法 | 否 |
| LogoText | string | Logo文本 | 否 |
| Examples | map[string]string | 示例使用 | 否 |
| Notes | []string | 注意事项 | 否 |
| SubCmds | []types.Command | 子命令列表 | 否 |
| MutexGroups | []types.MutexGroup | 互斥组列表 | 否 |

---

## 五、方法设计

### 5.1 NewCmdOpts

```go
// NewCmdOpts 创建新的命令选项
//
// 返回值:
//   - *CmdOpts: 初始化的命令选项
//
// 功能说明:
//   - 创建基本命令选项
//   - 初始化所有字段为零值
//   - 初始化 map 和 slice 避免空指针
func NewCmdOpts() *CmdOpts {
    return &CmdOpts{
        Examples:    make(map[string]string),
        Notes:       []string{},
        SubCmds:     []types.Command{},
        MutexGroups: []types.MutexGroup{},
    }
}
```

### 5.2 Cmd.ApplyOpts

```go
// ApplyOpts 应用选项到命令
//
// 参数:
//   - opts: 命令选项
//
// 返回值:
//   - error: 应用选项失败时返回错误
//
// 功能说明:
//   - 将选项结构体的所有属性应用到当前命令
//   - 支持部分配置（未设置的属性不会被修改）
//   - 使用defer捕获panic, 转换为错误返回
//
// 应用顺序:
//   1. 基本属性（Desc、RunFunc）
//   2. 配置选项（Version、UseChinese、EnvPrefix、UsageSyntax、LogoText）
//   3. 示例和说明（Examples、Notes）
//   4. 互斥组（MutexGroups）
//   5. 子命令（SubCmds）
//
// 错误处理:
//   - 选项为 nil: 返回 INVALID_CMDOPTS 错误
//   - 添加子命令失败: 返回 FAILED_TO_ADD_SUBCMDS 错误
//   - panic: 转换为 PANIC 错误
//
// 线程安全:
//   - 方法内部使用读写锁保护并发访问
//   - 可以安全地在多个 goroutine 中调用
func (c *Cmd) ApplyOpts(opts *CmdOpts) error
```

---

## 六、实现细节

### 6.1 Cmd.ApplyOpts 完整实现

```go
// ApplyOpts 应用选项到命令
func (c *Cmd) ApplyOpts(opts *CmdOpts) error {
    // 使用defer捕获panic, 转换为错误返回
    defer func() {
        if r := recover(); r != nil {
            // 将panic转换为错误
            switch x := r.(type) {
            case string:
                err = types.NewError("PANIC", x, nil)
            case error:
                err = types.NewError("PANIC", x.Error(), x)
            default:
                err = types.NewError("PANIC", fmt.Sprintf("%v", x), nil)
            }
        }
    }()

    // 验证选项
    if opts == nil {
        return types.NewError("INVALID_CMDOPTS", "cmd opts cannot be nil", nil)
    }

    // 1. 设置基本属性 - 调用现有方法
    if opts.Desc != "" {
        c.SetDesc(opts.Desc)
    }
    if opts.RunFunc != nil {
        c.SetRun(opts.RunFunc)
    }

    // 2. 设置配置选项 - 调用现有方法
    if opts.Version != "" {
        c.SetVersion(opts.Version)
    }
    c.SetChinese(opts.UseChinese)
    if opts.EnvPrefix != "" {
        c.SetEnvPrefix(opts.EnvPrefix)
    }
    if opts.UsageSyntax != "" {
        c.SetUsageSyntax(opts.UsageSyntax)
    }
    if opts.LogoText != "" {
        c.SetLogoText(opts.LogoText)
    }

    // 3. 添加示例和说明 - 调用现有方法
    if len(opts.Examples) > 0 {
        c.AddExamples(opts.Examples)
    }
    if len(opts.Notes) > 0 {
        c.AddNotes(opts.Notes)
    }

    // 4. 添加互斥组 - 调用现有方法
    for _, group := range opts.MutexGroups {
        c.AddMutexGroup(group.Name, group.Flags, group.AllowNone)
    }

    // 5. 添加子命令 - 调用现有方法
    if len(opts.SubCmds) > 0 {
        if err := c.AddSubCmds(opts.SubCmds...); err != nil {
            return types.WrapError(err, "FAILED_TO_ADD_SUBCMDS", "failed to add subcommands")
        }
    }

    return nil
}
```

### 6.2 设计说明

**为什么选择调用现有方法？**

1. **代码复用**
   - 利用现有的 `SetDesc`、`SetVersion`、`AddExamples` 等方法
   - 避免重复代码，降低维护成本

2. **行为一致性**
   - `ApplyOpts` 的行为与用户手动调用方法完全一致
   - 避免用户困惑，保持 API 的一致性

3. **封装性**
   - 现有方法内部可能有验证、事件通知等逻辑
   - 调用方法可以保留这些逻辑，避免绕过封装

4. **易于维护**
   - 如果方法内部逻辑改变，`ApplyOpts` 自动受益
   - 统一维护入口，降低 bug 风险

5. **测试友好**
   - 每个方法都有独立的测试覆盖
   - `ApplyOpts` 的测试可以更专注于组合逻辑

**性能考虑**

虽然每个方法都会加锁，但这在现代 Go 程序中通常不是瓶颈。如果确实需要优化性能，可以考虑：

1. 使用内部方法（如 `setDescLocked`）减少锁竞争
2. 批量操作时使用更粗粒度的锁

但在大多数情况下，代码质量和维护性比微小的性能提升更重要。

---

## 七、使用示例

### 7.1 基础使用

```go
package main

import (
    "log"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/internal/cmd"
)

func main() {
    // 创建命令
    myCmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)

    // 创建选项
    opts := &cmd.CmdOpts{
        Desc: "我的应用程序",
        Version: "1.0.0",
        UseChinese: true,
        EnvPrefix: "MYAPP",
    }

    // 应用选项
    if err := myCmd.ApplyOpts(opts); err != nil {
        log.Fatal(err)
    }

    // 解析参数
    if err := myCmd.Parse(os.Args[1:]); err != nil {
        log.Fatal(err)
    }
}
```

### 7.2 完整配置

```go
package main

import (
    "log"
    "os"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/internal/cmd"
    "gitee.com/MM-Q/qflag/internal/types"
)

func main() {
    // 创建命令
    myCmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)

    // 创建子命令
    listCmd := qflag.NewCmd("list", "ls", qflag.ExitOnError)
    listCmd.SetDesc("列出所有项目")

    addCmd := qflag.NewCmd("add", "a", qflag.ExitOnError)
    addCmd.SetDesc("添加新项目")

    // 创建选项
    opts := &cmd.CmdOpts{
        // 基本属性
        Desc: "我的应用程序",
        RunFunc: func(c types.Command) error {
            println("执行主命令")
            return nil
        },

        // 配置选项
        Version:     "1.0.0",
        UseChinese:  true,
        EnvPrefix:   "MYAPP",
        UsageSyntax: "myapp [options] [args...]",
        LogoText:    "MyApp v1.0.0",

        // 示例和说明
        Examples: map[string]string{
            "基本用法": "myapp --help",
            "详细模式": "myapp --verbose",
            "列出项目": "myapp list",
        },
        Notes: []string{
            "所有选项都可以通过环境变量设置",
            "使用 --help 查看详细帮助",
        },

        // 子命令和互斥组
        SubCmds: []types.Command{listCmd, addCmd},
        MutexGroups: []types.MutexGroup{
            {
                Name:      "format",
                Flags:     []string{"json", "xml"},
                AllowNone: false,
            },
        },
    }

    // 应用选项
    if err := myCmd.ApplyOpts(opts); err != nil {
        log.Fatal(err)
    }

    // 解析并执行
    if err := myCmd.ParseAndRoute(os.Args[1:]); err != nil {
        log.Fatal(err)
    }
}
```

### 7.3 部分配置

```go
package main

import (
    "log"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/internal/cmd"
)

func main() {
    // 创建命令
    myCmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)

    // 只配置部分属性
    opts := &cmd.CmdOpts{
        Desc: "我的应用程序",
        // 其他属性保持默认值
    }

    // 应用选项
    if err := myCmd.ApplyOpts(opts); err != nil {
        log.Fatal(err)
    }

    // 未设置的属性不会被修改
}
```

### 7.4 使用 NewCmdOpts

```go
package main

import (
    "log"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/internal/cmd"
)

func main() {
    // 创建命令
    myCmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)

    // 使用工厂函数创建选项
    opts := cmd.NewCmdOpts()
    opts.Desc = "我的应用程序"
    opts.Version = "1.0.0"
    opts.UseChinese = true

    // 应用选项
    if err := myCmd.ApplyOpts(opts); err != nil {
        log.Fatal(err)
    }
}
```

### 7.5 链式调用（可选）

```go
package main

import (
    "log"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/internal/cmd"
)

func main() {
    // 创建命令
    myCmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)

    // 链式配置
    if err := cmd.NewCmdOpts().
        WithDesc("我的应用程序").
        WithVersion("1.0.0").
        WithChinese(true).
        ApplyTo(myCmd); err != nil {
        log.Fatal(err)
    }
}
```

---

## 八、与现有代码的关系

### 8.1 结构对比

| 结构体 | 位置 | 用途 | 包含不可变属性 |
|--------|------|------|--------------|
| `CmdSpec` | `internal/cmd/cmdspec.go` | 创建新命令 | ✅ 是（LongName、ShortName） |
| `CmdOpts` | `internal/cmd/cmdopts.go` | 配置现有命令 | ❌ 否 |
| `types.CmdConfig` | `internal/types/config.go` | 存储命令内部配置 | ❌ 否 |

### 8.2 使用场景对比

```go
// 场景 1: 创建新命令 - 使用 CmdSpec
spec := &cmd.CmdSpec{
    LongName:  "myapp",
    ShortName: "m",
    Desc:      "我的应用程序",
    Version:   "1.0.0",
}
myCmd, err := cmd.NewCmdFromSpec(spec)

// 场景 2: 配置现有命令 - 使用 CmdOpts
myCmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)
opts := &cmd.CmdOpts{
    Desc:    "我的应用程序",
    Version: "1.0.0",
}
myCmd.ApplyOpts(opts)

// 场景 3: 访问命令配置 - 使用 types.CmdConfig
config := myCmd.Config()
version := config.Version
```

### 8.3 依赖关系

```
CmdOpts (cmd 包)
    ↓
    被 Cmd.ApplyOpts 使用
    ↓
Cmd (cmd 包)
    ↓
    使用
    ↓
types.CmdConfig (types 包)
```

### 8.4 命名呼应

| 结构体 | 方法 | 语义 |
|--------|------|------|
| `CmdOpts` | `ApplyOpts` | Apply Options = 应用选项 |

**命名优势**：
- ✅ 结构体名 `CmdOpts` 和方法名 `ApplyOpts` 形成呼应
- ✅ 语义清晰：Apply Options = 应用选项
- ✅ 易于记忆：Opts → ApplyOpts

---

## 九、注意事项

### 9.1 部分配置
- `CmdOpts` 支持部分配置，未设置的属性不会被修改
- 零值（如空字符串、空切片）不会被应用
- 如果需要重置属性为默认值，需要直接调用对应的方法

### 9.2 重复应用
- 多次调用 `ApplyOpts` 可能导致配置被覆盖
- 建议在命令初始化时一次性应用所有配置
- 如果需要增量配置，可以创建多个 `CmdOpts` 实例

### 9.3 错误处理
- `ApplyOpts` 方法会捕获 panic 并转换为错误
- 建议检查返回的错误
- 如果添加子命令失败，已应用的配置不会被回滚

### 9.4 并发安全
- `CmdOpts` 本身不是线程安全的
- `ApplyOpts` 方法内部使用 `Cmd` 的读写锁，保证并发安全
- 不要在多个 goroutine 中同时修改同一个 `CmdOpts` 实例

### 9.5 子命令处理
- 子命令会通过 `AddSubCmds` 方法添加到命令中
- 如果子命令已经存在，会返回错误
- 建议在应用之前检查子命令是否已存在

### 9.6 版本信息
- 只有根命令才能设置版本信息
- 如果对子命令调用 `ApplyOpts` 并设置 `Version`，该字段会被忽略

---

## 十、扩展性

### 10.1 添加新的配置选项

```go
// 1. 在 CmdOpts 结构体中添加新字段
type CmdOpts struct {
    // ... 现有字段
    NewField string // 新字段
}

// 2. 在 Cmd.ApplyOpts 方法中添加应用逻辑
func (c *Cmd) ApplyOpts(opts *CmdOpts) error {
    // ... 现有逻辑

    // 应用新字段
    if opts.NewField != "" {
        // 直接设置命令的属性
        c.newField = opts.NewField
    }

    return nil
}
```

### 10.2 链式调用（可选）

```go
// 可以考虑添加链式调用方法
func (o *CmdOpts) WithDesc(desc string) *CmdOpts {
    o.Desc = desc
    return o
}

func (o *CmdOpts) WithVersion(version string) *CmdOpts {
    o.Version = version
    return o
}

func (o *CmdOpts) WithChinese(useChinese bool) *CmdOpts {
    o.UseChinese = useChinese
    return o
}

// 使用示例
opts := cmd.NewCmdOpts().
    WithDesc("我的应用程序").
    WithVersion("1.0.0").
    WithChinese(true)
myCmd.ApplyOpts(opts)
```

### 10.3 配置验证（可选）

```go
// 可以添加验证方法
func (o *CmdOpts) Validate() error {
    if o.Version != "" && !isValidVersion(o.Version) {
        return types.NewError("INVALID_VERSION", "invalid version format", nil)
    }
    return nil
}

// 使用示例
if err := opts.Validate(); err != nil {
    return err
}
myCmd.ApplyOpts(opts)
```

---

## 十一、测试建议

### 11.1 单元测试

```go
func TestCmd_ApplyOpts(t *testing.T) {
    // 测试基本属性
    t.Run("Basic Properties", func(t *testing.T) {
        cmd := NewCmd("test", "t", ContinueOnError)
        opts := &CmdOpts{
            Desc: "测试命令",
            Version: "1.0.0",
        }

        err := cmd.ApplyOpts(opts)
        assert.NoError(t, err)
        assert.Equal(t, "测试命令", cmd.Desc())
        assert.Equal(t, "1.0.0", cmd.Config().Version)
    })

    // 测试部分配置
    t.Run("Partial Config", func(t *testing.T) {
        cmd := NewCmd("test", "t", ContinueOnError)
        cmd.SetDesc("原始描述")

        opts := &CmdOpts{
            Version: "1.0.0",
            // 不设置 Desc
        }

        err := cmd.ApplyOpts(opts)
        assert.NoError(t, err)
        assert.Equal(t, "原始描述", cmd.Desc()) // 未被修改
        assert.Equal(t, "1.0.0", cmd.Config().Version)
    })

    // 测试错误处理
    t.Run("Error Handling", func(t *testing.T) {
        cmd := NewCmd("test", "t", ContinueOnError)
        var opts *CmdOpts

        err := cmd.ApplyOpts(opts)
        assert.Error(t, err)
        assert.Equal(t, "INVALID_CMDOPTS", err.(*types.Error).Code)
    })

    // 测试并发安全
    t.Run("Concurrent Safety", func(t *testing.T) {
        cmd := NewCmd("test", "t", ContinueOnError)
        opts := &CmdOpts{
            Desc: "测试命令",
        }

        var wg sync.WaitGroup
        for i := 0; i < 100; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                cmd.ApplyOpts(opts)
            }()
        }
        wg.Wait()

        assert.Equal(t, "测试命令", cmd.Desc())
    })
}
```

### 11.2 集成测试

```go
func TestCmd_ApplyOpts_Integration(t *testing.T) {
    // 测试完整的配置流程
    cmd := NewCmd("test", "t", ContinueOnError)

    subCmd := NewCmd("sub", "s", ContinueOnError)
    opts := &CmdOpts{
        Desc: "测试命令",
        Version: "1.0.0",
        UseChinese: true,
        Examples: map[string]string{
            "示例1": "test --help",
        },
        SubCmds: []types.Command{subCmd},
    }

    err := cmd.ApplyOpts(opts)
    assert.NoError(t, err)

    // 验证所有属性
    assert.Equal(t, "测试命令", cmd.Desc())
    assert.Equal(t, "1.0.0", cmd.Config().Version)
    assert.True(t, cmd.Config().UseChinese)
    assert.Len(t, cmd.Config().Example, 1)
    assert.True(t, cmd.HasSubCmd("sub"))
}
```

---

## 十二、文档更新

### 12.1 API 文档

需要在以下文件中添加文档：
- `internal/cmd/APIDOC.md`
- `APIDOC.md`

### 12.2 示例代码

需要在以下文件中添加示例：
- `examples/cmdopts/` 目录

### 12.3 README 更新

在 `README.md` 中添加使用示例：
```go
## 使用 CmdOpts 配置现有命令

```go
cmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)

opts := &cmd.CmdOpts{
    Desc: "我的应用程序",
    Version: "1.0.0",
    UseChinese: true,
}

cmd.ApplyOpts(opts)
```
```

---

## 十三、总结

### 13.1 设计优势

1. **简洁性**：命名简洁，易于理解和使用
2. **语义清晰**：Opts = Options，明确表达配置选项的含义
3. **避免冲突**：不会与 `types.CmdConfig` 冲突
4. **形成对比**：与 `CmdSpec` 形成清晰的对比
5. **灵活性**：支持部分配置，易于扩展
6. **命名呼应**：`CmdOpts` 和 `ApplyOpts` 形成呼应，易于记忆
7. **语义明确**：Apply Options = 应用选项，一目了然

### 13.2 适用场景

- 已有命令实例，需要批量设置属性
- 需要结构化的配置管理
- 需要部分配置（未设置的属性不会被修改）
- 需要清晰的错误处理
- 需要语义明确的配置方式

### 13.3 后续优化

- 考虑添加链式调用方法
- 考虑添加配置验证
- 考虑添加配置合并功能
- 考虑添加配置序列化/反序列化

---

## 附录

### A. 完整代码示例

参见第七章"使用示例"

### B. 错误码参考

| 错误码 | 说明 |
|--------|------|
| INVALID_CMDOPTS | 命令选项不能为 nil |
| FAILED_TO_ADD_SUBCMDS | 添加子命令失败 |
| PANIC | 内部 panic |

### C. 相关文档

- [cmdspec.go](../internal/cmd/cmdspec.go) - CmdSpec 实现
- [cmd.go](../internal/cmd/cmd.go) - Cmd 实现
- [config.go](../internal/types/config.go) - types.CmdConfig 定义

### D. 命名对比

#### 结构体命名

| 命名 | 长度 | 语义清晰度 | 推荐度 |
|------|------|-----------|--------|
| CmdConfigurer | 12 字符 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **CmdOpts** | **7 字符** | ⭐⭐⭐⭐ | **⭐⭐⭐⭐⭐** |
| CmdOptions | 10 字符 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| CmdCustomizer | 12 字符 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

#### 方法命名

| 方法名 | 语义 | 优点 | 推荐度 |
|--------|------|------|--------|
| Configure(opts) | 配置命令 | 更简洁，符合常见习惯 | ⭐⭐⭐⭐⭐ |
| **ApplyOpts(opts)** | **应用选项到命令** | **更明确，与 CmdOpts 呼应** | **⭐⭐⭐⭐⭐** |

**最终选择：CmdOpts + ApplyOpts**

**理由**：
- ✅ 结构体名 `CmdOpts` 简洁明了
- ✅ 方法名 `ApplyOpts` 与结构体名呼应
- ✅ 语义清晰：Apply Options = 应用选项
- ✅ 易于记忆和使用
- ✅ 避免与现有类型冲突
