package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// TestCompletionPerformance 测试补全脚本生成性能
func TestCompletionPerformance(t *testing.T) {
	// 创建复杂命令结构
	rootCmd := NewCmd("root", "r", flag.ExitOnError)
	rootCmd.SetExitOnBuiltinFlags(false)
	rootCmd.SetEnableCompletion(true)

	// 添加多层子命令和大量选项
	for i := 0; i < 10; i++ {
		parentCmd := NewCmd(fmt.Sprintf("sub%d", i), fmt.Sprintf("s%d", i), flag.ExitOnError)
		if err := rootCmd.AddSubCmd(parentCmd); err != nil {
			t.Fatalf("Failed to add parent subcommand: %v", err)
		}

		// 为每个子命令添加选项
		for j := 0; j < 20; j++ {
			strVar := &flags.StringFlag{}
			parentCmd.StringVar(strVar, fmt.Sprintf("option%d", j), fmt.Sprintf("o%d", j), "", "test option")
		}

		// 添加孙子命令
		for k := 0; k < 5; k++ {
			childCmd := NewCmd(fmt.Sprintf("sub%d-grand%d", i, k), fmt.Sprintf("g%d", k), flag.ExitOnError)
			if err := parentCmd.AddSubCmd(childCmd); err != nil {
				t.Fatalf("Failed to add child subcommand: %v", err)
			}

			// 为孙子命令添加选项
			for l := 0; l < 15; l++ {
				intVar := &flags.IntFlag{}
				childCmd.IntVar(intVar, fmt.Sprintf("param%d", l), fmt.Sprintf("p%d", l), 0, "test parameter")
			}
		}
	}

	// 测试Bash补全生成速度
	start := time.Now()
	var err error
	_, err = rootCmd.generateShellCompletion("bash")
	if err != nil {
		t.Fatalf("Bash completion generation failed: %v", err)
	}
	bashDuration := time.Since(start)
	t.Logf("Bash completion generated in %v", bashDuration)
	if bashDuration > 100*time.Millisecond {
		t.Errorf("Bash completion generation took too long: %v", bashDuration)
	}

	// 测试PowerShell补全生成速度
	start = time.Now()
	_, err = rootCmd.generateShellCompletion("powershell")
	if err != nil {
		t.Fatalf("PowerShell completion generation failed: %v", err)
	}
	pwshDuration := time.Since(start)
	t.Logf("PowerShell completion generated in %v", pwshDuration)
	if pwshDuration > 150*time.Millisecond {
		t.Errorf("PowerShell completion generation took too long: %v", pwshDuration)
	}
}

// TestCompletionBash 测试自动补全生成
func TestCompletionBash(t *testing.T) {
	// 新建根命令
	cmd := NewCmd("root", "r", flag.ExitOnError)
	cmd.SetExitOnBuiltinFlags(false) // 禁止在解析命令行参数时退出
	cmd.SetUseChinese(true)          // 设置使用中文
	cmd.SetEnableCompletion(true)    // 启用自动补全功能

	cmd1 := NewCmd("cmd1", "c1", flag.ExitOnError)
	cmd1.String("str", "s", "", "test string")
	cmd2 := NewCmd("cmd2", "c2", flag.ExitOnError)
	cmd2.Int("int", "i", 0, "test int")

	if err := cmd.AddSubCmd(cmd1); err != nil {
		t.Fatal(err)
	}

	if err := cmd1.AddSubCmd(cmd2); err != nil {
		t.Fatal(err)
	}

	// 解析命令行参数
	if err := cmd.Parse([]string{"-gsc", "bash"}); err != nil {
		t.Fatal(err)
	}
}

// TestCompletionPwsh 测试PowerShell自动补全生成
func TestCompletionPwsh(t *testing.T) {
	// 新建根命令
	cmd := NewCmd("root", "r", flag.ExitOnError)
	cmd.SetExitOnBuiltinFlags(false) // 禁止在解析命令行参数时退出
	cmd.SetUseChinese(true)          // 设置使用中文
	cmd.SetEnableCompletion(true)

	cmd1 := NewCmd("cmd1", "c1", flag.ExitOnError)
	cmd1.String("str", "s", "", "test string")

	if err := cmd.AddSubCmd(cmd1); err != nil {
		t.Fatal(err)
	}

	// 解析命令行参数，指定PowerShell补全类型
	if err := cmd.Parse([]string{"-gsc", "pwsh"}); err != nil {
		t.Fatal(err)
	}
}

// TestCompletionHelp 测试启用自动补全的帮助信息
func TestCompletionHelp(t *testing.T) {
	// 新建根命令
	cmd := NewCmd("root", "r", flag.ExitOnError)
	cmd.SetExitOnBuiltinFlags(false) // 禁止在解析命令行参数时退出
	cmd.SetUseChinese(true)          // 设置使用中文

	cmd.SetEnableCompletion(true)

	// 解析命令行参数，指定帮助信息
	if err := cmd.Parse([]string{"-h"}); err != nil {
		t.Fatal(err)
	}
}

// TestCompletionShellNone 测试ShellNone模式下补全功能的行为
func TestCompletionShellNone(t *testing.T) {
	// 测试场景1: ShellNone模式下不应生成任何补全脚本
	t.Run("no_completion_script_generated", func(t *testing.T) {
		// 创建命令实例
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		cmd.SetExitOnBuiltinFlags(false)
		cmd.SetEnableCompletion(true)

		// 重定向标准输出以捕获可能的输出
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// 解析命令行参数，指定shell为none
		err := cmd.Parse([]string{"-gsc", "none"})
		if err != nil {
			t.Fatalf("解析参数失败: %v", err)
		}

		// 恢复标准输出
		if err = w.Close(); err != nil {
			t.Errorf("Failed to close writer: %v", err)
		}
		os.Stdout = oldStdout

		// 读取捕获的输出
		var buf bytes.Buffer
		// 修复: 检查io.Copy的错误返回值
		_, err = io.Copy(&buf, r)
		if err != nil {
			t.Fatalf("读取输出失败: %v", err)
		}
		output := buf.String()

		// 验证没有生成补全脚本
		if output != "" {
			t.Errorf("ShellNone模式下不应有输出，实际输出: %q", output)
		}
	})

	// 测试场景2: 验证exitOnBuiltinFlags标志的行为
	t.Run("exit_on_builtin_flags_behavior", func(t *testing.T) {
		// 创建命令实例
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		cmd.SetEnableCompletion(true)
		if err := cmd.completionShell.Set(flags.ShellNone); err != nil {
			t.Fatalf("设置Shell类型失败: %v", err)
		}

		// 测试当exitOnBuiltinFlags为true时不应返回退出信号（ShellNone模式）
		cmd.SetExitOnBuiltinFlags(true)
		shouldExit, err := cmd.handleBuiltinFlags()
		if err != nil {
			t.Fatalf("处理内置标志失败: %v", err)
		}
		if shouldExit {
			t.Error("当shell为ShellNone且exitOnBuiltinFlags为true时，不应返回需要退出")
		}

		// 测试当exitOnBuiltinFlags为false时不应返回退出信号
		cmd.SetExitOnBuiltinFlags(false)
		shouldExit, err = cmd.handleBuiltinFlags()
		if err != nil {
			t.Fatalf("处理内置标志失败: %v", err)
		}
		if shouldExit {
			t.Error("当exitOnBuiltinFlags为false时，不应返回需要退出")
		}
	})
}
