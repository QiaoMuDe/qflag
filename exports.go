// Package qflag 根包统一导出入口
// 本文件用于将各子包的核心功能导出到根包，简化外部使用
package qflag

import "gitee.com/MM-Q/qflag/cmd"

// 导出子包类型和函数

// Command 命令行解析器核心类型
 type Command = cmd.Cmd

// NewCmd 创建新的命令实例
var NewCmd = cmd.NewCmd

// ExampleInfo 示例信息结构体
type ExampleInfo = cmd.ExampleInfo

// 导出 flags 子包的核心类型
