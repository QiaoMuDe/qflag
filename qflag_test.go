package qflag

import (
	"flag"
	"os"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/internal/cmd"
	"gitee.com/MM-Q/qflag/internal/parser"
	"gitee.com/MM-Q/qflag/internal/types"
)

// createTestCmd 创建一个测试命令
//
// 返回值:
//   - *cmd.Cmd: 测试命令实例
func createTestCmd() *cmd.Cmd {
	testCmd := cmd.NewCmd("test", "t", types.ContinueOnError)
	testCmd.SetDesc("测试命令")
	testCmd.SetChinese(true)
	testCmd.AddExamples(map[string]string{
		"基本用法": "test -c 10 -o output.txt",
		"详细模式": "test -v -c 20 -o debug.log",
	})
	testCmd.AddNotes([]string{
		"注意: ",
		"1. 计数器和最大计数必须为非负数。",
		"2. 输出文件和配置文件路径可以是绝对路径或相对路径。",
	})

	_ = testCmd.Int("count", "c", "计数器", 10)
	_ = testCmd.Int("max-count", "m", "最大计数", 100)

	_ = testCmd.String("output", "o", "输出文件", "output.txt")
	_ = testCmd.String("config", "C", "配置文件", "config.yaml")

	_ = testCmd.Bool("verbose", "V", "详细输出", false)
	_ = testCmd.Bool("debug", "D", "调试模式", false)

	_ = testCmd.Int64("big-count", "B", "大计数器", int64(10000000000))
	_ = testCmd.Int64("total", "T", "总数", int64(9999999999))

	_ = testCmd.Uint("max-connections", "U", "最大连接数", uint(100))
	_ = testCmd.Uint("min-connections", "N", "最小连接数", uint(1))

	_ = testCmd.Uint8("port", "p", "端口号", uint8(250))
	_ = testCmd.Uint8("retries", "R", "重试次数", uint8(3))

	_ = testCmd.Uint16("timeout", "x", "超时时间(秒)", uint16(30))
	_ = testCmd.Uint16("interval", "I", "间隔时间(毫秒)", uint16(100))

	_ = testCmd.Uint32("buffer-size", "b", "缓冲区大小", uint32(1024))
	_ = testCmd.Uint32("max-buffer", "X", "最大缓冲区", uint32(8192))

	_ = testCmd.Uint64("file-size", "f", "文件大小", uint64(1048576))
	_ = testCmd.Uint64("max-file", "F", "最大文件大小", uint64(1073741824))

	_ = testCmd.Float64("ratio", "q", "比例", 0.75)
	_ = testCmd.Float64("threshold", "Q", "阈值", 0.5)

	_ = testCmd.Enum("mode", "A", "运行模式", "auto", []string{"auto", "manual", "debug"})
	_ = testCmd.Enum("level", "L", "日志级别", "info", []string{"debug", "info", "warn", "error"})

	_ = testCmd.Duration("duration", "j", "持续时间", time.Second*30)
	_ = testCmd.Duration("delay", "J", "延迟时间", time.Minute*5)

	_ = testCmd.Time("start-time", "S", "开始时间", time.Now())
	_ = testCmd.Time("end-time", "E", "结束时间", time.Time{})

	_ = testCmd.Size("max-size", "Z", "最大大小", int64(1024*1024))
	_ = testCmd.Size("min-size", "z", "最小大小", int64(512))

	_ = testCmd.StringSlice("paths", "P", "路径列表", nil)
	_ = testCmd.StringSlice("domains", "G", "域名列表", nil)

	_ = testCmd.IntSlice("ports", "W", "端口列表", nil)
	_ = testCmd.IntSlice("ids", "Y", "ID列表", nil)

	_ = testCmd.Int64Slice("large-numbers", "K", "大数字列表", nil)
	_ = testCmd.Int64Slice("codes", "k", "代码列表", nil)

	_ = testCmd.Map("headers", "H", "HTTP头部", nil)
	_ = testCmd.Map("tags", "w", "标签", nil)

	serverCmd := cmd.NewCmd("server", "s", types.ContinueOnError)
	serverCmd.SetDesc("服务器相关命令")
	_ = serverCmd.Bool("daemon", "d", "守护进程模式", false)
	_ = serverCmd.Uint16("listen", "l", "监听端口", uint16(8080))
	_ = serverCmd.String("host", "a", "主机地址", "localhost")

	databaseCmd := cmd.NewCmd("database", "a", types.ContinueOnError)
	databaseCmd.SetDesc("数据库相关命令")
	_ = databaseCmd.String("db-name", "n", "数据库名称", "testdb")
	_ = databaseCmd.String("db-host", "H", "数据库主机", "localhost")
	_ = databaseCmd.Uint16("db-port", "p", "数据库端口", uint16(3306))

	cacheCmd := cmd.NewCmd("cache", "c", types.ContinueOnError)
	cacheCmd.SetDesc("缓存相关命令")
	_ = cacheCmd.Uint32("ttl", "t", "缓存过期时间(秒)", uint32(3600))
	_ = cacheCmd.Uint64("max-items", "m", "最大缓存项数", uint64(10000))
	_ = cacheCmd.Bool("enabled", "e", "启用缓存", true)

	generateCmd := cmd.NewCmd("generate", "g", types.ContinueOnError)
	generateCmd.SetDesc("代码生成命令")
	_ = generateCmd.String("template", "T", "模板文件", "template.go")
	_ = generateCmd.String("output", "O", "输出文件", "output.go")
	_ = generateCmd.Bool("force", "f", "强制覆盖", false)

	deployCmd := cmd.NewCmd("deploy", "b", types.ContinueOnError)
	deployCmd.SetDesc("部署相关命令")
	_ = deployCmd.String("env", "E", "部署环境", "production")
	_ = deployCmd.String("version", "V", "版本号", "1.0.0")
	_ = deployCmd.Bool("rollback", "r", "回滚部署", false)

	if err := testCmd.AddSubCmds(serverCmd, databaseCmd, cacheCmd, generateCmd, deployCmd); err != nil {
		panic("Failed to add subcommands: " + err.Error())
	}

	return testCmd
}

