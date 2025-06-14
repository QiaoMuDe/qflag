// Package qflag 提供对标准库flag的封装, 自动实现长短标志互斥, 并默认绑定-h/--help标志打印帮助信息。
// 用户可通过Cmd.Help字段自定义帮助内容, 支持直接赋值字符串或从文件加载。
package qflag

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

// NewCmd 创建新的命令实例
// 参数:
// name: 命令名称
// shortName: 命令短名称
// errorHandling: 错误处理方式
// 返回值: *Cmd命令实例指针
// errorHandling可选值: flag.ContinueOnError、flag.ExitOnError、flag.PanicOnError
func NewCmd(name string, shortName string, errorHandling flag.ErrorHandling) *Cmd {
	// 检查命令名称是否为空
	if name == "" {
		panic("cmd name cannot be empty")
	}

	// 设置默认的错误处理方式
	if errorHandling == 0 {
		errorHandling = flag.ContinueOnError
	}

	cmd := &Cmd{
		fs:            flag.NewFlagSet(name, errorHandling), // 创建新的flag集
		shortToLong:   sync.Map{},                           // 存储长短标志的映射关系
		longToShort:   sync.Map{},                           // 存储长短标志的映射关系
		helpFlagName:  "help",                               // 默认的帮助标志名称
		helpShortName: "h",                                  // 默认的帮助标志短名称
		flagRegistry:  make(map[interface{}]Flag),           // 初始化指针注册表
		usage:         "",                                   // 允许用户直接设置帮助内容
		description:   "",                                   // 允许用户直接设置命令描述
		name:          name,                                 // 命令名称, 用于帮助信息中显示
		shortName:     shortName,                            // 命令短名称, 用于帮助信息中显示
	}

	cmd.bindHelpFlag() // 自动绑定帮助标志
	return cmd
}

// GetFlagByPtr 通过指针获取对应的Flag对象
// 参数: p - 绑定的指针
// 返回值: 对应的Flag对象和错误信息
// 如果指针未注册, 则返回错误
func (c *Cmd) GetFlagByPtr(p interface{}) (interface{}, error) {
	flag, exists := c.flagRegistry[p]
	if !exists {
		return nil, fmt.Errorf("指针未注册: %v", p)
	}
	return flag, nil
}

// AddSubCmd 关联一个或多个子命令到当前命令
// 参数:
// subCmds: 子命令的切片
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) {
	// 将子命令关联到当前命令并设置父命令指针
	for _, cmd := range subCmds {
		cmd.parentCmd = c
	}

	// 将子命令添加到当前命令的子命令列表中
	c.subCmds = append(c.subCmds, subCmds...)
}

// Parse 解析命令行参数, 自动检查长短标志互斥, 并处理帮助标志
func (c *Cmd) Parse(args []string) error {
	// 1. 调用flag库解析参数
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	// 2. 检查是否使用-h/--help标志
	if c.isHelpRequested() {
		c.printHelp()
		os.Exit(0)
	}

	// 3. 检查长短标志互斥（跳过帮助标志）
	var conflictMsg string
	c.longToShort.Range(func(longKey, shortVal interface{}) bool {
		longFlag := longKey.(string)   // 获取长标志名称
		shortFlag := shortVal.(string) // 获取短标志名称

		// 跳过帮助标志的互斥检查
		if longFlag == c.helpFlagName || shortFlag == c.helpShortName {
			return true
		}

		// 检查标志是否同时被设置
		longChanged := c.isFlagSet(longFlag)   // 检查长标志是否被设置
		shortChanged := c.isFlagSet(shortFlag) // 检查短标志是否被设置

		// 如果两个标志都发生变化, 则表示冲突
		if longChanged && shortChanged {
			conflictMsg = fmt.Sprintf("不能同时使用 --%s 和 -%s", longFlag, shortFlag)
			return false // 终止遍历
		}

		return true // 继续遍历
	})

	// 返回冲突错误（如果有）
	if conflictMsg != "" {
		return fmt.Errorf("%s", conflictMsg)
	}
	return nil
}

