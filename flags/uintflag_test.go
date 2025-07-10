package flags

import (
	"fmt"
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

// TestUint32Flag_BasicFunctionality 测试Uint32Flag的基本功能
func TestUint32Flag_BasicFunctionality(t *testing.T) {
	flag := &Uint32Flag{
		BaseFlag: BaseFlag[uint32]{
			initialValue: 0,
			value:        new(uint32),
		},
	}

	// 测试默认值
	if flag.Get() != 0 {
		t.Error("默认值应为0")
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试设置有效值
	testCases := []uint32{100, 4294967295, 0, 2147483648}
	for _, val := range testCases {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置值%d失败: %v", val, err)
		}
		if flag.Get() != val {
			t.Errorf("期望值%d, 实际值%d", val, flag.Get())
		}
	}

	// 测试重置功能
	flag.Reset()
	if flag.IsSet() {
		t.Error("重置后IsSet应返回false")
	}
	if flag.Get() != 0 {
		t.Error("重置后应返回默认值0")
	}
}

// TestUint32Flag_InvalidValue 测试设置无效值
func TestUint32Flag_InvalidValue(t *testing.T) {
	flag := &Uint32Flag{
		BaseFlag: BaseFlag[uint32]{
			value: new(uint32),
		},
	}

	invalidValues := []string{"4294967296", "-1", "abc"}
	for _, val := range invalidValues {
		if err := flag.Set(val); err == nil {
			t.Errorf("设置无效值%s应返回错误", val)
		}
	}
}

// TestUint32Flag_Type 验证Type()方法返回正确的标志类型
func TestUint32Flag_Type(t *testing.T) {
	flag := &Uint32Flag{}
	if flag.Type() != FlagTypeUint32 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeUint32, flag.Type())
	}
}

// TestUint64Flag_BasicFunctionality 测试Uint64Flag的基本功能
func TestUint64Flag_BasicFunctionality(t *testing.T) {
	flag := &Uint64Flag{
		BaseFlag: BaseFlag[uint64]{
			initialValue: 0,
			value:        new(uint64),
		},
	}

	// 测试默认值
	if flag.Get() != 0 {
		t.Error("默认值应为0")
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试设置有效值
	testCases := []uint64{100, 18446744073709551615, 0, 9223372036854775808}
	for _, val := range testCases {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置值%d失败: %v", val, err)
		}
		if flag.Get() != val {
			t.Errorf("期望值%d, 实际值%d", val, flag.Get())
		}
	}

	// 测试重置功能
	flag.Reset()
	if flag.IsSet() {
		t.Error("重置后IsSet应返回false")
	}
	if flag.Get() != 0 {
		t.Error("重置后应返回默认值0")
	}
}

// TestUint64Flag_InvalidValue 测试设置无效值
func TestUint64Flag_InvalidValue(t *testing.T) {
	flag := &Uint64Flag{
		BaseFlag: BaseFlag[uint64]{
			value: new(uint64),
		},
	}

	invalidValues := []string{"18446744073709551616", "-1", "not_a_number"}
	for _, val := range invalidValues {
		if err := flag.Set(val); err == nil {
			t.Errorf("设置无效值%s应返回错误", val)
		}
	}
}

// TestUint64Flag_Type 验证Type()方法返回正确的标志类型
func TestUint64Flag_Type(t *testing.T) {
	flag := &Uint64Flag{}
	if flag.Type() != FlagTypeUint64 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeUint64, flag.Type())
	}
}
