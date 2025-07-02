package flags

import (
	"fmt"
	"testing"
)

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
