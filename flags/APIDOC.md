package flags // import "gitee.com/MM-Q/qflag/flags"

Package flags 基础标志结构体定义 本文件定义了BaseFlag泛型结构体，为所有标志类型提供通用的字段和方法实现，
包括标志的初始化、值获取设置、验证器支持等核心功能。

Package flags 基本数据类型标志实现 本文件实现了整数、浮点数、布尔、字符串等基本数据类型的标志结构体， 提供了相应的解析、验证和类型转换功能。

Package flags 扩展数据类型标志实现 本文件实现了枚举、时间间隔、切片、时间、映射等扩展数据类型的标志结构体，
提供了复杂数据类型的解析、验证和格式化功能。

Package flags 标志类型定义和接口 本文件定义了所有标志类型的通用接口和基础标志结构体，包括标志类型枚举、
验证器接口、标志接口等核心定义，为整个标志系统提供基础类型支持。

Package flags 标志注册表和元数据管理 本文件实现了FlagRegistry标志注册表，提供标志的注册、查找、索引管理等功能，
支持按长名称、短名称进行标志查找和管理。

CONSTANTS

const (
	CompletionShellDescCN = "生成指定的 Shell 补全脚本, 可选类型: %v"
	CompletionShellDescEN = "Generate the specified Shell completion script, optional types: %v"
)
    定义中英文的补全标志的使用说明

const (
	ShellBash       = "bash"       // bash shell
	ShellPowershell = "powershell" // powershell shell
	ShellPwsh       = "pwsh"       // pwsh shell
	ShellNone       = "none"       // 无shell
)
    支持的Shell类型

const (
	// 逗号
	FlagSplitComma = ","

	// 分号
	FlagSplitSemicolon = ";"

	// 竖线
	FlagSplitPipe = "|"

	// 冒号
	FlagKVColon = ":"

	// 等号
	FlagKVEqual = "="
)
    定义标志的分隔符常量

const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"
    定义非法字符集常量, 防止非法的标志名称


VARIABLES

var (
	HelpFlagName                 = "help"                      // 帮助标志名称
	HelpFlagShortName            = "h"                         // 帮助标志短名称
	VersionFlagLongName          = "version"                   // 版本标志名称
	VersionFlagShortName         = "v"                         // 版本标志短名称
	CompletionShellFlagLongName  = "generate-shell-completion" // 生成shell补全标志长名称
	CompletionShellFlagShortName = "gsc"                       // 生成shell补全标志短名称
)
    内置标志名称

var (
	HelpFlagUsage    = "Show help"    // 帮助标志使用说明
	VersionFlagUsage = "Show version" // 版本标志使用说明
)
    内置标志使用说明

var FlagSplitSlice = []string{

	FlagSplitComma,

	FlagSplitSemicolon,

	FlagSplitPipe,

	FlagKVColon,
}
    Flag支持的标志分隔符切片

var ShellSlice = []string{ShellNone, ShellBash, ShellPowershell, ShellPwsh}
    支持的Shell类型切片


FUNCTIONS

func FlagTypeToString(flagType FlagType) string
    FlagTypeToString 将FlagType转换为带语义信息的字符串

    参数:
      - flagType: 需要转换的FlagType枚举值

    返回值:
      - 带语义信息的类型字符串，用于命令行帮助信息显示


TYPES

type BaseFlag[T any] struct {
	// Has unexported fields.
}
    BaseFlag 泛型基础标志结构体,封装所有标志的通用字段和方法

func (f *BaseFlag[T]) BindEnv(envName string) *BaseFlag[T]
    BindEnv 绑定环境变量到标志

    参数:
      - envName 环境变量名

    返回:
      - 标志对象本身,支持链式调用

func (f *BaseFlag[T]) Get() T
    Get 获取标志的实际值(泛型类型)

    返回值:
      - T: 标志值

func (f *BaseFlag[T]) GetDefault() T
    GetDefault 获取标志的初始默认值(泛型类型)

    返回值:
      - T: 初始默认值

