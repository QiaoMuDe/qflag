package qflag

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

// positiveIntValidator 验证整数是否为正数（首字母小写, 非导出）
type positiveIntValidator struct{}

// Validate 实现Validator接口, 检查值是否为正数
func (v *positiveIntValidator) Validate(value any) error {
	val, ok := value.(int)
	if !ok {
		return errors.New("value must be an integer")
	}
	if val <= 0 {
		return errors.New("value must be positive")
	}
	return nil
}

// stringLengthValidator 验证字符串长度是否在指定范围内（首字母小写, 非导出）
type stringLengthValidator struct {
	min, max int
}

// Validate 实现Validator接口, 检查字符串长度
func (v *stringLengthValidator) Validate(value any) error {
	val, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	if len(val) < v.min || len(val) > v.max {
		return fmt.Errorf("string length must be between %d and %d", v.min, v.max)
	}
	return nil
}

// TestIntFlag_Validator 测试IntFlag的验证器功能
func TestIntFlag_Validator(t *testing.T) {
	// 创建整数标志
	flag := &IntFlag{
		BaseFlag: BaseFlag[int]{
			defValue: 0,
			value:    new(int),
		},
	}

	// 设置正整数验证器
	flag.SetValidator(&positiveIntValidator{})

	// 测试用例：有效正值
	if err := flag.Set(100); err != nil {
		t.Errorf("expected no error for valid positive value, got %v", err)
	}

	// 测试用例：无效负值
	if err := flag.Set(-5); err == nil {
		t.Error("expected error for negative value, got nil")
	} else if err.Error() != "invalid value for : value must be positive" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestStringFlag_Validator 测试StringFlag的验证器功能
func TestStringFlag_Validator(t *testing.T) {
	// 创建字符串标志
	flag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			defValue: "",
			value:    new(string),
		},
	}

	// 设置字符串长度验证器（2-10个字符）
	flag.SetValidator(&stringLengthValidator{min: 2, max: 10})

	// 测试用例：有效长度
	validStr := "test"
	if err := flag.Set(validStr); err != nil {
		t.Errorf("expected no error for valid string length, got %v", err)
	}

	// 测试用例：太短的字符串
	shortStr := "a"
	if err := flag.Set(shortStr); err == nil {
		t.Error("expected error for too short string, got nil")
	}

	// 测试用例：太长的字符串
	longStr := "thisisaverylongstring"
	if err := flag.Set(longStr); err == nil {
		t.Error("expected error for too long string, got nil")
	}
}

// TestBaseFlag_GetPointer 验证GetPointer()方法的基本功能和指针访问有效性
func TestBaseFlag_GetPointer(t *testing.T) {
	// 1. 测试整数类型标志的指针行为
	intFlag := &IntFlag{
		BaseFlag: BaseFlag[int]{
			defValue: 10,
			value:    nil,
		},
	}

	// 未设置值时指针应为nil
	if ptr := intFlag.GetPointer(); ptr != nil {
		t.Error("IntFlag未设置值时, GetPointer()应返回nil")
	}

	// 设置值后验证指针有效性
	if err := intFlag.Set(20); err != nil {
		t.Fatalf("设置IntFlag值失败: %v", err)
	}

	ptr := intFlag.GetPointer()
	if ptr == nil {
		t.Fatal("IntFlag设置值后, GetPointer()不应返回nil")
	}

	if *ptr != 20 {
		t.Errorf("IntFlag指针值错误, 期望20, 实际%d", *ptr)
	}

	// 通过指针修改值并验证
	*ptr = 30
	if intFlag.Get() != 30 {
		t.Errorf("通过指针修改值失败, 期望30, 实际%d", intFlag.Get())
	}

	// 2. 测试字符串类型标志的指针行为
	strFlag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			defValue: "default",
		},
	}

	if err := strFlag.Set("test"); err != nil {
		t.Fatalf("设置StringFlag值失败: %v", err)
	}

	*strFlag.GetPointer() = "modified"
	if strFlag.Get() != "modified" {
		t.Errorf("StringFlag指针修改失败, 期望'modified', 实际'%s'", strFlag.Get())
	}

	// 3. 测试默认值场景（值未显式设置时）
	defaultFlag := &BoolFlag{
		BaseFlag: BaseFlag[bool]{
			defValue: true,
			value:    nil,
		},
	}

	// 未设置值时指针应为nil, Get()应返回默认值
	if ptr := defaultFlag.GetPointer(); ptr != nil {
		t.Error("BoolFlag未设置值时, GetPointer()应返回nil")
	}
	if defaultFlag.Get() != true {
		t.Error("BoolFlag未设置值时, Get()应返回默认值true")
	}
}

