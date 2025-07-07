package flags

import (
	"testing"
)

// TestEnumFlag_ValidInitialization 测试枚举标志的有效初始化
func TestEnumFlag_ValidInitialization(t *testing.T) {
	flag := &EnumFlag{}
	options := []string{"option1", "option2", "option3"}

	// 使用有效默认值初始化
	if err := flag.Init("enum", "e", "option1", "枚举测试", options); err != nil {
		t.Fatalf("初始化失败: %v", err)
	}

	// 验证默认值
	if flag.Get() != "option1" {
		t.Errorf("默认值应为'option1', 实际为'%s'", flag.Get())
	}
}

// TestEnumFlag_InvalidInitialization 测试枚举标志的无效初始化
func TestEnumFlag_InvalidInitialization(t *testing.T) {
	flag := &EnumFlag{}
	options := []string{"option1", "option2"}

	// 使用不在选项中的默认值
	if err := flag.Init("enum", "e", "invalid", "枚举测试", options); err == nil {
		t.Error("使用无效默认值时应返回错误")
	}
}

func TestEnumFlag_EmptyOptions(t *testing.T) {
	// 使用唯一名称避免冲突
	flag := &EnumFlag{}
	if err := flag.Init("enum_empty", "ee", "", "空选项枚举测试", []string{}); err != nil {
		t.Fatalf("空选项初始化失败: %v", err)
	}
	// 验证空选项时可以设置任意值
	if err := flag.Set("任意值"); err != nil {
		t.Errorf("空选项应允许任意值: %v", err)
	}
}

// TestEnumFlag_SetValidValues 测试设置有效枚举值
func TestEnumFlag_SetValidValues(t *testing.T) {
	flag := &EnumFlag{}
	options := []string{"apple", "banana", "cherry"}
	if err := flag.Init("fruit", "f", "apple", "水果枚举", options); err != nil {
		t.Fatalf("初始化失败: %v", err)
	}

	// 测试设置有效选项
	validValues := []string{"banana", "cherry", "apple"}
	for _, val := range validValues {
		if err := flag.Set(val); err != nil {
			t.Errorf("设置有效值'%s'失败: %v", val, err)
		}
		if flag.Get() != val {
			t.Errorf("设置后的值应为'%s', 实际为'%s'", val, flag.Get())
		}
	}
}

// TestEnumFlag_SetInvalidValues 测试设置无效枚举值
func TestEnumFlag_SetInvalidValues(t *testing.T) {
	flag := &EnumFlag{}
	options := []string{"red", "green", "blue"}
	if err := flag.Init("color", "c", "red", "颜色枚举", options); err != nil {
		t.Fatalf("初始化失败: %v", err)
	}

	// 测试设置无效选项
	invalidValues := []string{"yellow", "", "invalid"}
	for _, val := range invalidValues {
		if err := flag.Set(val); err == nil {
			t.Errorf("设置无效值'%s'应返回错误", val)
		}
	}
}

// TestEnumFlag_Type 验证Type()方法返回正确类型
func TestEnumFlag_Type(t *testing.T) {
	flag := &EnumFlag{}
	if flag.Type() != FlagTypeEnum {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeEnum, flag.Type())
	}
}

// TestEnumFlag_CaseInsensitive 测试不区分大小写模式下的枚举行为
func TestEnumFlag_CaseInsensitive(t *testing.T) {
	flag := &EnumFlag{}
	options := []string{"Apple", "Banana", "Cherry"}

	// 初始化枚举标志（默认不区分大小写）
	if err := flag.Init("fruit", "f", "Apple", "水果枚举测试", options); err != nil {
		t.Fatalf("初始化失败: %v", err)
	}

	// 测试不同大小写的有效值
	validInputs := []struct {
		input    string
		expected string
	}{{
		input:    "apple",
		expected: "apple",
	}, {
		input:    "BANANA",
		expected: "BANANA",
	}, {
		input:    "cHeRrY",
		expected: "cHeRrY",
	}}

	for _, test := range validInputs {
		t.Run(test.input, func(t *testing.T) {
			if err := flag.Set(test.input); err != nil {
				t.Errorf("设置值 '%s' 应该成功, 错误: %v", test.input, err)
			}
			if flag.Get() != test.expected {
				t.Errorf("获取值应为 '%s', 实际为 '%s'", test.expected, flag.Get())
			}
		})
	}
}

// TestEnumFlag_CaseSensitive 测试区分大小写模式下的枚举行为
func TestEnumFlag_CaseSensitive(t *testing.T) {
	flag := &EnumFlag{}
	options := []string{"Apple", "Banana", "Cherry"}

	// 初始化并设置为区分大小写
	if err := flag.Init("fruit", "f", "Apple", "水果枚举测试", options); err != nil {
		t.Fatalf("初始化失败: %v", err)
	}
	flag.SetCaseSensitive(true)

	// 测试大小写敏感的有效值
	validInputs := []string{"Apple", "Banana", "Cherry"}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			if err := flag.Set(input); err != nil {
				t.Errorf("设置值 '%s' 应该成功, 错误: %v", input, err)
			}
		})
	}

	// 测试大小写不匹配的无效值
	invalidInputs := []string{"apple", "BANANA", "cHeRrY", "grape"}
	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			if err := flag.Set(input); err == nil {
				t.Errorf("设置值 '%s' 应该失败", input)
			}
		})
	}
}
