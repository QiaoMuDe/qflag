package flags

import (
	"testing"
	"time"
)

// TestTimeFlag_BasicFunctionality 测试TimeFlag的基本功能
func TestTimeFlag_BasicFunctionality(t *testing.T) {
	defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	flag := &TimeFlag{
		BaseFlag: BaseFlag[time.Time]{
			defValue: defaultTime,
			value:    new(time.Time),
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
