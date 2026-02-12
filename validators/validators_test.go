package validators

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIntRange(t *testing.T) {
	validator := IntRange(1, 10)

	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{"有效值-最小边界", 1, false},
		{"有效值-中间值", 5, false},
		{"有效值-最大边界", 10, false},
		{"无效值-小于最小值", 0, true},
		{"无效值-大于最大值", 11, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUintRange(t *testing.T) {
	validator := UintRange(0, 100)

	tests := []struct {
		name    string
		value   uint
		wantErr bool
	}{
		{"有效值-最小边界", 0, false},
		{"有效值-中间值", 50, false},
		{"有效值-最大边界", 100, false},
		{"无效值-大于最大值", 101, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UintRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUint8Range(t *testing.T) {
	validator := Uint8Range(0, 255)

	tests := []struct {
		name    string
		value   uint8
		wantErr bool
	}{
		{"有效值-最小边界", 0, false},
		{"有效值-中间值", 128, false},
		{"有效值-最大边界", 255, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint8Range() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUint16Range(t *testing.T) {
	validator := Uint16Range(1, 65535)

	tests := []struct {
		name    string
		value   uint16
		wantErr bool
	}{
		{"有效值-最小边界", 1, false},
		{"有效值-中间值", 32768, false},
		{"有效值-最大边界", 65535, false},
		{"无效值-小于最小值", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint16Range() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUint32Range(t *testing.T) {
	validator := Uint32Range(0, 4294967295)

	tests := []struct {
		name    string
		value   uint32
		wantErr bool
	}{
		{"有效值-最小边界", 0, false},
		{"有效值-中间值", 2147483648, false},
		{"有效值-最大边界", 4294967295, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint32Range() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUint64Range(t *testing.T) {
	validator := Uint64Range(0, 18446744073709551615)

	tests := []struct {
		name    string
		value   uint64
		wantErr bool
	}{
		{"有效值-最小边界", 0, false},
		{"有效值-中间值", 9223372036854775808, false},
		{"有效值-最大边界", 18446744073709551615, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64Range() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt64Range(t *testing.T) {
	validator := Int64Range(0, 9999999999)

	tests := []struct {
		name    string
		value   int64
		wantErr bool
	}{
		{"有效值-最小边界", 0, false},
		{"有效值-中间值", 5000000000, false},
		{"有效值-最大边界", 9999999999, false},
		{"无效值-小于最小值", -1, true},
		{"无效值-大于最大值", 10000000000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64Range() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloat64Range(t *testing.T) {
	validator := Float64Range(-50.0, 100.0)

	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{"有效值-最小边界", -50.0, false},
		{"有效值-中间值", 25.0, false},
		{"有效值-最大边界", 100.0, false},
		{"无效值-小于最小值", -50.1, true},
		{"无效值-大于最大值", 100.1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Float64Range() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPositive(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"int-正数", 1, false},
		{"int-零", 0, true},
		{"int-负数", -1, true},
		{"int64-正数", int64(1), false},
		{"int64-零", int64(0), true},
		{"int64-负数", int64(-1), true},
		{"uint-正数", uint(1), false},
		{"uint-零", uint(0), true},
		{"uint64-正数", uint64(1), false},
		{"uint64-零", uint64(0), true},
		{"float64-正数", 1.0, false},
		{"float64-零", 0.0, true},
		{"float64-负数", -1.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch v := tt.value.(type) {
			case int:
				err = Positive[int]()(v)
			case int64:
				err = Positive[int64]()(v)
			case uint:
				err = Positive[uint]()(v)
			case uint64:
				err = Positive[uint64]()(v)
			case float64:
				err = Positive[float64]()(v)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Positive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNonNegative(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"int-正数", 1, false},
		{"int-零", 0, false},
		{"int-负数", -1, true},
		{"int64-正数", int64(1), false},
		{"int64-零", int64(0), false},
		{"int64-负数", int64(-1), true},
		{"float64-正数", 1.0, false},
		{"float64-零", 0.0, false},
		{"float64-负数", -1.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch v := tt.value.(type) {
			case int:
				err = NonNegative[int]()(v)
			case int64:
				err = NonNegative[int64]()(v)
			case float64:
				err = NonNegative[float64]()(v)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("NonNegative() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringLength(t *testing.T) {
	validator := StringLength(3, 20)

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-最小边界", "abc", false},
		{"有效值-中间值", "hello", false},
		{"有效值-最大边界", strings.Repeat("a", 20), false},
		{"无效值-小于最小值", "ab", true},
		{"无效值-大于最大值", strings.Repeat("a", 21), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringLengthNoLimit(t *testing.T) {
	validator := StringLength(0, 0)

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"空字符串", "", false},
		{"短字符串", "a", false},
		{"长字符串", strings.Repeat("a", 1000), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringMinLength(t *testing.T) {
	validator := StringMinLength(3)

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-最小边界", "abc", false},
		{"有效值-更长", "hello", false},
		{"无效值-太短", "ab", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringMinLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringMaxLength(t *testing.T) {
	validator := StringMaxLength(20)

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-最大边界", strings.Repeat("a", 20), false},
		{"有效值-更短", "hello", false},
		{"无效值-太长", strings.Repeat("a", 21), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringMaxLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringRegex(t *testing.T) {
	validator := StringRegex(`^[a-zA-Z0-9]+$`)

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-字母", "hello", false},
		{"有效值-数字", "123", false},
		{"有效值-字母数字", "abc123", false},
		{"无效值-包含特殊字符", "hello!", true},
		{"无效值-包含空格", "hello world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringRegex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringPrefix(t *testing.T) {
	validator := StringPrefix("APP_")

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-正确前缀", "APP_NAME", false},
		{"无效值-错误前缀", "NAME", true},
		{"无效值-无前缀", "NAME", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringPrefix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringSuffix(t *testing.T) {
	validator := StringSuffix(".txt")

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-正确后缀", "file.txt", false},
		{"无效值-错误后缀", "file.json", true},
		{"无效值-无后缀", "file", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringSuffix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringContains(t *testing.T) {
	validator := StringContains("/data/")

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-包含子串", "/data/file.txt", false},
		{"无效值-不包含子串", "/var/file.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringContains() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringNotEmpty(t *testing.T) {
	validator := StringNotEmpty()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-非空", "hello", false},
		{"无效值-空", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringNotEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringOneOf(t *testing.T) {
	validator := StringOneOf("debug", "info", "warn", "error")

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-debug", "debug", false},
		{"有效值-info", "info", false},
		{"有效值-warn", "warn", false},
		{"有效值-error", "error", false},
		{"无效值-不在列表", "trace", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringOneOf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringCharset(t *testing.T) {
	tests := []struct {
		name    string
		charset string
		value   string
		wantErr bool
	}{
		{"alnum-有效", "alnum", "abc123", false},
		{"alnum-无效", "alnum", "abc!", true},
		{"alpha-有效", "alpha", "abc", false},
		{"alpha-无效", "alpha", "abc1", true},
		{"digit-有效", "digit", "123", false},
		{"digit-无效", "digit", "123a", true},
		{"hex-有效", "hex", "abc123", false},
		{"hex-无效", "hex", "xyz", true},
		{"lower-有效", "lower", "abc", false},
		{"lower-无效", "lower", "ABC", true},
		{"upper-有效", "upper", "ABC", false},
		{"upper-无效", "upper", "abc", true},
		{"无效字符集", "invalid", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := StringCharset(tt.charset)
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringCharset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	validator := Email()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-简单邮箱", "user@example.com", false},
		{"有效值-带子域名", "user@mail.example.com", false},
		{"有效值-带加号", "user+tag@example.com", false},
		{"无效值-缺少@", "userexample.com", true},
		{"无效值-缺少域名", "user@", true},
		{"无效值-空", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		name    string
		scheme  string
		value   string
		wantErr bool
	}{
		{"无协议限制-http", "", "http://example.com", false},
		{"无协议限制-https", "", "https://example.com", false},
		{"限制http-匹配", "http", "http://example.com", false},
		{"限制http-不匹配", "http", "https://example.com", true},
		{"限制https-匹配", "https", "https://example.com", false},
		{"限制https-不匹配", "https", "http://example.com", true},
		{"无效URL", "", "not a url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := URL(tt.scheme)
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIPv4(t *testing.T) {
	validator := IPv4()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-标准", "192.168.1.1", false},
		{"有效值-最小值", "0.0.0.0", false},
		{"有效值-最大值", "255.255.255.255", false},
		{"无效值-IPv6", "2001:db8::1", true},
		{"无效值-格式错误", "192.168.1", true},
		{"无效值-超出范围", "256.168.1.1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("IPv4() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIPv6(t *testing.T) {
	validator := IPv6()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-完整", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		{"有效值-压缩", "2001:db8::1", false},
		{"有效值-本地", "::1", false},
		{"无效值-IPv4", "192.168.1.1", true},
		{"无效值-格式错误", "2001:db8", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("IPv6() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIP(t *testing.T) {
	validator := IP()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-IPv4", "192.168.1.1", false},
		{"有效值-IPv6", "2001:db8::1", false},
		{"无效值-格式错误", "not an ip", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("IP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPort(t *testing.T) {
	validator := Port()

	tests := []struct {
		name    string
		value   uint16
		wantErr bool
	}{
		{"有效值-最小边界", 1, false},
		{"有效值-常用端口", 8080, false},
		{"有效值-最大边界", 65535, false},
		{"无效值-零", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Port() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHostname(t *testing.T) {
	validator := Hostname()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-简单", "localhost", false},
		{"有效值-域名", "example.com", false},
		{"有效值-子域名", "sub.example.com", false},
		{"无效值-包含下划线", "my_host", true},
		{"无效值-包含空格", "my host", true},
		{"无效值-以连字符开头", "-host", true},
		{"无效值-以连字符结尾", "host-", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hostname() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDurationMin(t *testing.T) {
	validator := DurationMin(time.Second)

	tests := []struct {
		name    string
		value   time.Duration
		wantErr bool
	}{
		{"有效值-最小边界", time.Second, false},
		{"有效值-更大值", time.Minute, false},
		{"无效值-小于最小值", time.Millisecond, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurationMin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDurationMax(t *testing.T) {
	validator := DurationMax(time.Hour)

	tests := []struct {
		name    string
		value   time.Duration
		wantErr bool
	}{
		{"有效值-最大边界", time.Hour, false},
		{"有效值-更小值", time.Minute, false},
		{"无效值-大于最大值", 2 * time.Hour, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurationMax() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDurationRange(t *testing.T) {
	validator := DurationRange(time.Second, time.Minute)

	tests := []struct {
		name    string
		value   time.Duration
		wantErr bool
	}{
		{"有效值-最小边界", time.Second, false},
		{"有效值-中间值", 30 * time.Second, false},
		{"有效值-最大边界", time.Minute, false},
		{"无效值-小于最小值", time.Millisecond, true},
		{"无效值-大于最大值", 2 * time.Minute, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurationRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeAfter(t *testing.T) {
	now := time.Now()
	validator := TimeAfter(now)

	tests := []struct {
		name    string
		value   time.Time
		wantErr bool
	}{
		{"有效值-未来", now.Add(time.Hour), false},
		{"无效值-现在", now, true},
		{"无效值-过去", now.Add(-time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeAfter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeBefore(t *testing.T) {
	now := time.Now()
	validator := TimeBefore(now)

	tests := []struct {
		name    string
		value   time.Time
		wantErr bool
	}{
		{"有效值-过去", now.Add(-time.Hour), false},
		{"无效值-现在", now, true},
		{"无效值-未来", now.Add(time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeBefore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeRange(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour)
	validator := TimeRange(start, end)

	tests := []struct {
		name    string
		value   time.Time
		wantErr bool
	}{
		{"有效值-开始", start, false},
		{"有效值-中间", start.Add(30 * time.Minute), false},
		{"有效值-结束", end, false},
		{"无效值-早于开始", start.Add(-time.Minute), true},
		{"无效值-晚于结束", end.Add(time.Minute), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSliceLength(t *testing.T) {
	validator := SliceLength[string](1, 5)

	tests := []struct {
		name    string
		value   []string
		wantErr bool
	}{
		{"有效值-最小边界", []string{"a"}, false},
		{"有效值-中间值", []string{"a", "b", "c"}, false},
		{"有效值-最大边界", []string{"a", "b", "c", "d", "e"}, false},
		{"无效值-空", []string{}, true},
		{"无效值-超过最大值", []string{"a", "b", "c", "d", "e", "f"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSliceMinLength(t *testing.T) {
	validator := SliceMinLength[string](2)

	tests := []struct {
		name    string
		value   []string
		wantErr bool
	}{
		{"有效值-最小边界", []string{"a", "b"}, false},
		{"有效值-更长", []string{"a", "b", "c"}, false},
		{"无效值-太短", []string{"a"}, true},
		{"无效值-空", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceMinLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSliceMaxLength(t *testing.T) {
	validator := SliceMaxLength[string](5)

	tests := []struct {
		name    string
		value   []string
		wantErr bool
	}{
		{"有效值-最大边界", []string{"a", "b", "c", "d", "e"}, false},
		{"有效值-更短", []string{"a", "b", "c"}, false},
		{"无效值-太长", []string{"a", "b", "c", "d", "e", "f"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceMaxLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSliceNotEmpty(t *testing.T) {
	validator := SliceNotEmpty[string]()

	tests := []struct {
		name    string
		value   []string
		wantErr bool
	}{
		{"有效值-非空", []string{"a"}, false},
		{"无效值-空", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceNotEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSliceUnique(t *testing.T) {
	validator := SliceUnique[string]()

	tests := []struct {
		name    string
		value   []string
		wantErr bool
	}{
		{"有效值-唯一", []string{"a", "b", "c"}, false},
		{"无效值-重复", []string{"a", "b", "a"}, true},
		{"无效值-多个重复", []string{"a", "a", "a"}, true},
		{"有效值-空", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceUnique() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSliceContains(t *testing.T) {
	validator := SliceContains[string]("default")

	tests := []struct {
		name    string
		value   []string
		wantErr bool
	}{
		{"有效值-包含", []string{"default", "option"}, false},
		{"无效值-不包含", []string{"option", "other"}, true},
		{"无效值-空", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceContains() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	validator := MapKeys[string]("host", "port", "timeout")

	tests := []struct {
		name    string
		value   map[string]string
		wantErr bool
	}{
		{"有效值-所有键允许", map[string]string{"host": "localhost", "port": "8080"}, false},
		{"无效值-包含不允许的键", map[string]string{"host": "localhost", "invalid": "value"}, true},
		{"有效值-空", map[string]string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMapMinSize(t *testing.T) {
	validator := MapMinSize[string, int](2)

	tests := []struct {
		name    string
		value   map[string]int
		wantErr bool
	}{
		{"有效值-最小边界", map[string]int{"a": 1, "b": 2}, false},
		{"有效值-更大", map[string]int{"a": 1, "b": 2, "c": 3}, false},
		{"无效值-太小", map[string]int{"a": 1}, true},
		{"无效值-空", map[string]int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapMinSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMapMaxSize(t *testing.T) {
	validator := MapMaxSize[string, int](5)

	tests := []struct {
		name    string
		value   map[string]int
		wantErr bool
	}{
		{"有效值-最大边界", map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}, false},
		{"有效值-更小", map[string]int{"a": 1, "b": 2}, false},
		{"无效值-太大", map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapMaxSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMapRequiredKeys(t *testing.T) {
	validator := MapRequiredKeys[string]("host", "port")

	tests := []struct {
		name    string
		value   map[string]string
		wantErr bool
	}{
		{"有效值-包含所有必需键", map[string]string{"host": "localhost", "port": "8080"}, false},
		{"有效值-包含额外键", map[string]string{"host": "localhost", "port": "8080", "timeout": "30"}, false},
		{"无效值-缺少host", map[string]string{"port": "8080"}, true},
		{"无效值-缺少port", map[string]string{"host": "localhost"}, true},
		{"无效值-缺少所有", map[string]string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapRequiredKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnd(t *testing.T) {
	validator := And(
		StringMinLength(3),
		StringMaxLength(10),
	)

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-满足所有条件", "hello", false},
		{"无效值-太短", "hi", true},
		{"无效值-太长", strings.Repeat("a", 11), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("And() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOr(t *testing.T) {
	validator := Or(
		StringOneOf("admin", "root"),
		StringRegex(`^\d{11}$`),
	)

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-匹配第一个", "admin", false},
		{"有效值-匹配第二个", "13800138000", false},
		{"无效值-都不匹配", "user", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Or() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNot(t *testing.T) {
	validator := Not(StringOneOf("admin", "root"))

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-不在列表", "user", false},
		{"无效值-在列表", "admin", true},
		{"无效值-在列表", "root", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Not() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptional(t *testing.T) {
	validator := Optional(StringMinLength(3))

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-空", "", false},
		{"有效值-满足条件", "hello", false},
		{"无效值-不满足条件", "hi", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Optional() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	validator := FileExists()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-文件存在", tmpFile, false},
		{"无效值-文件不存在", filepath.Join(tmpDir, "notexist.txt"), true},
		{"无效值-是目录", tmpDir, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirExists(t *testing.T) {
	validator := DirExists()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-目录存在", tmpDir, false},
		{"无效值-目录不存在", filepath.Join(tmpDir, "notexist"), true},
		{"无效值-是文件", tmpFile, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DirExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileExtension(t *testing.T) {
	validator := FileExtension("json", "yaml", "yml")

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-json", "config.json", false},
		{"有效值-yaml", "config.yaml", false},
		{"有效值-yml", "config.yml", false},
		{"有效值-大写", "CONFIG.JSON", false},
		{"无效值-其他扩展名", "config.txt", true},
		{"无效值-无扩展名", "config", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileExtension() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	validator := IsNumeric()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-整数", "123", false},
		{"有效值-浮点数", "123.45", false},
		{"有效值-负数", "-123", false},
		{"无效值-非数字", "abc", true},
		{"无效值-混合", "123abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsNumeric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsInteger(t *testing.T) {
	validator := IsInteger()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-正整数", "123", false},
		{"有效值-负整数", "-123", false},
		{"无效值-浮点数", "123.45", true},
		{"无效值-非数字", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsInteger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsPositiveInteger(t *testing.T) {
	validator := IsPositiveInteger()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效值-正整数", "123", false},
		{"无效值-零", "0", true},
		{"无效值-负整数", "-123", true},
		{"无效值-浮点数", "123.45", true},
		{"无效值-非数字", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsPositiveInteger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
