// internal/subcmd/validator.go
package validator

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

	if parent == nil {
		return nil // 如果父命令为nil，跳过验证
	}

	// 检查父命令的SubCmdMap是否已初始化
	if parent.SubCmdMap == nil {
		panic("父命令的SubCmdMap未初始化")
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
// 检测两种循环情况：
// 1. 子命令的父命令链中包含当前命令
// 2. 子命令的子树中包含当前命令或其祖先
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

	// 检查1：子命令的父命令链中是否包含当前命令
	current := child.Parent
	depth := 0
	for current != nil && depth < 100 {
		if current == parent {
			return true
		}
		current = current.Parent
		depth++
	}

	// 检查2：子命令的子树中是否包含当前命令或其祖先
	// 收集parent及其所有祖先
	ancestors := make(map[*types.CmdContext]bool)
	current = parent
	depth = 0
	for current != nil && depth < 100 {
		ancestors[current] = true
		current = current.Parent
		depth++
	}

	// 检查child的子树中是否包含任何祖先
	visited := make(map[*types.CmdContext]bool)
	return containsAnyAncestor(child, ancestors, visited, 0)
}

// containsAnyAncestor 检查命令树中是否包含任何祖先命令
func containsAnyAncestor(root *types.CmdContext, ancestors map[*types.CmdContext]bool, visited map[*types.CmdContext]bool, depth int) bool {
	if root == nil || depth > 100 {
		return false
	}

	// 避免重复访问
	if visited[root] {
		return false
	}
	visited[root] = true

	// 如果找到任何祖先命令
	if ancestors[root] {
		return true
	}

	// 递归检查所有子命令
	if root.SubCmdMap != nil {
		for _, subCmd := range root.SubCmdMap {
			if subCmd == nil {
				continue
			}
			if containsAnyAncestor(subCmd, ancestors, visited, depth+1) {
				return true
			}
		}
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
