package flags

import (
	"testing"
)

// TestStringFlag_BasicFunctionality 测试StringFlag的基本功能
func TestStringFlag_BasicFunctionality(t *testing.T) {
	flag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			defValue: "default",
			value:    new(string),
		},
	}

	// 测试默认值
	if flag.GetDefault() != "default" {
		t.Errorf("默认值应为'default', 实际为'%s'", flag.Get())
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试设置值
	testValue := "test string"
	if err := flag.Set(testValue); err != nil {
		t.Fatalf("设置值失败: %v", err)
	}
	if flag.Get() != testValue {
		t.Errorf("期望值'%s', 实际值'%s'", testValue, flag.Get())
	}
	if !flag.IsSet() {
		t.Error("设置值后IsSet应返回true")
	}

	// 测试重置功能
	flag.Reset()
	if flag.Get() != "default" {
		t.Errorf("重置后应返回默认值'default', 实际为'%s'", flag.Get())
	}
}

// TestStringFlag_Methods 测试字符串特有方法
func TestStringFlag_Methods(t *testing.T) {
	flag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			value: new(string),
		},
	}

	// 设置测试值
	testValue := "Hello World"
	if err := flag.Set(testValue); err != nil {
		t.Fatalf("设置值失败: %v", err)
	}

	// 测试Len()
	if flag.Len() != len(testValue) {
		t.Errorf("Len()期望%d, 实际%d", len(testValue), flag.Len())
	}

	// 测试ToUpper()
	if flag.ToUpper() != "HELLO WORLD" {
		t.Errorf("ToUpper()期望'HELLO WORLD', 实际'%s'", flag.ToUpper())
	}

	// 测试ToLower()
	if flag.ToLower() != "hello world" {
		t.Errorf("ToLower()期望'hello world', 实际'%s'", flag.ToLower())
	}

	// 测试Contains()
	if !flag.Contains("World") {
		t.Error("Contains('World')应返回true")
	}
	if flag.Contains("Go") {
		t.Error("Contains('Go')应返回false")
	}
}

// TestStringFlag_TypeAndString 测试类型和字符串表示
func TestStringFlag_TypeAndString(t *testing.T) {
	flag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			value: new(string),
		},
	}

	// 测试Type()
	if flag.Type() != FlagTypeString {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeString, flag.Type())
	}

	// 测试String()带引号格式
	if err := flag.Set("test"); err != nil {
		t.Fatalf("设置值失败: %v", err)
	}
	if flag.String() != "\"test\"" {
		t.Errorf("String()期望'\"test\"', 实际'%s'", flag.String())
	}
}
