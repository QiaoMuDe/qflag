# qflag API文档

## 概述
qflag是一个Go语言命令行参数解析库，提供了比标准库flag更丰富的功能，包括长短标志绑定、子命令支持、自动帮助信息生成等特性。

## 核心类型

### Flag接口
`Flag`是所有标志类型的通用接口，定义了标志的元数据访问方法。

```go
type Flag interface {
    Name() string       // 获取标志的名称
    ShortName() string  // 获取标志的短名称
    Usage() string      // 获取标志的用法
    Type() FlagType     // 获取标志类型
    getDefaultAny() any // 获取默认值(内部使用)
}
```

### TypedFlag接口
`TypedFlag`是带类型的标志接口，继承自`Flag`并提供类型化的默认值和值访问方法。

```go
type TypedFlag[T any] interface {
    Flag
    GetDefault() T // 获取标志的默认值
    GetValue() T   // 获取标志的实际值
    SetValue(T)    // 设置标志的值
}
```

### 具体标志类型
qflag提供以下具体标志类型，均实现了`TypedFlag`接口：

- `StringFlag`: 字符串类型标志
- `IntFlag`: 整数类型标志
- `BoolFlag`: 布尔类型标志
- `FloatFlag`: 浮点数类型标志

每个标志类型都有对应的`GetDefault()`、`GetValue()`和`SetValue()`方法。

### Cmd结构体
`Cmd`是qflag库的核心结构体，实现了`Command`接口，用于管理命令行标志和子命令。

#### 主要字段
- `fs *flag.FlagSet`: 底层flag集合，处理参数解析
- `name string`: 命令名称
- `shortName string`: 命令短名称
- `description string`: 命令描述
- `usage string`: 自定义帮助内容
- `subCmds []*Cmd`: 子命令列表

#### 主要方法

##### Args
```go
func (c *Cmd) Args() []string
```
获取非标志参数切片。

返回值:
- `[]string`: 非标志参数切片

##### Arg
```go
func (c *Cmd) Arg(i int) string
```
获取指定索引的非标志参数。

参数:
- `i`: 参数索引

返回值:
- `string`: 指定索引的参数值，若索引无效则返回空字符串

##### NArg
```go
func (c *Cmd) NArg() int
```
获取非标志参数的数量。

返回值:
- `int`: 非标志参数数量

##### NFlag
```go
func (c *Cmd) NFlag() int
```
获取已设置的标志数量。

返回值:
- `int`: 已设置的标志数量

##### FlagExists
```go
func (c *Cmd) FlagExists(name string) bool
```
检查指定名称的标志是否存在。

参数:
- `name`: 标志名称（长标志或短标志）

返回值:
- `bool`: 标志是否存在

##### NewCmd
```go
func NewCmd(name string, shortName string, errorHandling flag.ErrorHandling) *Cmd
```
创建新的命令实例。

参数:
- `name`: 命令名称
- `shortName`: 命令短名称
- `errorHandling`: 错误处理方式(flag.ContinueOnError, flag.ExitOnError, flag.PanicOnError)

返回值:
- `*Cmd`: 新创建的命令实例

##### AddSubCmd
```go
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error
```
为当前命令添加子命令。

参数:
- `subCmds`: 一个或多个子命令

返回值:
- `error`: 如果检测到循环引用或nil子命令则返回错误

##### Parse
```go
func (c *Cmd) Parse(args []string) error
```
解析命令行参数。

参数:
- `args`: 命令行参数切片

返回值:
- `error`: 解析过程中遇到的错误

##### PrintUsage
```go
func (c *Cmd) PrintUsage()
```
打印命令的帮助信息，优先使用自定义帮助内容，否则自动生成。

## 标志操作函数

### 字符串标志

#### String
```go
func String(name, shortName, defValue, usage string) *StringFlag
```
创建字符串类型标志(全局默认命令)。

参数:
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

返回值:
- `*StringFlag`: 字符串标志实例

#### StringVar
```go
func StringVar(p *string, name, shortName, defValue, usage string)
```
绑定字符串类型标志到指针(全局默认命令)。

参数:
- `p`: 指向字符串变量的指针
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

### 整数标志

#### Int
```go
func Int(name, shortName string, defValue int, usage string) *IntFlag
```
创建整数类型标志(全局默认命令)。

参数:
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

返回值:
- `*IntFlag`: 整数标志实例

#### IntVar
```go
func IntVar(p *int, name, shortName string, defValue int, usage string)
```
绑定整数类型标志到指针(全局默认命令)。

参数:
- `p`: 指向整数变量的指针
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

### 布尔标志

#### Bool
```go
func Bool(name, shortName string, defValue bool, usage string) *BoolFlag
```
创建布尔类型标志(全局默认命令)。

参数:
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

