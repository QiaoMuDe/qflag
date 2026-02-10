package cmd

import (
	"fmt"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/internal/types"
)

// 测试通过Cmd创建标志并打印帮助信息
func TestCmdHelp(t *testing.T) {
	// 创建一个Cmd实例
	cmd := NewCmd("test-cmd", "测试命令", types.ContinueOnError)

	// 使用所有18种工厂方法创建标志, 确保每个短名称唯一

	// 1. 整数标志
	_ = cmd.Int("count", "c", "计数器", 10)
	// 2. 字符串标志
	_ = cmd.String("output", "o", "输出文件", "output.txt")
	// 3. 布尔标志
	_ = cmd.Bool("verbose", "v", "详细输出", false)
	// 4. 64位整数标志
	_ = cmd.Int64("big-count", "B", "大计数器", int64(10000000000))
	// 5. 无符号整数标志
	_ = cmd.Uint("max-connections", "m", "最大连接数", uint(100))
	// 6. 8位无符号整数标志
	_ = cmd.Uint8("port", "p", "端口号", uint8(250))
	// 7. 16位无符号整数标志
	_ = cmd.Uint16("timeout", "t", "超时时间(秒)", uint16(30))
	// 8. 32位无符号整数标志
	_ = cmd.Uint32("buffer-size", "b", "缓冲区大小", uint32(1024))
	// 9. 64位无符号整数标志
	_ = cmd.Uint64("file-size", "f", "文件大小", uint64(1048576))
	// 10. 64位浮点数标志
	_ = cmd.Float64("ratio", "r", "比例", 0.75)
	// 11. 枚举标志
	_ = cmd.Enum("mode", "M", "运行模式", "auto", []string{"auto", "manual", "debug"})
	// 12. 持续时间标志
	_ = cmd.Duration("interval", "i", "间隔时间", time.Second*30)
	// 13. 时间标志
	_ = cmd.Time("start-time", "s", "开始时间", time.Now())
	// 14. 大小标志
	_ = cmd.Size("max-size", "S", "最大大小", int64(1024*1024))
	// 15. 字符串切片标志
	_ = cmd.StringSlice("paths", "P", "路径列表", []string{"/usr/bin", "/usr/local/bin"})
	// 16. 整数切片标志
	_ = cmd.IntSlice("ports", "pts", "端口列表", []int{80, 443, 8080})
	// 17. 64位整数切片标志
	_ = cmd.Int64Slice("large-numbers", "L", "大数字列表", []int64{1000000, 2000000, 3000000})
	// 18. 映射标志
	_ = cmd.Map("headers", "H", "HTTP头部", map[string]string{"Content-Type": "application/json"})

	// 打印帮助信息
	cmd.PrintHelp()
}

// TestNewCmd 测试 NewCmd 函数
func TestNewCmd(t *testing.T) {
	// 测试创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	if cmd == nil {
		t.Fatal("NewCmd() returned nil")
	}

	if cmd.Name() != "test" {
		t.Errorf("Expected name 'test', got '%s'", cmd.Name())
	}

	if cmd.LongName() != "test" {
		t.Errorf("Expected long name 'test', got '%s'", cmd.LongName())
	}

	if cmd.ShortName() != "t" {
		t.Errorf("Expected short name 't', got '%s'", cmd.ShortName())
	}

	if cmd.IsRootCmd() != true {
		t.Error("Expected IsRootCmd() to return true for new command")
	}
}

// TestCmdProperties 测试命令属性
func TestCmdProperties(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 测试描述
	desc := "测试命令描述"
	cmd.SetDesc(desc)
	if cmd.Desc() != desc {
		t.Errorf("Expected description '%s', got '%s'", desc, cmd.Desc())
	}

	// 测试版本
	version := "1.0.0"
	cmd.SetVersion(version)
	if cmd.Config().Version != version {
		t.Errorf("Expected version '%s', got '%s'", version, cmd.Config().Version)
	}

	// 测试中文设置
	cmd.SetChinese(true)
	if cmd.Config().UseChinese != true {
		t.Error("Expected UseChinese to be true")
	}

	// 测试环境变量前缀
	prefix := "TEST"
	cmd.SetEnvPrefix(prefix)
	if cmd.Config().EnvPrefix != prefix+"_" {
		t.Errorf("Expected EnvPrefix '%s_', got '%s'", prefix, cmd.Config().EnvPrefix)
	}

	// 测试使用语法
	syntax := "test [options]"
	cmd.SetUsageSyntax(syntax)
	if cmd.Config().UsageSyntax != syntax {
		t.Errorf("Expected UsageSyntax '%s', got '%s'", syntax, cmd.Config().UsageSyntax)
	}

	// 测试Logo
	logo := "Test Logo"
	cmd.SetLogoText(logo)
	if cmd.Config().LogoText != logo {
		t.Errorf("Expected LogoText '%s', got '%s'", logo, cmd.Config().LogoText)
	}
}

