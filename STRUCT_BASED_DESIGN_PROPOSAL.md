# 提案：为 qflag 引入声明式（结构体）定义模式

## 1. 摘要

本提案建议在 `qflag` 库中增加对声明式 API 的支持，允许开发者通过定义 Go 结构体来描述命令行标志和命令。这将作为对现有命令式 API (`cmd.String(...)`) 的补充，旨在提升复杂应用的开发体验、可读性和可维护性，使其与 `Cobra`、`urfave/cli` 等现代 CLI 库的设计范式保持一致。

## 2. 必要性与优势分析

当前 `qflag` 的 API 是命令式的，这种方式在简单场景下直观易用。然而，随着应用复杂度的增加，其缺点也逐渐显现：

- **代码分散**：标志的定义、变量的声明、默认值的设置可能散落在代码的不同位置，不易于集中管理和审查。
- **可读性下降**：当一个命令包含大量标志时，一长串的 `cmd.String(...)`、`cmd.Int(...)` 调用会显得非常臃肿和重复。
- **数据绑定繁琐**：通常需要先声明一个变量，然后通过 `Var` 系列函数将其与标志绑定，增加了模板代码。

引入**结构体声明式定义**可以有效解决这些问题，其核心优势如下：

- **集中化与声明式**：将所有相关的标志聚合到一个结构体中，**代码即文档**。开发者一眼就能看出一个命令支持哪些配置项，以及它们的类型、默认值和用途。
- **自动数据绑定**：通过 Go 的结构体标签（Struct Tag），可以实现将命令行参数值自动解析并填充到结构体字段中，极大地简化了数据绑定过程。
- **强类型与可复用**：配置项被组织在强类型的结构体中，易于在程序的不同模块间安全地传递和复用，也方便进行单元测试。
- **符合现代 Go 语言习惯**：`json`、`yaml`、`gorm` 等大量流行库都采用结构体标签来处理数据映射。为 `qflag` 引入这种模式能让 Go 开发者感到非常亲切，降低学习成本。

**结论**：支持结构体定义**非常有必要**。它是一个强大的功能补充，能让 `qflag` 在处理复杂应用时更具竞争力，吸引更广泛的开发者。

## 3. 实现方案设计

我们建议分两步实现此功能，以确保平稳过渡和逐步迭代。

### 方案一：通过结构体标签定义和绑定标志 (核心功能)

这是最核心、最直接的改进。我们可以新增一个 `Cmd.AddFlags(interface{})` 方法，它接收一个结构体指针，通过反射来解析标签并注册标志。

#### 设计示例

用户可以这样定义他们的配置：

```go
// file: config.go
package main

// AppConfig 定义了应用的所有命令行标志
type AppConfig struct {
    // 通用配置
    Verbose    bool   `qflag:"name=verbose,short=v,usage='Enable verbose mode'"`
    ConfigFile string `qflag:"name=config,usage='Path to config file',default='config.yaml'"`

    // 服务器配置 (内嵌结构体)
    Server struct {
        Host string `qflag:"name=host,usage='Server host',default='127.0.0.1'"`
        Port int    `qflag:"name=port,usage='Server port',default=8080"`
    }

    // 数据库配置 (使用前缀)
    Database struct {
        DSN     string `qflag:"name=dsn,usage='Database DSN'"`
        Timeout int64  `qflag:"name=timeout,usage='Connection timeout',default=5"`
    } `qflag:"prefix=db."` // 为内嵌结构体的所有标志添加 "db." 前缀
}
```

#### 使用方式

```go
// file: main.go
package main

import (
    "fmt"
    "github.com/your/qflag"
    "os"
)

func main() {
    cmd := qflag.New("myapp", ...)
    
    // 声明配置结构体实例
    var config AppConfig

    // 新增方法，自动解析结构体并注册所有标志
    err := cmd.AddFlags(&config)
    if err != nil {
        panic(err)
    }

    // 正常解析
    cmd.Parse(os.Args[1:])

    // 解析后，可以直接使用填充好数据的结构体
    if config.Verbose {
        fmt.Println("Verbose mode is enabled.")
    }
    fmt.Printf("Server running on: %s:%d\n", config.Server.Host, config.Server.Port)
    fmt.Printf("Database DSN: %s\n", config.Database.DSN)
    
    // 命令行示例：./myapp --verbose --port=9090 --db.dsn="user:pass@..."
}
```

