# qflag 动态补全方案 V2 - 基于内置标志的设计

## 1. 设计思路

基于现有内置标志架构（HelpFlag、VersionFlag、CompletionFlag），新增 `DynamicCompletionFlag` 内置标志，通过标志值控制动态补全行为。

## 2. 核心设计

### 2.1 新增内置标志类型

```go
// types/builtin.go
const (
    HelpFlag BuiltinFlagType = iota
    VersionFlag
    CompletionFlag
    DynamicCompletionFlag  // 新增：动态补全标志
)

const (
    // ... 现有常量 ...
    
    // DynamicCompletionFlagName 动态补全标志名称
    DynamicCompletionFlagName = "dynamic-completion"
)
```

### 2.2 动态补全标志的行为

```go
// 使用方式：
// 1. 启用动态补全模式（生成动态补全脚本）
//    myapp --dynamic-completion=enable
//
// 2. 执行动态补全（内部使用，shell脚本调用）
//    myapp --dynamic-completion=complete -- <args...>
//
// 3. 生成动态补全脚本
//    myapp --completion=bash --dynamic-completion=enable
```

### 2.3 架构流程

```
用户操作:
┌─────────────────────────────────────────────────────────────┐
│  myapp --completion=bash --dynamic-completion=enable        │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  BuiltinFlagManager                          │
│  1. 检测到 --dynamic-completion=enable                      │
│  2. 设置全局状态：动态补全模式启用                           │
│  3. 继续处理 --completion 标志                               │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  CompletionHandler                           │
│  1. 检查是否处于动态补全模式                                 │
│  2. 如果是，生成动态补全脚本（调用程序）                     │
│  3. 如果不是，生成静态补全脚本（原有逻辑）                   │
└─────────────────────────────────────────────────────────────┘

Shell 补全过程:
┌─────────────────────────────────────────────────────────────┐
│  用户在 shell 中按 Tab                                      │
│  myapp --dynamic-completion=complete -- deploy --env        │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  DynamicCompletionHandler                    │
│  1. 解析传入的参数                                           │
│  2. 构建 CompletionContext                                   │
│  3. 查找注册的补全函数                                       │
│  4. 返回补全结果（格式：候选1\n候选2\n:指令）                 │
└─────────────────────────────────────────────────────────────┘
```

## 3. 核心组件

### 3.1 新增处理器：DynamicCompletionHandler

```go
package builtin

// DynamicCompletionHandler 动态补全处理器
type DynamicCompletionHandler struct {
    completers map[string]DynamicCompleter  // 注册的补全函数
    enabled    bool                         // 是否启用动态补全
}

// Handle 处理动态补全标志
func (h *DynamicCompletionHandler) Handle(cmd types.Command) error {
    flag := getDynamicCompletionFlag(cmd)
    value := flag.GetStr()
    
    switch value {
    case "enable":
        // 仅设置启用状态，不退出，让其他标志处理
        h.enabled = true
        return nil
        
    case "complete":
        // 执行补全逻辑
        return h.handleCompletion(cmd)
        
    default:
        return fmt.Errorf("invalid dynamic-completion value: %s", value)
    }
}

// handleCompletion 执行补全
func (h *DynamicCompletionHandler) handleCompletion(cmd types.Command) error {
    // 获取补全参数（-- 之后的参数）
    args := getCompletionArgs(cmd)
    
    // 构建 CompletionContext
    ctx := buildCompletionContext(cmd, args)
    
    // 查找并执行补全函数
    result := h.executeCompleter(ctx)
    
    // 输出结果
    printCompletionResult(result)
    os.Exit(0)
    return nil
}
```

### 3.2 用户 API

```go
package qflag

// EnableDynamicCompletion 启用动态补全模式
// 在程序初始化时调用，注册内置的动态补全标志
func (c *Cmd) EnableDynamicCompletion() {
    c.Config().DynamicCompletion = true
}

// RegisterCompleter 为特定标志注册自定义补全函数
func (c *Cmd) RegisterCompleter(flagName string, completer DynamicCompleter) {
    // 将补全函数注册到 DynamicCompletionHandler
    handler := c.builtinManager.GetHandler(types.DynamicCompletionFlag)
    if h, ok := handler.(*DynamicCompletionHandler); ok {
        h.Register(flagName, completer)
    }
}
```

## 4. 使用示例

### 4.1 基本使用