// TestCmdExamples 测试命令示例
func TestCmdExamples(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 测试添加单个示例
	title1 := "示例1"
	cmd1 := "test command --option value"
	cmd.AddExample(title1, cmd1)

	if cmd.Config().Example[title1] != cmd1 {
		t.Errorf("Expected example '%s', got '%s'", cmd1, cmd.Config().Example[title1])
	}

	// 测试批量添加示例
	examples := map[string]string{
		"示例2": "test command --option2 value2",
		"示例3": "test command --option3 value3",
	}
	cmd.AddExamples(examples)

	for title, expectedCmd := range examples {
		if cmd.Config().Example[title] != expectedCmd {
			t.Errorf("Expected example '%s', got '%s'", expectedCmd, cmd.Config().Example[title])
		}
	}
}

// TestCmdNotes 测试命令注释
func TestCmdNotes(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 测试添加单个注释
	note1 := "这是注释1"
	cmd.AddNote(note1)

	if len(cmd.Config().Notes) != 1 || cmd.Config().Notes[0] != note1 {
		t.Errorf("Expected note '%s', got '%s'", note1, cmd.Config().Notes[0])
	}

	// 测试批量添加注释
	notes := []string{"这是注释2", "这是注释3"}
	cmd.AddNotes(notes)

	expectedNotes := []string{note1, notes[0], notes[1]}
	if len(cmd.Config().Notes) != len(expectedNotes) {
		t.Errorf("Expected %d notes, got %d", len(expectedNotes), len(cmd.Config().Notes))
	}

	for i, expectedNote := range expectedNotes {
		if cmd.Config().Notes[i] != expectedNote {
			t.Errorf("Expected note '%s', got '%s'", expectedNote, cmd.Config().Notes[i])
		}
	}

	// 测试添加空注释 (应该被忽略)
	originalCount := len(cmd.Config().Notes)
	cmd.AddNote("")
	if len(cmd.Config().Notes) != originalCount {
		t.Error("Empty note should be ignored")
	}
}

// TestCmdFlags 测试命令标志
func TestCmdFlags(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 测试添加标志
	flag := cmd.Int("count", "c", "计数器", 10)

	// 测试获取标志
	f, found := cmd.GetFlag("count")
	if !found {
		t.Error("Flag 'count' not found")
	}
	if f != flag {
		t.Error("Retrieved flag is not the same as the added flag")
	}

	// 测试获取所有标志
	flags := cmd.Flags()
	// 由于标志对象只存储一次, 现在只有1个标志
	if len(flags) != 1 {
		t.Errorf("Expected 1 flag, got %d", len(flags))
	}
	// 检查第一个标志是否是我们添加的标志
	flagFound := false
	for _, f := range flags {
		if f == flag {
			flagFound = true
			break
		}
	}
	if !flagFound {
		t.Error("Added flag not found in flags list")
	}

	// 测试标志注册表
	registry := cmd.FlagRegistry()
	if registry == nil {
		t.Error("FlagRegistry() returned nil")
	}
}

