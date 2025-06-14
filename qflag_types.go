package qflag

import (
	"flag"
	"sync"
)

// Cmd 命令行标志管理结构体,封装参数解析、长短标志互斥及帮助系统。
type Cmd struct {
	/* 内部使用属性*/
	fs                           *flag.FlagSet        // 底层flag集合, 处理参数解析
	shortToLong                  sync.Map             // 短标志到长标志的映射（键：短标志,值：长标志）
	longToShort                  sync.Map             // 长标志到短标志的映射（键：长标志,值：短标志）
	helpFlagName                 string               // 帮助标志的长名称,默认"help"
	helpShortName                string               // 帮助标志的短名称,默认"h"
	helpFlag                     *bool                // 帮助标志指针,用于绑定和检查
	helpFlagBound                bool                 // 标记帮助标志是否已绑定
	helpOnce                     sync.Once            // 用于确保帮助标志只被绑定一次
	showInstallPathFlagName      string               // 安装路径标志的长名称,默认"show-install-path"
	showInstallPathFlagShortName string               // 安装路径标志的短名称,默认"sip"
	showInstallPathFlag          *bool                // 安装路径标志指针,用于绑定和检查
	subCmds                      []*Cmd               // 子命令列表, 用于关联子命令
	parentCmd                    *Cmd                 // 父命令指针,用于递归调用, 根命令的父命令为nil
	flagRegistry                 map[interface{}]Flag // 标志注册表,用于通过指针查找标志元数据
	usage                        string               // 自定义帮助内容,可由用户直接赋值
	description                  string               // 自定义描述,用于帮助信息中显示
	name                         string               // 命令名称,用于帮助信息中显示
	shortName                    string               // 命令短名称,用于帮助信息中显示
	args                         []string             // 命令行参数切片
	addMu                        sync.Mutex           // 互斥锁,确保并发安全操作
	parseOnce                    sync.Once            // 用于确保命令只被解析一次
	setMu                        sync.Mutex           // 互斥锁,确保并发安全操作
}

// 标志类型
type FlagType int

const (
	FlagTypeBool   FlagType = iota // 布尔类型
	FlagTypeInt                    // 整数类型
	FlagTypeString                 // 字符串类型
	FlagTypeFloat                  // 浮点数类型
)

// Command 命令接口定义,封装命令的核心功能
// 提供属性访问和子命令管理的标准方法
// 实现类应确保线程安全的标志操作
//
// 示例:
// cmd := NewCmd()
// cmd.SetName("app")
// cmd.SetDescription("示例应用")
type Command interface {
	// 属性访问方法
	Name() string               // 返回命令名称
	ShortName() string          // 返回命令短名称
	Description() string        // 获取命令描述
	SetDescription(desc string) // 设置命令描述
	Usage() string              // 获取命令用法
	SetUsage(usage string)      // 设置命令用法

	// 子命令管理
	AddSubCmd(subCmd *Cmd) // 添加子命令
	SubCmds() []*Cmd       // 获取子命令列表

	// 标志操作
	Parse(args []string) error                                            // 解析命令行参数
	GetFlagByPtr(ptr interface{}) (Flag, bool)                            // 通过标志指针获取标志元数据
	String(name, shortName, usage string, defValue string) *StringFlag    // 添加字符串类型标志
	Int(name, shortName, usage string, defValue int) *IntFlag             // 添加整数类型标志
	Bool(name, shortName, usage string, defValue bool) *BoolFlag          // 添加布尔类型标志
	Float(name, shortName, usage string, defValue float64) *FloatFlag     // 添加浮点数类型标志
	StringVar(name, shortName, usage string, defValue string) *StringFlag // 添加字符串类型标志
	IntVar(name, shortName, usage string, defValue int) *IntFlag          // 添加整数类型标志
	BoolVar(name, shortName, usage string, defValue bool) *BoolFlag       // 添加布尔类型标志
	FloatVar(name, shortName, usage string, defValue float64) *FloatFlag  // 添加浮点数类型标志
	Args() []string                                                       // 获取命令行参数切片
	Arg(i int) string                                                     // 获取命令行参数
}

// Name 命令名称
func (c *Cmd) Name() string { return c.name }

// ShortName 命令短名称
func (c *Cmd) ShortName() string { return c.shortName }

// Description 命令描述
func (c *Cmd) Description() string { return c.description }

// SetDescription 设置命令描述
func (c *Cmd) SetDescription(desc string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.description = desc
}