func (f *BaseFlag[T]) GetDefaultAny() any
    GetDefaultAny 获取标志的初始默认值(any类型)

    返回值:
      - any: 初始默认值

func (f *BaseFlag[T]) GetEnvVar() string
    GetEnvVar 获取绑定的环境变量名

    返回值:
      - string: 环境变量名

func (f *BaseFlag[T]) Init(longName, shortName string, usage string, value *T) error
    Init 初始化标志的元数据和值指针, 无需显式调用, 仅在创建标志对象时自动调用

    参数:
      - longName: 长标志名称
      - shortName: 短标志字符
      - usage: 帮助说明
      - value: 标志值指针

    返回值:
      - error: 初始化错误信息

func (f *BaseFlag[T]) IsSet() bool
    IsSet 判断标志是否已被设置值

    返回值:
      - bool: true表示已设置值, false表示未设置

func (f *BaseFlag[T]) LongName() string
    LongName 获取标志的长名称

    返回值:
      - string: 长标志名称

func (f *BaseFlag[T]) Name() string
    Name 获取标志的名称

    返回值:
      - string: 标志名称, 优先返回长名称, 如果长名称为空则返回短名称

func (f *BaseFlag[T]) Reset()
    Reset 将标志重置为初始默认值

func (f *BaseFlag[T]) Set(value T) error
    Set 设置标志的值(泛型类型)

    参数:
      - value T: 标志值

    返回:
      - error: 错误信息

func (f *BaseFlag[T]) SetValidator(validator Validator)
    SetValidator 设置标志的验证器(泛型类型)

    参数:
      - validator Validator: 验证器接口

func (f *BaseFlag[T]) ShortName() string
    ShortName 获取标志的短名称

    返回值:
      - string: 短标志字符

func (f *BaseFlag[T]) String() string
    String 返回标志的字符串表示

func (f *BaseFlag[T]) Type() FlagType
    Type 返回标志类型, 默认实现返回0, 需要子类重写

    注意：
      - 具体标志类型需要重写此方法返回正确的FlagType

func (f *BaseFlag[T]) Usage() string
    Usage 获取标志的用法说明

    返回值:
      - string: 用法说明

type BoolFlag struct {
	BaseFlag[bool]
	// Has unexported fields.
}
    BoolFlag 布尔类型标志结构体 继承BaseFlag[bool]泛型结构体,实现Flag接口

func (f *BoolFlag) IsBoolFlag() bool
    IsBoolFlag 实现flag.boolFlag接口,返回true

func (f *BoolFlag) Set(value string) error
    Set 实现flag.Value接口,解析并设置布尔值

    支持以下布尔值格式（大小写不敏感）:
      - 真值: "true", "1", "t", "T", "TRUE", "True"
      - 假值: "false", "0", "f", "F", "FALSE", "False"

    参数:
      - value: 待设置的布尔值字符串

    返回值:
      - error: 解析或验证失败时返回错误信息

    示例:
      - flag.Set("true") // ✅ 成功，值为 true
      - flag.Set("1") // ✅ 成功，值为 true
      - flag.Set("FALSE") // ✅ 成功，值为 false
      - flag.Set("0") // ✅ 成功，值为 false
      - flag.Set("yes") // ❌ 失败，返回解析错误

func (f *BoolFlag) String() string
    String 实现flag.Value接口,返回布尔值字符串

    返回值:
      - string: 布尔值字符串

func (f *BoolFlag) Type() FlagType
    Type 返回标志类型

    返回值:
      - FlagType: 标志类型枚举值

type DurationFlag struct {
	BaseFlag[time.Duration]
	// Has unexported fields.
}
    DurationFlag 时间间隔类型标志结构体 继承BaseFlag[time.Duration]泛型结构体,实现Flag接口

func (f *DurationFlag) Set(value string) error
    Set 实现flag.Value接口, 解析并设置时间间隔值

    参数:
      - value: 待设置的值

    返回值:
      - error: 解析或验证失败时返回错误信息

