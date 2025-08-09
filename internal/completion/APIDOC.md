package completion // import "gitee.com/MM-Q/qflag/internal/completion"

Package completion Bash Shell 自动补全实现 本文件实现了Bash Shell环境下的命令行自动补全功能，
生成Bash补全脚本，支持标志和子命令的智能补全。

Package completion 命令行自动补全功能 本文件实现了命令行自动补全的核心功能，包括标志补全、子命令补全、
参数值补全等，为用户提供便捷的命令行交互体验。

Package completion 自动补全内部实现 本文件包含了自动补全功能的内部实现逻辑，提供补全算法、 匹配策略等核心功能的底层支持。

Package completion PowerShell 自动补全实现 本文件实现了PowerShell环境下的命令行自动补全功能，
生成PowerShell补全脚本，支持标志和子命令的智能补全。

CONSTANTS

const (
	BashFlagParamItem = "{{.ProgramName}}_flag_params[%q]=%q\n"  // 标志参数项格式
	BashEnumOptions   = "{{.ProgramName}}_enum_options[%q]=%q\n" // 枚举选项格式
)
const (
	// DefaultFlagParamsCapacity 预估的标志参数初始容量
	// 基于常见CLI工具分析，大多数工具的标志数量在100-500之间
	DefaultFlagParamsCapacity = 256

	// NamesPerItem 每个标志/命令的名称数量(长名+短名)
	NamesPerItem = 2

	// MaxTraverseDepth 命令树遍历的最大深度限制
	// 防止循环引用导致的无限递归，一般CLI工具很少超过20层
	MaxTraverseDepth = 50
)
    补全脚本生成相关常量定义

