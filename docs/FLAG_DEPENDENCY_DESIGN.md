# FlagDependency 标志依赖关系设计方案

## 一、设计概述

### 1.1 背景与动机

现有的互斥组（MutexGroup）和必需组（RequiredGroup）解决了**组内标志之间**的约束关系，但无法表达**特定标志触发对其他标志的限制**这一常见场景。

### 1.2 核心需求

实现标志级别的条件依赖关系：
- **互斥依赖**：如果使用了标志A，则标志B、C不能使用
- **必需依赖**：如果使用了标志A，则标志B、C必须设置

### 1.3 使用场景

```bash
# 场景1: 使用 --remote 时，--local-path 和 --local-config 不能使用
myapp deploy --remote --local-path=/tmp  # ❌ 错误

# 场景2: 使用 --ssl 时，--cert 和 --key 必须同时设置
myapp server --ssl --cert=server.crt  # ❌ 错误：缺少 --key

# 场景3: 使用 --config 时，其他配置标志失效
myapp run --config=app.yaml --port=8080  # ❌ 错误：--config 与 --port 互斥
```

---

## 二、核心设计

### 2.1 类型定义

```go
// DepType 依赖关系类型
type DepType int

const (
	// DepMutex 互斥依赖
	// 当触发标志被设置时，目标标志不能被设置
	DepMutex DepType = iota

	// DepRequired 必需依赖
	// 当触发标志被设置时，所有目标标志必须被设置
	DepRequired
)

// String 返回依赖类型的字符串表示
func (d DepType) String() string {
	switch d {
	case DepMutex:
		return "mutex"
	case DepRequired:
		return "required"
	default:
		return "unknown"
	}
}

// FlagDependency 标志依赖关系定义
//
// FlagDependency 定义了当某个标志（触发标志）被设置时，
// 对其他标志（目标标志）的约束条件。
//
// 字段说明:
//   - Name: 依赖关系名称，用于错误提示和标识
//   - Trigger: 触发标志的名称，当此标志被设置时触发依赖检查
//   - Targets: 目标标志名称列表，这些标志会受到约束
//   - Type: 依赖关系类型（互斥或必需）
//
// 使用场景:
//   - 远程模式与本地路径互斥 (trigger="remote", targets=["local-path"], type=DepMutex)
//   - SSL模式需要证书和密钥 (trigger="ssl", targets=["cert","key"], type=DepRequired)
//   - 配置文件模式与其他配置互斥 (trigger="config", targets=["port","host"], type=DepMutex)
type FlagDependency struct {
	Name    string   // 依赖关系名称，用于错误提示和标识
	Trigger string   // 触发标志名称
	Targets []string // 目标标志名称列表
	Type    DepType  // 依赖关系类型
}
```

### 2.2 CmdConfig 扩展

```go
// CmdConfig 命令配置类型
type CmdConfig struct {
	Version           string            // 版本号
	UseChinese        bool              // 是否使用中文
	EnvPrefix         string            // 环境变量前缀
	UsageSyntax       string            // 命令使用语法
	Example           map[string]string // 示例使用
	Notes             []string          // 注意事项
	LogoText          string            // 命令logo文本
	MutexGroups       []MutexGroup      // 互斥组列表
	RequiredGroups    []RequiredGroup   // 必需组列表
	FlagDependencies  []FlagDependency  // 【新增】标志依赖关系列表
	Completion        bool              // 是否启用自动补全标志
	DynamicCompletion bool              // 是否启用动态补全
}

// NewCmdConfig 创建新的命令配置
func NewCmdConfig() *CmdConfig {
	return &CmdConfig{
		Version:           "",
		UseChinese:        false,
		EnvPrefix:         "",
		UsageSyntax:       "",
		Example:           map[string]string{},
		Notes:             []string{},
		LogoText:          "",
		MutexGroups:       []MutexGroup{},
		RequiredGroups:    []RequiredGroup{},
		FlagDependencies:  []FlagDependency{}, // 初始化
		Completion:        false,
		DynamicCompletion: false,
	}
}

// CmdOpts 扩展（在 internal/cmd/cmdopts.go 中）
//
// 添加 FlagDependencies 字段，支持通过 CmdOpts 批量配置依赖关系

type CmdOpts struct {
	// ... 现有字段 ...

	// 子命令和组配置
	SubCmds          []types.Command         // 子命令列表
	MutexGroups      []types.MutexGroup      // 互斥组列表
	RequiredGroups   []types.RequiredGroup   // 必需组列表
	FlagDependencies []types.FlagDependency  // 【新增】标志依赖关系列表
}

// NewCmdOpts 创建新的命令选项
func NewCmdOpts() *CmdOpts {
	return &CmdOpts{
		Examples:         make(map[string]string),
		Notes:            []string{},
		SubCmds:          []types.Command{},
		MutexGroups:      []types.MutexGroup{},
		RequiredGroups:   []types.RequiredGroup{},
		FlagDependencies: []types.FlagDependency{}, // 【新增】初始化
	}
}
```

