// internal/help/writers.go
package help

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// 帮助信息格式化常量
const (
	// 默认最大宽度，当计算失败时使用
	DefaultMaxWidth = 30

	// 描述信息与选项之间的间距
	DescriptionPadding = 5

	// 子命令名称分隔符长度 (", " 的长度)
	SubCmdSeparatorLen = 2

	// 子命令对齐额外空格数
	SubCmdAlignSpaces = 5

	// 最小填充空格数
	MinPadding = 1
)

// writeModuleHelps 写入自定义模块帮助信息
//
// 参数:
//   - ctx: 命令上下文
//   - buf: 输出缓冲区
func writeModuleHelps(ctx *types.CmdContext, buf *bytes.Buffer) {
	// 空指针检查
	if ctx == nil || buf == nil {
		return
	}

	// 如果存在自定义模块帮助信息，则写入
	if ctx.Config.ModuleHelps != "" {
		buf.WriteString("\n" + ctx.Config.ModuleHelps + "\n")
	}
}

// writeLogoText 写入Logo信息
//
// 参数:
//   - ctx: 命令上下文
//   - buf: 输出缓冲区
func writeLogoText(ctx *types.CmdContext, buf *bytes.Buffer) {
	// 空指针检查
	if ctx == nil || buf == nil {
		return
	}

	// 如果配置了Logo文本, 则写入
	if ctx.Config.LogoText != "" {
		buf.WriteString(ctx.Config.LogoText + "\n")
	}
}

// writeCommandHeader 写入命令名称和描述
//
// 参数:
//   - ctx: 命令上下文
//   - tpl: 模板
//   - buf: 输出缓冲区
func writeCommandHeader(ctx *types.CmdContext, tpl HelpTemplate, buf *bytes.Buffer) {
	// 空指针检查
	if ctx == nil || buf == nil {
		return
	}

	// 修改后的命令名称显示逻辑
	if ctx.LongName != "" && ctx.ShortName != "" {
		// 同时有长短名称
		fmt.Fprintf(buf, tpl.CmdNameWithShort, ctx.LongName, ctx.ShortName)
	} else if ctx.LongName != "" {
		// 只有长名称
		fmt.Fprintf(buf, tpl.CmdName, ctx.LongName)
	} else {
		// 只有短名称
		fmt.Fprintf(buf, tpl.CmdName, ctx.ShortName)
	}

	// 如果描述不为空, 则写入描述
	if ctx.Config.Description != "" {
		fmt.Fprintf(buf, tpl.CmdDescription, ctx.Config.Description)
	}
}

// writeUsageLine 写入用法说明
// ctx: 当前命令
// tpl: 模板实例
// buf: 输出缓冲区
func writeUsageLine(ctx *types.CmdContext, tpl HelpTemplate, buf *bytes.Buffer) {
	// 空指针检查
	if ctx == nil || buf == nil {
		return
	}

	// 使用模板中的用法说明前缀
	usageLinePrefix := tpl.UsagePrefix
	var usageLine string

	// 优先使用用户自定义用法
	if ctx.Config.UsageSyntax != "" {
		usageLine = usageLinePrefix + ctx.Config.UsageSyntax + "\n"
	} else {
		// 获取命令的完整路径
		fullCmdPath := getFullCommandPath(ctx)
		usageLine = usageLinePrefix + fullCmdPath

		// 如果是主命令(父命令为nil)，使用全局选项模板
		if ctx.Parent == nil {
			// 添加子命令部分
			if len(ctx.SubCmds) > 0 {
				usageLine += tpl.UsageGlobalOptions
				usageLine += tpl.UsageSubCmd
			}

			// 添加选项部分
			usageLine += tpl.UsageInfoWithOptions

		} else {
			// 子命令，使用子命令选项模板
			// 如果存在子命令，则添加子命令用法
			if len(ctx.SubCmds) > 0 {
				usageLine += tpl.UsageSubCmd
			}

			// 添加选项部分
			usageLine += tpl.UsageInfoWithOptions
		}
	}

	buf.WriteString(usageLine)
}