// TestMain 全局测试入口, 控制非verbose模式下的输出重定向
func TestMain(m *testing.M) {
	flag.Parse() // 解析命令行参数
	// 保存原始标准输出和错误输出
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	var nullFile *os.File
	var err error

	// 非verbose模式下重定向到空设备
	if !testing.Verbose() {
		nullFile, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
		if err != nil {
			panic("无法打开空设备文件: " + err.Error())
		}
		os.Stdout = nullFile
		os.Stderr = nullFile
	}

	// 运行所有测试
	exitCode := m.Run()

	// 恢复原始输出
	if !testing.Verbose() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
		_ = nullFile.Close()
	}

	os.Exit(exitCode)
}

func TestCmdCreation(t *testing.T) {
	rootCmd := cmd.NewCmd("root", "", types.ContinueOnError)
	rootCmd.SetDesc("Root cmd for testing")

	if rootCmd.Name() != "root" {
		t.Errorf("Expected name 'root', got '%s'", rootCmd.Name())
	}

	if rootCmd.Desc() != "Root cmd for testing" {
		t.Errorf("Description not set correctly")
	}

	t.Log("✓ Cmd creation test passed")
}

func TestFlagRegistration(t *testing.T) {
	testCmd := cmd.NewCmd("test", "t", types.ContinueOnError)
	testCmd.SetDesc("Test cmd")

	// 使用便捷函数创建标志, 自动添加到命令中
	_ = testCmd.String("config", "c", "Config file path", "config.yaml")
	_ = testCmd.Int("port", "p", "Port number", 8080)

	flags := testCmd.Flags()
	// 由于标志对象只存储一次, 现在只有2个标志
	if len(flags) != 2 {
		t.Errorf("Expected 2 flags, got %d", len(flags))
	}

	// 验证长名称可以获取
	retrievedFlag, exists := testCmd.GetFlag("config")
	if !exists {
		t.Error("Failed to retrieve 'config' flag")
	}

	if retrievedFlag.Name() != "config" {
		t.Errorf("Expected flag name 'config', got '%s'", retrievedFlag.Name())
	}

	// 验证短名称也可以获取
	retrievedFlag, exists = testCmd.GetFlag("c")
	if !exists {
		t.Error("Failed to retrieve 'c' flag")
	}

	if retrievedFlag.Name() != "config" {
		t.Errorf("Expected flag name 'config', got '%s'", retrievedFlag.Name())
	}

	t.Log("✓ Flag registration test passed")
}

