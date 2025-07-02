package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SetVersion 设置版本信息
func (c *Cmd) SetVersion(version string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.version = version
}

// GetVersion 获取版本信息
func (c *Cmd) GetVersion() string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.version
}

// SetModuleHelps 设置自定义模块帮助信息
func (c *Cmd) SetModuleHelps(moduleHelps string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.moduleHelps = moduleHelps
}

// GetModuleHelps 获取自定义模块帮助信息
func (c *Cmd) GetModuleHelps() string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.moduleHelps
}

// SetLogoText 设置logo文本
func (c *Cmd) SetLogoText(logoText string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.logoText = logoText
}

// GetLogoText 获取logo文本
func (c *Cmd) GetLogoText() string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.logoText
}

// GetUseChinese 获取是否使用中文帮助信息
func (c *Cmd) GetUseChinese() bool {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.useChinese
}

// SetUseChinese 设置是否使用中文帮助信息
func (c *Cmd) SetUseChinese(useChinese bool) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.useChinese = useChinese
}

// GetNotes 获取所有备注信息
func (c *Cmd) GetNotes() []string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	// 返回切片副本而非原始引用
	notes := make([]string, len(c.userInfo.notes))
	copy(notes, c.userInfo.notes)
	return notes
}

// LongName 返回命令长名称
func (c *Cmd) LongName() string { return c.userInfo.longName }

// ShortName 返回命令短名称
func (c *Cmd) ShortName() string { return c.userInfo.shortName }

// GetDescription 返回命令描述
func (c *Cmd) GetDescription() string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.description
}

// SetDescription 设置命令描述
func (c *Cmd) SetDescription(desc string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.description = desc
}

// GetHelp 返回命令用法帮助信息
func (c *Cmd) GetHelp() string {
	// 仅对 userInfo.help 访问加锁
	c.rwMu.RLock()
	customHelp := c.userInfo.help
	c.rwMu.RUnlock()

	// 如果用户自定义了帮助信息，则直接返回
	if customHelp != "" {
		return customHelp
	}

	// 自动生成帮助信息
	return generateHelpInfo(c)
}

// SetUsageSyntax 设置自定义命令用法
func (c *Cmd) SetUsageSyntax(usageSyntax string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.usageSyntax = usageSyntax
}

// GetUsageSyntax 获取自定义命令用法
func (c *Cmd) GetUsageSyntax() string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.usageSyntax
}

// SetHelp 设置用户自定义命令帮助信息
func (c *Cmd) SetHelp(help string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.help = help
}

// LoadHelp 从指定文件加载帮助信息
//
// 参数:
// filePath: 帮助信息文件路径
//
// 返回值:
// error: 如果文件不存在或读取文件失败，则返回错误信息
func (c *Cmd) LoadHelp(filePath string) error {
	// 检查是否为空
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// 清理路径并检查有效性
	cleanPath := filepath.Clean(filePath)
	if cleanPath == "" || strings.TrimSpace(cleanPath) == "" {
		return fmt.Errorf("file path cannot be empty or contain only whitespace")
	}

	// 直接读取文件内容并处理可能的错误（包括文件不存在的情况）
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File %s does not exist", filePath)
		}
		return fmt.Errorf("Failed to read file %s: %w", filePath, err)
	}

	// 设置帮助信息
	c.SetHelp(string(content))
	return nil
}

// AddNote 添加备注信息到命令
func (c *Cmd) AddNote(note string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.notes = append(c.userInfo.notes, note)
}

// AddExample 为命令添加使用示例
// description: 示例描述
// usage: 示例使用方式
func (c *Cmd) AddExample(e ExampleInfo) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	// 添加到使用示例列表中
	c.userInfo.examples = append(c.userInfo.examples, e)
}

// GetExamples 获取所有使用示例
// 返回示例切片的副本，防止外部修改
func (c *Cmd) GetExamples() []ExampleInfo {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	examples := make([]ExampleInfo, len(c.userInfo.examples))
	copy(examples, c.userInfo.examples)
	return examples
}

// Args 获取非标志参数切片
func (c *Cmd) Args() []string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	// 返回参数切片副本
	args := make([]string, len(c.args))
	copy(args, c.args)
	return args
}

// Arg 获取指定索引的非标志参数
func (c *Cmd) Arg(i int) string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	// 返回参数
	if i >= 0 && i < len(c.args) {
		return c.args[i]
	}
	return ""
}

// NArg 获取非标志参数的数量
func (c *Cmd) NArg() int {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return len(c.args)
}

// NFlag 获取标志的数量
func (c *Cmd) NFlag() int {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.fs.NFlag()
}

// FlagExists 检查指定名称的标志是否存在
func (c *Cmd) FlagExists(name string) bool {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()

	// 检查标志是否存在
	if _, exists := c.flagRegistry.GetByName(name); exists {
		return true
	}

	return false
}

// PrintHelp 打印命令的帮助信息, 优先打印用户的帮助信息, 否则自动生成帮助信息
func (c *Cmd) PrintHelp() {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	fmt.Println(generateHelpInfo(c))
}

// hasCycle 检测当前命令与待添加子命令间是否存在循环引用
// 循环引用场景包括：
// 1. 子命令直接或间接引用当前命令
// 2. 子命令的父命令链中包含当前命令
// 参数:
//
//	child: 待添加的子命令实例
//
// 返回值:
//
//	存在循环引用返回true，否则返回false
func (c *Cmd) hasCycle(child *Cmd) bool {
	if c == nil || child == nil {
		return false
	}

	visited := make(map[*Cmd]bool)

	// 添加初始深度参数0
	return c.dfs(child, visited, 0)
}

// dfs 深度优先搜索检测循环引用
// 递归检查当前节点及其子命令、父命令链中是否包含目标节点(q)
// 参数:
//
//		current: 当前遍历的命令节点
//		visited: 已访问节点集合，防止重复遍历
//	  depth: 当前递归深度，用于防止无限递归
//
// 返回值:
//
//	找到目标节点返回true, 否则返回false
func (c *Cmd) dfs(current *Cmd, visited map[*Cmd]bool, depth int) bool {
	// 添加递归深度限制(100层)
	if depth > 100 {
		panic(fmt.Sprintf("Potential circular reference detected (recursion depth exceeds %d), there may be circular dependencies between commands", depth))
		//return true // 视为存在循环以中断递归
	}

	// 已访问过当前节点，直接返回避免无限循环
	if visited[current] {
		return false
	}
	visited[current] = true

	// 找到目标节点，存在循环引用
	if current == c {
		return true
	}

	// 递归检查所有子命令
	for _, subCmd := range current.subCmds {
		if c.dfs(subCmd, visited, depth+1) {
			return true
		}
	}

	// 检查父命令链
	if current.parentCmd != nil {
		return c.dfs(current.parentCmd, visited, depth+1)
	}

	return false
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
