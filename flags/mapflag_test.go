package flags

import (
	"reflect"
	"strings"
	"testing"
)

// TestMapFlag_BasicParsing 测试基本的键值对解析功能
func TestMapFlag_BasicParsing(t *testing.T) {
	flag := &MapFlag{}
	flag.SetDelimiters(",", "=") // 显式设置分隔符
	err := flag.Set("name=test,env=dev")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// 验证Get()返回正确的map
	result := flag.Get()
	expectedMap := map[string]string{"name": "test", "env": "dev"}
	if !reflect.DeepEqual(result, expectedMap) {
		t.Errorf("Get() returned %v, expected %v", result, expectedMap)
	}

	// 验证String()输出正确
	actualMap := make(map[string]string)
	parts := strings.Split(flag.String(), ",")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			actualMap[kv[0]] = kv[1]
		}
	}

	if actualMap["name"] != "test" || actualMap["env"] != "dev" || len(actualMap) != 2 {
		t.Errorf("String() returned map %v, expected {name:test, env:dev}", actualMap)
	}
}

// TestMapFlag_CustomDelimiters 测试自定义分隔符
func TestMapFlag_CustomDelimiters(t *testing.T) {
	flag := &MapFlag{}
	flag.SetDelimiters("; ", ":")

	err := flag.Set("name:test; env:prod")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	result := flag.Get()
	expected := map[string]string{"name": "test", "env": "prod"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestMapFlag_IgnoreCase 测试忽略键的大小写
func TestMapFlag_IgnoreCase(t *testing.T) {
	flag := &MapFlag{}
	flag.SetDelimiters(",", "=")
	flag.SetIgnoreCase(true)

	err := flag.Set("Name=test,NAME=override")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	result := flag.Get()
	// 所有键应该被转换为小写
	expected := map[string]string{"name": "override"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestMapFlag_Errors 测试错误情况处理
func TestMapFlag_Errors(t *testing.T) {
	flag := &MapFlag{}
	flag.SetDelimiters(",", "=")

	// 测试空值
	err := flag.Set("")
	if err == nil || !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("Expected empty value error, got %v", err)
	}

	// 测试格式错误的键值对
	err = flag.Set("invalid-key")
	if err == nil || !strings.Contains(err.Error(), "invalid key-value pair format") {
		t.Errorf("Expected format error, got %v", err)
	}

	// 测试空键
	err = flag.Set("=emptykey")
	if err == nil || !strings.Contains(err.Error(), "empty key") {
		t.Errorf("Expected empty key error, got %v", err)
	}

	// 测试空值
	err = flag.Set("key=")
	if err == nil || !strings.Contains(err.Error(), "empty value") {
		t.Errorf("Expected empty value error, got %v", err)
	}
}
