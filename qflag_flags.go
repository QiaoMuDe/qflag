package qflag

import (
	"fmt"
	"strings"
	"sync"
)

// Flag 所有标志类型的通用接口,定义了标志的元数据访问方法
type Flag interface {
	Name() string      // 获取标志的名称
	ShortName() string // 获取标志的短名称
	Usage() string     // 获取标志的用法
	Type() FlagType    // 获取标志类型
}

// TypedFlag 所有标志类型的通用接口,定义了标志的元数据访问方法和默认值访问方法
type TypedFlag[T any] interface {
	Flag           // 继承标志接口
	GetDefault() T // 获取标志的默认值
	GetValue() T   // 获取标志的实际值
	SetValue(T)    // 设置标志的值
}

// IntFlag 整数类型标志结构体,包含标志元数据和值访问接口
type IntFlag struct {
	cmd       *Cmd       // 所属的命令实例
	name      string     // 长标志名称（如"port"）
	shortName string     // 短标志字符（如"p",空表示无短标志）
	defValue  int        // 默认值
	usage     string     // 帮助说明
	value     *int       // 标志值指针,通过flag库绑定
	mu        sync.Mutex // 并发访问锁
}

// 实现Flag接口
func (f *IntFlag) Name() string      { return f.name }
func (f *IntFlag) ShortName() string { return f.shortName }
func (f *IntFlag) Usage() string     { return f.usage }
func (f *IntFlag) GetDefault() int   { return f.defValue }
func (f *IntFlag) Type() FlagType    { return FlagTypeInt }

// GetValue 获取标志的实际值（带线程安全保护）
// 返回值优先级：解析值 > 默认值
func (f *IntFlag) GetValue() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.value != nil { // 优先返回解析值
		return *f.value
	}
	return f.defValue // 其次返回默认值
}

// SetValue 设置标志的值（带线程安全保护）
func (f *IntFlag) SetValue(value int) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.value = &value
}

// StringFlag 字符串类型标志结构体
type StringFlag struct {
	cmd       *Cmd       // 所属的命令实例
	name      string     // 长标志名称
	shortName string     // 短标志字符
	defValue  string     // 默认值
	usage     string     // 帮助说明
	value     *string    // 标志值指针
	mu        sync.Mutex // 并发访问锁
}

// 实现Flag接口
func (f *StringFlag) Name() string       { return f.name }
func (f *StringFlag) ShortName() string  { return f.shortName }
func (f *StringFlag) Usage() string      { return f.usage }
func (f *StringFlag) GetDefault() string { return f.defValue }
func (f *StringFlag) Type() FlagType     { return FlagTypeString }

// GetValue 获取标志的实际值（带线程安全保护）
// 返回值优先级：解析值 > 默认值
func (f *StringFlag) GetValue() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.value != nil { // 优先返回解析值
		return *f.value
	}
	return f.defValue // 其次返回默认值
}

// SetValue 设置标志的值（带线程安全保护）
func (f *StringFlag) SetValue(value string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.value = &value
}

// BoolFlag 布尔类型标志结构体
type BoolFlag struct {
	cmd       *Cmd       // 所属的命令实例
	name      string     // 长标志名称
	shortName string     // 短标志字符
	defValue  bool       // 默认值
	usage     string     // 帮助说明
	value     *bool      // 标志值指针
	mu        sync.Mutex // 并发访问锁
}

// 实现Flag接口
func (f *BoolFlag) Name() string      { return f.name }
func (f *BoolFlag) ShortName() string { return f.shortName }
func (f *BoolFlag) Usage() string     { return f.usage }
func (f *BoolFlag) GetDefault() bool  { return f.defValue }
func (f *BoolFlag) Type() FlagType    { return FlagTypeBool }

// GetValue 获取标志的实际值（带线程安全保护）
// 返回值优先级：解析值 > 默认值
func (f *BoolFlag) GetValue() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.value != nil { // 优先返回解析值
		return *f.value
	}
	return f.defValue // 其次返回默认值
}

// SetValue 设置标志的值（带线程安全保护）
func (f *BoolFlag) SetValue(value bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.value = &value
}

