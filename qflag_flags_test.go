package qflag

import (
	"strings"
	"testing"
	"time"
)

// IntRangeValidator 整数范围验证器
type IntRangeValidator struct {
	min, max int
}

// Validate 验证整数是否在指定范围内
func (v IntRangeValidator) Validate(value any) error {
	val, ok := value.(int)
	if !ok {
		return NewValidationError("value is not an integer")
	}

	if val < v.min || val > v.max {
		return NewValidationErrorf("must be between %d and %d", v.min, v.max)
	}
	return nil
}

// DurationRangeValidator 时间间隔范围验证器
type DurationRangeValidator struct {
	min, max time.Duration
}

// Validate 验证时间间隔是否在指定范围内
func (v DurationRangeValidator) Validate(value any) error {
	val, ok := value.(time.Duration)
	if !ok {
		return NewValidationError("value is not a duration")
	}

	if val < v.min || val > v.max {
		return NewValidationErrorf("must be between %v and %v", v.min, v.max)
	}
	return nil
}

// StringLengthValidator 字符串长度验证器
type StringLengthValidator struct {
	min, max int
}

// Validate 验证字符串长度是否在指定范围内
func (v StringLengthValidator) Validate(value any) error {
	val, ok := value.(string)
	if !ok {
		return NewValidationError("value is not a string")
	}

	if len(val) < v.min || len(val) > v.max {
		return NewValidationErrorf("length must be between %d and %d", v.min, v.max)
	}
	return nil
}

// PositiveFloatValidator 正浮点数验证器
type PositiveFloatValidator struct{}

// Validate 验证浮点数是否为正数
func (v PositiveFloatValidator) Validate(value any) error {
	val, ok := value.(float64)
	if !ok {
		return NewValidationError("value is not a float")
	}

	if val <= 0 {
		return NewValidationError("must be greater than 0")
	}
	return nil
}

// TestIntFlag_Validator 测试IntFlag的验证器功能
func TestIntFlag_Validator(t *testing.T) {
	// 创建整数范围验证器(1-100)
	validator := IntRangeValidator{min: 1, max: 100}

	// 创建IntFlag实例
	flag := &IntFlag{
		BaseFlag: BaseFlag[int]{
			longName:  "test-int",
			Validator: validator,
		},
	}

	// 测试1: 合法值(50)
	if err := flag.Set(50); err != nil {
		t.Errorf("Set(50) 应该成功, 实际错误: %v", err)
	}

	// 测试2: 小于最小值(0)
	if err := flag.Set(0); err == nil {
		t.Error("Set(0) 应该失败, 实际成功")
	} else if !strings.Contains(err.Error(), "invalid value for test-int") {
		t.Errorf("错误信息应该包含标志名称, 实际错误: %v", err)
	}

	// 测试3: 大于最大值(150)
	if err := flag.Set(150); err == nil {
		t.Error("Set(150) 应该失败, 实际成功")
	} else if !strings.Contains(err.Error(), "must be between 1 and 100") {
		t.Errorf("错误信息应该包含范围提示, 实际错误: %v", err)
	}
}

// TestStringFlag_Validator 测试StringFlag的验证器功能
func TestStringFlag_Validator(t *testing.T) {
	// 创建字符串长度验证器(3-10)
	validator := StringLengthValidator{min: 3, max: 10}

	// 创建StringFlag实例
	flag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			longName:  "test-string",
			Validator: validator,
		},
	}

	// 测试1: 合法值("valid")
	if err := flag.Set("valid"); err != nil {
		t.Errorf("Set(\"valid\") 应该成功, 实际错误: %v", err)
	}

	// 测试2: 太短("ab")
	if err := flag.Set("ab"); err == nil {
		t.Error("Set(\"ab\") 应该失败, 实际成功")
	}

	// 测试3: 太长("thisiswaytoolong")
	if err := flag.Set("thisiswaytoolong"); err == nil {
		t.Error("Set(\"thisiswaytoolong\") 应该失败, 实际成功")
	}
}