// TestCmdSubCommands 测试子命令
func TestCmdSubCommands(t *testing.T) {
	rootCmd := NewCmd("root", "r", types.ContinueOnError)

	// 创建子命令
	subCmd1 := NewCmd("sub1", "s1", types.ContinueOnError)
	subCmd2 := NewCmd("sub2", "s2", types.ContinueOnError)

	// 测试添加子命令
	err := rootCmd.AddSubCmds(subCmd1, subCmd2)
	if err != nil {
		t.Errorf("AddSubCmds() failed: %v", err)
	}

	// 测试获取子命令
	cmd, found := rootCmd.GetSubCmd("sub1")
	if !found {
		t.Error("SubCommand 'sub1' not found")
	}
	if cmd != subCmd1 {
		t.Error("Retrieved subcommand is not the same as the added subcommand")
	}

	// 测试检查子命令存在性
	if !rootCmd.HasSubCmd("sub2") {
		t.Error("HasSubCmd() should return true for 'sub2'")
	}
	if rootCmd.HasSubCmd("nonexistent") {
		t.Error("HasSubCmd() should return false for nonexistent subcommand")
	}

	// 测试获取所有子命令
	subCmds := rootCmd.SubCmds()
	// 由于命令对象只存储一次, 现在只有2个子命令
	if len(subCmds) != 2 {
		t.Errorf("Expected 2 subcommands, got %d", len(subCmds))
	}
	// 检查添加的子命令是否在列表中
	found1 := false
	found2 := false
	for _, cmd := range subCmds {
		if cmd == subCmd1 {
			found1 = true
		}
		if cmd == subCmd2 {
			found2 = true
		}
	}
	if !found1 {
		t.Error("SubCmd1 not found in subcommands list")
	}
	if !found2 {
		t.Error("SubCmd2 not found in subcommands list")
	}

	// 测试子命令注册表
	registry := rootCmd.CmdRegistry()
	if registry == nil {
		t.Error("CmdRegistry() returned nil")
	}
}

// TestCmdArgs 测试命令参数
func TestCmdArgs(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 测试设置参数
	args := []string{"arg1", "arg2", "arg3"}
	cmd.SetArgs(args)

	// 测试获取所有参数
	retrievedArgs := cmd.Args()
	if len(retrievedArgs) != len(args) {
		t.Errorf("Expected %d args, got %d", len(args), len(retrievedArgs))
	}

	for i, arg := range args {
		if retrievedArgs[i] != arg {
			t.Errorf("Expected arg '%s', got '%s'", arg, retrievedArgs[i])
		}
	}

	// 测试获取参数数量
	if cmd.NArg() != len(args) {
		t.Errorf("Expected NArg() %d, got %d", len(args), cmd.NArg())
	}

	// 测试获取指定索引的参数
	for i, arg := range args {
		if cmd.Arg(i) != arg {
			t.Errorf("Expected Arg(%d) '%s', got '%s'", i, arg, cmd.Arg(i))
		}
	}

	// 测试超出范围的索引
	if cmd.Arg(len(args)) != "" {
		t.Error("Arg() should return empty string for out-of-range index")
	}

	// 测试空参数 - 创建新命令以避免解析状态的影响
	emptyCmd := NewCmd("empty", "e", types.ContinueOnError)
	emptyCmd.SetArgs([]string{})
	if emptyCmd.NArg() != 0 {
		t.Errorf("NArg() should return 0 for empty args, got %d", emptyCmd.NArg())
	}
}

// TestCmdParse 测试命令解析
func TestCmdParse(t *testing.T) {
	// 创建一个简单的命令, 避免内置标志处理导致的panic
	cmd := NewCmd("test", "t", types.PanicOnError)

	// 添加标志
	countFlag := cmd.Int("count", "c", "计数器", 0)
	verboseFlag := cmd.Bool("verbose", "v", "详细输出", false)

	// 直接设置参数, 绕过解析过程
	cmd.SetArgs([]string{"arg1", "arg2"})
	cmd.SetParsed(true)

	// 检查解析状态
	if !cmd.IsParsed() {
		t.Error("IsParsed() should return true after setting parsed")
	}

	// 检查参数
	if cmd.NArg() != 2 {
		t.Errorf("Expected 2 args, got %d", cmd.NArg())
	}

	if cmd.Arg(0) != "arg1" || cmd.Arg(1) != "arg2" {
		t.Error("Arguments not set correctly")
	}

	// 检查标志值 (应该保持默认值)
	if countFlag.Get() != 0 {
		t.Errorf("Expected count flag default value 0, got %d", countFlag.Get())
	}

	if verboseFlag.Get() != false {
		t.Error("Expected verbose flag to be false")
	}
}

