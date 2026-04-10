package cmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/completion"
	"gitee.com/MM-Q/qflag/internal/types"
)

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
	cmd := NewCmd(types.CompleteCmdName, "", types.ExitOnError)
	opts := &CmdOpts{
		Desc:   "隐藏子命令，用于执行动态补全",
		Hidden: true,
		Examples: map[string]string{
			"执行模糊匹配补全": fmt.Sprintf("%s %s %s <模式> <候选1> [候选2] ...", root.Name(), types.CompleteCmdName, types.InstructionFuzzy),
		},
		RunFunc: func(c types.Command) error {
			args := c.Args()

			// 校验参数数量
			if len(args) < 2 {
				return fmt.Errorf("usage: %s %s <instruction> [args...]", root.Name(), types.CompleteCmdName)
			}

			// 拆分参数：第一个是指令，后续是参数列表
			instruction := args[0]
			params := args[1:]

			// 传递给 completion 包处理
			return completion.HandleDynamicComplete(instruction, params)
		},
	}

	// 设置选项
	if err := cmd.ApplyOpts(opts); err != nil {
		return err
	}

	// 设置补全子命令的父命令
	// 直接设置防止锁问题
	cmd.parent = root

	// 注册子命令
	if err := root.cmdRegistry.Register(cmd); err != nil {
		return fmt.Errorf("register subcommand '%s' failed in '%s': %w", cmd.Name(), root.Name(), err)
	}

	return nil
}
