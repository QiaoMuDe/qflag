---
name: qflag-cli
description: |
  使用 qflag 库开发 Go 语言命令行工具的专业技能。
  适用于需要创建 CLI 工具的场景，包括：
  (1) 创建新的命令行工具项目
  (2) 添加子命令和标志
  (3) 配置互斥组和必需组
  (4) 实现验证器和自动补全
  (5) 生成帮助文档和版本信息
  当用户使用 qflag 库或需要开发 Go CLI 工具时触发此技能。
---

# QFlag CLI 开发技能

## 快速开始

### 1. 创建新项目

```bash
# 初始化 Go 模块
go mod init your-project

# 添加 qflag 依赖
go get gitee.com/MM-Q/qflag
```

### 2. 基础项目结构

参见 [命令行工具开发规范](references/cli-tool-dev-spec.md) 获取完整目录结构规范。

```
your-project/
├── cmd/                          # 程序入口（固定）
│   └── yourapp/
│       └── main.go               # 唯一入口：调用 cli.InitAndRun()
├── internal/
│   ├── cli/                      # 所有命令定义（核心目录）
│   │   ├── root.go               # 根命令定义
│   │   ├── <command>.go          # 一级子命令
│   │   └── <command>/            # 二级子命令目录（按需创建）
│   │       └── <subcommand>.go
│   ├── utils/                    # 工具函数
│   └── service/                  # 业务服务
├── go.mod
└── README.md
```

### 3. 最小可用示例

**cmd/yourapp/main.go:**
```go
package main

import (
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 定义标志
    name := qflag.Root.String("name", "n", "用户名", "guest")
    verbose := qflag.Root.Bool("verbose", "v", "详细模式", false)
    
    // 解析参数
    if err := qflag.Parse(); err != nil {
        os.Exit(1)
    }
    
    // 使用标志值
    if verbose.Get() {
        println("详细模式已启用")
    }
    println("Hello, " + name.Get())
}
```

## 核心概念

### 标志类型

| 类型 | 方法 | 示例 |
|------|------|------|
| String | `Cmd.String()` | `cmd.String("output", "o", "输出文件", "")` |
| Bool | `Cmd.Bool()` | `cmd.Bool("force", "f", "强制", false)` |
| Int | `Cmd.Int()` | `cmd.Int("count", "c", "数量", 1)` |
| Enum | `Cmd.Enum()` | `cmd.Enum("format", "t", "格式", "json", []string{"json", "xml"})` |
| Duration | `Cmd.Duration()` | `cmd.Duration("timeout", "t", "超时", time.Second*30)` |
| Size | `Cmd.Size()` | `cmd.Size("limit", "l", "限制", 1024)` |
| StringSlice | `Cmd.StringSlice()` | `cmd.StringSlice("tags", "t", "标签", nil)` |
| Map | `Cmd.Map()` | `cmd.Map("env", "e", "环境变量", map[string]string{})` |

### 错误处理策略

```go
qflag.ContinueOnError  // 继续解析并返回错误
qflag.ExitOnError      // 解析错误时退出程序（推荐）
qflag.PanicOnError     // 解析错误时触发 panic
```

## 开发模式

### 模式1: 全局根命令（简单工具）

适用于单一功能的简单 CLI 工具。

```go
package main

import (
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 直接在 Root 上定义标志
    port := qflag.Root.Int("port", "p", "服务端口号", 8080)
    host := qflag.Root.String("host", "h", "服务器地址", "localhost")
    
    // 配置根命令
    qflag.Root.SetDesc("简单的 HTTP 服务器")
    qflag.Root.SetVersion("1.0.0")
    
    // 解析
    if err := qflag.Parse(); err != nil {
        return
    }
    
    // 使用
    println("Server:", host.Get(), ":", port.Get())
}
```

### 模式2: 子命令模式（复杂工具）

适用于具有多个子命令的 CLI 工具（如 git、docker）。

**目录结构:**
```
internal/cli/
├── root.go      # 根命令
├── run.go       # run 子命令
├── build.go     # build 子命令
└── config/      # config 子命令目录
    ├── get.go
    └── set.go
```

**internal/cli/root.go:**
```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

func InitAndRun() error {
    // 配置根命令
    opts := &qflag.CmdOpts{
        Desc:       "MyApp - 示例应用",
        Version:    "1.0.0",
        UseChinese: true,
        Completion: true,
        SubCmds: []qflag.Command{
            RunCmd,
            BuildCmd,
        },
    }
    
    if err := qflag.ApplyOpts(opts); err != nil {
        return err
    }
    
    return qflag.ParseAndRoute()
}
```

