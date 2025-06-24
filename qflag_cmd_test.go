package qflag

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
	"time"
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

	cmd.SetHelp(usage)
	cmd.SetDescription(desc)

	if cmd.GetHelp() != usage {
		t.Errorf("GetHelp() = %v, 期望 %v", cmd.GetHelp(), usage)
	}
	if cmd.GetDescription() != desc {
		t.Errorf("描述() = %v, 期望 %v", cmd.GetDescription(), desc)
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

// TestDurationFlag 测试时间间隔类型标志
func TestDurationFlag(t *testing.T) {
	// 测试默认值
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	durFlag := cmd.Duration("duration", "d", 5*time.Second, "test duration flag")
	if durFlag.Get() != 5*time.Second {
		t.Errorf("时间间隔标志默认值 = %v, 期望 %v", durFlag.Get(), 5*time.Second)
	}

	// 测试长标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		durFlag := cmd.Duration("duration", "d", 5*time.Second, "test duration flag")
		err := cmd.Parse([]string{"--duration", "1m30s"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		expected := 90 * time.Second
		if durFlag.Get() != expected {
			t.Errorf("时间间隔标志 = %v, 期望 %v", durFlag.Get(), expected)
		}
	}

	// 测试短标志
	{
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		durFlag := cmd.Duration("duration", "d", 5*time.Second, "test duration flag")
		err := cmd.Parse([]string{"-d", "2h"})
		if err != nil {
			t.Fatalf("解析() 错误 = %v", err)
		}
		expected := 2 * time.Hour
		if durFlag.Get() != expected {
			t.Errorf("时间间隔标志 = %v, 期望 %v", durFlag.Get(), expected)
		}
	}

	// 测试无效格式
	{
		// 捕获标准输出和错误输出
		oldStdout := os.Stdout
		oldStderr := os.Stderr
		rOut, wOut, _ := os.Pipe()
		rErr, wErr, _ := os.Pipe()
		os.Stdout = wOut
		os.Stderr = wErr
		defer func() {
			wOut.Close()
			wErr.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			// 仅在-v模式下输出捕获的内容
			if testing.Verbose() {
				outBuf := new(bytes.Buffer)
				errBuf := new(bytes.Buffer)
				if _, err := outBuf.ReadFrom(rOut); err != nil {
					t.Errorf("从标准输出管道读取数据时出错: %v", err)
				}
				if _, err := errBuf.ReadFrom(rErr); err != nil {
					t.Errorf("从标准错误输出管道读取数据时出错: %v", err)
				}
				t.Logf("捕获的标准输出:\n%s", outBuf.String())
				t.Logf("捕获的标准错误输出:\n%s", errBuf.String())
			}
		}()

		cmd := NewCmd("test", "t", flag.ContinueOnError)
		durFlag := cmd.Duration("duration", "d", 5*time.Second, "test duration flag")
		err := cmd.Parse([]string{"--duration", "invalid"})
		if err == nil {
			t.Fatal("无效时间格式期望报错，实际得到 nil")
		}
		// 验证无效输入后默认值不变
		if durFlag.GetDefault() != 5*time.Second {
			t.Errorf("无效输入后时间间隔标志应保持默认值，实际得到 %v", durFlag.GetDefault())
		}
	}
}

// TestStringFlagWithoutShort 测试无短标志的字符串标志
func TestStringFlagWithoutShort(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "string-flag"
	defValue := "default"
	usage := "测试无短标志的字符串标志"

	f := cmd.String(flagName, "", defValue, usage)
	err := cmd.Parse([]string{"--" + flagName, "test-value"})
	if err != nil {
		t.Errorf("解析参数时出错: %v", err)
	}

	if f.Get() != "test-value" {
		t.Errorf("字符串标志的值为 %q，期望为 %q", f.Get(), "test-value")
	}
}

// TestCmd_CustomUsage 测试自定义用法信息功能
// 验证当设置了自定义用法时，Help()方法是否返回自定义内容，且输出仅在-v模式可见
func TestCmd_CustomUsage(t *testing.T) {
	// 创建测试命令
	cmd := NewCmd("testcmd", "tc", flag.ContinueOnError)
	customUsage := "testcmd [全局选项] <操作> [参数]\n\n"

	// 设置自定义用法
	cmd.SetUsageSyntax(customUsage)

	// 获取帮助信息
	helpInfo := cmd.GetHelp()

	// 验证帮助信息是否包含自定义用法
	if !strings.Contains(helpInfo, customUsage) {
		t.Errorf("自定义用法测试失败\n期望包含: %q\n实际内容: %q", customUsage, helpInfo)
	}

	// 使用t.Log输出详细信息，仅在go test -v时可见
	if testing.Verbose() {
		t.Logf("自定义用法测试通过\n帮助信息内容:\n%s", helpInfo)
	}
}

// TestCmd_DefaultUsage 测试默认用法信息生成
// 验证未设置自定义用法时，是否能正确生成默认用法
func TestCmd_DefaultUsage(t *testing.T) {
	// 创建测试命令
	cmd := NewCmd("defaultcmd", "dc", flag.ContinueOnError)
	cmd.SetUseChinese(true)

	// 添加测试标志
	cmd.String("config", "c", "配置文件路径", "/etc/config.json")

	// 获取默认帮助信息
	helpInfo := cmd.GetHelp()

	// 验证默认用法格式
	if !strings.Contains(helpInfo, "defaultcmd [选项]") {
		t.Errorf("默认用法格式错误\n实际内容: %q", helpInfo)
	}

	// 验证标志信息是否正确生成
	if !strings.Contains(helpInfo, "--config") || !strings.Contains(helpInfo, "-c") {
		t.Errorf("标志信息未正确生成\n实际内容: %q", helpInfo)
	}
}

// TestParseVsParseFlagsOnly 测试解析函数是否正确处理子命令
func TestParseVsParseFlagsOnly(t *testing.T) {
	// 测试场景1: Parse函数应正确处理子命令
	t.Run("Parse处理子命令", func(t *testing.T) {
		// 创建独立的命令结构
		parent := NewCmd("parent", "p", flag.ContinueOnError)
		child := NewCmd("child", "c", flag.ContinueOnError)
		ct := child.String("child-flag", "cf", "", "子命令标志")

		if err := parent.AddSubCmd(child); err != nil {
			t.Fatalf("添加子命令失败: %v", err)
		}

		// 执行解析
		args := []string{"child", "--child-flag", "value"}
		if err := parent.Parse(args); err != nil {
			t.Fatalf("Parse解析失败: %v", err)
		}

		// 验证子命令参数
		if len(child.Args()) > 0 {
			t.Error("子命令参数未被正确解析")
		}

		// 验证子命令标志是否被正确解析
		flagValue := ct.Get()
		if flagValue != "value" {
			t.Errorf("子命令标志值错误, 期望 'value', 实际 %q", flagValue)
		}
	})

	// 测试场景2: ParseFlagsOnly函数应忽略子命令
	t.Run("ParseFlagsOnly忽略子命令", func(t *testing.T) {
		// 创建独立的命令结构
		parent := NewCmd("parent", "p", flag.ContinueOnError)
		child := NewCmd("child", "c", flag.ContinueOnError)
		ct := child.String("child-flag", "cf", "", "子命令标志")

		if err := parent.AddSubCmd(child); err != nil {
			t.Fatalf("添加子命令失败: %v", err)
		}

		// 执行解析
		args := []string{"child", "--child-flag", "value"}
		if err := parent.ParseFlagsOnly(args); err != nil {
			t.Fatalf("ParseFlagsOnly解析失败: %v", err)
		}

		// 验证子命令未被处理
		if len(child.Args()) > 0 {
			t.Errorf("ParseFlagsOnly不应处理子命令, 但接收到参数: %v", child.Args())
		}

		// 验证子命令标志未被设置
		flagValue := ct.Get()
		if flagValue != "" {
			t.Errorf("ParseFlagsOnly不应解析子命令标志, 实际值: %q", flagValue)
		}
	})
}

// // TestBuiltinFlags 测试内置标志-v/--version和-sip/--show-install-path的功能
// func TestBuiltinFlags(t *testing.T) {
// 	// 测试根命令的--version和-v标志
// 	t.Run("root command version flags", func(t *testing.T) {
// 		// 创建带有版本信息的根命令
// 		rootCmd := NewCmd("test", "t", flag.ContinueOnError)
// 		rootCmd.SetVersion("1.0.0")

// 		// 测试--version标志
// 		args := []string{"--version"}
// 		if err := rootCmd.Parse(args); err != nil {
// 			t.Fatalf("解析--version标志失败: %v", err)
// 		}
// 		if !rootCmd.versionFlag.Get() {
// 			t.Error("--version标志未被正确设置")
// 		}

// 		// 重置命令并测试-v短标志
// 		rootCmd = NewCmd("test", "t", flag.ContinueOnError)
// 		rootCmd.SetVersion("1.0.0")
// 		args = []string{"-v"}
// 		if err := rootCmd.Parse(args); err != nil {
// 			t.Fatalf("解析-v标志失败: %v", err)
// 		}
// 		if !rootCmd.versionFlag.Get() {
// 			t.Error("-v标志未被正确设置")
// 		}
// 	})

// 	// 测试根命令的--show-install-path和-sip标志
// 	t.Run("root command install path flags", func(t *testing.T) {
// 		// 创建根命令
// 		rootCmd := NewCmd("test", "t", flag.ContinueOnError)

// 		// 重置命令并测试-sip短标志
// 		rootCmd = NewCmd("test", "t", flag.ContinueOnError)
// 		args := []string{"-sip"}
// 		if err := rootCmd.Parse(args); err != nil {
// 			t.Fatalf("解析-sip标志失败: %v", err)
// 		}
// 		if !rootCmd.showInstallPathFlag.Get() {
// 			t.Error("-sip标志未被正确设置")
// 		}
// 	})

// 	// 测试ParseFlagsOnly也能正确处理这些标志
// 	t.Run("ParseFlagsOnly handles builtin flags", func(t *testing.T) {
// 		rootCmd := NewCmd("test", "t", flag.ContinueOnError)
// 		rootCmd.SetVersion("1.0.0")

// 		args := []string{"-v", "-sip"}
// 		if err := rootCmd.ParseFlagsOnly(args); err != nil {
// 			t.Fatalf("ParseFlagsOnly解析标志失败: %v", err)
// 		}
// 		if !rootCmd.versionFlag.Get() {
// 			t.Error("ParseFlagsOnly未正确设置versionFlag")
// 		}
// 		if !rootCmd.showInstallPathFlag.Get() {
// 			t.Error("ParseFlagsOnly未正确设置showInstallPathFlag")
// 		}
// 	})
// }

func TestBuiltinFlags(t *testing.T) {
	// 捕获标准输出和标准错误输出
	var stdout, stderr bytes.Buffer
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	// 由于 os.Stdout 类型为 *os.File，不能直接赋值 *bytes.Buffer，使用 os.NewFile 无法实现，这里通过自定义输出流重定向的方式
	// 创建一个管道
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("创建管道失败: %v", err)
	}
	os.Stdout = w

	// 启动一个 goroutine 将管道中的数据写入 stdout
	go func() {
		_, copyErr := io.Copy(&stdout, r)
		if copyErr != nil {
			t.Logf("从管道复制数据到 stdout 缓冲区失败: %v", copyErr)
		}
		r.Close()
	}()
	// 由于 os.Stderr 类型为 *os.File，不能直接赋值 *bytes.Buffer，使用与处理 stdout 相同的管道方式重定向
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatalf("创建标准错误输出管道失败: %v", err)
	}
	os.Stderr = wErr

	// 启动一个 goroutine 将管道中的数据写入 stderr
	go func() {
		_, err := io.Copy(&stderr, rErr)
		if err != nil {
			t.Logf("从标准错误输出管道复制数据到 stderr 缓冲区失败: %v", err)
		}
		rErr.Close()
	}()

	defer func() {
		// 恢复标准输出和标准错误输出
		os.Stdout = oldStdout
		os.Stderr = oldStderr

		// 如果testing.Verbose()为true，则打印捕获的内容
		if testing.Verbose() {
			t.Logf("stdout: \n%s", stdout.String())
			t.Logf("stderr: \n%s", stderr.String())
		}
	}()

	// 测试根命令的--version和-v标志
	t.Run("root command version flags", func(t *testing.T) {
		// 创建带有版本信息的根命令
		rootCmd := NewCmd("test", "t", flag.ContinueOnError)
		rootCmd.SetVersion("1.0.0")

		// 测试--version标志
		args := []string{"--version"}
		if err := rootCmd.Parse(args); err != nil {
			t.Fatalf("解析--version标志失败: %v", err)
		}
		if !rootCmd.versionFlag.Get() {
			t.Error("--version标志未被正确设置")
		}

		// 重置命令并测试-v短标志
		rootCmd = NewCmd("test", "t", flag.ContinueOnError)
		rootCmd.SetVersion("1.0.0")
		args = []string{"-v"}
		if err := rootCmd.Parse(args); err != nil {
			t.Fatalf("解析-v标志失败: %v", err)
		}
		if !rootCmd.versionFlag.Get() {
			t.Error("-v标志未被正确设置")
		}
	})

	// 测试根命令的--show-install-path和-sip标志
	t.Run("root command install path flags", func(t *testing.T) {
		// 创建根命令
		rootCmd := NewCmd("test", "t", flag.ContinueOnError)

		// 重置命令并测试-sip短标志
		rootCmd = NewCmd("test", "t", flag.ContinueOnError)
		args := []string{"-sip"}
		if err := rootCmd.Parse(args); err != nil {
			t.Fatalf("解析-sip标志失败: %v", err)
		}
		if !rootCmd.showInstallPathFlag.Get() {
			t.Error("-sip标志未被正确设置")
		}
	})

	// 测试ParseFlagsOnly也能正确处理这些标志
	t.Run("ParseFlagsOnly handles builtin flags", func(t *testing.T) {
		rootCmd := NewCmd("test", "t", flag.ContinueOnError)
		rootCmd.SetVersion("1.0.0")

		args := []string{"-v", "-sip"}
		if err := rootCmd.ParseFlagsOnly(args); err != nil {
			t.Fatalf("ParseFlagsOnly解析标志失败: %v", err)
		}
		if !rootCmd.versionFlag.Get() {
			t.Error("ParseFlagsOnly未正确设置versionFlag")
		}
		if !rootCmd.showInstallPathFlag.Get() {
			t.Error("ParseFlagsOnly未正确设置showInstallPathFlag")
		}
	})
}