func TestCommandAndFlagRegistration(t *testing.T) {
	// 测试用例1: 只有长名称的命令
	t.Run("Command with long name only", func(t *testing.T) {
		cmd := NewCmd("longcmd", "", flag.ContinueOnError)
		cmd.SetDescription("This command has only long name")

		// 添加只有长名称的标志
		cmd.String("aconfig", "", "Config file path", "/etc/app.conf")

		// 验证帮助信息生成
		help := cmd.GetHelp()
		if !strings.Contains(help, "longcmd") {
			t.Error("Command long name not found in help")
		}
		if strings.Contains(help, "-c") {
			t.Error("Unexpected short name in help")
		}
	})

	// 测试用例2: 只有短名称的命令
	t.Run("Command with short name only", func(t *testing.T) {
		cmd := NewCmd("", "s", flag.ContinueOnError)
		cmd.SetDescription("This command has only short name")

		// 跳过内置标志绑定，避免冲突
		cmd.initFlagBound = true

		// 添加只有短名称的标志
		cmd.String("", "c", "Config file path", "/etc/app.conf")

		// 验证帮助信息生成
		help := cmd.GetHelp()
		if !strings.Contains(help, "-c") {
			t.Error("Command short name not found in help")
		}
		if strings.Contains(help, "--config") {
			t.Error("Unexpected long name in help")
		}
	})

	// 测试用例3: 混合名称的命令和标志
	t.Run("Mixed name command and flags", func(t *testing.T) {
		cmd := NewCmd("mixed", "m", flag.ContinueOnError)
		cmd.SetDescription("This command has both long and short names")

		// 添加各种组合的标志
		cmd.String("config", "c", "Config file path", "/etc/app.conf")
		cmd.String("output", "", "Output directory", "./out")
		cmd.String("", "v", "Verbose mode", "false")

		// 验证帮助信息生成
		help := cmd.GetHelp()

		// 检查命令名称显示
		if !strings.Contains(help, "mixed, m") {
			t.Error("Command name display incorrect")
		}

		// 检查标志显示
		if !strings.Contains(help, "--config, -c") {
			t.Error("Flag with both names display incorrect")
		}
		if !strings.Contains(help, "--output") {
			t.Error("Flag with long name only display incorrect")
		}
		if !strings.Contains(help, "-v") {
			t.Error("Flag with short name only display incorrect")
		}
	})
}