---

## 三、API 设计

### 3.1 Command 接口扩展

```go
// Command 接口扩展（新增方法）
type Command interface {
	// ... 现有方法 ...

	// AddFlagDependency 添加标志依赖关系
	//
	// 参数:
	//   - name: 依赖关系名称
	//   - trigger: 触发标志名称
	//   - targets: 目标标志名称列表
	//   - depType: 依赖关系类型（DepMutex 或 DepRequired）
	//
	// 返回值:
	//   - error: 添加失败时返回错误
	AddFlagDependency(name, trigger string, targets []string, depType DepType) error
}
```

### 3.2 Cmd 实现

```go
// internal/cmd/cmd_flag_dependency.go

package cmd

import (
	"fmt"
	"gitee.com/MM-Q/qflag/internal/types"
)

// AddFlagDependency 添加标志依赖关系
//
// 参数:
//   - name: 依赖关系名称
//   - trigger: 触发标志名称
//   - targets: 目标标志名称列表
//   - depType: 依赖关系类型
//
// 返回值:
//   - error: 添加失败时返回错误
func (c *Cmd) AddFlagDependency(name, trigger string, targets []string, depType types.DepType) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 验证名称
	if name == "" {
		return fmt.Errorf("empty flag dependency name in '%s'", c.Name())
	}

	// 检查是否已存在
	for _, dep := range c.config.FlagDependencies {
		if dep.Name == name {
			return fmt.Errorf("duplicate flag dependency '%s' in '%s'", name, c.Name())
		}
	}

	// 验证触发标志
	if trigger == "" {
		return fmt.Errorf("empty trigger flag in '%s'", c.Name())
	}

	// 验证目标标志列表
	if len(targets) == 0 {
		return fmt.Errorf("empty target flags in '%s'", c.Name())
	}

	// 检查自依赖
	for _, target := range targets {
		if target == trigger {
			return fmt.Errorf("trigger flag '%s' cannot be in targets in '%s'", trigger, c.Name())
		}
	}

	// 验证触发标志是否存在
	if _, exists := c.flagRegistry.Get(trigger); !exists {
		return fmt.Errorf("trigger flag '%s' not found in '%s'", trigger, c.Name())
	}

	// 验证目标标志是否存在
	for _, target := range targets {
		if _, exists := c.flagRegistry.Get(target); !exists {
			return fmt.Errorf("target flag '%s' not found in '%s'", target, c.Name())
		}
	}

	// 创建依赖关系
	dep := types.FlagDependency{
		Name:    name,
		Trigger: trigger,
		Targets: targets,
		Type:    depType,
	}

	// 添加到配置
	c.config.FlagDependencies = append(c.config.FlagDependencies, dep)
	return nil
}

// ApplyOpts 中的 FlagDependencies 处理（在 internal/cmd/cmd_config.go 中）
//
// 在 ApplyOpts 方法中添加对 FlagDependencies 的处理逻辑：

func (c *Cmd) ApplyOpts(opts *CmdOpts) error {
	// ... 现有代码 ...

	// 4. 添加互斥组
	if len(opts.MutexGroups) > 0 {
		for _, group := range opts.MutexGroups {
			if err := c.AddMutexGroup(group.Name, group.Flags, group.AllowNone); err != nil {
				return fmt.Errorf("add mutex group '%s' failed in '%s': %w", group.Name, c.Name(), err)
			}
		}
	}

	// 5. 添加必需组
	if len(opts.RequiredGroups) > 0 {
		for _, group := range opts.RequiredGroups {
			if err := c.AddRequiredGroup(group.Name, group.Flags, group.Conditional); err != nil {
				return fmt.Errorf("add required group '%s' failed in '%s': %w", group.Name, c.Name(), err)
			}
		}
	}

	// 【新增】6. 添加标志依赖关系
	if len(opts.FlagDependencies) > 0 {
		for _, dep := range opts.FlagDependencies {
			if err := c.AddFlagDependency(dep.Name, dep.Trigger, dep.Targets, dep.Type); err != nil {
				return fmt.Errorf("add flag dependency '%s' failed in '%s': %w", dep.Name, c.Name(), err)
			}
		}
	}

	// 7. 添加子命令
	if len(opts.SubCmds) > 0 {
		if err := c.AddSubCmds(opts.SubCmds...); err != nil {
			return fmt.Errorf("add subcommands failed in '%s': %w", c.Name(), err)
		}
	}

	// ... 其余代码 ...
}
```