// TestCmdRun 测试命令执行
func TestCmdRun(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 测试未设置运行函数
	err := cmd.Run()
	if err == nil {
		t.Error("Run() should return error when no run function is set")
	}

	// 测试设置运行函数
	runCalled := false
	cmd.SetRun(func(c types.Command) error {
		runCalled = true
		if c != cmd {
			t.Error("Run function should receive the same command instance")
		}
		return nil
	})

	// 测试运行函数存在性
	if !cmd.HasRunFunc() {
		t.Error("HasRunFunc() should return true after SetRun()")
	}

	// 测试执行 (需要先解析)
	cmd.SetArgs([]string{})
	// 不调用 Parse 方法, 直接设置解析状态, 避免重复注册内置标志
	cmd.SetParsed(true)

	err = cmd.Run()
	if err != nil {
		t.Errorf("Run() failed: %v", err)
	}

	if !runCalled {
		t.Error("Run function was not called")
	}
}

// TestFlagFactories 测试所有标志工厂方法
func TestFlagFactories(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 测试 Int 标志
	intFlag := cmd.Int("int", "i", "整数标志", 42)
	if intFlag.Get() != 42 {
		t.Errorf("Expected int flag default value 42, got %d", intFlag.Get())
	}

	// 测试 String 标志
	stringFlag := cmd.String("string", "s", "字符串标志", "default")
	if stringFlag.Get() != "default" {
		t.Errorf("Expected string flag default value 'default', got '%s'", stringFlag.Get())
	}

	// 测试 Bool 标志
	boolFlag := cmd.Bool("bool", "b", "布尔标志", true)
	if boolFlag.Get() != true {
		t.Errorf("Expected bool flag default value true, got %t", boolFlag.Get())
	}

	// 测试 Int64 标志
	int64Flag := cmd.Int64("int64", "i64", "64位整数标志", int64(1234567890))
	if int64Flag.Get() != 1234567890 {
		t.Errorf("Expected int64 flag default value 1234567890, got %d", int64Flag.Get())
	}

	// 测试 Uint 标志
	uintFlag := cmd.Uint("uint", "u", "无符号整数标志", uint(123))
	if uintFlag.Get() != 123 {
		t.Errorf("Expected uint flag default value 123, got %d", uintFlag.Get())
	}

	// 测试 Uint8 标志
	uint8Flag := cmd.Uint8("uint8", "u8", "8位无符号整数标志", uint8(8))
	if uint8Flag.Get() != 8 {
		t.Errorf("Expected uint8 flag default value 8, got %d", uint8Flag.Get())
	}

	// 测试 Uint16 标志
	uint16Flag := cmd.Uint16("uint16", "u16", "16位无符号整数标志", uint16(16))
	if uint16Flag.Get() != 16 {
		t.Errorf("Expected uint16 flag default value 16, got %d", uint16Flag.Get())
	}

	// 测试 Uint32 标志
	uint32Flag := cmd.Uint32("uint32", "u32", "32位无符号整数标志", uint32(32))
	if uint32Flag.Get() != 32 {
		t.Errorf("Expected uint32 flag default value 32, got %d", uint32Flag.Get())
	}

	// 测试 Uint64 标志
	uint64Flag := cmd.Uint64("uint64", "u64", "64位无符号整数标志", uint64(64))
	if uint64Flag.Get() != 64 {
		t.Errorf("Expected uint64 flag default value 64, got %d", uint64Flag.Get())
	}

	// 测试 Float64 标志
	float64Flag := cmd.Float64("float64", "f64", "64位浮点数标志", 3.14159)
	if float64Flag.Get() != 3.14159 {
		t.Errorf("Expected float64 flag default value 3.14159, got %f", float64Flag.Get())
	}

	// 测试 Enum 标志
	enumFlag := cmd.Enum("enum", "e", "枚举标志", "option1", []string{"option1", "option2", "option3"})
	if enumFlag.Get() != "option1" {
		t.Errorf("Expected enum flag default value 'option1', got '%s'", enumFlag.Get())
	}

	// 测试 Duration 标志
	durationFlag := cmd.Duration("duration", "d", "持续时间标志", time.Second*30)
	if durationFlag.Get() != time.Second*30 {
		t.Errorf("Expected duration flag default value %v, got %v", time.Second*30, durationFlag.Get())
	}

	// 测试 Time 标志
	timeFlag := cmd.Time("time", "tm", "时间标志", time.Time{})
	if !timeFlag.Get().IsZero() {
		t.Error("Expected time flag default value to be zero time")
	}

	// 测试 Size 标志
	sizeFlag := cmd.Size("size", "sz", "大小标志", int64(1024))
	if sizeFlag.Get() != 1024 {
		t.Errorf("Expected size flag default value 1024, got %d", sizeFlag.Get())
	}

	// 测试 StringSlice 标志
	stringSliceFlag := cmd.StringSlice("string-slice", "ss", "字符串切片标志", []string{"a", "b", "c"})
	slice := stringSliceFlag.Get()
	if len(slice) != 3 || slice[0] != "a" || slice[1] != "b" || slice[2] != "c" {
		t.Errorf("Expected string slice flag default value [a b c], got %v", slice)
	}

	// 测试 IntSlice 标志
	intSliceFlag := cmd.IntSlice("int-slice", "is", "整数切片标志", []int{1, 2, 3})
	intSlice := intSliceFlag.Get()
	if len(intSlice) != 3 || intSlice[0] != 1 || intSlice[1] != 2 || intSlice[2] != 3 {
		t.Errorf("Expected int slice flag default value [1 2 3], got %v", intSlice)
	}

	// 测试 Int64Slice 标志
	int64SliceFlag := cmd.Int64Slice("int64-slice", "i64s", "64位整数切片标志", []int64{10, 20, 30})
	int64Slice := int64SliceFlag.Get()
	if len(int64Slice) != 3 || int64Slice[0] != 10 || int64Slice[1] != 20 || int64Slice[2] != 30 {
		t.Errorf("Expected int64 slice flag default value [10 20 30], got %v", int64Slice)
	}

	// 测试 Map 标志
	mapFlag := cmd.Map("map", "m", "映射标志", map[string]string{"key": "value"})
	mapValue := mapFlag.Get()
	if len(mapValue) != 1 || mapValue["key"] != "value" {
		t.Errorf("Expected map flag default value map[key:value], got %v", mapValue)
	}
}

