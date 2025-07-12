// cmd_completion.go - 自动补全命令的实现
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/qflag/flags"
)

// FlagParam 表示标志参数及其需求类型
type FlagParam struct {
	Name string // 标志名称(保留原始大小写)
	Type string // 参数需求类型: "required"|"optional"|"none"
}

// 生成标志的注意事项
var (
	// completionNotesCN 中文版本注意事项
	completionNotesCN = []string{
		"Windows环境: 需要PowerShell 5.1或更高版本以支持Register-ArgumentCompleter",
		"Linux环境: 需要bash 4.0或更高版本以支持关联数组特性",
		"请确保您的环境满足上述版本要求，否则自动补全功能可能无法正常工作",
	}

	// completionNotesEN 英文版本注意事项
	completionNotesEN = []string{
		"Windows environment: Requires PowerShell 5.1 or higher to support Register-ArgumentCompleter",
		"Linux environment: Requires bash 4.0 or higher to support associative array features",
		"Please ensure your environment meets the above version requirements, otherwise the auto-completion feature may not work properly",
	}
)

// 内置自动补全命令的示例使用（中文）
var completionExamplesCN = []ExampleInfo{
	{Description: "Linux环境 临时启用", Usage: "source <(%s --generate-shell-completion bash)"},
	{Description: "Linux环境 永久启用(添加到~/.bashrc)", Usage: "echo \"source <(%s --generate-shell-completion bash)\" >> ~/.bashrc"},
	{Description: "Linux环境 系统级安装(至/etc/profile.d)", Usage: "sudo %s --generate-shell-completion bash > /etc/profile.d/qflag_completion.bash"},
	{Description: "Windows环境 临时启用", Usage: "%s --generate-shell-completion powershell | Out-String | Invoke-Expression"},
	{Description: "Windows环境 永久启用(添加到PowerShell配置文件)", Usage: "echo \"%s --generate-shell-completion powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
	{Description: "Windows环境 系统级安装(至ProgramFiles)", Usage: "%s --generate-shell-completion powershell > $env:ProgramFiles\\qflag\\completion.ps1"},
}

// 内置自动补全命令的示例使用（英文）
var completionExamplesEN = []ExampleInfo{
	{Description: "Linux environment temporary activation", Usage: "source <(%s --generate-shell-completion bash)"},
	{Description: "Linux environment permanent activation (add to ~/.bashrc)", Usage: "echo \"source <(%s --generate-shell-completion bash)\" >> ~/.bashrc"},
	{Description: "Linux system-wide installation (to /etc/profile.d)", Usage: "sudo %s --generate-shell-completion bash > /etc/profile.d/qflag_completion.bash"},
	{Description: "Windows environment temporary activation", Usage: "%s --generate-shell-completion powershell | Out-String | Invoke-Expression"},
	{Description: "Windows environment permanent activation (add to PowerShell profile)", Usage: "echo \"%s --generate-shell-completion powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
	{Description: "Windows system-wide installation (to ProgramFiles)", Usage: "%s --generate-shell-completion powershell > $env:ProgramFiles\\qflag\\completion.ps1"},
}

// 补全脚本模板常量
const (
	// Bash补全模板
	BashFunctionHeader = `#!/usr/bin/env bash

# Static command tree definition - Pre-initialized outside the function
declare -A cmd_tree
cmd_tree[/]="%s"
%s

_%s() {
	local cur prev words cword context opts i arg
	COMPREPLY=()

	# Use _get_comp_words_by_ref to get completion parameters for better robustness
	if [[ -z "${_get_comp_words_by_ref}" ]]; then
		# Compatibility with older versions of Bash completion environment
		words=("${COMP_WORDS[@]}")
		cword=$COMP_CWORD
	else
		_get_comp_words_by_ref -n =: cur prev words cword
	fi

	cur="${words[cword]}"
	prev="${words[cword-1]}"

	# Find the current command context
	local context="/"
	local i
	for ((i=1; i < cword; i++)); do
		local arg="${words[i]}"
		if [[ -n "${cmd_tree[$context$arg/]}" ]]; then
			context="$context$arg/"
		fi
	done

	# Get the available options for the current context
	opts="${cmd_tree[$context]}"
	# Add -o filenames option to handle special characters and spaces
	COMPREPLY=($(compgen -o filenames -W "${opts}" -- "${cur}"))

	return 0
	}

complete -F _%s %s
`
	BashCommandTreeEntry = "cmd_tree[/%s/]=\"%s\"\n" // 命令树条目格式

	// PowerShell补全模板
	PwshFunctionHeader = `# Static flag parameter requirement definition - Pre-initialized outside the function
$script:flagParams = @(
%s      )

# Command tree definition - Pre-initialized outside the function
$script:cmdTree = @{
    '' = '%s'
%s      }
	
Register-ArgumentCompleter -CommandName %s -ScriptBlock {
		param($wordToComplete, $commandAst, $cursorPosition, $commandName, $parameterName)

		# Flag parameter requirement array (preserving original case)
		$flagParams = $script:flagParams
	
		# Command tree structure - Pre-initialized outside the function
		$cmdTree = $script:cmdTree
	
		# Parse command line arguments to get the current context
		$context = ''
		$args = $commandAst.CommandElements | Select-Object -Skip 1 | ForEach-Object { $_.Extent.Text.Trim('"') }
		$index = 0
		$count = $args.Count
	
		while ($index -lt $count) {
			$arg = $args[$index]
			# Use case-sensitive matching to find flags
			$paramInfo = $flagParams | Where-Object { $_.Name -ceq $arg } | Select-Object -First 1
			if ($paramInfo) {
				$paramType = $paramInfo.Type
				$index++
				
				# Determine whether to skip the next argument based on the parameter type
				if ($paramType -eq 'required' -or ($paramType -eq 'optional' -and $index -lt $count -and $args[$index] -notlike '-*')) {
					$index++
				}
				continue
			}
	
			$nextContext = if ($context) { "$context.$arg" } else { $arg }
			if ($cmdTree.ContainsKey($nextContext)) {
				$context = $nextContext
				$index++
			} else {
				break
			}
		}
	
		# Get the available options for the current context and filter
		$options = @()
		if ($cmdTree.ContainsKey($context)) {
			$options = $cmdTree[$context] -split '\|' | Where-Object { $_ -like "$wordToComplete*" }
		}
	
		$options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_) }
	}`
	PwshCommandTreeEntry  = "    '%s' = '%s'\n"                  // 命令树条目格式
	PwshCommandTreeOption = "    @{ Name = '%s'; Type = '%s'}\n" // 选项参数需求条目格式
)

// traverseCommandTree 通用的命令树遍历函数
//
// 参数:
//   - cmdTreeEntries - 用于存储命令树条目的缓冲区
//   - parentPath - 父命令路径
//   - cmds - 子命令列表
//   - shellType - shell类型 ("bash", "pwsh", "powershell")
func traverseCommandTree(cmdTreeEntries *bytes.Buffer, parentPath string, cmds []*Cmd, shellType string) {
	// 使用队列实现广度优先遍历
	type cmdNode struct {
		cmd        *Cmd
		parentPath string
	}
	queue := make([]cmdNode, 0, len(cmds))

	// 初始化队列
	for _, cmd := range cmds {
		if cmd != nil {
			queue = append(queue, cmdNode{cmd: cmd, parentPath: parentPath})
		}
	}

	// 定义处理命令名称的匿名函数
	processCmdName := func(name string, currentParentPath string, cmd *Cmd, queue *[]cmdNode) {
		// 检查命令名称和命令是否有效
		if name == "" || cmd == nil {
			return
		}

		// 获取子命令补全选项
		cmdOpts := cmd.collectCompletionOptions()

		// 根据shell类型构建命令路径
		var cmdPath string
		if currentParentPath != "" {
			switch shellType {
			case flags.ShellBash: // Bash
				cmdPath = filepath.Join(currentParentPath, name)
			case flags.ShellPwsh, flags.ShellPowershell: // Powershell ,和 Pwsh
				cmdPath = fmt.Sprintf("%s.%s", currentParentPath, name)
			}
		} else {
			cmdPath = name
		}

		// 根据shell类型写入命令树条目
		switch shellType {
		case flags.ShellBash: // Bash
			fmt.Fprintf(cmdTreeEntries, BashCommandTreeEntry, cmdPath, strings.Join(cmdOpts, " "))
		case flags.ShellPwsh, flags.ShellPowershell: // Powershell,和 Pwsh
			fmt.Fprintf(cmdTreeEntries, PwshCommandTreeEntry, cmdPath, strings.Join(cmdOpts, "|"))
		}

		// 将子命令加入队列
		for _, subCmd := range cmd.subCmds {
			if subCmd != nil {
				*queue = append(*queue, cmdNode{cmd: subCmd, parentPath: cmdPath})
			}
		}
	}

	// 处理队列中的所有命令
	for len(queue) > 0 {
		// 出队
		node := queue[0]                     // 获取当前命令节点
		queue = queue[1:]                    // 移除已处理的元素
		cmd := node.cmd                      // 获取当前命令
		currentParentPath := node.parentPath // 获取当前命令的父路径

		// 处理长命令和短命令
		processCmdName(cmd.LongName(), currentParentPath, cmd, &queue)
		processCmdName(cmd.ShortName(), currentParentPath, cmd, &queue)
	}
}

// generateShellCompletion 生成shell自动补全脚本
//
// 参数:
//   - shellType: shell类型 ("bash", "pwsh", "powershell")
//
// 返回值：
//   - string: 自动补全脚本
//   - error: 错误信息
func (c *Cmd) generateShellCompletion(shellType string) (string, error) {
	// 缓冲区
	var buf bytes.Buffer

	// 检查生成自动补全脚本的必要条件
	if checkErr := validateCompletionGeneration(c); checkErr != nil {
		return "", checkErr
	}

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 获取根命令的补全选项
	rootCmdOpts := c.collectCompletionOptions()

	// 构建命令树条目缓冲区
	var cmdTreeEntries bytes.Buffer

	// 从根命令的子命令开始添加条目
	traverseCommandTree(&cmdTreeEntries, "", c.subCmds, shellType)

	// 根据shell类型处理不同的逻辑
	switch shellType {
	case flags.ShellBash: // Bash特定处理
		// 写入Bash自动补全脚本头
		fmt.Fprintf(&buf, BashFunctionHeader, strings.Join(rootCmdOpts, " "), cmdTreeEntries.String(), programName, programName, programName)
	case flags.ShellPwsh, flags.ShellPowershell: // PowerShell特定处理
		var flagParamsBuf bytes.Buffer
		// 获取标志参数
		for _, param := range c.collectFlagParameters() {
			fmt.Fprintf(&flagParamsBuf, PwshCommandTreeOption, param.Name, param.Type)
		}
		// 写入PowerShell自动补全脚本头
		fmt.Fprintf(&buf, PwshFunctionHeader, flagParamsBuf.String(), strings.Join(rootCmdOpts, "|"), cmdTreeEntries.String(), programName)
	}

	// 返回自动补全脚本
	return buf.String(), nil
}

// collectFlagParameters 收集所有命令标志参数需求，返回标志名称到参数需求类型的映射
// 参数需求类型: "required"|"optional"|"none"
func (c *Cmd) collectFlagParameters() []FlagParam {
	params := make([]FlagParam, 0) // 使用切片存储标志参数需求
	seen := make(map[string]bool)  // 使用原始标志名称作为键，区分大小写

	// 定义匿名函数处理标志添加逻辑，包含参数类型判断
	addFlagParam := func(flag *flags.FlagMeta, prefix, opt string) {
		if opt == "" {
			return
		}

		// 拼接标志名称
		flagName := prefix + opt

		// 只有在标志名称未被添加过时才添加
		if !seen[flagName] {
			seen[flagName] = true // 标记为已添加

			// 根据标志类型设置参数类型
			paramType := "required"
			if flag.GetFlagType() == flags.FlagTypeBool {
				paramType = "none"
			}

			// 添加标志参数需求
			params = append(params, FlagParam{Name: flagName, Type: paramType})
		}
	}

	// 使用队列实现广度优先遍历
	queue := make([]*Cmd, 0, 10)
	queue = append(queue, c)

	// 遍历队列中的所有命令
	for len(queue) > 0 {
		cmd := queue[0]
		queue = queue[1:]

		// 收集当前命令的标志 - 同时处理长短选项
		for _, flag := range cmd.flagRegistry.GetAllFlagMetas() {
			// 处理短选项
			addFlagParam(flag, "-", flag.GetShortName())
			// 处理长选项
			addFlagParam(flag, "--", flag.GetLongName())
		}

		// 将子命令加入队列
		queue = append(queue, cmd.subCmds...)
	}

	return params
}

// collectCompletionOptions 收集命令的补全选项，包括标志和子命令
//
// 参数：
//   - c: 当前命令实例
//
// 返回值：
//   - 包含所有标志选项和子命令名称的字符串切片
func (c *Cmd) collectCompletionOptions() []string {
	// 防御性检查
	if c == nil || c.flagRegistry == nil {
		return []string{}
	}

	// 获取所有标志
	flags := c.flagRegistry.GetAllFlagMetas()

	// 获取所有子命令
	subCmdMap := c.subCmdMap

	// 预分配切片容量，减少动态扩容的开销
	capacity := len(flags)*2 + len(subCmdMap) // 每个标志和子命令(子命令map已包含长短名)最多2个选项
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
	for subCmd := range subCmdMap {
		opts = append(opts, subCmd)
	}

	// 返回所有选项
	return opts
}

// validateCompletionGeneration 验证补全脚本生成所需的命令状态
//
// 参数:
//   - c: 命令实例
//
// 返回值:
//   - error: 如果验证失败, 返回相应的错误信息; 否则返回nil
func validateCompletionGeneration(c *Cmd) error {
	if c == nil {
		return fmt.Errorf("command instance is nil")
	}
	if c.parentCmd != nil {
		return fmt.Errorf("invalid command state: not a root command")
	}
	if c.flagRegistry == nil {
		return fmt.Errorf("invalid command state: flag registry is nil")
	}
	return nil
}
