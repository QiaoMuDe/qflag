package qflag

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

// TestNewCmd 测试创建新命令
func TestNewCmd(t *testing.T) {
	tests := []struct {
		name      string
		shortName string
		errorMode flag.ErrorHandling
	}{
		{"test", "t", flag.ContinueOnError},
		{"app", "a", flag.ExitOnError},
		{"tool", "tl", flag.PanicOnError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmd(tt.name, tt.shortName, tt.errorMode)
			if cmd.LongName() != tt.name {
				t.Errorf("长名称() = %v, 期望 %v", cmd.LongName(), tt.name)
			}
			if cmd.ShortName() != tt.shortName {
				t.Errorf("短名称() = %v, 期望 %v", cmd.ShortName(), tt.shortName)
			}
		})
	}
}

// TestFlagBinding 测试标志绑定功能
func TestFlagBinding(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	// 测试各种类型标志绑定
	strFlag := cmd.String("string", "s", "test string flag", "default")
	intFlag := cmd.Int("int", "i", 123, "test int flag")
	boolFlag := cmd.Bool("bool", "b", false, "test bool flag")
	floatFlag := cmd.Float("float", "f", 3.14, "test float flag")

	// 测试标志解析
	err := cmd.Parse([]string{"--string", "value", "--int", "456", "--bool", "--float", "2.718"})
	if err != nil {
		t.Fatalf("解析() 错误 = %v", err)
	}

	if strFlag.Get() != "value" {
		t.Errorf("字符串标志 = %v, 期望 %v", strFlag.Get(), "value")
	}
	if intFlag.Get() != 456 {
		t.Errorf("整数标志 = %v, 期望 %v", intFlag.Get(), 456)
	}
	if boolFlag.Get() != true {
		t.Errorf("布尔标志 = %v, 期望 %v", boolFlag.Get(), true)
	}
	if floatFlag.Get() != 2.718 {
		t.Errorf("浮点数标志 = %v, 期望 %v", floatFlag.Get(), 2.718)
	}
}

// TestSubCommand 测试子命令功能
func TestSubCommand(t *testing.T) {
	parent := NewCmd("parent", "p", flag.ContinueOnError)
	child := NewCmd("child", "c", flag.ContinueOnError)

	// 添加子命令
	if err := parent.AddSubCmd(child); err != nil {
		t.Fatalf("添加子命令() 错误 = %v", err)
	}

	// 测试子命令解析
	err := parent.Parse([]string{"child", "arg1", "arg2"})
	if err != nil {
		t.Fatalf("解析() 错误 = %v", err)
	}

	if len(child.Args()) != 2 {
		t.Errorf("参数() 长度 = %v, 期望 %v", len(child.Args()), 2)
	}
}

// TestUsageAndDescription 测试用法和描述信息
func TestUsageAndDescription(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	usage := "Custom usage message"
	desc := "Test description"

	cmd.SetUsage(usage)
	cmd.SetDescription(desc)

	if cmd.Usage() != usage {
		t.Errorf("用法() = %v, 期望 %v", cmd.Usage(), usage)
	}
	if cmd.Description() != desc {
		t.Errorf("描述() = %v, 期望 %v", cmd.Description(), desc)
	}
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	// 重定向标准输出和错误输出
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	tests := []struct {
		name      string
		errorMode flag.ErrorHandling
		args      []string
		expectErr bool
	}{
		{"continue", flag.ContinueOnError, []string{"--invalid"}, true},
		{"exit", flag.ContinueOnError, []string{"--invalid"}, true},
		{"panic", flag.ContinueOnError, []string{"--invalid"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.errorMode == flag.PanicOnError {
					if r := recover(); r == nil {
						t.Error("期望发生 panic，但未发生")
					}
				}
			}()

			cmd := NewCmd("test", "t", tt.errorMode)
			err := cmd.Parse(tt.args)

			if (err != nil) != tt.expectErr {
				t.Errorf("解析() 错误 = %v, 期望错误 %v", err, tt.expectErr)
			}

			if err != nil && !strings.Contains(err.Error(), "not defined") {
				t.Errorf("错误信息应包含 '未定义', 实际得到: %v", err)
			}
		})
	}

	// 恢复标准输出和错误输出
	w.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// 读取缓冲区内容并打印
	var buf bytes.Buffer
	io.Copy(&buf, r)
	t.Logf("测试输出:\n%s", buf.String())
}

