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

	// 检测循环引用
	if HasCycle(parent, child) {
		return fmt.Errorf("cyclic reference detected: Command %s already exists in the command chain", GetCmdIdentifier(child))
	}

	// 检查名称冲突
	parent.Mutex.RLock()
	defer parent.Mutex.RUnlock()

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
// 循环引用场景包括：
//  1. 子命令直接或间接引用当前命令
//  2. 子命令的父命令链中包含当前命令
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

	// 创建一个已访问的命令集合
	visited := make(map[*types.CmdContext]bool)

	// 添加初始深度参数0
	return dfs(parent, child, visited, 0)
}

// dfs 深度优先搜索检测循环引用
// 递归检查当前节点及其子命令、父命令链中是否包含目标节点(q)
//
// 参数:
//   - current: 当前节点
//   - target: 目标节点
//   - visited: 已访问的节点集合
//   - depth: 当前递归深度
//
// 返回值:
//   - bool: 是否存在循环引用
func dfs(current, target *types.CmdContext, visited map[*types.CmdContext]bool, depth int) bool {
	if depth > 100 {
		return true // 防止无限递归
	}

	// 已访问过当前节点，直接返回避免无限循环
	if visited[current] {
		return false
	}
	visited[current] = true

	// 找到目标节点，存在循环引用
	if current == target {
		return true
	}

	// 检查子命令
	current.Mutex.RLock()
	subCmds := make([]*types.CmdContext, len(current.SubCmds))
	copy(subCmds, current.SubCmds)
	current.Mutex.RUnlock()

	// 递归检查子命令
	for _, subCmd := range subCmds {
		if dfs(subCmd, target, visited, depth+1) {
			return true
		}
	}

	// 检查父命令
	if current.Parent != nil {
		return dfs(current.Parent, target, visited, depth+1)
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
