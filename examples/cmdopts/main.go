package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/internal/cmd"
	"gitee.com/MM-Q/qflag/internal/types"
)

func main() {
	// 创建命令
	myCmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)

	// 创建子命令
	listCmd := qflag.NewCmd("list", "ls", qflag.ExitOnError)
	listCmd.SetDesc("列出所有项目")

	addCmd := qflag.NewCmd("add", "a", qflag.ExitOnError)
	addCmd.SetDesc("添加新项目")

	// 创建选项
	opts := &cmd.CmdOpts{
		// 基本属性
		Desc: "我的应用程序",
		RunFunc: func(c types.Command) error {
			fmt.Println("执行主命令")
			return nil
		},

		// 配置选项
		Version:     "1.0.0",
		UseChinese:  true,
		EnvPrefix:   "MYAPP",
		UsageSyntax: "myapp [options] [args...]",
		LogoText:    "MyApp v1.0.0",

		// 示例和说明
		Examples: map[string]string{
			"基本用法": "myapp --help",
			"详细模式": "myapp --verbose",
			"列出项目": "myapp list",
		},
		Notes: []string{
			"所有选项都可以通过环境变量设置",
			"使用 --help 查看详细帮助",
		},

		// 子命令和互斥组
		SubCmds: []types.Command{listCmd, addCmd},
		MutexGroups: []types.MutexGroup{
			{
				Name:      "format",
				Flags:     []string{"json", "xml"},
				AllowNone: false,
			},
		},
	}

	// 应用选项
	if err := myCmd.ApplyOpts(opts); err != nil {
		fmt.Printf("应用选项失败: %v\n", err)
		os.Exit(1)
	}

	// 解析并执行
	if err := myCmd.ParseAndRoute(os.Args[1:]); err != nil {
		fmt.Printf("执行失败: %v\n", err)
		os.Exit(1)
	}
}
