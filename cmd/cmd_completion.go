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
	Name        string         // 标志名称(保留原始大小写)
	Type        string         // 参数需求类型: "required"|"optional"|"none"
	FlagType    flags.FlagType // 标志数据类型
	EnumOptions []string       // 枚举类型的可选值列表
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

declare -A flag_types
%s

declare -A enum_options
%s

_%s() {
	local cur prev words cword context opts i arg flag_type enum_vals
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
	flag_type="${flag_types[$prev]}"
	enum_vals="${enum_options[$prev]}"

	# Handle different flag types
	case "$flag_type" in
		enum)
			COMPREPLY=($(compgen -W "$enum_vals" -- "${cur}"))
			;;
		path)
			COMPREPLY=($(compgen -o filenames -d -f -- "${cur}"))
			;;
		*)
			# Default to standard completion with filenames option
			COMPREPLY=($(compgen -o filenames -W "${opts}" -- "${cur}"))
			;;
	esac

	return $?
	}

complete -F _%s %s
`
	BashCommandTreeEntry = "cmd_tree[/%s/]=\"%s\"\n"   // 命令树条目格式
	BashFlagTypeEntry    = "flag_types[%s]=\"%s\"\n"   // 标志类型条目格式
	BashEnumOptionsEntry = "enum_options[%s]=\"%s\"\n" // 枚举选项条目格式
	BashFunctionFooter   = `_%s`

	// PowerShell补全模板
	PwshFunctionHeader = `# Static flag parameter requirement definition - Pre-initialized outside the function
$script:flagParams = @(
%s      )

# Flag type definition
$script:flagTypes = @(
%s      )

# Enum options definition
$script:enumOptions = @(
%s      )

# Command tree definition - Pre-initialized outside the function
$script:cmdTree = @(
%s      )