func (f *DurationFlag) String() string
    String 实现flag.Value接口, 返回当前值的字符串表示

    返回值:
      - string: 当前值的字符串表示

func (f *DurationFlag) Type() FlagType
    Type 返回标志类型

    返回值:
      - FlagType: 标志类型枚举值

type EnumFlag struct {
	BaseFlag[string]

	// Has unexported fields.
}
    EnumFlag 枚举类型标志结构体 继承BaseFlag[string]泛型结构体,增加枚举特有的选项验证

func (f *EnumFlag) GetOptions() []string
    GetOptions 返回枚举的所有可选值

    返回值:
      - []string: 枚举的所有可选值

func (f *EnumFlag) Init(longName, shortName string, defValue string, usage string, options []string) error
    Init 初始化枚举类型标志, 无需显式调用, 仅在创建标志对象时自动调用

    参数:
      - longName: 长标志名称
      - shortName: 短标志字符
      - defValue: 默认值
      - usage: 帮助说明
      - options: 枚举可选值列表

    返回值:
      - error: 初始化错误信息

func (f *EnumFlag) IsCheck(value string) error
    IsCheck 检查枚举值是否有效

    参数:
      - value: 待检查的枚举值

    返回值:
      - error: 为nil, 说明值有效,否则返回错误信息

func (f *EnumFlag) Set(value string) error
    Set 实现flag.Value接口, 解析并设置枚举值

    参数:
      - value: 待设置的值

    返回值:
      - error: 解析或验证失败时返回错误信息

func (f *EnumFlag) SetCaseSensitive(sensitive bool) *EnumFlag
    SetCaseSensitive 设置枚举值是否区分大小写

    参数:
      - sensitive - true表示区分大小写，false表示不区分（默认）

    返回值:
      - *EnumFlag - 返回自身以支持链式调用

func (f *EnumFlag) String() string
    String 实现flag.Value接口, 返回当前值的字符串表示

    返回值:
      - string: 当前值的字符串表示

func (f *EnumFlag) Type() FlagType
    Type 返回标志类型

    返回值:
      - FlagType: 标志类型枚举值

type Flag interface {
	LongName() string   // 获取标志的长名称
	ShortName() string  // 获取标志的短名称
	Usage() string      // 获取标志的用法
	Type() FlagType     // 获取标志类型
	GetDefaultAny() any // 获取标志的默认值(any类型)
	String() string     // 获取标志的字符串表示
	IsSet() bool        // 判断标志是否已设置值
	Reset()             // 重置标志值为默认值
	GetEnvVar() string  // 获取标志绑定的环境变量名称
}
    Flag 所有标志类型的通用接口,定义了标志的元数据访问方法

type FlagMeta struct {
	Flag Flag // 标志对象
}
    FlagMeta 统一存储标志的完整元数据

func (m *FlagMeta) GetDefault() any
    GetDefault 获取标志的默认值

func (m *FlagMeta) GetFlag() Flag
    GetFlag 获取标志对象

func (m *FlagMeta) GetFlagType() FlagType
    GetFlagType 获取标志的类型

func (m *FlagMeta) GetLongName() string
    GetLongName 获取标志的长名称

func (m *FlagMeta) GetName() string
    GetName 获取标志的名称

    优先返回长名称, 如果长名称为空, 则返回短名称

func (m *FlagMeta) GetShortName() string
    GetShortName 获取标志的短名称

func (m *FlagMeta) GetUsage() string
    GetUsage 获取标志的用法描述

type FlagRegistry struct {
	// Has unexported fields.
}
    FlagRegistry 集中管理所有标志元数据及索引

func NewFlagRegistry() *FlagRegistry
    创建一个空的标志注册表

    返回值:
      - *FlagRegistry: 创建的标志注册表指针

func (r *FlagRegistry) GetAllFlagsCount() int
    GetAllFlagsCount 获取所有标志数量(长标志+短标志)

    返回值:
      - int: 所有标志的数量

