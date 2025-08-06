package flags

import (
	"math"
	"strconv"
	"testing"
)

// TestIntFlag_EdgeCases 测试IntFlag的边界情况
func TestIntFlag_EdgeCases(t *testing.T) {
	t.Run("最大最小值", func(t *testing.T) {
		flag := &IntFlag{
			BaseFlag: BaseFlag[int]{
				initialValue: 0,
				value:        new(int),
			},
		}

		// 测试最大值
		err := flag.Set(strconv.Itoa(math.MaxInt32))
		if err != nil {
			t.Errorf("设置最大值失败: %v", err)
		}

		// 测试最小值
		err = flag.Set(strconv.Itoa(math.MinInt32))
		if err != nil {
			t.Errorf("设置最小值失败: %v", err)
		}
	})

	t.Run("无效输入", func(t *testing.T) {
		flag := &IntFlag{
			BaseFlag: BaseFlag[int]{
				initialValue: 0,
				value:        new(int),
			},
		}

		invalidInputs := []string{
			"abc",
			"12.34",
			"",
			"123abc",
			"∞",
		}

		for _, input := range invalidInputs {
			err := flag.Set(input)
			if err == nil {
				t.Errorf("无效输入 '%s' 应该返回错误", input)
			}
		}
	})

	t.Run("范围验证边界", func(t *testing.T) {
		flag := &IntFlag{
			BaseFlag: BaseFlag[int]{
				initialValue: 5,
				value:        new(int),
			},
		}

		// 设置范围 [1, 10]
		flag.SetRange(1, 10)

		// 测试边界值
		boundaryTests := []struct {
			value       string
			shouldError bool
			description string
		}{
			{"0", true, "小于最小值"},
			{"1", false, "等于最小值"},
			{"5", false, "中间值"},
			{"10", false, "等于最大值"},
			{"11", true, "大于最大值"},
		}

		for _, test := range boundaryTests {
			err := flag.Set(test.value)
			if test.shouldError && err == nil {
				t.Errorf("%s: 值 '%s' 应该返回错误", test.description, test.value)
			}
			if !test.shouldError && err != nil {
				t.Errorf("%s: 值 '%s' 不应该返回错误，但得到: %v", test.description, test.value, err)
			}
		}
	})
}

// TestInt64Flag_EdgeCases 测试Int64Flag的边界情况
func TestInt64Flag_EdgeCases(t *testing.T) {
	t.Run("极大值处理", func(t *testing.T) {
		flag := &Int64Flag{
			BaseFlag: BaseFlag[int64]{
				initialValue: 0,
				value:        new(int64),
			},
		}

		// 测试最大int64值
		maxInt64Str := strconv.FormatInt(math.MaxInt64, 10)
		err := flag.Set(maxInt64Str)
		if err != nil {
			t.Errorf("设置最大int64值失败: %v", err)
		}

		if flag.Get() != math.MaxInt64 {
			t.Errorf("期望 %d，实际 %d", math.MaxInt64, flag.Get())
		}

		// 测试最小int64值
		minInt64Str := strconv.FormatInt(math.MinInt64, 10)
		err = flag.Set(minInt64Str)
		if err != nil {
			t.Errorf("设置最小int64值失败: %v", err)
		}

		if flag.Get() != math.MinInt64 {
			t.Errorf("期望 %d，实际 %d", math.MinInt64, flag.Get())
		}
	})

	t.Run("超出范围的字符串", func(t *testing.T) {
		flag := &Int64Flag{
			BaseFlag: BaseFlag[int64]{
				initialValue: 0,
				value:        new(int64),
			},
		}

		// 测试超出int64范围的数字字符串
		oversizedInputs := []string{
			"9223372036854775808",  // MaxInt64 + 1
			"-9223372036854775809", // MinInt64 - 1
		}

		for _, input := range oversizedInputs {
			err := flag.Set(input)
			if err == nil {
				t.Errorf("超出范围的输入 '%s' 应该返回错误", input)
			}
		}
	})
}

