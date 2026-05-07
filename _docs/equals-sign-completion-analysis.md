# 等号赋值格式补全支持分析

## 问题描述

当前补全系统不支持 `--config=val` 这种等号赋值格式的补全。

### 场景示例

**用户输入**：`dynamic.exe --config=`

按 Tab 时期望行为：
- 识别为 `--config` 标志的赋值格式
- 返回 `--config` 的枚举值（如 `service`, `pod`, `deployment`）

**当前行为**：
- `$cur = "--config="`
- `cur` 以 `-` 开头，但不是有效标志名（因为包含 `=`）
- 被当作普通候选项匹配，无匹配结果
- 结果：无补全

## 实现方案

### 方案 1：在 Shell 模板中处理

在 Bash 和 PowerShell 模板中分别添加等号格式解析逻辑。

#### Bash 实现

```bash
# 检查是否是 = 赋值格式（标志名=，无值）
if [[ "$cur" =~ ^--[^=]+=$ ]]; then
    # 提取标志名（如 --config= 提取为 --config）
    local flag_name="${cur%=}"
    # 重新调用 all 指令，传入空 cur 和 flag_name= 作为 prev
    result=$({{.ProgramName}} __complete all "" "$flag_name=" "${cmd_args[@]}")
    # 解析结果，提取枚举值作为补全
    while IFS= read -r line; do
        case "$line" in
            ENUM:*) enum_values="${line#ENUM:}" ;;
        esac
    done <<< "$result"
    if [[ -n "$enum_values" ]]; then
        read -ra COMPREPLY <<< "$enum_values"
    fi
    return 0
# 检查是否是 = 赋值格式（标志名=值）
elif [[ "$cur" =~ ^--[^=]+= ]]; then
    # 提取标志名和当前值
    local flag_name="${cur%%=*}="
    local cur_value="${cur#*=}"
    # 重新调用 all 指令
    result=$({{.ProgramName}} __complete all "$cur_value" "$flag_name" "${cmd_args[@]}")
    # 解析 MATCHES 作为补全结果
    while IFS= read -r line; do
        case "$line" in
            MATCHES:*) matches="${line#MATCHES:}" ;;
        esac
    done <<< "$result"
    if [[ -n "$matches" ]]; then
        read -ra COMPREPLY <<< "$matches"
    fi
    return 0
fi
```

#### PowerShell 实现

```powershell
# 检查是否是 = 赋值格式（标志名=，无值）
if ($wordToComplete -match '^--[^=]+=$') {
    $flagName = $wordToComplete.TrimEnd('=')
    # 重新调用 all 指令
    $result = & {{.ProgramName}} __complete all "" "$flagName=" @cmdArgs
    # 解析结果，提取枚举值
    foreach ($line in $result) {
        if ($line -match '^ENUM:(.*)$') {
            $enumValues = $matches[1] -split ' ' | Where-Object { $_ }
        }
    }
    return $enumValues
}
# 检查是否是 = 赋值格式（标志名=值）
elseif ($wordToComplete -match '^--[^=]+=') {
    $flagName = ($wordToComplete -split '=')[0] + '='
    $curValue = ($wordToComplete -split '=')[1]
    # 重新调用 all 指令
    $result = & {{.ProgramName}} __complete all $curValue $flagName @cmdArgs
    # 解析 MATCHES
    foreach ($line in $result) {
        if ($line -match '^MATCHES:(.*)$') {
            $matchResults = $matches[1] -split ' ' | Where-Object { $_ }
        }
    }
    return $matchResults
}
```

#### 方案 1 优缺点

**优点**：
- 不修改 Go 代码，只修改模板
- 可以针对不同 Shell 优化细节

**缺点**：
- 两个模板都要修改，逻辑重复
- 增加了模板的复杂度
- 需要额外调用 `all` 指令（性能开销）
- 维护成本高，后续修改需要同步两个模板

---

### 方案 2：在 Go 代码中处理

在 `handleAll` 函数中统一处理等号赋值格式。

#### Go 实现

