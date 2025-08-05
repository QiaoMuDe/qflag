// cmd_internal 包含 Cmd 的内部实现细节，不对外暴露
// 变更需同步更新 cmd.go 中的公共接口文档
package cmd

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// // parseCommon 命令行参数解析公共逻辑
// //
// // 主要功能：
// //  1. 通用参数解析流程(标志解析、内置标志处理、错误处理)
// //  2. 枚举类型标志验证
// //  3. 可选的子命令解析支持
// //
// // 参数：
// //
// //	args: 原始命令行参数切片
// //	parseSubcommands: 是否解析子命令(true: 解析子命令, false: 忽略子命令)
// //
// // 返回值：
// //
// //   - 解析过程中遇到的错误(如标志格式错误、子命令解析失败等)
// //   - 是否需要退出程序, 用于处理内部选项标志的解析处理情况(true: 需要退出, false: 不需要退出)
// //
// // 注意事项：
// //   - 每个Cmd实例仅会被解析一次(线程安全)
// //   - 内置标志(-h/--help, -v/--version等)处理逻辑在此实现
// //   - 子命令解析仅在parseSubcommands=true时执行
// func (c *Cmd) parseCommon(args []string, parseSubcommands bool) (err error, shouldExit bool) {
// 	defer func() {
// 		// 添加panic捕获逻辑
// 		if r := recover(); r != nil {
// 			err = fmt.Errorf("%w: %v\nStack: %s", qerr.ErrPanicRecovered, r, debug.Stack())
// 		}
// 	}()

// 	// 如果命令为空, 则返回错误
// 	if c == nil {
// 		return fmt.Errorf("cmd cannot be nil"), false
// 	}

// 	// 调用提取的组件校验方法
// 	if err = c.validateComponents(); err != nil {
// 		return err, false
// 	}

// 	// 确保只解析一次
// 	c.parseOnce.Do(func() {
// 		defer c.parsed.Store(true) // 在返回时, 无论成功失败均标记为已解析

// 		// 调用内置标志注册方法
// 		c.registerBuiltinFlags()

// 		// 添加默认的注意事项
// 		if c.GetUseChinese() {
// 			c.AddNote(ChineseTemplate.DefaultNote)
// 		} else {
// 			c.AddNote(EnglishTemplate.DefaultNote)
// 		}

// 		// 设置底层flag库的Usage函数
// 		c.fs.Usage = func() {
// 			c.PrintHelp()
// 		}

// 		// 解析前加载环境变量的参数值
// 		if err = c.loadEnvVars(); err != nil {
// 			err = fmt.Errorf("%w: %v", qerr.ErrEnvLoadFailed, err)
// 			return
// 		}

// 		// 调用底层flag库解析参数
// 		if parseErr := c.fs.Parse(args); parseErr != nil {
// 			err = fmt.Errorf("%w: %w", qerr.ErrFlagParseFailed, parseErr)
// 			return
// 		}

// 		// 调用内置标志处理方法
// 		exit, handleErr := c.handleBuiltinFlags()
// 		if handleErr != nil {
// 			// 处理内置标志错误
// 			err = handleErr
// 			return
// 		}

// 		// 内置标志处理是否需要退出程序
// 		if exit {
// 			shouldExit = true
// 			return
// 		}

// 		// 设置当前命令的非标志参数
// 		c.args = append(c.args, c.fs.Args()...)

// 		// 如果允许解析子命令, 则进入子命令解析阶段, 否则跳过子命令解析
// 		if parseSubcommands {
// 			// 如果存在子命令并且非标志参数不为0
// 			if len(c.args) > 0 && (len(c.subCmdMap) > 0 && len(c.subCmds) > 0) {
// 				// 获取参数的第一个值(子命令名称: 长名或短名)
// 				arg := c.args[0]

// 				// 保存剩余参数
// 				remainingArgs := make([]string, len(c.args)-1)
// 				copy(remainingArgs, c.args[1:])

// 				// 直接通过参数0查找子命令, 如果存在则解析子命令
// 				if subCmd, ok := c.subCmdMap[arg]; ok {
// 					// 将剩余参数传递给子命令解析
// 					if parseErr := subCmd.Parse(remainingArgs); parseErr != nil {
// 						err = fmt.Errorf("%w for '%s': %v", qerr.ErrSubCommandParseFailed, arg, parseErr)
// 					}
// 					return
// 				}
// 			}
// 		}

// 		// 调用自定义解析阶段钩子函数
// 		if c.ParseHook != nil {
// 			// 执行钩子函数
// 			hookErr, hookExit := c.ParseHook(c)
// 			// 处理钩子函数错误
// 			if hookErr != nil {
// 				err = hookErr
// 				return
// 			}
// 			// 钩子函数是否需要退出程序
// 			if hookExit {
// 				shouldExit = true
// 				return
// 			}
// 		}
// 	})

// 	// 检查是否报错
// 	if err != nil {
// 		return err, false
// 	}

// 	// 根据内置标志处理结果决定是否退出
// 	return nil, shouldExit
// }

// validateComponents 校验核心组件和内置标志的初始化状态
//
// 返回值：
//   - 组件校验错误信息
func (c *Cmd) validateComponents() error {
	// 核心功能组件校验 (必须初始化)
	if c.ctx.FlagSet == nil {
		return fmt.Errorf("flag.FlagSet instance is not initialized")
	}
	if c.ctx.FlagRegistry == nil {
		return fmt.Errorf("FlagRegistry instance is not initialized")
	}
	if c.ctx.SubCmds == nil {
		return fmt.Errorf("subCmdMap cannot be nil")
	}

	// 内置标志校验
	if c.ctx.BuiltinFlags.Help == nil {
		return fmt.Errorf("help flag is not initialized")
	}
	if c.ctx.BuiltinFlags.Version == nil {
		return fmt.Errorf("version flag is not initialized")
	}

	return nil
}

