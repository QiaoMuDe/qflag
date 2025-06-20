package qflag

import (
	"bytes"
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
	// 根据语言选择模板实例
	var tpl HelpTemplate
	if cmd.useChinese {
		tpl = ChineseTemplate // 中文模板实例
	} else {
		tpl = EnglishTemplate // 英文模板实例
	}

	// 使用buffer提高字符串拼接性能
	var buf bytes.Buffer

	// 写入命令名称和描述
	writeCommandHeader(cmd, tpl, &buf)

	// 写入用法说明
	writeUsageLine(cmd, &buf)

	// 写入选项信息
	if err := writeOptions(cmd, tpl, &buf); err != nil {
		return fmt.Sprintf("生成帮助信息失败: %v", err)
	}

	// 写入子命令信息
	writeSubCmds(cmd, tpl, &buf)

	// 写入注意事项
	writeNotes(cmd, tpl, &buf)

	return buf.String()
}

// writeCommandHeader 写入命令名称和描述
func writeCommandHeader(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) {
	// 如果子命令不为空, 则使用带有短名称的模板
	if cmd.ShortName() != "" {
		fmt.Fprintf(buf, tpl.CmdNameWithShort, cmd.fs.Name(), cmd.shortName)
	} else {
		fmt.Fprintf(buf, tpl.CmdName, cmd.fs.Name())
	}

	// 如果描述不为空, 则写入描述
	if cmd.description != "" {
		fmt.Fprintf(buf, tpl.CmdDescription, cmd.description)
	}
}

// writeUsageLine 写入用法说明
func writeUsageLine(cmd *Cmd, buf *bytes.Buffer) {
	// 根据语言选择用法说明前缀
	usageLinePrefix := "Usage: "
	if cmd.useChinese {
		usageLinePrefix = "用法: "
	}

	// 获取命令的完整路径
	fullCmdPath := getFullCommandPath(cmd)
	usageLine := usageLinePrefix + fullCmdPath

	// 如果存在子命令，则需要添加子命令用法
	if len(cmd.subCmds) > 0 {
		usageLine += " [subcommand]"
	}

	// 添加用法说明
	usageLine += " [options] [arguments]\n\n"

	buf.WriteString(usageLine)
}

// writeOptions 写入选项信息
func writeOptions(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) error {
	// 写入选项头
	buf.WriteString(tpl.OptionsHeader)

	// 获取所有标志信息并排序
	flags := collectFlags(cmd)
	sortFlags(flags)

	// 计算描述信息对齐位置
	const descStartPos = 30 // 描述信息开始位置

	for _, flag := range flags {
		// 格式化选项部分
		optPart := ""
		if flag.shortFlag != "" {
			optPart = fmt.Sprintf(tpl.Option1, flag.shortFlag, flag.longFlag)
		} else {
			optPart = fmt.Sprintf(tpl.Option2, flag.longFlag)
		}

		// 计算选项部分需要的填充空格
		padding := descStartPos - len(optPart)
		if padding < 1 {
			padding = 1
		}

		// 格式化整行输出
		fmt.Fprintf(buf, tpl.OptionDefault,
			optPart, padding, "", flag.usage, flag.defValue)
	}

	return nil
}

// collectFlags 收集所有标志信息
func collectFlags(cmd *Cmd) []struct {
	longFlag  string
	shortFlag string
	usage     string
	defValue  string
} {
	var flags []struct {
		longFlag  string // 长标志
		shortFlag string // 短标志
		usage     string // 用法说明
		defValue  string // 默认值
	}

	// 遍历所有标志, 收集标志信息
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

	return flags
}

// sortFlags 按短标志字母顺序排序，有短标志的选项优先
func sortFlags(flags []struct {
	longFlag  string
	shortFlag string
	usage     string
	defValue  string
}) {
	sort.Slice(flags, func(i, j int) bool {
		// 组合短标志和长标志，以便按字母顺序排序
		a, b := flags[i], flags[j]

		// 优先处理有短标志的选项
		aHasShort := a.shortFlag != ""

		// 处理没有短标志的选项
		bHasShort := b.shortFlag != ""

		// 如果a有短标志而b没有，则a排在前面
		if aHasShort && !bHasShort {
			return true
		}

		// 如果a没有短标志而b有，则b排在前面
		if !aHasShort && bHasShort {
			return false
		}

		// 如果a和b都有短标志，则按短标志字母顺序排序
		if aHasShort && bHasShort {
			return a.shortFlag < b.shortFlag
		}

		return a.longFlag < b.longFlag
	})
}

