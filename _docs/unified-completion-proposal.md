# 统一补全指令方案

## 概述

将多次子进程调用合并为一次，提升动态补全性能。

## 当前问题

每次按 Tab 触发 3-4 次子进程调用：
1. `__complete context` - 计算上下文路径
2. `__complete candidates` - 获取候选项
3. `__complete enum` - 获取枚举值（如果是标志）
4. `__complete fuzzy` - 模糊匹配

## 新方案：统一 `all` 指令

### 调用方式

```bash
# 一次调用获取所有信息
# 注意：单独参数放前面，切片参数放后面
result=$(dy __complete all "$cur" "$prev" "${cmd_args[@]}")
```

### 参数说明

| 参数 | 位置 | 说明 | 示例 |
|------|------|------|------|
| `$cur` | args[0] | 当前输入的词 | `ser` 或空 |
| `$prev` | args[1] | 前一个词 | `-k` 或 `service` |
| `${cmd_args[@]}` | args[2:] | 已输入的子命令参数 | `service list` |

### 返回格式

```
CONTEXT:<上下文路径>
CUR:<当前输入>
PREV:<前一个词>
CANDIDATES:<候选项列表（空格分隔）>
ENUM:<枚举值列表（空格分隔，可能为空）>
MATCHES:<匹配结果（空格分隔）>
IS_FLAG:<true|false>
```

### 返回示例

**示例1：根命令补全**
```bash
$ dy __complete all "" "" ""

CONTEXT:/
CUR:
PREV:
CANDIDATES:--help -h --version -v --completion --config service deployment config completion
ENUM:
MATCHES:--help --version --completion --config service deployment config completion
IS_FLAG:false
```

**示例2：枚举标志补全（无输入）**
```bash
$ dy __complete all "" "-k" ""

CONTEXT:/
CUR:
PREV:-k
CANDIDATES:--help -h --version -v --completion --config service deployment config completion
ENUM:service deployment config
MATCHES:service deployment config
IS_FLAG:true
```

**示例3：枚举标志补全（有输入）**
```bash
$ dy __complete all "ser" "-k" ""

CONTEXT:/
CUR:ser
PREV:-k
CANDIDATES:--help -h --version -v --completion --config service deployment config completion
ENUM:service deployment config
MATCHES:service
IS_FLAG:true
```

**示例4：子命令补全**
```bash
$ dy __complete all "li" "service" "service"

CONTEXT:/service/
CUR:li
PREV:service
CANDIDATES:--help -h --namespace -n list create delete update logs
ENUM:
MATCHES:list
IS_FLAG:false
```

**示例5：嵌套子命令**
```bash
$ dy __complete all "" "list" "service" "list"

CONTEXT:/service/list/
CUR:
PREV:list
CANDIDATES:--help -h --namespace -n --all -a --format -f
ENUM:
MATCHES:--help -h --namespace -n --all -a --format -f
IS_FLAG:false
```

## Go 端实现

```go
// HandleAll 处理 all 指令，一次性返回所有补全信息
//
// 参数:
//   - root: 根命令实例
//   - args: [cur, prev, cmd_args...]
//     args[0]: cur - 当前输入的词
//     args[1]: prev - 前一个词
//     args[2:]: cmd_args - 已输入的子命令参数
//
// 返回值:
//   - error: 处理错误
func HandleAll(root types.Command, args []string) error {
    if len(args) < 2 {
        return fmt.Errorf("用法: __complete all <cur> <prev> [cmd_args...]")
    }

    // 解析参数
    cur := args[0]       // 当前输入
    prev := args[1]      // 前一个词
    cmdArgs := []string{}
    if len(args) > 2 {
        cmdArgs = args[2:] // 子命令参数
    }

    // 1. 计算上下文
    context := calculateContext(root, cmdArgs)

    // 2. 获取候选项
    candidates := getCandidates(root, context)

    // 3. 检查是否是标志值补全
    var enumValues []string
    var matches []string
    isFlag := false

    if strings.HasPrefix(prev, "-") {
        isFlag = true
        // 是标志，获取枚举值
        enumValues = getEnumValuesForFlag(root, context, prev)

        if len(enumValues) > 0 {
            // 枚举类型：对枚举值进行模糊匹配
            matches = fuzzyMatch(cur, enumValues)
        } else {
            // 非枚举类型：返回空（Shell 会处理为路径补全）
            matches = []string{}
        }
    } else {
        // 不是标志，对候选项进行模糊匹配
        matches = fuzzyMatch(cur, candidates)
    }

    // 4. 输出结果（带前缀的多行格式）
    fmt.Printf("CONTEXT:%s\n", context)
    fmt.Printf("CUR:%s\n", cur)
    fmt.Printf("PREV:%s\n", prev)
    fmt.Printf("CANDIDATES:%s\n", strings.Join(candidates, " "))
    fmt.Printf("ENUM:%s\n", strings.Join(enumValues, " "))
    fmt.Printf("MATCHES:%s\n", strings.Join(matches, " "))
    fmt.Printf("IS_FLAG:%v\n", isFlag && len(enumValues) > 0)

    return nil
}
```

