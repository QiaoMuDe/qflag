# 布尔标志补全修复方案

## 问题背景

### 布尔标志的特性

布尔标志（如 `--help`, `--verbose`）不需要指定值，存在即表示启用。例如：
```bash
dynamic.exe --help          # 显示帮助信息
dynamic.exe --verbose       # 启用详细输出
```

### 原逻辑的问题

当用户输入 `dynamic.exe --help <Tab>` 时，原补全逻辑无法正确处理：

```bash
# 原模板逻辑（问题版本）
if [[ "$is_flag" == "true" ]]; then
    # 枚举类型标志 → 显示枚举值
    read -ra COMPREPLY <<< "$matches"
elif [[ "$prev" =~ ^- ]]; then
    # 非枚举类型标志 → 路径补全  ❌ 布尔标志也走这里！
    COMPREPLY=($(compgen -f -d -- "$cur"))
elif [[ -n "$matches" ]]; then
    # 普通补全 → 显示候选项
    read -ra COMPREPLY <<< "$matches"
fi
```

**问题分析**：
1. `prev = "--help"`（布尔标志）
2. `IS_FLAG = false`（布尔标志没有枚举值）
3. `prev =~ ^-` 匹配成功（`--help` 以 `-` 开头）
4. 错误地执行了**路径补全**，而不是候选项补全

**期望行为**：`--help <Tab>` 应该补全其他标志（如 `--config`）或子命令，而不是文件路径。

## 解决方案

### 核心思想

通过 `FlagType` 区分标志类型：
- **布尔标志**（`FlagTypeBool`）：不需要值，按**普通候选项补全**处理
- **枚举标志**（`FlagTypeEnum`）：需要值，补全**枚举值**
- **其他类型**（String/Int等）：需要值，按**路径补全**处理

### 1. Go 代码层实现

#### 修改 `handleAll` 函数

```go
// 3. 执行补全逻辑
var matchStrings []string
var enumValues []string

// 判断是否是标志值补全上下文
isFlagValueCompletion := strings.HasPrefix(prev, "-") && prev != "--" && !strings.Contains(prev, "=")

if isFlagValueCompletion {
    // ========== 标志值补全 ==========
    flagType, found := getFlagType(root, context, prev)

    if !found {
        // 标志不存在，按普通候选项补全
        matchStrings = fuzzyMatch(candidates, cur)
    } else {
        switch flagType {
        case types.FlagTypeBool:
            // 布尔标志：不需要值，补全其他标志/子命令
            matchStrings = fuzzyMatch(candidates, cur)

        case types.FlagTypeEnum:
            // 枚举标志：获取枚举值并模糊匹配
            enumValues, _ = GetEnumValues(root, context, prev)
            matchStrings = fuzzyMatch(enumValues, cur)

        default:
            // 其他类型（String/Int/Duration/Size等）：需要值
            // matchStrings 保持为空，由 Shell 回退到路径补全
        }
    }
} else {
    // ========== 普通候选项补全 ==========
    matchStrings = fuzzyMatch(candidates, cur)
}
```

**关键点**：布尔标志返回 `candidates`（候选项列表），而不是空列表。

#### 新增 `getFlagType` 函数

```go
// getFlagType 获取指定上下文中标志的类型
func getFlagType(root types.Command, context string, flagName string) (types.FlagType, bool) {
    cmd := findCommandByContext(root, context)
    if cmd == nil {
        return types.FlagTypeUnknown, false
    }

    // 首先尝试从命令的 Flags() 中查找
    flag := findFlagByName(cmd, flagName)
    if flag != nil {
        return flag.Type(), true
    }

    // 如果没有找到，检查是否是内置标志
    // 内置标志在解析时动态注册，但我们可以根据名称识别其类型
    return getBuiltinFlagType(flagName, context, cmd)
}
```

#### 新增 `getBuiltinFlagType` 函数

内置标志（`--help`、`--version`）在解析时动态注册，不在 `cmd.Flags()` 中。通过名称识别：

```go
// getBuiltinFlagType 根据标志名称识别内置标志的类型
func getBuiltinFlagType(flagName string, context string, cmd types.Command) (types.FlagType, bool) {
    // 移除 "-" 或 "--" 前缀
    name := strings.TrimPrefix(flagName, "--")
    if name == flagName {
        name = strings.TrimPrefix(flagName, "-")
    }

    // 帮助标志：所有命令都有
    if name == types.HelpFlagName || name == types.HelpFlagShortName {
        return types.FlagTypeBool, true
    }

    // 根命令特有的内置标志
    if context == "/" {
        config := cmd.Config()

        // 版本标志
        if config.Version != "" {
            if name == types.VersionFlagName || name == types.VersionFlagShortName {
                return types.FlagTypeBool, true
            }
        }

        // 补全标志
        if config.Completion {
            if name == types.CompletionFlagName {
                return types.FlagTypeEnum, true
            }
        }
    }

    return types.FlagTypeUnknown, false
}
```

