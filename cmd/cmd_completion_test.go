package cmd

import (
	"flag"
	"testing"
)

// TestComPletion 测试自动补全生成
func TestComPletion(t *testing.T) {
	// 新建根命令
	cmd := NewCmd("root", "r", flag.ExitOnError)
	cmd.SetEnableCompletion(true)    // 启用自动补全
	cmd.SetExitOnBuiltinFlags(false) // 禁止在解析命令行参数时退出
	cmd.SetUseChinese(true)          // 设置使用中文

	// 解析命令行参数
	if err := cmd.Parse([]string{"comp", "-s", "bash"}); err != nil {
		t.Fatal(err)
	}
}

// TestCompletionPwsh 测试PowerShell自动补全生成
func TestCompletionPwsh(t *testing.T) {
	// 新建根命令
	cmd := NewCmd("root", "r", flag.ExitOnError)
	cmd.SetEnableCompletion(true)    // 启用自动补全
	cmd.SetExitOnBuiltinFlags(false) // 禁止在解析命令行参数时退出
	cmd.SetUseChinese(true)          // 设置使用中文

	// 解析命令行参数，指定PowerShell补全类型
	if err := cmd.Parse([]string{"comp", "-s", "pwsh"}); err != nil {
		t.Fatal(err)
	}
}
