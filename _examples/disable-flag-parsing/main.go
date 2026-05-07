// disable-flag-parsing 示例
//
// 本示例演示了 qflag 的 DisableFlagParsing 功能，该功能允许命令完全跳过标志解析，
// 将所有参数（包括 --flag 和 -f 形式）都作为位置参数处理。
//
// 使用场景：
//   - 包装外部命令（如 kubectl exec、docker run）
//   - Shell 脚本包装器
//   - 需要透传参数给子进程的场景
//
// 运行示例：
//
//	go run main.go exec mypod -- ls -la
//	go run main.go ssh user@host -- echo hello
//	go run main.go normal --verbose
package main

import (
	"fmt"
	"os"
	"strings"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建根命令
	root := qflag.NewCmd("demo", "d", qflag.ExitOnError)
	root.SetDesc("禁用标志解析功能演示")
	root.SetVersion("1.0.0")
	root.SetRun(func(cmd qflag.Command) error {
		// 没有子命令时显示帮助
		cmd.PrintHelp()
		return nil
	})

	// 添加全局标志（仅用于演示，实际不会被 exec 子命令解析）
	root.String("namespace", "n", "命名空间", "default")
	root.Bool("verbose", "v", "详细输出", false)

	// ========== 子命令 1: exec ==========
	// 模拟 kubectl exec 的功能，透传参数给容器
	execCmd := qflag.NewCmd("exec", "e", qflag.ExitOnError)
	execCmd.SetDesc("在容器中执行命令（透传所有参数）")
	execCmd.SetDisableFlagParsing(true) // 关键：禁用标志解析
	execCmd.SetRun(func(cmd qflag.Command) error {
		args := cmd.Args()
		fmt.Println("=== exec 子命令 ===")
		fmt.Printf("禁用标志解析: %v\n", cmd.IsDisableFlagParsing())
		fmt.Printf("原始参数: %v\n", args)
		fmt.Printf("参数个数: %d\n", len(args))
		fmt.Println()

		if len(args) == 0 {
			fmt.Println("用法: demo exec <pod名> [--] <命令> [参数...]")
			fmt.Println()
			fmt.Println("示例:")
			fmt.Println("  demo exec mypod -- ls -la")
			fmt.Println("  demo exec mypod -c mycontainer -- ps aux")
			return nil
		}

		// 查找 -- 分隔符
		separatorIdx := -1
		for i, arg := range args {
			if arg == "--" {
				separatorIdx = i
				break
			}
		}

		if separatorIdx == -1 {
			// 没有 --，第一个参数是 pod 名，其余是要执行的命令
			podName := args[0]
			command := args[1:]
			fmt.Printf("Pod 名称: %s\n", podName)
			fmt.Printf("执行命令: %v\n", command)
		} else {
			// 有 --，-- 前面是 pod 和选项，后面是要执行的命令
			podAndOpts := args[:separatorIdx]
			command := args[separatorIdx+1:]
			fmt.Printf("Pod 和选项: %v\n", podAndOpts)
			fmt.Printf("执行命令: %v\n", command)
		}

		fmt.Println("\n注意: 所有参数都作为位置参数处理，--namespace 等标志不会被解析")
		return nil
	})

	// ========== 子命令 2: ssh ==========
	// 模拟 SSH 包装器，透传所有参数给 SSH
	sshCmd := qflag.NewCmd("ssh", "s", qflag.ExitOnError)
	sshCmd.SetDesc("SSH 包装器（透传所有参数）")
	sshCmd.SetDisableFlagParsing(true)
	sshCmd.SetRun(func(cmd qflag.Command) error {
		args := cmd.Args()
		fmt.Println("=== ssh 子命令 ===")
		fmt.Printf("禁用标志解析: %v\n", cmd.IsDisableFlagParsing())
		fmt.Printf("原始参数: %v\n", args)
		fmt.Println()

		if len(args) == 0 {
			fmt.Println("用法: demo ssh <user@host> [--] <命令>")
			fmt.Println()
			fmt.Println("示例:")
			fmt.Println("  demo ssh user@192.168.1.1 -- ls -la")
			fmt.Println("  demo ssh user@host -p 2222 -- ps aux")
			return nil
		}

		// 模拟执行 SSH 命令
		sshArgs := append([]string{"ssh"}, args...)
		fmt.Printf("将执行: %s\n", strings.Join(sshArgs, " "))
		fmt.Println("(实际执行已跳过，仅演示)")

		return nil
	})

	// ========== 子命令 3: normal ==========
	// 正常解析标志的子命令（对比用）
	normalCmd := qflag.NewCmd("normal", "n", qflag.ExitOnError)
	normalCmd.SetDesc("正常解析标志（对比用）")
	// 不禁用标志解析（默认行为）
	normalCmd.Bool("verbose", "v", "详细输出", false)
	normalCmd.String("output", "o", "输出文件", "")
	normalCmd.SetRun(func(cmd qflag.Command) error {
		fmt.Println("=== normal 子命令 ===")
		fmt.Printf("禁用标志解析: %v\n", cmd.IsDisableFlagParsing())
		// 从 FlagRegistry 获取标志值
		if vFlag, exists := normalCmd.GetFlag("verbose"); exists {
			fmt.Printf("verbose: %v\n", vFlag.GetStr())
		}
		if oFlag, exists := normalCmd.GetFlag("output"); exists {
			fmt.Printf("output: %v\n", oFlag.GetStr())
		}
		fmt.Printf("位置参数: %v\n", cmd.Args())
		fmt.Println()
		fmt.Println("注意: 此命令正常解析 --verbose 和 --output 标志")
		return nil
	})

	// ========== 子命令 4: wrapper ==========
	// 使用 ApplyOpts 批量配置的示例
	wrapperCmd := qflag.NewCmd("wrapper", "w", qflag.ExitOnError)
	wrapperCmd.SetDesc("使用 ApplyOpts 配置")
	// 通过方法链设置
	wrapperCmd.SetDisableFlagParsing(true)
	wrapperCmd.SetRun(func(c qflag.Command) error {
		args := c.Args()
		fmt.Println("=== wrapper 子命令 (通过方法链配置) ===")
		fmt.Printf("禁用标志解析: %v\n", c.IsDisableFlagParsing())
		fmt.Printf("原始参数: %v\n", args)
		fmt.Println()
		fmt.Println("此命令通过 SetDisableFlagParsing(true) 禁用标志解析")
		return nil
	})

	// 添加所有子命令
	if err := root.AddSubCmds(execCmd, sshCmd, normalCmd, wrapperCmd); err != nil {
		fmt.Fprintf(os.Stderr, "添加子命令失败: %v\n", err)
		os.Exit(1)
	}

	// 添加示例说明
	root.AddExample("exec 子命令（透传参数）", "demo exec mypod -- ls -la")
	root.AddExample("ssh 子命令（透传参数）", "demo ssh user@host -- ps aux")
	root.AddExample("normal 子命令（解析标志）", "demo normal --verbose --output=file.txt arg1 arg2")
	root.AddExample("wrapper 子命令（方法链配置）", "demo wrapper --any-flag value")

	// 执行
	if err := root.ParseAndRoute(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