---

## 四、验证逻辑

### 4.1 Parser 验证实现

```go
// internal/parser/parser_validation.go

// validateFlagDependencies 验证标志依赖关系
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 如果依赖关系验证失败返回错误
//
// 功能说明:
//   - 遍历所有标志依赖关系
//   - 只有当触发标志被设置时才进行验证
//   - 根据依赖类型执行相应的验证逻辑
//   - 提供清晰的错误信息
func (p *DefaultParser) validateFlagDependencies(config *types.CmdConfig) error {
	if len(config.FlagDependencies) == 0 {
		return nil
	}

	// 使用缓存的已设置标志映射
	setFlags := p.setFlagsMap

	for _, dep := range config.FlagDependencies {
		// 只有当触发标志被设置时才检查
		if !setFlags[dep.Trigger] {
			continue
		}

		switch dep.Type {
		case types.DepMutex:
			// 互斥依赖：检查是否有目标标志被设置
			var conflictFlags []string
			seenConflicts := make(map[string]bool)

			for _, target := range dep.Targets {
				if setFlags[target] {
					displayName := p.flagDisplayNames[target]
					if !seenConflicts[displayName] {
						seenConflicts[displayName] = true
						conflictFlags = append(conflictFlags, displayName)
					}
				}
			}

			if len(conflictFlags) > 0 {
				triggerDisplay := p.flagDisplayNames[dep.Trigger]
				return fmt.Errorf("flag %s cannot be used with %v (dependency: %s)",
					triggerDisplay, conflictFlags, dep.Name)
			}

		case types.DepRequired:
			// 必需依赖：检查所有目标标志是否都被设置
			var missingFlags []string
			seenMissing := make(map[string]bool)

			for _, target := range dep.Targets {
				if !setFlags[target] {
					displayName := p.flagDisplayNames[target]
					if !seenMissing[displayName] {
						seenMissing[displayName] = true
						missingFlags = append(missingFlags, displayName)
					}
				}
			}

			if len(missingFlags) > 0 {
				triggerDisplay := p.flagDisplayNames[dep.Trigger]
				return fmt.Errorf("flag %s requires flags %v to be set (dependency: %s)",
					triggerDisplay, missingFlags, dep.Name)
			}
		}
	}

	return nil
}
```

