# 智能纠错提示功能设计方案

## 背景

当用户输入错误的子命令或标志时，当前 qflag 只会返回简单的错误信息。参考 Git 的行为，可以在报错的同时推荐最相似的命令，提升用户体验。

**Git 示例：**
```bash
$ git pu
git: 'pu' is not a git command. See 'git --help'.

The most similar commands are
        pull
        push
        p4
```

## 设计目标

1. **子命令纠错**：输入错误的子命令时，推荐相似的子命令
2. **标志纠错**：输入错误的标志时，推荐相似的标志
3. **默认启用**：无需配置，直接替换现有错误格式
4. **性能优化**：使用高效的模糊匹配算法

## 技术方案

### 1. 依赖库

使用 `gitee.com/MM-Q/go-kit/fuzzy` 库进行模糊匹配：
- 内部维护的库，与 qflag 项目同源
- 支持前缀匹配和模糊搜索
- 性能优秀

```bash
go get gitee.com/MM-Q/go-kit/fuzzy
```

### 2. 错误类型扩展

```go
// UnknownSubcommandError 未知子命令错误
type UnknownSubcommandError struct {
    Command     string   // 当前命令名
    Input       string   // 用户输入的错误子命令
    Suggestions []string // 相似子命令建议
}

func (e *UnknownSubcommandError) Error() string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%s: '%s' is not a valid command. See '%s --help'.\n",
        e.Command, e.Input, e.Command))

    if len(e.Suggestions) > 0 {
        sb.WriteString("\nThe most similar commands are\n")
        for _, sug := range e.Suggestions {
            sb.WriteString(fmt.Sprintf("        %s\n", sug))
        }
    }

    return sb.String()
}

// UnknownFlagError 未知标志错误
type UnknownFlagError struct {
    Command     string   // 当前命令名
    Input       string   // 用户输入的错误标志
    Suggestions []string // 相似标志建议
}

func (e *UnknownFlagError) Error() string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%s: unknown flag: '%s'\n",
        e.Command, e.Input))

    if len(e.Suggestions) > 0 {
        sb.WriteString("\nThe most similar flags are\n")
        for _, sug := range e.Suggestions {
            sb.WriteString(fmt.Sprintf("        %s\n", sug))
        }
    }

    return sb.String()
}
```

### 3. 模糊匹配实现

使用 go-kit/fuzzy 库的 `CompletePrefix` 函数，优先前缀匹配，再模糊匹配：

```go
// FindSimilarSubcommands 查找相似的子命令
func FindSimilarSubcommands(input string, subcommands []string, maxResults int) []string {
    if len(subcommands) == 0 {
        return nil
    }

    // 使用 go-kit/fuzzy 进行前缀优先的模糊匹配
    matches := fuzzy.CompletePrefix(input, subcommands)

    // 限制返回数量
    if len(matches) > maxResults {
        matches = matches[:maxResults]
    }

    return matches
}

// FindSimilarFlags 查找相似的标志
func FindSimilarFlags(input string, flags []*Flag, maxResults int) []string {
    if len(flags) == 0 {
        return nil
    }

    // 构建标志名称列表（包含长短名称）
    flagNames := make([]string, 0, len(flags)*2)
    for _, f := range flags {
        if f.LongName != "" {
            flagNames = append(flagNames, "--"+f.LongName)
        }
        if f.ShortName != "" {
            flagNames = append(flagNames, "-"+f.ShortName)
        }
    }

    // 移除输入中的横杠前缀进行匹配
    cleanInput := strings.TrimLeft(input, "-")

    // 使用 go-kit/fuzzy 进行前缀优先的模糊匹配
    matches := fuzzy.CompletePrefix(cleanInput, flagNames)

    // 限制返回数量
    if len(matches) > maxResults {
        matches = matches[:maxResults]
    }

    return matches
}
```

### 4. 智能纠错查找器封装

将建议查找逻辑封装到独立的查找器中，保持解析方法简洁：

