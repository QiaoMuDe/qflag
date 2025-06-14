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

### 基本用法

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建命令实例
	cmd := qflag.NewCmd("demo", "d", flag.ExitOnError)
	cmd.Description = "qflag示例程序"

	// 添加标志
	nameFlag := cmd.String("name", "n", "", "姓名")
	ageFlag := cmd.Int("age", "a", 0, "年龄")
	verboseFlag := cmd.Bool("verbose", "v", false, "详细输出")

	// 解析命令行参数
	if err := cmd.Parse(os.Args[1:]); err != nil {
		fmt.Println("参数解析错误:", err)
		return
	}

	// 使用标志值
	fmt.Printf("姓名: %s, 年龄: %d, 详细输出: %v\n", *nameFlag.value, *ageFlag.value, *verboseFlag.value)
}
```

### 子命令示例

```go
package main

import (
	"flag"
	"fmt"
	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建主命令
	mainCmd := qflag.NewCmd("main", "m", flag.ExitOnError)
	mainCmd.Description = "主命令"

	// 创建子命令
	subCmd := qflag.NewCmd("sub", "s", flag.ExitOnError)
	subCmd.Description = "子命令"
	subCmd.String("config", "c", "config.json", "配置文件路径")

	// 添加子命令到主命令
	mainCmd.AddSubCmd(subCmd)

	// 解析命令行参数
	if err := mainCmd.Parse(os.Args[1:]); err != nil {
		fmt.Println("参数解析错误:", err)
		return
	}
}
```





## API文档

详细API文档请参考[APIDOC.md](APIDOC.md)

## 许可证

[MIT](LICENSE)