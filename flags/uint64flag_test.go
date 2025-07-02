package flags

import (
	"fmt"
	"testing"
)

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
