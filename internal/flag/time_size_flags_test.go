package flag

import (
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/internal/types"
)

// TestTimeSizeDurationFlag 测试持续时间标志
func TestTimeSizeDurationFlag(t *testing.T) {
	// 创建持续时间标志
	defaultValue := 30 * time.Second
	flag := NewDurationFlag("timeout", "t", "超时时间", defaultValue)

	// 测试初始值
	if flag.Get() != defaultValue {
		t.Errorf("Expected default value %v, got %v", defaultValue, flag.Get())
	}

	// 测试设置秒
	err := flag.Set("30s")
	if err != nil {
		t.Errorf("Unexpected error setting seconds: %v", err)
	}

	expected := 30 * time.Second
	if flag.Get() != expected {
		t.Errorf("Expected %v, got %v", expected, flag.Get())
	}

	// 测试设置分钟
	err = flag.Set("5m")
	if err != nil {
		t.Errorf("Unexpected error setting minutes: %v", err)
	}

	expected = 5 * time.Minute
	if flag.Get() != expected {
		t.Errorf("Expected %v, got %v", expected, flag.Get())
	}

	// 测试设置小时
	err = flag.Set("2h")
	if err != nil {
		t.Errorf("Unexpected error setting hours: %v", err)
	}

	expected = 2 * time.Hour
	if flag.Get() != expected {
		t.Errorf("Expected %v, got %v", expected, flag.Get())
	}

	// 测试设置复合值
	err = flag.Set("1h30m")
	if err != nil {
		t.Errorf("Unexpected error setting composite value: %v", err)
	}

	expected = 1*time.Hour + 30*time.Minute
	if flag.Get() != expected {
		t.Errorf("Expected %v, got %v", expected, flag.Get())
	}

	// 测试设置毫秒
	err = flag.Set("500ms")
	if err != nil {
		t.Errorf("Unexpected error setting milliseconds: %v", err)
	}

	expected = 500 * time.Millisecond
	if flag.Get() != expected {
		t.Errorf("Expected %v, got %v", expected, flag.Get())
	}

	// 测试设置微秒
	err = flag.Set("1000us")
	if err != nil {
		t.Errorf("Unexpected error setting microseconds: %v", err)
	}

	expected = 1000 * time.Microsecond
	if flag.Get() != expected {
		t.Errorf("Expected %v, got %v", expected, flag.Get())
	}

	// 测试设置纳秒
	err = flag.Set("100ns")
	if err != nil {
		t.Errorf("Unexpected error setting nanoseconds: %v", err)
	}

	expected = 100 * time.Nanosecond
	if flag.Get() != expected {
		t.Errorf("Expected %v, got %v", expected, flag.Get())
	}

	// 测试负数
	err = flag.Set("-1h")
	if err != nil {
		t.Errorf("Unexpected error setting negative value: %v", err)
	}

	expected = -1 * time.Hour
	if flag.Get() != expected {
		t.Errorf("Expected %v, got %v", expected, flag.Get())
	}

	// 测试小数
	err = flag.Set("1.5h")
	if err != nil {
		t.Errorf("Unexpected error setting fractional value: %v", err)
	}

	expected = 90 * time.Minute
	if flag.Get() != expected {
		t.Errorf("Expected %v, got %v", expected, flag.Get())
	}

	// 测试无效值
	err = flag.Set("invalid")
	if err == nil {
		t.Error("Expected error for invalid duration value")
	}

	// 测试基本属性
	if flag.Name() != "timeout" {
		t.Errorf("Expected name 'timeout', got '%s'", flag.Name())
	}

	if flag.ShortName() != "t" {
		t.Errorf("Expected short name 't', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "超时时间" {
		t.Errorf("Expected description '超时时间', got '%s'", flag.Desc())
	}

	if flag.Type() != types.FlagTypeDuration {
		t.Errorf("Expected type 'duration', got '%s'", flag.Type())
	}

	// 测试重置
	flag.Reset()
	if flag.Get() != defaultValue {
		t.Errorf("Expected reset value %v, got %v", defaultValue, flag.Get())
	}

	if flag.IsSet() {
		t.Error("Expected IsSet() to be false after reset")
	}
}

// TestTimeSizeTimeFlag 测试时间标志
func TestTimeSizeTimeFlag(t *testing.T) {
	// 创建时间标志
	defaultValue := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	flag := NewTimeFlag("start-time", "s", "开始时间", defaultValue)

	// 测试初始值
	if !flag.Get().Equal(defaultValue) {
		t.Errorf("Expected default value %v, got %v", defaultValue, flag.Get())
	}

	// 测试设置RFC3339格式
	rfc3339Time := "2023-12-31T23:59:59Z"
	err := flag.Set(rfc3339Time)
	if err != nil {
		t.Errorf("Unexpected error setting RFC3339 time: %v", err)
	}

	parsedTime, _ := time.Parse(time.RFC3339, rfc3339Time)
	if !flag.Get().Equal(parsedTime) {
		t.Errorf("Expected %v, got %v", parsedTime, flag.Get())
	}

	// 测试设置日期格式
	dateTime := "2023-12-31"
	err = flag.Set(dateTime)
	if err != nil {
		t.Errorf("Unexpected error setting date time: %v", err)
	}

	parsedTime, _ = time.Parse("2006-01-02", dateTime)
	if !flag.Get().Equal(parsedTime) {
		t.Errorf("Expected %v, got %v", parsedTime, flag.Get())
	}

	// 测试设置日期时间格式
	dateTime = "2023-12-31 23:59:59"
	err = flag.Set(dateTime)
	if err != nil {
		t.Errorf("Unexpected error setting date time with space: %v", err)
	}

	parsedTime, _ = time.Parse("2006-01-02 15:04:05", dateTime)
	if !flag.Get().Equal(parsedTime) {
		t.Errorf("Expected %v, got %v", parsedTime, flag.Get())
	}

	// 测试设置时间格式
	timeOnly := "15:04:05"
	err = flag.Set(timeOnly)
	if err != nil {
		t.Errorf("Unexpected error setting time only: %v", err)
	}

	// 测试无效时间
	err = flag.Set("invalid time")
	if err == nil {
		t.Error("Expected error for invalid time value")
	}

	// 测试基本属性
	if flag.Name() != "start-time" {
		t.Errorf("Expected name 'start-time', got '%s'", flag.Name())
	}

	if flag.ShortName() != "s" {
		t.Errorf("Expected short name 's', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "开始时间" {
		t.Errorf("Expected description '开始时间', got '%s'", flag.Desc())
	}

	if flag.Type() != types.FlagTypeTime {
		t.Errorf("Expected type 'time', got '%s'", flag.Type())
	}

	// 测试重置
	flag.Reset()
	if !flag.Get().Equal(defaultValue) {
		t.Errorf("Expected reset value %v, got %v", defaultValue, flag.Get())
	}

	if flag.IsSet() {
		t.Error("Expected IsSet() to be false after reset")
	}

	// 测试 SetWithFormat 方法
	customTime := "2023-12-31"
	customFormat := "2006-01-02"
	err = flag.SetWithFormat(customTime, customFormat)
	if err != nil {
		t.Errorf("Unexpected error setting time with custom format: %v", err)
	}

	parsedTime, _ = time.Parse(customFormat, customTime)
	if !flag.Get().Equal(parsedTime) {
		t.Errorf("Expected %v, got %v", parsedTime, flag.Get())
	}

	// 测试 GetFormat 方法
	if flag.GetFormat() != customFormat {
		t.Errorf("Expected format '%s', got '%s'", customFormat, flag.GetFormat())
	}

	// 测试 FormatTime 方法
	formattedTime := flag.FormatTime(time.Now())
	if formattedTime == "" {
		t.Error("Expected non-empty formatted time")
	}
}

// TestTimeSizeSizeFlag 测试大小标志
func TestTimeSizeSizeFlag(t *testing.T) {
	// 创建大小标志
	defaultValue := int64(1024)
	flag := NewSizeFlag("max-size", "s", "最大大小", defaultValue)

	// 测试初始值
	if flag.Get() != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, flag.Get())
	}

	// 测试设置字节数
	err := flag.Set("1024")
	if err != nil {
		t.Errorf("Unexpected error setting bytes: %v", err)
	}

	if flag.Get() != 1024 {
		t.Errorf("Expected 1024, got %d", flag.Get())
	}

	// 测试设置KB
	err = flag.Set("1KB")
	if err != nil {
		t.Errorf("Unexpected error setting KB: %v", err)
	}

	if flag.Get() != 1000 {
		t.Errorf("Expected 1000, got %d", flag.Get())
	}

	// 测试设置MB
	err = flag.Set("1MB")
	if err != nil {
		t.Errorf("Unexpected error setting MB: %v", err)
	}

	if flag.Get() != 1000000 {
		t.Errorf("Expected 1000000, got %d", flag.Get())
	}

	// 测试设置GB
	err = flag.Set("1GB")
	if err != nil {
		t.Errorf("Unexpected error setting GB: %v", err)
	}

	if flag.Get() != 1000000000 {
		t.Errorf("Expected 1000000000, got %d", flag.Get())
	}

	// 测试设置TB
	err = flag.Set("1TB")
	if err != nil {
		t.Errorf("Unexpected error setting TB: %v", err)
	}

	if flag.Get() != 1000000000000 {
		t.Errorf("Expected 1000000000000, got %d", flag.Get())
	}

	// 测试设置PB
	err = flag.Set("1PB")
	if err != nil {
		t.Errorf("Unexpected error setting PB: %v", err)
	}

	if flag.Get() != 1000000000000000 {
		t.Errorf("Expected 1000000000000000, got %d", flag.Get())
	}

	// 测试设置二进制单位
	err = flag.Set("1KiB")
	if err != nil {
		t.Errorf("Unexpected error setting KiB: %v", err)
	}

	if flag.Get() != 1024 {
		t.Errorf("Expected 1024, got %d", flag.Get())
	}

	err = flag.Set("1MiB")
	if err != nil {
		t.Errorf("Unexpected error setting MiB: %v", err)
	}

	if flag.Get() != 1048576 {
		t.Errorf("Expected 1048576, got %d", flag.Get())
	}

	err = flag.Set("1GiB")
	if err != nil {
		t.Errorf("Unexpected error setting GiB: %v", err)
	}

	if flag.Get() != 1073741824 {
		t.Errorf("Expected 1073741824, got %d", flag.Get())
	}

	// 测试小数
	err = flag.Set("1.5MB")
	if err != nil {
		t.Errorf("Unexpected error setting fractional MB: %v", err)
	}

	if flag.Get() != 1500000 {
		t.Errorf("Expected 1500000, got %d", flag.Get())
	}

	// 测试大小写不敏感
	err = flag.Set("1kb")
	if err != nil {
		t.Errorf("Unexpected error setting lowercase kb: %v", err)
	}

	if flag.Get() != 1000 {
		t.Errorf("Expected 1000, got %d", flag.Get())
	}

	// 测试无效值
	err = flag.Set("invalid")
	if err == nil {
		t.Error("Expected error for invalid size value")
	}

	// 测试负数
	err = flag.Set("-1MB")
	if err == nil {
		t.Error("Expected error for negative size value")
	}

	// 测试空值
	err = flag.Set("")
	if err == nil {
		t.Error("Expected error for empty size value")
	}

	// 测试基本属性
	if flag.Name() != "max-size" {
		t.Errorf("Expected name 'max-size', got '%s'", flag.Name())
	}

	if flag.ShortName() != "s" {
		t.Errorf("Expected short name 's', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "最大大小" {
		t.Errorf("Expected description '最大大小', got '%s'", flag.Desc())
	}

	if flag.Type() != types.FlagTypeSize {
		t.Errorf("Expected type 'size', got '%s'", flag.Type())
	}

	// 测试重置
	flag.Reset()
	if flag.Get() != defaultValue {
		t.Errorf("Expected reset value %d, got %d", defaultValue, flag.Get())
	}

	if flag.IsSet() {
		t.Error("Expected IsSet() to be false after reset")
	}
}
