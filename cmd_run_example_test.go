package qflag

import (
	"fmt"
	"os"
	"strings"
)

// ExampleCmd_Run 演示如何使用Run函数字段手动执行命令逻辑
func ExampleCmd_Run() {
	// 创建命令
	cmd := NewCmd("server", "s", ExitOnError)

	// 定义标志
	port := cmd.Int("port", "p", 8080, "服务器端口")
	host := cmd.String("host", "H", "localhost", "服务器主机")
	debug := cmd.Bool("debug", "d", false, "调试模式")

	// 设置执行函数
	cmd.Run = func(c *Cmd) error {
		fmt.Printf("启动服务器: %s:%d (调试模式: %v)\n", host.Get(), port.Get(), debug.Get())
		// 这里可以放置实际的服务器启动逻辑
		return nil
	}

	// 模拟命令行参数
	args := []string{"--port", "3000", "--host", "0.0.0.0", "--debug"}

	// 解析参数（现在只解析，不自动执行）
	if err := cmd.Parse(args); err != nil {
		fmt.Printf("解析错误: %v\n", err)
		os.Exit(1)
	}

	// 手动执行Run函数
	if cmd.Run != nil {
		if err := cmd.Run(cmd); err != nil {
			fmt.Printf("执行错误: %v\n", err)
			os.Exit(1)
		}
	}

	// Output: 启动服务器: 0.0.0.0:3000 (调试模式: true)
}

// ExampleCmd_Run_subcommand 演示子命令的Run函数手动执行
func ExampleCmd_Run_subcommand() {
	// 创建根命令
	rootCmd := NewCmd("app", "a", ExitOnError)

	// 创建子命令 - start
	startCmd := NewCmd("start", "s", ExitOnError)
	port := startCmd.Int("port", "p", 8080, "服务端口")

	startCmd.Run = func(c *Cmd) error {
		fmt.Printf("启动服务，端口: %d\n", port.Get())
		return nil
	}

	// 添加子命令到根命令
	if err := rootCmd.AddSubCmd(startCmd); err != nil {
		fmt.Printf("添加子命令失败: %v\n", err)
		os.Exit(1)
	}

	// 测试start子命令 - 先解析，再手动执行
	args := []string{"start", "--port", "9000"}
	if err := rootCmd.Parse(args); err != nil {
		fmt.Printf("解析错误: %v\n", err)
		os.Exit(1)
	}
	// 手动执行start子命令的Run函数
	if startCmd.Run != nil {
		if err := startCmd.Run(startCmd); err != nil {
			fmt.Printf("执行错误: %v\n", err)
			os.Exit(1)
		}
	}

	// Output: 启动服务，端口: 9000
}

// ExampleCmd_Run_subcommand_stop 演示stop子命令的Run函数手动执行
func ExampleCmd_Run_subcommand_stop() {
	// 创建根命令
	rootCmd := NewCmd("app", "a", ExitOnError)

	// 创建子命令 - stop
	stopCmd := NewCmd("stop", "t", ExitOnError)
	force := stopCmd.Bool("force", "f", false, "强制停止")

	stopCmd.Run = func(c *Cmd) error {
		if force.Get() {
			fmt.Println("强制停止服务")
		} else {
			fmt.Println("优雅停止服务")
		}
		return nil
	}

	// 添加子命令到根命令
	if err := rootCmd.AddSubCmd(stopCmd); err != nil {
		fmt.Printf("添加子命令失败: %v\n", err)
		os.Exit(1)
	}

	// 测试stop子命令 - 先解析，再手动执行
	args := []string{"stop", "--force"}
	if err := rootCmd.Parse(args); err != nil {
		fmt.Printf("解析错误: %v\n", err)
		os.Exit(1)
	}
	// 手动执行stop子命令的Run函数
	if stopCmd.Run != nil {
		if err := stopCmd.Run(stopCmd); err != nil {
			fmt.Printf("执行错误: %v\n", err)
			os.Exit(1)
		}
	}

	// Output: 强制停止服务
}

