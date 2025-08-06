package flags

import (
	"math"
	"testing"
)

// TestMapFlag_EdgeCases_Fixed 修复MapFlag边界情况测试
func TestMapFlag_EdgeCases_Fixed(t *testing.T) {
	t.Run("键值包含分隔符", func(t *testing.T) {
		flag := &MapFlag{}
		flag.SetDelimiters(",", "=")

		// 测试键包含分隔符的情况 - 这应该失败
		err := flag.Set("key=with=equals=value")
		if err == nil {
			t.Error("包含多个分隔符的键值对应该返回错误")
		}

		// 测试正确的转义或引用方式
		err = flag.Set("key1=value1,key2=value with spaces")
		if err != nil {
			t.Errorf("正常键值对设置失败: %v", err)
		}
	})

	t.Run("空键值处理", func(t *testing.T) {
		flag := &MapFlag{}
		flag.SetDelimiters(",", "=")

		// 测试空键
		err := flag.Set("=value")
		if err == nil {
			t.Error("空键应该返回错误")
		}

		// 测试空值（这应该是允许的）
		err = flag.Set("key=")
		if err != nil {
			t.Errorf("空值应该被允许: %v", err)
		}

		result := flag.Get()
		if result["key"] != "" {
			t.Errorf("空值应该为空字符串，实际为 '%s'", result["key"])
		}
	})

	t.Run("默认分隔符处理", func(t *testing.T) {
		flag := &MapFlag{}
		// 不设置分隔符，使用默认值

		err := flag.Set("key1=value1,key2=value2")
		if err != nil {
			t.Errorf("使用默认分隔符失败: %v", err)
		}

		result := flag.Get()
		if len(result) != 2 {
			t.Errorf("应该有2个键值对，实际有 %d 个", len(result))
		}
	})
}

// TestBaseFlag_PointerAccess_Fixed 修复BaseFlag指针访问测试
func TestBaseFlag_PointerAccess_Fixed(t *testing.T) {
	t.Run("未设置值时指针行为", func(t *testing.T) {
		flag := &BaseFlag[int]{
			initialValue: 42,
			value:        new(int),
		}

		// 初始化标志
		flag.Init("test", "", "test flag", nil)

		// 获取指针 - 应该返回指向初始值的指针
		ptr := flag.GetPointer()
		if ptr == nil {
			t.Error("GetPointer()不应该返回nil")
			return
		}

		// 检查指针指向的值
		if *ptr != 42 {
			t.Errorf("指针应该指向初始值42，实际为 %d", *ptr)
		}
	})

	t.Run("设置值后指针有效", func(t *testing.T) {
		flag := &BaseFlag[int]{
			initialValue: 0,
			value:        new(int),
		}

		flag.Init("test", "", "test flag", nil)

		// 设置值
		*flag.value = 100
		flag.isSet = true

		ptr := flag.GetPointer()
		if ptr == nil {
			t.Error("设置值后GetPointer()不应该返回nil")
			return
		}

		if *ptr != 100 {
			t.Errorf("指针应该指向设置的值100，实际为 %d", *ptr)
		}
	})
}

// TestFloat64Flag_EdgeCases_Fixed 修复Float64Flag边界情况测试
func TestFloat64Flag_EdgeCases_Fixed(t *testing.T) {
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
			{"-0", math.Copysign(0, -1)},
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
			"∞",
			// 注意：Go的strconv.ParseFloat实际上可以解析"NaN"，所以我们移除这个测试
		}

		for _, input := range invalidInputs {
			err := flag.Set(input)
			if err == nil {
				t.Errorf("无效浮点输入 '%s' 应该返回错误", input)
			}
		}
	})

	t.Run("特殊IEEE754值", func(t *testing.T) {
		flag := &Float64Flag{
			BaseFlag: BaseFlag[float64]{
				initialValue: 0.0,
				value:        new(float64),
			},
		}

		// 测试Go能够解析的特殊值
		specialCases := []struct {
			input   string
			checkFn func(float64) bool
			desc    string
		}{
			{"NaN", math.IsNaN, "NaN值"},
			{"Inf", func(f float64) bool { return math.IsInf(f, 0) }, "无穷"},
			{"+Inf", func(f float64) bool { return math.IsInf(f, 1) }, "正无穷"},
			{"-Inf", func(f float64) bool { return math.IsInf(f, -1) }, "负无穷"},
		}

		for _, test := range specialCases {
			err := flag.Set(test.input)
			if err != nil {
				t.Errorf("设置%s失败: %v", test.desc, err)
				continue
			}

			result := flag.Get()
			if !test.checkFn(result) {
				t.Errorf("输入 '%s' 应该是%s", test.input, test.desc)
			}
		}
	})
}

// TestAllFlags_Comprehensive 综合测试所有标志类型
func TestAllFlags_Comprehensive(t *testing.T) {
	t.Run("所有标志类型的基本功能", func(t *testing.T) {
		testCases := []struct {
			name string
			flag interface {
				Set(string) error
				String() string
				Type() FlagType
			}
			setValue string
			checkFn  func(interface{}) bool
		}{
			{
				name:     "IntFlag",
				flag:     &IntFlag{BaseFlag: BaseFlag[int]{value: new(int)}},
				setValue: "42",
				checkFn: func(f interface{}) bool {
					return f.(*IntFlag).Get() == 42
				},
			},
			{
				name:     "StringFlag",
				flag:     &StringFlag{BaseFlag: BaseFlag[string]{value: new(string)}},
				setValue: "hello",
				checkFn: func(f interface{}) bool {
					return f.(*StringFlag).Get() == "hello"
				},
			},
			{
				name:     "BoolFlag",
				flag:     &BoolFlag{BaseFlag: BaseFlag[bool]{value: new(bool)}},
				setValue: "true",
				checkFn: func(f interface{}) bool {
					return f.(*BoolFlag).Get()
				},
			},
			{
				name:     "Float64Flag",
				flag:     &Float64Flag{BaseFlag: BaseFlag[float64]{value: new(float64)}},
				setValue: "3.14",
				checkFn: func(f interface{}) bool {
					return f.(*Float64Flag).Get() == 3.14
				},
			},
		}

		for _, test := range testCases {
			t.Run(test.name, func(t *testing.T) {
				err := test.flag.Set(test.setValue)
				if err != nil {
					t.Fatalf("设置值失败: %v", err)
				}

				if !test.checkFn(test.flag) {
					t.Errorf("%s 值检查失败", test.name)
				}

				// 测试String方法
				str := test.flag.String()
				if str == "" {
					t.Errorf("%s String()方法返回空字符串", test.name)
				}

				// 测试Type方法
				flagType := test.flag.Type()
				if flagType == 0 {
					t.Errorf("%s Type()方法返回无效类型", test.name)
				}
			})
		}
	})
}
