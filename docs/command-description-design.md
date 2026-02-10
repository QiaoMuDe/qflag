# CmdSpec 命令规格结构体设计方案

## 概述

本文档描述了通过命令规格结构体创建命令的设计方案, 旨在提供一种更直观、集中的命令配置方式, 替代当前的函数式配置方法。

## 设计目标

1. **配置集中化** - 将命令的所有属性集中在一个结构体中
2. **减少代码重复** - 避免多次调用Set方法
3. **提高可维护性** - 命令定义集中, 便于修改和扩展
4. **支持默认值** - 在结构体中设置默认值, 简化配置
5. **保持兼容性** - 与现有代码完全兼容

## 核心设计

### 1. 命令规格结构体

```go
// CmdSpec 命令规格结构体
//
// CmdSpec 提供了通过规格创建命令的方式, 包含命令的所有属性。
// 这种方式比函数式配置更加直观和集中。
type CmdSpec struct {
    // 基本属性
    LongName     string                    // 命令长名称
    ShortName    string                    // 命令短名称
    Desc         string                    // 命令描述
    ErrorHandling types.ErrorHandling      // 错误处理策略
    
    // 运行函数
    RunFunc      func(Command) error       // 命令执行函数
    
    // 配置选项
    Version      string                    // 版本号
    UseChinese   bool                      // 是否使用中文
    EnvPrefix    string                    // 环境变量前缀
    UsageSyntax  string                    // 命令使用语法
    LogoText     string                    // Logo文本
    
    // 示例和说明
    Examples     map[string]string         // 示例使用, key为描述, value为示例命令
    Notes        []string                  // 注意事项
    
    // 子命令和互斥组
    SubCmds      []Command                // 子命令列表
    MutexGroups  []MutexGroup             // 互斥组列表
}
```

### 2. 创建函数

```go
// NewCmdFromSpec 从规格创建命令
//
// 参数:
//   - spec: 命令规格结构体
//
// 返回值:
//   - *Cmd: 创建的命令实例
//   - error: 创建失败时返回错误
//
// 功能说明: 
//   - 根据规格结构体创建命令
//   - 自动设置所有属性和配置
//   - 添加子命令
//   - 支持默认值处理
func NewCmdFromSpec(spec *CmdSpec) (*Cmd, error) {
    // 验证规格结构体
    if spec == nil {
        return nil, errors.New("command spec cannot be nil")
    }
    
    // 创建基本命令
    cmd := NewCmd(spec.LongName, spec.ShortName, spec.ErrorHandling)
    
    // 设置基本属性
    cmd.SetDesc(spec.Desc)
    cmd.SetRun(spec.RunFunc)
    
    // 设置配置选项
    if spec.Version != "" {
        cmd.SetVersion(spec.Version)
    }
    cmd.SetChinese(spec.UseChinese)
    if spec.EnvPrefix != "" {
        cmd.SetEnvPrefix(spec.EnvPrefix)
    }
    if spec.UsageSyntax != "" {
        cmd.SetUsageSyntax(spec.UsageSyntax)
    }
    if spec.LogoText != "" {
        cmd.SetLogoText(spec.LogoText)
    }
    
    // 添加示例和说明
    if len(spec.Examples) > 0 {
        cmd.AddExamples(spec.Examples)
    }
    if len(spec.Notes) > 0 {
        cmd.AddNotes(spec.Notes)
    }
    
    // 添加互斥组
    for _, group := range spec.MutexGroups {
        cmd.AddMutexGroup(group.Name, group.Flags, group.AllowNone)
    }
    
    // 添加子命令
    if len(spec.SubCmds) > 0 {
        if err := cmd.AddSubCmds(spec.SubCmds...); err != nil {
            return nil, fmt.Errorf("failed to add subcommands: %w", err)
        }
    }
    
    return cmd, nil
}
```

### 3. 便捷构造函数

```go
// NewCmdSpec 创建新的命令规格
//
// 参数:
//   - longName: 命令长名称
//   - shortName: 命令短名称
//
// 返回值:
//   - *CmdSpec: 初始化的命令规格
func NewCmdSpec(longName, shortName string) *CmdSpec {
    return &CmdSpec{
        LongName:     longName,
        ShortName:    shortName,
        ErrorHandling: types.ExitOnError, // 默认错误处理策略
        UseChinese:   false,              // 默认不使用中文
        Examples:     make(map[string]string),
        Notes:        []string{},
        SubCmds:      []*CmdSpec{},
        MutexGroups:  []MutexGroup{},
    }
}
```

## 使用示例

### 基本用法

