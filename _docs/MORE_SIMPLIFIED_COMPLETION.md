# 更简化的补全脚本实现方案

## 概述

进一步简化实现, 去掉特殊字符转义等不必要的复杂逻辑, 只保留核心功能。

## 简化后的类型定义

只保留必要的类型定义: 

```go
package types

// FlagParam 表示标志参数及其需求类型和值类型
// 直接用于补全脚本生成, 避免类型转换
type FlagParam struct {
	CommandPath string   // 命令路径, 如 "/cmd/subcmd"
	Name        string   // 标志名称(保留原始大小写)
	Type        string   // 参数需求类型: "required"|"optional"|"none"
	ValueType   string   // 参数值类型: "path"|"string"|"number"|"enum"|"bool"等
	EnumOptions []string // 枚举类型的可选值列表
}

// CompletionGenerator 补全生成器接口
//
// CompletionGenerator 定义了补全脚本生成器的核心行为。
// 所有补全生成器都应该实现此接口。
type CompletionGenerator interface {
	// Generate 生成补全脚本
	//
	// 参数:
	//   - cmd: 要生成补全脚本的命令
	//   - shellType: Shell类型 (bash, pwsh) 
	//
	// 返回值:
	//   - string: 生成的补全脚本
	//   - error: 生成失败时返回错误
	Generate(cmd Command, shellType string) (string, error)
}
```

## 简化后的实现

### 1. 修改 completion.go

直接从 Command 收集 FlagParam, 不经过中间类型: 