Register-ArgumentCompleter -CommandName %s -ScriptBlock {
        param($wordToComplete, $commandAst, $cursorPosition, $commandName, $parameterName)

        # Flag parameter requirement array (preserving original case)
        $flagParams = $script:flagParams
        $flagTypes = $script:flagTypes
        $enumOptions = $script:enumOptions
        $cmdTree = $script:cmdTree

        # Parse command line arguments to get the current context
        $context = ''
        $commandArgs = $commandAst.CommandElements | Select-Object -Skip 1 | ForEach-Object { $_.ToString() }
        $index = 0
        $count = $commandArgs.Count

        while ($index -lt $count) {
            $arg = $commandArgs[$index]
            # Use case-sensitive matching to find flags
            $paramInfo = $flagParams | Where-Object { $_.Name -ceq $arg } | Select-Object -First 1
            if ($paramInfo) {
                $paramType = $paramInfo.Type
                $flagType = ($flagTypes | Where-Object { $_.Name -ceq $arg }).Type
                $index++

                # Handle enum type completion
                if ($flagType -eq 'enum' -and $index -eq $count) {
                    $options = ($enumOptions | Where-Object { $_.Name -ceq $arg }).Options -split ' ' | Where-Object { $_ -like "$wordToComplete*" }
                    $options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_) }
                    return
                }

                # Handle path type completion
                if ($flagType -eq 'path' -and $index -eq $count) {
                    $options = @(Get-ChildItem -Path "$wordToComplete*" -Directory -File | ForEach-Object { $_.Name })
                    $options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ProviderItem', $_) }
                    return
                }
                # Determine whether to skip the next argument based on the parameter type
                if ($paramType -eq 'required' -or ($paramType -eq 'optional' -and $index -lt $count -and $commandArgs[$index] -notlike '-*')) {
                    $index++
                }
                continue
            }

            $nextContext = if ($context) { "$context.$arg" } else { $arg }
            if ($cmdTree | Where-Object { $_.Path -ceq $nextContext }) {
                $context = $nextContext
                $index++
            } else {
                break
            }
        }

        # Get the available options for the current context and filter
        $options = @()
        if ($cmdTree.ContainsKey($context)) {
            $options = ($cmdTree | Where-Object { $_.Path -ceq $context }).Options -split ' ' | Where-Object { $_ -like "$wordToComplete*" }
        }

        $options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_) }
        }`
	PwshCommandTreeEntryRoot = "\t@{ Path = ''; Options = '%s' }\n"   // 根命令树条目格式
	PwshCommandTreeEntry     = "\t@{ Path = '%s'; Options = '%s' }\n" // 命令树条目格式
	PwshCommandTreeOption    = "\t@{ Name = '%s'; Type = '%s' }\n"    // 选项参数需求条目格式
)

// addSubCommandsBash 迭代方式添加子命令到命令树，替代递归实现
//
// 参数:
//   - cmdTreeEntries - 用于存储命令树条目的缓冲区
//   - parentPath - 父命令路径
//   - cmds - 子命令列表
func addSubCommandsBash(cmdTreeEntries *bytes.Buffer, parentPath string, cmds []*Cmd) {
	// 使用队列实现广度优先遍历
	type cmdNode struct {
		cmd        *Cmd
		parentPath string
	}
	queue := make([]cmdNode, 0, len(cmds))

	// 初始化队列
	for _, cmd := range cmds {
		queue = append(queue, cmdNode{cmd: cmd, parentPath: parentPath})
	}

	// 定义处理命令名称的匿名函数 (直接传入cmd对象)
	processCmdName := func(name string, currentParentPath string, cmd *Cmd, queue *[]cmdNode) {
		if name == "" || cmd == nil {
			return
		}
		// 使用strings.Builder优化路径拼接
		var cmdPathBuilder strings.Builder
		cmdPathBuilder.WriteString(currentParentPath)
		cmdPathBuilder.WriteString(name)
		cmdPathBuilder.WriteString("/")
		cmdPath := cmdPathBuilder.String()
		cmdPathBuilder.Reset()

		// 去除末尾的斜杠 (通过切片操作避免额外字符串分配)
		var trimmedPath string
		if len(cmdPath) > 0 && cmdPath[len(cmdPath)-1] == '/' {
			// 显式检查斜杠，避免极端情况下的切片越界
			trimmedPath = cmdPath[:len(cmdPath)-1]
		} else {
			trimmedPath = cmdPath // 保留原始路径，避免切片越界
		}

		// 获取子命令补全选项 (通过cmd对象直接调用方法)
		cmdOpts := cmd.collectCompletionOptions()

		// 写入命令树条目
		fmt.Fprintf(cmdTreeEntries, BashCommandTreeEntry, trimmedPath, strings.Join(cmdOpts, " "))

		// 将子命令加入队列 (通过cmd对象访问子命令)
		for _, subCmd := range cmd.subCmds {
			*queue = append(*queue, cmdNode{cmd: subCmd, parentPath: cmdPath})
		}
	}

	// 处理队列中的所有命令
	for len(queue) > 0 {
		// 出队
		node := queue[0]                     // 获取当前命令节点
		queue = queue[1:]                    // 移除已处理的元素
		cmd := node.cmd                      // 获取当前命令
		currentParentPath := node.parentPath // 获取当前命令的父路径

		// 直接传递cmd对象给匿名函数
		processCmdName(cmd.LongName(), currentParentPath, cmd, &queue)
		processCmdName(cmd.ShortName(), currentParentPath, cmd, &queue)
	}
}

// generateBashCompletion 生成Bash自动补全脚本
//
// 返回值：
//   - string: Bash自动补全脚本
func (c *Cmd) generateBashCompletion() (string, error) {
	// 缓冲区
	var buf bytes.Buffer

	// 检查命令树注册表是否为空
	if c == nil {
		return "", fmt.Errorf("command instance is nil")
	}

	// 如果不是父命令则返回错误
	if c.parentCmd != nil {
		return "", fmt.Errorf("invalid command state: not a root command")
	}

	// 检查标志注册表是否为空
	if c.flagRegistry == nil {
		return "", fmt.Errorf("invalid command state: flag registry is nil")
	}

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 获取根命令的补全选项
	rootCmdOpts := c.collectCompletionOptions()

	// 构建命令树条目缓冲区
	var cmdTreeEntries bytes.Buffer

	// 添加根命令选项
	rootOpts := strings.Join(rootCmdOpts, " ")

	// 从根命令的子命令开始添加条目
	addSubCommandsBash(&cmdTreeEntries, "", c.subCmds)

	// 收集标志类型和枚举选项
	var flagTypesBuf, enumOptionsBuf bytes.Buffer
	flagParams := c.collectFlagParameters()
	for _, param := range flagParams {
		if param.FlagType != flags.FlagTypeUnknown {
			fmt.Fprintf(&flagTypesBuf, BashFlagTypeEntry, param.Name, flags.FlagTypeToString(param.FlagType, false))
		}
		if param.FlagType == flags.FlagTypeEnum && len(param.EnumOptions) > 0 {
			opts := strings.Join(param.EnumOptions, " ")
			fmt.Fprintf(&enumOptionsBuf, BashEnumOptionsEntry, param.Name, opts)
		}
	}

	// 写入补全函数头部和命令树, 参数为: 根命令选项, 命令树, 标志类型, 枚举选项, 程序名称, 程序名称, 程序名称
	fmt.Fprintf(&buf, BashFunctionHeader, rootOpts, cmdTreeEntries.String(), flagTypesBuf.String(), enumOptionsBuf.String(), programName, programName, programName)

	return buf.String(), nil
}

// generatePwshCompletion 生成PowerShell自动补全脚本
//
// 返回值：
//   - string: PowerShell自动补全脚本
func (c *Cmd) generatePwshCompletion() (string, error) {
	// 缓冲区
	var buf bytes.Buffer

	// 检查命令树注册表是否为空
	if c == nil {
		return "", fmt.Errorf("command instance is nil")
	}

	// 如果不是父命令则返回错误
	if c.parentCmd != nil {
		return "", fmt.Errorf("invalid command state: not a root command")
	}

	// 检查标志注册表是否为空
	if c.flagRegistry == nil {
		return "", fmt.Errorf("invalid command state: flag registry is nil")
	}

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 构建命令树条目缓冲区
	var cmdTreeEntries bytes.Buffer

	// 获取根命令的补全选项
	rootCmdOpts := c.collectCompletionOptions()

	// 添加根命令选项
	rootOpts := strings.Join(rootCmdOpts, " ")

	// 从根命令的子命令开始添加条目
	addSubCommandsPwsh(&cmdTreeEntries, "", c.subCmds, rootOpts)

	// 构建标志参数需求映射
	var flagParamsBuf, flagTypesBuf, enumOptionsBuf bytes.Buffer

	// 收集当前命令的标志(从根命令开始)
	flagParams := c.collectFlagParameters() // 现在返回[]FlagParam

	var flagTypeEntries []string
	var enumOptionEntries []string

	// 写入标志参数需求条目 - 使用数组而非哈希表
	for _, param := range flagParams {
		fmt.Fprintf(&flagParamsBuf, PwshCommandTreeOption, param.Name, param.Type)

		// 收集标志类型
		if param.FlagType != flags.FlagTypeUnknown {
			entry := fmt.Sprintf(PwshCommandTreeOption, param.Name, flags.FlagTypeToString(param.FlagType, false))
			flagTypeEntries = append(flagTypeEntries, entry)
		}

		// 收集枚举选项
		if param.FlagType == flags.FlagTypeEnum && len(param.EnumOptions) > 0 {
			opts := strings.Join(param.EnumOptions, " ")
			entry := fmt.Sprintf(PwshCommandTreeEntry, param.Name, opts)
			enumOptionEntries = append(enumOptionEntries, entry)
		}
	}

	// 写入标志类型和枚举选项
	if len(flagTypeEntries) > 0 {
		fmt.Fprint(&flagTypesBuf, strings.Join(flagTypeEntries, " "))
	}
	if len(enumOptionEntries) > 0 {
		fmt.Fprint(&enumOptionsBuf, strings.Join(enumOptionEntries, " "))
	}

	// 写入补全函数头部和命令树, 参数为: 标志参数, 标志类型, 枚举选项, 命令树, 程序名称
	fmt.Fprintf(&buf, PwshFunctionHeader, flagParamsBuf.String(), flagTypesBuf.String(), enumOptionsBuf.String(), cmdTreeEntries.String(), programName)

	return buf.String(), nil
}

// addSubCommandsPwsh 迭代方式添加子命令到PowerShell命令树，替代递归实现
//
// 参数:
//   - cmdTreeEntries - 用于存储命令树条目的缓冲区
//   - parentPath - 父命令路径(使用.作为分隔符)
//   - cmds - 子命令列表
func addSubCommandsPwsh(cmdTreeEntries *bytes.Buffer, parentPath string, cmds []*Cmd, rootOpts string) {
	// 使用队列实现广度优先遍历
	type cmdNode struct {
		cmd        *Cmd
		parentPath string
	}
	queue := make([]cmdNode, 0, len(cmds))

	// 初始化队列
	for _, cmd := range cmds {
		queue = append(queue, cmdNode{cmd: cmd, parentPath: parentPath})
	}

	// 收集命令树条目
	var entries []string

	// 添加根命令条目
	entries = append(entries, fmt.Sprintf(PwshCommandTreeEntryRoot, rootOpts))

	// 定义处理命令名称的匿名函数 (抽取重复逻辑)
	processCmdName := func(name string, currentParentPath string, cmd *Cmd, cmdOpts []string) {
		if name == "" {
			return
		}
		// 构建命令路径(使用.作为分隔符)
		cmdPath := currentParentPath
		if cmdPath != "" {
			cmdPath += fmt.Sprintf(".%s", name)
		} else {
			cmdPath = name
		}

		// 添加命令树条目到切片
		entries = append(entries, fmt.Sprintf(PwshCommandTreeEntry, cmdPath, strings.Join(cmdOpts, " ")))

		// 将子命令加入队列
		for _, subCmd := range cmd.subCmds {
			queue = append(queue, cmdNode{cmd: subCmd, parentPath: cmdPath})
		}
	}

	// 处理队列中的所有命令
	for len(queue) > 0 {
		// 出队
		node := queue[0]                     // 获取当前命令节点
		queue = queue[1:]                    // 移除已处理的元素
		cmd := node.cmd                      // 获取当前命令
		currentParentPath := node.parentPath // 获取当前命令的父命令路径

		if cmd == nil {
			continue
		}

		// 获取子命令补全选项
		cmdOpts := cmd.collectCompletionOptions()

		// 处理长命令和短命令 (通过匿名函数消除重复)
		processCmdName(cmd.LongName(), currentParentPath, cmd, cmdOpts)  // 处理长命令
		processCmdName(cmd.ShortName(), currentParentPath, cmd, cmdOpts) // 处理短命令
	}

	// 写入所有条目，用逗号分隔以符合对象数组语法
	cmdTreeEntries.WriteString(strings.Join(entries, " "))
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
			flagType := flag.GetFlagType()
			enumOptions := []string{}

			// 根据标志类型设置参数类型
			switch flagType {
			case flags.FlagTypeBool:
				paramType = "none" // 布尔类型没有参数类型
			case flags.FlagTypeEnum:
				// 尝试将标志转换为EnumFlag以获取选项
				if enumFlag, ok := flag.GetFlag().(*flags.EnumFlag); ok {
					enumOptions = enumFlag.GetOptions()
				}
			}

			// 添加标志参数需求
			params = append(params, FlagParam{
				Name:        flagName,
				Type:        paramType,
				FlagType:    flagType,
				EnumOptions: enumOptions,
			})
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