返回值:
- `*BoolFlag`: 布尔标志实例

#### BoolVar
```go
func BoolVar(p *bool, name, shortName string, defValue bool, usage string)
```
绑定布尔类型标志到指针(全局默认命令)。

参数:
- `p`: 指向布尔变量的指针
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

### 浮点数标志

#### Float
```go
func Float(name, shortName string, defValue float64, usage string) *FloatFlag
```
创建浮点数类型标志(全局默认命令)。

参数:
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

返回值:
- `*FloatFlag`: 浮点数标志实例

#### FloatVar
```go
func FloatVar(p *float64, name, shortName string, defValue float64, usage string)
```
绑定浮点数类型标志到指针(全局默认命令)。

参数:
- `p`: 指向浮点数变量的指针
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

## 内置标志

qflag自动绑定了以下内置标志：

### 帮助标志
- 长标志: `--help`
- 短标志: `-h`
- 功能: 显示命令的帮助信息
- 使用示例: `myapp --help` 或 `myapp -h`

### 显示安装路径标志
- 长标志: `--show-install-path`
- 短标志: `-sip`
- 功能: 显示应用程序的安装路径
- 使用示例: `myapp --show-install-path` 或 `myapp -sip`

## 高级特性

### 线程安全
qflag库使用`sync.Mutex`和`sync.Once`确保所有标志操作和解析过程是线程安全的，可以在并发环境中安全使用。

### 循环引用检测
添加子命令时，qflag会自动检测命令间的循环引用，避免出现无限递归的命令结构。

### 动态帮助信息生成
qflag会根据命令和标志的定义自动生成格式化的帮助信息，包括命令描述、标志说明、子命令列表等。

## 使用示例

### 基本用法
```go
package main

import (
  "fmt"
  "qflag"
)

func main() {
  // 创建字符串标志
  nameFlag := qflag.String("name", "n", "", "Your name")
  // 创建整数标志
  ageFlag := qflag.Int("age", "a", 0, "Your age")
  // 创建布尔标志
  verboseFlag := qflag.Bool("verbose", "v", false, "Verbose output")

  // 解析命令行参数
  if err := qflag.Parse(); err != nil {
    fmt.Println("Error parsing flags:", err)
    return
  }

  // 使用标志值
  fmt.Printf("Hello, %s! You are %d years old.\n", nameFlag.GetValue(), ageFlag.GetValue())
  if verboseFlag.GetValue() {
    fmt.Println("Verbose mode enabled")
    // 访问默认值示例
    fmt.Printf("Default verbose value: %v\n", verboseFlag.GetDefault())
  }
}
```

### 自定义帮助信息
```go
package main

import (
  "fmt"
  "qflag"
  "flag"
)

func main() {
  // 创建自定义命令
  cmd := qflag.NewCmd("greet", "g", flag.ExitOnError)
  cmd.SetDescription("A simple greeting program")
  
  // 设置自定义帮助信息
  cmd.SetUsage(`Usage: greet [options] <name>

A simple program to greet people.

Options:
  -n, --name <name>   Your name (required)
  -a, --age <age>     Your age
  -h, --help          Show this help message`)
  
  // 创建标志
  name := cmd.String("name", "n", "", "Your name")
  age := cmd.Int("age", "a", 0, "Your age")
  
  // 解析参数
  if err := cmd.Parse(os.Args[1:]); err != nil {
    fmt.Println("Error parsing flags:", err)
    return
  }
  
  // 显示帮助信息（当--help/-h被使用时自动调用）
  if *cmd.helpFlag {
    cmd.PrintUsage()
    return
  }
  
  // 检查必填参数
  if *name == "" {
    fmt.Println("Error: name is required")
    cmd.PrintUsage()
    return
  }
  
  // 执行逻辑
  fmt.Printf("Hello, %s!", *name)
  if *age > 0 {
    fmt.Printf(" You are %d years old.", *age)
  }
  fmt.Println()
}
```

### 内置标志说明
```go
package main

import (
  "fmt"
  "qflag"
  "os"
)

func main() {
  // 创建标志
  name := qflag.String("name", "n", "", "Your name")
  age := qflag.Int("age", "a", 0, "Your age")

  // 解析命令行参数
  if err := qflag.Parse(); err != nil {
    fmt.Println("Error parsing flags:", err)
    os.Exit(1)
  }

  // 检查是否请求帮助
  if *qflag.QCommandLine.helpFlag {
    qflag.PrintUsage()
    return
  }

  // 检查是否请求显示安装路径
  if *qflag.QCommandLine.showInstallPathFlag {
    fmt.Println("Installation path: C:\Program Files\MyApp")
    return
  }

  // 使用标志值
  fmt.Printf("Hello, %s! You are %d years old.\n", *name, *age)
}
```

