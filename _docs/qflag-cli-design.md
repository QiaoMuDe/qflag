# qflag-cli 工具设计方案

> 基于 `qflag命令开发规范.md` 和 `qflag命令行工具开发规范.md` 设计

---

## 一、工具定位

### 1.1 核心目标

`qflag-cli` 是将 `qflag` 开发规范转化为可执行工具的脚手架，实现：

- **规范即代码** — 强制遵循最佳实践，减少人为错误
- **零配置启动** — 一行命令生成完整可运行的项目
- **交互式开发** — 向导式添加命令、标志和验证规则
- **双向验证** — 代码生成 + 规范检查，确保项目健康

### 1.2 目标用户

| 用户类型 | 需求 | 使用场景 |
|---------|------|----------|
| Go 新手 | 快速上手 qflag | `qflag-cli init` 生成模板 |
| 团队开发者 | 统一规范 | `qflag-cli validate` 检查合规性 |
| 开源作者 | 生成文档 | `qflag-cli gen docs` |

---

## 二、功能架构

```
qflag-cli
├── init          # 项目初始化
├── add           # 添加命令/标志/组
├── remove        # 删除命令/标志
├── list          # 列出项目结构
├── validate      # 验证规范合规性
├── generate      # 生成代码/文档/补全脚本
└── config        # 项目配置管理
```

---

## 三、命令详细设计

### 3.1 init - 项目初始化

**功能**: 创建新的 qflag 项目，生成符合规范的目录结构和示例代码。

**模式选择**:

| 模式 | 适用场景 | 目录结构 |
|------|----------|----------|
| `simple` | 小型工具、单文件命令 | 扁平结构，业务逻辑内联 |
| `layered` | 大型项目、多团队协作 | 分层结构，commands + cli 分离 |

**交互流程**:

```bash
$ qflag-cli init myapp

? 选择项目模式: layered
? Go 模块名: gitee.com/user/myapp
? 应用描述: 我的 CLI 工具
? 启用自动补全: Yes
? 启用动态补全: Yes
? 错误处理方式: ExitOnError

✓ 项目 myapp 创建成功
✓ 目录结构已生成
✓ 示例命令 'hello' 已创建
✓ 运行 'cd myapp && go run ./cmd/myapp hello' 测试
```

**生成的目录结构 (layered 模式)**:

```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go              # 入口文件
├── internal/
│   ├── cli/                     # CLI 层 (按规范)
│   │   ├── root.go              # 根命令
│   │   ├── hello.go             # 示例命令
│   │   └── completion.go        # 补全集成
│   └── commands/                # 业务逻辑层
│       └── hello/
│           └── run.go           # 命令执行逻辑
├── pkg/
│   └── utils/                   # 通用工具
├── scripts/
│   ├── install-completion.sh    # 补全安装脚本
│   └── build.sh                 # 构建脚本
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── .qflag-cli.yaml              # 项目配置 (可选)
```

**生成的目录结构 (simple 模式)**:

```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go
├── internal/
│   └── cli/
│       ├── root.go
│       └── hello.go             # 业务逻辑内联
├── go.mod
├── go.sum
└── README.md
```

---

### 3.2 add - 添加组件

#### 3.2.1 add cmd - 添加命令

**功能**: 添加新的子命令，自动更新父命令的 SubCmds 注册。

**用法示例**:

```bash
# 快捷添加
$ qflag-cli add cmd build --desc="构建项目" --short=b

# 多级命令（自动创建目录结构）
$ qflag-cli add cmd config/get --desc="获取配置值"
$ qflag-cli add cmd config/set --desc="设置配置值"

# 交互式添加（推荐）
$ qflag-cli add cmd

? 命令名称: deploy
? 短名称: d
? 描述: 部署应用到服务器
? 父命令: root
? 添加标志? Yes

  ? 标志名称: env
  ? 短名称: e
  ? 类型: enum
  ? 可选值: dev,staging,prod
  ? 默认值: dev
  ? 描述: 部署环境
  ? 必需? No
  ? 继续添加标志? Yes

  ? 标志名称: verbose
  ? 短名称: v
  ? 类型: bool
  ? 描述: 显示详细输出
  ? 继续添加标志? No

? 添加互斥组? Yes
  ? 组名: strategy
  ? 互斥标志: rolling,blue-green
  ? 允许都不设置? Yes

? 添加必需组? No

✓ 命令 'deploy' 已创建
✓ 文件: internal/cli/deploy.go
✓ 文件: internal/commands/deploy/run.go
✓ 已注册到 internal/cli/root.go
```

