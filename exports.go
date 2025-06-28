// Package qflag 根包统一导出入口
// 本文件用于将各子包的核心功能导出到根包，简化外部使用
package qflag

import (
	"gitee.com/MM-Q/qflag/cmd"
)

/*
项目地址: https://gitee.com/MM-Q/qflag
*/

// 导出子包类型和函数 //

// QCommandLine 导出cmd包的全局默认Command实例
var QCommandLine = cmd.QCommandLine

// NewCmd 导出cmd包中的NewCommand函数
var NewCmd = cmd.NewCommand
