// Package main 提供CmdSpec使用示例
//
// 本示例展示了如何使用CmdSpec结构体创建命令,
// 包括基本用法、子命令添加和互斥组设置。
package main

import (
	"fmt"
	"log"

	"gitee.com/MM-Q/qflag/internal/cmd"
	"gitee.com/MM-Q/qflag/internal/types"
)

func main() {
	// 基本用法示例
	basicExample()

	// 嵌套子命令示例
	nestedExample()

	// 复杂配置示例
	complexExample()
}

// basicExample 基本用法示例
func basicExample() {
	fmt.Println("=== 基本用法示例 ===")

	// 创建命令规格
	appSpec := cmd.NewCmdSpec("myapp", "app")
	appSpec.Desc = "我的应用程序"
	appSpec.Version = "1.0.0"
	appSpec.UseChinese = true
	appSpec.EnvPrefix = "MYAPP"
	appSpec.RunFunc = func(c types.Command) error {
		fmt.Println("运行应用程序")

		// 获取标志值
		if inputFile, exists := c.GetFlag("input-file"); exists {
			fmt.Printf("输入文件: %s\n", inputFile.GetStr())
		}

		if verboseFlag, exists := c.GetFlag("verbose-mode"); exists && verboseFlag.IsSet() {
			fmt.Println("详细模式已启用")
		}

		return nil
	}

	// 创建命令
	app, err := cmd.NewCmdFromSpec(appSpec)
	if err != nil {
		log.Fatal(err)
	}

	// 添加标志
	app.String("input-file", "i", "输入文件", "")
	app.Bool("verbose-mode", "V", "详细输出", false)

	// 解析并运行
	err = app.ParseAndRoute([]string{"--input-file", "test.txt", "--verbose-mode"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

// nestedExample 嵌套子命令示例
func nestedExample() {
	fmt.Println("=== 嵌套子命令示例 ===")

	// 创建主命令规格
	mainSpec := cmd.NewCmdSpec("main", "m")
	mainSpec.Desc = "主命令"
	mainSpec.RunFunc = func(c types.Command) error {
		fmt.Println("运行主命令")
		return nil
	}

	// 创建子命令
	subCmd := cmd.NewCmd("subcommand", "sub", types.ExitOnError)
	subCmd.SetDesc("子命令")
	subCmd.SetRun(func(c types.Command) error {
		fmt.Println("运行子命令")

		// 获取子命令选项
		if optionFlag, exists := c.GetFlag("sub-option"); exists {
			fmt.Printf("子命令选项: %s\n", optionFlag.GetStr())
		}

		return nil
	})

	// 添加到主命令
	mainSpec.SubCmds = []types.Command{subCmd}

	// 创建命令
	main, err := cmd.NewCmdFromSpec(mainSpec)
	if err != nil {
		log.Fatal(err)
	}

	// 获取子命令并添加标志
	retrievedSubCmd, _ := main.GetSubCmd("sub")
	subCmdInstance := retrievedSubCmd.(*cmd.Cmd)
	subCmdInstance.String("sub-option", "o", "子命令选项", "")

	// 解析并运行子命令
	err = main.ParseAndRoute([]string{"subcommand", "--sub-option", "test-value"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

// complexExample 复杂配置示例
func complexExample() {
	fmt.Println("=== 复杂配置示例 ===")

	// 复杂命令配置
	complexSpec := cmd.NewCmdSpec("complex", "cpx")
	complexSpec.Desc = "复杂命令示例"
	complexSpec.Version = "2.0.0"
	complexSpec.UseChinese = true
	complexSpec.EnvPrefix = "COMPLEX"
	complexSpec.UsageSyntax = "[options] <args>"
	complexSpec.LogoText = "Complex Command v2.0.0"
	complexSpec.RunFunc = func(c types.Command) error {
		fmt.Println("运行复杂命令")

		// 获取标志值
		if inputFlag, exists := c.GetFlag("main-input"); exists {
			fmt.Printf("主输入文件: %s\n", inputFlag.GetStr())
		}

		if outputFlag, exists := c.GetFlag("main-output"); exists && outputFlag.IsSet() {
			fmt.Printf("主输出文件: %s\n", outputFlag.GetStr())
		}

		if verboseFlag, exists := c.GetFlag("main-verbose"); exists && verboseFlag.IsSet() {
			fmt.Println("主详细模式已启用")
		}

		if jsonFlag, exists := c.GetFlag("main-json"); exists && jsonFlag.IsSet() {
			fmt.Println("使用主JSON格式输出")
		}

		if xmlFlag, exists := c.GetFlag("main-xml"); exists && xmlFlag.IsSet() {
			fmt.Println("使用主XML格式输出")
		}

		if limitFlag, exists := c.GetFlag("main-limit"); exists && limitFlag.IsSet() {
			fmt.Printf("主处理限制: %s\n", limitFlag.GetStr())
		}

		return nil
	}

	// 添加示例
	complexSpec.Examples = map[string]string{
		"基本用法":   "complex --input file.txt",
		"详细模式":   "complex --input file.txt --verbose",
		"输出JSON": "complex --input file.txt --json",
	}

	// 添加注意事项
	complexSpec.Notes = []string{
		"输入文件必须存在",
		"输出目录必须可写",
		"处理大文件时请增加内存限制",
	}

	// 添加子命令
	// processCmd := cmd.NewCmd("process", "proc", types.ExitOnError)
	// processCmd.SetDesc("处理数据")
	// processCmd.String("input", "I", "输入文件", "")
	// processCmd.String("output", "O", "输出文件", "")
	// processCmd.Bool("verbose", "V", "详细输出", false)
	// processCmd.Bool("json", "J", "JSON格式输出", false)
	// processCmd.Bool("xml", "X", "XML格式输出", false)
	// processCmd.Int("limit", "L", "处理限制", 1000)
	// processCmd.AddMutexGroup("output_format", []string{"json", "xml"}, true)
	// processCmd.SetRun(func(c types.Command) error {
	// 	fmt.Println("处理数据")

	// 	// 获取输入和输出文件
	// 	if inputFlag, exists := c.GetFlag("input"); exists {
	// 		fmt.Printf("输入文件: %s\n", inputFlag.GetStr())
	// 	}

	// 	if outputFlag, exists := c.GetFlag("output"); exists && outputFlag.IsSet() {
	// 		fmt.Printf("输出文件: %s\n", outputFlag.GetStr())
	// 	}

	// 	// 检查输出格式
	// 	if jsonFlag, exists := c.GetFlag("json"); exists && jsonFlag.IsSet() {
	// 		fmt.Println("使用JSON格式输出")
	// 	} else if xmlFlag, exists := c.GetFlag("xml"); exists && xmlFlag.IsSet() {
	// 		fmt.Println("使用XML格式输出")
	// 	}

	// 	return nil
	// })

	// validateCmd := cmd.NewCmd("validate", "val", types.ExitOnError)
	// validateCmd.SetDesc("验证数据")
	// validateCmd.String("input", "I", "输入文件", "")
	// validateCmd.SetRun(func(c types.Command) error {
	// 	fmt.Println("验证数据")

	// 	// 获取输入文件
	// 	if inputFlag, exists := c.GetFlag("input"); exists {
	// 		fmt.Printf("验证文件: %s\n", inputFlag.GetStr())
	// 	}

	// 	return nil
	// })

	// complexSpec.SubCmds = []types.Command{processCmd, validateCmd}

	// 创建命令
	complex, err := cmd.NewCmdFromSpec(complexSpec)
	if err != nil {
		log.Fatal(err)
	}

	// 添加多个标志
	complex.String("main-input", "I", "主输入文件", "")
	complex.String("main-output", "O", "主输出文件", "")
	complex.Bool("main-verbose", "V", "主详细输出", false)
	complex.Bool("main-json", "J", "主JSON格式输出", false)
	complex.Bool("main-xml", "X", "主XML格式输出", false)
	complex.Int("main-limit", "L", "主处理限制", 1000)

	// 添加互斥组
	// complex.AddMutexGroup("main_output_format", []string{"main-json", "main-xml"}, true)

	// 解析并运行主命令
	err = complex.ParseAndRoute([]string{"--main-input", "data.txt", "--main-output", "result.json", "--main-json", "--main-verbose"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}
