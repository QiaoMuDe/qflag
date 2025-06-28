package qflag

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestStringFlagLong 测试字符串类型长标志的注册和解析
func TestStringFlagLong(t *testing.T) {
	// 完全重定向标准输出到缓冲区
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 不输出任何捕获的内容
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Errorf("从管道读取数据时出错: %v", err)
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "string-flag"
	defValue := "default"
	usage := "测试字符串标志"

	// 测试String方法(仅长标志)
	f := cmd.String(flagName, "sf", defValue, usage)
	if f == nil {
		t.Fatal("String() 返回了 nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName, "test-value"})
	if err != nil {
		t.Errorf("解析参数时出错: %v", err)
	}

	// 验证值
	if f.Get() != "test-value" {
		t.Errorf("字符串标志的值为 %q，期望为 %q", f.Get(), "test-value")
	}
}

// TestStringFlagShort 测试字符串类型短标志的注册和解析
func TestStringFlagShort(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("从管道读取数据时出错: %v", err)
			} else {
				t.Logf("捕获的输出:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "s"
	defValue := "default"
	usage := "测试字符串标志"

	// 测试String方法(仅短标志)
	f := cmd.String("sf", shortName, defValue, usage)
	if f == nil {
		t.Fatal("String() 返回了 nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName, "test-value"})
	if err != nil {
		t.Errorf("解析参数时出错: %v", err)
	}

	// 验证值
	if f.Get() != "test-value" {
		t.Errorf("字符串标志的值为 %q，期望为 %q", f.Get(), "test-value")
	}
}

// TestIntFlagLong 测试整数类型长标志的注册和解析
func TestIntFlagLong(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("从管道读取数据时出错: %v", err)
			} else {
				t.Logf("捕获的输出:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "int-flag"
	defValue := 100
	usage := "测试整数标志"

	// 测试Int方法(仅长标志)
	f := cmd.Int(flagName, "if", defValue, usage)
	if f == nil {
		t.Fatal("Int() 返回了 nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName, "200"})
	if err != nil {
		t.Errorf("解析参数时出错: %v", err)
	}

	// 验证值
	if f.Get() != 200 {
		t.Errorf("整数标志的值为 %d，期望为 %d", f.Get(), 200)
	}
}

// TestIntFlagShort 测试整数类型短标志的注册和解析
func TestIntFlagShort(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("从管道读取数据时出错: %v", err)
			} else {
				t.Logf("捕获的输出:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "i"
	defValue := 100
	usage := "测试整数标志"

	// 测试Int方法(仅短标志)
	f := cmd.Int("ci", shortName, defValue, usage)
	if f == nil {
		t.Fatal("Int() 返回了 nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName, "200"})
	if err != nil {
		t.Errorf("解析参数时出错: %v", err)
	}

	// 验证值
	if f.Get() != 200 {
		t.Errorf("整数标志的值为 %d，期望为 %d", f.Get(), 200)
	}
}

// TestBoolFlagLong 测试布尔类型长标志的注册和解析
func TestBoolFlagLong(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("从管道读取数据时出错: %v", err)
			} else {
				t.Logf("捕获的输出:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "bool-flag"
	defValue := false
	usage := "测试布尔标志"

	// 测试Bool方法(仅长标志)
	f := cmd.Bool(flagName, "bl", defValue, usage)
	if f == nil {
		t.Fatal("Bool() 返回了 nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName})
	if err != nil {
		t.Errorf("解析参数时出错: %v", err)
	}

	// 验证值
	if f.Get() != true {
		t.Errorf("布尔标志的值为 %v，期望为 %v", f.Get(), true)
	}
}

// TestBoolFlagShort 测试布尔类型短标志的注册和解析
func TestBoolFlagShort(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("从管道读取数据时出错: %v", err)
			} else {
				t.Logf("捕获的输出:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "b"
	defValue := false
	usage := "测试布尔标志"

	// 测试Bool方法(仅短标志)
	f := cmd.Bool("ct", shortName, defValue, usage)
	if f == nil {
		t.Fatal("Bool() 返回了 nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName})
	if err != nil {
		t.Errorf("解析参数时出错: %v", err)
	}

	// 验证值
	if f.Get() != true {
		t.Errorf("布尔标志的值为 %v，期望为 %v", f.Get(), true)
	}
}

// TestFloatFlagLong 测试浮点数类型长标志的注册和解析
func TestFloatFlagLong(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("从管道读取数据时出错: %v", err)
			} else {
				t.Logf("捕获的输出:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "float-flag"
	defValue := 3.14
	usage := "测试浮点数标志"

	// 测试Float方法(仅长标志)
	f := cmd.Float(flagName, "ff", defValue, usage)
	if f == nil {
		t.Fatal("Int() 返回了 nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName, "6.28"})
	if err != nil {
		t.Errorf("解析参数时出错: %v", err)
	}

	// 验证值
	if f.Get() != 6.28 {
		t.Errorf("浮点数标志的值为 %v，期望为 %v", f.Get(), 6.28)
	}
}

// TestFloatFlagShort 测试浮点数类型短标志的注册和解析
func TestFloatFlagShort(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("从管道读取数据时出错: %v", err)
			} else {
				t.Logf("捕获的输出:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "f"
	defValue := 3.14
	usage := "测试浮点数标志"

	// 测试Float方法(仅短标志)
	f := cmd.Float("cf", shortName, defValue, usage)
	if f == nil {
		t.Fatal("Int() 返回了 nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName, "6.28"})
	if err != nil {
		t.Errorf("解析参数时出错: %v", err)
	}

	// 验证值
	if f.Get() != 6.28 {
		t.Errorf("浮点数标志的值为 %v，期望为 %v", f.Get(), 6.28)
	}
}

// TestParseError 测试参数解析错误
func TestParseError(t *testing.T) {
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

		// 只有在-v模式或测试失败时输出
		if testing.Verbose() || t.Failed() {
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
	cmd.Int("int-flag", "i", 0, "测试整数标志")

	// 测试无效参数
	err := cmd.Parse([]string{"--int-flag", "not-a-number"})
	if err == nil {
		t.Error("解析无效输入时 Parse() 应该返回错误")
	}
}

// TestHelpFlag 测试帮助标志
func TestHelpFlag(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r)
			if err != nil {
				t.Errorf("ReadFrom failed: %v", err)
			}
			t.Logf("捕获的输出:\n%s", buf.String())
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.String("string-flag", "s", "", "测试字符串标志")

	// 测试帮助标志
	err := cmd.Parse([]string{"--help"})
	if err != nil {
		t.Errorf("解析参数时出错: %v", err)
	}
}

// TestCmd_Name 测试Cmd的Name方法
func TestCmd_Name(t *testing.T) {
	cmd := NewCmd("testcmd", "t", flag.ContinueOnError)
	if cmd.LongName() != "testcmd" {
		t.Errorf("Name() 返回 %q，期望为 %q", cmd.LongName(), "testcmd")
	}
}

// TestCmd_ShortName 测试Cmd的ShortName方法
func TestCmd_ShortName(t *testing.T) {
	cmd := NewCmd("testcmd", "tc", flag.ContinueOnError)
	if cmd.ShortName() != "tc" {
		t.Errorf("ShortName() 返回 %q，期望为 %q", cmd.ShortName(), "tc")
	}
}

// TestCmd_Usage 测试Cmd的Usage和SetUsage方法
func TestCmd_Usage(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	usage := "测试用法"
	cmd.SetHelp(usage)
	if cmd.GetHelp() != usage {
		t.Errorf("GetHelp() 返回 %q，期望为 %q", cmd.GetHelp(), usage)
	}
}

// TestIntFlag_Interface 验证IntFlag实现了Flag接口
func TestIntFlag_Interface(t *testing.T) {
	var f flags.Flag = &flags.IntFlag{}
	_ = f // 若编译通过，则说明实现了接口
}

// TestIntFlag_Methods 测试IntFlag的各种方法
func TestIntFlag_Methods(t *testing.T) {
	defValue := 100

	// 新建子命令
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	f := cmd.Int("intflag", "i", defValue, "整数标志测试")

	if f.LongName() != "intflag" {
		t.Errorf("IntFlag.Name() 返回 %q，期望为 %q", f.LongName(), "intflag")
	}
	if f.ShortName() != "i" {
		t.Errorf("IntFlag.ShortName() 返回 %q，期望为 %q", f.ShortName(), "i")
	}
	if f.Usage() != "整数标志测试" {
		t.Errorf("IntFlag.Usage() 返回 %q，期望为 %q", f.Usage(), "整数标志测试")
	}
	if f.GetDefault() != defValue {
		t.Errorf("IntFlag.GetDefault() 返回 %v，期望为 %v", f.GetDefault(), defValue)
	}
	if f.Type() != flags.FlagTypeInt {
		t.Errorf("IntFlag.Type() 返回 %v，期望为 %v", f.Type(), flags.FlagTypeInt)
	}

	// 测试边界值
	if err := f.Set(0); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if f.Get() != 0 {
		t.Errorf("IntFlag.Get() 返回 %v，期望为 %v", f.Get(), 0)
	}

	if err := f.Set(-1); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if f.Get() != -1 {
		t.Errorf("IntFlag.Get() 返回 %v，期望为 %v", f.Get(), -1)
	}

	if err := f.Set(2147483647); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if f.Get() != 2147483647 {
		t.Errorf("IntFlag.Get() 返回 %v，期望为 %v", f.Get(), 2147483647)
	}
}

// TestStringFlag_Interface 验证StringFlag实现了Flag接口
func TestStringFlag_Interface(t *testing.T) {
	var f flags.Flag = &flags.StringFlag{}
	_ = f
}

// TestStringFlag_Methods 测试StringFlag的各种方法
func TestStringFlag_Methods(t *testing.T) {
	defValue := "default string"

	// 新建子命令
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	f := cmd.String("strflag", "s", defValue, "字符串标志测试")

	if f.LongName() != "strflag" {
		t.Errorf("StringFlag.Name() 返回 %q，期望为 %q", f.LongName(), "strflag")
	}
	if f.ShortName() != "s" {
		t.Errorf("StringFlag.ShortName() 返回 %q，期望为 %q", f.ShortName(), "s")
	}
	if f.Usage() != "字符串标志测试" {
		t.Errorf("StringFlag.Usage() 返回 %q，期望为 %q", f.Usage(), "字符串标志测试")
	}
	if f.GetDefault() != defValue {
		t.Errorf("StringFlag.GetDefault() 返回 %v，期望为 %v", f.GetDefault(), defValue)
	}
	if f.Type() != flags.FlagTypeString {
		t.Errorf("StringFlag.Type() 返回 %v，期望为 %v", f.Type(), flags.FlagTypeString)
	}

	// 测试边界值
	if err := f.Set(""); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if f.Get() != "" {
		t.Errorf("StringFlag.Get() 返回 %q，期望为 %q", f.Get(), "")
	}

	if err := f.Set("long_string_with_special_chars_!@#$%^&*()"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if f.Get() != "long_string_with_special_chars_!@#$%^&*()" {
		t.Errorf("StringFlag.Get() 返回 %q，期望为 %q", f.Get(), "long_string_with_special_chars_!@#$%^&*()")
	}
}

// TestBoolFlag_Interface 验证BoolFlag实现了Flag接口
func TestBoolFlag_Interface(t *testing.T) {
	var f flags.Flag = &flags.BoolFlag{}
	_ = f
}

// TestBoolFlag_Methods 测试BoolFlag的各种方法
func TestBoolFlag_Methods(t *testing.T) {
	defValue := true

	// 新建子命令
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	f := cmd.Bool("boolflag", "b", defValue, "布尔标志测试")

	if f.LongName() != "boolflag" {
		t.Errorf("BoolFlag.Name() 返回 %q，期望为 %q", f.LongName(), "boolflag")
	}
	if f.ShortName() != "b" {
		t.Errorf("BoolFlag.ShortName() 返回 %q，期望为 %q", f.ShortName(), "b")
	}
	if f.Usage() != "布尔标志测试" {
		t.Errorf("BoolFlag.Usage() 返回 %q，期望为 %q", f.Usage(), "布尔标志测试")
	}
	if f.GetDefault() != defValue {
		t.Errorf("BoolFlag.GetDefault() 返回 %v，期望为 %v", f.GetDefault(), defValue)
	}
	if f.Type() != flags.FlagTypeBool {
		t.Errorf("BoolFlag.Type() 返回 %v，期望为 %v", f.Type(), flags.FlagTypeBool)
	}

	// 测试切换布尔值
	if err := f.Set(false); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if f.Get() != false {
		t.Errorf("BoolFlag.Get() 返回 %v，期望为 %v", f.Get(), false)
	}

	if err := f.Set(true); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if f.Get() != true {
		t.Errorf("BoolFlag.Get() 返回 %v，期望为 %v", f.Get(), true)
	}
}

// TestFloatFlag_Interface 验证FloatFlag实现了Flag接口
func TestFloatFlag_Interface(t *testing.T) {
	var f flags.Flag = &flags.FloatFlag{}
	_ = f
}

// TestFloatFlag_Methods 测试FloatFlag的各种方法
func TestFloatFlag_Methods(t *testing.T) {
	defValue := 3.14

	// 新建子命令
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	f := cmd.Float("floatflag", "f", defValue, "浮点数标志测试")

	if f.LongName() != "floatflag" {
		t.Errorf("FloatFlag.Name() 返回 %q，期望为 %q", f.LongName(), "floatflag")
	}
	if f.ShortName() != "f" {
		t.Errorf("FloatFlag.ShortName() 返回 %q，期望为 %q", f.ShortName(), "f")
	}
	if f.Usage() != "浮点数标志测试" {
		t.Errorf("FloatFlag.Usage() 返回 %q，期望为 %q", f.Usage(), "浮点数标志测试")
	}
	if f.GetDefault() != defValue {
		t.Errorf("FloatFlag.GetDefault() 返回 %v，期望为 %v", f.GetDefault(), defValue)
	}
	if f.Type() != flags.FlagTypeFloat {
		t.Errorf("FloatFlag.Type() 返回 %v，期望为 %v", f.Type(), flags.FlagTypeFloat)
	}

	// 测试边界值
	if err := f.Set(0.0); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if f.Get() != 0.0 {
		t.Errorf("FloatFlag.Get() 返回 %v，期望为 %v", f.Get(), 0.0)
	}

	if err := f.Set(-1.5); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if f.Get() != -1.5 {
		t.Errorf("FloatFlag.Get() 返回 %v，期望为 %v", f.Get(), -1.5)
	}

	if err := f.Set(1.7976931348623157e+308); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if f.Get() != 1.7976931348623157e+308 {
		t.Errorf("FloatFlag.Get() 返回 %v，期望为 %v", f.Get(), 1.7976931348623157e+308)
	}
}

// TestPrintUsage 测试打印用法
func TestPrintUsage(t *testing.T) {
	// 测试自定义用法信息
	cmd1 := NewCmd("test", "t", flag.ExitOnError)
	cmd1.SetUsageSyntax("自定义用法信息")
	if testing.Verbose() {
		cmd1.PrintHelp()
	}

	// 测试自动生成的用法信息
	cmd2 := NewCmd("test2", "t2", flag.ExitOnError)
	cmd2.SetDescription("测试描述")
	cmd2.Bool("verbose", "v", false, "详细输出")
	cmd2.Int("count", "cc", 0, "重复次数")
	if testing.Verbose() {
		cmd2.PrintHelp()
	}

	// 测试带子命令的用法信息
	cmd3 := NewCmd("parent", "0t", flag.ExitOnError)
	subCmd := NewCmd("child", "xd", flag.ExitOnError)
	if err := cmd3.AddSubCmd(subCmd); err != nil {
		t.Errorf("添加子命令时出错: %v", err)
	}
	if testing.Verbose() {
		cmd3.PrintHelp()
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