const (
	// 标志参数条目(含枚举选项)
	PwshFlagParamItem = "	@{ Context = \"{{.Context}}\"; Parameter = \"{{.Parameter}}\"; ParamType = \"{{.ParamType}}\"; ValueType = \"{{.ValueType}}\"; Options = @({{.Options}}) }"
	// 命令树条目
	PwshCmdTreeItem = "	@{ Context = \"{{.Context}}\"; Options = @({{.Options}}) }"
)
const (
	// 优化的Bash补全模板 - 集成高性能模糊匹配功能
	BashFunctionHeader = `#!/usr/bin/env bash

# ==================== 模糊补全配置参数 ====================
# 模糊补全功能开关 (设置为0禁用，1启用)
readonly {{.ProgramName}}_FUZZY_COMPLETION_ENABLED=1

# 启用模糊补全的最大候选项数量阈值
# 超过此数量将回退到传统前缀匹配以保证性能
readonly {{.ProgramName}}_FUZZY_MAX_CANDIDATES=150

# 模糊匹配的最小输入长度 (小于此长度不启用模糊匹配)
readonly {{.ProgramName}}_FUZZY_MIN_PATTERN_LENGTH=2

# 模糊匹配分数阈值 (0-100，分数低于此值的匹配将被过滤)
readonly {{.ProgramName}}_FUZZY_SCORE_THRESHOLD=30

# 模糊匹配最大返回结果数
readonly {{.ProgramName}}_FUZZY_MAX_RESULTS=8

# 缓存大小控制参数
# 缓存条目数量超过此阈值时将清空缓存以防止内存无限增长
readonly {{.ProgramName}}_FUZZY_CACHE_MAX_SIZE=500

# ==================== 静态数据定义 ====================
# 静态命令树定义 - 在函数外预初始化
declare -A {{.ProgramName}}_cmd_tree
{{.ProgramName}}_cmd_tree[/]="{{.RootCmdOpts}}"
{{.CmdTreeEntries}}

# 标志参数定义 - 存储类型和值类型 (type|valueType)
declare -A {{.ProgramName}}_flag_params
{{.FlagParams}}

# 枚举选项定义 - 存储枚举标志的允许值
declare -A {{.ProgramName}}_enum_options
{{.EnumOptions}}

# 模糊匹配结果缓存 (格式: "pattern|candidate" -> score)
declare -A {{.ProgramName}}_fuzzy_cache

# ==================== 模糊匹配核心算法 ====================
# 高性能模糊评分函数 - 使用纯整数运算避免bc开销
# 参数: $1=输入模式, $2=候选字符串
# 返回: 0-100的整数分数 (通过echo输出)
_{{.ProgramName}}_fuzzy_score_fast() {
    local pattern="$1"
    local candidate="$2"
    local pattern_len=${#pattern}
    local candidate_len=${#candidate}
    
    # 性能优化1: 长度预检查 - 候选项太短直接返回0
    if [[ $candidate_len -lt $pattern_len ]]; then
        echo "0"
        return
    fi
    
    # 性能优化2: 完全匹配检查 - 避免不必要的复杂计算
    if [[ "$candidate" == "$pattern"* ]]; then
        echo "100"  # 前缀完全匹配给最高分
        return
    fi
    
    # 性能优化3: 字符存在性预检查 - 快速排除不可能的匹配
    local pattern_lower="${pattern,,}"  # 转小写用于大小写不敏感匹配
    local candidate_lower="${candidate,,}"
    local i
    for ((i=0; i<pattern_len; i++)); do
        local char="${pattern_lower:$i:1}"
        if [[ "$candidate_lower" != *"$char"* ]]; then
            echo "0"  # 必需字符不存在，直接返回
            return
        fi
    done
    
    # 核心匹配算法 - 计算字符匹配度和连续性
    local matched=0           # 匹配的字符数
    local consecutive=0       # 当前连续匹配长度
    local max_consecutive=0   # 最大连续匹配长度
    local candidate_pos=0     # 候选字符串当前搜索位置
    local start_bonus=0       # 起始位置奖励
    
    # 检查是否从开头匹配 (大小写不敏感)
    if [[ "$candidate_lower" == "$pattern_lower"* ]]; then
        start_bonus=20  # 起始匹配给20分奖励
    fi
    
    # 逐字符匹配算法
    for ((i=0; i<pattern_len; i++)); do
        local pattern_char="${pattern_lower:$i:1}"
        local found=0
        
        # 在候选字符串中查找当前模式字符
        for ((j=candidate_pos; j<candidate_len; j++)); do
            local candidate_char="${candidate_lower:$j:1}"
            if [[ "$pattern_char" == "$candidate_char" ]]; then
                ((matched++))
                found=1
                
                # 连续性检查 - 连续匹配的字符得分更高
                if [[ $j -eq $candidate_pos ]]; then
                    ((consecutive++))
                    if [[ $consecutive -gt $max_consecutive ]]; then
                        max_consecutive=$consecutive
                    fi
                else
                    consecutive=1  # 重置连续计数
                fi
                
                candidate_pos=$((j+1))  # 更新搜索位置
                break
            fi
        done
        
        # 如果某个字符未找到，重置连续计数
        if [[ $found -eq 0 ]]; then
            consecutive=0
        fi
    done
    
    # 评分计算 - 使用整数运算避免浮点数和bc
    # 基础分数: (匹配字符数 / 模式长度) * 60
    local base_score=$((matched * 60 / pattern_len))
    
    # 连续性奖励: (最大连续长度 / 模式长度) * 20
    local consecutive_bonus=$((max_consecutive * 20 / pattern_len))
    
    # 长度惩罚: 候选字符串越长，分数略微降低
    local length_penalty=$((candidate_len - pattern_len))
    if [[ $length_penalty -gt 10 ]]; then
        length_penalty=10  # 最大惩罚10分
    fi
    
    # 最终分数计算
    local final_score=$((base_score + consecutive_bonus + start_bonus - length_penalty))
    
    # 确保分数在0-100范围内
    if [[ $final_score -lt 0 ]]; then
        final_score=0
    elif [[ $final_score -gt 100 ]]; then
        final_score=100
    fi
    
    echo "$final_score"
}

# 带缓存的模糊评分函数 - 避免重复计算提高性能
# 参数: $1=输入模式, $2=候选字符串
_{{.ProgramName}}_fuzzy_score_cached() {
    local pattern="$1"
    local candidate="$2"
    local cache_key="${pattern}|${candidate}"
    
    # 缓存命中检查
    if [[ -n "${{.ProgramName}}_fuzzy_cache[$cache_key]" ]]; then
        echo "${{.ProgramName}}_fuzzy_cache[$cache_key]"
        return
    fi
    
    # 计算分数并缓存
    local score
    score=$(_{{.ProgramName}}_fuzzy_score_fast "$pattern" "$candidate")
    
    # 缓存大小控制 - 防止内存无限增长
    if [[ ${#{{.ProgramName}}_fuzzy_cache[@]} -gt ${{.ProgramName}}_FUZZY_CACHE_MAX_SIZE ]]; then
        {{.ProgramName}}_fuzzy_cache=()  # 清空缓存
    fi
    
    {{.ProgramName}}_fuzzy_cache[$cache_key]="$score"
    echo "$score"
}

# 智能补全匹配函数 - 分级匹配策略
# 参数: $1=输入模式, $2=候选选项字符串(用|分隔)
_{{.ProgramName}}_intelligent_match() {
    local pattern="$1"
    local options_str="$2"
    local pattern_len=${#pattern}
    
    # 解析候选选项到数组
    local opts_arr
    IFS='|' read -ra opts_arr <<< "$options_str"
    local total_candidates=${#opts_arr[@]}
    
    # 性能保护: 候选项过多时禁用模糊匹配
    if [[ $total_candidates -gt ${{.ProgramName}}_FUZZY_MAX_CANDIDATES ]]; then
        # 回退到传统compgen前缀匹配
        local opts
        printf -v opts '%s ' "${opts_arr[@]}"
        opts="${opts% }"
        COMPREPLY=($(compgen -W "$opts" -- "$pattern"))
        return 0
    fi
    
    local matches=()
    
    # 第1级: 精确前缀匹配 (最快，优先级最高)
    local exact_matches=()
    local opt
    for opt in "${opts_arr[@]}"; do
        if [[ "$opt" == "$pattern"* ]]; then
            exact_matches+=("$opt")
        fi
    done
    
    # 如果有精确匹配且数量合理，直接返回
    if [[ ${#exact_matches[@]} -gt 0 && ${#exact_matches[@]} -le 15 ]]; then
        COMPREPLY=("${exact_matches[@]}")
        return 0
    fi
    
    # 第2级: 大小写不敏感前缀匹配
    if [[ ${#exact_matches[@]} -eq 0 ]]; then
        local pattern_lower="${pattern,,}"
        local case_insensitive_matches=()
        
        for opt in "${opts_arr[@]}"; do
            local opt_lower="${opt,,}"
            if [[ "$opt_lower" == "$pattern_lower"* ]]; then
                case_insensitive_matches+=("$opt")
            fi
        done
        
        # 如果有大小写不敏感匹配，返回
        if [[ ${#case_insensitive_matches[@]} -gt 0 ]]; then
            COMPREPLY=("${case_insensitive_matches[@]}")
            return 0
        fi
    fi
    
    # 第3级: 模糊匹配 (最慢，仅在必要时使用)
    if [[ ${{.ProgramName}}_FUZZY_COMPLETION_ENABLED -eq 1 && $pattern_len -ge ${{.ProgramName}}_FUZZY_MIN_PATTERN_LENGTH ]]; then
        local fuzzy_results=()
        local scored_matches=()
        
        # 对所有候选项进行模糊评分
        for opt in "${opts_arr[@]}"; do
            local score
            score=$(_{{.ProgramName}}_fuzzy_score_cached "$pattern" "$opt")
            
            # 只保留分数达到阈值的匹配
            if [[ $score -ge ${{.ProgramName}}_FUZZY_SCORE_THRESHOLD ]]; then
                scored_matches+=("$score:$opt")
            fi
        done
        
        # 如果有模糊匹配结果，按分数排序并返回前N个
        if [[ ${#scored_matches[@]} -gt 0 ]]; then
            # 按分数降序排序
            local sorted_matches
            mapfile -t sorted_matches < <(printf '%s\n' "${scored_matches[@]}" | sort -t: -k1 -nr)
            
            # 提取选项名称，限制返回数量
            local count=0
            local match
            for match in "${sorted_matches[@]}"; do
                if [[ $count -ge ${{.ProgramName}}_FUZZY_MAX_RESULTS ]]; then
                    break
                fi
                fuzzy_results+=("${match#*:}")  # 移除分数前缀
                ((count++))
            done
            
            COMPREPLY=("${fuzzy_results[@]}")
            return 0
        fi
    fi
    
    # 第4级: 子字符串匹配 (最后的备选方案)
    local substring_matches=()
    local pattern_lower="${pattern,,}"
    
    for opt in "${opts_arr[@]}"; do
        local opt_lower="${opt,,}"
        if [[ "$opt_lower" == *"$pattern_lower"* ]]; then
            substring_matches+=("$opt")
        fi
    done
    
    if [[ ${#substring_matches[@]} -gt 0 ]]; then
        COMPREPLY=("${substring_matches[@]}")
        return 0
    fi
    
    # 如果所有匹配策略都失败，返回空结果
    COMPREPLY=()
    return 1
}

# ==================== 主补全函数 ====================
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

	# 根据参数类型动态生成补全
	if [[ -n "$prev_param_type" && "$prev_param_type" == "required" ]]; then
		case "$prev_value_type" in
			enum)
				local enum_key="${context}|${words[cword-1]}"
				local enum_opts="${{{.ProgramName}}_enum_options[$enum_key]:-}"
				
				if [[ -n "$enum_opts" ]]; then
					# 对枚举选项也使用智能匹配
					_{{.ProgramName}}_intelligent_match "$cur" "$enum_opts"
					return 0
				fi
				;;
			string)
				# 字符串类型 - 提供文件和目录路径补全
				COMPREPLY=($(compgen -f -d -- "$cur"))
				return 0
				;;
			*)
				# 默认值补全 - 使用智能匹配
				_{{.ProgramName}}_intelligent_match "$cur" "$current_context_opts"
				return 0
				;;
		esac
	fi

	# 主要的标志和命令补全 - 使用智能匹配算法
	_{{.ProgramName}}_intelligent_match "$cur" "$current_context_opts"
	return 0
}

# ==================== 调试和诊断函数 ====================
# 补全系统健康检查函数 (可选，用于调试)
_{{.ProgramName}}_completion_debug() {
    echo "=== {{.ProgramName}} 补全系统诊断 ==="
    echo "Bash版本: $BASH_VERSION"
    echo "补全函数状态: $(type -t _{{.ProgramName}})"
    echo "命令树条目数: ${#{{.ProgramName}}_cmd_tree[@]}"
    echo "标志参数数: ${#{{.ProgramName}}_flag_params[@]}"
    echo "枚举选项数: ${#{{.ProgramName}}_enum_options[@]}"
    echo "模糊补全状态: $([ ${{.ProgramName}}_FUZZY_COMPLETION_ENABLED -eq 1 ] && echo "启用" || echo "禁用")"
    echo "候选项阈值: ${{.ProgramName}}_FUZZY_MAX_CANDIDATES"
    echo "缓存条目数: ${#{{.ProgramName}}_fuzzy_cache[@]}"
    echo ""
    echo "使用方法: 在命令行输入 '_{{.ProgramName}}_completion_debug' 查看此信息"
}

# 注册补全函数
complete -F _{{.ProgramName}} {{.ProgramName}}
`
)
const (
	// PowerShell自动补全脚本头部
	PwshFunctionHeader = `# -------------------------- Configuration Area (Need to be modified according to actual commands) --------------------------
# Command Name
${{.SanitizedName}}_commandName = "{{.ProgramName}}"

# 1. Command Tree
${{.SanitizedName}}_cmdTree = @(
{{.CmdTree}}
)

# 2. Flag Parameter Definitions
${{.SanitizedName}}_flagParams = @(
{{.FlagParams}}
)

# -----------------------------------------------------------------------------------

# ==================== 模糊补全配置参数 ====================
# 模糊补全功能开关 (设置为$false禁用，$true启用)
$script:{{.SanitizedName}}_FUZZY_COMPLETION_ENABLED = $true

# 启用模糊补全的最大候选项数量阈值
# 超过此数量将回退到传统前缀匹配以保证性能
$script:{{.SanitizedName}}_FUZZY_MAX_CANDIDATES = 120

# 模糊匹配的最小输入长度 (小于此长度不启用模糊匹配)
$script:{{.SanitizedName}}_FUZZY_MIN_PATTERN_LENGTH = 2

# 模糊匹配分数阈值 (0-100，分数低于此值的匹配将被过滤)
$script:{{.SanitizedName}}_FUZZY_SCORE_THRESHOLD = 25

# 模糊匹配最大返回结果数
$script:{{.SanitizedName}}_FUZZY_MAX_RESULTS = 10

# 缓存大小控制参数
# 缓存条目数量超过此阈值时将清空缓存以防止内存无限增长
$script:{{.SanitizedName}}_FUZZY_CACHE_MAX_SIZE = 500

# 模糊匹配结果缓存 (格式: "pattern|candidate" -> score)
$script:{{.SanitizedName}}_fuzzyCache = @{}

# ==================== 模糊匹配核心算法 ====================

# 高性能模糊评分函数 - 使用优化的字符串操作
# 参数: $Pattern=输入模式, $Candidate=候选字符串
# 返回: 0-100的整数分数
function Get-{{.SanitizedName}}FuzzyScoreFast {
    param(
        [string]$Pattern,
        [string]$Candidate
    )
    
    $patternLen = $Pattern.Length
    $candidateLen = $Candidate.Length
    
    # 性能优化1: 长度预检查 - 候选项太短直接返回0
    if ($candidateLen -lt $patternLen) {
        return 0
    }
    
    # 性能优化2: 完全匹配检查 - 避免不必要的复杂计算
    if ($Candidate.StartsWith($Pattern, [System.StringComparison]::OrdinalIgnoreCase)) {
        return 100  # 前缀完全匹配给最高分
    }
    
    # 性能优化3: 字符存在性预检查 - 快速排除不可能的匹配
    $patternLower = $Pattern.ToLowerInvariant()
    $candidateLower = $Candidate.ToLowerInvariant()
    
    foreach ($char in $patternLower.ToCharArray()) {
        if ($candidateLower.IndexOf($char) -eq -1) {
            return 0  # 必需字符不存在，直接返回
        }
    }
    
    # 核心匹配算法 - 计算字符匹配度和连续性
    $matched = 0           # 匹配的字符数
    $consecutive = 0       # 当前连续匹配长度
    $maxConsecutive = 0    # 最大连续匹配长度
    $candidatePos = 0      # 候选字符串当前搜索位置
    $startBonus = 0        # 起始位置奖励
    
    # 检查是否从开头匹配 (大小写不敏感)
    if ($candidateLower.StartsWith($patternLower)) {
        $startBonus = 20  # 起始匹配给20分奖励
    }
    
    # 逐字符匹配算法
    for ($i = 0; $i -lt $patternLen; $i++) {
        $patternChar = $patternLower[$i]
        $found = $false
        
        # 在候选字符串中查找当前模式字符
        for ($j = $candidatePos; $j -lt $candidateLen; $j++) {
            if ($candidateLower[$j] -eq $patternChar) {
                $matched++
                $found = $true
                
                # 连续性检查 - 连续匹配的字符得分更高
                if ($j -eq $candidatePos) {
                    $consecutive++
                    if ($consecutive -gt $maxConsecutive) {
                        $maxConsecutive = $consecutive
                    }
                } else {
                    $consecutive = 1  # 重置连续计数
                }
                
                $candidatePos = $j + 1  # 更新搜索位置
                break
            }
        }
        
        # 如果某个字符未找到，重置连续计数
        if (-not $found) {
            $consecutive = 0
        }
    }
    
    # 评分计算 - 使用整数运算
    # 基础分数: (匹配字符数 / 模式长度) * 60
    $baseScore = [Math]::Floor(($matched * 60) / $patternLen)
    
    # 连续性奖励: (最大连续长度 / 模式长度) * 20
    $consecutiveBonus = [Math]::Floor(($maxConsecutive * 20) / $patternLen)
    
    # 长度惩罚: 候选字符串越长，分数略微降低
    $lengthPenalty = [Math]::Min(($candidateLen - $patternLen), 10)
    
    # 最终分数计算
    $finalScore = $baseScore + $consecutiveBonus + $startBonus - $lengthPenalty
    
    # 确保分数在0-100范围内
    return [Math]::Max(0, [Math]::Min(100, $finalScore))
}

# 带缓存的模糊评分函数 - 避免重复计算提高性能
# 参数: $Pattern=输入模式, $Candidate=候选字符串
function Get-{{.SanitizedName}}FuzzyScoreCached {
    param(
        [string]$Pattern,
        [string]$Candidate
    )
    
    $cacheKey = "$Pattern|$Candidate"
    
    # 缓存命中检查
    if ($script:{{.SanitizedName}}_fuzzyCache.ContainsKey($cacheKey)) {
        return $script:{{.SanitizedName}}_fuzzyCache[$cacheKey]
    }
    
    # 计算分数并缓存
    $score = Get-{{.SanitizedName}}FuzzyScoreFast -Pattern $Pattern -Candidate $Candidate
    
    # 缓存大小控制 - 防止内存无限增长
    if ($script:{{.SanitizedName}}_fuzzyCache.Count -gt $script:{{.SanitizedName}}_FUZZY_CACHE_MAX_SIZE) {
        $script:{{.SanitizedName}}_fuzzyCache.Clear()  # 清空缓存
    }
    
    $script:{{.SanitizedName}}_fuzzyCache[$cacheKey] = $score
    return $score
}

# 智能补全匹配函数 - 分级匹配策略
# 参数: $Pattern=输入模式, $Options=候选选项数组
function Get-{{.SanitizedName}}IntelligentMatches {
    param(
        [string]$Pattern,
        [array]$Options
    )
    
    $patternLen = $Pattern.Length
    $totalCandidates = $Options.Count
    
    # 性能保护: 候选项过多时禁用模糊匹配
    if ($totalCandidates -gt $script:{{.SanitizedName}}_FUZZY_MAX_CANDIDATES) {
        # 回退到传统前缀匹配
        $prefixMatches = @()
        foreach ($option in $Options) {
            if ($option -like "$Pattern*") {
                $prefixMatches += $option
            }
        }
        return $prefixMatches
    }
    
    # 第1级: 精确前缀匹配 (最快，优先级最高) - 使用ArrayList优化性能
    $exactMatches = [System.Collections.ArrayList]::new()
    foreach ($option in $Options) {
        if ($option.StartsWith($Pattern, [System.StringComparison]::Ordinal)) {
            [void]$exactMatches.Add($option)
        }
    }
    
    # 如果有精确匹配且数量合理，直接返回
    if ($exactMatches.Count -gt 0 -and $exactMatches.Count -le 12) {
        return $exactMatches.ToArray()
    }
    
    # 第2级: 大小写不敏感前缀匹配
    if ($exactMatches.Count -eq 0) {
        $caseInsensitiveMatches = [System.Collections.ArrayList]::new()
        foreach ($option in $Options) {
            if ($option.StartsWith($Pattern, [System.StringComparison]::OrdinalIgnoreCase)) {
                [void]$caseInsensitiveMatches.Add($option)
            }
        }
        
        # 如果有大小写不敏感匹配，返回
        if ($caseInsensitiveMatches.Count -gt 0) {
            return $caseInsensitiveMatches.ToArray()
        }
    }
    
    # 第3级: 模糊匹配 (最慢，仅在必要时使用) - 使用ArrayList优化性能
    if ($script:{{.SanitizedName}}_FUZZY_COMPLETION_ENABLED -and $patternLen -ge $script:{{.SanitizedName}}_FUZZY_MIN_PATTERN_LENGTH) {
        $scoredMatches = [System.Collections.ArrayList]::new()
        
        # 对所有候选项进行模糊评分
        foreach ($option in $Options) {
            $score = Get-{{.SanitizedName}}FuzzyScoreCached -Pattern $Pattern -Candidate $option
            
            # 只保留分数达到阈值的匹配
            if ($score -ge $script:{{.SanitizedName}}_FUZZY_SCORE_THRESHOLD) {
                [void]$scoredMatches.Add(@{
                    Option = $option
                    Score = $score
                })
            }
        }
        
        # 如果有模糊匹配结果，按分数排序并返回前N个
        if ($scoredMatches.Count -gt 0) {
            # 按分数降序排序
            $sortedMatches = $scoredMatches | Sort-Object Score -Descending
            
            # 提取选项名称，限制返回数量 - 使用ArrayList优化
            $fuzzyResults = [System.Collections.ArrayList]::new()
            $count = 0
            foreach ($match in $sortedMatches) {
                if ($count -ge $script:{{.SanitizedName}}_FUZZY_MAX_RESULTS) {
                    break
                }
                [void]$fuzzyResults.Add($match.Option)
                $count++
            }
            
            return $fuzzyResults.ToArray()
        }
    }
    
    # 第4级: 子字符串匹配 (最后的备选方案) - 使用ArrayList优化性能
    $substringMatches = [System.Collections.ArrayList]::new()
    $patternLower = $Pattern.ToLowerInvariant()
    
    foreach ($option in $Options) {
        $optionLower = $option.ToLowerInvariant()
        if ($optionLower.Contains($patternLower)) {
            [void]$substringMatches.Add($option)
        }
    }
    
    if ($substringMatches.Count -gt 0) {
        return $substringMatches.ToArray()
    }
    
    # 如果所有匹配策略都失败，返回空数组
    return @()
}

# -------------------------- Completion Logic Implementation ------------------------
$scriptBlock = {
    param(
        $wordToComplete,
        $commandAst,
        $cursorPosition
    )

    # 初始化缓存和索引（仅在首次调用时创建）
    if (-not $script:{{.SanitizedName}}_contextIndex) {
        $script:{{.SanitizedName}}_contextIndex = @{}
        $script:{{.SanitizedName}}_flagIndex = @{}
        
        # 预编译正则表达式以提高性能
        $script:{{.SanitizedName}}_flagRegex = [regex]::new('^-', [System.Text.RegularExpressions.RegexOptions]::Compiled)
        
        # 构建上下文索引以提高查找性能
        foreach ($item in ${{.SanitizedName}}_cmdTree) {
            if ($item.Context) {
                $script:{{.SanitizedName}}_contextIndex[$item.Context] = $item
            }
        }
        
        # 构建标志索引以提高查找性能
        foreach ($flag in ${{.SanitizedName}}_flagParams) {
            if ($flag.Context -and $flag.Parameter) {
                $key = "$($flag.Context)|$($flag.Parameter)"
                $script:{{.SanitizedName}}_flagIndex[$key] = $flag
            }
        }
    }

    try {
        # 1. 解析令牌
        $tokens = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
        if (-not $tokens -or $tokens.Count -eq 0) {
            return @()
        }
        
        $currentIndex = $tokens.Count - 1
        $prevElement = if ($currentIndex -ge 1) { $tokens[$currentIndex - 1] } else { $null }

        # 2. 计算当前命令上下文（优化版本）
        $context = "/"
        for ($i = 1; $i -le $currentIndex; $i++) {
            $elem = $tokens[$i]
            if ($script:{{.SanitizedName}}_flagRegex.IsMatch($elem)) { break }
            
            $nextContext = "$context$elem/"
            # 使用索引进行O(1)查找
            if ($script:{{.SanitizedName}}_contextIndex.ContainsKey($nextContext)) {
                $context = $nextContext
            } else {
                break
            }
        }

        # 3. 获取当前上下文的可用选项（优化版本）
        $currentContextItem = $script:{{.SanitizedName}}_contextIndex[$context]
        $currentOptions = if ($currentContextItem) { $currentContextItem.Options } else { @() }

        # 4. 优先补全当前级别的所有选项（子命令 + 标志）
        if ($currentOptions -and $currentOptions.Count -gt 0) {
            # 使用ArrayList提高数组操作性能
            $matchingOptions = [System.Collections.ArrayList]::new()
            $wordPattern = "$wordToComplete*"
            
            foreach ($option in $currentOptions) {
                if ($option -like $wordPattern) {
                    $result = if ($script:{{.SanitizedName}}_flagRegex.IsMatch($option)) { $option } else { "$option " }
                    [void]$matchingOptions.Add($result)
                }
            }
            if ($matchingOptions.Count -gt 0) {
                return $matchingOptions.ToArray()
            }
        }

        # 5. 补全标志本身（如 --ty -> --type）- 使用智能匹配
        if ($script:{{.SanitizedName}}_flagRegex.IsMatch($wordToComplete)) {
            # 收集当前上下文的所有标志 - 使用ArrayList优化性能
            $contextFlags = [System.Collections.ArrayList]::new()
            foreach ($flag in ${{.SanitizedName}}_flagParams) {
                if ($flag.Context -eq $context) {
                    [void]$contextFlags.Add($flag.Parameter)
                }
            }
            
            if ($contextFlags.Count -gt 0) {
                # 使用智能匹配获取最佳标志匹配
                $flagMatches = Get-{{.SanitizedName}}IntelligentMatches -Pattern $wordToComplete -Options $contextFlags.ToArray()
                if ($flagMatches.Count -gt 0) {
                    return $flagMatches
                }
            }
        }

        # 6. 枚举/预设值补全
        if ($prevElement -and $script:{{.SanitizedName}}_flagRegex.IsMatch($prevElement)) {
            $flagKey = "$context|$prevElement"
            $paramDef = $script:{{.SanitizedName}}_flagIndex[$flagKey]
            
            if ($paramDef) {
                switch ($paramDef.ValueType) {
                    'enum' {
                        if (-not $wordToComplete) {
                            # 当前单词为空 → 返回所有枚举值
                            return $paramDef.Options
                        } else {
                            # 使用智能匹配进行枚举值补全
                            $enumMatches = Get-{{.SanitizedName}}IntelligentMatches -Pattern $wordToComplete -Options $paramDef.Options
                            return $enumMatches
                        }
                    }
                    'string' {
                        # 字符串类型 - 提供文件和目录路径补全 - 使用ArrayList优化性能
                        $pathMatches = [System.Collections.ArrayList]::new()
                        
                        # 获取当前路径的目录部分
                        $basePath = if ($wordToComplete -and (Split-Path $wordToComplete -Parent)) {
                            Split-Path $wordToComplete -Parent
                        } else {
                            "."
                        }
                        
                        # 获取文件名部分用于过滤
                        $fileName = if ($wordToComplete) {
                            Split-Path $wordToComplete -Leaf
                        } else {
                            ""
                        }
                        
                        # 预编译文件名匹配模式
                        $filePattern = "$fileName*"
                        
                        try {
                            # 获取目录和文件
                            $items = Get-ChildItem -Path $basePath -ErrorAction SilentlyContinue | Where-Object {
                                $_.Name -like $filePattern
                            }
                            
                            foreach ($item in $items) {
                                $fullPath = if ($basePath -eq ".") {
                                    $item.Name
                                } else {
                                    Join-Path $basePath $item.Name
                                }
                                
                                # 目录添加路径分隔符
                                if ($item.PSIsContainer) {
                                    [void]$pathMatches.Add("$fullPath/")
                                } else {
                                    [void]$pathMatches.Add($fullPath)
                                }
                            }
                        }
                        catch {
                            # 路径访问失败时返回空数组
                        }
                        
                        return $pathMatches.ToArray()
                    }
                    default {
                        return @()
                    }
                }
            }
        }

        # 7. 无匹配
        return @()
    }
    catch {
        # 错误处理：返回空数组而不是抛出异常
        Write-Debug "PowerShell补全错误: $($_.Exception.Message)"
        return @()
    }
}

# ==================== 调试和诊断功能 ====================

# 补全系统健康检查函数 (可选，用于调试)
function Get-{{.SanitizedName}}CompletionDebug {
    Write-Host "=== {{.SanitizedName}} PowerShell补全系统诊断 ===" -ForegroundColor Cyan
    Write-Host "PowerShell版本: $($PSVersionTable.PSVersion)" -ForegroundColor Green
    Write-Host "补全函数状态: $(if (Get-Command Register-ArgumentCompleter -ErrorAction SilentlyContinue) { '已注册' } else { '未注册' })" -ForegroundColor Green
    Write-Host "命令树条目数: $(${{.SanitizedName}}_cmdTree.Count)" -ForegroundColor Green
    Write-Host "标志参数数: $(${{.SanitizedName}}_flagParams.Count)" -ForegroundColor Green
    Write-Host "模糊补全状态: $(if ($script:{{.SanitizedName}}_FUZZY_COMPLETION_ENABLED) { '启用' } else { '禁用' })" -ForegroundColor Green
    Write-Host "候选项阈值: $script:{{.SanitizedName}}_FUZZY_MAX_CANDIDATES" -ForegroundColor Green
    Write-Host "缓存条目数: $($script:{{.SanitizedName}}_fuzzyCache.Count)" -ForegroundColor Green
    Write-Host ""
    Write-Host "使用方法: 在PowerShell中输入 'Get-{{.SanitizedName}}CompletionDebug' 查看此信息" -ForegroundColor Yellow
}

# 模糊匹配测试函数 (用于调试和验证)
function Test-{{.SanitizedName}}FuzzyMatch {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Pattern,
        [Parameter(Mandatory=$true)]
        [string]$Candidate
    )
    
    $score = Get-{{.SanitizedName}}FuzzyScoreFast -Pattern $Pattern -Candidate $Candidate
    Write-Host "模式: '$Pattern' 匹配候选: '$Candidate' 得分: $score" -ForegroundColor Cyan
    
    # 详细分析
    if ($score -ge 80) {
        Write-Host "匹配质量: 优秀" -ForegroundColor Green
    } elseif ($score -ge 50) {
        Write-Host "匹配质量: 良好" -ForegroundColor Yellow
    } elseif ($score -ge $script:{{.SanitizedName}}_FUZZY_SCORE_THRESHOLD) {
        Write-Host "匹配质量: 可接受" -ForegroundColor DarkYellow
    } else {
        Write-Host "匹配质量: 不匹配" -ForegroundColor Red
    }
    
    return $score
}

Register-ArgumentCompleter -CommandName ${{.SanitizedName}}_commandName -ScriptBlock $scriptBlock
`
)

