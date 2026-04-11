// Package completion 自动补全内部实现
// 本文件包含动态补全逻辑，用于处理 __complete 子命令
package completion

import (
	"fmt"
	"strings"

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
	case types.InstructionAll:
		return handleAll(root, params)
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

	// 使用 fuzzy 包执行前缀补全
	matches := fuzzy.CompletePrefix(pattern, candidates)

	// 输出匹配结果（只输出匹配的字符串）
	for _, match := range matches {
		fmt.Println(match.Str)
	}

	return nil
}

// handleAll 处理 all 指令，一次性返回所有补全信息
//
// 参数:
//   - root: 根命令实例
//   - args: [cur, prev, cmd_args...]
//     args[0]: cur - 当前输入的词
//     args[1]: prev - 前一个词
//     args[2:]: cmd_args - 已输入的子命令参数
//
// 返回值:
//   - error: 处理错误
//
// 输出格式:
//
//	CONTEXT:<上下文路径>
//	CUR:<当前输入>
//	PREV:<前一个词>
//	CANDIDATES:<候选项列表>
//	ENUM:<枚举值列表>
//	MATCHES:<匹配结果>
//	IS_FLAG:<true|false>
func handleAll(root types.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: __complete all <cur> <prev> [cmd_args...]")
	}

	// 解析参数
	// 去除可能的引号 (PowerShell 等 shell 会保留引号在输入中)
	cur := strings.Trim(args[0], `"'`)  // 当前输入
	prev := strings.Trim(args[1], `"'`) // 前一个词
	cmdArgs := []string{}
	if len(args) > 2 {
		cmdArgs = args[2:] // 子命令参数
	}

	// 1. 计算上下文
	// CalculateContext 期望 tokens 包含程序名作为第一个元素 (从索引1开始遍历)
	// 当 cmdArgs 为空时 (如刚输入命令名后按 Tab) ,tokens = [""] 表示只有程序名
	// 当 cmdArgs = ["config"] 时, tokens = ["", "config"], CalculateContext 能正确识别子命令
	// 所以我们需要在 cmdArgs 前面添加一个空字符串作为占位符
	tokens := append([]string{""}, cmdArgs...)
	contextResult := CalculateContext(root, tokens, len(tokens))
	context := "/"
	if contextResult != nil {
		context = contextResult.Context
	}

	// 2. 获取候选项
	candidates, _ := GetCandidates(root, context)

	// 3. 检查是否是标志值补全
	var enumValues []string
	var matchStrings []string
	isFlag := false

	// 检查 prev 是否是标志：以 "-" 开头
	// 注意：当在根命令下刚输入命令名后按 Tab, prev 是程序名 (如 "dynamic.exe")
	// 正常情况下程序名不会以 "-" 开头，所以不会误判为标志
	// -- 是标志结束符，不是标志，应该按普通参数处理
	if strings.HasPrefix(prev, "-") && prev != "--" {
		isFlag = true
		// 是标志，获取枚举值
		enumValues, _ = GetEnumValues(root, context, prev)

		if len(enumValues) > 0 {
			// 枚举类型：对枚举值进行模糊匹配
			if cur == "" {
				// 空输入时返回所有枚举值
				matchStrings = enumValues
			} else {
				matches := fuzzy.CompletePrefix(cur, enumValues)
				matchStrings = make([]string, len(matches))
				for i, match := range matches {
					matchStrings[i] = match.Str
				}
			}
		}
		// 非枚举类型: matchStrings 保持为空，由 Shell 处理为路径补全
	} else {
		// 不是标志，对候选项进行模糊匹配
		if cur == "" {
			// 空输入时返回所有候选项
			matchStrings = candidates
		} else {
			matches := fuzzy.CompletePrefix(cur, candidates)
			matchStrings = make([]string, len(matches))
			for i, match := range matches {
				matchStrings[i] = match.Str
			}
		}
	}

	// 5. 输出结果（带前缀的多行格式）
	fmt.Printf("CONTEXT:%s\n", context)
	fmt.Printf("CUR:%s\n", cur)
	fmt.Printf("PREV:%s\n", prev)
	fmt.Printf("CANDIDATES:%s\n", strings.Join(candidates, " "))
	fmt.Printf("ENUM:%s\n", strings.Join(enumValues, " "))
	fmt.Printf("MATCHES:%s\n", strings.Join(matchStrings, " "))
	fmt.Printf("IS_FLAG:%v\n", isFlag && len(enumValues) > 0)

	return nil
}
