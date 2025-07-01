# Package flags



flags 定义了所有标志类型的通用接口和基础标志结构体。

## CONSTANTS

```go
const (
	FlagSplitComma     = "," // 逗号
	FlagSplitSemicolon = ";" // 分号
	FlagSplitPipe      = "|" // 竖线
	FlagKVColon        = ":" // 冒号
	FlagKVEqual        = "=" // 等号
)

const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"
```

定义标志的分隔符常量及非法字符集常量，防止非法的标志名称。

## VARIABLES

```go
var (
	HelpFlagName            = "help"    // 帮助标志名称
	HelpFlagShortName       = "h"       // 帮助标志短名称
	ShowInstallPathFlagName = "sip"     // 显示安装路径标志名称
	VersionFlagLongName     = "version" // 版本标志名称
	VersionFlagShortName    = "v"       // 版本标志短名称
)

var FlagSplitSlice = []string{
	FlagSplitComma,
	FlagSplitSemicolon,
	FlagSplitPipe,
	FlagKVColon,
}
```

内置标志名称及Flag支持的标志分隔符切片。

## FUNCTIONS

```go
func FlagTypeToString(flagType FlagType) string
```

FlagTypeToString 将FlagType转换为字符串。

## TYPES

### BaseFlag[T any]

```go
type BaseFlag[T any] struct {
	// Has unexported fields.
}
```

BaseFlag 泛型基础标志结构体，封装所有标志的通用字段和方法。

```go
func (f *BaseFlag[T]) Get() T
func (f *BaseFlag[T]) GetDefault() T
func (f *BaseFlag[T]) GetDefaultAny() any
func (f *BaseFlag[T]) GetPointer() *T
func (f *BaseFlag[T]) Init(longName, shortName string, defValue T, usage string, value *T) error
func (f *BaseFlag[T]) IsSet() bool
func (f *BaseFlag[T]) LongName() string
func (f *BaseFlag[T]) Reset()
func (f *BaseFlag[T]) Set(value T) error
func (f *BaseFlag[T]) SetValidator(validator Validator)
func (f *BaseFlag[T]) ShortName() string
func (f *BaseFlag[T]) String() string
func (f *BaseFlag[T]) Usage() string
```

方法说明：
- **Get**：获取标志的实际值，优先级：已设置的值 > 默认值。线程安全，使用互斥锁保证并发访问安全。
- **GetDefault**：获取标志的默认值。
- **GetDefaultAny**：获取标志的默认值(any类型)。
- **GetPointer**：返回标志值的指针。注意：获取指针过程受锁保护，但直接修改指针指向的值仍会绕过验证机制；多线程环境下修改时需额外同步措施，建议优先使用Set()方法。
- **Init**：初始化标志的元数据和值指针，无需显式调用，仅在创建标志对象时自动调用。
- **IsSet**：判断标志是否已被设置值。返回值：true表示已设置值，false表示未设置。
- **LongName**：获取标志的长名称。
- **Reset**：将标志值重置为默认值。线程安全，使用互斥锁保证并发安全。
- **Set**：设置标志的值。
- **SetValidator**：设置标志的验证器。
- **ShortName**：获取标志的短名称。
- **String**：返回标志的字符串表示。
- **Usage**：获取标志的用法说明。

### BoolFlag

```go
type BoolFlag struct {
	BaseFlag[bool]
}
```

BoolFlag 布尔类型标志结构体，继承BaseFlag[bool]泛型结构体，实现Flag接口。

```go
func (f *BoolFlag) Type() FlagType
```

**Type**：返回标志类型。

### DurationFlag

```go
type DurationFlag struct {
	BaseFlag[time.Duration]
}
```

DurationFlag 时间间隔类型标志结构体，继承BaseFlag[time.Duration]泛型结构体，实现Flag接口。

```go
func (f *DurationFlag) Set(value string) error
func (f *DurationFlag) String() string
func (f *DurationFlag) Type() FlagType
```

方法说明：
- **Set**：实现flag.Value接口，解析并设置时间间隔值。
- **String**：实现flag.Value接口，返回当前值的字符串表示。
- **Type**：返回标志类型。

### EnumFlag

```go
type EnumFlag struct {
	BaseFlag[string]
	// Has unexported fields.
}
```

