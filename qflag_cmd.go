package qflag

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Cmd 命令行标志管理结构体,封装参数解析、长短标志互斥及帮助系统。
type Cmd struct {
	/* 内部使用属性*/
	fs                           *flag.FlagSet // 底层flag集合, 处理参数解析
	flagRegistry                 *FlagRegistry // 标志注册表, 统一管理标志的元数据
	helpFlagName                 string        // 帮助标志的长名称,默认"help"
	helpFlagShortName            string        // 帮助标志的短名称,默认"h"
	helpFlag                     *BoolFlag     // 帮助标志指针,用于绑定和检查
	helpFlagBound                bool          // 标记帮助标志是否已绑定
	helpOnce                     sync.Once     // 用于确保帮助标志只被绑定一次
	showInstallPathFlagName      string        // 安装路径标志的长名称,默认"show-install-path"
	showInstallPathFlagShortName string        // 安装路径标志的短名称,默认"sip"
	showInstallPathFlag          *BoolFlag     // 安装路径标志指针,用于绑定和检查
	subCmds                      []*Cmd        // 子命令列表, 用于关联子命令
	parentCmd                    *Cmd          // 父命令指针,用于递归调用, 根命令的父命令为nil
	usage                        string        // 自定义帮助内容,可由用户直接赋值
	description                  string        // 自定义描述,用于帮助信息中显示
	name                         string        // 命令名称,用于帮助信息中显示
	shortName                    string        // 命令短名称,用于帮助信息中显示
	args                         []string      // 命令行参数切片
	addMu                        sync.Mutex    // 互斥锁,确保并发安全操作
	parseOnce                    sync.Once     // 用于确保命令只被解析一次
	setMu                        sync.Mutex    // 互斥锁,确保并发安全操作
	builtinFlagNameMap           sync.Map      // 用于存储内置标志名称的映射
}

// 内置标志名称
var (
	helpFlagName                 = "help"
	helpFlagShortName            = "h"
	showInstallPathFlagName      = "show-install-path"
	showInstallPathFlagShortName = "sip"
)