## Bash 端实现

```bash
#!/usr/bin/env bash

# ==================== 主补全函数 ====================
_{{.ProgramName}}_complete() {
    local cur prev words cword i
    COMPREPLY=()

    if declare -F _get_comp_words_by_ref >/dev/null 2>&1; then
        _get_comp_words_by_ref -n =: cur prev words cword
    else
        words=("${COMP_WORDS[@]}")
        cword=$COMP_CWORD
        cur="${words[cword]}"
        prev="${words[cword-1]}"
    fi

    # 路径快速路径
    if [[ "$cur" == *"/"* || "$cur" == *"."* || "$cur" == *"~"* ]]; then
        COMPREPLY=($(compgen -f -d -- "$cur"))
        return 0
    fi

    # 提取子命令参数（从索引1开始，跳过程序名）
    local cmd_args=()
    for ((i=1; i < cword; i++)); do
        cmd_args+=("${words[i]}")
    done

    # ========== 关键：一次调用获取所有信息 ==========
    # 注意参数顺序：cur prev cmd_args
    local result context candidates enum_values matches is_flag
    result=$({{.ProgramName}} __complete all "$cur" "$prev" "${cmd_args[@]}")

    # 解析结果（按行读取，根据前缀提取）
    while IFS= read -r line; do
        case "$line" in
            CONTEXT:*) context="${line#CONTEXT:}" ;;
            CUR:*) ;;  # 可选：用于调试
            PREV:*) ;;  # 可选：用于调试
            CANDIDATES:*) candidates="${line#CANDIDATES:}" ;;
            ENUM:*) enum_values="${line#ENUM:}" ;;
            MATCHES:*) matches="${line#MATCHES:}" ;;
            IS_FLAG:*) is_flag="${line#IS_FLAG:}" ;;
        esac
    done <<< "$result"

    # 根据结果决定补全行为
    if [[ "$is_flag" == "true" ]]; then
        # 枚举类型标志，显示匹配结果
        if [[ -n "$matches" ]]; then
            read -ra COMPREPLY <<< "$matches"
        fi
    elif [[ -z "$enum_values" && -n "$cur" && "$prev" =~ ^- ]]; then
        # 非枚举类型标志且当前输入非空，使用路径补全
        COMPREPLY=($(compgen -f -d -- "$cur"))
    else
        # 普通补全，显示匹配结果
        if [[ -n "$matches" ]]; then
            read -ra COMPREPLY <<< "$matches"
        fi
    fi

    return 0
}

complete -F _{{.ProgramName}}_complete {{.ProgramName}}
```

## PowerShell 端实现

