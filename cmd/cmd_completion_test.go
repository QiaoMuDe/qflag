package cmd

import (
	"flag"
	"fmt"
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
	if err := cmd.Parse([]string{"completion", "--shell", "bash"}); err != nil {
		t.Fatal(err)
	}
}

// TestComPletionHelp 测试自动补全生成
func TestComPletionHelp(t *testing.T) {
	// 新建根命令
	cmd := NewCmd("root", "r", flag.ExitOnError)
	cmd.SetEnableCompletion(true)    // 启用自动补全
	cmd.SetExitOnBuiltinFlags(false) // 禁止在解析命令行参数时退出
	cmd.SetUseChinese(true)          // 设置使用中文

	// 解析命令行参数
	if err := cmd.Parse([]string{}); err != nil {
		t.Fatal(err)
	}

	// 打印帮助信息
	fmt.Println("=====================================================")
	cmd.PrintHelp()
	fmt.Println("=====================================================")

	// 打印子命令帮助信息
	fmt.Println("=====================================================")
	cmd.SubCmds()[0].PrintHelp()
	fmt.Println("=====================================================")

	// 测试自动补全生成
	fmt.Println("=====================================================")
	fmt.Println(cmd.generateBashCompletion())
	fmt.Println("=====================================================")
}