// FloatFlag 浮点型标志结构体
type FloatFlag struct {
	cmd       *Cmd       // 所属的命令实例
	name      string     // 长标志名称
	shortName string     // 短标志字符
	defValue  float64    // 默认值
	usage     string     // 帮助说明
	value     *float64   // 标志值指针
	mu        sync.Mutex // 并发访问锁
}

// 实现Flag接口
func (f *FloatFlag) Name() string        { return f.name }
func (f *FloatFlag) ShortName() string   { return f.shortName }
func (f *FloatFlag) Usage() string       { return f.usage }
func (f *FloatFlag) GetDefault() float64 { return f.defValue }
func (f *FloatFlag) Type() FlagType      { return FlagTypeFloat }

// GetValue 获取标志的实际值（带线程安全保护）
// 返回值优先级：解析值 > 默认值
func (f *FloatFlag) GetValue() float64 {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.value != nil { // 优先返回解析值
		return *f.value
	}
	return f.defValue // 其次返回默认值
}

// SetValue 设置标志的值（带线程安全保护）
func (f *FloatFlag) SetValue(value float64) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.value = &value
}

// SliceFlag 表示字符串切片类型的命令行标志
type SliceFlag struct {
	cmd        *Cmd       // 命令对象
	name       string     // 长标志名
	shortName  string     // 短标志名
	defValue   []string   // 默认值
	usage      string     // 帮助说明
	value      *[]string  // 标志值指针
	mu         sync.Mutex // 并发访问锁
	hasBeenSet bool       // 标记是否已设置值
}

// 实现Flag接口
func (f *SliceFlag) Name() string         { return f.name }
func (f *SliceFlag) ShortName() string    { return f.shortName }
func (f *SliceFlag) Usage() string        { return f.usage }
func (f *SliceFlag) GetDefault() []string { return f.defValue }
func (f *SliceFlag) Type() FlagType       { return FlagTypeSlice }

// GetValue 获取标志的实际值（带线程安全保护）
// 返回值优先级：解析值 > 默认值
func (f *SliceFlag) GetValue() []string {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.value != nil { // 优先返回解析值
		return *f.value
	}
	return f.defValue // 其次返回默认值
}

// SetValue 设置标志的值（带线程安全保护）
func (f *SliceFlag) SetValue(value ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.hasBeenSet {
		// 首次设置，清空现有值（包括默认值）
		*f.value = make([]string, 0)
		f.hasBeenSet = true
	}

	// 添加多个值
	*f.value = append(*f.value, value...)
}

// Set 实现flag.Value接口, 解析并添加值到切片
func (f *SliceFlag) Set(value string) error {
	f.SetValue(value)
	return nil
}

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *SliceFlag) String() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return strings.Join(*f.value, ",")
}

// EnumFlag 枚举类型标志结构体
type EnumFlag struct {
	cmd       *Cmd            // 所属的命令实例
	name      string          // 长标志名称
	shortName string          // 短标志字符
	defValue  string          // 默认值
	usage     string          // 帮助说明
	value     *string         // 标志值指针
	options   []string        // 允许的枚举值列表
	optionMap map[string]bool // 允许的枚举值映射
	mu        sync.Mutex      // 并发访问锁
}

// 实现Flag接口
func (f *EnumFlag) Name() string       { return f.name }
func (f *EnumFlag) ShortName() string  { return f.shortName }
func (f *EnumFlag) Usage() string      { return f.usage }
func (f *EnumFlag) GetDefault() string { return f.defValue }
func (f *EnumFlag) Type() FlagType     { return FlagTypeEnum }

// GetValue 获取标志的实际值（带线程安全保护）
// 返回值优先级：解析值 > 默认值
func (f *EnumFlag) GetValue() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.value != nil { // 优先返回解析值
		return *f.value
	}
	return f.defValue // 其次返回默认值
}

// SetValue 设置标志的值（带线程安全保护）
func (f *EnumFlag) SetValue(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 使用map快速验证值是否在允许的枚举范围内
	if _, valid := f.optionMap[value]; !valid {
		return fmt.Errorf("invalid enum value '%s', options are %v", value, f.options)
	}

	// 设置当前值
	f.value = &value
	return nil
}

// Set 实现flag.Value接口, 解析并设置枚举值
func (f *EnumFlag) Set(value string) error {
	return f.SetValue(value)
}

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *EnumFlag) String() string {
	return f.GetValue()
}