// Command 命令接口定义，封装命令行程序的核心功能
// 提供统一的命令管理、参数解析和帮助系统
// 实现类需保证线程安全，所有方法应支持并发调用
//
// 示例用法:
// cmd := NewCmd("app", "a", flag.ContinueOnError)
// cmd.SetDescription("示例应用程序")
// cmd.String("config", "c", "配置文件路径", "/etc/app.conf")
type Command interface {
	Name() string               // 获取命令名称(长名称)，如"app"
	ShortName() string          // 获取命令短名称，如"a"
	Description() string        // 获取命令描述信息
	SetDescription(desc string) // 设置命令描述信息，用于帮助输出
	Usage() string              // 获取自定义用法说明，为空时自动生成
	SetUsage(usage string)      // 设置自定义用法说明，覆盖自动生成内容

	AddSubCmd(subCmd *Cmd) // 添加子命令，子命令会继承父命令的上下文
	SubCmds() []*Cmd       // 获取所有已注册的子命令列表

	Parse(args []string) error // 解析命令行参数，自动处理标志和子命令

	String(name, shortName, usage, defValue string) *StringFlag               // 添加字符串类型标志
	Int(name, shortName, usage string, defValue int) *IntFlag                 // 添加整数类型标志
	Bool(name, shortName, usage string, defValue bool) *BoolFlag              // 添加布尔类型标志
	Float(name, shortName, usage string, defValue float64) *FloatFlag         // 添加浮点数类型标志
	Slice(name, shortName string, defValue []string, usage string) *SliceFlag // 添加字符串切片类型标志

	StringVar(f *StringFlag, name, shortName, defValue, usage string)               // 绑定字符串标志到指定变量
	IntVar(f *IntFlag, name, shortName string, defValue int, usage string)          // 绑定整数标志到指定变量
	BoolVar(f *BoolFlag, name, shortName string, defValue bool, usage string)       // 绑定布尔标志到指定变量
	FloatVar(f *FloatFlag, name, shortName string, defValue float64, usage string)  // 绑定浮点数标志到指定变量
	SliceVar(f *SliceFlag, name, shortName string, defValue []string, usage string) // 绑定切片标志到指定变量

	Args() []string   // 获取所有非标志参数(未绑定到任何标志的参数)
	Arg(i int) string // 获取指定索引的非标志参数，索引越界返回空字符串
	NArg() int        // 获取非标志参数的数量
	NFlag() int       // 获取已解析的标志数量

	FlagExists(name string) bool // 检查指定名称的标志是否存在(支持长/短名称)
	PrintUsage()                 // 打印命令使用说明到标准输出

	bindHelpFlagAndShowInstallPathFlag() // 绑定帮助标志和显示安装路径标志
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

// Args 获取非标志参数切片
func (c *Cmd) Args() []string { return c.args }

// Arg 获取指定索引的非标志参数
func (c *Cmd) Arg(i int) string {
	if i >= 0 && i < len(c.args) {
		return c.args[i]
	}
	return ""
}

// NArg 获取非标志参数的数量
func (c *Cmd) NArg() int { return len(c.args) }

// NFlag 获取标志的数量
func (c *Cmd) NFlag() int { return c.fs.NFlag() }

// PrintUsage 打印命令的帮助信息, 优先打印用户的帮助信息, 否则自动生成帮助信息
func (c *Cmd) PrintUsage() {
	c.printUsage()
}

// FlagExists 检查指定名称的标志是否存在
func (c *Cmd) FlagExists(name string) bool {
	if _, exists := c.flagRegistry.GetByName(name); exists {
		return true
	}

	return false
}

// bindHelpFlagAndShowInstallPathFlag 绑定帮助标志和显示安装路径标志
func (c *Cmd) bindHelpFlagAndShowInstallPathFlag() {
	// 检查是否已绑定
	if c.helpFlagBound {
		return // 避免重复绑定
	}

	// 初始化帮助标志
	c.helpOnce.Do(func() {
		if c.helpFlag == nil {
			// 为空时自动初始化
			c.helpFlag = &BoolFlag{}
		}

		// 绑定帮助标志
		c.BoolVar(c.helpFlag, c.helpFlagName, c.helpFlagShortName, false, "Show help information")

		// 绑定显示安装路径标志
		if c.showInstallPathFlag == nil {
			c.showInstallPathFlag = &BoolFlag{}
		}

		// 绑定显示安装路径标志
		c.BoolVar(c.showInstallPathFlag, c.showInstallPathFlagName, c.showInstallPathFlagShortName, false, "Show install path")

		// 添加内置标志到检测映射
		c.builtinFlagNameMap.Store(helpFlagName, true)
		c.builtinFlagNameMap.Store(helpFlagShortName, true)
		c.builtinFlagNameMap.Store(showInstallPathFlagName, true)
		c.builtinFlagNameMap.Store(showInstallPathFlagShortName, true)

		// 设置帮助标志已绑定
		c.helpFlagBound = true
	})
}

// generateHelpInfo 生成命令帮助信息
// cmd: 当前命令
// 返回值: 命令帮助信息
func generateHelpInfo(cmd *Cmd) string {
	var helpInfo string

	// 命令名（支持短名称显示）
	if cmd.shortName != "" {
		helpInfo += fmt.Sprintf(cmdNameWithShortTemplate, cmd.fs.Name(), cmd.shortName)
	} else {
		helpInfo += fmt.Sprintf(cmdNameTemplate, cmd.fs.Name())
	}

	// 命令描述
	if cmd.description != "" {
		helpInfo += fmt.Sprintf(cmdDescriptionTemplate, cmd.description)
	}

	// 动态生成命令用法
	fullCmdPath := getFullCommandPath(cmd)
	usageLine := "Usage: " + fullCmdPath

	// 如果存在子命令，则需要添加子命令用法
	if len(cmd.subCmds) > 0 {
		usageLine += " [subcommand]"
	}

	// 添加选项用法
	usageLine += " [options] [arguments]\n\n"
	helpInfo += usageLine

	// 选项标题
	helpInfo += optionsHeaderTemplate

	// 收集所有标志信息
	var flags []struct {
		longFlag  string
		shortFlag string
		usage     string
		defValue  string
	}

	// 使用Flag接口统一访问标志属性
	for _, f := range cmd.flagRegistry.allFlags {
		flag := f
		flags = append(flags, struct {
			longFlag  string
			shortFlag string
			usage     string
			defValue  string
		}{
			longFlag:  flag.GetLongName(),
			shortFlag: flag.GetShortName(),
			usage:     flag.GetUsage(),
			defValue:  fmt.Sprintf("%v", flag.GetDefault()),
		})
	}

	// 按短标志字母顺序排序，有短标志的选项优先
	sort.Slice(flags, func(i, j int) bool {
		a, b := flags[i], flags[j]
		aHasShort := a.shortFlag != ""
		bHasShort := b.shortFlag != ""

		// 有短标志的选项排在前面
		if aHasShort && !bHasShort {
			return true
		}
		if !aHasShort && bHasShort {
			return false
		}

		// 都有短标志则按短标志排序，都没有则按长标志排序
		if aHasShort && bHasShort {
			return a.shortFlag < b.shortFlag
		}
		return a.longFlag < b.longFlag
	})

	// 生成排序后的标志信息
	for _, flag := range flags {
		if flag.shortFlag != "" {
			helpInfo += fmt.Sprintf(optionTemplate1, flag.shortFlag, flag.longFlag, flag.usage, flag.defValue)
		} else {
			helpInfo += fmt.Sprintf(optionTemplate2, flag.longFlag, flag.usage, flag.defValue)
		}
	}

	// 如果有子命令，添加子命令信息
	if len(cmd.subCmds) > 0 {
		helpInfo += subCmdsHeaderTemplate
		for _, subCmd := range cmd.subCmds {
			helpInfo += fmt.Sprintf(subCmdTemplate, subCmd.fs.Name(), subCmd.description)
		}
	}

	// 添加注意事项
	helpInfo += notesHeaderTemplate
	helpInfo += fmt.Sprintf(noteItemTemplate, 1, "In the case where both long options and short options are used at the same time, the option specified last shall take precedence.")

	return helpInfo
}

// printUsage 打印帮助内容, 优先显示用户自定义的Usage
func (c *Cmd) printUsage() {
	if c.usage != "" {
		fmt.Println(c.usage)
	} else {
		// 自动生成帮助信息
		fmt.Println(generateHelpInfo(c))
	}
}

// hasCycle 检测命令间是否存在循环引用
// 采用深度优先搜索(DFS)算法，通过访问标记避免重复检测
// 参数:
//
//	parent: 当前命令
//	child: 待添加的子命令
//
// 返回值:
//
//	如果存在循环引用则返回true
func hasCycle(parent, child *Cmd) bool {
	if parent == nil || child == nil {
		return false
	}

	visited := make(map[*Cmd]bool)
	return dfs(parent, child, visited)
}

// dfs 深度优先搜索检测循环引用
func dfs(target, current *Cmd, visited map[*Cmd]bool) bool {
	// 如果已访问过当前节点，直接返回避免无限循环
	if visited[current] {
		return false
	}
	visited[current] = true

	// 找到目标节点，存在循环引用
	if current == target {
		return true
	}

	// 递归检查所有子命令
	for _, subCmd := range current.subCmds {
		if dfs(target, subCmd, visited) {
			return true
		}
	}

	// 检查父命令链
	if current.parentCmd != nil {
		return dfs(target, current.parentCmd, visited)
	}

	return false
}

// joinErrors 将错误切片合并为单个错误
func joinErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}
	if len(errors) == 1 {
		return errors[0]
	}

	// 构建错误信息
	var b strings.Builder
	b.WriteString(fmt.Sprintf("A total of %d errors:\n", len(errors)))
	for i, err := range errors {
		b.WriteString(fmt.Sprintf("  %d. %v\n", i+1, err))
	}

	// 使用常量格式字符串，将错误信息作为参数传入
	return fmt.Errorf("Merged error message:\n%s", b.String())
}

