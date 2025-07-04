package cmd

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
)

// TestMain 全局测试入口，控制非verbose模式下的输出重定向
func TestMain(m *testing.M) {
	flag.Parse() // 解析命令行参数
	// 保存原始标准输出和错误输出
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	var nullFile *os.File
	var err error

	// 非verbose模式下重定向到空设备
	if !testing.Verbose() {
		nullFile, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
		if err != nil {
			panic("无法打开空设备文件: " + err.Error())
		}
		os.Stdout = nullFile
		os.Stderr = nullFile
	}

	// 运行所有测试
	exitCode := m.Run()

	// 恢复原始输出
	if !testing.Verbose() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
		nullFile.Close()
	}

	os.Exit(exitCode)
}

// 测试嵌套子命令生成的帮助信息样式
func TestNestedCommandHelp(t *testing.T) {
	// 定义子命令
	cmd1 := NewCmd("cmd1", "c1", flag.ContinueOnError)
	cmd1.SetDescription("cmd1 的描述信息")
	cmd1.SetVersion("1.0.0")
	cmd2 := NewCmd("cmd2", "", flag.ContinueOnError)
	cmd2.SetUseChinese(true)
	cmd2.SetDescription("cmd2 的描述信息")
	cmd3 := NewCmd("cmd3", "c3", flag.ContinueOnError)
	cmd3.SetUseChinese(true)
	cmd3.SetDescription("cmd3 的描述信息")

	// 定义标志
	cmd1.String("file", "f", "", "file")
	cmd1.Enum("enum", "e", "e1", "e2", []string{"e1", "e2"})
	cmd2.Bool("bool", "b", false, "bool")
	cmd2.Time("time", "t", time.Time{}, "time")
	cmd3.Int("int", "i", 0, "int")
	cmd3.Float64("float", "f", 0.0, "float")

	// 添加子命令
	if err := cmd1.AddSubCmd(cmd2); err != nil {
		t.Fatal(err)
	}
	if err := cmd2.AddSubCmd(cmd3); err != nil {
		t.Fatal(err)
	}

	// 解析命令行参数
	if err := cmd1.Parse([]string{}); err != nil {
		t.Fatal(err)
	}

	// 分隔符
	fmt.Println("=============================")

	// 打印 cmd1 帮助信息
	cmd1.PrintHelp()

	// 分隔符
	fmt.Println("=============================")

	// 打印 cmd2 帮助信息
	cmd2.PrintHelp()

	// 分隔符
	fmt.Println("=============================")

	// 打印 cmd3 帮助信息
	cmd3.PrintHelp()

	// 分隔符
	fmt.Println("=============================")
}
