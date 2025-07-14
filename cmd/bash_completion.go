// bash_completion.go - bash 补全脚本生成器
package cmd

import (
	"bytes"
	"fmt"
	"strings"
)

// generateBashCommandTreeEntry 生成Bash命令树条目
//
// 参数:
// - cmdTreeEntries: 命令树条目缓冲区
// - cmdPath: 命令路径
// - cmdOpts: 命令选项
func generateBashCommandTreeEntry(cmdTreeEntries *bytes.Buffer, cmdPath string, cmdOpts []string) {
	fmt.Fprintf(cmdTreeEntries, BashCommandTreeEntry, cmdPath, strings.Join(cmdOpts, "|"))
}

// generateBashCompletion 生成Bash自动补全脚本
//
// 参数:
// - buf: 输出缓冲区
// - params: 标志参数列表
// - rootCmdOpts: 根命令选项
// - cmdTreeEntries: 命令树条目
// - programName: 程序名称
func generateBashCompletion(buf *bytes.Buffer, params []FlagParam, rootCmdOpts []string, cmdTreeEntries string, programName string) {
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
	fmt.Fprintf(buf, BashFunctionHeader, strings.Join(rootCmdOpts, "|"), cmdTreeEntries, flagParamsBuf.String(), enumOptionsBuf.String(), programName, programName, programName)
}

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
)
