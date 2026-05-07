# 动态补全模板简化方案

## 目标

将 Bash 和 PowerShell 动态补全模板中的 `_intelligent_match` 函数简化，统一使用 `__complete fuzzy` 指令处理所有匹配逻辑。

## 当前问题

- Bash 模板：`~110` 行的匹配逻辑
- PowerShell 模板：`~80` 行的匹配逻辑
- 两者逻辑重复，只是语法不同
- 维护成本高，Bug 需要修复两次

## 简化方案

### 核心思想

使用 `__complete fuzzy` 一个指令替代原来的 4 级匹配逻辑：
- 精确前缀匹配
- 大小写不敏感前缀匹配
- 模糊匹配
- 子字符串匹配

`go-kit/fuzzy` 的 `Find` 函数本身就包含前缀匹配能力，前缀匹配会获得高分。

---

## Bash 模板修改

### 原代码（删除）

删除 `_{{.ProgramName}}_intelligent_match` 函数中第 44-154 行的所有分级匹配逻辑，保留：
- 函数定义和参数解析
- 空模式快速返回
- 候选项过多时的 `compgen` 回退
- 模糊匹配调用

### 新代码

```bash
# 智能补全匹配函数 - 使用 fuzzy 统一处理
# 参数: $1=输入模式, $2=候选选项字符串(用|分隔)
_{{.ProgramName}}_intelligent_match() {
    local pattern="$1"
    local options_str="$2"

    # 解析候选选项到数组
    local opts_arr
    IFS='|' read -ra opts_arr <<< "$options_str"
    local total_candidates=${#opts_arr[@]}

    # 快速路径: 模式为空时返回所有选项
    [[ -z "$pattern" ]] && {
        COMPREPLY=("${opts_arr[@]}")
        return 0
    }

    # 性能保护: 候选项过多时回退到传统前缀匹配
    if [[ $total_candidates -gt {{.ProgramName}}_FUZZY_MAX_CANDIDATES ]]; then
        local opts
        printf -v opts '%s ' "${opts_arr[@]}"
        opts="${opts% }"
        COMPREPLY=($(compgen -W "$opts" -- "$pattern"))
        return 0
    fi

    # 统一使用 fuzzy 处理所有匹配逻辑
    local fuzzy_output
    fuzzy_output=$({{.ProgramName}} __complete fuzzy "$pattern" "${opts_arr[@]}")

    # 读取匹配结果到 COMPREPLY
    while IFS= read -r match; do
        [[ -n "$match" ]] && COMPREPLY+=("$match")
    done <<< "$fuzzy_output"

    return 0
}
```

### 可删除的配置参数

以下参数在简化后不再需要，可以删除：
```bash
# 删除以下行（约第 10-13 行）
{{.ProgramName}}_FUZZY_MIN_PATTERN_LENGTH=2
{{.ProgramName}}_FUZZY_MAX_RESULTS=8
```

保留：
```bash
{{.ProgramName}}_FUZZY_COMPLETION_ENABLED=1  # 仍可用于总开关
{{.ProgramName}}_FUZZY_MAX_CANDIDATES=150     # 用于性能保护阈值
```

### 可选：删除 fuzzy_match 包装函数

原 `_{{.ProgramName}}_fuzzy_match` 函数（第 32-42 行）可以内联或保留：

**方案 A - 内联（更简洁）**：
直接调用 `{{.ProgramName}} __complete fuzzy`

**方案 B - 保留（更清晰）**：
保留包装函数，但简化注释

---

## PowerShell 模板修改

### 原代码（删除）

删除 `Get-{{.SanitizedName}}IntelligentMatches` 函数中第 60-120 行的所有分级匹配逻辑。

### 新代码