func TestArgumentParsing(t *testing.T) {
	t.Log("Starting TestArgumentParsing")

	testCmd := cmd.NewCmd("test", "", types.ContinueOnError)
	t.Log("Created base cmd")

	testCmd.SetDesc("Test cmd for parsing")
	t.Log("Set description")

	// 使用便捷函数创建标志, 自动添加到命令中
	_ = testCmd.String("config", "c", "Config file", "default.yaml")
	_ = testCmd.Int("port-num", "p", "Port", 8080)
	_ = testCmd.Bool("verbose", "V", "Verbose output", false)
	t.Log("Created flags")

	p := parser.NewDefaultParser(types.ContinueOnError)
	t.Log("Created parser")
	testCmd.SetParser(p)
	t.Log("Parser set")

	args := []string{"--config", "myconfig.yaml", "--port-num", "3000", "--verbose", "arg1", "arg2"}
	t.Logf("Parsing args: %v", args)

	if err := testCmd.Parse(args); err != nil {
		t.Errorf("Parse failed: %v", err)
	}
	t.Log("Parse completed")

	if !testCmd.IsParsed() {
		t.Error("Cmd should be marked as parsed")
	}

	if testCmd.NArg() != 2 {
		t.Errorf("Expected 2 args, got %d", testCmd.NArg())
	}

	if testCmd.Arg(0) != "arg1" {
		t.Errorf("Expected first arg 'arg1', got '%s'", testCmd.Arg(0))
	}

	if testCmd.Arg(1) != "arg2" {
		t.Errorf("Expected second arg 'arg2', got '%s'", testCmd.Arg(1))
	}

	t.Log("✓ Argument parsing test passed")
}

func TestSubCmd(t *testing.T) {
	rootCmd := cmd.NewCmd("app", "", types.ContinueOnError)
	rootCmd.SetDesc("Main application")

	subCmd := cmd.NewCmd("serve", "s", types.ContinueOnError)
	subCmd.SetDesc("Serve cmd")

	// 使用便捷函数创建标志, 自动添加到子命令中
	_ = subCmd.Int("port", "p", "Port to listen", 8080)

	if err := rootCmd.AddSubCmds(subCmd); err != nil {
		t.Errorf("Failed to add subcmd: %v", err)
	}

	SubCmds := rootCmd.SubCmds()
	if len(SubCmds) != 1 {
		t.Errorf("Expected 1 subcmd, got %d", len(SubCmds))
	}

	if !rootCmd.HasSubCmd("serve") {
		t.Error("Should have subcmd 'serve'")
	}

	retrievedSub, exists := rootCmd.GetSubCmd("serve")
	if !exists {
		t.Error("Failed to retrieve 'serve' subcmd")
	}

	if retrievedSub.Name() != "serve" {
		t.Errorf("Expected subcmd name 'serve', got '%s'", retrievedSub.Name())
	}

	t.Log("✓ Subcmd test passed")
}

func TestSubCmdParsing(t *testing.T) {
	rootCmd := cmd.NewCmd("app", "", types.ContinueOnError)
	rootCmd.SetDesc("Test app")

	// 使用便捷函数创建标志, 自动添加到根命令中
	_ = rootCmd.Bool("debug", "d", "Enable debug mode", false)

	subCmd := cmd.NewCmd("serve", "s", types.ContinueOnError)
	subCmd.SetDesc("Serve cmd")
	// 使用便捷函数创建标志, 自动添加到子命令中
	_ = subCmd.Int("port", "p", "Port", 8080)
	_ = rootCmd.AddSubCmds(subCmd)

	p := parser.NewDefaultParser(types.ContinueOnError)
	rootCmd.SetParser(p)

	args := []string{"--debug=true", "serve", "--port=9000"}

	if err := rootCmd.Parse(args); err != nil {
		t.Errorf("Parse failed: %v", err)
	}

	subCmdParsed, exists := rootCmd.GetSubCmd("serve")
	if !exists {
		t.Error("Failed to get subcmd 'serve'")
	}
	if !subCmdParsed.IsParsed() {
		t.Error("Subcmd should be marked as parsed")
	}

	t.Log("✓ Subcmd parsing test passed")
}

