package qflag

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestNewCmd 测试创建新Cmd实例
func TestNewCmd(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ExitOnError)
	if cmd == nil {
		t.Fatal("NewCmd returned nil")
	}
	if cmd.fs.Name() != "test" {
		t.Errorf("Expected fs name 'test', got '%s'", cmd.fs.Name())
	}
}

// TestBindFlags 测试标志绑定功能
func TestBindFlags(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ExitOnError)

	// 测试字符串标志
	strFlag := cmd.String("name", "n", "default", "name help")
	if *strFlag.value != "default" {
		t.Errorf("Expected default value 'default', got '%s'", *strFlag.value)
	}
	if strFlag.usage != "name help" {
		t.Errorf("Expected help 'name help', got '%s'", strFlag.usage)
	}

	// 测试整数标志
	intFlag := cmd.Int("port", "p", 8080, "port help")
	if *intFlag.value != 8080 {
		t.Errorf("Expected default value 8080, got %d", *intFlag.value)
	}
	if intFlag.usage != "port help" {
		t.Errorf("Expected help 'port help', got '%s'", intFlag.usage)
	}

	// 测试布尔标志
	boolFlag := cmd.Bool("verbose", "v", false, "verbose help")
	if *boolFlag.value != false {
		t.Errorf("Expected default value false, got %t", *boolFlag.value)
	}
	if boolFlag.usage != "verbose help" {
		t.Errorf("Expected help 'verbose help', got '%s'", boolFlag.usage)
	}
}

// TestHelpFlag 测试帮助标志功能
func TestHelpFlag(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ExitOnError)
	if !cmd.helpFlagBound {
		t.Error("helpFlagBound should be true")
	}
	if cmd.helpFlagName != "help" {
		t.Errorf("Expected helpFlagName 'help', got '%s'", cmd.helpFlagName)
	}
	if cmd.helpShortName != "h" {
		t.Errorf("Expected helpShortName 'h', got '%s'", cmd.helpShortName)
	}
}

// TestFlagConflict 测试长短标志冲突检测
func TestFlagConflict(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ExitOnError)
	cmd.String("name", "n", "default", "name help")

	// 模拟命令行参数同时使用长短标志并过滤-test.*标志
	args := []string{"--name", "value", "-n", "value"}
	filteredArgs := []string{}
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-test.") {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	err := cmd.Parse(filteredArgs)
	if err == nil {
		t.Error("Expected error for flag conflict")
	} else if err.Error() != "Cannot use both --name and -n" {
		t.Errorf("Expected error message 'Cannot use both --name and -n', got '%s'", err.Error())
	}
}

// TestIsFlagSet 测试标志设置检测
func TestIsFlagSet(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ExitOnError)
	cmd.String("name", "n", "default", "name help")

	// 测试未设置标志
	if cmd.isFlagSet("name") {
		t.Error("Expected isFlagSet to return false before parsing")
	}

	// 测试设置标志并过滤-test.*标志
	args := []string{"--name", "value"}
	filteredArgs := []string{}
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-test.") {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	err := cmd.Parse(filteredArgs)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !cmd.isFlagSet("name") {
		t.Error("Expected isFlagSet to return true after parsing")
	}
}

// TestHelpContent 测试自定义帮助内容
func TestHelpContent(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ExitOnError)
	cmd.SetUsage("Custom help content")

	if cmd.Usage() != "Custom help content" {
		t.Errorf("Expected Help 'Custom help content', got '%s'", cmd.Usage())
	}
}

// TestGlobalDefaultCmd 测试全局默认命令功能
func TestGlobalDefaultCmd(t *testing.T) {
	// 测试全局String函数
	strFlag := String("gname", "gn", "gdefault", "global name help")
	if *strFlag.value != "gdefault" {
		t.Errorf("全局String标志默认值错误, 预期 'gdefault', 实际 '%s'", *strFlag.value)
	}

	// 测试全局Int函数
	intFlag := Int("gport", "gp", 9090, "global port help")
	if *intFlag.value != 9090 {
		t.Errorf("全局Int标志默认值错误, 预期 9090, 实际 %d", *intFlag.value)
	}

	// 测试全局Bool函数
	boolFlag := Bool("gverbose", "gv", true, "global verbose help")
	if *boolFlag.value != true {
		t.Errorf("全局Bool标志默认值错误, 预期 true, 实际 %t", *boolFlag.value)
	}

	// 测试全局Float函数
	floatFlag := Float("gpi", "gπ", 3.14, "global pi help")
	if *floatFlag.value != 3.14 {
		t.Errorf("全局Float标志默认值错误, 预期 3.14, 实际 %f", *floatFlag.value)
	}
}