VARIABLES

var (
	// CompletionNotesCN 中文版本注意事项
	CompletionNotesCN = []string{
		"Windows环境: 需要PowerShell 5.1或更高版本以支持Register-ArgumentCompleter",
		"Linux环境: 需要bash 4.0或更高版本以支持关联数组特性",
		"请确保您的环境满足上述版本要求，否则自动补全功能可能无法正常工作",
	}

	// CompletionNotesEN 英文版本注意事项
	CompletionNotesEN = []string{
		"Windows environment: Requires PowerShell 5.1 or higher to support Register-ArgumentCompleter",
		"Linux environment: Requires bash 4.0 or higher to support associative array features",
		"Please ensure your environment meets the above version requirements, otherwise the auto-completion feature may not work properly",
	}
)
    生成标志的注意事项

var CompletionExamplesCN = []types.ExampleInfo{
	{Description: "Linux环境 临时启用", Usage: "source <(%s --generate-shell-completion bash)"},
	{Description: "Linux环境 永久启用(添加到~/.bashrc)", Usage: "echo \"source <(%s --generate-shell-completion bash)\" >> ~/.bashrc"},

	{Description: "Windows环境 临时启用", Usage: "%s --generate-shell-completion powershell | Out-String | Invoke-Expression"},
	{Description: "Windows环境 永久启用(添加到PowerShell配置文件)", Usage: "echo \"%s --generate-shell-completion powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
}
    内置自动补全命令的示例使用（中文）

var CompletionExamplesEN = []types.ExampleInfo{
	{Description: "Linux environment temporary activation", Usage: "source <(%s --generate-shell-completion bash)"},
	{Description: "Linux environment permanent activation (add to ~/.bashrc)", Usage: "echo \"source <(%s --generate-shell-completion bash)\" >> ~/.bashrc"},

	{Description: "Windows environment temporary activation", Usage: "%s --generate-shell-completion powershell | Out-String | Invoke-Expression"},
	{Description: "Windows environment permanent activation (add to PowerShell profile)", Usage: "echo \"%s --generate-shell-completion powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
}
    内置自动补全命令的示例使用（英文）


FUNCTIONS

func GenerateShellCompletion(ctx *types.CmdContext, shellType string) (string, error)
    GenerateShellCompletion 生成shell自动补全脚本

    参数:
      - ctx: 命令上下文
      - shellType: shell类型 ("bash", "pwsh", "powershell")

    返回值：
      - string: 自动补全脚本
      - error: 错误信息


TYPES

type FlagParam struct {
	CommandPath string   // 命令路径，如 "/cmd/subcmd"
	Name        string   // 标志名称(保留原始大小写)
	Type        string   // 参数需求类型: "required"|"optional"|"none"
	ValueType   string   // 参数值类型: "path"|"string"|"number"|"enum"|"bool"等
	EnumOptions []string // 枚举类型的可选值列表
}
    FlagParam 表示标志参数及其需求类型和值类型

