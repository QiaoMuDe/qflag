# qflag

qflag是一个功能强大的Go语言命令行参数解析库，基于标准库flag进行封装，提供更便捷、更灵活的命令行参数处理体验。主要特点包括长短标志绑定、内置帮助系统、子命令支持以及线程安全的标志操作。

## 项目地址

[https://gitee.com/MM-Q/qflag.git](https://gitee.com/MM-Q/qflag.git)

## 安装

```bash
go get gitee.com/MM-Q/qflag
```

## 特性

- 支持长标志和短标志绑定
- 自动生成主命令及其子命令的帮助信息，支持多级嵌套子命令
- 默认绑定-h/--help标志，自动生成格式化的帮助信息
- 支持字符串、整数、布尔、浮点数等多种标志类型
- 线程安全的标志操作，支持并发解析
- 完善的子命令系统，支持循环引用检测
- 允许自定义帮助内容和命令描述
- 内置安装路径显示功能(-sip/--show-install-path)
- 支持标志默认值动态设置
- 支持枚举类型标志（EnumFlag），可限制输入值为预定义选项
- 标志互斥组功能，允许定义不能同时使用的标志组
- 提供丰富的API文档和测试用例

## 使用示例
### 基本使用示例

```go
package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 定义一个启动app的Bool型标志
	runF := qflag.Bool("run", "r", false, "run app")
	nameF := qflag.String("name", "n", "", "app name")
	pathF := qflag.String("path", "p", qflag.GetExecutablePath(), "app path")

	// 解析参数
	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	// 获取参数值
	if runF.GetValue() {
		fmt.Printf("启动app: %s\n", nameF.GetValue())
		fmt.Printf("app路径: %s\n", pathF.GetValue())
	}
}

```

### 子命令示例

```go
package main

import (
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 定义一个启动app的Bool型标志
	runF := qflag.Bool("run", "r", false, "运行app")
	nameF := qflag.String("name", "n", "", "指定app名称")
	pathF := qflag.String("path", "p", qflag.GetExecutablePath(), "指定app路径")

	// 定义子命令
	stopCmd := qflag.NewCmd("stop", "st", flag.ExitOnError)
	stopF := stopCmd.Bool("stop", "s", false, "停止app")
	stopCmdName := stopCmd.String("name", "n", "", "指定app名称")

	// 添加子命令描述
	stopCmd.SetDescription("停止app")

	// 添加子命令到全局QCommandLine
	qflag.AddSubCmd(stopCmd)

	// 解析命令
	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	// 获取参数值
	if runF.GetValue() {
		fmt.Printf("启动app: %s\n", nameF.GetValue())
		fmt.Printf("app路径: %s\n", pathF.GetValue())
	}

	// 判断是否执行了stop子命令
	if stopF.GetValue() {
		fmt.Printf("停止app: %s\n", stopCmdName.GetValue())
	}
}
```

## API文档

详细API文档请参考[APIDOC.md](APIDOC.md)

## 许可证

[MIT](LICENSE)