EnumFlag 枚举类型标志结构体，继承BaseFlag[string]泛型结构体，增加枚举特有的选项验证。

```go
func (f *EnumFlag) Init(longName, shortName string, defValue string, usage string, options []string) error
func (f *EnumFlag) IsCheck(value string) error
func (f *EnumFlag) Set(value string) error
func (f *EnumFlag) String() string
func (f *EnumFlag) Type() FlagType
```

方法说明：
- **Init**：初始化枚举类型标志，无需显式调用，仅在创建标志对象时自动调用。
- **IsCheck**：检查枚举值是否有效。返回值：为nil，说明值有效，否则返回错误信息。
- **Set**：实现flag.Value接口，解析并设置枚举值。
- **String**：实现flag.Value接口，返回当前值的字符串表示。
- **Type**：实现Flag接口。

### Flag

```go
type Flag interface {
	LongName() string   // 获取标志的长名称
	ShortName() string  // 获取标志的短名称
	Usage() string      // 获取标志的用法
	Type() FlagType     // 获取标志类型
	GetDefaultAny() any // 获取标志的默认值(any类型)
	String() string     // 获取标志的字符串表示
	IsSet() bool        // 判断标志是否已设置值
	Reset()             // 重置标志值为默认值
}
```

Flag 所有标志类型的通用接口，定义了标志的元数据访问方法。

### FlagMeta

```go
type FlagMeta struct {
	Flag Flag // 标志对象
}
```

FlagMeta 统一存储标志的完整元数据。

```go
func (m *FlagMeta) GetDefault() any
func (m *FlagMeta) GetFlag() Flag
func (m *FlagMeta) GetFlagType() FlagType
func (m *FlagMeta) GetLongName() string
func (m *FlagMeta) GetShortName() string
func (m *FlagMeta) GetUsage() string
```

### FlagMetaInterface

```go
type FlagMetaInterface interface {
	GetFlagType() FlagType // 获取标志类型
	GetFlag() Flag         // 获取标志对象
	GetLongName() string   // 获取标志的长名称
	GetShortName() string  // 获取标志的短名称
	GetUsage() string      // 获取标志的用法描述
	GetDefault() any       // 获取标志的默认值
	GetValue() any         // 获取标志的当前值
}
```

FlagMetaInterface 标志元数据接口，定义了标志元数据的获取方法。

### FlagRegistry

```go
type FlagRegistry struct {
	// Has unexported fields.
}
```

FlagRegistry 集中管理所有标志元数据及索引。

```go
func NewFlagRegistry() *FlagRegistry
func (r *FlagRegistry) GetAllFlags() []*FlagMeta
func (r *FlagRegistry) GetByLong(longName string) (*FlagMeta, bool)
func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool)
func (r *FlagRegistry) GetByShort(shortName string) (*FlagMeta, bool)
func (r *FlagRegistry) GetLongFlags() map[string]*FlagMeta
func (r *FlagRegistry) GetShortFlags() map[string]*FlagMeta
func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error
```

方法说明：
- **NewFlagRegistry**：创建一个空的标志注册表。
- **GetAllFlags**：获取所有标志元数据列表。
- **GetByLong**：通过长标志名称查找对应的标志元数据。
- **GetByName**：通过标志名称查找标志元数据，参数name可以是长名称或短名称。
- **GetByShort**：通过短标志名称查找对应的标志元数据。
- **GetLongFlags**：获取长标志映射。
- **GetShortFlags**：获取短标志映射。
- **RegisterFlag**：注册一个新的标志元数据到注册表中。该方法会检查长名称和短名称是否已存在，并将标志添加到长名称索引、短名称索引以及所有标志列表。该方法线程安全，但发现重复标志时会panic。

### FlagRegistryInterface

```go
type FlagRegistryInterface interface {
	GetAllFlags() []*FlagMeta                      // 获取所有标志元数据列表
	GetLongFlags() map[string]*FlagMeta            // 获取长标志映射
	GetShortFlags() map[string]*FlagMeta           // 获取短标志映射
	RegisterFlag(meta *FlagMeta) error             // 注册一个新的标志元数据到注册表中
	GetByLong(longName string) (*FlagMeta, bool)   // 通过长标志名称查找对应的标志元数据
	GetByShort(shortName string) (*FlagMeta, bool) // 通过短标志名称查找对应的标志元数据
	GetByName(name string) (*FlagMeta, bool)       // 通过标志名称查找标志元数据
}
```

