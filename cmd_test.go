// package qflag 命令结构体功能测试
// 本文件包含了Cmd结构体的单元测试，测试命令创建、解析、子命令管理等核心功能，
// 确保命令行处理逻辑的正确性和稳定性。
package qflag

import (
	"flag"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/qerr"
)

// 测试嵌套子命令生成的帮助信息样式
func TestNestedCommandHelp(t *testing.T) {
	// 创建根命令
	rootCmd := NewCmd("myapp", "", flag.ContinueOnError)
	rootCmd.SetDesc("这是一个演示应用程序")
	rootCmd.SetVersion("1.0.0")
	rootCmd.SetChinese(true)
	rootCmd.SetNoFgExit(true)
	rootCmd.SetCompletion(true)

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
		level2Commands[i].SetDesc(level2Descriptions[i])
		level2Commands[i].SetChinese(true)
		level2Commands[i].SetNoFgExit(false)

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
			level3Cmd.SetDesc(fmt.Sprintf("%s的子命令%d", level2Descriptions[i], j+1))
			level3Cmd.SetChinese(true)
			level3Cmd.SetNoFgExit(false)

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
	serverSubCmdMap := level2Commands[0].SubCmdMap()
	if len(serverSubCmdMap) > 0 {
		fmt.Println("\n========== Server-Sub1命令帮助信息 ==========")
		// 从SubCmdMap中获取第一个子命令
		for _, cmd := range serverSubCmdMap {
			cmd.PrintHelp()
			break // 只打印第一个子命令的帮助信息
		}
	}

	fmt.Println("\n=============================")
}

// TestCmdRunFunction 测试命令的run函数方法
func TestCmdRunFunction(t *testing.T) {
	// 创建测试命令
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.SetNoFgExit(true)

	// 测试1: 未解析命令时调用Run应该返回错误
	t.Run("未解析命令时调用Run", func(t *testing.T) {
		err := cmd.Run()
		if err == nil {
			t.Errorf("期望返回错误，但没有返回错误")
		}

		expectedErr := "validation failed: command must be parsed before execution"
		if err.Error() != expectedErr {
			t.Errorf("错误信息不匹配，期望: %s, 实际: %s", expectedErr, err.Error())
		}
	})

	// 测试2: 解析命令但未设置run函数时调用Run应该返回错误
	t.Run("解析命令但未设置run函数时调用Run", func(t *testing.T) {
		// 先解析命令
		if err := cmd.Parse([]string{}); err != nil {
			t.Fatalf("解析命令失败: %v", err)
		}

		err := cmd.Run()
		if err == nil {
			t.Errorf("期望返回错误，但没有返回错误")
		}

		expectedErr := "validation failed: no run function set for command"
		if err.Error() != expectedErr {
			t.Errorf("错误信息不匹配，期望: %s, 实际: %s", expectedErr, err.Error())
		}
	})

	// 测试2: 设置nil run函数应该panic
	t.Run("设置nil run函数应该panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("期望panic，但没有panic")
			} else {
				expectedPanic := "run function cannot be nil"
				if r != expectedPanic {
					t.Errorf("panic信息不匹配，期望: %s, 实际: %v", expectedPanic, r)
				}
			}
		}()

		cmd.SetRun(nil)
	})

	// 测试3: 设置正常run函数并执行
	t.Run("设置正常run函数并执行", func(t *testing.T) {
		executed := false
		testError := qerr.NewValidationError("test error")

		// 设置run函数
		cmd.SetRun(func(c *Cmd) error {
			// 验证传入的命令实例
			if c == nil {
				t.Errorf("期望传入非nil的Cmd实例")
				return nil
			}

			// 验证传入的命令是正确的实例
			if c.Name() != "test" {
				t.Errorf("命令名称不匹配，期望: test, 实际: %s", c.Name())
			}

			executed = true
			return testError
		})

		// 执行run函数
		err := cmd.Run()

		// 验证run函数被执行
		if !executed {
			t.Errorf("run函数未被执行")
		}

		// 验证返回的错误
		if err == nil {
			t.Errorf("期望返回错误，但没有返回错误")
		}

		if err != testError {
			t.Errorf("返回的错误不匹配，期望: %v, 实际: %v", testError, err)
		}
	})

	// 测试4: 测试run函数的并发安全性
	t.Run("测试run函数的并发安全性", func(t *testing.T) {
		var executedCount int64
		iterations := 100

		// 设置run函数，使用原子操作确保计数器的线程安全
		cmd.SetRun(func(c *Cmd) error {
			// 使用原子操作递增计数器
			atomic.AddInt64(&executedCount, 1)
			return nil
		})

		// 并发执行run函数
		done := make(chan bool, iterations)
		for i := 0; i < iterations; i++ {
			go func() {
				_ = cmd.Run() // 忽略错误，因为在这个测试中我们只关心执行次数
				done <- true
			}()
		}

		// 等待所有goroutine完成
		for i := 0; i < iterations; i++ {
			<-done
		}

		// 验证run函数被执行了正确的次数
		actualCount := atomic.LoadInt64(&executedCount)
		if actualCount != int64(iterations) {
			t.Errorf("run函数执行次数不匹配，期望: %d, 实际: %d", iterations, actualCount)
		}
	})

	// 测试5: 测试SetRun的并发安全性
	t.Run("测试SetRun的并发安全性", func(t *testing.T) {
		iterations := 50
		done := make(chan bool, iterations)

		// 并发设置run函数
		for i := 0; i < iterations; i++ {
			go func(index int) {
				cmd.SetRun(func(c *Cmd) error {
					// 简单验证，不执行复杂逻辑
					return nil
				})
				done <- true
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < iterations; i++ {
			<-done
		}

		// 验证最后一次设置仍然有效
		err := cmd.Run()
		if err != nil {
			t.Errorf("最后一次Run调用返回了意外错误: %v", err)
		}
	})

	// 测试6: 测试在子命令中使用run函数
	t.Run("测试在子命令中使用run函数", func(t *testing.T) {
		// 创建父命令
		parentCmd := NewCmd("parent", "p", flag.ContinueOnError)

		// 创建子命令
		childCmd := NewCmd("child", "c", flag.ContinueOnError)

		// 为子命令设置run函数
		childExecuted := false
		childCmd.SetRun(func(c *Cmd) error {
			childExecuted = true
			return nil
		})

		// 添加子命令到父命令
		if err := parentCmd.AddSubCmd(childCmd); err != nil {
			t.Fatalf("添加子命令失败: %v", err)
		}

		// 获取子命令并执行
		retrievedChild := parentCmd.GetSubCmd("child")
		if retrievedChild == nil {
			t.Fatalf("无法获取子命令")
		}

		// 执行子命令的run函数
		// 先解析子命令
		if err := retrievedChild.Parse([]string{}); err != nil {
			t.Fatalf("解析子命令失败: %v", err)
		}

		if err := retrievedChild.Run(); err != nil {
			t.Errorf("执行子命令run函数失败: %v", err)
		}

		// 验证子命令的run函数被执行
		if !childExecuted {
			t.Errorf("子命令的run函数未被执行")
		}
	})
}
