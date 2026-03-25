# QFlag 完整示例

## 示例1: 文件处理工具

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/validators"
)

func main() {
    // 定义标志
    input := qflag.Root.String("input", "i", "输入文件", "")
    output := qflag.Root.String("output", "o", "输出文件", "")
    format := qflag.Root.Enum("format", "f", "输出格式", "txt", []string{"txt", "json", "csv"})
    verbose := qflag.Root.Bool("verbose", "v", "详细输出", false)
    
    // 添加验证器
    input.SetValidator(validators.FileExists())
    
    // 配置命令
    qflag.Root.SetDesc("文件格式转换工具")
    qflag.Root.SetVersion("1.0.0")
    
    opts := &qflag.CmdOpts{
        UseChinese: true,
        Completion: true,
        Examples: map[string]string{
            "转换文本文件":   "filetool -i input.txt -o output.json -f json",
            "详细模式转换":   "filetool -i data.csv -o result.txt -v",
        },
        Notes: []string{
            "支持 txt、json、csv 格式互转",
            "输出文件如果不存在会自动创建",
        },
    }
    
    if err := qflag.ApplyOpts(opts); err != nil {
        fmt.Println("配置错误:", err)
        os.Exit(1)
    }
    
    // 解析参数
    if err := qflag.Parse(); err != nil {
        os.Exit(1)
    }
    
    // 执行业务逻辑
    if verbose.Get() {
        fmt.Printf("输入: %s\n", input.Get())
        fmt.Printf("输出: %s\n", output.Get())
        fmt.Printf("格式: %s\n", format.Get())
    }
    
    // 转换文件...
    fmt.Println("转换完成!")
}
```

## 示例2: HTTP 服务器

```go
package main

import (
    "fmt"
    "net/http"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/validators"
)

func main() {
    // 定义标志
    host := qflag.Root.String("host", "h", "服务器地址", "0.0.0.0")
    port := qflag.Root.Int("port", "p", "服务器端口", 8080)
    dir := qflag.Root.String("dir", "d", "服务目录", ".")
    cors := qflag.Root.Bool("cors", "c", "启用 CORS", false)
    
    // 验证器
    port.SetValidator(validators.IntRange(1, 65535))
    
    // 配置
    qflag.Root.SetDesc("简单的静态文件服务器")
    qflag.Root.SetVersion("1.0.0")
    
    opts := &qflag.CmdOpts{
        UseChinese: true,
        Completion: true,
        Examples: map[string]string{
            "默认启动":     "httpserver",
            "指定端口":     "httpserver -p 3000",
            "指定目录":     "httpserver -d /var/www -p 80",
        },
    }
    
    if err := qflag.ApplyOpts(opts); err != nil {
        fmt.Println("配置错误:", err)
        return
    }
    
    if err := qflag.Parse(); err != nil {
        return
    }
    
    // 启动服务器
    addr := fmt.Sprintf("%s:%d", host.Get(), port.Get())
    fmt.Printf("服务器启动: http://%s\n", addr)
    
    handler := http.FileServer(http.Dir(dir.Get()))
    if cors.Get() {
        handler = withCORS(handler)
    }
    
    http.ListenAndServe(addr, handler)
}

func withCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next.ServeHTTP(w, r)
    })
}
```

## 示例3: 数据库迁移工具

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 数据库连接标志
    dbHost := qflag.Root.String("db-host", "", "数据库主机", "localhost")
    dbPort := qflag.Root.Int("db-port", "", "数据库端口", 3306)
    dbUser := qflag.Root.String("db-user", "u", "数据库用户", "")
    dbPass := qflag.Root.String("db-pass", "p", "数据库密码", "")
    dbName := qflag.Root.String("db-name", "", "数据库名称", "")
    
    // 迁移配置
    migrateUp := qflag.Root.Bool("up", "", "向上迁移", false)
    migrateDown := qflag.Root.Bool("down", "", "向下迁移", false)
    version := qflag.Root.Int("version", "v", "指定版本", 0)
    
    // 配置
    opts := &qflag.CmdOpts{
        Desc:       "数据库迁移工具",
        UseChinese: true,
        Completion: true,
        MutexGroups: []qflag.MutexGroup{
            {
                Name:      "direction",
                Flags:     []string{"up", "down"},
                AllowNone: false,  // 必须选择一个方向
            },
        },
        RequiredGroups: []qflag.RequiredGroup{
            {
                Name:        "connection",
                Flags:       []string{"db-host", "db-port", "db-user", "db-pass", "db-name"},
                Conditional: false,
            },
        },
        Examples: map[string]string{
            "向上迁移":     "migrate --up -u root -p secret --db-name myapp",
            "向下迁移":     "migrate --down -u root -p secret --db-name myapp",
            "指定版本":     "migrate -v 5 -u root -p secret --db-name myapp",
        },
    }
    
    if err := qflag.ApplyOpts(opts); err != nil {
        fmt.Println("配置错误:", err)
        return
    }
    
    if err := qflag.Parse(); err != nil {
        return
    }
    
    // 执行迁移
    fmt.Printf("连接到 %s@%s:%d/%s\n", dbUser.Get(), dbHost.Get(), dbPort.Get(), dbName.Get())
    
    if migrateUp.Get() {
        fmt.Println("执行向上迁移...")
    } else if migrateDown.Get() {
        fmt.Println("执行向下迁移...")
    }
    
    if version.Get() > 0 {
        fmt.Printf("目标版本: %d\n", version.Get())
    }
}
```

## 示例4: 多子命令工具（Git 风格）