**生成的代码示例** (layered 模式):

```go
// internal/cli/deploy.go
package cli

import (
    "gitee.com/MM-Q/qflag"
    "gitee.com/user/myapp/internal/commands/deploy"
)

var (
    // deployCmd 部署命令
    deployCmd = qflag.NewCmd("deploy", "d", qflag.ExitOnError)

    // deployEnv 部署环境
    deployEnv *qflag.EnumFlag[string]

    // deployVerbose 详细输出
    deployVerbose *qflag.BoolFlag
)

func init() {
    deployCmd.SetDesc("部署应用到服务器")
    deployCmd.SetUsageSyntax("[OPTIONS]")

    // 注册标志
    deployEnv = deployCmd.Enum("env", "e", []string{"dev", "staging", "prod"}, "dev", "部署环境")
    deployVerbose = deployCmd.Bool("verbose", "v", false, "显示详细输出")

    // 注册互斥组
    deployCmd.AddMutexGroup("strategy", []string{"rolling", "blue-green"}, false)

    // 注册执行函数
    deployCmd.Action(deploy.Run)
}
```

```go
// internal/commands/deploy/run.go
package deploy

import (
    "fmt"
    "gitee.com/MM-Q/qflag/types"
)

// Run 执行部署命令
//
// 参数:
//   - ctx: 命令执行上下文
//
// 返回值:
//   - error: 执行过程中的错误
func Run(ctx types.CmdContext) error {
    // 获取标志值
    env := ctx.Flags().Get("env").Value().(string)
    verbose := ctx.Flags().Get("verbose").Value().(bool)

    fmt.Printf("部署到环境: %s\n", env)
    if verbose {
        fmt.Println("详细模式已启用")
    }

    // TODO: 实现部署逻辑
    return nil
}
```

---

#### 3.2.2 add flag - 添加标志

**功能**: 为现有命令添加标志，支持所有 qflag 标志类型。

**用法示例**:

```bash
# 基础标志
$ qflag-cli add flag build --name=output --short=o --type=string --default="dist"

# 带验证器的标志
$ qflag-cli add flag server --name=port --short=p --type=int --validator=port

# 枚举类型
$ qflag-cli add flag deploy --name=env --type=enum --values=dev,staging,prod --default=dev

# 切片类型
$ qflag-cli add flag build --name=tags --type=strings --desc="构建标签"

# 计数器类型
$ qflag-cli add flag test --name=verbose --short=v --type=count --desc="详细级别"
```

**支持的标志类型**:

| 类型 | 说明 | 示例 |
|------|------|------|
| `string` | 字符串 | `--name=value` |
| `strings` | 字符串切片 | `--tag=a --tag=b` |
| `int` | 整数 | `--count=10` |
| `ints` | 整数切片 | `--port=80 --port=443` |
| `int64` | 64位整数 | `--size=1073741824` |
| `float64` | 浮点数 | `--rate=0.95` |
| `bool` | 布尔值 | `--force` |
| `count` | 计数器 | `-vvv` |
| `enum` | 枚举值 | `--env=prod` |
| `enums` | 枚举切片 | `--color=red --color=blue` |
| `duration` | 时间间隔 | `--timeout=5m` |
| `ip` | IP地址 | `--bind=127.0.0.1` |
| `ips` | IP地址切片 | `--allow=10.0.0.1 --allow=10.0.0.2` |
| `url` | URL | `--endpoint=https://api.example.com` |
| `filepath` | 文件路径 | `--config=/etc/app.conf` |
| `filepath-exists` | 存在的文件 | `--input=data.csv` |

**内置验证器快捷方式**:

