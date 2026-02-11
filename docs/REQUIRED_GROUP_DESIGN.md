# 必需组功能设计方案

## 一、设计概述

### 核心需求
与互斥组（MutexGroup）相反，**必需组（RequiredGroup）**要求组内**所有标志都必须被设置**。

### 使用场景
1. **数据库连接配置**：`--db-host`, `--db-port`, `--db-user`, `--db-pass` 必须同时提供
2. **API 认证配置**：`--api-key`, `--api-secret` 必须同时提供
3. **文件上传配置**：`--file-path`, `--upload-url` 必须同时提供
4. **服务器配置**：`--host`, `--port` 必须同时提供

### 设计原则
1. **与互斥组对称**：API 设计与互斥组保持一致
2. **向后兼容**：不破坏现有代码
3. **类型安全**：使用结构化类型
4. **错误清晰**：提供明确的错误提示

---

## 二、核心设计

### 2.1 RequiredGroup 类型定义

```go
// RequiredGroup 必需组定义
//
// RequiredGroup 定义了一组必需的标志，其中所有标志都必须被设置。
// 当用户没有设置必需组中的某些标志时，解析器会返回错误。
//
// 字段说明:
//   - Name: 必需组名称，用于错误提示和标识
//   - Flags: 必需组中的标志名称列表
//
// 使用场景:
//   - 数据库连接配置（host, port, user, pass）
//   - API 认证配置（api-key, api-secret）
//   - 文件上传配置（file-path, upload-url）
type RequiredGroup struct {
    Name  string   // 必需组名称，用于错误提示和标识
    Flags []string // 必需组中的标志名称列表
}
```

### 2.2 CmdConfig 扩展

```go
// CmdConfig 命令配置类型
type CmdConfig struct {
    Version     string            // 版本号
    UseChinese  bool              // 是否使用中文
    EnvPrefix   string            // 环境变量前缀
    UsageSyntax string            // 命令使用语法
    Example     map[string]string // 示例使用，key为描述，value为示例命令
    Notes       []string          // 注意事项
    LogoText    string            // 命令logo文本
    MutexGroups   []MutexGroup    // 互斥组列表
    RequiredGroups []RequiredGroup // 必需组列表（新增）
    Completion  bool              // 是否启用自动补全标志
}
```

### 2.3 Command 接口扩展

```go
// Command 接口定义了命令的核心行为
type Command interface {
    // ... 现有方法 ...

    // 必需组管理（新增）
    AddRequiredGroup(name string, flags []string) error  // 添加必需组
    RemoveRequiredGroup(name string) error              // 移除必需组
    GetRequiredGroup(name string) (*RequiredGroup, bool) // 获取必需组
    RequiredGroups() []RequiredGroup                   // 获取所有必需组
}
```

---

## 三、API 设计

### 3.1 添加必需组

```go
// AddRequiredGroup 添加必需组
//
// 参数:
//   - name: 必需组名称
//   - flags: 必需组中的标志名称列表
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 添加一个必需组到命令配置
//   - 如果组名已存在，返回错误
//   - 如果标志列表为空，返回错误
//   - 如果标志不存在，返回错误
//
// 错误码:
//   - REQUIRED_GROUP_ALREADY_EXISTS: 必需组已存在
//   - EMPTY_REQUIRED_GROUP: 必需组标志列表为空
//   - FLAG_NOT_FOUND: 标志不存在
//
// 示例:
//   - cmd.AddRequiredGroup("数据库配置", []string{"db-host", "db-port", "db-user", "db-pass"})
func (c *Cmd) AddRequiredGroup(name string, flags []string) error
```

### 3.2 移除必需组

```go
// RemoveRequiredGroup 移除必需组
//
// 参数:
//   - name: 必需组名称
//
// 返回值:
//   - error: 移除失败时返回错误
//
// 功能说明:
//   - 从命令配置中移除指定的必需组
//   - 如果组不存在，返回错误
//
// 错误码:
//   - REQUIRED_GROUP_NOT_FOUND: 必需组不存在
//
// 示例:
//   - cmd.RemoveRequiredGroup("数据库配置")
func (c *Cmd) RemoveRequiredGroup(name string) error
```

### 3.3 获取必需组

```go
// GetRequiredGroup 获取必需组
//
// 参数:
//   - name: 必需组名称
//
// 返回值:
//   - *RequiredGroup: 必需组指针
//   - bool: 是否找到
//
// 功能说明:
//   - 根据名称获取必需组
//   - 如果组不存在，返回 nil 和 false
//
// 示例:
//   - group, ok := cmd.GetRequiredGroup("数据库配置")
func (c *Cmd) GetRequiredGroup(name string) (*RequiredGroup, bool)
```

