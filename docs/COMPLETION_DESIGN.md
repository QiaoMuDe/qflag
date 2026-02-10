# 补全脚本生成器设计文档

## 概述

本文档描述了 qflag 项目的命令行自动补全脚本生成器的设计方案, 采用简化的架构实现。

## 设计目标

1. 为 qflag 项目提供命令行自动补全功能
2. 支持 Bash 和 PowerShell 两种主流 Shell
3. 与现有架构无缝集成
4. 保持代码简洁、可维护

## 核心设计

### 1. 核心接口

```go
// CompletionGenerator 补全生成器接口
type CompletionGenerator interface {
    Generate(cmd types.Command, shellType string) (string, error)
}
```

### 2. 数据结构

```go
// FlagInfo 标志信息
type FlagInfo struct {
    Name        string   // 标志长名称
    ShortName   string   // 标志短名称
    Desc        string   // 标志描述
    Type        string   // 标志类型
    EnumValues  []string // 枚举值 (如果有) 
}

// CommandInfo 命令信息
type CommandInfo struct {
    Name        string        // 命令长名称
    ShortName   string        // 命令短名称
    Desc        string        // 命令描述
    Flags       []FlagInfo    // 命令标志列表
    SubCommands []CommandInfo // 子命令列表
}
```

### 3. 实现结构

```go
// DefaultCompletionGenerator 默认补全生成器
type DefaultCompletionGenerator struct{}

// Generate 生成补全脚本
func (g *DefaultCompletionGenerator) Generate(cmd types.Command, shellType string) (string, error) {
    // 收集命令信息
    cmdInfo := g.collectCommandInfo(cmd)
    
    // 根据shell类型生成脚本
    switch shellType {
    case "bash":
        return g.generateBashScript(cmdInfo)
    case "pwsh":
        return g.generatePwshScript(cmdInfo)
    default:
        return "", fmt.Errorf("unsupported shell type: %s", shellType)
    }
}

// collectCommandInfo 收集命令信息
func (g *DefaultCompletionGenerator) collectCommandInfo(cmd types.Command) CommandInfo {
    // 递归收集命令和子命令信息
    // 转换标志类型为 FlagInfo
}

// generateBashScript 生成Bash补全脚本
func (g *DefaultCompletionGenerator) generateBashScript(cmdInfo CommandInfo) (string, error) {
    // 使用模板生成Bash脚本
}

// generatePwshScript 生成PowerShell补全脚本
func (g *DefaultCompletionGenerator) generatePwshScript(cmdInfo CommandInfo) (string, error) {
    // 使用模板生成PowerShell脚本
}
```

## 集成方式

### 1. 内置标志系统

```go
// CompletionHandler 补全标志处理器
type CompletionHandler struct {
    generator CompletionGenerator
}

// Handle 处理补全标志
func (h *CompletionHandler) Handle(cmd types.Command) error {
    // 获取shell类型参数
    shellType := getShellTypeFromArgs()
    
    // 生成脚本
    script, err := h.generator.Generate(cmd, shellType)
    if err != nil {
        return err
    }
    
    // 输出脚本
    fmt.Println(script)
    os.Exit(0)
    return nil
}

// ShouldRegister 判断是否注册补全标志
func (h *CompletionHandler) ShouldRegister(cmd types.Command) bool {
    return true // 所有命令都支持补全
}

// Type 返回标志类型
func (h *CompletionHandler) Type() types.BuiltinFlagType {
    return types.CompletionFlag // 新增的补全标志类型
}
```

### 2. 内置标志类型扩展

```go
// 在 types/builtin.go 中添加
const (
    // ...
    CompletionFlag BuiltinFlagType = "completion"
)
```

## 文件结构

```
internal/completion/
├── completion.go       # 核心接口和数据结构
├── bash.go           # Bash脚本生成
├── pwsh.go           # PowerShell脚本生成
└── templates/        # 脚本模板
    ├── bash.tmpl
    └── pwsh.tmpl
```

## 实现步骤

1. 创建核心接口和数据结构
2. 实现命令信息收集逻辑
3. 创建Bash和PowerShell模板
4. 实现脚本生成逻辑
5. 集成到内置标志系统
6. 添加测试用例

## 模板处理

使用 Go 标准库的 `text/template` 处理模板, 将收集的命令信息填充到预定义的模板中。

## 国际化支持

根据命令的 `UseChinese` 配置生成相应语言的描述信息。

## 优势

1. **简洁性**: 最少的接口和数据结构
2. **可扩展性**: 易于添加新的Shell支持
3. **集成性**: 与现有架构无缝集成
4. **可维护性**: 代码结构清晰, 职责明确