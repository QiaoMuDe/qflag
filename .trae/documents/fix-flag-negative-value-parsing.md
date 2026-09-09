# 修复：标志值为负号开头时被误判为未知标志

## Summary（概述）

修复 qflag 解析器中「标志的取值以 `-` 开头（如负数 `-8080`、`-1.5`、`-foo`）时被误认为未知标志并报错」的缺陷。根因在于 [suggestion.go](file:///d:/资源池/下水道/Dev/本地项目/qflag/internal/parser/suggestion.go#L175-L218) 的 `checkUnknownFlags` 预扫描逻辑：它把任意以 `-` 开头的 token 都当作独立标志名，却没有维护「前一个取值标志已经消费了该 token 作为其值」的状态。

修复后补齐常见、冷门、边界三类标志解析测试用例，覆盖全部标志类型与特殊命令行形式。

## Current State Analysis（现状分析）

### 缺陷触发链路
[parser.go](file:///d:/资源池/下水道/Dev/本地项目/qflag/internal/parser/parser.go#L105-L108) 在调用标准库 `flag.FlagSet.Parse()` **之前**，先执行 `checkUnknownFlags(cmd, args)` 做未知标志预扫描。

`checkUnknownFlags` 的核心逻辑（[suggestion.go](file:///d:/资源池/下水道/Dev/本地项目/qflag/internal/parser/suggestion.go#L189-L215)）：

```go
for i := 0; i < len(args); i++ {
    arg := args[i]
    ...
    if !strings.HasPrefix(arg, "-") { ...continue }
    // 处理 --flag=value 格式
    flagName := arg
    if idx := strings.Index(arg, "="); idx != -1 { flagName = arg[:idx] }
    if !registeredFlags[flagName] { return newUnknownFlagError(cmd, flagName) }
}
```

现象：`mycmd --port -8080`
- 扫描到 `--port` → 命中已注册标志，通过；
- 扫描到 `-8080` → 以 `-` 开头，被当作标志名 `flagName="-8080"`；
- `-8080` 不在 `registeredFlags`（`map[string]bool`）中 → 直接返回 `UnknownFlagError`。

`--flag=value` 形式正常（值被 `=` 剥掉只比对 `--flag`）；空格分隔的负值则必现该 bug。

### 标准库对照
[parser_register.go](file:///d:/资源池/下水道/Dev/本地项目/qflag/internal/parser/parser_register.go#L96-L98) 中 wrapper 已实现 `IsBoolFlag()`（`Type() == types.FlagTypeBool`）。标准库 `flag` 的 `parseOne` 对非布尔标志会无条件把下一个参数作为其值（不管是否以 `-` 开头）。因此**只需让预扫描放行，真实 `Parse` 必然成功**，本次修复是与标准库口径对齐，而非引入新规则。

### 可用信息（无需新增 API）
- `types.Flag.Type()` 暴露 `FlagType`（[flag.go](file:///d:/资源池/下水道/Dev/本地项目/qflag/internal/types/flag.go#L288-L297)），`FlagTypeBool` 用于判别布尔；
- `flagValueWrapper.IsBoolFlag()` 已用同一判定（[parser_register.go](file:///d:/资源池/下水道/Dev/本地项目/qflag/internal/parser/parser_register.go#L96-L98)）。

### 现有测试情况
- `internal/flag/*_test.go`：直接对 `Set()` 做单元测试（数值边界、验证器等），**已覆盖单值类型的解析边界**，本次无需重复。
- 根 `parser_test.go`：`package qflag`，用 `cmd.NewCmd` + `c.Parse` 做集成解析。
- `internal/parser/*_test.go`：用 `mock.NewMockCommand` + `mock.NewMockFlag` + `parser.ParseOnly` 直接测解析器内部逻辑。
- **当前没有针对 `checkUnknownFlags` 预扫描的测试文件**（`internal/parser` 无 `suggestion_test.go`）。

## Proposed Changes（改动方案）

### 改动 1：修复 `checkUnknownFlags`（核心）

文件：[suggestion.go](file:///d:/资源池/下水道/Dev/本地项目/qflag/internal/parser/suggestion.go#L175-L218)

- 把 `registeredFlags` 从 `map[string]bool` 改为 `map[string]types.Flag`（拿到实例以调 `Type()`）。
- 新增局部状态 `skipNextAsValue bool`：当命中「非布尔标志且未内联 `=` 取值」时置 `true`，下一个循环迭代在其处理各种分支跳转动作**之前**直接 `continue` 跳过（该 token 是当前标志被消费掉的值）。
- 保持现有 `--` 终止符 break、子命令 break、位置参数 continue 三条分支顺序与小写逻辑不变。

伪代码（已在分析阶段充分验证，实现时严格按此）：

```go
func checkUnknownFlags(cmd types.Command, args []string) error {
    // 已注册标志名 -> 标志实例的映射（用于判断是否取值）
    registeredFlags := make(map[string]types.Flag)
    for _, f := range cmd.FlagRegistry().List() {
        if l := f.LongName(); l != "" {
            registeredFlags["--"+l] = f
            registeredFlags["-"+l] = f // 单横杠长名称
        }
        if s := f.ShortName(); s != "" {
            registeredFlags["-"+s] = f
        }
    }

    // 标记下一个 token 是否被上一个取值标志消费为「值」
    skipNextAsValue := false

    for i := 0; i < len(args); i++ {
        arg := args[i]

        // 上一个非布尔取值标志已消费该 token 作为值，跳过不检查
        if skipNextAsValue {
            skipNextAsValue = false
            continue
        }

        // 遇到 -- 停止扫描，后面都视为位置参数
        if arg == "--" {
            break
        }

        // 不是标志格式，检查是否为子命令
        if !strings.HasPrefix(arg, "-") {
            if _, isSubCmd := cmd.CmdRegistry().Get(arg); isSubCmd {
                break
            }
            continue
        }

        // 处理 --flag=value 格式
        flagName := arg
        hasInlineValue := false
        if idx := strings.Index(arg, "="); idx != -1 {
            flagName = arg[:idx]
            hasInlineValue = true
        }

        f, ok := registeredFlags[flagName]
        if !ok {
            return newUnknownFlagError(cmd, flagName)
        }

        // 非布尔标志且未内联取值 → 下一个 token 是它的值，需跳过
        if f.Type() != types.FlagTypeBool && !hasInlineValue {
            skipNextAsValue = true
        }
    }

    return nil
}
```

### 改动 2：新增 `internal/parser/suggestion_test.go`（白盒，主测试）

用 `mock.NewMockCommand` + `mock.NewMockFlag` / `mock.NewMockBoolFlag` + 直接调用 `checkUnknownFlags(cmd, args)`（同 `parser_validation_test.go` 的构造方式），表驱动断言 `nil` 或 `*types.UnknownFlagError`。

测试用例按「常见 / 冷门 / 边界」组织：

**A. 常见（缺陷回归 + 常规）**
| 用例 | 期望 |
|------|------|
| `--port -8080`（int 空格负值） | nil（修复目标） |
| `--ratio -1.5`（float 空格负值） | nil |
| `--name -foo`（string 负号开头值） | nil |
| `--port 8080`（普通正值） | nil |
| `--name=hello`（`=` 形式） | nil |
| `--verbose`（布尔单独使用） | nil |

**B. 冷门 / 少见形式**
| 用例 | 期望 |
|------|------|
| `--port=-8080`（负值内联） | nil |
| `--name=-foo`（负号开头内联） | nil |
| `-p -8080`（短名称取负值） | nil |
| `--num -9223372036854775808`（int64 极值负值） | nil |
| `--filter -[a-z]+`（特殊字符串值） | nil |
| `--vals -1,-2`（slice 负值，Serve 后实际在两段解析验证） | nil |
| `--ratio -Inf`（float 特殊值） | nil |

**C. 边界 / 行为保持**
| 用例 | 期望 |
|------|------|
| `--verbose -x`（布尔后跟未知标志） | 报错 `-x`（布尔不消费，行为不变） |
| `--port -8080 --unknown`（消费负值后仍发现未知标志） | 报错 `--unknown` |
| `--port --verbose`（取值标志后跟另一合法标志充当值） | nil（`--verbose` 被当作 port 的值消费） |
| `-- --port`（`--` 终止符后按位置参数） | nil（不检查） |
| `--port 8080 sub`（子命令名中止扫描；mock 需有名为 sub 的子命令） | nil |
| `--nope`（纯未知标志） | 报错 `--nope`（回归保护） |

### 改动 3：扩展集成解析测试（真实标志类型，验证最终取值）

在根 `parser_test.go`（`package qflag`，沿用 `cmd.NewCmd` + `c.Parse` 模式）追加或新增测试函数，end-to-end 断言**解析后的实际值**（`Get()`），确保不是「只放行但没正确赋值」：

- `TestParse_NegativeValues`：`--count -8080` → `intFlag.Get() == -8080`；`--ratio -1.5` → `floatFlag.Get() == -1.5`；`--name -foo` → `stringFlag.Get() == "-foo"`。
- `TestParse_InlineNegative`：`--count=-8080`、`--ratio=-1.5` → 对应值。
- `TestParse_BoolFollowedByUnknown`：`--verbose -x` → `errors.As(err, &types.UnknownFlagError{})` 为真。
- `TestParse_BoundaryValues`：int/int64 极值（最大/最小）、uint 负数（报错）、float `NaN/+Inf/-Inf`、slice 含负数成员、duration 负值 `--timeout -1m30s`（可选，验证 size/duration 若支持负值则在集成层覆盖）。

> 注：`internal/flag/*_test.go` 已单测各类型 `Set()` 的数值边界（溢出、负数非法等），集成层只做「命令行空格/内联/布尔 三形态」的值贯通验证，避免重复。

## Assumptions & Decisions（假设与决策）

1. **布尔标志不消费下一个参数**：与标准库 `flag` 行为一致。布尔后用 `-x` 仍按未知标志报错。
2. **非布尔取值标志（含数字 `-` 值）**：始终消费紧随的下一个 token 作为值（即便以 `-` 开头），与标准库 `parseOne` 一致。
3. `--flag=`（空值内联）视为已内联取值，不启动 skip；空字符串即该标志的值。
4. 值为最后一项的取值标志（如仅 `--port`）：预扫描放行，由真实 `flag.Parse` 报「flag needs an argument」，不做额外改动（行为与标准库一致，非本次目标）。
5. mock 标志的 `Set` 不做类型校验，故冷门取值类用例在「白盒层」只验证「预扫描是否放行」，实际数值贯通在「集成层」用真实标志验证。
6. 测试文件归属：白盒放 `internal/parser`，集成放根 `parser_test.go`，遵循项目既有分层。

## Verification（验证步骤）

1. `go build ./...` — 编译通过。
2. `go vet ./...` — 静态检查通过。
3. `go test ./...` — 全部既有测试 + 新增测试通过（无回归）。
4. 针对改动 2/3 单独跑：
   - `go test ./internal/parser -run TestCheckUnknownFlags` / 新增测试名
   - `go test . -run TestParse_NegativeValues`
5. 手动冒烟（可选）：`app --port -8080` 不再报「-8080 不存在」，`Get()` 得到 `-8080`；`app --verbose -bogus` 仍报未知标志错误。