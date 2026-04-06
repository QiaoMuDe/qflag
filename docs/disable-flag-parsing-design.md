# DisableFlagParsing 功能设计文档

## 1. 功能概述

实现类似 Cobra 的 `DisableFlagParsing` 功能，允许命令禁用标志解析，将所有参数（包括以 `-` 或 `--` 开头的参数）都作为位置参数处理。

## 2. 使用场景

- **特殊子命令**: 如 `__complete` 补全命令，需要接收 `--` 开头的参数而不解析为标志
- **透传参数**: 需要将参数原样传递给下游命令的场景
- **自定义解析**: 命令需要自行处理所有参数，不使用标准标志解析

## 3. 设计方案

### 3.1 数据结构修改

#### Cmd 结构体 (internal/cmd/cmd.go)
```go
type Cmd struct {
    // ... 现有字段 ...
    disableFlagParsing bool  // 新增：禁用标志解析
}
```

#### CmdOpts 结构体 (internal/cmd/cmdopts.go)
```go
type CmdOpts struct {
    // ... 现有字段 ...
    DisableFlagParsing bool  // 新增：禁用标志解析
}
```

#### Command 接口 (internal/types/command.go)
```go
type Command interface {
    // ... 现有方法 ...
    SetDisableFlagParsing(disable bool)  // 新增：设置禁用标志解析
    IsDisableFlagParsing() bool          // 新增：检查是否禁用标志解析
}
```

### 3.2 方法实现

#### SetDisableFlagParsing
```go
func (c *Cmd) SetDisableFlagParsing(disable bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.disableFlagParsing = disable
}
```

#### IsDisableFlagParsing
```go
func (c *Cmd) IsDisableFlagParsing() bool {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.disableFlagParsing
}
```

### 3.3 解析器修改

#### DefaultParser.ParseOnly (internal/parser/parser.go)

在解析前检查命令是否禁用了标志解析：

```go
func (p *DefaultParser) ParseOnly(cmd types.Command, args []string) error {
    // 如果禁用标志解析，直接设置参数并返回
    if cmd.IsDisableFlagParsing() {
        cmd.SetParsed(true)
        cmd.SetArgs(args)
        return nil
    }
    
    // ... 原有的标志解析逻辑 ...
}
```

#### DefaultParser.Parse (internal/parser/parser.go)

同样需要在解析前检查：

```go
func (p *DefaultParser) Parse(cmd types.Command, args []string) error {
    // 如果禁用标志解析，直接设置参数并返回
    if cmd.IsDisableFlagParsing() {
        cmd.SetParsed(true)
        cmd.SetArgs(args)
        return nil
    }
    
    // ... 原有的解析逻辑 ...
}
```

#### DefaultParser.ParseAndRoute (internal/parser/parser.go)

在路由子命令时，如果当前命令禁用标志解析，需要特殊处理：

```go
func (p *DefaultParser) ParseAndRoute(cmd types.Command, args []string) error {
    // 如果禁用标志解析，直接设置参数并尝试执行当前命令
    if cmd.IsDisableFlagParsing() {
        cmd.SetParsed(true)
        cmd.SetArgs(args)
        
        // 如果有运行函数则执行，否则返回错误
        if cmd.HasRunFunc() {
            return cmd.Run()
        }
        return fmt.Errorf("command '%s' has no run function", cmd.Name())
    }
    
    // ... 原有的解析和路由逻辑 ...
}
```

### 3.4 ApplyOpts 修改

在 `ApplyOpts` 方法中添加对 `DisableFlagParsing` 的处理：

```go
func (c *Cmd) ApplyOpts(opts *CmdOpts) error {
    // ... 现有代码 ...
    
    c.SetDisableFlagParsing(opts.DisableFlagParsing)
    
    // ... 现有代码 ...
}
```

### 3.5 Mock 实现更新

需要在 MockCommandBasic 中添加相应字段和方法：

```go
type MockCommandBasic struct {
    // ... 现有字段 ...
    disableFlagParsing bool
}

func (c *MockCommandBasic) SetDisableFlagParsing(disable bool) {
    c.disableFlagParsing = disable
}

func (c *MockCommandBasic) IsDisableFlagParsing() bool {
    return c.disableFlagParsing
}
```

## 4. 使用示例

### 4.1 使用 SetDisableFlagParsing 方法