// Usage 命令用法
func (c *Cmd) Usage() string { return c.usage }

// SetUsage 设置命令用法
func (c *Cmd) SetUsage(usage string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.usage = usage
}

// SubCmds 子命令列表
func (c *Cmd) SubCmds() []*Cmd { return c.subCmds }

// Args 命令行参数切片
func (c *Cmd) Args() []string { return c.args }

// Arg 获取命令行参数
func (c *Cmd) Arg(i int) string {
	if i >= 0 && i < len(c.args) {
		return c.args[i]
	}

	return ""
}

// PrintUsage 打印命令的帮助信息, 优先打印用户的帮助信息, 否则自动生成帮助信息
func (c *Cmd) PrintUsage() {
	c.printUsage()
}

// 帮助信息模板常量
const (
	cmdNameTemplate            = "Command: %s\n\n"                                  // 命令名称
	cmdNameWithShortTemplate   = "Command: %s(%s)\n\n"                              // 命令名称和短名称
	cmdDescriptionTemplate     = "Description: %s\n\n"                              // 命令描述
	cmdUsageTemplate           = "Usage: %s [options] [arguments]\n\n"              // 命令用法
	cmdUsageWithSubCmdTemplate = "Usage: %s [subcommand] [options] [arguments]\n\n" // 命令用法(带子命令)
	cmdUsageSubCmdTemplate     = "Usage: %s %s [options] [arguments]\n\n"           // 命令用法(带子命令)
	optionsHeaderTemplate      = "Options:\n"                                       // 选项头部
	optionTemplate1            = "  -%s, --%s\t%s (default: %s)\n"                  // 选项模板1
	optionTemplate2            = "  --%s\t%s (默认值: %s)\n"                           // 选项模板2
	subCmdsHeaderTemplate      = "\nSubcommands:\n"                                 // 子命令头部
	subCmdTemplate             = "  %s\t%s\n"                                       // 子命令模板
	notesHeaderTemplate        = "\nNotes:\n"                                       // 注意事项头部
	noteItemTemplate           = "  %d、%s\n"                                        // 注意事项项
)

// Flag 所有标志类型的通用接口,定义了标志的元数据访问方法
type Flag interface {
	Name() string       // 获取标志的名称
	ShortName() string  // 获取标志的短名称
	Usage() string      // 获取标志的用法
	Type() FlagType     // 获取标志类型
	getDefaultAny() any // 新增通用方法(内部使用)
}

// TypedFlag 所有标志类型的通用接口,定义了标志的元数据访问方法和默认值访问方法
type TypedFlag[T any] interface {
	Flag           // 继承标志接口
	GetDefault() T // 获取标志的默认值
	GetValue() T   // 获取标志的实际值
	SetValue(T)    // 设置标志的值
}

// IntFlag 整数类型标志结构体,包含标志元数据和值访问接口
type IntFlag struct {
	cmd       *Cmd         // 所属的命令实例
	name      string       // 长标志名称（如"port"）
	shortName string       // 短标志字符（如"p",空表示无短标志）
	defValue  int          // 默认值
	usage     string       // 帮助说明
	value     *int         // 标志值指针,通过flag库绑定
	getRu     sync.RWMutex // 用于确保标志值访问的线程安全
	setMu     sync.Mutex   // 用于确保标志值设置操作的线程安全
}

// 实现Flag接口
func (f *IntFlag) Name() string       { return f.name }
func (f *IntFlag) ShortName() string  { return f.shortName }
func (f *IntFlag) Usage() string      { return f.usage }
func (f *IntFlag) GetDefault() int    { return f.defValue }
func (f *IntFlag) Type() FlagType     { return FlagTypeInt }
func (f *IntFlag) getDefaultAny() any { return f.defValue }

// GetValue 获取标志的实际值（带线程安全保护）
// 返回值优先级：解析值 > 默认值
func (f *IntFlag) GetValue() int {
	f.getRu.RLock()         // 加读锁
	defer f.getRu.RUnlock() // 确保锁释放

	if f.value != nil { // 优先返回解析值
		return *f.value
	}
	return f.defValue // 其次返回默认值
}

// SetValue 设置标志的值（带线程安全保护）
func (f *IntFlag) SetValue(value int) {
	f.setMu.Lock()
	defer f.setMu.Unlock()

	f.value = &value
}

