# QFlag-CLI开发规范指南

本文档定义了使用 qflag 库开发命令行工具的标准规范，强调简洁、统一、易维护的目录结构。

---

## 目录

1. [概述](#概述)
2. [目录结构规范](#目录结构规范)
3. [命名规范](#命名规范)
4. [根命令规范](#根命令规范)
5. [子命令规范](#子命令规范)
6. [二级子命令规范](#二级子命令规范)
7. [程序入口规范](#程序入口规范)
8. [标志类型与配置](#标志类型与配置)
9. [高级功能配置](#高级功能配置)
10. [环境变量绑定](#环境变量绑定)
11. [代码风格规范](#代码风格规范)
12. [测试规范](#测试规范)
13. [最佳实践](#最佳实践)
14. [常见问题](#常见问题)
15. [完整示例项目](#完整示例项目)
16. [开发检查清单](#开发检查清单)

---

## 概述

### 设计原则

1. **入口统一**: 程序入口固定在 `cmd/` 目录
2. **命令集中**: 所有命令定义统一放在 `internal/cli/` 目录
3. **单文件原则**: 每个命令 = 单个文件（包含：初始化 + 标志 + run 函数）
4. **按需创建**: 只有两级子命令才创建同名子目录
5. **简洁统一**: 使用统一的模式和命名规范

### 两种开发模式

根据项目复杂度，可选择以下两种模式：

| 模式 | 特点 | 适用场景 |
|------|------|----------|
| **简单模式** | 命令定义和业务逻辑在同一文件 | 简单命令、快速开发、小型项目 |
| **分离模式** | 命令定义和业务逻辑分离 | 复杂命令、需要单元测试、大型项目 |

---

## 目录结构规范

### 标准目录结构

```
your-project/
├── cmd/                          # 程序入口（固定）
│   └── yourapp/                  # 应用名称
│       └── main.go               # 唯一入口
├── internal/
│   ├── cli/                      # 所有命令定义（核心目录）
│   │   ├── root.go               # 根命令定义
│   │   ├── run.go                # 一级子命令
│   │   ├── build.go              # 一级子命令
│   │   ├── config.go             # 一级命令
│   │   └── config/               # 二级子命令目录
│   │       ├── get.go            # 二级子命令
│   │       └── set.go            # 二级子命令
│   ├── commands/                 # 业务逻辑层（分离模式使用）
│   │   └── <command>/
│   │       └── cmd_<command>.go
│   ├── utils/                    # 工具函数
│   └── config/                   # 配置相关
├── go.mod
└── README.md
```

### 目录说明

| 目录/文件 | 说明 | 必需 | 模式 |
|----------|------|------|------|
| `cmd/yourapp/main.go` | 程序入口，调用 `cli.InitAndRun()` | ✅ | 两种模式 |
| `internal/cli/root.go` | 根命令定义和全局标志 | ✅ | 两种模式 |
| `internal/cli/*.go` | 一级子命令定义 | 按需 | 简单模式 |
| `internal/cli/*/` | 二级子命令目录 | 按需 | 两种模式 |
| `internal/commands/` | 业务逻辑分离 | 按需 | 分离模式 |

---

## 命名规范

### 文件命名

| 类型 | 命名规则 | 示例 |
|------|----------|------|
| 命令目录 | 小写 | `mkdir`, `config` |
| 命令文件 | `<command>.go` | `run.go`, `build.go` |
| 业务逻辑文件 | `cmd_<command>.go` | `cmd_run.go` |
| 根命令 | `root.go` | - |

### 代码命名

| 类型 | 命名规则 | 示例 |
|------|----------|------|
| 包名 | 与目录同名，小写 | `package cli`, `package config` |
| 命令变量 | `<Command>Cmd`（必须导出） | `RunCmd`, `BuildCmd` |
| 标志变量 | `<command><Flag>`（小写前缀） | `runInput`, `buildOutput` |
| 配置结构体 | `<Command>Config` | `RunConfig`, `BuildConfig` |
| 主函数 | `<Command>CmdMain` 或 `Execute` | `RunCmdMain`, `Execute` |
| run 函数 | `run<Command>` | `runRun`, `runBuild` |
| 常量 | 全大写 | `MAX_DEPTH` |

---

## 根命令规范

### 核心区别

**根命令 ≠ 子命令**：
- 子命令：在 `init()` 中创建，导出命令对象供上级注册
- 根命令：没有上级，使用 `InitAndRun()` 模式直接操作 `qflag.Root`

### 文件位置

`internal/cli/root.go`

### 完整示例

```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

// ============================================
// 1. 全局标志变量（定义在根命令上）
// ============================================
var (
    verboseFlag    *qflag.BoolFlag
    configFlag     *qflag.StringFlag
    debugFlag      *qflag.BoolFlag
)

// ============================================
// 2. InitAndRun 初始化并运行根命令
// ============================================
// InitAndRun 初始化并运行根命令
//
// 返回值:
//   - err: 初始化或运行命令时可能发生的错误
func InitAndRun() (err error) {
    // defer 捕获 panic
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()

    // 注册根命令的标志（直接在 qflag.Root 上定义）
    verboseFlag = qflag.Root.Bool("verbose", "v", "详细输出", false)
    configFlag = qflag.Root.String("config", "c", "配置文件路径", "")
    debugFlag = qflag.Root.Bool("debug", "d", "调试模式", false)

    // 配置根命令
    rootCmdOpts := &qflag.CmdOpts{
        Desc:       "应用描述",
        Version:    "1.0.0",
        UseChinese: true,
        Completion: true,
        RunFunc:    runRoot,
        Examples: map[string]string{
            "列出所有任务": fmt.Sprintf("%s -l", qflag.Root.Name()),
            "运行指定任务": fmt.Sprintf("%s -r deploy", qflag.Root.Name()),
        },
        Notes: []string{
            "默认查找的配置文件: config.toml",
            "未指定配置文件时，按优先级查找当前目录下的配置文件",
        },
        SubCmds: []qflag.Command{
            RunCmd,     // yourapp run
            BuildCmd,   // yourapp build
            ConfigCmd,  // yourapp config
        },
    }

    // 应用根命令配置
    if err = qflag.ApplyOpts(rootCmdOpts); err != nil {
        return fmt.Errorf("apply opts failed: %w", err)
    }

    // 解析并自动路由到子命令
    if err = qflag.ParseAndRoute(); err != nil {
        return fmt.Errorf("parse and route failed: %w", err)
    }

    return nil
}

// ============================================
// 3. run() 函数：根命令的业务逻辑
// ============================================
// runRoot 是根命令的执行函数
//
// 参数:
//   - cmd: 根命令接口
//
// 返回值:
//   - error: 执行时可能遇到的错误
func runRoot(cmd qflag.Command) error {
    // 默认显示帮助
    cmd.PrintHelp()
    return nil
}
```

### 根命令要点

1. **InitAndRun() 函数**: 统一入口函数，由 `main.go` 调用
2. **直接操作 qflag.Root**: 在 `InitAndRun()` 中直接定义标志
3. **defer 捕获 panic**: 保证 panic 转换为 error 返回
4. **命名返回值**: 使用 `err` 作为命名返回值，方便 defer 和错误处理
5. **SubCmds 注册**: 所有一级子命令在这里注册
6. **自动路由**: 使用 `qflag.ParseAndRoute()` 自动路由到子命令

---

## 子命令规范

### 模式一：简单模式（定义+逻辑合一）

`internal/cli/run.go`：

```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

// ============================================
// 1. 全局命令变量（供注册到根命令）
// ============================================
var (
    // RunCmd yourapp run 命令（必须导出）
    RunCmd *qflag.Cmd
)

// ============================================
// 2. 全局标志变量
// ============================================
var (
    runInput    *qflag.StringFlag
    runOutput   *qflag.StringFlag
    runParallel *qflag.BoolFlag
)

// ============================================
// 3. init() 初始化命令、定义标志
// ============================================
func init() {
    // 初始化命令
    RunCmd = qflag.NewCmd("run", "r", qflag.ExitOnError)
    
    // 定义标志
    runInput = RunCmd.String("input", "i", "输入文件路径", "")
    runOutput = RunCmd.String("output", "o", "输出文件路径", "")
    runParallel = RunCmd.Bool("parallel", "p", "并行执行", false)
    
    // 应用命令配置
    cmdOpts := &qflag.CmdOpts{
        Desc:        "运行任务",
        UsageSyntax: fmt.Sprintf("%s run [选项] [参数]", qflag.Root.Name()),
        UseChinese:  true,
    }
    
    if err := RunCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }
    
    // 设置运行函数
    RunCmd.SetRun(runRun)
}

// ============================================
// 4. run() 函数：业务逻辑
// ============================================
// runRun 执行 run 命令
//
// 参数:
//   - cmd: 命令接口
//
// 返回值:
//   - error: 执行错误
func runRun(cmd qflag.Command) error {
    // 获取标志值
    input := runInput.Get()
    output := runOutput.Get()
    parallel := runParallel.Get()
    verbose := verboseFlag.Get() // 访问全局标志
    
    // 执行业务逻辑
    fmt.Printf("运行任务: input=%s, output=%s, parallel=%v\n", 
        input, output, parallel)
    
    return nil
}
```

### 模式二：分离模式（定义与逻辑分离）

**CLI 定义** `internal/cli/run.go`：

```go
package cli

import (
    "fmt"
    "your-project/internal/commands/run"
    "gitee.com/MM-Q/qflag"
)

var RunCmd *qflag.Cmd

var (
    runInput    *qflag.StringFlag
    runOutput   *qflag.StringFlag
    runParallel *qflag.BoolFlag
)

func init() {
    RunCmd = qflag.NewCmd("run", "r", qflag.ExitOnError)
    
    runInput = RunCmd.String("input", "i", "输入文件路径", "")
    runOutput = RunCmd.String("output", "o", "输出文件路径", "")
    runParallel = RunCmd.Bool("parallel", "p", "并行执行", false)
    
    cmdOpts := &qflag.CmdOpts{
        Desc:        "运行任务",
        UsageSyntax: fmt.Sprintf("%s run [选项] [参数]", qflag.Root.Name()),
        UseChinese:  true,
    }
    
    if err := RunCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }
    
    RunCmd.SetRun(runRun)
}

func runRun(cmd qflag.Command) error {
    config := run.Config{
        Input:   runInput.Get(),
        Output:  runOutput.Get(),
        Parallel: runParallel.Get(),
        Verbose: verboseFlag.Get(),
        Args:    cmd.Args(),
    }
    
    return run.Execute(config)
}
```

**业务逻辑** `internal/commands/run/cmd_run.go`：

```go
package run

import (
    "fmt"
)

// Config run 命令配置
type Config struct {
    Input    string
    Output   string
    Parallel bool
    Verbose  bool
    Args     []string
}

// Stats 操作统计（可选）
type Stats struct {
    Processed int
    Errors    int
}

// Execute 执行 run 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func Execute(config Config) error {
    if config.Input == "" {
        return fmt.Errorf("未指定输入文件")
    }
    
    stats := &Stats{}
    
    // 执行业务逻辑
    if config.Verbose {
        fmt.Printf("输入: %s\n", config.Input)
        fmt.Printf("输出: %s\n", config.Output)
    }
    
    // 处理逻辑...
    stats.Processed++
    
    if config.Verbose {
        fmt.Printf("处理完成: %d 个\n", stats.Processed)
    }
    
    return nil
}
```

### 文件结构要点

1. **包声明**: 所有命令文件统一使用 `package cli`
2. **命令变量**: 必须导出（首字母大写），命名为 `<Command>Cmd`
3. **标志变量**: 使用 `<command><Flag>` 格式命名，私有变量
4. **init() 函数**: 负责初始化命令、定义标志、设置配置
5. **run() 函数**: 负责业务逻辑执行，返回 error

---

## 二级子命令规范

### 目录结构

```
internal/cli/
├── config.go           # 一级命令: yourapp config
└── config/             # 二级子命令目录
    ├── get.go          # yourapp config get
    └── set.go          # yourapp config set
```

### 一级命令示例

`internal/cli/config.go`：

```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

var (
    // ConfigCmd yourapp config 命令
    ConfigCmd *qflag.Cmd
)

func init() {
    ConfigCmd = qflag.NewCmd("config", "c", qflag.ExitOnError)
    
    cmdOpts := &qflag.CmdOpts{
        Desc:        "配置管理",
        UsageSyntax: fmt.Sprintf("%s config [命令]", qflag.Root.Name()),
        UseChinese:  true,
        SubCmds: []qflag.Command{
            ConfigGetCmd, // yourapp config get
            ConfigSetCmd, // yourapp config set
        },
    }
    
    if err := ConfigCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }
    
    ConfigCmd.SetRun(runConfig)
}

func runConfig(cmd qflag.Command) error {
    cmd.PrintHelp()
    return nil
}
```

### 二级命令示例

`internal/cli/config/get.go`：

```go
package config

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

var (
    // ConfigGetCmd yourapp config get 命令
    ConfigGetCmd *qflag.Cmd
)

var (
    getKey    *qflag.StringFlag
    getGlobal *qflag.BoolFlag
)

func init() {
    ConfigGetCmd = qflag.NewCmd("get", "g", qflag.ExitOnError)
    
    getKey = ConfigGetCmd.String("key", "k", "配置键名", "")
    getGlobal = ConfigGetCmd.Bool("global", "G", "全局配置", false)
    
    cmdOpts := &qflag.CmdOpts{
        Desc:        "获取配置项",
        UsageSyntax: fmt.Sprintf("%s config get [选项]", qflag.Root.Name()),
        UseChinese:  true,
    }
    
    if err := ConfigGetCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }
    
    ConfigGetCmd.SetRun(runConfigGet)
}

func runConfigGet(cmd qflag.Command) error {
    key := getKey.Get()
    global := getGlobal.Get()
    
    fmt.Printf("获取配置: key=%s, global=%v\n", key, global)
    return nil
}
```

### 二级命令要点

1. **包名**: 使用命令名作为包名（如 `package config`）
2. **UsageSyntax**: 使用 `qflag.Root.Name()` 拼接完整命令路径
3. **注册位置**: 在一级命令的 `SubCmds` 中注册
4. **命名规范**: `一级命令 + 二级命令 + Cmd`（如 `ConfigGetCmd`）

---

## 程序入口规范

### 文件位置

`cmd/yourapp/main.go`

### 完整示例

```go
package main

import (
    "fmt"
    "os"
    "your-project/internal/cli"
)

func main() {
    if err := cli.InitAndRun(); err != nil {
        fmt.Fprintf(os.Stderr, "错误: %v\n", err)
        os.Exit(1)
    }
}
```

### 入口文件要点

1. **极简原则**: 只负责调用 `cli.InitAndRun()`
2. **错误处理**: 输出错误信息到 stderr
3. **退出码**: 错误时使用 `os.Exit(1)`

---

## 标志类型与配置

### 支持的标志类型

| 类型 | 方法 | 示例 |
|------|------|------|
| Bool | `Cmd.Bool()` | `Cmd.Bool("force", "f", "强制", false)` |
| String | `Cmd.String()` | `Cmd.String("output", "o", "输出", "")` |
| Int | `Cmd.Int()` | `Cmd.Int("count", "c", "数量", 0)` |
| Int64 | `Cmd.Int64()` | `Cmd.Int64("size", "s", "大小", 0)` |
| Uint | `Cmd.Uint()` | `Cmd.Uint("port", "p", "端口", 8080)` |
| Uint8 | `Cmd.Uint8()` | `Cmd.Uint8("level", "l", "级别", 1)` |
| Uint16 | `Cmd.Uint16()` | `Cmd.Uint16("code", "c", "代码", 100)` |
| Uint32 | `Cmd.Uint32()` | `Cmd.Uint32("id", "i", "ID", 0)` |
| Uint64 | `Cmd.Uint64()` | `Cmd.Uint64("total", "t", "总数", 0)` |
| Float64 | `Cmd.Float64()` | `Cmd.Float64("rate", "r", "比率", 0.0)` |
| Enum | `Cmd.Enum()` | `Cmd.Enum("type", "t", "类型", "md5", []string{"md5", "sha1"})` |
| Duration | `Cmd.Duration()` | `Cmd.Duration("timeout", "t", "超时", 0)` |
| Time | `Cmd.Time()` | `Cmd.Time("start", "s", "开始时间", time.Time{})` |
| Size | `Cmd.Size()` | `Cmd.Size("limit", "l", "限制", 0)` |
| StringSlice | `Cmd.StringSlice()` | `Cmd.StringSlice("tags", "t", "标签", []string{})` |
| IntSlice | `Cmd.IntSlice()` | `Cmd.IntSlice("ids", "i", "ID列表", []int{})` |
| Int64Slice | `Cmd.Int64Slice()` | `Cmd.Int64Slice("values", "v", "值列表", []int64{})` |
| Map | `Cmd.Map()` | `Cmd.Map("env", "e", "环境变量", map[string]string{})` |

### CmdOpts 配置项

| 配置项 | 类型 | 说明 | 示例 |
|--------|------|------|------|
| `Desc` | string | 命令描述 | `"创建目录"` |
| `RunFunc` | `func(Command) error` | 命令执行函数 | `run` |
| `Version` | string | 版本号（根命令） | `"1.0.0"` |
| `UseChinese` | bool | 使用中文帮助 | `true` |
| `EnvPrefix` | string | 环境变量前缀 | `"MYAPP"` |
| `UsageSyntax` | string | 命令使用语法 | `fmt.Sprintf("%s [选项]", qflag.Root.Name())` |
| `LogoText` | string | Logo文本 | `"MyApp"` |
| `Completion` | bool | 启用自动补全 | `true` |
| `DynamicCompletion` | bool | 启用动态补全 | `true` |
| `DisableFlagParsing` | bool | 禁用标志解析 | `true` |
| `Hidden` | bool | 隐藏命令 | `true` |
| `AutoBindEnv` | bool | 自动绑定环境变量 | `true` |
| `Examples` | map[string]string | 使用示例 | `map[string]string{"示例": "cmd"}` |
| `Notes` | []string | 注意事项 | `[]string{"说明1"}` |
| `SubCmds` | []Command | 子命令列表 | `[]qflag.Command{RunCmd}` |
| `MutexGroups` | []MutexGroup | 互斥组 | - |
| `RequiredGroups` | []RequiredGroup | 必需组 | - |
| `FlagDependencies` | []FlagDependency | 标志依赖 | - |

---

## 高级功能配置

### 互斥组（MutexGroup）

确保组内**最多只有一个**标志被设置。

**方式一：通过 CmdOpts 配置（推荐）**

```go
cmdOpts := &qflag.CmdOpts{
    Desc: "部署命令",
    MutexGroups: []qflag.MutexGroup{
        {
            Name:      "deploy-target",
            Flags:     []string{"dev", "staging", "prod"},
            AllowNone: true,
        },
    },
}
```

**方式二：通过 Root 直接添加**

```go
qflag.Root.AddMutexGroup("format", []string{"json", "xml", "yaml"}, false)
```

### 必需组（RequiredGroup）

确保组内**所有**标志都被设置。

**方式一：通过 CmdOpts 配置（推荐）**

```go
cmdOpts := &qflag.CmdOpts{
    Desc: "数据库配置",
    RequiredGroups: []qflag.RequiredGroup{
        {
            Name:        "database",
            Flags:       []string{"db-host", "db-port", "db-name"},
            Conditional: false,  // false = 普通必需组，true = 条件性必需组
        },
    },
}
```

**方式二：通过 Root 直接添加**

```go
// 普通必需组
qflag.Root.AddRequiredGroup("connection", []string{"host", "port"}, false)

// 条件性必需组
qflag.Root.AddRequiredGroup("database", []string{"dbhost", "dbport"}, true)
```

**普通必需组 vs 条件性必需组**：

- **普通必需组**（`Conditional: false`）：组中所有标志都必须被设置
- **条件性必需组**（`Conditional: true`）：只有设置了其中一个，才要求设置全部

### 标志依赖关系（FlagDependency）

定义标志之间的依赖约束。

**方式一：通过 CmdOpts 配置（推荐）**

```go
cmdOpts := &qflag.CmdOpts{
    Desc: "服务器配置",
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

**方式二：通过 Root 直接添加**

```go
// 必需依赖：SSL模式需要证书和密钥
qflag.Root.AddFlagDependency("ssl_requires_cert", "ssl", []string{"cert", "key"}, qflag.DepRequired)

// 互斥依赖：远程模式与本地路径互斥
qflag.Root.AddFlagDependency("remote_mutex_local", "remote", []string{"local-path"}, qflag.DepMutex)
```

**依赖类型**：

- **DepRequired**：当触发标志被设置时，目标标志必须被设置
- **DepMutex**：当触发标志被设置时，目标标志不能被设置

---

## 环境变量绑定

### 绑定方式

| 方式 | 方法 | 适用场景 | 特点 |
|------|------|----------|------|
| 手动指定 | `BindEnv("NAME")` | 需要自定义名称 | 灵活，可指定任意名称 |
| 标志自动绑定 | `AutoBindEnv()` | 单个标志自动绑定 | 使用长名称大写 |
| 命令批量绑定 | `AutoBindAllEnv()` | 批量绑定所有标志 | 一次性绑定 |
| CmdOpts 配置 | `AutoBindEnv: true` | 配置化管理 | 与其他配置一起设置 |

### 示例

```go
// 方式1：手动指定
hostFlag := cmd.String("host", "H", "主机", "localhost")
hostFlag.BindEnv("HOST_ADDR")  // 绑定到 MYAPP_HOST_ADDR

// 方式2：自动绑定
portFlag := cmd.Int("port", "p", "端口", 8080)
portFlag.AutoBindEnv()  // 绑定到 MYAPP_PORT

// 方式3：批量绑定
cmd.AutoBindAllEnv()

// 方式4：CmdOpts 配置（推荐）
opts := &qflag.CmdOpts{
    EnvPrefix:   "MYAPP",
    AutoBindEnv: true,
}
```

### 注意事项

1. **前缀设置**：使用 `SetEnvPrefix()` 或 `CmdOpts.EnvPrefix`
2. **命名规则**：环境变量名 = 前缀 + _ + 标志名（大写）
3. **优先级**：命令行参数 > 环境变量 > 默认值
4. **长名称要求**：`AutoBindEnv()` 和 `AutoBindAllEnv()` 要求标志必须有长名称

---

## 代码风格规范

### 注释规范

#### 1. 函数注释

**所有函数**（无论公有还是私有）都必须添加函数级注释：

```go
// FunctionName 函数功能简述
//
// 参数:
//   - param1: 参数1描述
//   - param2: 参数2描述
//
// 返回值:
//   - returnType1: 返回值1描述
//   - error: 错误信息
func FunctionName(param1 string, param2 int) (returnType1, error) {
    // 实现
}

// privateFunction 私有函数同样需要注释
//
// 参数:
//   - input: 输入值
//
// 返回值:
//   - string: 处理后的结果
func privateFunction(input string) string {
    // 实现
}
```

#### 2. 结构体注释

**所有结构体及其字段**都必须添加注释：

```go
// Config 应用配置结构体
// 包含应用运行所需的所有配置项
type Config struct {
    // Debug 是否启用调试模式
    // 为 true 时输出详细日志
    Debug bool
    
    // Port 服务监听端口
    // 默认为 8080
    Port int
    
    // Host 服务绑定地址
    // 默认为 "localhost"
    Host string
}

// CommandConfig 命令配置结构体
type CommandConfig struct {
    // Input 输入文件路径
    Input string
    
    // Output 输出文件路径
    Output string
    
    // Parallel 是否并行执行
    Parallel bool
}
```

#### 3. 常量注释

**所有常量**都必须添加注释：

```go
// 应用相关常量
const (
    // DefaultPort 默认服务端口
    DefaultPort = 8080
    
    // DefaultHost 默认绑定地址
    DefaultHost = "localhost"
    
    // MaxRetries 最大重试次数
    MaxRetries = 3
)

// 错误信息常量
const (
    // ErrNotFound 资源未找到错误
    ErrNotFound = "resource not found"
    
    // ErrInvalidInput 输入无效错误
    ErrInvalidInput = "invalid input"
)
```

#### 4. 变量注释

**包级变量**（尤其是导出的变量）必须添加注释：

```go
// 全局标志变量
var (
    // verboseFlag 详细输出标志
    // 控制是否输出详细日志信息
    verboseFlag *qflag.BoolFlag
    
    // configFlag 配置文件路径标志
    configFlag *qflag.StringFlag
    
    // debugFlag 调试模式标志
    debugFlag *qflag.BoolFlag
)

// RunCmd yourapp run 命令
// 用于运行指定任务
var RunCmd *qflag.Cmd
```

#### 5. 接口注释

**所有接口**及其方法都必须添加注释：

```go
// Service 业务服务接口
// 定义了业务处理的核心方法
type Service interface {
    // Process 处理数据
    //
    // 参数:
    //   - data: 输入数据
    //
    // 返回值:
    //   - string: 处理结果
    //   - error: 处理错误
    Process(data string) (string, error)
    
    // Validate 验证数据有效性
    //
    // 参数:
    //   - data: 待验证数据
    //
    // 返回值:
    //   - error: 验证错误，nil 表示验证通过
    Validate(data string) error
}
```

### 错误处理

```go
// 包装错误
if err != nil {
    return fmt.Errorf("操作失败: %w", err)
}

// 创建新错误
if len(args) == 0 {
    return fmt.Errorf("缺少必需参数")
}
```

### 错误处理策略

| 策略 | 使用场景 |
|------|---------|
| `qflag.ExitOnError` | 生产环境（默认） |
| `qflag.ContinueOnError` | 测试环境 |
| `qflag.PanicOnError` | 开发调试 |

### 导入顺序

```go
import (
    // 标准库
    "fmt"
    "os"
    
    // 第三方库
    "gitee.com/MM-Q/qflag"
    
    // 内部包
    "your-project/internal/utils"
    "your-project/internal/config"
)
```

---

## 测试规范

### 单元测试示例

```go
package cli

import (
    "testing"
    "gitee.com/MM-Q/qflag"
)

func TestRunCommand(t *testing.T) {
    // 创建测试命令
    testCmd := qflag.NewCmd("test", "t", qflag.ContinueOnError)
    
    inputFlag := testCmd.String("input", "i", "输入", "")
    
    // 解析测试参数
    args := []string{"--input", "test.txt"}
    err := testCmd.Parse(args)
    
    if err != nil {
        t.Fatalf("解析失败: %v", err)
    }
    
    // 验证标志值
    if inputFlag.Get() != "test.txt" {
        t.Errorf("期望 input=test.txt，实际得到 %s", inputFlag.Get())
    }
}

func TestMutexGroup(t *testing.T) {
    testCmd := qflag.NewCmd("test", "t", qflag.ContinueOnError)
    
    devFlag := testCmd.Bool("dev", "d", "开发环境", false)
    prodFlag := testCmd.Bool("prod", "p", "生产环境", false)
    
    cmdOpts := &qflag.CmdOpts{
        Desc: "测试互斥组",
        MutexGroups: []qflag.MutexGroup{
            {
                Name:      "env",
                Flags:     []string{"dev", "prod"},
                AllowNone: true,
            },
        },
    }
    
    if err := testCmd.ApplyOpts(cmdOpts); err != nil {
        t.Fatalf("应用配置失败: %v", err)
    }
    
    // 测试互斥标志冲突
    args := []string{"--dev", "--prod"}
    err := testCmd.Parse(args)
    
    if err == nil {
        t.Error("期望返回互斥标志冲突错误，但解析成功")
    }
}
```

### 测试覆盖率

```bash
# 运行测试并生成覆盖率报告
go test -cover ./internal/cli/...

# 生成详细覆盖率报告
go test -coverprofile=coverage.out ./internal/cli/...
go tool cover -html=coverage.out
```

---

## 最佳实践

### 1. 全局配置管理

```go
// internal/config/global.go
package config

import "sync"

var (
    once     sync.Once
    instance *GlobalConfig
)

type GlobalConfig struct {
    Verbose    bool
    Debug      bool
    ConfigPath string
}

func GetGlobalConfig() *GlobalConfig {
    once.Do(func() {
        instance = &GlobalConfig{}
    })
    return instance
}
```

### 2. 配置文件处理

```go
// internal/config/loader.go
package config

import (
    "fmt"
    "os"
    "github.com/pelletier/go-toml/v2"
)

func LoadFromFile(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("读取配置文件失败: %w", err)
    }
    
    var cfg Config
    if err := toml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("解析配置文件失败: %w", err)
    }
    
    return &cfg, nil
}
```

### 3. 标志验证器

```go
import (
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/validators"
)

// 使用内置验证器
portFlag := Cmd.Int("port", "p", "端口", 8080).
    SetValidator(validators.Range(1, 65535))

emailFlag := Cmd.String("email", "e", "邮箱", "").
    SetValidator(validators.Email())

// 自定义验证器
pathFlag := Cmd.String("path", "p", "文件路径", "").
    SetValidator(func(value string) error {
        if value == "" {
            return nil
        }
        if _, err := os.Stat(value); os.IsNotExist(err) {
            return fmt.Errorf("文件不存在: %s", value)
        }
        return nil
    })
```

### 4. 优雅的错误处理

```go
func run(cmd qflag.Command) error {
    if runTask := runFlag.Get(); runTask != "" {
        cfg, err := loadConfig(filePathFlag.Get())
        if err != nil {
            return fmt.Errorf("加载配置失败: %w", err)
        }
        
        if err := executeTask(runTask, cfg); err != nil {
            if errors.Is(err, ErrTaskNotFound) {
                fmt.Printf("任务不存在: %s\n", runTask)
                fmt.Println("使用 --list 查看所有可用任务")
                return nil
            }
            return fmt.Errorf("执行任务失败: %w", err)
        }
    }
    
    return nil
}
```

---

## 常见问题

### Q1: 如何处理可选参数？

```go
// 使用 IsSet() 判断是否设置
if optionalFlag.IsSet() {
    value := optionalFlag.Get()
} else {
    // 使用默认值或其他逻辑
}
```

### Q2: 如何实现全局标志？

```go
// internal/cli/root.go
var (
    verboseFlag *qflag.BoolFlag
    debugFlag   *qflag.BoolFlag
)

func InitAndRun() (err error) {
    verboseFlag = qflag.Root.Bool("verbose", "v", "详细输出", false)
    debugFlag = qflag.Root.Bool("debug", "d", "调试模式", false)
    // ...
}

// internal/cli/run.go
func runRun(cmd qflag.Command) error {
    verbose := verboseFlag.Get()  // 访问全局标志
    debug := debugFlag.Get()
    // ...
}
```

### Q3: 如何自定义帮助信息？

```go
cmdOpts := &qflag.CmdOpts{
    Desc: "构建项目",
    UsageSyntax: fmt.Sprintf("%s build [选项]", qflag.Root.Name()),
    Examples: map[string]string{
        "构建当前项目": fmt.Sprintf("%s build", qflag.Root.Name()),
        "构建指定平台": fmt.Sprintf("%s build --platform linux/amd64", qflag.Root.Name()),
    },
    Notes: []string{
        "默认构建当前平台的二进制文件",
        "支持的平台: linux/amd64, windows/amd64, darwin/amd64",
    },
}
```

### Q4: 如何处理子命令的参数？

```go
func runRun(cmd qflag.Command) error {
    args := cmd.Args()
    
    if len(args) == 0 {
        return fmt.Errorf("缺少必需参数")
    }
    
    taskName := cmd.Arg(0)  // 获取第一个参数
    
    for i, arg := range args {
        fmt.Printf("参数 %d: %s\n", i, arg)
    }
    
    return nil
}
```

---

## 完整示例项目

### 项目结构

```
gob/
├── cmd/
│   └── gob/
│       └── main.go
├── internal/
│   ├── cli/
│   │   ├── root.go
│   │   ├── run.go
│   │   ├── build.go
│   │   └── config/
│   │       ├── get.go
│   │       └── set.go
│   ├── utils/
│   │   └── log.go
│   └── config/
│       └── config.go
├── go.mod
└── README.md
```

### main.go

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/gob/internal/cli"
)

func main() {
    if err := cli.InitAndRun(); err != nil {
        fmt.Fprintf(os.Stderr, "错误: %v\n", err)
        os.Exit(1)
    }
}
```

### internal/cli/root.go

```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

const logoText = `
  ____ _           
 / ___| |__   ___  
| |   | '_ \ / _ \ 
| |___| | | |  __/ 
 \____|_| |_|\___| 
`

var (
    verboseFlag *qflag.BoolFlag
    configFlag  *qflag.StringFlag
)

func InitAndRun() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()

    verboseFlag = qflag.Root.Bool("verbose", "v", "详细输出", false)
    configFlag = qflag.Root.String("config", "c", "配置文件路径", "")

    rootCmdOpts := &qflag.CmdOpts{
        Version:    "1.0.0",
        Desc:       "Go 项目构建工具",
        LogoText:   logoText,
        UseChinese: true,
        Completion: true,
        RunFunc:    runRoot,
        Examples: map[string]string{
            "初始化项目": fmt.Sprintf("%s init", qflag.Root.Name()),
            "运行任务":   fmt.Sprintf("%s run -i main.go", qflag.Root.Name()),
        },
        SubCmds: []qflag.Command{
            RunCmd,
            BuildCmd,
            ConfigCmd,
        },
    }

    if err = qflag.ApplyOpts(rootCmdOpts); err != nil {
        return fmt.Errorf("apply opts failed: %w", err)
    }

    if err = qflag.ParseAndRoute(); err != nil {
        return fmt.Errorf("parse and route failed: %w", err)
    }

    return nil
}

func runRoot(cmd qflag.Command) error {
    cmd.PrintHelp()
    return nil
}
```

---

## 开发检查清单

在完成命令行工具开发后，检查以下项目：

- [ ] 目录结构符合规范
- [ ] 命名符合规范（文件、变量、函数）
- [ ] 根命令使用 `InitAndRun()` 模式
- [ ] 子命令导出 `Cmd` 变量
- [ ] 标志变量使用正确的前缀命名
- [ ] 子命令正确注册到父命令的 `SubCmds`
- [ ] 全局标志定义在根命令上
- [ ] 添加了必要的验证器（如需要）
- [ ] 配置了互斥组/必需组（如需要）
- [ ] 配置了标志依赖关系（如需要）
- [ ] 配置了环境变量绑定（如需要）
- [ ] 启用了自动补全 `Completion: true`
- [ ] 添加了函数级注释
- [ ] 编译通过 `go build ./...`
- [ ] 编写了单元测试
- [ ] 更新了 README 文档

---

## 核心原则总结

1. ✅ **入口统一**: `cmd/yourapp/main.go` 只调用 `cli.InitAndRun()`
2. ✅ **命令集中**: 所有命令定义放在 `internal/cli/`
3. ✅ **单文件原则**: 每个命令 = 单个文件（简单模式）
4. ✅ **按需创建**: 只有两级子命令才创建目录
5. ✅ **根命令特殊**: 使用 `InitAndRun()` 模式，直接操作 `qflag.Root`
6. ✅ **命名统一**: 遵循 `<Command>Cmd`、`<command><Flag>`、`run<Command>` 规范
7. ✅ **注释完整**: 所有公共函数添加函数级注释
8. ✅ **错误处理**: 使用 `fmt.Errorf("...: %w", err)` 包装错误

---

## 附录

### qflag 库信息

- **仓库**: `gitee.com/MM-Q/qflag`
- **版本要求**: v0.5.17+

### 相关文档

- [qflag 使用指南](https://gitee.com/MM-Q/qflag)
- [FLAG_USAGE.md](FLAG_USAGE.md) - 标志使用语法
