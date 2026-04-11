# 基于 FlagType 的补全方案

## 问题背景

当前补全逻辑无法正确处理布尔标志。当用户输入 `dynamic.exe --verbose <Tab>` 时：
- 当前逻辑认为 `prev = "--verbose"` 是标志，需要补全值
- 但布尔标志不需要值，存在即表示启用
- 结果：只能补全路径，无法补全其他标志/子命令

## 解决方案

通过 `flag.Type()` 判断标志类型，不同类型的标志采用不同的补全策略。

## FlagType 类型定义

```go
const (
    FlagTypeUnknown FlagType = iota // 未知类型
    FlagTypeString                   // 字符串标志
    FlagTypeInt                      // 整数标志
    FlagTypeInt64                    // 64位整数标志
    FlagTypeUint                     // 无符号整数标志
    FlagTypeUint8                    // 8位无符号整数标志
    FlagTypeUint16                   // 16位无符号整数标志
    FlagTypeUint32                   // 32位无符号整数标志
    FlagTypeUint64                   // 64位无符号整数标志
    FlagTypeFloat64                  // 64位浮点数标志
    FlagTypeBool                     // 布尔标志
    FlagTypeEnum                     // 枚举标志
    FlagTypeDuration                 // 持续时间标志
    FlagTypeTime                     // 时间标志
    FlagTypeSize                     // 大小标志
    FlagTypeMap                      // 映射标志
    FlagTypeStringSlice              // 字符串切片标志
    FlagTypeIntSlice                 // 整数切片标志
    FlagTypeInt64Slice               // 64位整数切片标志
)
```

## 补全策略矩阵

| 标志类型 | 需要值 | 补全行为 | 示例 |
|---------|--------|---------|------|
| `FlagTypeBool` | 否 | 补全其他标志/子命令 | `--verbose <Tab>` → `--config` |
| `FlagTypeEnum` | 是 | 补全枚举值 | `--kind <Tab>` → `service` |
| `FlagTypeString` | 是 | 路径补全 | `--config <Tab>` → `file.txt` |
| `FlagTypeInt` | 是 | 路径补全（数字输入） | `--port <Tab>` → 任意数字 |
| `FlagTypeDuration` | 是 | 路径补全 | `--timeout <Tab>` → `10s` |
| `FlagTypeSize` | 是 | 路径补全 | `--max-size <Tab>` → `1GB` |
| 其他数值类型 | 是 | 路径补全 | 同上 |
| 切片类型 | 是 | 路径补全 | `--tags <Tab>` → `value` |

## 核心实现

### 1. 新增辅助函数

```go
// getFlagType 获取指定上下文中标志的类型
//
// 参数:
//   - root: 根命令实例
//   - context: 上下文路径
//   - flagName: 标志名称
//
// 返回值:
//   - FlagType: 标志类型
//   - bool: 是否找到标志
func getFlagType(root types.Command, context string, flagName string) (types.FlagType, bool) {
    cmd := findCommandByContext(root, context)
    if cmd == nil {
        return types.FlagTypeUnknown, false
    }

    flag := findFlagByName(cmd, flagName)
    if flag == nil {
        return types.FlagTypeUnknown, false
    }

    return flag.Type(), true
}

// isFlagNeedValue 判断标志是否需要值
//
// 参数:
//   - flagType: 标志类型
//
// 返回值:
//   - bool: 是否需要值
//
// 说明:
//   - 只有布尔标志不需要值
//   - 其他所有类型都需要值
func isFlagNeedValue(flagType types.FlagType) bool {
    return flagType != types.FlagTypeBool
}
```

### 2. 修改 handleAll 函数

```go
func handleAll(root types.Command, args []string) error {
    // ... 前面的代码不变 ...

    // 3. 执行补全逻辑
    var matchStrings []string
    var enumValues []string

    // 判断是否是标志值补全上下文
    // 条件：prev 是待补全值的标志（以 - 开头，不是 --，不包含 =）
    isFlagValueCompletion := strings.HasPrefix(prev, "-") && 
                             prev != "--" && 
                             !strings.Contains(prev, "=")

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
                // 这些类型不是枚举类型，不需要获取 enumValues
            }
        }
    } else {
        // ========== 普通候选项补全 ==========
        matchStrings = fuzzyMatch(candidates, cur)
    }

    // 4. 输出结果
    fmt.Printf("CONTEXT:%s\n", context)
    fmt.Printf("CUR:%s\n", cur)
    fmt.Printf("PREV:%s\n", prev)
    fmt.Printf("CANDIDATES:%s\n", strings.Join(candidates, " "))
    fmt.Printf("ENUM:%s\n", strings.Join(enumValues, " "))
    fmt.Printf("MATCHES:%s\n", strings.Join(matchStrings, " "))
    fmt.Printf("IS_FLAG:%v\n", isFlagValueCompletion && len(enumValues) > 0)

    return nil
}

// fuzzyMatch 对候选列表进行模糊匹配
func fuzzyMatch(candidates []string, cur string) []string {
    if cur == "" {
        return candidates
    }
    matches := fuzzy.CompletePrefix(cur, candidates)
    result := make([]string, len(matches))
    for i, match := range matches {
        result[i] = match.Str
    }
    return result
}
```

### 3. 优化 GetEnumValues

```go
// GetEnumValues 获取指定标志的枚举值
//
// 参数:
//   - root: 根命令实例
//   - context: 上下文路径
//   - flagName: 标志名称
//
// 返回值:
//   - []string: 枚举值列表
//   - error: 错误信息
//
// 说明:
//   - 只对 FlagTypeEnum 类型返回有效值
//   - 其他类型返回空列表
func GetEnumValues(root types.Command, context string, flagName string) ([]string, error) {
    cmd := findCommandByContext(root, context)
    if cmd == nil {
        return nil, fmt.Errorf("无效的上下文: %s", context)
    }

    flag := findFlagByName(cmd, flagName)
    if flag == nil {
        return nil, fmt.Errorf("标志不存在: %s", flagName)
    }

    // 只有枚举类型才有枚举值
    if flag.Type() != types.FlagTypeEnum {
        return nil, nil
    }

    // 获取枚举值
    return flag.EnumValues(), nil
}
```

## 测试用例

| 输入 | 上下文 | 预期结果 |
|------|--------|---------|
| `dynamic --verbose <Tab>` | `/` | 补全 `--config`, `--kind` 等其他标志 |
| `dynamic --kind <Tab>` | `/` | 补全 `service`, `pod`, `deployment` |
| `dynamic --config <Tab>` | `/` | Shell 路径补全 |
| `dynamic --port <Tab>` | `/` | Shell 路径补全（数字输入） |
| `dynamic config --verbose <Tab>` | `/config/` | 补全 `get`, `set`, `list` 子命令 |

## 优势

1. **类型驱动**：根据标志类型决定补全行为，逻辑清晰
2. **统一处理**：所有标志类型都有明确的补全策略
3. **易于扩展**：新增标志类型时，只需添加对应的 case
4. **向后兼容**：不影响现有功能，只修复布尔标志问题

## 实施步骤

1. 新增 `getFlagType` 和 `isFlagNeedValue` 辅助函数
2. 修改 `handleAll` 函数，使用 switch 处理不同类型
3. 优化 `GetEnumValues`，只对枚举类型返回有效值
4. 添加单元测试覆盖各种标志类型
5. 编译测试并验证