```go
package completion

import (
	"bytes"
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// DefaultCompletionGenerator 默认补全生成器
type DefaultCompletionGenerator struct{}

// NewDefaultCompletionGenerator 创建默认补全生成器
func NewDefaultCompletionGenerator() *DefaultCompletionGenerator {
	return &DefaultCompletionGenerator{}
}

// Generate 生成补全脚本
func (g *DefaultCompletionGenerator) Generate(cmd types.Command, shellType string) (string, error) {
	// 直接从命令收集标志参数
	params := g.collectFlagParams(cmd, "")
	
	// 收集命令选项
	rootCmdOpts := g.collectCommandOptions(cmd)
	
	// 构建命令树条目
	var cmdTreeEntries strings.Builder
	g.traverseCommandTree(&cmdTreeEntries, cmd, "", rootCmdOpts)
	
	// 程序名称
	programName := cmd.Name()
	
	// 根据shell类型生成脚本
	switch shellType {
	case "bash":
		return g.generateBashScript(params, rootCmdOpts, cmdTreeEntries.String(), programName)
	case "pwsh":
		return g.generatePwshScript(params, rootCmdOpts, cmdTreeEntries.String(), programName)
	default:
		return "", fmt.Errorf("unsupported shell type: %s", shellType)
	}
}

// collectFlagParams 收集标志参数
func (g *DefaultCompletionGenerator) collectFlagParams(cmd types.Command, commandPath string) []types.FlagParam {
	var params []types.FlagParam
	
	// 收集当前命令的标志
	flags := cmd.Flags()
	for _, flag := range flags {
		param := types.FlagParam{
			CommandPath: commandPath,
			Name:        g.buildFlagName(flag),
			Type:        g.getFlagType(flag.Type()),
			ValueType:   g.getFlagValueType(flag.Type()),
		}
		
		// 如果是枚举类型, 收集枚举值
		if flag.Type() == types.FlagTypeEnum {
			// 使用新添加的 EnumValues 方法获取枚举值
			param.EnumOptions = flag.EnumValues()
		}
		
		params = append(params, param)
	}
	
	// 递归收集子命令的标志
	subCmds := cmd.SubCmds()
	for _, subCmd := range subCmds {
		var subCommandPath string
		if commandPath == "" {
			subCommandPath = "/" + subCmd.Name()
		} else {
			subCommandPath = commandPath + "/" + subCmd.Name()
		}
		
		// 如果有短名称, 添加到路径中
		if subCmd.ShortName() != "" {
			subCommandPath += "|" + subCmd.ShortName()
		}
		
		subParams := g.collectFlagParams(subCmd, subCommandPath)
		params = append(params, subParams...)
	}
	
	return params
}

// buildFlagName 构建标志名称
func (g *DefaultCompletionGenerator) buildFlagName(flag types.Flag) string {
	longName := flag.Name()
	shortName := flag.ShortName()
	
	if longName != "" && shortName != "" {
		return longName + "|" + shortName
	} else if longName != "" {
		return longName
	} else {
		return shortName
	}
}

// getFlagType 根据标志类型获取参数需求类型
func (g *DefaultCompletionGenerator) getFlagType(flagType types.FlagType) string {
	switch flagType {
	case types.FlagTypeBool:
		return "none" // 布尔标志不需要值
	default:
		return "required" // 其他类型需要值
	}
}

// getFlagValueType 根据标志类型获取参数值类型
func (g *DefaultCompletionGenerator) getFlagValueType(flagType types.FlagType) string {
	switch flagType {
	case types.FlagTypeString:
		return "string"
	case types.FlagTypeInt:
		return "number"
	case types.FlagTypeFloat64:
		return "number"
	case types.FlagTypeBool:
		return "bool"
	case types.FlagTypeEnum:
		return "enum"
	case types.FlagTypeDuration:
		return "string"
	default:
		return "string"
	}
}

// collectCommandOptions 收集命令选项
func (g *DefaultCompletionGenerator) collectCommandOptions(cmd types.Command) []string {
	var options []string
	
	// 添加标志选项
	flags := cmd.Flags()
	for _, flag := range flags {
		longName := flag.Name()
		shortName := flag.ShortName()
		
		if longName != "" {
			options = append(options, "--"+longName)
		}
		if shortName != "" {
			options = append(options, "-"+shortName)
		}
	}
	
	// 添加子命令选项
	subCmds := cmd.SubCmds()
	for _, subCmd := range subCmds {
		longName := subCmd.Name()
		shortName := subCmd.ShortName()
		
		if longName != "" {
			options = append(options, longName)
		}
		if shortName != "" {
			options = append(options, shortName)
		}
	}
	
	return options
}

// traverseCommandTree 遍历命令树并生成命令树条目
func (g *DefaultCompletionGenerator) traverseCommandTree(cmdTreeEntries *strings.Builder, cmd types.Command, commandPath string, rootCmdOpts []string) {
	// 添加根命令条目
	if commandPath == "" {
		cmdTreeEntries.WriteString(cmd.Name())
		cmdTreeEntries.WriteString("_cmd_tree[\"/\"]=\"")
		cmdTreeEntries.WriteString(strings.Join(rootCmdOpts, "|"))
		cmdTreeEntries.WriteString("\"\n")
	}
	
	// 处理子命令
	subCmds := cmd.SubCmds()
	for _, subCmd := range subCmds {
		var subCommandPath string
		if commandPath == "" {
			subCommandPath = "/" + subCmd.Name()
		} else {
			subCommandPath = commandPath + "/" + subCmd.Name()
		}
		
		// 如果有短名称, 添加到路径中
		if subCmd.ShortName() != "" {
			subCommandPath += "|" + subCmd.ShortName()
		}
		
		// 收集子命令选项
		subCmdOpts := g.collectCommandOptions(subCmd)
		
		// 生成命令树条目
		cmdTreeEntries.WriteString(cmd.Name())
		cmdTreeEntries.WriteString("_cmd_tree[\"")
		cmdTreeEntries.WriteString(subCommandPath)
		cmdTreeEntries.WriteString("\"]=\"")
		cmdTreeEntries.WriteString(strings.Join(subCmdOpts, "|"))
		cmdTreeEntries.WriteString("\"\n")
		
		// 递归处理子命令的子命令
		g.traverseCommandTree(cmdTreeEntries, subCmd, subCommandPath, subCmdOpts)
	}
}
```

### 2. 简化 bash.go

直接使用 FlagParam 生成 Bash 脚本, 去掉特殊字符转义: 

