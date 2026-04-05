# qflag 动态补全详细实施方案

## 1. 各包修改清单

### 1.1 internal/types/config.go

添加动态补全配置字段：

```go
// CmdConfig 命令配置
type CmdConfig struct {
    // ... 现有字段 ...
    
    Completion     bool              // 是否启用自动补全标志
    DynamicCompletion bool           // 【新增】是否启用动态补全
    
    // ... 其他字段 ...
}
```

### 1.2 internal/cmd/cmdopts.go

添加动态补全配置方法：

```go
// CmdOptions 命令选项配置结构体
type CmdOptions struct {
    // ... 现有字段 ...
    
    Completion  bool   // 是否启用自动补全标志
    DynamicCompletion bool  // 【新增】是否启用动态补全
    
    // ... 其他字段 ...
}

// apply 应用选项到配置
func (o *CmdOptions) apply(c *types.CmdConfig) {
    // ... 现有逻辑 ...
    
    c.Completion = o.Completion
    c.DynamicCompletion = o.DynamicCompletion  // 【新增】
}
```

### 1.3 internal/cmd/cmd.go

添加启用动态补全的方法：

```go
// SetDynamicCompletion 设置是否启用动态补全
//
// 参数:
//   - enable: 是否启用动态补全
func (c *Cmd) SetDynamicCompletion(enable bool) {
    c.config.DynamicCompletion = enable
}

// SetCompletion 设置是否启用自动补全标志
//
// 参数:
//   - enable: 是否启用自动补全
//
// 注意:
//   - 此设置与动态补全完全独立
//   - 启用补全不会自动启用动态补全
func (c *Cmd) SetCompletion(enable bool) {
    c.config.Completion = enable
}
```

### 1.4 internal/types/builtin.go

新增动态补全相关类型和常量：

```go
const (
    HelpFlag BuiltinFlagType = iota
    VersionFlag
    CompletionFlag
    // 【新增】不需要新增 BuiltinFlagType，复用 CompletionFlag
)

// 【新增】动态补全相关常量
const (
    // CompModeFlagName 补全模式标志名称
    CompModeFlagName = "comp-mode"
    
    // CompModeStatic 静态补全模式
    CompModeStatic = "static"
    
    // CompModeDynamic 动态补全模式
    CompModeDynamic = "dynamic"
)

// 【新增】动态补全函数类型
type DynamicCompleter func(ctx CompletionContext) CompletionResult

// 【新增】补全上下文
type CompletionContext struct {
    Command        Command  // 当前命令
    WordToComplete string   // 当前输入的词
    PreviousWord   string   // 前一个词
    Args           []string // 所有参数
    CommandPath    string   // 命令路径
}

// 【新增】补全项
type Completion struct {
    Value       string
    Description string
}

// 【新增】补全结果
type CompletionResult struct {
    Completions []Completion
    Directive   int  // 指令位图
}

// 【新增】指令常量
const (
    CompDirectiveError = 1 << iota
    CompDirectiveNoSpace
    CompDirectiveNoFileComp
    CompDirectiveFilterFileExt
    CompDirectiveFilterDirs
    CompDirectiveKeepOrder
)
```

### 1.5 internal/cmd/cmd.go（Cmd 结构体）

添加动态补全函数注册表：

```go
// Cmd 命令结构体
type Cmd struct {
    // ... 现有字段 ...
    
    completers map[string]types.DynamicCompleter  // 【新增】动态补全函数注册表
}

// NewCmd 创建新命令时初始化 completers
func NewCmd(name, desc string, errorHandling ErrorHandling) *Cmd {
    // ... 现有逻辑 ...
    
    cmd.completers = make(map[string]types.DynamicCompleter)  // 【新增】
    
    return cmd
}

// RegisterCompleter 注册动态补全函数
//
// 参数:
//   - flagName: 标志名称（如 "env"）
//   - completer: 补全函数
func (c *Cmd) RegisterCompleter(flagName string, completer types.DynamicCompleter) {
    if c.completers == nil {
        c.completers = make(map[string]types.DynamicCompleter)
    }
    c.completers[flagName] = completer
}

// GetCompleter 获取动态补全函数
func (c *Cmd) GetCompleter(flagName string) (types.DynamicCompleter, bool) {
    completer, ok := c.completers[flagName]
    return completer, ok
}

// HasDynamicCompletion 是否启用了动态补全
func (c *Cmd) HasDynamicCompletion() bool {
    return c.config.DynamicCompletion
}
```

