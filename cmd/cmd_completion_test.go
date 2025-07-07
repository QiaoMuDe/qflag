package cmd

import (
	"flag"
	"fmt"
	"testing"
)

func TestComPletion(t *testing.T) {
	// 新建根命令
	cmd := NewCmd("root", "r", flag.ExitOnError).SetEnableCompletion(true)
	cmd.SetUseChinese(true)

	// 解析命令行参数
	if err := cmd.Parse([]string{"completion", "--shell", "bash"}); err != nil {
		t.Fatal(err)
	}

	fmt.Println(cmd.Args())

	fmt.Println(cmd.enableCompletion)

	// 获取自动补全子命令
	completionCmd := cmd.subCmdMap["completion"]
	if completionCmd == nil {
		completionCmd = cmd.subCmdMap["comp"]
	}
	if completionCmd != nil {
		fmt.Println(completionCmd.completionShell.Get())
	} else {
		fmt.Println("Completion subcommand not found")
	}

	// 打印帮助信息
	fmt.Println("=====================================================")
	//cmd.PrintHelp()
	fmt.Println("=====================================================")

	// cmd.SubCmds()[0].Parse([]string{"test"})

	// // 打印子命令帮助信息
	// fmt.Println("=====================================================")
	// cmd.SubCmds()[0].PrintHelp()
	// fmt.Println("=====================================================")

	// 测试自动补全生成
	fmt.Println("=====================================================")
	//fmt.Println(cmd.generateBashCompletion())
	fmt.Println("=====================================================")
}
