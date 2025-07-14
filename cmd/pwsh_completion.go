// pwsh_completion.go - PowerShell 补全脚本生成器
package cmd

import (
	"bytes"
	"fmt"
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
//
// 参数:
// - cmdTreeEntries: 命令树条目缓冲区
// - cmdPath: 命令路径
// - cmdOpts: 命令选项
func generatePwshCommandTreeEntry(cmdTreeEntries *bytes.Buffer, cmdPath string, cmdOpts []string) {
	// 预分配缓冲区容量，减少内存分配
	optsBuf := bytes.NewBuffer(make([]byte, 0, len(cmdOpts)*20)) // 假设平均每个选项20字节

	// 格式化命令选项
	formatOptions(optsBuf, cmdOpts, escapePwshString)

	// 添加命令树条目
	fmt.Fprintf(cmdTreeEntries, PwshCmdTreeItem, cmdPath, optsBuf.String())
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
		// 生成标志参数条目
		// 生成带枚举选项的标志参数条目
		enumOptions := ""
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			optionsBuf := bytes.NewBuffer(make([]byte, 0, len(param.EnumOptions)*15))
			formatOptions(optionsBuf, param.EnumOptions, escapePwshString)
			enumOptions = optionsBuf.String()
		}
		fmt.Fprintf(flagParamsBuf, PwshFlagParamItem,
			param.CommandPath,
			param.Name,
			param.Type,
			param.ValueType,
			enumOptions,
		)

		// 条目之间添加逗号，非最后一个条目
		if i < len(params)-1 {
			flagParamsBuf.WriteString(",\n")
		}
	}

	// 生成根命令条目
	rootCmdEntry := fmt.Sprintf(PwshCmdTreeItem, "/", rootOptsBuf.String())
	if cmdTreeEntries != "" {
		rootCmdEntry += ",\n" + cmdTreeEntries
	}

	// 写入PowerShell自动补全脚本
	fmt.Fprintf(buf, PwshFunctionHeader,
		programName,
		rootCmdEntry,
		flagParamsBuf.String(),
	) // 移除独立的枚举选项数组
}

// escapePwshString 转义PowerShell字符串中的特殊字符
// 优化：单次循环处理所有转义，减少字符串分配
func escapePwshString(s string) string {
	// 预计算所需容量：最坏情况下每个字符都需要转义
	buf := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\'':
			buf = append(buf, '\'', '\'') // 单引号转义为两个单引号
		case '\\':
			buf = append(buf, '\\', '\\') // 反斜杠转义为两个反斜杠
		default:
			buf = append(buf, c)
		}
	}
	return string(buf)
}