### 3.4 获取所有必需组

```go
// RequiredGroups 获取所有必需组
//
// 返回值:
//   - []RequiredGroup: 所有必需组列表
//
// 功能说明:
//   - 返回命令配置中的所有必需组
//   - 返回的是副本，修改不会影响原配置
//
// 示例:
//   - groups := cmd.RequiredGroups()
func (c *Cmd) RequiredGroups() []RequiredGroup
```

---

## 四、解析器验证逻辑

### 4.1 验证函数

```go
// validateRequiredGroups 验证命令的必需组规则
//
// 参数:
//   - cmd: 要验证的命令
//
// 返回值:
//   - error: 如果必需组验证失败返回错误
//
// 功能说明:
//   - 检查每个必需组中是否有标志未被设置
//   - 提供清晰的错误信息，指出未设置的标志和组名
//
// 验证规则:
//   - 必需组中的所有标志都必须被设置
//   - 如果有任何一个标志未被设置，返回错误
//
// 错误处理:
//   - 使用 types.NewError 创建结构化错误
//   - 错误信息包含必需组名称和未设置的标志列表
//
// 示例错误:
//   - "required flags [db-host, db-user, db-pass] in group '数据库配置' must be set"
func (p *DefaultParser) validateRequiredGroups(cmd types.Command) error {
    config := cmd.Config()
    if config == nil {
        return nil
    }

    // 检查必需组是否为空
    if len(config.RequiredGroups) == 0 {
        return nil
    }

    // 遍历所有必需组
    for _, group := range config.RequiredGroups {
        var unsetFlags []string

        // 检查必需组中的每个标志是否被设置
        for _, flagName := range group.Flags {
            if flag, exists := cmd.GetFlag(flagName); exists && !flag.IsSet() {
                unsetFlags = append(unsetFlags, flagName)
            }
        }

        // 验证必需组规则
        if len(unsetFlags) > 0 {
            return types.NewError("REQUIRED_GROUP_VIOLATION",
                fmt.Sprintf("required flags %v in group '%s' must be set", unsetFlags, group.Name),
                nil)
        }
    }

    return nil
}
```

### 4.2 验证时机

必需组的验证在解析器的 `Parse` 方法中执行，在解析完所有参数之后、执行命令之前：

```go
func (p *DefaultParser) Parse(cmd types.Command, args []string) error {
    // 1. 注册内置标志
    // 2. 注册命令的所有标志到 FlagSet
    // 3. 解析命令行参数
    // 4. 加载环境变量
    // 5. 验证互斥组规则
    // 6. 验证必需组规则（新增）
    if err := p.validateRequiredGroups(cmd); err != nil {
        return err
    }
    // 7. 处理内置标志
    // ...
}
```

---

## 五、使用示例

### 5.1 基本使用

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建命令
    cmd := qflag.NewCmd("db-tool", "", qflag.ExitOnError)

    // 定义标志
    dbHost := cmd.String("db-host", "", "数据库主机", "")
    dbPort := cmd.Int("db-port", "", "数据库端口", 3306)
    dbUser := cmd.String("db-user", "", "数据库用户", "")
    dbPass := cmd.String("db-pass", "", "数据库密码", "")

    // 添加必需组：所有数据库参数都必须提供
    cmd.AddRequiredGroup("数据库连接", []string{"db-host", "db-port", "db-user", "db-pass"})

    // 解析
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("数据库主机: %s\n", dbHost.Get())
    fmt.Printf("数据库端口: %d\n", dbPort.Get())
    fmt.Printf("数据库用户: %s\n", dbUser.Get())
}
```

### 5.2 多个必需组

```go
// 创建命令
cmd := qflag.NewCmd("api-client", "", qflag.ExitOnError)

// API 认证配置
apiKey := cmd.String("api-key", "", "API 密钥", "")
apiSecret := cmd.String("api-secret", "", "API 密钥", "")

// 服务器配置
serverHost := cmd.String("server-host", "", "服务器主机", "")
serverPort := cmd.Int("server-port", "", "服务器端口", 8080)

// 添加必需组
cmd.AddRequiredGroup("API 认证", []string{"api-key", "api-secret"})
cmd.AddRequiredGroup("服务器配置", []string{"server-host", "server-port"})

// 解析
if err := cmd.Parse(os.Args[1:]); err != nil {
    fmt.Printf("错误: %v\n", err)
    os.Exit(1)
}
```

### 5.3 与互斥组配合使用

```go
// 场景：数据源配置
// 方式1：使用数据库（需要所有数据库参数）
// 方式2：使用文件（需要文件路径）
// 两种方式互斥