### 1.6 internal/builtin/manager.go

修改 RegisterBuiltinFlags，条件注册 comp-mode 标志：

```go
func (m *BuiltinFlagManager) RegisterBuiltinFlags(cmd types.Command) error {
    for _, handler := range m.handlers {
        if !handler.ShouldRegister(cmd) {
            continue
        }
        
        if handler.ShouldSkipRegistration(cmd) {
            continue
        }
        
        switch handler.Type() {
        // ... 现有 case ...
        
        case types.CompletionFlag:
            // 注册原有的 completion 标志
            // ... 现有逻辑 ...
            
            // 【新增】如果启用了动态补全，同时注册 comp-mode 标志
            if cmd.Config().DynamicCompletion {
                if err := m.registerCompModeFlag(cmd); err != nil {
                    return err
                }
            }
        }
    }
    
    return nil
}

// 【新增】注册 comp-mode 标志
func (m *BuiltinFlagManager) registerCompModeFlag(cmd types.Command) error {
    // 检查是否已存在
    if _, exists := cmd.GetFlag(types.CompModeFlagName); exists {
        return nil
    }
    
    var desc string
    if cmd.Config().UseChinese {
        desc = fmt.Sprintf("补全模式: %s(默认) 或 %s", types.CompModeStatic, types.CompModeDynamic)
    } else {
        desc = fmt.Sprintf("Completion mode: %s(default) or %s", types.CompModeStatic, types.CompModeDynamic)
    }
    
    modeFlag := flag.NewEnumFlag(
        types.CompModeFlagName, 
        "", 
        desc, 
        types.CompModeStatic,
        []string{types.CompModeStatic, types.CompModeDynamic},
    )
    
    return cmd.AddFlag(modeFlag)
}
```

### 1.7 internal/builtin/completion_handler.go

修改 Handle 方法，调用 completion 包的 API：

```go
func (h *CompletionHandler) Handle(cmd types.Command) error {
    shellType := getShellTypeFromArgs(cmd)
    
    // 【新增】检查是否启用了动态补全
    if cmd.Config().DynamicCompletion {
        mode := getCompModeFromArgs(cmd)
        
        if mode == types.CompModeDynamic {
            // 获取位置参数（-- 之后的参数）
            posArgs := getPositionalArgs(cmd)
            
            if len(posArgs) > 0 {
                // 【调用 completion 包 API】执行动态补全
                result := completion.ExecuteDynamic(cmd, shellType, posArgs)
                completion.PrintResult(result)
                os.Exit(0)
                return nil
            }
            
            // 【调用 completion 包 API】生成动态补全脚本
            completion.GenDynamicAndPrint(cmd, shellType)
            os.Exit(0)
            return nil
        }
    }
    
    // 原有逻辑：生成静态补全脚本
    completion.GenAndPrint(cmd, shellType)
    os.Exit(0)
    return nil
}

// 【新增】获取 comp-mode 标志值
func getCompModeFromArgs(cmd types.Command) string {
    if f, ok := cmd.GetFlag(types.CompModeFlagName); ok {
        return f.GetStr()
    }
    return types.CompModeStatic
}

// 【新增】获取位置参数（-- 之后的参数）
func getPositionalArgs(cmd types.Command) []string {
    args := cmd.Args()
    for i, arg := range args {
        if arg == "--" && i+1 < len(args) {
            return args[i+1:]
        }
    }
    return nil
}
```

### 1.8 internal/completion/dynamic.go（【核心】新增文件）

动态补全的核心实现，提供 API 供 builtin 调用：

