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
    Completion: true,  // 启用自动补全
}

// 生成补全脚本
// yourapp --completion bash > /etc/bash_completion.d/yourapp
// yourapp --completion pwsh > yourapp-completion.ps1
```

### 环境变量绑定

通过 `BindEnv()` 设置要解析的环境变量名：

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

**规则**: 设置 `EnvPrefix` 后，`BindEnv("NAME")` 会解析 `前缀_NAME` 格式的环境变量。

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
3. **添加示例**: 在 `Examples` 中提供常用用法示例
4. **验证输入**: 使用验证器确保输入有效性
5. **合理分组**: 使用互斥组和必需组管理标志关系
