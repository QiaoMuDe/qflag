# 智能纠错功能配置项设计方案

## 背景

智能纠错功能在用户输入错误的子命令或标志时，自动提供相似的建议。这是一个提升用户体验的功能，但可能会影响到一些特殊的使用场景：

- 用户没有设置运行函数，自己处理剩余参数（Args）进行自定义路由
- 用户希望保持与之前版本完全一致的行为

因此，需要添加一个配置项来控制是否启用智能纠错功能。

## 设计目标

1. **向后兼容** - 默认行为应与之前版本一致，或提供明确的迁移路径
2. **简单易用** - 配置项命名清晰，易于理解和使用
3. **灵活控制** - 支持全局配置和单个命令配置

## 配置项设计

### 配置项名称

```go
SmartSuggestion bool // 是否启用智能纠错（默认 false，即禁用）
```

**命名理由：**
- 简洁直观，不带冗余的 Enable/Disable 前缀
- 默认值为 `false`（禁用），保持与之前版本行为一致
- 用户通过 `SetSmartSuggestion(true)` 显式启用

### 配置位置

在 `types.CmdConfig` 结构体中添加：

```go
type CmdConfig struct {
    // ... 现有字段 ...
    Completion              bool // 是否启用自动补全标志
    DynamicCompletion bool // 是否启用动态补全
    SmartSuggestion         bool // 是否启用智能纠错（默认 false）
}
```

## 实现方案

### 1. 修改 types/config.go

```go
// NewCmdConfig 创建新的命令配置
func NewCmdConfig() *CmdConfig {
    return &CmdConfig{
        // ... 现有字段初始化 ...
        SmartSuggestion: false, // 默认禁用智能纠错
    }
}

// Clone 克隆命令配置
func (c *CmdConfig) Clone() *CmdConfig {
    // ... 现有克隆逻辑 ...
    clone.SmartSuggestion = c.SmartSuggestion
    // ...
}
```

### 2. 修改 parser/parser.go

在 `Parse()` 和 `ParseAndRoute()` 方法中，根据配置判断是否启用纠错：

```go
func (p *DefaultParser) Parse(cmd types.Command, args []string) error {
    // ... 前置逻辑 ...

    if len(remainingArgs) > 0 {
        firstArg := remainingArgs[0]

        if subCmd, ok := cmdRegistry.Get(firstArg); ok {
            return subCmd.Parse(remainingArgs[1:])
        }

        // 智能纠错判断：有子命令、未匹配、不以 - 开头、且启用了纠错
        config := cmd.Config()
        if len(cmd.SubCmds()) > 0 && 
            !strings.HasPrefix(firstArg, "-") && 
            config.SmartSuggestion {
            return newUnknownSubcommandError(cmd, firstArg)
        }
    }

    return nil
}
```

### 3. 修改 cmd/cmd_config.go

添加设置方法：

```go
// SetSmartSuggestion 设置是否启用智能纠错
//
// 参数:
//   - enable: true 表示启用智能纠错，false 表示禁用（默认）
//
// 功能说明:
//   - 控制当用户输入错误的子命令或标志时是否提供建议
//   - 默认禁用智能纠错，保持与之前版本行为一致
//   - 启用后，错误输入将返回带建议的错误
func (c *Cmd) SetSmartSuggestion(enable bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.config.SmartSuggestion = enable
}
```

### 4. 支持 CmdOpts 批量配置

在 `CmdOpts` 结构体中添加对应字段（如果存在该结构体）：

```go
type CmdOpts struct {
    // ... 现有字段 ...
    SmartSuggestion bool // 是否启用智能纠错
}
```

在 `ApplyOpts` 方法中应用配置：

```go
func (c *Cmd) ApplyOpts(opts *CmdOpts) error {
    // ... 现有逻辑 ...
    
    c.SetSmartSuggestion(opts.SmartSuggestion)
    
    // ...
}
```

## 使用示例

### 方式一：直接设置

```go
cmd := qflag.NewCmd("myapp", "", qflag.ContinueOnError)

// 启用智能纠错
cmd.SetSmartSuggestion(true)

// 解析参数
if err := cmd.Parse(os.Args[1:]); err != nil {
    log.Fatal(err)
}
```

### 方式二：通过 CmdOpts 配置

```go
opts := &qflag.CmdOpts{
    Desc:            "我的应用",
    SmartSuggestion: true, // 启用智能纠错
}

cmd.ApplyOpts(opts)
```

### 方式三：全局根命令

```go
// 启用全局根命令的智能纠错
qflag.Root.SetSmartSuggestion(true)
```

## 影响范围

### 受影响的文件

1. `internal/types/config.go` - 添加配置字段
2. `internal/parser/parser.go` - 根据配置判断是否启用纠错
3. `internal/cmd/cmd_config.go` - 添加设置方法
4. `internal/types/opts.go`（如果存在）- CmdOpts 添加字段
5. `internal/cmd/cmd_opts.go`（如果存在）- ApplyOpts 应用配置

### 向后兼容性

- **完全兼容** - 默认禁用智能纠错，保持与之前版本完全一致的行为
- **显式启用** - 需要智能纠错功能的用户可以显式启用

## 测试建议

1. **默认行为测试** - 验证默认情况下智能纠错正常工作
2. **禁用测试** - 验证禁用后不再返回建议错误
3. **配置继承测试** - 验证子命令是否正确继承父命令配置（如果需要）
4. **边界情况测试** - 验证无子命令、禁用标志解析等场景

## 文档更新

1. 更新 README.md - 添加智能纠错配置说明
2. 更新 APIDOC.md - 添加配置项和方法文档
3. 更新示例代码 - 添加禁用智能纠错的示例

## 实施计划

1. 修改 `types/config.go` - 添加配置字段
2. 修改 `parser/parser.go` - 集成配置判断
3. 修改 `cmd/cmd_config.go` - 添加设置方法
4. 更新 `CmdOpts` 相关代码（如果存在）
5. 编写测试用例
6. 更新文档
7. 验证示例代码

## 决策记录

### 为什么配置项名称不带 Enable/Disable 前缀？

- **简洁直观**：`SmartSuggestion` 比 `EnableSmartSuggestion` 更简洁
- **语义清晰**：bool 类型的配置项，true 表示启用，false 表示禁用，约定俗成
- **与现有配置保持一致**：如 `Completion`、`DynamicCompletion` 等

### 为什么默认禁用？

- **向后兼容**：保持与之前版本完全一致的行为，不破坏现有代码
- **显式启用**：用户需要时主动启用，避免意外行为改变
- **安全优先**：不强制改变用户的错误处理逻辑

## 相关文档

- [智能纠错功能影响分析](smart-error-suggestion-impact-analysis.md)
- [智能纠错功能设计方案](smart-error-suggestion-proposal.md)
