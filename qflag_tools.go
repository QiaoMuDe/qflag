package qflag

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// generateHelpInfo 生成命令帮助信息
// cmd: 当前命令
// 返回值: 命令帮助信息
func generateHelpInfo(cmd *Cmd) string {
	var helpInfo string

	// 根据语言选择模板
	var (
		nameTpl, descTpl, optionsHeader, option1Tpl, option2Tpl,
		subCmdsHeader, subCmdTpl, notesHeader, noteItemTpl string
	)

	// 根据语言选择模板
	if cmd.useChinese {
		if cmd.shortName != "" {
			nameTpl = cmdNameWithShortTemplateCN // 命令名（支持短名称显示）
		} else {
			nameTpl = cmdNameTemplateCN // 命令名
		}
		descTpl = cmdDescriptionTemplateCN      // 命令描述
		optionsHeader = optionsHeaderTemplateCN // 选项标题
		option1Tpl = optionTemplate1CN          // 选项模板1
		option2Tpl = optionTemplate2CN          // 选项模板2
		subCmdsHeader = subCmdsHeaderTemplateCN // 子命令标题
		subCmdTpl = subCmdTemplateCN            // 子命令模板
		notesHeader = notesHeaderTemplateCN     // 注意事项标题
		noteItemTpl = noteItemTemplateCN        // 注意事项模板
	} else {
		if cmd.shortName != "" {
			nameTpl = cmdNameWithShortTemplate // 命令名（支持短名称显示）
		} else {
			nameTpl = cmdNameTemplate // 命令名
		}
		descTpl = cmdDescriptionTemplate      // 命令描述
		optionsHeader = optionsHeaderTemplate // 选项标题
		option1Tpl = optionTemplate1          // 选项模板1
		option2Tpl = optionTemplate2          // 选项模板2
		subCmdsHeader = subCmdsHeaderTemplate // 子命令标题
		subCmdTpl = subCmdTemplate            // 子命令模板
		notesHeader = notesHeaderTemplate     // 提示标题
		noteItemTpl = noteItemTemplate        // 提示项模板
	}

	// 命令名（支持短名称显示）
	if cmd.ShortName() != "" {
		helpInfo += fmt.Sprintf(nameTpl, cmd.fs.Name(), cmd.shortName)
	} else {
		helpInfo += fmt.Sprintf(nameTpl, cmd.fs.Name())
	}

	// 命令描述
	if cmd.description != "" {
		helpInfo += fmt.Sprintf(descTpl, cmd.description)
	}

	// 动态生成命令用法
	fullCmdPath := getFullCommandPath(cmd)
	usageLinePrefix := "Usage: "
	if cmd.useChinese {
		usageLinePrefix = "用法: "
	}
	usageLine := usageLinePrefix + fullCmdPath

	// 如果存在子命令，则需要添加子命令用法
	if len(cmd.subCmds) > 0 {
		usageLine += " [subcommand]"
	}

	// 添加选项用法
	usageLine += " [options] [arguments]\n\n"
	helpInfo += usageLine

	// 选项标题
	helpInfo += optionsHeader

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

	// 计算最大标志名称宽度, 用于对齐
	maxWidth := 0
	for _, flag := range flags {
		var nameLength int
		if flag.shortFlag != "" {
			// 格式: "-s, --longname"
			nameLength = len(flag.shortFlag) + len(flag.longFlag) + 5 // 1('-') + 2(', ') + 2('--')
		} else {
			// 格式: "--longname"
			nameLength = len(flag.longFlag) + 2 // 2('--')
		}

		// 更新最大宽度
		if nameLength > maxWidth {
			maxWidth = nameLength
		}
	}

	// 生成排序后的标志信息
	for _, flag := range flags {
		if flag.shortFlag != "" {
			helpInfo += fmt.Sprintf(option1Tpl, flag.shortFlag, maxWidth, flag.longFlag, flag.usage, flag.defValue)
		} else {
			helpInfo += fmt.Sprintf(option2Tpl, maxWidth, flag.longFlag, flag.usage, flag.defValue)
		}
	}

	// 如果有子命令，添加子命令信息
	if len(cmd.subCmds) > 0 {
		helpInfo += subCmdsHeader
		for _, subCmd := range cmd.subCmds {
			helpInfo += fmt.Sprintf(subCmdTpl, subCmd.fs.Name(), subCmd.description)
		}
	}

	// 添加备注
	if len(cmd.GetNotes()) > 0 {
		helpInfo += notesHeader
		for i, note := range cmd.GetNotes() {
			helpInfo += fmt.Sprintf(noteItemTpl, i+1, note)
		}
	}

	return helpInfo
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

// joinErrors 将错误切片合并为单个错误，并去除重复错误
func joinErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}
	if len(errors) == 1 {
		return errors[0]
	}

	// 使用map去重
	uniqueErrors := make(map[string]error)
	for _, err := range errors {
		errStr := err.Error()
		if _, exists := uniqueErrors[errStr]; !exists {
			uniqueErrors[errStr] = err
		}
	}

	// 构建错误信息
	var b strings.Builder
	b.WriteString(fmt.Sprintf("A total of %d unique errors:\n", len(uniqueErrors)))
	i := 1
	for _, err := range uniqueErrors {
		b.WriteString(fmt.Sprintf("  %d. %v\n", i, err))
		i++
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
