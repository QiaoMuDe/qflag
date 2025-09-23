// Package cmd 命令结构体功能测试
// 本文件包含了Cmd结构体的单元测试，测试命令创建、解析、子命令管理等核心功能，
// 确保命令行处理逻辑的正确性和稳定性。
package cmd

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

// 测试嵌套子命令生成的帮助信息样式
func TestNestedCommandHelp(t *testing.T) {
	// 创建根命令
	rootCmd := NewCmd("myapp", "", flag.ContinueOnError)
	rootCmd.SetDescription("这是一个演示应用程序")
	rootCmd.SetVersion("1.0.0")
	rootCmd.SetUseChinese(true)
	rootCmd.SetExitOnBuiltinFlags(false)
	rootCmd.SetEnableCompletion(true)

	// 为根命令添加所有13种类型的标志
	rootCmd.Bool("verbose", "vv", false, "启用详细输出")
	rootCmd.String("config", "c", "config.json", "配置文件路径")
	rootCmd.Int("port", "p", 8080, "服务端口号")
	rootCmd.Int64("max-size", "", 1024000, "最大文件大小")
	rootCmd.Uint16("workers", "w", 10, "工作线程数")
	rootCmd.Uint32("buffer-size", "", 4096, "缓冲区大小")
	rootCmd.Uint64("memory-limit", "", 1073741824, "内存限制")
	rootCmd.Float64("timeout", "t", 30.5, "超时时间(秒)")
	rootCmd.Enum("log-level", "l", "info", "日志级别", []string{"debug", "info", "warn", "error"})
	rootCmd.Duration("retry-interval", "", 5*time.Second, "重试间隔")
	rootCmd.Time("start-time", "", "now", "开始时间")
	rootCmd.StringSlice("tags", "", []string{"default"}, "标签列表")
	rootCmd.Map("env", "e", map[string]string{"NODE_ENV": "production"}, "环境变量")

	// 创建5个二级命令
	level2Commands := make([]*Cmd, 5)
	level2Names := []string{"server", "client", "database", "monitor", "deploy"}
	level2Descriptions := []string{
		"服务器相关命令",
		"客户端相关命令",
		"数据库相关命令",
		"监控相关命令",
		"部署相关命令",
	}

	for i := 0; i < 5; i++ {
		level2Commands[i] = NewCmd(level2Names[i], "", flag.ContinueOnError)
		level2Commands[i].SetDescription(level2Descriptions[i])
		level2Commands[i].SetUseChinese(true)
		level2Commands[i].SetExitOnBuiltinFlags(false)

		// 为每个二级命令添加3个标志
		switch i {
		case 0: // server
			level2Commands[i].String("host", "hs", "localhost", "服务器主机地址")
			level2Commands[i].Int("port", "p", 3000, "服务器端口")
			level2Commands[i].Bool("ssl", "s", false, "启用SSL")
		case 1: // client
			level2Commands[i].String("endpoint", "e", "http://localhost:3000", "服务端点")
			level2Commands[i].Duration("timeout", "t", 30*time.Second, "请求超时")
			level2Commands[i].Int("retries", "r", 3, "重试次数")
		case 2: // database
			level2Commands[i].String("", "d", "mysql", "数据库驱动")
			level2Commands[i].String("dsn", "", "user:pass@tcp(localhost:3306)/db", "数据源名称")
			level2Commands[i].Int("max-connections", "m", 10, "最大连接数")
		case 3: // monitor
			level2Commands[i].Duration("interval", "i", 60*time.Second, "监控间隔")
			level2Commands[i].StringSlice("metrics", "m", []string{"cpu", "memory"}, "监控指标")
			level2Commands[i].Bool("alert", "a", true, "启用告警")
		case 4: // deploy
			level2Commands[i].String("target", "t", "production", "部署目标")
			level2Commands[i].Bool("dry-run", "d", false, "试运行模式")
			level2Commands[i].Map("vars", "vv", map[string]string{}, "部署变量")
		}

		// 将二级命令添加到根命令
		if err := rootCmd.AddSubCmd(level2Commands[i]); err != nil {
			t.Fatal(err)
		}
	}

	// 为每个二级命令创建3个三级命令
	for i, level2Cmd := range level2Commands {
		for j := 0; j < 3; j++ {
			level3Name := fmt.Sprintf("%s-sub%d", level2Names[i], j+1)
			level3Cmd := NewCmd(level3Name, "", flag.ContinueOnError)
			level3Cmd.SetDescription(fmt.Sprintf("%s的子命令%d", level2Descriptions[i], j+1))
			level3Cmd.SetUseChinese(true)
			level3Cmd.SetExitOnBuiltinFlags(false)

			// 为每个三级命令添加3个标志
			level3Cmd.String("input", "i", "", "输入文件")
			level3Cmd.String("output", "o", "", "输出文件")
			level3Cmd.Bool("force", "f", false, "强制执行")

			// 将三级命令添加到对应的二级命令
			if err := level2Cmd.AddSubCmd(level3Cmd); err != nil {
				t.Fatal(err)
			}
		}
	}

	// 解析命令行参数
	if err := rootCmd.Parse([]string{"-h"}); err != nil {
		t.Fatal(err)
	}

	// 打印根命令帮助信息
	fmt.Println("========== 根命令帮助信息 ==========")
	rootCmd.PrintHelp()

	// 打印部分二级命令帮助信息
	fmt.Println("\n========== Server命令帮助信息 ==========")
	level2Commands[0].PrintHelp()

	fmt.Println("\n========== Database命令帮助信息 ==========")
	level2Commands[2].PrintHelp()

	// 打印部分三级命令帮助信息
	serverSubCmds := level2Commands[0].SubCmds()
	if len(serverSubCmds) > 0 {
		fmt.Println("\n========== Server-Sub1命令帮助信息 ==========")
		serverSubCmds[0].PrintHelp()
	}

	fmt.Println("\n=============================")
}
