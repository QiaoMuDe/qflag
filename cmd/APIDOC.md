# Package cmd

cmd 命令行标志管理结构体，封装参数解析、长短标志互斥及帮助系统。

## Constants

```go
const (
    BashFunctionHeader = `#!/usr/bin/env bash

_%s() {
    local cur prev words cword context opts i arg
    COMPREPLY=()

    # 使用_get_comp_words_by_ref获取补全参数, 提高健壮性
    if [[ -z "${_get_comp_words_by_ref}" ]]; then
        # 兼容旧版本Bash补全环境
        words=("${COMP_WORDS[@]}")
        cword=$COMP_CWORD
    else
        _get_comp_words_by_ref -n =: cur prev words cword
    fi

    cur="${words[cword]}"
    prev="${words[cword-1]}"

    # 构建命令树结构
    declare -A cmd_tree
    cmd_tree[/]="%s"
%s

    # 查找当前命令上下文
    local context="/"
    local i
    for ((i=1; i < cword; i++)); do
        local arg="${words[i]}"
        if [[ -n "${cmd_tree[$context$arg/]}" ]]; then
            context="$context$arg/"
        fi
    done

    # 获取当前上下文可用选项
    opts="${cmd_tree[$context]}"
    # 添加-o filenames选项处理特殊字符和空格
    COMPREPLY=($(compgen -o filenames -W "${opts}" -- ${cur}))

    # 模糊匹配与纠错提示：当无精确匹配时，从所有选项中查找包含关键词的项
    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
        local all_opts=()
        # 使用循环安全收集所有选项，避免空格分割问题
        for path in "${!cmd_tree[@]}"; do
            for opt in ${cmd_tree[$path]}; do
                all_opts+=($opt)
            done
        done
        # 去重并生成补全结果
        local unique_opts=($(printf "%%s\n" "${all_opts[@]}" | sort -u))
        COMPREPLY=($(compgen -o filenames -W "${unique_opts[*]}" -- ${cur}))
    fi

    return 0
    }

complete -F _%s %s
`

    BashCommandTreeEntry = "\tcmd_tree[/%s/]=\"%s\"\n"

    PwshFunctionHeader = `Register-ArgumentCompleter -CommandName %s -ScriptBlock {
        param($wordToComplete, $commandAst, $cursorPosition, $commandName, $parameterName)
    
        # 标志参数需求数组(保留原始大小写)
        $flagParams = @(
%s      )
    
        # 构建命令树结构
        $cmdTree = @{
%s      }
    
        # 解析命令行参数获取当前上下文
        $context = ''
        $args = $commandAst.CommandElements | Select-Object -Skip 1 | ForEach-Object { $_.ToString() }
        $index = 0
        $count = $args.Count
    
        while ($index -lt $count) {
            $arg = $args[$index]
            # 使用大小写敏感匹配查找标志
            $paramInfo = $flagParams | Where-Object { $_.Name -ceq $arg } | Select-Object -First 1
            if ($paramInfo) {
                $paramType = $paramInfo.Type
                $index++
                
                # 根据参数类型决定是否跳过下一个参数
                if ($paramType -eq 'required' -or ($paramType -eq 'optional' -and $index -lt $count -and $args[$index] -notlike '-*')) {
                    $index++
                }
                continue
            }
    
            $nextContext = if ($context) { "$context.$arg" } else { $arg }
            if ($cmdTree.ContainsKey($nextContext)) {
                $context = $nextContext
                $index++
            } else {
                break
            }
        }
    
        # 获取当前上下文可用选项并过滤
        $options = @()
        if ($cmdTree.ContainsKey($context)) {
            $options = $cmdTree[$context] -split ' ' | Where-Object { $_ -like "$wordToComplete*" }
        }
    
        # 模糊匹配与纠错提示：当无精确匹配时，从所有选项中查找包含关键词的项
        if (-not $options) {
            # 递归收集所有层级的选项
            $allOptions = @()
            $cmdTree.Values | ForEach-Object { $allOptions += $_ -split ' ' }
            $options = $allOptions | Select-Object -Unique | Where-Object { $_ -like "*$wordToComplete*" }
        }
    
        $options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_) }
    }`

    PwshCommandTreeEntryRoot = "\t\t'' = '%s'\n"
    PwshCommandTreeEntry     = "\t\t'%s' = '%s'\n"
    PwshCommandTreeOption    = "\t\t@{ Name = '%s'; Type = '%s'}\n"
)
```

