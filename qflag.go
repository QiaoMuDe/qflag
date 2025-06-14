// Package qflag 提供对标准库flag的封装, 自动实现长短标志互斥, 并默认绑定-h/--help标志打印帮助信息。
// 用户可通过Cmd.Help字段自定义帮助内容, 支持直接赋值字符串或从文件加载。
package qflag

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// QCommandLine 全局默认Cmd实例
var QCommandLine *Cmd

// parseOnce 保证Parse函数只会执行一次
var parseOnce sync.Once

// 在包初始化时创建全局默认Cmd实例
func init() {
	// 处理可能的空os.Args情况
	if len(os.Args) == 0 {
		// 如果os.Args为空,则创建一个新的Cmd对象,命令行参数为空,错误处理方式为ExitOnError
		QCommandLine = NewCmd("", "", flag.ExitOnError)
	} else {
		// 如果os.Args不为空,则创建一个新的Cmd对象,命令行参数为filepath.Base(os.Args[0]),错误处理方式为ExitOnError
		QCommandLine = NewCmd(filepath.Base(os.Args[0]), "", flag.ExitOnError)
	}
}

// String 创建字符串类型标志（全局默认命令）
// 参数依次为：长标志名、短标志、默认值、帮助说明
func String(name, shortName, defValue, usage string) *StringFlag {
	return QCommandLine.String(name, shortName, defValue, usage)
}

// Int 创建整数类型标志（全局默认命令）
// 参数依次为：长标志名、短标志、默认值、帮助说明
func Int(name, shortName string, defValue int, usage string) *IntFlag {
	return QCommandLine.Int(name, shortName, defValue, usage)
}

// Bool 创建布尔类型标志（全局默认命令）
// 参数依次为：长标志名、短标志、默认值、帮助说明
func Bool(name, shortName string, defValue bool, usage string) *BoolFlag {
	return QCommandLine.Bool(name, shortName, defValue, usage)
}

// Float 创建浮点数类型标志（全局默认命令）
// 参数依次为：长标志名、短标志、默认值、帮助说明
func Float(name, shortName string, defValue float64, usage string) *FloatFlag {
	return QCommandLine.Float(name, shortName, defValue, usage)
}

// StringVar 绑定字符串类型标志到指针（全局默认命令）
// 参数依次为：指针、长标志名、短标志、默认值、帮助说明
func StringVar(p *string, name, shortName, defValue, usage string) {
	QCommandLine.StringVar(p, name, shortName, defValue, usage)
}

// IntVar 绑定整数类型标志到指针（全局默认命令）
// 参数依次为：指针、长标志名、短标志、默认值、帮助说明
func IntVar(p *int, name, shortName string, defValue int, usage string) {
	QCommandLine.IntVar(p, name, shortName, defValue, usage)
}

// BoolVar 绑定布尔类型标志到指针（全局默认命令）
// 参数依次为：指针、长标志名、短标志、默认值、帮助说明
func BoolVar(p *bool, name, shortName string, defValue bool, usage string) {
	QCommandLine.BoolVar(p, name, shortName, defValue, usage)
}

// FloatVar 绑定浮点数类型标志到指针（全局默认命令）
// 参数依次为：指针、长标志名、短标志、默认值、帮助说明
func FloatVar(p *float64, name, shortName string, defValue float64, usage string) {
	QCommandLine.FloatVar(p, name, shortName, defValue, usage)
}

// Parse 解析命令行参数（全局默认命令）
// 该函数保证只会执行一次解析操作
func Parse() error {
	var err error
	parseOnce.Do(func() {
		// 解析命令行参数
		err = QCommandLine.Parse(os.Args[1:])
	})
	return err
}

// AddSubCmd 给全局默认命令添加子命令
// 返回值: 错误信息, 如果检测到循环引用或nil子命令
func AddSubCmd(subCmds ...*Cmd) error {
	return QCommandLine.AddSubCmd(subCmds...)
}

// GetFlagByPtr 通过指针获取标志（全局默认命令）
func GetFlagByPtr(p any) (any, error) {
	return QCommandLine.GetFlagByPtr(p)
}

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

	// 设置默认的错误处理方式为ContinueOnError, 避免测试时意外退出
	if errorHandling == 0 {
		errorHandling = flag.ContinueOnError
	}

	cmd := &Cmd{
		fs:                           flag.NewFlagSet(name, errorHandling), // 创建新的flag集
		shortToLong:                  sync.Map{},                           // 存储长短标志的映射关系
		longToShort:                  sync.Map{},                           // 存储长短标志的映射关系
		helpFlagName:                 "help",                               // 默认的帮助标志名称
		helpShortName:                "h",                                  // 默认的帮助标志短名称
		showInstallPathFlagName:      "show-install-path",                  // 默认的显示安装路径标志名称
		showInstallPathFlagShortName: "sip",                                // 默认的显示安装路径标志短名称
		flagRegistry:                 make(map[interface{}]Flag),           // 初始化指针注册表
		usage:                        "",                                   // 允许用户直接设置帮助内容
		description:                  "",                                   // 允许用户直接设置命令描述
		name:                         name,                                 // 命令名称, 用于帮助信息中显示
		shortName:                    shortName,                            // 命令短名称, 用于帮助信息中显示
		args:                         []string{},                           // 命令行参数
		helpFlag:                     new(bool),                            // 默认的帮助标志
		showInstallPathFlag:          new(bool),                            // 默认的显示安装路径标志
	}

	// 自动绑定帮助标志和显示安装路径标志
	cmd.bindHelpFlagAndShowInstallPathFlag()
	return cmd
}

