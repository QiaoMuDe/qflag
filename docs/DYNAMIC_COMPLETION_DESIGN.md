# qflag 动态补全方案设计文档

## 1. 设计目标

在保留现有静态补全方案的基础上，新增类似 Cobra 的动态补全模式，让用户可以选择：
- **静态模式**（默认）：性能好，无需程序支持
- **动态模式**（可选）：完全灵活，可补全动态数据

## 2. 核心架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        用户程序 (使用 qflag)                      │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  cmd := qflag.NewCmd("myapp", "", qflag.ExitOnError)    │   │
│  │  cmd.AddSubCmds(sub1, sub2)                             │   │
│  │                                                         │   │
│  │  // 启用动态补全                                        │   │
│  │  cmd.EnableDynamicCompletion()                          │   │
│  │                                                         │   │
│  │  // 或注册自定义补全函数                                │   │
│  │  cmd.RegisterCompleter("user", func(ctx CompletionContext) []string {
│  │      return fetchUsersFromDB()  // 从数据库获取用户列表
│  │  })                                                     │   │
│  │                                                         │   │
│  │  cmd.ParseAndRoute(os.Args[1:])                         │   │
│  └─────────────────────────────────────────────────────────┘   │
└──────────────────────────────────┬────────────────────────────┘
                                   │
                                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                      qflag 内部实现                              │
│  ┌─────────────────────┐    ┌───────────────────────────────┐  │
│  │   静态补全生成器     │    │      动态补全处理器            │  │
│  │  (现有功能)         │    │  (新增功能)                    │  │
│  │                     │    │                               │  │
│  │  • 命令树遍历       │    │  • __complete 子命令处理      │  │
│  │  • 标志收集         │    │  • CompletionContext 构建     │  │
│  │  • 模板渲染         │    │  • 自定义补全函数调用          │  │
│  │  • 脚本输出         │    │  • 指令系统 (Directive)       │  │
│  └─────────────────────┘    └───────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                   │
                                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                      生成的补全脚本                              │
│  ┌─────────────────────────┐    ┌─────────────────────────────┐ │
│  │    静态补全脚本          │    │     动态补全脚本             │ │
│  │  (bash/pwsh.tmpl)       │    │  (bash_dynamic.tmpl)        │ │
│  │                         │    │                             │ │
│  │  所有逻辑在脚本中        │    │  脚本只是代理:               │ │
│  │  • 命令树数组           │    │  myapp __complete arg1 arg2 │ │
│  │  • 标志数组             │    │  • 调用程序获取结果         │ │
│  │  • 本地模糊匹配         │    │  • 解析指令并输出           │ │
│  └─────────────────────────┘    └─────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## 3. 核心组件设计

### 3.1 指令系统 (Directive)

```go
// CompletionDirective 补全指令，控制补全行为
type CompletionDirective int

const (
    // CompDirectiveError 错误，无补全
    CompDirectiveError CompletionDirective = 1 << iota
    
    // CompDirectiveNoSpace 补全后不加空格
    CompDirectiveNoSpace
    
    // CompDirectiveNoFileComp 禁用文件补全
    CompDirectiveNoFileComp
    
    // CompDirectiveFilterFileExt 过滤文件扩展名
    CompDirectiveFilterFileExt
    
    // CompDirectiveFilterDirs 只补全目录
    CompDirectiveFilterDirs
    
    // CompDirectiveKeepOrder 保持原始顺序，不排序
    CompDirectiveKeepOrder
)
```

### 3.2 补全上下文 (CompletionContext)

```go
// CompletionContext 包含补全请求的完整上下文信息
type CompletionContext struct {
    // 当前命令
    Command Command
    
    // 当前输入的词
    WordToComplete string
    
    // 前一个词
    PreviousWord string
    
    // 所有已解析的参数
    Args []string
    
    // 当前命令路径，如 "/cmd/subcmd"
    CommandPath string
    
    // 当前是否在补全标志值
    CompletingFlagValue bool
    
    // 当前正在补全的标志名（如果在补全标志值）
    CurrentFlag string
}
```

### 3.3 补全结果 (CompletionResult)

```go
// Completion 单个补全项
type Completion struct {
    Value       string // 补全值
    Description string // 描述（可选）
}

// CompletionResult 补全结果
type CompletionResult struct {
    Completions []Completion        // 补全列表
    Directive   CompletionDirective // 指令
}
```

### 3.4 动态补全接口

```go
// DynamicCompleter 动态补全函数类型
type DynamicCompleter func(ctx CompletionContext) CompletionResult

// Command 接口扩展

// EnableDynamicCompletion 启用动态补全模式
func (c *Cmd) EnableDynamicCompletion()

// RegisterCompleter 为特定标志或路径注册自定义补全函数
func (c *Cmd) RegisterCompleter(flagOrPath string, completer DynamicCompleter)

// SetDefaultCompleter 设置默认补全函数（当没有特定匹配时调用）
func (c *Cmd) SetDefaultCompleter(completer DynamicCompleter)
```