// getFullCommandPath 递归构建完整的命令路径，从根命令到当前命令
func getFullCommandPath(cmd *Cmd) string {
	if cmd.parentCmd == nil {
		return cmd.fs.Name()
	}
	return getFullCommandPath(cmd.parentCmd) + " " + cmd.fs.Name()
}

// validateFlag 通用标志验证逻辑
func (c *Cmd) validateFlag(name, shortName string) {
	// 新增格式校验
	if strings.ContainsAny(name, invalidFlagChars) {
		panic(fmt.Sprintf("The flag name '%s' contains illegal characters", name))
	}

	// 检查标志名称和短名称是否为空
	if name == "" {
		panic("Flag name cannot be empty")
	}
	if shortName == "" {
		panic("Flag short name cannot be empty")
	}

	// 检查标志是否已存在
	if _, exists := c.flagRegistry.GetByName(name); exists {
		panic(fmt.Sprintf("Flag name %s already exists", name))
	}

	if _, exists := c.flagRegistry.GetByName(shortName); exists {
		panic(fmt.Sprintf("Flag short name %s already exists", shortName))
	}

	// 检查标志是否为内置标志
	if _, ok := c.builtinFlagNameMap.Load(name); ok {
		panic(fmt.Sprintf("Flag name %s is reserved", name))
	}
	if _, ok := c.builtinFlagNameMap.Load(shortName); ok {
		panic(fmt.Sprintf("Flag short name %s is reserved", shortName))
	}
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

	// 检查命令短名称是否为空
	if shortName == "" {
		panic("cmd short name cannot be empty")
	}

	// 设置默认的错误处理方式为ContinueOnError, 避免测试时意外退出
	if errorHandling == 0 {
		errorHandling = flag.ContinueOnError
	}

	// 创建标志注册表
	flags := &FlagRegistry{
		mu:       sync.RWMutex{},             // 并发安全锁
		byLong:   make(map[string]*FlagMeta), // 存储长标志的映射
		byShort:  make(map[string]*FlagMeta), // 存储短标志的映射
		allFlags: []*FlagMeta{},              // 存储所有标志的切片
	}

	// 创建新的Cmd实例
	cmd := &Cmd{
		fs:                           flag.NewFlagSet(name, errorHandling), // 创建新的flag集
		helpFlagName:                 helpFlagName,                         // 默认的帮助标志名称
		helpFlagShortName:            helpFlagShortName,                    // 默认的帮助标志短名称
		showInstallPathFlagName:      showInstallPathFlagName,              // 默认的显示安装路径标志名称
		showInstallPathFlagShortName: showInstallPathFlagShortName,         // 默认的显示安装路径标志短名称
		usage:                        "",                                   // 允许用户直接设置帮助内容
		description:                  "",                                   // 允许用户直接设置命令描述
		name:                         name,                                 // 命令名称, 用于帮助信息中显示
		shortName:                    shortName,                            // 命令短名称, 用于帮助信息中显示
		args:                         []string{},                           // 命令行参数
		flagRegistry:                 flags,                                // 初始化标志注册表
		helpFlag:                     &BoolFlag{},                          // 初始化帮助标志
		showInstallPathFlag:          &BoolFlag{},                          // 初始化显示安装路径标志
	}

	// 自动绑定帮助标志和显示安装路径标志
	cmd.bindHelpFlagAndShowInstallPathFlag()
	return cmd
}

