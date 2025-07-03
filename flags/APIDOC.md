# Package flags

Package flags 定义了所有标志类型的通用接口和基础标志结构体。

```go
package flags // import "gitee.com/MM-Q/qflag/flags"
```

## 常量

定义标志的分隔符常量：

```go
const (
    FlagSplitComma     = "," // 逗号
    FlagSplitSemicolon = ";" // 分号
    FlagSplitPipe      = "|" // 竖线
    FlagKVColon        = ":" // 冒号
    FlagKVEqual        = "=" // 等号
)
```

定义非法字符集常量，防止非法的标志名称：

```go
const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"
```

## 变量

内置标志名称：

```go
var (
    HelpFlagName            = "help"    // 帮助标志名称
    HelpFlagShortName       = "h"       // 帮助标志短名称
    ShowInstallPathFlagName = "sip"     // 显示安装路径标志名称
    VersionFlagLongName     = "version" // 版本标志名称
    VersionFlagShortName    = "v"       // 版本标志短名称
)
```

Flag 支持的标志分隔符切片：

```go
var FlagSplitSlice = []string{
    FlagSplitComma,
    FlagSplitSemicolon,
    FlagSplitPipe,
    FlagKVColon,
}
```

## 函数

将 FlagType 转换为字符串：

```go
func FlagTypeToString(flagType FlagType) string
```

## 类型

### BaseFlag[T any]

泛型基础标志结构体，封装所有标志的通用字段和方法。

```go
type BaseFlag[T any] struct {
    // Has unexported fields.
}
```

方法：

```go
func (f *BaseFlag[T]) Get() T
func (f *BaseFlag[T]) GetDefault() T
func (f *BaseFlag[T]) GetDefaultAny() any
func (f *BaseFlag[T]) GetPointer() *T
func (f *BaseFlag[T]) Init(longName, shortName string, usage string, value *T) error
func (f *BaseFlag[T]) IsSet() bool
func (f *BaseFlag[T]) LongName() string
func (f *BaseFlag[T]) Reset()
func (f *BaseFlag[T]) Set(value T) error
func (f *BaseFlag[T]) SetValidator(validator Validator)
func (f *BaseFlag[T]) ShortName() string
func (f *BaseFlag[T]) String() string
func (f *BaseFlag[T]) Type() FlagType
func (f *BaseFlag[T]) Usage() string
```

### BoolFlag

布尔类型标志结构体，继承 BaseFlag[bool] 泛型结构体，实现 Flag 接口。

```go
type BoolFlag struct {
    BaseFlag[bool]
    // Has unexported fields.
}
```

方法：

```go
func (f *BoolFlag) IsBoolFlag() bool
func (f *BoolFlag) Set(value string) error
func (f *BoolFlag) String() string
func (f *BoolFlag) Type() FlagType
```

### DurationFlag

时间间隔类型标志结构体，继承 BaseFlag[time.Duration] 泛型结构体，实现 Flag 接口。

```go
type DurationFlag struct {
    BaseFlag[time.Duration]
    // Has unexported fields.
}
```

方法：

```go
func (f *DurationFlag) Set(value string) error
func (f *DurationFlag) String() string
func (f *DurationFlag) Type() FlagType
```

### EnumFlag

枚举类型标志结构体，继承 BaseFlag[string] 泛型结构体，增加枚举特有的选项验证。

```go
type EnumFlag struct {
    BaseFlag[string]
    // Has unexported fields.
}
```

方法：

```go
func (f *EnumFlag) Init(longName, shortName string, defValue string, usage string, options []string) error
func (f *EnumFlag) IsCheck(value string) error
func (f *EnumFlag) Set(value string) error
func (f *EnumFlag) String() string
func (f *EnumFlag) Type() FlagType
```

### Flag

所有标志类型的通用接口，定义了标志的元数据访问方法。

```go
type Flag interface {
    LongName() string   // 获取标志的长名称
    ShortName() string  // 获取标志的短名称
    Usage() string      // 获取标志的用法
    Type() FlagType     // 获取标志类型
    GetDefaultAny() any // 获取标志的默认值(any 类型)
    String() string     // 获取标志的字符串表示
    IsSet() bool        // 判断标志是否已设置值
    Reset()             // 重置标志值为默认值
}
```

### FlagMeta

统一存储标志的完整元数据。

```go
type FlagMeta struct {
    Flag Flag // 标志对象
}
```

方法：

```go
func (m *FlagMeta) GetDefault() any
func (m *FlagMeta) GetFlag() Flag
func (m *FlagMeta) GetFlagType() FlagType
func (m *FlagMeta) GetLongName() string
func (m *FlagMeta) GetShortName() string
func (m *FlagMeta) GetUsage() string
```

### FlagMetaInterface

标志元数据接口，定义了标志元数据的获取方法。

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

### FlagRegistry

集中管理所有标志元数据及索引。

```go
type FlagRegistry struct {
    // Has unexported fields.
}
```

方法：

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

### FlagRegistryInterface

标志注册表接口，定义了标志元数据的增删改查操作。

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

### FlagType

标志类型。

```go
type FlagType int
```

枚举值：

