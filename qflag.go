// Package qflag 提供对标准库flag的封装，自动实现长短标志互斥，并默认绑定-h/--help标志打印帮助信息。
// 用户可通过Command.HelpContent字段自定义帮助内容，支持直接赋值字符串或从文件加载。
package qflag

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

// Command 命令行标志管理结构体，封装参数解析、长短标志互斥及帮助系统。
type Command struct {
	fs            *flag.FlagSet // 底层flag集合，处理参数解析
	shortToLong   sync.Map      // 短标志到长标志的映射（键：短标志，值：长标志）
	longToShort   sync.Map      // 长标志到短标志的映射（键：长标志，值：短标志）
	HelpContent   string        // 自定义帮助内容，可由用户直接赋值
	helpFlagName  string        // 帮助标志的长名称，默认"help"
	helpShortName string        // 帮助标志的短名称，默认"h"
	helpFlagBound bool          // 标记帮助标志是否已绑定
}

// NewCommand 创建新的命令实例，参数name为命令名称，errorHandling指定错误处理方式。
// errorHandling可选值：flag.ContinueOnError、flag.ExitOnError、flag.PanicOnError
func NewCommand(name string, errorHandling flag.ErrorHandling) *Command {
	cmd := &Command{
		fs:            flag.NewFlagSet(name, errorHandling),
		shortToLong:   sync.Map{},
		longToShort:   sync.Map{},
		helpFlagName:  "help",
		helpShortName: "h",
		HelpContent:   "", // 允许用户直接设置帮助内容
	}
	cmd.bindHelpFlag() // 自动绑定帮助标志
	return cmd
}

// bindHelpFlag 绑定-h/--help标志到显示帮助信息的逻辑
func (c *Command) bindHelpFlag() {
	if c.helpFlagBound {
		return
	}
	var showHelp bool
	// 绑定长帮助标志
	c.fs.BoolVar(&showHelp, c.helpFlagName, false, "显示帮助信息")
	// 绑定短帮助标志（若设置）
	if c.helpShortName != "" {
		c.fs.BoolVar(&showHelp, c.helpShortName, false, "显示帮助信息")
		c.shortToLong.Store(c.helpShortName, c.helpFlagName)
		c.longToShort.Store(c.helpFlagName, c.helpShortName)
	}
	c.helpFlagBound = true
}

// IntFlag 整数类型标志结构体，包含标志元数据和值访问接口
type IntFlag struct {
	cmd       *Command // 所属的命令实例
	name      string   // 长标志名称（如"port"）
	shortName string   // 短标志字符（如"p"，空表示无短标志）
	defValue  int      // 默认值
	help      string   // 帮助说明
	value     *int     // 标志值指针，通过flag库绑定
}

// Int 添加整数类型标志，返回标志对象。参数依次为长标志名、短标志、默认值、帮助说明
func (c *Command) Int(name, shortName string, defValue int, help string) *IntFlag {
	value := c.fs.Int(name, defValue, help)
	f := &IntFlag{
		cmd:       c,
		name:      name,
		shortName: shortName,
		defValue:  defValue,
		help:      help,
		value:     value,
	}
	// 非帮助标志才记录映射（避免覆盖帮助标志）
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)
		c.longToShort.Store(name, shortName)
		c.fs.Int(shortName, defValue, help)
	}
	return f
}

// StringFlag 字符串类型标志结构体
type StringFlag struct {
	cmd       *Command // 所属的命令实例
	name      string   // 长标志名称
	shortName string   // 短标志字符
	defValue  string   // 默认值
	help      string   // 帮助说明
	value     *string  // 标志值指针
}

// String 添加字符串类型标志，参数含义与Int方法一致
func (c *Command) String(name, shortName, defValue, help string) *StringFlag {
	value := c.fs.String(name, defValue, help)
	f := &StringFlag{
		cmd:       c,
		name:      name,
		shortName: shortName,
		defValue:  defValue,
		help:      help,
		value:     value,
	}
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)
		c.longToShort.Store(name, shortName)
		c.fs.String(shortName, defValue, help)
	}
	return f
}

