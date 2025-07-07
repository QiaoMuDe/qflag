package cmd

import (
	"fmt"
	"os"
)

// 支持的Shell类型
const (
	ShellBash = "bash" // bash shell
	// ShellZsh        = "zsh"        // zsh shell
	// ShellFish       = "fish"       // fish shell
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

// 内置自动补全命令的示例使用
var (
	// linux环境 临时启用
	linuxTempExample = ExampleInfo{
		Description: "Linux环境 临时启用",
		Usage:       fmt.Sprintf("source <(%s completion bash)", os.Args[0]),
	}

	// linux环境 永久启用
	linuxPermanentExample = ExampleInfo{
		Description: "Linux环境 永久启用(添加到~/.bashrc)",
		Usage:       fmt.Sprintf("echo \"source <(%s completion bash)\" >> ~/.bashrc", os.Args[0]),
	}

	// windows环境 临时启用
	windowsTempExample = ExampleInfo{
		Description: "Windows环境 临时启用",
		Usage:       fmt.Sprintf("%s completion powershell | Out-String | Invoke-Expression", os.Args[0]),
	}

	// windows环境 永久启用
	windowsPermanentExample = ExampleInfo{
		Description: "Windows环境 永久启用(添加到PowerShell配置文件)",
		Usage:       fmt.Sprintf("echo \"%s completion powershell | Out-String | Invoke-Expression\" >> $PROFILE", os.Args[0]),
	}
)

// BashCompletionTemplate 定义Bash自动补全脚本的模块化模板结构
// 用于统一管理补全脚本的各个组成部分
type BashCompletionTemplate struct {
	FunctionHeader string // 补全函数头部模板
	RootCommand    string // 根命令补全模板
	SubCommand     string // 子命令补全模板
	ShortOption    string // 短选项补全模板
	LongOption     string // 长选项补全模板
	BindingCommand string // 命令绑定模板
}

// bashCompletionTemplateString 定义Bash自动补全脚本的模块化模板字符串
var bashCompletionTemplate = BashCompletionTemplate{
	// 补全函数头部模板
	FunctionHeader: ``,
}

// generateBashCompletion 生成Bash自动补全脚本
//
// 返回值：
//   - string: Bash自动补全脚本
func (c *Cmd) generateBashCompletion() string {
	return ``
}