// TestStringFlag 测试字符串类型标志
func TestStringFlag(t *testing.T) {
	// 测试默认值
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	strFlag := cmd.String("string", "s", "default", "test string flag")
	if strFlag.GetDefault() != "default" {
		t.Errorf("字符串标志默认值 = %v, 期望 %v", strFlag.GetDefault(), "default")
	}

	// 测试长标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		strFlag := cmd.String("string", "s", "default", "test string flag")
		err := cmd.Parse([]string{"--string", "value"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		if strFlag.Get() != "value" {
			t.Errorf("字符串标志 = %v, 期望 %v", strFlag.Get(), "value")
		}
	}

	// 测试短标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		strFlag := cmd.String("string", "s", "default", "test string flag")
		err := cmd.Parse([]string{"-s", "short"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		if strFlag.Get() != "short" {
			t.Errorf("字符串标志 = %s, 期望 %s", strFlag.Get(), "short")
		}
	}
}

// TestIntFlag 测试整数类型标志
func TestIntFlag(t *testing.T) {
	// 测试默认值
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	intFlag := cmd.Int("int", "i", 123, "test int flag")
	if intFlag.Get() != 123 {
		t.Errorf("整数标志默认值 = %v, 期望 %v", intFlag.Get(), 123)
	}

	// 测试长标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		intFlag := cmd.Int("int", "i", 123, "test int flag")
		err := cmd.Parse([]string{"--int", "456"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		if intFlag.Get() != 456 {
			t.Errorf("整数标志 = %v, 期望 %v", intFlag.Get(), 456)
		}
	}

	// 测试短标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		intFlag := cmd.Int("int", "i", 123, "test int flag")
		err := cmd.Parse([]string{"-i", "789"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		if intFlag.Get() != 789 {
			t.Errorf("整数标志 = %v, 期望 %v", intFlag.Get(), 789)
		}
	}
}

// TestBoolFlag 测试布尔类型标志
func TestBoolFlag(t *testing.T) {
	// 测试默认值
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	boolFlag := cmd.Bool("bool", "b", false, "test bool flag")
	if boolFlag.Get() != false {
		t.Errorf("布尔标志默认值 = %v, 期望 %v", boolFlag.Get(), false)
	}

	// 测试长标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		boolFlag := cmd.Bool("bool", "b", false, "test bool flag")
		err := cmd.Parse([]string{"--bool"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		if boolFlag.Get() != true {
			t.Errorf("布尔标志 = %v, 期望 %v", boolFlag.Get(), true)
		}
	}

	// 测试短标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		boolFlag := cmd.Bool("bool", "b", false, "test bool flag")
		err := cmd.Parse([]string{"-b"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		if boolFlag.Get() != true {
			t.Errorf("布尔标志 = %v, 期望 %v", boolFlag.Get(), true)
		}
	}
}

// TestFloatFlag 测试浮点数类型标志
func TestFloatFlag(t *testing.T) {
	// 测试默认值
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	floatFlag := cmd.Float("float", "f", 3.14, "test float flag")
	if floatFlag.Get() != 3.14 {
		t.Errorf("浮点数标志默认值 = %v, 期望 %v", floatFlag.Get(), 3.14)
	}

	// 测试长标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		floatFlag := cmd.Float("float", "f", 3.14, "test float flag")
		err := cmd.Parse([]string{"--float", "2.718"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		if floatFlag.Get() != 2.718 {
			t.Errorf("浮点数标志 = %v, 期望 %v", floatFlag.Get(), 2.718)
		}
	}

	// 测试短标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		floatFlag := cmd.Float("float", "f", 3.14, "test float flag")
		err := cmd.Parse([]string{"-f", "1.618"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		if floatFlag.Get() != 1.618 {
			t.Errorf("浮点数标志 = %v, 期望 %v", floatFlag.Get(), 1.618)
		}
	}
}

// TestEnumFlag 测试枚举类型标志
func TestEnumFlag(t *testing.T) {
	// 测试默认值
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		enumFlag := cmd.Enum("mode", "m", "test", "test", []string{"debug", "test", "prod"})
		if enumFlag.GetDefault() != "test" {
			t.Errorf("枚举标志默认值 = %v, 期望 %v", enumFlag.GetDefault(), "test")
		}
	}

	// 测试长标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		enumFlag := cmd.Enum("mode", "m", "test", "test", []string{"debug", "test", "prod"})
		err := cmd.Parse([]string{"--mode", "prod"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		if enumFlag.Get() != "prod" {
			t.Errorf("枚举标志 = %v, 期望 %v", enumFlag.Get(), "prod")
		}
	}

	// 测试短标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		enumFlag := cmd.Enum("mode", "m", "test", "test", []string{"debug", "test", "prod"})
		err := cmd.Parse([]string{"-m", "debug"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		if enumFlag.Get() != "debug" {
			t.Errorf("枚举标志 = %v, 期望 %v", enumFlag.Get(), "debug")
		}
	}

	// 测试无效值
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		enumFlag := cmd.Enum("mode", "m", "test", "test", []string{"debug", "test", "prod"})
		err := cmd.Parse([]string{"--mode", "invalid"})
		if err == nil {
			t.Fatal("无效枚举值期望报错，实际得到 nil")
		}
		if enumFlag.GetDefault() != "test" {
			t.Errorf("无效输入后枚举标志应保持默认值，实际得到 = %v, 期望 %v", enumFlag.Get(), "test")
		}
	}
}