// ExampleCmd_Run_error 演示手动执行Run函数时错误处理
func ExampleCmd_Run_error() {
	// 测试1: 空值应该返回错误
	cmd1 := NewCmd("test", "t", ContinueOnError)
	value1 := cmd1.String("value", "v", "", "测试值")

	cmd1.Run = func(c *Cmd) error {
		if strings.TrimSpace(value1.Get()) == "" {
			return fmt.Errorf("值不能为空")
		}
		fmt.Printf("处理值: %s\n", value1.Get())
		return nil
	}

	args1 := []string{"--value", "   "}
	if err := cmd1.Parse(args1); err != nil {
		fmt.Printf("解析错误: %v\n", err)
		return
	}
	// 手动执行Run函数并处理错误
	if cmd1.Run != nil {
		if err := cmd1.Run(cmd1); err != nil {
			fmt.Printf("执行错误: %v\n", err)
		}
	}
}

// ExampleCmd_Run_success 演示成功执行的情况
func ExampleCmd_Run_success() {
	cmd := NewCmd("test", "t", ContinueOnError)
	value := cmd.String("value", "v", "", "测试值")

	cmd.Run = func(c *Cmd) error {
		if strings.TrimSpace(value.Get()) == "" {
			return fmt.Errorf("值不能为空")
		}
		fmt.Printf("处理值: %s\n", value.Get())
		return nil
	}

	args := []string{"--value", "hello"}
	if err := cmd.Parse(args); err != nil {
		fmt.Printf("解析错误: %v\n", err)
		return
	}

	// 手动执行Run函数
	if cmd.Run != nil {
		if err := cmd.Run(cmd); err != nil {
			fmt.Printf("执行错误: %v\n", err)
		}
	}

	// Output: 处理值: hello
}

// ExampleCmd_Run_nil 演示Run函数为nil时手动检查并执行
func ExampleCmd_Run_nil() {
	cmd := NewCmd("normal", "n", ExitOnError)
	flag := cmd.Bool("flag", "f", false, "测试标志")

	// 不设置Run函数

	args := []string{"--flag"}
	if err := cmd.Parse(args); err != nil {
		fmt.Printf("解析错误: %v\n", err)
		os.Exit(1)
	}

	// 手动检查Run函数是否存在并执行
	if cmd.Run != nil {
		if err := cmd.Run(cmd); err != nil {
			fmt.Printf("执行错误: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("没有定义Run函数，跳过执行")
	}

	fmt.Printf("标志值: %v\n", flag.Get())

	// Output:
	// 没有定义Run函数，跳过执行
	// 标志值: true
}

// ExampleCmd_Run_manualWorkflow 演示完整的手动执行工作流程
func ExampleCmd_Run_manualWorkflow() {
	// 创建命令
	cmd := NewCmd("deploy", "d", ExitOnError)

	// 定义标志
	env := cmd.String("env", "e", "dev", "部署环境")
	tag := cmd.String("tag", "t", "latest", "镜像标签")
	dryRun := cmd.Bool("dry-run", "", false, "干运行模式")

	// 设置执行函数
	cmd.Run = func(c *Cmd) error {
		if dryRun.Get() {
			fmt.Printf("[干运行] 准备部署到 %s 环境，使用标签 %s\n", env.Get(), tag.Get())
			return nil
		}
		fmt.Printf("正在部署到 %s 环境，使用标签 %s\n", env.Get(), tag.Get())
		// 这里放置实际的部署逻辑
		return nil
	}

	// 模拟命令行参数
	args := []string{"--env", "prod", "--tag", "v1.2.3"}

	// 步骤1: 解析参数
	if err := cmd.Parse(args); err != nil {
		fmt.Printf("解析失败: %v\n", err)
		os.Exit(1)
	}

	// 步骤2: 验证参数（可选）
	if env.Get() == "" {
		fmt.Println("错误: 环境不能为空")
		os.Exit(1)
	}

	// 步骤3: 手动执行Run函数
	if cmd.Run != nil {
		if err := cmd.Run(cmd); err != nil {
			fmt.Printf("执行失败: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("警告: 未定义执行函数")
	}

	// Output: 正在部署到 prod 环境，使用标签 v1.2.3
}
