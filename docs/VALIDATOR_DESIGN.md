# 验证器功能设计方案

## 一、设计概述

### 核心思路
1. **创建独立的验证器类型**（泛型函数类型）
2. **在 `BaseFlag[T]` 中添加单个验证器字段**（不是切片）
3. **重复设置验证器会覆盖之前的验证器**
4. **在具体 Flag 的 `Set` 方法中调用验证**
5. **不修改 `Flag` 接口**

---

## 二、核心设计

### 2.1 验证器类型定义

```go
// Validator 验证器函数类型
//
// Validator 是一个泛型函数类型，用于验证标志值的有效性。
// 验证器接收一个类型为 T 的值，返回错误信息。
//
// 参数:
//   - value: 要验证的值
//
// 返回值:
//   - error: 验证失败时返回错误，验证通过返回 nil
type Validator[T any] func(value T) error
```

### 2.2 BaseFlag[T] 修改

```go
// BaseFlag 泛型基础标志结构体
type BaseFlag[T any] struct {
    mu       sync.RWMutex // 读写锁
    value    *T           // 当前值指针
    default_ T            // 默认值
    isSet    bool         // 标志是否已被设置
    envVar   string       // 关联的环境变量名

    // 新增：验证器（单个，不是切片）
    validator Validator[T]

    // 不可变属性，无需挂锁
    longName  string         // 长选项名称
    shortName string         // 短选项名称
    desc      string         // 标志描述信息
    flagType  types.FlagType // 标志类型枚举值
}

// SetValidator 设置验证器
//
// 参数:
//   - validator: 验证器函数
//
// 功能说明:
//   - 设置标志的验证器
//   - 如果之前已设置验证器，会被覆盖
//   - 验证器会在 Set 方法中解析完值后被调用
func (f *BaseFlag[T]) SetValidator(validator Validator[T]) {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.validator = validator
}

// ClearValidator 清除验证器
//
// 功能说明:
//   - 移除标志的验证器
//   - 之后调用 Set 方法将不会进行验证
func (f *BaseFlag[T]) ClearValidator() {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.validator = nil
}

// HasValidator 检查是否设置了验证器
//
// 返回值:
//   - bool: 是否设置了验证器
func (f *BaseFlag[T]) HasValidator() bool {
    f.mu.RLock()
    defer f.mu.RUnlock()
    return f.validator != nil
}

// validate 执行验证（私有方法）
//
// 参数:
//   - value: 要验证的值
//
// 返回值:
//   - error: 验证失败时返回错误
//
// 功能说明:
//   - 调用设置的验证器进行验证
//   - 如果未设置验证器，直接返回 nil
func (f *BaseFlag[T]) validate(value T) error {
    f.mu.RLock()
    defer f.mu.RUnlock()

    if f.validator == nil {
        return nil
    }
    return f.validator(value)
}
```

### 2.3 具体类型的 Set 方法修改

以 `IntFlag` 为例：

```go
// IntFlag 整数标志
type IntFlag struct {
    *BaseFlag[int]
}

// Set 设置整数标志的值
//
// 参数:
//   - value: 要设置的整数字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 strconv.ParseInt 解析字符串
//   - 使用平台相关的位数(IntSize)
//   - 解析成功后调用验证器（如果设置了）
//   - 如果值超出平台int范围，返回解析错误
func (f *IntFlag) Set(value string) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    // 1. 解析
    n, err := strconv.ParseInt(value, 10, IntSize)
    if err != nil {
        return types.WrapParseError(err, "int", value)
    }

    // 2. 验证（新增：调用 BaseFlag 的 validate 方法）
    if err := f.validate(int(n)); err != nil {
        return err
    }

    // 3. 设置值
    *f.value = int(n)
    f.isSet = true

    return nil
}
```

---

## 三、预置验证器库

### 3.1 数值验证器

