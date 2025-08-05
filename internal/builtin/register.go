// internal/builtin/register.go
package builtin

import (
	"fmt"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// RegisterBuiltinFlags 注册内置标志
// 纯函数设计，通过回调函数注册标志
//
// 参数：
//   - ctx: 命令上下文
//   - registerFunc: 注册标志的回调函数
//     参数：
//   - ctx: 命令上下文
//   - flag: 标志
//   - longName: 标志的长名称
//   - shortName: 标志的短名称
//   - defaultVal: 标志的默认值
//   - usage: 标志的使用说明
func RegisterBuiltinFlags(ctx *types.CmdContext, registerFunc func(*types.CmdContext, flags.Flag, string, string, any, string)) {
	// 注册帮助标志
	registerHelpFlag(ctx, registerFunc)

	// 仅在根命令注册版本和补全标志
	if ctx.Parent == nil {
		registerVersionFlag(ctx, registerFunc)
		registerCompletionFlag(ctx, registerFunc)
	}
}

// registerHelpFlag 注册帮助标志
func registerHelpFlag(ctx *types.CmdContext, registerFunc func(*types.CmdContext, flags.Flag, string, string, any, string)) {
	// 注册帮助标志
	registerFunc(ctx, ctx.BuiltinFlags.Help, flags.HelpFlagName, flags.HelpFlagShortName, false, flags.HelpFlagUsageEn)

	// 标记为内置标志
	ctx.BuiltinFlags.MarkAsBuiltin(flags.HelpFlagName, flags.HelpFlagShortName)
}

// registerVersionFlag 注册版本标志
func registerVersionFlag(ctx *types.CmdContext, registerFunc func(*types.CmdContext, flags.Flag, string, string, any, string)) {
	// 如果没有设置版本信息，则不注册
	if ctx.Config.Version == "" {
		return
	}

	// 获取版本信息英文使用说明
	versionUsage := flags.VersionFlagUsageEn
	if ctx.Config.UseChinese {
		// 如果使用中文，则使用中文使用说明
		versionUsage = flags.VersionFlagUsageZh
	}

	// 注册版本标志
	registerFunc(ctx, ctx.BuiltinFlags.Version, flags.VersionFlagLongName, flags.VersionFlagShortName, false, versionUsage)

	// 标记为内置标志
	ctx.BuiltinFlags.MarkAsBuiltin(flags.VersionFlagLongName, flags.VersionFlagShortName)
}

// registerCompletionFlag 注册补全标志
func registerCompletionFlag(ctx *types.CmdContext, registerFunc func(*types.CmdContext, flags.Flag, string, string, any, string)) {
	// 如果禁用了自动补全，则不注册
	if !ctx.Config.EnableCompletion {
		return
	}

	// 获取补全标志的英文使用说明
	shellDesc := flags.CompletionShellDescEN
	if ctx.Config.UseChinese {
		// 如果启用了自动补全但是使用中文，则使用中文使用说明
		shellDesc = flags.CompletionShellDescCN
	}

	// 这里需要特殊处理枚举标志
	enumFlag := ctx.BuiltinFlags.Completion
	// 设置枚举选项等...

	registerFunc(ctx, enumFlag, flags.CompletionShellFlagLongName, flags.CompletionShellFlagShortName, flags.ShellNone, fmt.Sprintf(shellDesc, flags.ShellSlice))
	ctx.BuiltinFlags.MarkAsBuiltin(flags.CompletionShellFlagLongName, flags.CompletionShellFlagShortName)
}
