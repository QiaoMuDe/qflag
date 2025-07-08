package cmd

import (
	"flag"
	"fmt"
	"testing"
)

func TestComPletion(t *testing.T) {
	// 新建根命令
	cmd := NewCmd("root", "r", flag.ExitOnError)
	cmd.SetEnableCompletion(true)
	cmd.SetExitOnBuiltinFlags(false)
	cmd.SetUseChinese(true)

	// 解析命令行参数
	if err := cmd.Parse([]string{"completion", "--shell", "bash"}); err != nil {
		t.Fatal(err)
	}

	fmt.Println(cmd.Args())

	// 打印帮助信息
	//fmt.Println("=====================================================")
	//cmd.PrintHelp()
	//fmt.Println("=====================================================")

	// cmd.SubCmds()[0].Parse([]string{"test"})

	// // 打印子命令帮助信息
	// fmt.Println("=====================================================")
	// cmd.SubCmds()[0].PrintHelp()
	// fmt.Println("=====================================================")

	// 测试自动补全生成
	//fmt.Println("=====================================================")
	//fmt.Println(cmd.generateBashCompletion())
	//fmt.Println("=====================================================")
}
