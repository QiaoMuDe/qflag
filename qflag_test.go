package qflag

import (
	"flag"
	"testing"
)

// TestNewCmd 测试创建新Cmd实例
func TestNewCmd(t *testing.T) {
	cmd := NewCmd("test", flag.ExitOnError)
	if cmd == nil {
		t.Fatal("NewCmd returned nil")
	}
	if cmd.fs.Name() != "test" {
		t.Errorf("Expected fs name 'test', got '%s'", cmd.fs.Name())
	}
}

// TestBindFlags 测试标志绑定功能
func TestBindFlags(t *testing.T) {
	cmd := NewCmd("test", flag.ExitOnError)

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
	cmd := NewCmd("test", flag.ExitOnError)
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
	cmd := NewCmd("test", flag.ExitOnError)
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
	cmd := NewCmd("test", flag.ExitOnError)
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
	cmd := NewCmd("test", flag.ExitOnError)
	cmd.Help = "Custom help content"

	if cmd.Help != "Custom help content" {
		t.Errorf("Expected Help 'Custom help content', got '%s'", cmd.Help)
	}
}
