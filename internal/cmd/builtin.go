package cmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
)

const completeCmdName = "__complete"

// createCompleteCmd 创建动态补全子命令
//
// 返回值:
//   - types.Command: 动态补全子命令实例
//
// 功能说明:
//   - 创建一个隐藏的子命令，用于执行动态补全
//   - 接收当前命令行作为参数
//   - 计算并返回补全结果
func createCompleteCmd() types.Command {
	cmd := NewCmd(completeCmdName, "", types.ExitOnError)
	cmd.SetDesc("内部命令：执行动态补全")
	cmd.SetHidden(true) // 隐藏子命令，不在帮助中显示

	// 设置执行函数
	cmd.SetRun(func(c types.Command) error {
		args := c.Args()
		if len(args) < 2 {
			return fmt.Errorf("用法: __complete <wordToComplete> <previousWord> [args...]")
		}

		return nil
	})

	// 将子命令添加到根命令
	///xxx

	return cmd
}
