package completion

import (
	"fmt"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// getValueTypeByFlagType 根据标志类型获取值类型
//
// 参数:
//   - flagType - 标志类型
//
// 返回值:
//   - string: 值类型
func getValueTypeByFlagType(flagType flags.FlagType) string {
	switch flagType {
	case flags.FlagTypeBool:
		return "bool"
	case flags.FlagTypeEnum:
		return "enum"
	default:
		return "string"
	}
}

// getParamTypeByFlagType 根据标志类型获取参数需求类型
//
// 参数:
//   - flagType - 标志类型
//
// 返回值:
//   - string: 参数需求类型
func getParamTypeByFlagType(flagType flags.FlagType) string {
	if flagType == flags.FlagTypeBool {
		return "none"
	}
	return "required"
}

// validateCompletionGeneration 验证补全脚本生成所需的命令状态
//
// 参数:
//   - c: 命令实例
//
// 返回值:
//   - error: 如果验证失败, 返回相应的错误信息; 否则返回nil
func validateCompletionGeneration(ctx *types.CmdContext) error {
	if ctx == nil {
		return fmt.Errorf("command instance is nil")
	}
	if ctx.Parent != nil {
		return fmt.Errorf("invalid command state: not a root command")
	}
	if ctx.FlagRegistry == nil {
		return fmt.Errorf("invalid command state: flag registry is nil")
	}
	return nil
}