FlagRegistryInterface 标志注册表接口，定义了标志元数据的增删改查操作。

### FlagType

```go
type FlagType int
```

标志类型。

```go
const (
	FlagTypeInt      FlagType = iota + 1 // 整数类型
	FlagTypeInt64                        // 64位整数类型
	FlagTypeUint16                       // 16位无符号整数类型
	FlagTypeString                       // 字符串类型
	FlagTypeBool                         // 布尔类型
	FlagTypeFloat64                      // 64位浮点数类型
	FlagTypeEnum                         // 枚举类型
	FlagTypeDuration                     // 时间间隔类型
	FlagTypeSlice                        // 切片类型
	FlagTypeTime                         // 时间类型
	FlagTypeMap                          // 映射类型
	FlagTypePath                         // 路径类型
)
```

### Float64Flag

```go
type Float64Flag struct {
	BaseFlag[float64]
}
```

Float64Flag 浮点型标志结构体，继承BaseFlag[float64]泛型结构体，实现Flag接口。

```go
func (f *Float64Flag) Type() FlagType
```

**Type**：返回标志类型。

### Int64Flag

```go
type Int64Flag struct {
	BaseFlag[int64]
}
```

Int64Flag 64位整数类型标志结构体，继承BaseFlag[int64]泛型结构体，实现Flag接口。

```go
func (f *Int64Flag) SetRange(min, max int64)
func (f *Int64Flag) Type() FlagType
```

方法说明：
- **SetRange**：设置64位整数的有效范围。
- **Type**：返回标志类型。

### IntFlag

```go
type IntFlag struct {
	BaseFlag[int]
}
```

IntFlag 整数类型标志结构体，继承BaseFlag[int]泛型结构体，实现Flag接口。

```go
func (f *IntFlag) SetRange(min, max int)
func (f *IntFlag) Type() FlagType
```

方法说明：
- **SetRange**：设置整数的有效范围。
- **Type**：返回标志类型。

### MapFlag

```go
type MapFlag struct {
	BaseFlag[map[string]string]

	IgnoreCase bool // 是否忽略键的大小写
	// Has unexported fields.
}
```

MapFlag 键值对类型标志结构体，继承BaseFlag[map[string]string]泛型结构体，实现Flag接口。

```go
func (f *MapFlag) Set(value string) error
func (f *MapFlag) SetDelimiters(keyDelimiter, valueDelimiter string)
func (f *MapFlag) SetIgnoreCase(enable bool)
func (f *MapFlag) String() string
func (f *MapFlag) Type() FlagType
```

方法说明：
- **Set**：实现flag.Value接口，解析并设置键值对。
- **SetDelimiters**：设置键值对分隔符。
- **SetIgnoreCase**：设置是否忽略键的大小写。
- **String**：实现flag.Value接口，返回当前值的字符串表示。
- **Type**：返回标志类型。

### PathFlag

```go
type PathFlag struct {
	BaseFlag[string]
}
```

PathFlag 路径类型标志结构体，继承BaseFlag[string]泛型结构体，实现Flag接口。

```go
func (f *PathFlag) Init(longName, shortName string, defValue string, usage string) error
func (f *PathFlag) Set(value string) error
func (f *PathFlag) String() string
func (f *PathFlag) Type() FlagType
```

方法说明：
- **Init**：初始化路径标志。
- **Set**：实现flag.Value接口，解析并验证路径。
- **String**：实现flag.Value接口，返回当前值的字符串表示。
- **Type**：返回标志类型。

### SliceFlag

```go
type SliceFlag struct {
	BaseFlag[[]string] // 基类

	SkipEmpty bool // 是否跳过空元素
	// Has unexported fields.
}
```

SliceFlag 切片类型标志结构体，继承BaseFlag[[]string]泛型结构体，实现Flag接口。