func TestParseAndRoute(t *testing.T) {
	testCmd := cmd.NewCmd("test", "", types.ContinueOnError)
	testCmd.SetDesc("Test cmd")

	executed := false
	testCmd.SetRun(func(c types.Command) error {
		executed = true
		return nil
	})

	p := parser.NewDefaultParser(types.ContinueOnError)
	testCmd.SetParser(p)

	if err := testCmd.ParseAndRoute([]string{}); err != nil {
		t.Errorf("ParseAndRoute failed: %v", err)
	}

	if !executed {
		t.Error("Run function should have been executed")
	}

	t.Log("✓ ParseAndRoute test passed")
}

func TestEnvironmentVariables(t *testing.T) {
	testCmd := cmd.NewCmd("test", "", types.ContinueOnError)
	testCmd.SetDesc("Test cmd with env vars")

	// 使用便捷函数创建标志, 自动添加到命令中
	configFlag := testCmd.String("config", "c", "Config file", "default.yaml")
	configFlag.BindEnv("CONFIG_FILE")

	testCmd.SetEnvPrefix("MYAPP")

	p := parser.NewDefaultParser(types.ContinueOnError)
	testCmd.SetParser(p)

	_ = os.Setenv("MYAPP_CONFIG_FILE", "/env/config.yaml")
	defer func() { _ = os.Unsetenv("MYAPP_CONFIG_FILE") }()

	if err := testCmd.Parse([]string{}); err != nil {
		t.Errorf("Parse failed: %v", err)
	}

	t.Log("✓ Environment variable test passed")
}

func TestParseOnly(t *testing.T) {
	testCmd := cmd.NewCmd("test", "", types.ContinueOnError)
	testCmd.SetDesc("Test cmd")

	// 使用便捷函数创建标志, 自动添加到命令中
	_ = testCmd.String("name", "n", "Name", "default")
	_ = testCmd.Bool("verbose", "v", "Verbose", false)

	p := parser.NewDefaultParser(types.ContinueOnError)
	testCmd.SetParser(p)

	if err := testCmd.ParseOnly([]string{"--name=test", "arg1", "arg2"}); err != nil {
		t.Errorf("ParseOnly failed: %v", err)
	}

	if len(testCmd.Args()) != 2 {
		t.Errorf("Expected 2 args, got %d", len(testCmd.Args()))
	}

	if testCmd.Arg(0) != "arg1" || testCmd.Arg(1) != "arg2" {
		t.Errorf("Unexpected args: %v", testCmd.Args())
	}

	if !testCmd.IsParsed() {
		t.Error("Cmd should be marked as parsed")
	}

	t.Log("✓ ParseOnly test passed")
}

func TestFlagNotFound(t *testing.T) {
	testCmd := cmd.NewCmd("test", "", types.ContinueOnError)
	testCmd.SetDesc("Test cmd")

	// 使用便捷函数创建标志, 自动添加到命令中
	_ = testCmd.String("name", "n", "Name", "default")

	p := parser.NewDefaultParser(types.ContinueOnError)
	testCmd.SetParser(p)

	err := testCmd.Parse([]string{"--unknown=value"})
	if err == nil {
		t.Error("Should fail on unknown flag")
	}

	t.Logf("Expected error for unknown flag: %v", err)
}

func TestIntFlagParsing(t *testing.T) {
	testCmd := cmd.NewCmd("test", "", types.ContinueOnError)
	testCmd.SetDesc("Test cmd")

	// 使用便捷函数创建标志, 自动添加到命令中
	_ = testCmd.Int("port", "p", "Port number", 8080)

	p := parser.NewDefaultParser(types.ContinueOnError)
	testCmd.SetParser(p)

	if err := testCmd.Parse([]string{"--port=3000"}); err != nil {
		t.Errorf("Parse failed: %v", err)
	}

	t.Log("✓ Int flag parsing test passed")
}

func TestBoolFlagParsing(t *testing.T) {
	testCmd := cmd.NewCmd("test", "", types.ContinueOnError)
	testCmd.SetDesc("Test cmd")

	// 使用便捷函数创建标志, 自动添加到命令中
	// 使用一个不会与内置标志冲突的名称
	_ = testCmd.Bool("debug", "d", "Debug output", false)

	p := parser.NewDefaultParser(types.ContinueOnError)
	testCmd.SetParser(p)

	if err := testCmd.Parse([]string{"--debug=true"}); err != nil {
		t.Errorf("Parse failed: %v", err)
	}

	t.Log("✓ Bool flag parsing test passed")
}

