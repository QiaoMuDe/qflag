package qflag

import (
	"flag"
	"fmt"
	"sort"
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

// generateHelpInfo 生成命令帮助信息
// cmd: 当前命令
// isMainCommand: 是否是主命令
// 返回值: 命令帮助信息
func generateHelpInfo(cmd *Cmd, isMainCommand bool) string {
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

	// 命令用法
	if isMainCommand && len(cmd.subCmds) > 0 {
		helpInfo += fmt.Sprintf(cmdUsageWithSubCmdTemplate, cmd.fs.Name())
	} else if cmd.parentCmd != nil {
		helpInfo += fmt.Sprintf(cmdUsageSubCmdTemplate, cmd.parentCmd.fs.Name(), cmd.fs.Name())
	} else {
		helpInfo += fmt.Sprintf(cmdUsageTemplate, cmd.fs.Name())
	}

	// 选项标题
	helpInfo += optionsHeaderTemplate

	// 收集所有标志信息
	var flags []struct {
		longFlag  string
		shortFlag string
		usage     string
		defValue  interface{}
	}

	// 使用Flag接口统一访问标志属性
	for _, f := range cmd.flagRegistry {
		flag := f
		flags = append(flags, struct {
			longFlag  string
			shortFlag string
			usage     string
			defValue  interface{}
		}{
			longFlag:  flag.Name(),
			shortFlag: flag.ShortName(),
			usage:     flag.Usage(),
			defValue:  flag.DefaultValue(),
		})
	}

	cmd.fs.VisitAll(func(f *flag.Flag) {
		// 如果是短标志，跳过处理（会在对应的长标志处理时一并处理）
		if _, ok := cmd.shortToLong.Load(f.Name); ok {
			return
		}

		// 获取短标志（如果存在）
		shortFlag := ""
		if v, ok := cmd.longToShort.Load(f.Name); ok {
			shortFlag = v.(string)
		}

		// 收集标志信息
		flags = append(flags, struct {
			longFlag  string
			shortFlag string
			usage     string
			defValue  interface{}
		}{f.Name, shortFlag, f.Usage, f.DefValue})
	})

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
			helpInfo += fmt.Sprintf(optionTemplate, flag.shortFlag, flag.longFlag, flag.usage, flag.defValue)
		} else {
			helpInfo += fmt.Sprintf("  --%s\t%s (默认值: %s)\n", flag.longFlag, flag.usage, flag.defValue)
		}
	}

	// 如果有子命令，添加子命令信息
	if isMainCommand && len(cmd.subCmds) > 0 {
		helpInfo += subCmdsHeaderTemplate
		for _, subCmd := range cmd.subCmds {
			helpInfo += fmt.Sprintf(subCmdTemplate, subCmd.fs.Name(), subCmd.description)
		}
	}

	return helpInfo
}

// printHelp 打印帮助内容，优先显示用户自定义的HelpContent
func (c *Cmd) printHelp() {
	if c.usage != "" {
		fmt.Println(c.usage)
	} else {
		// 自动生成帮助信息
		fmt.Println(generateHelpInfo(c, c.parentCmd == nil))
	}
	fmt.Println()
}

// errorIf 辅助函数，将非空字符串转为error，空字符串返回nil
func errorIf(cond bool, msg string) error {
	if !cond {
		return nil
	}
	// 使用 %s 格式化字符串，避免非常量格式字符串的问题
	return fmt.Errorf("%s", msg)
}