| 快捷方式 | 实际验证器 | 适用类型 |
|---------|-----------|----------|
| `port` | Range(1, 65535) | int |
| `email` | Email() | string |
| `url` | URL() | string |
| `ip` | IP() | string |
| `filepath` | FileExists() | string |
| `dirpath` | DirExists() | string |

---

#### 3.2.3 add mutex-group - 添加互斥组

**功能**: 为一组标志添加互斥约束。

**用法示例**:

```bash
# 添加互斥组
$ qflag-cli add mutex-group deploy --name=strategy --flags=rolling,blue-green --optional

# 必须选择一个
$ qflag-cli add mutex-group build --name=format --flags=json,xml,yaml --required
```

---

#### 3.2.4 add required-group - 添加必需组

**功能**: 为一组标志添加必需约束（至少选择一个）。

**用法示例**:

```bash
$ qflag-cli add required-group deploy --name=target --flags=server,client --desc="至少选择一个部署目标"
```

---

### 3.3 remove - 删除组件

**功能**: 删除命令、标志或组，自动清理引用。

**用法示例**:

```bash
# 删除命令
$ qflag-cli remove cmd build

# 删除标志
$ qflag-cli remove flag build --name=output

# 删除互斥组
$ qflag-cli remove mutex-group deploy --name=strategy

# 强制删除（不提示确认）
$ qflag-cli remove cmd build --force
```

---

### 3.4 list - 列出项目结构

**功能**: 可视化展示当前项目的命令结构、标志和约束。

**用法示例**:

```bash
$ qflag-cli list

项目: myapp
模块: gitee.com/user/myapp
模式: layered
路径: /home/user/projects/myapp

命令结构:
  myapp
  ├── build (b)              # 构建项目
  │   ├── output (o)         # string  输出目录
  │   ├── target (t)         # enum    目标平台 [linux,windows,darwin]
  │   └── verbose (v)        # bool    详细输出
  ├── deploy (d)             # 部署应用到服务器
  │   ├── env (e)            # enum    部署环境 [dev,staging,prod]
  │   ├── strategy           # mutex   部署策略
  │   │   ├── rolling        # bool    滚动部署
  │   │   └── blue-green     # bool    蓝绿部署
  │   └── dry-run            # bool    模拟运行
  └── config                 # 配置管理
      ├── get (g)            # 获取配置值
      │   └── key            # string  配置键
      └── set (s)            # 设置配置值
          ├── key            # string  配置键
          └── value          # string  配置值

互斥组:
  - deploy.strategy: rolling, blue-green (可选)

统计:
  命令数: 6
  标志数: 12
  互斥组: 1
```

---

### 3.5 validate - 规范验证

**功能**: 检查项目是否符合 qflag 开发规范。

**验证项**:

1. **目录结构检查**
   - 是否符合 layered/simple 模式
   - 必要文件是否存在

2. **命名规范检查**
   - 命令变量名: `xxxCmd`
   - 标志变量名: `cmdName` 或 `xxxFlag`
   - run 函数名: `Run`
   - 包名规范

3. **代码结构检查**
   - 是否包含必要的注释（文件头、CmdOpts 注释）
   - 是否使用 `qflag.ExitOnError` 或 `qflag.ContinueOnError`
   - 标志变量是否为包级变量

4. **一致性检查**
   - SubCmds 中注册的命令是否存在
   - 标志引用的验证器是否存在
   - 互斥组中的标志是否存在

**用法示例**:

```bash
$ qflag-cli validate

✓ 目录结构检查通过
✓ 命名规范检查通过
✓ 代码结构检查通过
⚠ 发现 2 个警告:

  1. internal/cli/build.go:23
     标志 'output' 缺少短名称
     建议: 添加短名称以提高易用性
     修复: qflag-cli add flag build --name=output --short=o

  2. internal/cli/deploy.go:45
     建议为 'env' 标志添加环境变量绑定
     修复: 手动添加 .Env("APP_ENV") 调用

✓ 一致性检查通过

验证结果: 通过 (2个警告)
```

---

### 3.6 generate - 代码生成

#### 3.6.1 gen completion - 生成补全脚本

**功能**: 生成 shell 补全脚本。

**用法示例**:

