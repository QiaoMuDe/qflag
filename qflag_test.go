package qflag

import (
	"bytes"
	"flag"
	"os"
	"testing"
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
			t.Errorf("ReadFrom error: %v", err)
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "string-flag"
	defValue := "default"
	usage := "test string flag"

	// 测试String方法(仅长标志)
	f := cmd.String(flagName, "sf", defValue, usage)
	if f == nil {
		t.Error("String() returned nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName, "test-value"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != "test-value" {
		t.Errorf("String flag value = %q, want %q", *f.value, "test-value")
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
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "s"
	defValue := "default"
	usage := "test string flag"

	// 测试String方法(仅短标志)
	f := cmd.String("sf", shortName, defValue, usage)
	if f == nil {
		t.Error("String() returned nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName, "test-value"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != "test-value" {
		t.Errorf("String flag value = %q, want %q", *f.value, "test-value")
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
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "int-flag"
	defValue := 100
	usage := "test int flag"

	// 测试Int方法(仅长标志)
	f := cmd.Int(flagName, "if", defValue, usage)
	if f == nil {
		t.Error("Int() returned nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName, "200"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != 200 {
		t.Errorf("Int flag value = %d, want %d", *f.value, 200)
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
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "i"
	defValue := 100
	usage := "test int flag"

	// 测试Int方法(仅短标志)
	f := cmd.Int("ci", shortName, defValue, usage)
	if f == nil {
		t.Error("Int() returned nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName, "200"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != 200 {
		t.Errorf("Int flag value = %d, want %d", *f.value, 200)
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
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "bool-flag"
	defValue := false
	usage := "test bool flag"

	// 测试Bool方法(仅长标志)
	f := cmd.Bool(flagName, "bl", defValue, usage)
	if f == nil {
		t.Error("Bool() returned nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != true {
		t.Errorf("Bool flag value = %v, want %v", *f.value, true)
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
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "b"
	defValue := false
	usage := "test bool flag"

	// 测试Bool方法(仅短标志)
	f := cmd.Bool("ct", shortName, defValue, usage)
	if f == nil {
		t.Error("Bool() returned nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != true {
		t.Errorf("Bool flag value = %v, want %v", *f.value, true)
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
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "float-flag"
	defValue := 3.14
	usage := "test float flag"

	// 测试Float方法(仅长标志)
	f := cmd.Float(flagName, "ff", defValue, usage)
	if f == nil {
		t.Error("Float() returned nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName, "6.28"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != 6.28 {
		t.Errorf("Float flag value = %v, want %v", *f.value, 6.28)
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
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "f"
	defValue := 3.14
	usage := "test float flag"

	// 测试Float方法(仅短标志)
	f := cmd.Float("cf", shortName, defValue, usage)
	if f == nil {
		t.Error("Float() returned nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName, "6.28"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != 6.28 {
		t.Errorf("Float flag value = %v, want %v", *f.value, 6.28)
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
				t.Errorf("ReadFrom stdout error: %v", err)
			}
			if _, err := errBuf.ReadFrom(rErr); err != nil {
				t.Errorf("ReadFrom stderr error: %v", err)
			}
			t.Logf("Captured stdout:\n%s", outBuf.String())
			t.Logf("Captured stderr:\n%s", errBuf.String())
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.Int("int-flag", "i", 0, "test int flag")

	// 测试无效参数
	err := cmd.Parse([]string{"--int-flag", "not-a-number"})
	if err == nil {
		t.Error("Parse() should return error for invalid input")
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
			t.Logf("Captured output:\n%s", buf.String())
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.String("string-flag", "s", "", "test string flag")

	// 测试帮助标志
	err := cmd.Parse([]string{"--help"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}
}

// TestCmd_Name 测试Cmd的Name方法
func TestCmd_Name(t *testing.T) {
	cmd := &Cmd{longName: "testcmd"}
	if cmd.LongName() != "testcmd" {
		t.Errorf("Name() = %q, want %q", cmd.LongName(), "testcmd")
	}
}

// TestCmd_ShortName 测试Cmd的ShortName方法
func TestCmd_ShortName(t *testing.T) {
	cmd := &Cmd{shortName: "tc"}
	if cmd.ShortName() != "tc" {
		t.Errorf("ShortName() = %q, want %q", cmd.ShortName(), "tc")
	}
}

// TestCmd_Description 测试Cmd的Description和SetDescription方法
func TestCmd_Description(t *testing.T) {
	cmd := &Cmd{}
	desc := "test description"
	cmd.SetDescription(desc)
	if cmd.Description() != desc {
		t.Errorf("Description() = %q, want %q", cmd.Description(), desc)
	}
}

// TestCmd_Usage 测试Cmd的Usage和SetUsage方法
func TestCmd_Usage(t *testing.T) {
	cmd := &Cmd{}
	usage := "test usage"
	cmd.SetUsage(usage)
	if cmd.Usage() != usage {
		t.Errorf("Usage() = %q, want %q", cmd.Usage(), usage)
	}
}

// TestCmd_Args 测试Cmd的Args方法
func TestCmd_Args(t *testing.T) {
	args := []string{"arg1", "arg2"}
	cmd := &Cmd{args: args}
	result := cmd.Args()
	// 检查长度
	if len(result) != len(args) {
		t.Fatalf("Args() length = %d, want %d", len(result), len(args))
	}
	// 检查每个元素
	for i, arg := range args {
		if result[i] != arg {
			t.Errorf("Args()[%d] = %q, want %q", i, result[i], arg)
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
		name: "valid index 0",
		i:    0,
		want: "arg0",
	}, {
		name: "valid index 1",
		i:    1,
		want: "arg1",
	}, {
		name: "index out of range",
		i:    3,
		want: "",
	}, {
		name: "negative index",
		i:    -1,
		want: "",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmd.Arg(tt.i); got != tt.want {
				t.Errorf("Arg(%d) = %q, want %q", tt.i, got, tt.want)
			}
		})
	}
}

// TestIntFlag_Interface 验证IntFlag实现了Flag接口
func TestIntFlag_Interface(t *testing.T) {
	var f Flag = &IntFlag{}
	_ = f // 若编译通过，则说明实现了接口
}

// TestIntFlag_Methods 测试IntFlag的各种方法
func TestIntFlag_Methods(t *testing.T) {
	defValue := 100

	// 新建子命令
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	f := cmd.Int("intflag", "i", defValue, "integer flag test")

	if f.longName != "intflag" {
		t.Errorf("IntFlag.Name() = %q, want %q", f.longName, "intflag")
	}
	if f.ShortName() != "i" {
		t.Errorf("IntFlag.ShortName() = %q, want %q", f.ShortName(), "i")
	}
	if f.Usage() != "integer flag test" {
		t.Errorf("IntFlag.Usage() = %q, want %q", f.Usage(), "integer flag test")
	}
	if f.GetDefault() != defValue {
		t.Errorf("IntFlag.GetDefault() = %v, want %v", f.GetDefault(), defValue)
	}
	if f.Type() != FlagTypeInt {
		t.Errorf("IntFlag.Type() = %v, want %v", f.Type(), FlagTypeInt)
	}

	// 测试边界值
	f.Set(0)
	if f.Get() != 0 {
		t.Errorf("IntFlag.Get() = %v, want %v", f.Get(), 0)
	}

	f.Set(-1)
	if f.Get() != -1 {
		t.Errorf("IntFlag.Get() = %v, want %v", f.Get(), -1)
	}

	f.Set(2147483647)
	if f.Get() != 2147483647 {
		t.Errorf("IntFlag.Get() = %v, want %v", f.Get(), 2147483647)
	}
}

// TestStringFlag_Interface 验证StringFlag实现了Flag接口
func TestStringFlag_Interface(t *testing.T) {
	var f Flag = &StringFlag{}
	_ = f
}

// TestStringFlag_Methods 测试StringFlag的各种方法
func TestStringFlag_Methods(t *testing.T) {
	defValue := "default string"

	// 新建子命令
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	f := cmd.String("strflag", "s", defValue, "string flag test")

	if f.longName != "strflag" {
		t.Errorf("StringFlag.Name() = %q, want %q", f.longName, "strflag")
	}
	if f.ShortName() != "s" {
		t.Errorf("StringFlag.ShortName() = %q, want %q", f.ShortName(), "s")
	}
	if f.Usage() != "string flag test" {
		t.Errorf("StringFlag.Usage() = %q, want %q", f.Usage(), "string flag test")
	}
	if f.GetDefault() != defValue {
		t.Errorf("StringFlag.GetDefault() = %v, want %v", f.GetDefault(), defValue)
	}
	if f.Type() != FlagTypeString {
		t.Errorf("StringFlag.Type() = %v, want %v", f.Type(), FlagTypeString)
	}

	// 测试边界值
	f.Set("")
	if f.Get() != "" {
		t.Errorf("StringFlag.Get() = %q, want %q", f.Get(), "")
	}

	f.Set("long_string_with_special_chars_!@#$%^&*()")
	if f.Get() != "long_string_with_special_chars_!@#$%^&*()" {
		t.Errorf("StringFlag.Get() = %q, want %q", f.Get(), "long_string_with_special_chars_!@#$%^&*()")
	}
}

// TestBoolFlag_Interface 验证BoolFlag实现了Flag接口
func TestBoolFlag_Interface(t *testing.T) {
	var f Flag = &BoolFlag{}
	_ = f
}

// TestBoolFlag_Methods 测试BoolFlag的各种方法
func TestBoolFlag_Methods(t *testing.T) {
	defValue := true

	// 新建子命令
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	f := cmd.Bool("boolflag", "b", defValue, "bool flag test")

	if f.longName != "boolflag" {
		t.Errorf("BoolFlag.Name() = %q, want %q", f.longName, "boolflag")
	}
	if f.ShortName() != "b" {
		t.Errorf("BoolFlag.ShortName() = %q, want %q", f.ShortName(), "b")
	}
	if f.Usage() != "bool flag test" {
		t.Errorf("BoolFlag.Usage() = %q, want %q", f.Usage(), "bool flag test")
	}
	if f.GetDefault() != defValue {
		t.Errorf("BoolFlag.GetDefault() = %v, want %v", f.GetDefault(), defValue)
	}
	if f.Type() != FlagTypeBool {
		t.Errorf("BoolFlag.Type() = %v, want %v", f.Type(), FlagTypeBool)
	}

	// 测试切换布尔值
	f.Set(false)
	if f.Get() != false {
		t.Errorf("BoolFlag.Get() = %v, want %v", f.Get(), false)
	}

	f.Set(true)
	if f.Get() != true {
		t.Errorf("BoolFlag.Get() = %v, want %v", f.Get(), true)
	}
}

// TestFloatFlag_Interface 验证FloatFlag实现了Flag接口
func TestFloatFlag_Interface(t *testing.T) {
	var f Flag = &FloatFlag{}
	_ = f
}

// TestFloatFlag_Methods 测试FloatFlag的各种方法
func TestFloatFlag_Methods(t *testing.T) {
	defValue := 3.14

	// 新建子命令
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	f := cmd.Float("floatflag", "f", defValue, "float flag test")

	if f.longName != "floatflag" {
		t.Errorf("FloatFlag.Name() = %q, want %q", f.longName, "floatflag")
	}
	if f.ShortName() != "f" {
		t.Errorf("FloatFlag.ShortName() = %q, want %q", f.ShortName(), "f")
	}
	if f.Usage() != "float flag test" {
		t.Errorf("FloatFlag.Usage() = %q, want %q", f.Usage(), "float flag test")
	}
	if f.GetDefault() != defValue {
		t.Errorf("FloatFlag.GetDefault() = %v, want %v", f.GetDefault(), defValue)
	}
	if f.Type() != FlagTypeFloat {
		t.Errorf("FloatFlag.Type() = %v, want %v", f.Type(), FlagTypeFloat)
	}

	// 测试边界值
	f.Set(0.0)
	if f.Get() != 0.0 {
		t.Errorf("FloatFlag.Get() = %v, want %v", f.Get(), 0.0)
	}

	f.Set(-1.5)
	if f.Get() != -1.5 {
		t.Errorf("FloatFlag.Get() = %v, want %v", f.Get(), -1.5)
	}

	f.Set(1.7976931348623157e+308)
	if f.Get() != 1.7976931348623157e+308 {
		t.Errorf("FloatFlag.Get() = %v, want %v", f.Get(), 1.7976931348623157e+308)
	}
}

// TestBindHelpFlag 测试绑定帮助标志
func TestBindHelpFlag(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ExitOnError)
	cmd.initBuiltinFlags()
	// 验证帮助标志已绑定
	if !cmd.initFlagBound {
		t.Error("help flag should be bound")
	}
	if _, ok := cmd.flagRegistry.GetByName(helpFlagName); !ok {
		t.Error("help flag should be registered")
	}

	// 当短帮助标志名存在时，检查该标志是否已注册，若未注册则报错。
	_, ok := cmd.flagRegistry.GetByName(helpFlagShortName)
	if helpFlagShortName != "" && !ok {
		t.Error("short help flag should be registered")
	}
}

// TestPrintUsage 测试打印用法
func TestPrintUsage(t *testing.T) {
	// 测试自定义用法信息
	cmd1 := NewCmd("test", "t", flag.ExitOnError)
	cmd1.SetUsage("Custom usage message")
	if testing.Verbose() {
		cmd1.PrintUsage()
	}

	// 测试自动生成的用法信息
	cmd2 := NewCmd("test2", "t2", flag.ExitOnError)
	cmd2.SetDescription("Test description")
	cmd2.Bool("verbose", "v", false, "verbose output")
	cmd2.Int("count", "cc", 0, "number of times to repeat")
	if testing.Verbose() {
		cmd2.PrintUsage()
	}

	// 测试带子命令的用法信息
	cmd3 := NewCmd("parent", "0t", flag.ExitOnError)
	subCmd := NewCmd("child", "xd", flag.ExitOnError)
	if err := cmd3.AddSubCmd(subCmd); err != nil {
		t.Errorf("AddSubCmd error: %v", err)
	}
	if testing.Verbose() {
		cmd3.PrintUsage()
	}
}

// TestHasCycle 测试检测循环引用
func TestHasCycle(t *testing.T) {
	cmd1 := NewCmd("cmd1", "c1", flag.ExitOnError)
	cmd2 := NewCmd("cmd2", "c2", flag.ExitOnError)
	cmd3 := NewCmd("cmd3", "c3", flag.ExitOnError)
	cmd4 := NewCmd("cmd4", "c4", flag.ExitOnError)

	// 无循环情况
	if hasCycle(cmd1, cmd2) {
		t.Error("should not have cycle initially")
	}

	// 添加子命令
	if err := cmd1.AddSubCmd(cmd2); err != nil {
		t.Errorf("AddSubCmd error: %v", err)
	}
	cmd2.parentCmd = cmd1
	if err := cmd2.AddSubCmd(cmd3); err != nil {
		t.Errorf("AddSubCmd error: %v", err)
	}
	cmd3.parentCmd = cmd2

	// 检测循环
	if hasCycle(cmd1, cmd4) {
		t.Error("should not have cycle with unrelated cmd")
	}
	if !hasCycle(cmd1, cmd1) { // 自引用
		t.Error("should detect self cycle")
	}
	if !hasCycle(cmd2, cmd1) { // 反向引用
		t.Error("should detect reverse cycle")
	}
	if !hasCycle(cmd3, cmd1) { // 多级反向引用
		t.Error("should detect multi-level reverse cycle")
	}
}

// TestNestedCmdHelp 测试嵌套子命令的帮助信息生成
func TestNestedCmdHelp(t *testing.T) {
	// 创建三级嵌套命令结构
	cmd1 := NewCmd("cmd1", "c1", flag.ExitOnError)
	cmd1.SetDescription("一级命令描述")
	cmd1.String("config", "c", "config.json", "配置文件路径")

	cmd2 := NewCmd("cmd2", "c2", flag.ExitOnError)
	cmd2.SetDescription("二级命令描述")
	cmd2.Int("port", "p", 8080, "服务端口号")

	cmd3 := NewCmd("cmd3", "c3", flag.ExitOnError)
	cmd3.SetDescription("三级命令描述")
	cmd3.Bool("verbose", "v", false, "详细输出模式")

	// 构建命令层级
	cmd1.AddSubCmd(cmd2)
	cmd2.AddSubCmd(cmd3)

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
			cmd.PrintUsage()
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