// TestFloatFlag_Validator 测试FloatFlag的验证器功能
func TestFloatFlag_Validator(t *testing.T) {
	// 创建正浮点数验证器
	validator := PositiveFloatValidator{}

	// 创建FloatFlag实例
	flag := &FloatFlag{
		BaseFlag: BaseFlag[float64]{
			longName:  "test-float",
			Validator: validator,
		},
	}

	// 测试1: 合法值(3.14)
	if err := flag.Set(3.14); err != nil {
		t.Errorf("Set(3.14) 应该成功, 实际错误: %v", err)
	}

	// 测试2: 零值(0.0)
	if err := flag.Set(0.0); err == nil {
		t.Error("Set(0.0) 应该失败, 实际成功")
	}

	// 测试3: 负值(-2.5)
	if err := flag.Set(-2.5); err == nil {
		t.Error("Set(-2.5) 应该失败, 实际成功")
	}
}

// TestDurationFlag_Validator 测试DurationFlag的验证器功能
func TestDurationFlag_Validator(t *testing.T) {
	// 创建时间间隔范围验证器(100ms-10s)
	validator := DurationRangeValidator{min: 100 * time.Millisecond, max: 10 * time.Second}

	// 创建DurationFlag实例
	flag := &DurationFlag{
		BaseFlag: BaseFlag[time.Duration]{
			longName:  "test-duration",
			Validator: validator,
		},
	}

	// 测试1: 合法值(500ms)
	if err := flag.Set("500ms"); err != nil {
		t.Errorf("Set(\"500ms\") 应该成功, 实际错误: %v", err)
	}

	// 测试2: 小于最小值(50ms)
	if err := flag.Set("50ms"); err == nil {
		t.Error("Set(\"50ms\") 应该失败, 实际成功")
	} else if !strings.Contains(err.Error(), "must be between 100ms and 10s") {
		t.Errorf("错误信息应该包含范围提示, 实际错误: %v", err)
	}

	// 测试3: 大于最大值(15s)
	if err := flag.Set("15s"); err == nil {
		t.Error("Set(\"15s\") 应该失败, 实际成功")
	} else if !strings.Contains(err.Error(), "must be between 100ms and 10s") {
		t.Errorf("错误信息应该包含范围提示, 实际错误: %v", err)
	}
}

// TestEnumFlag_Validator 测试EnumFlag的验证器功能(同时测试内置验证和自定义验证器)
func TestEnumFlag_Validator(t *testing.T) {
	// 创建字符串长度验证器(3-5)
	validator := StringLengthValidator{min: 3, max: 5}

	// 创建EnumFlag实例
	flag := &EnumFlag{
		BaseFlag: BaseFlag[string]{
			longName:  "test-enum",
			Validator: validator,
		},
		optionMap: map[string]bool{"apple": true, "banana": true, "cherry": true},
	}

	// 测试1: 合法值("apple")
	if err := flag.Set("apple"); err != nil {
		t.Errorf("Set(\"apple\") 应该成功, 实际错误: %v", err)
	}

	// 测试2: 枚举不合法("date")
	if err := flag.Set("date"); err == nil {
		t.Error("Set(\"date\") 应该失败, 实际成功")
	} else if !strings.Contains(err.Error(), "invalid enum value") {
		t.Errorf("错误信息应该包含枚举验证失败, 实际错误: %v", err)
	}

	// 测试3: 枚举合法但长度不合法("blueberry")
	flag.optionMap["blueberry"] = true // 添加到枚举选项
	if err := flag.Set("blueberry"); err == nil {
		t.Error("Set(\"blueberry\") 应该失败, 实际成功")
	} else if !strings.Contains(err.Error(), "length must be between 3 and 5") {
		t.Errorf("错误信息应该包含长度验证失败, 实际错误: %v", err)
	}
}

// TestNoValidator 测试没有设置验证器的情况
func TestNoValidator(t *testing.T) {
	// 创建没有验证器的IntFlag
	flag := &IntFlag{
		BaseFlag: BaseFlag[int]{
			longName: "test-no-validator",
		},
	}

	// 任何值都应该被接受
	if err := flag.Set(-100); err != nil {
		t.Errorf("Set(-100) 应该成功, 实际错误: %v", err)
	}

	if err := flag.Set(999); err != nil {
		t.Errorf("Set(999) 应该成功, 实际错误: %v", err)
	}
}