补全脚本模板常量

## Variables

```go
var ChineseTemplate = HelpTemplate{
    CmdName:              "名称: %s\n\n",
    UsagePrefix:          "用法: ",
    UsageSubCmd:          " [子命令]",
    UsageInfoWithOptions: " [选项]\n\n",
    UsageGlobalOptions:   " [全局选项]",
    CmdNameWithShort:     "名称: %s, %s\n\n",
    CmdDescription:       "描述: %s\n\n",
    OptionsHeader:        "选项:\n",
    Option1:              "  --%s, -%s %s",
    Option2:              "  --%s %s",
    Option3:              "  -%s %s",
    OptionDefault:        "%s%*s%s (默认值: %s)\n",
    SubCmdsHeader:        "\n子命令:\n",
    SubCmd:               "  %s\t%s\n",
    SubCmdWithShort:      "  %s, %s\t%s\n",
    NotesHeader:          "\n注意事项:\n",
    NoteItem:             "  %d、%s\n",
    DefaultNote:          "当长选项和短选项同时使用时，最后指定的选项将优先生效。",
    ExamplesHeader:       "\n示例:\n",
    ExampleItem:          "  %d、%s\n     %s\n",
}
```

中文模板实例

```go
var EnglishTemplate = HelpTemplate{
    CmdName:              "Name: %s\n\n",
    UsagePrefix:          "Usage: ",
    UsageSubCmd:          " [subcmd]",
    UsageInfoWithOptions: " [options]\n\n",
    UsageGlobalOptions:   " [global options]",
    CmdNameWithShort:     "Name: %s, %s\n\n",
    CmdDescription:       "Desc: %s\n\n",
    OptionsHeader:        "Options:\n",
    Option1:              "  --%s, -%s %s",
    Option2:              "  --%s %s",
    Option3:              "  -%s %s",
    OptionDefault:        "%s%*s%s (default: %s)\n",
    SubCmdsHeader:        "\nSubCmds:\n",
    SubCmd:               "  %s\t%s\n",
    SubCmdWithShort:      "  %s, %s\t%s\n",
    NotesHeader:          "\nNotes:\n",
    NoteItem:             "  %d. %s\n",
    DefaultNote:          "In the case where both long options and short options are used at the same time,\n the option specified last shall take precedence.",
    ExamplesHeader:       "\nExamples:\n",
    ExampleItem:          "  %d. %s\n     %s\n",
}
```

英文模板实例

## Types

### Cmd

```go
type Cmd struct {
    // 解析阶段钩子函数
    // 在标志解析完成后、子命令参数处理后调用
    //
    // 参数:
    //   - 当前命令实例
    //
    // 返回值:
    //   - error: 错误信息, 非nil时会中断解析流程
    //   - bool: 是否需要退出程序
    ParseHook func(*Cmd) (error, bool)
    // Has unexported fields.
}
```

Cmd 命令行标志管理结构体，封装参数解析、长短标志互斥及帮助系统。

### QCommandLine

```go
var QCommandLine *Cmd
```

QCommandLine 全局默认 Command 实例

### NewCmd

```go
func NewCmd(longName string, shortName string, errorHandling flag.ErrorHandling) *Cmd
```

NewCmd 创建新的命令实例

**参数:**

  * `longName` : 命令长名称
  * `shortName` : 命令短名称
  * `errorHandling` : 错误处理方式

**返回值:**

  * `*Cmd` : 新的命令实例指针

**errorHandling 可选值:**

  * `flag.ContinueOnError` : 解析标志时遇到错误继续解析, 并返回错误信息
  * `flag.ExitOnError` : 解析标志时遇到错误立即退出程序, 并返回错误信息
  * `flag.PanicOnError` : 解析标志时遇到错误立即触发 panic

