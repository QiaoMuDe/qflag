package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 自动补全命令的标志名称
const (
	CompletionShellFlagLongName  = "shell" // shell 标志名称
	CompletionShellFlagShortName = "s"     // shell 标志名称
)

// 支持的Shell类型
const (
	ShellBash       = "bash"       // bash shell
	ShellZsh        = "zsh"        // zsh shell
	ShellFish       = "fish"       // fish shell
	ShellPowershell = "powershell" // powershell shell
	ShellPwsh       = "pwsh"       // pwsh shell
	ShellNone       = "none"       // 无shell
)

// 支持的Shell类型切片
var ShellSlice = []string{ShellNone, ShellBash, ShellZsh, ShellFish, ShellPowershell, ShellPwsh}

// 内置子命令名称
var (
	CompletionShellLongName  = "completion" // 补全shell命令长名称
	CompletionShellShortName = "comp"       // 补全shell命令短名称
)

// 内置子命令使用说明
var (
	CompletionShellUsageEn = "Generate the autocompletion script for the specified shell" // 补全shell命令英文使用说明
	CompletionShellUsageZh = "生成指定 shell 的自动补全脚本"                                         // 补全shell命令中文使用说明
)

// 内置自动补全命令的示例使用
var completionExamples = []ExampleInfo{
	{"Linux环境 临时启用", "source <(%s completion --shell bash)"},
	{"Linux环境 永久启用(添加到~/.bashrc)", "echo \"source <(%s completion --shell bash)\" >> ~/.bashrc"},
	{"Windows环境 临时启用", "%s completion --shell powershell | Out-String | Invoke-Expression"},
	{"Windows环境 永久启用(添加到PowerShell配置文件)", "echo \"%s completion --shell powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
}

// Bash自动补全脚本模板常量
const (
	// 补全函数头部模板
	FunctionHeader = "#!/usr/bin/env bash\n\n_%s() {\n\tlocal cur prev opts\n\tCOMPREPLY=()\n\tcur=\"${COMP_WORDS[COMP_CWORD]}\"\n\tprev=\"${COMP_WORDS[COMP_CWORD - 1]}\"\n\n\tcase \"${prev}\" in\n"
	// 带子命令别名的模板
	CommandTemplate1 = "\t\t%s|%s)\n\t\t\topts=\"%s\"\n\t\t\tCOMPREPLY=($(compgen -W \"${opts}\" -- ${cur}))\n\t\t\treturn 0\n\t\t\t;;\n"
	// 不带子命令别名的模板
	CommandTemplate2 = "\t\t%s)\n\t\t\topts=\"%s\"\n\t\t\tCOMPREPLY=($(compgen -W \"${opts}\" -- ${cur}))\n\t\t\treturn 0\n\t\t\t;;\n"
	// 命令绑定模板
	BindingCommand = "\t\t*)\n\t\t\t;;\n\tesac\n\t}\n\ncomplete -F _%s %s\n"
)

// generateBashCompletion 生成Bash自动补全脚本
//
// 返回值：
//   - string: Bash自动补全脚本
func (c *Cmd) generateBashCompletion() string {
	// 缓冲区
	var buf bytes.Buffer

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 补全函数头部
	fmt.Fprintf(&buf, FunctionHeader, programName)

	// 获取根命令的补全选项
	rootCmdOpts := c.collectCompletionOptions()

	// 写入根命令的补全选项
	fmt.Fprintf(&buf, CommandTemplate2, programName, strings.Join(rootCmdOpts, " "))

	// 遍历子命令
	if subCmds := c.subCmds; len(subCmds) > 0 {
		for _, cmd := range subCmds {
			// 获取子命令的补全选项
			cmdOpts := cmd.collectCompletionOptions()

			// 写入子命令的补全选项
			fmt.Fprintf(&buf, CommandTemplate1, cmd.LongName(), cmd.ShortName(), strings.Join(cmdOpts, " "))
		}
	}

	// 写入命令绑定模板
	fmt.Fprintf(&buf, BindingCommand, programName, programName)

	return buf.String()
}

// collectCompletionOptions 收集命令的补全选项，包括标志和子命令
// 参数：
//
//	c - 当前命令实例
//
// 返回值：
//
//	包含所有标志选项和子命令名称的字符串切片
func (c *Cmd) collectCompletionOptions() []string {
	// 防御性检查
	if c == nil || c.flagRegistry == nil {
		return []string{}
	}

	// 获取所有标志
	flags := c.flagRegistry.GetAllFlagMetas()

	// 获取所有子命令
	subCmds := c.subCmds

	// 预分配切片容量，减少动态扩容
	capacity := len(flags)*2 + len(subCmds)*2 // 每个标志和子命令最多2个选项
	opts := make([]string, 0, capacity)

	// 获取所有长选项和短选项(为空时不会循环)
	for _, m := range flags {
		if m.GetLongName() != "" {
			opts = append(opts, fmt.Sprint("--", m.GetLongName()))
		}
		if m.GetShortName() != "" {
			opts = append(opts, fmt.Sprint("-", m.GetShortName()))
		}
	}

	// 获取所有子命令(为空时不会循环)
	for _, cmd := range subCmds {
		if cmd.LongName() != "" {
			opts = append(opts, cmd.LongName())
		}
		if cmd.ShortName() != "" {
			opts = append(opts, cmd.ShortName())
		}
	}

	// 返回所有选项
	return opts
}

// HandleCompletionHook 自动补全钩子实现
//
// 功能：处理自动补全子命令逻辑，生成指定shell的补全脚本
//
// 参数：
//   - c: 当前命令实例
//
// 返回值：
//   - error: 处理过程中的错误信息
//   - bool: 是否需要退出程序
func HandleCompletionHook(c *Cmd) (error, bool) {
	// 检查是否启用自动补全
	if !c.enableCompletion {
		return nil, false
	}

	// 获取补全子命令
	rootCmd := c
	for rootCmd.parentCmd != nil { // 追溯到根命令
		rootCmd = rootCmd.parentCmd
	}

	s, ok := rootCmd.subCmdMap[CompletionShellLongName]
	if !ok {
		return nil, false
	}

	// 获取shell类型
	shell := s.completionShell.Get()
	if shell == ShellNone {
		return nil, false
	}

	// 生成对应shell的补全脚本
	switch shell {
	case ShellBash:
		fmt.Println(c.generateBashCompletion())
	case ShellZsh:
		// 实现Zsh补全逻辑
	case ShellFish:
		// 实现Fish补全逻辑
	case ShellPowershell, ShellPwsh:
		// 实现PowerShell补全逻辑
	default:
		return fmt.Errorf("unsupported shell: %s", shell), false
	}

	// 判断是否需要退出
	if c.exitOnBuiltinFlags {
		return nil, true
	}

	return nil, false
}