func (r *FlagRegistry) GetByLong(longName string) (*FlagMeta, bool)
    GetByLong 通过长标志名称查找对应的标志元数据

    参数:
      - longName: 标志的长名称(如"help")

    返回值:
      - *FlagMeta: 找到的标志元数据指针, 未找到时为nil
      - bool: 是否找到标志, true表示找到

func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool)
    GetByName 通过标志名称查找标志元数据

    参数:
      - name可以是长名称(如"help")或短名称(如"h")

    返回值:
      - *FlagMeta: 找到的标志元数据指针, 未找到时为nil
      - bool: 是否找到标志, true表示找到

func (r *FlagRegistry) GetByShort(shortName string) (*FlagMeta, bool)
    GetByShort 通过短标志名称查找对应的标志元数据

    参数:
      - shortName: 标志的短名称(如"h"对应"help")

    返回值:
      - *FlagMeta: 找到的标志元数据指针, 未找到时为nil
      - bool: 是否找到标志, true表示找到

func (r *FlagRegistry) GetFlagMetaCount() int
    GetFlagMetaCount 获取标志元数据数量

    返回值:
      - int: 标志元数据的数量

func (r *FlagRegistry) GetFlagMetaList() []*FlagMeta
    GetFlagMetaList 获取所有标志元数据列表

    返回值:
      - []*FlagMeta: 所有标志元数据的切片

func (r *FlagRegistry) GetFlagNameMap() map[string]*FlagMeta
    GetFlagNameMap 获取所有标志映射(长标志+短标志)

    返回值:
      - map[string]*FlagMeta: 长短标志名称到标志元数据的映射

func (r *FlagRegistry) GetLongFlagMap() map[string]*FlagMeta
    GetLongFlagMap 获取长标志映射

    返回值:
      - map[string]*FlagMeta: 长标志名称到标志元数据的映射

func (r *FlagRegistry) GetLongFlagsCount() int
    GetLongFlagsCount 获取长标志数量

    返回值:
      - int: 长标志的数量

func (r *FlagRegistry) GetShortFlagMap() map[string]*FlagMeta
    GetShortFlagMap 获取短标志映射

    返回值:
      - map[string]*FlagMeta: 短标志名称到标志元数据的映射

func (r *FlagRegistry) GetShortFlagsCount() int
    GetShortFlagsCount 获取短标志数量

    返回值:
      - int: 短标志的数量

func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error
    RegisterFlag 注册一个新的标志元数据到注册表中

    参数:
      - meta: 要注册的标志元数据

    该方法会执行以下操作:
      - 1.检查长名称和短名称是否已存在
      - 2.将标志添加到长名称索引
      - 3.将标志添加到短名称索引
      - 4.将标志添加到所有标志列表

    返回值:
      - error: 错误信息, 无错误时为nil

type FlagType int
    标志类型

const (
	FlagTypeUnknown  FlagType = iota // 未知类型
	FlagTypeInt                      // 整数类型
	FlagTypeInt64                    // 64位整数类型
	FlagTypeUint16                   // 16位无符号整数类型
	FlagTypeUint32                   // 32位无符号整数类型
	FlagTypeUint64                   // 64位无符号整数类型
	FlagTypeString                   // 字符串类型
	FlagTypeBool                     // 布尔类型
	FlagTypeFloat64                  // 64位浮点数类型
	FlagTypeEnum                     // 枚举类型
	FlagTypeDuration                 // 时间间隔类型
	FlagTypeSlice                    // 切片类型
	FlagTypeTime                     // 时间类型
	FlagTypeMap                      // 映射类型
)
type Float64Flag struct {
	BaseFlag[float64]
	// Has unexported fields.
}
    Float64Flag 浮点型标志结构体 继承BaseFlag[float64]泛型结构体,实现Flag接口

func (f *Float64Flag) Set(value string) error
    Set 实现flag.Value接口,解析并设置浮点值

    参数:
      - value: 待解析的浮点值

    返回值:
      - error: 解析错误或验证错误

