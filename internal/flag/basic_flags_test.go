package flag

import (
	"testing"
)

// TestStringFlag 测试字符串标志
func TestStringFlag(t *testing.T) {
	// 创建字符串标志
	flag := NewStringFlag("output", "o", "输出文件", "default.txt")

	// 测试初始值
	if flag.Get() != "default.txt" {
		t.Errorf("Expected default value 'default.txt', got '%s'", flag.Get())
	}

	// 测试设置值
	err := flag.Set("test.txt")
	if err != nil {
		t.Errorf("Unexpected error setting value: %v", err)
	}

	if flag.Get() != "test.txt" {
		t.Errorf("Expected 'test.txt', got '%s'", flag.Get())
	}

	// 测试是否已设置
	if !flag.IsSet() {
		t.Error("Expected flag to be set")
	}

	// 测试重置
	flag.Reset()
	if flag.Get() != "default.txt" {
		t.Errorf("Expected default value after reset, got '%s'", flag.Get())
	}

	if flag.IsSet() {
		t.Error("Expected flag to be not set after reset")
	}

	// 测试基本属性
	if flag.Name() != "output" {
		t.Errorf("Expected name 'output', got '%s'", flag.Name())
	}

	if flag.ShortName() != "o" {
		t.Errorf("Expected short name 'o', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "输出文件" {
		t.Errorf("Expected description '输出文件', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "string" {
		t.Errorf("Expected type 'string', got '%s'", flag.Type().String())
	}
}

// TestIntFlag 测试整数标志
func TestIntFlag(t *testing.T) {
	// 创建整数标志
	flag := NewIntFlag("count", "c", "计数器", 10)

	// 测试初始值
	if flag.Get() != 10 {
		t.Errorf("Expected default value 10, got %d", flag.Get())
	}

	// 测试设置正数
	err := flag.Set("42")
	if err != nil {
		t.Errorf("Unexpected error setting positive value: %v", err)
	}

	if flag.Get() != 42 {
		t.Errorf("Expected 42, got %d", flag.Get())
	}

	// 测试设置负数
	err = flag.Set("-10")
	if err != nil {
		t.Errorf("Unexpected error setting negative value: %v", err)
	}

	if flag.Get() != -10 {
		t.Errorf("Expected -10, got %d", flag.Get())
	}

	// 测试无效值
	err = flag.Set("invalid")
	if err == nil {
		t.Error("Expected error for invalid integer value")
	}

	// 测试重置
	flag.Reset()
	if flag.Get() != 10 {
		t.Errorf("Expected default value after reset, got %d", flag.Get())
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

	if flag.Type().String() != "int" {
		t.Errorf("Expected type 'int', got '%s'", flag.Type().String())
	}
}

// TestBoolFlag 测试布尔标志
func TestBoolFlag(t *testing.T) {
	// 创建布尔标志
	flag := NewBoolFlag("verbose", "v", "详细输出", false)

	// 测试初始值
	if flag.Get() != false {
		t.Errorf("Expected default value false, got %v", flag.Get())
	}

	// 测试设置true
	err := flag.Set("true")
	if err != nil {
		t.Errorf("Unexpected error setting true: %v", err)
	}

	if flag.Get() != true {
		t.Errorf("Expected true, got %v", flag.Get())
	}

	// 测试设置false
	err = flag.Set("false")
	if err != nil {
		t.Errorf("Unexpected error setting false: %v", err)
	}

	if flag.Get() != false {
		t.Errorf("Expected false, got %v", flag.Get())
	}

	// 测试其他有效值
	validValues := []string{"1", "0", "TRUE", "FALSE", "t", "f", "T", "F", "true", "false"}
	expectedResults := []bool{true, false, true, false, true, false, true, false, true, false}

	for i, value := range validValues {
		err = flag.Set(value)
		if err != nil {
			t.Errorf("Unexpected error setting '%s': %v", value, err)
		}

		if flag.Get() != expectedResults[i] {
			t.Errorf("Expected %v for '%s', got %v", expectedResults[i], value, flag.Get())
		}
	}

	// 测试重置
	flag.Reset()
	if flag.Get() != false {
		t.Errorf("Expected default value after reset, got %v", flag.Get())
	}

	// 测试基本属性
	if flag.Name() != "verbose" {
		t.Errorf("Expected name 'verbose', got '%s'", flag.Name())
	}

	if flag.ShortName() != "v" {
		t.Errorf("Expected short name 'v', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "详细输出" {
		t.Errorf("Expected description '详细输出', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "bool" {
		t.Errorf("Expected type 'bool', got '%s'", flag.Type().String())
	}
}

// TestFloat64Flag 测试浮点数标志
func TestFloat64Flag(t *testing.T) {
	// 创建浮点数标志
	flag := NewFloat64Flag("ratio", "r", "比例", 0.5)

	// 测试初始值
	if flag.Get() != 0.5 {
		t.Errorf("Expected default value 0.5, got %f", flag.Get())
	}

	// 测试设置整数
	err := flag.Set("42")
	if err != nil {
		t.Errorf("Unexpected error setting integer value: %v", err)
	}

	if flag.Get() != 42.0 {
		t.Errorf("Expected 42.0, got %f", flag.Get())
	}

	// 测试设置小数
	err = flag.Set("3.14159")
	if err != nil {
		t.Errorf("Unexpected error setting decimal value: %v", err)
	}

	if flag.Get() != 3.14159 {
		t.Errorf("Expected 3.14159, got %f", flag.Get())
	}

	// 测试设置负数
	err = flag.Set("-2.5")
	if err != nil {
		t.Errorf("Unexpected error setting negative value: %v", err)
	}

	if flag.Get() != -2.5 {
		t.Errorf("Expected -2.5, got %f", flag.Get())
	}

	// 测试无效值
	err = flag.Set("invalid")
	if err == nil {
		t.Error("Expected error for invalid float value")
	}

	// 测试重置
	flag.Reset()
	if flag.Get() != 0.5 {
		t.Errorf("Expected default value after reset, got %f", flag.Get())
	}

	// 测试基本属性
	if flag.Name() != "ratio" {
		t.Errorf("Expected name 'ratio', got '%s'", flag.Name())
	}

	if flag.ShortName() != "r" {
		t.Errorf("Expected short name 'r', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "比例" {
		t.Errorf("Expected description '比例', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "float64" {
		t.Errorf("Expected type 'float64', got '%s'", flag.Type().String())
	}
}