### AddExample

```go
func (c *Cmd) AddExample(e ExampleInfo)
```

AddExample 为命令添加使用示例

**参数:**

  * `e` : 示例信息，包含 description 和 usage

### AddNote

```go
func (c *Cmd) AddNote(note string)
```

AddNote 添加备注信息到命令

### AddSubCmd

```go
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error
```

AddSubCmd 添加外部子命令到当前命令 支持批量添加多个子命令, 遇到错误时收集所有错误并返回

**参数:**

  * `subCmds` : 一个或多个子命令实例指针

**返回值:**

  * 错误信息, 如果所有子命令添加成功则返回 nil

### Arg

```go
func (c *Cmd) Arg(i int) string
```

Arg 获取指定索引的非标志参数

### Args

```go
func (c *Cmd) Args() []string
```

Args 获取非标志参数切片

### Bool

```go
func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag
```

Bool 添加布尔类型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明

**返回值:** 布尔标志对象指针

### BoolVar

```go
func (c *Cmd) BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)
```

BoolVar 绑定布尔类型标志到指针并内部注册 Flag 对象

**参数依次为:** 布尔标志指针、长标志名、短标志、默认值、帮助说明

### CmdExists

```go
func (c *Cmd) CmdExists(cmdName string) bool
```

CmdExists 检查子命令是否存在

**参数:**

  * `cmdName` : 子命令名称

**返回:**

  * `bool` : 子命令是否存在

### Duration

```go
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag
```

Duration 添加时间间隔类型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明

**返回值:** 时间间隔标志对象指针

### DurationVar

```go
func (c *Cmd) DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```

DurationVar 绑定时间间隔类型标志到指针并内部注册 Flag 对象

**参数依次为:** 时间间隔标志指针、长标志名、短标志、默认值、帮助说明

### Enum

```go
func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag
```

Enum 添加枚举类型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片

**返回值:** 枚举标志对象指针

### EnumVar

```go
func (c *Cmd) EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string)
```

EnumVar 绑定枚举类型标志到指针并内部注册 Flag 对象

**参数依次为:** 枚举标志指针、长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片

### FlagExists

```go
func (c *Cmd) FlagExists(name string) bool
```

FlagExists 检查指定名称的标志是否存在

### FlagRegistry

```go
func (c *Cmd) FlagRegistry() *flags.FlagRegistry
```

FlagRegistry 获取标志注册表的只读访问

**返回值:**

  * `*flags.FlagRegistry` : 标志注册表的只读访问

### Float64

```go
func (c *Cmd) Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag
```

Float64 添加浮点型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明

**返回值:** 浮点型标志对象指针

### Float64Var

```go
func (c *Cmd) Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)
```

Float64Var 绑定浮点型标志到指针并内部注册 Flag 对象

**参数依次为:** 浮点数标志指针、长标志名、短标志、默认值、帮助说明

### GetDescription

```go
func (c *Cmd) GetDescription() string
```

GetDescription 返回命令描述

### GetExamples

```go
func (c *Cmd) GetExamples() []ExampleInfo
```

GetExamples 获取所有使用示例 返回示例切片的副本，防止外部修改

### GetHelp

```go
func (c *Cmd) GetHelp() string
```

GetHelp 返回命令用法帮助信息

### GetLogoText

```go
func (c *Cmd) GetLogoText() string
```

GetLogoText 获取 logo 文本

### GetModuleHelps

```go
func (c *Cmd) GetModuleHelps() string
```

GetModuleHelps 获取自定义模块帮助信息

### GetNotes

```go
func (c *Cmd) GetNotes() []string
```

GetNotes 获取所有备注信息

### GetUsageSyntax

```go
func (c *Cmd) GetUsageSyntax() string
```

GetUsageSyntax 获取自定义命令用法

### GetUseChinese

```go
func (c *Cmd) GetUseChinese() bool
```

GetUseChinese 获取是否使用中文帮助信息

### GetVersion