**internal/cli/run.go:**
```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

var RunCmd *qflag.Cmd

var (
    runInput  *qflag.StringFlag
    runOutput *qflag.StringFlag
)

func init() {
    RunCmd = qflag.NewCmd("run", "r", qflag.ExitOnError)
    
    runInput = RunCmd.String("input", "i", "输入文件", "")
    runOutput = RunCmd.String("output", "o", "输出文件", "")
    
    opts := &qflag.CmdOpts{
        Desc:       "运行任务",
        UseChinese: true,
    }
    
    if err := RunCmd.ApplyOpts(opts); err != nil {
        panic(err)
    }
    
    RunCmd.SetRun(func(cmd qflag.Command) error {
        fmt.Println("Running with:", runInput.Get(), "->", runOutput.Get())
        return nil
    })
}
```

**cmd/myapp/main.go:**
```go
package main

import (
    "os"
    "your-project/internal/cli"
)

func main() {
    if err := cli.InitAndRun(); err != nil {
        os.Exit(1)
    }
}
```

## 高级特性

### 验证器

```go
import (
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/validators"
)

// 范围验证
port := cmd.Int("port", "p", "端口号", 8080)
port.SetValidator(validators.IntRange(1, 65535))

// 正则验证
email := cmd.String("email", "e", "邮箱", "")
email.SetValidator(validators.MatchRegex(`^[\w.-]+@[\w.-]+\.\w+$`))

// 文件验证
config := cmd.String("config", "c", "配置文件", "")
config.SetValidator(validators.FileExists())

// 自定义验证
name := cmd.String("name", "n", "名称", "")
name.SetValidator(func(value string) error {
    if len(value) < 3 {
        return fmt.Errorf("名称长度必须大于3")
    }
    return nil
})
```

### 互斥组

确保组内最多只有一个标志被设置。

```go
opts := &qflag.CmdOpts{
    Desc: "部署命令",
    MutexGroups: []qflag.MutexGroup{
        {
            Name:      "environment",
            Flags:     []string{"dev", "staging", "prod"},
            AllowNone: true,  // 允许都不设置
        },
    },
}
```

### 必需组

确保组内所有标志都被设置。

```go
opts := &qflag.CmdOpts{
    Desc: "数据库连接",
    RequiredGroups: []qflag.RequiredGroup{
        {
            Name:        "database",
            Flags:       []string{"host", "port", "user", "password"},
            Conditional: false,  // 普通必需组
        },
        {
            Name:        "auth",
            Flags:       []string{"username", "password"},
            Conditional: true,   // 条件性必需：设置了其中一个就必须设置全部
        },
    },
}
```

### 自动补全

```go
opts := &qflag.CmdOpts{
    Completion:              true,  // 启用自动补全
    DynamicCompletion: true,  // 启用动态补全（可选）
}

// 生成补全脚本
// yourapp --completion bash > /etc/bash_completion.d/yourapp
// yourapp --completion pwsh > yourapp-completion.ps1
```

**动态补全**：将跨平台补全逻辑统一到内部 `__complete` 子命令实现，提升一致性并降低脚本体积。需先启用 `Completion` 才能启用 `DynamicCompletion`。

### 禁用标志解析

当需要将所有参数（包括 `--flag` 形式）作为位置参数处理时，可以禁用标志解析。

```go
// 方式1: 通过 CmdOpts 配置
opts := &qflag.CmdOpts{
    Desc:               "执行外部命令",
    DisableFlagParsing: true,  // 禁用标志解析
}

// 方式2: 通过方法设置
cmd := qflag.NewCmd("exec", "e", qflag.ExitOnError)
cmd.SetDisableFlagParsing(true)

// 使用场景：透传参数给外部命令
// yourapp exec -- kubectl get pods -n default
// 所有参数（包括 --flag）都会作为位置参数传递给 kubectl
```

### 隐藏命令

隐藏命令不会显示在帮助信息的子命令列表中，但仍可通过命令行正常调用。

```go
// 方式1: 通过 CmdOpts 配置
opts := &qflag.CmdOpts{
    Desc:   "调试命令",
    Hidden: true,  // 隐藏命令
}

// 方式2: 通过方法设置
cmd := qflag.NewCmd("debug", "d", qflag.ExitOnError)
cmd.SetHidden(true)

// 使用场景：
// - 内部调试命令
// - 已弃用但仍需兼容的命令
// - 高级或实验性功能
```

### 环境变量绑定

QFlag 提供了四种环境变量绑定方式，满足不同场景需求。