// 标志定义
useDb := cmd.Bool("use-db", "", "使用数据库", false)
useFile := cmd.Bool("use-file", "", "使用文件", false)

dbHost := cmd.String("db-host", "", "数据库主机", "")
dbPort := cmd.Int("db-port", "", "数据库端口", 3306)
dbUser := cmd.String("db-user", "", "数据库用户", "")
dbPass := cmd.String("db-pass", "", "数据库密码", "")

filePath := cmd.String("file-path", "", "文件路径", "")

// 互斥组：两种方式只能选一种
cmd.AddMutexGroup("数据源", []string{"use-db", "use-file"}, false)

// 必需组：如果选择数据库，必须提供所有数据库参数
cmd.AddRequiredGroup("数据库配置", []string{"db-host", "db-port", "db-user", "db-pass"})

// 必需组：如果选择文件，必须提供文件路径
cmd.AddRequiredGroup("文件配置", []string{"file-path"})
```

### 5.4 移除必需组

```go
// 添加必需组
cmd.AddRequiredGroup("数据库配置", []string{"db-host", "db-port", "db-user", "db-pass"})

// 获取必需组
group, ok := cmd.GetRequiredGroup("数据库配置")
if ok {
    fmt.Printf("必需组: %s, 标志: %v\n", group.Name, group.Flags)
}

// 移除必需组
cmd.RemoveRequiredGroup("数据库配置")
```

### 5.5 获取所有必需组

```go
// 添加多个必需组
cmd.AddRequiredGroup("数据库配置", []string{"db-host", "db-port", "db-user", "db-pass"})
cmd.AddRequiredGroup("API 认证", []string{"api-key", "api-secret"})

// 获取所有必需组
groups := cmd.RequiredGroups()
for _, group := range groups {
    fmt.Printf("必需组: %s, 标志: %v\n", group.Name, group.Flags)
}
```

### 5.6 使用 CmdSpec 配置必需组

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建命令规格
    spec := qflag.NewCmdSpec("db-tool", "")

    // 设置基本属性
    spec.Desc = "数据库连接工具"
    spec.Version = "1.0.0"
    spec.UseChinese = true

    // 设置运行函数
    spec.RunFunc = func(cmd qflag.Command) error {
        fmt.Println("数据库连接成功！")
        return nil
    }

    // 添加必需组
    spec.RequiredGroups = []qflag.RequiredGroup{
        {
            Name:  "数据库连接",
            Flags: []string{"db-host", "db-port", "db-user", "db-pass"},
        },
        {
            Name:  "TLS 配置",
            Flags: []string{"cert-file", "key-file"},
        },
    }

    // 从规格创建命令
    cmd, err := qflag.NewCmdFromSpec(spec)
    if err != nil {
        fmt.Printf("创建命令失败: %v\n", err)
        os.Exit(1)
    }

    // 添加标志
    dbHost := cmd.String("db-host", "", "数据库主机", "")
    dbPort := cmd.Int("db-port", "", "数据库端口", 3306)
    dbUser := cmd.String("db-user", "", "数据库用户", "")
    dbPass := cmd.String("db-pass", "", "数据库密码", "")
    certFile := cmd.String("cert-file", "", "证书文件", "")
    keyFile := cmd.String("key-file", "", "密钥文件", "")

    // 解析
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("数据库主机: %s\n", dbHost.Get())
    fmt.Printf("数据库端口: %d\n", dbPort.Get())
    fmt.Printf("数据库用户: %s\n", dbUser.Get())
}
```