// calculateMaxWidth 计算最大标志名称宽度用于对齐
func calculateMaxWidth(flags []struct {
	longFlag  string
	shortFlag string
	usage     string
	defValue  string
}) (int, error) {
	// 如果没有标志，则返回0
	if len(flags) == 0 {
		return 0, nil
	}

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

		// 如果名称长度大于最大宽度，则更新最大宽度
		if nameLength > maxWidth {
			maxWidth = nameLength
		}
	}

	return maxWidth, nil
}

// // writeSubCmds 写入子命令信息
// func writeSubCmds(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) {
// 	// 如果没有子命令，则返回
// 	if len(cmd.subCmds) == 0 {
// 		return
// 	}

// 	// 添加子命令标题
// 	buf.WriteString(tpl.SubCmdsHeader)

// 	// 遍历所有子命令，生成子命令信息
// 	for _, subCmd := range cmd.subCmds {
// 		if subCmd.shortName != "" {
// 			fmt.Fprintf(buf, tpl.SubCmdWithShort, subCmd.fs.Name(), subCmd.shortName, subCmd.description)
// 		} else {
// 			fmt.Fprintf(buf, tpl.SubCmd, subCmd.fs.Name(), subCmd.description)
// 		}
// 	}
// }

// writeSubCmds 写入子命令信息
func writeSubCmds(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) {
	if len(cmd.subCmds) == 0 {
		return
	}

	// 添加子命令标题
	buf.WriteString(tpl.SubCmdsHeader)

	// 排序子命令：
	// 1. 按长命令名首字母排序
	// 2. 有短命令名的优先
	// 3. 只有长命令名的排最后
	sortedSubCmds := make([]*Cmd, len(cmd.subCmds))
	copy(sortedSubCmds, cmd.subCmds)
	// sort.Slice(sortedSubCmds, func(i, j int) bool {
	// 	a, b := sortedSubCmds[i], sortedSubCmds[j]

	// 	// 首先按长命令名排序
	// 	if a.fs.Name() != b.fs.Name() {
	// 		return a.fs.Name() < b.fs.Name()
	// 	}

	// 	// 然后有短命令名的优先
	// 	aHasShort := a.shortName != ""
	// 	bHasShort := b.shortName != ""

	// 	if aHasShort && !bHasShort {
	// 		return true
	// 	}
	// 	if !aHasShort && bHasShort {
	// 		return false
	// 	}

	// 	// 最后按短命令名排序
	// 	return a.shortName < b.shortName
	// })
	sort.Slice(sortedSubCmds, func(i, j int) bool {
		a, b := sortedSubCmds[i], sortedSubCmds[j]

		// 1. 有短命令名的优先
		aHasShort := a.shortName != ""
		bHasShort := b.shortName != ""

		//  如果a有短命令名而b没有，则a排在前面
		if aHasShort && !bHasShort {
			return true
		}

		//   如果a没有短命令名而b有，则b排在前面
		if !aHasShort && bHasShort {
			return false
		}

		// 2. 按长命令名首字母排序
		if a.fs.Name() != b.fs.Name() {
			return a.fs.Name() < b.fs.Name()
		}

		// 3. 只有长命令名的排最后
		// 如果到这里，说明 a 和 b 的长命令名相同，且要么都有短命令名，要么都没有短命令名
		return a.shortName < b.shortName
	})

	// 计算最大命令名长度用于对齐
	maxNameLen := 0
	for _, subCmd := range sortedSubCmds {
		nameLen := len(subCmd.fs.Name())
		if subCmd.shortName != "" {
			nameLen += len(subCmd.shortName) + 2 // 加逗号和空格
		}
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}
	}

	// 生成对齐的子命令信息
	for _, subCmd := range sortedSubCmds {
		namePart := subCmd.fs.Name()
		if subCmd.shortName != "" {
			namePart = fmt.Sprintf("%s, %s", subCmd.fs.Name(), subCmd.shortName)
		}

		// 格式化输出，确保描述信息对齐
		fmt.Fprintf(buf, "  %-*s\t%s\n", maxNameLen, namePart, subCmd.description)
	}
}

// writeNotes 写入注意事项
func writeNotes(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) {
	// 如果没有注意事项，则返回
	notes := cmd.GetNotes()
	if len(notes) == 0 {
		return
	}

	// 添加注意事项标题
	buf.WriteString(tpl.NotesHeader)

	// 遍历添加注意事项
	for i, note := range notes {
		fmt.Fprintf(buf, tpl.NoteItem, i+1, note)
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
