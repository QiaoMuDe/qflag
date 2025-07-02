package flags

import (
	"testing"
)

// TestUint16Flag_ValidParsing 测试有效的uint16值解析
func TestUint16Flag_ValidParsing(t *testing.T) {
	flag := &Uint16Flag{
		BaseFlag: BaseFlag[uint16]{
			initialValue: 0,
			value:        new(uint16),
		},
	}

	// 测试用例集合
	testCases := []struct {
		input    string
		expected uint16
	}{{
		input:    "0",
		expected: 0,
	}, {
		input:    "65535",
		expected: 65535,
	}, {
		input:    "32768",
		expected: 32768,
	}, {
		input:    "1000",
		expected: 1000,
	}}

	for _, tc := range testCases {
		if err := flag.Set(tc.input); err != nil {
			t.Errorf("解析 %s 失败: %v", tc.input, err)
			continue
		}
		if flag.Get() != tc.expected {
			t.Errorf("%s 期望 %d, 实际 %d", tc.input, tc.expected, flag.Get())
		}
	}
}

// TestUint16Flag_InvalidParsing 测试无效的uint16值解析
func TestUint16Flag_InvalidParsing(t *testing.T) {
	flag := &Uint16Flag{
		BaseFlag: BaseFlag[uint16]{
			initialValue: 0,
			value:        new(uint16),
		},
	}

	// 测试用例集合
	invalidInputs := []string{
		"-1",    // 负值
		"65536", // 超出最大值
		"abc",   // 非数字
		"12.34", // 浮点数
		" 123 ", // 带空格
	}

	for _, input := range invalidInputs {
		if err := flag.Set(input); err == nil {
			t.Errorf("解析无效值 '%s' 应返回错误", input)
		}
	}
}

// TestUint16Flag_TypeAndString 测试类型和字符串表示
func TestUint16Flag_TypeAndString(t *testing.T) {
	flag := &Uint16Flag{}

	// 测试Type()
	if flag.Type() != FlagTypeUint16 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeUint16, flag.Type())
	}

	// 测试String()
	if err := flag.Set("1234"); err != nil {
		t.Fatalf("设置值失败: %v", err)
	}
	if flag.String() != "1234" {
		t.Errorf("String()期望'1234', 实际'%s'", flag.String())
	}
}
