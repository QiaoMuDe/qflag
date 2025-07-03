package cmd // import "gitee.com/MM-Q/qflag/cmd"


VARIABLES

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
    中文模板实例

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
    英文模板实例

var NewCmd = NewCommand
    保持兼容API 支持 NewCmd 别名


FUNCTIONS

func GetExecutablePath() string
    GetExecutablePath 获取程序的绝对安装路径 如果无法通过 os.Executable 获取路径,则使用 os.Args[0] 作为替代
    返回：程序的绝对路径字符串


TYPES

type Cmd struct {
	// Has unexported fields.
}
    Cmd 命令行标志管理结构体,封装参数解析、长短标志互斥及帮助系统。

var QCommandLine *Cmd
    QCommandLine 全局默认Command实例

func NewCommand(longName string, shortName string, errorHandling flag.ErrorHandling) *Cmd
    NewCommand 创建新的命令实例 参数: longName: 命令长名称 shortName:
    命令短名称 errorHandling: 错误处理方式 返回值: *Cmd命令实例指针 errorHandling可选值:
    flag.ContinueOnError、flag.ExitOnError、flag.PanicOnError

func (c *Cmd) AddExample(e ExampleInfo)
    AddExample 为命令添加使用示例 description: 示例描述 usage: 示例使用方式

func (c *Cmd) AddNote(note string)
    AddNote 添加备注信息到命令

func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error
    AddSubCmd 关联一个或多个子命令到当前命令 支持批量添加多个子命令，遇到错误时收集所有错误并返回 参数:

        subCmds: 一个或多个子命令实例指针

    返回值:

        错误信息列表, 如果所有子命令添加成功则返回nil

func (c *Cmd) Arg(i int) string
    Arg 获取指定索引的非标志参数

func (c *Cmd) Args() []string
    Args 获取非标志参数切片

func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag
    Bool 添加布尔类型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明

    返回值: 布尔标志对象指针

func (c *Cmd) BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)
    BoolVar 绑定布尔类型标志到指针并内部注册Flag对象

    参数依次为: 布尔标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag
    Duration 添加时间间隔类型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明

    返回值: 时间间隔标志对象指针

func (c *Cmd) DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)
    DurationVar 绑定时间间隔类型标志到指针并内部注册Flag对象

    参数依次为: 时间间隔标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag
    Enum 添加枚举类型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片

    返回值: 枚举标志对象指针

func (c *Cmd) EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string)
    EnumVar 绑定枚举类型标志到指针并内部注册Flag对象

    参数依次为: 枚举标志指针、长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片

func (c *Cmd) FlagExists(name string) bool
    FlagExists 检查指定名称的标志是否存在

func (c *Cmd) Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag
    Float64 添加浮点型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明

    返回值: 浮点型标志对象指针

func (c *Cmd) Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)
    Float64Var 绑定浮点型标志到指针并内部注册Flag对象

    参数依次为: 浮点数标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) GetDescription() string
    GetDescription 返回命令描述

func (c *Cmd) GetExamples() []ExampleInfo
    GetExamples 获取所有使用示例 返回示例切片的副本，防止外部修改

func (c *Cmd) GetHelp() string
    GetHelp 返回命令用法帮助信息

func (c *Cmd) GetLogoText() string
    GetLogoText 获取logo文本

func (c *Cmd) GetModuleHelps() string
    GetModuleHelps 获取自定义模块帮助信息

func (c *Cmd) GetNotes() []string
    GetNotes 获取所有备注信息

func (c *Cmd) GetUsageSyntax() string
    GetUsageSyntax 获取自定义命令用法

func (c *Cmd) GetUseChinese() bool
    GetUseChinese 获取是否使用中文帮助信息

func (c *Cmd) GetVersion() string
    GetVersion 获取版本信息

func (c *Cmd) IP4(longName, shortName string, defValue string, usage string) *flags.IP4Flag
    IP4 添加IPv4地址类型标志, 返回标志对象指针 参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: IPv4地址标志对象指针

func (c *Cmd) IP4Var(f *flags.IP4Flag, longName, shortName string, defValue string, usage string)
    IP4Var 绑定IPv4地址类型标志到指针并内部注册Flag对象 参数依次为: IPv4标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) IP6(longName, shortName string, defValue string, usage string) *flags.IP6Flag
    IP6 添加IPv6地址类型标志, 返回标志对象指针 参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: IPv6地址标志对象指针

func (c *Cmd) IP6Var(f *flags.IP6Flag, longName, shortName string, defValue string, usage string)
    IP6Var 绑定IPv6地址类型标志到指针并内部注册Flag对象 参数依次为: IPv6标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *flags.IntFlag
    Int 添加整数类型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: 整数标志对象指针

func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag
    Int64 添加64位整数类型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明

    返回值: 64位整数标志对象指针

func (c *Cmd) Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)
    Int64Var 绑定64位整数类型标志到指针并内部注册Flag对象

    参数依次为: 64位整数标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)
    IntVar 绑定整数类型标志到指针并内部注册Flag对象

    参数依次为: 整数标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) LoadHelp(filePath string) error
    LoadHelp 从指定文件加载帮助信息

    参数: filePath: 帮助信息文件路径

    返回值: error: 如果文件不存在或读取文件失败，则返回错误信息

func (c *Cmd) LongName() string
    LongName 返回命令长名称

func (c *Cmd) Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag
    Map 添加键值对类型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: 键值对标志对象指针

func (c *Cmd) MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)
    MapVar 绑定键值对类型标志到指针并内部注册Flag对象

    参数依次为: 键值对标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) NArg() int
    NArg 获取非标志参数的数量

func (c *Cmd) NFlag() int
    NFlag 获取标志的数量

