# qflag 动态补全方案 V3 - 复用现有补全标志的设计

## 1. 设计思路

复用现有的 `--completion` 标志，新增 `--comp-mode` 内置标志来控制补全模式：
- `--comp-mode=static` (默认): 生成静态补全脚本
- `--comp-mode=dynamic`: 启用动态补全机制

当启用动态模式且有位置参数时，执行补全逻辑而不是生成脚本。

## 2. 使用方式

```bash
# 1. 生成静态补全脚本（默认行为，不变）
myapp --completion=bash

# 2. 生成动态补全脚本
myapp --completion=bash --comp-mode=dynamic

# 3. 执行动态补全（Shell脚本内部调用）
myapp --completion=bash --comp-mode=dynamic -- deploy --env <tab>
```

## 3. 核心设计

### 3.1 新增内置标志

```go
// types/builtin.go
const (
    HelpFlag BuiltinFlagType = iota
    VersionFlag
    CompletionFlag
    CompModeFlag  // 新增：补全模式标志
)

const (
    // ... 现有常量 ...
    
    // CompModeFlagName 补全模式标志名称
    CompModeFlagName = "comp-mode"
    
    // CompModeStatic 静态补全模式
    CompModeStatic = "static"
    
    // CompModeDynamic 动态补全模式
    CompModeDynamic = "dynamic"
)
```

### 3.2 修改 CompletionHandler

```go
// builtin/completion_handler.go

func (h *CompletionHandler) Handle(cmd types.Command) error {
    shellType := getShellTypeFromArgs(cmd)
    mode := getCompModeFromArgs(cmd) // 获取模式：static 或 dynamic
    
    // 获取位置参数（-- 之后的参数）
    posArgs := getPositionalArgs(cmd)
    
    if mode == CompModeDynamic && len(posArgs) > 0 {
        // 动态补全模式 + 有位置参数 = 执行补全
        return h.handleDynamicCompletion(cmd, shellType, posArgs)
    }
    
    // 生成补全脚本（静态或动态）
    if mode == CompModeDynamic {
        completion.GenDynamicAndPrint(cmd, shellType)
    } else {
        completion.GenAndPrint(cmd, shellType) // 原有静态逻辑
    }
    
    os.Exit(0)
    return nil
}

// handleDynamicCompletion 执行动态补全
func (h *CompletionHandler) handleDynamicCompletion(cmd types.Command, shellType string, args []string) error {
    // 构建 CompletionContext
    ctx := buildCompletionContext(cmd, args)
    
    // 查找并执行补全函数
    result := executeCompleter(ctx)
    
    // 输出结果（格式：候选1\n候选2\n:指令）
    printCompletionResult(result)
    os.Exit(0)
    return nil
}
```

## 4. 架构流程

