# qflag

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/qflag)

qflag 是一个用于解析命令行参数的 Go 语言库。它提供了丰富的功能，包括多种类型的标志（flag）、子命令支持、帮助信息生成等。

## 项目地址

该项目托管在 Gitee 和 GitHub 上，您可以访问以下链接查看源代码和最新动态：

- [Gitee](https://gitee.com/MM-Q/qflag.git)
- [GitHub](https://github.com/QiaoMuDe/qflag.git)

## 安装

要使用 qflag，您需要先安装 Go 环境。然后可以通过以下命令安装：

```bash
go get -u gitee.com/MM-Q/qflag
```

## 特性

- 支持多种类型的标志：字符串、整数、布尔值、64位浮点数、枚举、切片、时间间隔、64位整数、无符号16位整数、映射、路径、时间、IPv4地址、IPv6地址和URL。
- 支持子命令。
- 提供内置的帮助信息标志和版本信息标志。
- 支持自定义帮助信息。
- 支持环境变量自动映射，可从环境变量加载配置并与命令行参数协同工作。
- 线程安全设计。
- 循环引用检测。
- 动态帮助信息生成。
- 标志命名规则。
- 内置参数验证器，支持路径存在性、数值范围等验证。

## 标志类型

以下是qflag支持的所有标志类型及其说明：

| 标志类型 | 描述 | 示例 | 注意事项 |
|----------|------|------|--------|
| String | 字符串类型标志 | `--name "example"` | 无 |
| Int | 整数类型标志 | `--port 8080` | 无 |
| Int64 | 64位整数类型标志 | `--size 1073741824` | 无 |
| Uint16 | 无符号16位整数类型标志 | `--timeout 300` | 无 |
| Bool | 布尔类型标志 | `--debug` | 无 |
| Float64 | 64位浮点数类型标志 | `--threshold 0.95` | 无 |
| Enum | 枚举类型标志 | `--mode "debug"` | 支持大小写敏感设置，通过`SetCaseSensitive(true)`启用 |
| Slice | 切片类型标志，支持自定义分隔符 | `--files file1.txt,file2.txt` | 默认使用逗号分隔，可通过`SetSeparator`修改 |
| Duration | 时间间隔类型标志 | `--timeout 30s` | 无 |
| Time | 时间类型标志 | `--start-time "2024-01-01T00:00:00"` | 无 |
| Map | 映射类型标志 | `--config key=value,key2=value2` | 支持`=`和`:`作为键值分隔符，可通过`SetDelimiters`修改 |
| Path | 路径类型标志，支持路径验证 | `--log-path "./logs"` | 可通过`IsDirectory(true)`和`MustExist(true)`设置验证规则 |
| IP4 | IPv4地址类型标志 | `--server-ip 192.168.1.1` | 无 |
| IP6 | IPv6地址类型标志 | `--server-ipv6 ::1` | 无 |
| URL | URL类型标志 | `--api-url https://api.example.com` |	无 |

## 使用示例

### SliceFlag 示例

展示如何使用切片标志并自定义分隔符：

```go
package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建切片标志，默认使用逗号分隔
	filesF := qflag.Slice("files", "f", []string{}, "要处理的文件列表")
	// 设置自定义分隔符为分号
	filesF.SetSeparator(";")

	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	// 获取切片值
	files := filesF.Get()
	fmt.Printf("要处理的文件: %v\n", files)
}
```

使用方式：`./app --files file1.txt;file2.txt;file3.txt`

### MapFlag 示例

展示如何使用映射标志并自定义分隔符：

```go
package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建映射标志
	configF := qflag.Map("config", "c", map[string]string{}, "配置键值对")
	// 设置键值对分隔符为逗号，键值分隔符为冒号
	configF.SetDelimiters(",", ":")

	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	// 获取映射值
	config := configF.Get()
	fmt.Printf("配置: %v\n", config)
}
```

使用方式：`./app --config server:localhost,port:8080,timeout:30s`

### PathFlag 示例

展示路径标志的路径验证功能：

```go
package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建路径标志并设置验证规则
	logPathF := qflag.Path("log-path", "l", "/var/log/app", "日志目录")
		.IsDirectory(true)  // 必须是目录
		.MustExist(false)   // 路径不存在时自动创建

	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("日志目录: %s\n", logPathF.Get())
}
```

### EnumFlag 大小写敏感示例

展示枚举标志的大小写敏感设置：

```go
package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建枚举标志并启用大小写敏感
	modeF := qflag.Enum("mode", "m", "debug", []string{"Debug", "Release", "Test"}, "运行模式")
	modeF.SetCaseSensitive(true)

	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("运行模式: %s\n", modeF.Get())
}
```

注意：启用大小写敏感后，`--mode debug`会报错，必须使用`--mode Debug`

### 基本使用示例

以下是一个简单的使用示例，展示如何创建一个命令并添加一些标志：

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
	if runF.Get() {
		fmt.Printf("启动app: %s\n", nameF.Get())
		fmt.Printf("app路径: %s\n", pathF.Get())
	}
}

```

### 子命令示例

以下是一个使用子命令的示例：

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
	if runF.Get() {
		fmt.Printf("启动app: %s\n", nameF.GetValue())
		fmt.Printf("app路径: %s\n", pathF.GetValue())
	}

	// 判断是否执行了stop子命令
	if stopF.Get() {
		fmt.Printf("停止app: %s\n", stopCmdName.Get())
	}
}
```

## API文档

qflag提供了完善的API文档，按模块组织如下：

- **全局命令处理**: [qflag包文档](./APIDOC.md) - 包含全局命令创建、参数解析和子命令管理相关API
- **核心命令处理**: [cmd包文档](./cmd/APIDOC.md) - 包含命令创建、参数解析和子命令管理相关API
- **标志类型定义**: [flags包文档](./flags/APIDOC.md) - 包含所有标志类型的详细定义和使用方法
- **错误处理**: [qerr包文档](./qerr/APIDOC.md) - 包含错误类型和处理相关API
- **参数验证**: [validator包文档](./validator/APIDOC.md) - 包含参数验证器接口和实现

完整的API文档可通过访问各模块对应的APIDOC.md文件查看。

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。
