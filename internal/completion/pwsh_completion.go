package completion

import (
	"bytes"
	"path/filepath"
	"strings"
)

// formatOptions 将选项列表格式化为PowerShell数组字符串
//
// 参数:
// - buf: 输出缓冲区
// - options: 选项列表
// - escape: 字符串转义函数
func formatOptions(buf *bytes.Buffer, options []string, escape func(string) string) {
	for i, opt := range options {
		// 只有不为空的选项才添加到缓冲区
		if opt == "" {
			continue
		}

		// 如果不是第一个选项，则添加逗号
		if i > 0 {
			buf.WriteString(", ")
		}

		// 添加选项
		buf.WriteByte('\'')
		buf.WriteString(escape(opt))
		buf.WriteByte('\'')
	}
}

// generatePwshCommandTreeEntry 生成PowerShell命令树条目
// 使用对象池优化内存分配，避免创建临时缓冲区和Replacer
//
// 参数:
// - cmdTreeEntries: 命令树条目缓冲区
// - cmdPath: 命令路径
// - cmdOpts: 命令选项
func generatePwshCommandTreeEntry(cmdTreeEntries *bytes.Buffer, cmdPath string, cmdOpts []string) {
	// 使用对象池构建命令树条目，避免创建临时缓冲区和strings.NewReplacer的开销
	cmdTreeItem := buildString(func(builder *strings.Builder) {
		builder.WriteString("\t@{ Context = \"")
		builder.WriteString(cmdPath)
		builder.WriteString("\"; Options = @(")

		// 直接在builder中格式化选项，避免额外的字符串分配
		first := true
		for _, opt := range cmdOpts {
			if opt == "" {
				continue
			}

			if !first {
				builder.WriteString(", ")
			}
			first = false

			builder.WriteByte('\'')
			builder.WriteString(escapePwshString(opt))
			builder.WriteByte('\'')
		}

		builder.WriteString(") }")
	})

	cmdTreeEntries.WriteString(cmdTreeItem)
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
	// 构建标志参数和枚举选项
	flagParamsBuf := bytes.NewBuffer(make([]byte, 0, len(params)*100)) // 预分配容量

	// 处理根命令选项
	rootOptsBuf := bytes.NewBuffer(make([]byte, 0, len(rootCmdOpts)*20))
	formatOptions(rootOptsBuf, rootCmdOpts, escapePwshString)

	// 处理标志参数
	for i, param := range params {
		// 生成带枚举选项的标志参数条目
		enumOptions := ""
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			optionsBuf := bytes.NewBuffer(make([]byte, 0, len(param.EnumOptions)*15))
			formatOptions(optionsBuf, param.EnumOptions, escapePwshString)
			enumOptions = optionsBuf.String()
		}

		// 使用命名占位符替换位置参数
		flagReplacer := strings.NewReplacer(
			"{{.Context}}", param.CommandPath,
			"{{.Parameter}}", param.Name,
			"{{.ParamType}}", param.Type,
			"{{.ValueType}}", param.ValueType,
			"{{.Options}}", enumOptions,
		)
		flagParamsBuf.WriteString(flagReplacer.Replace(PwshFlagParamItem))

		// 条目之间添加逗号，非最后一个条目
		if i < len(params)-1 {
			flagParamsBuf.WriteString(",\n")
		}
	}

	// 清理程序名，去除可能的后缀
	sanitizedProgramName := strings.TrimSuffix(programName, filepath.Ext(programName))

	// 生成根命令条目
	rootReplacer := strings.NewReplacer(
		"{{.Context}}", "/",
		"{{.Options}}", rootOptsBuf.String(),
	)
	// 生成根命令条目
	rootCmdEntry := rootReplacer.Replace(PwshCmdTreeItem)

	// 如果命令树条目不为空，则添加逗号
	if cmdTreeEntries != "" {
		rootCmdEntry += ",\n" + cmdTreeEntries
	}

	// 使用命名占位符替换位置参数
	completionReplacer := strings.NewReplacer(
		"{{.SanitizedName}}", sanitizedProgramName, // 替换程序名称
		"{{.ProgramName}}", programName, // 替换程序名称
		"{{.CmdTree}}", rootCmdEntry, // 替换命令树条目
		"{{.FlagParams}}", flagParamsBuf.String(), // 替换标志参数
	)

	// 写入PowerShell自动补全脚本
	_, _ = buf.WriteString(completionReplacer.Replace(PwshFunctionHeader))
}