```powershell
# 智能补全匹配函数 - 使用 fuzzy 统一处理
# 参数: $Pattern=输入模式, $Options=候选选项数组
function Get-{{.SanitizedName}}IntelligentMatches {
    param(
        [string]$Pattern,
        [array]$Options
    )

    $totalCandidates = $Options.Count

    # 快速路径: 空模式时返回所有选项
    if ([string]::IsNullOrEmpty($Pattern)) {
        return $Options
    }

    # 性能保护: 候选项过多时回退到传统前缀匹配
    if ($totalCandidates -gt $script:{{.SanitizedName}}_FUZZY_MAX_CANDIDATES) {
        return $Options | Where-Object { $_ -like "$Pattern*" }
    }

    # 统一使用 fuzzy 处理所有匹配逻辑
    $output = & {{.ProgramName}} __complete fuzzy $Pattern $Options
    return $output -split "`n" | Where-Object { $_ }
}
```

### 可删除的配置参数

```powershell
# 删除以下变量（约第 20-23 行）
$script:{{.SanitizedName}}_FUZZY_MIN_PATTERN_LENGTH = 2
$script:{{.SanitizedName}}_FUZZY_MAX_RESULTS = 10
```

保留：
```powershell
$script:{{.SanitizedName}}_FUZZY_COMPLETION_ENABLED = $true
$script:{{.SanitizedName}}_FUZZY_MAX_CANDIDATES = 120
```

### 可选：删除 Get-FuzzyMatch 包装函数

原 `Get-{{.SanitizedName}}FuzzyMatch` 函数（第 35-45 行）可以内联或保留。

---

## 代码量对比

| 模板 | 原代码行数 | 新代码行数 | 减少比例 |
|------|-----------|-----------|---------|
| Bash | ~110 行 | ~25 行 | ~77% |
| PowerShell | ~80 行 | ~20 行 | ~75% |

---

## 实现步骤

### 步骤 1: 修改 Bash 模板

1. 打开 `internal/completion/templates/bash_dynamic.tmpl`
2. 删除第 10-13 行（MIN_PATTERN_LENGTH, MAX_RESULTS）
3. 替换第 44-154 行的 `_intelligent_match` 函数体
4. （可选）删除或简化第 32-42 行的 `fuzzy_match` 函数
5. 测试验证

### 步骤 2: 修改 PowerShell 模板

1. 打开 `internal/completion/templates/pwsh_dynamic.tmpl`
2. 删除第 20-23 行（MIN_PATTERN_LENGTH, MAX_RESULTS）
3. 替换第 60-120 行的 `Get-IntelligentMatches` 函数体
4. （可选）删除或简化第 35-45 行的 `Get-FuzzyMatch` 函数
5. 测试验证

### 步骤 3: 测试验证

1. 生成新的补全脚本
2. 在 Bash 和 PowerShell 中测试补全功能
3. 验证以下场景：
   - 空输入返回所有候选
   - 前缀匹配正常工作
   - 模糊匹配正常工作
   - 候选项过多时性能正常

---

## 潜在问题及解决方案

### 问题 1: fuzzy 算法排序与原来不完全一致

**现象**：前缀匹配不再绝对优先于模糊匹配

**影响**：低，fuzzy 算法前缀匹配分数通常很高

**解决**：如需要严格前缀优先，可在 Go 层修改 fuzzy 分数权重

### 问题 2: 子进程调用开销

**现象**：每次补全都要启动子进程调用 Go

**影响**：中等，约 10-50ms 延迟

**缓解**：保留 150 候选项阈值，大量候选时用 Shell 原生匹配

### 问题 3: 结果数量控制

**现象**：原来通过 MAX_RESULTS 限制返回数量

**解决**：
- 方案 A: 在 Go 层 `handleFuzzy` 中添加限制逻辑
- 方案 B: 在 Shell 层用 `head` 或数组切片限制

推荐方案 B（Shell 层控制更灵活）：
```bash
# Bash
fuzzy_output=$({{.ProgramName}} __complete fuzzy "$pattern" "${opts_arr[@]}" | head -n 20)

# PowerShell
$output = & {{.ProgramName}} __complete fuzzy $Pattern $Options | Select-Object -First 20
```

---

## 备选方案

如果希望保留严格的分级匹配逻辑，可以新增 `match` 指令：

```bash
# 用法
yourapp __complete match <模式> <候选1> [候选2] ...

# 功能
# 1. 精确前缀匹配
# 2. 大小写不敏感前缀匹配  
# 3. 子字符串匹配
# 4. 模糊匹配（可选）
```

但会增加维护成本，建议先用简化方案验证效果。

---

## 总结

通过统一使用 `__complete fuzzy` 指令：

1. **代码量大幅减少**：Bash ~77%，PowerShell ~75%
2. **逻辑统一**：所有匹配算法在 Go 层实现
3. **易于维护**：Bug 只需修复一次
4. **易于扩展**：新增 Shell 支持只需写 ~20 行代码

建议先实现简化方案，根据实际使用反馈决定是否进一步优化。
