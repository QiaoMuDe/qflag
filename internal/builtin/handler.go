// internal/builtin/handler.go
package builtin

import (
	"fmt"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// HandleBuiltinFlags 处理内置标志
// 纯函数设计，返回是否需要退出和错误信息
func HandleBuiltinFlags(ctx *types.CmdContext, printHelpFunc func(*types.CmdContext), generateCompletionFunc func(*types.CmdContext, string) (string, error)) (bool, error) {
	// 处理帮助标志
	if ctx.BuiltinFlags.Help.Get() {
		printHelpFunc(ctx)
		return ctx.Config.ExitOnBuiltinFlags, nil
	}

	// 仅在根命令处理版本和补全标志
	if ctx.Parent == nil {
		// 处理版本标志
		if ctx.BuiltinFlags.Version.Get() && ctx.Config.Version != "" {
			fmt.Println(ctx.Config.Version)
			return ctx.Config.ExitOnBuiltinFlags, nil
		}

		// 处理自动补全标志
		if ctx.Config.EnableCompletion {
			shell := ctx.BuiltinFlags.Completion.Get()
			if shell != flags.ShellNone {
				completion, err := generateCompletionFunc(ctx, shell)
				if err != nil {
					return false, err
				}
				fmt.Println(completion)
				return ctx.Config.ExitOnBuiltinFlags, nil
			}
		}
	}

	// 验证枚举标志
	return false, ValidateEnumFlags(ctx)
}

// ValidateEnumFlags 验证枚举类型标志
func ValidateEnumFlags(ctx *types.CmdContext) error {
	for _, meta := range ctx.FlagRegistry.GetAllFlagMetas() {
		// 跳过非枚举类型标志
		if meta.GetFlagType() != flags.FlagTypeEnum {
			continue
		}

		// 遍历并检查枚举类型标志
		enumFlag, ok := meta.GetFlag().(*flags.EnumFlag)
		if !ok {
			// 获取不到枚举类型标志，跳过
			continue
		}

		// 调用IsCheck方法进行验证
		if err := enumFlag.IsCheck(enumFlag.Get()); err != nil {
			return fmt.Errorf("flag %s: %w", meta.GetName(), err)
		}
	}

	return nil
}