```bash
# 生成 Bash 补全
$ qflag-cli gen completion bash > /etc/bash_completion.d/myapp

# 生成 Zsh 补全
$ qflag-cli gen completion zsh > /usr/local/share/zsh/site-functions/_myapp

# 生成 PowerShell 补全
$ qflag-cli gen completion powershell > myapp.ps1

# 生成 Fish 补全
$ qflag-cli gen completion fish > ~/.config/fish/completions/myapp.fish
```

---

#### 3.6.2 gen docs - 生成文档

**功能**: 生成命令文档。

**用法示例**:

```bash
# Markdown 格式
$ qflag-cli gen docs --format=markdown --output=docs/commands.md

# Man page 格式
$ qflag-cli gen docs --format=man --output=man/

# JSON 格式（用于自动化处理）
$ qflag-cli gen docs --format=json --output=docs/commands.json
```

---

#### 3.6.3 gen schema - 生成 JSON Schema

**功能**: 生成 JSON Schema，用于 IDE 自动补全和验证。

**用法示例**:

```bash
$ qflag-cli gen schema --output=.qflag-schema.json
```

---

### 3.7 config - 配置管理

**功能**: 管理项目级别的 qflag-cli 配置。

**用法示例**:

```bash
# 查看当前配置
$ qflag-cli config list

# 设置默认值
$ qflag-cli config set defaults.errorHandling ExitOnError
$ qflag-cli config set defaults.useChinese true

# 设置命名规范
$ qflag-cli config set naming.commandVar "{{.Name}}Cmd"
$ qflag-cli config set naming.flagVar "{{.Cmd}}{{.Name}}"
```

---

## 四、项目配置文件

### 4.1 配置文件位置

- 项目级: `.qflag-cli.yaml` (项目根目录)
- 用户级: `~/.config/qflag-cli/config.yaml`
- 系统级: `/etc/qflag-cli/config.yaml`

### 4.2 配置文件格式

```yaml
# .qflag-cli.yaml
version: "1.0"

project:
  name: myapp
  module: gitee.com/user/myapp
  mode: layered  # simple 或 layered
  created_at: "2026-04-12T10:00:00Z"

defaults:
  # 默认 CmdOpts
  errorHandling: ExitOnError  # ExitOnError, ContinueOnError, PanicOnError
  
  # 默认功能开关
  useChinese: true
  enableCompletion: true
  enableDynamicCompletion: false
  
  # 默认标志配置
  flagDefaults:
    required: false
    hidden: false

naming:
  # 命令变量命名模板
  commandVar: "{{.Name}}Cmd"           # 例如: buildCmd
  
  # 标志变量命名模板  
  flagVar: "{{.Cmd}}{{.Name}}"         # 例如: deployEnv
  
  # run 函数命名模板
  runFunc: "Run"                        # 固定为 Run
  
  # 文件名命名模板
  fileName: "{{.Name}}.go"              # 例如: deploy.go

templates:
  # 自定义模板路径 (相对于项目根目录)
  command: "templates/custom_cmd.go.tmpl"
  run: "templates/custom_run.go.tmpl"
  root: "templates/custom_root.go.tmpl"
  main: "templates/custom_main.go.tmpl"

hooks:
  # 生命周期钩子
  post-init: "scripts/post-init.sh"     # init 后执行
  post-add-cmd: "scripts/post-add-cmd.sh"  # add cmd 后执行

plugins:
  # 插件配置
  - name: "docker"
    enabled: false
  - name: "github-actions"
    enabled: true
```

---

## 五、实现架构

### 5.1 项目结构