```go
package completion

import (
    "bytes"
    "fmt"
    "os"
    "path/filepath"
    "text/template"
    
    "gitee.com/MM-Q/qflag/internal/types"
)

// ========================================
// 对外 API（供 builtin/completion_handler 调用）
// ========================================

// GenDynamicAndPrint 生成并打印动态补全脚本
//
// 参数:
//   - cmd: 命令实例
//   - shellType: Shell 类型
func GenDynamicAndPrint(cmd types.Command, shellType string) {
    script, err := GenerateDynamic(cmd, shellType)
    if err != nil {
        fmt.Printf("Error generating dynamic completion script: %v\n", err)
        return
    }
    fmt.Println(script)
}

// ExecuteDynamic 执行动态补全
//
// 参数:
//   - cmd: 命令实例
//   - shellType: Shell 类型
//   - args: 位置参数（用户输入的命令行参数）
//
// 返回值:
//   - CompletionResult: 补全结果
func ExecuteDynamic(cmd types.Command, shellType string, args []string) types.CompletionResult {
    // 构建补全上下文
    ctx := buildCompletionContext(cmd, args)
    
    // 查找补全函数
    if completer, ok := cmd.GetCompleter(ctx.CurrentFlag); ok {
        // 使用注册的自定义补全函数
        return completer(ctx)
    }
    
    // 使用默认补全逻辑
    return getDefaultCompletions(cmd, ctx)
}

// PrintResult 打印补全结果（供 shell 脚本解析）
//
// 参数:
//   - result: 补全结果
//
// 输出格式:
//   候选值1
//   候选值2
//   :指令
func PrintResult(result types.CompletionResult) {
    for _, comp := range result.Completions {
        if comp.Description != "" {
            fmt.Printf("%s\t%s\n", comp.Value, comp.Description)
        } else {
            fmt.Println(comp.Value)
        }
    }
    fmt.Printf(":%d\n", result.Directive)
}

// ========================================
// 内部实现
// ========================================

// buildCompletionContext 构建补全上下文
func buildCompletionContext(cmd types.Command, args []string) types.CompletionContext {
    ctx := types.CompletionContext{
        Command: cmd,
        Args:    args,
    }
    
    // 解析命令路径
    ctx.CommandPath = parseCommandPath(cmd, args)
    
    // 确定当前补全状态
    ctx.WordToComplete, ctx.PreviousWord, ctx.CurrentFlag, ctx.CompletingFlagValue = 
        parseCompletionState(cmd, args)
    
    return ctx
}

// parseCommandPath 解析命令路径
func parseCommandPath(cmd types.Command, args []string) string {
    path := "/"
    for _, arg := range args {
        if strings.HasPrefix(arg, "-") {
            break
        }
        if subCmd, ok := cmd.GetSubCmd(arg); ok {
            path += arg + "/"
            cmd = subCmd
        } else {
            break
        }
    }
    return path
}

// parseCompletionState 解析补全状态
func parseCompletionState(cmd types.Command, args []string) (wordToComplete, previousWord, currentFlag string, completingFlagValue bool) {
    if len(args) == 0 {
        return "", "", "", false
    }
    
    // 最后一个参数是当前正在补全的词
    wordToComplete = args[len(args)-1]
    
    // 前一个词
    if len(args) > 1 {
        previousWord = args[len(args)-2]
    }
    
    // 检查是否正在补全标志值
    if strings.HasPrefix(previousWord, "-") {
        currentFlag = previousWord
        completingFlagValue = true
    }
    
    return
}

// getDefaultCompletions 获取默认补全（子命令、标志、枚举值等）
func getDefaultCompletions(cmd types.Command, ctx types.CompletionContext) types.CompletionResult {
    var completions []types.Completion
    var directive int
    
    if ctx.CompletingFlagValue {
        // 补全标志值
        completions = getFlagValueCompletions(cmd, ctx)
    } else if strings.HasPrefix(ctx.WordToComplete, "-") {
        // 补全标志名
        completions = getFlagCompletions(cmd, ctx)
    } else {
        // 补全子命令
        completions = getSubCommandCompletions(cmd, ctx)
    }
    
    return types.CompletionResult{
        Completions: completions,
        Directive:   directive,
    }
}

// getFlagValueCompletions 获取标志值补全
func getFlagValueCompletions(cmd types.Command, ctx types.CompletionContext) []types.Completion {
    var completions []types.Completion
    
    // 获取标志定义
    if flag, ok := cmd.GetFlag(ctx.CurrentFlag); ok {
        // 根据标志类型提供补全
        switch flag.ValueType() {
        case "enum":
            // 枚举类型：返回枚举值
            for _, opt := range flag.EnumOptions() {
                if strings.HasPrefix(opt, ctx.WordToComplete) {
                    completions = append(completions, types.Completion{Value: opt})
                }
            }
        case "path", "string":
            // 路径/字符串类型：返回文件路径（由 shell 处理）
            // 不返回具体补全，设置 NoFileComp 指令
        }
    }
    
    return completions
}

// getFlagCompletions 获取标志名补全
func getFlagCompletions(cmd types.Command, ctx types.CompletionContext) []types.Completion {
    var completions []types.Completion
    
    for _, flag := range cmd.Flags() {
        name := flag.LongName()
        if strings.HasPrefix(name, ctx.WordToComplete) {
            completions = append(completions, types.Completion{
                Value:       name,
                Description: flag.Desc(),
            })
        }
    }
    
    return completions
}

// getSubCommandCompletions 获取子命令补全
func getSubCommandCompletions(cmd types.Command, ctx types.CompletionContext) []types.Completion {
    var completions []types.Completion
    
    for _, subCmd := range cmd.SubCmds() {
        name := subCmd.Name()
        if strings.HasPrefix(name, ctx.WordToComplete) {
            completions = append(completions, types.Completion{
                Value:       name,
                Description: subCmd.Desc(),
            })
        }
    }
    
    return completions
}

// ========================================
// 动态补全脚本生成
// ========================================

//go:embed templates/bash_dynamic.tmpl
var bashDynamicTemplate string

//go:embed templates/pwsh_dynamic.tmpl
var pwshDynamicTemplate string

// GenerateDynamic 生成动态补全脚本
func GenerateDynamic(cmd types.Command, shellType string) (string, error) {
    programName := filepath.Base(os.Args[0])
    
    switch shellType {
    case types.BashShell:
        return generateDynamicBash(programName)
    case types.PwshShell, types.PowershellShell:
        return generateDynamicPwsh(programName)
    default:
        return "", fmt.Errorf("unsupported shell type '%s'", shellType)
    }
}

func generateDynamicBash(programName string) (string, error) {
    tmpl, err := template.New("bash_dynamic").Parse(bashDynamicTemplate)
    if err != nil {
        return "", err
    }
    
    var buf bytes.Buffer
    data := map[string]string{
        "ProgramName": programName,
    }
    
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    
    return buf.String(), nil
}

func generateDynamicPwsh(programName string) (string, error) {
    tmpl, err := template.New("pwsh_dynamic").Parse(pwshDynamicTemplate)
    if err != nil {
        return "", err
    }
    
    var buf bytes.Buffer
    data := map[string]string{
        "ProgramName": programName,
    }
    
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    
    return buf.String(), nil
}
```