```go
// internal/parser/suggestion.go

// SuggestionFinder 智能纠错查找器
type SuggestionFinder struct {
    maxSuggestions int
}

// NewSuggestionFinder 创建查找器
func NewSuggestionFinder(maxSuggestions int) *SuggestionFinder {
    return &SuggestionFinder{maxSuggestions: maxSuggestions}
}

// FindForSubcommand 查找子命令建议
func (f *SuggestionFinder) FindForSubcommand(input string, cmd types.Command) []string {
    subCmds := cmd.SubCmds()
    if len(subCmds) == 0 {
        return nil
    }
    
    names := make([]string, len(subCmds))
    for i, sc := range subCmds {
        names[i] = sc.Name()
    }
    
    return f.findSimilar(input, names)
}

// FindForFlag 查找标志建议
func (f *SuggestionFinder) FindForFlag(input string, cmd types.Command) []string {
    flags := cmd.FlagRegistry().List()
    if len(flags) == 0 {
        return nil
    }
    
    names := make([]string, 0, len(flags)*2)
    for _, fl := range flags {
        if fl.LongName != "" {
            names = append(names, "--"+fl.LongName)
        }
        if fl.ShortName != "" {
            names = append(names, "-"+fl.ShortName)
        }
    }
    
    cleanInput := strings.TrimLeft(input, "-")
    return f.findSimilar(cleanInput, names)
}

// 内部模糊匹配
func (f *SuggestionFinder) findSimilar(input string, candidates []string) []string {
    matches := fuzzy.CompletePrefix(input, candidates)
    if len(matches) > f.maxSuggestions {
        matches = matches[:f.maxSuggestions]
    }
    result := make([]string, len(matches))
    for i, m := range matches {
        result[i] = m.Str
    }
    return result
}

// newUnknownSubcommandError 创建未知子命令错误（带建议）
func newUnknownSubcommandError(cmd types.Command, input string) error {
    finder := NewSuggestionFinder(3)
    suggestions := finder.FindForSubcommand(input, cmd)
    
    return &UnknownSubcommandError{
        Command:     cmd.Name(),
        Input:       input,
        Suggestions: suggestions,
    }
}

// newUnknownFlagError 创建未知标志错误（带建议）
func newUnknownFlagError(cmd types.Command, input string) error {
    finder := NewSuggestionFinder(3)
    suggestions := finder.FindForFlag(input, cmd)
    
    return &UnknownFlagError{
        Command:     cmd.Name(),
        Input:       input,
        Suggestions: suggestions,
    }
}
```

### 5. 集成点分析

#### 5.1 子命令纠错 - `parser.go`

**判断逻辑：**
- 如果命令**有子命令**，第一个参数**必须**是子命令
- 如果不是已知的子命令，就认为是输错了，进行纠错
- 如果命令**没有子命令**，第一个参数是普通位置参数，不纠错

```go
// Parse 方法修改（约第 158-172 行）
func (p *DefaultParser) Parse(cmd types.Command, args []string) error {
    // ... 前面的代码不变 ...

    // 检查剩余参数是否为子命令
    cmdRegistry := cmd.CmdRegistry()
    remainingArgs := cmd.Args()

    // 如果有剩余参数, 检查是否为子命令
    if len(remainingArgs) > 0 {
        firstArg := remainingArgs[0]

        // 检查是否为子命令
        if subCmd, ok := cmdRegistry.Get(firstArg); ok {
            return subCmd.Parse(remainingArgs[1:])
        }

        // 有子命令但没匹配上 → 一定是输错了，纠错
        if len(cmd.SubCmds()) > 0 {
            return newUnknownSubcommandError(cmd, firstArg)
        }
        
        // 没有子命令 → 是普通参数，正常处理
    }

    return nil
}

// ParseAndRoute 方法修改（约第 197-222 行）
func (p *DefaultParser) ParseAndRoute(cmd types.Command, args []string) error {
    // ... 前面的代码不变 ...

    // 如果是子命令, 递归解析并执行子命令
    if len(remainingArgs) > 0 {
        firstArg := remainingArgs[0]

        // 检查是否为子命令
        if subCmd, ok := cmdRegistry.Get(firstArg); ok {
            return subCmd.ParseAndRoute(remainingArgs[1:])
        }

        // 有子命令但没匹配上 → 一定是输错了，纠错
        if len(cmd.SubCmds()) > 0 {
            return newUnknownSubcommandError(cmd, firstArg)
        }
        
        // 没有子命令 → 是普通参数，正常处理
    }

    // ... 后面的代码不变 ...
}
```

#### 5.2 标志纠错 - 预解析检查方案

**判断逻辑：**
- 以 `-` 或 `--` 开头的参数**一定是标志**
- 如果不在已注册标志列表中，就是错误的标志，进行纠错
- 遇到 `--` 停止扫描，后面的都视为位置参数（即使以 `-` 开头）

```go
// ParseOnly 方法修改（约第 105 行附近）
func (p *DefaultParser) ParseOnly(cmd types.Command, args []string) error {
    // ... 前面的代码不变 ...

    // 预检查：扫描未知标志（调用封装函数）
    if err := checkUnknownFlags(cmd, args); err != nil {
        return err
    }

    // 调用标准库解析
    if err := p.flagSet.Parse(args); err != nil {
        return err
    }

    // ... 后面的代码不变 ...
}

// checkUnknownFlags 预扫描参数，检查未知标志（独立函数）
func checkUnknownFlags(cmd types.Command, args []string) error {
    // 获取所有已注册的标志名（长短名称都包括）
    registeredFlags := make(map[string]bool)
    for _, f := range cmd.FlagRegistry().List() {
        if f.LongName != "" {
            registeredFlags["--"+f.LongName] = true
            registeredFlags["-"+f.LongName] = true  // 支持单横杠长名称
        }
        if f.ShortName != "" {
            registeredFlags["-"+f.ShortName] = true
        }
    }

    // 扫描参数
    for i := 0; i < len(args); i++ {
        arg := args[i]

        // 【关键】遇到 -- 停止扫描，后面的都视为位置参数
        if arg == "--" {
            break
        }

        // 不是标志格式，跳过
        if !strings.HasPrefix(arg, "-") {
            continue
        }

        // 处理 --flag=value 格式
        flagName := arg
        if idx := strings.Index(arg, "="); idx != -1 {
            flagName = arg[:idx]
        }

        // 不是已注册的标志 → 纠错
        if !registeredFlags[flagName] {
            return newUnknownFlagError(cmd, flagName)
        }
    }

    return nil
}
```

