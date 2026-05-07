# 泛型验证器设计方案

## 概述

本文档提供了使用泛型重构验证器的完整方案，解决当前使用 `any` 类型导致的类型安全性和性能问题。

## 1. 核心接口定义

### 1.1 泛型验证器接口

```go
// Validator 泛型验证器接口
//
// Validator 定义了类型安全的值验证接口, 使用泛型确保类型安全。
//
// 使用场景:
//   - 标志值范围验证
//   - 格式检查 (如邮箱、URL等)
//   - 业务规则验证
//   - 复杂条件判断
//
// 实现说明:
//   - 实现类型应该是线程安全的
//   - 验证逻辑应该尽可能高效
//   - 错误信息应该清晰明确
//   - 泛型参数 T 确保类型安全
type Validator[T any] interface {
	// Validate 验证给定值是否符合要求
	//
	// 参数:
	//   - value: 要验证的值, 类型为 T
	//
	// 返回值:
	//   - error: 验证失败时返回错误, 成功时返回nil
	//
	// 功能说明:
	//   - 对传入值进行业务规则验证
	//   - 返回详细的错误信息
	//   - 类型安全, 无需类型断言
	Validate(value T) error
}
```

## 2. BaseFlag 修改

### 2.1 修改后的 BaseFlag 结构

```go
// BaseFlag 泛型基础标志结构体
type BaseFlag[T any] struct {
	mu         sync.RWMutex      // 读写锁
	value      *T                // 当前值指针
	default_   T                 // 默认值
	isSet      bool             // 标志是否已被设置
	validator_ Validator[T]     // 泛型值验证器
	envVar     string           // 关联的环境变量名

	// 不可变属性, 无需挂锁
	longName  string         // 长选项名称
	shortName string         // 短选项名称
	desc      string         // 标志描述信息
	flagType  types.FlagType // 标志类型枚举值
}
```

### 2.2 修改后的 SetValidator 方法

```go
// SetValidator 设置标志的验证器
//
// 参数:
//   - v: 验证器实例, 实现泛型 Validator[T] 接口
func (f *BaseFlag[T]) SetValidator(v Validator[T]) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.validator_ = v
}
```

### 2.3 修改后的 Validate 方法

```go
// Validate 验证标志的当前值
//
// 返回值:
//   - error: 如果验证失败返回错误, 否则返回nil
//
// 注意事项:
//   - 如果设置了验证器, 则使用验证器验证当前值
//   - 如果未设置验证器, 则直接返回nil
func (f *BaseFlag[T]) Validate() error {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// 快速失败：没有验证器时直接返回
	if f.validator_ == nil {
		return nil
	}

	// 调用验证器验证当前值，类型安全，无需断言
	return f.validator_.Validate(f.Get())
}
```

## 3. 验证器实现示例

### 3.1 字符串验证器

```go
// StringValidator 字符串验证器
type StringValidator struct {
	// MinLength 最小长度
	MinLength int
	// MaxLength 最大长度
	MaxLength int
	// Pattern 正则表达式模式
	Pattern string
	// regex 编译后的正则表达式
	regex *regexp.Regexp
}

// Validate 验证字符串
func (v *StringValidator) Validate(value string) error {
	// 检查最小长度
	if v.MinLength > 0 && len(value) < v.MinLength {
		return fmt.Errorf("string length %d is smaller than minimum %d", len(value), v.MinLength)
	}

	// 检查最大长度
	if v.MaxLength > 0 && len(value) > v.MaxLength {
		return fmt.Errorf("string length %d is larger than maximum %d", len(value), v.MaxLength)
	}

	// 检查正则表达式
	if v.Pattern != "" {
		if v.regex == nil {
			regex, err := regexp.Compile(v.Pattern)
			if err != nil {
				return fmt.Errorf("invalid regex pattern: %v", err)
			}
			v.regex = regex
		}

		if !v.regex.MatchString(value) {
			return fmt.Errorf("string does not match pattern '%s': %s", v.Pattern, value)
		}
	}

	return nil
}

// 便捷创建函数
func String(minLength, maxLength int, pattern string) *StringValidator {
	return &StringValidator{
		MinLength: minLength,
		MaxLength: maxLength,
		Pattern:   pattern,
	}
}
```

