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
	c.fs.BoolVar(&showHelp, c.helpFlagName, false, "Show help information")
	// 绑定短帮助标志（若设置）
	if c.helpShortName != "" {
		c.fs.BoolVar(&showHelp, c.helpShortName, false, "Show help information")
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

// isFlagSet 检查指定名称的标志是否被用户显式设置
// 参数:
//
//	name - 标志名称(长格式或短格式)
//
// 返回值:
//
//	bool - true表示标志被显式设置，false表示未设置或使用默认值
func (c *Cmd) isFlagSet(name string) bool {
	// 首先尝试直接查找标志
	flag := c.fs.Lookup(name)
	if flag == nil {
		// 如果找不到，检查是否是短标志
		if longName, ok := c.shortToLong.Load(name); ok {
			flag = c.fs.Lookup(longName.(string))
		} else {
			// 检查是否是长标志的短形式
			if shortName, ok := c.longToShort.Load(name); ok {
				flag = c.fs.Lookup(shortName.(string))
			}
		}
		if flag == nil {
			return false
		}
	}

	// 通过注册表获取标志元数据
	if f, ok := c.flagRegistry[flag.Value]; ok {
		switch f.Type() {
		case FlagTypeBool:
			// 布尔标志特殊处理
			currentVal := flag.Value.String()
			return (currentVal == "true" && flag.DefValue == "false") ||
				(currentVal == "false" && flag.DefValue == "true")
		case FlagTypeInt:
			// 修改整数标志判断逻辑：只要当前值与默认值不同就认为已设置
			return flag.Value.String() != flag.DefValue
		case FlagTypeString:
			// 字符串标志：当前值与默认值不同
			return flag.Value.String() != flag.DefValue
		case FlagTypeFloat:
			// 浮点标志：当前值与默认值不同
			return flag.Value.String() != flag.DefValue
		default:
			// 未知类型：使用默认比较
			return flag.Value.String() != flag.DefValue
		}
	}

	// 默认处理
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
			defValue  string
		}{
			longFlag:  f.Name,
			shortFlag: shortFlag,
			usage:     f.Usage,
			defValue:  fmt.Sprintf("%v", f.DefValue),
		})
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
			helpInfo += fmt.Sprintf(optionTemplate1, flag.shortFlag, flag.longFlag, flag.usage, flag.defValue)
		} else {
			helpInfo += fmt.Sprintf(optionTemplate2, flag.longFlag, flag.usage, flag.defValue)
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

// printUsage 打印帮助内容, 优先显示用户自定义的Usage
func (c *Cmd) printUsage() {
	if c.usage != "" {
		fmt.Println(c.usage)
	} else {
		// 自动生成帮助信息
		fmt.Println(generateHelpInfo(c, c.parentCmd == nil))
	}
}

// hasCycle 检测命令间是否存在循环引用
// parent: 当前命令
// child: 待添加的子命令
// 返回值: 如果存在循环引用则返回true
func hasCycle(parent, child *Cmd) bool {
	current := parent
	for current != nil {
		if current == child {
			return true
		}
		current = current.parentCmd
	}
	return false
}