```go
func (f *SliceFlag) Clear()
func (f *SliceFlag) Contains(element string) bool
func (f *SliceFlag) GetDelimiters() []string
func (f *SliceFlag) Init(longName, shortName string, defValue []string, usage string) error
func (f *SliceFlag) Len() int
func (f *SliceFlag) Remove(element string) error
func (f *SliceFlag) Set(value string) error
func (f *SliceFlag) SetDelimiters(delimiters []string)
func (f *SliceFlag) SetSkipEmpty(skip bool)
func (f *SliceFlag) Sort() error
func (f *SliceFlag) String() string
func (f *SliceFlag) Type() FlagType
```

方法说明：
- **Clear**：清空切片所有元素。注意：该方法会改变切片的指针。
- **Contains**：检查切片是否包含指定元素。当切片未设置值时，将使用默认值进行检查。
- **GetDelimiters**：获取当前分隔符列表。
- **Init**：初始化切片类型标志。
- **Len**：获取切片长度。
- **Remove**：从切片中移除指定元素（支持移除空字符串元素）。
- **Set**：实现flag.Value接口，解析并设置切片值。
- **SetDelimiters**：设置切片解析的分隔符列表。
- **SetSkipEmpty**：设置是否跳过空元素。
- **Sort**：对切片进行排序。
- **String**：实现flag.Value接口，返回当前值的字符串表示。
- **Type**：返回标志类型。

### StringFlag

```go
type StringFlag struct {
	BaseFlag[string]
}
```

StringFlag 字符串类型标志结构体，继承BaseFlag[string]泛型结构体，实现Flag接口。

```go
func (f *StringFlag) Contains(substr string) bool
func (f *StringFlag) Len() int
func (f *StringFlag) String() string
func (f *StringFlag) ToLower() string
func (f *StringFlag) ToUpper() string
func (f *StringFlag) Type() FlagType
```

方法说明：
- **Contains**：检查字符串是否包含指定子串。
- **Len**：获取字符串标志的长度。
- **String**：返回带引号的字符串值。
- **ToLower**：将字符串标志值转换为小写。
- **ToUpper**：将字符串标志值转换为大写。
- **Type**：返回标志类型。

### TimeFlag

```go
type TimeFlag struct {
	BaseFlag[time.Time]

	// Has unexported fields.
}
```

TimeFlag 时间类型标志结构体，继承BaseFlag[time.Time]泛型结构体，实现Flag接口。

```go
func (f *TimeFlag) Set(value string) error
func (f *TimeFlag) SetOutputFormat(format string)
func (f *TimeFlag) String() string
func (f *TimeFlag) Type() FlagType
```

方法说明：
- **Set**：实现flag.Value接口，解析并设置时间值。
- **SetOutputFormat**：设置时间输出格式。
- **String**：实现flag.Value接口，返回当前时间的字符串表示。
- **Type**：返回标志类型。

### TypedFlag[T any]

```go
type TypedFlag[T any] interface {
	Flag                    // 继承标志接口
	GetDefault() T          // 获取标志的具体类型默认值
	Get() T                 // 获取标志的具体类型值
	GetPointer() *T         // 获取标志值的指针
	Set(T) error            // 设置标志的具体类型值
	SetValidator(Validator) // 设置标志的验证器
}
```

TypedFlag 所有标志类型的通用接口，定义了标志的元数据访问方法和默认值访问方法。

### Uint16Flag

```go
type Uint16Flag struct {
	BaseFlag[uint16] // 基类
	// Has unexported fields.
}
```

Uint16Flag 16位无符号整数类型标志结构体，继承BaseFlag[uint16]泛型结构体，实现Flag接口。

```go
func (f *Uint16Flag) Set(value string) error
func (f *Uint16Flag) String() string
func (f *Uint16Flag) Type() FlagType
```

方法说明：
- **Set**：实现flag.Value接口，解析并设置16位无符号整数值。验证值是否在uint16范围内(0-65535)。
- **String**：实现flag.Value接口，返回当前值的字符串表示。
- **Type**：返回标志类型。

### Validator

```go
type Validator interface {
	// Validate 验证参数值是否合法
	// value: 待验证的参数值
	// 返回值: 验证通过返回nil, 否则返回错误信息
	Validate(value any) error
}
```

Validator 验证器接口，所有自定义验证器需实现此接口。
