// validator_test.go
// 验证器测试用例
package validator

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

// TestStringLengthValidator 测试字符串长度验证器
func TestStringLengthValidator(t *testing.T) {
	tests := []struct {
		name     string
		min      int
		max      int
		value    any
		expected error
	}{
		{"valid length", 2, 5, "test", nil},
		{"min boundary", 3, 5, "abc", nil},
		{"max boundary", 2, 4, "test", nil},
		{"below min", 3, 5, "ab", fmt.Errorf("string length must be at least 3")},
		{"above max", 2, 3, "test", fmt.Errorf("string length must be at most 3")},
		{"unlimited max", 1, 0, "long string without max limit", nil},
		{"non-string type", 1, 5, 123, fmt.Errorf("value is not a string")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &StringLengthValidator{Min: tt.min, Max: tt.max}
			err := v.Validate(tt.value)
			if (err == nil && tt.expected != nil) || (err != nil && tt.expected == nil) {
				t.Errorf("expected error %v, got %v", tt.expected, err)
				return
			}
			if err != nil && err.Error() != tt.expected.Error() {
				t.Errorf("expected error message '%s', got '%s'", tt.expected.Error(), err.Error())
			}
		})
	}
}

// TestStringRegexValidator 测试字符串正则表达式验证器
func TestStringRegexValidator(t *testing.T) {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	validEmail := "test@example.com"
	invalidEmail := "invalid-email"

	tests := []struct {
		name     string
		pattern  string
		value    any
		expected error
	}{
		{"valid match", emailRegex, validEmail, nil},
		{"invalid match", emailRegex, invalidEmail, fmt.Errorf("string does not match pattern: %s", emailRegex)},
		{"empty pattern", "", validEmail, fmt.Errorf("regex pattern is empty")},
		{"invalid regex", "[a-z", validEmail, fmt.Errorf("invalid regex pattern: error parsing regexp: missing closing ]: `[a-z`")},
		{"non-string type", emailRegex, 123, fmt.Errorf("value is not a string")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &StringRegexValidator{Pattern: tt.pattern}
			err := v.Validate(tt.value)
			if (err == nil && tt.expected != nil) || (err != nil && tt.expected == nil) {
				t.Errorf("expected error %v, got %v", tt.expected, err)
				return
			}
			if err != nil && err.Error() != tt.expected.Error() {
				t.Errorf("expected error message '%s', got '%s'", tt.expected.Error(), err.Error())
			}
		})
	}

	// 测试已编译的正则表达式复用
	t.Run("compiled regex reuse", func(t *testing.T) {
		compiledRegex := regexp.MustCompile(emailRegex)
		v := &StringRegexValidator{Regex: compiledRegex}
		err := v.Validate(validEmail)
		if err != nil {
			t.Errorf("expected no error with compiled regex, got %v", err)
		}
	})
}

// TestIntRangeValidator 测试整数范围验证器
func TestIntRangeValidator(t *testing.T) {
	tests := []struct {
		name     string
		min      int64
		max      int64
		value    any
		expected error
	}{
		{"valid int", 10, 20, 15, nil},
		{"min boundary", 5, 10, 5, nil},
		{"max boundary", 5, 10, 10, nil},
		{"below min", 5, 10, 3, fmt.Errorf("value must be at least 5")},
		{"above max", 5, 10, 12, fmt.Errorf("value must be at most 10")},
		{"int8 type", 0, 100, int8(50), nil},
		{"uint type", 0, 100, uint(50), nil},
		{"uint64 type", 100, 200, uint64(150), nil},
		{"non-integer type", 5, 10, "15", fmt.Errorf("value is not an int64-compatible integer type")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &IntRangeValidator64{Min: tt.min, Max: tt.max}
			err := v.Validate(tt.value)
			if (err == nil && tt.expected != nil) || (err != nil && tt.expected == nil) {
				t.Errorf("expected error %v, got %v", tt.expected, err)
				return
			}
			if err != nil && err.Error() != tt.expected.Error() {
				t.Errorf("expected error message '%s', got '%s'", tt.expected.Error(), err.Error())
			}
		})
	}
}

// TestFloatRangeValidator 测试浮点数范围验证器
func TestFloatRangeValidator(t *testing.T) {
	tests := []struct {
		name     string
		min      float64
		max      float64
		value    any
		expected error
	}{
		{"valid float64", 0.5, 2.5, 1.5, nil},
		{"valid float32", 0.5, 2.5, float32(1.5), nil},
		{"min boundary", 1.0, 3.0, 1.0, nil},
		{"max boundary", 1.0, 3.0, 3.0, nil},
		{"below min", 1.0, 3.0, 0.5, fmt.Errorf("value must be at least 1.000000")},
		{"above max", 1.0, 3.0, 3.5, fmt.Errorf("value must be at most 3.000000")},
		{"non-float type", 1.0, 3.0, "2.5", fmt.Errorf("value is not a float type")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &FloatRangeValidator{Min: tt.min, Max: tt.max}
			err := v.Validate(tt.value)
			if (err == nil && tt.expected != nil) || (err != nil && tt.expected == nil) {
				t.Errorf("expected error %v, got %v", tt.expected, err)
				return
			}
			if err != nil && err.Error() != tt.expected.Error() {
				t.Errorf("expected error message '%s', got '%s'", tt.expected.Error(), err.Error())
			}
		})
	}
}