### 4.2 集成到解析流程

```go
// internal/parser/parser.go

// ParseOnly 解析命令行参数（仅解析，不执行）
func (p *DefaultParser) ParseOnly(cmd types.Command, args []string) error {
	// ... 现有代码：解析参数 ...

	// 构建已设置标志映射（缓存）
	p.buildSetFlagsMap(cmd)

	// 验证互斥组
	if err := p.validateMutexGroups(config); err != nil {
		return err
	}

	// 验证必需组
	if err := p.validateRequiredGroups(config); err != nil {
		return err
	}

	// 【新增】验证标志依赖关系
	if err := p.validateFlagDependencies(config); err != nil {
		return err
	}

	// ... 现有代码 ...
}
```

---

## 五、使用示例

### 5.1 基础用法

```go
package main

import (
	"fmt"
	"os"
	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建命令
	deployCmd := qflag.NewCmd("deploy", "d", qflag.ExitOnError)

	// 定义标志
	deployCmd.Bool("remote", "r", "使用远程部署", false)
	deployCmd.String("local-path", "p", "本地路径", "")
	deployCmd.String("local-config", "c", "本地配置文件", "")

	// 添加互斥依赖：使用 --remote 时，--local-path 和 --local-config 不能使用
	deployCmd.AddFlagDependency(
		"remote-local-mutex",           // 依赖关系名称
		"remote",                       // 触发标志
		[]string{"local-path", "local-config"}, // 目标标志
		qflag.DepMutex,                 // 依赖类型：互斥
	)

	// 解析
	if err := deployCmd.Parse(os.Args[1:]); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
```

### 5.2 SSL 证书场景

```go
package main

import (
	"gitee.com/MM-Q/qflag"
)

func main() {
	serverCmd := qflag.NewCmd("server", "s", qflag.ExitOnError)

	// 定义标志
	serverCmd.Bool("ssl", "", "启用SSL", false)
	serverCmd.String("cert", "", "证书文件路径", "")
	serverCmd.String("key", "", "密钥文件路径", "")
	serverCmd.String("ca-cert", "", "CA证书路径", "")

	// 添加必需依赖：使用 --ssl 时，--cert 和 --key 必须设置
	serverCmd.AddFlagDependency(
		"ssl-requires-cert",
		"ssl",
		[]string{"cert", "key"},
		qflag.DepRequired,
	)

	// 可以添加多个依赖关系
	serverCmd.AddFlagDependency(
		"ssl-optional-ca",
		"ssl",
		[]string{"ca-cert"},
		qflag.DepRequired,
	)

	// 解析执行...
}
```

### 5.3 配置文件场景

```go
package main

import (
	"gitee.com/MM-Q/qflag"
)

func main() {
	runCmd := qflag.NewCmd("run", "r", qflag.ExitOnError)

	// 定义标志
	runCmd.String("config", "c", "配置文件路径", "")
	runCmd.String("host", "H", "服务器地址", "localhost")
	runCmd.Int("port", "p", "端口号", 8080)
	runCmd.String("log-level", "l", "日志级别", "info")

	// 使用配置文件时，其他配置项通过配置文件指定，命令行不能重复设置
	runCmd.AddFlagDependency(
		"config-mutex",
		"config",
		[]string{"host", "port", "log-level"},
		qflag.DepMutex,
	)

	// 解析执行...
}
```

### 5.4 使用 CmdOpts 配置

```go
package main

import (
	"gitee.com/MM-Q/qflag"
)

func main() {
	cmd := qflag.NewCmd("app", "a", qflag.ExitOnError)

	// 定义标志
	cmd.Bool("verbose", "v", "详细模式", false)
	cmd.String("log-file", "l", "日志文件", "")
	cmd.String("output", "o", "输出文件", "")

	// 使用 CmdOpts 配置依赖关系
	opts := &qflag.CmdOpts{
		Desc: "示例应用",
		FlagDependencies: []qflag.FlagDependency{
			{
				Name:    "verbose-requires-log",
				Trigger: "verbose",
				Targets: []string{"log-file"},
				Type:    qflag.DepRequired,
			},
			{
				Name:    "output-mutex-verbose",
				Trigger: "output",
				Targets: []string{"verbose"},
				Type:    qflag.DepMutex,
			},
		},
	}

	cmd.ApplyOpts(opts)

	// 解析执行...
}
```