// registerBuiltinFlags 注册内置标志
// 仅在顶级命令中注册
func (c *Cmd) registerBuiltinFlags() {
	// 仅在顶级命令中注册内置标志
	if c.ctx.Parent != nil {
		return
	}

	// 如果启用了自动补全功能
	if c.ctx.Config.EnableCompletion {
		// 语言配置结构体：集中管理所有语言相关资源
		type languageConfig struct {
			notes     []string      // 注意事项
			shellDesc string        // 用法描述
			examples  []ExampleInfo // 示例
		}

		// 一次性判断语言并初始化所有相关资源
		useChinese := c.ctx.Config.UseChinese
		var langConfig languageConfig

		// 根据语言选择对应的资源
		if useChinese {
			langConfig = languageConfig{
				notes:     completionNotesCN,
				shellDesc: fmt.Sprintf(flags.CompletionShellDescCN, flags.ShellSlice),
				examples:  completionExamplesCN,
			}
		} else {
			langConfig = languageConfig{
				notes:     completionNotesEN,
				shellDesc: fmt.Sprintf(flags.CompletionShellDescEN, flags.ShellSlice),
				examples:  completionExamplesEN,
			}
		}

		// 注册补全标志
		c.EnumVar(c.ctx.BuiltinFlags.Completion, flags.CompletionShellFlagLongName, flags.CompletionShellFlagShortName, flags.ShellNone, langConfig.shellDesc, flags.ShellSlice)

		// 注册到内置的标志名称映射
		c.ctx.BuiltinFlags.NameMap.Store(flags.CompletionShellFlagLongName, true)
		c.ctx.BuiltinFlags.NameMap.Store(flags.CompletionShellFlagShortName, true)

		// 添加自动补全子命令的注意事项
		c.ctx.Config.Notes = append(c.ctx.Config.Notes, langConfig.notes...)

		// 获取运行的程序名
		cmdName := os.Args[0]

		// 添加自动补全子命令的示例
		for _, ex := range langConfig.examples {
			// 直接添加到底层切片中
			c.ctx.Config.Examples = append(c.ctx.Config.Examples, types.ExampleInfo{
				Description: ex.Description,
				Usage:       fmt.Sprintf(ex.Usage, cmdName),
			})
		}
	}

	// 仅在版本信息不为空时注册(-v/--version)
	if c.ctx.Config.Version != "" {
		// 获取版本信息的使用说明
		versionUsage := flags.VersionFlagUsageEn
		if c.ctx.Config.UseChinese {
			versionUsage = flags.VersionFlagUsageZh
		}

		// 注册版本信息标志
		c.BoolVar(c.ctx.BuiltinFlags.Version, flags.VersionFlagLongName, flags.VersionFlagShortName, false, versionUsage)

		// 添加到内置标志名称映射
		c.ctx.BuiltinFlags.NameMap.Store(flags.VersionFlagLongName, true)
		c.ctx.BuiltinFlags.NameMap.Store(flags.VersionFlagShortName, true)
	}
}

// handleBuiltinFlags 处理内置标志(-h/--help, -v/--version, --generate-shell-completion/-gsc等)的逻辑
//
// 返回值:
//   - 是否需要退出程序
//   - 处理过程中遇到的错误
func (c *Cmd) handleBuiltinFlags() (bool, error) {
	// 检查是否使用-h/--help标志
	if c.ctx.BuiltinFlags.Help.Get() {
		//c.PrintHelp()
		if c.ctx.Config.ExitOnBuiltinFlags {
			return true, nil // 标记需要退出
		}
		return false, nil
	}

	// 只有在顶级命令中处理-v/--version标志
	if c.ctx.Parent == nil {
		// 检查是否使用-v/--version标志
		if c.ctx.BuiltinFlags.Version.Get() && c.ctx.Config.Version != "" {
			fmt.Println(c.ctx.Config.Version)
			if c.ctx.Config.ExitOnBuiltinFlags {
				return true, nil // 标记需要退出
			}
			return false, nil
		}
	}

	// 检查是否启用补全功能
	if c.ctx.Config.EnableCompletion {
		// 获取shell类型
		shell := c.ctx.BuiltinFlags.Completion.Get()

		// 只有不是默认值时才生成补全脚本
		if shell != flags.ShellNone {
			// 生成对应shell的补全脚本
			switch shell {
			case flags.ShellBash, flags.ShellPowershell, flags.ShellPwsh: // 受支持的shell类型
				shellCompletion, err := c.generateShellCompletion(c.ctx, shell)
				if err != nil {
					return false, err
				}
				fmt.Println(shellCompletion)
			default:
				return false, fmt.Errorf("unsupported shell: %s. Supported shells are: %v", shell, flags.ShellSlice)
			}

			// 仅在生成补全脚本后检查退出标志
			if c.ctx.Config.ExitOnBuiltinFlags {
				return true, nil // 标记需要退出
			}
			return false, nil
		}
	}

	// 检查枚举类型标志是否有效
	for _, meta := range c.ctx.FlagRegistry.GetAllFlagMetas() {
		// 不是枚举类型标志，跳过
		if meta.GetFlagType() != flags.FlagTypeEnum {
			continue
		}

		// 获取枚举类型标志失败，跳过
		enumFlag, ok := meta.GetFlag().(*flags.EnumFlag)
		if !ok {
			continue
		}

		// 调用IsCheck方法进行验证
		if checkErr := enumFlag.IsCheck(enumFlag.Get()); checkErr != nil {
			// 添加标志名称到错误信息, 便于定位问题
			return false, fmt.Errorf("flag %s: %w", meta.GetName(), checkErr)
		}
	}

	return false, nil
}
