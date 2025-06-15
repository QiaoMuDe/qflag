# qflag

qflag是一个Go语言命令行参数解析库，对标准库flag进行封装，提供更便捷的使用体验，自动实现长短标志互斥，并默认绑定-h/--help标志打印帮助信息。

## 项目地址

[https://gitee.com/MM-Q/qflag.git](https://gitee.com/MM-Q/qflag.git)

## 安装

```bash
go get gitee.com/MM-Q/qflag
```

## 特性

- 支持长短标志自动互斥
- 默认绑定-h/--help标志，自动生成帮助信息
- 支持字符串、整数、布尔等多种标志类型
- 支持子命令功能
- 允许自定义帮助内容和命令描述

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
	"gitee.com/MM-Q/qflag"
	"os"
)

func main() {
	// 使用全局实例创建主命令
	qflag.QCommandLine.Description = "主命令"

	// 创建子命令
	subCmd := qflag.NewCmd("sub", "s", flag.ExitOnError)
	subCmd.Description = "子命令"
	subCmd.String("config", "c", "config.json", "配置文件路径")

	// 添加子命令到全局实例
	qflag.AddSubCmd(subCmd)

	// 解析命令行参数
	if err := qflag.Parse(); err != nil {
		fmt.Println("参数解析错误:", err)
		os.Exit(1)
	}

	// 检查是否执行子命令
	if len(qflag.QCommandLine.Args()) > 0 && qflag.QCommandLine.Args()[0] == "sub" {
		// 获取子命令配置参数
		configFile := subCmd.String("config", "c", "config.json", "配置文件路径")
		fmt.Printf("执行子命令，配置文件路径: %s\n", *configFile)
	}
}
```





## API文档

详细API文档请参考[APIDOC.md](APIDOC.md)

## 许可证

[MIT](LICENSE)