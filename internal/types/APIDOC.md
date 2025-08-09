package types // import "gitee.com/MM-Q/qflag/internal/types"

Package types 内置类型和数据结构定义 本文件定义了qflag包内部使用的内置类型和数据结构， 包括内置标志、配置选项等核心数据类型的定义。

Package types 配置结构体和选项定义 本文件定义了命令配置相关的结构体和选项，包括命令的各种配置参数、
帮助信息设置、版本信息等配置数据的定义和管理。

Package types 命令上下文和状态管理 本文件定义了命令上下文结构体，用于管理命令的状态、子命令、标志注册表等信息，
提供命令执行过程中的状态维护和数据共享功能。

TYPES

type BuiltinFlags struct {
	Help       *flags.BoolFlag // 标志-帮助
	Version    *flags.BoolFlag // 标志-版本
	Completion *flags.EnumFlag // 标志-自动完成
	NameMap    sync.Map        // 内置标志名称映射
}
    BuiltinFlags 内置标志结构体

func NewBuiltinFlags() *BuiltinFlags
    NewBuiltinFlags 创建内置标志实例

func (bf *BuiltinFlags) IsBuiltinFlag(name string) bool
    IsBuiltinFlag 检查是否为内置标志

    参数:
      - name: 标志名称

    返回值:
      - bool: 是否为内置标志

func (bf *BuiltinFlags) MarkAsBuiltin(names ...string)
    MarkAsBuiltin 标记为内置标志

    参数:
      - names: 标志名称列表

type CmdConfig struct {
	// 版本信息
	Version string

	// 自定义描述
	Description string

	// 自定义的完整命令行帮助信息
	Help string

	// 自定义用法格式说明
	UsageSyntax string

	// 模块帮助信息
	ModuleHelps string

	// logo文本
	LogoText string

	// 备注内容切片
	Notes []string

	// 示例信息切片
	Examples []ExampleInfo

	// 是否使用中文帮助信息
	UseChinese bool

	// 控制内置标志是否自动退出
	ExitOnBuiltinFlags bool

	// 控制是否启用自动补全功能
	EnableCompletion bool
}
    CmdConfig 命令行配置

func NewCmdConfig() *CmdConfig
    NewCmdConfig 创建一个新的CmdConfig实例

type CmdContext struct {
	// 长命令名称
	LongName string
	// 短命令名称
	ShortName string

	// 标志注册表, 统一管理标志的元数据
	FlagRegistry *flags.FlagRegistry
	// 底层flag集合, 处理参数解析
	FlagSet *flag.FlagSet

	// 命令行参数(非标志参数)
	Args []string
	// 是否已经解析过参数
	Parsed atomic.Bool
	// 用于确保参数解析只执行一次
	ParseOnce sync.Once
	// 读写锁
	Mutex sync.RWMutex

	// 子命令上下文切片
	SubCmds []*CmdContext
	// 子命令映射表
	SubCmdMap map[string]*CmdContext
	// 父命令上下文
	Parent *CmdContext

	// 配置信息
	Config *CmdConfig

	// 内置标志结构体
	BuiltinFlags *BuiltinFlags

	// ParseHook 解析阶段钩子函数
	// 在标志解析完成后、子命令参数处理后调用
	//
	// 参数:
	//   - 当前命令上下文
	//
	// 返回值:
	//   - error: 错误信息, 非nil时会中断解析流程
	//   - bool: 是否需要退出程序
	ParseHook func(*CmdContext) (error, bool)
}
    CmdContext 命令上下文，包含所有必要的状态信息 这是所有函数操作的核心数据结构

func NewCmdContext(longName, shortName string, errorHandling flag.ErrorHandling) *CmdContext
    NewCmdContext 创建新的命令上下文

    参数:
      - longName: 长命令名称
      - shortName: 短命令名称
      - errorHandling: 错误处理方式

    返回值:
      - *CmdContext: 新创建的命令上下文

    errorHandling可选参数:
      - flag.ContinueOnError: 解析标志时遇到错误继续解析, 并返回错误信息
      - flag.ExitOnError: 解析标志时遇到错误立即退出程序, 并返回错误信息
      - flag.PanicOnError: 解析标志时遇到错误立即触发panic

func (ctx *CmdContext) GetName() string
    GetName 获取命令名称 如果长命令名称不为空则返回长命令名称, 否则返回短命令名称

    返回值:
      - string: 命令名称

type ExampleInfo struct {
	Description string // 示例描述
	Usage       string // 示例使用方式
}
    ExampleInfo 示例信息结构体 用于存储命令的使用示例，包括描述和示例内容

    字段:
      - Description: 示例描述
      - Usage: 示例使用方式