```go
cmd := qflag.NewCmd("myapp", "", qflag.ExitOnError)

// 创建禁用标志解析的子命令
completeCmd := qflag.NewCmd("__complete", "", qflag.ExitOnError)
completeCmd.SetDisableFlagParsing(true)
completeCmd.SetRun(func(c qflag.Command) error {
    // 所有参数都在 c.Args() 中，包括 -- 开头的参数
    args := c.Args()
    // 处理补全逻辑...
    return nil
})

cmd.AddSubCmds(completeCmd)
```

### 4.2 使用 CmdOpts

```go
cmd := qflag.NewCmd("myapp", "", qflag.ExitOnError)

// 使用 CmdOpts 创建禁用标志解析的子命令
completeCmd := qflag.NewCmd("__complete", "", qflag.ExitOnError)
completeCmd.ApplyOpts(&qflag.CmdOpts{
    DisableFlagParsing: true,
    RunFunc: func(c qflag.Command) error {
        args := c.Args()
        // 处理补全逻辑...
        return nil
    },
})

cmd.AddSubCmds(completeCmd)
```

### 4.3 命令行使用

```bash
# 正常命令，标志会被解析
myapp --verbose config get --key name

# __complete 命令，所有参数都作为位置参数
myapp __complete config get -- --key
# args = ["config", "get", "--", "--key"]
```

## 5. 注意事项

1. **与内置标志的冲突**: 禁用标志解析后，`--help` 和 `--version` 等内置标志也不会被解析
2. **子命令路由**: 在 `ParseAndRoute` 中，如果父命令禁用标志解析，子命令不会被自动路由，需要自行处理
3. **环境变量**: 禁用标志解析后，环境变量绑定也不会生效
4. **互斥组和必需组**: 禁用标志解析后，这些验证也不会执行

## 6. 实现步骤

1. [ ] 在 `Cmd` 结构体添加 `disableFlagParsing` 字段
2. [ ] 添加 `SetDisableFlagParsing` 和 `IsDisableFlagParsing` 方法
3. [ ] 在 `CmdOpts` 结构体添加 `DisableFlagParsing` 字段
4. [ ] 在 `ApplyOpts` 方法中处理 `DisableFlagParsing`
5. [ ] 在 `Command` 接口添加方法声明
6. [ ] 修改 `DefaultParser.ParseOnly` 支持禁用标志解析
7. [ ] 修改 `DefaultParser.Parse` 支持禁用标志解析
8. [ ] 修改 `DefaultParser.ParseAndRoute` 支持禁用标志解析
9. [ ] 更新 `MockCommandBasic` 实现
10. [ ] 编写测试用例
11. [ ] 验证功能

## 7. 测试用例

### 7.1 基本功能测试

```go
func TestDisableFlagParsing(t *testing.T) {
    cmd := NewCmd("test", "", ExitOnError)
    cmd.SetDisableFlagParsing(true)
    
    // 添加一个标志（应该不会被解析）
    cmd.String("config", "c", "config file", "")
    
    // 解析包含标志的参数
    err := cmd.Parse([]string{"--config", "test.conf", "arg1"})
    assert.NoError(t, err)
    
    // 所有参数都应该作为位置参数
    assert.Equal(t, []string{"--config", "test.conf", "arg1"}, cmd.Args())
    
    // 标志不应该被设置
    flag, _ := cmd.GetFlag("config")
    assert.Equal(t, "", flag.Get())
}
```

### 7.2 使用 CmdOpts 测试

```go
func TestDisableFlagParsingWithOpts(t *testing.T) {
    cmd := NewCmd("test", "", ExitOnError)
    cmd.ApplyOpts(&CmdOpts{
        DisableFlagParsing: true,
    })
    
    err := cmd.Parse([]string{"--foo", "bar"})
    assert.NoError(t, err)
    assert.Equal(t, []string{"--foo", "bar"}, cmd.Args())
}
```

## 8. 与 Cobra 的对比

| 特性 | Cobra | qflag (实现后) |
|------|-------|----------------|
| 禁用标志解析 | `DisableFlagParsing` | `DisableFlagParsing` |
| 设置方式 | 结构体字段 | 方法 + CmdOpts |
| 影响范围 | 当前命令 | 当前命令 |
| 子命令继承 | 不继承 | 不继承 |
| 帮助标志 | 也被禁用 | 也被禁用 |

## 9. 结论

该功能完全可以实现，且实现复杂度较低。主要修改集中在：
- 命令结构体添加字段和方法
- 解析器在解析前检查标志
- 接口和 Mock 的更新

预计实现时间：30-60 分钟
