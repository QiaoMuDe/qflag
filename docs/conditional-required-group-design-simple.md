# 条件性必需组设计方案（简化版）

## 概述

本方案旨在为 QFlag 项目添加条件性必需组功能，实现"只要使用了组中的任何一个标志，就必须使用组中的所有标志"的逻辑。此方案通过直接在 `AddRequiredGroup` 方法中添加一个参数，保持与 `AddMutexGroup` 方法的一致性。

## 设计目标

1. **最小改动**：只添加一个字段和一个参数
2. **API一致性**：与 `AddMutexGroup` 方法保持一致的参数模式
3. **向后兼容**：现有的必需组功能不受影响
4. **简洁明了**：直接通过参数控制是否为条件性必需组

## 实现方案

### 1. 修改 RequiredGroup 结构

在 `internal/types/config.go` 中，给 `RequiredGroup` 结构添加一个 `Conditional` 字段：

```go
// RequiredGroup 必需组定义
//
// RequiredGroup 定义了一组必需的标志，其中所有标志都必须被设置。
// 当用户没有设置必需组中的某些标志时，解析器会返回错误。
//
// 字段说明:
//   - Name: 必需组名称，用于错误提示和标识
//   - Flags: 必需组中的标志名称列表
//   - Conditional: 是否为条件性必需组，如果为true，则只有当组中任何一个标志被设置时，才要求所有标志都被设置
//
// 使用场景:
//   - 数据库连接配置（host, port, user, pass）
//   - API 认证配置（api-key, api-secret）
//   - 文件上传配置（file-path, upload-url）
//   - 条件性配置（如果使用了任何一个标志，则必须使用所有标志）
type RequiredGroup struct {
    Name        string   // 必需组名称，用于错误提示和标识
    Flags       []string // 必需组中的标志名称列表
    Conditional bool     // 是否为条件性必需组
}
```

### 2. 修改 AddRequiredGroup 方法

在 `internal/cmd/cmd.go` 中，修改现有的 `AddRequiredGroup` 方法，添加一个条件性参数：

```go
// AddRequiredGroup 添加必需组
//
// 参数:
//   - name: 必需组名称
//   - flags: 必需组中的标志名称列表
//   - conditional: 是否为条件性必需组，如果为true，则只有当组中任何一个标志被设置时，才要求所有标志都被设置
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 添加一个必需组到命令配置
//   - 如果组名已存在，返回错误
//   - 如果标志列表为空，返回错误
//   - 如果标志不存在，返回错误
//   - 支持条件性必需组，当conditional为true时，只有当组中任何一个标志被设置时，才要求所有标志都被设置
//
// 错误码:
//   - REQUIRED_GROUP_ALREADY_EXISTS: 必需组已存在
//   - EMPTY_REQUIRED_GROUP: 必需组标志列表为空
//   - FLAG_NOT_FOUND: 标志不存在
func (c *Cmd) AddRequiredGroup(name string, flags []string, conditional bool) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    // 必需组名称不能为空
    if name == "" {
        return types.NewError("EMPTY_REQUIRED_GROUP_NAME",
            "required group name cannot be empty", nil)
    }

    // 检查必需组名称是否已存在
    for _, group := range c.config.RequiredGroups {
        if group.Name == name {
            return types.NewError("REQUIRED_GROUP_ALREADY_EXISTS",
                fmt.Sprintf("required group '%s' already exists", name), nil)
        }
    }

    // 必需组标志列表不能为空
    if len(flags) == 0 {
        return types.NewError("EMPTY_REQUIRED_GROUP",
            "required group cannot be empty", nil)
    }

    // 检查必需组标志是否存在
    for _, flagName := range flags {
        if _, exists := c.flagRegistry.Get(flagName); !exists {
            return types.NewError("FLAG_NOT_FOUND",
                fmt.Sprintf("flag '%s' not found", flagName), nil)
        }
    }

    // 添加必需组
    group := types.RequiredGroup{
        Name:        name,
        Flags:       flags,
        Conditional: conditional,
    }

    c.config.RequiredGroups = append(c.config.RequiredGroups, group)
    return nil
}
```

### 3. 修改验证逻辑

在 `internal/parser/parser_validation.go` 中，修改 `validateRequiredGroups` 函数：