// StringFlag 字符串类型标志结构体
type StringFlag struct {
	cmd       *Cmd         // 所属的命令实例
	name      string       // 长标志名称
	shortName string       // 短标志字符
	defValue  string       // 默认值
	usage     string       // 帮助说明
	value     *string      // 标志值指针
	getRu     sync.RWMutex // 用于确保标志值访问的线程安全
	setMu     sync.Mutex   // 用于确保标志值设置操作的线程安全
}

// 实现Flag接口
func (f *StringFlag) Name() string       { return f.name }
func (f *StringFlag) ShortName() string  { return f.shortName }
func (f *StringFlag) Usage() string      { return f.usage }
func (f *StringFlag) GetDefault() string { return f.defValue }
func (f *StringFlag) Type() FlagType     { return FlagTypeString }
func (f *StringFlag) getDefaultAny() any { return f.defValue }

// GetValue 获取标志的实际值（带线程安全保护）
// 返回值优先级：解析值 > 默认值
func (f *StringFlag) GetValue() string {
	f.getRu.RLock()         // 加读锁
	defer f.getRu.RUnlock() // 确保锁释放

	if f.value != nil { // 优先返回解析值
		return *f.value
	}
	return f.defValue // 其次返回默认值
}

// SetValue 设置标志的值（带线程安全保护）
func (f *StringFlag) SetValue(value string) {
	f.setMu.Lock()
	defer f.setMu.Unlock()

	f.value = &value
}

// BoolFlag 布尔类型标志结构体
type BoolFlag struct {
	cmd       *Cmd         // 所属的命令实例
	name      string       // 长标志名称
	shortName string       // 短标志字符
	defValue  bool         // 默认值
	usage     string       // 帮助说明
	value     *bool        // 标志值指针
	getRu     sync.RWMutex // 用于确保标志值访问的线程安全
	setMu     sync.Mutex   // 用于确保标志值设置操作的线程安全
}

// 实现Flag接口
func (f *BoolFlag) Name() string       { return f.name }
func (f *BoolFlag) ShortName() string  { return f.shortName }
func (f *BoolFlag) Usage() string      { return f.usage }
func (f *BoolFlag) GetDefault() bool   { return f.defValue }
func (f *BoolFlag) Type() FlagType     { return FlagTypeBool }
func (f *BoolFlag) getDefaultAny() any { return f.defValue }

// GetValue 获取标志的实际值（带线程安全保护）
// 返回值优先级：解析值 > 默认值
func (f *BoolFlag) GetValue() bool {
	f.getRu.RLock()         // 加读锁
	defer f.getRu.RUnlock() // 确保锁释放

	if f.value != nil { // 优先返回解析值
		return *f.value
	}
	return f.defValue // 其次返回默认值
}

// SetValue 设置标志的值（带线程安全保护）
func (f *BoolFlag) SetValue(value bool) {
	f.setMu.Lock()
	defer f.setMu.Unlock()

	f.value = &value
}

// FloatFlag 浮点型标志结构体
type FloatFlag struct {
	cmd       *Cmd         // 所属的命令实例
	name      string       // 长标志名称
	shortName string       // 短标志字符
	defValue  float64      // 默认值
	usage     string       // 帮助说明
	value     *float64     // 标志值指针
	getRu     sync.RWMutex // 用于确保标志值访问的线程安全
	setMu     sync.Mutex   // 用于确保标志值设置操作的线程安全
}

// 实现Flag接口
func (f *FloatFlag) Name() string        { return f.name }
func (f *FloatFlag) ShortName() string   { return f.shortName }
func (f *FloatFlag) Usage() string       { return f.usage }
func (f *FloatFlag) GetDefault() float64 { return f.defValue }
func (f *FloatFlag) Type() FlagType      { return FlagTypeFloat }
func (f *FloatFlag) getDefaultAny() any  { return f.defValue }

// GetValue 获取标志的实际值（带线程安全保护）
// 返回值优先级：解析值 > 默认值
func (f *FloatFlag) GetValue() float64 {
	f.getRu.RLock()         // 加读锁
	defer f.getRu.RUnlock() // 确保锁释放

	if f.value != nil { // 优先返回解析值
		return *f.value
	}
	return f.defValue // 其次返回默认值
}

// SetValue 设置标志的值（带线程安全保护）
func (f *FloatFlag) SetValue(value float64) {
	f.setMu.Lock()
	defer f.setMu.Unlock()

	f.value = &value
}