// AddSubCmd 关联一个或多个子命令到当前命令
// 支持批量添加多个子命令，遇到错误时收集所有错误并返回
// 参数:
//
//	subCmds: 子命令的切片
//
// 返回值:
//
//	错误信息列表，如果所有子命令添加成功则返回nil
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error {
	c.addMu.Lock()
	defer c.addMu.Unlock()

	// 检查子命令是否为空
	if len(subCmds) == 0 {
		return fmt.Errorf("subcommand list cannot be empty")
	}

	var errors []error
	addedCmds := make([]*Cmd, 0, len(subCmds))

	// 第一阶段：验证所有子命令
	for _, cmd := range subCmds {
		if cmd == nil {
			errors = append(errors, fmt.Errorf("Subcommand cannot be nil"))
			continue
		}

		// 检测循环引用
		if hasCycle(c, cmd) {
			errors = append(errors, fmt.Errorf("Cyclic reference detected: Command %s already exists in the command chain", cmd.name))
			continue
		}

		addedCmds = append(addedCmds, cmd)
	}

	// 如果有验证错误，返回所有错误信息
	if len(errors) > 0 {
		return fmt.Errorf("Failed to add subcommands: %w", joinErrors(errors))
	}

	// 第二阶段：批量添加子命令
	for _, cmd := range addedCmds {
		cmd.parentCmd = c                  // 设置父命令指针
		c.subCmds = append(c.subCmds, cmd) // 添加到子命令列表
	}

	return nil
}