```
qflag-cli/
├── cmd/
│   └── qflag-cli/
│       └── main.go                    # 入口
├── internal/
│   ├── cli/                           # 命令定义 (自身使用 qflag)
│   │   ├── root.go                    # 根命令
│   │   ├── init.go                    # init 子命令
│   │   ├── add.go                     # add 子命令
│   │   ├── remove.go                  # remove 子命令
│   │   ├── list.go                    # list 子命令
│   │   ├── validate.go                # validate 子命令
│   │   ├── generate.go                # generate 子命令
│   │   └── config.go                  # config 子命令
│   │
│   ├── core/                          # 核心业务逻辑
│   │   ├── project/                   # 项目管理
│   │   │   ├── project.go             # 项目结构定义
│   │   │   ├── scanner.go             # 项目扫描器
│   │   │   └── loader.go              # 项目加载器
│   │   │
│   │   ├── generator/                 # 代码生成
│   │   │   ├── generator.go           # 生成器接口
│   │   │   ├── templates/             # 代码模板
│   │   │   │   ├── root.go.tmpl
│   │   │   │   ├── cmd.go.tmpl
│   │   │   │   ├── run.go.tmpl
│   │   │   │   └── main.go.tmpl
│   │   │   ├── parser.go              # AST 解析
│   │   │   └── writer.go              # 文件写入
│   │   │
│   │   └── validator/                 # 规范验证
│   │       ├── validator.go           # 验证器接口
│   │       ├── naming.go              # 命名规范检查
│   │       ├── structure.go           # 结构检查
│   │       └── consistency.go         # 一致性检查
│   │
│   ├── interactive/                   # 交互式向导
│   │   ├── wizard.go                  # 向导框架
│   │   ├── cmd_wizard.go              # 命令向导
│   │   └── flag_wizard.go             # 标志向导
│   │
│   └── utils/                         # 工具函数
│       ├── file.go                    # 文件操作
│       ├── string.go                  # 字符串处理
│       └── template.go                # 模板处理
│
├── pkg/                               # 公共库
│   ├── scaffold/                      # 脚手架模板
│   │   ├── simple/                    # simple 模式模板
│   │   └── layered/                   # layered 模式模板
│   │
│   └── astutil/                       # AST 工具
│       └── parser.go
│
├── templates/                         # 默认模板
│   ├── simple/
│   │   ├── root.go.tmpl
│   │   ├── cmd.go.tmpl
│   │   └── main.go.tmpl
│   └── layered/
│       ├── root.go.tmpl
│       ├── cmd.go.tmpl
│       ├── run.go.tmpl
│       └── main.go.tmpl
│
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

### 5.2 关键数据结构

```go
// 项目配置
package project

type Project struct {
    Name        string
    Module      string
    Mode        ProjectMode          // simple / layered
    RootPath    string
    Config      *Config
    Commands    []Command
}

type ProjectMode string

const (
    ModeSimple  ProjectMode = "simple"
    ModeLayered ProjectMode = "layered"
)

type Command struct {
    Name        string
    ShortName   string
    Description string
    Parent      string              // 父命令名称
    Flags       []Flag
    MutexGroups []MutexGroup
    RequiredGroups []RequiredGroup
    FilePath    string              // 生成的文件路径
}

type Flag struct {
    Name        string
    ShortName   string
    Type        FlagType
    Default     interface{}
    Description string
    Required    bool
    Hidden      bool
    Validators  []string
}

type FlagType string

const (
    TypeString       FlagType = "string"
    TypeStrings      FlagType = "strings"
    TypeInt          FlagType = "int"
    TypeInts         FlagType = "ints"
    TypeInt64        FlagType = "int64"
    TypeFloat64      FlagType = "float64"
    TypeBool         FlagType = "bool"
    TypeCount        FlagType = "count"
    TypeEnum         FlagType = "enum"
    TypeEnums        FlagType = "enums"
    TypeDuration     FlagType = "duration"
    TypeIP           FlagType = "ip"
    TypeIPs          FlagType = "ips"
    TypeURL          FlagType = "url"
    TypeFilePath     FlagType = "filepath"
    TypeFilePathExists FlagType = "filepath-exists"
)

type MutexGroup struct {
    Name        string
    Flags       []string
    Optional    bool           // 是否允许都不选
}