```go
package main

import (
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    cmd := qflag.NewCmd("myapp", "示例应用", qflag.ExitOnError)
    
    // 启用动态补全（注册内置标志）
    cmd.EnableDynamicCompletion()
    
    // 添加子命令和标志
    deploy := cmd.SubCmd("deploy", "部署应用")
    deploy.String("env", "", "部署环境")
    
    // 注册自定义补全
    cmd.RegisterCompleter("env", func(ctx qflag.CompletionContext) qflag.CompletionResult {
        return qflag.CompletionResult{
            Completions: []qflag.Completion{
                {Value: "dev", Description: "开发环境"},
                {Value: "prod", Description: "生产环境"},
            },
        }
    })
    
    cmd.ParseAndRoute(os.Args[1:])
}
```

### 4.2 生成动态补全脚本

```bash
# 生成动态补全脚本（Bash）
myapp --completion=bash --dynamic-completion=enable

# 生成动态补全脚本（PowerShell）
myapp --completion=pwsh --dynamic-completion=enable
```

### 4.3 内部调用（Shell 脚本自动调用）

```bash
# Shell 补全脚本内部调用
myapp --dynamic-completion=complete -- deploy --env <tab>
```

## 5. 动态补全脚本模板

### 5.1 Bash 动态补全模板

```bash
#!/usr/bin/env bash

# 动态补全脚本
_{{.ProgramName}}_complete() {
    local cur prev words cword
    COMPREPLY=()
    
    # 获取当前输入
    _get_comp_words_by_ref -n =: cur prev words cword
    
    # 构造参数
    local args=()
    for ((i=1; i < ${#words[@]}; i++)); do
        args+=("${words[i]}")
    done
    
    # 调用程序获取补全结果
    local out
    out=$("${words[0]}" "--dynamic-completion=complete" "--" "${args[@]}" 2>/dev/null)
    
    # 解析结果（最后一行是指令）
    local directive="${out##*:}"
    out="${out%%:*}"
    
    # 处理指令
    if (( directive & 2 )); then
        compopt -o nospace 2>/dev/null
    fi
    
    # 输出补全结果
    while IFS= read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "$out" -- "$cur")
}

complete -F _{{.ProgramName}}_complete {{.ProgramName}}
```

### 5.2 PowerShell 动态补全模板

```powershell
# 动态补全脚本
$scriptBlock = {
    param($wordToComplete, $commandAst, $cursorPosition)
    
    $tokens = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
    $program = $tokens[0]
    $args = $tokens[1..($tokens.Count-1)]
    
    # 调用程序获取补全结果
    $out = & $program "--dynamic-completion=complete" "--" $args 2>$null
    
    # 解析结果
    $directive = [int]($out[-1] -replace '^:','')
    $completions = $out[0..($out.Count-2)]
    
    # 构建结果
    foreach ($comp in $completions) {
        if ($comp -like "$wordToComplete*") {
            [System.Management.Automation.CompletionResult]::new($comp, $comp, 'ParameterValue', $comp)
        }
    }
}

Register-ArgumentCompleter -CommandName "{{.ProgramName}}" -ScriptBlock $scriptBlock
```

## 6. 与现有架构的集成

### 6.1 修改点

1. **types/builtin.go**: 新增 `DynamicCompletionFlag` 类型和常量
2. **builtin/manager.go**: 注册新的处理器
3. **builtin/completion_handler.go**: 检查动态补全模式，生成对应脚本
4. **新增 builtin/dynamic_completion_handler.go**: 实现动态补全逻辑

### 6.2 关键代码

```go
// builtin/completion_handler.go 修改
func (h *CompletionHandler) Handle(cmd types.Command) error {
    shellType := getShellTypeFromArgs(cmd)
    
    // 检查是否启用动态补全
    if isDynamicCompletionEnabled(cmd) {
        // 生成动态补全脚本
        completion.GenDynamicAndPrint(cmd, shellType)
    } else {
        // 生成静态补全脚本（原有逻辑）
        completion.GenAndPrint(cmd, shellType)
    }
    
    os.Exit(0)
    return nil
}
```

## 7. 优势

1. **符合现有架构**: 完全基于内置标志机制，无需修改解析逻辑
2. **向后兼容**: 默认行为不变，用户需要显式启用
3. **灵活切换**: 通过标志值控制，可以在静态和动态模式间切换
4. **易于扩展**: 可以添加更多动态补全相关的标志值

## 8. 总结

这个方案通过新增 `DynamicCompletionFlag` 内置标志，将动态补全功能集成到现有的内置标志处理框架中：

- `--dynamic-completion=enable`: 启用动态模式（配合 `--completion` 生成动态脚本）
- `--dynamic-completion=complete`: 执行补全（Shell 脚本内部调用）

这样避免了新增内置子命令的复杂性，同时保持了与现有架构的一致性。