// TestFloat64Flag_EdgeCases 测试Float64Flag的边界情况
func TestFloat64Flag_EdgeCases(t *testing.T) {
	t.Run("特殊浮点值", func(t *testing.T) {
		flag := &Float64Flag{
			BaseFlag: BaseFlag[float64]{
				initialValue: 0.0,
				value:        new(float64),
			},
		}

		specialValues := []struct {
			input    string
			expected float64
		}{
			{"0", 0.0},
			{"-0", math.Copysign(0, -1)}, // 使用math.Copysign创建负零
			{"3.14159", 3.14159},
			{"-3.14159", -3.14159},
			{"1e10", 1e10},
			{"1e-10", 1e-10},
		}

		for _, test := range specialValues {
			err := flag.Set(test.input)
			if err != nil {
				t.Errorf("设置浮点值 '%s' 失败: %v", test.input, err)
				continue
			}

			result := flag.Get()
			if result != test.expected {
				t.Errorf("输入 '%s'，期望 %f，实际 %f", test.input, test.expected, result)
			}
		}
	})

	t.Run("无效浮点输入", func(t *testing.T) {
		flag := &Float64Flag{
			BaseFlag: BaseFlag[float64]{
				initialValue: 0.0,
				value:        new(float64),
			},
		}

		invalidInputs := []string{
			"abc",
			"12.34.56",
			"",
			"12.34abc",
			"∞", // 注意：Go可以解析"Inf"但不能解析"∞"符号
		}

		for _, input := range invalidInputs {
			err := flag.Set(input)
			if err == nil {
				t.Errorf("无效浮点输入 '%s' 应该返回错误", input)
			}
		}
	})

	t.Run("特殊浮点值解析", func(t *testing.T) {
		flag := &Float64Flag{
			BaseFlag: BaseFlag[float64]{
				initialValue: 0.0,
				value:        new(float64),
			},
		}

		// Go可以成功解析这些特殊值
		validSpecialInputs := []string{
			"NaN",
			"Inf",
			"+Inf",
			"-Inf",
		}

		for _, input := range validSpecialInputs {
			err := flag.Set(input)
			if err != nil {
				t.Errorf("特殊浮点输入 '%s' 应该解析成功，但得到错误: %v", input, err)
			}
		}
	})
}

// TestBoolFlag_EdgeCases 测试BoolFlag的边界情况
func TestBoolFlag_EdgeCases(t *testing.T) {
	t.Run("各种布尔值表示", func(t *testing.T) {
		flag := &BoolFlag{
			BaseFlag: BaseFlag[bool]{
				initialValue: false,
				value:        new(bool),
			},
		}

		trueValues := []string{"true", "TRUE", "True", "1", "t", "T"}
		falseValues := []string{"false", "FALSE", "False", "0", "f", "F"}

		// 测试true值
		for _, val := range trueValues {
			err := flag.Set(val)
			if err != nil {
				t.Errorf("设置true值 '%s' 失败: %v", val, err)
				continue
			}
			if !flag.Get() {
				t.Errorf("输入 '%s' 应该解析为true", val)
			}
		}

		// 测试false值
		for _, val := range falseValues {
			err := flag.Set(val)
			if err != nil {
				t.Errorf("设置false值 '%s' 失败: %v", val, err)
				continue
			}
			if flag.Get() {
				t.Errorf("输入 '%s' 应该解析为false", val)
			}
		}
	})

	t.Run("无效布尔输入", func(t *testing.T) {
		flag := &BoolFlag{
			BaseFlag: BaseFlag[bool]{
				initialValue: false,
				value:        new(bool),
			},
		}

		invalidInputs := []string{
			"yes",
			"no",
			"on",
			"off",
			"",
			"maybe",
			"2",
		}

		for _, input := range invalidInputs {
			err := flag.Set(input)
			if err == nil {
				t.Errorf("无效布尔输入 '%s' 应该返回错误", input)
			}
		}
	})

	t.Run("IsBoolFlag接口", func(t *testing.T) {
		flag := &BoolFlag{}
		if !flag.IsBoolFlag() {
			t.Error("BoolFlag应该实现IsBoolFlag接口并返回true")
		}
	})
}

