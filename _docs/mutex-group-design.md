# 互斥组功能设计方案

## 概述

本文档详细描述了为 qflag 参数解析库添加互斥组 (mutex group) 功能的完整设计方案。互斥组是一组标志, 其中最多只能有一个被设置。当用户设置了互斥组中的多个标志时, 应该返回错误。

## 设计目标

1. **最小侵入性**: 不破坏现有接口和实现, 保持向后兼容性
2. **架构一致性**: 利用现有的验证机制, 与 `Validator` 接口设计理念一致
3. **易于使用**: 提供简单直观的 API, 一行代码即可添加互斥组
4. **错误友好**: 提供清晰的错误信息, 指出冲突的标志和组名
5. **扩展性强**: 未来可以轻松添加更复杂的组关系 (如依赖组、条件组等) 

## 架构分析

### 当前架构特点

1. **标志管理**: 通过 `FlagRegistry` 接口管理标志, 支持注册、查找和遍历
2. **解析流程**: 在 `DefaultParser.validateFlags()` 方法中验证所有已设置的标志
3. **验证机制**: 每个标志支持通过 `Validator` 接口进行自定义验证
4. **命令结构**: `Cmd` 结构体包含标志注册器和配置信息

### 实现方式选择

经过分析, 我们选择**基于验证器的实现方式**, 原因如下: 

1. 不需要修改 `types.Command` 接口, 保持向后兼容性
2. 利用现有的验证机制, 与 `Validator` 接口设计理念一致
3. 只需要添加验证逻辑, 不需要扩展核心接口
4. 可以轻松扩展为更复杂的关系 (依赖组、条件组等) 

## 详细设计

### 1. 核心数据结构

```go
// MutexGroup 互斥组定义
type MutexGroup struct {
    Name      string   // 互斥组名称, 用于错误提示和标识
    Flags     []string // 互斥组中的标志名称列表
    AllowNone bool     // 是否允许一个都不设置
}

// 在 CmdConfig 中添加互斥组支持
type CmdConfig struct {
    // ... 现有字段 ...
    MutexGroups []MutexGroup // 互斥组列表
}
```

### 2. 核心验证逻辑

```go
// validateMutexGroups 验证互斥组规则
func (p *DefaultParser) validateMutexGroups(cmd types.Command) error {
    config := cmd.Config()
    if config == nil {
        return nil
    }
    
    for _, group := range config.MutexGroups {
        setCount := 0
        var setFlags []string
        
        // 检查互斥组中的每个标志是否被设置
        for _, flagName := range group.Flags {
            if flag, exists := cmd.GetFlag(flagName); exists && flag.IsSet() {
                setCount++
                setFlags = append(setFlags, flagName)
            }
        }
        
        // 验证互斥组规则
        if setCount > 1 {
            return fmt.Errorf("mutually exclusive flags %v in group '%s' cannot be used together", setFlags, group.Name)
        }
        
        if !group.AllowNone && setCount == 0 {
            return fmt.Errorf("one of the mutually exclusive flags %v in group '%s' must be set", group.Flags, group.Name)
        }
    }
    
    return nil
}
```

### 3. 解析流程集成

```go
// 在 DefaultParser.ParseOnly() 方法中集成验证
func (p *DefaultParser) ParseOnly(cmd types.Command, args []string) error {
    // ... 现有代码 ...
    
    // 验证所有设置了验证器的标志
    if err := p.validateFlags(cmd); err != nil {
        return err
    }
    
    // 验证互斥组规则
    if err := p.validateMutexGroups(cmd); err != nil {
        return err
    }
    
    // ... 其余代码 ...
}
```

### 4. 便捷方法

```go
// AddMutexGroup 添加互斥组到命令
func (c *Cmd) AddMutexGroup(name string, flags []string, allowNone bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    group := MutexGroup{
        Name:      name,
        Flags:     flags,
        AllowNone: allowNone,
    }
    
    c.config.MutexGroups = append(c.config.MutexGroups, group)
}

// GetMutexGroups 获取命令的所有互斥组
func (c *Cmd) GetMutexGroups() []MutexGroup {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    // 返回副本以防止外部修改
    groups := make([]MutexGroup, len(c.config.MutexGroups))
    copy(groups, c.config.MutexGroups)
    return groups
}

// RemoveMutexGroup 移除指定名称的互斥组
func (c *Cmd) RemoveMutexGroup(name string) bool {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    for i, group := range c.config.MutexGroups {
        if group.Name == name {
            c.config.MutexGroups = append(c.config.MutexGroups[:i], c.config.MutexGroups[i+1:]...)
            return true
        }
    }
    return false
}
```