### 3.2 整数验证器

```go
// IntValidator 整数验证器
type IntValidator struct {
	// Min 最小值 (包含)
	Min int
	// Max 最大值 (包含)
	Max int
}

// Validate 验证整数
func (v *IntValidator) Validate(value int) error {
	if value < v.Min {
		return fmt.Errorf("value %d is smaller than minimum %d", value, v.Min)
	}

	if value > v.Max {
		return fmt.Errorf("value %d is larger than maximum %d", value, v.Max)
	}

	return nil
}

// 便捷创建函数
func IntRange(min, max int) *IntValidator {
	return &IntValidator{Min: min, Max: max}
}
```

### 3.3 切片长度验证器

```go
// SliceLengthValidator 切片长度验证器
type SliceLengthValidator[T any] struct {
	// Min 最小长度
	Min int
	// Max 最大长度
	Max int
}

// Validate 验证切片长度
func (v *SliceLengthValidator[T]) Validate(value []T) error {
	length := len(value)

	if v.Min > 0 && length < v.Min {
		return fmt.Errorf("slice length %d is smaller than minimum %d", length, v.Min)
	}

	if v.Max > 0 && length > v.Max {
		return fmt.Errorf("slice length %d is larger than maximum %d", length, v.Max)
	}

	return nil
}

// 便捷创建函数
func LengthRange[T any](min, max int) *SliceLengthValidator[T] {
	return &SliceLengthValidator[T]{Min: min, Max: max}
}
```

### 3.4 Map 长度验证器

```go
// MapLengthValidator Map 长度验证器
type MapLengthValidator[K comparable, V any] struct {
	// Min 最小长度
	Min int
	// Max 最大长度
	Max int
}

// Validate 验证 Map 长度
func (v *MapLengthValidator[K, V]) Validate(value map[K]V) error {
	length := len(value)

	if v.Min > 0 && length < v.Min {
		return fmt.Errorf("map length %d is smaller than minimum %d", length, v.Min)
	}

	if v.Max > 0 && length > v.Max {
		return fmt.Errorf("map length %d is larger than maximum %d", length, v.Max)
	}

	return nil
}

// 便捷创建函数
func MapLengthRange[K comparable, V any](min, max int) *MapLengthValidator[K, V] {
	return &MapLengthValidator[K, V]{Min: min, Max: max}
}
```

### 3.5 文件存在验证器

```go
// FileExistsValidator 文件存在性验证器
type FileExistsValidator struct {
	// AllowEmpty 是否允许空值
	AllowEmpty bool
}

// Validate 验证文件是否存在
func (v *FileExistsValidator) Validate(value string) error {
	// 允许空值
	if value == "" {
		if v.AllowEmpty {
			return nil
		}
		return fmt.Errorf("file path cannot be empty")
	}

	// 检查文件是否存在
	if _, err := os.Stat(value); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", value)
	}

	return nil
}

// 便捷创建函数
func FileExists(allowEmpty bool) *FileExistsValidator {
	return &FileExistsValidator{AllowEmpty: allowEmpty}
}
```

## 4. 使用示例

### 4.1 StringFlag 使用验证器

```go
// 创建字符串标志并设置验证器
flag := NewStringFlag("email", "e", "邮箱地址", "")

// 设置验证器：要求邮箱格式，长度 5-100
flag.SetValidator(&StringValidator{
	MinLength: 5,
	MaxLength: 100,
	Pattern:   `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
})

// 使用便捷函数
flag.SetValidator(String(5, 100, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))
```

### 4.2 IntFlag 使用验证器

```go
// 创建整数标志并设置验证器
flag := NewIntFlag("port", "p", "端口号", 8080)

// 设置验证器：要求端口在 1-65535 范围内
flag.SetValidator(IntRange(1, 65535))
```

### 4.3 StringSliceFlag 使用验证器

```go
// 创建字符串切片标志并设置验证器
flag := NewStringSliceFlag("paths", "p", "路径列表", nil)