**方案优势：**
- 不依赖标准库错误字符串格式，向后兼容
- 完全控制错误信息，易于国际化
- 实现简洁，不拷贝源码或嵌入类型
- 性能好，只扫描一次参数列表
- 解析方法中只需一行调用，逻辑清晰

**判断逻辑总结：**

| 场景 | 判断条件 | 处理方式 |
|------|---------|---------|
| **子命令纠错** | 命令有子命令，但第一个参数不匹配任何子命令 | 报错并推荐相似子命令 |
| **子命令正常** | 命令没有子命令 | 第一个参数视为普通位置参数，不纠错 |
| **标志纠错** | 以 `-` 开头，且不在已注册标志列表中 | 报错并推荐相似标志 |
| **标志终止** | 遇到 `--` | 停止扫描，后面所有参数视为位置参数 |

### 6. 使用示例

```go
package main

import (
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建根命令
    root := qflag.NewCmd("myapp", "m", qflag.ExitOnError)

    // 添加子命令
    root.AddSubCmds(
        qflag.NewCmd("config", "c", qflag.ExitOnError),
        qflag.NewCmd("build", "b", qflag.ExitOnError),
        qflag.NewCmd("deploy", "d", qflag.ExitOnError),
    )

    // 添加标志
    root.String("output", "o", "输出文件", "")
    root.Bool("verbose", "v", "详细模式", false)

    // 应用配置
    root.ApplyOpts(&qflag.CmdOpts{
        Desc:       "My Application",
        UseChinese: true,
    })

    // 解析并执行
    if err := root.ParseAndRun(); err != nil {
        // 错误信息已经包含智能纠错建议
        println(err.Error())
    }
}
```

### 7. 预期输出效果

**子命令纠错：**
```bash
$ myapp cnfig
myapp: 'cnfig' is not a valid command. See 'myapp --help'.

The most similar commands are
        config

$ myapp bld
myapp: 'bld' is not a valid command. See 'myapp --help'.

The most similar commands are
        build
```

**标志纠错：**
```bash
$ myapp --verboose
myapp: unknown flag: '--verboose'

The most similar flags are
        --verbose
        -v

$ myapp --outpu
myapp: unknown flag: '--outpu'

The most similar flags are
        --output
        -o
```

## 实施步骤

### 第一阶段：基础实现
1. 添加 `gitee.com/MM-Q/go-kit/fuzzy` 依赖
2. 在 `internal/types/error.go` 创建 `UnknownSubcommandError` 和 `UnknownFlagError` 错误类型
3. 在 `internal/parser` 包创建 `suggestion.go`，实现 `SuggestionFinder` 和封装函数

### 第二阶段：子命令纠错集成
1. 修改 `parser.go` 的 `Parse` 方法，调用 `newUnknownSubcommandError`
2. 修改 `parser.go` 的 `ParseAndRoute` 方法，调用 `newUnknownSubcommandError`
3. 添加单元测试

### 第三阶段：标志纠错集成
1. 实现 `checkUnknownFlags` 函数预扫描未知标志
2. 修改 `ParseOnly` 方法，在 `flagSet.Parse()` 之前调用 `checkUnknownFlags`
3. 添加单元测试

### 第四阶段：文档更新
1. 更新 README.md 添加智能纠错特性说明
2. 更新 API 文档
3. 添加使用示例

## 注意事项

1. **性能考虑**：模糊匹配在子命令/标志数量较多时可能影响性能，限制最大建议数量为 3
2. **国际化**：错误信息需要支持中英文切换（根据 `UseChinese` 配置）
3. **向后兼容**：直接替换错误格式，行为更加友好
4. **测试覆盖**：需要充分测试各种边缘情况（空输入、特殊字符等）
5. **预扫描逻辑**：`checkUnknownFlags` 需要正确处理 `--flag=value` 格式、`--` 终止符和子命令边界

## 相关文件

- `internal/types/error.go` - 错误类型定义（需要新建或扩展）
- `internal/parser/suggestion.go` - 智能纠错查找器（新建）
- `internal/parser/parser.go` - 解析器修改（Parse、ParseAndRoute、ParseOnly 方法）
- `go.mod` - 添加 fuzzy 依赖