```go
// RangeValidator 数值范围验证器
//
// 参数:
//   - min: 最小值（包含）
//   - max: 最大值（包含）
//
// 返回值:
//   - Validator[T]: 范围验证器
//
// 功能说明:
//   - 验证值是否在指定范围内
//   - 支持所有可比较的类型
//
// 示例:
//   - port.SetValidator(qflag.RangeValidator(1, 65535))
func RangeValidator[T constraints.Ordered](min, max T) Validator[T] {
    return func(value T) error {
        if value < min || value > max {
            return fmt.Errorf("值 %v 超出范围 [%v, %v]", value, min, max)
        }
        return nil
    }
}

// PositiveValidator 正数验证器
//
// 返回值:
//   - Validator[T]: 正数验证器
//
// 功能说明:
//   - 验证值是否为正数（大于 0）
//
// 示例:
//   - count.SetValidator(qflag.PositiveValidator[int]())
func PositiveValidator[T constraints.Signed]() Validator[T] {
    return func(value T) error {
        if value <= 0 {
            return fmt.Errorf("值必须为正数")
        }
        return nil
    }
}

// NonNegativeValidator 非负数验证器
//
// 返回值:
//   - Validator[T]: 非负数验证器
//
// 功能说明:
//   - 验证值是否为非负数（大于等于 0）
//
// 示例:
//   - size.SetValidator(qflag.NonNegativeValidator[int]())
func NonNegativeValidator[T constraints.Signed]() Validator[T] {
    return func(value T) error {
        if value < 0 {
            return fmt.Errorf("值不能为负数")
        }
        return nil
    }
}
```

### 3.2 字符串验证器

```go
// LengthValidator 字符串长度验证器
//
// 参数:
//   - min: 最小长度（包含）
//   - max: 最大长度（包含）
//
// 返回值:
//   - Validator[string]: 长度验证器
//
// 功能说明:
//   - 验证字符串长度是否在指定范围内
//
// 示例:
//   - username.SetValidator(qflag.LengthValidator(3, 20))
func LengthValidator(min, max int) Validator[string] {
    return func(value string) error {
        if len(value) < min || len(value) > max {
            return fmt.Errorf("字符串长度 %d 超出范围 [%d, %d]", len(value), min, max)
        }
        return nil
    }
}

// RegexValidator 正则表达式验证器
//
// 参数:
//   - pattern: 正则表达式模式
//
// 返回值:
//   - Validator[string]: 正则验证器
//
// 功能说明:
//   - 验证字符串是否匹配指定的正则表达式
//
// 示例:
//   - email.SetValidator(qflag.RegexValidator(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))
func RegexValidator(pattern string) Validator[string] {
    re := regexp.MustCompile(pattern)
    return func(value string) error {
        if !re.MatchString(value) {
            return fmt.Errorf("值 %q 不匹配正则表达式 %q", value, pattern)
        }
        return nil
    }
}

// NonEmptyValidator 非空字符串验证器
//
// 返回值:
//   - Validator[string]: 非空验证器
//
// 功能说明:
//   - 验证字符串是否为空
//
// 示例:
//   - name.SetValidator(qflag.NonEmptyValidator[string]())
func NonEmptyValidator[T comparable]() Validator[T] {
    var zero T
    return func(value T) error {
        if value == zero {
            return fmt.Errorf("值不能为空")
        }
        return nil
    }
}
```

### 3.3 枚举验证器

```go
// EnumValidator 枚举验证器
//
// 参数:
//   - allowed: 允许的值列表
//
// 返回值:
//   - Validator[T]: 枚举验证器
//
// 功能说明:
//   - 验证值是否在允许的列表中
//
// 示例:
//   - mode.SetValidator(qflag.EnumValidator("dev", "test", "prod"))
func EnumValidator[T comparable](allowed ...T) Validator[T] {
    allowedSet := make(map[T]bool)
    for _, v := range allowed {
        allowedSet[v] = true
    }
    return func(value T) error {
        if !allowedSet[value] {
            return fmt.Errorf("值 %v 不在允许的列表中", value)
        }
        return nil
    }
}
```

### 3.4 文件验证器

```go
// FileExistsValidator 文件存在验证器
//
// 参数:
//   - requireExists: 是否要求文件必须存在
//
// 返回值:
//   - Validator[string]: 文件存在验证器
//
// 功能说明:
//   - 验证文件是否存在
//   - 如果 requireExists 为 true，文件不存在时返回错误
//
// 示例:
//   - configFile.SetValidator(qflag.FileExistsValidator(true))
func FileExistsValidator(requireExists bool) Validator[string] {
    return func(path string) error {
        if requireExists {
            if _, err := os.Stat(path); os.IsNotExist(err) {
                return fmt.Errorf("文件 %q 不存在", path)
            }
        }
        return nil
    }
}

// DirExistsValidator 目录存在验证器
//
// 参数:
//   - requireExists: 是否要求目录必须存在
//
// 返回值:
//   - Validator[string]: 目录存在验证器
//
// 功能说明:
//   - 验证目录是否存在
//   - 如果 requireExists 为 true，目录不存在时返回错误
//
// 示例:
//   - outputDir.SetValidator(qflag.DirExistsValidator(true))
func DirExistsValidator(requireExists bool) Validator[string] {
    return func(path string) error {
        if requireExists {
            info, err := os.Stat(path)
            if os.IsNotExist(err) {
                return fmt.Errorf("目录 %q 不存在", path)
            }
            if !info.IsDir() {
                return fmt.Errorf("路径 %q 不是目录", path)
            }
        }
        return nil
    }
}
```