---

## 六、与现有功能的关系

### 6.1 功能对比

| 功能 | 解决的问题 | 关系类型 | 触发条件 |
|------|-----------|----------|----------|
| **MutexGroup** | 多选一（如 --json 和 --xml） | 组内互斥 | 任意标志被设置 |
| **RequiredGroup** | 必须全选（如 --host, --port） | 组内必需 | 无条件或任意标志被设置 |
| **FlagDependency** | 条件依赖（如 --ssl 需要 --cert） | 单向依赖 | **特定标志被设置** |

### 6.2 组合使用示例

```go
package main

import (
	"gitee.com/MM-Q/qflag"
)

func main() {
	deployCmd := qflag.NewCmd("deploy", "d", qflag.ExitOnError)

	// 定义标志
	deployCmd.Bool("dev", "", "开发环境", false)
	deployCmd.Bool("prod", "", "生产环境", false)
	deployCmd.Bool("remote", "r", "远程部署", false)
	deployCmd.String("ssh-key", "k", "SSH密钥", "")
	deployCmd.String("password", "p", "密码", "")

	opts := &qflag.CmdOpts{
		Desc: "部署命令",
		// 互斥组：只能选择一个环境
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "environment",
				Flags:     []string{"dev", "prod"},
				AllowNone: false, // 必须选择一个
			},
		},
		// 标志依赖：远程部署需要认证方式
		FlagDependencies: []qflag.FlagDependency{
			{
				Name:    "remote-requires-auth",
				Trigger: "remote",
				Targets: []string{"ssh-key", "password"},
				Type:    qflag.DepRequired,
			},
		},
	}

	deployCmd.ApplyOpts(opts)

	// 验证场景：
	// ✅ --dev --remote --ssh-key=key.pem
	// ❌ --dev --prod（互斥组错误）
	// ❌ --prod --remote（缺少认证方式）
}
```

---

## 七、边界情况处理

### 7.1 循环依赖检测

```go
// 潜在问题：A 依赖 B，B 又依赖 A
// 当前设计不直接支持，但可以通过验证避免

// 方案：添加循环依赖检测（可选增强）
func (c *Cmd) detectCircularDependencies() error {
	// 构建依赖图
	graph := make(map[string][]string)
	for _, dep := range c.config.FlagDependencies {
		if dep.Type == types.DepRequired {
			// 必需依赖：trigger -> targets
			graph[dep.Trigger] = append(graph[dep.Trigger], dep.Targets...)
		}
	}

	// 检测循环（使用DFS）
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}

		recStack[node] = false
		return false
	}

	for node := range graph {
		if !visited[node] {
			if dfs(node) {
				return fmt.Errorf("circular dependency detected involving flag '%s'", node)
			}
		}
	}

	return nil
}
```

### 7.2 冲突检测