### 1.9 internal/completion/templates/（新增模板文件）

新增两个模板文件（使用 `//go:embed` 嵌入）：
- `bash_dynamic.tmpl` - Bash 动态补全脚本模板
- `pwsh_dynamic.tmpl` - PowerShell 动态补全脚本模板

#### Bash 动态补全脚本模板 (bash_dynamic.tmpl)

```bash
# {{.ProgramName}} dynamic completion script
# Generated by qflag

_{{.ProgramName}}_dynamic_completion() {
    local cur prev words cword
    _init_completion || return

    # 获取当前命令行内容（从第一个参数开始）
    local cmd_line="${words[*]}"
    
    # 调用程序获取补全结果
    # 格式: program --completion=bash --comp-mode=dynamic -- <当前命令行>
    local completions
    completions=$({{.ProgramName}} --completion=bash --comp-mode=dynamic -- ${words[@]} 2>/dev/null)
    
    # 解析结果（最后一行是 :指令）
    local directive
    directive=$(echo "$completions" | tail -n 1)
    
    # 提取候选值（去掉最后一行指令）
    local candidates
    candidates=$(echo "$completions" | head -n -1)
    
    # 设置补全结果
    COMPREPLY=($(compgen -W "$candidates" -- "$cur"))
    
    return 0
}

# 注册补全函数
complete -F _{{.ProgramName}}_dynamic_completion {{.ProgramName}}
```

#### PowerShell 动态补全脚本模板 (pwsh_dynamic.tmpl)