// 设置验证器：要求至少 1 个元素，最多 5 个元素
flag.SetValidator(LengthRange[string](1, 5))
```

### 4.4 MapFlag 使用验证器

```go
// 创建 Map 标志并设置验证器
flag := NewMapFlag("headers", "h", "HTTP头部", nil)

// 设置验证器：要求至少 1 个键值对，最多 3 个键值对
flag.SetValidator(MapLengthRange[string, string](1, 3))
```

## 5. 组合验证器

### 5.1 组合验证器实现

```go
// CompositeValidator 组合验证器
//
// CompositeValidator 允许组合多个验证器, 所有验证器都必须通过。
type CompositeValidator[T any] struct {
	validators []Validator[T]
}

// Validate 验证值
func (v *CompositeValidator[T]) Validate(value T) error {
	for _, validator := range v.validators {
		if err := validator.Validate(value); err != nil {
			return err
		}
	}
	return nil
}

// Add 添加验证器
func (v *CompositeValidator[T]) Add(validator Validator[T]) {
	v.validators = append(v.validators, validator)
}

// 便捷创建函数
func Compose[T any](validators ...Validator[T]) Validator[T] {
	return &CompositeValidator[T]{validators: validators}
}
```

### 5.2 组合验证器使用示例

```go
// 创建组合验证器
flag := NewStringFlag("username", "u", "用户名", "")

// 组合多个验证器
flag.SetValidator(Compose[string](
	&StringValidator{MinLength: 3, MaxLength: 20},
	&StringValidator{Pattern: `^[a-zA-Z0-9_]+$`},
))
```

## 6. 向后兼容方案（可选）

如果需要保持向后兼容，可以保留旧的 `Validator` 接口：

```go
// Validator 旧版验证器接口 (保持向后兼容)
type Validator interface {
	Validate(value any) error
}

// Validator 泛型验证器接口 (推荐使用)
type Validator[T any] interface {
	Validate(value T) error
}

// ValidatorAdapter 适配器，将泛型验证器适配为旧接口
type ValidatorAdapter[T any] struct {
	inner Validator[T]
}

func (a *ValidatorAdapter[T]) Validate(value any) error {
	typed, ok := value.(T)
	if !ok {
		return fmt.Errorf("expected type %T, got %T", *new(T), value)
	}
	return a.inner.Validate(typed)
}

// WrapValidator 包装泛型验证器为旧接口
func WrapValidator[T any](v Validator[T]) Validator {
	return &ValidatorAdapter[T]{inner: v}
}
```

## 7. 迁移步骤

### 7.1 修改 types 包

1. 修改 `internal/types/validator.go`，添加泛型验证器接口
2. 保留旧接口（如果需要向后兼容）

### 7.2 修改 flag 包

1. 修改 `internal/flag/base_flag.go`，将 `validator_` 字段类型改为 `Validator[T]`
2. 更新 `SetValidator` 方法签名
3. 更新 `Validate` 方法实现

### 7.3 重构 validator 包

1. 将所有验证器改为泛型实现
2. 移除类型断言代码
3. 更新便捷创建函数

### 7.4 更新测试用例

1. 更新所有使用验证器的测试用例
2. 确保类型安全

## 8. 优势总结

### 8.1 类型安全

- 编译时类型检查
- 无需运行时类型断言
- 避免 panic 风险

### 8.2 性能优化

- 消除类型断言开销
- 更好的编译器优化
- 减少运行时检查

### 8.3 代码质量

- 更清晰的 API
- 更好的 IDE 支持
- 更容易维护

### 8.4 扩展性

- 支持任意类型的验证器
- 易于组合验证器
- 灵活的验证逻辑

## 9. 注意事项

1. **泛型参数推断**：Go 编译器通常可以自动推断泛型参数，但有时需要显式指定
2. **空接口**：如果验证器需要支持多种类型，可以考虑使用接口约束
3. **性能考虑**：泛型在 Go 1.18+ 中性能已经很好，但在极端性能场景下需要测试
4. **向后兼容**：如果需要保持向后兼容，可以使用适配器模式

## 10. 完整示例代码

### 10.1 types/validator.go

```go
package types

// Validator 泛型验证器接口
type Validator[T any] interface {
	Validate(value T) error
}
```

### 10.2 validator/string_validator.go

```go
package validator

