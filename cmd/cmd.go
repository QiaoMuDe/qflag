// cmd/cmd.go
package cmd

import (
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag/internal/builtin"
	"gitee.com/MM-Q/qflag/internal/help"
	"gitee.com/MM-Q/qflag/internal/parser"
	"gitee.com/MM-Q/qflag/internal/types"
)

// Cmd 简化的命令结构体，作为适配器连接内部函数式API和外部面向对象API
type Cmd struct {
	ctx *types.CmdContext // 内部上下文，包含所有状态
}

// NewCmd 创建新命令
func NewCmd(longName, shortName string, errorHandling flag.ErrorHandling) *Cmd {
	// 创建内部上下文
	ctx := types.NewCmdContext(longName, shortName, errorHandling)

	// 创建命令实例
	cmd := &Cmd{ctx: ctx}

	// 注册内置标志

	return cmd
}

// Parse 解析命令行参数
func (c *Cmd) Parse(args []string) (err error) {
	// 检查命令是否为nil
	if c == nil {
		return fmt.Errorf("cmd: nil command")
	}

	// 调用提取的组件校验方法
	if err := c.validateComponents(); err != nil {
		return err
	}

	c.ctx.ParseOnce.Do(func() {
		defer c.ctx.Parsed.Store(true)

		// 调用内置标志注册函数
		// 调用内置标志注册方法
		c.registerBuiltinFlags()

		// 添加默认注意事项
		if c.ctx.Config.UseChinese {
			c.ctx.Config.Notes = append(c.ctx.Config.Notes, ChineseTemplate.DefaultNote)
		} else {
			c.ctx.Config.Notes = append(c.ctx.Config.Notes, EnglishTemplate.DefaultNote)
		}

		// 解析参数
		if err = parser.ParseArgs(c.ctx, args, true); err != nil {
			return
		}

		// 处理内置标志
		shouldExit, handleErr := builtin.HandleBuiltinFlags(c.ctx, help.PrintHelp, c.generateShellCompletion)
		if handleErr != nil {
			err = handleErr
			return
		}

		// 内置标志处理是否需要退出程序
		if shouldExit {
			if c.ctx.Config.ExitOnBuiltinFlags {
				os.Exit(0)
			}
			return
		}

		// 执行解析钩子
		if c.ctx.ParseHook != nil {
			hookErr, hookExit := c.ctx.ParseHook(c.ctx)
			if hookErr != nil {
				err = hookErr
				return
			}
			if hookExit {
				if c.ctx.Config.ExitOnBuiltinFlags {
					os.Exit(0)
				}
				return
			}
		}
	})

	return err
}

func (c *Cmd) generateShellCompletion(ctx *types.CmdContext, shell string) (string, error) {
	// 实现生成补全脚本的逻辑
	return "", nil // 简化示例
}