```go
package completion

import (
	"bytes"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// generateBashScript 生成Bash补全脚本
func (g *DefaultCompletionGenerator) generateBashScript(params []types.FlagParam, rootCmdOpts []string, cmdTreeEntries string, programName string) (string, error) {
	var buf bytes.Buffer

	// 构建标志参数映射
	var flagParamsBuf strings.Builder
	var enumOptionsBuf strings.Builder

	// 遍历标志参数并生成相应的Bash自动补全脚本
	for _, param := range params {
		var key string
		if param.CommandPath == "" {
			key = param.Name
		} else {
			key = param.CommandPath + "|" + param.Name
		}

		// 构建参数值字符串
		paramValue := param.Type + "|" + param.ValueType

		// 写入标志参数项
		flagParamsBuf.WriteString(programName)
		flagParamsBuf.WriteString("_flag_params[\"")
		flagParamsBuf.WriteString(key)
		flagParamsBuf.WriteString("\"]=\"")
		flagParamsBuf.WriteString(paramValue)
		flagParamsBuf.WriteString("\"\n")

		// 如果参数类型为枚举, 则生成枚举选项
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			// 将枚举选项转换为字符串, 使用|分隔符与其他选项保持一致
			options := strings.Join(param.EnumOptions, "|")

			// 写入枚举选项
			enumOptionsBuf.WriteString(programName)
			enumOptionsBuf.WriteString("_enum_options[\"")
			enumOptionsBuf.WriteString(key)
			enumOptionsBuf.WriteString("\"]=\"")
			enumOptionsBuf.WriteString(options)
			enumOptionsBuf.WriteString("\"\n")
		}
	}

	// 使用字符串替换生成Bash自动补全脚本
	tmpl := strings.NewReplacer(
		"{{.RootCmdOpts}}", strings.Join(rootCmdOpts, "|"), // 根命令选项
		"{{.CmdTreeEntries}}", cmdTreeEntries, // 命令树条目
		"{{.FlagParams}}", flagParamsBuf.String(), // 标志参数
		"{{.EnumOptions}}", enumOptionsBuf.String(), // 枚举选项
		"{{.ProgramName}}", programName, // 程序名称
	)

	// 写入Bash函数头部
	_, _ = tmpl.WriteString(&buf, BashFunctionHeader)

	return buf.String(), nil
}

// BashFunctionHeader Bash补全脚本模板
const BashFunctionHeader = `#!/usr/bin/env bash

# ==================== 静态数据定义 ====================
# 静态命令树定义 - 在函数外预初始化
declare -A {{.ProgramName}}_cmd_tree
{{.ProgramName}}_cmd_tree["/"]="{{.RootCmdOpts}}"
{{.CmdTreeEntries}}

# 标志参数定义 - 存储类型和值类型 (type|valueType)
declare -A {{.ProgramName}}_flag_params
{{.FlagParams}}

# 枚举选项定义 - 存储枚举标志的允许值
declare -A {{.ProgramName}}_enum_options
{{.EnumOptions}}