```go
func (c *Cmd) GetVersion() string
```

GetVersion 获取版本信息

### IP4

```go
func (c *Cmd) IP4(longName, shortName string, defValue string, usage string) *flags.IP4Flag
```

IP4 添加 IPv4 地址类型标志, 返回标志对象指针 参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: IPv4 地址标志对象指针

### IP4Var

```go
func (c *Cmd) IP4Var(f *flags.IP4Flag, longName, shortName string, defValue string, usage string)
```

IP4Var 绑定 IPv4 地址类型标志到指针并内部注册 Flag 对象 参数依次为: IPv4 标志指针、长标志名、短标志、默认值、帮助说明

### IP6

```go
func (c *Cmd) IP6(longName, shortName string, defValue string, usage string) *flags.IP6Flag
```

IP6 添加 IPv6 地址类型标志, 返回标志对象指针 参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: IPv6 地址标志对象指针

### IP6Var

```go
func (c *Cmd) IP6Var(f *flags.IP6Flag, longName, shortName string, defValue string, usage string)
```

IP6Var 绑定 IPv6 地址类型标志到指针并内部注册 Flag 对象 参数依次为: IPv6 标志指针、长标志名、短标志、默认值、帮助说明

### Int

```go
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *flags.IntFlag
```

Int 添加整数类型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明 返回值: 整数标志对象指针

### Int64

```go
func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag
```

Int64 添加 64 位整数类型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明

**返回值:** 64 位整数标志对象指针

### Int64Var

```go
func (c *Cmd) Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)
```

Int64Var 绑定 64 位整数类型标志到指针并内部注册 Flag 对象

**参数依次为:** 64 位整数标志指针、长标志名、短标志、默认值、帮助说明

### IntVar

```go
func (c *Cmd) IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)
```

IntVar 绑定整数类型标志到指针并内部注册 Flag 对象

**参数依次为:** 整数标志指针、长标志名、短标志、默认值、帮助说明

### IsParsed

```go
func (c *Cmd) IsParsed() bool
```

IsParsed 检查命令是否已完成解析

**返回值:**

  * `bool` : 解析状态, `true` 表示已解析 (无论成功失败), `false` 表示未解析

### LoadHelp

```go
func (c *Cmd) LoadHelp(filePath string) error
```

LoadHelp 从指定文件加载帮助信息

**参数:** `filePath` : 帮助信息文件路径

**返回值:** `error` : 如果文件不存在或读取文件失败，则返回错误信息

### LongName

```go
func (c *Cmd) LongName() string
```

LongName 返回命令长名称

### Map

```go
func (c *Cmd) Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag
```

Map 添加键值对类型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明 返回值: 键值对标志对象指针

### MapVar

```go
func (c *Cmd) MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)
```

MapVar 绑定键值对类型标志到指针并内部注册 Flag 对象

**参数依次为:** 键值对标志指针、长标志名、短标志、默认值、帮助说明

### NArg

```go
func (c *Cmd) NArg() int
```

NArg 获取非标志参数的数量

### NFlag

```go
func (c *Cmd) NFlag() int
```

NFlag 获取标志的数量

### Name

```go
func (c *Cmd) Name() string
```

Name 获取命令名称

**返回值:** - 优先返回长名称, 如果长名称不存在则返回短名称

### Parse

```go
func (c *Cmd) Parse(args []string) (err error)
```

Parse 完整解析命令行参数 (含子命令处理)

**主要功能：**

  1. 解析当前命令的长短标志及内置标志
  2. 自动检测并解析子命令及其参数 (若存在)
  3. 验证枚举类型标志的有效性

**参数：**

  * `args` : 原始命令行参数切片 (包含可能的子命令及参数)

**返回值：**

  * 解析过程中遇到的错误 (如标志格式错误、子命令解析失败等)

**注意事项：**

  * 每个 Cmd 实例仅会被解析一次 (线程安全)
  * 若检测到子命令, 会将剩余参数传递给子命令的 Parse 方法
  * 处理内置标志执行逻辑

### ParseFlagsOnly

