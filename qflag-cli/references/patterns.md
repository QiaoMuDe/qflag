# QFlag 设计模式

## 命令组织结构模式

### 模式1: 单层命令（简单工具）

适用于功能单一的 CLI 工具。

```
mytool/
├── cmd/
│   └── mytool/
│       └── main.go
└── go.mod
```

**main.go:**
```go
package main

import (
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 直接定义标志
    input := qflag.Root.String("input", "i", "输入文件", "")
    output := qflag.Root.String("output", "o", "输出文件", "")
    
    qflag.Root.SetDesc("文件转换工具")
    qflag.Root.SetVersion("1.0.0")
    
    if err := qflag.Parse(); err != nil {
        return
    }
    
    // 执行业务逻辑
    convert(input.Get(), output.Get())
}
```

### 模式2: 多层子命令（复杂工具）

适用于具有多个子命令的 CLI 工具（如 git、docker）。

```
mytool/
├── cmd/
│   └── mytool/
│       └── main.go
├── internal/
│   └── cli/
│       ├── root.go
│       ├── run.go
│       ├── build.go
│       └── config/
│           ├── get.go
│           └── set.go
└── go.mod
```

**internal/cli/root.go:**
```go
package cli

import "gitee.com/MM-Q/qflag"

func InitAndRun() error {
    opts := &qflag.CmdOpts{
        Desc:       "MyTool - 开发工具集",
        Version:    "1.0.0",
        UseChinese: true,
        Completion: true,
        SubCmds: []qflag.Command{
            RunCmd,
            BuildCmd,
            ConfigCmd,
        },
    }
    
    if err := qflag.ApplyOpts(opts); err != nil {
        return err
    }
    
    return qflag.ParseAndRoute()
}
```

### 模式3: 业务逻辑分离

将 CLI 定义与业务逻辑分离，便于测试和维护。

```
mytool/
├── internal/
│   ├── cli/           # CLI 定义层
│   │   ├── root.go
│   │   └── deploy.go
│   └── commands/      # 业务逻辑层
│       └── deploy/
│           └── cmd_deploy.go
```

**internal/commands/deploy/cmd_deploy.go:**
```go
package deploy

import "fmt"

type DeployConfig struct {
    Env      string
    Force    bool
    Timeout  int
}

func DeployCmdMain(config DeployConfig) error {
    fmt.Printf("部署到 %s 环境\n", config.Env)
    // 业务逻辑...
    return nil
}
```

**internal/cli/deploy.go:**
```go
package cli

import (
    "gitee.com/MM-Q/qflag"
    "mytool/internal/commands/deploy"
)

var DeployCmd *qflag.Cmd

var (
    deployEnv     *qflag.StringFlag
    deployForce   *qflag.BoolFlag
    deployTimeout *qflag.IntFlag
)

func init() {
    DeployCmd = qflag.NewCmd("deploy", "d", qflag.ExitOnError)
    
    deployEnv = DeployCmd.String("env", "e", "部署环境", "dev")
    deployForce = DeployCmd.Bool("force", "f", "强制部署", false)
    deployTimeout = DeployCmd.Int("timeout", "t", "超时时间", 300)
    
    DeployCmd.SetRun(func(cmd qflag.Command) error {
        config := deploy.DeployConfig{
            Env:     deployEnv.Get(),
            Force:   deployForce.Get(),
            Timeout: deployTimeout.Get(),
        }
        return deploy.DeployCmdMain(config)
    })
}
```

## 配置管理模式

### 使用配置文件

```go
package main

import (
    "os"
    "gopkg.in/yaml.v3"
    "gitee.com/MM-Q/qflag"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        URL string `yaml:"url"`
    } `yaml:"database"`
}

func main() {
    configFile := qflag.Root.String("config", "c", "配置文件", "config.yaml")
    
    if err := qflag.Parse(); err != nil {
        return
    }
    
    // 加载配置文件
    var cfg Config
    data, err := os.ReadFile(configFile.Get())
    if err == nil {
        yaml.Unmarshal(data, &cfg)
    }
    
    // 命令行参数覆盖配置文件
    host := qflag.Root.String("host", "h", "服务器地址", cfg.Server.Host)
    port := qflag.Root.Int("port", "p", "服务器端口", cfg.Server.Port)
    
    // 重新解析以获取覆盖后的值
    qflag.Parse()
    
    println("Host:", host.Get())
    println("Port:", port.Get())
}
```

## 中间件模式

### 前置处理

```go
func WithLogging(next func(qflag.Command) error) func(qflag.Command) error {
    return func(cmd qflag.Command) error {
        println("开始执行命令...")
        err := next(cmd)
        if err != nil {
            println("执行失败:", err)
        } else {
            println("执行成功")
        }
        return err
    }
}

// 使用
cmd.SetRun(WithLogging(func(cmd qflag.Command) error {
    // 业务逻辑
    return nil
}))
```

### 权限检查

```go
func WithAuth(requiredRole string) func(func(qflag.Command) error) func(qflag.Command) error {
    return func(next func(qflag.Command) error) func(qflag.Command) error {
        return func(cmd qflag.Command) error {
            // 检查权限
            if !hasRole(requiredRole) {
                return fmt.Errorf("权限不足，需要 %s 角色", requiredRole)
            }
            return next(cmd)
        }
    }
}
```

## 插件模式

```go
// 定义插件接口
type Plugin interface {
    Name() string
    Register(cmd *qflag.Cmd)
}

// 插件注册表
var plugins []Plugin

func RegisterPlugin(p Plugin) {
    plugins = append(plugins, p)
}

// 在根命令中加载所有插件
func InitAndRun() error {
    for _, p := range plugins {
        p.Register(qflag.Root)
    }
    // ...
}
```