#### 方式一：手动指定环境变量名

通过 `BindEnv()` 方法手动指定环境变量名称：

```go
port := cmd.Int("port", "p", "服务端口号", 8080)

// 不设置前缀：解析 SERVER_PORT 环境变量
port.BindEnv("SERVER_PORT")

// 设置前缀后：解析前缀_环境变量名
opts := &qflag.CmdOpts{
    EnvPrefix: "MYAPP",
}
port.BindEnv("SERVER_PORT")  // 解析 MYAPP_SERVER_PORT
```

#### 方式二：标志自动绑定

通过 `AutoBindEnv()` 方法自动使用标志长名称的大写形式作为环境变量名：

```go
cmd := qflag.NewCmd("run", "r", qflag.ExitOnError)
cmd.SetEnvPrefix("MYAPP")

hostFlag := cmd.String("host", "H", "主机地址", "localhost")
portFlag := cmd.Int("port", "p", "端口号", 8080)

// 自动绑定：host -> MYAPP_HOST, port -> MYAPP_PORT
hostFlag.AutoBindEnv()
portFlag.AutoBindEnv()
```

#### 方式三：命令批量自动绑定

通过 `AutoBindAllEnv()` 方法一次性为命令的所有标志自动绑定环境变量：

```go
cmd := qflag.NewCmd("run", "r", qflag.ExitOnError)
cmd.SetEnvPrefix("MYAPP")

cmd.String("host", "H", "主机地址", "localhost")
cmd.Int("port", "p", "端口号", 8080)
cmd.String("user", "u", "用户名", "admin")

// 批量自动绑定所有标志
cmd.AutoBindAllEnv()
```

#### 方式四：通过 CmdOpts 配置

在 `CmdOpts` 中设置 `AutoBindEnv` 字段：

```go
cmd := qflag.NewCmd("run", "r", qflag.ExitOnError)

cmd.String("host", "H", "主机地址", "localhost")
cmd.Int("port", "p", "端口号", 8080)

cmdOpts := &qflag.CmdOpts{
    Desc:        "运行服务",
    EnvPrefix:   "MYAPP",
    AutoBindEnv: true,  // 自动绑定所有标志的环境变量
    UseChinese:  true,
}

cmd.ApplyOpts(cmdOpts)
```

#### 四种方式对比

| 方式 | 方法 | 适用场景 | 特点 |
|------|------|----------|------|
| 手动指定 | `BindEnv("NAME")` | 需要自定义环境变量名 | 灵活，可指定任意名称 |
| 标志自动绑定 | `AutoBindEnv()` | 单个标志自动绑定 | 使用长名称大写，简洁 |
| 命令批量绑定 | `AutoBindAllEnv()` | 批量绑定所有标志 | 一次性绑定，高效 |
| CmdOpts 配置 | `AutoBindEnv: true` | 配置化管理 | 与其他配置一起设置 |

#### 环境变量绑定注意事项

1. **前缀设置**：使用 `SetEnvPrefix()` 或 `CmdOpts.EnvPrefix` 设置环境变量前缀
2. **命名规则**：环境变量名 = 前缀 + _ + 标志名（大写）
3. **优先级**：命令行参数 > 环境变量 > 默认值
4. **长名称要求**：`AutoBindEnv()` 和 `AutoBindAllEnv()` 要求标志必须有长名称，否则会 panic

## 完整示例

参见 [references/examples.md](references/examples.md) 获取完整示例代码。

## 参考文档

### 核心参考
- [标志类型详解](references/flag-types.md) - 15+ 种标志类型的完整说明
- [验证器列表](references/validators.md) - 内置验证器和自定义验证器
- [设计模式](references/patterns.md) - 命令组织结构设计模式
- [完整示例](references/examples.md) - 5 个完整示例项目

### 开发规范
- [命令行工具开发规范](references/cli-tool-dev-spec.md) - 目录结构、文件组织、命名规范
- [命令开发规范](references/command-dev-spec.md) - 命令开发详细规范、代码模板

### 官方资源
- [API 文档](https://gitee.com/MM-Q/qflag) - qflag 官方 API 文档

## 最佳实践

1. **使用中文帮助**: `UseChinese: true` 提供更友好的中文帮助
2. **启用自动补全**: `Completion: true` 提升用户体验
3. **启用动态补全**: `DynamicCompletion: true` 统一补全逻辑，降低脚本体积
4. **添加示例**: 在 `Examples` 中提供常用用法示例
5. **验证输入**: 使用验证器确保输入有效性
6. **合理分组**: 使用互斥组和必需组管理标志关系