```powershell
# {{.ProgramName}} dynamic completion script
# Generated by qflag

Register-ArgumentCompleter -Native -CommandName {{.ProgramName}} -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    
    # 获取完整的命令行
    $cmdLine = $commandAst.ToString()
    $words = $commandAst.CommandElements | ForEach-Object { $_.ToString() }
    
    # 调用程序获取补全结果
    # 格式: program --completion=pwsh --comp-mode=dynamic -- <当前命令行>
    $completions = & {{.ProgramName}} --completion=pwsh --comp-mode=dynamic -- @words 2>$null
    
    if (-not $completions) {
        return
    }
    
    # 解析结果（最后一行是 :指令）
    $lines = $completions -split "`n"
    $directive = $lines[-1]
    $candidates = $lines[0..($lines.Count - 2)]
    
    # 返回候选值
    foreach ($candidate in $candidates) {
        if ($candidate -match "^`$wordToComplete") {
            # 处理带描述的格式: value\tdescription
            $parts = $candidate -split "`t", 2
            [System.Management.Automation.CompletionResult]::new(
                $parts[0],
                $parts[0],
                'ParameterValue',
                if ($parts[1]) { $parts[1] } else { $parts[0] }
            )
        }
    }
}
```

#### internal/completion/dynamic.go 中的模板引用

```go
package completion

import (
    "_embed"
    // ... 其他导入
)

//go:embed templates/bash_dynamic.tmpl
var bashDynamicTemplate string

//go:embed templates/pwsh_dynamic.tmpl  
var pwshDynamicTemplate string

// GenerateDynamic 生成动态补全脚本
func GenerateDynamic(cmd types.Command, shellType string) (string, error) {
    programName := filepath.Base(os.Args[0])
    
    switch shellType {
    case types.BashShell:
        return generateFromTemplate(bashDynamicTemplate, programName)
    case types.PwshShell, types.PowershellShell:
        return generateFromTemplate(pwshDynamicTemplate, programName)
    default:
        return "", fmt.Errorf("unsupported shell type '%s'", shellType)
    }
}

func generateFromTemplate(templateStr, programName string) (string, error) {
    tmpl, err := template.New("dynamic").Parse(templateStr)
    if err != nil {
        return "", err
    }
    
    var buf bytes.Buffer
    data := map[string]string{
        "ProgramName": programName,
    }
    
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    
    return buf.String(), nil
}
```

## 2. 用户 API

```go
package main

import (
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    cmd := qflag.NewCmd("myapp", "示例应用", qflag.ExitOnError)
    
    // 【新增】启用动态补全（会注册 --comp-mode 标志）
    cmd.SetDynamicCompletion(true)
    
    // 或者通过选项启用
    // cmd := qflag.NewCmdWithOptions("myapp", "", qflag.CmdOptions{
    //     DynamicCompletion: true,
    // })
    
    // 添加子命令
    deploy := cmd.SubCmd("deploy", "部署应用")
    deploy.String("env", "", "部署环境")
    
    // 【新增】注册自定义补全函数
    cmd.RegisterCompleter("env", func(ctx qflag.CompletionContext) qflag.CompletionResult {
        return qflag.CompletionResult{
            Completions: []qflag.Completion{
                {Value: "dev", Description: "开发环境"},
                {Value: "prod", Description: "生产环境"},
            },
            Directive: qflag.CompDirectiveNoSpace,
        }
    })
    
    cmd.ParseAndRoute(os.Args[1:])
}
```

## 3. 使用方式与场景

### 场景 1: 首次运行生成补全脚本

用户首次使用时，需要生成并安装补全脚本：

```bash
# 生成动态补全脚本并保存到文件
myapp --completion=bash --comp-mode=dynamic > /etc/bash_completion.d/myapp

# 或手动加载
source <(myapp --completion=bash --comp-mode=dynamic)
```

**程序内部处理流程：**

```
用户输入: myapp --completion=bash --comp-mode=dynamic
         ↓
    builtin/completion_handler.go
         ↓
    检测到 --comp-mode=dynamic 且没有位置参数
         ↓
    调用 completion.GenDynamicAndPrint(cmd, "bash")
         ↓
    internal/completion/dynamic.go
         ↓
    生成动态补全脚本（包含调用程序的命令）
         ↓
    输出脚本内容到 stdout
         ↓
    用户将脚本保存到系统目录或 source 加载