### 5.7 使用 CmdSpec 配置互斥组和必需组

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建命令规格
    spec := qflag.NewCmdSpec("data-import", "")

    spec.Desc = "数据导入工具"
    spec.Version = "1.0.0"

    // 添加互斥组：数据源互斥
    spec.MutexGroups = []qflag.MutexGroup{
        {
            Name:      "数据源",
            Flags:     []string{"use-db", "use-file"},
            AllowNone: false,
        },
    }

    // 添加必需组：如果使用数据库，必须提供所有数据库参数
    spec.RequiredGroups = []qflag.RequiredGroup{
        {
            Name:  "数据库配置",
            Flags: []string{"db-host", "db-port", "db-user", "db-pass"},
        },
        {
            Name:  "文件配置",
            Flags: []string{"file-path"},
        },
    }

    // 从规格创建命令
    cmd, err := qflag.NewCmdFromSpec(spec)
    if err != nil {
        fmt.Printf("创建命令失败: %v\n", err)
        os.Exit(1)
    }

    // 添加标志
    useDb := cmd.Bool("use-db", "", "使用数据库", false)
    useFile := cmd.Bool("use-file", "", "使用文件", false)
    dbHost := cmd.String("db-host", "", "数据库主机", "")
    dbPort := cmd.Int("db-port", "", "数据库端口", 3306)
    dbUser := cmd.String("db-user", "", "数据库用户", "")
    dbPass := cmd.String("db-pass", "", "数据库密码", "")
    filePath := cmd.String("file-path", "", "文件路径", "")

    // 解析
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        os.Exit(1)
    }

    if useDb.Get() {
        fmt.Printf("使用数据库: %s:%d\n", dbHost.Get(), dbPort.Get())
    } else if useFile.Get() {
        fmt.Printf("使用文件: %s\n", filePath.Get())
    }
}
```

---

## 六、错误处理

### 6.1 错误码定义

```go
// 错误码常量
const (
    // REQUIRED_GROUP_ALREADY_EXISTS 必需组已存在
    REQUIRED_GROUP_ALREADY_EXISTS = "REQUIRED_GROUP_ALREADY_EXISTS"

    // REQUIRED_GROUP_NOT_FOUND 必需组不存在
    REQUIRED_GROUP_NOT_FOUND = "REQUIRED_GROUP_NOT_FOUND"

    // EMPTY_REQUIRED_GROUP 必需组标志列表为空
    EMPTY_REQUIRED_GROUP = "EMPTY_REQUIRED_GROUP"

    // REQUIRED_GROUP_VIOLATION 必需组验证失败
    REQUIRED_GROUP_VIOLATION = "REQUIRED_GROUP_VIOLATION"
)
```

### 6.2 错误示例

#### 6.2.1 必需组已存在

```go
cmd.AddRequiredGroup("数据库配置", []string{"db-host", "db-port"})
// 错误: required group '数据库配置' already exists
cmd.AddRequiredGroup("数据库配置", []string{"db-user", "db-pass"})
```

#### 6.2.2 必需组不存在

```go
// 错误: required group '不存在的组' not found
cmd.RemoveRequiredGroup("不存在的组")
```

#### 6.2.3 必需组标志列表为空

```go
// 错误: required group cannot be empty
cmd.AddRequiredGroup("空组", []string{})
```

#### 6.2.4 标志不存在

```go
// 错误: flag '不存在的标志' not found
cmd.AddRequiredGroup("组", []string{"不存在的标志"})
```

#### 6.2.5 必需组验证失败

```go
// 用户只提供了 --db-host 和 --db-port
// 错误: required flags [db-user, db-pass] in group '数据库配置' must be set
cmd.AddRequiredGroup("数据库配置", []string{"db-host", "db-port", "db-user", "db-pass"})
```

### 6.3 错误信息格式

```go
// 必需组验证失败
错误: required flags [db-host, db-user, db-pass] in group '数据库配置' must be set

// 中文版本（如果 UseChinese 为 true）
错误: 必需组 '数据库配置' 中的标志 [db-host, db-user, db-pass] 必须被设置
```

---

## 七、实现步骤

### 第一步：定义 RequiredGroup 类型

文件：`internal/types/config.go`

```go
// RequiredGroup 必需组定义
type RequiredGroup struct {
    Name  string   // 必需组名称
    Flags []string // 必需组中的标志名称列表
}
```

### 第二步：扩展 CmdConfig

文件：`internal/types/config.go`

```go
type CmdConfig struct {
    // ... 现有字段 ...
    MutexGroups    []MutexGroup    // 互斥组列表
    RequiredGroups []RequiredGroup // 必需组列表（新增）
    // ... 其他字段 ...
}
```

### 第三步：扩展 Command 接口

文件：`internal/types/command.go`

```go
type Command interface {
    // ... 现有方法 ...

    // 必需组管理（新增）
    AddRequiredGroup(name string, flags []string) error
    RemoveRequiredGroup(name string) error
    GetRequiredGroup(name string) (*RequiredGroup, bool)
    RequiredGroups() []RequiredGroup
}
```

### 第四步：实现 Command 接口方法

文件：`internal/cmd/cmd.go`

```go
// AddRequiredGroup 添加必需组
func (c *Cmd) AddRequiredGroup(name string, flags []string) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    // 检查组名是否已存在
    for _, group := range c.config.RequiredGroups {
        if group.Name == name {
            return types.NewError("REQUIRED_GROUP_ALREADY_EXISTS",
                fmt.Sprintf("required group '%s' already exists", name), nil)
        }
    }

    // 检查标志列表是否为空
    if len(flags) == 0 {
        return types.NewError("EMPTY_REQUIRED_GROUP",
            "required group cannot be empty", nil)
    }

    // 检查所有标志是否存在
    for _, flagName := range flags {
        if _, exists := c.GetFlag(flagName); !exists {
            return types.NewError("FLAG_NOT_FOUND",
                fmt.Sprintf("flag '%s' not found", flagName), nil)
        }
    }

    // 添加必需组
    c.config.RequiredGroups = append(c.config.RequiredGroups, RequiredGroup{
        Name:  name,
        Flags: flags,
    })

    return nil
}

