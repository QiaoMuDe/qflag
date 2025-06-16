// flags 定义了所有标志类型的通用接口和基础标志结构体
package qflag

import (
	"fmt"
	"strings"
	"sync"
)

// Flag 所有标志类型的通用接口,定义了标志的元数据访问方法
type Flag interface {
	LongName() string   // 获取标志的长名称
	ShortName() string  // 获取标志的短名称
	Usage() string      // 获取标志的用法
	Type() FlagType     // 获取标志类型
	GetDefaultAny() any // 获取标志的默认值
}

// TypedFlag 所有标志类型的通用接口,定义了标志的元数据访问方法和默认值访问方法
type TypedFlag[T any] interface {
	Flag           // 继承标志接口
	GetDefault() T // 获取标志的具体类型默认值
	Get() T        // 获取标志的具体类型值
	Set(T) error   // 设置标志的具体类型值
}

// BaseFlag 泛型基础标志结构体,封装所有标志的通用字段和方法
type BaseFlag[T any] struct {
	cmd       *Cmd       // 所属的命令实例
	longName  string     // 长标志名称
	shortName string     // 短标志字符
	defValue  T          // 默认值
	usage     string     // 帮助说明
	value     *T         // 标志值指针
	mu        sync.Mutex // 并发访问锁
}

// LongName 获取标志的长名称
func (f *BaseFlag[T]) LongName() string { return f.longName }

// ShortName 获取标志的短名称
func (f *BaseFlag[T]) ShortName() string { return f.shortName }

// Usage 获取标志的用法说明
func (f *BaseFlag[T]) Usage() string { return f.usage }

// GetDefault 获取标志的默认值
func (f *BaseFlag[T]) GetDefault() T { return f.defValue }

// GetDefaultAny 获取标志的默认值(any类型)
func (f *BaseFlag[T]) GetDefaultAny() any { return f.defValue }

// Get 获取标志的实际值
func (f *BaseFlag[T]) Get() T {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.value != nil {
		return *f.value
	}
	return f.defValue
}

// GetValueAny 获取标志的实际值(any类型)
func (f *BaseFlag[T]) GetValueAny() any {
	return f.Get()
}

// Set 设置标志的值
func (f *BaseFlag[T]) Set(value T) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 创建一个副本
	v := value
	f.value = &v
	return nil
}

// IntFlag 整数类型标志结构体
// 继承BaseFlag[int]泛型结构体,实现Flag接口
type IntFlag struct {
	BaseFlag[int]
}

// Type 返回标志类型
func (f *IntFlag) Type() FlagType { return FlagTypeInt }

// StringFlag 字符串类型标志结构体
// 继承BaseFlag[string]泛型结构体,实现Flag接口
type StringFlag struct {
	BaseFlag[string]
}

// Type 返回标志类型
func (f *StringFlag) Type() FlagType { return FlagTypeString }

// BoolFlag 布尔类型标志结构体
// 继承BaseFlag[bool]泛型结构体,实现Flag接口
type BoolFlag struct {
	BaseFlag[bool]
}

// Type 返回标志类型
func (f *BoolFlag) Type() FlagType { return FlagTypeBool }

// FloatFlag 浮点型标志结构体
// 继承BaseFlag[float64]泛型结构体,实现Flag接口
type FloatFlag struct {
	BaseFlag[float64]
}

// Type 返回标志类型
func (f *FloatFlag) Type() FlagType { return FlagTypeFloat }

// EnumFlag 枚举类型标志结构体
// 继承BaseFlag[string]泛型结构体,增加枚举特有的选项验证
type EnumFlag struct {
	BaseFlag[string]
	optionMap map[string]bool // 允许的枚举值映射
}

// 实现Flag接口
func (f *EnumFlag) Type() FlagType { return FlagTypeEnum }

// Check 检查标志值是否在允许的枚举范围内
func (f *EnumFlag) Check(value string) error {
	// 如果枚举map为空,则不需要检查
	if len(f.optionMap) == 0 {
		return nil
	}

	// 转换为小写
	value = strings.ToLower(value)

	// 检查值是否在枚举map中
	if _, valid := f.optionMap[value]; !valid {
		var options []string
		for k := range f.optionMap {
			options = append(options, k)
		}
		return fmt.Errorf("invalid enum value '%s', options are %v", value, options)
	}
	return nil
}

// Set 实现flag.Value接口, 解析并设置枚举值
func (f *EnumFlag) Set(value string) error {
	// 先验证值是否有效
	if err := f.Check(value); err != nil {
		return err
	}
	// 调用基类方法设置值
	return f.BaseFlag.Set(value)
}

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *EnumFlag) String() string { return f.Get() }
