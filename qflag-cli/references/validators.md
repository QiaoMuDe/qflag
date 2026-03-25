# QFlag 验证器列表

## 数值验证器

### 整数范围

```go
import "gitee.com/MM-Q/qflag/validators"

// Int 范围
port := cmd.Int("port", "p", "端口号", 8080)
port.SetValidator(validators.IntRange(1, 65535))

// Uint 范围
percentage := cmd.Uint("percentage", "p", "百分比", 0)
percentage.SetValidator(validators.UintRange(0, 100))

// Uint16 范围（常用端口范围）
port16 := cmd.Uint16("port", "p", "端口号", 8080)
port16.SetValidator(validators.Uint16Range(1, 65535))

// 其他整数类型
validators.Uint8Range(min, max)
validators.Uint32Range(min, max)
validators.Uint64Range(min, max)
```

### 浮点数范围

```go
// Float64 范围
rate := cmd.Float64("rate", "r", "比率", 0.5)
rate.SetValidator(validators.Float64Range(0.0, 1.0))
```

## 字符串验证器

### 正则匹配

```go
// 邮箱验证
email := cmd.String("email", "e", "邮箱地址", "")
email.SetValidator(validators.MatchRegex(`^[\w.-]+@[\w.-]+\.\w+$`))

// URL 验证
url := cmd.String("url", "u", "URL", "")
url.SetValidator(validators.MatchRegex(`^https?://`))

// IP 地址验证
ip := cmd.String("ip", "i", "IP地址", "")
ip.SetValidator(validators.MatchRegex(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`))
```

### 长度验证

```go
// 最小长度
name := cmd.String("name", "n", "名称", "")
name.SetValidator(validators.MinLength(3))

// 最大长度
desc := cmd.String("desc", "d", "描述", "")
desc.SetValidator(validators.MaxLength(100))

// 长度范围
username := cmd.String("username", "u", "用户名", "")
username.SetValidator(validators.LengthRange(3, 20))
```

### 格式验证

```go
// 邮箱格式
email := cmd.String("email", "e", "邮箱", "")
email.SetValidator(validators.IsEmail())

// URL 格式
url := cmd.String("url", "u", "URL", "")
url.SetValidator(validators.IsURL())

// IP 格式
ip := cmd.String("ip", "i", "IP地址", "")
ip.SetValidator(validators.IsIP())

// MAC 地址格式
mac := cmd.String("mac", "m", "MAC地址", "")
mac.SetValidator(validators.IsMAC())
```

## 文件系统验证器

```go
// 文件必须存在
config := cmd.String("config", "c", "配置文件", "")
config.SetValidator(validators.FileExists())

// 必须是目录
dir := cmd.String("dir", "d", "目录路径", "")
dir.SetValidator(validators.IsDir())

// 路径可写
output := cmd.String("output", "o", "输出路径", "")
output.SetValidator(validators.IsWritable())

// 路径可读
input := cmd.String("input", "i", "输入路径", "")
input.SetValidator(validators.IsReadable())
```

## 网络验证器

```go
// 验证主机名或 IP
host := cmd.String("host", "h", "主机地址", "")
host.SetValidator(validators.IsHost())

// 验证端口号
port := cmd.Int("port", "p", "端口号", 8080)
port.SetValidator(validators.IsPort())
```

## 自定义验证器

### 基础自定义验证器

```go
// String 验证器
name := cmd.String("name", "n", "名称", "")
name.SetValidator(func(value string) error {
    if value == "" {
        return fmt.Errorf("名称不能为空")
    }
    if len(value) < 3 {
        return fmt.Errorf("名称长度至少为3个字符")
    }
    return nil
})

// Int 验证器
age := cmd.Int("age", "a", "年龄", 0)
age.SetValidator(func(value int) error {
    if value < 0 || value > 150 {
        return fmt.Errorf("年龄必须在 0-150 之间")
    }
    return nil
})

// Bool 验证器（较少使用）
confirm := cmd.Bool("confirm", "c", "确认", false)
confirm.SetValidator(func(value bool) error {
    if !value {
        return fmt.Errorf("必须确认才能继续")
    }
    return nil
})
```

### 组合验证器

```go
// 多个验证条件
password := cmd.String("password", "p", "密码", "")
password.SetValidator(func(value string) error {
    if len(value) < 8 {
        return fmt.Errorf("密码长度至少为8位")
    }
    hasUpper := false
    hasLower := false
    hasDigit := false
    for _, c := range value {
        switch {
        case c >= 'A' && c <= 'Z':
            hasUpper = true
        case c >= 'a' && c <= 'z':
            hasLower = true
        case c >= '0' && c <= '9':
            hasDigit = true
        }
    }
    if !hasUpper || !hasLower || !hasDigit {
        return fmt.Errorf("密码必须包含大小写字母和数字")
    }
    return nil
})
```

## 验证器链

```go
// 使用多个验证器（需要自定义实现）
func Chain(validators ...func(string) error) func(string) error {
    return func(value string) error {
        for _, v := range validators {
            if err := v(value); err != nil {
                return err
            }
        }
        return nil
    }
}

// 使用
email := cmd.String("email", "e", "邮箱", "")
email.SetValidator(Chain(
    validators.MinLength(5),
    validators.IsEmail(),
))
```
