package qflag

import (
	"fmt"
	"sort"
	"strings"
)

// bindHelpFlagAndShowInstallPathFlag 绑定-h/--help标志到显示帮助信息的逻辑
func (c *Cmd) bindHelpFlagAndShowInstallPathFlag() {
	// 检查是否已绑定
	if c.helpFlagBound {
		return // 避免重复绑定
	}

	// 初始化帮助标志
	c.helpOnce.Do(func() {
		if c.helpFlag == nil {
			c.helpFlag = new(bool) // 为空时自动初始化
		}

		// 绑定帮助标志
		c.BoolVar(c.helpFlag, c.helpFlagName, c.helpFlagShortName, false, "Show help information")

		// 绑定显示安装路径标志
		if c.showInstallPathFlag == nil {
			c.showInstallPathFlag = new(bool)
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
	for _, f := range cmd.flagRegistry {
		flag := f
		flags = append(flags, struct {
			longFlag  string
			shortFlag string
			usage     string
			defValue  string
		}{
			longFlag:  flag.Name(),
			shortFlag: flag.ShortName(),
			usage:     flag.Usage(),
			defValue:  fmt.Sprintf("%v", flag.getDefaultAny()),
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
	if c.FlagExists(name) {
		panic(fmt.Sprintf("Flag name %s already exists", name))
	}
	if c.FlagExists(shortName) {
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
