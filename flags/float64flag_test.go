package flags

import (
	"fmt"
	"testing"
)

// TestFloatFlag_BasicFunctionality 测试FloatFlag的基本功能
func TestFloatFlag_BasicFunctionality(t *testing.T) {
	flag := &Float64Flag{
		BaseFlag: BaseFlag[float64]{
			initialValue: 0.0,
			value:        new(float64),
		},
	}

	// 测试默认值
	if flag.Get() != 0.0 {
		t.Error("默认值应为0.0")
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试设置有效值
	testCases := []float64{3.14, -2.5, 0.0, 100.0}
	for _, val := range testCases {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置值%.2f失败: %v", val, err)
		}
		if flag.Get() != val {
			t.Errorf("期望值%.2f, 实际值%.2f", val, flag.Get())
		}
	}

	// 测试重置功能
	flag.Reset()
	if flag.IsSet() {
		t.Error("重置后IsSet应返回false")
	}
	if flag.Get() != 0.0 {
		t.Error("重置后应返回默认值0.0")
	}
}

// TestFloatFlag_Type 验证Type()方法返回正确的标志类型
func TestFloatFlag_Type(t *testing.T) {
	flag := &Float64Flag{}
	if flag.Type() != FlagTypeFloat64 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeFloat64, flag.Type())
	}
}