## 使用示例

### 基本用法

```go
// 创建命令
cmd := cmd.NewCmd("myapp", "app", types.ContinueOnError)

// 添加标志
cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json"))
cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", ""))
cmd.AddFlag(flag.NewBoolFlag("verbose", "v", "Verbose output", false))

// 添加互斥组 - format 和 output 不能同时使用, 但可以都不使用
cmd.AddMutexGroup("output_format", []string{"format", "output"}, true)

// 添加互斥组 - 必须使用 input 或 source 中的一个
cmd.AddFlag(flag.NewStringFlag("input", "i", "Input file", ""))
cmd.AddFlag(flag.NewStringFlag("source", "s", "Source file", ""))
cmd.AddMutexGroup("input_source", []string{"input", "source"}, false)

// 解析参数
err := cmd.Parse(os.Args[1:])
if err != nil {
    fmt.Printf("Error: %v\n", err)
    os.Exit(1)
}
```

### 错误示例

```bash
# 互斥组冲突错误
$ myapp --format json --output result.txt
Error: mutually exclusive flags [format output] in group 'output_format' cannot be used together

# 必须设置一个标志错误
$ myapp
Error: one of the mutually exclusive flags [input source] in group 'input_source' must be set
```

## 实现步骤

### 第一阶段: 核心功能实现

1. 在 `types` 包中添加 `MutexGroup` 结构体定义
2. 在 `CmdConfig` 中添加 `MutexGroups` 字段
3. 在 `DefaultParser` 中实现 `validateMutexGroups` 方法
4. 在 `ParseOnly` 方法中集成互斥组验证
5. 在 `Cmd` 中实现 `AddMutexGroup` 方法

### 第二阶段: 完善功能

1. 实现 `GetMutexGroups` 和 `RemoveMutexGroup` 方法
2. 添加互斥组到帮助信息显示
3. 在补全脚本中考虑互斥组关系
4. 添加单元测试和集成测试

### 第三阶段: 扩展功能

1. 支持依赖组 (dependency group) 
2. 支持条件组 (conditional group) 
3. 支持嵌套组关系
4. 提供更丰富的错误信息和建议

## 测试计划

### 单元测试

1. **基本功能测试**
   - 测试互斥组的基本验证逻辑
   - 测试允许和不允许都不设置的情况
   - 测试错误消息的正确性

2. **边界情况测试**
   - 空的互斥组
   - 不存在的标志
   - 重复的标志名称
   - 嵌套的互斥组

3. **并发安全测试**
   - 多个 goroutine 同时添加互斥组
   - 并发访问互斥组信息

### 集成测试

1. **解析流程测试**
   - 测试在完整解析流程中的互斥组验证
   - 测试与其他验证器的交互

2. **命令行工具测试**
   - 创建使用互斥组的示例命令
   - 测试各种参数组合的行为

## 性能考虑

1. **验证时机**: 互斥组验证只在解析阶段进行一次, 不影响运行时性能
2. **数据结构**: 使用切片存储互斥组, 查找效率为 O(n), 但通常互斥组数量较少
3. **内存占用**: 互斥组只存储标志名称字符串, 内存占用很小
4. **并发安全**: 使用读写锁保护, 支持并发读取

## 向后兼容性

1. **接口兼容**: 不修改现有接口, 所有现有代码无需修改
2. **行为兼容**: 不使用互斥组时, 行为与原来完全一致
3. **配置兼容**: 现有配置文件和代码无需修改

## 扩展性设计

### 未来可能的扩展

1. **依赖组**: 一组标志中, 设置某个标志时必须同时设置其他标志
2. **条件组**: 根据某个标志的值, 决定其他标志是否可用
3. **嵌套组**: 组与组之间的复杂关系
4. **动态组**: 根据运行时条件动态确定组关系

### 扩展预留

1. 在 `MutexGroup` 中预留字段, 用于未来扩展
2. 使用接口抽象组关系, 便于实现不同类型的组
3. 提供扩展点, 允许用户自定义组关系验证逻辑

## 总结

本设计方案通过在现有验证机制基础上添加互斥组功能, 实现了最小侵入性的扩展。该方案保持了架构的一致性, 提供了简单易用的 API, 同时具备良好的扩展性。实现后, 用户可以轻松定义互斥组, 确保命令行参数的正确使用, 提升用户体验和工具的健壮性。