package flag

import (
	"testing"
)

// TestInt64Flag 测试64位整数标志
func TestInt64Flag(t *testing.T) {
	// 创建64位整数标志
	flag := NewInt64Flag("big-number", "b", "大数字", 123456789012345)

	// 测试初始值
	if flag.Get() != 123456789012345 {
		t.Errorf("Expected default value 123456789012345, got %d", flag.Get())
	}

	// 测试设置大数
	err := flag.Set("9223372036854775807") // int64最大值
	if err != nil {
		t.Errorf("Unexpected error setting max int64 value: %v", err)
	}

	if flag.Get() != 9223372036854775807 {
		t.Errorf("Expected max int64 value, got %d", flag.Get())
	}

	// 测试设置小数
	err = flag.Set("-9223372036854775808") // int64最小值
	if err != nil {
		t.Errorf("Unexpected error setting min int64 value: %v", err)
	}

	if flag.Get() != -9223372036854775808 {
		t.Errorf("Expected min int64 value, got %d", flag.Get())
	}

	// 测试超出范围的值
	err = flag.Set("9223372036854775808") // 超出int64范围
	if err == nil {
		t.Error("Expected error for value exceeding int64 range")
	}

	// 测试无效值
	err = flag.Set("invalid")
	if err == nil {
		t.Error("Expected error for invalid integer value")
	}

	// 测试基本属性
	if flag.Name() != "big-number" {
		t.Errorf("Expected name 'big-number', got '%s'", flag.Name())
	}

	if flag.ShortName() != "b" {
		t.Errorf("Expected short name 'b', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "大数字" {
		t.Errorf("Expected description '大数字', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "int64" {
		t.Errorf("Expected type 'int64', got '%s'", flag.Type().String())
	}
}

// TestUintFlag 测试无符号整数标志
func TestUintFlag(t *testing.T) {
	// 创建无符号整数标志
	flag := NewUintFlag("count", "c", "计数器", 100)

	// 测试初始值
	if flag.Get() != 100 {
		t.Errorf("Expected default value 100, got %d", flag.Get())
	}

	// 测试设置正数
	err := flag.Set("42")
	if err != nil {
		t.Errorf("Unexpected error setting positive value: %v", err)
	}

	if flag.Get() != 42 {
		t.Errorf("Expected 42, got %d", flag.Get())
	}

	// 测试设置大数
	err = flag.Set("4294967295") // 接近uint最大值
	if err != nil {
		t.Errorf("Unexpected error setting large value: %v", err)
	}

	// 测试设置负数
	err = flag.Set("-10")
	if err == nil {
		t.Error("Expected error for negative value")
	}

	// 测试无效值
	err = flag.Set("invalid")
	if err == nil {
		t.Error("Expected error for invalid integer value")
	}

	// 测试基本属性
	if flag.Name() != "count" {
		t.Errorf("Expected name 'count', got '%s'", flag.Name())
	}

	if flag.ShortName() != "c" {
		t.Errorf("Expected short name 'c', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "计数器" {
		t.Errorf("Expected description '计数器', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "uint" {
		t.Errorf("Expected type 'uint', got '%s'", flag.Type().String())
	}
}

// TestUint8Flag 测试8位无符号整数标志
func TestUint8Flag(t *testing.T) {
	// 创建8位无符号整数标志
	flag := NewUint8Flag("port", "p", "端口号", 80)

	// 测试初始值
	if flag.Get() != 80 {
		t.Errorf("Expected default value 80, got %d", flag.Get())
	}

	// 测试设置有效值
	err := flag.Set("255") // uint8最大值
	if err != nil {
		t.Errorf("Unexpected error setting max uint8 value: %v", err)
	}

	if flag.Get() != 255 {
		t.Errorf("Expected 255, got %d", flag.Get())
	}

	// 测试超出范围的值
	err = flag.Set("256") // 超出uint8范围
	if err == nil {
		t.Error("Expected error for value exceeding uint8 range")
	}

	// 测试设置负数
	err = flag.Set("-1")
	if err == nil {
		t.Error("Expected error for negative value")
	}

	// 测试基本属性
	if flag.Name() != "port" {
		t.Errorf("Expected name 'port', got '%s'", flag.Name())
	}

	if flag.ShortName() != "p" {
		t.Errorf("Expected short name 'p', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "端口号" {
		t.Errorf("Expected description '端口号', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "uint8" {
		t.Errorf("Expected type 'uint8', got '%s'", flag.Type().String())
	}
}

// TestUint16Flag 测试16位无符号整数标志
func TestUint16Flag(t *testing.T) {
	// 创建16位无符号整数标志
	flag := NewUint16Flag("timeout", "t", "超时时间", 30)

	// 测试初始值
	if flag.Get() != 30 {
		t.Errorf("Expected default value 30, got %d", flag.Get())
	}

	// 测试设置有效值
	err := flag.Set("65535") // uint16最大值
	if err != nil {
		t.Errorf("Unexpected error setting max uint16 value: %v", err)
	}

	if flag.Get() != 65535 {
		t.Errorf("Expected 65535, got %d", flag.Get())
	}

	// 测试超出范围的值
	err = flag.Set("65536") // 超出uint16范围
	if err == nil {
		t.Error("Expected error for value exceeding uint16 range")
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

	if flag.Type().String() != "uint16" {
		t.Errorf("Expected type 'uint16', got '%s'", flag.Type().String())
	}
}

// TestUint32Flag 测试32位无符号整数标志
func TestUint32Flag(t *testing.T) {
	// 创建32位无符号整数标志
	flag := NewUint32Flag("buffer-size", "b", "缓冲区大小", 1024)

	// 测试初始值
	if flag.Get() != 1024 {
		t.Errorf("Expected default value 1024, got %d", flag.Get())
	}

	// 测试设置有效值
	err := flag.Set("4294967295") // 接近uint32最大值
	if err != nil {
		t.Errorf("Unexpected error setting large value: %v", err)
	}

	// 测试超出范围的值
	err = flag.Set("4294967296") // 超出uint32范围
	if err == nil {
		t.Error("Expected error for value exceeding uint32 range")
	}

	// 测试基本属性
	if flag.Name() != "buffer-size" {
		t.Errorf("Expected name 'buffer-size', got '%s'", flag.Name())
	}

	if flag.ShortName() != "b" {
		t.Errorf("Expected short name 'b', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "缓冲区大小" {
		t.Errorf("Expected description '缓冲区大小', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "uint32" {
		t.Errorf("Expected type 'uint32', got '%s'", flag.Type().String())
	}
}

// TestUint64Flag 测试64位无符号整数标志
func TestUint64Flag(t *testing.T) {
	// 创建64位无符号整数标志
	flag := NewUint64Flag("file-size", "f", "文件大小", 1048576)

	// 测试初始值
	if flag.Get() != 1048576 {
		t.Errorf("Expected default value 1048576, got %d", flag.Get())
	}

	// 测试设置大数
	err := flag.Set("18446744073709551615") // 接近uint64最大值
	if err != nil {
		t.Errorf("Unexpected error setting large value: %v", err)
	}

	// 测试超出范围的值
	err = flag.Set("18446744073709551616") // 超出uint64范围
	if err == nil {
		t.Error("Expected error for value exceeding uint64 range")
	}

	// 测试设置负数
	err = flag.Set("-1")
	if err == nil {
		t.Error("Expected error for negative value")
	}

	// 测试基本属性
	if flag.Name() != "file-size" {
		t.Errorf("Expected name 'file-size', got '%s'", flag.Name())
	}

	if flag.ShortName() != "f" {
		t.Errorf("Expected short name 'f', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "文件大小" {
		t.Errorf("Expected description '文件大小', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "uint64" {
		t.Errorf("Expected type 'uint64', got '%s'", flag.Type().String())
	}
}