// Parse 解析命令行参数, 自动检查长短标志, 并处理帮助标志
// 工作流程:
//  1. 调用flag库解析参数
//  2. 处理内置标志(-h/--help和-sip/--show-install-path)
//  3. 检测并处理子命令:当第一个非标志参数匹配子命令名称时,
//     将剩余参数传递给子命令解析
//
// 注意: 该方法保证每个Cmd实例只会解析一次
func (c *Cmd) Parse(args []string) error {
	var err error

	// 确保只解析一次
	c.parseOnce.Do(func() {
		// 1调用flag库解析参数
		if parseErr := c.fs.Parse(args); parseErr != nil {
			err = fmt.Errorf("Parameter parsing error: %w", parseErr)
			return
		}

		// 检查是否使用-h/--help标志
		if c.helpFlag.GetValue() {
			if c.fs.ErrorHandling() != flag.ContinueOnError {
				c.printUsage() // 只有在ExitOnError或PanicOnError时才打印使用说明
				os.Exit(0)
			}
			return
		}

		// 检查是否使用-sip/--show-install-path标志
		if c.showInstallPathFlag.GetValue() {
			if c.fs.ErrorHandling() != flag.ContinueOnError {
				// 只有在ExitOnError或PanicOnError时才打印安装路径
				fmt.Println(GetExecutablePath())
				os.Exit(0)
			}
			return
		}

		// 设置非标志参数
		c.args = append(c.args, c.fs.Args()...)

		// 检查是否有子命令
		if len(c.args) > 0 {
			for _, subCmd := range c.subCmds {
				if c.args[0] == subCmd.name || c.args[0] == subCmd.shortName {
					// 将剩余参数传递给子命令解析
					if err = subCmd.Parse(c.args[1:]); err != nil {
						err = fmt.Errorf("Subcommand parsing error: %w", err)
					}
					return
				}
			}
		}
	})

	// 检查是否报错
	if err != nil {
		return err
	}
	return nil
}

// String 添加字符串类型标志, 返回标志对象,参数依次为长标志名、短标志、默认值、帮助说明
func (c *Cmd) String(name, shortName, defValue, usage string) *StringFlag {
	// 参数校验（复用公共函数）
	c.validateFlag(name, shortName)

	// 显式初始化当前值的默认值
	currentStr := defValue

	// 创建StringFlag对象
	f := &StringFlag{
		cmd:       c,           // 命令对象
		name:      name,        // 长标志名
		shortName: shortName,   // 短标志名
		defValue:  defValue,    // 默认值
		usage:     usage,       // 帮助说明
		value:     &currentStr, // 当前值
	}

	// 创建FlagMeta对象
	meta := &FlagMeta{
		longName:  name,           // 长标志名
		shortName: shortName,      // 短标志名
		flagType:  FlagTypeString, // 标志类型
		usage:     usage,          // 帮助说明
		defValue:  defValue,       // 默认值
		isBuiltin: false,          // 是否为内置标志
	}

	// 绑定短标志到同一个变量
	c.fs.StringVar(&currentStr, shortName, defValue, usage)

	// 绑定长标志到变量
	c.fs.StringVar(&currentStr, name, defValue, usage)

	// 注册Flag元数据
	c.flagRegistry.RegisterFlag(meta)

	return f
}

