package parser

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
)

// validateMutexGroups 验证命令的互斥组规则
//
// 参数:
//   - cmd: 要验证的命令
//
// 返回值:
//   - error: 如果互斥组验证失败返回错误
//
// 功能说明:
//   - 检查每个互斥组中是否有多个标志被设置
//   - 检查不允许为空的互斥组中是否有至少一个标志被设置
//   - 提供清晰的错误信息, 指出冲突的标志和组名
//
// 验证规则:
//   - 互斥组中最多只能有一个标志被设置
//   - 如果 AllowNone 为 false, 则必须至少有一个标志被设置
//
// 错误处理:
//   - 使用 types.NewError 创建结构化错误
//   - 错误信息包含互斥组名称和冲突的标志列表
func (p *DefaultParser) validateMutexGroups(cmd types.Command) error {
	// 获取命令配置
	config := cmd.Config()
	if config == nil {
		return nil
	}

	// 检查互斥组是否为空
	if len(config.MutexGroups) == 0 {
		return nil
	}

	// 遍历所有互斥组
	for _, group := range config.MutexGroups {
		setCount := 0
		var setFlags []string

		// 检查互斥组中的每个标志是否被设置
		for _, flagName := range group.Flags {
			if flag, exists := cmd.GetFlag(flagName); exists && flag.IsSet() {
				setCount++
				setFlags = append(setFlags, flagName)
			}
		}

		// 验证互斥组规则
		if setCount > 1 {
			return types.NewError("MUTEX_GROUP_VIOLATION",
				fmt.Sprintf("mutually exclusive flags %v in group '%s' cannot be used together", setFlags, group.Name),
				nil)
		}

		// 验证不允许为空的互斥组
		if !group.AllowNone && setCount == 0 {
			return types.NewError("MUTEX_GROUP_REQUIRED",
				fmt.Sprintf("one of the mutually exclusive flags %v in group '%s' must be set", group.Flags, group.Name),
				nil)
		}
	}

	return nil
}

// validateRequiredGroups 验证命令的必需组规则
//
// 参数:
//   - cmd: 要验证的命令
//
// 返回值:
//   - error: 如果必需组验证失败返回错误
//
// 功能说明:
//   - 检查每个必需组中是否有标志未被设置
//   - 提供清晰的错误信息，指出未设置的标志和组名
//
// 验证规则:
//   - 必需组中的所有标志都必须被设置
//   - 如果有任何一个标志未被设置，返回错误
//
// 错误处理:
//   - 使用 types.NewError 创建结构化错误
//   - 错误信息包含必需组名称和未设置的标志列表
func (p *DefaultParser) validateRequiredGroups(cmd types.Command) error {
	config := cmd.Config()
	if config == nil {
		return nil
	}

	if len(config.RequiredGroups) == 0 {
		return nil
	}

	for _, group := range config.RequiredGroups {
		var unsetFlags []string

		for _, flagName := range group.Flags {
			if flag, exists := cmd.GetFlag(flagName); exists && !flag.IsSet() {
				unsetFlags = append(unsetFlags, flagName)
			}
		}

		if len(unsetFlags) > 0 {
			return types.NewError("REQUIRED_GROUP_VIOLATION",
				fmt.Sprintf("required flags %v in group '%s' must be set", unsetFlags, group.Name),
				nil)
		}
	}

	return nil
}
