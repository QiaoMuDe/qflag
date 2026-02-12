# Validators 包

Validators 包提供了丰富的预置验证器函数，用于验证命令行标志的值。

## 目录

- [数值验证器](#数值验证器)
- [字符串验证器](#字符串验证器)
- [网络验证器](#网络验证器)
- [时间验证器](#时间验证器)
- [集合验证器](#集合验证器)
- [映射验证器](#映射验证器)
- [组合验证器](#组合验证器)
- [文件验证器](#文件验证器)
- [其他验证器](#其他验证器)

## 数值验证器

### IntRange / UintRange / Int64Range / Uint64Range / Uint8Range / Uint16Range / Uint32Range

验证数值是否在指定范围内。

```go
import "gitee.com/MM-Q/qflag/validators"

// 端口号验证：1-65535
port.SetValidator(validators.Uint16Range(1, 65535))

// 年龄验证：0-150
age.SetValidator(validators.IntRange(0, 150))

// 百分比验证：0-100
percentage.SetValidator(validators.UintRange(0, 100))
```

### Float64Range

验证浮点数是否在指定范围内。

```go
// 温度验证：-50.0 到 100.0
temperature.SetValidator(validators.Float64Range(-50.0, 100.0))
```

### Positive

验证数值是否大于 0。

```go
// 数量必须大于 0
count.SetValidator(validators.Positive[int]())

// 价格必须大于 0
price.SetValidator(validators.Positive[float64]())
```

### NonNegative

验证数值是否大于等于 0。

```go
// 超时时间不能为负数
timeout.SetValidator(validators.NonNegative[int]())

// 延迟不能为负数
delay.SetValidator(validators.NonNegative[float64]())
```

## 字符串验证器

### StringLength

验证字符串长度是否在指定范围内。

```go
// 用户名长度：3-20 个字符
username.SetValidator(validators.StringLength(3, 20))
```

### StringMinLength / StringMaxLength

验证字符串最小或最大长度。

```go
// 密码至少 8 个字符
password.SetValidator(validators.StringMinLength(8))

// 描述最多 500 个字符
description.SetValidator(validators.StringMaxLength(500))
```

### StringRegex

使用正则表达式验证字符串格式。

```go
// 邮箱格式验证
email.SetValidator(validators.StringRegex(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))

// 手机号验证（中国）
phone.SetValidator(validators.StringRegex(`^1[3-9]\d{9}$`))
```

### StringPrefix / StringSuffix

验证字符串是否以指定前缀或后缀开头/结尾。

```go
// 环境变量名以 "APP_" 开头
envVar.SetValidator(validators.StringPrefix("APP_"))

// 文件名以 ".txt" 结尾
filename.SetValidator(validators.StringSuffix(".txt"))
```

### StringContains

验证字符串是否包含指定子串。

```go
// 路径包含 "/data/"
path.SetValidator(validators.StringContains("/data/"))
```

### StringNotEmpty

验证字符串非空。

```go
// 用户名不能为空
username.SetValidator(validators.StringNotEmpty())
```

### StringOneOf

验证字符串是否在允许的值列表中。

```go
// 日志级别
logLevel.SetValidator(validators.StringOneOf("debug", "info", "warn", "error"))

// 颜色选项
color.SetValidator(validators.StringOneOf("red", "green", "blue"))
```

### StringCharset

验证字符串是否只包含指定字符集的字符。

```go
// 只允许字母和数字
username.SetValidator(validators.StringCharset("alnum"))

// 只允许小写字母
code.SetValidator(validators.StringCharset("lower"))

// 只允许十六进制字符
hash.SetValidator(validators.StringCharset("hex"))
```

支持的字符集：
- `alnum`: 字母和数字 (a-zA-Z0-9)
- `alpha`: 字母 (a-zA-Z)
- `digit`: 数字 (0-9)
- `hex`: 十六进制字符 (0-9a-fA-F)
- `lower`: 小写字母 (a-z)
- `upper`: 大写字母 (A-Z)

## 网络验证器

### Email

验证邮箱格式。

```go
email.SetValidator(validators.Email())
```

### URL

验证 URL 格式，可以指定必须的协议。

```go
// 任意 URL
apiURL.SetValidator(validators.URL(""))

// 只允许 HTTPS
secureURL.SetValidator(validators.URL("https"))
```

### IPv4 / IPv6 / IP

验证 IP 地址格式。

```go
// IPv4 地址
ipv4.SetValidator(validators.IPv4())

// IPv6 地址
ipv6.SetValidator(validators.IPv6())

// IPv4 或 IPv6 地址
ip.SetValidator(validators.IP())
```

### Port

验证端口号范围 (1-65535)。

```go
port.SetValidator(validators.Port())
```

### Hostname

验证主机名格式。

```go
host.SetValidator(validators.Hostname())
```

## 时间验证器

### DurationMin / DurationMax / DurationRange

验证持续时间范围。

```go
// 超时至少 1 秒
timeout.SetValidator(validators.DurationMin(time.Second))

// 延迟最多 1 小时
delay.SetValidator(validators.DurationMax(time.Hour))

// 重试间隔：1-60 秒
retryInterval.SetValidator(validators.DurationRange(time.Second, time.Minute))
```

### TimeAfter / TimeBefore / TimeRange

验证时间范围。

```go
// 截止时间必须在未来
deadline.SetValidator(validators.TimeAfter(time.Now()))

// 开始时间必须在过去
startTime.SetValidator(validators.TimeBefore(time.Now()))

// 预约时间在工作时间内
appointmentTime.SetValidator(validators.TimeRange(startTime, endTime))
```

## 集合验证器

### SliceLength / SliceMinLength / SliceMaxLength

验证切片长度范围。

```go
// 至少选择 1 个，最多 5 个
tags.SetValidator(validators.SliceLength[string](1, 5))

// 至少选择 1 个
options.SetValidator(validators.SliceMinLength[string](1))

// 最多选择 10 个
items.SetValidator(validators.SliceMaxLength[string](10))
```

### SliceNotEmpty

验证切片非空。

```go
selectedItems.SetValidator(validators.SliceNotEmpty[string]())
```

### SliceUnique

验证切片元素唯一性。

```go
tags.SetValidator(validators.SliceUnique[string]())
```

### SliceContains

验证切片是否包含指定元素。

```go
options.SetValidator(validators.SliceContains[string]("default"))
```

## 映射验证器

### MapKeys

验证映射的所有键是否都在允许的键列表中。

```go
config.SetValidator(validators.MapKeys[string]("host", "port", "timeout"))
```

### MapMinSize / MapMaxSize

验证映射大小范围。

```go
// 至少提供 1 个配置项
config.SetValidator(validators.MapMinSize[string, int](1))

// 最多提供 10 个配置项
config.SetValidator(validators.MapMaxSize[string, int](10))
```

### MapRequiredKeys

验证映射是否包含所有必需的键。

```go
config.SetValidator(validators.MapRequiredKeys[string]("host", "port"))
```

## 组合验证器

### And

所有验证器必须全部通过。

```go
username.SetValidator(validators.And(
    validators.StringMinLength(3),
    validators.StringMaxLength(20),
    validators.StringCharset("alnum"),
))
```

### Or

至少有一个验证器通过。

```go
contact.SetValidator(validators.Or(
    validators.Email(),
    validators.StringRegex(`^\d{11}$`), // 手机号
))
```

### Not

验证器必须失败。

```go
username.SetValidator(validators.Not(validators.StringOneOf("admin", "root")))
```

### Optional

如果值非空则验证，空值跳过验证。

```go
email.SetValidator(validators.Optional(validators.Email()))
```

## 文件验证器

### FileExists

验证文件是否存在。

```go
configFile.SetValidator(validators.FileExists())
```

### DirExists

验证目录是否存在。

```go
outputDir.SetValidator(validators.DirExists())
```

### FileExtension

验证文件扩展名。

```go
filename.SetValidator(validators.FileExtension("json", "yaml", "yml"))
```

## 其他验证器

### IsNumeric

验证字符串是否为有效的数字。

```go
number.SetValidator(validators.IsNumeric())
```

### IsInteger

验证字符串是否为有效的整数。

```go
intNumber.SetValidator(validators.IsInteger())
```

### IsPositiveInteger

验证字符串是否为有效的正整数。

```go
count.SetValidator(validators.IsPositiveInteger())
```

## 完整示例

```go
package main

import (
    "fmt"
    "time"

    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/validators"
)

func main() {
    // 创建端口号标志，验证范围 1-65535
    port := qflag.Root.Int("port", "p", "端口号", 8080)
    port.SetValidator(validators.Port())

    // 创建用户名标志，验证长度和字符集
    username := qflag.Root.String("username", "u", "用户名", "")
    username.SetValidator(validators.And(
        validators.StringLength(3, 20),
        validators.StringCharset("alnum"),
    ))

    // 创建邮箱标志，可选但提供时必须格式正确
    email := qflag.Root.String("email", "e", "邮箱地址", "")
    email.SetValidator(validators.Optional(validators.Email()))

    // 创建超时标志，验证范围 1-3600 秒
    timeout := qflag.Root.Int("timeout", "t", "超时时间（秒）", 30)
    timeout.SetValidator(validators.IntRange(1, 3600))

    // 创建标签标志，验证长度和唯一性
    tags := qflag.Root.StringSlice("tags", "", "标签", []string{})
    tags.SetValidator(validators.And(
        validators.SliceMaxLength[string](5),
        validators.SliceUnique[string](),
    ))

    // 解析命令行参数
    if err := qflag.Root.Parse(); err != nil {
        fmt.Printf("参数解析失败: %v\n", err)
        return
    }

    fmt.Printf("端口: %d\n", port.Get())
    fmt.Printf("用户名: %s\n", username.Get())
    fmt.Printf("邮箱: %s\n", email.Get())
    fmt.Printf("超时: %d 秒\n", timeout.Get())
    fmt.Printf("标签: %v\n", tags.Get())
}
```

## 注意事项

1. **空值处理**：
   - `StringFlag`: 空字符串不经过验证器，直接设置
   - `BoolFlag`: 不经过验证器（无空值概念）
   - 集合类型 (MapFlag, StringSliceFlag, IntSliceFlag, Int64SliceFlag): 空字符串不经过验证器，创建空集合
   - 其他类型: 空字符串直接返回错误，不经过验证器

2. **验证器性能**：
   - 验证器应该快速执行，避免耗时操作
   - 避免在验证器中进行网络请求或文件 I/O

3. **错误消息**：
   - 验证器返回的错误应该清晰描述失败原因
   - 使用组合验证器时，错误消息会按顺序返回

4. **线程安全**：
   - 验证器执行时已经持有锁，验证器本身不需要处理并发