// TestBoolValidator 测试布尔值验证器
func TestBoolValidator(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected error
	}{
		{"valid true", true, nil},
		{"valid false", false, nil},
		{"non-bool type", 1, fmt.Errorf("value is not a boolean")},
		{"non-bool string", "true", fmt.Errorf("value is not a boolean")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &BoolValidator{}
			err := v.Validate(tt.value)
			if (err == nil && tt.expected != nil) || (err != nil && tt.expected == nil) {
				t.Errorf("expected error %v, got %v", tt.expected, err)
				return
			}
			if err != nil && err.Error() != tt.expected.Error() {
				t.Errorf("expected error message '%s', got '%s'", tt.expected.Error(), err.Error())
			}
		})
	}
}

// TestDurationValidator 测试时间间隔验证器
func TestDurationValidator(t *testing.T) {
	minute := time.Minute
	fiveMinutes := 5 * time.Minute
	invalidDurationStr := "invalid"

	tests := []struct {
		name     string
		min      time.Duration
		max      time.Duration
		value    any
		expected error
	}{
		{"valid duration string", minute, fiveMinutes, "3m", nil},
		{"valid duration type", minute, fiveMinutes, 3 * minute, nil},
		{"min boundary", minute, fiveMinutes, "1m", nil},
		{"max boundary", minute, fiveMinutes, "5m", nil},
		{"below min", minute, fiveMinutes, "30s", fmt.Errorf("duration must be at least 1m0s")},
		{"above max", minute, fiveMinutes, "6m", fmt.Errorf("duration must be at most 5m0s")},
		{"invalid duration string", minute, fiveMinutes, invalidDurationStr, fmt.Errorf("invalid duration string: time: invalid duration \"invalid\"")},
		{"non-duration type", minute, fiveMinutes, 123, fmt.Errorf("value is not a duration string or time.Duration")},
		{"unlimited max", minute, 0, "10m", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &DurationValidator{Min: tt.min, Max: tt.max}
			err := v.Validate(tt.value)
			if (err == nil && tt.expected != nil) || (err != nil && tt.expected == nil) {
				t.Errorf("expected error %v, got %v", tt.expected, err)
				return
			}
			if err != nil && err.Error() != tt.expected.Error() {
				t.Errorf("expected error message '%s', got '%s'", tt.expected.Error(), err.Error())
			}
		})
	}
}

// TestSliceLengthValidator 测试切片长度验证器
func TestSliceLengthValidator(t *testing.T) {
	validSlice := []int{1, 2, 3}
	longSlice := []int{1, 2, 3, 4, 5, 6}

	tests := []struct {
		name     string
		min      int
		max      int
		value    any
		expected error
	}{
		{"valid slice", 2, 5, validSlice, nil},
		{"min boundary", 3, 5, validSlice, nil},
		{"max boundary", 1, 3, validSlice, nil},
		{"below min", 4, 5, validSlice, fmt.Errorf("slice length must be at least 4")},
		{"above max", 1, 3, longSlice, fmt.Errorf("slice length must be at most 3")},
		{"unlimited max", 2, 0, longSlice, nil},
		{"non-slice type", 1, 5, "not a slice", fmt.Errorf("value is not a slice or array")},
		{"array type", 2, 3, [3]int{1, 2, 3}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &SliceLengthValidator{Min: tt.min, Max: tt.max}
			err := v.Validate(tt.value)
			if (err == nil && tt.expected != nil) || (err != nil && tt.expected == nil) {
				t.Errorf("expected error %v, got %v", tt.expected, err)
				return
			}
			if err != nil && err.Error() != tt.expected.Error() {
				t.Errorf("expected error message '%s', got '%s'", tt.expected.Error(), err.Error())
			}
		})
	}
}

// TestEnumValidator 测试枚举验证器
func TestEnumValidator(t *testing.T) {
	allowedValues := []any{"option1", "option2", 3}

	tests := []struct {
		name     string
		allowed  []any
		value    any
		expected error
	}{
		{"valid string value", allowedValues, "option1", nil},
		{"valid int value", allowedValues, 3, nil},
		{"invalid value", allowedValues, "option3", fmt.Errorf("value option3 is not in allowed values list")},
		{"empty allowed list", []any{}, "option1", fmt.Errorf("no allowed values specified")},
		{"different type", allowedValues, 1, fmt.Errorf("value 1 is not in allowed values list")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &EnumValidator{AllowedValues: tt.allowed}
			err := v.Validate(tt.value)
			if (err == nil && tt.expected != nil) || (err != nil && tt.expected == nil) {
				t.Errorf("expected error %v, got %v", tt.expected, err)
				return
			}
			if err != nil && err.Error() != tt.expected.Error() {
				t.Errorf("expected error message '%s', got '%s'", tt.expected.Error(), err.Error())
			}
		})
	}
}
