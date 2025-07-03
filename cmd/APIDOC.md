# Package cmd

Package cmd 提供命令行标志管理、参数解析和帮助系统。

```go
package cmd // import "gitee.com/MM-Q/qflag/cmd"
```

## 变量

### 中文模板实例

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

### 英文模板实例

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

### 保持兼容 API

```go
var NewCmd = NewCommand
```

## 函数

### GetExecutablePath

获取程序的绝对安装路径。如果无法通过 `os.Executable` 获取路径，则使用 `os.Args[0]` 作为替代。

```go
func GetExecutablePath() string
```

### NewCommand

创建新的命令实例。

```go
func NewCommand(longName string, shortName string, errorHandling flag.ErrorHandling) *Cmd
```

## 类型

### Cmd

命令行标志管理结构体，封装参数解析、长短标志互斥及帮助系统。

```go
type Cmd struct {
    // Has unexported fields.
}
```

### CmdInterface

命令接口定义，封装命令行程序的核心功能，提供统一的命令管理、参数解析和帮助系统。实现类需保证线程安全，所有方法应支持并发调用。

```go
type CmdInterface interface {
    // 元数据操作方法
    LongName() string                         // 获取命令名称(长名称)，如"app"
    ShortName() string                        // 获取命令短名称，如"a"
    GetDescription() string                   // 获取命令描述信息
    SetDescription(desc string)               // 设置命令描述信息，用于帮助输出
    GetHelp() string                          // 获取自定义帮助信息
    SetHelp(help string)                      // 设置自定义帮助信息，覆盖自动生成内容
    LoadHelp(filePath string) error           // 加载自定义帮助信息，从文件中读取
    SetUsageSyntax(usageSyntax string)        // 设置自定义用法说明，覆盖自动生成内容
    GetUsageSyntax() string                   // 获取自定义用法说明
    GetUseChinese() bool                      // 获取是否使用中文帮助信息
    SetUseChinese(useChinese bool)            // 设置是否使用中文帮助信息
    AddSubCmd(subCmd *Cmd)                    // 添加子命令，子命令会继承父命令的上下文
    SubCmds() []*Cmd                          // 获取所有已注册的子命令列表
    Parse(args []string) error                // 解析命令行参数，自动处理标志和子命令
    ParseFlagsOnly(args []string) (err error) // 仅解析标志参数，不处理子命令
    Args() []string                           // 获取所有非标志参数(未绑定到任何标志的参数)
    Arg(i int) string                         // 获取指定索引的非标志参数，索引越界返回空字符串
    NArg() int                                // 获取非标志参数的数量
    NFlag() int                               // 获取已解析的标志数量
    FlagExists(name string) bool              // 检查指定名称的标志是否存在(支持长/短名称)
    PrintHelp()                               // 打印命令帮助信息
    AddNote(note string)                      // 添加备注信息
    GetNotes() []string                       // 获取所有备注信息
    AddExample(e ExampleInfo)                 // 添加示例信息
    GetExamples() []ExampleInfo               // 获取所有示例信息
    SetVersion(version string)                // 设置版本信息
    GetVersion() string                       // 获取版本信息
    SetLogoText(logoText string)              // 设置logo文本
    GetLogoText() string                      // 获取logo文本
    SetModuleHelps(moduleHelps string)        // 设置自定义模块帮助信息
    GetModuleHelps() string                   // 获取自定义模块帮助信息

    // 添加标志方法
    String(longName, shortName, usage, defValue string) *flags.StringFlag                             // 添加字符串类型标志
    Int(longName, shortName, usage string, defValue int) *flags.IntFlag                               // 添加整数类型标志
    Int64(longName, shortName, usage string, defValue int64) *flags.Int64Flag                         // 添加64位整数类型标志
    Bool(longName, shortName, usage string, defValue bool) *flags.BoolFlag                            // 添加布尔类型标志
    Float64(longName, shortName, usage string, defValue float64) *flags.Float64Flag                   // 添加浮点数类型标志
    Duration(longName, shortName, usage string, defValue time.Duration) *flags.DurationFlag           // 添加时间间隔类型标志
    Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag // 添加枚举类型标志
    Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag               // 添加字符串切片类型标志

    Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag       // 添加时间类型标志
    Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag // 添加Map标志
    Path(longName, shortName string, defValue string, usage string) *flags.PathFlag          // 添加路径标志

    // 绑定标志方法
    StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)                             // 绑定字符串标志到指定变量
    IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)                        // 绑定整数标志到指定变量
    Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)                  // 绑定64位整数标志到指定变量
    BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)                     // 绑定布尔标志到指定变量
    Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)            // 绑定浮点数标志到指定变量
    DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)    // 绑定时间间隔类型标志到指定变量
    EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string) // 绑定枚举标志到指定变量
    SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)               // 绑定字符串切片标志到指定变量
    Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)               // 绑定无符号16位整型标志到指定变量
    TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)                // 绑定时间类型标志到指定变量
    MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)          // 绑定字符串映射标志到指定变量
    PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string) // 绑定路径标志到指定变量
    // Has unexported methods.
}
```

### ExampleInfo

示例信息结构体，用于存储命令的使用示例，包括描述和示例内容。

```go
type ExampleInfo struct {
    Description string // 示例描述
    Usage       string // 示例使用方式
}
```

### HelpTemplate

帮助信息模板结构体。

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
    Option1              string // 选项模板(带短选项)
    Option2              string // 选项模板(无短选项)
    Option3              string // 选项模板(无长选项)
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

### UserInfo

存储用户自定义信息的嵌套结构体。

```go
type UserInfo struct {
    // Has unexported fields.
}
```