### 3.5 自定义验证器

```go
// CustomValidator 自定义验证器
//
// 参数:
//   - validateFunc: 自定义验证函数
//
// 返回值:
//   - Validator[T]: 自定义验证器
//
// 功能说明:
//   - 允许用户传入自定义的验证逻辑
//
// 示例:
//   - port.SetValidator(qflag.CustomValidator(func(value int) error {
//         if isPortInUse(value) {
//             return fmt.Errorf("端口 %d 已被占用", value)
//         }
//         return nil
//     }))
func CustomValidator[T any](validateFunc func(value T) error) Validator[T] {
    return validateFunc
}
```

---

## 四、使用示例

### 4.1 基本使用

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 端口号：1-65535
    port := qflag.Root.Int("port", "p", "端口号", 8080)
    port.SetValidator(qflag.RangeValidator(1, 65535))

    // 用户名：3-20个字符
    username := qflag.Root.String("username", "u", "用户名", "")
    username.SetValidator(qflag.LengthValidator(3, 20))

    // 邮箱格式
    email := qflag.Root.String("email", "e", "邮箱", "")
    email.SetValidator(qflag.RegexValidator(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))

    // 解析
    if err := qflag.Parse(); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("端口: %d\n", port.Get())
    fmt.Printf("用户名: %s\n", username.Get())
    fmt.Printf("邮箱: %s\n", email.Get())
}
```

### 4.2 自定义验证器

```go
// 端口号：1-65535 且未被占用
port := qflag.Root.Int("port", "p", "端口号", 8080)
port.SetValidator(func(value int) error {
    // 范围验证
    if value < 1 || value > 65535 {
        return fmt.Errorf("端口 %d 超出范围 [1, 65535]", value)
    }
    // 端口占用检查
    if isPortInUse(value) {
        return fmt.Errorf("端口 %d 已被占用", value)
    }
    return nil
})
```

### 4.3 枚举验证

```go
// 模式：dev, test, prod
mode := qflag.Root.String("mode", "m", "运行模式", "dev")
mode.SetValidator(qflag.EnumValidator("dev", "test", "prod"))
```

### 4.4 文件验证

```go
// 配置文件：必须存在
configFile := qflag.Root.String("config", "c", "配置文件", "")
configFile.SetValidator(qflag.FileExistsValidator(true))

// 输出目录：必须存在
outputDir := qflag.Root.String("output", "o", "输出目录", "")
outputDir.SetValidator(qflag.DirExistsValidator(true))
```

### 4.5 覆盖验证器

```go
// 先设置一个验证器
port.SetValidator(qflag.RangeValidator(1, 1024))

// 覆盖为新的验证器
port.SetValidator(qflag.RangeValidator(1, 65535))

// 清除验证器
port.ClearValidator()
```

---

## 五、实现步骤

### 第一步：修改 BaseFlag[T]

文件：`internal/flag/base_flag.go`

1. 添加 `validator` 字段
2. 实现 `SetValidator` 方法
3. 实现 `ClearValidator` 方法
4. 实现 `HasValidator` 方法
5. 实现 `validate` 私有方法

### 第二步：修改各类型的 Set 方法

需要修改的文件：
- `internal/flag/basic_flags.go` (StringFlag, BoolFlag)
- `internal/flag/numeric_flags.go` (IntFlag, Int64Flag, UintFlag, Uint8Flag, Uint16Flag, Uint32Flag, Uint64Flag, Float64Flag)
- `internal/flag/special_flags.go` (EnumFlag)
- `internal/flag/time_size_flags.go` (DurationFlag, TimeFlag, SizeFlag)
- `internal/flag/collection_flags.go` (StringSliceFlag, IntSliceFlag, Int64SliceFlag, MapFlag)

在每个类型的 `Set` 方法中，解析完值后调用 `f.validate(value)`

### 第三步：实现预置验证器

文件：`internal/flag/validators.go`（新建）

实现所有预置验证器函数

### 第四步：导出验证器类型和函数

文件：`exports.go`

导出 `Validator` 类型和预置验证器函数

### 第五步：编写测试

文件：`internal/flag/validators_test.go`（新建）

为所有验证器编写单元测试

### 第六步：编写示例

文件：`examples/validators/main.go`（新建）

提供验证器使用示例

---

## 六、错误处理

### 6.1 验证器返回的错误

验证器返回的错误应该：
1. **描述清晰**：告诉用户哪个值验证失败，以及失败的原因
2. **包含上下文**：包含标志名称和值
3. **易于理解**：使用用户友好的语言

### 6.2 错误示例

```go
// 范围验证失败
错误: 值 70000 超出范围 [1, 65535]