```powershell
function Get-{{.SanitizedName}}PathCompletions {
    param([string]$WordToComplete)

    $pathMatches = [System.Collections.ArrayList]::new()
    $basePath = if ($WordToComplete -and (Split-Path $WordToComplete -Parent)) {
        Split-Path $WordToComplete -Parent
    } else { "." }
    $fileName = if ($WordToComplete) { Split-Path $WordToComplete -Leaf } else { "" }
    $filePattern = "$fileName*"

    try {
        $items = Get-ChildItem -Path $basePath -ErrorAction SilentlyContinue |
            Where-Object { $_.Name -like $filePattern }

        foreach ($item in $items) {
            $fullPath = if ($basePath -eq ".") { $item.Name } else { Join-Path $basePath $item.Name }
            if ($item.PSIsContainer) {
                [void]$pathMatches.Add("$fullPath/")
            } else {
                [void]$pathMatches.Add($fullPath)
            }
        }
    }
    catch {
        Write-Debug "路径访问失败: $($_.Exception.Message)"
    }

    return $pathMatches.ToArray()
}

$scriptBlock = {
    param($wordToComplete, $commandAst, $cursorPosition)

    try {
        $tokens = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
        if (-not $tokens -or $tokens.Count -eq 0) {
            return @()
        }

        $currentIndex = $tokens.Count - 1
        $prevElement = if ([string]::IsNullOrEmpty($wordToComplete)) {
            if ($tokens.Count -gt 0) { $tokens[$tokens.Count - 1] } else { $null }
        } else {
            if ($currentIndex -ge 1) { $tokens[$currentIndex - 1] } else { $null }
        }

        # 路径快速路径
        if ($wordToComplete -match '[/\~\.]' -or $wordToComplete -like './*' -or $wordToComplete -like '../*') {
            return Get-{{.SanitizedName}}PathCompletions -WordToComplete $wordToComplete
        }

        # 提取子命令参数（从索引1开始，跳过程序名）
        $cmdArgs = @()
        for ($i = 1; $i -lt $currentIndex; $i++) {
            $cmdArgs += $tokens[$i]
        }

        # ========== 关键：一次调用获取所有信息 ==========
        # 注意参数顺序：wordToComplete prevElement cmdArgs
        $result = & {{.ProgramName}} __complete all $wordToComplete $prevElement @cmdArgs

        # 解析结果
        $context = ""
        $candidates = @()
        $enumValues = @()
        $matches = @()
        $isFlag = $false

        foreach ($line in $result) {
            if ($line -match '^CONTEXT:(.+)$') {
                $context = $matches[1]
            }
            elseif ($line -match '^CUR:(.*)$') {
                # 可选：用于调试
            }
            elseif ($line -match '^PREV:(.*)$') {
                # 可选：用于调试
            }
            elseif ($line -match '^CANDIDATES:(.*)$') {
                $candidates = $matches[1] -split ' ' | Where-Object { $_ }
            }
            elseif ($line -match '^ENUM:(.*)$') {
                $enumValues = $matches[1] -split ' ' | Where-Object { $_ }
            }
            elseif ($line -match '^MATCHES:(.*)$') {
                $matches = $matches[1] -split ' ' | Where-Object { $_ }
            }
            elseif ($line -match '^IS_FLAG:(.+)$') {
                $isFlag = [bool]::Parse($matches[1])
            }
        }

        # 根据结果决定补全行为
        if ($isFlag -and $enumValues.Count -gt 0) {
            # 枚举类型标志
            return $matches
        }
        elseif ($enumValues.Count -eq 0 -and -not [string]::IsNullOrEmpty($wordToComplete) -and $prevElement -match '^-') {
            # 非枚举类型标志，使用路径补全
            return Get-{{.SanitizedName}}PathCompletions -WordToComplete $wordToComplete
        }
        else {
            # 普通补全
            $matchingOptions = [System.Collections.ArrayList]::new()
            $flagRegex = [regex]::new('^-')
            foreach ($match in $matches) {
                $result = if ($flagRegex.IsMatch($match)) { $match } else { "$match " }
                [void]$matchingOptions.Add($result)
            }
            return $matchingOptions.ToArray()
        }
    }
    catch {
        Write-Debug "补全错误: $($_.Exception.Message)"
        return @()
    }
}

Register-ArgumentCompleter -CommandName '{{.ProgramName}}' -ScriptBlock $scriptBlock

${{SanitizedName}}_withoutExt = [System.IO.Path]::GetFileNameWithoutExtension('{{.ProgramName}}')
if (${{SanitizedName}}_withoutExt -ne '{{.ProgramName}}') {
    Register-ArgumentCompleter -CommandName ${{SanitizedName}}_withoutExt -ScriptBlock $scriptBlock
}
```

## 优势

| 优势 | 说明 |
|------|------|
| **性能提升** | 从 3-4 次子进程调用减少到 1 次 |
| **简单解析** | 按行读取，前缀匹配，无需 JSON 解析器 |
| **可扩展** | 需要新增字段只需加一行 `KEY:value` |
| **向后兼容** | 不影响现有指令，新增 `all` 指令即可 |
| **跨平台** | Bash 和 PowerShell 解析逻辑一致 |
| **自描述** | 返回包含 CUR/PREV，便于调试 |

## 实现步骤

1. 在 `dynamic.go` 中添加 `handleAll` 函数
2. 在 `HandleDynamicComplete` 中添加 `all` 指令路由
3. 修改 `bash_dynamic.tmpl` 使用新的 `all` 指令（替换多次调用为一次调用）
4. 修改 `pwsh_dynamic.tmpl` 使用新的 `all` 指令（替换多次调用为一次调用）
5. 测试验证

## 修改范围

### Go 代码修改
- `internal/completion/dynamic.go` - 添加 `handleAll` 函数和路由

### 模板修改
- `internal/completion/templates/bash_dynamic.tmpl` - 简化补全逻辑
- `internal/completion/templates/pwsh_dynamic.tmpl` - 简化补全逻辑

**注意**：直接在现有动态模板上修改，不创建新模板文件

## 注意事项

- **参数顺序**：`cur` 和 `prev` 必须放在 `cmd_args` 前面
- **向后兼容**：保持现有指令不变，确保旧模板仍能工作
- **空值处理**：空字符串用空值表示，不省略字段
- **调试友好**：返回包含输入参数，便于排查问题
- **错误处理**：Go 端出错时返回非零退出码，Shell 端静默处理