func TestHelpGeneration(t *testing.T) {
	testLogo := `_                   __   __            
FJ___       ____     LJ   LJ     ____   
J  __.    F __ J    FJ   FJ    F __ J  
| |--| |   | _____J  J  L J  L  | |--| | 
F L  J J   F L___--. J  L J  L  F L__J J 
J__L  J__L J\______/F J__L J__L J\______/F
|__L  J__|  J______F  |__| |__|  J______F                                     
	`
	testCmd := cmd.NewCmd("test", "t", types.ContinueOnError)
	testCmd.SetChinese(true)
	testCmd.SetDesc("Test cmd for help generation with many flags and SubCmds")
	testCmd.AddNote("9999")
	testCmd.AddNote("8888")
	testCmd.AddExample("测试", "111")
	testCmd.SetLogoText(testLogo)
	testCmd.SetVersion("1.0.0")

	testCmd1 := cmd.NewCmd("cest1", "t1", types.ContinueOnError)
	testCmd1.SetDesc("Test cmd 1 for help")
	testCmd2 := cmd.NewCmd("gest2", "", types.ContinueOnError)
	testCmd2.SetDesc("Test cmd 2 for help")
	testCmd3 := cmd.NewCmd("", "pt3", types.ContinueOnError)
	testCmd3.SetDesc("Test cmd 3 for help")
	testCmd4 := cmd.NewCmd("hest4", "", types.ContinueOnError)
	testCmd4.SetDesc("Test cmd 4 for help")
	testCmd5 := cmd.NewCmd("", "ot5", types.ContinueOnError)
	testCmd5.SetDesc("Test cmd 5 for help")

	_ = testCmd.AddSubCmds(testCmd1, testCmd2, testCmd3, testCmd4, testCmd5)

	// 使用便捷函数创建标志, 自动添加到命令中
	_ = testCmd.String("config", "c", "Config file", "config.yaml")
	_ = testCmd.Int("port", "p", "Port number", 8080)
	_ = testCmd.Bool("verbose", "V", "Enable verbose output", false)
	_ = testCmd.String("output", "o", "Output file path", "output.txt")
	_ = testCmd.Int("level", "l", "Log level (1-5)", 3)
	_ = testCmd.Bool("debug", "bbc", "Enable debug mode", false)
	_ = testCmd.String("input", "", "Input file path", "")
	_ = testCmd.String("", "abc", "Request timeout in seconds", "30")
	_ = testCmd.String("threshold", "th", "Threshold value", "0.75")
	_ = testCmd.String("format", "f", "Output format (json/yaml/text)", "text")

	_ = testCmd.Parse([]string{})

	helpText := testCmd.Help()
	if helpText == "" {
		t.Error("Help text should not be empty")
	}
	testCmd.PrintHelp()
	testCmd1.PrintHelp()

	// 获取标志的值并打印
	// fmt.Printf("Flag name: %s\n", flag1.Name())
	// fmt.Printf("Flag long name: %s\n", flag1.LongName())
	// fmt.Printf("Flag short name: %s\n", flag1.ShortName())
	// fmt.Printf("Flag description: %s\n", flag1.Desc())
	// fmt.Printf("Flag type: %v\n", flag1.Type())
	// fmt.Printf("Flag value: %v\n", flag1.Get())
	// fmt.Printf("Flag default value: %v\n", flag1.GetDef())
	// fmt.Printf("Flag is set: %t\n", flag1.IsSet())
	// fmt.Printf("Flag formatted name: %s\n", flag1.String())

	//completion.GenAndPrint(testCmd, types.BashShell)
}