// StringVar 绑定字符串类型标志到指针并内部注册Flag对象
// 参数依次为: 字符串标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) StringVar(f *StringFlag, name, shortName, defValue, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("StringFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	c.validateFlag(name, shortName)

	// 显式初始化当前值的默认值
	currentStr := defValue

	// 创建StringFlag对象
	f = &StringFlag{
		cmd:       c,           // 命令对象
		name:      name,        // 长标志名
		shortName: shortName,   // 短标志名
		defValue:  defValue,    // 默认值
		usage:     usage,       // 帮助说明
		value:     &currentStr, // 标志对象
	}

	// 创建FlagMeta对象
	meta := &FlagMeta{
		longName:  name,           // 长标志名
		shortName: shortName,      // 短标志名
		flagType:  FlagTypeString, // 标志类型
		usage:     usage,          // 帮助说明
		defValue:  defValue,       // 默认值
		isBuiltin: false,          // 是否为内置标志
	}

	// 绑定短标志
	c.fs.StringVar(&currentStr, shortName, defValue, usage)

	// 绑定长标志
	c.fs.StringVar(&currentStr, name, defValue, usage)

	// 注册Flag对象
	c.flagRegistry.RegisterFlag(meta)
}

// IntVar 绑定整数类型标志到指针并内部注册Flag对象
// 参数依次为: 整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) IntVar(f *IntFlag, name, shortName string, defValue int, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("IntFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	c.validateFlag(name, shortName)

	// 初始化默认值
	currentInt := defValue

	// 创建IntFlag对象
	f = &IntFlag{
		cmd:       c,           // 命令对象
		name:      name,        // 长标志名
		shortName: shortName,   // 短标志名
		defValue:  defValue,    // 默认值
		usage:     usage,       // 帮助说明
		value:     &currentInt, // 标志对象
	}

	// 创建FlagMeta对象
	meta := &FlagMeta{
		longName:  name,        // 长标志名
		shortName: shortName,   // 短标志名
		flagType:  FlagTypeInt, // 标志类型
		usage:     usage,       // 帮助说明
		defValue:  defValue,    // 默认值
		isBuiltin: false,       // 是否为内置标志
	}

	// 绑定短标志
	c.fs.IntVar(&currentInt, shortName, defValue, usage)

	// 绑定长标志
	c.fs.IntVar(&currentInt, name, defValue, usage)

	// 注册Flag对象
	c.flagRegistry.RegisterFlag(meta)
}

// Int 添加整数类型标志, 返回标志对象。参数依次为长标志名、短标志、默认值、帮助说明
func (c *Cmd) Int(name, shortName string, defValue int, usage string) *IntFlag {
	// 参数校验（复用公共函数）
	c.validateFlag(name, shortName)

	// 显式初始化默认值，提高可读性
	currentInt := defValue

	// 创建IntFlag对象
	f := &IntFlag{
		cmd:       c,           // 命令对象
		name:      name,        // 长标志名
		shortName: shortName,   // 短标志名
		defValue:  defValue,    // 默认值
		usage:     usage,       // 帮助说明
		value:     &currentInt, // 当前值
	}

	// 创建FlagMeta对象
	meta := &FlagMeta{
		longName:  name,        // 长标志名
		shortName: shortName,   // 短标志名
		flagType:  FlagTypeInt, // 标志类型
		usage:     usage,       // 帮助说明
		defValue:  defValue,    // 默认值
		isBuiltin: false,       // 是否为内置标志
	}

	// 绑定短标志到同一个变量
	c.fs.IntVar(&currentInt, shortName, defValue, usage)

	// 绑定长标志到变量
	c.fs.IntVar(&currentInt, name, defValue, usage)

	// 注册Flag元数据
	c.flagRegistry.RegisterFlag(meta)

	return f
}

// BoolVar 绑定布尔类型标志到指针并内部注册Flag对象
// 参数依次为: 布尔标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) BoolVar(f *BoolFlag, name, shortName string, defValue bool, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("BoolFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	c.validateFlag(name, shortName)

	// 显式初始化默认值，提高可读性
	currentBool := defValue

	// 创建BoolFlag对象
	f = &BoolFlag{
		cmd:       c,            // 命令对象
		name:      name,         // 长标志名
		shortName: shortName,    // 短标志名
		defValue:  defValue,     // 默认值
		usage:     usage,        // 帮助说明
		value:     &currentBool, // 当前值
	}

	// 创建FlagMeta对象
	meta := &FlagMeta{
		longName:  name,         // 长标志名
		shortName: shortName,    // 短标志名
		flagType:  FlagTypeBool, // 标志类型
		usage:     usage,        // 帮助说明
		defValue:  defValue,     // 默认值
		isBuiltin: false,        // 是否为内置标志
	}

	// 绑定短标志
	c.fs.BoolVar(&currentBool, shortName, defValue, usage)

	// 绑定长标志
	c.fs.BoolVar(&currentBool, name, defValue, usage)

	// 注册Flag对象
	c.flagRegistry.RegisterFlag(meta)
}