// TestSliceFlag 测试SliceFlag的功能
func TestSliceFlag(t *testing.T) {
	// 测试基本切片解析功能
	t.Run("BasicSliceParsing", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				defValue: []string{},
				value:    new([]string),
			},
			delimiters: []string{","},
		}

		// 测试正常分割
		if err := flag.Set("a,b,c"); err != nil {
			t.Errorf("Set failed: %v", err)
		}
		result := flag.Get()
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// 测试无分隔符情况
		if err := flag.Set("d"); err != nil {
			t.Errorf("Set failed: %v", err)
		}
		result = flag.Get()
		expected = []string{"a", "b", "c", "d"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// 测试空元素过滤功能
	t.Run("EmptyElementFiltering", func(t *testing.T) {
		// 测试SkipEmpty=true情况
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				defValue: []string{},
				value:    new([]string),
			},
			delimiters: []string{","},
			SkipEmpty:  true,
		}

		if err := flag.Set("a,,b,,c"); err != nil {
			t.Errorf("Set failed: %v", err)
		}
		result := flag.Get()
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// 测试SkipEmpty=false情况
		flag = &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				defValue: []string{},
				value:    new([]string),
			},
			delimiters: []string{","},
			SkipEmpty:  false,
		}

		if err := flag.Set("a,,b,,c"); err != nil {
			t.Errorf("Set failed: %v", err)
		}
		result = flag.Get()
		expected = []string{"a", "", "b", "", "c"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// 测试SetSkipEmpty方法
	t.Run("SetSkipEmptyMethod", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				defValue: []string{},
				value:    new([]string),
			},
			delimiters: []string{","},
		}

		// 设置SkipEmpty=true
		flag.SetSkipEmpty(true)
		if err := flag.Set("x,,y"); err != nil {
			t.Errorf("Set failed: %v", err)
		}
		result := flag.Get()
		expected := []string{"x", "y"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// 动态修改为SkipEmpty=false
		flag.SetSkipEmpty(false)
		if err := flag.Set("z,,w"); err != nil {
			t.Errorf("Set failed: %v", err)
		}
		result = flag.Get()
		expected = []string{"x", "y", "z", "", "w"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// 测试错误情况
	t.Run("ErrorHandling", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				defValue: []string{},
				value:    new([]string),
			},
			delimiters: []string{","},
		}

		// 测试空输入
		if err := flag.Set(""); err == nil {
			t.Error("Expected error for empty input, got nil")
		} else if !strings.Contains(err.Error(), "slice cannot be empty") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// 测试新增的Len方法返回长度是否符合预期
	t.Run("Len", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				defValue: []string{"a", "b", "c"},
				value:    new([]string),
			},
			delimiters: []string{","},
		}

		if err := flag.Set("a,b,c"); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if flag.Len() != 3 {
			t.Errorf("Expected length 3, got %d", flag.Len())
		}
	})
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
					longName:  "intFlag",
					shortName: "i",
					defValue:  0,
					usage:     "整数标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "IntFlag已设置值",
			f: &IntFlag{
				BaseFlag: BaseFlag[int]{
					longName:  "intFlag",
					shortName: "i",
					defValue:  0,
					usage:     "整数标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*IntFlag).Set(100) },
		},
		{
			name: "IntFlag重置后",
			f: &IntFlag{
				BaseFlag: BaseFlag[int]{
					longName:  "intFlag",
					shortName: "i",
					defValue:  0,
					usage:     "整数标志测试",
				},
			},
			setValue: func(f Flag) error {
				if err := f.(*IntFlag).Set(100); err != nil {
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
					longName:  "strFlag",
					shortName: "s",
					defValue:  "default",
					usage:     "字符串标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "StringFlag已设置值",
			f: &StringFlag{
				BaseFlag: BaseFlag[string]{
					longName:  "strFlag",
					shortName: "s",
					defValue:  "default",
					usage:     "字符串标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*StringFlag).Set("test") },
		},

		// BoolFlag测试用例
		{
			name: "BoolFlag未设置值",
			f: &BoolFlag{
				BaseFlag: BaseFlag[bool]{
					longName:  "boolFlag",
					shortName: "b",
					defValue:  false,
					usage:     "布尔标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "BoolFlag已设置值",
			f: &BoolFlag{
				BaseFlag: BaseFlag[bool]{
					longName:  "boolFlag",
					shortName: "b",
					defValue:  false,
					usage:     "布尔标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*BoolFlag).Set(true) },
		},

		// FloatFlag测试用例
		{
			name: "FloatFlag未设置值",
			f: &FloatFlag{
				BaseFlag: BaseFlag[float64]{
					longName:  "floatFlag",
					shortName: "f",
					defValue:  0.0,
					usage:     "浮点标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "FloatFlag已设置值",
			f: &FloatFlag{
				BaseFlag: BaseFlag[float64]{
					longName:  "floatFlag",
					shortName: "f",
					defValue:  0.0,
					usage:     "浮点标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*FloatFlag).Set(3.14) },
		},

		// DurationFlag测试用例
		{
			name: "DurationFlag未设置值",
			f: &DurationFlag{
				BaseFlag: BaseFlag[time.Duration]{
					longName:  "durationFlag",
					shortName: "d",
					defValue:  0,
					usage:     "时间间隔标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "DurationFlag已设置值",
			f: &DurationFlag{
				BaseFlag: BaseFlag[time.Duration]{
					longName:  "durationFlag",
					shortName: "d",
					defValue:  0,
					usage:     "时间间隔标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*DurationFlag).Set((5 * time.Second).String()) },
		},

		// EnumFlag测试用例
		{
			name: "EnumFlag未设置值",
			f: &EnumFlag{
				BaseFlag: BaseFlag[string]{
					longName:  "enumFlag",
					shortName: "e",
					defValue:  "default",
					usage:     "枚举标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "EnumFlag已设置值",
			f: &EnumFlag{
				BaseFlag: BaseFlag[string]{
					longName:  "enumFlag",
					shortName: "e",
					defValue:  "default",
					usage:     "枚举标志测试",
				},
			},
			setValue: func(f Flag) error { return f.(*EnumFlag).Set("option1") },
		},

		// SliceFlag测试用例
		{
			name: "SliceFlag未设置值",
			f: &SliceFlag{
				BaseFlag: BaseFlag[[]string]{
					longName:  "sliceFlag",
					shortName: "sl",
					defValue:  []string{"default"},
					usage:     "切片标志测试",
				},
			},
			setValue: func(f Flag) error { return nil },
		},
		{
			name: "SliceFlag已设置值",
			f: &SliceFlag{
				BaseFlag: BaseFlag[[]string]{
					longName:  "sliceFlag",
					shortName: "sl",
					defValue:  []string{"default"},
					usage:     "切片标志测试",
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
			shouldBeSet := true
			if strings.Contains(tc.name, "未设置值") || strings.Contains(tc.name, "重置后") {
				shouldBeSet = false
			}

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