```go
// 使用规格结构体创建命令
appSpec := NewCmdSpec("myapp", "app")
appSpec.Desc = "我的应用程序"
appSpec.Version = "1.0.0"
appSpec.UseChinese = true
appSpec.EnvPrefix = "MYAPP"
appSpec.RunFunc = func(cmd types.Command) error {
    fmt.Println("运行应用程序")
    return nil
}

// 创建命令
app, err := NewCmdFromSpec(appSpec)
if err != nil {
    log.Fatal(err)
}

// 添加标志
app.String("input", "i", "输入文件", "")
app.Bool("verbose", "v", "详细输出", false)

// 添加互斥组
app.AddMutexGroup("format", []string{"json", "xml"}, true)
}

// 创建命令
app, err := NewCmdFromSpec(appSpec)
if err != nil {
    log.Fatal(err)
}
```

### 嵌套子命令

```go
// 创建子命令
subCmd := NewCmd("subcommand", "sub", types.ExitOnError)
subCmd.SetDesc("子命令")

// 添加到主命令
appSpec.SubCmds = []types.Command{subCmd}

// 创建命令
app, err := NewCmdFromSpec(appSpec)
if err != nil {
    log.Fatal(err)
}

// 获取子命令并添加标志
retrievedSubCmd, _ := app.GetSubCmd("sub")
retrievedSubCmd.String("option", "o", "子命令选项", "")
```

### 复杂配置示例

```go
// 复杂命令配置
complexSpec := NewCmdSpec("complex", "cpx")
complexSpec.Desc = "复杂命令示例"
complexSpec.Version = "2.0.0"
complexSpec.UseChinese = true
complexSpec.EnvPrefix = "COMPLEX"
complexSpec.UsageSyntax = "[options] <args>"
complexSpec.LogoText = "Complex Command v2.0.0"

// 添加示例
complexSpec.Examples = map[string]string{
    "基本用法": "complex --input file.txt",
    "详细模式": "complex --input file.txt --verbose",
    "输出JSON": "complex --input file.txt --json",
}

// 添加注意事项
complexSpec.Notes = []string{
    "输入文件必须存在",
    "输出目录必须可写",
    "处理大文件时请增加内存限制",
}

// 添加子命令
processCmd := NewCmd("process", "proc", types.ExitOnError)
processCmd.SetDesc("处理数据")

validateCmd := NewCmd("validate", "val", types.ExitOnError)
validateCmd.SetDesc("验证数据")

complexSpec.SubCmds = []types.Command{processCmd, validateCmd}

// 创建命令
complex, err := NewCmdFromSpec(complexSpec)
if err != nil {
    log.Fatal(err)
}

// 添加多个标志
complex.String("input", "i", "输入文件", "")
complex.String("output", "o", "输出文件", "")
complex.Bool("verbose", "v", "详细输出", false)
complex.Bool("json", "j", "JSON格式输出", false)
complex.Bool("xml", "x", "XML格式输出", false)
complex.Int("limit", "l", "处理限制", 1000)

// 添加互斥组
complex.AddMutexGroup("output_format", []string{"json", "xml"}, true)
```

## 实现优势

1. **配置集中化** - 所有配置在一个结构体中, 便于查看和维护
2. **支持嵌套** - 可以递归创建子命令, 保持结构清晰
3. **默认值支持** - 可以在构造函数中设置合理的默认值
4. **错误处理** - 集中的错误处理, 便于调试
5. **类型安全** - 编译时检查, 减少运行时错误
6. **易于测试** - 可以轻松创建测试用的命令描述
7. **代码可读性** - 命令结构一目了然, 便于理解

## 与现有代码的兼容性

这个方案完全兼容现有代码: 

1. **保留现有API** - 保留现有的 `NewCmd` 函数和所有方法
2. **新增功能** - 新增 `NewCmdFromSpec` 函数和 `CmdSpec` 结构体
3. **并存方式** - 两种创建方式可以并存使用
4. **渐进迁移** - 可以逐步将现有代码迁移到新方式

## 实施建议

1. **第一阶段** - 实现 `CmdSpec` 结构体和 `NewCmdFromSpec` 函数
2. **第二阶段** - 更新示例代码, 展示新用法
3. **第三阶段** - 逐步迁移现有代码到新方式
4. **第四阶段** - 考虑将新方式设为推荐用法

## 扩展可能性

1. **配置文件支持** - 可以从JSON/YAML文件加载 `CmdSpec`
2. **代码生成** - 可以从配置文件生成命令代码
3. **验证器** - 可以添加 `CmdSpec` 的验证逻辑
4. **构建器模式** - 可以结合构建器模式提供更灵活的配置

## 总结

通过命令规格结构体创建命令的方案提供了一种更直观、集中的命令配置方式, 同时保持了与现有代码的完全兼容性。这种设计既符合Go语言的惯例, 又提高了代码的可读性和可维护性, 是一个值得实施的改进方案。