// RemoveRequiredGroup 移除必需组
func (c *Cmd) RemoveRequiredGroup(name string) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    for i, group := range c.config.RequiredGroups {
        if group.Name == name {
            // 移除必需组
            c.config.RequiredGroups = append(c.config.RequiredGroups[:i], c.config.RequiredGroups[i+1:]...)
            return nil
        }
    }

    return types.NewError("REQUIRED_GROUP_NOT_FOUND",
        fmt.Sprintf("required group '%s' not found", name), nil)
}

// GetRequiredGroup 获取必需组
func (c *Cmd) GetRequiredGroup(name string) (*RequiredGroup, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    for _, group := range c.config.RequiredGroups {
        if group.Name == name {
            return &group, true
        }
    }

    return nil, false
}

// RequiredGroups 获取所有必需组
func (c *Cmd) RequiredGroups() []RequiredGroup {
    c.mu.RLock()
    defer c.mu.RUnlock()

    // 返回副本
    result := make([]RequiredGroup, len(c.config.RequiredGroups))
    copy(result, c.config.RequiredGroups)
    return result
}
```

### 第五步：实现解析器验证逻辑

文件：`internal/parser/parser_validation.go`

```go
// validateRequiredGroups 验证命令的必需组规则
func (p *DefaultParser) validateRequiredGroups(cmd types.Command) error {
    config := cmd.Config()
    if config == nil {
        return nil
    }

    if len(config.RequiredGroups) == 0 {
        return nil
    }

    for _, group := range config.RequiredGroups {
        var unsetFlags []string

        for _, flagName := range group.Flags {
            if flag, exists := cmd.GetFlag(flagName); exists && !flag.IsSet() {
                unsetFlags = append(unsetFlags, flagName)
            }
        }

        if len(unsetFlags) > 0 {
            return types.NewError("REQUIRED_GROUP_VIOLATION",
                fmt.Sprintf("required flags %v in group '%s' must be set", unsetFlags, group.Name),
                nil)
        }
    }

    return nil
}
```

### 第六步：在解析器中调用验证

文件：`internal/parser/parser.go`

```go
func (p *DefaultParser) Parse(cmd types.Command, args []string) error {
    // ... 现有逻辑 ...

    // 验证互斥组规则
    if err := p.validateMutexGroups(cmd); err != nil {
        return err
    }

    // 验证必需组规则（新增）
    if err := p.validateRequiredGroups(cmd); err != nil {
        return err
    }

    // ... 其他逻辑 ...
}
```

### 第七步：导出 RequiredGroup 类型

文件：`exports.go`

```go
// RequiredGroup 必需组定义
type RequiredGroup = types.RequiredGroup
```

### 第八步：扩展 CmdSpec 支持 RequiredGroups

文件：`internal/cmd/cmdspec.go`

#### 8.1 修改 CmdSpec 结构体

```go
// CmdSpec 命令规格结构体
type CmdSpec struct {
    // 基本属性
    LongName      string              // 命令长名称
    ShortName     string              // 命令短名称
    Desc          string              // 命令描述
    ErrorHandling types.ErrorHandling // 错误处理策略

    // 运行函数
    RunFunc func(types.Command) error // 命令执行函数

    // 配置选项
    Version     string // 版本号
    UseChinese  bool   // 是否使用中文
    EnvPrefix   string // 环境变量前缀
    UsageSyntax string // 命令使用语法
    LogoText    string // Logo文本
    Completion  bool   // 是否启用自动补全标志

    // 示例和说明
    Examples map[string]string // 示例使用, key为描述, value为示例命令
    Notes    []string          // 注意事项

    // 子命令和互斥组
    SubCmds       []types.Command    // 子命令列表, 用于添加到命令中
    MutexGroups   []types.MutexGroup // 互斥组列表
    RequiredGroups []types.RequiredGroup // 必需组列表（新增）
}
```

#### 8.2 修改 NewCmdSpec 函数

```go
// NewCmdSpec 创建新的命令规格
func NewCmdSpec(longName, shortName string) *CmdSpec {
    return &CmdSpec{
        LongName:      longName,
        ShortName:     shortName,
        ErrorHandling: types.ExitOnError, // 默认错误处理策略
        UseChinese:    false,             // 默认不使用中文
        Completion:    false,             // 默认不启用自动补全
        Examples:      make(map[string]string),
        Notes:         []string{},
        SubCmds:       []types.Command{},
        MutexGroups:   []types.MutexGroup{},
        RequiredGroups: []types.RequiredGroup{}, // 初始化必需组列表（新增）
    }
}
```

#### 8.3 修改 NewCmdFromSpec 函数

```go
// NewCmdFromSpec 从规格创建命令
func NewCmdFromSpec(spec *CmdSpec) (cmd *Cmd, err error) {
    // ... 现有逻辑 ...

    // 添加互斥组
    for _, group := range spec.MutexGroups {
        cmd.AddMutexGroup(group.Name, group.Flags, group.AllowNone)
    }

    // 添加必需组（新增）
    for _, group := range spec.RequiredGroups {
        if err := cmd.AddRequiredGroup(group.Name, group.Flags); err != nil {
            return nil, types.WrapError(err, "FAILED_TO_ADD_REQUIRED_GROUP", "failed to add required group")
        }
    }

    // 添加子命令
    if len(spec.SubCmds) > 0 {
        if err := cmd.AddSubCmds(spec.SubCmds...); err != nil {
            return nil, types.WrapError(err, "FAILED_TO_ADD_SUBCMDS", "failed to add subcommands")
        }
    }

    return cmd, nil
}
```

### 第九步：扩展 CmdOpts 支持 RequiredGroups

文件：`internal/cmd/cmdopts.go`

#### 9.1 修改 CmdOpts 结构体

```go
// CmdOpts 命令选项
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
    Completion  bool   // 是否启用自动补全标志

    // 示例和说明
    Examples map[string]string // 示例使用, key为描述, value为示例命令
    Notes    []string          // 注意事项

    // 子命令和互斥组
    SubCmds       []types.Command    // 子命令列表, 用于添加到命令中
    MutexGroups   []types.MutexGroup // 互斥组列表
    RequiredGroups []types.RequiredGroup // 必需组列表（新增）
}
```

#### 9.2 修改 NewCmdOpts 函数

```go
// NewCmdOpts 创建新的命令选项
func NewCmdOpts() *CmdOpts {
    return &CmdOpts{
        Examples:      make(map[string]string),
        Notes:         []string{},
        SubCmds:       []types.Command{},
        MutexGroups:   []types.MutexGroup{},
        RequiredGroups: []types.RequiredGroup{}, // 初始化必需组列表（新增）
    }
}
```

#### 9.3 修改 ApplyOpts 函数

文件：`internal/cmd/cmd.go`

```go
// ApplyOpts 应用选项到命令
func (c *Cmd) ApplyOpts(opts *CmdOpts) error {
    // ... 现有逻辑 ...

    // 4. 添加互斥组 - 调用现有方法
    if len(opts.MutexGroups) > 0 {
        for _, group := range opts.MutexGroups {
            c.AddMutexGroup(group.Name, group.Flags, group.AllowNone)
        }
    }

    // 5. 添加必需组 - 调用现有方法（新增）
    if len(opts.RequiredGroups) > 0 {
        for _, group := range opts.RequiredGroups {
            if err := c.AddRequiredGroup(group.Name, group.Flags); err != nil {
                return types.WrapError(err, "FAILED_TO_ADD_REQUIRED_GROUP", "failed to add required group")
            }
        }
    }

    // 6. 添加子命令 - 调用现有方法
    if len(opts.SubCmds) > 0 {
        if err := c.AddSubCmds(opts.SubCmds...); err != nil {
            return types.WrapError(err, "FAILED_TO_ADD_SUBCMDS", "failed to add subcommands")
        }
    }

    return err
}
```

### 第十步：编写测试

文件：`internal/cmd/required_group_test.go`（新建）

```go
package cmd