```go
// handleAll 处理 all 指令，一次性返回所有补全信息
func handleAll(root types.Command, args []string) error {
    if len(args) < 2 {
        return fmt.Errorf("用法: __complete all <cur> <prev> [cmd_args...]")
    }

    // 解析参数
    cur := strings.Trim(args[0], `"'`)
    prev := strings.Trim(args[1], `"'`)
    cmdArgs := []string{}
    if len(args) > 2 {
        cmdArgs = args[2:]
    }

    // ========== 新增：处理等号赋值格式 ==========
    // 支持 --flag= 和 --flag=value 两种格式
    if strings.Contains(cur, "=") {
        parts := strings.SplitN(cur, "=", 2)
        flagName := parts[0]
        cur = parts[1]  // = 后的部分作为 cur（可能为空）
        prev = flagName + "="  // 标志名= 作为 prev
        
        // 检查是否是有效标志
        if strings.HasPrefix(flagName, "-") {
            // 标记为标志上下文
            // 后续逻辑会正常处理枚举值匹配
        }
    }
    // ===========================================

    // 1. 计算上下文
    tokens := append([]string{""}, cmdArgs...)
    contextResult := CalculateContext(root, tokens, len(tokens))
    context := "/"
    if contextResult != nil {
        context = contextResult.Context
    }

    // 2. 获取候选项
    candidates, _ := GetCandidates(root, context)

    // 3. 检查是否是标志值补全
    var enumValues []string
    var matchStrings []string
    isFlag := false

    // 检查 prev 是否是标志：以 "-" 开头
    // 注意：当在根命令下刚输入命令名后按 Tab, prev 是程序名 (如 "dynamic.exe")
    // 正常情况下程序名不会以 "-" 开头，所以不会误判为标志
    // -- 是标志结束符，不是标志，应该按普通参数处理
    if strings.HasPrefix(prev, "-") && prev != "--" {
        isFlag = true
        // 是标志，获取枚举值
        // 注意：prev 可能包含末尾的 =（如 --config=）
        flagName := strings.TrimSuffix(prev, "=")
        enumValues, _ = GetEnumValues(root, context, flagName)

        if len(enumValues) > 0 {
            // 枚举类型：对枚举值进行模糊匹配
            if cur == "" {
                // 空输入时返回所有枚举值
                matchStrings = enumValues
            } else {
                matches := fuzzy.CompletePrefix(cur, enumValues)
                matchStrings = make([]string, len(matches))
                for i, match := range matches {
                    matchStrings[i] = match.Str
                }
            }
        }
        // 非枚举类型: matchStrings 保持为空，由 Shell 处理为路径补全
    } else {
        // 不是标志，对候选项进行模糊匹配
        if cur == "" {
            // 空输入时返回所有候选项
            matchStrings = candidates
        } else {
            matches := fuzzy.CompletePrefix(cur, candidates)
            matchStrings = make([]string, len(matches))
            for i, match := range matches {
                matchStrings[i] = match.Str
            }
        }
    }

    // 5. 输出结果（带前缀的多行格式）
    fmt.Printf("CONTEXT:%s\n", context)
    fmt.Printf("CUR:%s\n", cur)
    fmt.Printf("PREV:%s\n", prev)
    fmt.Printf("CANDIDATES:%s\n", strings.Join(candidates, " "))
    fmt.Printf("ENUM:%s\n", strings.Join(enumValues, " "))
    fmt.Printf("MATCHES:%s\n", strings.Join(matchStrings, " "))
    fmt.Printf("IS_FLAG:%v\n", isFlag && len(enumValues) > 0)

    return nil
}
```

#### 需要同时修改 GetEnumValues

```go
// GetEnumValues 获取指定标志的枚举值
func GetEnumValues(root types.Command, context string, flagName string) ([]string, error) {
    // 移除可能的 = 后缀
    flagName = strings.TrimSuffix(flagName, "=")
    
    // ... 原有逻辑
}
```

#### 方案 2 优缺点

**优点**：
- 一处修改，跨平台统一（Bash 和 PowerShell 自动支持）
- 逻辑集中，易于维护和测试
- 不需要额外调用 `all` 指令（性能好）
- 代码简洁清晰

**缺点**：
- 需要重新编译 Go 程序
- 需要修改两处（`handleAll` 和 `GetEnumValues`）

---

## 边界情况分析

| 场景 | 期望行为 | 处理方案 |
|------|----------|----------|
| `--config=`（空值） | 返回所有枚举值 | `cur=""`，`prev="--config="`，返回全部枚举值 |
| `--config=va`（部分值） | 模糊匹配枚举值 | `cur="va"`，`prev="--config="`，模糊匹配 |
| `--config=/path`（路径值） | 路径补全 | 非枚举标志，返回空，Shell 回退到路径补全 |
| `-c=val`（短标志） | 同样支持 | 逻辑相同，支持单横杠短标志 |
| `--config==`（值以=开头） | 值就是 `=` | `SplitN(cur, "=", 2)` 确保只分割一次 |
| `--`（标志结束符） | 不处理 | `prev != "--"` 已经排除 |

---

## 推荐方案

**推荐方案 2（Go 代码）**，理由：

1. **维护成本低**：一处修改，两个 Shell 自动受益
2. **逻辑清晰**：等号解析在 Go 层统一处理，Shell 模板保持简洁
3. **性能好**：不需要额外调用 `all` 指令
4. **易于测试**：可以在 Go 单元测试中覆盖等号格式
5. **符合设计**：动态补全的核心逻辑应该在 Go 层，Shell 只负责调用和展示

---

## 实施步骤（方案 2）

1. 修改 `handleAll` 函数，添加等号格式解析
2. 修改 `GetEnumValues` 函数，支持带 `=` 的标志名
3. 添加单元测试覆盖等号格式
4. 编译测试
5. 重新生成补全脚本
6. 手动测试各种等号场景
