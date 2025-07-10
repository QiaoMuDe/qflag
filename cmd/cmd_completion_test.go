package cmd

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// TestCompletionPerformance 测试补全脚本生成性能
func TestCompletionPerformance(t *testing.T) {
	// 创建复杂命令结构
	rootCmd := NewCmd("root", "r", flag.ExitOnError)
	rootCmd.SetExitOnBuiltinFlags(false)
	if err := rootCmd.SetEnableCompletion(true); err != nil {
		t.Fatal(err)
	}

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
	_, err = rootCmd.SubCmdMap()["comp"].generateBashCompletion()
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
	_, err = rootCmd.SubCmdMap()["comp"].generatePwshCompletion()
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

	if err := cmd.SetEnableCompletion(true); err != nil {
		t.Fatal(err)
	}

	// 解析命令行参数
	if err := cmd.Parse([]string{"completion", "-s", "bash"}); err != nil {
		t.Fatal(err)
	}
}

// TestCompletionPwsh 测试PowerShell自动补全生成
func TestCompletionPwsh(t *testing.T) {
	// 新建根命令
	cmd := NewCmd("root", "r", flag.ExitOnError)
	cmd.SetExitOnBuiltinFlags(false) // 禁止在解析命令行参数时退出
	cmd.SetUseChinese(true)          // 设置使用中文

	if err := cmd.SetEnableCompletion(true); err != nil {
		t.Fatal(err)
	}

	// 解析命令行参数，指定PowerShell补全类型
	if err := cmd.Parse([]string{"comp", "-s", "pwsh"}); err != nil {
		t.Fatal(err)
	}
}