### 2. 模板层实现

调整判断顺序，优先检查 `matches`：

#### Bash 模板

```bash
# 根据结果决定补全行为
# IS_FLAG 为 true 表示当前是枚举类型标志的值补全
if [[ "$is_flag" == "true" ]]; then
    # 枚举类型标志, 显示匹配结果
    if [[ -n "$matches" ]]; then
        read -ra COMPREPLY <<< "$matches"
    fi
elif [[ -n "$matches" ]]; then
    # 普通补全（包括布尔标志后的候选项补全）, 显示匹配结果
    read -ra COMPREPLY <<< "$matches"
elif [[ "$prev" =~ ^- ]]; then
    # 非枚举类型标志（如 String/Int 等）, 使用路径补全
    COMPREPLY=($(compgen -f -d -- "$cur"))
fi
```

#### PowerShell 模板

```powershell
# 根据结果决定补全行为
# IS_FLAG 为 true 表示当前是枚举类型标志的值补全
if ($isFlag -and $enumValues.Count -gt 0) {
    # 枚举类型标志, 返回匹配结果
    return $matchResults
}
elseif ($matchResults.Count -gt 0) {
    # 普通补全（包括布尔标志后的候选项补全）, 处理匹配结果
    $matchingOptions = [System.Collections.ArrayList]::new()
    $flagRegex = [regex]::new('^-')
    foreach ($match in $matchResults) {
        $result = if ($flagRegex.IsMatch($match)) { $match } else { "$match " }
        [void]$matchingOptions.Add($result)
    }
    return $matchingOptions.ToArray()
}
elseif ($prevElement -match '^-') {
    # 非枚举类型标志（如 String/Int 等）, 使用路径补全
    return Get-{{.SanitizedName}}PathCompletions -WordToComplete $wordToComplete
}
else {
    # 没有匹配结果, 返回空数组
    return @()
}
```

### 判断优先级

| 优先级 | 条件 | 补全行为 | 适用场景 |
|--------|------|----------|----------|
| 1 | `IS_FLAG=true` | 枚举值补全 | `--kind <Tab>` → `service pod deployment` |
| 2 | `MATCHES` 非空 | 候选项补全 | `--help <Tab>` → `--config --verbose ...` |
| 3 | `prev` 是标志 | 路径补全 | `--config <Tab>` → `file.txt` |

## 修复效果

| 场景 | 修复前 | 修复后 |
|------|--------|--------|
| `--help <Tab>` | 路径补全 ❌ | 候选项补全 ✅ |
| `--verbose <Tab>` | 路径补全 ❌ | 候选项补全 ✅ |
| `-h <Tab>` | 路径补全 ❌ | 候选项补全 ✅ |
| `--version <Tab>` | 路径补全 ❌ | 候选项补全 ✅ |
| `--kind <Tab>` | 枚举值补全 ✅ | 枚举值补全 ✅ |
| `--config <Tab>` | 路径补全 ✅ | 路径补全 ✅ |
| `--output <Tab>` | 路径补全 ✅ | 路径补全 ✅ |

## 关键设计决策

### 1. 为什么用 `FlagType` 而不是 `EnumValues` 判断？

原方案通过 `GetEnumValues` 是否返回空来判断，但：
- 布尔标志返回空（无枚举值）→ 误判为非枚举标志
- 需要通过 `FlagType` 明确知道是布尔类型

### 2. 为什么调整模板判断顺序？

原顺序：`IS_FLAG` → `prev` 是标志 → `matches` 非空

问题：布尔标志的 `IS_FLAG=false`，`prev` 是标志，直接走路径补全。

修复后：`IS_FLAG` → `matches` 非空 → `prev` 是标志

优势：布尔标志返回了 `matches`（候选项），优先匹配第二个条件。

### 3. 为什么需要特殊处理内置标志？

内置标志（`--help`、`--version`）在**解析时**动态注册，不是在**构建命令树时**注册。所以 `cmd.Flags()` 中找不到这些标志，需要通过名称匹配来识别。

## 代码文件变更

- `internal/completion/dynamic.go`: 修改 `handleAll`，新增 `getFlagType` 和 `getBuiltinFlagType`
- `internal/completion/templates/bash_dynamic.tmpl`: 调整判断顺序
- `internal/completion/templates/pwsh_dynamic.tmpl`: 调整判断顺序