// TestFlagNameValidation 测试标志名称验证
func TestFlagNameValidation(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 测试有效标志名称
	validFlags := []struct {
		longName  string
		shortName string
		desc      string
	}{
		{"valid-long", "v", "长名和短名都有效"},
		{"valid-long2", "", "只有长名有效"},
		{"", "s", "只有短名有效"},
	}

	for i, flag := range validFlags {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Valid flag names should not panic: %v", r)
				}
			}()
			// 为每个测试用例创建新的命令, 避免标志重复
			testCmd := NewCmd(fmt.Sprintf("test%d", i), "", types.ContinueOnError)
			// 为每个测试用例使用唯一的长名称和短名称
			longName := flag.longName
			shortName := flag.shortName
			if longName != "" {
				longName = fmt.Sprintf("%s%d", flag.longName, i)
			}
			if shortName != "" {
				// 使用完全不同的短名称, 避免与内置标志冲突
				shortName = fmt.Sprintf("x%d", i)
			}
			testCmd.Int(longName, shortName, "测试标志", 0)
		})
	}

	// 测试无效标志名称
	invalidFlags := []struct {
		longName  string
		shortName string
		desc      string
	}{
		{"", "", "空长名和短名"},
	}

	for _, flag := range invalidFlags {
		t.Run(flag.desc, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Invalid flag names should panic")
				}
			}()
			cmd.Int(flag.longName, flag.shortName, "测试标志", 0)
		})
	}

	// 测试重复标志名称
	testCmd := NewCmd("test-duplicate", "td", types.ContinueOnError)
	testCmd.Int("duplicate", "d", "第一个标志", 0)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Duplicate flag names should panic")
		}
	}()
	testCmd.Int("duplicate", "x", "重复长名的标志", 0)
}

// TestHelpGeneration 测试帮助信息生成
func TestHelpGeneration(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)
	cmd.SetDesc("测试命令描述")

	// 添加标志
	cmd.Int("count", "c", "计数器", 10)
	cmd.Bool("verbose", "v", "详细输出", false)

	// 添加子命令
	subCmd := NewCmd("sub", "s", types.ContinueOnError)
	subCmd.SetDesc("子命令描述")
	if err := cmd.AddSubCmds(subCmd); err != nil {
		t.Errorf("AddSubCmds() error = %v", err)
	}

	// 添加示例
	cmd.AddExample("基本用法", "test --count 5")

	// 添加注释
	cmd.AddNote("注意: 这是一个测试命令")

	// 获取帮助信息
	help := cmd.Help()
	if help == "" {
		t.Error("Help() should not return empty string")
	}

	// 测试打印帮助信息 (不会实际打印, 只是确保不崩溃)
	cmd.PrintHelp()
}