```go
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

浮点型标志结构体，继承 BaseFlag[float64] 泛型结构体，实现 Flag 接口。

```go
type Float64Flag struct {
    BaseFlag[float64]
    // Has unexported fields.
}
```

方法：

```go
func (f *Float64Flag) Set(value string) error
func (f *Float64Flag) Type() FlagType
```

### IP4Flag

IPv4 地址类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```go
type IP4Flag struct {
    BaseFlag[string]
    // Has unexported fields.
}
```

方法：

```go
func (f *IP4Flag) Set(value string) error
func (f *IP4Flag) String() string
func (f *IP4Flag) Type() FlagType
```

### IP6Flag

IPv6 地址类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```go
type IP6Flag struct {
    BaseFlag[string]
    // Has unexported fields.
}
```

方法：

```go
func (f *IP6Flag) Set(value string) error
func (f *IP6Flag) String() string
func (f *IP6Flag) Type() FlagType
```

### Int64Flag

64 位整数类型标志结构体，继承 BaseFlag[int64] 泛型结构体，实现 Flag 接口。

```go
type Int64Flag struct {
    BaseFlag[int64]
    // Has unexported fields.
}
```

方法：

```go
func (f *Int64Flag) Set(value string) error
func (f *Int64Flag) SetRange(min, max int64)
func (f *Int64Flag) Type() FlagType
```

### IntFlag

整数类型标志结构体，继承 BaseFlag[int] 泛型结构体，实现 Flag 接口。

```go
type IntFlag struct {
    BaseFlag[int]
    // Has unexported fields.
}
```

方法：

```go
func (f *IntFlag) Set(value string) error
func (f *IntFlag) SetRange(min, max int)
func (f *IntFlag) String() string
func (f *IntFlag) Type() FlagType
```

### MapFlag

键值对类型标志结构体，继承 BaseFlag[map[string]string] 泛型结构体，实现 Flag 接口。

```go
type MapFlag struct {
    BaseFlag[map[string]string]
    // Has unexported fields.
}
```

方法：

```go
func (f *MapFlag) Set(value string) error
func (f *MapFlag) SetDelimiters(keyDelimiter, valueDelimiter string)
func (f *MapFlag) SetIgnoreCase(enable bool)
func (f *MapFlag) String() string
func (f *MapFlag) Type() FlagType
```

### PathFlag

路径类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```go
type PathFlag struct {
    BaseFlag[string]
    // Has unexported fields.
}
```

方法：

```go
func (f *PathFlag) Init(longName, shortName string, defValue string, usage string) error
func (f *PathFlag) IsDirectory(isDir bool) *PathFlag
func (f *PathFlag) MustExist(mustExist bool) *PathFlag
func (f *PathFlag) Set(value string) error
func (f *PathFlag) String() string
func (f *PathFlag) Type() FlagType
```

### SliceFlag

切片类型标志结构体，继承 BaseFlag[[]string] 泛型结构体，实现 Flag 接口。

```go
type SliceFlag struct {
    BaseFlag[[]string] // 基类
    // Has unexported fields.
}
```

方法：

```go
func (f *SliceFlag) Clear() error
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

### StringFlag

字符串类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```go
type StringFlag struct {
    BaseFlag[string]
}
```

方法：

```go
func (f *StringFlag) Contains(substr string) bool
func (f *StringFlag) Len() int
func (f *StringFlag) Set(value string) error
func (f *StringFlag) String() string
func (f *StringFlag) ToLower() string
func (f *StringFlag) ToUpper() string
func (f *StringFlag) Type() FlagType
```

### TimeFlag

时间类型标志结构体，继承 BaseFlag[time.Time] 泛型结构体，实现 Flag 接口。

```go
type TimeFlag struct {
    BaseFlag[time.Time]
    // Has unexported fields.
}
```

方法：

```go
func (f *TimeFlag) Set(value string) error
func (f *TimeFlag) SetOutputFormat(format string)
func (f *TimeFlag) String() string
func (f *TimeFlag) Type() FlagType
```

### TypedFlag[T any]

所有标志类型的通用接口，定义了标志的元数据访问方法和默认值访问方法。

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

### URLFlag

URL 类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```go
type URLFlag struct {
    BaseFlag[string]
    // Has unexported fields.
}
```

方法：

```go
func (f *URLFlag) Set(value string) error
func (f *URLFlag) String() string
func (f *URLFlag) Type() FlagType
```

### Uint16Flag

16 位无符号整数类型标志结构体，继承 BaseFlag[uint16] 泛型结构体，实现 Flag 接口。

```go
type Uint16Flag struct {
    BaseFlag[uint16] // 基类
    // Has unexported fields.
}
```

方法：

```go
func (f *Uint16Flag) Set(value string) error
func (f *Uint16Flag) String() string
func (f *Uint16Flag) Type() FlagType
```

### Uint32Flag

32 位无符号整数类型标志结构体，继承 BaseFlag[uint32] 泛型结构体，实现 Flag 接口。

```go
type Uint32Flag struct {
    BaseFlag[uint32] // 基类
    // Has unexported fields.
}
```

方法：

```go
func (f *Uint32Flag) Set(value string) error
func (f *Uint32Flag) String() string
func (f *Uint32Flag) Type() FlagType
```

### Uint64Flag

64 位无符号整数类型标志结构体，继承 BaseFlag[uint64] 泛型结构体，实现 Flag 接口。

```go
type Uint64Flag struct {
    BaseFlag[uint64] // 基类
    // Has unexported fields.
}
```

方法：

```go
func (f *Uint64Flag) Set(value string) error
func (f *Uint64Flag) String() string
func (f *Uint64Flag) Type() FlagType
```

### Validator

验证器接口，所有自定义验证器需实现此接口。

```go
type Validator interface {
    // Validate 验证参数值是否合法
    // value: 待验证的参数值
    // 返回值: 验证通过返回 nil, 否则返回错误信息
    Validate(value any) error
}
```