// Bool 添加布尔类型标志, 返回标志对象。参数依次为长标志名、短标志、默认值、帮助说明
func (c *Cmd) Bool(name, shortName string, defValue bool, usage string) *BoolFlag {
	// 参数校验（复用公共函数）
	c.validateFlag(name, shortName)

	// 显式初始化默认值，提高可读性
	currentBool := defValue

	// 创建标志对象
	f := &BoolFlag{
		cmd:       c,            // 命令对象
		name:      name,         // 长标志名
		shortName: shortName,    // 短标志名
		defValue:  defValue,     // 默认值
		usage:     usage,        // 帮助说明
		value:     &currentBool, // 当前值
	}

	// 创建FlagMeta对象
	meta := &FlagMeta{
		longName:  name,         // 长标志名
		shortName: shortName,    // 短标志名
		flagType:  FlagTypeBool, // 标志类型
		usage:     usage,        // 帮助说明
		defValue:  defValue,     // 默认值
		isBuiltin: false,        // 是否为内置标志
	}

	// 绑定长标志到变量
	c.fs.BoolVar(&currentBool, name, defValue, usage)

	// 绑定短标志到同一个变量
	c.fs.BoolVar(&currentBool, shortName, defValue, usage)

	// 注册Flag元数据
	c.flagRegistry.RegisterFlag(meta)

	return f
}

// Float 添加浮点型标志, 返回标志对象。参数依次为长标志名、短标志、默认值、帮助说明
func (c *Cmd) Float(name, shortName string, defValue float64, usage string) *FloatFlag {
	// 参数校验（复用公共函数）
	c.validateFlag(name, shortName)

	// 显式初始化默认值
	currentFloat := new(float64) // 显式堆分配
	*currentFloat = defValue

	// 创建标志对象
	f := &FloatFlag{
		cmd:       c,            // 命令对象
		name:      name,         // 长标志名
		shortName: shortName,    // 短标志名
		defValue:  defValue,     // 默认值
		usage:     usage,        // 帮助说明
		value:     currentFloat, // 标志对象
	}

	// 创建FlagMeta对象
	meta := &FlagMeta{
		longName:  name,          // 长标志名
		shortName: shortName,     // 短标志名
		flagType:  FlagTypeFloat, // 标志类型
		usage:     usage,         // 帮助说明
		defValue:  defValue,      // 默认值
		isBuiltin: false,         // 是否为内置标志
	}

	// 绑定短标志到同一个变量
	c.fs.Float64Var(currentFloat, shortName, defValue, usage)

	// 绑定长标志到变量
	c.fs.Float64Var(currentFloat, name, defValue, usage)

	// 注册Flag元数据
	c.flagRegistry.RegisterFlag(meta)

	return f
}