```go
func (c *Cmd) ParseFlagsOnly(args []string) (err error)
```

ParseFlagsOnly 仅解析当前命令的标志参数 (忽略子命令)

**主要功能：**

  1. 解析当前命令的长短标志及内置标志
  2. 验证枚举类型标志的有效性
  3. 明确忽略所有子命令及后续参数

**参数：**

  * `args` : 原始命令行参数切片 (子命令及后续参数会被忽略)

**返回值：**

  * 解析过程中遇到的错误 (如标志格式错误等)

**注意事项：**

  * 每个 Cmd 实例仅会被解析一次 (线程安全)
  * 不会处理任何子命令, 所有参数均视为当前命令的标志或位置参数
  * 处理内置标志逻辑

### Path

```go
func (c *Cmd) Path(longName, shortName string, defValue string, usage string) *flags.PathFlag
```

Path 添加路径类型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明 返回值: 路径标志对象指针

### PathVar

```go
func (c *Cmd) PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)
```

PathVar 绑定路径类型标志到指针并内部注册 Flag 对象

**参数依次为:** 路径标志指针、长标志名、短标志、默认值、帮助说明

### PrintHelp

```go
func (c *Cmd) PrintHelp()
```

PrintHelp 打印命令的帮助信息, 优先打印用户的帮助信息, 否则自动生成帮助信息

**注意:**

  * 打印帮助信息时, 不会自动退出程序

### SetDescription

```go
func (c *Cmd) SetDescription(desc string)
```

SetDescription 设置命令描述

### SetEnableCompletion

```go
func (c *Cmd) SetEnableCompletion(enable bool)
```

SetEnableCompletion 设置是否启用自动补全, 只能在根命令上启用

**参数:**

  * `enable` : `true` 表示启用补全, `false` 表示禁用

### SetExitOnBuiltinFlags

```go
func (c *Cmd) SetExitOnBuiltinFlags(exit bool)
```

SetExitOnBuiltinFlags 设置是否在解析内置参数时退出 默认情况下为 `true`, 当解析到内置参数时, QFlag 将退出程序

**参数:**

  * `exit` : 是否退出

### SetHelp

```go
func (c *Cmd) SetHelp(help string)
```

SetHelp 设置用户自定义命令帮助信息

### SetLogoText

```go
func (c *Cmd) SetLogoText(logoText string)
```

SetLogoText 设置 logo 文本

### SetModuleHelps

```go
func (c *Cmd) SetModuleHelps(moduleHelps string)
```

SetModuleHelps 设置自定义模块帮助信息

### SetUsageSyntax

```go
func (c *Cmd) SetUsageSyntax(usageSyntax string)
```

SetUsageSyntax 设置自定义命令用法

### SetUseChinese

```go
func (c *Cmd) SetUseChinese(useChinese bool)
```

SetUseChinese 设置是否使用中文帮助信息

### SetVersion

```go
func (c *Cmd) SetVersion(version string)
```

SetVersion 设置版本信息

### ShortName

```go
func (c *Cmd) ShortName() string
```

ShortName 返回命令短名称

### Slice

```go
func (c *Cmd) Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag
```

Slice 绑定字符串切片类型标志并内部注册 Flag 对象

**参数依次为:** 长标志名、短标志、默认值、帮助说明

**返回值:** 字符串切片标志对象指针

### SliceVar

```go
func (c *Cmd) SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)
```

SliceVar 绑定字符串切片类型标志到指针并内部注册 Flag 对象

**参数依次为:** 字符串切片标志指针、长标志名、短标志、默认值、帮助说明

### String

```go
func (c *Cmd) String(longName, shortName, defValue, usage string) *flags.StringFlag
```

String 添加字符串类型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明

**返回值:** 字符串标志对象指针

### StringVar

```go
func (c *Cmd) StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)
```

StringVar 绑定字符串类型标志到指针并内部注册 Flag 对象

**参数依次为:** 字符串标志指针、长标志名、短标志、默认值、帮助说明

### SubCmdMap

```go
func (c *Cmd) SubCmdMap() map[string]*Cmd
```

