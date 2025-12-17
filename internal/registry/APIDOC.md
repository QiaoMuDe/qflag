# Package registry

**Import Path:** `gitee.com/MM-Q/qflag/internal/registry`

Package registry 内部注册表管理。本包实现了内部组件的注册表管理功能，提供统一的组件注册、查找和管理机制，支持模块化的架构设计。

## 功能模块

- **内部注册表管理** - 实现了内部组件的注册表管理功能，提供统一的组件注册、查找和管理机制，支持模块化的架构设计

## 目录

- [函数](#函数)
  - [RegisterFlag](#registerflag)
  - [ValidateFlagNames](#validateflagnames)

## 函数

### RegisterFlag

```go
func RegisterFlag(ctx *types.CmdContext, flag flags.Flag, longName, shortName string) error
```

RegisterFlag 注册标志。纯函数设计，通过参数传递所有必要信息。

**参数:**
- `ctx`: 命令上下文，用于存储注册的标志信息
- `flag`: 要注册的标志对象，实现了 `flags.Flag` 接口
- `longName`: 标志的长名称（如 "help", "version"）
- `shortName`: 标志的短名称（如 "h", "v"）

**返回值:**
- `error`: 如果注册失败（如名称冲突、无效名称等），返回错误信息；成功时返回 nil

**功能特点:**
- 纯函数设计，无副作用
- 支持长短名称的双重注册
- 自动进行名称冲突检测
- 提供统一的标志注册接口

**注册流程:**
1. 验证标志名称的有效性
2. 检查名称是否已被占用
3. 将标志添加到上下文的注册表中
4. 建立名称到标志的映射关系

### ValidateFlagNames

```go
func ValidateFlagNames(ctx *types.CmdContext, longName, shortName string) error
```

ValidateFlagNames 验证标志名称的有效性和唯一性。

**参数:**
- `ctx`: 命令上下文，包含已注册标志的信息
- `longName`: 要验证的长标志名称
- `shortName`: 要验证的短标志名称

**返回值:**
- `error`: 如果标志名称无效或已存在，则返回错误；否则返回 nil

**验证规则:**
- **名称格式验证**: 检查名称是否符合标志命名规范
- **唯一性验证**: 确保名称在当前上下文中唯一
- **长度验证**: 检查名称长度是否在合理范围内
- **字符验证**: 验证名称只包含允许的字符

**验证内容:**
- 长名称不能为空且符合命名规范
- 短名称应为单个字符（如果提供）
- 名称不能与已注册的标志冲突
- 名称不能使用保留字符或关键字

## 使用示例

### 基本标志注册

```go
package main

import (
    "gitee.com/MM-Q/qflag/internal/registry"
    "gitee.com/MM-Q/qflag/internal/types"
    "gitee.com/MM-Q/qflag/flags"
)

func main() {
    // 创建命令上下文
    ctx := &types.CmdContext{
        Flags: make(map[string]flags.Flag),
    }
    
    // 创建一个布尔标志
    helpFlag := flags.NewBoolFlag(false)
    
    // 注册标志
    err := registry.RegisterFlag(ctx, helpFlag, "help", "h")
    if err != nil {
        panic(err)
    }
    
    // 现在可以通过 "help" 或 "h" 访问这个标志
}
```

### 批量注册标志

```go
func registerCommonFlags(ctx *types.CmdContext) error {
    // 定义要注册的标志
    flagsToRegister := []struct {
        flag      flags.Flag
        longName  string
        shortName string
    }{
        {flags.NewBoolFlag(false), "help", "h"},
        {flags.NewBoolFlag(false), "version", "v"},
        {flags.NewStringFlag(""), "config", "c"},
        {flags.NewIntFlag(0), "port", "p"},
        {flags.NewBoolFlag(false), "verbose", ""},
    }
    
    // 批量注册
    for _, f := range flagsToRegister {
        if err := registry.RegisterFlag(ctx, f.flag, f.longName, f.shortName); err != nil {
            return fmt.Errorf("注册标志 %s 失败: %w", f.longName, err)
        }
    }
    
    return nil
}
```

### 名称验证示例

```go
func validateAndRegister(ctx *types.CmdContext, flag flags.Flag, longName, shortName string) error {
    // 先验证名称
    if err := registry.ValidateFlagNames(ctx, longName, shortName); err != nil {
        return fmt.Errorf("标志名称验证失败: %w", err)
    }
    
    // 验证通过后注册
    if err := registry.RegisterFlag(ctx, flag, longName, shortName); err != nil {
        return fmt.Errorf("标志注册失败: %w", err)
    }
    
    return nil
}

func main() {
    ctx := &types.CmdContext{}
    
    // 注册第一个标志
    debugFlag := flags.NewBoolFlag(false)
    if err := validateAndRegister(ctx, debugFlag, "debug", "d"); err != nil {
        panic(err)
    }
    
    // 尝试注册冲突的标志（会失败）
    verboseFlag := flags.NewBoolFlag(false)
    if err := validateAndRegister(ctx, verboseFlag, "debug", "v"); err != nil {
        fmt.Printf("预期的错误: %v\n", err)
    }
}
```

### 动态标志注册

```go
type FlagConfig struct {
    LongName    string
    ShortName   string
    FlagType    string
    DefaultValue interface{}
    Desc string
}

func registerFlagsFromConfig(ctx *types.CmdContext, configs []FlagConfig) error {
    for _, config := range configs {
        // 根据配置创建标志
        var flag flags.Flag
        switch config.FlagType {
        case "bool":
            flag = flags.NewBoolFlag(config.DefaultValue.(bool))
        case "string":
            flag = flags.NewStringFlag(config.DefaultValue.(string))
        case "int":
            flag = flags.NewIntFlag(config.DefaultValue.(int))
        default:
            return fmt.Errorf("不支持的标志类型: %s", config.FlagType)
        }
        
        // 注册标志
        if err := registry.RegisterFlag(ctx, flag, config.LongName, config.ShortName); err != nil {
            return fmt.Errorf("注册标志 %s 失败: %w", config.LongName, err)
        }
    }
    
    return nil
}
```

## 注册规则

### 名称规范

- **长名称**: 
  - 必须以字母开头
  - 可包含字母、数字、连字符（-）
  - 长度建议在 2-20 个字符之间
  - 使用小写字母和连字符的组合

- **短名称**:
  - 必须是单个字符
  - 通常使用字母（a-z, A-Z）
  - 避免使用数字和特殊字符

### 冲突检测

- 长名称在同一上下文中必须唯一
- 短名称在同一上下文中必须唯一
- 不区分大小写的重复检测
- 支持层级上下文的名称隔离

### 保留名称

以下名称为系统保留，不能用于自定义标志：
- `help`, `h` - 帮助信息
- `version`, `v` - 版本信息
- 以 `_` 开头的名称 - 内部使用

## 错误处理

注册器会返回以下类型的错误：

- **名称冲突错误**: 标志名称已被占用
- **无效名称错误**: 名称不符合命名规范
- **保留名称错误**: 尝试使用系统保留名称
- **空名称错误**: 长名称为空或无效
- **类型错误**: 标志对象类型不正确

## 设计特点

1. **纯函数设计** - 所有函数都是纯函数，不依赖全局状态
2. **上下文隔离** - 通过 `CmdContext` 实现标志的作用域隔离
3. **类型安全** - 使用接口确保标志对象的类型安全
4. **错误友好** - 提供详细的错误信息和建议
5. **扩展性强** - 支持自定义标志类型和验证规则

## 性能考虑

- 注册操作为 O(1) 时间复杂度（基于哈希表）
- 名称验证为 O(n) 时间复杂度，其中 n 为已注册标志数量
- 内存使用与注册标志数量成正比
- 建议在程序启动时一次性完成所有标志注册

## 最佳实践

1. **统一注册**: 在程序启动时统一注册所有标志
2. **名称规范**: 遵循一致的命名规范
3. **错误处理**: 妥善处理注册过程中的错误
4. **文档化**: 为每个标志提供清晰的描述
5. **测试覆盖**: 确保注册逻辑有充分的测试覆盖

## 注意事项

- 标志注册应在解析之前完成
- 注册顺序不影响解析结果
- 上下文销毁时会自动清理注册的标志
- 不支持运行时动态注销标志
- 标志名称一旦注册不可修改