// cmd_completion.go - 自动补全命令的实现
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/qflag/flags"
)

// FlagParam 表示标志参数及其需求类型和值类型
type FlagParam struct {
	CommandPath string   // 命令路径，如 "/cmd/subcmd"
	Name        string   // 标志名称(保留原始大小写)
	Type        string   // 参数需求类型: "required"|"optional"|"none"
	ValueType   string   // 参数值类型: "path"|"string"|"number"|"enum"|"bool"等
	EnumOptions []string // 枚举类型的可选值列表
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

# Flag parameters definition - stores type and value type (type|valueType)
declare -A flag_params
%s

# Enum options definition - stores allowed values for enum flags
declare -A enum_options
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
	IFS='|' read -ra opts_arr <<< "${cmd_tree[$context]}"
	opts=$(IFS=' '; echo "${opts_arr[*]}")
	
	# 检查前一个参数是否需要值并获取其类型
	prev_param_type=""
	prev_value_type=""
	if [[ $cword -gt 1 ]]; then
		prev_arg="${words[cword-1]}"
		key="${context}|${prev_arg}"
		prev_param_info=${flag_params[$key]}
		IFS='|' read -r prev_param_type prev_value_type <<< "${prev_param_info}"
	fi

	# 根据参数类型动态生成补全
	if [[ -n "$prev_param_type" && ($prev_param_type == "required" || $prev_param_type == "optional") ]]; then
		case "$prev_value_type" in
			path)
				# 路径类型参数，使用文件和目录补全
				COMPREPLY=($(compgen -f -d -- "${cur}"))
				;;
			number)
				# 数字类型参数，提供基本数字补全
				COMPREPLY=($(compgen -W "$(seq 1 10)" -- "${cur}"))
				;;
			ip)
				# IP地址类型参数，提供基本IP补全
				COMPREPLY=($(compgen -W "192.168. 10.0. 172.16." -- "${cur}"))
				;;
			enum)
				# 枚举类型参数，使用预定义的枚举选项
				COMPREPLY=($(compgen -W "${enum_options[$key]}" -- "${cur}"))
				;;

			url)
				# URL类型参数，提供常见URL前缀补全
				COMPREPLY=($(compgen -W "http:// https:// ftp://" -- "${cur}"))
				;;
			*)
				# 默认值补全
				COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
				;;
			esac
	elif [[ "${cur}" == -* ]]; then
		# 输入以-开头，只显示标志补全
		COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
	else
		# 命令补全，包含文件和目录
		COMPREPLY=($(compgen -W "${opts}" -f -d -- "${cur}"))
	fi

	return $?
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
    '/' = @(%s)
%s      }
	
