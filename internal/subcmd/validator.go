// internal/subcmd/validator.go
package subcmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
)

// ValidateSubCommand 验证单个子命令的有效性
//
// 参数：
//   - parent: 当前上下文实例
//   - child: 待添加的上下文实例
//
// 返回值：
//   - error: 验证失败时返回的错误信息, 否则返回nil
func ValidateSubCommand(parent, child *types.CmdContext) error {
	if child == nil {
		return fmt.Errorf("subcmd %s is nil", GetCmdIdentifier(child))
	}

	// 检测循环引用 - 在获取锁之前进行，避免死锁
	if HasCycle(parent, child) {
		return fmt.Errorf("cyclic reference detected: Command %s already exists in the command chain", GetCmdIdentifier(child))
	}

	// 检查名称冲突（外部API层已提供锁保护）
	// 检查长名称冲突
	if child.LongName != "" {
		if _, exists := parent.SubCmdMap[child.LongName]; exists {
			return fmt.Errorf("long name '%s' already exists", child.LongName)
		}
	}

	// 检查短名称冲突
	if child.ShortName != "" {
		if _, exists := parent.SubCmdMap[child.ShortName]; exists {
			return fmt.Errorf("short name '%s' already exists", child.ShortName)
		}
	}

	return nil
}

// HasCycle 检测当前命令与待添加子命令间是否存在循环引用
//
// 使用简化的无锁算法：只检查父命令链，避免复杂的锁操作和递归调用
// 循环引用场景：子命令的父命令链中包含当前命令
//
// 参数:
//   - parent: 当前上下文实例
//   - child: 待添加的上下文实例
//
// 返回值:
//   - bool: 是否存在循环引用
func HasCycle(parent, child *types.CmdContext) bool {
	if parent == nil || child == nil {
		return false
	}

	// 检查是否子命令本身就是父命令（直接循环）
	if child == parent {
		return true
	}

	// 简化的循环检测：只检查父命令链
	// 这避免了复杂的锁操作和递归调用
	current := child.Parent
	depth := 0
	
	for current != nil && depth < 100 {
		// 如果在子命令的父命令链中找到了要添加到的父命令，则存在循环
		if current == parent {
			return true
		}
		current = current.Parent
		depth++
	}

	return false
}

// GetCmdIdentifier 获取命令的标识字符串，用于错误信息
//
// 参数：
//   - cmd: 命令对象
//
// 返回：
//   - 命令标识字符串, 如果为空则返回 <nil>
func GetCmdIdentifier(cmd *types.CmdContext) string {
	if cmd == nil {
		return "<nil>"
	}
	return cmd.GetName()
}