```go
// 检测 FlagDependency 与 MutexGroup 的潜在冲突
// 例如：MutexGroup 要求 --a 和 --b 互斥
//       FlagDependency 要求 --a 触发时 --b 必需
// 这是一个逻辑矛盾

func (c *Cmd) detectDependencyConflicts() error {
	// 获取所有互斥关系
	mutexPairs := make(map[string]map[string]bool)
	for _, mg := range c.config.MutexGroups {
		for i, f1 := range mg.Flags {
			for _, f2 := range mg.Flags[i+1:] {
				if mutexPairs[f1] == nil {
					mutexPairs[f1] = make(map[string]bool)
				}
				if mutexPairs[f2] == nil {
					mutexPairs[f2] = make(map[string]bool)
				}
				mutexPairs[f1][f2] = true
				mutexPairs[f2][f1] = true
			}
		}
	}

	// 检查 FlagDependency 是否与互斥关系冲突
	for _, dep := range c.config.FlagDependencies {
		if dep.Type == types.DepRequired {
			// 必需依赖：trigger 要求 targets
			// 如果 trigger 和某个 target 在互斥组中，则冲突
			for _, target := range dep.Targets {
				if mutexPairs[dep.Trigger] != nil && mutexPairs[dep.Trigger][target] {
					return fmt.Errorf("conflict: flag '%s' and '%s' are mutually exclusive but '%s' requires '%s'",
						dep.Trigger, target, dep.Trigger, target)
				}
			}
		}
	}

	return nil
}
```

---

## 八、测试策略

### 8.1 单元测试要点

```go
// internal/parser/parser_validation_test.go

func TestValidateFlagDependencies(t *testing.T) {
	tests := []struct {
		name        string
		deps        []types.FlagDependency
		setFlags    map[string]bool
		wantErr     bool
		errContains string
	}{
		{
			name: "mutex dependency - no conflict",
			deps: []types.FlagDependency{
				{Name: "test", Trigger: "a", Targets: []string{"b", "c"}, Type: types.DepMutex},
			},
			setFlags: map[string]bool{"a": true},
			wantErr:  false,
		},
		{
			name: "mutex dependency - conflict",
			deps: []types.FlagDependency{
				{Name: "test", Trigger: "a", Targets: []string{"b", "c"}, Type: types.DepMutex},
			},
			setFlags:    map[string]bool{"a": true, "b": true},
			wantErr:     true,
			errContains: "cannot be used with",
		},
		{
			name: "required dependency - all set",
			deps: []types.FlagDependency{
				{Name: "test", Trigger: "a", Targets: []string{"b", "c"}, Type: types.DepRequired},
			},
			setFlags: map[string]bool{"a": true, "b": true, "c": true},
			wantErr:  false,
		},
		{
			name: "required dependency - missing",
			deps: []types.FlagDependency{
				{Name: "test", Trigger: "a", Targets: []string{"b", "c"}, Type: types.DepRequired},
			},
			setFlags:    map[string]bool{"a": true, "b": true},
			wantErr:     true,
			errContains: "requires flags",
		},
		{
			name: "trigger not set - skip validation",
			deps: []types.FlagDependency{
				{Name: "test", Trigger: "a", Targets: []string{"b", "c"}, Type: types.DepRequired},
			},
			setFlags: map[string]bool{"b": true}, // a 未设置
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行测试...
		})
	}
}
```

### 8.2 集成测试

```go
// 完整的使用场景测试
func TestFlagDependencyIntegration(t *testing.T) {
	cmd := NewCmd("test", "t", ContinueOnError)

	cmd.Bool("ssl", "", "启用SSL", false)
	cmd.String("cert", "", "证书", "")
	cmd.String("key", "", "密钥", "")

	err := cmd.AddFlagDependency("ssl-requires-cert", "ssl", []string{"cert", "key"}, DepRequired)
	if err != nil {
		t.Fatalf("添加依赖失败: %v", err)
	}

	// 测试场景1: 只设置 --ssl，应该失败
	err = cmd.Parse([]string{"--ssl"})
	if err == nil {
		t.Error("期望错误但未返回")
	}

	// 测试场景2: 设置 --ssl --cert --key，应该成功
	cmd.Reset() // 重置解析状态
	err = cmd.Parse([]string{"--ssl", "--cert=server.crt", "--key=server.key"})
	if err != nil {
		t.Errorf("不期望错误但返回: %v", err)
	}
}
```

---

## 九、实现清单

### 9.1 文件修改清单