import (
	"fmt"
	"regexp"
)

// StringValidator 字符串验证器
type StringValidator struct {
	MinLength int
	MaxLength int
	Pattern   string
	regex     *regexp.Regexp
}

func (v *StringValidator) Validate(value string) error {
	if v.MinLength > 0 && len(value) < v.MinLength {
		return fmt.Errorf("string length %d is smaller than minimum %d", len(value), v.MinLength)
	}

	if v.MaxLength > 0 && len(value) > v.MaxLength {
		return fmt.Errorf("string length %d is larger than maximum %d", len(value), v.MaxLength)
	}

	if v.Pattern != "" {
		if v.regex == nil {
			regex, err := regexp.Compile(v.Pattern)
			if err != nil {
				return fmt.Errorf("invalid regex pattern: %v", err)
			}
			v.regex = regex
		}

		if !v.regex.MatchString(value) {
			return fmt.Errorf("string does not match pattern '%s': %s", v.Pattern, value)
		}
	}

	return nil
}

func String(minLength, maxLength int, pattern string) *StringValidator {
	return &StringValidator{
		MinLength: minLength,
		MaxLength: maxLength,
		Pattern:   pattern,
	}
}
```

### 10.3 validator/int_validator.go

```go
package validator

import "fmt"

// IntValidator 整数验证器
type IntValidator struct {
	Min int
	Max int
}

func (v *IntValidator) Validate(value int) error {
	if value < v.Min {
		return fmt.Errorf("value %d is smaller than minimum %d", value, v.Min)
	}

	if value > v.Max {
		return fmt.Errorf("value %d is larger than maximum %d", value, v.Max)
	}

	return nil
}

func IntRange(min, max int) *IntValidator {
	return &IntValidator{Min: min, Max: max}
}
```

### 10.4 validator/slice_validator.go

```go
package validator

import "fmt"

// SliceLengthValidator 切片长度验证器
type SliceLengthValidator[T any] struct {
	Min int
	Max int
}

func (v *SliceLengthValidator[T]) Validate(value []T) error {
	length := len(value)

	if v.Min > 0 && length < v.Min {
		return fmt.Errorf("slice length %d is smaller than minimum %d", length, v.Min)
	}

	if v.Max > 0 && length > v.Max {
		return fmt.Errorf("slice length %d is larger than maximum %d", length, v.Max)
	}

	return nil
}

func LengthRange[T any](min, max int) *SliceLengthValidator[T] {
	return &SliceLengthValidator[T]{Min: min, Max: max}
}
```

### 10.5 validator/map_validator.go

```go
package validator

import "fmt"

// MapLengthValidator Map 长度验证器
type MapLengthValidator[K comparable, V any] struct {
	Min int
	Max int
}

func (v *MapLengthValidator[K, V]) Validate(value map[K]V) error {
	length := len(value)

	if v.Min > 0 && length < v.Min {
		return fmt.Errorf("map length %d is smaller than minimum %d", length, v.Min)
	}

	if v.Max > 0 && length > v.Max {
		return fmt.Errorf("map length %d is larger than maximum %d", length, v.Max)
	}

	return nil
}

func MapLengthRange[K comparable, V any](min, max int) *MapLengthValidator[K, V] {
	return &MapLengthValidator[K, V]{Min: min, Max: max}
}
```

### 10.6 validator/composite_validator.go

```go
package validator

// CompositeValidator 组合验证器
type CompositeValidator[T any] struct {
	validators []Validator[T]
}

func (v *CompositeValidator[T]) Validate(value T) error {
	for _, validator := range v.validators {
		if err := validator.Validate(value); err != nil {
			return err
		}
	}
	return nil
}

func (v *CompositeValidator[T]) Add(validator Validator[T]) {
	v.validators = append(v.validators, validator)
}

func Compose[T any](validators ...Validator[T]) Validator[T] {
	return &CompositeValidator[T]{validators: validators}
}
```

## 总结

这个泛型验证器方案提供了类型安全、高性能、易扩展的验证器实现。通过使用 Go 的泛型特性，我们可以在编译时确保类型安全，避免运行时类型断言，同时保持代码的简洁性和可读性。
