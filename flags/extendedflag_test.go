package flags

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestTimeFlag_BasicFunctionality 测试TimeFlag的基本功能
func TestTimeFlag_BasicFunctionality(t *testing.T) {
	defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	flag := &TimeFlag{
		BaseFlag: BaseFlag[time.Time]{
			initialValue: defaultTime,
			value:        new(time.Time),
		},
	}

	// 测试默认值
	if flag.GetDefault() != defaultTime {
		t.Errorf("默认值应为%v, 实际为%v", defaultTime, flag.Get())
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试设置有效值
	testTime := time.Date(2024, time.March, 15, 14, 30, 0, 0, time.UTC)
	if err := flag.Set(testTime.Format(time.RFC3339)); err != nil {
		t.Fatalf("设置值失败: %v", err)
	}
	if !flag.Get().Equal(testTime) {
		t.Errorf("期望值%v, 实际值%v", testTime, flag.Get())
	}
	if !flag.IsSet() {
		t.Error("设置值后IsSet应返回true")
	}

	// 测试重置功能
	flag.Reset()
	if !flag.Get().Equal(defaultTime) {
		t.Errorf("重置后应返回默认值%v, 实际为%v", defaultTime, flag.Get())
	}
	if flag.IsSet() {
		t.Error("重置后IsSet应返回false")
	}
}

// TestTimeFlag_InvalidFormat 测试无效时间格式
func TestTimeFlag_InvalidFormat(t *testing.T) {
	flag := &TimeFlag{
		BaseFlag: BaseFlag[time.Time]{
			value: new(time.Time),
		},
	}

	invalidInputs := []string{
		"2024/03/15",
		"15-03-2024",
		"invalid-date",
		"",
	}

	for _, input := range invalidInputs {
		if err := flag.Set(input); err == nil {
			t.Errorf("输入'%s'应返回错误", input)
		}
	}
}

// TestTimeFlag_SupportedFormats 测试支持的时间格式
func TestTimeFlag_SupportedFormats(t *testing.T) {
	flag := &TimeFlag{
		BaseFlag: BaseFlag[time.Time]{
			value: new(time.Time),
		},
	}

	testCases := []struct {
		input    string
		expected time.Time
	}{{
		input:    "2024-03-15",
		expected: time.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC),
	}, {
		input:    "2024-03-15 14:30",
		expected: time.Date(2024, time.March, 15, 14, 30, 0, 0, time.UTC),
	}, {
		input:    "2024-03-15 14:30:45",
		expected: time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC),
	}, {
		input:    "2024-03-15T14:30:45Z",
		expected: time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC),
	}}

	for _, tc := range testCases {
		if err := flag.Set(tc.input); err != nil {
			t.Errorf("解析'%s'失败: %v", tc.input, err)
			continue
		}
		if !flag.Get().Equal(tc.expected) {
			t.Errorf("输入'%s'期望%v, 实际%v", tc.input, tc.expected, flag.Get())
		}
	}
}

// TestTimeFlag_TypeAndString 测试类型和字符串表示
func TestTimeFlag_TypeAndString(t *testing.T) {
	flag := &TimeFlag{}

	// 测试Type()
	if flag.Type() != FlagTypeTime {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeTime, flag.Type())
	}

	// 测试String()
	testTime := time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC)
	if err := flag.Set(testTime.Format(time.RFC3339)); err != nil {
		t.Fatalf("设置值失败: %v", err)
	}
	if flag.String() != testTime.Format(time.RFC3339) {
		t.Errorf("String()期望'%s', 实际'%s'", testTime.Format(time.RFC3339), flag.String())
	}
}

// TestDurationFlag_ValidParsing 测试有效的时间格式解析
func TestDurationFlag_ValidParsing(t *testing.T) {
	flag := &DurationFlag{
		BaseFlag: BaseFlag[time.Duration]{
			initialValue: 0,
			value:        new(time.Duration),
		},
	}

	// 测试用例集合
	testCases := []struct {
		input    string
		expected time.Duration
	}{{
		input:    "5s",
		expected: 5 * time.Second,
	}, {
		input:    "1m30s",
		expected: 90 * time.Second,
	}, {
		input:    "2h",
		expected: 2 * time.Hour,
	}, {
		input:    "100ms",
		expected: 100 * time.Millisecond,
	}}

	for _, tc := range testCases {
		if err := flag.Set(tc.input); err != nil {
			t.Errorf("解析 %s 失败: %v", tc.input, err)
			continue
		}
		if flag.Get() != tc.expected {
			t.Errorf("%s 期望 %v, 实际 %v", tc.input, tc.expected, flag.Get())
		}
	}
}