#### 实现步骤

1.  **定义 `qflag` 标签格式**：设计一套清晰的标签规则，例如 `qflag:"name=...,short=...,usage=...,default=...,prefix=..."`。
2.  **创建 `Cmd.AddFlags(interface{})` 方法**：
    - 该方法接收一个 `interface{}`，并使用 `reflect` 包验证它是否为指向结构体的指针。
    - 遍历结构体的所有字段，解析其 `qflag` 标签。
    - 根据字段的类型（`string`, `int`, `bool`, `time.Duration` 等），调用对应的 `cmd.StringVar`, `cmd.IntVar`, `cmd.BoolVar` 等方法。关键是将字段的地址传递给这些 `Var` 方法，从而实现自动绑定。
    - 支持递归解析内嵌的结构体，并能正确处理 `prefix` 标签，实现标志名称的层级化。
3.  **错误处理**：在反射解析过程中，如果遇到标签格式错误、字段类型不支持等情况，应返回详细的错误信息。

### 方案二：完整的命令结构体定义 (进阶功能)

在方案一的基础上，我们可以更进一步，用结构体来完整描述一个命令，包括它的子命令和要执行的动作。

#### 设计示例

```go
// 定义 qflag.Command 结构体
type Command struct {
    Name        string
    ShortName   string
    Usage       string
    Description string
    Flags       interface{} // 指向一个标志结构体的指针，用方案一的方法解析
    Subcommands []*Command  // 子命令列表
    Action      func(c *Cmd) error // 命令执行的动作
}

// 用户使用 Command 结构体定义整个应用
var rootCommand = &Command{
    Name:  "docker",
    Usage: "A self-sufficient runtime for containers",
    Flags: &DockerFlags{}, // 假设 DockerFlags 是一个标志结构体
    Subcommands: []*Command{
        {
            Name:  "run",
            Usage: "Run a command in a new container",
            Flags: &RunFlags{}, // 'run' 命令的标志
            Action: func(c *Cmd) error {
                // 'docker run' 的业务逻辑
                return nil
            },
        },
        {
            Name:  "ps",
            Usage: "List containers",
            Flags: &PsFlags{},
            Action: func(c *Cmd) error {
                // 'docker ps' 的业务逻辑
                return nil
            },
        },
    },
}

// 启动应用
func main() {
    // NewAppFromCommand 会将 Command 结构体递归地“编译”成 qflag 内部的 Cmd 树
    app := qflag.NewAppFromCommand(rootCommand)
    app.Run(os.Args)
}
```

#### 实现步骤

1.  **定义 `qflag.Command` 结构体**：如上所示，包含命令的元数据、标志、子命令和动作。
2.  **创建 `NewAppFromCommand(*Command)` 函数**：
    - 这是一个高级构造函数，它接收用户定义的 `Command` 树作为蓝图。
    - 它会递归地遍历 `Command` 结构体，为每一个 `Command` 创建一个 `cmd.Cmd` 实例。
    - 调用方案一中实现的 `cmd.AddFlags(cmd.Flags)` 来注册标志。
    - 将 `Action` 函数包装成 `cmd.Cmd` 的一部分，以便在解析后执行。
    - 递归处理 `Subcommands`，并将它们添加为 `cmd.Cmd` 的子命令。
    - 最终返回一个完全配置好的、可直接运行的根 `*cmd.Cmd` 对象。

## 4. 结论

通过上述两步方案，`qflag` 将能够同时支持命令式和声明式两种 API 风格，满足从简单到复杂各种场景的需求。我们建议优先实现方案一，因为它能以最小的成本带来最大的体验提升。在此基础上，再实现方案二，为大型应用的构建提供一个更加优雅和结构化的解决方案。