func (f *Float64Flag) Type() FlagType
    Type 返回标志类型

    返回值:
      - FlagType: 标志类型枚举值

type Int64Flag struct {
	BaseFlag[int64]
	// Has unexported fields.
}
    Int64Flag 64位整数类型标志结构体 继承BaseFlag[int64]泛型结构体,实现Flag接口

func (f *Int64Flag) Set(value string) error
    Set 实现flag.Value接口,解析并设置64位整数值

    参数:
      - value: 待解析的64位整数值

    返回值:
      - error: 解析错误或验证错误

func (f *Int64Flag) SetRange(min, max int64)
    SetRange 设置64位整数的有效范围

    参数:
      - min: 最小值
      - max: 最大值

func (f *Int64Flag) Type() FlagType
    Type 返回标志类型

    返回值:
      - FlagType: 标志类型枚举值

type IntFlag struct {
	BaseFlag[int]
	// Has unexported fields.
}
    IntFlag 整数类型标志结构体 继承BaseFlag[int]泛型结构体,实现Flag接口

func (f *IntFlag) Set(value string) error
    Set 实现flag.Value接口,解析并验证整数值

    参数:
      - value: 待解析的整数值

    返回值:
      - error: 解析错误或验证错误

func (f *IntFlag) SetRange(min, max int)
    SetRange 设置整数的有效范围

    参数:
      - min: 最小值
      - max: 最大值

func (f *IntFlag) String() string
    String 实现flag.Value接口,返回当前整数值的字符串表示

    返回值:
      - string: 当前整数值的字符串表示

func (f *IntFlag) Type() FlagType
    Type 返回标志类型

    返回值:
      - FlagType: 标志类型枚举值

type MapFlag struct {
	BaseFlag[map[string]string]

	// Has unexported fields.
}
    MapFlag 键值对类型标志结构体 继承BaseFlag[map[string]string]泛型结构体,实现Flag接口

func (f *MapFlag) Set(value string) error
    Set 实现flag.Value接口,解析并设置键值对

    参数:
      - value: 待设置的值

    返回值:
      - error: 解析或验证失败时返回错误信息

func (f *MapFlag) SetDelimiters(keyDelimiter, valueDelimiter string)
    SetDelimiters 设置键值对分隔符

    参数：
      - keyDelimiter 键值对分隔符
      - valueDelimiter 键值分隔符

func (f *MapFlag) SetIgnoreCase(enable bool)
    SetIgnoreCase 设置是否忽略键的大小写

    参数:
      - enable: 是否忽略键的大小写

    注意:
      - 当enable为true时,所有键将转换为小写进行存储和比较

func (f *MapFlag) String() string
    String 实现flag.Value接口,返回当前值的字符串表示

    返回值:
      - string: 当前值的字符串表示

func (f *MapFlag) Type() FlagType
    Type 返回标志类型

    返回值:
      - FlagType: 标志类型枚举值

type SliceFlag struct {
	BaseFlag[[]string] // 基类

	// Has unexported fields.
}
    SliceFlag 切片类型标志结构体 继承BaseFlag[[]string]泛型结构体,实现Flag接口

func (f *SliceFlag) Clear() error
    Clear 清空切片所有元素

    返回值:
      - 操作成功返回nil, 否则返回错误信息

    注意：
      - 该方法会改变切片的指针

func (f *SliceFlag) Contains(element string) bool
    Contains 检查切片是否包含指定元素

    参数:
      - element 待检查的元素

    返回:
      - 若切片包含指定元素, 返回true, 否则返回false

    注意:
      - 当切片未设置值时,将使用默认值进行检查

func (f *SliceFlag) GetDelimiters() []string
    GetDelimiters 获取当前分隔符列表

func (f *SliceFlag) Init(longName, shortName string, defValue []string, usage string) error
    Init 初始化切片类型标志

    参数:
      - longName: 长标志名称
      - shortName: 短标志字符
      - defValue: 默认值（切片类型）
      - usage: 帮助说明

    返回值:
      - error: 初始化错误信息

