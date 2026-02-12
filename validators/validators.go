package validators

import (
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// IntRange 创建整数范围验证器
//
// 参数:
//   - min: 最小值（包含）
//   - max: 最大值（包含）
//
// 返回值:
//   - func(int) error: 整数验证器函数
//
// 功能说明:
//   - 验证整数是否在 [min, max] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	port.SetValidator(validators.IntRange(1, 65535))
func IntRange(min, max int) func(int) error {
	return func(value int) error {
		if value < min || value > max {
			return fmt.Errorf("整数 %d 超出范围 [%d, %d]", value, min, max)
		}
		return nil
	}
}

// UintRange 创建无符号整数范围验证器
//
// 参数:
//   - min: 最小值（包含）
//   - max: 最大值（包含）
//
// 返回值:
//   - func(uint) error: 无符号整数验证器函数
//
// 功能说明:
//   - 验证无符号整数是否在 [min, max] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	percentage.SetValidator(validators.UintRange(0, 100))
func UintRange(min, max uint) func(uint) error {
	return func(value uint) error {
		if value < min || value > max {
			return fmt.Errorf("无符号整数 %d 超出范围 [%d, %d]", value, min, max)
		}
		return nil
	}
}

// Uint8Range 创建8位无符号整数范围验证器
//
// 参数:
//   - min: 最小值（包含）
//   - max: 最大值（包含）
//
// 返回值:
//   - func(uint8) error: 8位无符号整数验证器函数
//
// 功能说明:
//   - 验证8位无符号整数是否在 [min, max] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	byteValue.SetValidator(validators.Uint8Range(0, 255))
func Uint8Range(min, max uint8) func(uint8) error {
	return func(value uint8) error {
		if value < min || value > max {
			return fmt.Errorf("8位无符号整数 %d 超出范围 [%d, %d]", value, min, max)
		}
		return nil
	}
}

// Uint16Range 创建16位无符号整数范围验证器
//
// 参数:
//   - min: 最小值（包含）
//   - max: 最大值（包含）
//
// 返回值:
//   - func(uint16) error: 16位无符号整数验证器函数
//
// 功能说明:
//   - 验证16位无符号整数是否在 [min, max] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	port.SetValidator(validators.Uint16Range(1, 65535))
func Uint16Range(min, max uint16) func(uint16) error {
	return func(value uint16) error {
		if value < min || value > max {
			return fmt.Errorf("16位无符号整数 %d 超出范围 [%d, %d]", value, min, max)
		}
		return nil
	}
}

// Uint32Range 创建32位无符号整数范围验证器
//
// 参数:
//   - min: 最小值（包含）
//   - max: 最大值（包含）
//
// 返回值:
//   - func(uint32) error: 32位无符号整数验证器函数
//
// 功能说明:
//   - 验证32位无符号整数是否在 [min, max] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	ipAddr.SetValidator(validators.Uint32Range(0, 4294967295))
func Uint32Range(min, max uint32) func(uint32) error {
	return func(value uint32) error {
		if value < min || value > max {
			return fmt.Errorf("32位无符号整数 %d 超出范围 [%d, %d]", value, min, max)
		}
		return nil
	}
}

// Uint64Range 创建64位无符号整数范围验证器
//
// 参数:
//   - min: 最小值（包含）
//   - max: 最大值（包含）
//
// 返回值:
//   - func(uint64) error: 64位无符号整数验证器函数
//
// 功能说明:
//   - 验证64位无符号整数是否在 [min, max] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	largeNumber.SetValidator(validators.Uint64Range(0, 18446744073709551615))
func Uint64Range(min, max uint64) func(uint64) error {
	return func(value uint64) error {
		if value < min || value > max {
			return fmt.Errorf("64位无符号整数 %d 超出范围 [%d, %d]", value, min, max)
		}
		return nil
	}
}

// Int64Range 创建64位整数范围验证器
//
// 参数:
//   - min: 最小值（包含）
//   - max: 最大值（包含）
//
// 返回值:
//   - func(int64) error: 64位整数验证器函数
//
// 功能说明:
//   - 验证64位整数是否在 [min, max] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	timestamp.SetValidator(validators.Int64Range(0, 9999999999))
func Int64Range(min, max int64) func(int64) error {
	return func(value int64) error {
		if value < min || value > max {
			return fmt.Errorf("64位整数 %d 超出范围 [%d, %d]", value, min, max)
		}
		return nil
	}
}

// Float64Range 创建64位浮点数范围验证器
//
// 参数:
//   - min: 最小值（包含）
//   - max: 最大值（包含）
//
// 返回值:
//   - func(float64) error: 64位浮点数验证器函数
//
// 功能说明:
//   - 验证64位浮点数是否在 [min, max] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	temperature.SetValidator(validators.Float64Range(-50.0, 100.0))
func Float64Range(min, max float64) func(float64) error {
	return func(value float64) error {
		if value < min || value > max {
			return fmt.Errorf("浮点数 %.2f 超出范围 [%.2f, %.2f]", value, min, max)
		}
		return nil
	}
}

// Positive 创建正数验证器
//
// 返回值:
//   - func(int) error: 整数正数验证器函数
//
// 功能说明:
//   - 验证整数是否大于 0
//   - 小于等于 0 返回错误
//
// 使用示例:
//
//	count.SetValidator(validators.Positive[int]())
func Positive[T ~int | ~int64 | ~uint | ~uint64 | ~float64]() func(T) error {
	return func(value T) error {
		switch v := any(value).(type) {
		case int:
			if v <= 0 {
				return fmt.Errorf("整数 %d 必须大于 0", v)
			}
		case int64:
			if v <= 0 {
				return fmt.Errorf("64位整数 %d 必须大于 0", v)
			}
		case uint:
			if v == 0 {
				return fmt.Errorf("无符号整数 %d 必须大于 0", v)
			}
		case uint64:
			if v == 0 {
				return fmt.Errorf("64位无符号整数 %d 必须大于 0", v)
			}
		case float64:
			if v <= 0 {
				return fmt.Errorf("浮点数 %.2f 必须大于 0", v)
			}
		}
		return nil
	}
}

// NonNegative 创建非负数验证器
//
// 返回值:
//   - func(int) error: 整数非负数验证器函数
//
// 功能说明:
//   - 验证整数是否大于等于 0
//   - 小于 0 返回错误
//
// 使用示例:
//
//	timeout.SetValidator(validators.NonNegative[int]())
func NonNegative[T ~int | ~int64 | ~float64]() func(T) error {
	return func(value T) error {
		switch v := any(value).(type) {
		case int:
			if v < 0 {
				return fmt.Errorf("整数 %d 必须大于等于 0", v)
			}
		case int64:
			if v < 0 {
				return fmt.Errorf("64位整数 %d 必须大于等于 0", v)
			}
		case float64:
			if v < 0 {
				return fmt.Errorf("浮点数 %.2f 必须大于等于 0", v)
			}
		}
		return nil
	}
}

// StringLength 创建字符串长度范围验证器
//
// 参数:
//   - minLength: 最小长度（包含），0 表示无限制
//   - maxLength: 最大长度（包含），0 表示无限制
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串长度是否在 [minLength, maxLength] 范围内
//   - 超出范围返回错误
//   - minLength 或 maxLength 为 0 表示该方向无限制
//
// 使用示例:
//
//	username.SetValidator(validators.StringLength(3, 20))
func StringLength(minLength, maxLength int) func(string) error {
	return func(value string) error {
		length := len(value)
		if minLength > 0 && length < minLength {
			return fmt.Errorf("字符串长度 %d 小于最小长度 %d", length, minLength)
		}
		if maxLength > 0 && length > maxLength {
			return fmt.Errorf("字符串长度 %d 大于最大长度 %d", length, maxLength)
		}
		return nil
	}
}

// StringMinLength 创建字符串最小长度验证器
//
// 参数:
//   - minLength: 最小长度（包含）
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串长度是否大于等于 minLength
//   - 小于 minLength 返回错误
//
// 使用示例:
//
//	password.SetValidator(validators.StringMinLength(8))
func StringMinLength(minLength int) func(string) error {
	return func(value string) error {
		if len(value) < minLength {
			return fmt.Errorf("字符串长度 %d 小于最小长度 %d", len(value), minLength)
		}
		return nil
	}
}

// StringMaxLength 创建字符串最大长度验证器
//
// 参数:
//   - maxLength: 最大长度（包含）
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串长度是否小于等于 maxLength
//   - 大于 maxLength 返回错误
//
// 使用示例:
//
//	description.SetValidator(validators.StringMaxLength(500))
func StringMaxLength(maxLength int) func(string) error {
	return func(value string) error {
		if len(value) > maxLength {
			return fmt.Errorf("字符串长度 %d 大于最大长度 %d", len(value), maxLength)
		}
		return nil
	}
}

// StringRegex 创建正则表达式验证器
//
// 参数:
//   - pattern: 正则表达式模式
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 使用正则表达式验证字符串格式
//   - 不匹配返回错误
//
// 使用示例:
//
//	email.SetValidator(validators.StringRegex(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))
func StringRegex(pattern string) func(string) error {
	regex := regexp.MustCompile(pattern)
	return func(value string) error {
		if !regex.MatchString(value) {
			return fmt.Errorf("字符串 '%s' 不匹配正则表达式 '%s'", value, pattern)
		}
		return nil
	}
}

// StringPrefix 创建字符串前缀验证器
//
// 参数:
//   - prefix: 必须包含的前缀
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否以指定前缀开头
//   - 不以该前缀开头返回错误
//
// 使用示例:
//
//	envVar.SetValidator(validators.StringPrefix("APP_"))
func StringPrefix(prefix string) func(string) error {
	return func(value string) error {
		if !strings.HasPrefix(value, prefix) {
			return fmt.Errorf("字符串 '%s' 必须以 '%s' 开头", value, prefix)
		}
		return nil
	}
}

// StringSuffix 创建字符串后缀验证器
//
// 参数:
//   - suffix: 必须包含的后缀
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否以指定后缀结尾
//   - 不以该后缀结尾返回错误
//
// 使用示例:
//
//	filename.SetValidator(validators.StringSuffix(".txt"))
func StringSuffix(suffix string) func(string) error {
	return func(value string) error {
		if !strings.HasSuffix(value, suffix) {
			return fmt.Errorf("字符串 '%s' 必须以 '%s' 结尾", value, suffix)
		}
		return nil
	}
}

// StringContains 创建字符串包含验证器
//
// 参数:
//   - substring: 必须包含的子串
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否包含指定子串
//   - 不包含返回错误
//
// 使用示例:
//
//	path.SetValidator(validators.StringContains("/data/"))
func StringContains(substring string) func(string) error {
	return func(value string) error {
		if !strings.Contains(value, substring) {
			return fmt.Errorf("字符串 '%s' 必须包含 '%s'", value, substring)
		}
		return nil
	}
}

// StringNotEmpty 创建非空字符串验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否非空
//   - 空字符串返回错误
//
// 注意事项:
//   - StringFlag 的空字符串不经过验证器，此验证器主要用于其他场景
//
// 使用示例:
//
//	username.SetValidator(validators.StringNotEmpty())
func StringNotEmpty() func(string) error {
	return func(value string) error {
		if value == "" {
			return errors.New("字符串不能为空")
		}
		return nil
	}
}

// StringOneOf 创建字符串枚举验证器
//
// 参数:
//   - allowedValues: 允许的值列表
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否在允许的值列表中
//   - 不在列表中返回错误
//
// 使用示例:
//
//	logLevel.SetValidator(validators.StringOneOf("debug", "info", "warn", "error"))
func StringOneOf(allowedValues ...string) func(string) error {
	allowedSet := make(map[string]bool)
	for _, v := range allowedValues {
		allowedSet[v] = true
	}
	return func(value string) error {
		if !allowedSet[value] {
			return fmt.Errorf("字符串 '%s' 不在允许的值列表中: %v", value, allowedValues)
		}
		return nil
	}
}

// StringCharset 创建字符集验证器
//
// 参数:
//   - charset: 字符集类型，可选值: "alnum", "alpha", "digit", "hex", "lower", "upper"
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否只包含指定字符集的字符
//   - 包含其他字符返回错误
//
// 字符集说明:
//   - "alnum": 字母和数字 (a-zA-Z0-9)
//   - "alpha": 字母 (a-zA-Z)
//   - "digit": 数字 (0-9)
//   - "hex": 十六进制字符 (0-9a-fA-F)
//   - "lower": 小写字母 (a-z)
//   - "upper": 大写字母 (A-Z)
//
// 使用示例:
//
//	username.SetValidator(validators.StringCharset("alnum"))
func StringCharset(charset string) func(string) error {
	var regex *regexp.Regexp
	switch charset {
	case "alnum":
		regex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	case "alpha":
		regex = regexp.MustCompile(`^[a-zA-Z]+$`)
	case "digit":
		regex = regexp.MustCompile(`^[0-9]+$`)
	case "hex":
		regex = regexp.MustCompile(`^[0-9a-fA-F]+$`)
	case "lower":
		regex = regexp.MustCompile(`^[a-z]+$`)
	case "upper":
		regex = regexp.MustCompile(`^[A-Z]+$`)
	default:
		return func(value string) error {
			return fmt.Errorf("不支持的字符集类型: %s", charset)
		}
	}
	return func(value string) error {
		if !regex.MatchString(value) {
			return fmt.Errorf("字符串 '%s' 包含非 %s 字符集的字符", value, charset)
		}
		return nil
	}
}

// Email 创建邮箱格式验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否符合邮箱格式
//   - 不符合返回错误
//
// 使用示例:
//
//	email.SetValidator(validators.Email())
func Email() func(string) error {
	return func(value string) error {
		_, err := mail.ParseAddress(value)
		if err != nil {
			return fmt.Errorf("邮箱格式无效: %s", value)
		}
		return nil
	}
}

// URL 创建 URL 格式验证器
//
// 参数:
//   - scheme: 可选的 URL 协议，如 "http", "https"，空字符串表示不限制协议
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否符合 URL 格式
//   - 可以指定必须的协议
//   - 不符合返回错误
//
// 使用示例:
//
//	apiURL.SetValidator(validators.URL("https"))
func URL(scheme string) func(string) error {
	return func(value string) error {
		u, err := url.Parse(value)
		if err != nil {
			return fmt.Errorf("URL 格式无效: %s", value)
		}
		if u.Scheme == "" {
			return fmt.Errorf("URL 必须包含协议（scheme）: %s", value)
		}
		if u.Host == "" && u.Path == "" {
			return fmt.Errorf("URL 必须包含主机或路径: %s", value)
		}
		if scheme != "" && u.Scheme != scheme {
			return fmt.Errorf("URL 必须使用 %s 协议", scheme)
		}
		return nil
	}
}

// IPv4 创建 IPv4 地址格式验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否符合 IPv4 地址格式
//   - 不符合返回错误
//
// 使用示例:
//
//	ip.SetValidator(validators.IPv4())
func IPv4() func(string) error {
	return func(value string) error {
		ip := net.ParseIP(value)
		if ip == nil || ip.To4() == nil {
			return fmt.Errorf("IPv4 地址格式无效: %s", value)
		}
		return nil
	}
}

// IPv6 创建 IPv6 地址格式验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否符合 IPv6 地址格式
//   - 不符合返回错误
//
// 使用示例:
//
//	ip.SetValidator(validators.IPv6())
func IPv6() func(string) error {
	return func(value string) error {
		ip := net.ParseIP(value)
		if ip == nil || ip.To4() != nil {
			return fmt.Errorf("IPv6 地址格式无效: %s", value)
		}
		return nil
	}
}

// IP 创建 IP 地址格式验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否符合 IPv4 或 IPv6 地址格式
//   - 不符合返回错误
//
// 使用示例:
//
//	ip.SetValidator(validators.IP())
func IP() func(string) error {
	return func(value string) error {
		ip := net.ParseIP(value)
		if ip == nil {
			return fmt.Errorf("IP 地址格式无效: %s", value)
		}
		return nil
	}
}

// Port 创建端口号验证器
//
// 返回值:
//   - func(uint16) error: 端口号验证器函数
//
// 功能说明:
//   - 验证端口号是否在有效范围 [1, 65535] 内
//   - 超出范围返回错误
//
// 使用示例:
//
//	port.SetValidator(validators.Port())
func Port() func(uint16) error {
	return Uint16Range(1, 65535)
}

// Hostname 创建主机名格式验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否符合主机名格式
//   - 不符合返回错误
//
// 使用示例:
//
//	host.SetValidator(validators.Hostname())
func Hostname() func(string) error {
	hostnameRegex := regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`)
	return func(value string) error {
		if !hostnameRegex.MatchString(value) {
			return fmt.Errorf("主机名格式无效: %s", value)
		}
		return nil
	}
}

// DurationMin 创建持续时间最小值验证器
//
// 参数:
//   - minDuration: 最小持续时间
//
// 返回值:
//   - func(time.Duration) error: 持续时间验证器函数
//
// 功能说明:
//   - 验证持续时间是否大于等于 minDuration
//   - 小于 minDuration 返回错误
//
// 使用示例:
//
//	timeout.SetValidator(validators.DurationMin(time.Second))
func DurationMin(minDuration time.Duration) func(time.Duration) error {
	return func(value time.Duration) error {
		if value < minDuration {
			return fmt.Errorf("持续时间 %v 小于最小值 %v", value, minDuration)
		}
		return nil
	}
}

// DurationMax 创建持续时间最大值验证器
//
// 参数:
//   - maxDuration: 最大持续时间
//
// 返回值:
//   - func(time.Duration) error: 持续时间验证器函数
//
// 功能说明:
//   - 验证持续时间是否小于等于 maxDuration
//   - 大于 maxDuration 返回错误
//
// 使用示例:
//
//	delay.SetValidator(validators.DurationMax(time.Hour))
func DurationMax(maxDuration time.Duration) func(time.Duration) error {
	return func(value time.Duration) error {
		if value > maxDuration {
			return fmt.Errorf("持续时间 %v 大于最大值 %v", value, maxDuration)
		}
		return nil
	}
}

// DurationRange 创建持续时间范围验证器
//
// 参数:
//   - minDuration: 最小持续时间
//   - maxDuration: 最大持续时间
//
// 返回值:
//   - func(time.Duration) error: 持续时间验证器函数
//
// 功能说明:
//   - 验证持续时间是否在 [minDuration, maxDuration] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	retryInterval.SetValidator(validators.DurationRange(time.Second, time.Minute))
func DurationRange(minDuration, maxDuration time.Duration) func(time.Duration) error {
	return func(value time.Duration) error {
		if value < minDuration {
			return fmt.Errorf("持续时间 %v 小于最小值 %v", value, minDuration)
		}
		if value > maxDuration {
			return fmt.Errorf("持续时间 %v 大于最大值 %v", value, maxDuration)
		}
		return nil
	}
}

// TimeAfter 创建时间范围验证器（指定时间之后）
//
// 参数:
//   - t: 参考时间
//
// 返回值:
//   - func(time.Time) error: 时间验证器函数
//
// 功能说明:
//   - 验证时间是否在指定时间之后
//   - 不在之后返回错误
//
// 使用示例:
//
//	deadline.SetValidator(validators.TimeAfter(time.Now()))
func TimeAfter(t time.Time) func(time.Time) error {
	return func(value time.Time) error {
		if !value.After(t) {
			return fmt.Errorf("时间 %v 必须在 %v 之后", value, t)
		}
		return nil
	}
}

// TimeBefore 创建时间范围验证器（指定时间之前）
//
// 参数:
//   - t: 参考时间
//
// 返回值:
//   - func(time.Time) error: 时间验证器函数
//
// 功能说明:
//   - 验证时间是否在指定时间之前
//   - 不在之前返回错误
//
// 使用示例:
//
//	startTime.SetValidator(validators.TimeBefore(time.Now()))
func TimeBefore(t time.Time) func(time.Time) error {
	return func(value time.Time) error {
		if !value.Before(t) {
			return fmt.Errorf("时间 %v 必须在 %v 之前", value, t)
		}
		return nil
	}
}

// TimeRange 创建时间范围验证器
//
// 参数:
//   - startTime: 开始时间
//   - endTime: 结束时间
//
// 返回值:
//   - func(time.Time) error: 时间验证器函数
//
// 功能说明:
//   - 验证时间是否在 [startTime, endTime] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	appointmentTime.SetValidator(validators.TimeRange(startTime, endTime))
func TimeRange(startTime, endTime time.Time) func(time.Time) error {
	return func(value time.Time) error {
		if value.Before(startTime) {
			return fmt.Errorf("时间 %v 必须在 %v 之后", value, startTime)
		}
		if value.After(endTime) {
			return fmt.Errorf("时间 %v 必须在 %v 之前", value, endTime)
		}
		return nil
	}
}

// SliceLength 创建切片长度范围验证器
//
// 参数:
//   - minLength: 最小长度（包含），0 表示无限制
//   - maxLength: 最大长度（包含），0 表示无限制
//
// 返回值:
//   - func([]T) error: 切片验证器函数
//
// 功能说明:
//   - 验证切片长度是否在 [minLength, maxLength] 范围内
//   - 超出范围返回错误
//
// 使用示例:
//
//	tags.SetValidator(validators.SliceLength[string](1, 5))
func SliceLength[T any](minLength, maxLength int) func([]T) error {
	return func(value []T) error {
		length := len(value)
		if minLength > 0 && length < minLength {
			return fmt.Errorf("切片长度 %d 小于最小长度 %d", length, minLength)
		}
		if maxLength > 0 && length > maxLength {
			return fmt.Errorf("切片长度 %d 大于最大长度 %d", length, maxLength)
		}
		return nil
	}
}

// SliceMinLength 创建切片最小长度验证器
//
// 参数:
//   - minLength: 最小长度（包含）
//
// 返回值:
//   - func([]T) error: 切片验证器函数
//
// 功能说明:
//   - 验证切片长度是否大于等于 minLength
//   - 小于 minLength 返回错误
//
// 使用示例:
//
//	options.SetValidator(validators.SliceMinLength[string](1))
func SliceMinLength[T any](minLength int) func([]T) error {
	return func(value []T) error {
		if len(value) < minLength {
			return fmt.Errorf("切片长度 %d 小于最小长度 %d", len(value), minLength)
		}
		return nil
	}
}

// SliceMaxLength 创建切片最大长度验证器
//
// 参数:
//   - maxLength: 最大长度（包含）
//
// 返回值:
//   - func([]T) error: 切片验证器函数
//
// 功能说明:
//   - 验证切片长度是否小于等于 maxLength
//   - 大于 maxLength 返回错误
//
// 使用示例:
//
//	items.SetValidator(validators.SliceMaxLength[string](10))
func SliceMaxLength[T any](maxLength int) func([]T) error {
	return func(value []T) error {
		if len(value) > maxLength {
			return fmt.Errorf("切片长度 %d 大于最大长度 %d", len(value), maxLength)
		}
		return nil
	}
}

// SliceNotEmpty 创建非空切片验证器
//
// 返回值:
//   - func([]T) error: 切片验证器函数
//
// 功能说明:
//   - 验证切片是否非空
//   - 空切片返回错误
//
// 使用示例:
//
//	selectedItems.SetValidator(validators.SliceNotEmpty[string]())
func SliceNotEmpty[T any]() func([]T) error {
	return func(value []T) error {
		if len(value) == 0 {
			return errors.New("切片不能为空")
		}
		return nil
	}
}

// SliceUnique 创建唯一性验证器（适用于可比较的类型）
//
// 返回值:
//   - func([]T) error: 切片验证器函数
//
// 功能说明:
//   - 验证切片元素是否唯一
//   - 有重复元素返回错误
//
// 使用示例:
//
//	tags.SetValidator(validators.SliceUnique[string]())
func SliceUnique[T comparable]() func([]T) error {
	return func(value []T) error {
		seen := make(map[T]bool)
		for _, item := range value {
			if seen[item] {
				return fmt.Errorf("切片包含重复元素: %v", item)
			}
			seen[item] = true
		}
		return nil
	}
}

// SliceContains 创建切片包含验证器
//
// 参数:
//   - element: 必须包含的元素
//
// 返回值:
//   - func([]T) error: 切片验证器函数
//
// 功能说明:
//   - 验证切片是否包含指定元素
//   - 不包含返回错误
//
// 使用示例:
//
//	options.SetValidator(validators.SliceContains[string]("default"))
func SliceContains[T comparable](element T) func([]T) error {
	return func(value []T) error {
		for _, item := range value {
			if item == element {
				return nil
			}
		}
		return fmt.Errorf("切片必须包含元素: %v", element)
	}
}

// MapKeys 创建映射键验证器
//
// 参数:
//   - allowedKeys: 允许的键列表
//
// 返回值:
//   - func(map[string]T) error: 映射验证器函数
//
// 功能说明:
//   - 验证映射的所有键是否都在允许的键列表中
//   - 有不允许的键返回错误
//
// 使用示例:
//
//	config.SetValidator(validators.MapKeys[string]("host", "port", "timeout"))
func MapKeys[T any](allowedKeys ...string) func(map[string]T) error {
	allowedSet := make(map[string]bool)
	for _, k := range allowedKeys {
		allowedSet[k] = true
	}
	return func(value map[string]T) error {
		for key := range value {
			if !allowedSet[key] {
				return fmt.Errorf("映射包含不允许的键: %s", key)
			}
		}
		return nil
	}
}

// MapMinSize 创建映射最小大小验证器
//
// 参数:
//   - minSize: 最小大小（包含）
//
// 返回值:
//   - func(map[K]V) error: 映射验证器函数
//
// 功能说明:
//   - 验证映射大小是否大于等于 minSize
//   - 小于 minSize 返回错误
//
// 使用示例:
//
//	config.SetValidator(validators.MapMinSize[string, int](1))
func MapMinSize[K comparable, V any](minSize int) func(map[K]V) error {
	return func(value map[K]V) error {
		if len(value) < minSize {
			return fmt.Errorf("映射大小 %d 小于最小值 %d", len(value), minSize)
		}
		return nil
	}
}

// MapMaxSize 创建映射最大大小验证器
//
// 参数:
//   - maxSize: 最大大小（包含）
//
// 返回值:
//   - func(map[K]V) error: 映射验证器函数
//
// 功能说明:
//   - 验证映射大小是否小于等于 maxSize
//   - 大于 maxSize 返回错误
//
// 使用示例:
//
//	config.SetValidator(validators.MapMaxSize[string, int](10))
func MapMaxSize[K comparable, V any](maxSize int) func(map[K]V) error {
	return func(value map[K]V) error {
		if len(value) > maxSize {
			return fmt.Errorf("映射大小 %d 大于最大值 %d", len(value), maxSize)
		}
		return nil
	}
}

// MapRequiredKeys 创建必需键验证器
//
// 参数:
//   - requiredKeys: 必须包含的键列表
//
// 返回值:
//   - func(map[string]T) error: 映射验证器函数
//
// 功能说明:
//   - 验证映射是否包含所有必需的键
//   - 缺少必需键返回错误
//
// 使用示例:
//
//	config.SetValidator(validators.MapRequiredKeys[string]("host", "port"))
func MapRequiredKeys[T any](requiredKeys ...string) func(map[string]T) error {
	return func(value map[string]T) error {
		for _, key := range requiredKeys {
			if _, ok := value[key]; !ok {
				return fmt.Errorf("映射缺少必需的键: %s", key)
			}
		}
		return nil
	}
}

// And 创建组合验证器（逻辑与）
//
// 参数:
//   - validators: 验证器函数列表
//
// 返回值:
//   - func(T) error: 组合验证器函数
//
// 功能说明:
//   - 所有验证器必须全部通过
//   - 有一个验证器失败则返回该错误
//
// 使用示例:
//
//	username.SetValidator(validators.And(
//	    validators.StringMinLength(3),
//	    validators.StringMaxLength(20),
//	    validators.StringCharset("alnum"),
//	))
func And[T any](validators ...func(T) error) func(T) error {
	return func(value T) error {
		for _, validator := range validators {
			if err := validator(value); err != nil {
				return err
			}
		}
		return nil
	}
}

// Or 创建组合验证器（逻辑或）
//
// 参数:
//   - validators: 验证器函数列表
//
// 返回值:
//   - func(T) error: 组合验证器函数
//
// 功能说明:
//   - 至少有一个验证器通过
//   - 所有验证器都失败则返回最后一个错误
//
// 使用示例:
//
//	contact.SetValidator(validators.Or(
//	    validators.Email(),
//	    validators.StringRegex(`^\d{11}$`), // 手机号
//	))
func Or[T any](validators ...func(T) error) func(T) error {
	return func(value T) error {
		var lastErr error
		for _, validator := range validators {
			if err := validator(value); err == nil {
				return nil
			} else {
				lastErr = err
			}
		}
		return lastErr
	}
}

// Not 创建反向验证器
//
// 参数:
//   - validator: 要反向的验证器
//
// 返回值:
//   - func(T) error: 反向验证器函数
//
// 功能说明:
//   - 验证器必须失败
//   - 验证器通过则返回错误
//
// 使用示例:
//
//	username.SetValidator(validators.Not(validators.StringOneOf("admin", "root")))
func Not[T any](validator func(T) error) func(T) error {
	return func(value T) error {
		if err := validator(value); err != nil {
			return nil
		}
		return fmt.Errorf("值 %v 不允许", value)
	}
}

// Optional 创建可选验证器
//
// 参数:
//   - validator: 条件验证器
//
// 返回值:
//   - func(T) error: 可选验证器函数
//
// 功能说明:
//   - 如果值非空则验证，空值跳过验证
//   - 适用于字符串和切片类型
//
// 使用示例:
//
//	email.SetValidator(validators.Optional(validators.Email()))
func Optional[T any](validator func(T) error) func(T) error {
	return func(value T) error {
		switch v := any(value).(type) {
		case string:
			if v == "" {
				return nil
			}
		case []T:
			if len(v) == 0 {
				return nil
			}
		}
		return validator(value)
	}
}

// FileExtension 创建文件扩展名验证器
//
// 参数:
//   - extensions: 允许的扩展名列表（不包含点号）
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证文件扩展名是否在允许的列表中
//   - 不在列表中返回错误
//
// 使用示例:
//
//	filename.SetValidator(validators.FileExtension("json", "yaml", "yml"))
func FileExtension(extensions ...string) func(string) error {
	allowedExts := make(map[string]bool)
	for _, ext := range extensions {
		allowedExts[strings.ToLower(ext)] = true
	}
	return func(value string) error {
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(value)), ".")
		if !allowedExts[ext] {
			return fmt.Errorf("文件扩展名 '.%s' 不在允许的列表中: %v", ext, extensions)
		}
		return nil
	}
}

// FileExists 创建文件存在性验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证文件是否存在
//   - 文件不存在返回错误
//
// 使用示例:
//
//	configFile.SetValidator(validators.FileExists())
func FileExists() func(string) error {
	return func(value string) error {
		info, err := os.Stat(value)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("文件 '%s' 不存在", value)
			}
			return fmt.Errorf("无法访问文件 '%s': %v", value, err)
		}
		if info.IsDir() {
			return fmt.Errorf("'%s' 是目录，不是文件", value)
		}
		return nil
	}
}

// DirExists 创建目录存在性验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证目录是否存在
//   - 目录不存在返回错误
//
// 使用示例:
//
//	outputDir.SetValidator(validators.DirExists())
func DirExists() func(string) error {
	return func(value string) error {
		info, err := os.Stat(value)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("目录 '%s' 不存在", value)
			}
			return fmt.Errorf("无法访问目录 '%s': %v", value, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("'%s' 是文件，不是目录", value)
		}
		return nil
	}
}

// IsNumeric 创建数字字符串验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否为有效的数字（整数或浮点数）
//   - 不是有效数字返回错误
//
// 使用示例:
//
//	number.SetValidator(validators.IsNumeric())
func IsNumeric() func(string) error {
	return func(value string) error {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("字符串 '%s' 不是有效的数字", value)
		}
		return nil
	}
}

// IsInteger 创建整数字符串验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否为有效的整数
//   - 不是有效整数返回错误
//
// 使用示例:
//
//	intNumber.SetValidator(validators.IsInteger())
func IsInteger() func(string) error {
	return func(value string) error {
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Errorf("字符串 '%s' 不是有效的整数", value)
		}
		return nil
	}
}

// IsPositiveInteger 创建正整数字符串验证器
//
// 返回值:
//   - func(string) error: 字符串验证器函数
//
// 功能说明:
//   - 验证字符串是否为有效的正整数
//   - 不是有效正整数返回错误
//
// 使用示例:
//
//	count.SetValidator(validators.IsPositiveInteger())
func IsPositiveInteger() func(string) error {
	return func(value string) error {
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("字符串 '%s' 不是有效的整数", value)
		}
		if num <= 0 {
			return fmt.Errorf("字符串 '%s' 不是正整数", value)
		}
		return nil
	}
}
