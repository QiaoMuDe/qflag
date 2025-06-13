package qflag

import (
	"flag"
	"fmt"
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
	if strFlag.help != "name help" {
		t.Errorf("Expected help 'name help', got '%s'", strFlag.help)
	}

	// 测试整数标志
	intFlag := cmd.Int("port", "p", 8080, "port help")
	if *intFlag.value != 8080 {
		t.Errorf("Expected default value 8080, got %d", *intFlag.value)
	}
	if intFlag.help != "port help" {
		t.Errorf("Expected help 'port help', got '%s'", intFlag.help)
	}

	// 测试布尔标志
	boolFlag := cmd.Bool("verbose", "v", false, "verbose help")
	if *boolFlag.value != false {
		t.Errorf("Expected default value false, got %t", *boolFlag.value)
	}
	if boolFlag.help != "verbose help" {
		t.Errorf("Expected help 'verbose help', got '%s'", boolFlag.help)
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

	// 模拟命令行参数同时使用长短标志
	err := cmd.Parse([]string{"--name", "value", "-n", "value"})
	if err == nil {
		t.Error("Expected error for flag conflict")
	} else if err.Error() != "不能同时使用 --name 和 -n" {
		t.Errorf("Expected error message '不能同时使用 --name 和 -n', got '%s'", err.Error())
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

	// 测试设置标志
	err := cmd.Parse([]string{"--name", "value"})
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
	cmd.Help = "Custom help content"

	if cmd.Help != "Custom help content" {
		t.Errorf("Expected Help 'Custom help content', got '%s'", cmd.Help)
	}
}

// TestGenerateHelpInfo 测试帮助信息生成函数
func TestGenerateHelpInfo(t *testing.T) {
	// 场景1: 只有主命令，没有标志和子命令
	t.Run("OnlyMainCommand", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ExitOnError)
		fmt.Println("=== 场景: 只有主命令，没有标志和子命令 ===")
		helpInfo := generateHelpInfo(cmd, true)
		fmt.Println(helpInfo)

		// 验证命令名
		expectedHeader := "命令: test(t)\n"
		if !strings.Contains(helpInfo, expectedHeader) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedHeader, helpInfo)
		}

		// 验证用法说明
		expectedUsage := "用法: test [选项] [参数]"
		if !strings.Contains(helpInfo, expectedUsage) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedUsage, helpInfo)
		}
	})

	// 场景2: 主命令有标志
	t.Run("MainCommandWithFlags", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ExitOnError)
		cmd.String("name", "n", "default", "name help")
		cmd.Int("port", "p", 8080, "port help")

		fmt.Println("=== 场景: 主命令带标志 ===")
		helpInfo := generateHelpInfo(cmd, true)
		fmt.Println(helpInfo)

		// 验证选项标题
		expectedOptionsHeader := "选项:\n"
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
		mainCmd.Description = "主命令描述"
		subCmd := NewCmd("sub", "s", flag.ExitOnError)
		subCmd.Description = "子命令描述"
		mainCmd.AddSubCmd(subCmd)

		fmt.Println("=== 场景: 主命令带子命令 ===")
		helpInfo := generateHelpInfo(mainCmd, true)
		fmt.Println(helpInfo)

		// 验证子命令标题
		expectedSubHeader := "子命令:\n"
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

		fmt.Println("=== 场景: 子命令帮助信息 ===")
		helpInfo := generateHelpInfo(subCmd, false)
		fmt.Println(helpInfo)

		// 验证子命令用法说明
		expectedUsage := "用法: main sub [选项] [参数]"
		if !strings.Contains(helpInfo, expectedUsage) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedUsage, helpInfo)
		}

		// 验证子命令标志
		expectedConfigFlag := "-c, --config"
		if !strings.Contains(helpInfo, expectedConfigFlag) {
			t.Errorf("Expected help info to contain '%s', got '%s'", expectedConfigFlag, helpInfo)
		}
	})
}

// TestGenerateHelpInfoWithSubcommands 测试包含子命令的帮助信息生成
func TestGenerateHelpInfoWithSubcommands(t *testing.T) {
	// 创建主命令
	mainCmd := NewCmd("main", "m", flag.ExitOnError)
	mainCmd.Description = "主命令描述"
	mainCmd.String("config", "c", "config.json", "配置文件路径")

	// 添加子命令
	subCmd1 := NewCmd("sub1", "s1", flag.ExitOnError)
	subCmd1.Description = "子命令1描述"
	subCmd1.Int("port", "p", 8080, "监听端口")
	mainCmd.AddSubCmd(subCmd1)

	subCmd2 := NewCmd("sub2", "s2", flag.ExitOnError)
	subCmd2.Description = "子命令2描述"
	subCmd2.Bool("verbose", "v", false, "详细输出")
	mainCmd.AddSubCmd(subCmd2)

	// 生成帮助信息
	fmt.Println("=== 场景: 主命令带子命令 ===")
	helpInfo := generateHelpInfo(mainCmd, true)
	fmt.Println(helpInfo)

	// 验证命令描述
	expectedDesc := "主命令描述\n"
	if !strings.Contains(helpInfo, expectedDesc) {
		t.Errorf("Expected help info to contain '%s', got '%s'", expectedDesc, helpInfo)
	}

	// 验证子命令标题
	expectedSubHeader := "子命令:\n"
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