// TestDurationFlag_InvalidCases 测试无效输入和边界条件
func TestDurationFlag_InvalidCases(t *testing.T) {
	flag := &DurationFlag{
		BaseFlag: BaseFlag[time.Duration]{
			initialValue: 0,
			value:        new(time.Duration),
		},
	}

	// 测试空输入
	if err := flag.Set(""); err == nil {
		t.Error("空输入应返回错误")
	}

	// 测试无效格式
	if err := flag.Set("invalid"); err == nil {
		t.Error("无效格式应返回错误")
	}

	// 测试负值
	if err := flag.Set("-5s"); err == nil {
		t.Error("负值应返回错误")
	}
}

// TestDurationFlag_TypeAndString 测试类型和字符串表示
func TestDurationFlag_TypeAndString(t *testing.T) {
	flag := &DurationFlag{
		BaseFlag: BaseFlag[time.Duration]{
			initialValue: 5 * time.Second,
			value:        new(time.Duration),
		},
	}

	// 测试类型
	if flag.Type() != FlagTypeDuration {
		t.Errorf("Type() 期望 %d, 实际 %d", FlagTypeDuration, flag.Type())
	}

	// 测试String()
	if err := flag.Set("2m"); err != nil {
		t.Fatalf("设置值失败: %v", err)
	}
	if flag.String() != "2m0s" {
		t.Errorf("String() 期望 '2m0s', 实际 '%s'", flag.String())
	}
}

// TestMapFlag_BasicParsing 测试基本的键值对解析功能
func TestMapFlag_BasicParsing(t *testing.T) {
	flag := &MapFlag{}
	flag.SetDelimiters(",", "=") // 显式设置分隔符
	err := flag.Set("name=test,env=dev")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// 验证Get()返回正确的map
	result := flag.Get()
	expectedMap := map[string]string{"name": "test", "env": "dev"}
	if !reflect.DeepEqual(result, expectedMap) {
		t.Errorf("Get() returned %v, expected %v", result, expectedMap)
	}

	// 验证String()输出正确
	actualMap := make(map[string]string)
	parts := strings.Split(flag.String(), ",")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			actualMap[kv[0]] = kv[1]
		}
	}

	if actualMap["name"] != "test" || actualMap["env"] != "dev" || len(actualMap) != 2 {
		t.Errorf("String() returned map %v, expected {name:test, env:dev}", actualMap)
	}
}