```
mygit/
├── cmd/
│   └── mygit/
│       └── main.go
├── internal/
│   └── cli/
│       ├── root.go
│       ├── init.go
│       ├── add.go
│       ├── commit.go
│       ├── push.go
│       └── remote/
│           ├── add.go
│           └── remove.go
```

**internal/cli/root.go:**
```go
package cli

import "gitee.com/MM-Q/qflag"

func InitAndRun() error {
    opts := &qflag.CmdOpts{
        Desc:       "MyGit - 简化版 Git 工具",
        Version:    "1.0.0",
        UseChinese: true,
        Completion: true,
        SubCmds: []qflag.Command{
            InitCmd,
            AddCmd,
            CommitCmd,
            PushCmd,
            RemoteCmd,
        },
    }
    
    if err := qflag.ApplyOpts(opts); err != nil {
        return err
    }
    
    return qflag.ParseAndRoute()
}
```

**internal/cli/init.go:**
```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

var InitCmd *qflag.Cmd

var initBare *qflag.BoolFlag

func init() {
    InitCmd = qflag.NewCmd("init", "", qflag.ExitOnError)
    
    initBare = InitCmd.Bool("bare", "", "创建裸仓库", false)
    
    opts := &qflag.CmdOpts{
        Desc:       "初始化 Git 仓库",
        UseChinese: true,
    }
    
    if err := InitCmd.ApplyOpts(opts); err != nil {
        panic(err)
    }
    
    InitCmd.SetRun(func(cmd qflag.Command) error {
        if initBare.Get() {
            fmt.Println("创建裸仓库...")
        } else {
            fmt.Println("初始化仓库...")
        }
        return nil
    })
}
```

**internal/cli/commit.go:**
```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/validators"
)

var CommitCmd *qflag.Cmd

var (
    commitMessage *qflag.StringFlag
    commitAmend   *qflag.BoolFlag
    commitAll     *qflag.BoolFlag
)

func init() {
    CommitCmd = qflag.NewCmd("commit", "c", qflag.ExitOnError)
    
    commitMessage = CommitCmd.String("message", "m", "提交信息", "")
    commitAmend = CommitCmd.Bool("amend", "", "修改上次提交", false)
    commitAll = CommitCmd.Bool("all", "a", "提交所有修改", false)
    
    // 验证提交信息不为空
    commitMessage.SetValidator(func(value string) error {
        if value == "" {
            return fmt.Errorf("提交信息不能为空，请使用 -m 指定")
        }
        return nil
    })
    
    opts := &qflag.CmdOpts{
        Desc:       "提交更改",
        UseChinese: true,
        Examples: map[string]string{
            "普通提交": "mygit commit -m \"修复 bug\"",
            "提交所有": "mygit commit -am \"更新\"",
            "修改提交": "mygit commit --amend -m \"新信息\"",
        },
    }
    
    if err := CommitCmd.ApplyOpts(opts); err != nil {
        panic(err)
    }
    
    CommitCmd.SetRun(func(cmd qflag.Command) error {
        fmt.Printf("提交: %s\n", commitMessage.Get())
        if commitAmend.Get() {
            fmt.Println("修改上次提交")
        }
        if commitAll.Get() {
            fmt.Println("包含所有修改")
        }
        return nil
    })
}
```

**internal/cli/remote/add.go:**
```go
package remote

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

var RemoteAddCmd *qflag.Cmd

var (
    addName *qflag.StringFlag
    addURL  *qflag.StringFlag
)

func init() {
    RemoteAddCmd = qflag.NewCmd("add", "", qflag.ExitOnError)
    
    addName = RemoteAddCmd.String("name", "n", "远程名称", "origin")
    addURL = RemoteAddCmd.String("url", "u", "远程地址", "")
    
    opts := &qflag.CmdOpts{
        Desc:       "添加远程仓库",
        UseChinese: true,
    }
    
    if err := RemoteAddCmd.ApplyOpts(opts); err != nil {
        panic(err)
    }
    
    RemoteAddCmd.SetRun(func(cmd qflag.Command) error {
        fmt.Printf("添加远程仓库: %s -> %s\n", addName.Get(), addURL.Get())
        return nil
    })
}
```

**cmd/mygit/main.go:**
```go
package main

import (
    "os"
    "mygit/internal/cli"
)

func main() {
    if err := cli.InitAndRun(); err != nil {
        os.Exit(1)
    }
}
```

## 示例5: 带配置文件的复杂工具

```go
package main

import (
    "fmt"
    "os"
    "gopkg.in/yaml.v3"
    "gitee.com/MM-Q/qflag"
)

type AppConfig struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Log struct {
        Level  string `yaml:"level"`
        Format string `yaml:"format"`
    } `yaml:"log"`
}

func loadConfig(path string) (*AppConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var cfg AppConfig
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}

func main() {
    configFile := qflag.Root.String("config", "c", "配置文件", "config.yaml")
    
    if err := qflag.Parse(); err != nil {
        return
    }
    
    // 加载配置
    cfg, err := loadConfig(configFile.Get())
    if err != nil {
        fmt.Println("加载配置失败:", err)
        cfg = &AppConfig{}
    }
    
    // 定义标志（使用配置文件的值作为默认值）
    host := qflag.Root.String("host", "h", "服务器地址", cfg.Server.Host)
    port := qflag.Root.Int("port", "p", "服务器端口", cfg.Server.Port)
    logLevel := qflag.Root.Enum("log-level", "", "日志级别", cfg.Log.Level, []string{"debug", "info", "warn", "error"})
    
    // 重新解析（命令行参数覆盖配置文件）
    if err := qflag.Parse(); err != nil {
        return
    }
    
    fmt.Printf("Server: %s:%d\n", host.Get(), port.Get())
    fmt.Printf("Log Level: %s\n", logLevel.Get())
}
```
