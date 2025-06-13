package qflag

import (
	"flag"
	"sync"
)

// Cmd 命令行标志管理结构体，封装参数解析、长短标志互斥及帮助系统。
type Cmd struct {
	/* 内部使用属性*/
	fs            *flag.FlagSet // 底层flag集合, 处理参数解析
	shortToLong   sync.Map      // 短标志到长标志的映射（键：短标志，值：长标志）
	longToShort   sync.Map      // 长标志到短标志的映射（键：长标志，值：短标志）
	helpFlagName  string        // 帮助标志的长名称，默认"help"
	helpShortName string        // 帮助标志的短名称，默认"h"
	helpFlagBound bool          // 标记帮助标志是否已绑定
	subCmds       []*Cmd        // 子命令列表, 用于关联子命令
	parentCmd     *Cmd          // 父命令指针，用于递归调用, 根命令的父命令为nil

	/* 外部可访问属性 */
	Help        string // 自定义帮助内容，可由用户直接赋值
	Description string // 自定义描述，用于帮助信息中显示
	Name        string // 命令名称，用于帮助信息中显示
	ShortName   string // 命令短名称，用于帮助信息中显示
}

// 帮助信息模板常量
const (
	cmdNameTemplate            = "命令: %s\n\n"                  // 命令名称
	cmdNameWithShortTemplate   = "命令: %s(%s)\n\n"              // 命令名称和短名称
	cmdDescriptionTemplate     = "描述: %s\n\n"                  // 命令描述
	cmdUsageTemplate           = "用法: %s [选项] [参数]\n\n"        // 命令用法
	cmdUsageWithSubCmdTemplate = "用法: %s [子命令] [选项] [参数]\n\n"  // 命令用法(带子命令)
	cmdUsageSubCmdTemplate     = "用法: %s %s [选项] [参数]\n\n"     // 命令用法(带子命令)
	optionsHeaderTemplate      = "选项:\n"                       // 选项头部
	optionTemplate             = "  -%s, --%s\t%s (默认值: %s)\n" // 选项模板
	subCmdsHeaderTemplate      = "\n子命令:\n"                    // 子命令头部
	subCmdTemplate             = "  %s\t%s\n"                  // 子命令模板
)

// IntFlag 整数类型标志结构体，包含标志元数据和值访问接口
type IntFlag struct {
	cmd       *Cmd   // 所属的命令实例
	name      string // 长标志名称（如"port"）
	shortName string // 短标志字符（如"p"，空表示无短标志）
	defValue  int    // 默认值
	help      string // 帮助说明
	value     *int   // 标志值指针，通过flag库绑定
}

// StringFlag 字符串类型标志结构体
type StringFlag struct {
	cmd       *Cmd    // 所属的命令实例
	name      string  // 长标志名称
	shortName string  // 短标志字符
	defValue  string  // 默认值
	help      string  // 帮助说明
	value     *string // 标志值指针
}

// BoolFlag 布尔类型标志结构体
type BoolFlag struct {
	cmd       *Cmd   // 所属的命令实例
	name      string // 长标志名称
	shortName string // 短标志字符
	defValue  bool   // 默认值
	help      string // 帮助说明
	value     *bool  // 标志值指针
}