const (
	// PowerShell补全模板
	PwshFunctionHeader = `# -------------------------- 配置区域(需根据实际命令修改) --------------------------
# 命令名称(替换为你的实际命令名，如"mycmd")
$commandName = "%s"

# 1. 命令树定义(数组形式，每个元素包含 Context 和 Options)
# Context: 上下文路径(如"/"、"/init/");Options: 该层级可用选项数组(不区分大小写)
$cmdTree = @(
%s  
)

# 2. 标志参数定义(数组形式，每个元素包含 Context、Parameter、ParamType、ValueType)
# ParamType: required(必须带值)、optional(可选值);ValueType: path|number|ip|enum|url等
$flagParams = @(
%s  
)

# -----------------------------------------------------------------------------------


# -------------------------- 补全逻辑实现(核心处理大小写不敏感) --------------------------
$scriptBlock = {
    param(
        $wordToComplete,       # 当前输入的待补全文本
        $commandAst,           # 命令抽象语法树
        $cursorPosition        # 光标位置
    )

    # 1. 解析当前命令参数列表
    $commandElements = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
    $currentIndex = $commandElements.Count - 1  # 当前补全位置索引
    $prevElement = if ($currentIndex -ge 1) { $commandElements[$currentIndex - 1] } else { $null }

    # 2. 计算当前命令上下文(不区分大小写匹配)
    $context = "/"
    for ($i = 1; $i -lt $currentIndex; $i++) {  # 跳过命令名本身(索引0)
        $elem = $commandElements[$i]
        # 遍历命令树数组，查找匹配当前上下文+元素的子上下文(不区分大小写)
        $match = $cmdTree | Where-Object {
            $_.Context -eq "$context$elem/" -or $_.Context -eq "$context$($elem.ToLower())/" -or $_.Context -eq "$context$($elem.ToUpper())/"
        }
        if ($match) {
            $context = $match.Context  # 更新上下文为匹配的路径
        }
    }

    # 3. 获取当前上下文可用的选项(从数组中查找)
    $currentOptions = @()
    $contextMatch = $cmdTree | Where-Object { $_.Context -eq $context }
    if ($contextMatch) {
        $currentOptions = $contextMatch.Options
    }

    # 4. 处理前一个参数需要值的情况(不区分大小写匹配参数)
    $completionItems = @()
    if ($prevElement) {
        # 查找前一个参数的定义(上下文+参数名，不区分大小写)
        $paramDef = $flagParams | Where-Object {
            $_.Context -eq $context -and 
            ($_.Parameter -eq $prevElement -or $_.Parameter -eq $prevElement.ToLower() -or $_.Parameter -eq $prevElement.ToUpper())
        }
        if ($paramDef) {
            $paramType = $paramDef.ParamType
            $valueType = $paramDef.ValueType
            # 处理需要值的参数
            if ($paramType -in "required", "optional") {
                switch ($valueType) {
                    "path" {
                        $completionItems = Get-ChildItem -Path "$wordToComplete*" -ErrorAction SilentlyContinue | 
                            ForEach-Object { $_.FullName }
                    }
                    "number" {
                        $completionItems = 1..10 | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object { "$_" }
                    }
                    "ip" {
                        $completionItems = @("192.168.", "10.0.", "172.16.", "127.0.0.") | 
                            Where-Object { $_ -like "$wordToComplete*" }
                    }
                    "enum" {
                        # 直接从标志参数获取枚举选项
                        $completionItems = $paramDef.Options | Where-Object { $_ -like "$wordToComplete*" }
                    }
                    "url" {
                        $completionItems = @("http://", "https://", "ftp://") | Where-Object { $_ -like "$wordToComplete*" }
                    }
                    default {
                        $completionItems = @()
                    }
                }
                return $completionItems
            }
        }
    }

    # 5. 处理普通选项补全(子命令或参数，不区分大小写匹配)
    if ($currentOptions) {
        # 不区分大小写过滤匹配项(支持部分输入)
        $matchingOptions = $currentOptions | Where-Object {
            $_ -like "$wordToComplete*" -or 
            $_.ToLower() -like "$($wordToComplete.ToLower())*"
        }

        # 处理补全项格式(子命令加空格，参数直接返回)
        $completionItems = $matchingOptions | ForEach-Object {
            if ($_ -match '^-') {  # 参数(如 -force)
                $_
            } else {  # 子命令(如 start)
                "$_ "  # 补全后加空格，方便后续输入
            }
        }
    }

    return $completionItems
}

# 注册补全函数
Register-ArgumentCompleter -CommandName $commandName -ScriptBlock $scriptBlock
`
	PwshEnmuOtpsItem  = "	@{ Context = \"%s\"; Parameter = \"%s\"; Options = @(%s) }"                                         // 枚举选项条目
	PwshFlagParamItem = "	@{ Context = \"%s\"; Parameter = \"%s\"; ParamType = \"%s\"; ValueType = \"%s\"; Options = @(%s) }" // 标志参数条目(含枚举选项)
	PwshCmdTreeItem   = "	@{ Context = \"%s\"; Options = @(%s) }"                                                             // 命令树条目
)