// BoolFlag 布尔类型标志结构体
type BoolFlag struct {
	cmd       *Command // 所属的命令实例
	name      string   // 长标志名称
	shortName string   // 短标志字符
	defValue  bool     // 默认值
	help      string   // 帮助说明
	value     *bool    // 标志值指针
}

// Bool 添加布尔类型标志，参数含义与Int方法一致
func (c *Command) Bool(name, shortName string, defValue bool, help string) *BoolFlag {
	value := c.fs.Bool(name, defValue, help)
	f := &BoolFlag{
		cmd:       c,
		name:      name,
		shortName: shortName,
		defValue:  defValue,
		help:      help,
		value:     value,
	}
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)
		c.longToShort.Store(name, shortName)
		c.fs.Bool(shortName, defValue, help)
	}
	return f
}

// Parse 解析命令行参数，自动检查长短标志互斥，并处理帮助标志
func (c *Command) Parse(args []string) error {
	// 1. 调用flag库解析参数
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	// 2. 检查是否请求帮助
	if c.isHelpRequested() {
		c.printHelp()
		os.Exit(0)
	}

	// 3. 检查长短标志互斥（跳过帮助标志）
	var conflictMsg string
	c.longToShort.Range(func(longKey, shortVal interface{}) bool {
		longFlag := longKey.(string)
		shortFlag := shortVal.(string)
		// 跳过帮助标志的互斥检查
		if longFlag == c.helpFlagName || shortFlag == c.helpShortName {
			return true
		}

		// 检查标志是否同时被设置
		longChanged := c.isFlagSet(longFlag)
		shortChanged := c.isFlagSet(shortFlag)
		if longChanged && shortChanged {
			conflictMsg = fmt.Sprintf("不能同时使用 --%s 和 -%s", longFlag, shortFlag)
			return false // 终止遍历
		}
		return true // 继续遍历
	})
	// 返回冲突错误（如果有）
	return errorIf(conflictMsg != "", conflictMsg)
}

// isHelpRequested 检测帮助标志是否被用户设置
func (c *Command) isHelpRequested() bool {
	// 检查长帮助标志
	if c.isFlagSet(c.helpFlagName) {
		return true
	}
	// 检查短帮助标志
	if c.helpShortName != "" {
		return c.isFlagSet(c.helpShortName)
	}
	return false
}

// isFlagSet 检查标志是否被用户显式设置
func (c *Command) isFlagSet(name string) bool {
	// 获取标志对象
	flag := c.fs.Lookup(name)
	if flag == nil {
		return false
	}

	// 特殊处理布尔标志
	if b, ok := flag.Value.(interface{ IsBoolFlag() bool }); ok && b.IsBoolFlag() {
		// 布尔标志被设置的条件：
		// 1. 值为true且默认值为false（用户显式启用）
		// 2. 值为false且默认值为true（用户显式禁用）
		currentVal := flag.Value.String()
		return (currentVal == "true" && flag.DefValue == "false") ||
			(currentVal == "false" && flag.DefValue == "true")
	}

	// 处理其他类型标志（int/string等）
	// 只要当前值与默认值不同，即认为被设置
	return flag.Value.String() != flag.DefValue
}

// printHelp 打印帮助内容，优先显示用户自定义的HelpContent
func (c *Command) printHelp() {
	if c.HelpContent != "" {
		fmt.Println(c.HelpContent)
	} else {
		fmt.Println("未设置帮助内容，请通过cmd.HelpContent赋值")
	}
	fmt.Println()
	c.fs.Usage() // 打印flag原生帮助信息
}

// errorIf 辅助函数，将非空字符串转为error，空字符串返回nil
func errorIf(cond bool, msg string) error {
	if !cond {
		return nil
	}
	// 使用 %s 格式化字符串，避免非常量格式字符串的问题
	return fmt.Errorf("%s", msg)
}
