# Package flags

flags 定义了所有标志类型的通用接口和基础标志结构体。

## CONSTANTS

```markdown
const (
  FlagSplitComma     = "," // 逗号
  FlagSplitSemicolon = ";" // 分号
  FlagSplitPipe      = "|" // 竖线
  FlagKVColon        = ":" // 冒号
  FlagKVEqual        = "=" // 等号
)
```

定义标志的分隔符常量。

```markdown
const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"
```

定义非法字符集常量，防止非法的标志名称。

## VARIABLES

```markdown
var (
  HelpFlagName         = "help"    // 帮助标志名称
  HelpFlagShortName    = "h"       // 帮助标志短名称
  VersionFlagLongName  = "version" // 版本标志名称
  VersionFlagShortName = "v"       // 版本标志短名称
)
```

内置标志名称。

```markdown
var (
  HelpFlagUsageEn    = "Show help"                       // 帮助标志英文使用说明
  VersionFlagUsageEn = "Show the version of the program" // 版本标志英文使用说明
  VersionFlagUsageZh = "显示程序的版本信息"                       // 版本标志中文使用说明
)
```

内置标志使用说明。

```markdown
var FlagSplitSlice = []string{
  FlagSplitComma,
  FlagSplitSemicolon,
  FlagSplitPipe,
  FlagKVColon,
}
```

Flag 支持的标志分隔符切片。

## FUNCTIONS

```markdown
func FlagTypeToString(flagType FlagType) string
```

FlagTypeToString 将 FlagType 转换为字符串。

## TYPES

### BaseFlag

```markdown
type BaseFlag[T any] struct {
  // Has unexported fields.
}
```

BaseFlag 泛型基础标志结构体，封装所有标志的通用字段和方法。

```markdown
func (f *BaseFlag[T]) BindEnv(envName string) *BaseFlag[T]
```

BindEnv 绑定环境变量到标志。

- 参数: `envName` 环境变量名。
- 返回: 标志对象本身，支持链式调用。

```markdown
func (f *BaseFlag[T]) Get() T
```

Get 获取标志的实际值。

```markdown
func (f *BaseFlag[T]) GetDefault() T
```

GetDefault 获取标志的初始默认值。

```markdown
func (f *BaseFlag[T]) GetDefaultAny() any
```

GetDefaultAny 获取标志的初始默认值 (any 类型)。

```markdown
func (f *BaseFlag[T]) GetEnvVar() string
```

GetEnvVar 获取绑定的环境变量名。

```markdown
func (f *BaseFlag[T]) GetPointer() *T
```

GetPointer 返回标志值的指针。

注意:

1. 获取指针过程受锁保护，但直接修改指针指向的值仍会绕过验证机制。
2. 多线程环境下修改时需额外同步措施，建议优先使用 Set() 方法。

```markdown
func (f *BaseFlag[T]) Init(longName, shortName string, usage string, value *T) error
```

Init 初始化标志的元数据和值指针，无需显式调用，仅在创建标志对象时自动调用。

- 参数:
  - `longName`: 长标志名称。
  - `shortName`: 短标志字符。
  - `usage`: 帮助说明。
  - `value`: 标志值指针。
- 返回值: 初始化错误信息。

```markdown
func (f *BaseFlag[T]) IsSet() bool
```

IsSet 判断标志是否已被设置值。

返回值: `true` 表示已设置值，`false` 表示未设置。

```markdown
func (f *BaseFlag[T]) LongName() string
```

LongName 获取标志的长名称。

```markdown
func (f *BaseFlag[T]) Reset()
```

Reset 将标志重置为初始默认值。

```markdown
func (f *BaseFlag[T]) Set(value T) error
```

Set 设置标志的值。

- 参数: `value` 标志值。
- 返回: 错误信息。

```markdown
func (f *BaseFlag[T]) SetValidator(validator Validator)
```

SetValidator 设置标志的验证器。

- 参数: `validator` 验证器接口。

```markdown
func (f *BaseFlag[T]) ShortName() string
```

ShortName 获取标志的短名称。

```markdown
func (f *BaseFlag[T]) String() string
```

String 返回标志的字符串表示。

```markdown
func (f *BaseFlag[T]) Type() FlagType
```

Type 返回标志类型。注意：具体标志类型需要重写此方法返回正确的 FlagType。

```markdown
func (f *BaseFlag[T]) Usage() string
```

Usage 获取标志的用法说明。

### BoolFlag

