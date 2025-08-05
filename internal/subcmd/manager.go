// internal/subcmd/manager.go
package subcmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
)

// AddSubCommands 添加子命令
// 纯函数设计，通过参数传递父子命令上下文
//
// 参数:
//   - parent: 父命令上下文
//   - children: 子命令上下文列表
//
// 返回值:
//   - error: 错误信息, 无错误时为nil
func AddSubCommands(parent *types.CmdContext, children ...*types.CmdContext) error {
	if len(children) == 0 {
		return fmt.Errorf("子命令列表不能为空")
	}

	// 验证所有子命令
	for i, child := range children {
		if child == nil {
			return fmt.Errorf("索引 %d 的子命令不能为空", i)
		}

		if err := ValidateSubCommand(parent, child); err != nil {
			return fmt.Errorf("无效的子命令 %s: %w", child.GetName(), err)
		}
	}

	// 预分配临时切片(容量=validCmds长度, 避免多次扩容)
	tempList := make([]*types.CmdContext, 0, len(children))

	// 添加子命令到父命令
	parent.Mutex.Lock()
	defer parent.Mutex.Unlock()
	for _, child := range children {
		// 设置子命令的父命令
		child.Parent = parent

		// 如果长标志不为空则添加到子命令映射
		if child.LongName != "" {
			parent.SubCmdMap[child.LongName] = child
		}

		// 如果短标志不为空则添加到子命令映射
		if child.ShortName != "" {
			parent.SubCmdMap[child.ShortName] = child
		}

		// 添加子命令到临时切片
		tempList = append(tempList, child)
	}

	// 一次性添加到目标切片
	parent.SubCmds = append(parent.SubCmds, tempList...)

	return nil
}

// GetSubCommand 获取子命令上下文
//
// 参数:
//   - ctx: 命令上下文
//   - name: 子命令名称
//
// 返回值:
//   - *types.CmdContext: 子命令上下文, 如果不存在则返回nil
func GetSubCommand(ctx *types.CmdContext, name string) *types.CmdContext {
	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()
	return ctx.SubCmdMap[name]
}

// GetAllSubCommands 获取所有子命令
//
// 参数:
//   - ctx: 命令上下文
//
// 返回值:
//   - []*types.CmdContext: 所有子命令
func GetAllSubCommands(ctx *types.CmdContext) []*types.CmdContext {
	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	// 创建一个切片并复制子命令
	result := make([]*types.CmdContext, len(ctx.SubCmds))
	copy(result, ctx.SubCmds)
	return result
}

// SubCommandExists 检查子命令是否存在
//
// 参数:
//   - ctx: 命令上下文
//   - name: 子命令名称
//
// 返回值:
//   - bool: 如果子命令存在则返回true, 否则返回false
func SubCommandExists(ctx *types.CmdContext, name string) bool {
	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()
	_, exists := ctx.SubCmdMap[name]
	return exists
}