// pwshEscapeMap PowerShell特殊字符转义映射表
// 使用全局map提高转义性能，避免重复的switch判断
var pwshEscapeMap = map[byte][]byte{
	'\'': {'\'', '\''}, // 单引号转义为两个单引号
	'\\': {'\\', '\\'}, // 反斜杠转义为两个反斜杠
	'$':  {'`', '$'},   // 美元符号转义
	'`':  {'`', '`'},   // 反引号转义
	'"':  {'`', '"'},   // 双引号转义
	'&':  {'`', '&'},   // 与符号转义
	'|':  {'`', '|'},   // 管道符转义
	';':  {'`', ';'},   // 分号转义
	'<':  {'`', '<'},   // 小于号转义
	'>':  {'`', '>'},   // 大于号转义
	'(':  {'`', '('},   // 左括号转义
	')':  {'`', ')'},   // 右括号转义
}

// escapePwshString 转义PowerShell字符串中的特殊字符
// 优化版本：使用全局map进行O(1)查找，提升性能
//
// 参数:
// - s: 需要转义的字符串
//
// 返回:
// - 转义后的字符串
func escapePwshString(s string) string {
	// 预计算所需容量：最坏情况下每个字符都需要转义
	buf := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if escaped, exists := pwshEscapeMap[c]; exists {
			buf = append(buf, escaped...)
		} else {
			buf = append(buf, c)
		}
	}
	return string(buf)
}

const (
	// 标志参数条目(含枚举选项)
	PwshFlagParamItem = "	@{ Context = \"{{.Context}}\"; Parameter = \"{{.Parameter}}\"; ParamType = \"{{.ParamType}}\"; ValueType = \"{{.ValueType}}\"; Options = @({{.Options}}) }"
	// 命令树条目
	PwshCmdTreeItem = "	@{ Context = \"{{.Context}}\"; Options = @({{.Options}}) }"
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

# -------------------------- Completion Logic Implementation ------------------------
$scriptBlock = {
    param(
        $wordToComplete,
        $commandAst,
        $cursorPosition
    )

    # 初始化缓存（仅在首次调用时创建）
    if (-not $script:{{.SanitizedName}}_contextIndex) {
        $script:{{.SanitizedName}}_contextIndex = @{}
        $script:{{.SanitizedName}}_flagIndex = @{}
        
        # 构建上下文索引以提高查找性能
        foreach ($item in ${{{.SanitizedName}}_cmdTree}) {
            if ($item.Context) {
                $script:{{.SanitizedName}}_contextIndex[$item.Context] = $item
            }
        }
        
        # 构建标志索引以提高查找性能
        foreach ($flag in ${{{.SanitizedName}}_flagParams}) {
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
            if ($elem -match '^-') { break }
            
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
            $matchingOptions = @()
            foreach ($option in $currentOptions) {
                if ($option -like "$wordToComplete*") {
                    $matchingOptions += if ($option -match '^-') { $option } else { "$option " }
                }
            }
            if ($matchingOptions.Count -gt 0) {
                return $matchingOptions
            }
        }

        # 5. 补全标志本身（如 --ty -> --type）
        if ($wordToComplete -match '^-') {
            $flagMatches = @()
            foreach ($flag in ${{{.SanitizedName}}_flagParams}) {
                if ($flag.Context -eq $context -and $flag.Parameter -like "$wordToComplete*") {
                    $flagMatches += $flag.Parameter
                }
            }
            if ($flagMatches.Count -gt 0) {
                return $flagMatches
            }
        }

        # 6. 枚举/预设值补全
        if ($prevElement -match '^-') {
            $flagKey = "$context|$prevElement"
            $paramDef = $script:{{.SanitizedName}}_flagIndex[$flagKey]
            
            if ($paramDef) {
                switch ($paramDef.ValueType) {
                    'enum' {
                        if (-not $wordToComplete) {
                            # 当前单词为空 → 返回所有枚举值
                            return $paramDef.Options
                        } else {
                            # 前缀过滤
                            $enumMatches = @()
                            foreach ($option in $paramDef.Options) {
                                if ($option -like "$wordToComplete*") {
                                    $enumMatches += $option
                                }
                            }
                            return $enumMatches
                        }
                    }
                    'string' {
                        # 字符串类型 - 提供文件和目录路径补全
                        $pathMatches = @()
                        
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
                        
                        try {
                            # 获取目录和文件
                            $items = Get-ChildItem -Path $basePath -ErrorAction SilentlyContinue | Where-Object {
                                $_.Name -like "$fileName*"
                            }
                            
                            foreach ($item in $items) {
                                $fullPath = if ($basePath -eq ".") {
                                    $item.Name
                                } else {
                                    Join-Path $basePath $item.Name
                                }
                                
                                # 目录添加路径分隔符
                                if ($item.PSIsContainer) {
                                    $pathMatches += "$fullPath/"
                                } else {
                                    $pathMatches += $fullPath
                                }
                            }
                        }
                        catch {
                            # 路径访问失败时返回空数组
                        }
                        
                        return $pathMatches
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

Register-ArgumentCompleter -CommandName ${{{.SanitizedName}}_commandName} -ScriptBlock $scriptBlock
`
)