```markdown
type BoolFlag struct {
  BaseFlag[bool]
  // Has unexported fields.
}
```

BoolFlag 布尔类型标志结构体，继承 BaseFlag[bool] 泛型结构体，实现 Flag 接口。

```markdown
func (f *BoolFlag) IsBoolFlag() bool
```

IsBoolFlag 实现 flag.Value 接口，返回布尔值。

```markdown
func (f *BoolFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置布尔值。

```markdown
func (f *BoolFlag) String() string
```

String 实现 flag.Value 接口，返回布尔值字符串。

```markdown
func (f *BoolFlag) Type() FlagType
```

Type 返回标志类型。

### DurationFlag

```markdown
type DurationFlag struct {
  BaseFlag[time.Duration]
  // Has unexported fields.
}
```

DurationFlag 时间间隔类型标志结构体，继承 BaseFlag[time.Duration] 泛型结构体，实现 Flag 接口。

```markdown
func (f *DurationFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置时间间隔值。

```markdown
func (f *DurationFlag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *DurationFlag) Type() FlagType
```

Type 返回标志类型。

### EnumFlag

```markdown
type EnumFlag struct {
  BaseFlag[string]

  // Has unexported fields.
}
```

EnumFlag 枚举类型标志结构体，继承 BaseFlag[string] 泛型结构体，增加枚举特有的选项验证。

```markdown
func (f *EnumFlag) Init(longName, shortName string, defValue string, usage string, options []string) error
```

Init 初始化枚举类型标志，无需显式调用，仅在创建标志对象时自动调用。

- 参数:
  - `longName`: 长标志名称。
  - `shortName`: 短标志字符。
  - `defValue`: 默认值。
  - `usage`: 帮助说明。
  - `options`: 枚举可选值列表。
- 返回值: 初始化错误信息。

```markdown
func (f *EnumFlag) IsCheck(value string) error
```

IsCheck 检查枚举值是否有效。

返回值: 为 `nil`，说明值有效；否则返回错误信息。

```markdown
func (f *EnumFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置枚举值。

```markdown
func (f *EnumFlag) SetCaseSensitive(sensitive bool) *EnumFlag
```

SetCaseSensitive 设置枚举值是否区分大小写。

- 参数:
  - `sensitive`: `true` 表示区分大小写，`false` 表示不区分（默认）。
- 返回值: `*EnumFlag` 返回自身以支持链式调用。

```markdown
func (f *EnumFlag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *EnumFlag) Type() FlagType
```

Type 实现 Flag 接口。

### Flag

```markdown
type Flag interface {
  LongName() string   // 获取标志的长名称
  ShortName() string  // 获取标志的短名称
  Usage() string      // 获取标志的用法
  Type() FlagType     // 获取标志类型
  GetDefaultAny() any // 获取标志的默认值 (any 类型)
  String() string     // 获取标志的字符串表示
  IsSet() bool        // 判断标志是否已被设置值
  Reset()             // 重置标志值为默认值
  GetEnvVar() string  // 获取标志绑定的环境变量名称
}
```

Flag 所有标志类型的通用接口，定义了标志的元数据访问方法。

### FlagMeta

```markdown
type FlagMeta struct {
  Flag Flag // 标志对象
}
```

FlagMeta 统一存储标志的完整元数据。

```markdown
func (m *FlagMeta) GetDefault() any
```

GetDefault 获取标志的默认值。

```markdown
func (m *FlagMeta) GetFlag() Flag
```

GetFlag 获取标志对象。

```markdown
func (m *FlagMeta) GetFlagType() FlagType
```

GetFlagType 获取标志的类型。

```markdown
func (m *FlagMeta) GetLongName() string
```

GetLongName 获取标志的长名称。

```markdown
func (m *FlagMeta) GetName() string
```

GetName 获取标志的名称。

优先返回长名称，如果长名称为空，则返回短名称。

```markdown
func (m *FlagMeta) GetShortName() string
```

GetShortName 获取标志的短名称。

```markdown
func (m *FlagMeta) GetUsage() string
```

GetUsage 获取标志的用法描述。

### FlagMetaInterface

```markdown
type FlagMetaInterface interface {
  GetFlagType() FlagType // 获取标志类型
  GetFlag() Flag         // 获取标志对象
  GetLongName() string   // 获取标志的长名称
  GetShortName() string  // 获取标志的短名称
  GetName() string       // 获取标志的名称
  GetUsage() string      // 获取标志的用法描述
  GetDefault() any       // 获取标志的默认值
  GetValue() any         // 获取标志的当前值
}
```

FlagMetaInterface 标志元数据接口，定义了标志元数据的获取方法。

### FlagRegistry

```markdown
type FlagRegistry struct {
  // Has unexported fields.
}
```

FlagRegistry 集中管理所有标志元数据及索引。

```markdown
func NewFlagRegistry() *FlagRegistry
```

创建一个空的标志注册表。

```markdown
func (r *FlagRegistry) GetALLFlagsCount() int
```

GetALLFlagsCount 获取所有标志数量（长标志 + 短标志）。

返回值:
- `int`: 所有标志的数量。

```markdown
func (r *FlagRegistry) GetAllFlagMetas() []*FlagMeta
```

GetAllFlagMetas 获取所有标志元数据列表。

返回值:
- `[]*FlagMeta`: 所有标志元数据的切片。

```markdown
func (r *FlagRegistry) GetAllFlags() map[string]*FlagMeta
```

GetALLFlags 获取所有标志映射（长标志 + 短标志）。

返回值:
- `map[string]*FlagMeta`: 长短标志名称到标志元数据的映射。

```markdown
func (r *FlagRegistry) GetByLong(longName string) (*FlagMeta, bool)
```

GetByLong 通过长标志名称查找对应的标志元数据。

- 参数:
  - `longName`: 标志的长名称（如 "help"）。
- 返回值:
  - `*FlagMeta`: 找到的标志元数据指针，未找到时为 `nil`。
  - `bool`: 是否找到标志，`true` 表示找到。

```markdown
func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool)
```

GetByName 通过标志名称查找标志元数据。

- 参数:
  - `name`: 可以是长名称（如 "help"）或短名称（如 "h"）。
- 返回值:
  - `*FlagMeta`: 找到的标志元数据指针，未找到时为 `nil`。
  - `bool`: 是否找到标志，`true` 表示找到。

```markdown
func (r *FlagRegistry) GetByShort(shortName string) (*FlagMeta, bool)
```

GetByShort 通过短标志名称查找对应的标志元数据。

- 参数:
  - `shortName`: 标志的短名称（如 "h" 对应 "help"）。
- 返回值:
  - `*FlagMeta`: 找到的标志元数据指针，未找到时为 `nil`。
  - `bool`: 是否找到标志，`true` 表示找到。

```markdown
func (r *FlagRegistry) GetLongFlags() map[string]*FlagMeta
```

GetLongFlags 获取长标志映射。

返回值:
- `map[string]*FlagMeta`: 长标志名称到标志元数据的映射。

```markdown
func (r *FlagRegistry) GetLongFlagsCount() int
```

GetLongFlagsCount 获取长标志数量。

返回值:
- `int`: 长标志的数量。

```markdown
func (r *FlagRegistry) GetShortFlags() map[string]*FlagMeta
```

GetShortFlags 获取短标志映射。

返回值:
- `map[string]*FlagMeta`: 短标志名称到标志元数据的映射。

```markdown
func (r *FlagRegistry) GetShortFlagsCount() int
```

GetShortFlagsCount 获取短标志数量。

返回值:
- `int`: 短标志的数量。

```markdown
func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error
```

RegisterFlag 注册一个新的标志元数据到注册表中。

- 参数:
  - `meta`: 要注册的标志元数据。
- 该方法会执行以下操作:
  1. 检查长名称和短名称是否已存在。
  2. 将标志添加到长名称索引。
  3. 将标志添加到短名称索引。
  4. 将标志添加到所有标志列表。
- 注意: 该方法线程安全，但发现重复标志时会 panic。

### FlagRegistryInterface

```markdown
type FlagRegistryInterface interface {
  GetAllFlagMetas() []*FlagMeta                  // 获取所有标志元数据列表
  GetLongFlags() map[string]*FlagMeta            // 获取长标志映射
  GetShortFlags() map[string]*FlagMeta           // 获取短标志映射
  GetAllFlags() map[string]*FlagMeta             // 获取所有标志映射（长 + 短）
  GetLongFlagsCount() int                        // 获取长标志数量
  GetShortFlagsCount() int                       // 获取短标志数量
  GetALLFlagsCount() int                         // 获取所有标志数量（长 + 短）
  RegisterFlag(meta *FlagMeta) error             // 注册一个新的标志元数据到注册表中
  GetByLong(longName string) (*FlagMeta, bool)   // 通过长标志名称查找对应的标志元数据
  GetByShort(shortName string) (*FlagMeta, bool) // 通过短标志名称查找对应的标志元数据
  GetByName(name string) (*FlagMeta, bool)       // 通过标志名称查找标志元数据
}
```

FlagRegistryInterface 标志注册表接口，定义了标志元数据的增删改查操作。

### FlagType

```markdown
type FlagType int
```

标志类型。

```markdown
const (
  FlagTypeInt      FlagType = iota + 1 // 整数类型
  FlagTypeInt64                        // 64 位整数类型
  FlagTypeUint16                       // 16 位无符号整数类型
  FlagTypeUint32                       // 32 位无符号整数类型
  FlagTypeUint64                       // 64 位无符号整数类型
  FlagTypeString                       // 字符串类型
  FlagTypeBool                         // 布尔类型
  FlagTypeFloat64                      // 64 位浮点数类型
  FlagTypeEnum                         // 枚举类型
  FlagTypeDuration                     // 时间间隔类型
  FlagTypeSlice                        // 切片类型
  FlagTypeTime                         // 时间类型
  FlagTypeMap                          // 映射类型
  FlagTypePath                         // 路径类型
  FlagTypeIP4                          // IPv4 地址类型
  FlagTypeIP6                          // IPv6 地址类型
  FlagTypeURL                          // URL 类型
)
```

### Float64Flag

```markdown
type Float64Flag struct {
  BaseFlag[float64]
  // Has unexported fields.
}
```

Float64Flag 浮点型标志结构体，继承 BaseFlag[float64] 泛型结构体，实现 Flag 接口。

```markdown
func (f *Float64Flag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置浮点值。

```markdown
func (f *Float64Flag) Type() FlagType
```

Type 返回标志类型。

### IP4Flag

```markdown
type IP4Flag struct {
  BaseFlag[string]
  // Has unexported fields.
}
```

IP4Flag IPv4 地址类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```markdown
func (f *IP4Flag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并验证 IPv4 地址。

```markdown
func (f *IP4Flag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *IP4Flag) Type() FlagType
```

Type 返回标志类型。

### IP6Flag

```markdown
type IP6Flag struct {
  BaseFlag[string]
  // Has unexported fields.
}
```

IP6Flag IPv6 地址类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```markdown
func (f *IP6Flag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并验证 IPv6 地址。

```markdown
func (f *IP6Flag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *IP6Flag) Type() FlagType
```

Type 返回标志类型。

### Int64Flag

```markdown
type Int64Flag struct {
  BaseFlag[int64]
  // Has unexported fields.
}
```

Int64Flag 64 位整数类型标志结构体，继承 BaseFlag[int64] 泛型结构体，实现 Flag 接口。

```markdown
func (f *Int64Flag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置 64 位整数值。

```markdown
func (f *Int64Flag) SetRange(min, max int64)
```

SetRange 设置 64 位整数的有效范围。

- 参数:
  - `min`: 最小值。
  - `max`: 最大值。

```markdown
func (f *Int64Flag) Type() FlagType
```

Type 返回标志类型。

### IntFlag

```markdown
type IntFlag struct {
  BaseFlag[int]
  // Has unexported fields.
}
```

IntFlag 整数类型标志结构体，继承 BaseFlag[int] 泛型结构体，实现 Flag 接口。

```markdown
func (f *IntFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并验证整数值。

```markdown
func (f *IntFlag) SetRange(min, max int)
```

SetRange 设置整数的有效范围。

- 参数:
  - `min`: 最小值。
  - `max`: 最大值。

```markdown
func (f *IntFlag) String() string
```

String 实现 flag.Value 接口，返回当前整数值的字符串表示。

```markdown
func (f *IntFlag) Type() FlagType
```

Type 返回标志类型。

### MapFlag

```markdown
type MapFlag struct {
  BaseFlag[map[string]string]

  // Has unexported fields.
}
```

MapFlag 键值对类型标志结构体，继承 BaseFlag[map[string]string] 泛型结构体，实现 Flag 接口。

```markdown
func (f *MapFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置键值对。

```markdown
func (f *MapFlag) SetDelimiters(keyDelimiter, valueDelimiter string)
```

SetDelimiters 设置键值对分隔符。

- 参数:
  - `keyDelimiter`: 键值对分隔符。
  - `valueDelimiter`: 键值分隔符。

```markdown
func (f *MapFlag) SetIgnoreCase(enable bool)
```

SetIgnoreCase 设置是否忽略键的大小写。

- 参数:
  - `enable`: 为 `true` 时，所有键将转换为小写进行存储和比较。

```markdown
func (f *MapFlag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *MapFlag) Type() FlagType
```

Type 返回标志类型。

### PathFlag

```markdown
type PathFlag struct {
  BaseFlag[string]
  // Has unexported fields.
}
```

PathFlag 路径类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```markdown
func (f *PathFlag) Init(longName, shortName string, defValue string, usage string) error
```

Init 初始化路径标志。

```markdown
func (f *PathFlag) IsDirectory(isDir bool) *PathFlag
```

IsDirectory 设置路径是否必须是目录。

示例: `cmd.Path("log-dir", "l", "/var/log/app", "日志目录").IsDirectory(true)`

```markdown
func (f *PathFlag) MustExist(mustExist bool) *PathFlag
```

MustExist 设置路径是否必须存在。

示例: `cmd.Path("output", "o", "/tmp/output", "输出目录").MustExist(false)`

```markdown
func (f *PathFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并验证路径。

```markdown
func (f *PathFlag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *PathFlag) Type() FlagType
```

Type 返回标志类型。

### SliceFlag

```markdown
type SliceFlag struct {
  BaseFlag[[]string] // 基类

  // Has unexported fields.
}
```

SliceFlag 切片类型标志结构体，继承 BaseFlag[[]string] 泛型结构体，实现 Flag 接口。

```markdown
func (f *SliceFlag) Clear() error
```

Clear 清空切片所有元素。

返回值: 操作成功返回 `nil`，否则返回错误信息。

注意：该方法会改变切片的指针。

```markdown
func (f *SliceFlag) Contains(element string) bool
```

Contains 检查切片是否包含指定元素。当切片未设置值时，将使用默认值进行检查。

```markdown
func (f *SliceFlag) GetDelimiters() []string
```

GetDelimiters 获取当前分隔符列表。

```markdown
func (f *SliceFlag) Init(longName, shortName string, defValue []string, usage string) error
```

Init 初始化切片类型标志。

- 参数:
  - `longName`: 长标志名称。
  - `shortName`: 短标志字符。
  - `defValue`: 默认值（切片类型）。
  - `usage`: 帮助说明。
- 返回值: 初始化错误信息。

```markdown
func (f *SliceFlag) Len() int
```

Len 获取切片长度。

返回: 获取切片长度。

```markdown
func (f *SliceFlag) Remove(element string) error
```

Remove 从切片中移除指定元素（支持移除空字符串元素）。

- 参数: `element` 待移除的元素（支持空字符串）。
- 返回值: 操作成功返回 `nil`，否则返回错误信息。

```markdown
func (f *SliceFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置切片值。

- 参数: `value` 待解析的切片值。
- 注意: 如果切片中包含分隔符，则根据分隔符进行分割，否则将整个值作为单个元素。例如: `"a,b,c" -> ["a", "b", "c"]`

```markdown
func (f *SliceFlag) SetDelimiters(delimiters []string)
```

SetDelimiters 设置切片解析的分隔符列表。

- 参数: `delimiters` 分隔符列表。
- 线程安全的分隔符更新。

```markdown
func (f *SliceFlag) SetSkipEmpty(skip bool)
```

SetSkipEmpty 设置是否跳过空元素。

- 参数: `skip` 为 `true` 时跳过空元素，为 `false` 时保留空元素。

```markdown
func (f *SliceFlag) Sort() error
```

Sort 对切片进行排序。

功能：对当前切片标志的值进行原地排序，修改原切片内容。

排序规则:
- 采用 Go 标准库的 `sort.Strings()` 函数进行字典序排序（按 Unicode 代码点升序排列）。

注意事项：
1. 排序会直接修改当前标志的值，而非返回新切片。
2. 排序区分大小写，遵循 Unicode 代码点比较规则（如 `'A' < 'a' < 'z'`）。
3. 若切片未设置值，将使用默认值进行排序。

返回值：
- 排序成功返回 `nil`，若排序过程中发生错误则返回错误信息。

```markdown
func (f *SliceFlag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *SliceFlag) Type() FlagType
```

Type 返回标志类型。

### StringFlag

```markdown
type StringFlag struct {
  BaseFlag[string]
}
```

StringFlag 字符串类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```markdown
func (f *StringFlag) Contains(substr string) bool
```

Contains 检查字符串是否包含指定子串。

- 参数: `substr` 子串。
- 返回值: `true` 表示包含，`false` 表示不包含。

```markdown
func (f *StringFlag) Len() int
```

Len 获取字符串标志的长度。

返回值：字符串的字符数（按 UTF-8 编码计算）。

```markdown
func (f *StringFlag) Set(value string) error
```

Set 实现 flag.Value 接口的 Set 方法，将字符串值解析并设置到标志中。

```markdown
func (f *StringFlag) String() string
```

String 返回带引号的字符串值。

```markdown
func (f *StringFlag) ToLower() string
```

ToLower 将字符串标志值转换为小写。

```markdown
func (f *StringFlag) ToUpper() string
```

ToUpper 将字符串标志值转换为大写。

```markdown
func (f *StringFlag) Type() FlagType
```

Type 返回标志类型。

### TimeFlag

```markdown
type TimeFlag struct {
  BaseFlag[time.Time]

  // Has unexported fields.
}
```

TimeFlag 时间类型标志结构体，继承 BaseFlag[time.Time] 泛型结构体，实现 Flag 接口。

```markdown
func (f *TimeFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置时间值。

```markdown
func (f *TimeFlag) SetOutputFormat(format string)
```

SetOutputFormat 设置时间输出格式。

```markdown
func (f *TimeFlag) String() string
```

String 实现 flag.Value 接口，返回当前时间的字符串表示。加锁保证 outputFormat 和 value 的并发安全访问。

```markdown
func (f *TimeFlag) Type() FlagType
```

Type 返回标志类型。

### TypedFlag

```markdown
type TypedFlag[T any] interface {
  Flag                                 // 继承标志接口
  GetDefault() T                       // 获取标志的具体类型默认值
  Get() T                              // 获取标志的具体类型值
  GetPointer() *T                      // 获取标志值的指针
  Set(T) error                         // 设置标志的具体类型值
  SetValidator(Validator)              // 设置标志的验证器
  BindEnv(envName string) *BaseFlag[T] // 绑定环境变量
}
```

TypedFlag 所有标志类型的通用接口，定义了标志的元数据访问方法和默认值访问方法。

### URLFlag

```markdown
type URLFlag struct {
  BaseFlag[string]
  // Has unexported fields.
}
```

URLFlag URL 类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```markdown
func (f *URLFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并验证 URL 格式。

```markdown
func (f *URLFlag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *URLFlag) Type() FlagType
```

Type 返回标志类型。

### Uint16Flag

```markdown
type Uint16Flag struct {
  BaseFlag[uint16] // 基类
  // Has unexported fields.
}
```

Uint16Flag 16 位无符号整数类型标志结构体，继承 BaseFlag[uint16] 泛型结构体，实现 Flag 接口。

```markdown
func (f *Uint16Flag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置 16 位无符号整数值。验证值是否在 uint16 范围内（0-65535）。

- 参数:
  - `value`: 待设置的值（0-65535）。
- 返回值:
  - `error`: 错误信息。

```markdown
func (f *Uint16Flag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *Uint16Flag) Type() FlagType
```

Type 返回标志类型。

### Uint32Flag

```markdown
type Uint32Flag struct {
  BaseFlag[uint32] // 基类
  // Has unexported fields.
}
```

Uint32Flag 32 位无符号整数类型标志结构体，继承 BaseFlag[uint32] 泛型结构体，实现 Flag 接口。

```markdown
func (f *Uint32Flag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置 32 位无符号整数值。验证值是否在 uint32 范围内（0-4294967295）。

```markdown
func (f *Uint32Flag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *Uint32Flag) Type() FlagType
```

Type 返回标志类型。

### Uint64Flag

```markdown
type Uint64Flag struct {
  BaseFlag[uint64] // 基类
  // Has unexported fields.
}
```

Uint64Flag 64 位无符号整数类型标志结构体，继承 BaseFlag[uint64] 泛型结构体，实现 Flag 接口。

```markdown
func (f *Uint64Flag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置 64 位无符号整数值。验证值是否在 uint64 范围内（0-18446744073709551615）。

```markdown
func (f *Uint64Flag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```markdown
func (f *Uint64Flag) Type() FlagType
```

Type 返回标志类型。

### Validator

```markdown
type Validator interface {
  // Validate 验证参数值是否合法
  // value: 待验证的参数值
  // 返回值: 验证通过返回 nil, 否则返回错误信息
  Validate(value any) error
}
```

Validator 验证器接口，所有自定义验证器需实现此接口。