// GetFlagByPtr 通过指针获取对应的Flag对象
// 参数: p - 绑定的指针
// 返回值: 对应的Flag对象和错误信息
// 如果指针未注册, 则返回错误
func (c *Cmd) GetFlagByPtr(p any) (any, error) {
	flag, exists := c.flagRegistry[p]
	if !exists {
		return nil, fmt.Errorf("指针未注册: %v", p)
	}
	return flag, nil
}

// AddSubCmd 关联一个或多个子命令到当前命令
// 参数:
// subCmds: 子命令的切片
// 返回值: 错误信息, 如果检测到循环引用或nil子命令
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error {
	c.addMu.Lock()
	defer c.addMu.Unlock()

	// 检查子命令是否为空或nil
	if len(subCmds) == 0 {
		return fmt.Errorf("没有可添加的子命令")
	}

	// 合并处理：设置父命令指针并添加到子命令列表
	for _, cmd := range subCmds {
		if cmd == nil {
			return fmt.Errorf("子命令不能为nil")
		}

		// 检测循环引用
		if hasCycle(c, cmd) {
			return fmt.Errorf("检测到循环引用: 命令 %s 已存在于父命令链中", cmd.name)
		}

		cmd.parentCmd = c                  // 设置父命令指针
		c.subCmds = append(c.subCmds, cmd) // 添加到子命令列表
	}
	return nil
}

// Parse 解析命令行参数, 自动检查长短标志互斥, 并处理帮助标志
func (c *Cmd) Parse(args []string) error {
	var err error

	// 确保只解析一次
	c.parseOnce.Do(func() {
		// 1. 调用flag库解析参数
		if parseErr := c.fs.Parse(args); parseErr != nil {
			err = fmt.Errorf("Parameter parsing error: %v", parseErr)
			return
		}

		// 2. 检查是否使用-h/--help标志
		if *c.helpFlag {
			if c.fs.ErrorHandling() != flag.ContinueOnError {
				c.printUsage() // 只有在ExitOnError或PanicOnError时才打印使用说明
				os.Exit(0)
			}
			return
		}

		// 3. 设置命令行参数
		c.args = append(c.args, c.fs.Args()...)
	})

	// 检查是否报错
	if err != nil {
		return err
	}
	return nil
}

// String 添加字符串类型标志, 参数含义与Int方法一致
func (c *Cmd) String(name, shortName, defValue, usage string) *StringFlag {
	var value string
	c.fs.StringVar(&value, name, defValue, usage) // 绑定长标志到变量
	f := &StringFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     &value,    // 标志对象
	}
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)               // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)               // 存储长到短的映射关系
		c.fs.StringVar(&value, shortName, defValue, usage) // 绑定短标志到同一个变量
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
	var value int
	c.fs.IntVar(&value, name, defValue, usage) // 绑定长标志到变量
	f := &IntFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     &value,    // 标志对象
	}
	// 非帮助标志才记录映射（避免覆盖帮助标志）
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)            // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)            // 存储长到短的映射关系
		c.fs.IntVar(&value, shortName, defValue, usage) // 绑定短标志到同一个变量
	}
	return f
}

// Bool 添加布尔类型标志, 参数含义与Int方法一致
func (c *Cmd) Bool(name, shortName string, defValue bool, usage string) *BoolFlag {
	var value bool
	c.fs.BoolVar(&value, name, defValue, usage) // 绑定长标志到变量
	f := &BoolFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     &value,    // 标志对象
	}
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)             // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)             // 存储长到短的映射关系
		c.fs.BoolVar(&value, shortName, defValue, usage) // 绑定短标志到同一个变量
	}
	return f
}

// Float 添加浮点型标志, 返回标志对象。参数依次为长标志名、短标志、默认值、帮助说明
func (c *Cmd) Float(name, shortName string, defValue float64, usage string) *FloatFlag {
	var value float64
	c.fs.Float64Var(&value, name, defValue, usage) // 绑定长标志到变量
	f := &FloatFlag{
		cmd:       c,         // 命令对象
		name:      name,      // 长标志名
		shortName: shortName, // 短标志名
		defValue:  defValue,  // 默认值
		usage:     usage,     // 帮助说明
		value:     &value,    // 标志对象
	}
	// 非帮助标志才记录映射（避免覆盖帮助标志）
	if shortName != "" && shortName != c.helpShortName {
		c.shortToLong.Store(shortName, name)                // 存储短到长的映射关系
		c.longToShort.Store(name, shortName)                // 存储长到短的映射关系
		c.fs.Float64Var(&value, shortName, defValue, usage) // 绑定短标志到同一个变量
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
