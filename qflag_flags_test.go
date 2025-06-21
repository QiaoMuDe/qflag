package qflag

import (
	"errors"
	"fmt"
	"testing"
)

// positiveIntValidator 验证整数是否为正数（首字母小写，非导出）
type positiveIntValidator struct{}

// Validate 实现Validator接口，检查值是否为正数
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

// stringLengthValidator 验证字符串长度是否在指定范围内（首字母小写，非导出）
type stringLengthValidator struct {
	min, max int
}

// Validate 实现Validator接口，检查字符串长度
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