// FloatVar 绑定浮点型标志到指针并内部注册Flag对象
// 参数依次为: 浮点数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) FloatVar(f *FloatFlag, name, shortName string, defValue float64, usage string) {
	// 检查指针是否为空
	if f == nil {
		panic("FloatFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	c.validateFlag(name, shortName)

	// 显式初始化默认值
	currentFloat := new(float64) // 显式堆分配
	*currentFloat = defValue

	// 创建标志对象
	f = &FloatFlag{
		cmd:       c,            // 命令对象
		name:      name,         // 长标志名
		shortName: shortName,    // 短标志名
		defValue:  defValue,     // 默认值
		usage:     usage,        // 帮助说明
		value:     currentFloat, // 标志对象
	}

	// 创建FlagMeta对象
	meta := &FlagMeta{
		longName:  name,          // 长标志名
		shortName: shortName,     // 短标志名
		flagType:  FlagTypeFloat, // 标志类型
		usage:     usage,         // 帮助说明
		defValue:  defValue,      // 默认值
		isBuiltin: false,         // 是否为内置标志
	}

	// 绑定短标志
	c.fs.Float64Var(currentFloat, shortName, defValue, usage)

	// 绑定长标志
	c.fs.Float64Var(currentFloat, name, defValue, usage)

	// 注册Flag对象
	c.flagRegistry.RegisterFlag(meta)
}

// Slice 为命令添加字符串切片类型标志, 返回标志对象。参数依次为: 长标志名、短标志、默认值切片、帮助说明
func (c *Cmd) Slice(name, shortName string, defValue []string, usage string) *SliceFlag {
	// 参数校验（复用公共函数）
	c.validateFlag(name, shortName)

	// 初始化默认值
	if defValue == nil {
		defValue = make([]string, 0)
	}

	// 初始化默认值（修复当前实现的空切片问题）
	currentSlice := make([]string, len(defValue))
	copy(currentSlice, defValue) // 创建一个副本

	// 创建标志对象
	f := &SliceFlag{
		cmd:       c,             // 命令对象
		name:      name,          // 长标志名
		shortName: shortName,     // 短标志名
		defValue:  defValue,      // 默认值
		usage:     usage,         // 帮助说明
		value:     &currentSlice, // 当前值
	}

	// 创建FlagMeta对象
	meta := &FlagMeta{
		longName:  name,          // 长标志名
		shortName: shortName,     // 短标志名
		flagType:  FlagTypeSlice, // 标志类型
		usage:     usage,         // 帮助说明
		defValue:  defValue,      // 默认值
		isBuiltin: false,         // 是否为内置标志
	}

	// 绑定短标志到同一个变量
	c.fs.Var(f, shortName, usage)

	// 绑定长标志到变量
	c.fs.Var(f, name, usage)

	// 注册Flag元数据
	c.flagRegistry.RegisterFlag(meta)

	return f
}

// SliceVar 为命令添加字符串切片类型标志, 返回标志对象。参数依次为: 指针切片、长标志名、短标志、默认值切片、帮助说明
func (c *Cmd) SliceVar(f *SliceFlag, name, shortName string, defValue []string, usage string) {
	// 检查指针是否为空
	if f == nil {
		panic("SliceFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	c.validateFlag(name, shortName)

	// 初始化默认值（修复当前实现的空切片问题）
	if defValue == nil {
		defValue = make([]string, 0)
	}

	// 在创建SliceFlag前添加默认值初始化
	currentSlice := make([]string, len(defValue))
	copy(currentSlice, defValue) // 创建一个副本

	// 创建SliceFlag
	f = &SliceFlag{
		cmd:       c,             // 命令对象
		name:      name,          // 长标志名
		shortName: shortName,     // 短标志名
		defValue:  defValue,      // 默认值
		usage:     usage,         // 帮助说明
		value:     &currentSlice, // 标志对象
	}

	// 绑定长标志
	c.fs.Var(f, name, usage)

	// 绑定短标志
	c.fs.Var(f, shortName, usage)

	// 创建FlagMeta对象
	meta := &FlagMeta{
		longName:  name,          // 长标志名
		shortName: shortName,     // 短标志名
		flagType:  FlagTypeSlice, // 标志类型
		usage:     usage,         // 帮助说明
		defValue:  defValue,      // 默认值
		isBuiltin: false,         // 是否为内置标志
	}

	// 注册Flag元数据
	c.flagRegistry.RegisterFlag(meta)
}

// GetExecutablePath 获取程序的绝对安装路径
// 如果无法通过 os.Executable 获取路径,则使用 os.Args[0] 作为替代
// 返回：程序的绝对路径字符串
func GetExecutablePath() string {
	// 尝试使用 os.Executable 获取可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		// 如果 os.Executable 报错,使用 os.Args[0] 作为替代
		exePath = os.Args[0]
	}
	// 使用 filepath.Abs 确保路径是绝对路径
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		// 如果 filepath.Abs 报错,直接返回原始路径
		return exePath
	}
	return absPath
}
