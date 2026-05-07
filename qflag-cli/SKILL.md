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

```
your-project/
├── cmd/yourapp/
│   └── main.go          # 入口：调用 cli.InitAndRun()
├── internal/cli/
│   ├── root.go          # 根命令定义
│   ├── <command>.go     # 一级子命令
│   └── <command>/       # 二级子命令目录（按需）
├── go.mod
└── README.md
```

### 3. 最小可用示例

**cmd/yourapp/main.go:**
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

**internal/cli/root.go:**
```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

var (
    verboseFlag *qflag.BoolFlag
)

func InitAndRun() error {
    verboseFlag = qflag.Root.Bool("verbose", "v", "详细输出", false)
    
    opts := &qflag.CmdOpts{
        Desc:       "MyApp",
        Version:    "1.0.0",
        UseChinese: true,
        Completion: true,
        SubCmds:    []qflag.Command{RunCmd},
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
var runInput *qflag.StringFlag

func init() {
    RunCmd = qflag.NewCmd("run", "r", qflag.ExitOnError)
    runInput = RunCmd.String("input", "i", "输入文件", "")
    
    RunCmd.ApplyOpts(&qflag.CmdOpts{
        Desc:       "运行任务",
        UseChinese: true,
    })
    RunCmd.SetRun(func(cmd qflag.Command) error {
        fmt.Println("Running:", runInput.Get())
        return nil
    })
}
```

## 核心概念速查

### 标志类型

| 类型 | 方法 | 示例 |
|------|------|------|
| String | `Cmd.String()` | `cmd.String("output", "o", "输出", "")` |
| Bool | `Cmd.Bool()` | `cmd.Bool("force", "f", "强制", false)` |
| Int | `Cmd.Int()` | `cmd.Int("count", "c", "数量", 1)` |
| Enum | `Cmd.Enum()` | `cmd.Enum("format", "t", "格式", "json", []string{"json", "xml"})` |
| Duration | `Cmd.Duration()` | `cmd.Duration("timeout", "t", "超时", time.Second*30)` |
| StringSlice | `Cmd.StringSlice()` | `cmd.StringSlice("tags", "t", "标签", nil)` |

### 错误处理策略

```go
qflag.ContinueOnError  // 继续解析并返回错误
qflag.ExitOnError      // 解析错误时退出程序（推荐）
qflag.PanicOnError     // 解析错误时触发 panic
```

## 高级特性速查

### 验证器

```go
import "gitee.com/MM-Q/qflag/validators"

port.SetValidator(validators.IntRange(1, 65535))
email.SetValidator(validators.Email())
config.SetValidator(validators.FileExists())
```

### 互斥组

```go
opts := &qflag.CmdOpts{
    MutexGroups: []qflag.MutexGroup{
        {
            Name:      "env",
            Flags:     []string{"dev", "prod"},
            AllowNone: true,
        },
    },
}
```

### 必需组

```go
opts := &qflag.CmdOpts{
    RequiredGroups: []qflag.RequiredGroup{
        {
            Name:        "db",
            Flags:       []string{"host", "port"},
            Conditional: false,  // false=普通必需, true=条件性必需
        },
    },
}
```

### 标志依赖

```go
opts := &qflag.CmdOpts{
    FlagDependencies: []qflag.FlagDependency{
        {
            Name:    "ssl_requires_cert",
            Trigger: "ssl",
            Targets: []string{"cert", "key"},
            Type:    qflag.DepRequired,  // DepRequired 或 DepMutex
        },
    },
}
```

### 环境变量绑定

```go
// 方式1: 手动指定
flag.BindEnv("MY_VAR")

// 方式2: 自动绑定（使用标志名大写）
flag.AutoBindEnv()

// 方式3: 批量绑定
cmd.AutoBindAllEnv()

// 方式4: CmdOpts 配置
opts := &qflag.CmdOpts{
    EnvPrefix:   "MYAPP",
    AutoBindEnv: true,
}
```

### 自动补全

```go
opts := &qflag.CmdOpts{
    Completion:        true,  // 启用自动补全
    DynamicCompletion: true,  // 启用动态补全（可选）
}

// 生成补全脚本
// yourapp --completion bash > /etc/bash_completion.d/yourapp
```

## 参考文档

详细规范、完整示例和最佳实践请参见：

- **[QFlag-CLI开发规范指南](references/QFlag-CLI开发规范指南.md)** - 完整开发规范、目录结构、命名规范、代码模板
- **[标志使用语法指南](references/FLAG_USAGE.md)** - 各种标志类型的使用语法详解
- **[完整示例](references/examples.md)** - 5个完整示例项目（文件处理、HTTP服务器、数据库迁移、多子命令工具等）
