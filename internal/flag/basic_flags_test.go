package flag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// TestStringFlag 测试字符串标志
func TestStringFlag(t *testing.T) {
	// 创建字符串标志
	flag := NewStringFlag("output", "o", "输出文件", "default.txt")

	// 测试初始值
	if flag.Get() != "default.txt" {
		t.Errorf("Expected default value 'default.txt', got '%s'", flag.Get())
	}

	// 测试设置值
	err := flag.Set("test.txt")
	if err != nil {
		t.Errorf("Unexpected error setting value: %v", err)
	}

	if flag.Get() != "test.txt" {
		t.Errorf("Expected 'test.txt', got '%s'", flag.Get())
	}

	// 测试是否已设置
	if !flag.IsSet() {
		t.Error("Expected flag to be set")
	}

	// 测试重置
	flag.Reset()
	if flag.Get() != "default.txt" {
		t.Errorf("Expected default value after reset, got '%s'", flag.Get())
	}

	if flag.IsSet() {
		t.Error("Expected flag to be not set after reset")
	}

	// 测试基本属性
	if flag.Name() != "output" {
		t.Errorf("Expected name 'output', got '%s'", flag.Name())
	}

	if flag.ShortName() != "o" {
		t.Errorf("Expected short name 'o', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "输出文件" {
		t.Errorf("Expected description '输出文件', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "string" {
		t.Errorf("Expected type 'string', got '%s'", flag.Type().String())
	}
}

// TestIntFlag 测试整数标志
func TestIntFlag(t *testing.T) {
	// 创建整数标志
	flag := NewIntFlag("count", "c", "计数器", 10)

	// 测试初始值
	if flag.Get() != 10 {
		t.Errorf("Expected default value 10, got %d", flag.Get())
	}

	// 测试设置正数
	err := flag.Set("42")
	if err != nil {
		t.Errorf("Unexpected error setting positive value: %v", err)
	}

	if flag.Get() != 42 {
		t.Errorf("Expected 42, got %d", flag.Get())
	}

	// 测试设置负数
	err = flag.Set("-10")
	if err != nil {
		t.Errorf("Unexpected error setting negative value: %v", err)
	}

	if flag.Get() != -10 {
		t.Errorf("Expected -10, got %d", flag.Get())
	}

	// 测试无效值
	err = flag.Set("invalid")
	if err == nil {
		t.Error("Expected error for invalid integer value")
	}

	// 测试重置
	flag.Reset()
	if flag.Get() != 10 {
		t.Errorf("Expected default value after reset, got %d", flag.Get())
	}

	// 测试基本属性
	if flag.Name() != "count" {
		t.Errorf("Expected name 'count', got '%s'", flag.Name())
	}

	if flag.ShortName() != "c" {
		t.Errorf("Expected short name 'c', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "计数器" {
		t.Errorf("Expected description '计数器', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "int" {
		t.Errorf("Expected type 'int', got '%s'", flag.Type().String())
	}
}

// TestBoolFlag 测试布尔标志
func TestBoolFlag(t *testing.T) {
	// 创建布尔标志
	flag := NewBoolFlag("verbose", "v", "详细输出", false)

	// 测试初始值
	if flag.Get() != false {
		t.Errorf("Expected default value false, got %v", flag.Get())
	}

	// 测试设置true
	err := flag.Set("true")
	if err != nil {
		t.Errorf("Unexpected error setting true: %v", err)
	}

	if flag.Get() != true {
		t.Errorf("Expected true, got %v", flag.Get())
	}

	// 测试设置false
	err = flag.Set("false")
	if err != nil {
		t.Errorf("Unexpected error setting false: %v", err)
	}

	if flag.Get() != false {
		t.Errorf("Expected false, got %v", flag.Get())
	}

	// 测试其他有效值
	validValues := []string{"1", "0", "TRUE", "FALSE", "t", "f", "T", "F", "true", "false"}
	expectedResults := []bool{true, false, true, false, true, false, true, false, true, false}

	for i, value := range validValues {
		err = flag.Set(value)
		if err != nil {
			t.Errorf("Unexpected error setting '%s': %v", value, err)
		}

		if flag.Get() != expectedResults[i] {
			t.Errorf("Expected %v for '%s', got %v", expectedResults[i], value, flag.Get())
		}
	}

	// 测试重置
	flag.Reset()
	if flag.Get() != false {
		t.Errorf("Expected default value after reset, got %v", flag.Get())
	}

	// 测试基本属性
	if flag.Name() != "verbose" {
		t.Errorf("Expected name 'verbose', got '%s'", flag.Name())
	}

	if flag.ShortName() != "v" {
		t.Errorf("Expected short name 'v', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "详细输出" {
		t.Errorf("Expected description '详细输出', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "bool" {
		t.Errorf("Expected type 'bool', got '%s'", flag.Type().String())
	}
}

// TestFloat64Flag 测试浮点数标志
func TestFloat64Flag(t *testing.T) {
	// 创建浮点数标志
	flag := NewFloat64Flag("ratio", "r", "比例", 0.5)

	// 测试初始值
	if flag.Get() != 0.5 {
		t.Errorf("Expected default value 0.5, got %f", flag.Get())
	}

	// 测试设置整数
	err := flag.Set("42")
	if err != nil {
		t.Errorf("Unexpected error setting integer value: %v", err)
	}

	if flag.Get() != 42.0 {
		t.Errorf("Expected 42.0, got %f", flag.Get())
	}

	// 测试设置小数
	err = flag.Set("3.14159")
	if err != nil {
		t.Errorf("Unexpected error setting decimal value: %v", err)
	}

	if flag.Get() != 3.14159 {
		t.Errorf("Expected 3.14159, got %f", flag.Get())
	}

	// 测试设置负数
	err = flag.Set("-2.5")
	if err != nil {
		t.Errorf("Unexpected error setting negative value: %v", err)
	}

	if flag.Get() != -2.5 {
		t.Errorf("Expected -2.5, got %f", flag.Get())
	}

	// 测试无效值
	err = flag.Set("invalid")
	if err == nil {
		t.Error("Expected error for invalid float value")
	}

	// 测试重置
	flag.Reset()
	if flag.Get() != 0.5 {
		t.Errorf("Expected default value after reset, got %f", flag.Get())
	}

	// 测试基本属性
	if flag.Name() != "ratio" {
		t.Errorf("Expected name 'ratio', got '%s'", flag.Name())
	}

	if flag.ShortName() != "r" {
		t.Errorf("Expected short name 'r', got '%s'", flag.ShortName())
	}

	if flag.Desc() != "比例" {
		t.Errorf("Expected description '比例', got '%s'", flag.Desc())
	}

	if flag.Type().String() != "float64" {
		t.Errorf("Expected type 'float64', got '%s'", flag.Type().String())
	}
}

// TestStringFlagValidator 测试字符串标志的验证器功能
func TestStringFlagValidator(t *testing.T) {
	// 创建字符串标志
	flag := NewStringFlag("username", "u", "用户名", "guest")

	// 测试1：检查初始状态
	if flag.HasValidator() {
		t.Error("Expected no validator initially")
	}

	// 测试2：设置长度验证器（3-20个字符）
	flag.SetValidator(func(value string) error {
		if len(value) < 3 {
			return fmt.Errorf("用户名长度 %d 小于最小值 3", len(value))
		}
		if len(value) > 20 {
			return fmt.Errorf("用户名长度 %d 大于最大值 20", len(value))
		}
		return nil
	})

	// 验证验证器已设置
	if !flag.HasValidator() {
		t.Error("Expected validator to be set")
	}

	// 测试3：验证通过的情况
	err := flag.Set("john")
	if err != nil {
		t.Errorf("Unexpected error for valid username: %v", err)
	}

	if flag.Get() != "john" {
		t.Errorf("Expected 'john', got '%s'", flag.Get())
	}

	// 测试4：验证失败的情况 - 太短
	err = flag.Set("ab")
	if err == nil {
		t.Error("Expected error for username that is too short")
	}

	if flag.Get() != "john" {
		t.Errorf("Expected value to remain 'john' after failed validation, got '%s'", flag.Get())
	}

	// 测试5：验证失败的情况 - 太长
	longName := "this_is_a_very_long_username_that_exceeds_limit"
	err = flag.Set(longName)
	if err == nil {
		t.Error("Expected error for username that is too long")
	}

	if flag.Get() != "john" {
		t.Errorf("Expected value to remain 'john' after failed validation, got '%s'", flag.Get())
	}

	// 测试6：清除验证器
	flag.ClearValidator()

	if flag.HasValidator() {
		t.Error("Expected validator to be cleared")
	}

	// 清除后应该可以设置任意值
	err = flag.Set("ab")
	if err != nil {
		t.Errorf("Unexpected error after clearing validator: %v", err)
	}

	if flag.Get() != "ab" {
		t.Errorf("Expected 'ab', got '%s'", flag.Get())
	}

	// 测试7：覆盖验证器
	flag.SetValidator(func(value string) error {
		if !strings.HasPrefix(value, "user_") {
			return fmt.Errorf("用户名必须以 'user_' 开头")
		}
		return nil
	})

	// 新验证器应该生效
	err = flag.Set("john")
	if err == nil {
		t.Error("Expected error for username without prefix")
	}

	err = flag.Set("user_john")
	if err != nil {
		t.Errorf("Unexpected error for valid username with prefix: %v", err)
	}

	if flag.Get() != "user_john" {
		t.Errorf("Expected 'user_john', got '%s'", flag.Get())
	}

	// 测试8：空字符串不经过验证器
	flag.SetValidator(func(value string) error {
		return fmt.Errorf("不应该被调用")
	})

	// 空字符串应该直接设置，不经过验证器
	err = flag.Set("")
	if err != nil {
		t.Errorf("Unexpected error for empty string: %v", err)
	}

	if flag.Get() != "" {
		t.Errorf("Expected empty string, got '%s'", flag.Get())
	}

	// 测试9：重置不影响验证器
	flag.Reset()

	if !flag.HasValidator() {
		t.Error("Expected validator to remain after reset")
	}

	// 重新设置长度验证器以测试重置后的验证器
	flag.SetValidator(func(value string) error {
		if len(value) < 3 {
			return fmt.Errorf("用户名长度 %d 小于最小值 3", len(value))
		}
		return nil
	})

	// 验证器仍然生效
	err = flag.Set("ab")
	if err == nil {
		t.Error("Expected error for short username after reset")
	}
}

// TestStringFlagValidatorEmail 测试邮箱格式验证器
func TestStringFlagValidatorEmail(t *testing.T) {
	// 简单的邮箱格式验证器
	emailValidator := func(value string) error {
		if !strings.Contains(value, "@") {
			return fmt.Errorf("邮箱格式无效: 缺少 @ 符号")
		}
		parts := strings.Split(value, "@")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("邮箱格式无效: %s", value)
		}
		if !strings.Contains(parts[1], ".") {
			return fmt.Errorf("邮箱格式无效: 域名缺少 . 符号")
		}
		return nil
	}

	// 创建邮箱标志
	email := NewStringFlag("email", "e", "邮箱地址", "")
	email.SetValidator(emailValidator)

	// 测试有效邮箱
	validEmails := []string{
		"user@example.com",
		"test.user@test.co.uk",
		"admin@domain.org",
	}

	for _, validEmail := range validEmails {
		err := email.Set(validEmail)
		if err != nil {
			t.Errorf("Unexpected error for valid email '%s': %v", validEmail, err)
		}
	}

	// 测试无效邮箱
	invalidEmails := []string{
		"invalid",      // 缺少 @
		"@example.com", // 缺少用户名
		"user@",        // 缺少域名
		"user@domain",  // 缺少域名后缀
	}

	for _, invalidEmail := range invalidEmails {
		err := email.Set(invalidEmail)
		if err == nil {
			t.Errorf("Expected error for invalid email '%s'", invalidEmail)
		}
	}

	// 测试空字符串不经过验证器
	err := email.Set("")
	if err != nil {
		t.Errorf("Unexpected error for empty string: %v", err)
	}
}

// TestStringFlagValidatorPort 测试端口号验证器
func TestStringFlagValidatorPort(t *testing.T) {
	// 端口号验证器（1-65535）
	portValidator := func(value string) error {
		port, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("端口号必须是整数: %s", value)
		}
		if port < 1 || port > 65535 {
			return fmt.Errorf("端口号 %d 超出范围 [1, 65535]", port)
		}
		return nil
	}

	// 创建端口标志
	port := NewStringFlag("port", "p", "端口号", "8080")
	port.SetValidator(portValidator)

	// 测试有效端口
	validPorts := []string{"1", "80", "8080", "65535"}
	for _, validPort := range validPorts {
		err := port.Set(validPort)
		if err != nil {
			t.Errorf("Unexpected error for valid port '%s': %v", validPort, err)
		}
	}

	// 测试无效端口
	invalidPorts := []string{"0", "-1", "65536", "abc", "80.5"}
	for _, invalidPort := range invalidPorts {
		err := port.Set(invalidPort)
		if err == nil {
			t.Errorf("Expected error for invalid port '%s'", invalidPort)
		}
	}
}

// TestStringFlagValidatorNil 测试 nil 验证器
func TestStringFlagValidatorNil(t *testing.T) {
	flag := NewStringFlag("test", "t", "测试", "default")

	// 设置 nil 验证器
	flag.SetValidator(nil)

	// 应该可以设置任意值
	err := flag.Set("any value")
	if err != nil {
		t.Errorf("Unexpected error with nil validator: %v", err)
	}

	if flag.Get() != "any value" {
		t.Errorf("Expected 'any value', got '%s'", flag.Get())
	}
}

// TestStringFlagValidatorMultiple 测试多个验证器场景
func TestStringFlagValidatorMultiple(t *testing.T) {
	flag := NewStringFlag("password", "p", "密码", "")

	// 复杂的密码验证器
	passwordValidator := func(value string) error {
		if len(value) < 8 {
			return fmt.Errorf("密码长度必须至少 8 个字符")
		}
		if len(value) > 32 {
			return fmt.Errorf("密码长度不能超过 32 个字符")
		}
		hasUpper := false
		hasLower := false
		hasDigit := false
		for _, c := range value {
			switch {
			case c >= 'A' && c <= 'Z':
				hasUpper = true
			case c >= 'a' && c <= 'z':
				hasLower = true
			case c >= '0' && c <= '9':
				hasDigit = true
			}
		}
		if !hasUpper {
			return fmt.Errorf("密码必须包含至少一个大写字母")
		}
		if !hasLower {
			return fmt.Errorf("密码必须包含至少一个小写字母")
		}
		if !hasDigit {
			return fmt.Errorf("密码必须包含至少一个数字")
		}
		return nil
	}

	flag.SetValidator(passwordValidator)

	// 测试有效密码
	validPasswords := []string{
		"Password123",
		"Abcdefg1",
		"Xyz789abc",
	}

	for _, validPassword := range validPasswords {
		err := flag.Set(validPassword)
		if err != nil {
			t.Errorf("Unexpected error for valid password '%s': %v", validPassword, err)
		}
	}

	// 测试无效密码 - 太短
	err := flag.Set("Short1")
	if err == nil {
		t.Error("Expected error for password that is too short")
	}

	// 测试无效密码 - 缺少大写字母
	err = flag.Set("password123")
	if err == nil {
		t.Error("Expected error for password without uppercase letter")
	}

	// 测试无效密码 - 缺少小写字母
	err = flag.Set("PASSWORD123")
	if err == nil {
		t.Error("Expected error for password without lowercase letter")
	}

	// 测试无效密码 - 缺少数字
	err = flag.Set("Passwordabc")
	if err == nil {
		t.Error("Expected error for password without digit")
	}
}
