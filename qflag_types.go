package qflag

import (
	"flag"
	"sync"
)

// Cmd 命令行标志管理结构体，封装参数解析、长短标志互斥及帮助系统。
type Cmd struct {
	fs            *flag.FlagSet // 底层flag集合，处理参数解析
	shortToLong   sync.Map      // 短标志到长标志的映射（键：短标志，值：长标志）
	longToShort   sync.Map      // 长标志到短标志的映射（键：长标志，值：短标志）
	Help          string        // 自定义帮助内容，可由用户直接赋值
	Description   string        // 自定义描述，用于帮助信息中显示
	helpFlagName  string        // 帮助标志的长名称，默认"help"
	helpShortName string        // 帮助标志的短名称，默认"h"
	helpFlagBound bool          // 标记帮助标志是否已绑定
	SubCmds       []*Cmd        // 子命令列表, 用于关联子命令
}

// 帮助信息模板常量
const (
	helpHeaderTemplate      = "命令: %s\n\n"
	helpDescriptionTemplate = "描述: %s\n\n"
	helpUsageTemplate       = "用法: %s [选项] [参数]\n\n"
	helpOptionsHeader       = "选项:\n"
	helpOptionTemplate      = "  -%s, --%s\t%s (默认值: %s)\n"
	helpSubCommandsHeader   = "\n子命令:\n"
	helpSubCommandTemplate  = "  %s\t%s\n"
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