Register-ArgumentCompleter -CommandName %s -ScriptBlock {
		param($wordToComplete, $commandAst, $cursorPosition, $commandName, $parameterName)

		# Flag parameter requirement array (preserving original case)
		$flagParams = $script:flagParams
	
		# Command tree structure - Pre-initialized outside the function
		$cmdTree = $script:cmdTree
	
		# Parse command line arguments to get the current context
		$context = '/'
		$args = $commandAst.CommandElements | Select-Object -Skip 1 | ForEach-Object { $_.Extent.Text.Trim('"') }
		$index = 0
		$count = $args.Count
	
		while ($index -lt $count) {
			$arg = $args[$index]
			# Use case-sensitive matching to find flags
			$key = "$context|$arg"
			$paramInfo = $flagParams | Where-Object { $_.Name -eq $key } | Select-Object -First 1
			if ($paramInfo) {
				$paramType = $paramInfo.Type
				$valueType = $paramInfo.ValueType
				$index++
				
				# Determine whether to skip the next argument based on the parameter type
				if ($paramType -eq 'required' -or ($paramType -eq 'optional' -and $index -lt $count -and $args[$index] -notlike '-*')) {
					$prevParamInfo = $paramInfo
					$index++
				}
				continue
			}
	
			$nextContext = if ($context) { "$context/$arg" } else { $arg }
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
			$options = $cmdTree[$context] | Where-Object { $_ -ilike "$($wordToComplete.Trim())*" }
		}
	
		# 根据参数类型提供值层补全
			if ($prevParamInfo) {
				$valueType = $prevParamInfo.ValueType
				switch ($valueType) {
					'path' {
						# 路径类型补全
						// 包含隐藏文件并使用完整路径补全
Get-ChildItem -Directory -File -Force | Where-Object { $_.Name -like "$($wordToComplete.Trim())*" } | ForEach-Object {
							[System.Management.Automation.CompletionResult]::new($_.FullName, $_.Name, 'ProviderItem', $_.FullName)
						}
						break
					}
					'number' {
						# 数字类型补全
						1..10 | Where-Object { $_ -like "$($wordToComplete.Trim())*" } | ForEach-Object {
							[System.Management.Automation.CompletionResult]::new($_, $_, 'Number', $_)
						}
						break
					}
					'ip' {
						# IP地址类型补全
						@('192.168.', '10.0.', '172.16.', '127.0.0.') | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
							[System.Management.Automation.CompletionResult]::new($_, $_, 'Text', $_)
						}
						break
					}
					'enum' {
			# 枚举类型参数，使用预定义的枚举选项
			$prevParamInfo.EnumOptions -split '\|' | Where-Object { $_ -and ($_.Trim() -ilike "*$($wordToComplete.Trim())*") } | ForEach-Object {
				[System.Management.Automation.CompletionResult]::new($_, $_, 'Text', $_)
			}
			break
		}

		'url' {
			# URL类型参数，提供常见URL前缀补全
			@('http://', 'https://', 'ftp://') | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
				[System.Management.Automation.CompletionResult]::new($_, $_, 'Text', $_)
			}
			break
		}
		default: {
					# 默认值补全
					$options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_) }
				}
				}
			} else {
				# 参数层补全
				$options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_) }
			}
	}`
	PwshCommandTreeEntry = "    '/%s/' = @(%s)\n"
	// 命令树条目格式
	PwshCommandTreeOption = "    @{ Name = '%s'; Type = '%s'; ValueType = '%s'; EnumOptions = '%s'}\n" // 选项参数需求条目格式
)

// getValueTypeByFlagType 根据标志类型获取值类型
//
// 参数:
//   - flagType - 标志类型
//
// 返回值:
//   - string: 值类型
func getValueTypeByFlagType(flagType flags.FlagType) string {
	switch flagType {
	case flags.FlagTypeBool:
		return "bool"
	case flags.FlagTypeInt, flags.FlagTypeUint16, flags.FlagTypeUint32, flags.FlagTypeUint64, flags.FlagTypeInt64, flags.FlagTypeFloat64:
		return "number"
	case flags.FlagTypePath:
		return "path"
	case flags.FlagTypeEnum:
		return "enum"
	case flags.FlagTypeDuration, flags.FlagTypeTime:
		return "time"
	case flags.FlagTypeIP4, flags.FlagTypeIP6:
		return "ip"
	case flags.FlagTypeURL:
		return "url"
	default:
		return "string"
	}
}

// getParamTypeByFlagType 根据标志类型获取参数需求类型
//
// 参数:
//   - flagType - 标志类型
//
// 返回值:
//   - string: 参数需求类型
func getParamTypeByFlagType(flagType flags.FlagType) string {
	if flagType == flags.FlagTypeBool {
		return "none"
	}
	return "required"
}

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

		// 构建命令路径
		var cmdPath string
		if currentParentPath != "" {
			cmdPath = path.Join(currentParentPath, name, "/")
		} else {
			cmdPath = name
		}

		// 根据shell类型写入命令树条目
		switch shellType {
		case flags.ShellBash: // Bash
			fmt.Fprintf(cmdTreeEntries, BashCommandTreeEntry, cmdPath, strings.Join(cmdOpts, "|"))
		case flags.ShellPwsh, flags.ShellPowershell: // Powershell,和 Pwsh
			// 使用strings.Builder优化字符串拼接性能
			var quotedOptsBuilder strings.Builder
			// 预分配缓冲区减少动态扩容 (假设平均每个选项20字符)
			quotedOptsBuilder.Grow(len(cmdOpts) * 20)
			for i, opt := range cmdOpts {
				if i > 0 {
					quotedOptsBuilder.WriteString(", ")
				}
				quotedOptsBuilder.WriteByte('\'')
				// 高效替换单引号
				quotedOptsBuilder.WriteString(strings.ReplaceAll(opt, "'", "''"))
				quotedOptsBuilder.WriteByte('\'')
			}
			fmt.Fprintf(cmdTreeEntries, PwshCommandTreeEntry, cmdPath, quotedOptsBuilder.String())
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

	// 缓存标志参数，避免重复计算
	params := c.collectFlagParameters()

	// 根据shell类型处理不同的逻辑
	switch shellType {
	case flags.ShellBash: // Bash特定处理
		// 构建标志参数映射
		var flagParamsBuf bytes.Buffer
		var enumOptionsBuf bytes.Buffer
		for _, param := range params {
			var key string
			if param.CommandPath == "" {
				key = param.Name
			} else {
				key = fmt.Sprintf("%s|%s", param.CommandPath, param.Name)
			}
			fmt.Fprintf(&flagParamsBuf, "flag_params[%q]=%q\n", key, param.Type+"|"+param.ValueType)
			if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
				// 使用bytes.Buffer减少内存分配
				var optionsBuf bytes.Buffer
				for i, opt := range param.EnumOptions {
					if i > 0 {
						optionsBuf.WriteString(" ")
					}
					// 增强特殊字符转义处理: 引号、反斜杠和空格
					escapedOpt := strings.ReplaceAll(opt, "\\", "\\\\")
					escapedOpt = strings.ReplaceAll(escapedOpt, "\"", "\\\"")
					escapedOpt = strings.ReplaceAll(escapedOpt, " ", "\\ ")
					optionsBuf.WriteString(escapedOpt)
				}
				options := optionsBuf.String()
				fmt.Fprintf(&enumOptionsBuf, "enum_options[%q]=%q\n", key, options)
			}
		}

		// 写入Bash自动补全脚本头
		fmt.Fprintf(&buf, BashFunctionHeader, strings.Join(rootCmdOpts, "|"), cmdTreeEntries.String(), flagParamsBuf.String(), enumOptionsBuf.String(), programName, programName, programName)

	case flags.ShellPwsh, flags.ShellPowershell: // PowerShell特定处理
		var flagParamsBuf bytes.Buffer
		// 使用缓存的标志参数
		for _, param := range params {
			key := fmt.Sprintf("%s|%s", param.CommandPath, param.Name)
			// 使用strings.Builder优化枚举选项拼接性能
			var enumBuf strings.Builder
			// 预分配缓冲区 (假设平均每个选项20字符)
			enumBuf.Grow(len(param.EnumOptions) * 20)
			first := true
			for _, opt := range param.EnumOptions {
				opt = strings.TrimSpace(opt)
				if opt == "" {
					continue
				}
				if !first {
					enumBuf.WriteString("|")
				}
				// 增强特殊字符转义处理: 引号、反斜杠和空格
				escapedOpt := strings.ReplaceAll(opt, "\\", "\\\\")
				escapedOpt = strings.ReplaceAll(escapedOpt, "\"", "\\\"")
				escapedOpt = strings.ReplaceAll(escapedOpt, " ", "\\ ")
				enumBuf.WriteString(escapedOpt)
				first = false
			}
			enumOptions := enumBuf.String()
			fmt.Fprintf(&flagParamsBuf, PwshCommandTreeOption, key, param.Type, param.ValueType, enumOptions)
		}
		// 写入PowerShell自动补全脚本头
		// 根命令选项数组化处理 - 使用strings.Builder优化
		var rootOptsBuilder strings.Builder
		rootOptsBuilder.Grow(len(rootCmdOpts) * 20) // 预分配缓冲区
		firstOpt := true
		for _, opt := range rootCmdOpts {
			if !firstOpt {
				rootOptsBuilder.WriteString(", ")
			}
			rootOptsBuilder.WriteByte('\'')
			rootOptsBuilder.WriteString(strings.ReplaceAll(opt, "'", "''"))
			rootOptsBuilder.WriteByte('\'')
			firstOpt = false
		}
		fmt.Fprintf(&buf, PwshFunctionHeader, flagParamsBuf.String(), rootOptsBuilder.String(), cmdTreeEntries.String(), programName)
	}

	// 返回自动补全脚本
	return buf.String(), nil
}

// collectFlagParameters 收集所有命令标志参数需求，返回标志名称到参数需求类型的映射
// 参数需求类型: "required"|"optional"|"none"
func (c *Cmd) collectFlagParameters() []FlagParam {
	// 基于命令数量和平均标志密度预分配容量 (每个命令约8个标志)
	initialCapacity := len(c.subCmds)*8 + 16        // 额外+16应对根命令标志
	params := make([]FlagParam, 0, initialCapacity) // 使用切片存储标志参数需求
	seen := make(map[string]bool)                   // 使用原始标志名称作为键，区分大小写

	// 定义匿名函数处理标志添加逻辑，包含参数类型判断和命令路径
	addFlagParam := func(flag *flags.FlagMeta, prefix, opt string, cmdPath string) {
		if opt == "" {
			return
		}

		// 拼接标志名称
		flagName := prefix + opt

		// 只有在标志名称未被添加过时才添加
		if !seen[flagName] {
			seen[flagName] = true // 标记为已添加

			// 根据标志类型获取参数类型和值类型
			flagType := flag.GetFlagType()
			paramType := getParamTypeByFlagType(flagType)
			valueType := getValueTypeByFlagType(flagType)
			var enumOptions []string

			if flagType == flags.FlagTypeEnum {
				if enumFlag, ok := flag.GetFlag().(*flags.EnumFlag); ok {
					enumOptions = enumFlag.GetOptions()
				}
			}

			// 添加标志参数需求，包含命令路径
			params = append(params, FlagParam{CommandPath: cmdPath, Name: flagName, Type: paramType, ValueType: valueType, EnumOptions: enumOptions})
		}
	}

	// 使用队列实现广度优先遍历，记录命令路径
	type cmdNode struct {
		cmd        *Cmd
		parentPath string
	}
	queue := []cmdNode{{cmd: c, parentPath: ""}}

	// 循环遍历命令树
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		cmd := node.cmd
		currentParentPath := node.parentPath

		// 构建当前命令路径，拼接父路径和当前命令名
		cmdPath := "/"
		if currentParentPath != "" {
			cmdPath = currentParentPath + cmd.Name() + "/"
		}

		// 收集当前命令的标志 - 同时处理长短选项
		for _, flag := range cmd.flagRegistry.GetAllFlagMetas() {
			// 处理短选项
			addFlagParam(flag, "-", flag.GetShortName(), cmdPath)
			// 处理长选项
			addFlagParam(flag, "--", flag.GetLongName(), cmdPath)
		}

		// 将子命令加入队列
		for _, subCmd := range cmd.subCmds {
			queue = append(queue, cmdNode{cmd: subCmd, parentPath: cmdPath})
		}
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
			opts = append(opts, "--"+m.GetLongName())
		}
		if m.GetShortName() != "" {
			opts = append(opts, "-"+m.GetShortName())
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