| 文件路径 | 修改类型 | 说明 |
|---------|---------|------|
| `internal/types/config.go` | 新增 | FlagDependency 类型定义、DepType 枚举 |
| `internal/types/config.go` | 修改 | CmdConfig 添加 FlagDependencies 字段 |
| `internal/cmd/cmd_group.go` | 修改 | 添加 AddFlagDependency 方法 |
| `internal/cmd/cmd_group.go` | **简化** | **移除 GetMutexGroup/RemoveMutexGroup/GetRequiredGroup/RemoveRequiredGroup 等方法** |
| `internal/parser/parser_validation.go` | 新增 | validateFlagDependencies 函数 |
| `internal/parser/parser.go` | 修改 | ParseOnly 中集成依赖验证 |
| `internal/cmd/cmdopts.go` | 修改 | CmdOpts 添加 FlagDependencies 字段 |
| `internal/cmd/cmd_config.go` | 修改 | ApplyOpts 支持 FlagDependencies |
| `internal/types/command.go` | 修改 | Command 接口添加 AddFlagDependency 方法，**移除组管理相关方法** |
| `exports.go` | 修改 | 导出 FlagDependency、DepType、DepMutex、DepRequired |

### 9.2 简化现有 API（重要）

为保持 API 简洁性，本次实现将同步简化 `cmd_group.go` 中现有的"花哨"方法：

**MutexGroup 相关 - 只保留：**
- ✅ `AddMutexGroup(name string, flags []string, allowNone bool) error`
- ✅ `MutexGroups() []MutexGroup`

**移除以下方法：**
- ❌ `GetMutexGroup(name string) (*MutexGroup, bool)`
- ❌ `RemoveMutexGroup(name string) error`
- ❌ `GetMutexGroups() []MutexGroup`（与 `MutexGroups()` 重复）

**RequiredGroup 相关 - 只保留：**
- ✅ `AddRequiredGroup(name string, flags []string, conditional bool) error`
- ✅ `RequiredGroups() []RequiredGroup`

**移除以下方法：**
- ❌ `GetRequiredGroup(name string) (*RequiredGroup, bool)`
- ❌ `RemoveRequiredGroup(name string) error`

**Command 接口同步调整：**
- 只保留 `AddMutexGroup` 和 `AddRequiredGroup`
- 移除 `GetMutexGroup`、`RemoveMutexGroup`、`GetRequiredGroup`、`RemoveRequiredGroup`

### 9.3 测试文件

| 文件路径 | 说明 |
|---------|------|
| `internal/parser/parser_validation_test.go` | 添加 FlagDependency 验证测试 |
| `examples/flag-dependency/` | 新增使用示例 |

### 9.4 文档文件

| 文件路径 | 说明 |
|---------|------|
| `docs/FLAG_DEPENDENCY_DESIGN.md` | 本设计文档 |
| `APIDOC.md` | 更新 API 文档 |
| `README.md` | 添加功能说明 |

---

## 十、总结

### 10.1 设计亮点

1. **语义清晰**：专门解决标志级别的条件依赖，与组级别的约束形成互补
2. **API 简洁**：一行代码即可定义复杂的依赖关系
3. **类型安全**：使用枚举类型定义依赖关系类型
4. **错误友好**：提供清晰的错误信息，包含依赖关系名称
5. **向后兼容**：不影响现有功能，可逐步采用

### 10.2 使用建议

- **MutexGroup**：用于"多选一"场景（如输出格式选择）
- **RequiredGroup**：用于"必须全选"场景（如数据库连接配置）
- **FlagDependency**：用于"如果A则必须/不能B"场景（如 SSL 需要证书）

### 10.3 未来扩展

- 支持更复杂的条件表达式（如 `--a && --b` 触发对 `--c` 的依赖）
- 支持默认值设置（当触发标志设置时，自动设置目标标志的默认值）
- 支持条件值验证（如 `--mode=advanced` 时才要求 `--token`）

---

*设计文档版本: 1.0*  
*创建日期: 2026-05-07*  
*作者: QFlag 开发团队*