```

### 场景 2: 补全时调用程序

用户在 Shell 中按 Tab 键时，Shell 脚本会调用程序获取补全结果：

```bash
# 用户在 Shell 中输入
myapp deploy --env <tab>

# Shell 脚本内部调用程序
myapp --completion=bash --comp-mode=dynamic -- deploy --env ""
```

**程序内部处理流程：**

```
用户按 Tab 键
         ↓
    Shell 脚本执行:
    myapp --completion=bash --comp-mode=dynamic -- deploy --env ""
         ↓
    builtin/completion_handler.go
         ↓
    检测到 --comp-mode=dynamic 且有位置参数 ("deploy", "--env", "")
         ↓
    调用 completion.ExecuteDynamic(cmd, "bash", ["deploy", "--env", ""])
         ↓
    internal/completion/dynamic.go
         ↓
    1. 构建 CompletionContext
       - CommandPath: "/deploy"
       - WordToComplete: "" (当前输入为空)
       - PreviousWord: "--env"
       - CurrentFlag: "--env"
       - CompletingFlagValue: true
         ↓
    2. 查找注册的补全函数
       - 如果注册了 "env" 的补全函数，调用它
       - 否则使用默认补全逻辑
         ↓
    3. 返回 CompletionResult
       Completions: [{"dev", "开发环境"}, {"prod", "生产环境"}]
       Directive: 0
         ↓
    调用 completion.PrintResult(result)
         ↓
    输出格式:
    dev	开发环境
    prod	生产环境
    :0
         ↓
    Shell 脚本解析输出，显示补全候选
```

### 完整使用流程

```bash
# ========== 第一步：生成并安装补全脚本（只需执行一次）==========

# Bash 用户
myapp --completion=bash --comp-mode=dynamic > ~/.bash_completion.d/myapp
source ~/.bash_completion.d/myapp

# PowerShell 用户  
myapp --completion=pwsh --comp-mode=dynamic > $PROFILE
. $PROFILE

# ========== 第二步：使用补全（日常操作）==========

# 补全子命令
myapp <tab>
# 显示: deploy, status, logs

# 补全标志
myapp deploy --<tab>
# 显示: --env, --version, --region

# 补全标志值（调用自定义补全函数）
myapp deploy --env <tab>
# 显示: dev, prod, staging

# 补全标志值（枚举类型）
myapp deploy --region <tab>
# 显示: cn-north-1, cn-south-1
```

## 4. 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                  builtin/completion_handler.go               │
│  - 检测 --comp-mode 标志                                     │
│  - 调用 completion 包 API                                    │
│     • GenDynamicAndPrint()  生成脚本                         │
│     • ExecuteDynamic()      执行补全                         │
│     • PrintResult()         输出结果                         │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  internal/completion/dynamic.go              │
│  【核心实现】                                                 │
│  - buildCompletionContext()  构建补全上下文                  │
│  - ExecuteDynamic()          执行动态补全                    │
│  - getDefaultCompletions()   默认补全逻辑                    │
│  - GenerateDynamic()         生成动态脚本                    │
└─────────────────────────────────────────────────────────────┘
```

## 5. 修改的文件清单

1. **internal/types/config.go** - 添加 `DynamicCompletion bool` 字段
2. **internal/types/builtin.go** - 添加动态补全相关类型和常量
3. **internal/cmd/cmdopts.go** - 添加 `DynamicCompletion` 选项字段
4. **internal/cmd/cmd.go** - 添加 `SetDynamicCompletion()` 和 `RegisterCompleter()` 方法
5. **internal/builtin/manager.go** - 条件注册 `comp-mode` 标志
6. **internal/builtin/completion_handler.go** - 调用 completion 包 API
7. **internal/completion/dynamic.go** - 【核心】动态补全实现
8. **internal/completion/templates/*.tmpl** - 动态补全脚本模板

## 6. 关键设计点

- **completion 包是核心**：所有动态补全逻辑在 `internal/completion/dynamic.go` 实现
- **builtin 只负责调用**：`completion_handler.go` 只检测标志并调用 API
- **清晰的 API 边界**：
  - `GenDynamicAndPrint()` - 生成脚本
  - `ExecuteDynamic()` - 执行补全
  - `PrintResult()` - 输出结果
