package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/internal/cmd"
	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

func main() {
	// 创建命令
	app := cmd.NewCmd("ssl-server", "ss", types.ContinueOnError)
	app.SetDesc("SSL服务器配置工具, 演示标志依赖关系功能")

	// 添加SSL开关标志
	if err := app.AddFlag(flag.NewBoolFlag("ssl", "s", "启用SSL/TLS加密", false)); err != nil {
		fmt.Printf("添加ssl标志失败: %v\n", err)
		os.Exit(1)
	}

	// 添加证书文件标志
	if err := app.AddFlag(flag.NewStringFlag("cert", "c", "SSL证书文件路径", "")); err != nil {
		fmt.Printf("添加cert标志失败: %v\n", err)
		os.Exit(1)
	}

	// 添加私钥文件标志
	if err := app.AddFlag(flag.NewStringFlag("key", "k", "SSL私钥文件路径", "")); err != nil {
		fmt.Printf("添加key标志失败: %v\n", err)
		os.Exit(1)
	}

	// 添加CA证书标志
	if err := app.AddFlag(flag.NewStringFlag("ca-cert", "a", "CA证书文件路径", "")); err != nil {
		fmt.Printf("添加ca-cert标志失败: %v\n", err)
		os.Exit(1)
	}

	// 添加端口标志
	if err := app.AddFlag(flag.NewIntFlag("port", "p", "服务器端口", 8080)); err != nil {
		fmt.Printf("添加port标志失败: %v\n", err)
		os.Exit(1)
	}

	// 添加调试模式标志
	if err := app.AddFlag(flag.NewBoolFlag("debug", "d", "启用调试模式", false)); err != nil {
		fmt.Printf("添加debug标志失败: %v\n", err)
		os.Exit(1)
	}

	// 添加标志依赖关系
	// 当启用SSL时, 必须提供证书和私钥
	if err := app.AddFlagDependency("ssl_requires_cert", "ssl", []string{"cert", "key"}, qflag.DepRequired); err != nil {
		fmt.Printf("添加ssl_requires_cert依赖失败: %v\n", err)
		os.Exit(1)
	}

	// 当启用调试模式时, 不能使用SSL (互斥关系)
	if err := app.AddFlagDependency("debug_mutex_ssl", "debug", []string{"ssl"}, qflag.DepMutex); err != nil {
		fmt.Printf("添加debug_mutex_ssl依赖失败: %v\n", err)
		os.Exit(1)
	}

	// 解析参数
	err := app.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("参数错误: %v\n", err)
		os.Exit(1)
	}

	// 处理逻辑
	fmt.Println("参数解析成功！")
	fmt.Println()

	// 显示配置
	if sslFlag, _ := app.GetFlag("ssl"); sslFlag.IsSet() {
		fmt.Println("SSL: 已启用")
		if certFlag, _ := app.GetFlag("cert"); certFlag.IsSet() {
			fmt.Printf("  证书文件: %s\n", certFlag.GetStr())
		}
		if keyFlag, _ := app.GetFlag("key"); keyFlag.IsSet() {
			fmt.Printf("  私钥文件: %s\n", keyFlag.GetStr())
		}
		if caCertFlag, _ := app.GetFlag("ca-cert"); caCertFlag.IsSet() {
			fmt.Printf("  CA证书: %s\n", caCertFlag.GetStr())
		}
	} else {
		fmt.Println("SSL: 未启用")
	}

	if debugFlag, _ := app.GetFlag("debug"); debugFlag.IsSet() {
		fmt.Println("调试模式: 已启用")
	} else {
		fmt.Println("调试模式: 未启用")
	}

	if portFlag, _ := app.GetFlag("port"); portFlag.IsSet() {
		fmt.Printf("端口: %s\n", portFlag.GetStr())
	}

	// 显示所有依赖关系
	fmt.Println()
	fmt.Println("已配置的标志依赖关系:")
	for _, dep := range app.FlagDependencies() {
		depTypeStr := "互斥"
		if dep.Type == qflag.DepRequired {
			depTypeStr = "必需"
		}
		fmt.Printf("  - %s: 当 --%s 设置时, %s --%v\n",
			dep.Name, dep.Trigger, depTypeStr, dep.Targets)
	}
}