// TestStringFlag_EdgeCases 测试StringFlag的边界情况
func TestStringFlag_EdgeCases(t *testing.T) {
	t.Run("特殊字符串", func(t *testing.T) {
		flag := &StringFlag{
			BaseFlag: BaseFlag[string]{
				initialValue: "",
				value:        new(string),
			},
		}

		specialStrings := []string{
			"",             // 空字符串
			" ",            // 空格
			"\n",           // 换行符
			"\t",           // 制表符
			"中文字符串",        // 中文
			"🚀🎉",           // emoji
			"\"quoted\"",   // 带引号
			"line1\nline2", // 多行
			"very long string " + string(make([]byte, 1000)), // 长字符串
		}

		for _, str := range specialStrings {
			err := flag.Set(str)
			if err != nil {
				t.Errorf("设置字符串失败: %v", err)
				continue
			}

			if flag.Get() != str {
				t.Errorf("字符串不匹配，期望 '%s'，实际 '%s'", str, flag.Get())
			}
		}
	})

	t.Run("字符串方法测试", func(t *testing.T) {
		flag := &StringFlag{
			BaseFlag: BaseFlag[string]{
				initialValue: "",
				value:        new(string),
			},
		}

		testString := "Hello, 世界! 🌍"
		err := flag.Set(testString)
		if err != nil {
			t.Fatalf("设置标志失败: %v", err)
		}

		// 测试Len方法
		if flag.Len() != len(testString) {
			t.Errorf("Len()期望 %d，实际 %d", len(testString), flag.Len())
		}

		// 测试ToUpper方法
		expectedUpper := "HELLO, 世界! 🌍"
		if flag.ToUpper() != expectedUpper {
			t.Errorf("ToUpper()期望 '%s'，实际 '%s'", expectedUpper, flag.ToUpper())
		}

		// 测试ToLower方法
		expectedLower := "hello, 世界! 🌍"
		if flag.ToLower() != expectedLower {
			t.Errorf("ToLower()期望 '%s'，实际 '%s'", expectedLower, flag.ToLower())
		}

		// 测试Contains方法
		if !flag.Contains("世界") {
			t.Error("Contains('世界')应该返回true")
		}
		if flag.Contains("不存在") {
			t.Error("Contains('不存在')应该返回false")
		}
	})

	t.Run("String方法带引号", func(t *testing.T) {
		flag := &StringFlag{
			BaseFlag: BaseFlag[string]{
				initialValue: "",
				value:        new(string),
			},
		}

		testCases := []struct {
			input    string
			expected string
		}{
			{"hello", "\"hello\""},
			{"", "\"\""},
			{"with\"quotes", "\"with\\\"quotes\""},
			{"line1\nline2", "\"line1\\nline2\""},
		}

		for _, test := range testCases {
			_ = flag.Set(test.input)
			result := flag.String()
			if result != test.expected {
				t.Errorf("输入 '%s'，String()期望 '%s'，实际 '%s'", test.input, test.expected, result)
			}
		}
	})
}

