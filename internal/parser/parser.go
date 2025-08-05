package parser

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
	"gitee.com/MM-Q/qflag/qerr"
)

// ParseArgs 解析命令行参数
//
// 参数:
//   - ctx: 命令上下文
//   - args: 命令行参数
//   - parseSubcmds: 是否解析子命令
//
// 返回值:
//   - error: 如果解析失败，返回错误信息
func ParseArgs(ctx *types.CmdContext, args []string, parseSubcmds bool) error {
	// 加载环境变量
	if err := LoadEnvVars(ctx); err != nil {
		return fmt.Errorf("%w: %v", qerr.ErrEnvLoadFailed, err)
	}

	// 解析标志
	if err := ctx.FlagSet.Parse(args); err != nil {
		return fmt.Errorf("%w: %w", qerr.ErrFlagParseFailed, err)
	}

	// 获取解析后的参数
	parsedArgs := ctx.FlagSet.Args()

	// 一次性更新参数，减少锁持有时间
	ctx.Mutex.Lock()
	ctx.Args = append(ctx.Args, parsedArgs...)     // 设置非标志参数
	argsToProcess := make([]string, len(ctx.Args)) // 获取参数的副本, 降低锁持有时间
	copy(argsToProcess, ctx.Args)
	ctx.Mutex.Unlock()

	// 解析子命令: 如果存在子命令并且存在非标志参数, 则解析子命令
	if !parseSubcmds {
		return nil
	}
	if len(argsToProcess) > 0 && (len(ctx.SubCmdMap) > 0 && len(ctx.SubCmds) > 0) {
		return ParseSubCommandSafe(ctx, argsToProcess)
	}

	return nil
}

// ParseSubCommandSafe 安全的子命令解析
//
// 参数:
//   - ctx: 命令上下文
//   - args: 命令行参数
//
// 返回值:
//   - error: 如果解析失败，返回错误信息
func ParseSubCommandSafe(ctx *types.CmdContext, args []string) error {
	// 非标志参数为空，直接返回
	if len(args) == 0 {
		return nil
	}

	// 获取非标志参数的第一个参数(子命令名称)
	subCmdName := args[0]

	// 检查子命令是否存在
	subCmd, exists := ctx.SubCmdMap[subCmdName]
	if exists {
		remainingArgs := args[1:] // 获取除子命令名称外的剩余参数
		return ParseArgs(subCmd, remainingArgs, true)
	}

	return nil
}