// writeOptions 写入选项信息
//
// 参数:
//   - ctx: 命令上下文
//   - tpl: 模板实例
//   - buf: 输出缓冲区
func writeOptions(ctx *types.CmdContext, tpl HelpTemplate, buf *bytes.Buffer) {
	// 空指针检查
	if ctx == nil || buf == nil {
		return
	}

	// 获取所有标志信息并排序
	flags := collectFlags(ctx)
	// 如果没有标志，不显示选项部分
	if len(flags) == 0 {
		return
	}
	// 写入选项头
	buf.WriteString(tpl.OptionsHeader)
	sortFlags(flags)

	// 计算描述信息对齐位置，使用常量替代魔法数字
	maxWidth := calculateMaxWidth(flags)
	if maxWidth == 0 {
		maxWidth = DefaultMaxWidth
	}
	descStartPos := maxWidth + DescriptionPadding

	// 遍历所有标志
	for _, flag := range flags {
		// 格式化选项部分
		optPart := ""

		// 根据标志生成选项部分
		if flag.longFlag != "" && flag.shortFlag != "" {
			// 同时有长短选项
			optPart = fmt.Sprintf(tpl.Option1, flag.longFlag, flag.shortFlag, flag.typeStr)
		} else if flag.longFlag != "" {
			// 只有长选项
			optPart = fmt.Sprintf(tpl.Option2, flag.longFlag, flag.typeStr)
		} else {
			// 只有短选项
			optPart = fmt.Sprintf(tpl.Option3, flag.shortFlag, flag.typeStr)
		}

		// 计算选项部分需要的填充空格，使用常量替代魔法数字
		padding := descStartPos - len(optPart)
		if padding < MinPadding {
			padding = MinPadding
		}

		// 格式化整行输出
		fmt.Fprintf(buf, tpl.OptionDefault,
			optPart, padding, "", flag.usage, flag.defValue)
	}
}

// collectFlags 收集所有标志信息
//
// 参数:
//   - cmd: 命令上下文
//
// 返回:
//   - []flagInfo: 标志信息列表
func collectFlags(cmd *types.CmdContext) []flagInfo {
	// 空指针检查
	if cmd == nil || cmd.FlagRegistry == nil {
		return []flagInfo{}
	}

	var flagInfos []flagInfo

	// 遍历所有标志, 收集标志信息
	for _, f := range cmd.FlagRegistry.GetFlagMetaList() {
		flag := f // 获取标志类型

		// 收集默认值
		defValue := fmt.Sprintf("%v", flag.GetDefault())

		// 对Duration类型进行特殊格式化
		if flag.GetFlagType() == flags.FlagTypeDuration {
			if duration, ok := flag.GetDefault().(time.Duration); ok {
				defValue = duration.String() // 获取时间间隔标志的默认值的字符串表示
			}
		}

		// 创建标志元数据
		flagInfos = append(flagInfos, flagInfo{
			longFlag:  flag.GetLongName(),                               // 长标志名
			shortFlag: flag.GetShortName(),                              // 短标志
			usage:     flag.GetUsage(),                                  // 使用说明
			defValue:  defValue,                                         // 默认值
			typeStr:   flags.FlagTypeToString(flag.GetFlagType(), true), // 标志类型字符串
		})
	}

	return flagInfos
}

// sortFlags 按短标志字母顺序排序，有短标志的选项优先
//
// 参数：
//
//	flags - 需要排序的标志列表
func sortFlags(flags []flagInfo) {
	// 根据标志的短标志名和长标志名进行排序
	sort.Slice(flags, func(i, j int) bool {
		a, b := flags[i], flags[j]
		return sortWithShortNamePriority(
			a.shortFlag != "",
			b.shortFlag != "",
			a.longFlag,
			b.longFlag,
			a.shortFlag,
			b.shortFlag,
		)
	})
}

// sortWithShortNamePriority 通用排序函数
//
// 排序优先级: 1.有短名称的优先 2.按长名称字母序 3.短名称字母序
//
// 参数：
//   - aHasShort: a是否有短名称
//   - bHasShort: b是否有短名称
//   - aName: a的长名称
//   - bName: b的长名称
//   - aShort: a的短名称
//   - bShort: b的短名称
//
// 返回：
//   - bool: a是否应该排在b之前
func sortWithShortNamePriority(aHasShort, bHasShort bool, aName, bName, aShort, bShort string) bool {
	// 1. 有短名称的优先
	if aHasShort != bHasShort {
		return aHasShort
	}

	// 2. 按长名称字母顺序排序
	if aName != bName {
		return aName < bName
	}

	// 3. 都有短名称则按短名称字母顺序排序
	return aShort < bShort
}

