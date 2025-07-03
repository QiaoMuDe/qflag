# Package validator

Package validator 提供常用的参数验证器实现。

```go
package validator // import "gitee.com/MM-Q/qflag/validator"
```

## 类型

### BoolValidator

验证布尔值。

```go
type BoolValidator struct{}
```

方法：

```go
func (v *BoolValidator) Validate(value any) error
```

### DurationValidator

验证时间间隔是否有效。

```go
type DurationValidator struct {
    Min time.Duration // 最小时间间隔, 包含在内
    Max time.Duration // 最大时间间隔, 包含在内
}
```

方法：

```go
func (v *DurationValidator) Validate(value any) error
```

### EnumValidator

验证值是否在枚举列表中。

```go
type EnumValidator struct {
    AllowedValues []any // 允许的值列表
}
```

方法：

```go
func (v *EnumValidator) Validate(value any) error
```

### FloatRangeValidator

验证浮点数是否在指定范围内。

```go
type FloatRangeValidator struct {
    Min float64 // 最小值, 包含在内
    Max float64 // 最大值, 包含在内
}
```

方法：

```go
func (v *FloatRangeValidator) Validate(value any) error
```

### IntRangeValidator

验证整数是否在指定范围内。

```go
type IntRangeValidator struct {
    Min int // 最小值, 包含在内
    Max int // 最大值, 包含在内
}
```

方法：

```go
func (v *IntRangeValidator) Validate(value any) error
```

### IntRangeValidator64

验证整数是否在指定范围内。

```go
type IntRangeValidator64 struct {
    Min int64 // 最小值, 包含在内
    Max int64 // 最大值, 包含在内
}
```

方法：

```go
func (v *IntRangeValidator64) Validate(value any) error
```

### IntValueValidator

验证整数是否为指定值之一。

```go
type IntValueValidator struct {
    AllowedValues []int // 允许的整数值列表
}
```

方法：

```go
func (v *IntValueValidator) Validate(value any) error
```

### PathValidator

路径验证器。

```go
type PathValidator struct {
    MustExist   bool // 是否必须存在，默认true为
    IsDirectory bool // 是否必须是目录，默认为false
}
```

方法：

```go
func (v *PathValidator) Validate(value any) error
```

### SliceLengthValidator

验证切片长度是否在指定范围内。

```go
type SliceLengthValidator struct {
    Min int // 最小长度, 包含在内
    Max int // 最大长度, 包含在内, 0表示不限制
}
```

方法：

```go
func (v *SliceLengthValidator) Validate(value any) error
```

### StringLengthValidator

验证字符串长度是否在指定范围内。

```go
type StringLengthValidator struct {
    Min int // 最小长度, 包含在内
    Max int // 最大长度, 包含在内, 0表示不限制
}
```

方法：

```go
func (v *StringLengthValidator) Validate(value any) error
```

### StringRegexValidator

验证字符串是否匹配正则表达式。

```go
type StringRegexValidator struct {
    Pattern string         // 正则表达式模式
    Regex   *regexp.Regexp // 编译后的正则表达式
}
```

方法：

```go
func (v *StringRegexValidator) Validate(value any) error
```