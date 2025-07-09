package cmd

import (
	"bytes"
	"flag"
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

// 自动补全命令的标志名称
const (
	CompletionShellFlagLongName  = "shell" // shell 标志名称
	CompletionShellFlagShortName = "s"     // shell 标志名称
)

// 支持的Shell类型
const (
	ShellBash       = "bash"       // bash shell
	ShellPowershell = "powershell" // powershell shell
	ShellPwsh       = "pwsh"       // pwsh shell
	ShellNone       = "none"       // 无shell
)

// 支持的Shell类型切片
var ShellSlice = []string{ShellNone, ShellBash, ShellPowershell, ShellPwsh}

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

// 补全注意事项
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

// 内置自动补全命令的示例使用
var completionExamples = []ExampleInfo{
	{"Linux环境 临时启用", "source <(%s completion --shell bash)"},
	{"Linux环境 永久启用(添加到~/.bashrc)", "echo \"source <(%s completion --shell bash)\" >> ~/.bashrc"},
	{"Windows环境 临时启用", "%s completion --shell powershell | Out-String | Invoke-Expression"},
	{"Windows环境 永久启用(添加到PowerShell配置文件)", "echo \"%s completion --shell powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
}

// 补全脚本模板常量
const (
	// Bash补全模板
	//BashFunctionHeader   = "#!/usr/bin/env bash\n\n_%s() {\n\tlocal cur prev words cword context opts i arg\n\tCOMPREPLY=()\n\n\t# 使用_get_comp_words_by_ref获取补全参数, 提高健壮性\n\tif [[ -z \"${_get_comp_words_by_ref}\" ]]; then\n\t\t# 兼容旧版本Bash补全环境\n\t\twords=(\"${COMP_WORDS[@]}\")\n\t\tcword=$COMP_CWORD\n\telse\n\t\t_get_comp_words_by_ref -n =: cur prev words cword\n\tfi\n\n\tcur=\"${words[cword]}\"\n\tprev=\"${words[cword-1]}\"\n\n\t# 构建命令树结构\n\tdeclare -A cmd_tree\n\tcmd_tree[/]=\"%s\"\n%s\n\n\t# 查找当前命令上下文\n\tlocal context=\"/\"\n\tlocal i\n\tfor ((i=1; i < cword; i++)); do\n\t\tlocal arg=\"${words[i]}\"\n\t\tif [[ -n \"${cmd_tree[$context$arg/]}\" ]]; then\n\t\t\tcontext=\"$context$arg/\"\n\t\tfi\n\tdone\n\n\t# 获取当前上下文可用选项\n\topts=\"${cmd_tree[$context]}\"\n\tCOMPREPLY=($(compgen -W \"${opts}\" -- ${cur}))\n\treturn 0\n\t}\n\ncomplete -F _%s %s\n" // Bash补全函数头部
	BashFunctionHeader = `#!/usr/bin/env bash

_%s() {
	local cur prev words cword context opts i arg
	COMPREPLY=()

	# 使用_get_comp_words_by_ref获取补全参数, 提高健壮性
	if [[ -z "${_get_comp_words_by_ref}" ]]; then
		# 兼容旧版本Bash补全环境
		words=("${COMP_WORDS[@]}")
		cword=$COMP_CWORD
	else
		_get_comp_words_by_ref -n =: cur prev words cword
	fi

	cur="${words[cword]}"
	prev="${words[cword-1]}"

	# 构建命令树结构
	declare -A cmd_tree
	cmd_tree[/]="%s"
%s

	# 查找当前命令上下文
	local context="/"
	local i
	for ((i=1; i < cword; i++)); do
		local arg="${words[i]}"
		if [[ -n "${cmd_tree[$context$arg/]}" ]]; then
			context="$context$arg/"
		fi
	done

	# 获取当前上下文可用选项
	opts="${cmd_tree[$context]}"
	# 添加-o filenames选项处理特殊字符和空格
	COMPREPLY=($(compgen -o filenames -W "${opts}" -- ${cur}))

	# 模糊匹配与纠错提示：当无精确匹配时，从所有选项中查找包含关键词的项
	if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
		local all_opts=()
		# 使用循环安全收集所有选项，避免空格分割问题
		for path in "${!cmd_tree[@]}"; do
			for opt in ${cmd_tree[$path]}; do
				all_opts+=($opt)
			done
		done
		# 去重并生成补全结果
		local unique_opts=($(printf "%%s\n" "${all_opts[@]}" | sort -u))
		COMPREPLY=($(compgen -o filenames -W "${unique_opts[*]}" -- ${cur}))
	fi

	return 0
	}

complete -F _%s %s
`
	BashCommandTreeEntry = "\tcmd_tree[/%s/]=\"%s\"\n" // 命令树条目格式

	// PowerShell补全模板
	//PwshFunctionHeader   = "Register-ArgumentCompleter -CommandName %s -ScriptBlock {\n    param($wordToComplete, $commandAst, $cursorPosition, $commandName, $parameterName)\n\n    # 标志参数需求映射\n    $flagParams = @{\n%s\n    }\n\n    # 构建命令树结构\n    $cmdTree = @{\n        '' = '%s'\n%s\n    }\n\n    # 解析命令行参数获取当前上下文\n    $context = ''\n    $args = $commandAst.CommandElements | Select-Object -Skip 1 | ForEach-Object { $_.ToString() }\n    $index = 0\n    $count = $args.Count\n\n    while ($index -lt $count) {\n        $arg = $args[$index]\n        # 处理选项参数及其值\n        if ($arg -like '-*' -and $flagParams.ContainsKey($arg)) {\n            $paramType = $flagParams[$arg]\n            $index++\n            \n            # 根据参数类型决定是否跳过下一个参数\n            if ($paramType -eq 'required' -or ($paramType -eq 'optional' -and $index -lt $count -and $args[$index] -notlike '-*')) {\n                $index++\n            }\n            continue\n        }\n\n        $nextContext = if ($context) { \"$context.$arg\" } else { $arg }\n        if ($cmdTree.ContainsKey($nextContext)) {\n            $context = $nextContext\n            $index++\n        } else {\n            break\n        }\n    }\n\n    # 获取当前上下文可用选项并过滤\n    $options = @()\n    if ($cmdTree.ContainsKey($context)) {\n        $options = $cmdTree[$context] -split ' ' | Where-Object { $_ -like \"$wordToComplete*\" }\n    }\n\n    $options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_) }\n}\n" // PowerShell补全函数头部
	PwshFunctionHeader = `Register-ArgumentCompleter -CommandName %s -ScriptBlock {
		param($wordToComplete, $commandAst, $cursorPosition, $commandName, $parameterName)
	
		# 标志参数需求数组(保留原始大小写)
		$flagParams = @(
%s      )
	
		# 构建命令树结构
		$cmdTree = @{
%s      }
	
		# 解析命令行参数获取当前上下文
		$context = ''
		$args = $commandAst.CommandElements | Select-Object -Skip 1 | ForEach-Object { $_.ToString() }
		$index = 0
		$count = $args.Count
	
		while ($index -lt $count) {
			$arg = $args[$index]
			# 使用大小写敏感匹配查找标志
			$paramInfo = $flagParams | Where-Object { $_.Name -ceq $arg } | Select-Object -First 1
			if ($paramInfo) {
				$paramType = $paramInfo.Type
				$index++
				
				# 根据参数类型决定是否跳过下一个参数
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
	
		# 获取当前上下文可用选项并过滤
		$options = @()
		if ($cmdTree.ContainsKey($context)) {
			$options = $cmdTree[$context] -split ' ' | Where-Object { $_ -like "$wordToComplete*" }
		}
	
		# 模糊匹配与纠错提示：当无精确匹配时，从所有选项中查找包含关键词的项
		if (-not $options) {
			# 递归收集所有层级的选项
			$allOptions = @()
			$cmdTree.Values | ForEach-Object { $allOptions += $_ -split ' ' }
			$options = $allOptions | Select-Object -Unique | Where-Object { $_ -like "*$wordToComplete*" }
		}
	
		$options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_) }
	}`
	PwshCommandTreeEntryRoot = "\t\t'' = '%s'\n"                    // 根命令树条目格式
	PwshCommandTreeEntry     = "\t\t'%s' = '%s'\n"                  // 命令树条目格式
	PwshCommandTreeOption    = "\t\t@{ Name = '%s'; Type = '%s'}\n" // 选项参数需求条目格式
)

// addSubCommands 迭代方式添加子命令到命令树，替代递归实现
//
// 参数:
//   - cmdTreeEntries - 用于存储命令树条目的缓冲区
//   - parentPath - 父命令路径
//   - cmds - 子命令列表
func addSubCommands(cmdTreeEntries *bytes.Buffer, parentPath string, cmds []*Cmd) {
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

	// 处理队列中的所有命令
	for len(queue) > 0 {
		// 出队
		node := queue[0]
		queue = queue[1:]
		cmd := node.cmd
		currentParentPath := node.parentPath

		if cmd == nil {
			continue
		}

		// 获取子命令补全选项
		cmdOpts := cmd.collectCompletionOptions()

		// 处理长命令
		longName := cmd.LongName()
		if longName != "" {
			// 构建命令路径
			cmdLongPath := currentParentPath + longName + "/"
			trimmedPath := strings.TrimSuffix(cmdLongPath, "/")

			// 写入命令树条目
			fmt.Fprintf(cmdTreeEntries, BashCommandTreeEntry, trimmedPath, strings.Join(cmdOpts, " "))

			// 将子命令加入队列
			for _, subCmd := range cmd.subCmds {
				queue = append(queue, cmdNode{cmd: subCmd, parentPath: cmdLongPath})
			}
		}

		// 处理短命令
		shortName := cmd.ShortName()
		if shortName != "" {
			// 构建命令路径
			cmdShortPath := currentParentPath + shortName + "/"
			trimmedPath := strings.TrimSuffix(cmdShortPath, "/")

			// 写入命令树条目
			fmt.Fprintf(cmdTreeEntries, BashCommandTreeEntry, trimmedPath, strings.Join(cmdOpts, " "))

			// 将子命令加入队列
			for _, subCmd := range cmd.subCmds {
				queue = append(queue, cmdNode{cmd: subCmd, parentPath: cmdShortPath})
			}
		}
	}
}

// generateBashCompletion 生成Bash自动补全脚本
//
// 返回值：
//   - string: Bash自动补全脚本
func (c *Cmd) generateBashCompletion() (string, error) {
	// 缓冲区
	var buf bytes.Buffer

	// 检查父命令和命令树注册表是否为空
	if c.parentCmd == nil || c.flagRegistry == nil {
		return "", fmt.Errorf("invalid command state: parent command or flag registry is nil")
	}

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 获取根命令的补全选项
	rootCmdOpts := c.parentCmd.collectCompletionOptions()

	// 构建命令树条目缓冲区
	var cmdTreeEntries bytes.Buffer

	// 添加根命令选项
	rootOpts := strings.Join(rootCmdOpts, " ")

	// 从根命令的子命令开始添加
	addSubCommands(&cmdTreeEntries, "", c.parentCmd.subCmds)

	// 写入补全函数头部和命令树
	fmt.Fprintf(&buf, BashFunctionHeader, programName, rootOpts, cmdTreeEntries.String(), programName, programName)

	return buf.String(), nil
}

// generatePwshCompletion 生成PowerShell自动补全脚本
//
// 返回值：
//   - string: PowerShell自动补全脚本
func (c *Cmd) generatePwshCompletion() (string, error) {
	// 缓冲区
	var buf bytes.Buffer

	// 检查父命令和命令树注册表是否为空
	if c.parentCmd == nil || c.flagRegistry == nil {
		return "", fmt.Errorf("invalid command state: parent command or flag registry is nil")
	}

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 构建命令树条目缓冲区
	var cmdTreeEntries bytes.Buffer

	// 获取根命令的补全选项
	rootCmdOpts := c.parentCmd.collectCompletionOptions()

	// 添加根命令选项
	rootOpts := strings.Join(rootCmdOpts, " ")

	// 从根命令的子命令开始添加条目
	addSubCommandsPwsh(&cmdTreeEntries, "", c.parentCmd.subCmds, rootOpts)

	// 构建标志参数需求映射
	var flagParamsBuf bytes.Buffer

	// 收集当前命令的标志(从根命令开始)
	flagParams := c.parentCmd.collectFlagParameters() // 现在返回[]FlagParam

	// 写入标志参数需求条目 - 使用数组而非哈希表
	for _, param := range flagParams {
		fmt.Fprintf(&flagParamsBuf, PwshCommandTreeOption, param.Name, param.Type)
	}

	// 写入补全函数头部和命令树
	fmt.Fprintf(&buf, PwshFunctionHeader, programName, flagParamsBuf.String(), cmdTreeEntries.String())

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

	// 写入根命令条目
	fmt.Fprintf(cmdTreeEntries, PwshCommandTreeEntryRoot, rootOpts)

	// 处理队列中的所有命令
	for len(queue) > 0 {
		// 出队
		node := queue[0]
		queue = queue[1:]
		cmd := node.cmd
		currentParentPath := node.parentPath

		if cmd == nil {
			continue
		}

		// 获取子命令补全选项
		cmdOpts := cmd.collectCompletionOptions()

		// 处理长命令
		longName := cmd.LongName()
		if longName != "" {
			// 构建命令路径(使用.作为分隔符)
			cmdLongPath := currentParentPath
			if cmdLongPath != "" {
				cmdLongPath += fmt.Sprintf(".%s", longName)
			} else {
				cmdLongPath = longName
			}

			// 写入命令树条目
			fmt.Fprintf(cmdTreeEntries, PwshCommandTreeEntry, cmdLongPath, strings.Join(cmdOpts, " "))

			// 将子命令加入队列
			for _, subCmd := range cmd.subCmds {
				queue = append(queue, cmdNode{cmd: subCmd, parentPath: cmdLongPath})
			}
		}

		// 处理短命令
		shortName := cmd.ShortName()
		if shortName != "" {
			// 构建命令路径(使用.作为分隔符)
			cmdShortPath := currentParentPath
			if cmdShortPath != "" {
				cmdShortPath += fmt.Sprintf(".%s", shortName)
			} else {
				cmdShortPath = shortName
			}

			// 写入命令树条目
			fmt.Fprintf(cmdTreeEntries, PwshCommandTreeEntry, cmdShortPath, strings.Join(cmdOpts, " "))

			// 将子命令加入队列
			for _, subCmd := range cmd.subCmds {
				queue = append(queue, cmdNode{cmd: subCmd, parentPath: cmdShortPath})
			}
		}
	}
}

// collectFlagParameters 收集所有命令标志参数需求，返回标志名称到参数需求类型的映射
// 参数需求类型: "required"|"optional"|"none"
func (c *Cmd) collectFlagParameters() []FlagParam { // 修改返回类型为切片
	params := make([]FlagParam, 0)
	lowercaseKeys := make(map[string]bool)

	// 使用队列实现广度优先遍历替代递归
	queue := make([]*Cmd, 0, 10) // 预分配队列容量
	queue = append(queue, c)

	for len(queue) > 0 {
		// 出队
		cmd := queue[0]
		queue = queue[1:]

		// 收集当前命令的标志 - 同时处理长短选项
		for _, flag := range cmd.flagRegistry.GetAllFlagMetas() {
			// 处理短选项
			shortOpt := flag.GetShortName()
			if shortOpt != "" {
				processFlagOption("-"+shortOpt, flag, &params, lowercaseKeys)
			}

			// 处理长选项
			longOpt := flag.GetLongName()
			if longOpt != "" {
				processFlagOption("--"+longOpt, flag, &params, lowercaseKeys)
			}
		}

		// 将子命令加入队列
		queue = append(queue, cmd.subCmds...)
	}

	return params
}

// processFlagOption 处理单个标志选项并添加到参数列表
func processFlagOption(opt string, flag *flags.FlagMeta, params *[]FlagParam, lowercaseKeys map[string]bool) {
	// 根据标志类型确定参数需求
	paramType := "required" // 默认为必需参数
	if flag.GetFlagType() == flags.FlagTypeBool {
		paramType = "none" // 布尔标志不需要参数
	}

	// 检查是否已存在相同小写键
	lowerKey := strings.ToLower(opt)
	if !lowercaseKeys[lowerKey] {
		lowercaseKeys[lowerKey] = true
		*params = append(*params, FlagParam{Name: opt, Type: paramType})
	}
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

// HandleCompletionHook 自动补全钩子实现
//
// 功能: 处理自动补全子命令逻辑, 生成指定shell的补全脚本
//
// 参数:
//   - c: 当前命令实例
//
// 返回值:
//   - bool: 是否需要退出程序
//   - error: 处理过程中的错误信息
//
// 注意事项:
//   - Linux环境: 需要bash 4.0或更高版本以支持关联数组特性
//   - Windows环境: 需要PowerShell 5.1或更高版本以支持Register-ArgumentCompleter
//   - 请确保您的环境满足上述版本要求，否则自动补全功能可能无法正常工作
func HandleCompletionHook(c *Cmd) (bool, error) {
	// 检查是否启用自动补全
	if !c.enableCompletion {
		return false, nil
	}

	// 获取补全子命令
	rootCmd := c
	for rootCmd.parentCmd != nil { // 追溯到根命令
		rootCmd = rootCmd.parentCmd
	}

	s, ok := rootCmd.subCmdMap[CompletionShellLongName]
	if !ok {
		return false, nil
	}

	// 获取shell类型
	shell := s.completionShell.Get()
	if shell == ShellNone {
		return false, nil
	}

	// 生成对应shell的补全脚本
	switch shell {
	case ShellBash: // 生成Bash补全脚本
		bashCompletion, err := c.generateBashCompletion()
		if err != nil {
			return false, err
		}
		fmt.Println(bashCompletion)
	case ShellPowershell, ShellPwsh: // 兼容Powershell和Pwsh
		pwshCompletion, err := c.generatePwshCompletion()
		if err != nil {
			return false, err
		}
		fmt.Println(pwshCompletion)
	default:
		return false, fmt.Errorf("unsupported shell: %s. Supported shells are: %v", shell, ShellSlice)
	}

	// 判断是否需要退出
	if c.exitOnBuiltinFlags {
		return true, nil
	}

	return false, nil
}

// createCompletionSubcommand 创建自动补全子命令
//
// 返回值：
//   - 自动补全子命令实例
//   - 错误信息
func (c *Cmd) createCompletionSubcommand() (*Cmd, error) {
	// 创建自动补全子命令
	completionCmd := NewCmd(CompletionShellLongName, CompletionShellShortName, flag.ExitOnError)
	completionCmd.SetEnableCompletion(true) // 启用自动补全

	// 根据语言设置自动补全子命令的注意事项
	var completionNotes []string

	// 根据父命令的语言设置自动补全子命令的语言
	if c.GetUseChinese() {
		// 添加中文提示
		completionNotes = completionNotesCN
	} else { // 默认为英文提示
		completionNotes = completionNotesEN
	}

	// 添加自动补全子命令的注意事项
	for _, note := range completionNotes {
		completionCmd.AddNote(note)
	}

	// 设置自动补全子命令的内置标志退出策略
	if !c.exitOnBuiltinFlags {
		completionCmd.SetExitOnBuiltinFlags(false)
	}

	// 为子命令定义并绑定自动补全标志
	completionCmd.EnumVar(completionCmd.completionShell, CompletionShellFlagLongName, CompletionShellFlagShortName, ShellNone, fmt.Sprintf("指定要生成的shell补全脚本类型, 可选值: %v", ShellSlice), ShellSlice)

	// 根据父命令的语言设置自动设置子命令的语言
	if c.GetUseChinese() {
		completionCmd.SetDescription(CompletionShellUsageZh)
		completionCmd.SetUseChinese(true)
	} else {
		completionCmd.SetDescription(CompletionShellUsageEn)
		completionCmd.SetUseChinese(false)
	}

	// 添加自动补全子命令的示例
	cmdName := os.Args[0]

	// 遍历示例
	for _, ex := range completionExamples {
		completionCmd.AddExample(ExampleInfo{
			Description: ex.Description,
			Usage:       fmt.Sprintf(ex.Usage, cmdName),
		})
	}

	return completionCmd, nil
}