# ==================== 主补全函数 ====================
_{{.ProgramName}}() {
	local cur prev words cword
	COMPREPLY=()

	# 使用 _get_comp_words_by_ref 获取补全参数以提高健壮性
	if declare -F _get_comp_words_by_ref >/dev/null 2>&1; then
		_get_comp_words_by_ref -n =: cur prev words cword
	else
		# 与旧版本 Bash 补全环境的兼容性
		words=("${COMP_WORDS[@]}")
		cword=$COMP_CWORD
		cur="${words[cword]}"
		prev="${words[cword-1]}"
	fi

	# 输入验证
	if [[ $cword -lt 0 || ${#words[@]} -eq 0 ]]; then
		return 1
	fi

	# 快速路径: 如果当前输入看起来像是路径, 优先提供路径补全
	if [[ "$cur" == *"/"* || "$cur" == *"."* || "$cur" == *"~"* ]]; then
		COMPREPLY=($(compgen -f -d -- "$cur"))
		return 0
	fi

	# 查找当前命令上下文
	local context="/"
	local i
	for ((i=1; i < cword; i++)); do
		local arg="${words[i]}"
		# 如果遇到标志, 停止上下文构建
		if [[ "$arg" =~ ^- ]]; then
			break
		fi
		
		local next_context="$context$arg/"
		# 验证上下文是否存在
		if [[ -n "${{{.ProgramName}}_cmd_tree[$next_context]}" ]]; then
			context="$next_context"
		else
			break
		fi
	done

	# 获取当前上下文的可用选项
	local current_context_opts="${{{.ProgramName}}_cmd_tree[$context]}"
	if [[ -z "$current_context_opts" ]]; then
		return 1
	fi
	
	# 检查前一个参数是否需要值并获取其类型
	local prev_param_type=""
	local prev_value_type=""
	if [[ $cword -gt 1 ]]; then
		local prev_arg="${words[cword-1]}"
		local key="${context}|${prev_arg}"
		local prev_param_info="${{{.ProgramName}}_flag_params[$key]:-}"
		
		if [[ -n "$prev_param_info" ]]; then
			IFS='|' read -r prev_param_type prev_value_type <<< "$prev_param_info"
		fi
	fi

	# 根据参数类型动态生成补全 - 仅处理需要值的参数
	if [[ -n "$prev_param_type" && "$prev_param_type" == "required" ]]; then
		if [[ "$prev_value_type" == "enum" ]]; then
			local enum_key="${context}|${words[cword-1]}"
			local enum_opts="${{{.ProgramName}}_enum_options[$enum_key]:-}"
			
			if [[ -n "$enum_opts" ]]; then
				# 对枚举选项使用compgen
				local opts
				printf -v opts '%s ' "${enum_opts//|/ }"
				COMPREPLY=($(compgen -W "$opts" -- "$cur"))
				return 0
			fi
		else
			# 非枚举类型 - 统一使用文件和目录路径补全
			COMPREPLY=($(compgen -f -d -- "$cur"))
			return 0
		fi
	fi

	# 主要的标志和命令补全 - 使用compgen
	local opts
	printf -v opts '%s ' "${current_context_opts//|/ }"
	COMPREPLY=($(compgen -W "$opts" -- "$cur"))
	return 0
}

# 注册补全函数
complete -F _{{.ProgramName}} {{.ProgramName}}`
```

### 3. 简化 pwsh.go

类似地实现 PowerShell 补全脚本生成: 

```go
package completion

import (
	"bytes"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// generatePwshScript 生成PowerShell补全脚本
func (g *DefaultCompletionGenerator) generatePwshScript(params []types.FlagParam, rootCmdOpts []string, cmdTreeEntries string, programName string) (string, error) {
	var buf bytes.Buffer

	// 构建标志参数映射
	var flagParamsBuf strings.Builder
	var enumOptionsBuf strings.Builder

	// 遍历标志参数并生成相应的PowerShell自动补全脚本
	for _, param := range params {
		var key string
		if param.CommandPath == "" {
			key = param.Name
		} else {
			key = param.CommandPath + "|" + param.Name
		}

		// 构建参数值字符串
		paramValue := param.Type + "|" + param.ValueType

		// 写入标志参数项
		flagParamsBuf.WriteString("$")
		flagParamsBuf.WriteString(programName)
		flagParamsBuf.WriteString(":FlagParams[\"")
		flagParamsBuf.WriteString(key)
		flagParamsBuf.WriteString("\"]=\"")
		flagParamsBuf.WriteString(paramValue)
		flagParamsBuf.WriteString("\"\n")

		// 如果参数类型为枚举, 则生成枚举选项
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			// 将枚举选项转换为字符串, 使用|分隔符与其他选项保持一致
			options := strings.Join(param.EnumOptions, "|")

			// 写入枚举选项
			enumOptionsBuf.WriteString("$")
			enumOptionsBuf.WriteString(programName)
			enumOptionsBuf.WriteString(":EnumOptions[\"")
			enumOptionsBuf.WriteString(key)
			enumOptionsBuf.WriteString("\"]=\"")
			enumOptionsBuf.WriteString(options)
			enumOptionsBuf.WriteString("\"\n")
		}
	}

	// 使用字符串替换生成PowerShell自动补全脚本
	tmpl := strings.NewReplacer(
		"{{.RootCmdOpts}}", strings.Join(rootCmdOpts, "|"), // 根命令选项
		"{{.CmdTreeEntries}}", cmdTreeEntries, // 命令树条目
		"{{.FlagParams}}", flagParamsBuf.String(), // 标志参数
		"{{.EnumOptions}}", enumOptionsBuf.String(), // 枚举选项
		"{{.ProgramName}}", programName, // 程序名称
	)

	// 写入PowerShell函数头部
	_, _ = tmpl.WriteString(&buf, PwshFunctionHeader)

	return buf.String(), nil
}

// PwshFunctionHeader PowerShell补全脚本模板
const PwshFunctionHeader = `# PowerShell 补全脚本 for {{.ProgramName}}

# ==================== 静态数据定义 ====================
# 静态命令树定义
${{.ProgramName}}:CmdTree = @{}
${{.ProgramName}}:CmdTree["/"] = "{{.RootCmdOpts}}"
{{.CmdTreeEntries}}

# 标志参数定义 - 存储类型和值类型 (type|valueType)
${{.ProgramName}}:FlagParams = @{}
{{.FlagParams}}

# 枚举选项定义 - 存储枚举标志的允许值
${{.ProgramName}}:EnumOptions = @{}
{{.EnumOptions}}

# ==================== 补全函数 ====================
function _{{.ProgramName}} {
    param(
        [string]$commandName,
        [string]$parameterName,
        [string]$wordToComplete,
        [System.Management.Automation.CommandAst]$commandAst,
        [System.Collections.Hashtable]$parameters
    )

    # 获取当前命令上下文
    $context = "/"
    $currentIndex = 0
    for ($i = 1; $i -lt $commandAst.CommandElements.Count; $i++) {
        $element = $commandAst.CommandElements[$i]
        if ($element -is [System.Management.Automation.CommandParameterAst]) {
            # 遇到标志, 停止上下文构建
            break
        }
        
        $nextContext = $context + $element.Value + "/"
        if (${{.ProgramName}}:CmdTree.ContainsKey($nextContext)) {
            $context = $nextContext
            $currentIndex = $i
        } else {
            break
        }
    }

    # 获取当前上下文的可用选项
    $currentContextOpts = ${{.ProgramName}}:CmdTree[$context]
    if (-not $currentContextOpts) {
        return
    }

    # 检查前一个参数是否需要值并获取其类型
    $prevParamType = ""
    $prevValueType = ""
    if ($currentIndex -gt 0) {
        $prevElement = $commandAst.CommandElements[$currentIndex - 1]
        if ($prevElement -is [System.Management.Automation.CommandParameterAst]) {
            $prevParamName = $prevElement.ParameterName
            $key = $context + "|" + $prevParamName
            $prevParamInfo = ${{.ProgramName}}:FlagParams[$key]
            
            if ($prevParamInfo) {
                $parts = $prevParamInfo -split '\|'
                $prevParamType = $parts[0]
                $prevValueType = $parts[1]
            }
        }
    }

    # 根据参数类型动态生成补全
    if ($prevParamType -eq "required") {
        if ($prevValueType -eq "enum") {
            $enumKey = $context + "|" + $commandAst.CommandElements[$currentIndex - 1].ParameterName
            $enumOpts = ${{.ProgramName}}:EnumOptions[$enumKey]
            
            if ($enumOpts) {
                $options = $enumOpts -split '\|'
                $options | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object { 
                    [System.Management.Automation.CompletionResult]::new($_)
                }
                return
            }
        } else {
            # 非枚举类型 - 提供文件和目录路径补全
            return (Get-ChildItem).Name | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object { 
                [System.Management.Automation.CompletionResult]::new($_)
            }
        }
    }

    # 主要的标志和命令补全
    $options = $currentContextOpts -split '\|'
    $options | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object { 
        [System.Management.Automation.CompletionResult]::new($_)
    }
}

# 注册补全函数
Register-ArgumentCompleter -Native -CommandName {{.ProgramName}} -ScriptBlock $_{{.ProgramName}}`
```

## 总结

进一步简化后的方案: 

1. **去掉特殊字符转义**: 不再处理特殊字符转义, 简化代码
2. **简化补全逻辑**: 使用标准的 `compgen` 和 `Where-Object` 进行补全, 而不是复杂的模糊匹配
3. **只保留核心功能**: 只保留基本的补全功能, 去掉高级特性如模糊匹配、缓存等
4. **直接使用 FlagParam**: 避免类型转换, 直接从 Command 收集 FlagParam

这样代码更简洁, 逻辑更清晰, 同时满足基本的补全需求。