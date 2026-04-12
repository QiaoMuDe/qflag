// smart-suggestion - 智能纠错功能测试示例
//
// 该示例演示 qflag 的智能纠错功能，包括子命令纠错和标志纠错
//
// 使用方法:
//   go run main.go                    # 显示帮助
//   go run main.go cnfig              # 测试子命令纠错（cnfig -> config）
//   go run main.go cfg                # 测试短名称子命令
//   go run main.go config --verboose  # 测试标志纠错（verboose -> verbose）
//   go run main.go config -o          # 测试短标志
//   go run main.go --help             # 显示帮助

package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建根命令
	root := qflag.NewCmd("smart-suggestion", "ss", qflag.ContinueOnError)
	root.SetDesc("智能纠错功能测试工具")

	// 添加全局标志
	_ = root.Bool("debug", "d", "启用调试模式", false)
	_ = root.String("output", "o", "输出文件", "")

	// 创建 config 子命令（带长短名称）
	configCmd := qflag.NewCmd("config", "cfg", qflag.ContinueOnError)
	configCmd.SetDesc("配置管理命令")
	_ = configCmd.Bool("verbose", "v", "详细输出", false)
	_ = configCmd.String("format", "f", "输出格式", "json")
	configCmd.SetRun(func(cmd qflag.Command) error {
		fmt.Println("执行 config 命令")
		if flag, exists := configCmd.GetFlag("verbose"); exists && flag.IsSet() {
			fmt.Println("  详细模式: 开启")
		}
		if formatFlag, exists := configCmd.GetFlag("format"); exists {
			fmt.Printf("  输出格式: %s\n", formatFlag.GetStr())
		}
		return nil
	})

	// 创建 status 子命令（带长短名称）
	statusCmd := qflag.NewCmd("status", "st", qflag.ContinueOnError)
	statusCmd.SetDesc("查看状态命令")
	_ = statusCmd.Bool("watch", "w", "持续监控", false)
	statusCmd.SetRun(func(cmd qflag.Command) error {
		fmt.Println("执行 status 命令")
		return nil
	})

	// 创建 init 子命令（只有长名称）
	initCmd := qflag.NewCmd("init", "", qflag.ContinueOnError)
	initCmd.SetDesc("初始化命令")
	initCmd.SetRun(func(cmd qflag.Command) error {
		fmt.Println("执行 init 命令")
		return nil
	})

	// 添加子命令到根命令
	if err := root.AddSubCmds(configCmd, statusCmd, initCmd); err != nil {
		fmt.Fprintf(os.Stderr, "添加子命令失败: %v\n", err)
		os.Exit(1)
	}

	// 设置根命令的运行函数（当没有子命令时执行）
	root.SetRun(func(cmd qflag.Command) error {
		fmt.Println("智能纠错功能测试工具")
		fmt.Println()
		fmt.Println("可用命令:")
		fmt.Println("  config, cfg    配置管理")
		fmt.Println("  status, st     查看状态")
		fmt.Println("  init           初始化")
		fmt.Println()
		fmt.Println("测试建议:")
		fmt.Println("  1. 输入 'cnfig' 查看子命令纠错（推荐 config）")
		fmt.Println("  2. 输入 'config --verboose' 查看标志纠错（推荐 --verbose）")
		fmt.Println("  3. 输入 'st' 查看短名称子命令")
		return nil
	})

	// 解析并执行命令
	if err := root.ParseAndRoute(os.Args[1:]); err != nil {
		// 打印错误信息（智能纠错功能已格式化）
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