func (f *SliceFlag) Len() int
    Len 获取切片长度

    返回:
      - 获取切片长度

func (f *SliceFlag) Remove(element string) error
    Remove 从切片中移除指定元素（支持移除空字符串元素）

    参数:
      - element 待移除的元素（支持空字符串）

    返回值:
      - 操作成功返回nil, 否则返回错误信息

func (f *SliceFlag) Set(value string) error
    Set 实现flag.Value接口, 解析并设置切片值

    参数:
      - value 待解析的切片值

    注意:
      - 如果切片中包含分隔符,则根据分隔符进行分割, 否则将整个值作为单个元素
      - 例如: "a,b,c" -> ["a", "b", "c"]

func (f *SliceFlag) SetDelimiters(delimiters []string)
    SetDelimiters 设置切片解析的分隔符列表

    参数:
      - delimiters 分隔符列表

func (f *SliceFlag) SetSkipEmpty(skip bool)
    SetSkipEmpty 设置是否跳过空元素

    参数:
      - skip - 为true时跳过空元素, 为false时保留空元素

    线程安全的空元素跳过更新

func (f *SliceFlag) Sort() error
    Sort 对切片进行排序 对当前切片标志的值进行原地排序，修改原切片内容
    采用Go标准库的sort.Strings()函数进行字典序排序(按Unicode代码点升序排列)

    注意：
      - 排序会直接修改当前标志的值，而非返回新切片
      - 排序区分大小写, 遵循Unicode代码点比较规则(如'A' < 'a' < 'z')
      - 若切片未设置值，将使用默认值进行排序

    返回值：
      - 排序成功返回nil, 若排序过程中发生错误则返回错误信息

func (f *SliceFlag) String() string
    String 实现flag.Value接口, 返回当前值的字符串表示

func (f *SliceFlag) Type() FlagType
    Type 返回标志类型

type StringFlag struct {
	BaseFlag[string]
}
    StringFlag 字符串类型标志结构体 继承BaseFlag[string]泛型结构体,实现Flag接口

func (f *StringFlag) Contains(substr string) bool
    Contains 检查字符串是否包含指定子串

    参数:
      - substr 子串

    返回值:
      - bool: 如果包含子串则返回true,否则返回false

func (f *StringFlag) Len() int
    Len 获取字符串标志的长度

    返回值：
      - 字符串的字符数(按UTF-8编码计算)

func (f *StringFlag) Set(value string) error
    Set 实现flag.Value接口的Set方法 将字符串值解析并设置到标志中

    参数:
      - value: 待设置的字符串值

    返回值:
      - error: 设置失败时返回错误信息

func (f *StringFlag) String() string
    String 返回带引号的字符串值

    返回值:
      - string: 带引号的字符串值

func (f *StringFlag) ToLower() string
    ToLower 将字符串标志值转换为小写

func (f *StringFlag) ToUpper() string
    ToUpper 将字符串标志值转换为大写

func (f *StringFlag) Type() FlagType
    Type 返回标志类型

    返回值:
      - FlagType: 标志类型枚举值

type TimeFlag struct {
	BaseFlag[time.Time]

	// Has unexported fields.
}
    TimeFlag 时间类型标志结构体 继承BaseFlag[time.Time]泛型结构体,实现Flag接口

func (f *TimeFlag) Init(longName, shortName string, defValue string, usage string) error
    Init 初始化时间类型标志，支持字符串类型默认值

    参数:
      - longName: 长标志名称
      - shortName: 短标志字符
      - defValue: 默认值（字符串格式，支持多种时间表达）
      - usage: 帮助说明

    返回值:
      - error: 初始化错误信息

    支持的默认值格式:
      - "now" 或 "" : 当前时间
      - "zero" : 零时间 (time.Time{})
      - "1h", "30m", "-2h" : 相对时间（基于当前时间的偏移）
      - "2006-01-02", "2006-01-02 15:04:05" : 绝对时间格式
      - RFC3339等标准格式