// TestMapFlag_CustomDelimiters 测试自定义分隔符
func TestMapFlag_CustomDelimiters(t *testing.T) {
	flag := &MapFlag{}
	flag.SetDelimiters("; ", ":")

	err := flag.Set("name:test; env:prod")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	result := flag.Get()
	expected := map[string]string{"name": "test", "env": "prod"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestMapFlag_IgnoreCase 测试忽略键的大小写
func TestMapFlag_IgnoreCase(t *testing.T) {
	flag := &MapFlag{}
	flag.SetDelimiters(",", "=")
	flag.SetIgnoreCase(true)

	err := flag.Set("Name=test,NAME=override")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	result := flag.Get()
	// 所有键应该被转换为小写
	expected := map[string]string{"name": "override"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestMapFlag_Errors 测试错误情况处理
func TestMapFlag_Errors(t *testing.T) {
	flag := &MapFlag{}
	flag.SetDelimiters(",", "=")

	// 测试空值
	err := flag.Set("")
	if err == nil || !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("Expected empty value error, got %v", err)
	}

	// 测试格式错误的键值对
	err = flag.Set("invalid-key")
	if err == nil || !strings.Contains(err.Error(), "invalid key-value pair format") {
		t.Errorf("Expected format error, got %v", err)
	}

	// 测试空键
	err = flag.Set("=emptykey")
	if err == nil || !strings.Contains(err.Error(), "empty key") {
		t.Errorf("Expected empty key error, got %v", err)
	}

	// 测试空值
	err = flag.Set("key=")
	if err == nil || !strings.Contains(err.Error(), "empty value") {
		t.Errorf("Expected empty value error, got %v", err)
	}
}

// TestSliceFlag_BasicParsing 测试基本切片解析功能
func TestSliceFlag_BasicParsing(t *testing.T) {
	flag := &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{},
			value:        new([]string),
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
	expected = []string{"d"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	result = flag.Get()
	expected = []string{"d"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestSliceFlag_SkipEmpty 测试空元素过滤功能
func TestSliceFlag_SkipEmpty(t *testing.T) {
	// 测试SkipEmpty=true情况
	flag := &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{},
			value:        new([]string),
		},
		delimiters: []string{","},
		skipEmpty:  true,
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
			initialValue: []string{},
			value:        new([]string),
		},
		delimiters: []string{","},
		skipEmpty:  false,
	}

	if err := flag.Set("a,,b,,c"); err != nil {
		t.Errorf("Set failed: %v", err)
	}
	result = flag.Get()
	expected = []string{"a", "", "b", "", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestSliceFlag_LenAndContains 测试Len和Contains方法
func TestSliceFlag_LenAndContains(t *testing.T) {
	flag := &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{"x", "y"},
			value:        new([]string),
		},
		delimiters: []string{","},
	}

	// 测试Len
	if flag.Len() != 2 {
		t.Errorf("Expected length 2, got %d", flag.Len())
	}

	// 设置值后测试
	if err := flag.Set("a,b,c"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if flag.Len() != 3 {
		t.Errorf("Expected length 3, got %d", flag.Len())
	}

	// 测试Contains
	if !flag.Contains("b") {
		t.Error("Should contain 'b'")
	}
	if flag.Contains("z") {
		t.Error("Should not contain 'z'")
	}
}

// TestSliceFlag_ClearAndRemove 测试Clear和Remove方法
func TestSliceFlag_ClearAndRemove(t *testing.T) {
	flag := &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{},
			value:        new([]string),
		},
		delimiters: []string{","},
	}

	// 设置初始值
	if err := flag.Set("a,b,c,d"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// 测试Remove
	if err := flag.Remove("b"); err != nil {
		t.Errorf("Remove failed: %v", err)
	}
	result := flag.Get()
	expected := []string{"a", "c", "d"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After Remove, expected %v, got %v", expected, result)
	}

	// 测试Clear
	if err := flag.Clear(); err != nil {
		t.Errorf("Clear failed: %v", err)
	}
	if flag.Len() != 0 {
		t.Errorf("After Clear, expected length 0, got %d", flag.Len())
	}
}

// TestSliceFlag_Sort 测试Sort方法
func TestSliceFlag_Sort(t *testing.T) {
	flag := &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{},
			value:        new([]string),
		},
		delimiters: []string{","},
	}

	if err := flag.Set("c,a,b"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if err := flag.Sort(); err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	result := flag.Get()
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After Sort, expected %v, got %v", expected, result)
	}
}

// TestSliceFlag 测试SliceFlag的功能
func TestSliceFlag(t *testing.T) {
	// 测试基本切片解析功能
	t.Run("BasicSliceParsing", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				initialValue: []string{},
				value:        new([]string),
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
		expected = []string{"d"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// 测试空元素过滤功能
	t.Run("EmptyElementFiltering", func(t *testing.T) {
		// 测试SkipEmpty=true情况
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				initialValue: []string{},
				value:        new([]string),
			},
			delimiters: []string{","},
			skipEmpty:  true,
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
				initialValue: []string{},
				value:        new([]string),
			},
			delimiters: []string{","},
			skipEmpty:  false,
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
				initialValue: []string{},
				value:        new([]string),
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
		expected = []string{"z", "", "w"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// 测试错误情况
	t.Run("ErrorHandling", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				initialValue: []string{},
				value:        new([]string),
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
				initialValue: []string{"a", "b", "c"},
				value:        new([]string),
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
