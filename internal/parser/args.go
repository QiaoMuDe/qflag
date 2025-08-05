// internal/parser/args.go
package parser

import "gitee.com/MM-Q/qflag/internal/types"

// GetArgs 获取参数列表
//
// 参数:
//   - ctx: 命令上下文
//
// 返回值:
//   - []string: 参数列表
func GetArgs(ctx *types.CmdContext) []string {
	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	result := make([]string, len(ctx.Args))
	copy(result, ctx.Args)
	return result
}

// GetArg 获取指定索引的参数
//
// 参数:
//   - ctx: 命令上下文
//   - i: 参数索引
//
// 返回值:
//   - string: 参数值
func GetArg(ctx *types.CmdContext, i int) string {
	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	if i >= 0 && i < len(ctx.Args) {
		return ctx.Args[i]
	}
	return ""
}

// GetArgCount 获取参数数量
//
// 参数:
//   - ctx: 命令上下文
//
// 返回值:
//   - int: 参数数量
func GetArgCount(ctx *types.CmdContext) int {
	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()
	return len(ctx.Args)
}

// SetArgs 设置参数列表（内部使用）
//
// 参数:
//   - ctx: 命令上下文
//   - args: 参数列表
func SetArgs(ctx *types.CmdContext, args []string) {
	ctx.Mutex.Lock()
	defer ctx.Mutex.Unlock()

	ctx.Args = make([]string, len(args))
	copy(ctx.Args, args)
}