// writeSubCmds 写入子命令信息
//
// 参数:
//   - ctx: 命令上下文
//   - tpl: 模板实例
//   - buf: 输出缓冲区
func writeSubCmds(ctx *types.CmdContext, tpl HelpTemplate, buf *bytes.Buffer) {
	// 空指针检查
	if ctx == nil || buf == nil {
		return
	}

	// 没有子命令则返回
	if len(ctx.SubCmds) == 0 {
		return
	}

	// 添加子命令标题
	buf.WriteString(tpl.SubCmdsHeader)

	// 排序子命令：
	// 1. 按长命令名首字母排序
	// 2. 有短命令名的优先
	// 3. 只有长命令名的排最后

	// 获取子命令列表
	sortedSubCmds := ctx.SubCmds

	// 排序子命令
	sort.Slice(sortedSubCmds, func(i, j int) bool {
		a, b := sortedSubCmds[i], sortedSubCmds[j]
		return sortWithShortNamePriority(
			a.ShortName != "",
			b.ShortName != "",
			a.LongName,
			b.LongName,
			a.ShortName,
			b.ShortName,
		)
	})

	// 计算最大命令名长度用于对齐，使用常量替代魔法数字
	maxNameLen := 0
	for _, subCmd := range sortedSubCmds {
		nameLen := len(subCmd.LongName) // 计算长命令名长度
		// 如果子命令有短名称，则计算短名称长度
		if subCmd.ShortName != "" {
			nameLen += len(subCmd.ShortName) + SubCmdAlignSpaces // 使用常量替代魔法数字
		}
		// 更新最大命令名长度
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}
	}

	// 生成对齐的子命令信息
	for _, subCmd := range sortedSubCmds {
		// 构建子命令名称
		var namePart string

		// 根据子命令名称和短名称生成名称部分
		if subCmd.LongName != "" && subCmd.ShortName != "" {
			// 如果长短名称都存在, 则同时显示
			namePart = fmt.Sprintf("%s, %s", subCmd.LongName, subCmd.ShortName)
		} else if subCmd.LongName != "" {
			// 如果只有长短名称中的一个存在, 则只显示一个
			namePart = subCmd.LongName
		} else if subCmd.ShortName != "" {
			// 如果只有长短名称中的一个存在, 则只显示一个
			namePart = subCmd.ShortName
		} else {
			// 如果长短名称都不存在, 则使用默认的注册名
			namePart = subCmd.GetName()
		}

		// 格式化输出，确保描述信息对齐
		fmt.Fprintf(buf, "  %-*s\t%s\n", maxNameLen, namePart, subCmd.Config.Description)
	}
}

// calculateMaxWidth 计算最大标志名称宽度用于对齐
//
// 参数:
//   - flags: 标志信息列表
//
// 返回:
//   - int: 最大标志名称宽度
func calculateMaxWidth(flags []flagInfo) int {
	// 如果没有标志，则返回0
	if len(flags) == 0 {
		return 0
	}

	maxWidth := 0
	for _, flag := range flags {
		var nameLength int
		if flag.longFlag != "" && flag.shortFlag != "" {
			// 同时有长短选项: --longFlag, -shortFlag <type>
			nameLength = len(flag.longFlag) + len(flag.shortFlag) + len(flag.typeStr) + 8 // 2('--') + 1(' ') + 2('<type>') + 2(', ') + 1('-')
		} else if flag.longFlag != "" {
			// 只有长选项: --longFlag <type>
			nameLength = len(flag.longFlag) + len(flag.typeStr) + 5 // 2('--') + 1(' ') + 2('<type>')
		} else {
			// 只有短选项: -shortFlag <type>
			nameLength = len(flag.shortFlag) + len(flag.typeStr) + 4 // 1('-') + 1(' ') + 2('<type>')
		}

		// 如果名称长度大于最大宽度，则更新最大宽度
		if nameLength > maxWidth {
			maxWidth = nameLength
		}
	}

	return maxWidth
}

// writeExamples 写入示例信息
//
// 参数:
//   - ctx: 命令上下文
//   - tpl: 模板实例
//   - buf: 输出缓冲区
func writeExamples(ctx *types.CmdContext, tpl HelpTemplate, buf *bytes.Buffer) {
	// 空指针检查
	if ctx == nil || buf == nil {
		return
	}

	// 如果没有示例信息，则返回
	examples := ctx.Config.Examples
	if len(examples) == 0 {
		return
	}

	// 添加示例信息标题
	buf.WriteString(tpl.ExamplesHeader)

	// 遍历添加示例信息
	for i, example := range examples {
		// 格式化示例信息
		fmt.Fprintf(buf, tpl.ExampleItem, i+1, example.Description, example.Usage)

		// 如果不是最后一个示例，添加空行
		if i < len(examples)-1 {
			fmt.Fprintln(buf)
		}
	}
}

// writeNotes 写入注意事项
//
// 参数:
//   - ctx: 命令上下文
//   - tpl: 模板实例
//   - buf: 输出缓冲区
func writeNotes(ctx *types.CmdContext, tpl HelpTemplate, buf *bytes.Buffer) {
	// 空指针检查
	if ctx == nil || buf == nil {
		return
	}

	// 如果没有注意事项，则返回
	notes := ctx.Config.Notes
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

// getFullCommandPath 递归构建完整的命令路径，从根命令到当前命令
//
// 参数:
//   - ctx: 命令上下文
//
// 返回:
//   - string: 完整的命令路径
func getFullCommandPath(ctx *types.CmdContext) string {
	// 空指针检查
	if ctx == nil || ctx.FlagSet == nil {
		return ""
	}

	if ctx.Parent == nil {
		// 如果没有父命令，则直接返回命令名
		return ctx.FlagSet.Name()
	}

	// 递归获取上一层命令的完整路径: 上一层的完整路径 + " " + 当前的命令名
	return getFullCommandPath(ctx.Parent) + " " + ctx.FlagSet.Name()
}
