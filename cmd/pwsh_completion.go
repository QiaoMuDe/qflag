// pwsh_completion.go - PowerShell 补全脚本生成器
package cmd

import (
	"bytes"
	"fmt"
	"strings"
)

// generatePwshCommandTreeEntry 生成PowerShell命令树条目
//
// 参数:
// - cmdTreeEntries: 命令树条目缓冲区
// - cmdPath: 命令路径
// - cmdOpts: 命令选项
func generatePwshCommandTreeEntry(cmdTreeEntries *bytes.Buffer, cmdPath string, cmdOpts []string) {
	// PowerShell使用数组存储命令选项，避免大小写敏感的映射键
	escapedOpts := make([]string, len(cmdOpts))
	for i, opt := range cmdOpts {
		escapedOpt := strings.ReplaceAll(opt, "'", "''")
		escapedOpts[i] = escapedOpt
	}
	opts := strings.Join(escapedOpts, "', '")
	fmt.Fprintf(cmdTreeEntries, PwshCommandTreeEntry, cmdPath, opts)
}

// generatePwshCompletion 生成PowerShell自动补全脚本
//
// 参数:
// - buf: 输出缓冲区
// - params: 标志参数列表
// - rootCmdOpts: 根命令选项
// - cmdTreeEntries: 命令树条目
// - programName: 程序名称
func generatePwshCompletion(buf *bytes.Buffer, params []FlagParam, rootCmdOpts []string, cmdTreeEntries string, programName string) {
	// 构建标志参数和枚举选项（使用字符串数组而非映射）
	var flagParamsBuf bytes.Buffer
	var enumOptionsBuf bytes.Buffer

	for i, param := range params {
		key := param.CommandPath + "|" + param.Name
		flagParamsBuf.WriteString(fmt.Sprintf("$flagParams[%d] = @('%s', '%s', '%s')\n", i, key, param.Type, param.ValueType))

		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			escapedOptions := make([]string, len(param.EnumOptions))
			for j, opt := range param.EnumOptions {
				escapedOpt := strings.ReplaceAll(opt, "'", "''")
				escapedOptions[j] = escapedOpt
			}
			options := strings.Join(escapedOptions, "', '")
			enumOptionsBuf.WriteString(fmt.Sprintf("$enumOptions[%d] = @('%s', '%s')\n", i, key, options))
		}
	}

	// 根命令选项处理
	rootOpts := strings.Join(rootCmdOpts, "', '")

	// 写入PowerShell自动补全脚本
	fmt.Fprintf(buf, PwshFunctionHeader, programName, rootOpts, cmdTreeEntries, flagParamsBuf.String(), enumOptionsBuf.String())
}

const (
	// PwshCommandTreeEntry 命令树条目格式
	PwshCommandTreeEntry = "$script:cmdTree['/%s/'] = @('%s')\n"

	// PwshFunctionHeader PowerShell补全函数模板
	PwshFunctionHeader = `using namespace System.Management.Automation
using namespace System.Management.Automation.Language

# 全局命令树定义
$script:cmdTree = @{
	'/' = @('%s', %s)	
}

# 全局标志参数定义 (索引 => @(命令路径|标志名, 参数类型, 值类型))
$script:flagParams = @()
%s

# 全局枚举选项定义 (索引 => @(命令路径|标志名, 选项列表))
$script:enumOptions = @()
%s

Register-ArgumentCompleter -CommandName '%s' -ScriptBlock {
	param($wordToComplete, $commandAst, $cursorPosition)

	# 解析命令上下文
	$context = '/'
	$args = $commandAst.CommandElements | Select-Object -Skip 1 | ForEach-Object { $_.ToString() }

	for ($i = 0; $i -lt $args.Count; $i++) {
		$arg = $args[$i]
		$potentialContext = $context + $arg + '/'
		if ($script:cmdTree.ContainsKey($potentialContext)) {
			$context = $potentialContext
		}
	}

	# 获取当前上下文选项
	$currentOpts = $script:cmdTree[$context]

	# 处理标志值补全
	$prevArg = if ($args.Count -gt 0) { $args[-1] } else { $null }
	$completions = @()

	if ($prevArg -match '^-') {
		# 查找匹配的标志参数
		$targetKey = "$context|$prevArg"
		$paramInfo = $script:flagParams | Where-Object { $_[0] -eq $targetKey } | Select-Object -First 1

		if ($paramInfo) {
			$paramType, $valueType = $paramInfo[1], $paramInfo[2]

			switch ($valueType) {
				'path' {
					$completions = @(Get-ChildItem -Path "$wordToComplete*" | ForEach-Object { $_.FullName })
					break
				}
				'number' {
					$completions = 1..10 | ForEach-Object { "$_" } | Where-Object { $_ -like "$wordToComplete*" }
					break
				}
				'ip' {
					$completions = @('192.168.', '10.0.', '172.16.') | Where-Object { $_ -like "$wordToComplete*" }
					break
				}
				'enum' {
					$enumEntry = $script:enumOptions | Where-Object { $_[0] -eq $targetKey } | Select-Object -First 1
					if ($enumEntry) {
						$completions = $enumEntry[1..($enumEntry.Length-1)] | Where-Object { $_ -like "$wordToComplete*" }
					}
					break
				}
				'url' {
					$completions = @('http://', 'https://', 'ftp://') | Where-Object { $_ -like "$wordToComplete*" }
					break
				}
			}
		}
	}

	# 处理命令和标志补全
	if (-not $completions) {
		if ($wordToComplete -match '^-') {
			$completions = $currentOpts | Where-Object { $_ -like "$wordToComplete*" }
		} else {
			$completions = $currentOpts | Where-Object { $_ -like "$wordToComplete*" }
			# 添加文件系统补全
			$fileCompletions = @(Get-ChildItem -Path "$wordToComplete*" | ForEach-Object { $_.Name })
			$completions = $completions + $fileCompletions | Select-Object -Unique
		}
	}

	# 返回补全结果
	$completions | ForEach-Object {
		[System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_)
	}
}
`
)