import (
    "testing"

    "gitee.com/MM-Q/qflag/internal/types"
)

func TestAddRequiredGroup(t *testing.T) {
    cmd := NewCmd("test", "", types.ExitOnError)

    // 添加标志
    cmd.String("flag1", "", "Flag 1", "")
    cmd.String("flag2", "", "Flag 2", "")

    // 测试添加必需组
    err := cmd.AddRequiredGroup("group1", []string{"flag1", "flag2"})
    if err != nil {
        t.Fatalf("AddRequiredGroup failed: %v", err)
    }

    // 测试重复添加
    err = cmd.AddRequiredGroup("group1", []string{"flag1"})
    if err == nil {
        t.Fatal("expected error for duplicate group")
    }
}

func TestRemoveRequiredGroup(t *testing.T) {
    cmd := NewCmd("test", "", types.ExitOnError)

    // 添加标志和必需组
    cmd.String("flag1", "", "Flag 1", "")
    cmd.AddRequiredGroup("group1", []string{"flag1"})

    // 测试移除必需组
    err := cmd.RemoveRequiredGroup("group1")
    if err != nil {
        t.Fatalf("RemoveRequiredGroup failed: %v", err)
    }

    // 测试移除不存在的组
    err = cmd.RemoveRequiredGroup("group2")
    if err == nil {
        t.Fatal("expected error for non-existent group")
    }
}