// String 添加字符串类型标志, 参数含义与Int方法一致
func (c *Cmd) String(name, shortName, defValue, usage string) *StringFlag {
	value := c.fs.String(name, defValue, usage) // 获取标志对象(绑定长标志)
	f := &StringFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     value,     // 标志对象
	}
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)    // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)    // 存储长到短的映射关系
		c.fs.String(shortName, defValue, usage) // 绑定短标志
	}
	return f
}

// StringVar 绑定字符串类型标志到指针并内部注册Flag对象
// 参数依次为: 指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) StringVar(p *string, name, shortName, defValue, usage string) {
	c.fs.StringVar(p, name, defValue, usage) // 绑定长标志
	f := &StringFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     p,         // 标志对象
	}
	c.flagRegistry[p] = f // 注册Flag对象

	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)          // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)          // 存储长到短的映射关系
		c.fs.StringVar(p, shortName, defValue, usage) // 绑定短标志
	}
}

// IntVar 绑定整数类型标志到指针并内部注册Flag对象
// 参数依次为: 指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) IntVar(p *int, name, shortName string, defValue int, usage string) {
	c.fs.IntVar(p, name, defValue, usage) // 绑定长标志
	f := &IntFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     p,         // 标志对象
	}
	c.flagRegistry[p] = f // 注册Flag对象

	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)       // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)       // 存储长到短的映射关系
		c.fs.IntVar(p, shortName, defValue, usage) // 绑定短标志
	}
}

// BoolVar 绑定布尔类型标志到指针并内部注册Flag对象
// 参数依次为: 指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) BoolVar(p *bool, name, shortName string, defValue bool, usage string) {
	c.fs.BoolVar(p, name, defValue, usage) // 绑定长标志
	f := &BoolFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     p,         // 标志对象
	}
	c.flagRegistry[p] = f // 注册Flag对象

	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)        // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)        // 存储长到短的映射关系
		c.fs.BoolVar(p, shortName, defValue, usage) // 绑定短标志
	}
}

// Int 添加整数类型标志, 返回标志对象。参数依次为长标志名、短标志、默认值、帮助说明
func (c *Cmd) Int(name, shortName string, defValue int, usage string) *IntFlag {
	value := c.fs.Int(name, defValue, usage) // 获取标志对象(绑定长标志)
	f := &IntFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     value,     // 标志对象
	}
	// 非帮助标志才记录映射（避免覆盖帮助标志）
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name) // 存储短到长的映射关系
		c.longToShort.Store(name, shortName) // 存储长到短的映射关系
		c.fs.Int(shortName, defValue, usage) // 绑定短标志
	}
	return f
}

// Bool 添加布尔类型标志, 参数含义与Int方法一致
func (c *Cmd) Bool(name, shortName string, defValue bool, usage string) *BoolFlag {
	value := c.fs.Bool(name, defValue, usage) // 获取标志对象(绑定长标志)
	f := &BoolFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     value,     // 标志对象
	}
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)  // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)  // 存储长到短的映射关系
		c.fs.Bool(shortName, defValue, usage) // 绑定短标志
	}
	return f
}

// Float 添加浮点型标志, 返回标志对象。参数依次为长标志名、短标志、默认值、帮助说明
func (c *Cmd) Float(name, shortName string, defValue float64, usage string) *FloatFlag {
	value := c.fs.Float64(name, defValue, usage) // 获取标志对象(绑定长标志)
	f := &FloatFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     value,     // 标志对象
	}
	// 非帮助标志才记录映射（避免覆盖帮助标志）
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)     // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)     // 存储长到短的映射关系
		c.fs.Float64(shortName, defValue, usage) // 绑定短标志
	}
	return f
}

// FloatVar 绑定浮点型标志到指针并内部注册Flag对象
// 参数依次为: 指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) FloatVar(p *float64, name, shortName string, defValue float64, usage string) {
	c.fs.Float64Var(p, name, defValue, usage) // 绑定长标志
	f := &FloatFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     p,         // 标志对象
	}
	c.flagRegistry[p] = f // 注册Flag对象

	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)           // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)           // 存储长到短的映射关系
		c.fs.Float64Var(p, shortName, defValue, usage) // 绑定短标志
	}
}
