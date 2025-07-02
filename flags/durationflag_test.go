package flags

import (
	"testing"
	"time"
)

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
