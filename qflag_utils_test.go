package qflag

import (
	"flag"
	"testing"
)

// TestBindHelpFlag 测试绑定帮助标志
func TestBindHelpFlag(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ExitOnError)
	cmd.bindHelpFlag()

	// 验证帮助标志已绑定
	if !cmd.helpFlagBound {
		t.Error("help flag should be bound")
	}
	if cmd.fs.Lookup(cmd.helpFlagName) == nil {
		t.Error("help flag should be registered")
	}
	if cmd.helpShortName != "" && cmd.fs.Lookup(cmd.helpShortName) == nil {
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
	cmd2 := NewCmd("test2", "", flag.ExitOnError)
	cmd2.SetDescription("Test description")
	cmd2.Bool("verbose", "v", false, "verbose output")
	cmd2.Int("count", "", 0, "number of times to repeat")
	if testing.Verbose() {
		cmd2.PrintUsage()
	}

	// 测试带子命令的用法信息
	cmd3 := NewCmd("parent", "", flag.ExitOnError)
	subCmd := NewCmd("child", "", flag.ExitOnError)
	if err := cmd3.AddSubCmd(subCmd); err != nil {
		t.Errorf("AddSubCmd error: %v", err)
	}
	if testing.Verbose() {
		cmd3.PrintUsage()
	}
}

// TestHasCycle 测试检测循环引用
func TestHasCycle(t *testing.T) {
	cmd1 := NewCmd("cmd1", "", flag.ExitOnError)
	cmd2 := NewCmd("cmd2", "", flag.ExitOnError)
	cmd3 := NewCmd("cmd3", "", flag.ExitOnError)
	cmd4 := NewCmd("cmd4", "", flag.ExitOnError)

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
