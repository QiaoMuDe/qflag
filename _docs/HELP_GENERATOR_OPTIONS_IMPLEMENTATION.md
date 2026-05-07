# 帮助信息生成器 - 选项和子命令实现方案

## 1. 概述

当前 `gen.go` 中的 `writeOptions` 和 `writeSubCmds` 函数为空, 需要实现这两个函数以生成格式化的选项和子命令帮助信息。

参考 `qflag\internal\help\writers.go` 的实现思路, 但适配新的 `types.Flag` 和 `types.Cmd` 接口。

## 2. 设计目标

- **自动对齐**: 计算选项/子命令名称的最大宽度, 确保描述信息对齐
- **格式美观**: 类似 `--flag, -f <type>` 的格式, 带有描述和默认值
- **中英文支持**: 根据 `UseChinese` 配置输出不同标题
- **空处理**: 无选项/子命令时不输出相应部分

## 3. 数据结构设计

### 3.1 选项信息结构

```go
type optionInfo struct {
    shortName  string // 短名称, 如 "v"
    longName   string // 长名称, 如 "verbose"
    namePart   string // 格式化后的名称部分, 如 "-v, --verbose"
    typ        string // 类型字符串, 如 "string"
    desc       string // 描述
    defValue   string // 默认值字符串
}
```

### 3.2 子命令信息结构

```go
type subCmdInfo struct {
    name string // 格式化后的名称, 如 "init, i"
    desc string // 描述
}
```

## 4. 输出格式设计

### 4.1 选项输出格式

```
Options:
  -v, --verbose <string>    输出详细日志信息 (default: false)
  -f, --file <string>       指定输入文件路径 (default: "")
  --output <string>         输出格式 (txt/json/yaml) (default: "txt")
```

**格式说明**: 
- `-s, --long <type>`: 短名称和长名称组合
- `--long <type>`: 仅长名称
- `-s <type>`: 仅短名称
- 描述后用括号显示默认值, 如 `(default: false)`
- 无默认值时不显示

### 4.2 子命令输出格式

```
Cmds:
  init, i           初始化项目结构
  build, b           编译项目
  run               运行程序
```

**格式说明**: 
- 2 空格缩进
- 名称右对齐到最大宽度
- 描述左对齐

## 5. 实现方案

### 5.1 writeOptions 实现步骤

1. 获取标志列表: `flags := cmd.Flags()`
2. 空检查: 无标志则返回
3. 写入标题: 根据 `UseChinese` 写入 `types.HelpOptionsCN` 或 `types.HelpOptionsEN`
4. 收集选项信息: 遍历标志, 构建 `optionInfo` 列表
5. 计算最大名称宽度
6. 遍历输出: 格式化每个选项

### 5.2 writeSubCmds 实现步骤

1. 获取子命令列表: `subCmds := cmd.SubCmds()`
2. 空检查: 无子命令则返回
3. 写入标题: 根据 `UseChinese` 写入 `types.HelpSubCmdsCN` 或 `types.HelpSubCmdsEN`
4. 收集子命令信息: 遍历子命令, 构建 `subCmdInfo` 列表
5. 计算最大名称长度
6. 遍历输出: 格式化每个子命令

### 5.3 辅助函数

#### calcOptionMaxWidth

```go
func calcOptionMaxWidth(options []optionInfo) int
```

计算选项名称部分的最大长度, 用于对齐。

计算规则: 
- `-s, --long <type>`: `len(shortName) + len(longName) + len(type) + 7`
  - 1(`-`) + 1(`,) + 1(空格) + 2(`--`) + 1(空格) + 3(`<>`)
- `--long <type>`: `len(longName) + len(type) + 5`
  - 2(`--`) + 1(空格) + 3(`<>`)
- `-s <type>`: `len(shortName) + len(type) + 4`
  - 1(`-`) + 1(空格) + 3(`<>`)

#### calcSubCmdMaxLen

```go
func calcSubCmdMaxLen(subCmds []subCmdInfo) int
```

计算子命令名称的最大长度。

### 5.4 格式化默认值的处理

Flag 接口的 `GetDefault()` 返回 `any`, 需要转换为字符串: 

```go
func formatDefaultValue(defValue any) string {
    if defValue == nil {
        return ""
    }
    switch v := defValue.(type) {
    case string:
        if v == "" {
            return `""`
        }
        return v
    case bool:
        return strconv.FormatBool(v)
    case int, int64, uint, uint64:
        return fmt.Sprintf("%d", v)
    case float64:
        return fmt.Sprintf("%.2f", v)
    case []string:
        return fmt.Sprintf("%v", v)
    case time.Duration:
        return v.String()
    default:
        return fmt.Sprintf("%v", v)
    }
}
```

## 6. 常量定义

在 `constants.go` 或 `gen.go` 中添加: 

```go
const (
    HelpOptionsCN   = "\n选项:\n"
    HelpOptionsEN   = "\nOptions:\n"
    HelpSubCmdsCN   = "\n子命令:\n"
    HelpSubCmdsEN   = "\nCmds:\n"
    DefaultPadding  = 4  // 描述与选项之间的最小间距
    OptionIndent    = 2   // 选项缩进
)
```

## 7. 输出示例

### 7.1 中文模式

```
选项:
  -v, --verbose <string>    输出详细日志信息 (default: false)
  -f, --file <string>       指定输入文件路径 (default: "")
      --level <int>         日志级别 (1-5) (default: 3)
```

### 7.2 英文模式

```
Options:
  -v, --verbose <string>    Output detailed log information (default: false)
  -f, --file <string>      Specify input file path (default: "")
      --level <int>        Log level (1-5) (default: 3)
```

## 8. 实现优先级

1. 添加常量定义
2. 实现 `formatDefaultValue` 辅助函数
3. 实现 `writeOptions` 函数
4. 实现 `writeSubCmds` 函数
5. 添加单元测试

## 9. 注意事项

- 使用 `strings.Builder` 而非 `bytes.Buffer`
- 保持与现有代码风格一致
- 考虑添加 `sync.Once` 缓存计算结果 (如需要) 
- 排序: 选项和子命令应按名称排序输出
