package flags

import "testing"

// TestBoolFlag_基本功能测试 验证BoolFlag的设置值、默认值和重置功能
func TestBoolFlag_BasicFunctionality(t *testing.T) {
	// 创建BoolFlag实例
	flag := &BoolFlag{
		BaseFlag: BaseFlag[bool]{
			initialValue: false,
			value:        new(bool),
		},
	}

	// 测试1: 验证默认值
	if flag.GetDefault() != false {
		t.Error("默认值应为false")
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试2: 设置并验证true值
	if err := flag.Set("true"); err != nil {
		t.Fatalf("设置true值失败: %v", err)
	}
	if flag.Get() != true {
		t.Error("设置后的值应为true")
	}
	if !flag.IsSet() {
		t.Error("设置值后IsSet应返回true")
	}

	// 测试3: 设置并验证false值
	if err := flag.Set("false"); err != nil {
		t.Fatalf("设置false值失败: %v", err)
	}
	if flag.Get() != false {
		t.Error("设置后的值应为false")
	}

	// 测试4: 重置功能
	flag.Reset()
	if flag.IsSet() {
		t.Error("重置后IsSet应返回false")
	}
	if flag.Get() != false {
		t.Error("重置后应返回默认值false")
	}
}

// TestBoolFlag_Type 验证Type()方法返回正确的标志类型
func TestBoolFlag_Type(t *testing.T) {
	flag := &BoolFlag{}
	if flag.Type() != FlagTypeBool {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeBool, flag.Type())
	}
}