func (c *Cmd) Parse(args []string) (err error)
    Parse 完整解析命令行参数（含子命令处理） 主要功能：
     1. 解析当前命令的长短标志及内置标志
     2. 自动检测并解析子命令及其参数（若存在）
     3. 验证枚举类型标志的有效性

    参数：

        args: 原始命令行参数切片（包含可能的子命令及参数）

    返回值：

        解析过程中遇到的错误（如标志格式错误、子命令解析失败等）

    注意事项：
      - 每个Cmd实例仅会被解析一次(线程安全)
      - 若检测到子命令, 会将剩余参数传递给子命令的Parse方法
      - 处理内置标志执行逻辑

func (c *Cmd) ParseFlagsOnly(args []string) (err error)
    ParseFlagsOnly 仅解析当前命令的标志参数（忽略子命令） 主要功能：
     1. 解析当前命令的长短标志及内置标志
     2. 验证枚举类型标志的有效性
     3. 明确忽略所有子命令及后续参数

    参数：

        args: 原始命令行参数切片（子命令及后续参数会被忽略）

    返回值：

        解析过程中遇到的错误（如标志格式错误等）

    注意事项：
      - 每个Cmd实例仅会被解析一次（线程安全）
      - 不会处理任何子命令，所有参数均视为当前命令的标志或位置参数
      - 处理内置标志逻辑

func (c *Cmd) Path(longName, shortName string, defValue string, usage string) *flags.PathFlag
    Path 添加路径类型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: 路径标志对象指针

func (c *Cmd) PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)
    PathVar 绑定路径类型标志到指针并内部注册Flag对象

    参数依次为: 路径标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) PrintHelp()
    PrintHelp 打印命令的帮助信息, 优先打印用户的帮助信息, 否则自动生成帮助信息

func (c *Cmd) SetDescription(desc string)
    SetDescription 设置命令描述

func (c *Cmd) SetHelp(help string)
    SetHelp 设置用户自定义命令帮助信息

func (c *Cmd) SetLogoText(logoText string)
    SetLogoText 设置logo文本

func (c *Cmd) SetModuleHelps(moduleHelps string)
    SetModuleHelps 设置自定义模块帮助信息

func (c *Cmd) SetUsageSyntax(usageSyntax string)
    SetUsageSyntax 设置自定义命令用法

func (c *Cmd) SetUseChinese(useChinese bool)
    SetUseChinese 设置是否使用中文帮助信息

func (c *Cmd) SetVersion(version string)
    SetVersion 设置版本信息

func (c *Cmd) ShortName() string
    ShortName 返回命令短名称

func (c *Cmd) Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag
    Slice 绑定字符串切片类型标志并内部注册Flag对象

    参数依次为: 长标志名、短标志、默认值、帮助说明

    返回值: 字符串切片标志对象指针

func (c *Cmd) SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)
    SliceVar 绑定字符串切片类型标志到指针并内部注册Flag对象

    参数依次为: 字符串切片标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) String(longName, shortName, defValue, usage string) *flags.StringFlag
    String 添加字符串类型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明

    返回值: 字符串标志对象指针

func (c *Cmd) StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)
    StringVar 绑定字符串类型标志到指针并内部注册Flag对象

    参数依次为: 字符串标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) SubCmds() []*Cmd
    SubCmds 返回子命令列表

func (c *Cmd) Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag
    Time 添加时间类型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明

    返回值: 时间标志对象指针

func (c *Cmd) TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)
    TimeVar 绑定时间类型标志到指针并内部注册Flag对象

    参数依次为: 时间标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) URL(longName, shortName string, defValue string, usage string) *flags.URLFlag
    URL 添加URL类型标志, 返回标志对象指针 参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: URL标志对象指针

func (c *Cmd) URLVar(f *flags.URLFlag, longName, shortName string, defValue string, usage string)
    URLVar 绑定URL类型标志到指针并内部注册Flag对象 参数依次为: URL标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag
    Uint16 添加16位无符号整数类型标志, 返回标志对象指针

    参数依次为: 长标志名、短标志、默认值、帮助说明

    返回值: 16位无符号整数标志对象指针

func (c *Cmd) Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)
    Uint16Var 绑定16位无符号整数类型标志到指针并内部注册Flag对象

    参数依次为: 16位无符号整数标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag
    Uint32 添加32位无符号整数类型标志, 返回标志对象指针 参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: 32位无符号整数标志对象指针

func (c *Cmd) Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string)
    Uint32Var 绑定32位无符号整数类型标志到指针并内部注册Flag对象 参数依次为: 32位无符号整数标志指针、长标志名、短标志、默认值、帮助说明

func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag
    Uint64 添加64位无符号整数类型标志, 返回标志对象指针 参数依次为: 长标志名、短标志、默认值、帮助说明 返回值: 64位无符号整数标志对象指针

func (c *Cmd) Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string)
    Uint64Var 绑定64位无符号整数类型标志到指针并内部注册Flag对象 参数依次为: 64位无符号整数标志指针、长标志名、短标志、默认值、帮助说明

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
    CmdInterface 命令接口定义，封装命令行程序的核心功能 提供统一的命令管理、参数解析和帮助系统 实现类需保证线程安全，所有方法应支持并发调用

    示例用法: cmd := NewCommand("app", "a", flag.ContinueOnError)
    cmd.SetDescription("示例应用程序") cmd.String("config", "c", "配置文件路径",
    "/etc/app.conf")

type ExampleInfo struct {
	Description string // 示例描述
	Usage       string // 示例使用方式
}
    ExampleInfo 示例信息结构体 用于存储命令的使用示例，包括描述和示例内容

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
    HelpTemplate 帮助信息模板结构体

type UserInfo struct {
	// Has unexported fields.
}
    UserInfo 存储用户自定义信息的嵌套结构体

