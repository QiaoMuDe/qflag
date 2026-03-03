package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	fmt.Println("=== 测试三层子命令递归解析 ===")
	fmt.Println()

	// 创建主命令
	mainCmd := qflag.NewCmd("main", "m", qflag.ExitOnError)
	mainCmd.SetDesc("主命令")
	mainCmd.SetRun(func(qflag.Command) error {
		fmt.Println("执行主命令")
		return nil
	})

	// 创建子命令
	subCmd := qflag.NewCmd("sub", "s", qflag.ExitOnError)
	subCmd.String("rr", "r", "111", "")
	subCmd.SetDesc("子命令")
	subCmd.SetRun(func(qflag.Command) error {
		fmt.Println("执行子命令")

		if flag, exists := subCmd.GetFlag("rr"); exists {
			fmt.Printf("子命令选项: %s\n", flag.GetStr())
		}
		return nil
	})

	// 创建子子命令
	subSubCmd := qflag.NewCmd("subsub", "ss", qflag.ExitOnError)
	// 添加标志
	subSubCmd.String("option", "o", "选项", "")
	subSubCmd.SetDesc("子子命令")
	subSubCmd.SetRun(func(qflag.Command) error {
		fmt.Println("执行子子命令")

		// 获取标志值
		if flag, exists := subSubCmd.GetFlag("option"); exists {
			fmt.Printf("选项值: %s\n", flag.GetStr())
		}

		return nil
	})

	// 添加子命令
	if err := mainCmd.AddSubCmds(subCmd); err != nil {
		fmt.Printf("添加子命令失败: %v\n", err)
		os.Exit(1)
	}

	if err := subCmd.AddSubCmds(subSubCmd); err != nil {
		fmt.Printf("添加子子命令失败: %v\n", err)
		os.Exit(1)
	}

	if err := mainCmd.ParseAndRoute(os.Args[1:]); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}