// 长度验证失败
错误: 字符串长度 2 超出范围 [3, 20]

// 正则验证失败
错误: 值 "invalid-email" 不匹配正则表达式 "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$"

// 枚举验证失败
错误: 值 staging 不在允许的列表中

// 文件验证失败
错误: 文件 "/path/to/config.yaml" 不存在
```

---

## 七、设计优势

| 优势 | 说明 |
|------|------|
| ✅ **不修改 Flag 接口** | 保持向后兼容，现有代码无需改动 |
| ✅ **类型安全** | `BaseFlag[T]` 是泛型的，验证器也是泛型的 |
| ✅ **实现简单** | 只需修改 `BaseFlag[T]` 和各类型的 `Set` 方法 |
| ✅ **用户友好** | 直接调用 `SetValidator`，无需类型断言 |
| ✅ **可选功能** | 不需要验证器的标志可以不使用 |
| ✅ **可覆盖** | 重复设置验证器会覆盖之前的验证器 |
| ✅ **可清除** | 可以随时清除验证器 |
| ✅ **可组合** | 通过自定义验证器可以实现复杂的验证逻辑 |

---

## 八、注意事项

1. **验证器执行时机**：在 `Set` 方法中解析完值后立即执行
2. **验证器线程安全**：验证器执行时已经持有锁，验证器本身不需要处理并发
3. **验证器性能**：验证器应该快速执行，避免耗时操作
4. **验证器错误**：验证器返回的错误会被直接传递给调用者
5. **验证器覆盖**：重复调用 `SetValidator` 会覆盖之前的验证器
6. **验证器清除**：调用 `ClearValidator` 后，后续的 `Set` 调用将不会进行验证

---

## 九、与必需组的配合

验证器和必需组可以配合使用：

```go
// 场景：数据库连接配置
cmd := qflag.NewCmd("db-tool", "", qflag.ExitOnError)

// 标志定义
host := cmd.String("host", "h", "数据库主机", "")
port := cmd.Int("port", "p", "数据库端口", 3306)
user := cmd.String("user", "u", "数据库用户", "")
pass := cmd.String("pass", "", "数据库密码", "")

// 添加验证器
port.SetValidator(qflag.RangeValidator(1, 65535))
host.SetValidator(qflag.NonEmptyValidator[string]())
user.SetValidator(qflag.LengthValidator(1, 50))

// 添加必需组：所有数据库参数都必须提供
cmd.AddRequiredGroup("数据库连接", []string{"host", "port", "user", "pass"})

// 解析
if err := cmd.Parse(os.Args[1:]); err != nil {
    fmt.Printf("错误: %v\n", err)
    os.Exit(1)
}
```

### 验证顺序

1. **必需组验证**：检查所有必需的标志是否被设置
2. **标志解析**：解析命令行参数
3. **验证器验证**：对每个设置了验证器的标志进行验证

---

## 十、总结

本设计方案提供了一个灵活、类型安全、易于使用的验证器机制：

1. **最小侵入性**：只修改 `BaseFlag[T]` 和各类型的 `Set` 方法
2. **完全向后兼容**：不修改 `Flag` 接口
3. **类型安全**：泛型保证类型安全
4. **易于使用**：直接调用 `SetValidator`，无需类型断言
5. **可扩展**：用户可以轻松创建自定义验证器
6. **可覆盖**：重复设置验证器会覆盖之前的验证器
7. **可清除**：可以随时清除验证器

这个方案完美地满足了需求，同时保持了代码的简洁性和可维护性。