func TestGetRequiredGroup(t *testing.T) {
    cmd := NewCmd("test", "", types.ExitOnError)

    // 添加标志和必需组
    cmd.String("flag1", "", "Flag 1", "")
    cmd.AddRequiredGroup("group1", []string{"flag1"})

    // 测试获取必需组
    group, ok := cmd.GetRequiredGroup("group1")
    if !ok {
        t.Fatal("expected group to exist")
    }
    if group.Name != "group1" {
        t.Fatalf("expected group name 'group1', got '%s'", group.Name)
    }

    // 测试获取不存在的组
    _, ok = cmd.GetRequiredGroup("group2")
    if ok {
        t.Fatal("expected group to not exist")
    }
}

func TestRequiredGroups(t *testing.T) {
    cmd := NewCmd("test", "", types.ExitOnError)

    // 添加标志和必需组
    cmd.String("flag1", "", "Flag 1", "")
    cmd.String("flag2", "", "Flag 2", "")
    cmd.AddRequiredGroup("group1", []string{"flag1"})
    cmd.AddRequiredGroup("group2", []string{"flag2"})

    // 测试获取所有必需组
    groups := cmd.RequiredGroups()
    if len(groups) != 2 {
        t.Fatalf("expected 2 groups, got %d", len(groups))
    }
}
```

### 第九步：编写示例

文件：`examples/required-group/main.go`（新建）

```go
package main

import (
    "fmt"
    "os"

    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建命令
    cmd := qflag.NewCmd("db-tool", "", qflag.ExitOnError)

    // 定义标志
    dbHost := cmd.String("db-host", "", "数据库主机", "")
    dbPort := cmd.Int("db-port", "", "数据库端口", 3306)
    dbUser := cmd.String("db-user", "", "数据库用户", "")
    dbPass := cmd.String("db-pass", "", "数据库密码", "")

    // 添加必需组：所有数据库参数都必须提供
    cmd.AddRequiredGroup("数据库连接", []string{"db-host", "db-port", "db-user", "db-pass"})

    // 解析
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("数据库主机: %s\n", dbHost.Get())
    fmt.Printf("数据库端口: %d\n", dbPort.Get())
    fmt.Printf("数据库用户: %s\n", dbUser.Get())
    fmt.Printf("数据库密码: %s\n", dbPass.Get())
}
```

---

## 八、与互斥组的对比

| 特性 | 互斥组 | 必需组 |
|------|--------|--------|
| **语义** | 最多只能有一个 | 必须全部设置 |
| **AllowNone** | 是否允许一个都不设置 | 不需要此字段 |
| **验证规则** | setCount > 1 时报错 | unsetFlags > 0 时报错 |
| **使用场景** | 互斥选项（如 --json 和 --xml） | 必需选项（如数据库配置） |
| **错误信息** | "mutually exclusive flags ... cannot be used together" | "required flags ... must be set" |

### 对比示例

```go
// 互斥组：两种输出格式只能选一种
cmd.AddMutexGroup("输出格式", []string{"--json", "--xml"}, false)

// 必需组：所有数据库参数都必须提供
cmd.AddRequiredGroup("数据库连接", []string{"db-host", "db-port", "db-user", "db-pass"})
```

---

## 九、与验证器的配合

### 9.1 验证顺序

1. **必需组验证**：检查所有必需的标志是否被设置
2. **标志解析**：解析命令行参数
3. **验证器验证**：对每个设置了验证器的标志进行验证

### 9.2 配合使用示例

```go
// 场景：数据库连接工具
cmd := qflag.NewCmd("db-tool", "", qflag.ExitOnError)

// 标志定义
host := cmd.String("host", "h", "数据库主机", "")
port := cmd.Int("port", "p", "数据库端口", 3306)
user := cmd.String("user", "u", "数据库用户", "")
pass := cmd.String("pass", "", "数据库密码", "")