func (f *TimeFlag) Set(value string) error
    Set 实现flag.Value接口, 解析并设置时间值

    参数:
      - value: 待解析的时间字符串

    返回值:
      - error: 解析或验证失败时返回错误信息

func (f *TimeFlag) SetOutputFormat(format string)
    SetOutputFormat 设置时间输出格式

    参数:
      - format: 时间格式化字符串

    注意: 此方法线程安全

func (f *TimeFlag) String() string
    String 实现flag.Value接口, 返回当前时间的字符串表示

    返回值:
      - string: 格式化后的时间字符串

    注意: 加锁保证outputFormat和value的并发安全访问

func (f *TimeFlag) Type() FlagType
    Type 返回标志类型

    返回值:
      - FlagType: 标志类型枚举值

type TypedFlag[T any] interface {
	Flag                                 // 继承标志接口
	GetDefault() T                       // 获取标志的具体类型默认值
	Get() T                              // 获取标志的具体类型值
	GetPointer() *T                      // 获取标志值的指针
	Set(T) error                         // 设置标志的具体类型值
	SetValidator(Validator)              // 设置标志的验证器
	BindEnv(envName string) *BaseFlag[T] // 绑定环境变量
}
    TypedFlag 所有标志类型的通用接口,定义了标志的元数据访问方法和默认值访问方法

type Uint16Flag struct {
	BaseFlag[uint16] // 基类
	// Has unexported fields.
}
    Uint16Flag 16位无符号整数类型标志结构体 继承BaseFlag[uint16]泛型结构体,实现Flag接口

func (f *Uint16Flag) Set(value string) error
    Set 实现flag.Value接口, 解析并设置16位无符号整数值 验证值是否在uint16范围内(0-65535)

    参数:
      - value: 待设置的值(0-65535)

    返回值:
      - error: 解析或验证失败时返回错误信息

func (f *Uint16Flag) String() string
    String 实现flag.Value接口, 返回当前值的字符串表示

    返回值:
      - string: 当前值的字符串表示

func (f *Uint16Flag) Type() FlagType
    Type 返回标志类型

    返回值:
      - FlagType: 标志类型枚举值

type Uint32Flag struct {
	BaseFlag[uint32] // 基类
	// Has unexported fields.
}
    Uint32Flag 32位无符号整数类型标志结构体 继承BaseFlag[uint32]泛型结构体,实现Flag接口

func (f *Uint32Flag) Set(value string) error
    Set 实现flag.Value接口, 解析并设置32位无符号整数值 验证值是否在uint32范围内(0-4294967295)

    参数:
      - value: 待设置的值(0-4294967295)

    返回值:
      - error: 解析或验证失败时返回错误信息

func (f *Uint32Flag) String() string
    String 实现flag.Value接口, 返回当前值的字符串表示

func (f *Uint32Flag) Type() FlagType
    Type 返回标志类型

type Uint64Flag struct {
	BaseFlag[uint64] // 基类
	// Has unexported fields.
}
    Uint64Flag 64位无符号整数类型标志结构体 继承BaseFlag[uint64]泛型结构体,实现Flag接口

func (f *Uint64Flag) Set(value string) error
    Set 实现flag.Value接口, 解析并设置64位无符号整数值 验证值是否在uint64范围内(0-18446744073709551615)

    参数:
      - value: 待设置的值(0-18446744073709551615)

    返回值:
      - error: 解析或验证失败时返回错误信息

func (f *Uint64Flag) String() string
    String 实现flag.Value接口, 返回当前值的字符串表示

func (f *Uint64Flag) Type() FlagType
    Type 返回标志类型

type Validator interface {
	// Validate 验证参数值是否合法
	// value: 待验证的参数值
	// 返回值: 验证通过返回nil, 否则返回错误信息
	Validate(value any) error
}
    Validator 验证器接口, 所有自定义验证器需实现此接口

