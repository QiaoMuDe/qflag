package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag/internal/cmd"
	"gitee.com/MM-Q/qflag/internal/types"
)

func main() {
	// 创建根命令
	root := cmd.NewCmd("example", "ex", types.ContinueOnError)
	root.SetDesc("内置标志示例程序")
	root.SetVersion("1.0.0")
	root.SetChinese(true) // 使用中文

	// 添加一些自定义标志
	nameFlag := root.String("name", "n", "用户名", "guest")
	ageFlag := root.Int("age", "a", "年龄", 18)
	debugFlag := root.Bool("debug", "d", "调试模式", false)

	// 设置运行函数
	root.SetRun(func(c types.Command) error {
		fmt.Println("=== 程序运行 ===")
		fmt.Printf("用户名: %s\n", nameFlag.Get())
		fmt.Printf("年龄: %d\n", ageFlag.Get())
		fmt.Printf("调试模式: %t\n", debugFlag.Get())
		return nil
	})

	// 解析并执行
	err := root.ParseAndRoute(os.Args[1:])
	if err != nil {
		fmt.Println("错误:", err.Error())
	}
}
