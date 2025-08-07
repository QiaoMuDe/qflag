package completion

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

	// 遍历标志参数并生成相应的Bash自动补全脚本
	for _, param := range params {
		var key string
		// 如果命令路径为空，则使用参数名称作为键
		if param.CommandPath == "" {
			key = param.Name
		} else {
			key = fmt.Sprintf("%s|%s", param.CommandPath, param.Name)
		}

		// 写入标志参数项
		fmt.Fprintf(&flagParamsBuf, BashFlagParamItem, key, param.Type+"|"+param.ValueType)

		// 如果参数类型为枚举，则生成枚举选项
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			// 预分配切片容量以提高性能
			escapedOpts := make([]string, len(param.EnumOptions))
			for i, opt := range param.EnumOptions {
				escapedOpts[i] = escapeSpecialChars(opt)
			}

			// 将枚举选项转换为字符串
			options := strings.Join(escapedOpts, " ")

			// 写入枚举选项
			fmt.Fprintf(&enumOptionsBuf, BashEnumOptions, key, options)
		}
	}

	// 使用命名模板生成Bash自动补全脚本
	tmpl := strings.NewReplacer(
		"{{.RootCmdOpts}}", strings.Join(rootCmdOpts, "|"), // 根命令选项
		"{{.CmdTreeEntries}}", cmdTreeEntries, // 命令树条目
		"{{.FlagParams}}", flagParamsBuf.String(), // 标志参数
		"{{.EnumOptions}}", enumOptionsBuf.String(), // 枚举选项
		"{{.ProgramName}}", programName, // 程序名称
	)

	// 写入Bash函数头部
	_, _ = tmpl.WriteString(buf, BashFunctionHeader)
}

// escapeSpecialChars 处理字符串中的特殊字符转义
//
// 参数:
//   - s: 需要处理的字符串
//
// 返回值:
//   - 转义后的字符串
func escapeSpecialChars(s string) string {
	var builder strings.Builder
	builder.Grow(len(s) * 2) // 预分配容量以减少重新分配

	for _, r := range s {
		switch r {
		case '\\':
			builder.WriteString("\\\\")
		case '"':
			builder.WriteString("\\\"")
		case ' ':
			builder.WriteString("\\ ")
		case '$':
			builder.WriteString("\\$")
		case '`':
			builder.WriteString("\\`")
		case '|':
			builder.WriteString("\\|")
		case '&':
			builder.WriteString("\\&")
		case ';':
			builder.WriteString("\\;")
		case '(':
			builder.WriteString("\\(")
		case ')':
			builder.WriteString("\\)")
		case '<':
			builder.WriteString("\\<")
		case '>':
			builder.WriteString("\\>")
		case '*':
			builder.WriteString("\\*")
		case '?':
			builder.WriteString("\\?")
		case '[':
			builder.WriteString("\\[")
		case ']':
			builder.WriteString("\\]")
		case '{':
			builder.WriteString("\\{")
		case '}':
			builder.WriteString("\\}")
		case '~':
			builder.WriteString("\\~")
		case '#':
			builder.WriteString("\\#")
		default:
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

const (
	BashCommandTreeEntry = "cmd_tree[%s]=\"%s\"\n" // 命令树条目格式
	BashFlagParamItem    = "flag_params[%q]=%q\n"  // 标志参数项格式
	BashEnumOptions      = "enum_options[%q]=%q\n" // 枚举选项格式
)

const (
	// Bash补全模板
	BashFunctionHeader = `#!/usr/bin/env bash

# 静态命令树定义 - 在函数外预初始化
declare -A cmd_tree
cmd_tree[/]="{{.RootCmdOpts}}"
{{.CmdTreeEntries}}

# 标志参数定义 - 存储类型和值类型 (type|valueType)
declare -A flag_params
{{.FlagParams}}

# 枚举选项定义 - 存储枚举标志的允许值
declare -A enum_options
{{.EnumOptions}}

_{{.ProgramName}}() {
	local cur prev words cword context opts i arg
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

	# 查找当前命令上下文
	local context="/"
	local i
	for ((i=1; i < cword; i++)); do
		local arg="${words[i]}"
		# 如果遇到标志，停止上下文构建
		if [[ "$arg" =~ ^- ]]; then
			break
		fi
		
		local next_context="$context$arg/"
		# 验证上下文是否存在
		if [[ -n "${cmd_tree[$next_context]}" ]]; then
			context="$next_context"
		else
			break
		fi
	done

	# 获取当前上下文的可用选项
	local current_context_opts="${cmd_tree[$context]}"
	if [[ -z "$current_context_opts" ]]; then
		return 1
	fi
	
	# 安全地解析选项
	local opts_arr
	IFS='|' read -ra opts_arr <<< "$current_context_opts"
	local opts
	printf -v opts '%s ' "${opts_arr[@]}"
	opts="${opts% }" # 移除尾部空格
	
	# 检查前一个参数是否需要值并获取其类型
	local prev_param_type=""
	local prev_value_type=""
	if [[ $cword -gt 1 ]]; then
		local prev_arg="${words[cword-1]}"
		local key="${context}|${prev_arg}"
		local prev_param_info="${flag_params[$key]}"
		
		if [[ -n "$prev_param_info" ]]; then
			IFS='|' read -r prev_param_type prev_value_type <<< "$prev_param_info"
		fi
	fi

	# 根据参数类型动态生成补全
	if [[ -n "$prev_param_type" && "$prev_param_type" == "required" ]]; then
		case "$prev_value_type" in
			enum)
				local enum_key="${context}|${words[cword-1]}"
				local enum_opts="${enum_options[$enum_key]}"
				
				if [[ -n "$enum_opts" ]]; then
					# 使用内置功能进行枚举补全，避免外部命令
					local enum_arr
					read -ra enum_arr <<< "$enum_opts"
					
					if [[ -z "$cur" ]]; then
						# 当前单词为空 → 返回所有枚举值
						COMPREPLY=("${enum_arr[@]}")
					else
						# 前缀过滤（大小写敏感，性能更好）
						local matches=()
						local opt
						for opt in "${enum_arr[@]}"; do
							if [[ "$opt" == "$cur"* ]]; then
								matches+=("$opt")
							fi
						done
						COMPREPLY=("${matches[@]}")
					fi
					return 0
				fi
				;;
			string)
				# 字符串类型 - 提供文件和目录路径补全
				COMPREPLY=($(compgen -f -d -- "$cur"))
				return 0
				;;
			*)
				# 默认值补全 - 回退到标准选项
				COMPREPLY=($(compgen -W "$opts" -- "$cur"))
				return 0
				;;
		esac
	fi

	# 标志参数补全
	COMPREPLY=($(compgen -W "$opts" -- "$cur"))
	return 0
}

# 注册补全函数
complete -F _{{.ProgramName}} {{.ProgramName}}
`
)
