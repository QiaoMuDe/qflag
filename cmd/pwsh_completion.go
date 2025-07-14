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
	// 格式化选项为PowerShell数组格式
	var optsBuf bytes.Buffer
	for i, opt := range cmdOpts {
		if i > 0 {
			// 非首个元素前添加逗号和空格
			optsBuf.WriteString(", ")
		}
		optsBuf.WriteString(fmt.Sprintf("'%s'", escapePwshString(opt)))
	}
	fmt.Fprintf(cmdTreeEntries, PwshCmdTreeItem+",\n", cmdPath, optsBuf.String())
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
	var flagParamsBuf bytes.Buffer
	var enumOptionsBuf bytes.Buffer

	// 处理根命令选项
	var rootOptsBuf bytes.Buffer
	for i, opt := range rootCmdOpts {
		if i > 0 {
			// 非首个元素前添加逗号和空格
			rootOptsBuf.WriteString(", ")
		}
		rootOptsBuf.WriteString(fmt.Sprintf("'%s'", escapePwshString(opt)))
	}

	// 处理标志参数
	for _, param := range params {
		// 生成标志参数条目
		fmt.Fprintf(&flagParamsBuf, PwshFlagParamItem+",\n",
			param.CommandPath,
			param.Name,
			param.Type,
			param.ValueType)

		// 处理枚举类型选项
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			var optionsBuf bytes.Buffer
			for i, opt := range param.EnumOptions {
				if i > 0 {
					// 非首个元素前添加逗号和空格
					optionsBuf.WriteString(", ")
				}
				optionsBuf.WriteString(fmt.Sprintf("'%s'", escapePwshString(opt)))
			}
			fmt.Fprintf(&enumOptionsBuf, PwshEnmuOtpsItem+",\n",
				param.CommandPath,
				param.Name,
				optionsBuf.String())
		}
	}

	// 写入PowerShell自动补全脚本
	fmt.Fprintf(buf, PwshFunctionHeader,
		programName,
		// 根命令树条目
		fmt.Sprintf(PwshCmdTreeItem+",\n", "/", rootOptsBuf.String())+cmdTreeEntries,
		flagParamsBuf.String(),
		enumOptionsBuf.String())
}

// escapePwshString 转义PowerShell字符串中的特殊字符
func escapePwshString(s string) string {
	// 替换单引号为'' (PowerShell转义方式)
	escaped := strings.ReplaceAll(s, "'", "''")
	// 替换反斜杠为\
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
	return escaped
}

const (
	// PowerShell补全模板
	PwshFunctionHeader = `# -------------------------- 配置区域(需根据实际命令修改) --------------------------
# 命令名称(替换为你的实际命令名，如"mycmd")
$commandName = "%s"

# 1. 命令树定义(数组形式，每个元素包含 Context 和 Options)
# Context: 上下文路径(如"/"、"/init/");Options: 该层级可用选项数组(不区分大小写)
$cmdTree = @(
%s  )

# 2. 标志参数定义(数组形式，每个元素包含 Context、Parameter、ParamType、ValueType)
# ParamType: required(必须带值)、optional(可选值);ValueType: path|number|ip|enum|url等
$flagParams = @(
%s  )

# 3. 枚举选项定义(数组形式，每个元素包含 Context、Parameter、Options)
$enumOptions = @(
%s  )

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
                        # 查找枚举选项(不区分大小写)
                        $enumDef = $enumOptions | Where-Object {
                            $_.Context -eq $context -and
                            ($_.Parameter -eq $prevElement -or $_.Parameter -eq $prevElement.ToLower() -or $_.Parameter -eq $prevElement.ToUpper())
                        }
                        if ($enumDef) {
                            $completionItems = $enumDef.Options | Where-Object { $_ -like "$wordToComplete*" }
                        }
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
            }
            else {  # 子命令(如 start)
                "$_ "  # 补全后加空格，方便后续输入
            }
        }
    }

    return $completionItems
}

# 注册补全函数
Register-ArgumentCompleter -CommandName $commandName -ScriptBlock $scriptBlock
`
	PwshEnmuOtpsItem  = "\t@{ Context = \"%s\"; Parameter = \"%s\"; Options = @(%s) }"                        // 枚举选项条目
	PwshFlagParamItem = "\t@{ Context = \"%s\"; Parameter = \"%s\"; ParamType = \"%s\"; ValueType = \"%s\" }" // 标志参数条目
	PwshCmdTreeItem   = "\t@{ Context = \"%s\"; Options = @(%s) }"                                            // 命令树条目
)
