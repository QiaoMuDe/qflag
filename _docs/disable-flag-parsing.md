# 禁用标志解析功能 (DisableFlagParsing)

## 概述

`DisableFlagParsing` 是 qflag 提供的一个高级功能，允许命令完全跳过标志解析阶段，将所有命令行参数（包括 `--flag` 和 `-f` 形式）都作为**位置参数**处理。

这个功能类似于 Cobra 的 `DisableFlagParsing` 选项，适用于需要透传参数给子进程或外部命令的场景。

## 原理

### 核心机制

当命令的 `DisableFlagParsing` 设置为 `true` 时，解析器会：

1. **跳过标志解析**：不创建 `flag.FlagSet`，不调用 `flagSet.Parse()`
2. **保留原始参数**：将所有参数原样保存为位置参数
3. **不影响子命令路由**：子命令识别和路由逻辑正常工作
4. **支持嵌套控制**：每个命令可以独立设置是否禁用标志解析

### 代码实现

在 `ParseOnly` 方法中，首先检查是否禁用标志解析：

```go
func (p *DefaultParser) ParseOnly(cmd types.Command, args []string) error {
    // 如果禁用标志解析，直接设置参数并返回
    if cmd.IsDisableFlagParsing() {
        cmd.SetParsed(true)
        cmd.SetArgs(args)  // 原样保留所有参数
        return nil
    }
    // ... 正常标志解析逻辑
}
```

`Parse` 和 `ParseAndRoute` 方法依赖 `ParseOnly`，因此自动继承此行为：

```go
func (p *DefaultParser) Parse(cmd types.Command, args []string) error {
    // 先解析参数（ParseOnly 会处理禁用标志解析的情况）
    if err := p.ParseOnly(cmd, args); err != nil {
        return err
    }
    // ... 子命令路由逻辑（不受影响）
}
```

## 使用场景

### 场景一：包装外部命令

实现类似 `kubectl exec` 或 `docker run` 的功能，需要将参数透传给子进程：

```go
// myapp exec podname -- ls -la
// 其中 --namespace 是 myapp 的标志，-- ls -la 要透传给 exec

func main() {
    root := qflag.NewCmd("myapp", "", qflag.ExitOnError)
    
    // 添加全局标志
    root.String("namespace", "n", "命名空间", "default")
    
    // exec 子命令禁用标志解析
    execCmd := qflag.NewCmd("exec", "", qflag.ExitOnError)
    execCmd.SetDisableFlagParsing(true)
    execCmd.SetRun(func(cmd types.Command) error {
        args := cmd.Args()
        // args = ["podname", "--", "ls", "-la"]
        // 直接透传给 kubectl exec
        return runKubectlExec(args)
    })
    
    root.AddSubCmds(execCmd)
    root.Execute()
}
```

### 场景二：Shell 脚本包装器

命令只是作为其他程序的包装，不需要解析任何标志：

```go
// myapp ssh user@host -- some-command
// 所有参数都透传给 ssh

sshCmd := qflag.NewCmd("ssh", "", qflag.ExitOnError)
sshCmd.SetDisableFlagParsing(true)
sshCmd.SetDesc("通过 SSH 执行远程命令")
sshCmd.SetRun(func(cmd types.Command) error {
    args := cmd.Args()
    // args = ["user@host", "--", "some-command"]
    return runSSH(args)
})
```

### 场景三：嵌套命令中的差异化控制

父命令解析标志，子命令禁用解析：

```go
// parent --verbose child --flag value
// --verbose 被 parent 解析
// --flag value 作为位置参数传给 child

parent := qflag.NewCmd("parent", "", qflag.ExitOnError)
verbose := parent.Bool("verbose", "v", "详细输出", false)

child := qflag.NewCmd("child", "", qflag.ExitOnError)
child.SetDisableFlagParsing(true)  // 子命令禁用解析
child.SetRun(func(cmd types.Command) error {
    args := cmd.Args()
    // args = ["--flag", "value"]
    // 可以手动处理这些参数
    return nil
})

parent.AddSubCmds(child)
```

### 场景四：通过 CmdOpts 批量配置

```go
cmd := qflag.NewCmd("wrapper", "", qflag.ExitOnError)
opts := &qflag.CmdOpts{
    Desc:               "命令包装器",
    DisableFlagParsing: true,
    RunFunc: func(c types.Command) error {
        args := c.Args()
        // 处理原始参数
        return nil
    },
}
cmd.ApplyOpts(opts)
```

## API 参考

### 方法

| 方法 | 说明 |
|------|------|
| `SetDisableFlagParsing(disable bool)` | 设置是否禁用标志解析 |
| `IsDisableFlagParsing() bool` | 检查是否禁用标志解析 |

### CmdOpts 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `DisableFlagParsing` | `bool` | 是否禁用标志解析，默认 `false` |

## 注意事项

1. **子命令路由不受影响**：禁用标志解析只影响当前命令的标志解析，不影响子命令的识别和路由

2. **嵌套独立控制**：每个命令可以独立设置 `DisableFlagParsing`，父命令禁用不影响子命令

3. **内置标志也被禁用**：当禁用标志解析时，`--help` 和 `--version` 等内置标志也不会被特殊处理

4. **环境变量绑定**：禁用标志解析时，环境变量绑定也会被跳过

5. **互斥组和必需组验证**：禁用标志解析时，这些验证也会被跳过

## 示例对比

### 不禁用标志解析（默认行为）

```bash
$ myapp --flag value arg1
# --flag 被解析为标志
# arg1 作为位置参数
```

### 禁用标志解析

```bash
$ myapp --flag value arg1
# --flag 作为位置参数
# value 作为位置参数
# arg1 作为位置参数
# 所有参数 = ["--flag", "value", "arg1"]
```

## 测试覆盖

测试文件：`internal/cmd/disable_flag_parsing_test.go`

测试用例包括：
- 基本禁用功能（各种参数形式）
- 子命令路由验证
- 嵌套子命令控制
- ParseAndRoute 执行逻辑
- 不禁用时的正常行为
- 三种解析方法的一致性
- 通过 CmdOpts 设置
