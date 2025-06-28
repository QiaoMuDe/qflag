package qflag

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
	"time"

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
		t.Error("String() 返回了 nil")
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
		t.Error("String() 返回了 nil")
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
		t.Error("Int() 返回了 nil")
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
		t.Error("Int() 返回了 nil")
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
		t.Error("Bool() 返回了 nil")
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
		t.Error("Bool() 返回了 nil")
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
		t.Error("Float() 返回了 nil")
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
		t.Error("Float() 返回了 nil")
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
			buf.ReadFrom(r)
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

// TestCmd_Description 测试Cmd的Description和SetDescription方法
func TestCmd_Description(t *testing.T) {
	cmd := &Cmd{}
	desc := "测试描述"
	cmd.SetDescription(desc)
	if cmd.GetDescription() != desc {
		t.Errorf("Description() 返回 %q，期望为 %q", cmd.GetDescription(), desc)
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

// TestCmd_Args 测试Cmd的Args方法
func TestCmd_Args(t *testing.T) {
	args := []string{"arg1", "arg2"}
	cmd := &Cmd{args: args}
	result := cmd.Args()
	// 检查长度
	if len(result) != len(args) {
		t.Fatalf("Args() 返回的长度为 %d，期望为 %d", len(result), len(args))
	}
	// 检查每个元素
	for i, arg := range args {
		if result[i] != arg {
			t.Errorf("Args()[%d] 返回 %q，期望为 %q", i, result[i], arg)
		}
	}
}

// TestCmd_Arg 测试Cmd的Arg方法
func TestCmd_Arg(t *testing.T) {
	cmd := &Cmd{args: []string{"arg0", "arg1", "arg2"}}
	tests := []struct {
		name string
		i    int
		want string
	}{{
		name: "有效索引 0",
		i:    0,
		want: "arg0",
	}, {
		name: "有效索引 1",
		i:    1,
		want: "arg1",
	}, {
		name: "索引越界",
		i:    3,
		want: "",
	}, {
		name: "负索引",
		i:    -1,
		want: "",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmd.Arg(tt.i); got != tt.want {
				t.Errorf("Arg(%d) 返回 %q，期望为 %q", tt.i, got, tt.want)
			}
		})
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
	f.Set(0)
	if f.Get() != 0 {
		t.Errorf("IntFlag.Get() 返回 %v，期望为 %v", f.Get(), 0)
	}

	f.Set(-1)
	if f.Get() != -1 {
		t.Errorf("IntFlag.Get() 返回 %v，期望为 %v", f.Get(), -1)
	}

	f.Set(2147483647)
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
	f.Set("")
	if f.Get() != "" {
		t.Errorf("StringFlag.Get() 返回 %q，期望为 %q", f.Get(), "")
	}

	f.Set("long_string_with_special_chars_!@#$%^&*()")
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
	f.Set(false)
	if f.Get() != false {
		t.Errorf("BoolFlag.Get() 返回 %v，期望为 %v", f.Get(), false)
	}

	f.Set(true)
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
	f.Set(0.0)
	if f.Get() != 0.0 {
		t.Errorf("FloatFlag.Get() 返回 %v，期望为 %v", f.Get(), 0.0)
	}

	f.Set(-1.5)
	if f.Get() != -1.5 {
		t.Errorf("FloatFlag.Get() 返回 %v，期望为 %v", f.Get(), -1.5)
	}

	f.Set(1.7976931348623157e+308)
	if f.Get() != 1.7976931348623157e+308 {
		t.Errorf("FloatFlag.Get() 返回 %v，期望为 %v", f.Get(), 1.7976931348623157e+308)
	}
}

// TestBindHelpFlag 测试绑定帮助标志
func TestBindHelpFlag(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ExitOnError)
	cmd.initBuiltinFlags()
	// 验证帮助标志已绑定
	if !cmd.initFlagBound {
		t.Error("帮助标志应该已绑定")
	}
	if _, ok := cmd.flagRegistry.GetByName(flags.HelpFlagName); !ok {
		t.Error("帮助标志应该已注册")
	}

	// 当短帮助标志名存在时，检查该标志是否已注册，若未注册则报错。
	_, ok := cmd.flagRegistry.GetByName(flags.HelpFlagShortName)
	if flags.HelpFlagShortName != "" && !ok {
		t.Error("短帮助标志应该已注册")
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

// TestHasCycle 测试检测循环引用
func TestHasCycle(t *testing.T) {
	cmd1 := NewCmd("cmd1", "", flag.ExitOnError)
	cmd2 := NewCmd("", "c2", flag.ExitOnError)
	cmd3 := NewCmd("cmd3", "c3", flag.ExitOnError)
	cmd4 := NewCmd("", "c4", flag.ExitOnError)

	// 无循环情况
	if cmd1.hasCycle(cmd2) {
		t.Error("初始时不应存在循环引用")
	}

	// 添加子命令
	if err := cmd1.AddSubCmd(cmd2); err != nil {
		t.Errorf("添加子命令时出错: %v", err)
	}
	cmd2.parentCmd = cmd1
	if err := cmd2.AddSubCmd(cmd3); err != nil {
		t.Errorf("添加子命令时出错: %v", err)
	}
	cmd3.parentCmd = cmd2

	// 检测循环
	if cmd1.hasCycle(cmd4) {
		t.Error("与不相关的命令不应存在循环引用")
	}
	if !cmd1.hasCycle(cmd1) { // 自引用
		t.Error("应检测到自循环引用")
	}
	if !cmd2.hasCycle(cmd1) { // 反向引用
		t.Error("应检测到反向循环引用")
	}
	if !cmd3.hasCycle(cmd1) { // 多级反向引用
		t.Error("应检测到多级反向循环引用")
	}
}

// TestNestedCmdHelp 测试嵌套子命令的帮助信息生成
func TestNestedCmdHelp(t *testing.T) {
	// 创建三级嵌套命令结构
	cmd1 := NewCmd("cmd1", "", flag.ExitOnError)
	cmd1.SetDescription("一级命令描述")
	cmd1.String("config", "c", "config.json", "配置文件路径")

	cmd2 := NewCmd("", "c2", flag.ExitOnError)
	cmd2.SetDescription("二级命令描述")
	cmd2.Int("port", "p", 8080, "服务端口号")

	cmd3 := NewCmd("cmd3", "", flag.ExitOnError)
	cmd3.SetDescription("三级命令描述")
	cmd3.Bool("verbose", "", false, "详细输出模式")
	cmd3.SetUseChinese(true)
	cmd2.SetUseChinese(true)
	cmd3.String("output", "o", "", "输出文件路径")
	cmd3.Float("timeout", "t", 5.0, "超时时间")
	cmd3.Duration("duration", "d", 10*time.Second, "持续时间")
	cmd3.Enum("format", "f", "json", "输出格式", []string{"json", "xml", "yaml"})

	cmd4 := NewCmd("ssssssscmd4", "ccccc4", flag.ExitOnError)
	cmd4.SetDescription("四级命令描述")

	cmd5 := NewCmd("acmd5", "ccccc5", flag.ExitOnError)
	cmd5.SetDescription("五级命令描述")

	// 新增子命令用于测试帮助信息生成
	cmd6 := NewCmd("randomizer", "rz", flag.ExitOnError)
	cmd6.SetDescription("新增六级命令描述")
	cmd6.Float("timeout", "t", 5.0, "超时时间")

	cmd7 := NewCmd("generator", "gn", flag.ExitOnError)
	cmd7.SetDescription("新增七级命令描述")
	cmd7.String("format", "f", "json", "输出格式")

	cmd8 := NewCmd("processor", "ps", flag.ExitOnError)
	cmd8.SetDescription("新增八级命令描述")
	cmd8.Int("retry", "r", 3, "重试次数")

	// 添加示例
	cmd3.AddExample(ExampleInfo{"示例1", "echo 111"})
	cmd3.AddExample(ExampleInfo{"示例2", "echo 222"})

	// 构建命令层级
	cmd1.AddSubCmd(cmd2)
	cmd2.AddSubCmd(cmd3)
	cmd2.AddSubCmd(cmd4, cmd5)
	cmd3.AddSubCmd(cmd6, cmd7, cmd8)

	// 添加注意事项
	cmd1.AddNote("注意事项1")
	cmd1.AddNote("注意事项2")
	cmd2.AddNote("注意事项3")
	cmd3.AddNote("注意事项4")

	// 解析命令行参数
	if err := cmd1.Parse([]string{}); err != nil {
		t.Errorf("解析命令行参数时出错: %v", err)
	}

	// 测试帮助信息生成
	// 使用t.Log()替代fmt.Println()，并添加testing.Verbose()条件控制
	printSection := func(section string) {
		if testing.Verbose() {
			t.Log(section)
		}
	}

	printSeparator := func() {
		if testing.Verbose() {
			t.Log("========================")
		}
	}

	printUsage := func(cmd *Cmd) {
		if testing.Verbose() {
			// 重定向cmd.PrintUsage()输出到t.Log
			o := cmd.fs.Output()
			cmd.fs.SetOutput(&testLogWriter{t: t})
			cmd.PrintHelp()
			cmd.fs.SetOutput(o)
		}
	}

	// 一级命令帮助信息
	printSection("=== 一级命令帮助信息 ===")
	printUsage(cmd1)
	printSeparator()

	// 二级命令帮助信息
	printSection("=== 二级命令帮助信息 ===")
	printUsage(cmd2)
	printSeparator()

	// 三级命令帮助信息
	printSection("=== 三级命令帮助信息 ===")
	printUsage(cmd3)
	printSeparator()
}

// testLogWriter 用于将flag.FlagSet的输出重定向到testing.T的Log方法
type testLogWriter struct {
	t *testing.T
}

func (w *testLogWriter) Write(p []byte) (n int, err error) {
	if testing.Verbose() {
		w.t.Log(string(p))
	}
	return len(p), nil
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

		// 跳过内置标志绑定，避免冲突
		cmd.initFlagBound = true

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
