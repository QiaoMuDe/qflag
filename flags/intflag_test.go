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