```
用户生成脚本:
┌─────────────────────────────────────────────────────────────┐
│  myapp --completion=bash --comp-mode=dynamic               │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  CompletionHandler                           │
│  1. 检测到 --comp-mode=dynamic                              │
│  2. 检查是否有位置参数（-- 之后的参数）                      │
│  3. 无位置参数 → 生成动态补全脚本                           │
└─────────────────────────────────────────────────────────────┘

Shell 动态补全:
┌─────────────────────────────────────────────────────────────┐
│  myapp --completion=bash --comp-mode=dynamic -- deploy --e  │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  CompletionHandler                           │
│  1. 检测到 --comp-mode=dynamic                              │
│  2. 有位置参数 ["deploy", "--e"]                            │
│  3. 执行动态补全逻辑                                        │
│  4. 返回补全结果                                            │
└─────────────────────────────────────────────────────────────┘
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
    
    # 构造参数（-- 之后的所有参数）
    local args=()
    for ((i=1; i < ${#words[@]}; i++)); do
        args+=("${words[i]}")
    done
    
    # 调用程序获取补全结果
    # 格式: myapp --completion=bash --comp-mode=dynamic -- arg1 arg2
    local out
    out=$("${words[0]}" "--completion=bash" "--comp-mode=dynamic" "--" "${args[@]}" 2>/dev/null)
    
    # 解析结果（最后一行是指令）
    local directive="${out##*:}"
    out="${out%%:*}"
    
    if [[ "$directive" == "$out" ]]; then
        directive=0
    fi
    
    # 处理指令
    if (( directive & 2 )); then  # NoSpace
        if [[ $(type -t compopt) = "builtin" ]]; then
            compopt -o nospace
        fi
    fi
    if (( directive & 4 )); then  # NoFileComp
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
    if [[ ${#COMPREPLY[@]} -eq 0 && $((directive & 4)) -eq 0 ]]; then
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
    $program = $tokens[0]
    
    # 构造参数（所有位置参数）
    $args = $tokens[1..($tokens.Count-1)]
    
    # 调用程序获取补全结果
    # 格式: myapp --completion=pwsh --comp-mode=dynamic -- arg1 arg2
    $out = & $program "--completion=pwsh" "--comp-mode=dynamic" "--" $args 2>$null
    
    # 解析结果（最后一行是指令）
    $directive = [int]($out[-1] -replace '^:','')
    $completions = $out[0..($out.Count-2)]
    
    # 处理指令
    $noSpace = ($directive -band 2) -ne 0
    
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

## 6. 用户 API

```go
package main

import (
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    cmd := qflag.NewCmd("myapp", "示例应用", qflag.ExitOnError)
    
    // 启用动态补全（注册 --comp-mode 标志）
    cmd.EnableDynamicCompletion()
    
    // 添加子命令和标志
    deploy := cmd.SubCmd("deploy", "部署应用")
    deploy.String("env", "", "部署环境")
    
    // 注册自定义补全函数
    cmd.RegisterCompleter("env", func(ctx qflag.CompletionContext) qflag.CompletionResult {
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

## 7. 需要修改的文件

### 7.1 types/builtin.go

```go
const (
    HelpFlag BuiltinFlagType = iota
    VersionFlag
    CompletionFlag
    CompModeFlag  // 新增
)

const (
    CompModeFlagName = "comp-mode"
    CompModeStatic   = "static"
    CompModeDynamic  = "dynamic"
)
```

### 7.2 builtin/manager.go

在 `RegisterBuiltinFlags` 中新增 `CompModeFlag` 的处理：

```go
case types.CompModeFlag:
    var desc string
    if cmd.Config().UseChinese {
        desc = fmt.Sprintf("补全模式: %s(默认) 或 %s", types.CompModeStatic, types.CompModeDynamic)
    } else {
        desc = fmt.Sprintf("Completion mode: %s(default) or %s", types.CompModeStatic, types.CompModeDynamic)
    }
    modeFlag := flag.NewEnumFlag(types.CompModeFlagName, "", desc, types.CompModeStatic, 
        []string{types.CompModeStatic, types.CompModeDynamic})
    cmd.AddFlag(modeFlag)
```

### 7.3 builtin/completion_handler.go

修改 `Handle` 方法，支持动态补全逻辑。

### 7.4 新增 internal/completion/dynamic.go

实现动态补全的核心逻辑：
- `buildCompletionContext` - 构建补全上下文
- `executeCompleter` - 执行注册的补全函数
- `GenDynamicAndPrint` - 生成动态补全脚本

## 8. 优势

1. **复用现有标志**: 不需要新增复杂的内置子命令
2. **向后兼容**: `--completion` 行为完全不变，新增 `--comp-mode` 可选
3. **统一入口**: 所有补全相关功能都通过 `--completion` 标志
4. **清晰分离**: 生成脚本 vs 执行补全，通过位置参数区分

## 9. 总结

这个方案通过新增 `--comp-mode` 内置标志，复用现有的 `--completion` 标志：

- **生成脚本**: `myapp --completion=bash --comp-mode=dynamic`
- **执行补全**: `myapp --completion=bash --comp-mode=dynamic -- arg1 arg2`

完全复用现有架构，只需修改 `CompletionHandler` 和新增 `CompModeFlag` 处理器。
