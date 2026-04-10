// Package completion 自动补全内部实现
// 本文件包含动态补全逻辑，用于处理 __complete 子命令
package completion

import (
	"fmt"

	"gitee.com/MM-Q/go-kit/fuzzy"
	"gitee.com/MM-Q/qflag/internal/types"
)

// HandleDynamicComplete 处理 __complete 子命令的路由
//
// 参数:
//   - root: 根命令实例，用于查询命令树
//   - instruction: 指令名称
//   - params: 指令参数列表
//
// 返回值:
//   - error: 处理错误
func HandleDynamicComplete(root types.Command, instruction string, params []string) error {
	switch instruction {
	case types.InstructionFuzzy:
		return handleFuzzy(params)
	case types.InstructionContext:
		return HandleContext(root, params)
	case types.InstructionCandidates:
		return HandleCandidates(root, params)
	case types.InstructionEnum:
		return HandleEnum(root, params)
	default:
		return fmt.Errorf("unknown instruction: %s", instruction)
	}
}

// handleFuzzy 处理 fuzzy 指令
//
// 参数:
//   - args: 参数列表，第一个是模式，后面是候选列表
//
// 返回值:
//   - error: 处理错误
//
// 输出格式: 每行一个匹配结果（按匹配质量降序）
func handleFuzzy(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("%s instruction requires at least 2 arguments: pattern and candidates", types.InstructionFuzzy)
	}

	pattern := args[0]     // 模式参数
	candidates := args[1:] // 候选参数

	// 使用 fuzzy 包执行模糊匹配
	matches := fuzzy.Find(pattern, candidates)

	// 输出匹配结果（只输出匹配的字符串）
	for _, match := range matches {
		fmt.Println(match.Str)
	}

	return nil
}