type RequiredGroup struct {
    Name        string
    Flags       []string
    Description string
}
```

### 5.3 关键技术点

1. **AST 解析**
   - 使用 `go/parser` 和 `go/ast` 解析现有代码
   - 提取命令结构、标志定义、SubCmds 注册
   - 修改 AST 并生成新代码

2. **代码生成**
   - 使用 `text/template` 生成代码
   - 支持自定义模板覆盖
   - 保持代码格式（使用 `gofmt`）

3. **文件操作**
   - 修改前备份原文件
   - 错误时回滚
   - 并发安全

4. **交互界面**
   - 使用 `github.com/AlecAivazis/survey/v2` 实现交互式向导
   - 支持非交互模式（CI/CD）

5. **配置管理**
   - 使用 Viper 管理配置
   - 支持多层级配置覆盖

---

## 六、开发路线图

### Phase 1: MVP (1 周)

目标: 核心功能可用

- [x] 项目初始化 (`init`)
  - simple 和 layered 模式
  - 生成完整可运行项目
- [x] 添加命令 (`add cmd`)
  - 支持一级命令
  - 自动生成代码和注册
- [x] 添加标志 (`add flag`)
  - 支持基础类型
  - 自动更新命令文件

### Phase 2: 完整功能 (1-2 周)

目标: 支持复杂场景

- [ ] 多级子命令 (`add cmd parent/child`)
- [ ] 互斥组和必需组
- [ ] 删除功能 (`remove`)
- [ ] 项目列表 (`list`)
- [ ] 交互式向导

### Phase 3: 质量保障 (1 周)

目标: 确保项目质量

- [ ] 规范验证 (`validate`)
- [ ] 代码生成 (`gen completion`, `gen docs`)
- [ ] 测试覆盖率 > 80%
- [ ] 文档完善

### Phase 4: 生态建设 (持续)

- [ ] IDE 插件 (VS Code)
- [ ] 更多项目模板
- [ ] 社区贡献模板
- [ ] CI/CD 集成

---

## 七、与 cobra-cli 对比

| 特性 | cobra-cli | qflag-cli |
|------|-----------|-----------|
| **初始化** | 基础结构 | 完整可运行项目 + 示例 |
| **项目模式** | 单一 | simple / layered 可选 |
| **交互式** | 弱 | 强（向导式） |
| **标志类型** | 基础 | 17种类型 |
| **验证器** | ❌ | ✅ 内置验证器 |
| **互斥组** | ❌ | ✅ |
| **必需组** | ❌ | ✅ |
| **动态补全** | ❌ | ✅ |
| **智能纠错** | ❌ | ✅ |
| **规范检查** | ❌ | ✅ |
| **文档生成** | 基础 | 多格式支持 |
| **并发安全** | ❌ | ✅ |

---

## 八、使用示例

### 快速开始

```bash
# 1. 安装工具
go install gitee.com/MM-Q/qflag/qflag-cli@latest

# 2. 创建项目
qflag-cli init myapp
cd myapp

# 3. 添加命令
qflag-cli add cmd server --desc="启动服务器"
qflag-cli add flag server --name=port --short=p --type=int --default=8080 --validator=port

# 4. 运行
qflag run ./cmd/myapp server --port=9000

# 5. 验证规范
qflag-cli validate
```

### 完整示例

```bash
# 创建项目
qflag-cli init todo-cli --mode=layered

# 添加命令
qflag-cli add cmd add --desc="添加待办事项"
qflag-cli add flag add --name=priority --short=p --type=enum --values=low,medium,high --default=medium
qflag-cli add flag add --name=due --type=duration --desc="截止日期"

qflag-cli add cmd list --desc="列出待办事项"
qflag-cli add flag list --name=status --short=s --type=enum --values=all,pending,done --default=all
qflag-cli add flag list --name=limit --short=n --type=int --default=10

qflag-cli add cmd done --desc="标记为完成"
qflag-cli add mutex-group done --name=target --flags=id,all

# 查看结构
qflag-cli list

# 生成文档
qflag-cli gen docs --format=markdown --output=README.md

# 验证
qflag-cli validate
```

---

## 九、总结

`qflag-cli` 将成为 `qflag` 生态系统的重要组成部分，通过提供：

1. **标准化** — 强制遵循最佳实践
2. **自动化** — 减少重复劳动
3. **智能化** — 向导式交互
4. **可验证** — 确保项目健康

这将大大降低 `qflag` 的使用门槛，提升开发效率，形成与 `cobra` 的差异化竞争优势。

---

*文档版本: v1.0*
*最后更新: 2026-04-12*