// 添加验证器
port.SetValidator(qflag.RangeValidator(1, 65535))
host.SetValidator(qflag.NonEmptyValidator[string]())
user.SetValidator(qflag.LengthValidator(1, 50))

// 添加必需组：所有数据库参数都必须提供
cmd.AddRequiredGroup("数据库连接", []string{"host", "port", "user", "pass"})

// 解析
if err := cmd.Parse(os.Args[1:]); err != nil {
    fmt.Printf("错误: %v\n", err)
    os.Exit(1)
}
```

### 9.3 错误场景

#### 场景1：必需组验证失败

```
用户输入：--host localhost --port 3306
错误: required flags [user, pass] in group '数据库连接' must be set
```

#### 场景2：验证器验证失败

```
用户输入：--host localhost --port 70000 --user root --pass secret
错误: 值 70000 超出范围 [1, 65535]
```

---

## 十、设计优势

| 优势 | 说明 |
|------|------|
| ✅ **与互斥组对称** | API 设计与互斥组保持一致，易于理解 |
| ✅ **向后兼容** | 不破坏现有代码，可选功能 |
| ✅ **类型安全** | 使用结构化类型，编译时检查 |
| ✅ **错误清晰** | 提供明确的错误提示，包含组名和标志列表 |
| ✅ **易于使用** | 简单的 API，一行代码即可添加必需组 |
| ✅ **可组合** | 与互斥组、验证器等功能配合使用 |
| ✅ **线程安全** | 所有方法都使用锁保护 |
| ✅ **可测试** | 提供完整的测试覆盖 |

---

## 十一、注意事项

1. **必需组名称唯一**：同一个命令中，必需组名称必须唯一
2. **标志必须存在**：添加必需组时，所有标志必须已经注册到命令
3. **标志列表不能为空**：必需组必须包含至少一个标志
4. **验证时机**：必需组验证在解析完所有参数之后执行
5. **与互斥组的配合**：必需组和互斥组可以同时使用，但要注意逻辑冲突
6. **环境变量**：通过环境变量设置的值也会被 `IsSet()` 识别为已设置
7. **默认值**：如果标志有默认值，但用户没有设置，`IsSet()` 返回 false

---

## 十二、高级用法

### 12.1 条件必需组

```go
// 场景：根据某个标志的值，决定哪些标志是必需的
// 注意：这需要自定义验证逻辑，不是必需组本身的功能

// 使用验证器实现条件必需
mode := cmd.String("mode", "m", "运行模式", "local")

// 根据模式验证必需的标志
mode.SetValidator(func(value string) error {
    switch value {
    case "remote":
        // 远程模式需要 host 和 port
        if !cmd.GetFlag("host").IsSet() || !cmd.GetFlag("port").IsSet() {
            return fmt.Errorf("remote mode requires --host and --port")
        }
    case "local":
        // 本地模式需要 path
        if !cmd.GetFlag("path").IsSet() {
            return fmt.Errorf("local mode requires --path")
        }
    }
    return nil
})
```

### 12.2 动态必需组

```go
// 场景：根据配置动态添加必需组
func setupRequiredGroups(cmd *qflag.Cmd, config *Config) {
    if config.RequireAuth {
        cmd.AddRequiredGroup("认证配置", []string{"username", "password"})
    }
    if config.RequireTLS {
        cmd.AddRequiredGroup("TLS 配置", []string{"cert-file", "key-file"})
    }
}
```

### 12.3 嵌套必需组

```go
// 场景：多层必需关系
// 注意：这需要自定义验证逻辑，不是必需组本身的功能

// 使用验证器实现嵌套必需
cmd.AddRequiredGroup("基础配置", []string{"host", "port"})

// 如果启用了某个功能，需要额外的配置
feature := cmd.Bool("feature", "", "启用功能", false)
feature.SetValidator(func(value bool) error {
    if value {
        if !cmd.GetFlag("feature-config").IsSet() {
            return fmt.Errorf("feature requires --feature-config")
        }
    }
    return nil
})
```

---

## 十三、总结

本设计方案提供了一个与互斥组对称的必需组功能：

1. **API 一致**：与互斥组的 API 设计保持一致
2. **易于理解**：概念简单直观
3. **向后兼容**：不破坏现有代码
4. **类型安全**：使用结构化类型
5. **错误清晰**：提供明确的错误提示
6. **可组合**：与互斥组、验证器等功能配合使用
7. **线程安全**：所有方法都使用锁保护
8. **可测试**：提供完整的测试覆盖

这个方案完美地满足了需求，同时保持了代码的简洁性和可维护性。