// 测试内置标志
func TestBuiltinFlags(t *testing.T) {
	// 测试帮助标志 (中文)
	t.Run("HelpFlag_Chinese", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic due to os.Exit, but got none")
			}
		}()

		cmd := cmd.NewCmd("test", "t", types.ContinueOnError)
		cmd.SetDesc("测试命令")
		cmd.SetChinese(true)

		// 解析帮助标志
		err := cmd.Parse([]string{"--help"})
		// 由于帮助标志会调用os.Exit, 这里应该不会执行到
		// 所以我们检查错误类型
		if err == nil {
			t.Error("Expected error due to os.Exit, but got nil")
		}
	})

	// 测试帮助标志 (英文)
	t.Run("HelpFlag_English", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic due to os.Exit, but got none")
			}
		}()

		cmd := cmd.NewCmd("test", "t", types.ContinueOnError)
		cmd.SetDesc("Test command")
		cmd.SetChinese(false)

		// 解析帮助标志
		err := cmd.Parse([]string{"-h"})
		// 由于帮助标志会调用os.Exit, 这里应该不会执行到
		// 所以我们检查错误类型
		if err == nil {
			t.Error("Expected error due to os.Exit, but got nil")
		}
	})

	// 测试版本标志 (有版本信息)
	t.Run("VersionFlag_WithVersion", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic due to os.Exit, but got none")
			}
		}()

		cmd := cmd.NewCmd("test", "t", types.ContinueOnError)
		cmd.SetDesc("测试命令")
		cmd.SetVersion("1.0.0")
		cmd.SetChinese(true)

		// 解析版本标志
		err := cmd.Parse([]string{"--version"})
		// 由于版本标志会调用os.Exit, 这里应该不会执行到
		// 所以我们检查错误类型
		if err == nil {
			t.Error("Expected error due to os.Exit, but got nil")
		}
	})

	// 测试版本标志 (无版本信息)
	t.Run("VersionFlag_WithoutVersion", func(t *testing.T) {
		cmd := cmd.NewCmd("test", "t", types.ContinueOnError)
		cmd.SetDesc("测试命令")
		// 不设置版本信息

		// 解析版本标志, 应该报错, 因为标志不存在
		err := cmd.Parse([]string{"--version"})
		if err == nil {
			t.Error("Expected error for unknown flag, but got nil")
		}
	})

	// 测试标志注册
	t.Run("BuiltinFlags_Registration", func(t *testing.T) {
		// 测试有版本信息的命令
		cmd1 := cmd.NewCmd("test1", "t1", types.ContinueOnError)
		cmd1.SetDesc("测试命令1")
		cmd1.SetVersion("1.0.0")
		cmd1.SetChinese(true)

		// 解析空参数, 只注册标志不执行
		err := cmd1.ParseOnly([]string{})
		if err != nil {
			t.Errorf("ParseOnly failed: %v", err)
		}

		// 检查注册的标志
		flags := cmd1.Flags()
		flagNames := make(map[string]bool)
		for _, f := range flags {
			flagNames[f.Name()] = true
		}

		// 应该有help和version两个标志
		if !flagNames["help"] {
			t.Error("Help flag not registered")
		}
		if !flagNames["version"] {
			t.Error("Version flag not registered")
		}

		// 检查标志描述是中文
		flagDescs := make(map[string]string)
		for _, f := range flags {
			flagDescs[f.Name()] = f.Desc()
		}

		if flagDescs["help"] != "显示帮助信息" {
			t.Errorf("Expected Chinese help description, got %s", flagDescs["help"])
		}

		if flagDescs["version"] != "显示版本信息" {
			t.Errorf("Expected Chinese version description, got %s", flagDescs["version"])
		}

		// 测试没有版本信息的命令
		cmd2 := cmd.NewCmd("test2", "t2", types.ContinueOnError)
		cmd2.SetDesc("Test command 2")
		cmd2.SetChinese(false)
		// 不设置版本信息

		// 解析空参数, 只注册标志不执行
		err = cmd2.ParseOnly([]string{})
		if err != nil {
			t.Errorf("ParseOnly failed: %v", err)
		}

		// 检查注册的标志
		flags2 := cmd2.Flags()
		flagNames2 := make(map[string]bool)
		for _, f := range flags2 {
			flagNames2[f.Name()] = true
		}

		// 应该只有help标志
		if !flagNames2["help"] {
			t.Error("Help flag not registered")
		}

		if flagNames2["version"] {
			t.Error("Version flag should not be registered when no version is set")
		}

		// 检查标志描述是英文
		flagDescs2 := make(map[string]string)
		for _, f := range flags2 {
			flagDescs2[f.Name()] = f.Desc()
		}

		if flagDescs2["help"] != "Show help information" {
			t.Errorf("Expected English help description, got %s", flagDescs2["help"])
		}
	})
}
