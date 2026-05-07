# 新增 --install-completion 标志设计方案

## 背景

Windows 系统首次运行 Go 程序时启动较慢，导致用户在终端配置中直接调用 `--completion` 生成脚本的方式会造成终端启动卡顿。

## 目标

新增 `--install-completion` 标志，实现：
1. 自动生成补全脚本到用户家目录的 `.qflag_completions/` 目录
2. 自动将加载命令添加到对应 Shell 的配置文件中
3. 避免每次终端启动都执行 Go 程序

## 设计方案

### 1. 新增常量定义 (types/builtin.go)

```go
// 内置标志名称常量 - 新增安装补全标志
const (
    // ... 原有常量 ...
    
    // InstallCompletionFlagName 安装补全标志名称
    InstallCompletionFlagName = "install-completion"
)

// 新增内置标志类型
const (
    // ... 原有类型 ...
    
    // InstallCompletionFlag 安装补全标志
    // 用于自动生成补全脚本并配置到 Shell
    InstallCompletionFlag BuiltinFlagType = iota + 3
)

// 补全安装相关常量
const (
    // CompletionsDirName 补全脚本存放目录名
    CompletionsDirName = ".qflag_completions"
)
```

### 2. 新增处理器 (builtin/install_completion_handler.go)

```go
package builtin

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "strings"

    "gitee.com/MM-Q/qflag/internal/completion"
    "gitee.com/MM-Q/qflag/internal/types"
)

// InstallCompletionHandler 安装补全标志处理器
//
// 负责处理 --install-completion 标志，自动完成：
// 1. 创建 ~/.qflag_completions/ 目录
// 2. 生成补全脚本到该目录
// 3. 将加载命令添加到 Shell 配置文件
type InstallCompletionHandler struct{}

// Handle 处理安装补全标志
func (h *InstallCompletionHandler) Handle(cmd types.Command) error {
    shellType := getShellTypeFromArgs(cmd)
    programName := filepath.Base(os.Args[0])
    
    // 1. 获取用户家目录
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("无法获取用户家目录: %w", err)
    }
    
    // 2. 创建补全目录
    completionsDir := filepath.Join(homeDir, types.CompletionsDirName)
    if err := os.MkdirAll(completionsDir, 0755); err != nil {
        return fmt.Errorf("创建补全目录失败: %w", err)
    }
    
    // 3. 确定脚本文件名和路径
    scriptName := programName
    if shellType == types.PwshShell || shellType == types.PowershellShell {
        scriptName += ".ps1"
    } else {
        scriptName += ".sh"
    }
    scriptPath := filepath.Join(completionsDir, scriptName)
    
    // 4. 生成补全脚本到文件
    scriptContent, err := completion.Generate(cmd, shellType)
    if err != nil {
        return fmt.Errorf("生成补全脚本失败: %w", err)
    }
    
    if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
        return fmt.Errorf("写入补全脚本失败: %w", err)
    }
    
    // 5. 添加加载命令到配置文件
    if err := h.addLoadCommandToProfile(homeDir, scriptPath, shellType); err != nil {
        return fmt.Errorf("添加加载命令失败: %w", err)
    }
    
    // 6. 输出成功信息
    fmt.Printf("✓ 补全脚本已生成: %s\n", scriptPath)
    fmt.Printf("✓ 加载命令已添加到 %s\n", h.getProfilePath(homeDir, shellType))
    fmt.Println("\n请重新打开终端或执行以下命令使补全生效:")
    fmt.Printf("  source %s\n", h.getProfilePath(homeDir, shellType))
    
    os.Exit(0)
    return nil
}

// addLoadCommandToProfile 添加加载命令到配置文件
func (h *InstallCompletionHandler) addLoadCommandToProfile(homeDir, scriptPath, shellType string) error {
    profilePath := h.getProfilePath(homeDir, shellType)
    loadCommand := h.generateLoadCommand(scriptPath, shellType)
    
    // 检查是否已存在
    content, err := os.ReadFile(profilePath)
    if err == nil && strings.Contains(string(content), scriptPath) {
        // 已存在，不需要重复添加
        return nil
    }
    
    // 追加到配置文件
    f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()
    
    // 添加换行和注释
    if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
        f.WriteString("\n")
    }
    f.WriteString(fmt.Sprintf("\n# qflag completion for %s\n", filepath.Base(os.Args[0])))
    f.WriteString(loadCommand + "\n")
    
    return nil
}

// getProfilePath 获取配置文件路径
func (h *InstallCompletionHandler) getProfilePath(homeDir, shellType string) string {
    switch shellType {
    case types.PwshShell, types.PowershellShell:
        // PowerShell 配置文件
        if runtime.GOOS == "windows" {
            return filepath.Join(homeDir, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
        }
        return filepath.Join(homeDir, ".config", "powershell", "Microsoft.PowerShell_profile.ps1")
    default:
        // Bash 配置文件
        if runtime.GOOS == "darwin" {
            return filepath.Join(homeDir, ".bash_profile")
        }
        return filepath.Join(homeDir, ".bashrc")
    }
}

// generateLoadCommand 生成加载命令
func (h *InstallCompletionHandler) generateLoadCommand(scriptPath, shellType string) string {
    switch shellType {
    case types.PwshShell, types.PowershellShell:
        // PowerShell: 使用变量避免重复路径，检查文件存在
        return fmt.Sprintf("$__qflag_comp = '%s'; if (Test-Path $__qflag_comp) { . $__qflag_comp }", scriptPath)
    default:
        // Bash: 使用 -f 检查文件存在
        return fmt.Sprintf("[ -f '%s' ] && source '%s'", scriptPath, scriptPath)
    }
}

// Type 返回标志类型
func (h *InstallCompletionHandler) Type() types.BuiltinFlagType {
    return types.InstallCompletionFlag
}

// ShouldRegister 判断是否应该注册
func (h *InstallCompletionHandler) ShouldRegister(cmd types.Command) bool {
    // 只在根命令注册，且需要启用补全功能
    return cmd.IsRootCmd() && cmd.Config().Completion
}

// ShouldSkipRegistration 判断是否应该跳过注册
func (h *InstallCompletionHandler) ShouldSkipRegistration(cmd types.Command) bool {
    _, exists := cmd.GetFlag(types.InstallCompletionFlagName)
    return exists
}
```

