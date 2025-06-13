package qflag

import (
	"flag"
	"fmt"
)

// bindHelpFlag 绑定-h/--help标志到显示帮助信息的逻辑
func (c *Cmd) bindHelpFlag() {
	if c.helpFlagBound {
		return // 避免重复绑定
	}

	var showHelp bool

	// 绑定长帮助标志
	c.fs.BoolVar(&showHelp, c.helpFlagName, false, "显示帮助信息")
	// 绑定短帮助标志（若设置）
	if c.helpShortName != "" {
		c.fs.BoolVar(&showHelp, c.helpShortName, false, "显示帮助信息")
		c.shortToLong.Store(c.helpShortName, c.helpFlagName) // 存储短到长的映射关系
		c.longToShort.Store(c.helpFlagName, c.helpShortName) // 存储长到短的映射关系
	}

	// 设置帮助标志已绑定
	c.helpFlagBound = true
}

// isHelpRequested 检测帮助标志是否被用户设置
func (c *Cmd) isHelpRequested() bool {
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
func (c *Cmd) isFlagSet(name string) bool {
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

// generateHelpInfo 自动生成帮助信息
// 参数:
//
//	cmd: 命令指针
//	isMainCommand: 是否为主命令，如果为true则会包含子命令信息
//
// 返回:
//
//	生成的帮助信息字符串
func generateHelpInfo(cmd *Cmd, isMainCommand bool) string {
	var helpInfo string

	// 第一行：命令名
	helpInfo += fmt.Sprintf(helpHeaderTemplate, cmd.fs.Name())

	// 第二行：命令描述
	helpInfo += fmt.Sprintf(helpDescriptionTemplate, cmd.Description)

	// 第五行：选项标题
	helpInfo += helpOptionsHeader

	// 添加标志信息
	cmd.fs.VisitAll(func(f *flag.Flag) {
		// 跳过帮助标志
		if f.Name == cmd.helpFlagName || f.Name == cmd.helpShortName {
			return
		}

		// 获取短标志（如果存在）
		shortFlag := ""
		if v, ok := cmd.longToShort.Load(f.Name); ok {
			shortFlag = v.(string)
		}

		// 格式化标志行
		helpInfo += fmt.Sprintf(helpOptionTemplate, shortFlag, f.Name, f.Usage, f.DefValue)
	})

	// 如果有子命令，添加子命令信息
	if isMainCommand && len(cmd.SubCmds) > 0 {
		helpInfo += helpSubCommandsHeader
		for _, subCmd := range cmd.SubCmds {
			helpInfo += fmt.Sprintf(helpSubCommandTemplate, subCmd.fs.Name(), subCmd.Help)
		}
	}

	return helpInfo
}

// printHelp 打印帮助内容，优先显示用户自定义的HelpContent
func (c *Cmd) printHelp() {
	if c.Help != "" {
		fmt.Println(c.Help)
	} else {
		// 自动生成帮助信息，假设当前命令是主命令
		fmt.Println(generateHelpInfo(c, true))
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