SubCmdMap 返回子命令映射表

### SubCmds

```go
func (c *Cmd) SubCmds() []*Cmd
```

SubCmds 返回子命令切片

### Time

```go
func (c *Cmd) Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag
```

Time 添加时间类型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明

**返回值:** 时间标志对象指针

### TimeVar

```go
func (c *Cmd) TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)
```

TimeVar 绑定时间类型标志到指针并内部注册 Flag 对象

**参数依次为:** 时间标志指针、长标志名、短标志、默认值、帮助说明

### URL

```go
func (c *Cmd) URL(longName, shortName string, defValue string, usage string) *flags.URLFlag
```

URL 添加 URL 类型标志, 返回标志对象指针 参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: URL 标志对象指针

### URLVar

```go
func (c *Cmd) URLVar(f *flags.URLFlag, longName, shortName string, defValue string, usage string)
```

URLVar 绑定 URL 类型标志到指针并内部注册 Flag 对象 参数依次为: URL 标志指针、长标志名、短标志、默认值、帮助说明

### Uint16

```go
func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag
```

Uint16 添加 16 位无符号整数类型标志, 返回标志对象指针

**参数依次为:** 长标志名、短标志、默认值、帮助说明

**返回值:** 16 位无符号整数标志对象指针

### Uint16Var

```go
func (c *Cmd) Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)
```

Uint16Var 绑定 16 位无符号整数类型标志到指针并内部注册 Flag 对象

**参数依次为:** 16 位无符号整数标志指针、长标志名、短标志、默认值、帮助说明

### Uint32

```go
func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag
```

Uint32 添加 32 位无符号整数类型标志, 返回标志对象指针 参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: 32 位无符号整数标志对象指针

### Uint32Var

```go
func (c *Cmd) Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string)
```

Uint32Var 绑定 32 位无符号整数类型标志到指针并内部注册 Flag 对象 参数依次为: 32 位无符号整数标志指针、长标志名、短标志、默认值、帮助说明

### Uint64

```go
func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag
```

Uint64 添加 64 位无符号整数类型标志, 返回标志对象指针 参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: 64 位无符号整数标志对象指针

### Uint64Var

```go
func (c *Cmd) Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string)
```

Uint64Var 绑定 64 位无符号整数类型标志到指针并内部注册 Flag 对象 参数依次为: 64 位无符号整数标志指针、长标志名、短标志、默认值、帮助说明

### ExampleInfo

```go
type ExampleInfo struct {
    Description string // 示例描述
    Usage       string // 示例使用方式
}
```

ExampleInfo 示例信息结构体 用于存储命令的使用示例，包括描述和示例内容

### FlagParam

```go
type FlagParam struct {
    Name string // 标志名称 (保留原始大小写)
    Type string // 参数需求类型: "required"|"optional"|"none"
}
```

FlagParam 表示标志参数及其需求类型

### HelpTemplate

```go
type HelpTemplate struct {
    CmdName              string // 命令名称模板
    CmdNameWithShort     string // 命令名称带短名称模板
    CmdDescription       string // 命令描述模板
    UsagePrefix          string // 用法说明前缀模板
    UsageSubCmd          string // 用法说明子命令模板
    UsageInfoWithOptions string // 带选项的用法说明信息模板
    UsageGlobalOptions   string // 全局选项部分
    OptionsHeader        string // 选项头部模板
    Option1              string // 选项模板 (带短选项)
    Option2              string // 选项模板 (无短选项)
    Option3              string // 选项模板 (无长选项)
    OptionDefault        string // 选项模板的默认值
    SubCmdsHeader        string // 子命令头部模板
    SubCmd               string // 子命令模板
    SubCmdWithShort      string // 子命令带短名称模板
    NotesHeader          string // 注意事项头部模板
    NoteItem             string // 注意事项项模板
    DefaultNote          string // 默认注意事项
    ExamplesHeader       string // 示例信息头部模板
    ExampleItem          string // 示例信息项模板
}
```

HelpTemplate 帮助信息模板结构体