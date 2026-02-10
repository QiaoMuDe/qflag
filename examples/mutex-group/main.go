package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag/internal/cmd"
	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

func main() {
	// 创建命令
	app := cmd.NewCmd("format-converter", "fc", types.ContinueOnError)
	app.SetDesc("一个格式转换工具, 演示互斥组功能")

	// 添加输出格式标志 (互斥组)
	if err := app.AddFlag(flag.NewBoolFlag("json", "j", "输出为JSON格式", false)); err != nil {
		fmt.Printf("添加json标志失败: %v\n", err)
		os.Exit(1)
	}
	if err := app.AddFlag(flag.NewBoolFlag("xml", "x", "输出为XML格式", false)); err != nil {
		fmt.Printf("添加xml标志失败: %v\n", err)
		os.Exit(1)
	}
	if err := app.AddFlag(flag.NewBoolFlag("yaml", "y", "输出为YAML格式", false)); err != nil {
		fmt.Printf("添加yaml标志失败: %v\n", err)
		os.Exit(1)
	}

	// 添加输入源标志 (互斥组, 必须选择一个)
	if err := app.AddFlag(flag.NewStringFlag("file", "f", "从文件读取输入", "")); err != nil {
		fmt.Printf("添加file标志失败: %v\n", err)
		os.Exit(1)
	}
	if err := app.AddFlag(flag.NewStringFlag("url", "u", "从URL读取输入", "")); err != nil {
		fmt.Printf("添加url标志失败: %v\n", err)
		os.Exit(1)
	}
	if err := app.AddFlag(flag.NewBoolFlag("stdin", "s", "从标准输入读取", false)); err != nil {
		fmt.Printf("添加stdin标志失败: %v\n", err)
		os.Exit(1)
	}

	// 添加互斥组
	// 输出格式互斥组: 可以选择一种格式, 也可以都不选择 (使用默认格式)
	app.AddMutexGroup("output_format", []string{"json", "xml", "yaml"}, true)

	// 输入源互斥组: 必须选择一个输入源
	app.AddMutexGroup("input_source", []string{"file", "url", "stdin"}, false)

	// 解析参数
	err := app.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("参数错误: %v\n", err)
		os.Exit(1)
	}

	// 处理逻辑
	fmt.Println("参数解析成功！")

	// 检查输出格式
	if app.FlagRegistry().Has("json") {
		if flag, _ := app.GetFlag("json"); flag.IsSet() {
			fmt.Println("输出格式: JSON")
		}
	}
	if app.FlagRegistry().Has("xml") {
		if flag, _ := app.GetFlag("xml"); flag.IsSet() {
			fmt.Println("输出格式: XML")
		}
	}
	if app.FlagRegistry().Has("yaml") {
		if flag, _ := app.GetFlag("yaml"); flag.IsSet() {
			fmt.Println("输出格式: YAML")
		}
	}

	// 检查输入源
	if app.FlagRegistry().Has("file") {
		if flag, _ := app.GetFlag("file"); flag.IsSet() {
			fmt.Printf("输入源: 文件 (%s)\n", flag.GetStr())
		}
	}
	if app.FlagRegistry().Has("url") {
		if flag, _ := app.GetFlag("url"); flag.IsSet() {
			fmt.Printf("输入源: URL (%s)\n", flag.GetStr())
		}
	}
	if app.FlagRegistry().Has("stdin") {
		if flag, _ := app.GetFlag("stdin"); flag.IsSet() {
			fmt.Println("输入源: 标准输入")
		}
	}
}
