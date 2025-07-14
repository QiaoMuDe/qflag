// pwsh_completion.go - powershell 补全脚本生成器
package cmd

import (
	"bytes"
	"fmt"
	"strings"
)

// generatePwshCommandTreeEntry 生成PowerShell命令树条目
//
// 参数:
//   - cmdTreeEntries: 命令树条目字符串
//   - cmdPath: 命令路径
//   - cmdOpts: 命令选项数组
func generatePwshCommandTreeEntry(cmdTreeEntries *bytes.Buffer, cmdPath string, cmdOpts []string) {
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

// generatePwshCompletion 生成PowerShell自动补全脚本
//
// 参数:
//   - buf: 存储生成的PowerShell自动补全脚本的缓冲区
//   - params: 标志参数数组
//   - rootCmdOpts: 根命令选项数组
//   - cmdTreeEntries: 命令树条目字符串
//   - programName: 程序名称
func generatePwshCompletion(buf *bytes.Buffer, params []FlagParam, rootCmdOpts []string, cmdTreeEntries string, programName string) {
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
	fmt.Fprintf(buf, PwshFunctionHeader, flagParamsBuf.String(), rootOptsBuilder.String(), cmdTreeEntries, programName)
}

const (
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