// TestGlobalFlagConflict 测试全局命令标志冲突
func TestGlobalFlagConflict(t *testing.T) {
	// 重置默认命令以避免测试污染
	defaultCmd = NewCmd("main", "", flag.ExitOnError)
	String("gname", "gn", "", "global name help")

	// 同时使用长短标志应返回错误
	// 准备测试参数
	testArgs := []string{"--gname", "value", "-gn", "value"}

	// 保存原始os.Args并在测试后恢复
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// 设置os.Args为测试参数
	os.Args = append([]string{os.Args[0]}, testArgs...)

	err := Parse()
	if err == nil {
		t.Error("Expected error for global flag conflict")
	} else if err.Error() != "Cannot use both --gname and -gn" {
		t.Errorf("Global flag conflict message incorrect, expected 'Cannot use both --gname and -gn', got '%s'", err.Error())
	}
}

// TestGenerateHelpInfo 测试帮助信息生成函数
func TestGenerateHelpInfo(t *testing.T) {
	// 场景1: 只有主命令，没有标志和子命令
	t.Run("OnlyMainCommand", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ExitOnError)
		fmt.Println("=== Scenario: Main command only, no flags or subcommands ===")
		helpInfo := generateHelpInfo(cmd, true)
		fmt.Println(helpInfo)

		// 验证命令名
		expectedHeader := "Command: test(t)\n"
		if !strings.Contains(helpInfo, expectedHeader) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedHeader, helpInfo)
		}

		// 验证用法说明
		expectedUsage := "Usage: test [options] [arguments]"
		if !strings.Contains(helpInfo, expectedUsage) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedUsage, helpInfo)
		}
	})

	// 场景2: 主命令有标志
	t.Run("MainCommandWithFlags", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ExitOnError)
		cmd.String("name", "n", "default", "name help")
		cmd.Int("port", "p", 8080, "port help")

		fmt.Println("=== Scenario: Main command with flags ===")
		helpInfo := generateHelpInfo(cmd, true)
		fmt.Println(helpInfo)

		// 验证选项标题
		expectedOptionsHeader := "Options:\n"
		if !strings.Contains(helpInfo, expectedOptionsHeader) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedOptionsHeader, helpInfo)
		}

		// 验证标志信息
		expectedNameFlag := "-n, --name"
		if !strings.Contains(helpInfo, expectedNameFlag) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedNameFlag, helpInfo)
		}

		expectedPortFlag := "-p, --port"
		if !strings.Contains(helpInfo, expectedPortFlag) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedPortFlag, helpInfo)
		}
	})

	// 场景3: 主命令有子命令
	t.Run("MainCommandWithSubcommands", func(t *testing.T) {
		mainCmd := NewCmd("main", "m", flag.ExitOnError)
		mainCmd.SetDescription("Main command description")
		subCmd := NewCmd("sub", "s", flag.ExitOnError)
		subCmd.SetDescription("Subcommand description")
		mainCmd.AddSubCmd(subCmd)

		fmt.Println("=== Scenario: Main command with subcommands ===")
		helpInfo := generateHelpInfo(mainCmd, true)
		fmt.Println(helpInfo)

		// 验证子命令标题
		expectedSubHeader := "Subcommands:\n"
		if !strings.Contains(helpInfo, expectedSubHeader) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedSubHeader, helpInfo)
		}

		// 验证子命令信息
		expectedSubCmd := "sub"
		if !strings.Contains(helpInfo, expectedSubCmd) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedSubCmd, helpInfo)
		}
	})

	// 场景4: 子命令的帮助信息
	t.Run("SubcommandHelpInfo", func(t *testing.T) {
		mainCmd := NewCmd("main", "m", flag.ExitOnError)
		subCmd := NewCmd("sub", "s", flag.ExitOnError)
		subCmd.String("config", "c", "config.json", "配置文件路径")
		mainCmd.AddSubCmd(subCmd)

		fmt.Println("=== Scenario: Subcommand help information ===")
		helpInfo := generateHelpInfo(subCmd, false)
		fmt.Println(helpInfo)

		// 验证子命令用法说明
		expectedUsage := "Usage: main sub [options] [arguments]"
		if !strings.Contains(helpInfo, expectedUsage) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedUsage, helpInfo)
		}
	})

	// 场景5: 全局命令帮助信息
	t.Run("GlobalCommandHelpInfo", func(t *testing.T) {
		// 重置默认命令
		defaultCmd = NewCmd("main", "", flag.ContinueOnError)
		String("gname", "gn", "gdefault", "global name help")
		Int("gport", "gp", 9090, "global port help")

		fmt.Println("=== Scenario: Global command help information ===")
		// 设置测试参数，避免依赖os.Args
		testArgs := []string{"--gname", "test", "-gp", "8888"}
		// 忽略未知标志以避免测试框架干扰
		defaultCmd.Parse(testArgs)
		helpInfo := generateHelpInfo(defaultCmd, true)
		fmt.Println(helpInfo)

		// 验证全局命令标志
		expectedNameFlag := "-gn, --gname"
		if !strings.Contains(helpInfo, expectedNameFlag) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedNameFlag, helpInfo)
		}

		expectedPortFlag := "-gp, --gport"
		if !strings.Contains(helpInfo, expectedPortFlag) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedPortFlag, helpInfo)
		}
	})
	// 测试子命令标志
	t.Run("SubcommandFlags", func(t *testing.T) {
		mainCmd := NewCmd("main", "m", flag.ExitOnError)
		subCmd := NewCmd("sub", "s", flag.ExitOnError)
		subCmd.String("config", "c", "config.json", "配置文件路径")
		mainCmd.AddSubCmd(subCmd)

		helpInfo := generateHelpInfo(subCmd, false)
		expectedConfigFlag := "-c, --config"
		if !strings.Contains(helpInfo, expectedConfigFlag) {
			t.Errorf("Subcommand help info should contain '%s', got '%s'", expectedConfigFlag, helpInfo)
		}
	})
}