```go
// validateRequiredGroups 验证命令的必需组规则
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 如果必需组验证失败返回错误
//
// 功能说明:
//   - 检查每个必需组中是否有标志未被设置
//   - 对于条件性必需组，只有当组中任何一个标志被设置时，才检查所有标志
//   - 提供清晰的错误信息，指出未设置的标志和组名
//
// 验证规则:
//   - 普通必需组：所有标志都必须被设置
//   - 条件性必需组：如果任何一个标志被设置，则所有标志都必须被设置
//
// 错误处理:
//   - 使用 types.NewError 创建结构化错误
//   - 错误信息包含必需组名称和未设置的标志列表
//   - 条件性必需组的错误信息会明确指出是因为使用了某个标志而要求其他标志
//
// 性能优化:
//   - 使用缓存的已设置标志映射，避免重复的 GetFlag() 和 IsSet() 调用
func (p *DefaultParser) validateRequiredGroups(config *types.CmdConfig) error {
    if len(config.RequiredGroups) == 0 {
        return nil
    }

    // 使用缓存的已设置标志映射
    setFlags := p.setFlagsMap

    // 遍历所有必需组
    for _, group := range config.RequiredGroups {
        var unsetFlags []string
        var setFlagsCount int
        
        // 检查组中是否有任何标志被设置
        for _, flagName := range group.Flags {
            if setFlags[flagName] {
                setFlagsCount++
            }
        }
        
        // 如果是条件性必需组，且没有任何标志被设置，则跳过验证
        if group.Conditional && setFlagsCount == 0 {
            continue
        }
        
        seenUnsetDisplayNames := make(map[string]bool, len(group.Flags)) // 去重map，防止重复显示相同的标志

        // 遍历组中的每个标志
        for _, flagName := range group.Flags {
            // 如果拿组里的标志没有获取到显示名称, 则表示为不是一个有效标志, 返回错误
            displayName, ok := p.flagDisplayNames[flagName]
            if !ok {
                return types.NewError("INVALID_FLAG_NAME",
                    fmt.Sprintf("invalid flag name '%s' in required group '%s'", flagName, group.Name),
                    nil)
            }

            // 如果标志未被设置, 添加到未设置列表
            if !setFlags[flagName] {
                if !seenUnsetDisplayNames[displayName] {
                    // 添加去重检查，避免同一个标志的多个名称重复显示
                    seenUnsetDisplayNames[displayName] = true
                    unsetFlags = append(unsetFlags, displayName)
                }
            }
        }

        // 如果组中有未设置的标志, 返回错误
        if len(unsetFlags) > 0 {
            var errorMsg string
            if group.Conditional {
                errorMsg = fmt.Sprintf("since one of the flags in group '%s' is used, all flags %v must be set", group.Name, unsetFlags)
            } else {
                errorMsg = fmt.Sprintf("required flags %v in group '%s' must be set", unsetFlags, group.Name)
            }
            
            return types.NewError("REQUIRED_GROUP_VIOLATION", errorMsg, nil)
        }
    }

    return nil
}
```

### 4. 更新相关结构体

在 `internal/cmd/cmdspec.go` 和 `internal/cmd/cmdopts.go` 中，更新 `CmdSpec` 和 `CmdOpts` 结构，支持条件性必需组：

```go
// CmdSpec 命令规格结构体
//
// CmdSpec 提供了通过规格创建命令的方式, 包含命令的所有属性。
// 这种方式比函数式配置更加直观和集中。
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
    SubCmds        []types.Command       // 子命令列表, 用于添加到命令中
    MutexGroups    []types.MutexGroup    // 互斥组列表
    RequiredGroups []types.RequiredGroup // 必需组列表（支持条件性）
}
```

```go
// CmdOpts 命令选项
//
// CmdOpts 提供了配置现有命令的方式，包含命令的所有可配置属性。
// 与 CmdSpec 不同，CmdOpts 用于配置已存在的命令，而不是创建新命令。
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
    SubCmds        []types.Command       // 子命令列表, 用于添加到命令中
    MutexGroups    []types.MutexGroup    // 互斥组列表
    RequiredGroups []types.RequiredGroup // 必需组列表（支持条件性）
}
```

### 5. 使用示例

```go
// 创建命令
cmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)

// 添加标志
hostFlag := cmd.String("host", "h", "数据库主机地址", "")
portFlag := cmd.String("port", "p", "数据库端口", "")
dbFlag := cmd.String("database", "d", "数据库名称", "")
userFlag := cmd.String("username", "u", "用户名", "")
passFlag := cmd.String("password", "P", "密码", "")

// 添加普通必需组（所有标志都必须设置）
cmd.AddRequiredGroup("auth", []string{"username", "password"}, false)

// 添加条件性必需组（如果使用了任何一个标志，则所有标志都必须设置）
cmd.AddRequiredGroup("database", []string{"host", "port", "database"}, true)

// 解析参数
err := cmd.Parse(os.Args[1:])
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}
```

## 实际应用场景

