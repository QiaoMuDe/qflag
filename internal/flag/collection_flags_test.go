package flag

import (
	"testing"
)

// TestStringSliceFlag 测试字符串切片标志
func TestStringSliceFlag(t *testing.T) {
	// 创建字符串切片标志
	defaultValue := []string{"default1", "default2"}
	flag := NewStringSliceFlag("paths", "p", "路径列表", defaultValue)

	// 测试初始值
	paths := flag.Get()
	if len(paths) != len(defaultValue) {
		t.Errorf("Expected %d default values, got %d", len(defaultValue), len(paths))
	}

	for i, v := range defaultValue {
		if paths[i] != v {
			t.Errorf("Expected default value '%s', got '%s'", v, paths[i])
		}
	}

	// 测试设置单个值
	err := flag.Set("single")
	if err != nil {
		t.Errorf("Unexpected error setting single value: %v", err)
	}

	paths = flag.Get()
	if len(paths) != 1 || paths[0] != "single" {
		t.Errorf("Expected ['single'], got %v", paths)
	}

	// 测试设置逗号分隔的值
	err = flag.Set("a,b,c")
	if err != nil {
		t.Errorf("Unexpected error setting comma-separated values: %v", err)
	}

	paths = flag.Get()
	if len(paths) != 3 || paths[0] != "a" || paths[1] != "b" || paths[2] != "c" {
		t.Errorf("Expected ['a', 'b', 'c'], got %v", paths)
	}

	// 测试设置空值
	err = flag.Set("")
	if err != nil {
		t.Errorf("Unexpected error setting empty value: %v", err)
	}

	paths = flag.Get()
	if len(paths) != 0 {
		t.Errorf("Expected empty slice, got %v", paths)
	}

	// 测试重置
	flag.Reset()
	paths = flag.Get()
	if len(paths) != len(defaultValue) {
		t.Errorf("Expected %d default values after reset, got %d", len(defaultValue), len(paths))
	}

	// 测试基本属性
	if flag.Name() != "paths" {
		t.Errorf("Expected name 'paths', got '%s'", flag.Name())
	}

	if flag.ShortName() != "p" {
		t.Errorf("Expected short name 'p', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "路径列表" {
		t.Errorf("Expected description '路径列表', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "[]string" {
		t.Errorf("Expected type '[]string', got '%s'", flag.Type().String())
	}
}

// TestIntSliceFlag 测试整数切片标志
func TestIntSliceFlag(t *testing.T) {
	// 创建整数切片标志
	defaultValue := []int{80, 443}
	flag := NewIntSliceFlag("ports", "p", "端口列表", defaultValue)

	// 测试初始值
	ports := flag.Get()
	if len(ports) != len(defaultValue) {
		t.Errorf("Expected %d default values, got %d", len(defaultValue), len(ports))
	}

	for i, v := range defaultValue {
		if ports[i] != v {
			t.Errorf("Expected default value %d, got %d", v, ports[i])
		}
	}

	// 测试设置单个值
	err := flag.Set("8080")
	if err != nil {
		t.Errorf("Unexpected error setting single value: %v", err)
	}

	ports = flag.Get()
	if len(ports) != 1 || ports[0] != 8080 {
		t.Errorf("Expected [8080], got %v", ports)
	}

	// 测试设置逗号分隔的值
	err = flag.Set("80,443,8080")
	if err != nil {
		t.Errorf("Unexpected error setting comma-separated values: %v", err)
	}

	ports = flag.Get()
	if len(ports) != 3 || ports[0] != 80 || ports[1] != 443 || ports[2] != 8080 {
		t.Errorf("Expected [80, 443, 8080], got %v", ports)
	}

	// 测试设置无效值
	err = flag.Set("80,invalid,8080")
	if err == nil {
		t.Error("Expected error for invalid integer value")
	}

	// 测试基本属性
	if flag.Name() != "ports" {
		t.Errorf("Expected name 'ports', got '%s'", flag.Name())
	}

	if flag.ShortName() != "p" {
		t.Errorf("Expected short name 'p', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "端口列表" {
		t.Errorf("Expected description '端口列表', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "[]int" {
		t.Errorf("Expected type '[]int', got '%s'", flag.Type().String())
	}
}

// TestInt64SliceFlag 测试64位整数切片标志
func TestInt64SliceFlag(t *testing.T) {
	// 创建64位整数切片标志
	defaultValue := []int64{1000000, 2000000}
	flag := NewInt64SliceFlag("large-numbers", "l", "大数字列表", defaultValue)

	// 测试初始值
	numbers := flag.Get()
	if len(numbers) != len(defaultValue) {
		t.Errorf("Expected %d default values, got %d", len(defaultValue), len(numbers))
	}

	for i, v := range defaultValue {
		if numbers[i] != v {
			t.Errorf("Expected default value %d, got %d", v, numbers[i])
		}
	}

	// 测试设置大数
	err := flag.Set("9223372036854775807")
	if err != nil {
		t.Errorf("Unexpected error setting large value: %v", err)
	}

	numbers = flag.Get()
	if len(numbers) != 1 || numbers[0] != 9223372036854775807 {
		t.Errorf("Expected [9223372036854775807], got %v", numbers)
	}

	// 测试设置多个大数
	err = flag.Set("1000000,2000000,3000000")
	if err != nil {
		t.Errorf("Unexpected error setting multiple large values: %v", err)
	}

	numbers = flag.Get()
	if len(numbers) != 3 || numbers[0] != 1000000 || numbers[1] != 2000000 || numbers[2] != 3000000 {
		t.Errorf("Expected [1000000, 2000000, 3000000], got %v", numbers)
	}

	// 测试基本属性
	if flag.Name() != "large-numbers" {
		t.Errorf("Expected name 'large-numbers', got '%s'", flag.Name())
	}

	if flag.ShortName() != "l" {
		t.Errorf("Expected short name 'l', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "大数字列表" {
		t.Errorf("Expected description '大数字列表', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "[]int64" {
		t.Errorf("Expected type '[]int64', got '%s'", flag.Type().String())
	}
}

// TestMapFlag 测试映射标志
func TestMapFlag(t *testing.T) {
	// 创建映射标志
	defaultValue := map[string]string{"key1": "value1", "key2": "value2"}
	flag := NewMapFlag("headers", "h", "HTTP头部", defaultValue)

	// 测试初始值
	headers := flag.Get()
	if len(headers) != len(defaultValue) {
		t.Errorf("Expected %d default values, got %d", len(defaultValue), len(headers))
	}

	for k, v := range defaultValue {
		if headers[k] != v {
			t.Errorf("Expected default value '%s' for key '%s', got '%s'", v, k, headers[k])
		}
	}

	// 测试设置单个键值对
	err := flag.Set("key3=value3")
	if err != nil {
		t.Errorf("Unexpected error setting single key-value pair: %v", err)
	}

	headers = flag.Get()
	if len(headers) != 1 || headers["key3"] != "value3" {
		t.Errorf("Expected map with key3=value3, got %v", headers)
	}

	// 测试设置多个键值对
	err = flag.Set("key1=value1,key2=value2,key3=value3")
	if err != nil {
		t.Errorf("Unexpected error setting multiple key-value pairs: %v", err)
	}

	headers = flag.Get()
	if len(headers) != 3 ||
		headers["key1"] != "value1" ||
		headers["key2"] != "value2" ||
		headers["key3"] != "value3" {
		t.Errorf("Expected map with key1=value1, key2=value2, key3=value3, got %v", headers)
	}

	// 测试设置无效格式
	err = flag.Set("invalid-format")
	if err == nil {
		t.Error("Expected error for invalid format")
	}

	// 测试设置空值
	err = flag.Set("")
	if err != nil {
		t.Errorf("Unexpected error setting empty value: %v", err)
	}

	headers = flag.Get()
	if len(headers) != 0 {
		t.Errorf("Expected empty map, got %v", headers)
	}

	// 测试基本属性
	if flag.Name() != "headers" {
		t.Errorf("Expected name 'headers', got '%s'", flag.Name())
	}

	if flag.ShortName() != "h" {
		t.Errorf("Expected short name 'h', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "HTTP头部" {
		t.Errorf("Expected description 'HTTP头部', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "map" {
		t.Errorf("Expected type 'map', got '%s'", flag.Type().String())
	}
}