// TestGenerateHelpInfoWithSubcommands 测试包含子命令的帮助信息生成
func TestGenerateHelpInfoWithSubcommands(t *testing.T) {
	// 创建主命令，使用ContinueOnError模式以忽略未知标志
	mainCmd := NewCmd("main", "m", flag.ContinueOnError)
	mainCmd.SetDescription("Main command description")
	mainCmd.String("config", "c", "config.json", "Configuration file path")

	// 添加子命令
	subCmd1 := NewCmd("sub1", "s1", flag.ContinueOnError)
	subCmd1.SetDescription("Subcommand 1 description")
	subCmd1.Int("port", "p", 8080, "Listening port")
	mainCmd.AddSubCmd(subCmd1)

	subCmd2 := NewCmd("sub2", "s2", flag.ContinueOnError)
	subCmd2.SetDescription("Subcommand 2 description")
	subCmd2.Bool("verbose", "v", false, "Verbose output")
	mainCmd.AddSubCmd(subCmd2)

	// 生成帮助信息
	fmt.Println("=== Scenario: Main command with subcommands ===")
	helpInfo := generateHelpInfo(mainCmd, true)
	fmt.Println(helpInfo)

	// 验证命令描述
	expectedDesc := "Description: Main command description\n\n"
	if !strings.Contains(helpInfo, expectedDesc) {
		t.Errorf("Expected help info to contain '%s', got '%s'", expectedDesc, helpInfo)
	}

	// 验证子命令标题
	expectedSubHeader := "Subcommands:\n"
	if !strings.Contains(helpInfo, expectedSubHeader) {
		t.Errorf("Expected help info to contain '%s', got '%s'", expectedSubHeader, helpInfo)
	}

	// 验证子命令1
	expectedSub1 := "sub1"
	if !strings.Contains(helpInfo, expectedSub1) {
		t.Errorf("Expected help info to contain '%s', got '%s'", expectedSub1, helpInfo)
	}

	// 验证子命令2
	expectedSub2 := "sub2"
	if !strings.Contains(helpInfo, expectedSub2) {
		t.Errorf("Expected help info to contain '%s', got '%s'", expectedSub2, helpInfo)
	}
}
