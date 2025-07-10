package flags

import (
	"fmt"
	"testing"
)

// TestIntFlag_BasicFunctionality 测试IntFlag的基本功能
func TestIntFlag_BasicFunctionality(t *testing.T) {
	flag := &IntFlag{
		BaseFlag: BaseFlag[int]{
			initialValue: 0,
			value:        new(int),
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
	testCases := []int{10, -5, 0, 1000}
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

// TestIntFlag_RangeValidation 测试整数范围验证
func TestIntFlag_RangeValidation(t *testing.T) {
	flag := &IntFlag{
		BaseFlag: BaseFlag[int]{
			initialValue: 5,
			value:        new(int),
		},
	}

	// 设置范围为1-10
	flag.SetRange(1, 10)

	// 测试有效范围内的值
	validValues := []int{1, 5, 10}
	for _, val := range validValues {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置有效值%d失败: %v", val, err)
		}
	}

	// 测试超出范围的值
	invalidValues := []int{0, 11, -5, 100}
	for _, val := range invalidValues {
		if err := flag.Set(fmt.Sprint(val)); err == nil {
			t.Errorf("设置无效值%d应返回错误", val)
		}
	}
}

// TestIntFlag_Type 验证Type()方法返回正确的标志类型
func TestIntFlag_Type(t *testing.T) {
	flag := &IntFlag{}
	if flag.Type() != FlagTypeInt {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeInt, flag.Type())
	}
}

// TestInt64Flag_BasicFunctionality 测试Int64Flag的基本功能
func TestInt64Flag_BasicFunctionality(t *testing.T) {
	flag := &Int64Flag{
		BaseFlag: BaseFlag[int64]{
			initialValue: 0,
			value:        new(int64),
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
	testCases := []int64{100, -50, 0, 9223372036854775807}
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

// TestInt64Flag_RangeValidation 测试64位整数范围验证
func TestInt64Flag_RangeValidation(t *testing.T) {
	flag := &Int64Flag{
		BaseFlag: BaseFlag[int64]{
			initialValue: 100,
			value:        new(int64),
		},
	}

	// 设置范围为-1000到1000
	flag.SetRange(-1000, 1000)

	// 测试有效范围内的值
	validValues := []int64{-1000, 0, 500, 1000}
	for _, val := range validValues {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置有效值%d失败: %v", val, err)
		}
	}

	// 测试超出范围的值
	invalidValues := []int64{-1001, 1001, -9223372036854775808, 9223372036854775807}
	for _, val := range invalidValues {
		if err := flag.Set(fmt.Sprint(val)); err == nil {
			t.Errorf("设置无效值%d应返回错误", val)
		}
	}
}

// TestInt64Flag_Type 验证Type()方法返回正确的标志类型
func TestInt64Flag_Type(t *testing.T) {
	flag := &Int64Flag{}
	if flag.Type() != FlagTypeInt64 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeInt64, flag.Type())
	}
}