## 4. 使用示例

### 4.1 基本动态补全

```go
package main

import (
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    cmd := qflag.NewCmd("myapp", "示例应用", qflag.ExitOnError)
    
    // 启用动态补全
    cmd.EnableDynamicCompletion()
    
    // 添加子命令
    deploy := cmd.SubCmd("deploy", "部署应用")
    
    // 添加标志
    deploy.String("env", "", "部署环境")
    deploy.String("version", "", "版本号")
    
    // 为 --env 标志注册自定义补全
    deploy.RegisterCompleter("env", func(ctx qflag.CompletionContext) qflag.CompletionResult {
        return qflag.CompletionResult{
            Completions: []qflag.Completion{
                {Value: "dev", Description: "开发环境"},
                {Value: "staging", Description: "预发布环境"},
                {Value: "prod", Description: "生产环境"},
            },
            Directive: qflag.CompDirectiveNoSpace,
        }
    })
    
    cmd.ParseAndRoute(os.Args[1:])
}
```

### 4.2 从数据库动态补全

```go
// 为 --user-id 标志注册数据库查询补全
cmd.RegisterCompleter("user-id", func(ctx qflag.CompletionContext) qflag.CompletionResult {
    // 从数据库获取用户列表
    users := db.QueryUsers(ctx.WordToComplete) // 模糊查询
    
    var completions []qflag.Completion
    for _, user := range users {
        completions = append(completions, qflag.Completion{
            Value:       user.ID,
            Description: user.Name,
        })
    }
    
    return qflag.CompletionResult{
        Completions: completions,
        Directive:   qflag.CompDirectiveNoFileComp, // 禁用文件补全
    }
})
```

### 4.3 根据上下文动态补全

```go
// 根据已选择的 region 补全可用 zone
cmd.RegisterCompleter("zone", func(ctx qflag.CompletionContext) qflag.CompletionResult {
    // 获取已选择的 region
    region := ctx.Command.GetFlagValue("region")
    
    zones := getZonesForRegion(region)
    
    var completions []qflag.Completion
    for _, zone := range zones {
        completions = append(completions, qflag.Completion{
            Value:       zone.ID,
            Description: zone.Name,
        })
    }
    
    return qflag.CompletionResult{
        Completions: completions,
    }
})
```

## 5. 动态补全脚本模板

### 5.1 Bash 动态补全模板

```bash
#!/usr/bin/env bash

# 动态补全脚本 - 调用程序获取补全结果

_{{.ProgramName}}_complete() {
    local cur prev words cword
    COMPREPLY=()
    
    # 获取当前输入
    if declare -F _get_comp_words_by_ref >/dev/null 2>&1; then
        _get_comp_words_by_ref -n =: cur prev words cword
    else
        words=("${COMP_WORDS[@]}")
        cword=$COMP_CWORD
        cur="${words[cword]}"
        prev="${words[cword-1]}"
    fi
    
    # 构造调用命令
    # 传递所有参数给程序的 __complete 命令
    local args=()
    for ((i=1; i < ${#words[@]}; i++)); do
        args+=("${words[i]}")
    done
    
    # 如果最后一个参数已完成（有空格），添加空字符串
    if [[ -z "$cur" && "${words[cword-1]}" != *= ]]; then
        args+=("")
    fi
    
    # 调用程序获取补全结果
    local out directive
    out=$("${words[0]}" "{{.CompleteCmd}}" "${args[@]}" 2>/dev/null)
    
    # 解析指令（最后一行）
    directive="${out##*:}"
    out="${out%%:*}"
    
    if [[ "$directive" == "$out" ]]; then
        directive=0
    fi
    
    # 处理指令
    local noSpace=0 noFileComp=0
    if (( directive & 1 )); then  # Error
        return
    fi
    if (( directive & 2 )); then  # NoSpace
        noSpace=1
        if [[ $(type -t compopt) = "builtin" ]]; then
            compopt -o nospace
        fi
    fi
    if (( directive & 4 )); then  # NoFileComp
        noFileComp=1
        if [[ $(type -t compopt) = "builtin" ]]; then
            compopt +o default
        fi
    fi
    
    # 输出补全结果
    if [[ -n "$out" ]]; then
        while IFS= read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "$out" -- "$cur")
    fi
    
    # 如果没有补全且允许文件补全
    if [[ ${#COMPREPLY[@]} -eq 0 && $noFileComp -eq 0 ]]; then
        COMPREPLY=($(compgen -f -d -- "$cur"))
    fi
}

complete -F _{{.ProgramName}}_complete {{.ProgramName}}
```

