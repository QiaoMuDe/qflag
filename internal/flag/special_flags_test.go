package flag

import (
	"testing"
	"time"
)

// TestEnumFlag 测试枚举标志
func TestEnumFlag(t *testing.T) {
	// 创建枚举标志
	options := []string{"auto", "manual", "debug"}
	flag := NewEnumFlag("mode", "m", "运行模式", "auto", options)

	// 测试初始值
	if flag.Get() != "auto" {
		t.Errorf("Expected default value 'auto', got '%s'", flag.Get())
	}

	// 测试设置有效值
	for _, option := range options {
		err := flag.Set(option)
		if err != nil {
			t.Errorf("Unexpected error setting valid option '%s': %v", option, err)
		}

		if flag.Get() != option {
			t.Errorf("Expected '%s', got '%s'", option, flag.Get())
		}
	}

	// 测试设置无效值
	err := flag.Set("invalid")
	if err == nil {
		t.Error("Expected error for invalid enum value")
	}

	// 测试枚举值列表
	enumOptions := flag.EnumValues()
	if len(enumOptions) != len(options) {
		t.Errorf("Expected %d enum options, got %d", len(options), len(enumOptions))
	}

	// 检查所有选项都存在, 不检查顺序
	optionMap := make(map[string]bool)
	for _, option := range options {
		optionMap[option] = true
	}

	for _, enumOption := range enumOptions {
		if !optionMap[enumOption] {
			t.Errorf("Unexpected enum option '%s'", enumOption)
		}
	}

	// 测试基本属性
	if flag.Name() != "mode" {
		t.Errorf("Expected name 'mode', got '%s'", flag.Name())
	}

	if flag.ShortName() != "m" {
		t.Errorf("Expected short name 'm', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "运行模式" {
		t.Errorf("Expected description '运行模式', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "enum" {
		t.Errorf("Expected type 'enum', got '%s'", flag.Type().String())
	}
}

// TestDurationFlag 测试持续时间标志
func TestDurationFlag(t *testing.T) {
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

	expected = time.Hour + 30*time.Minute
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

	if flag.Type().String() != "duration" {
		t.Errorf("Expected type 'duration', got '%s'", flag.Type().String())
	}
}

// TestTimeFlag 测试时间标志
func TestTimeFlag(t *testing.T) {
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

	// 测试设置其他格式
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"Jan 2, 2006 15:04:05",
	}

	for _, layout := range layouts {
		testTime := "2023-12-31 23:59:59"
		err = flag.Set(testTime)
		if err != nil {
			t.Errorf("Unexpected error setting time with format '%s': %v", layout, err)
		}
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

	if flag.Type().String() != "time" {
		t.Errorf("Expected type 'time', got '%s'", flag.Type().String())
	}
}

// TestSizeFlag 测试大小标志
func TestSizeFlag(t *testing.T) {
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

	// 测试无效值
	err = flag.Set("invalid")
	if err == nil {
		t.Error("Expected error for invalid size value")
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

	if flag.Type().String() != "size" {
		t.Errorf("Expected type 'size', got '%s'", flag.Type().String())
	}
}