使用示例:
```bash
# 显示帮助信息
myapp -h
myapp --help

# 显示安装路径
myapp -sip
myapp --show-install-path

# 正常运行
myapp -n Alice -a 30
```

### 错误处理示例
```go
package main

import (
  "fmt"
  "qflag"
  "flag"
)

func main() {
  // 创建自定义命令并设置错误处理方式
  cmd := qflag.NewCmd("myapp", "ma", flag.ContinueOnError)
  cmd.Int("port", "p", 8080, "Server port")

  // 解析参数
  err := cmd.Parse(os.Args[1:])
  if err != nil {
    // 处理解析错误
    if err == flag.ErrHelp {
      // 帮助信息已显示，无需额外处理
      return
    }
    fmt.Printf("Error parsing command line arguments: %v\n", err)
    cmd.PrintUsage()
    os.Exit(1)
  }

  // 程序逻辑...
}
```

## API摘要

| 类别 | 名称 | 描述 |
|------|------|------|
| 命令创建 | NewCmd | 创建新的命令实例 |
| 标志操作 | String, Int, Bool, Float | 创建指定类型的标志 |
| 标志绑定 | StringVar, IntVar, BoolVar, FloatVar | 绑定标志到变量指针 |
| 参数解析 | Parse | 解析命令行参数 |
| 帮助信息 | PrintUsage | 打印命令帮助信息 |
| 子命令管理 | AddSubCmd | 为命令添加子命令 |
| 参数访问 | Args, Arg, NArg | 获取非标志参数 |
| 标志查询 | NFlag, FlagExists | 查询标志状态 |

## 参数优先级规则

当长标志和短标志同时使用时，后指定的标志将覆盖先指定的标志。例如：
```bash
myapp --name Alice -n Bob
```
上述命令中，最终的name值为"Bob"，因为短标志`-n`在长标志`--name`之后指定。

## 使用注意事项

1. 标志名称和短名称不能为空，且不能与内置标志冲突
2. 添加子命令时要避免循环引用
3. Parse方法应只调用一次
4. 布尔类型标志不需要显式赋值， presence即为true
5. 使用自定义错误处理方式时，需手动处理帮助标志和错误信息

### 子命令用法
```go
package main

import (
  "fmt"
  "os"
  "qflag"
  "flag"
)

func main() {
  // 创建主命令
  rootCmd := qflag.NewCmd("myapp", "ma", flag.ExitOnError)
  rootCmd.SetDescription("My application description")
  rootCmd.String("log-level", "l", "info", "Set log level (debug, info, warn, error)")

  // 创建子命令
  serverCmd := qflag.NewCmd("server", "s", flag.ExitOnError)
  serverCmd.SetDescription("Start the server")
  serverCmd.Int("port", "p", 8080, "Server port")
  serverCmd.Bool("tls", "t", false, "Enable TLS")

  // 创建嵌套子命令
  configCmd := qflag.NewCmd("config", "c", flag.ExitOnError)
  configCmd.SetDescription("Configure the server")
  configCmd.String("file", "f", "config.json", "Config file path")

  // 构建命令层级
  serverCmd.AddSubCmd(configCmd)
  rootCmd.AddSubCmd(serverCmd)

  // 解析命令行参数
  if err := rootCmd.Parse(os.Args[1:]); err != nil {
    fmt.Println("Error parsing flags:", err)
    return
  }

  // 处理命令
  if len(rootCmd.Args()) > 0 {
    switch rootCmd.Args()[0] {
    case "server":
      fmt.Printf("Starting server on port %d\n", serverCmd.Arg(0))
      if len(serverCmd.Args()) > 0 && serverCmd.Args()[0] == "config" {
        fmt.Printf("Using config file: %s\n", configCmd.Arg(0))
      }
    }
  }
}
```

### 内置标志说明
```go
package main

import (
  "fmt"
  "qflag"
  "flag"
)

func main() {
  // 创建主命令
  rootCmd := qflag.NewCmd("myapp", "ma", flag.ExitOnError)
  rootCmd.SetDescription("My application description")

  // 创建子命令
  subCmd := qflag.NewCmd("sub", "s", flag.ExitOnError)
  subCmd.SetDescription("Subcommand description")
  subCmd.String("config", "c", "config.json", "Config file path")

  // 添加子命令
  rootCmd.AddSubCmd(subCmd)

  // 解析命令行参数
  if err := rootCmd.Parse(os.Args[1:]); err != nil {
    fmt.Println("Error parsing flags:", err)
    return
  }

  // 处理子命令
  if len(rootCmd.Args()) > 0 && rootCmd.Args()[0] == "sub" {
    fmt.Println("Running subcommand")
    // 可以通过subCmd访问子命令的标志
  }
}
```