### 5.2 PowerShell 动态补全模板

```powershell
# 动态补全脚本 - 调用程序获取补全结果

$scriptBlock = {
    param($wordToComplete, $commandAst, $cursorPosition)
    
    # 获取命令元素
    $tokens = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
    
    # 构造调用命令
    $program = $tokens[0]
    $args = $tokens[1..($tokens.Count-1)]
    
    # 如果最后一个参数已完成，添加空字符串
    if ($wordToComplete -eq "" -and $tokens[-1] -notlike "*=*") {
        $args += ""
    }
    
    # 调用程序获取补全结果
    $out = & $program "{{.CompleteCmd}}" $args 2>$null
    
    # 解析指令（最后一行）
    $directive = [int]($out[-1] -replace '^:','')
    $completions = $out[0..($out.Count-2)]
    
    # 处理指令
    $noSpace = ($directive -band 2) -ne 0
    $noFileComp = ($directive -band 4) -ne 0
    
    # 构建补全结果
    $results = @()
    foreach ($comp in $completions) {
        if ($comp -like "$wordToComplete*") {
            $result = New-Object System.Management.Automation.CompletionResult(
                $comp,
                $comp,
                'ParameterValue',
                $comp
            )
            $results += $result
        }
    }
    
    return $results
}

Register-ArgumentCompleter -CommandName "{{.ProgramName}}" -ScriptBlock $scriptBlock
```

## 6. 内部实现细节

### 6.1 __complete 命令处理流程

```
用户输入: myapp __complete deploy --env dev --region

1. 解析参数
   - 识别出这是补全请求
   - 提取命令路径: ["deploy"]
   - 提取已解析的标志: {"env": "dev"}
   - 当前补全位置: "--region" 的值

2. 构建 CompletionContext
   - Command: deploy 命令对象
   - WordToComplete: "" (空，等待输入)
   - PreviousWord: "--region"
   - Args: ["deploy", "--env", "dev", "--region"]
   - CommandPath: "/deploy"
   - CompletingFlagValue: true
   - CurrentFlag: "region"

3. 查找补全函数
   - 检查是否有为 "region" 标志注册的补全函数
   - 如果有，调用它
   - 如果没有，使用默认补全（枚举值、文件路径等）

4. 返回结果
   - 格式: "候选1\n候选2\n候选3\n:6"
   - 最后一行是指令（6 = NoSpace | NoFileComp）
```

### 6.2 默认补全行为

当没有注册自定义补全函数时，qflag 提供以下默认补全：

1. **子命令补全**: 列出当前命令的所有子命令
2. **标志补全**: 列出当前命令的所有标志
3. **枚举值补全**: 如果标志是枚举类型，列出所有枚举值
4. **文件路径补全**: 对于非枚举类型的标志值，提供文件路径补全

### 6.3 性能优化

1. **延迟初始化**: 只在检测到 `__complete` 命令时才初始化补全相关数据结构
2. **缓存**: 对于不频繁变化的数据（如命令树），可以缓存解析结果
3. **超时控制**: 自定义补全函数可以设置超时，防止长时间阻塞

## 7. 与静态方案的对比

| 特性 | 静态方案 | 动态方案 |
|------|---------|---------|
| **性能** | 快（纯脚本） | 稍慢（需启动程序） |
| **灵活性** | 低（生成后固定） | 高（运行时计算） |
| **动态数据** | 不支持 | 支持（数据库、API等） |
| **上下文感知** | 有限 | 完全 |
| **部署复杂度** | 低 | 中（需确保程序可执行） |
| **适用场景** | 简单工具 | 复杂工具、动态数据 |

## 8. 迁移路径

现有用户可以通过以下方式迁移到动态补全：

```go
// 原有代码
cmd := qflag.NewCmd("myapp", "", qflag.ExitOnError)
// ... 配置命令 ...
cmd.ParseAndRoute(os.Args[1:])

// 新增一行启用动态补全
cmd.EnableDynamicCompletion()

// 可选：注册自定义补全
cmd.RegisterCompleter("user", userCompleter)
```

生成补全脚本时，用户可以选择模式：

```bash
# 生成静态补全脚本（默认）
myapp completion bash

# 生成动态补全脚本
myapp completion bash --dynamic
```

## 9. 总结

动态补全方案为 qflag 提供了与 Cobra 类似的灵活性，同时保持了与现有静态方案的兼容性。用户可以根据需求选择：

- **静态方案**: 简单、快速、无需额外配置
- **动态方案**: 灵活、强大、支持动态数据

这种设计让 qflag 既能满足简单工具的需求，也能支持复杂的企业级应用。