### 1. 数据库连接配置
```go
// 普通必需组：必须提供用户名和密码
cmd.AddRequiredGroup("auth", []string{"username", "password"}, false)

// 条件性必需组：如果使用数据库，则必须提供完整连接信息
cmd.AddRequiredGroup("database", []string{"host", "port", "database"}, true)

// 使用示例：
// ./app --username admin --password secret                     // 成功：提供认证信息
// ./app --host localhost --port 5432 --database mydb        // 成功：提供数据库信息
// ./app --host localhost                                    // 失败：需要提供port和database
// ./app                                                     // 成功：不使用数据库和认证功能
```

### 2. 文件处理
```go
// 条件性必需组：如果处理文件，则必须提供输入和输出文件
cmd.AddRequiredGroup("file_processing", []string{"input-file", "output-file"}, true)

// 使用示例：
// ./app --input-file input.txt --output-file output.txt  // 成功
// ./app --input-file input.txt                        // 失败：需要提供output-file
// ./app                                              // 成功：不处理文件
```

### 3. 网络请求
```go
// 条件性必需组：如果发送请求，则必须提供完整请求信息
cmd.AddRequiredGroup("http_request", []string{"url", "method", "timeout"}, true)

// 使用示例：
// ./app --url https://api.example.com --method GET --timeout 30s  // 成功
// ./app --url https://api.example.com                           // 失败：需要提供method和timeout
// ./app                                                          // 成功：不发送网络请求
```

## 设计优势

1. **最小改动**：只添加一个字段和一个参数，不添加新方法
2. **API一致性**：与 `AddMutexGroup(name, flags, allowNone)` 方法保持一致的参数模式
3. **向后兼容**：现有的必需组功能不受影响，只需在调用时添加 `false` 参数
4. **简洁明了**：直接通过布尔参数控制是否为条件性必需组
5. **统一体验**：所有相关结构体都支持条件性必需组

## 测试用例

需要添加的测试用例：

1. **普通必需组**：所有标志都必须设置，应该按原逻辑工作
2. **条件性必需组，不使用任何标志**：应该成功
3. **条件性必需组，使用部分标志**：应该失败
4. **条件性必需组，使用所有标志**：应该成功
5. **普通必需组和条件性必需组混合使用**：应该正常工作
6. **多个条件性必需组同时使用**：应该正常工作
7. **条件性必需组与互斥组组合使用**：应该正常工作
8. **CmdSpec和CmdOpts支持条件性必需组**：应该正常工作

### 测试用例示例

```go
func TestConditionalRequiredGroup(t *testing.T) {
    // 测试1: 普通必需组，应该按原逻辑工作
    func() {
        cmd := NewCmd("test1", "t1", types.ContinueOnError)
        cmd.String("username", "u", "Username", "")
        cmd.String("password", "p", "Password", "")
        cmd.AddRequiredGroup("auth", []string{"username", "password"}, false) // 普通必需组
        
        err := cmd.Parse([]string{})
        if err == nil {
            t.Error("Expected error when no flags are set in normal required group")
        }
    }()
    
    // 测试2: 条件性必需组，不使用任何标志，应该成功
    func() {
        cmd := NewCmd("test2", "t2", types.ContinueOnError)
        cmd.String("host", "h", "Host", "")
        cmd.String("port", "p", "Port", "")
        cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
        
        err := cmd.Parse([]string{})
        if err != nil {
            t.Errorf("Expected no error when no flags are used in conditional group, got: %v", err)
        }
    }()
    
    // 测试3: 条件性必需组，使用部分标志，应该失败
    func() {
        cmd := NewCmd("test3", "t3", types.ContinueOnError)
        cmd.String("host", "h", "Host", "")
        cmd.String("port", "p", "Port", "")
        cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
        
        err := cmd.Parse([]string{"--host", "localhost"})
        if err == nil {
            t.Error("Expected error when only some flags in conditional group are used")
        }
    }()
    
    // 测试4: 条件性必需组，使用所有标志，应该成功
    func() {
        cmd := NewCmd("test4", "t4", types.ContinueOnError)
        cmd.String("host", "h", "Host", "")
        cmd.String("port", "p", "Port", "")
        cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
        
        err := cmd.Parse([]string{"--host", "localhost", "--port", "5432"})
        if err != nil {
            t.Errorf("Expected no error when all flags in conditional group are used, got: %v", err)
        }
    }()
}
```

## 总结

本方案通过最小改动实现了条件性必需组功能，提供了更加灵活的标志验证机制。主要特点：

1. **API一致性**：与 `AddMutexGroup(name, flags, allowNone)` 方法保持一致的参数模式
2. **最小改动**：只添加一个字段和一个参数，不添加新方法
3. **向后兼容**：现有的必需组功能不受影响，只需在调用时添加 `false` 参数
4. **简洁明了**：直接通过布尔参数控制是否为条件性必需组

这种设计可以满足更多实际场景的需求，如数据库连接、文件处理、网络请求等，同时保持了代码的简洁性和API的一致性。