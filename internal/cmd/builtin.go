package cmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
)

const completeCmdName = "__complete"

// createCompleteCmd 创建动态补全子命令
//
// 参数:
//   - root: 根命令实例
//
// 返回值:
//   - error: 注册子命令时失败的错误信息
//
// 功能说明:
//   - 创建一个隐藏的子命令，用于执行动态补全
//   - 接收当前命令行作为参数
//   - 计算并返回补全结果
func createCompleteCmd(root *Cmd) error {
	cmd := NewCmd(completeCmdName, "", types.ExitOnError)
	cmd.SetDesc("内部命令：执行动态补全")
	cmd.SetHidden(true) // 隐藏子命令，不在帮助中显示

	// 设置执行函数
	cmd.SetRun(func(c types.Command) error {
		args := c.Args()
		if len(args) < 2 {
			return fmt.Errorf("usage: __complete <wordToComplete> <previousWord> [args...]")
		}

		return nil
	})

	// 设置补全子命令的父命令
	// 直接设置防止锁问题
	cmd.parent = root

	// 注册子命令
	if err := root.cmdRegistry.Register(cmd); err != nil {
		return fmt.Errorf("register subcommand '%s' failed in '%s': %w", cmd.Name(), root.Name(), err)
	}

	return nil
}