### 3. 修改 completion 包接口

需要在 `internal/completion/completion.go` 中新增 `Generate` 函数：

```go
// Generate 生成补全脚本内容
//
// 参数:
//   - cmd: 命令实例
//   - shellType: Shell 类型
//
// 返回值:
//   - string: 补全脚本内容
//   - error: 生成错误
func Generate(cmd *Cmd, shellType string) (string, error) {
    // 根据 shellType 调用对应的生成函数
    // 返回脚本内容而不是直接打印
}
```

### 4. 注册新处理器 (builtin/manager.go)

在管理器中注册新的处理器：

```go
func init() {
    defaultManager = NewManager()
    // 注册内置处理器
    defaultManager.Register(&HelpHandler{})
    defaultManager.Register(&VersionHandler{})
    defaultManager.Register(&CompletionHandler{})
    defaultManager.Register(&InstallCompletionHandler{}) // 新增
}
```

### 5. 修改帮助信息 (types/help.go)

新增安装补全的示例：

```go
// 补全安装示例信息 - Windows 中文
var HelpInstallCompletionExampleWinCN = map[string]string{
    "安装补全脚本": fmt.Sprintf("%s --install-completion pwsh", filepath.Base(os.Args[0])),
}

// 补全安装示例信息 - Linux 中文
var HelpInstallCompletionExampleLinuxCN = map[string]string{
    "安装补全脚本": fmt.Sprintf("%s --install-completion bash", filepath.Base(os.Args[0])),
}

// 补全安装示例信息 - macOS 中文
var HelpInstallCompletionExampleMacCN = map[string]string{
    "安装补全脚本(Bash)": fmt.Sprintf("%s --install-completion bash", filepath.Base(os.Args[0])),
    "安装补全脚本(Zsh)":  fmt.Sprintf("%s --install-completion zsh", filepath.Base(os.Args[0])),
}
```

## 使用方式

### 一次性安装

```bash
# Windows PowerShell
yourapp --install-completion pwsh

# Linux Bash
yourapp --install-completion bash

# macOS Bash
yourapp --install-completion bash

# macOS Zsh
yourapp --install-completion zsh
```

### 输出示例

```
✓ 补全脚本已生成: /home/user/.qflag_completions/yourapp.sh
✓ 加载命令已添加到 /home/user/.bashrc

请重新打开终端或执行以下命令使补全生效:
  source /home/user/.bashrc
```

## 目录结构

```
~/.qflag_completions/
├── yourapp.sh      # Bash 补全脚本
├── yourapp.ps1     # PowerShell 补全脚本
└── anotherapp.sh   # 其他应用补全脚本
```

## 优势

1. **解决性能问题**：终端启动不再调用 Go 程序
2. **一键安装**：用户只需执行一条命令
3. **自动配置**：自动检测 Shell 类型并配置
4. **安全加载**：配置文件中使用条件判断，文件不存在不报错
5. **统一管理**：所有 qflag 应用的补全脚本集中管理

## 实现步骤

1. **新增常量** (`types/builtin.go`)
   - `InstallCompletionFlagName`
   - `InstallCompletionFlag`
   - `CompletionsDirName`

2. **新增处理器** (`builtin/install_completion_handler.go`)
   - 实现 `InstallCompletionHandler` 结构体
   - 实现所有接口方法

3. **修改 completion 包** (`completion/completion.go`)
   - 新增 `Generate` 函数

4. **注册处理器** (`builtin/manager.go`)
   - 在 `init` 函数中注册

5. **更新帮助信息** (`types/help.go`)
   - 新增安装补全示例

6. **测试验证**
   - 各平台测试安装流程
   - 验证终端启动速度
   - 验证补全功能正常