// TestUintFlags_EdgeCases 测试无符号整数标志的边界情况
func TestUintFlags_EdgeCases(t *testing.T) {
	t.Run("Uint16Flag边界值", func(t *testing.T) {
		flag := &Uint16Flag{
			BaseFlag: BaseFlag[uint16]{
				initialValue: 0,
				value:        new(uint16),
			},
		}

		// 测试有效范围
		validValues := []string{"0", "32767", "65535"}
		for _, val := range validValues {
			err := flag.Set(val)
			if err != nil {
				t.Errorf("设置有效uint16值 '%s' 失败: %v", val, err)
			}
		}

		// 测试无效值
		invalidValues := []string{"-1", "65536", "abc"}
		for _, val := range invalidValues {
			err := flag.Set(val)
			if err == nil {
				t.Errorf("无效uint16值 '%s' 应该返回错误", val)
			}
		}
	})

	t.Run("Uint32Flag边界值", func(t *testing.T) {
		flag := &Uint32Flag{
			BaseFlag: BaseFlag[uint32]{
				initialValue: 0,
				value:        new(uint32),
			},
		}

		// 测试最大值
		err := flag.Set("4294967295")
		if err != nil {
			t.Errorf("设置uint32最大值失败: %v", err)
		}

		// 测试超出范围
		err = flag.Set("4294967296")
		if err == nil {
			t.Error("超出uint32范围的值应该返回错误")
		}
	})

	t.Run("Uint64Flag边界值", func(t *testing.T) {
		flag := &Uint64Flag{
			BaseFlag: BaseFlag[uint64]{
				initialValue: 0,
				value:        new(uint64),
			},
		}

		// 测试最大值
		maxUint64Str := strconv.FormatUint(math.MaxUint64, 10)
		err := flag.Set(maxUint64Str)
		if err != nil {
			t.Errorf("设置uint64最大值失败: %v", err)
		}

		if flag.Get() != math.MaxUint64 {
			t.Errorf("期望 %d，实际 %d", uint64(math.MaxUint64), flag.Get())
		}
	})
}

// TestAllFlags_StringRepresentation 测试所有标志类型的字符串表示
func TestAllFlags_StringRepresentation(t *testing.T) {
	testCases := []struct {
		name     string
		flag     Flag
		setValue func(Flag) error
		expected string
	}{
		{
			name: "IntFlag",
			flag: &IntFlag{BaseFlag: BaseFlag[int]{value: new(int)}},
			setValue: func(f Flag) error {
				return f.(*IntFlag).Set("42")
			},
			expected: "42",
		},
		{
			name: "StringFlag",
			flag: &StringFlag{BaseFlag: BaseFlag[string]{value: new(string)}},
			setValue: func(f Flag) error {
				return f.(*StringFlag).Set("hello")
			},
			expected: "\"hello\"",
		},
		{
			name: "BoolFlag",
			flag: &BoolFlag{BaseFlag: BaseFlag[bool]{value: new(bool)}},
			setValue: func(f Flag) error {
				return f.(*BoolFlag).Set("true")
			},
			expected: "true",
		},
		{
			name: "Float64Flag",
			flag: &Float64Flag{BaseFlag: BaseFlag[float64]{value: new(float64)}},
			setValue: func(f Flag) error {
				return f.(*Float64Flag).Set("3.14")
			},
			expected: "3.14",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := test.setValue(test.flag)
			if err != nil {
				t.Fatalf("设置值失败: %v", err)
			}

			result := test.flag.String()
			if result != test.expected {
				t.Errorf("String()期望 '%s'，实际 '%s'", test.expected, result)
			}
		})
	}
}

// TestAllFlags_TypeMethods 测试所有标志类型的Type方法
func TestAllFlags_TypeMethods(t *testing.T) {
	testCases := []struct {
		name         string
		flag         Flag
		expectedType FlagType
	}{
		{"IntFlag", &IntFlag{}, FlagTypeInt},
		{"Int64Flag", &Int64Flag{}, FlagTypeInt64},
		{"Uint16Flag", &Uint16Flag{}, FlagTypeUint16},
		{"Uint32Flag", &Uint32Flag{}, FlagTypeUint32},
		{"Uint64Flag", &Uint64Flag{}, FlagTypeUint64},
		{"StringFlag", &StringFlag{}, FlagTypeString},
		{"BoolFlag", &BoolFlag{}, FlagTypeBool},
		{"Float64Flag", &Float64Flag{}, FlagTypeFloat64},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result := test.flag.Type()
			if result != test.expectedType {
				t.Errorf("%s.Type()期望 %d，实际 %d", test.name, test.expectedType, result)
			}
		})
	}
}
