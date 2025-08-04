package flags

import (
	"strings"
	"testing"
	"time"
)

// TestStringFlag_BasicFunctionality 测试StringFlag的基本功能
func TestStringFlag_BasicFunctionality(t *testing.T) {
	flag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			initialValue: "default",
			value:        new(string),
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

// TestIsSetMethods 测试所有标志类型的IsSet()方法行为
func TestIsSetMethods(t *testing.T) {
	// 测试用例结构体：包含标志实例、设置值函数和测试名称
	type testCase struct {
		name     string
		f        Flag
		setValue func(f Flag) error
	}

	// 创建测试用例集合
	testCases := []testCase{
		// IntFlag测试用例
		{
			name: "IntFlag未设置值",
			f: &IntFlag{
				BaseFlag: BaseFlag[int]{
					longName:     "intFlag",
					shortName:    "i",
					initialValue: 0,
					usage:        "整数标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "IntFlag已设置值",
			f: &IntFlag{
				BaseFlag: BaseFlag[int]{
					longName:     "intFlag",
					shortName:    "i",
					initialValue: 0,
					usage:        "整数标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*IntFlag).Set("100") },
		},
		{
			name: "IntFlag重置后",
			f: &IntFlag{
				BaseFlag: BaseFlag[int]{
					longName:     "intFlag",
					shortName:    "i",
					initialValue: 0,
					usage:        "整数标志测试",
				},
			},
			setValue: func(f Flag) error {
				if err := f.(*IntFlag).Set("100"); err != nil {
					return err
				}
				f.Reset()
				return nil
			},
		},

		// StringFlag测试用例
		{
			name: "StringFlag未设置值",
			f: &StringFlag{
				BaseFlag: BaseFlag[string]{
					longName:     "strFlag",
					shortName:    "s",
					initialValue: "default",
					usage:        "字符串标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "StringFlag已设置值",
			f: &StringFlag{
				BaseFlag: BaseFlag[string]{
					longName:     "strFlag",
					shortName:    "s",
					initialValue: "default",
					usage:        "字符串标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*StringFlag).Set("test") },
		},

		// BoolFlag测试用例
		{
			name: "BoolFlag未设置值",
			f: &BoolFlag{
				BaseFlag: BaseFlag[bool]{
					longName:     "boolFlag",
					shortName:    "b",
					initialValue: false,
					usage:        "布尔标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "BoolFlag已设置值",
			f: &BoolFlag{
				BaseFlag: BaseFlag[bool]{
					longName:     "boolFlag",
					shortName:    "b",
					initialValue: false,
					usage:        "布尔标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*BoolFlag).Set("true") },
		},

		// FloatFlag测试用例
		{
			name: "FloatFlag未设置值",
			f: &Float64Flag{
				BaseFlag: BaseFlag[float64]{
					longName:     "floatFlag",
					shortName:    "f",
					initialValue: 0.0,
					usage:        "浮点标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "FloatFlag已设置值",
			f: &Float64Flag{
				BaseFlag: BaseFlag[float64]{
					longName:     "floatFlag",
					shortName:    "f",
					initialValue: 0.0,
					usage:        "浮点标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*Float64Flag).Set("3.14") },
		},

		// DurationFlag测试用例
		{
			name: "DurationFlag未设置值",
			f: &DurationFlag{
				BaseFlag: BaseFlag[time.Duration]{
					longName:     "durationFlag",
					shortName:    "d",
					initialValue: 0,
					usage:        "时间间隔标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "DurationFlag已设置值",
			f: &DurationFlag{
				BaseFlag: BaseFlag[time.Duration]{
					longName:     "durationFlag",
					shortName:    "d",
					initialValue: 0,
					usage:        "时间间隔标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*DurationFlag).Set((5 * time.Second).String()) },
		},

		// EnumFlag测试用例
		{
			name: "EnumFlag未设置值",
			f: &EnumFlag{
				BaseFlag: BaseFlag[string]{
					longName:     "enumFlag",
					shortName:    "e",
					initialValue: "default",
					usage:        "枚举标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "EnumFlag已设置值",
			f: &EnumFlag{
				BaseFlag: BaseFlag[string]{
					longName:     "enumFlag",
					shortName:    "e",
					initialValue: "default",
					usage:        "枚举标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*EnumFlag).Set("option1") },
		},

		// SliceFlag测试用例
		{
			name: "SliceFlag未设置值",
			f: &SliceFlag{
				BaseFlag: BaseFlag[[]string]{
					longName:     "sliceFlag",
					shortName:    "sl",
					initialValue: []string{"default"},
					usage:        "切片标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "SliceFlag已设置值",
			f: &SliceFlag{
				BaseFlag: BaseFlag[[]string]{
					longName:     "sliceFlag",
					shortName:    "sl",
					initialValue: []string{"default"},
					usage:        "切片标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*SliceFlag).Set("item1,item2") },
		},
	}

	// 执行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始状态检查
			if tc.f.IsSet() {
				t.Errorf("%s: 初始状态下IsSet()应为false, 实际为%v", tc.name, tc.f.IsSet())
			}

			// 设置值
			if err := tc.setValue(tc.f); err != nil {
				t.Fatalf("%s: 设置值失败: %v", tc.name, err)
			}

			// 根据测试类型判断预期结果
			shouldBeSet := !strings.Contains(tc.name, "未设置值") && !strings.Contains(tc.name, "重置后")

			// 检查设置后状态
			if tc.f.IsSet() != shouldBeSet {
				// 修复重置后状态的预期值
				if strings.Contains(tc.name, "重置后") {
					shouldBeSet = false
				}
				t.Errorf("%s: 设置后IsSet()应为%v, 实际为%v", tc.name, shouldBeSet, tc.f.IsSet())
			}

			// 重置标志
			tc.f.Reset()

			// 检查重置后状态
			if tc.f.IsSet() {
				t.Errorf("%s: 重置后IsSet()应为false, 实际为true", tc.name)
			}
		})
	}
}
